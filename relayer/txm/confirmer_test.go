//go:build unit

package txm_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/mocks"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

func TestConfirmerRoutine_GasBump(t *testing.T) {
	t.Parallel()
	// Set up logger.
	lggr := logger.Test(t)

	// Use the real in-memory store.
	store := txm.NewTxmStoreImpl(lggr)

	// Create a fake retry manager that marks errors as retryable with the GasBump strategy.
	nrRetries := 3
	retryManager := txm.NewDefaultRetryManager(nrRetries)

	// Create gomock controller and mock client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockSuiPTBClient(ctrl)

	// For this test, we simulate a failure with error "ErrGasBudgetTooHigh".
	// The confirmer will then invoke the retry logic.
	mockClient.EXPECT().
		GetTransactionStatus(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(client.TransactionResult{
			Status: "failure",
			Error:  "ErrGasBudgetTooHigh",
		}, nil)

	// Return a SUI coin with sufficient balance for gas
	objectID := "0x1234567890abcdef1234567890abcdef12345678"
	digest := "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT"
	coinType := "0x2::sui::SUI"
	version := uint64(1)
	balance := uint64(100000000)
	testCoin := &suirpcv2.Object{
		ObjectId:   &objectID,
		Version:    &version,
		Digest:     &digest,
		ObjectType: &coinType,
		Balance:    &balance,
	}

	mockClient.EXPECT().
		QueryCoinsByAddress(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return([]*suirpcv2.Object{testCoin}, nil)

	mockClient.EXPECT().
		GetSUIBalance(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&suirpcv2.Balance{Balance: &balance}, nil)

	gasPrice := uint64(1000)
	mockClient.EXPECT().
		GetReferenceGasPrice(gomock.Any()).
		AnyTimes().
		Return(big.NewInt(int64(gasPrice)), nil)

	mockClient.EXPECT().
		HashTxBytes(gomock.Any()).
		AnyTimes().
		Return([]byte("hashed-tx-bytes"))

	mockClient.EXPECT().
		SendTransaction(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&suirpcv2.ExecuteTransactionResponse{}, nil)

	// Create a fake gas manager that returns an updated gas value.
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, mockClient, *maxGasBudget, 0)
	coinManager := txm.NewGasCoinManager(lggr, mockClient)

	// For the confirmer, the keystore is not used; create a dummy signer.
	keystoreInstance := testutils.NewTestKeystore(t)

	// Use the default configuration.
	conf := txm.DefaultConfigSet

	// Create the TXM.
	txmInstance, err := txm.NewSuiTxm(lggr, mockClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = txmInstance.Start(ctx)

	// Generate a real Ed25519 public key for testing
	publicKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	keystoreInstance.AddKey(privKey)

	// Convert public key to bytes
	publicKeyBytes := []byte(publicKey)

	address, err := client.GetAddressFromPublicKey(publicKeyBytes)
	require.NoError(t, err)

	// Create a minimal PTB for testing
	ptb := transaction.NewTransaction()
	ptb.SetGasBudget(10000000)
	ptb.SetSender(models.SuiAddress(address))
	ptb.SetGasOwner(models.SuiAddress(address))
	ptb.SetGasPrice(10000000)

	coinObjectIdBytes, _ := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(address))
	versionUint, _ := strconv.ParseUint("1", 10, 64)
	digestBytes, _ := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest("9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT"))

	ptb.SetGasPayment([]transaction.SuiObjectRef{
		{
			ObjectId: *coinObjectIdBytes,
			Version:  versionUint,
			Digest:   *digestBytes,
		},
	})

	// Add a transaction in StateSubmitted with a known digest ("test-digest").
	txID := "tx-gasbump-test"
	tx := txm.SuiTx{
		TransactionID: txID,
		Sender:        "dummy-sender",
		PublicKey:     publicKeyBytes,
		Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
		Timestamp:     txm.GetCurrentUnixTimestamp(),
		Payload:       "payload",
		Signatures:    []string{"signature"},
		RequestType:   "WaitForEffectsCert",
		Attempt:       3,
		State:         txm.StateSubmitted,
		Digest:        "test-digest",
		LastUpdatedAt: txm.GetCurrentUnixTimestamp(),
		TxError:       nil,
		GasBudget:     maxGasBudget.Uint64(),
		Ptb:           ptb,
		CoinManager:   coinManager,
	}
	err = store.AddTransaction(tx)
	require.NoError(t, err)
	err = store.ChangeState(txID, txm.StateSubmitted)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		updatedTx, e := store.GetTransaction(txID)
		if e != nil {
			return false
		}

		return updatedTx.Attempt > 2
	}, 15*time.Second, 1*time.Second, "Transaction did not retry as expected")

	require.Eventually(t, func() bool {
		updatedTx, e := store.GetTransaction(txID)
		if e != nil {
			return false
		}

		return updatedTx.State == txm.StateFailed
	}, 15*time.Second, 1000*time.Millisecond, "Transaction did not reach Failed state")

	// Check that the transaction was retried and the gas limit was updated.
	updatedTx, err := store.GetTransaction(txID)
	require.NoError(t, err)
	require.Equal(t, 4, updatedTx.Attempt)
	require.Equal(t, suierrors.ErrGasBudgetTooHigh, updatedTx.TxError)

	txmInstance.Close()
}

func TestConfirmerRoutine_SuccessfulGasBumpAfterTwoAttempts(t *testing.T) {
	t.Parallel()
	// Set up logger.
	lggr := logger.Test(t)

	// Use the real in-memory store.
	store := txm.NewTxmStoreImpl(lggr)

	// Create a fake retry manager that marks errors as retryable with the GasBump strategy.
	nrRetries := 3
	retryManager := txm.NewDefaultRetryManager(nrRetries)

	// Create gomock controller and mock client
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockSuiPTBClient(ctrl)

	// Set up expectations: first two calls fail with GasBudgetTooLow, third succeeds
	// Use Any() for digest since it may change after re-broadcast
	gomock.InOrder(
		mockClient.EXPECT().
			GetTransactionStatus(gomock.Any(), gomock.Any()).
			Return(client.TransactionResult{Status: "failure", Error: "GasBudgetTooLow"}, nil),
		mockClient.EXPECT().
			GetTransactionStatus(gomock.Any(), gomock.Any()).
			Return(client.TransactionResult{Status: "failure", Error: "GasBudgetTooLow"}, nil),
		mockClient.EXPECT().
			GetTransactionStatus(gomock.Any(), gomock.Any()).
			Return(client.TransactionResult{Status: "success"}, nil),
	)

	// Return a SUI coin with sufficient balance for gas
	objectID := "0x1234567890abcdef1234567890abcdef12345678"
	digest := "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT"
	coinType := "0x2::sui::SUI"
	version := uint64(1)
	balance := uint64(100000000)
	testCoin := &suirpcv2.Object{
		ObjectId:   &objectID,
		Version:    &version,
		Digest:     &digest,
		ObjectType: &coinType,
		Balance:    &balance,
	}

	// Allow these methods to be called any number of times with default returns
	mockClient.EXPECT().
		QueryCoinsByAddress(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return([]*suirpcv2.Object{testCoin}, nil)

	mockClient.EXPECT().
		GetSUIBalance(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&suirpcv2.Balance{Balance: &balance}, nil)

	gasPrice := uint64(1000)
	mockClient.EXPECT().
		GetReferenceGasPrice(gomock.Any()).
		AnyTimes().
		Return(big.NewInt(int64(gasPrice)), nil)

	mockClient.EXPECT().
		HashTxBytes(gomock.Any()).
		AnyTimes().
		Return([]byte("hashed-tx-bytes"))

	mockClient.EXPECT().
		SendTransaction(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&suirpcv2.ExecuteTransactionResponse{}, nil)

	// Create a gas manager with lower max budget and percentage increase
	maxGasBudget := big.NewInt(10000000)
	percentIncrease := int64(120) // 120% (20% increase) per bump
	gasManager := txm.NewSuiGasManager(lggr, mockClient, *maxGasBudget, percentIncrease)
	coinManager := txm.NewGasCoinManager(lggr, mockClient)

	// Create keystore
	keystoreInstance := testutils.NewTestKeystore(t)

	// Use the default configuration.
	conf := txm.DefaultConfigSet

	// Create the TXM.
	txmInstance, err := txm.NewSuiTxm(lggr, mockClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = txmInstance.Start(ctx)

	// Generate a real Ed25519 public key for testing
	publicKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	keystoreInstance.AddKey(privKey)

	// Convert public key to bytes
	publicKeyBytes := []byte(publicKey)

	address, err := client.GetAddressFromPublicKey(publicKeyBytes)
	require.NoError(t, err)

	// Create a minimal PTB for testing with low initial gas budget
	initialGasBudget := uint64(6000000) // Start with low gas budget
	ptb := transaction.NewTransaction()
	ptb.SetGasBudget(initialGasBudget)
	ptb.SetSender(models.SuiAddress(address))
	ptb.SetGasOwner(models.SuiAddress(address))
	ptb.SetGasPrice(gasPrice)

	coinObjectIdBytes, _ := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(address))
	versionUint, _ := strconv.ParseUint("1", 10, 64)
	digestBytes, _ := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest("9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT"))

	ptb.SetGasPayment([]transaction.SuiObjectRef{
		{
			ObjectId: *coinObjectIdBytes,
			Version:  versionUint,
			Digest:   *digestBytes,
		},
	})

	// Add a transaction in StateSubmitted with a known digest
	txID := "tx-gasbump-success-test"
	tx := txm.SuiTx{
		TransactionID: txID,
		Sender:        address,
		PublicKey:     publicKeyBytes,
		Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(int64(initialGasBudget))},
		Timestamp:     txm.GetCurrentUnixTimestamp(),
		Payload:       "payload",
		Signatures:    []string{"signature"},
		RequestType:   "WaitForEffectsCert",
		Attempt:       1,
		State:         txm.StateSubmitted,
		Digest:        "test-digest-success",
		LastUpdatedAt: txm.GetCurrentUnixTimestamp(),
		TxError:       nil,
		GasBudget:     maxGasBudget.Uint64(), // Use max budget to allow for gas bumps
		Ptb:           ptb,
		CoinManager:   coinManager,
	}
	err = store.AddTransaction(tx)
	require.NoError(t, err)
	err = store.ChangeState(txID, txm.StateSubmitted)
	require.NoError(t, err)

	// Wait for the transaction to eventually succeed after gas bumps
	require.Eventually(t, func() bool {
		updatedTx, e := store.GetTransaction(txID)
		if e != nil {
			return false
		}

		return updatedTx.State == txm.StateFinalized
	}, 10*time.Second, 100*time.Millisecond, "Transaction did not succeed after gas bumps")

	// Check that the transaction was retried twice and then succeeded
	updatedTx, err := store.GetTransaction(txID)
	require.NoError(t, err)
	require.Equal(t, txm.StateFinalized, updatedTx.State)
	require.Equal(t, 3, updatedTx.Attempt)
	require.Nil(t, updatedTx.TxError)

	// Verify that the gas budget was increased appropriately
	// After 2 bumps: 6M * 1.2 * 1.2 = 8.64M (should be > threshold of 8M)
	expectedMinGas := uint64(8000000)
	require.GreaterOrEqual(t, updatedTx.GasBudget, expectedMinGas, "Gas budget should have been bumped sufficiently")

	txmInstance.Close()
}

func TestConfirmerRoutine_ExponentialBackoffRetry(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	store := txm.NewTxmStoreImpl(lggr)

	nrRetries := 3
	retryManager := txm.NewDefaultRetryManager(nrRetries)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockSuiPTBClient(ctrl)

	checkpointNotFoundErr := suierrors.ErrVerifiedCheckpointNotFound.Error()
	var broadcastCount atomic.Int32

	// Exponential backoff polls status repeatedly while waiting for the delay to elapse,
	// so use a dynamic mock keyed off re-broadcast count instead of a fixed call order.
	mockClient.EXPECT().
		GetTransactionStatus(gomock.Any(), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(context.Context, string) (client.TransactionResult, error) {
			if broadcastCount.Load() < 2 {
				return client.TransactionResult{Status: "failure", Error: checkpointNotFoundErr}, nil
			}

			return client.TransactionResult{Status: "success"}, nil
		})

	objectID := "0x1234567890abcdef1234567890abcdef12345678"
	digest := "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT"
	coinType := "0x2::sui::SUI"
	version := uint64(1)
	balance := uint64(100000000)
	testCoin := &suirpcv2.Object{
		ObjectId:   &objectID,
		Version:    &version,
		Digest:     &digest,
		ObjectType: &coinType,
		Balance:    &balance,
	}

	mockClient.EXPECT().
		QueryCoinsByAddress(gomock.Any(), gomock.Any(), gomock.Any()).
		AnyTimes().
		Return([]*suirpcv2.Object{testCoin}, nil)

	mockClient.EXPECT().
		GetSUIBalance(gomock.Any(), gomock.Any()).
		AnyTimes().
		Return(&suirpcv2.Balance{Balance: &balance}, nil)

	gasPrice := uint64(1000)
	mockClient.EXPECT().
		GetReferenceGasPrice(gomock.Any()).
		AnyTimes().
		Return(big.NewInt(int64(gasPrice)), nil)

	mockClient.EXPECT().
		HashTxBytes(gomock.Any()).
		AnyTimes().
		Return([]byte("hashed-tx-bytes"))

	mockClient.EXPECT().
		SendTransaction(gomock.Any(), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(context.Context, *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error) {
			broadcastCount.Add(1)
			return &suirpcv2.ExecuteTransactionResponse{}, nil
		})

	maxGasBudget := big.NewInt(10000000)
	gasManager := txm.NewSuiGasManager(lggr, mockClient, *maxGasBudget, 0)
	coinManager := txm.NewGasCoinManager(lggr, mockClient)

	keystoreInstance := testutils.NewTestKeystore(t)
	conf := txm.DefaultConfigSet

	txmInstance, err := txm.NewSuiTxm(lggr, mockClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
		txmInstance.Close()
	})

	_ = txmInstance.Start(ctx)

	publicKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	keystoreInstance.AddKey(privKey)
	publicKeyBytes := []byte(publicKey)

	address, err := client.GetAddressFromPublicKey(publicKeyBytes)
	require.NoError(t, err)

	gasBudget := uint64(6000000)
	ptb := transaction.NewTransaction()
	ptb.SetGasBudget(gasBudget)
	ptb.SetSender(models.SuiAddress(address))
	ptb.SetGasOwner(models.SuiAddress(address))
	ptb.SetGasPrice(gasPrice)

	coinObjectIdBytes, _ := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(address))
	versionUint, _ := strconv.ParseUint("1", 10, 64)
	digestBytes, _ := transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest("9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT"))

	ptb.SetGasPayment([]transaction.SuiObjectRef{
		{
			ObjectId: *coinObjectIdBytes,
			Version:  versionUint,
			Digest:   *digestBytes,
		},
	})

	txID := "tx-exponential-backoff-retry-test"
	tx := txm.SuiTx{
		TransactionID: txID,
		Sender:        address,
		PublicKey:     publicKeyBytes,
		Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(int64(gasBudget))},
		Timestamp:     txm.GetCurrentUnixTimestamp(),
		RequestType:   "WaitForEffectsCert",
		Attempt:       1,
		State:         txm.StateSubmitted,
		Digest:        "test-digest-exponential-backoff-retry",
		LastUpdatedAt: txm.GetCurrentUnixTimestamp(),
		GasBudget:     gasBudget,
		Ptb:           ptb,
		CoinManager:   coinManager,
	}

	// Exponential backoff re-broadcasts the existing payload without rebuilding it
	// (unlike gas bump), so the tx must already have valid base64 payload and signatures.
	err = tx.UpdateBSCPayload(ctx, lggr, keystoreInstance, mockClient)
	require.NoError(t, err)

	err = store.AddTransaction(tx)
	require.NoError(t, err)
	err = store.ChangeState(txID, txm.StateSubmitted)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		updatedTx, e := store.GetTransaction(txID)
		if e != nil {
			return false
		}

		return updatedTx.State == txm.StateFinalized
	}, 30*time.Second, 100*time.Millisecond, "Transaction did not finalize after exponential backoff retries")

	updatedTx, err := store.GetTransaction(txID)
	require.NoError(t, err)
	require.Equal(t, 3, updatedTx.Attempt)
	require.Nil(t, updatedTx.TxError)
}
