package txm

import (
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
)

// RetryStrategy defines different strategies for retrying a transaction after a failure.
type RetryStrategy int

const (
	// NoRetry indicates that the transaction should not be retried.
	NoRetry RetryStrategy = iota
	// ExponentialBackoff indicates that the transaction should be retried with exponential backoff delays.
	ExponentialBackoff
	// GasBump indicates that the transaction should be retried after bumping its gas budget.
	GasBump
)

// RetryStrategyFunc defines a function signature for evaluating if a transaction error
// is retryable and determining which retry strategy to use.
//
// Parameters:
//   - tx: the transaction that encountered an error.
//   - txErrorMsg: the error message returned for the transaction.
//   - maxRetries: the maximum number of allowable retries.
//
// Returns:
//   - bool: true if the transaction error is retryable, false otherwise.
//   - RetryStrategy: the retry strategy to use if the transaction is retryable.
type RetryStrategyFunc func(tx *SuiTx, txErrorMsg string, maxRetries int) (bool, RetryStrategy)

// RetryManager is an interface for evaluating transaction errors and registering a custom strategy function.
// It allows the caller to determine whether a given transaction error is retryable and to obtain a recommended retry strategy.
type RetryManager interface {
	// IsRetryable checks if the error for the given transaction is retryable.
	// It returns a bool indicating if a retry is allowed and a RetryStrategy value recommending the retry approach.
	IsRetryable(tx *SuiTx, errMessage string) (bool, RetryStrategy)
	// RegisterStrategyFunc sets a custom strategy function for retry evaluation.
	// The provided function will be used for future retry evaluations.
	RegisterStrategyFunc(strategyFunc RetryStrategyFunc)
	// GetMaxNumberRetries returns the maximum number of retries allowed.
	GetMaxNumberRetries() int
}

// DefaultRetryManager is a concrete implementation of RetryManager that uses an injected strategy function.
type DefaultRetryManager struct {
	strategyFunc     RetryStrategyFunc
	maxNumberRetries int
}

// NewDefaultRetryManager creates a new DefaultRetryManager with the given maximum number of retries.
// The default strategy function (defaultRetryCallback) is used unless a custom function is registered later.
//
// Parameters:
//   - retries: maximum number of retries allowed for a transaction.
//
// Returns:
//   - *DefaultRetryManager: a new instance of DefaultRetryManager.
func NewDefaultRetryManager(retries int) *DefaultRetryManager {
	return &DefaultRetryManager{
		maxNumberRetries: retries,
		strategyFunc: func(tx *SuiTx, errorMsg string, retryNumber int) (bool, RetryStrategy) {
			return defaultRetryStrategy(tx, errorMsg, retryNumber)
		},
	}
}

// GetMaxNumberRetries returns the maximum number of retries allowed.
//
// Returns:
//   - int: the maximum number of retries configured for the manager.
func (rm *DefaultRetryManager) GetMaxNumberRetries() int {
	return rm.maxNumberRetries
}

// RegisterStrategyFunc registers the provided strategy function for retry evaluation.
// After registration, the function is used by IsRetryable to determine whether a transaction error is retryable.
//
// Parameters:
//   - strategyFunc: a RetryStrategyFunc function to use for evaluating transaction errors.
func (rm *DefaultRetryManager) RegisterStrategyFunc(strategyFunc RetryStrategyFunc) {
	rm.strategyFunc = strategyFunc
}

// IsRetryable uses the registered strategy function to determine if a transaction error is retryable.
// If no function is registered, it returns false and NoRetry.
//
// Parameters:
//   - tx: a pointer to the SuiTx that encountered an error.
//   - errMessage: the error message returned for the transaction.
//
// Returns:
//   - bool: true if the error is considered retryable, false otherwise.
//   - RetryStrategy: the strategy to use when retrying the transaction.
func (rm *DefaultRetryManager) IsRetryable(tx *SuiTx, errMessage string) (bool, RetryStrategy) {
	if rm.strategyFunc == nil {
		return false, NoRetry
	}

	return rm.strategyFunc(tx, errMessage, rm.maxNumberRetries)
}

// defaultRetryStrategy is the default implementation of the RetryStrategyFunc.
// It inspects the transaction error message and returns a strategy based on the error category.
// Specifically, if the error pertains to gas issues and the transaction has not exceeded the maximum retry count,
// it returns true with the GasBump strategy, else it returns ExponentialBackoff for other errors.
//
// Parameters:
//   - tx: the transaction that encountered an error.
//   - txErrorMsg: the error message.
//   - maxRetries: the maximum allowable retries for the transaction.
//
// Returns:
//   - bool: true if the error is retryable and the retry count has not been exceeded.
//   - RetryStrategy: the recommended retry strategy (GasBump or ExponentialBackoff), or NoRetry if not retryable.
func defaultRetryStrategy(tx *SuiTx, txErrorMsg string, maxRetries int) (bool, RetryStrategy) {
	txError := suierrors.ParseSuiErrorMessage(txErrorMsg)

	if !suierrors.IsRetryable(txError) {
		return false, NoRetry
	}

	// Check if the transaction has exceeded the number of retries allowed.
	if tx.Attempt >= maxRetries {
		return false, NoRetry
	}

	// nolint:exhaustive
	switch txError.Category {
	case suierrors.GasErrors:
		return true, GasBump
	default:
		return true, ExponentialBackoff
	}
}
