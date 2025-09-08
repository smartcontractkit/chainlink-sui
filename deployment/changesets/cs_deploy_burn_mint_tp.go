package changesets

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
)

type DeploySuiBurnMintTpConfig struct {
	CCIPPackageId string
	MCMSPackageId string

	ChainSelector             uint64
	RemoteChainSelector       uint64
	TokenPoolAddress          string
	CCIPObjectRefId           string
	LinkTokenPkgId            string
	LinkTokenObjectMetadataId string
	LinkTokenTreasuryCapId    string

	EVMToken     common.Address
	EVMTokenPool common.Address
}

var _ cldf.ChangeSetV2[DeploySuiBurnMintTpConfig] = DeploySuiBurnMintTp{}

// DeploySuiBurnMintTp deploys Sui burn mint token pool
type DeploySuiBurnMintTp struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeploySuiBurnMintTp) Apply(e cldf.Environment, config DeploySuiBurnMintTpConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[config.ChainSelector]
	suiSigner := suiChain.Signer

	signerAddr, err := suiSigner.GetAddress()
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiSigner,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	// Deploy BurnMint TP on SUI
	deployBurnMintTp, err := operations.ExecuteSequence(e.OperationsBundle, burnminttokenpoolops.DeployAndInitBurnMintTokenPoolSequence, deps,
		burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
			BurnMintTokenPoolDeployInput: burnminttokenpoolops.BurnMintTokenPoolDeployInput{
				CCIPPackageId:          config.CCIPPackageId,
				CCIPTokenPoolPackageId: config.TokenPoolAddress,
				MCMSAddress:            config.MCMSPackageId,
				MCMSOwnerAddress:       signerAddr,
			},

			CoinObjectTypeArg:      config.LinkTokenPkgId + "::link::LINK",
			CCIPObjectRefObjectId:  config.CCIPObjectRefId,
			CoinMetadataObjectId:   config.LinkTokenObjectMetadataId,
			TreasuryCapObjectId:    config.LinkTokenTreasuryCapId,
			TokenPoolAdministrator: signerAddr,

			// apply dest chain updates
			RemoteChainSelectorsToRemove: []uint64{},
			RemoteChainSelectorsToAdd:    []uint64{config.RemoteChainSelector},
			RemotePoolAddressesToAdd:     [][]string{{config.EVMTokenPool.String()}},
			RemoteTokenAddressesToAdd: []string{
				config.EVMToken.String(),
			},
			// set chain rate limiter configs
			RemoteChainSelectors: []uint64{config.RemoteChainSelector},
			OutboundIsEnableds:   []bool{false},
			OutboundCapacities:   []uint64{10000000000},
			OutboundRates:        []uint64{100},
			InboundIsEnableds:    []bool{false},
			InboundCapacities:    []uint64{10000000000},
			InboundRates:         []uint64{100},
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy BurnMintTP for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save BnM TokenPool to addressbook
	typeAndVersionBurnMintTokenPool := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, deployBurnMintTp.Output.BurnMintTPPackageID, typeAndVersionBurnMintTokenPool)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BurnMintTokenPool address %s for Sui chain %d: %w", deployBurnMintTp.Output.BurnMintTPPackageID, config.ChainSelector, err)
	}

	// save BnM TokenPool State to addressbook
	typeAndVersionBnMTpState := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolStateType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, deployBurnMintTp.Output.Objects.StateObjectId, typeAndVersionBnMTpState)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BurnMintTokenPoolState address %s for Sui chain %d: %w", deployBurnMintTp.Output.Objects.StateObjectId, config.ChainSelector, err)
	}
	seqReports = append(seqReports, deployBurnMintTp.ExecutionReports...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeploySuiBurnMintTp) VerifyPreconditions(e cldf.Environment, config DeploySuiBurnMintTpConfig) error {
	return nil
}
