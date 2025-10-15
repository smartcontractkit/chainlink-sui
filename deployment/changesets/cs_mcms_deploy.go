package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/mcms"
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

	err = storeMCMSInAddressBook(ab, config.ChainSelector, mcmsReport.Output)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to store MCMS in address book for Sui chain %d: %w", config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		Reports:               seqReports,
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsReport.Output.AcceptOwnershipProposal},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployMCMS) VerifyPreconditions(e cldf.Environment, config mcmsops.DeployMCMSSeqInput) error {
	return nil
}

func storeMCMSInAddressBook(ab *cldf.AddressBookMap, chainSelector uint64, mcmsReport mcmsops.DeployMCMSSeqOutput) error {
	// save MCMS address to the addressbook
	typeAndVersionMCMS := cldf.NewTypeAndVersion(deployment.SuiMcmsPackageIDType, deployment.Version1_0_0)
	err := ab.Save(chainSelector, mcmsReport.PackageId, typeAndVersionMCMS)
	if err != nil {
		return fmt.Errorf("failed to save MCMS address %s for Sui chain %d: %w", mcmsReport.PackageId, chainSelector, err)
	}

	// save MCMS MultisigState object ID to the addressbook
	typeAndVersionMCMSObject := cldf.NewTypeAndVersion(deployment.SuiMcmsObjectIDType, deployment.Version1_0_0)
	err = ab.Save(chainSelector, mcmsReport.Objects.McmsMultisigStateObjectId, typeAndVersionMCMSObject)
	if err != nil {
		return fmt.Errorf("failed to save MCMS MultisigState object ID %s for Sui chain %d: %w", mcmsReport.Objects.McmsMultisigStateObjectId, chainSelector, err)
	}

	// save MCMS Registry object ID to the addressbook
	typeAndVersionMCMSRegistry := cldf.NewTypeAndVersion(deployment.SuiMcmsRegistryObjectIDType, deployment.Version1_0_0)
	err = ab.Save(chainSelector, mcmsReport.Objects.McmsRegistryObjectId, typeAndVersionMCMSRegistry)
	if err != nil {
		return fmt.Errorf("failed to save MCMS Registry object ID %s for Sui chain %d: %w", mcmsReport.Objects.McmsRegistryObjectId, chainSelector, err)
	}

	// save MCMS AccountState object ID to the addressbook
	typeAndVersionMCMSAccountState := cldf.NewTypeAndVersion(deployment.SuiMcmsAccountStateObjectIDType, deployment.Version1_0_0)
	err = ab.Save(chainSelector, mcmsReport.Objects.McmsAccountStateObjectId, typeAndVersionMCMSAccountState)
	if err != nil {
		return fmt.Errorf("failed to save MCMS AccountState object ID %s for Sui chain %d: %w", mcmsReport.Objects.McmsAccountStateObjectId, chainSelector, err)
	}

	// save MCMS AccountOwnerCap object ID to the addressbook
	typeAndVersionMCMSAccountOwnerCap := cldf.NewTypeAndVersion(deployment.SuiMcmsAccountOwnerCapObjectIDType, deployment.Version1_0_0)
	err = ab.Save(chainSelector, mcmsReport.Objects.McmsAccountOwnerCapObjectId, typeAndVersionMCMSAccountOwnerCap)
	if err != nil {
		return fmt.Errorf("failed to save MCMS AccountOwnerCap object ID %s for Sui chain %d: %w", mcmsReport.Objects.McmsAccountOwnerCapObjectId, chainSelector, err)
	}

	// save MCMS Timelock object ID to the addressbook
	typeAndVersionMCMSTimelock := cldf.NewTypeAndVersion(deployment.SuiMcmsTimelockObjectIDType, deployment.Version1_0_0)
	err = ab.Save(chainSelector, mcmsReport.Objects.TimelockObjectId, typeAndVersionMCMSTimelock)
	if err != nil {
		return fmt.Errorf("failed to save MCMS Timelock object ID %s for Sui chain %d: %w", mcmsReport.Objects.TimelockObjectId, chainSelector, err)
	}

	return nil
}
