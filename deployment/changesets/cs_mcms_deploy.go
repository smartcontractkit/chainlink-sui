package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[mcmsops.DeployMCMSSeqInput] = DeployMCMS{}

type DeployMCMS struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployMCMS) Apply(e cldf.Environment, config mcmsops.DeployMCMSSeqInput) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]

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
	}

	// Run DeployMCMS Sequence
	mcmsReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps, config)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save MCMS address to the addressbook
	typeAndVersionMCMS := cldf.NewTypeAndVersion(deployment.SuiMcmsPackageIDType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, mcmsReport.Output.PackageId, typeAndVersionMCMS)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS address %s for Sui chain %d: %w", mcmsReport.Output.PackageId, config.ChainSelector, err)
	}

	// save MCMS MultisigState object ID to the addressbook
	typeAndVersionMCMSObject := cldf.NewTypeAndVersion(deployment.SuiMcmsObjectIDType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, mcmsReport.Output.Objects.McmsMultisigStateObjectId, typeAndVersionMCMSObject)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS MultisigState object ID %s for Sui chain %d: %w", mcmsReport.Output.Objects.McmsMultisigStateObjectId, config.ChainSelector, err)
	}

	// save MCMS Registry object ID to the addressbook
	typeAndVersionMCMSRegistry := cldf.NewTypeAndVersion(deployment.SuiMcmsRegistryObjectIdType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, mcmsReport.Output.Objects.McmsRegistryObjectId, typeAndVersionMCMSRegistry)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS Registry object ID %s for Sui chain %d: %w", mcmsReport.Output.Objects.McmsRegistryObjectId, config.ChainSelector, err)
	}

	// save MCMS AccountState object ID to the addressbook
	typeAndVersionMCMSAccountState := cldf.NewTypeAndVersion(deployment.SuiMcmsAccountStateObjectIdType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, mcmsReport.Output.Objects.McmsAccountStateObjectId, typeAndVersionMCMSAccountState)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS AccountState object ID %s for Sui chain %d: %w", mcmsReport.Output.Objects.McmsAccountStateObjectId, config.ChainSelector, err)
	}

	// save MCMS AccountOwnerCap object ID to the addressbook
	typeAndVersionMCMSAccountOwnerCap := cldf.NewTypeAndVersion(deployment.SuiMcmsAccountOwnerCapObjectIdType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, mcmsReport.Output.Objects.McmsAccountOwnerCapObjectId, typeAndVersionMCMSAccountOwnerCap)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS AccountOwnerCap object ID %s for Sui chain %d: %w", mcmsReport.Output.Objects.McmsAccountOwnerCapObjectId, config.ChainSelector, err)
	}

	// save MCMS Timelock object ID to the addressbook
	typeAndVersionMCMSTimelock := cldf.NewTypeAndVersion(deployment.SuiMcmsTimelockObjectIdType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, mcmsReport.Output.Objects.TimelockObjectId, typeAndVersionMCMSTimelock)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS Timelock object ID %s for Sui chain %d: %w", mcmsReport.Output.Objects.TimelockObjectId, config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployMCMS) VerifyPreconditions(e cldf.Environment, config mcmsops.DeployMCMSSeqInput) error {
	return nil
}
