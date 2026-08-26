package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
)

type DeployLinkTokenConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	// ReplaceExisting allows this changeset to take datastore keys that are already
	// recorded, as a redeploy of this token does. Without it, an occupied key is an error
	// raised before anything is deployed.
	ReplaceExisting bool `yaml:"replaceExisting"`
}

var _ cldf.ChangeSetV2[DeployLinkTokenConfig] = DeployLinkToken{}

// DeployLinkToken deploys Sui chain packages and modules
type DeployLinkToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployLinkToken) Apply(e cldf.Environment, config DeployLinkTokenConfig) (cldf.ChangesetOutput, error) {
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

	// Run DeployLinkToken Operation
	linkTokenReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy LinkToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save LinkToken address to the addressbook
	typeAndVersionLinkToken := cldf.NewTypeAndVersion(deployment.SuiLinkTokenType, deployment.Version1_0_0)
	typeAndVersionLinkToken.AddLabel("LINK")
	typeAndVersionLinkToken.AddLabel("coinType=" + deployment.SuiLinkCoinTypeSuffix)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, linkTokenReport.Output.PackageId, typeAndVersionLinkToken, deployment.TokenQualifier("LINK")) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken address %s for Sui chain %d: %w", linkTokenReport.Output.PackageId, config.ChainSelector, err)
	}

	// save LinkTokenCoinMetadataId address to the addressbook
	typeAndVersionCoinMetadataId := cldf.NewTypeAndVersion(deployment.SuiLinkTokenObjectMetadataID, deployment.Version1_0_0)
	typeAndVersionCoinMetadataId.AddLabel("LINK")
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, linkTokenReport.Output.Objects.CoinMetadataObjectId, typeAndVersionCoinMetadataId, deployment.TokenQualifier("LINK")) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken CoinmetadataObjectId address %s for Sui chain %d: %w", linkTokenReport.Output.Objects.CoinMetadataObjectId, config.ChainSelector, err)
	}

	// save LinkTokenTreasuryCapId address to the addressbook
	typeAndVersionTreasuryCapId := cldf.NewTypeAndVersion(deployment.SuiLinkTokenTreasuryCapID, deployment.Version1_0_0)
	typeAndVersionTreasuryCapId.AddLabel("LINK")
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, linkTokenReport.Output.Objects.TreasuryCapObjectId, typeAndVersionTreasuryCapId, deployment.TokenQualifier("LINK")) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken TreasuryCapObjectId address %s for Sui chain %d: %w", linkTokenReport.Output.Objects.TreasuryCapObjectId, config.ChainSelector, err)
	}

	// save LinkTokenUpgradeCapId address to the addressbook
	typeAndVersionUpgradeCapID := cldf.NewTypeAndVersion(deployment.SuiLinkTokenUpgradeCapID, deployment.Version1_0_0)
	typeAndVersionUpgradeCapID.AddLabel("LINK")
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, linkTokenReport.Output.Objects.UpgradeCapObjectId, typeAndVersionUpgradeCapID, deployment.TokenQualifier("LINK")) // token-scoped; the label keeps the display name, the key uses the symbol form
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken UpgradeCapObjectId address %s for Sui chain %d: %w", linkTokenReport.Output.Objects.UpgradeCapObjectId, config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployLinkToken) VerifyPreconditions(e cldf.Environment, config DeployLinkTokenConfig) error {
	qualifier := deployment.TokenQualifier("LINK")

	return deployment.ValidateNoDatastoreConflicts(e, config.ChainSelector, config.ReplaceExisting,
		func() ([]deployment.PlannedRef, error) {
			return []deployment.PlannedRef{
				{Type: deployment.SuiLinkTokenType, Qualifier: qualifier},
				{Type: deployment.SuiLinkTokenObjectMetadataID, Qualifier: qualifier},
				{Type: deployment.SuiLinkTokenTreasuryCapID, Qualifier: qualifier},
				{Type: deployment.SuiLinkTokenUpgradeCapID, Qualifier: qualifier},
			}, nil
		})
}
