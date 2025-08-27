// This module implements gas management functionality. It defines the GasManager
// interface for estimating gas budgets and applying gas bumps (increasing gas limits) to transactions,
// as well as a concrete implementation (SuiGasManager) that uses a simple fixed-percentage increase
// heuristic.
package txm

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const (
	// gasLimitPercentualIncrease is the fixed percentage increase applied during a gas bump.
	gasLimitPercentualIncrease = 120
	// percentualNormalization is used to normalize the percentage calculation.
	percentualNormalization = 100
)

// GasManager defines the interface for managing and adjusting gas budgets for a transaction.
// It provides methods to estimate the gas budget needed for a Sui transaction and to compute a new
// gas budget (gas bump) when a transaction encounters a gas-related issue.
type GasManager interface {
	// EstimateGasBudget estimates the gas budget required for the given transaction.
	//
	// Parameters:
	//   - ctx: Context allowing cancellation and timeouts.
	//   - tx: The Sui transaction for which to estimate the gas budget.
	//
	// Returns:
	//   - uint64: The estimated gas budget.
	//   - error: An error if estimation is not implemented or fails.
	EstimateGasBudget(ctx context.Context, tx *SuiTx) (uint64, error)

	// GasBump calculates a new gas budget for the given transaction by increasing its current gas limit.
	// The new budget is computed by increasing the current gas limit by a specified percentage (defaults to 20%)
	// If the new budget would exceed the maximum allowed budget, it is capped, and if the current gas
	// limit is already at or above the maximum, an error is returned.
	//
	// Parameters:
	//   - ctx: Context allowing cancellation and timeouts.
	//   - tx: The Sui transaction whose gas budget needs to be bumped.
	//
	// Returns:
	//   - big.Int: The new gas budget.
	//   - error: An error if the gas bump operation cannot be performed.
	GasBump(ctx context.Context, tx *SuiTx) (big.Int, error)
}

// SuiGasManager is a concrete implementation of the GasManager interface.
// It manages gas budgets for transactions by applying a fixed-percent increase to the current budget,
// while ensuring that the new budget does not exceed the specified maximum gas budget.
type SuiGasManager struct {
	lggr               logger.Logger
	maxGasBudget       big.Int
	percentualIncrease int64
	ptbClient          client.SuiPTBClient
}

var _ GasManager = (*SuiGasManager)(nil)

// NewSuiGasManager creates a new SuiGasManager instance.
//
// Parameters:
//   - lggr: Logger instance for recording gas management events.
//   - maxGasBudget: The maximum gas budget permitted (as a big.Int).
//   - percentualIncrase: The percentage increase to apply when bumping the gas budget (not used).
//
// Returns:
//   - *SuiGasManager: A pointer to the initialized SuiGasManager.
func NewSuiGasManager(lggr logger.Logger, ptbClient client.SuiPTBClient, maxGasBudget big.Int, percentualIncrase int64) *SuiGasManager {
	if percentualIncrase == 0 {
		percentualIncrase = gasLimitPercentualIncrease
	}

	return &SuiGasManager{
		lggr:               lggr,
		maxGasBudget:       maxGasBudget,
		percentualIncrease: percentualIncrase,
		ptbClient:          ptbClient,
	}
}

// EstimateGasBudget estimates the gas budget for a transaction. Note that this is not an entirely
// gas-less operation, it requires a very small amount of gas to estimate the gas budget.
//
// Parameters:
//   - ctx: Context allowing cancellation and timeouts.
//   - tx: The Sui transaction for which the gas budget is being estimated.
//
// Returns:
//   - uint64: The estimated gas budget.
//   - error: An error if the estimation fails.
func (s *SuiGasManager) EstimateGasBudget(ctx context.Context, tx *SuiTx) (uint64, error) {
	gasBudget, err := s.ptbClient.EstimateGas(ctx, tx.Payload)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas budget: %w", err)
	}

	return gasBudget, nil
}

// GasBump computes a new gas budget for a transaction by increasing its current gas limit.
// The new budget is determined by multiplying the current budget by gasLimitPercentualIncrease and
// then dividing by percentualNormalization. If the current gas limit already meets or exceeds the
// maximum allowed budget, the function returns an error. If the computed budget exceeds the maximum,
// it is capped at maxGasBudget.
//
// Parameters:
//   - ctx: Context allowing cancellation and timeouts.
//   - tx: The Sui transaction whose gas limit is being increased.
//
// Returns:
//   - big.Int: The new gas budget.
//   - error: An error if the gas budget is already at or above the maximum allowed.
func (s *SuiGasManager) GasBump(ctx context.Context, tx *SuiTx) (big.Int, error) {
	// Check if the current gas limit is at or above the maximum allowed budget.
	if tx.Metadata.GasLimit.Cmp(&s.maxGasBudget) >= 0 {
		return *big.NewInt(0), errors.New("gas budget is already at max gas budget")
	}

	// Calculate the new gas budget: newBudget = currentGasLimit * gasLimitPercentualIncrease / percentualNormalization.
	newBudget := new(big.Int).Mul(tx.Metadata.GasLimit, big.NewInt(gasLimitPercentualIncrease))
	newBudget.Div(newBudget, big.NewInt(percentualNormalization))

	// Cap the new budget at maxGasBudget if it exceeds the allowed maximum.
	if newBudget.Cmp(&s.maxGasBudget) > 0 {
		newBudget.Set(&s.maxGasBudget)
	}

	return *newBudget, nil
}
