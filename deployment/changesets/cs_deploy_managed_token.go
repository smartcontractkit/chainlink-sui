package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
)

type DeployManagedTokenConfig struct {
	managedtokenops.DeployAndInitManagedTokenInput
	ChainSelector uint64 `yaml:"chainSelector"`
	// ReplaceExisting allows this changeset to take datastore keys that are already recorded,
	// as a redeploy of an existing managed token does. Without it, an occupied key is an error
	// raised before anything is deployed.
	ReplaceExisting bool `yaml:"replaceExisting"`
}

var _ cldf.ChangeSetV2[DeployManagedTokenConfig] = DeployManagedToken{}

// DeployAptosChain deploys Sui chain packages and modules
type DeployManagedToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployManagedToken) Apply(e cldf.Environment, config DeployManagedTokenConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]

	deployerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
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

	config.MCMSAddress = state[config.ChainSelector].MCMSPackageID
	config.MCMSOwnerAddress = deployerAddr

	// Run DeployManagedToken Operation
	managedTokenReport, err := operations.ExecuteSequence(e.OperationsBundle, managedtokenops.DeployAndInitManagedTokenSequence, deps, config.DeployAndInitManagedTokenInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ManagedToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save the new managed token package id address to the addressbook
	typeAndVersionManagedTokenPackageID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPackageIDType, deployment.Version1_0_0)
	typeAndVersionManagedTokenPackageID.AddLabel(managedTokenReport.Output.TokenSymbol)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, managedTokenReport.Output.ManagedTokenPackageId, typeAndVersionManagedTokenPackageID, deployment.TokenQualifier(managedTokenReport.Output.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken address %s for Sui chain %d: %w", managedTokenReport.Output.ManagedTokenPackageId, config.ChainSelector, err)
	}

	// save ManagedTokenOwnerCapObjectID address to the addressbook
	typeAndVersionOwnerCapObjectID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenOwnerCapObjectID, deployment.Version1_0_0)
	typeAndVersionOwnerCapObjectID.AddLabel(managedTokenReport.Output.TokenSymbol)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, managedTokenReport.Output.Objects.OwnerCapObjectId, typeAndVersionOwnerCapObjectID, deployment.TokenQualifier(managedTokenReport.Output.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken OwnerCapObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.OwnerCapObjectId, config.ChainSelector, err)
	}

	if config.MinterAddress != "" {
		// save ManagedTokenMinterCapID address to the addressbook
		typeAndVersionMinterCapID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenMinterCapID, deployment.Version1_0_0)
		typeAndVersionMinterCapID.AddLabel(managedTokenReport.Output.TokenSymbol)
		// holder-scoped: this token can have a cap per minter, so the symbol alone is not a key
		minterCapQualifier := deployment.MinterCapQualifier(managedTokenReport.Output.TokenSymbol, config.MinterAddress)
		err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, managedTokenReport.Output.Objects.MinterCapObjectId, typeAndVersionMinterCapID, minterCapQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken MinterCapObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.MinterCapObjectId, config.ChainSelector, err)
		}
	}

	// save ManagedTokenStateObjectID address to the addressbook
	typeAndVersionStateObjectID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenStateObjectID, deployment.Version1_0_0)
	typeAndVersionStateObjectID.AddLabel(managedTokenReport.Output.TokenSymbol)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, managedTokenReport.Output.Objects.StateObjectId, typeAndVersionStateObjectID, deployment.TokenQualifier(managedTokenReport.Output.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken StateObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.StateObjectId, config.ChainSelector, err)
	}

	// save PublisherObjectId address to the addressbook
	typeAndVersionPublisherObjectId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPublisherObjectId, deployment.Version1_0_0)
	typeAndVersionPublisherObjectId.AddLabel(managedTokenReport.Output.TokenSymbol)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, managedTokenReport.Output.Objects.PublisherObjectId, typeAndVersionPublisherObjectId, deployment.TokenQualifier(managedTokenReport.Output.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken PublisherObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.PublisherObjectId, config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployManagedToken) VerifyPreconditions(e cldf.Environment, config DeployManagedTokenConfig) error {
	return deployment.ValidateNoDatastoreConflicts(e, config.ChainSelector, config.ReplaceExisting,
		func() ([]deployment.PlannedRef, error) {
			symbol, err := coinSymbol(e, config.ChainSelector, config.CoinObjectTypeArg)
			if err != nil {
				return nil, err
			}
			qualifier := deployment.TokenQualifier(symbol)

			planned := []deployment.PlannedRef{
				{Type: deployment.SuiManagedTokenPackageIDType, Qualifier: qualifier},
				{Type: deployment.SuiManagedTokenOwnerCapObjectID, Qualifier: qualifier},
				{Type: deployment.SuiManagedTokenStateObjectID, Qualifier: qualifier},
				{Type: deployment.SuiManagedTokenPublisherObjectId, Qualifier: qualifier},
			}
			if config.MinterAddress != "" {
				planned = append(planned, deployment.PlannedRef{
					Type:      deployment.SuiManagedTokenMinterCapID,
					Qualifier: deployment.MinterCapQualifier(symbol, config.MinterAddress),
				})
			}

			return planned, nil
		})
}
