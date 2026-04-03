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
	SuiChainSelector      uint64 `yaml:"suiChainSelector"`
	CCIPPackageId         string `yaml:"ccipPackageId"`
	CCIPObjectRefObjectId string `yaml:"ccipObjectRefObjectId"`
	OwnerCapObjectId      string `yaml:"ownerCapObjectId"`
	ModuleName            string `yaml:"moduleName"`
	FunctionName          string `yaml:"functionName"`
	Version               uint8  `yaml:"version"`
}

var _ cldf.ChangeSetV2[BlockFunctionConfig] = BlockFunction{}

type BlockFunction struct{}

func (d BlockFunction) Apply(e cldf.Environment, config BlockFunctionConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()

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

	report, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.BlockFunctionOp, deps, ccipops.BlockFunctionInput{
		CCIPPackageId:    config.CCIPPackageId,
		StateObjectId:    config.CCIPObjectRefObjectId,
		OwnerCapObjectId: config.OwnerCapObjectId,
		ModuleName:       config.ModuleName,
		FunctionName:     config.FunctionName,
		Version:          config.Version,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to block function for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

func (d BlockFunction) VerifyPreconditions(e cldf.Environment, config BlockFunctionConfig) error {
	return nil
}
