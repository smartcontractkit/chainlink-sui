package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type UpgradeRegistryConfig struct {
	SuiChainSelector uint64
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
}

var _ cldf.ChangeSetV2[UpgradeRegistryConfig] = UpgradeRegistry{}

type UpgradeRegistry struct{}

// Apply implements deployment.ChangeSetV2.
func (d UpgradeRegistry) Apply(e cldf.Environment, config UpgradeRegistryConfig) (cldf.ChangesetOutput, error) {
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

	UpgradeRegistryInitializeOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.UpgradeRegistryInitializeOp, deps, ccipops.InitUpgradeRegistryInput{
		CCIPPackageId:    config.CCIPPackageId,
		StateObjectId:    config.StateObjectId,
		OwnerCapObjectId: config.OwnerCapObjectId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize upgrade registry for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	// save UpgradeRegistryObjectId address to the addressbook
	typeAndVersionUpgradeRegistryObjectId := cldf.NewTypeAndVersion(deployment.SuiUpgradeRegistryObjectId, deployment.Version1_0_0)
	err = ab.Save(config.SuiChainSelector, UpgradeRegistryInitializeOp.Output.Objects.UpgradeRegistryObjectId, typeAndVersionUpgradeRegistryObjectId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save UpgradeRegistryInitializeOp address %s for Sui chain %d: %w", UpgradeRegistryInitializeOp.Output.Objects.UpgradeRegistryObjectId, config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{UpgradeRegistryInitializeOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d UpgradeRegistry) VerifyPreconditions(e cldf.Environment, config UpgradeRegistryConfig) error {
	return nil
}
