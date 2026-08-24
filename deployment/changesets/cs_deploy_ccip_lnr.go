package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	lnrops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token"
)

const CCIPLnRSymbol = "CCIP LnR"

type DeployCCIPLnRTokenConfig struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	MintAmount    uint64 `yaml:"mintAmount"`
	MintToAddress string `yaml:"mintToAddress"`
}

var _ cldf.ChangeSetV2[DeployCCIPLnRTokenConfig] = DeployCCIPLnRToken{}

// DeployCCIPLnRToken deploys Sui chain packages and modules
type DeployCCIPLnRToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployCCIPLnRToken) Apply(e cldf.Environment, config DeployCCIPLnRTokenConfig) (cldf.ChangesetOutput, error) {
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

	// Run DeployCCIPLnRToken Operation
	ccipLnRTokenReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, lnrops.DeployLnROp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy CCIPLnRToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save CCIPLnRToken package ID to the addressbook
	typeAndVersionCCIPLnRToken := cldf.NewTypeAndVersion(deployment.SuiManagedTokenType, deployment.Version1_0_0)
	typeAndVersionCCIPLnRToken.AddLabel(CCIPLnRSymbol)
	typeAndVersionCCIPLnRToken.AddLabel("coinType=" + deployment.SuiCCIPLnRCoinTypeSuffix)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, ccipLnRTokenReport.Output.PackageId, typeAndVersionCCIPLnRToken)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPLnRToken package ID %s for Sui chain %d: %w", ccipLnRTokenReport.Output.PackageId, config.ChainSelector, err)
	}

	// save CCIPLnRTokenCoinMetadataId address to the addressbook
	typeAndVersionCoinMetadataId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenCoinMetadataIDType, deployment.Version1_0_0)
	typeAndVersionCoinMetadataId.AddLabel(CCIPLnRSymbol)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, ccipLnRTokenReport.Output.Objects.CoinMetadataObjectId, typeAndVersionCoinMetadataId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPLnRToken CoinmetadataObjectId address %s for Sui chain %d: %w", ccipLnRTokenReport.Output.Objects.CoinMetadataObjectId, config.ChainSelector, err)
	}

	// save CCIPLnRTokenTreasuryCapId address to the addressbook
	typeAndVersionTreasuryCapId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenTreasuryCapIDType, deployment.Version1_0_0)
	typeAndVersionTreasuryCapId.AddLabel(CCIPLnRSymbol)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, ccipLnRTokenReport.Output.Objects.TreasuryCapObjectId, typeAndVersionTreasuryCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPLnRToken TreasuryCapObjectId address %s for Sui chain %d: %w", ccipLnRTokenReport.Output.Objects.TreasuryCapObjectId, config.ChainSelector, err)
	}

	// save CCIPLnRTokenUpgradeCapId address to the addressbook
	typeAndVersionUpgradeCapId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenUpgradeCapIDType, deployment.Version1_0_0)
	typeAndVersionUpgradeCapId.AddLabel(CCIPLnRSymbol)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.ChainSelector, ccipLnRTokenReport.Output.Objects.UpgradeCapObjectId, typeAndVersionUpgradeCapId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save CCIPLnRToken UpgradeCapObjectId address %s for Sui chain %d: %w", ccipLnRTokenReport.Output.Objects.UpgradeCapObjectId, config.ChainSelector, err)
	}

	if config.MintAmount != 0 || config.MintToAddress != "" {
		// Run MintCCIPLnRToken Operation
		_, err = cld_ops.ExecuteOperation(e.OperationsBundle, lnrops.MintLnROp, deps, lnrops.MintLnRTokenInput{
			LnRTokenPackageId: ccipLnRTokenReport.Output.PackageId,
			TreasuryCapId:     ccipLnRTokenReport.Output.Objects.TreasuryCapObjectId,
			Amount:            config.MintAmount,
			ToAddress:         config.MintToAddress,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to mint CCIPLnRToken for Sui chain %d: %w", config.ChainSelector, err)
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployCCIPLnRToken) VerifyPreconditions(e cldf.Environment, config DeployCCIPLnRTokenConfig) error {
	return nil
}
