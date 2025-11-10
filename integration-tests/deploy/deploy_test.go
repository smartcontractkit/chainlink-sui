//go:build integration

package e2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"slices"
	"testing"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/suite"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

type RoleConfig struct {
	Role   suisdk.TimelockRole
	Count  int
	Quorum uint8
	Keys   []*ecdsa.PrivateKey
	Config *types.Config
}

func CreateConfig(role suisdk.TimelockRole, count int, quorum uint8) *RoleConfig {
	signers := make([]common.Address, count)
	signerKeys := make([]*ecdsa.PrivateKey, count)

	for i := range signers {
		signerKeys[i], _ = crypto.GenerateKey()
		signers[i] = crypto.PubkeyToAddress(signerKeys[i].PublicKey)
	}
	slices.SortFunc(signers[:], func(a, b common.Address) int {
		return a.Cmp(b)
	})

	return &RoleConfig{
		Role:   role,
		Count:  count,
		Quorum: quorum,
		Keys:   signerKeys,
		Config: &types.Config{
			Quorum:  quorum,
			Signers: signers[:],
		},
	}
}

type DeployTestSuite struct {
	suite.Suite
	lggr   logger.Logger
	signer bindutils.SuiSigner
	client sui.ISuiAPI
	env    cldf.Environment
}

func (s *DeployTestSuite) SetupSuite() {
	s.signer, s.client = testenv.SetupEnvironment(s.T())
	s.lggr = logger.Test(s.T())

	// Setup operation registry
	ops := make([]*cld_ops.Operation[any, any, any], len(opregistry.AllOperations))
	for i := range opregistry.AllOperations {
		ops[i] = &opregistry.AllOperations[i]
	}
	registry := cld_ops.NewOperationRegistry(ops...)

	bundle := cld_ops.NewBundle(
		func() context.Context { return s.T().Context() },
		s.lggr,
		cld_ops.NewMemoryReporter(),
		cld_ops.WithOperationRegistry(registry),
	)

	s.env = cldf.Environment{
		Name:              "test",
		Logger:            s.lggr,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: chain.NewBlockChains(
			map[uint64]chain.BlockChain{
				cselectors.SUI_LOCALNET.Selector: cldfsui.Chain{
					ChainMetadata: cldfsui.ChainMetadata{
						Selector: cselectors.SUI_LOCALNET.Selector,
					},
					Client: s.client,
					Signer: s.signer,
				},
			}),
		OperationsBundle: bundle,
	}
}

// TestDeployAndConfigureSuiChain tests the deployment and configuration of a Sui chain, lane and BnM TP
// using views to verify the deployments.
func TestDeployAndConfigureSuiChain(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DeployTestSuite))
}

func (s *DeployTestSuite) TestDeployAndConfigureSuiChain() {
	// "Phase 1: Deploy MCMS"
	s.DeployMCMS()
	// "Phase 2: Deploy LINK Token"
	s.DeployLink()
	// "Phase 3: Deploy Core CCIP Infrastructure"
	s.DeployCCIPCore()
	// "Phase 4: Connect Lanes"
	s.ConnectLanes()
	// "Phase 5: Deploy Token Pools"
	s.DeployTokenPools()

	// Load view and check deployments
	state, err := deployment.LoadOnchainStatesui(s.env)
	s.Require().NoError(err, "failed to load on-chain state")
	view, err := state[cselectors.SUI_LOCALNET.Selector].GenerateView(&s.env, cselectors.SUI_LOCALNET.Selector, "sui_localnet")
	s.Require().NoError(err, "failed to generate on-chain view")
	_ = view
	// TODO: Add assertions to verify deployed contracts and configurations
}

func (s *DeployTestSuite) DeployMCMS() {
	s.T().Log("Phase 1: Deploying MCMS...")

	deployInput := mcmsops.DeployMCMSSeqInput{
		ChainSelector: cselectors.SUI_LOCALNET.Selector,
	}

	// Execute MCMS deployment changeset
	out, err := changesets.DeployMCMS{}.Apply(s.env, deployInput)
	s.Require().NoError(err, "failed to deploy MCMS")
	// Update existing addresses
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge MCMS addresses")
}

func (s *DeployTestSuite) DeployLink() {
	s.T().Log("Phase 2: Deploying LINK Token...")

	deployInput := changesets.DeployLinkTokenConfig{
		ChainSelector: cselectors.SUI_LOCALNET.Selector,
	}

	// Execute LINK token deployment changeset
	out, err := changesets.DeployLinkToken{}.Apply(s.env, deployInput)
	// Update existing addresses
	s.Require().NoError(err, "failed to deploy LINK token")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge LINK token addresses")
}

func (s *DeployTestSuite) DeployCCIPCore() {
	s.T().Log("Phase 3: Deploying Core CCIP Infrastructure...")

	// Get LINK token CoinMetadata from address book
	addresses, err := s.env.ExistingAddresses.AddressesForChain(cselectors.SUI_LOCALNET.Selector)
	s.Require().NoError(err, "failed to get addresses")

	var linkCoinMetadataObjectId string
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiLinkTokenObjectMetadataID {
			linkCoinMetadataObjectId = addr
			break
		}
	}
	s.Require().NotEmpty(linkCoinMetadataObjectId, "LINK CoinMetadata not found")

	// For testing, use dummy values for destination chain
	// In real deployment, these would come from the EVM chain deployment
	destChainSelector := uint64(11155111) // Sepolia
	// Use a valid dummy EVM address for the destination onramp (20 bytes padded to 32 bytes)
	destChainOnRampBytes := common.Hex2Bytes("000000000000000000000000a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")

	deployInput := changesets.DeploySuiChainConfig{
		SuiChainSelector:              cselectors.SUI_LOCALNET.Selector,
		DestChainSelector:             destChainSelector,
		DestChainOnRampAddressBytes:   destChainOnRampBytes,
		LinkTokenCoinMetadataObjectId: linkCoinMetadataObjectId,
	}

	// Execute CCIP deployment changeset
	out, err := changesets.DeploySuiChain{}.Apply(s.env, deployInput)
	s.Require().NoError(err, "failed to deploy CCIP")
	// Update existing addresses
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge CCIP addresses")
}

func (s *DeployTestSuite) ConnectLanes() {
	s.T().Log("Phase 4: Connecting Lanes...")

	destChainSelector := uint64(11155111) // Sepolia

	// Use default configs from deployment package
	ccipConfig := deployment.DefaultCCIPSeqConfig
	offrampConfig := deployment.DefaultOffRampSeqConfig

	// Dummy onramp address bytes for testing
	destChainOnRampBytes := common.Hex2Bytes("000000000000000000000000a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")

	connectInput := changesets.ConnectSuiToEVMConfig{
		SuiChainSelector: cselectors.SUI_LOCALNET.Selector,
		FeeQuoterApplyTokenTransferFeeConfigUpdatesInput: ccipops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
			DestChainSelector:    destChainSelector,
			AddTokens:            []string{}, // Empty - tokens will be configured after pool deployment
			AddMinFeeUsdCents:    []uint32{}, // Empty
			AddMaxFeeUsdCents:    []uint32{}, // Empty
			AddDeciBps:           []uint16{}, // Empty
			AddDestGasOverhead:   []uint32{}, // Empty
			AddDestBytesOverhead: []uint32{}, // Empty
			AddIsEnabled:         []bool{},   // Empty
			RemoveTokens:         []string{}, // Empty
		},
		FeeQuoterApplyDestChainConfigUpdatesInput: ccipops.FeeQuoterApplyDestChainConfigUpdatesInput{
			DestChainSelector:                 destChainSelector,
			IsEnabled:                         ccipConfig.IsEnabled,
			MaxNumberOfTokensPerMsg:           ccipConfig.MaxNumberOfTokensPerMsg,
			MaxDataBytes:                      ccipConfig.MaxDataBytes,
			MaxPerMsgGasLimit:                 ccipConfig.MaxPerMsgGasLimit,
			DestGasOverhead:                   ccipConfig.DestGasOverhead,
			DestGasPerPayloadByteBase:         ccipConfig.DestGasPerPayloadByteBase,
			DestGasPerPayloadByteHigh:         ccipConfig.DestGasPerPayloadByteHigh,
			DestGasPerPayloadByteThreshold:    ccipConfig.DestGasPerPayloadByteThreshold,
			DestDataAvailabilityOverheadGas:   ccipConfig.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    ccipConfig.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: ccipConfig.DestDataAvailabilityMultiplierBps,
			DefaultTokenFeeUsdCents:           ccipConfig.DefaultTokenFeeUsdCents,
			DefaultTokenDestGasOverhead:       ccipConfig.DefaultTokenDestGasOverhead,
			DefaultTxGasLimit:                 ccipConfig.DefaultTxGasLimit,
			GasMultiplierWeiPerEth:            ccipConfig.GasMultiplierWeiPerEth,
			GasPriceStalenessThreshold:        ccipConfig.GasPriceStalenessThreshold,
			NetworkFeeUsdCents:                ccipConfig.NetworkFeeUsdCents,
			EnforceOutOfOrder:                 ccipConfig.EnforceOutOfOrder,
			ChainFamilySelector:               ccipConfig.ChainFamilySelector,
		},
		FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput: ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
			Tokens:                     []string{}, // Empty - tokens will be configured after pool deployment
			PremiumMultiplierWeiPerEth: []uint64{}, // Empty - must match Tokens length
		},
		ApplyDestChainConfigureOnRampInput: onrampops.ApplyDestChainConfigureOnRampInput{
			DestChainSelector:         []uint64{destChainSelector},
			DestChainAllowListEnabled: []bool{false},                                                                  // Single destination chain, allowlist disabled
			DestChainRouters:          []string{"0x0000000000000000000000000000000000000000000000000000000000000001"}, // Placeholder router address
		},
		ApplySourceChainConfigUpdateInput: offrampops.ApplySourceChainConfigUpdateInput{
			SourceChainsSelectors:                 []uint64{destChainSelector},
			SourceChainsOnRamp:                    [][]byte{destChainOnRampBytes},
			SourceChainsIsEnabled:                 offrampConfig.InitializeOffRampInput.SourceChainsIsEnabled,
			SourceChainsIsRMNVerificationDisabled: offrampConfig.InitializeOffRampInput.SourceChainsIsRMNVerificationDisabled,
		},
	}

	// Execute lane connection changeset
	_, err := changesets.ConnectSuiToEVM{}.Apply(s.env, connectInput)
	s.Require().NoError(err, "failed to connect lanes")
}

func (s *DeployTestSuite) DeployTokenPools() {
	s.T().Log("Phase 5: Deploying Token Pools...")

	// Get LINK token details from address book
	s.T().Log("Retrieving LINK token information...")

	addresses, err := s.env.ExistingAddresses.AddressesForChain(cselectors.SUI_LOCALNET.Selector)
	s.Require().NoError(err, "failed to get addresses")

	var (
		linkTokenPackageId     string
		linkTokenMetadataId    string
		linkTokenTreasuryCapId string
		ccipPackageId          string
		ccipObjectRef          string
		mcmsPackageId          string
	)

	// Extract addresses from address book
	for addr, typeAndVersion := range addresses {
		switch typeAndVersion.Type {
		case deployment.SuiLinkTokenType:
			linkTokenPackageId = addr
		case deployment.SuiLinkTokenObjectMetadataID:
			linkTokenMetadataId = addr
		case deployment.SuiLinkTokenTreasuryCapID:
			linkTokenTreasuryCapId = addr
		case deployment.SuiCCIPType:
			ccipPackageId = addr
		case deployment.SuiCCIPObjectRefType:
			ccipObjectRef = addr
		case deployment.SuiMcmsPackageIDType:
			mcmsPackageId = addr
		}
	}

	s.Require().NotEmpty(linkTokenPackageId, "LINK token package ID not found")
	s.Require().NotEmpty(linkTokenMetadataId, "LINK token metadata ID not found")
	s.Require().NotEmpty(linkTokenTreasuryCapId, "LINK token treasury cap ID not found")
	s.Require().NotEmpty(ccipPackageId, "CCIP package ID not found")
	s.Require().NotEmpty(ccipObjectRef, "CCIP object ref not found")
	s.Require().NotEmpty(mcmsPackageId, "MCMS package ID not found")

	// Get deployer address
	deployerAddr, err := s.signer.GetAddress()
	s.Require().NoError(err, "failed to get deployer address")

	// Get the coin type argument for LINK token
	coinTypeArg := fmt.Sprintf("%s::link::LINK", linkTokenPackageId)

	// Deploy BurnMint token pool for LINK
	s.T().Log("Deploying BurnMint token pool for LINK...")

	destChainSelector := uint64(11155111) // Sepolia

	// Use real EVM addresses for pool and token (examples from Sepolia testnet)
	// In production, these would be the actual deployed pool and token addresses on the destination chain
	evmPoolAddress := "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"  // Example EVM pool address
	evmTokenAddress := "0x779877a7b0d9e8603169ddbd7836e478b4624789" // Example LINK token address on Sepolia

	tokenPoolInput := changesets.DeployTPAndConfigureConfig{
		SuiChainSelector: cselectors.SUI_LOCALNET.Selector,
		TokenPoolTypes:   []string{"bnm"},
		BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
			BurnMintTokenPoolDeployInput: burnminttokenpoolops.BurnMintTokenPoolDeployInput{
				CCIPPackageId:    ccipPackageId,
				MCMSAddress:      mcmsPackageId,
				MCMSOwnerAddress: deployerAddr,
			},
			CoinObjectTypeArg:      coinTypeArg,
			CCIPObjectRefObjectId:  ccipObjectRef,
			CoinMetadataObjectId:   linkTokenMetadataId,
			TreasuryCapObjectId:    linkTokenTreasuryCapId,
			TokenPoolAdministrator: deployerAddr,
			// Configure remote chain connections
			RemoteChainSelectorsToRemove: []uint64{},
			RemoteChainSelectorsToAdd:    []uint64{destChainSelector},
			RemotePoolAddressesToAdd:     [][]string{{evmPoolAddress}}, // Real EVM pool address
			RemoteTokenAddressesToAdd:    []string{evmTokenAddress},    // Real EVM token address
			// Rate limiter configs
			RemoteChainSelectors: []uint64{destChainSelector},
			OutboundIsEnableds:   []bool{false},                 // Disabled for testing
			OutboundCapacities:   []uint64{1000000000000000000}, // 1B LINK (18 decimals)
			OutboundRates:        []uint64{100000000000000},     // 100k LINK per second
			InboundIsEnableds:    []bool{false},                 // Disabled for testing
			InboundCapacities:    []uint64{1000000000000000000}, // 1B LINK (18 decimals)
			InboundRates:         []uint64{100000000000000},     // 100k LINK per second
		},
	}

	tokenPoolOut, err := changesets.DeployTPAndConfigure{}.Apply(s.env, tokenPoolInput)
	s.Require().NoError(err, "failed to deploy LINK token pool")
	err = s.env.ExistingAddresses.Merge(tokenPoolOut.AddressBook)
	s.Require().NoError(err, "failed to merge LINK token pool addresses")
}
