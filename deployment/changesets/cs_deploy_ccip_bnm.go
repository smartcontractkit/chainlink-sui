package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	bnmops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token"
)

const CCIPBnMSymbol = "CCIP BnM"

type DeployCCIPBnMTokenConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	MintAmount    uint64 `yaml:"mintAmount"`
	MintToAddress string `yaml:"mintToAddress"`
}

var _ cldf.ChangeSetV2[DeployCCIPBnMTokenConfig] = DeployCCIPBnMToken{}

// DeployCCIPBnMToken deploys Sui chain packages and modules
type DeployCCIPBnMToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployCCIPBnMToken) Apply(e cldf.Environment, config DeployCCIPBnMTokenConfig) (cldf.ChangesetOutput, error) {
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
		SuiRPC: suiChain.URL,
	}

	// Run DeployCCIPBnMToken Operation
	ccipBnMTokenReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, bnmops.DeployBnMOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIPBnMToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save CCIPBnMToken package ID to the addressbook
	typeAndVersionCCIPBnMToken := cldf.NewTypeAndVersion(deployment.SuiManagedTokenType, deployment.Version1_0_0)
	typeAndVersionCCIPBnMToken.AddLabel(CCIPBnMSymbol)
	err = ab.Save(config.ChainSelector, ccipBnMTokenReport.Output.PackageId, typeAndVersionCCIPBnMToken)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPBnMToken package ID %s for Sui chain %d: %w", ccipBnMTokenReport.Output.PackageId, config.ChainSelector, err)
	}

	// save CCIPBnMTokenCoinMetadataId address to the addressbook
	typeAndVersionCoinMetadataId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenCoinMetadataIDType, deployment.Version1_0_0)
	typeAndVersionCoinMetadataId.AddLabel(CCIPBnMSymbol)
	err = ab.Save(config.ChainSelector, ccipBnMTokenReport.Output.Objects.CoinMetadataObjectId, typeAndVersionCoinMetadataId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPBnMToken CoinmetadataObjectId address %s for Sui chain %d: %w", ccipBnMTokenReport.Output.Objects.CoinMetadataObjectId, config.ChainSelector, err)
	}

	// save CCIPBnMTokenTreasuryCapId address to the addressbook
	typeAndVersionTreasuryCapId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenTreasuryCapIDType, deployment.Version1_0_0)
	typeAndVersionTreasuryCapId.AddLabel(CCIPBnMSymbol)
	err = ab.Save(config.ChainSelector, ccipBnMTokenReport.Output.Objects.TreasuryCapObjectId, typeAndVersionTreasuryCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPBnMToken TreasuryCapObjectId address %s for Sui chain %d: %w", ccipBnMTokenReport.Output.Objects.TreasuryCapObjectId, config.ChainSelector, err)
	}

	// save CCIPBnMTokenUpgradeCapId address to the addressbook
	typeAndVersionUpgradeCapId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenUpgradeCapIDType, deployment.Version1_0_0)
	typeAndVersionUpgradeCapId.AddLabel(CCIPBnMSymbol)
	err = ab.Save(config.ChainSelector, ccipBnMTokenReport.Output.Objects.UpgradeCapObjectId, typeAndVersionUpgradeCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPBnMToken UpgradeCapObjectId address %s for Sui chain %d: %w", ccipBnMTokenReport.Output.Objects.UpgradeCapObjectId, config.ChainSelector, err)
	}

	if config.MintAmount != 0 || config.MintToAddress != "" {
		// Run MintCCIPBnMToken Operation
		_, err = cld_ops.ExecuteOperation(e.OperationsBundle, bnmops.MintBnMOp, deps, bnmops.MintBnMTokenInput{
			BnMTokenPackageId: ccipBnMTokenReport.Output.PackageId,
			TreasuryCapId:     ccipBnMTokenReport.Output.Objects.TreasuryCapObjectId,
			Amount:            config.MintAmount,
			ToAddress:         config.MintToAddress,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to mint CCIPBnMToken for Sui chain %d: %w", config.ChainSelector, err)
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployCCIPBnMToken) VerifyPreconditions(e cldf.Environment, config DeployCCIPBnMTokenConfig) error {
	return nil
}
