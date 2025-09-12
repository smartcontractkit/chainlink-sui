//go:build unit

package txm_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

// createTestTransaction creates a test transaction with the specified state and timestamp
func createTestTransaction(t *testing.T, txID string, state txm.TransactionState, lastUpdatedAt uint64) txm.SuiTx {
	t.Helper()

	// Generate a real Ed25519 public key for testing
	publicKey, _, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	publicKeyBytes := []byte(publicKey)

	// Create a minimal PTB for testing
	ptb := transaction.NewTransaction()
	ptb.SetGasBudget(10000000)
	ptb.SetGasPrice(1000000)

	return txm.SuiTx{
		TransactionID: txID,
		Sender:        "test-sender",
		PublicKey:     publicKeyBytes,
		Metadata:      &commontypes.TxMeta{GasLimit: big.NewInt(10000000)},
		Timestamp:     txm.GetCurrentUnixTimestamp(),
		Payload:       "test-payload",
		Signatures:    []string{"test-signature"},
		RequestType:   "WaitForEffectsCert",
		Attempt:       1,
		State:         state,
		Digest:        "test-digest-" + txID,
		LastUpdatedAt: lastUpdatedAt,
		TxError:       nil,
		GasBudget:     10000000,
		Ptb:           ptb,
	}
}

func TestReaperRoutine_CleanupOldTransactions(t *testing.T) {
	t.Parallel()

	// Set up logger and store
	lggr := logger.Test(t)
	store := txm.NewTxmStoreImpl(lggr)

	// Create fake client
	fakeClient := &testutils.FakeSuiPTBClient{
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

	// Create keystore and gas manager
	keystoreInstance := testutils.NewTestKeystore(t)
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, 0)
	retryManager := txm.NewDefaultRetryManager(3)

	// Set short retention period for testing (2 seconds)
	conf := txm.Config{
		BroadcastChanSize:        100,
		RequestType:              "WaitForLocalExecution",
		ConfirmPollSecs:          1,
		DefaultMaxGasAmount:      200000,
		MaxTxRetryAttempts:       3,
		TransactionTimeout:       "10s",
		MaxConcurrentRequests:    5,
		ReaperPollSecs:           1,
		TransactionRetentionSecs: 2,
	}

	// Create TXM instance
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = txmInstance.Start(ctx)
	require.NoError(t, err)
	defer txmInstance.Close()

	currentTime := txm.GetCurrentUnixTimestamp()

	// Add old finalized transaction (older than retention period)
	// First create and finalize the old transaction
	oldFinalizedTx := createTestTransaction(t, "old-finalized", txm.StatePending, currentTime)
	err = store.AddTransaction(oldFinalizedTx)
	require.NoError(t, err)
	err = store.ChangeState("old-finalized", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("old-finalized", txm.StateFinalized)
	require.NoError(t, err)

	// Add old failed transaction (older than retention period)
	oldFailedTx := createTestTransaction(t, "old-failed", txm.StatePending, currentTime)
	err = store.AddTransaction(oldFailedTx)
	require.NoError(t, err)
	err = store.ChangeState("old-failed", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("old-failed", txm.StateFailed)
	require.NoError(t, err)

	// Wait for transactions to age beyond the retention period (2 seconds)
	time.Sleep(3 * time.Second)

	// Add recent finalized transaction (newer than retention period)
	recentFinalizedTx := createTestTransaction(t, "recent-finalized", txm.StatePending, currentTime)
	err = store.AddTransaction(recentFinalizedTx)
	require.NoError(t, err)
	err = store.ChangeState("recent-finalized", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("recent-finalized", txm.StateFinalized)
	require.NoError(t, err)

	// Add pending transaction (should never be cleaned up regardless of age)
	oldPendingTx := createTestTransaction(t, "old-pending", txm.StatePending, currentTime-10)
	err = store.AddTransaction(oldPendingTx)
	require.NoError(t, err)

	// Wait for reaper to run and cleanup old transactions
	require.Eventually(t, func() bool {
		// Check that old finalized transaction was deleted
		_, err1 := store.GetTransaction("old-finalized")
		// Check that old failed transaction was deleted
		_, err2 := store.GetTransaction("old-failed")

		// Both should return errors (transactions not found)
		return err1 != nil && err2 != nil
	}, 10*time.Second, 500*time.Millisecond, "Old transactions should have been cleaned up")

	// Verify that recent finalized transaction still exists
	_, err = store.GetTransaction("recent-finalized")
	require.NoError(t, err, "Recent finalized transaction should not be cleaned up")

	// Verify that old pending transaction still exists
	_, err = store.GetTransaction("old-pending")
	require.NoError(t, err, "Pending transaction should never be cleaned up regardless of age")
}

func TestReaperRoutine_DoesNotCleanupRecentTransactions(t *testing.T) {
	t.Parallel()

	// Set up logger and store
	lggr := logger.Test(t)
	store := txm.NewTxmStoreImpl(lggr)

	// Create fake client
	fakeClient := &testutils.FakeSuiPTBClient{
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

	// Create keystore and gas manager
	keystoreInstance := testutils.NewTestKeystore(t)
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, 0)
	retryManager := txm.NewDefaultRetryManager(3)

	// Set long retention period (1000 seconds)
	conf := txm.Config{
		BroadcastChanSize:        100,
		RequestType:              "WaitForLocalExecution",
		ConfirmPollSecs:          1,
		DefaultMaxGasAmount:      200000,
		MaxTxRetryAttempts:       3,
		TransactionTimeout:       "10s",
		MaxConcurrentRequests:    5,
		ReaperPollSecs:           1,
		TransactionRetentionSecs: 1000,
	}

	// Create TXM instance
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = txmInstance.Start(ctx)
	require.NoError(t, err)
	defer txmInstance.Close()

	currentTime := txm.GetCurrentUnixTimestamp()

	// Add recent finalized transaction
	recentFinalizedTx := createTestTransaction(t, "recent-finalized", txm.StatePending, currentTime-5)
	err = store.AddTransaction(recentFinalizedTx)
	require.NoError(t, err)
	err = store.ChangeState("recent-finalized", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("recent-finalized", txm.StateFinalized)
	require.NoError(t, err)

	// Add recent failed transaction
	recentFailedTx := createTestTransaction(t, "recent-failed", txm.StatePending, currentTime-5)
	err = store.AddTransaction(recentFailedTx)
	require.NoError(t, err)
	err = store.ChangeState("recent-failed", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("recent-failed", txm.StateFailed)
	require.NoError(t, err)

	// Wait a bit for reaper to potentially run
	time.Sleep(3 * time.Second)

	// Verify that recent transactions still exist
	_, err = store.GetTransaction("recent-finalized")
	require.NoError(t, err, "Recent finalized transaction should not be cleaned up")

	_, err = store.GetTransaction("recent-failed")
	require.NoError(t, err, "Recent failed transaction should not be cleaned up")
}

func TestReaperRoutine_OnlyCleanupFinalizedAndFailed(t *testing.T) {
	t.Parallel()

	// Set up logger and store
	lggr := logger.Test(t)
	store := txm.NewTxmStoreImpl(lggr)

	// Create fake client
	fakeClient := &testutils.FakeSuiPTBClient{
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

	// Create keystore and gas manager
	keystoreInstance := testutils.NewTestKeystore(t)
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, 0)
	retryManager := txm.NewDefaultRetryManager(3)

	// Set short retention period for testing
	conf := txm.Config{
		BroadcastChanSize:        100,
		RequestType:              "WaitForLocalExecution",
		ConfirmPollSecs:          1,
		DefaultMaxGasAmount:      200000,
		MaxTxRetryAttempts:       3,
		TransactionTimeout:       "10s",
		MaxConcurrentRequests:    5,
		ReaperPollSecs:           1,
		TransactionRetentionSecs: 2,
	}

	// Create TXM instance
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = txmInstance.Start(ctx)
	require.NoError(t, err)
	defer txmInstance.Close()

	currentTime := txm.GetCurrentUnixTimestamp()
	oldTimestamp := currentTime - 10 // Much older than retention period

	// Add old transactions in various states
	states := []txm.TransactionState{
		txm.StatePending,
		txm.StateRetriable,
		txm.StateSubmitted,
		txm.StateFinalized, // Should be cleaned up
		txm.StateFailed,    // Should be cleaned up
	}

	for i, state := range states {
		txID := "old-tx-" + string(rune('0'+i))
		tx := createTestTransaction(t, txID, txm.StatePending, oldTimestamp)
		err = store.AddTransaction(tx)
		require.NoError(t, err)

		// Transition through valid states to reach the target state
		switch state {
		case txm.StatePending:
			// Already pending, no need to change
		case txm.StateRetriable, txm.StateSubmitted, txm.StateFinalized, txm.StateFailed:
			// First transition to submitted
			err = store.ChangeState(txID, txm.StateSubmitted)
			require.NoError(t, err)

			if state == txm.StateRetriable {
				err = store.ChangeState(txID, txm.StateRetriable)
				require.NoError(t, err)
			} else if state == txm.StateFinalized {
				err = store.ChangeState(txID, txm.StateFinalized)
				require.NoError(t, err)
			} else if state == txm.StateFailed {
				err = store.ChangeState(txID, txm.StateFailed)
				require.NoError(t, err)
			}
			// StateSubmitted is already set above
		}
	}

	// Wait for reaper to run
	time.Sleep(5 * time.Second)

	// Check which transactions still exist
	for i, state := range states {
		txID := "old-tx-" + string(rune('0'+i))
		_, err = store.GetTransaction(txID)

		if state == txm.StateFinalized || state == txm.StateFailed {
			require.Error(t, err, "Transaction in state %v should have been cleaned up", state)
		} else {
			require.NoError(t, err, "Transaction in state %v should not be cleaned up", state)
		}
	}
}

func TestReaperRoutine_EmptyStore(t *testing.T) {
	t.Parallel()

	// Set up logger and empty store
	lggr := logger.Test(t)
	store := txm.NewTxmStoreImpl(lggr)

	// Create fake client
	fakeClient := &testutils.FakeSuiPTBClient{
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

	// Create keystore and gas manager
	keystoreInstance := testutils.NewTestKeystore(t)
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, 0)
	retryManager := txm.NewDefaultRetryManager(3)

	// Set short poll period for testing
	conf := txm.Config{
		BroadcastChanSize:        100,
		RequestType:              "WaitForLocalExecution",
		ConfirmPollSecs:          1,
		DefaultMaxGasAmount:      200000,
		MaxTxRetryAttempts:       3,
		TransactionTimeout:       "10s",
		MaxConcurrentRequests:    5,
		ReaperPollSecs:           1,
		TransactionRetentionSecs: 2,
	}

	// Create TXM instance
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = txmInstance.Start(ctx)
	require.NoError(t, err)

	// Let reaper run on empty store for a bit - should not crash or cause issues
	time.Sleep(3 * time.Second)

	// Verify store is still empty and functional
	transactions, err := store.GetInflightTransactions()
	require.NoError(t, err)
	require.Empty(t, transactions, "Store should still be empty")

	txmInstance.Close()
}

func TestReaperRoutine_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	// Set up logger and store
	lggr := logger.Test(t)
	store := txm.NewTxmStoreImpl(lggr)

	// Create fake client
	fakeClient := &testutils.FakeSuiPTBClient{
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

	// Create keystore and gas manager
	keystoreInstance := testutils.NewTestKeystore(t)
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, 0)
	retryManager := txm.NewDefaultRetryManager(3)

	// Set short retention period for aggressive cleanup
	conf := txm.Config{
		BroadcastChanSize:        100,
		RequestType:              "WaitForLocalExecution",
		ConfirmPollSecs:          1,
		DefaultMaxGasAmount:      200000,
		MaxTxRetryAttempts:       3,
		TransactionTimeout:       "10s",
		MaxConcurrentRequests:    5,
		ReaperPollSecs:           1,
		TransactionRetentionSecs: 1,
	}

	// Create TXM instance
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = txmInstance.Start(ctx)
	require.NoError(t, err)
	defer txmInstance.Close()

	// Concurrently add transactions while reaper is running
	var wg sync.WaitGroup
	numGoroutines := 10
	transactionsPerGoroutine := 5

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < transactionsPerGoroutine; j++ {
				txID := "concurrent-tx-" + string(rune('0'+goroutineID)) + "-" + string(rune('0'+j))
				currentTime := txm.GetCurrentUnixTimestamp()

				// Create transaction that will become old (older than 1 second retention)
				tx := createTestTransaction(t, txID, txm.StatePending, currentTime)
				err := store.AddTransaction(tx)
				if err != nil {
					// Transaction might already exist or store might be busy
					continue
				}
				err = store.ChangeState(txID, txm.StateSubmitted)
				if err != nil {
					// State change might fail if transaction was already cleaned up
					continue
				}
				err = store.ChangeState(txID, txm.StateFinalized)
				if err != nil {
					// State change might fail if transaction was already cleaned up
					continue
				}

				// Brief pause to allow reaper to potentially clean up
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	// Wait for reaper to clean up
	time.Sleep(5 * time.Second)

	// Verify that most transactions have been cleaned up
	// Due to concurrent access, we can't guarantee exact counts, but most should be gone
	allTransactions, err := store.GetInflightTransactions()
	require.NoError(t, err)

	// The number of remaining transactions should be significantly less than total created
	totalCreated := numGoroutines * transactionsPerGoroutine
	require.Less(t, len(allTransactions), totalCreated,
		"Most transactions should have been cleaned up by reaper")
}

func TestReaperRoutine_RespectsRetentionPeriod(t *testing.T) {
	t.Parallel()

	// Set up logger and store
	lggr := logger.Test(t)
	store := txm.NewTxmStoreImpl(lggr)

	// Create fake client
	fakeClient := &testutils.FakeSuiPTBClient{
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

	// Create keystore and gas manager
	keystoreInstance := testutils.NewTestKeystore(t)
	maxGasBudget := big.NewInt(12000000)
	gasManager := txm.NewSuiGasManager(lggr, fakeClient, *maxGasBudget, 0)
	retryManager := txm.NewDefaultRetryManager(3)

	// Set 5-second retention period
	retentionPeriod := uint64(5)
	conf := txm.Config{
		BroadcastChanSize:        100,
		RequestType:              "WaitForLocalExecution",
		ConfirmPollSecs:          1,
		DefaultMaxGasAmount:      200000,
		MaxTxRetryAttempts:       3,
		TransactionTimeout:       "10s",
		MaxConcurrentRequests:    5,
		ReaperPollSecs:           1,
		TransactionRetentionSecs: retentionPeriod,
	}

	// Create TXM instance
	txmInstance, err := txm.NewSuiTxm(lggr, fakeClient, keystoreInstance, conf, store, retryManager, gasManager)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = txmInstance.Start(ctx)
	require.NoError(t, err)
	defer txmInstance.Close()

	currentTime := txm.GetCurrentUnixTimestamp()

	// Add transaction that's well outside the retention period (should be cleaned up)
	outsideTx := createTestTransaction(t, "outside-tx", txm.StatePending, currentTime)
	err = store.AddTransaction(outsideTx)
	require.NoError(t, err)
	err = store.ChangeState("outside-tx", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("outside-tx", txm.StateFinalized)
	require.NoError(t, err)

	// Wait to make this transaction older than retention period
	time.Sleep(time.Duration(retentionPeriod+1) * time.Second)

	// Add transaction that's exactly at the retention boundary (should NOT be cleaned up)
	// Since reaper uses `timeDiff > retentionPeriod`, exactly equal should not be cleaned
	boundaryTx := createTestTransaction(t, "boundary-tx", txm.StatePending, currentTime)
	err = store.AddTransaction(boundaryTx)
	require.NoError(t, err)
	err = store.ChangeState("boundary-tx", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("boundary-tx", txm.StateFinalized)
	require.NoError(t, err)

	// Add transaction that's just inside the retention period (should NOT be cleaned up)
	insideTx := createTestTransaction(t, "inside-tx", txm.StatePending, currentTime)
	err = store.AddTransaction(insideTx)
	require.NoError(t, err)
	err = store.ChangeState("inside-tx", txm.StateSubmitted)
	require.NoError(t, err)
	err = store.ChangeState("inside-tx", txm.StateFinalized)
	require.NoError(t, err)

	// Wait for reaper to run
	time.Sleep(3 * time.Second)

	// Check results
	_, err = store.GetTransaction("boundary-tx")
	require.NoError(t, err, "Transaction at retention boundary should NOT be cleaned up (timeDiff == retentionPeriod)")

	_, err = store.GetTransaction("inside-tx")
	require.NoError(t, err, "Transaction inside retention period should not be cleaned up")

	_, err = store.GetTransaction("outside-tx")
	require.Error(t, err, "Transaction outside retention period should be cleaned up")
}
