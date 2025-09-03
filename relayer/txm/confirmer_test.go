//go:build unit

package txm_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
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

	// For this test, we simulate a failure with error "simulated gas error".
	// The confirmer will then invoke the retry logic.
	fakeClient := &testutils.FakeSuiPTBClient{
		Status: client.TransactionResult{
			Status: "failure",
			Error:  "ErrGasBudgetTooHigh",
		},
		CoinsData: []models.CoinData{
			{
				CoinType:     "0x2::sui::SUI",
				Balance:      "100000000",
				CoinObjectId: "0x1234567890abcdef1234567890abcdef12345678",
				Version:      "1",
				Digest:       "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT",
			},
		},
	}

	// Create a fake gas manager that returns an updated gas value.
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, 0)

	// For the confirmer, the keystore is not used; create a dummy signer.
	keystoreInstance := testutils.NewTestKeystore(t)

	// Use the default configuration.
	conf := txm.DefaultConfigSet

	// Create the TXM.
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
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
		Attempt:       1,
		State:         txm.StateSubmitted,
		Digest:        "test-digest",
		LastUpdatedAt: txm.GetCurrentUnixTimestamp(),
		TxError:       nil,
		GasBudget:     maxGasBudget.Uint64(),
		Ptb:           ptb,
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

		return updatedTx.State == txm.StateFailed
	}, 5*time.Second, 100*time.Millisecond, "Transaction did not retry as expected")

	// Check that the transaction was retried and the gas limit was updated.
	updatedTx, err := store.GetTransaction(txID)
	require.NoError(t, err)
	require.Equal(t, 3, updatedTx.Attempt)
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
	nrRetries := 5
	retryManager := txm.NewDefaultRetryManager(nrRetries)

	// Create a stateful fake client that changes behavior based on gas budget
	fakeClient := &testutils.StatefulFakeSuiPTBClient{
		CoinsData: []models.CoinData{
			{
				CoinType:     "0x2::sui::SUI",
				Balance:      "100000000",
				CoinObjectId: "0x1234567890abcdef1234567890abcdef12345678",
				Version:      "1",
				Digest:       "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT",
			},
		},
		// Track gas budgets and return appropriate status
		GasBudgetThreshold: 8000000, // Minimum gas budget required for success
		CallCount:          0,
	}

	// Create a gas manager with lower max budget and percentage increase
	maxGasBudget := big.NewInt(10000000)
	percentIncrease := int64(120) // 120% (20% increase) per bump
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, percentIncrease)

	// Create keystore
	keystoreInstance := testutils.NewTestKeystore(t)

	// Use the default configuration.
	conf := txm.DefaultConfigSet

	// Create the TXM.
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
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
	ptb.SetGasPrice(uint64(1000))

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
	require.Equal(t, 5, updatedTx.Attempt) // Initial attempt + 4 gas bump cycles
	require.Nil(t, updatedTx.TxError)

	// Verify that the gas budget was increased appropriately
	// After 2 bumps: 6M * 1.2 * 1.2 = 8.64M (should be > threshold of 8M)
	expectedMinGas := uint64(8000000)
	require.GreaterOrEqual(t, updatedTx.GasBudget, expectedMinGas, "Gas budget should have been bumped sufficiently")

	txmInstance.Close()
}
