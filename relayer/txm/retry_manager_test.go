//go:build unit

package txm_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

// dummySuiTx returns a dummy SuiTx with the given number of retries.
func dummySuiTx(numRetries int) *txm.SuiTx {
	return &txm.SuiTx{
		TransactionID: "dummy-tx",
		Attempt:       numRetries,
	}
}

func TestCallbackRetryManager_IsRetryable_Scenarios(t *testing.T) {
	t.Parallel()
	// Define test scenarios
	type testScenario struct {
		name          string
		txRetries     int    // Current transaction retry count
		errMessage    string // Input error message
		maxRetries    int    // Configured maximum retries
		expectedRetry bool   // Expected flag: retryable or not
		expectedStrat txm.RetryStrategy
	}

	scenarios := []testScenario{
		{
			name:          "Gas error with retries available returns GasBump",
			txRetries:     0,
			errMessage:    "Transaction failed: GasPriceTooHigh", // Should map to GasErrors by suierrors mapping
			maxRetries:    3,
			expectedRetry: true,
			expectedStrat: txm.GasBump,
		},
		{
			name:          "Non gas retryable error returns ExponentialBackoff",
			txRetries:     0,
			errMessage:    "Transaction failed: PackageVerificationTimeout", // Should fall under default (non gas) retryable
			maxRetries:    3,
			expectedRetry: true,
			expectedStrat: txm.ExponentialBackoff,
		},
		{
			name:          "Exceeded max retries returns NoRetry",
			txRetries:     3,
			errMessage:    "Transaction failed: GasPriceTooHigh",
			maxRetries:    3,
			expectedRetry: false,
			expectedStrat: txm.NoRetry,
		},
		{
			name:          "Non-retryable error returns NoRetry",
			txRetries:     0,
			errMessage:    "Transaction failed: SomeRandomNonRetryableError",
			maxRetries:    3,
			expectedRetry: false,
			expectedStrat: txm.NoRetry,
		},
	}

	for _, tc := range scenarios {
		// capture range var
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Create a dummy transaction with specified retry count
			tx := dummySuiTx(tc.txRetries)

			// Create a new Strategy with maxRetries
			rm := txm.NewDefaultRetryManager(tc.maxRetries)

			// Call IsRetryable with our error message
			retryable, start := rm.IsRetryable(tx, tc.errMessage)

			// Validate the results
			assert.Equal(t, tc.expectedRetry, retryable, "retryable flag mismatch")
			assert.Equal(t, tc.expectedStrat, start, "retry strategy mismatch")
		})
	}
}

func TestDefaultRetryManager_RegisterStrategy_OverridesDefault(t *testing.T) {
	t.Parallel()
	// This test verifies that RegisterStrategy properly overrides the default strategy.
	rm := txm.NewDefaultRetryManager(5)
	called := false

	customStrategy := func(tx *txm.SuiTx, txErrorMsg string, maxRetries int) (bool, txm.RetryStrategy) {
		called = true
		return true, txm.ExponentialBackoff
	}
	rm.RegisterStrategyFunc(customStrategy)

	tx := dummySuiTx(0)
	retryable, start := rm.IsRetryable(tx, "any error message")

	require.True(t, called, "custom strategy should be invoked")
	assert.True(t, retryable, "custom strategy should return retryable true")
	assert.Equal(t, txm.ExponentialBackoff, start, "expected ImmediateRetry as strategy from custom strategy")
}
