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

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
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
	// MaxGasBudget returns the maximum gas budget permitted.
	MaxGasBudget() *big.Int

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

	// EstimateGasBudgetFromCCIPMessage estimates the gas budget required for the given CCIP message.
	//
	// Parameters:
	//   - ctx: Context allowing cancellation and timeouts.
	//   - message: The CCIP message for which to estimate the gas budget.
	//
	// Returns:
	//   - uint64: The estimated gas budget.
	//   - error: An error if estimation is not implemented or fails.
	CalculateOfframpExecuteGasBudget(ctx context.Context, arguments offramp.SuiOffRampExecCallArgs) (*big.Int, error)

	// GasBump calculates a new gas budget for the given transaction by increasing its current gas limit
	// by the configured percentage (defaults to 20%). If the new budget would exceed the manager's
	// maximum allowed budget, it is capped at that maximum, and if the current gas limit is already at
	// or above the maximum, an error is returned.
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
//   - percentualIncrase: The percentage (relative to 100) to scale the gas budget by when bumping.
//     Values at or below 100 would shrink or keep the budget, so they fall back to the default.
//
// Returns:
//   - *SuiGasManager: A pointer to the initialized SuiGasManager.
func NewSuiGasManager(lggr logger.Logger, ptbClient client.SuiPTBClient, maxGasBudget big.Int, percentualIncrase int64) *SuiGasManager {
	if percentualIncrase <= percentualNormalization {
		percentualIncrase = gasLimitPercentualIncrease
	}

	return &SuiGasManager{
		lggr:               lggr,
		maxGasBudget:       maxGasBudget,
		percentualIncrease: percentualIncrase,
		ptbClient:          ptbClient,
	}
}

// MaxGasBudget returns the maximum gas budget permitted.
func (s *SuiGasManager) MaxGasBudget() *big.Int {
	return &s.maxGasBudget
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
	gasBudget, err := s.ptbClient.EstimateGas(ctx, tx.Ptb)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas budget: %w", err)
	}

	return gasBudget, nil
}

// EstimateGasBudgetFromCCIPMessage estimates the gas budget required for the given CCIP message.
// If the message contains ExtraArgs, it attempts to ABI decode them as SuiExtraArgsV1 struct
// and extract the gas limit from the decoded data.
//
// Parameters:
//   - ctx: Context allowing cancellation and timeouts.
//   - message: The CCIP message for which to estimate the gas budget.
//
// Returns:
//   - uint64: The estimated gas budget.
//   - error: An error if estimation is not implemented or fails.
func (s *SuiGasManager) CalculateOfframpExecuteGasBudget(ctx context.Context, arguments offramp.SuiOffRampExecCallArgs) (*big.Int, error) {
	gasLimit := big.NewInt(0)
	// we can't assume extraArgs gasLimit to be bigInt all the time
	if val, ok := arguments.ExtraData.ExtraArgsDecoded["gasLimit"]; ok {
		var gl *big.Int

		switch v := val.(type) {
		case *big.Int:
			gl = v
		case uint64:
			gl = new(big.Int).SetUint64(v)
		default:
			return nil, fmt.Errorf("gasLimit is not *big.Int or uint64, got %T", val)
		}

		if gl == nil {
			return nil, fmt.Errorf("gasLimit is nil")
		}

		gasLimit.Add(gasLimit, gl)
	}

	for _, destExecData := range arguments.ExtraData.DestExecDataDecoded {
		if val, ok := destExecData["destGasAmount"]; ok {
			if destGasAmount, ok := val.(uint64); ok {
				gasLimit.Add(gasLimit, big.NewInt(int64(destGasAmount)))
			} else {
				return nil, fmt.Errorf("destGasAmount in DestExecDataDecoded is not uint64, got %T", val)
			}
		}
	}

	return gasLimit, nil
}

// GasBump increases the gas budget for a given transaction by applying the configured percentage bump.
// The new gas budget is calculated as:
//
//	newBudget = currentGasLimit * percentualIncrease / percentualNormalization
//
// The new budget is capped at the gas manager's configured maximum (s.maxGasBudget). If the current
// gas limit is already at or above that maximum, an error is returned and no bump occurs — retrying
// at an unchanged budget would deterministically fail again.
//
// Parameters:
//   - ctx: Context for cancellation and timeouts.
//   - tx:  The SuiTx transaction whose gas limit should be increased.
//
// Returns:
//   - big.Int: The new gas budget to use.
//   - error:   An error if the gas budget is already at or above the allowed maximum.
func (s *SuiGasManager) GasBump(ctx context.Context, tx *SuiTx) (big.Int, error) {
	txGasLimit := tx.Metadata.GasLimit
	maxGasLimit := &s.maxGasBudget

	s.lggr.Debugw("GasBump", "txGasLimit", txGasLimit, "maxGasLimit", maxGasLimit)

	// A budget already at (or above) the maximum cannot be increased any further.
	if txGasLimit.Cmp(maxGasLimit) >= 0 {
		return *big.NewInt(0), errors.New("gas budget is already at max gas limit")
	}

	newBudget := new(big.Int).Mul(txGasLimit, big.NewInt(s.percentualIncrease))
	newBudget.Div(newBudget, big.NewInt(percentualNormalization))

	// Cap the new budget at maxGasBudget if it exceeds the allowed maximum.
	if newBudget.Cmp(maxGasLimit) > 0 {
		newBudget.Set(maxGasLimit)
	}

	return *newBudget, nil
}
