package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type BlockVersionConfig struct {
	SuiChainSelector uint64
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	Version          uint8
}

var _ cldf.ChangeSetV2[BlockVersionConfig] = BlockVersion{}

type BlockVersion struct{}

// Apply implements deployment.ChangeSetV2.
func (d BlockVersion) Apply(e cldf.Environment, config BlockVersionConfig) (cldf.ChangesetOutput, error) {
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

	BlockVersionInitializeOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.BlockVersionOp, deps, ccipops.BlockVersionInput{
		CCIPPackageId:    config.CCIPPackageId,
		StateObjectId:    config.StateObjectId,
		OwnerCapObjectId: config.OwnerCapObjectId,
		ModuleName:       config.ModuleName,
		Version:          config.Version,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize upgrade registry for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{BlockVersionInitializeOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d BlockVersion) VerifyPreconditions(e cldf.Environment, config BlockVersionConfig) error {
	return nil
}
