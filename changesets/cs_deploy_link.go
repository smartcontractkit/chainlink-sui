package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/ops/link"
)

type DeployLinkTokenConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
}

var _ cldf.ChangeSetV2[DeployLinkTokenConfig] = DeployLinkToken{}

// DeployAptosChain deploys Aptos chain packages and modules
type DeployLinkToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployLinkToken) Apply(e cldf.Environment, config DeployLinkTokenConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

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
	linkTokenReport, err := operations.ExecuteOperation(e.OperationsBundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy LinkToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save LinkToken address to the addressbook
	typeAndVersionLinkToken := cldf.NewTypeAndVersion(SuiLinkTokenType, Version1_0_0)
	err = ab.Save(config.ChainSelector, linkTokenReport.Output.PackageId, typeAndVersionLinkToken)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken address %s for Sui chain %d: %w", linkTokenReport.Output.PackageId, config.ChainSelector, err)
	}

	// save LinkTokenCoinMetadataId address to the addressbook
	typeAndVersionCoinMetadataId := cldf.NewTypeAndVersion(SuiLinkTokenObjectMetadataId, Version1_0_0)
	err = ab.Save(config.ChainSelector, linkTokenReport.Output.Objects.CoinMetadataObjectId, typeAndVersionCoinMetadataId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken CoinmetadataObjectId address %s for Sui chain %d: %w", linkTokenReport.Output.Objects.CoinMetadataObjectId, config.ChainSelector, err)
	}

	// save LinkTokenTreasuryCapId address to the addressbook
	typeAndVersionTreasuryCapId := cldf.NewTypeAndVersion(SuiLinkTokenTreasuryCapId, Version1_0_0)
	err = ab.Save(config.ChainSelector, linkTokenReport.Output.Objects.TreasuryCapObjectId, typeAndVersionTreasuryCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LinkToken TreasuryCapObjectId address %s for Sui chain %d: %w", linkTokenReport.Output.Objects.TreasuryCapObjectId, config.ChainSelector, err)
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
