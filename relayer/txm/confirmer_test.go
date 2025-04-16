//go:build unit

package txm_test

import (
	"context"
	"crypto/ed25519"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

func TestConfirmerRoutine_GasBump(t *testing.T) {
	t.Parallel()
	// Set up logger.
	lggr := logger.Test(t)

	// Use the real in-memory store.
	store := txm.NewTxmStoreImpl()

	// Create a fake retry manager that marks errors as retryable with the GasBump strategy.
	nrRetries := 3
	retryManager := txm.NewDefaultRetryManager(nrRetries)

	// Create a fake gas manager that returns an updated gas value.
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, *maxGasBudget, 0)

	// For this test, we simulate a failure with error "simulated gas error".
	// The confirmer will then invoke the retry logic.
	fakeClient := &testutils.FakeSuiPTBClient{
		Status: client.TransactionResult{
			Status: "failure",
			Error:  "ErrGasBudgetTooHigh",
		},
	}

	// For the confirmer, the keystore is not used; create a dummy signer.
	dummyPrivateKey := make([]byte, ed25519.PrivateKeySize)
	dummySigner := signer.NewPrivateKeySigner(dummyPrivateKey)

	// Use the default configuration.
	conf := txm.DefaultConfigSet

	// Create the TXM.
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, nil, conf, dummySigner, store, retryManager, gasManager)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = txmInstance.Start(ctx)

	// Add a transaction in StateSubmitted with a known digest ("test-digest").
	txID := "tx-gasbump-test"
	tx := txm.SuiTx{
		TransactionID: txID,
		Sender:        "dummy-sender",
		Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
		Timestamp:     txm.GetCurrentUnixTimestamp(),
		Payload:       []byte("payload"),
		Signatures:    []string{"signature"},
		RequestType:   "WaitForEffectsCert",
		Attempt:       1,
		State:         txm.StateSubmitted,
		Digest:        "test-digest",
		LastUpdatedAt: txm.GetCurrentUnixTimestamp(),
		TxError:       nil,
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
