//go:build unit

package txm

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

func GetTransaction() SuiTx {
	return SuiTx{
		TransactionID: "1",
		Sender:        "0x123",
		Metadata:      nil,
		Timestamp:     0,
		Payload:       []byte("payload"),
		Attempt:       0,
		State:         StatePending,
	}
}

func TestAddTransaction(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	tx := GetTransaction()
	err := store.AddTransaction(tx)
	require.NoError(t, err, "expected no error when adding a new transaction")
	storeTx, err := store.GetTransaction("1")
	require.NoError(t, err, "expected no error when retrieving the added transaction")
	assert.Equal(t, tx, storeTx, "transaction should be added to the store")
	assert.Contains(t, store.stateBuckets[StatePending], "1", "transaction ID should be in the Pending state bucket")
}

func TestAddDuplicateTransaction(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	tx := GetTransaction()
	_ = store.AddTransaction(tx)

	err := store.AddTransaction(tx)
	assert.Error(t, err, "expected an error when adding a duplicate transaction")
}

func TestGetTransaction(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	tx := GetTransaction()
	_ = store.AddTransaction(tx)

	retrievedTx, err := store.GetTransaction("1")
	require.NoError(t, err, "expected no error when retrieving an existing transaction")
	assert.Equal(t, tx, retrievedTx, "retrieved transaction should match the added transaction")
}

func TestGetNonExistentTransaction(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	_, err := store.GetTransaction("1")
	require.Error(t, err, "expected an error when retrieving a non-existent transaction")
}

func TestChangeState(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	tx := GetTransaction()
	_ = store.AddTransaction(tx)

	err := store.ChangeState("1", StateSubmitted)
	require.NoError(t, err, "expected no error when changing the state of a transaction")
	updatedTx, err := store.GetTransaction("1")
	require.NoError(t, err, "expected no error when retrieving the updated transaction")
	assert.Equal(t, StateSubmitted, updatedTx.State, "transaction state should be updated")
	txs, err := store.GetTransactionsByState(StatePending)
	require.NoError(t, err, "expected no error when retrieving transactions by state")
	assert.Empty(t, txs, "expected no transactions in the old state bucket")
	assert.Contains(t, store.stateBuckets[StateSubmitted], "1", "transaction ID should be added to the new state bucket")
}

func TestInvalidStateTransition(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	tx := GetTransaction()
	_ = store.AddTransaction(tx)

	err := store.ChangeState("1", StatePending)
	assert.Error(t, err, "expected an error when attempting an invalid state transition")
}

func TestDeleteTransaction(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	tx := GetTransaction()
	_ = store.AddTransaction(tx)

	err := store.DeleteTransaction("1")
	require.NoError(t, err, "expected no error when deleting a transaction")
	assert.NotContains(t, store.transactions, "1", "transaction should be removed from the store")
	assert.NotContains(t, store.stateBuckets[StatePending], "1", "transaction ID should be removed from the state bucket")
}

func TestDeleteNonExistentTransaction(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	err := store.DeleteTransaction("1")
	assert.Error(t, err, "expected an error when deleting a non-existent transaction")
}

func TestGetTransactionsByState(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	tx1 := GetTransaction()
	tx2 := GetTransaction()
	tx2.TransactionID = "2"

	tx3 := GetTransaction()
	tx3.TransactionID = "3"

	_ = store.AddTransaction(tx1)
	_ = store.AddTransaction(tx2)
	_ = store.AddTransaction(tx3)

	_ = store.ChangeState("3", StateSubmitted)

	pendingTxs, err := store.GetTransactionsByState(StatePending)
	require.NoError(t, err, "expected no error when retrieving transactions by state")
	assert.Len(t, pendingTxs, 2, "expected two transactions in the Pending state")
	assert.Contains(t, pendingTxs, tx1, "expected tx1 to be in the Pending state")
	assert.Contains(t, pendingTxs, tx2, "expected tx2 to be in the Pending state")

	submittedTxs, err := store.GetTransactionsByState(StateSubmitted)
	require.NoError(t, err, "expected no error when retrieving transactions by state")
	tx3, err = store.GetTransaction("3")
	require.NoError(t, err, "expected no error when retrieving the submitted transaction")
	assert.Len(t, submittedTxs, 1, "expected one transaction in the Submitted state")
	assert.Contains(t, submittedTxs, tx3, "expected tx3 to be in the Submitted state")
}

func TestGetTransactionsByInvalidState(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	_, err := store.GetTransactionsByState(999) // Invalid state
	assert.Error(t, err, "expected an error when retrieving transactions by an invalid state")
}

func TestConcurrentReadAndWrite(t *testing.T) {
	t.Parallel()
	store := NewTxmStoreImpl()

	// Create a transaction
	tx := GetTransaction()

	// Add the transaction to the store
	err := store.AddTransaction(tx)
	require.NoError(t, err, "expected no error when adding a transaction")

	// Use a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Number of goroutines to spawn
	numReaders := 50
	numWriters := 10

	// Concurrently read the transaction
	wg.Add(numReaders)
	for range numReaders {
		go func() {
			defer wg.Done()
			for range 10 { // Simulate multiple reads per goroutine
				localTx, localErr := store.GetTransaction("1")
				require.NoError(t, localErr, "expected no error when retrieving the transaction")
				assert.NotNil(t, localTx, "transaction should not be nil")
			}
		}()
	}

	// Create a channel to coordinate transitions between writers
	transitionCh := make(chan struct{}, 1)
	// Initial signal to start transitions
	transitionCh <- struct{}{}

	// Use atomic counter to track transitions
	var transitionCounter int32 = 0

	// Concurrently write (update the state of the transaction)
	wg.Add(numWriters)
	for i := range numWriters {
		go func(writerID int) {
			defer wg.Done()

			// Each writer waits for its turn
			<-transitionCh

			// Get current transaction state
			localTx, localErr := store.GetTransaction("1")
			require.NoError(t, localErr, "expected no error when retrieving the transaction")

			// Perform state transition based on current state
			switch localTx.State {
			case StatePending:
				localErr := store.ChangeState("1", StateSubmitted)
				require.NoError(t, localErr, "expected no error when changing the state to Submitted")
				atomic.AddInt32(&transitionCounter, 1)
			case StateSubmitted:
				localErr := store.ChangeState("1", StateRetriable)
				require.NoError(t, localErr, "expected no error when changing the state to Retriable")
				atomic.AddInt32(&transitionCounter, 1)
			case StateRetriable:
				localErr := store.ChangeState("1", StateFinalized)
				require.NoError(t, localErr, "expected no error when changing the state to Finalized")
				atomic.AddInt32(&transitionCounter, 1)
			case StateFinalized:
				// No more transitions possible
			case StateFailed:
				// No more transitions possible
			}

			// Signal the next writer
			transitionCh <- struct{}{}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Ensure the transaction is still in the store and has the correct state
	finalTx, err := store.GetTransaction("1")
	require.NoError(t, err, "expected no error when retrieving the transaction after concurrent access")
	assert.NotNil(t, finalTx, "transaction should not be nil")

	// Assert that at least some transitions occurred
	assert.Positive(t, transitionCounter, "expected some state transitions to occur")

	// The transaction should end up in one of these valid states
	validFinalStates := []TransactionState{StateSubmitted, StateRetriable, StateFinalized}
	assert.Contains(t, validFinalStates, finalTx.State, "expected the transaction to be in a valid state")
}
