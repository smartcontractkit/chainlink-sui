package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsuserops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms_user"
)

type DeployMCMSUserConfig struct {
	mcmsuserops.DeployMCMSUserSeqInput
	ChainSelector uint64 `json:"chainSelector"`
	// ReplaceExisting allows this changeset to take datastore keys that are already recorded,
	// as redeploying the MCMS user module does. Without it, an occupied key is an error
	// raised before anything is deployed.
	ReplaceExisting bool `json:"replaceExisting"`
}

var _ cldf.ChangeSetV2[DeployMCMSUserConfig] = DeployMCMSUser{}

type DeployMCMSUser struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployMCMSUser) Apply(e cldf.Environment, config DeployMCMSUserConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
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
		SuiRPC: suiChain.URL,
	}

	// Run DeployMCMSUser Sequence
	mcmsUserReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsuserops.DeployMCMSUserSequence, deps, config.DeployMCMSUserSeqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS User for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save MCMS User package ID to the addressbook
	typeAndVersionMCMSUserPackage := cldf.NewTypeAndVersion(deployment.SuiMcmsUserPackageIDType, deployment.Version1_0_0)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, mcmsUserReport.Output.PackageId, typeAndVersionMCMSUserPackage, deployment.ChainSingletonQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS User package ID %s for Sui chain %d: %w", mcmsUserReport.Output.PackageId, config.ChainSelector, err)
	}

	// save MCMS User Data object ID to the addressbook
	typeAndVersionMCMSUserData := cldf.NewTypeAndVersion(deployment.SuiMcmsUserDataObjectIDType, deployment.Version1_0_0)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, mcmsUserReport.Output.Objects.McmsUserDataObjectID, typeAndVersionMCMSUserData, deployment.ChainSingletonQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS User Data object ID %s for Sui chain %d: %w", mcmsUserReport.Output.Objects.McmsUserDataObjectID, config.ChainSelector, err)
	}

	// save MCMS User OwnerCap object ID to the addressbook
	typeAndVersionMCMSUserOwnerCap := cldf.NewTypeAndVersion(deployment.SuiMcmsUserOwnerCapObjectIDType, deployment.Version1_0_0)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, mcmsUserReport.Output.Objects.McmsUserOwnerCapObjectID, typeAndVersionMCMSUserOwnerCap, deployment.ChainSingletonQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save MCMS User OwnerCap object ID %s for Sui chain %d: %w", mcmsUserReport.Output.Objects.McmsUserOwnerCapObjectID, config.ChainSelector, err)
	}

	// Convert the specific report type to the generic type needed for seqReports
	genericReport := cld_ops.Report[any, any]{
		Input:  mcmsUserReport.Input,
		Output: mcmsUserReport.Output,
	}
	seqReports = append(seqReports, genericReport)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployMCMSUser) VerifyPreconditions(e cldf.Environment, config DeployMCMSUserConfig) error {
	return deployment.ValidateNoDatastoreConflicts(e, config.ChainSelector, config.ReplaceExisting,
		func() ([]deployment.PlannedRef, error) {
			return []deployment.PlannedRef{
				{Type: deployment.SuiMcmsUserPackageIDType, Qualifier: deployment.ChainSingletonQualifier},
				{Type: deployment.SuiMcmsUserDataObjectIDType, Qualifier: deployment.ChainSingletonQualifier},
				{Type: deployment.SuiMcmsUserOwnerCapObjectIDType, Qualifier: deployment.ChainSingletonQualifier},
			}, nil
		})
}
