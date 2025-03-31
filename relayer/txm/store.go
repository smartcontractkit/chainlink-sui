package txm

import (
	"fmt"
	"sync"
)

// Define transaction states as constants.
// These states represent the lifecycle of a transaction in the system.
const (
	// StatePending represents a transaction that has been created but not yet submitted.
	StatePending = "Pending"

	// StateSubmitted represents a transaction that has been submitted to the network.
	StateSubmitted = "Submitted"

	// StateFinalized represents a transaction that has been successfully executed and finalized.
	StateFinalized = "Finalized"

	// StateRetriable represents a transaction that encountered an issue but can be retried.
	StateRetriable = "Retriable"

	// StateFailed represents a transaction that has failed permanently.
	StateFailed = "Failed"
)

// TxmStore defines the interface for managing transaction lifecycle.
// It provides methods for adding, retrieving, updating, and deleting transactions,
// as well as querying transactions by their current state.
type TxmStore interface {
	// AddTransaction adds a new transaction to the store.
	// Returns an error if a transaction with the same ID already exists.
	AddTransaction(tx SuiTx) error

	// GetTransaction retrieves a transaction by its ID.
	// Returns the transaction and nil if found, otherwise returns an empty transaction and an error.
	GetTransaction(transactionID string) (SuiTx, error)

	// ChangeState updates the state of a transaction.
	// Returns an error if the transaction is not found or if the state transition is invalid.
	ChangeState(transactionID string, state string) error

	// DeleteTransaction removes a transaction from the store.
	// Returns an error if the transaction is not found.
	DeleteTransaction(transactionID string) error

	// GetTransactionsByState retrieves all transactions in a given state.
	// Returns a slice of transactions and nil if successful, otherwise returns nil and an error.
	GetTransactionsByState(state string) ([]SuiTx, error)
}

// InMemoryStore implements the TxmStore interface using in-memory data structures.
// It provides thread-safe operations on transactions using a read-write mutex.
// The implementation is optimized for memory efficiency and performance:
//   - The main 'transactions' map stores pointers to avoid duplicating the transaction data
//   - The 'stateBuckets' use nested maps with empty struct values (map[string]struct{})
//     which consume zero additional memory while enabling O(1) lookups by state
//
// This design allows for efficient filtering of transactions by state without
// maintaining duplicate copies of transaction data or performing expensive iterations.
type InMemoryStore struct {
	mu           sync.RWMutex                   // Mutex to control concurrent access to the data structures
	transactions map[string]*SuiTx              // Main map to store pointers to transactions by ID
	stateBuckets map[string]map[string]struct{} // Auxiliary maps to store transaction IDs by state for efficient lookups
}

// NewTxmStoreImpl creates and initializes a new InMemoryStore instance.
// It initializes the transactions map and state buckets for all possible transaction states.
func NewTxmStoreImpl() *InMemoryStore {
	return &InMemoryStore{
		transactions: make(map[string]*SuiTx),
		stateBuckets: map[string]map[string]struct{}{
			StatePending:   make(map[string]struct{}),
			StateSubmitted: make(map[string]struct{}),
			StateFinalized: make(map[string]struct{}),
			StateRetriable: make(map[string]struct{}),
			StateFailed:    make(map[string]struct{}),
		},
	}
}

// AddTransaction adds a new transaction to the store.
// It sets the initial state to StatePending and updates both the transactions map
// and the state buckets accordingly.
// Returns an error if a transaction with the same ID already exists.
func (s *InMemoryStore) AddTransaction(tx SuiTx) error {
	id := tx.transactionID

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the transaction already exists
	_, exists := s.transactions[id]
	if exists {
		return fmt.Errorf("transaction already exists")
	}

	tx.state = StatePending

	// Add to the main transactions map
	s.transactions[id] = &tx

	// Add the transaction ID to the appropriate state bucket
	s.stateBuckets[StatePending][id] = struct{}{}

	return nil
}

// GetTransaction retrieves a transaction by its ID.
// It acquires a read lock to ensure thread safety.
// Returns the transaction and nil if found, otherwise returns an empty transaction and an error.
func (s *InMemoryStore) GetTransaction(transactionID string) (SuiTx, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, exists := s.transactions[transactionID]
	if !exists {
		var empty SuiTx
		return empty, fmt.Errorf("transaction not found")
	}

	return *tx, nil
}

// ChangeState updates the state of a transaction.
// It validates the state transition according to the allowed transitions:
// - Pending -> Submitted
// - Submitted -> Finalized, Retriable, or Failed
// - Retriable -> Submitted, Failed, or Finalized
// - Finalized and Failed are terminal states
// Returns an error if the transaction is not found or if the state transition is invalid.
func (s *InMemoryStore) ChangeState(transactionID string, newState string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, exists := s.transactions[transactionID]
	if !exists {
		return fmt.Errorf("transaction not found")
	}

	oldState := tx.state

	// Check if the state transition is valid
	switch oldState {
	case StatePending:
		if newState != StateSubmitted {
			return fmt.Errorf("pending state must transition to submitted")
		}
	case StateSubmitted:
		if newState == StatePending {
			return fmt.Errorf("submitted state cannot transition to pending")
		}
	case StateFinalized:
		return fmt.Errorf("finalized state cannot transition to any other state")
	case StateRetriable:
		if newState != StateSubmitted && newState != StateFailed && newState != StateFinalized {
			return fmt.Errorf("invalid state transition from %s to %s", oldState, newState)
		}
	case StateFailed:
		return fmt.Errorf("invalid state transition from %s to %s", oldState, newState)
	default:
		return fmt.Errorf("invalid state: %s", oldState)
	}

	// Remove from the old state bucket
	delete(s.stateBuckets[oldState], transactionID)

	// Update the transaction's state
	tx.state = newState

	// Add the transaction ID to the new state bucket
	s.stateBuckets[newState][transactionID] = struct{}{}

	// Update the transaction in the main transactions map
	delete(s.transactions, transactionID)
	s.transactions[transactionID] = tx

	return nil
}

// DeleteTransaction removes a transaction from the store.
// It removes the transaction from both the transactions map and the appropriate state bucket.
// Returns an error if the transaction is not found.
func (s *InMemoryStore) DeleteTransaction(transactionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, exists := s.transactions[transactionID]
	if !exists {
		return fmt.Errorf("transaction not found")
	}

	// Get transaction state
	state := tx.state

	// Remove from the main transactions map
	delete(s.transactions, transactionID)

	// Remove from the state bucket
	delete(s.stateBuckets[state], transactionID)

	return nil
}

// GetTransactionsByState retrieves all transactions in a given state.
// It acquires a read lock to ensure thread safety.
// Returns a slice of transactions and nil if successful, otherwise returns nil and an error.
func (s *InMemoryStore) GetTransactionsByState(state string) ([]SuiTx, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stateMap, exists := s.stateBuckets[state]
	if !exists {
		return nil, fmt.Errorf("invalid state: %s", state)
	}

	// Collect transaction pointers from the main transactions map
	transactions := make([]SuiTx, 0, len(stateMap))
	for id := range stateMap {
		tx := s.transactions[id]
		transactions = append(transactions, *tx)
	}

	return transactions, nil
}
