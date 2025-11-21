package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type BlockFunctionConfig struct {
	SuiChainSelector uint64
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	FunctionName     string
	Version          uint8
}

var _ cldf.ChangeSetV2[BlockFunctionConfig] = BlockFunction{}

type BlockFunction struct{}

// Apply implements deployment.ChangeSetV2.
func (d BlockFunction) Apply(e cldf.Environment, config BlockFunctionConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	BlockFunctionInitializeOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.BlockFunctionOp, deps, ccipops.BlockFunctionInput{
		CCIPPackageId:    config.CCIPPackageId,
		StateObjectId:    config.StateObjectId,
		OwnerCapObjectId: config.OwnerCapObjectId,
		ModuleName:       config.ModuleName,
		FunctionName:     config.FunctionName,
		Version:          config.Version,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to block Functionname for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{BlockFunctionInitializeOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d BlockFunction) VerifyPreconditions(e cldf.Environment, config BlockFunctionConfig) error {
	return nil
}
