package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/ops/link"
)

type DeployLinkTokenConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
}

var _ cldf.ChangeSetV2[DeployLinkTokenConfig] = DeployLinkToken{}

type DeployLinkToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployLinkToken) Apply(e cldf.Environment, config DeployLinkTokenConfig) (cldf.ChangesetOutput, error) {
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

	// Run DeployLinkToken Operation
	linkTokenReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy LinkToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save LinkToken address to the addressbook
	typeAndVersionLinkToken := cldf.NewTypeAndVersion(deployment.SuiLinkTokenType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, linkTokenReport.Output.PackageId, typeAndVersionLinkToken)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken address %s for Sui chain %d: %w", linkTokenReport.Output.PackageId, config.ChainSelector, err)
	}

	// save LinkTokenCoinMetadataID address to the addressbook
	typeAndVersionCoinMetadataID := cldf.NewTypeAndVersion(deployment.SuiLinkTokenObjectMetadataID, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, linkTokenReport.Output.Objects.CoinMetadataObjectId, typeAndVersionCoinMetadataID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken CoinmetadataObjectID address %s for Sui chain %d: %w", linkTokenReport.Output.Objects.CoinMetadataObjectId, config.ChainSelector, err)
	}

	// save LinkTokenTreasuryCapID address to the addressbook
	typeAndVersionTreasuryCapID := cldf.NewTypeAndVersion(deployment.SuiLinkTokenTreasuryCapID, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, linkTokenReport.Output.Objects.TreasuryCapObjectId, typeAndVersionTreasuryCapID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken TreasuryCapObjectID address %s for Sui chain %d: %w", linkTokenReport.Output.Objects.TreasuryCapObjectId, config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployLinkToken) VerifyPreconditions(e cldf.Environment, config DeployLinkTokenConfig) error {
	return nil
}
