package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
)

// Match the gas budget used by the existing Sui changesets.
const suiPreconditionGasBudget uint64 = 400_000_000

// suiOpDeps builds dependencies for precondition read operations.
func suiOpDeps(e cldf.Environment, chainSelector uint64) (sui_ops.OpTxDeps, error) {
	suiChain, ok := e.BlockChains.SuiChains()[chainSelector]
	if !ok {
		return sui_ops.OpTxDeps{}, fmt.Errorf("sui chain %d not found", chainSelector)
	}

	return sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := suiPreconditionGasBudget
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}, nil
}

// coinSymbol reads the symbol used to qualify token datastore refs.
func coinSymbol(e cldf.Environment, chainSelector uint64, coinObjectTypeArg string) (string, error) {
	deps, err := suiOpDeps(e, chainSelector)
	if err != nil {
		return "", err
	}
	report, err := cld_ops.ExecuteOperation(e.OperationsBundle, coin_ops.GetCoinSymbolOp, deps, coinObjectTypeArg)
	if err != nil {
		return "", fmt.Errorf("failed to get coin symbol for %s: %w", coinObjectTypeArg, err)
	}
	if report.Output.Symbol == "" {
		return "", fmt.Errorf("coin %s has no symbol in its on-chain metadata", coinObjectTypeArg)
	}

	return report.Output.Symbol, nil
}
