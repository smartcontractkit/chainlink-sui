//go:build integration

package e2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/ethereum/go-ethereum/common"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
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
	"github.com/smartcontractkit/chainlink-sui/deployment/view"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// Test configuration constants
var (
	// Chain selectors
	SuiChainSelector = cselectors.SUI_LOCALNET.Selector
	EVMChainSelector = cselectors.ETHEREUM_TESTNET_SEPOLIA.Selector

	// EVM addresses for destination chain (examples from Sepolia testnet)
	DestChainOnRampAddress      = "000000000000000000000000a0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	DestChainOnRampAddressBytes = common.Hex2Bytes(DestChainOnRampAddress)
	DestChainRouterAddress      = "0x0000000000000000000000000000000000000000000000000000000000000001"
	EVMPoolAddress              = "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"
	EVMTokenAddress             = "779877a7b0d9e8603169ddbd7836e478b4624789" // LINK on Sepolia

	// Token pool rate limiter configs
	RateLimiterCapacity = uint64(1000000000000000000) // 1B LINK (18 decimals)
	RateLimiterRate     = uint64(100000000000000)     // 100k LINK per second

	// Mock MCMS Signers
	MCMSMockSigners = []common.Address{
		common.HexToAddress("0xa000000000000000000000000000000000000001"),
		common.HexToAddress("0xa000000000000000000000000000000000000002"),
	}
)

type DeployTestSuite struct {
	suite.Suite
	lggr   logger.Logger
	signer bindutils.SuiSigner
	client sui.ISuiAPI
	env    cldf.Environment

	// Cached deployment addresses
	linkTokenPackageID     string
	linkTokenMetadataID    string
	linkTokenTreasuryCapID string
	ccipPackageID          string
	ccipObjectRef          string
	mcmsPackageID          string
	deployerAddr           string
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
	states, err := deployment.LoadOnchainStatesui(s.env)
	state := states[cselectors.SUI_LOCALNET.Selector]
	s.Require().NoError(err, "failed to load on-chain state")
	actualView, err := state.GenerateView(&s.env, cselectors.SUI_LOCALNET.Selector, "sui_localnet")
	s.Require().NoError(err, "failed to generate on-chain view")
	owner, err := s.signer.GetAddress()
	s.Require().NoError(err, "failed to get signer address")
	_ = owner
	expectedView := deployment.SuiChainView{
		ChainSelector: SuiChainSelector,
		ChainID:       "",
		MCMSWithTimelock: view.MCMSWithTimelockView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.MCMSPackageID,
				Owner:          owner,
				TypeAndVersion: "MCMS 1.6.0",
			},
			Bypasser: types.Config{
				Quorum:       1,
				Signers:      MCMSMockSigners,
				GroupSigners: []types.Config{},
			},
			Proposer: types.Config{
				Quorum:       1,
				Signers:      MCMSMockSigners,
				GroupSigners: []types.Config{},
			},
			Canceller: types.Config{
				Quorum:       2,
				Signers:      MCMSMockSigners,
				GroupSigners: []types.Config{},
			},
			TimelockMinDelay:         0,
			TimelockBlockedFunctions: []view.TimelockBlockedFunction{},
		},
		CCIP: view.CCIPView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.CCIPAddress,
				Owner:          owner,
				TypeAndVersion: "",
			},
			FeeQuoter: view.FeeQuoterView{
				ContractMetaData: view.ContractMetaData{
					Address:        state.CCIPAddress,
					Owner:          owner,
					TypeAndVersion: "FeeQuoter 1.6.0",
				},
				FeeTokens: []string{state.LinkTokenCoinMetadataId},
				StaticConfig: view.FeeQuoterStaticConfig{
					MaxFeeJuelsPerMsg:            deployment.DefaultCCIPSeqConfig.MaxFeeJuelsPerMsg,
					LinkToken:                    state.LinkTokenCoinMetadataId,
					TokenPriceStalenessThreshold: deployment.DefaultCCIPSeqConfig.TokenPriceStalenessThreshold,
				},
				DestinationChainConfigs: map[uint64]view.FeeQuoterDestChainConfig{}, // TODO: this is not being populated somehow, fix
			},
			RMNRemote: view.RMNRemoteView{
				ContractMetaData: view.ContractMetaData{
					Address:        state.CCIPAddress,
					Owner:          owner,
					TypeAndVersion: "RMNRemote 1.6.0",
				},
				IsCursed:             false,
				Config:               view.RMNRemoteVersionedConfig{},
				CursedSubjectEntries: []view.RMNRemoteCurseEntry{},
			},
			TokenAdminRegistry: view.TokenAdminRegistryView{
				ContractMetaData: view.ContractMetaData{
					Address:        state.CCIPAddress,
					Owner:          owner,
					TypeAndVersion: "TokenAdminRegistry 1.6.0",
				},
				TokenConfigs: map[string]view.TokenConfigView{
					state.LinkTokenCoinMetadataId: {
						TokenPoolPackageId:  state.BnMTokenPools["LINK"].PackageID,
						TokenPoolModule:     "burn_mint_token_pool",
						TokenType:           fmt.Sprintf("%s::link::LINK", strings.Replace(state.LinkTokenAddress, "0x", "", 1)),
						Administrator:       owner,
						TokenPoolTypeProof:  fmt.Sprintf("%s::burn_mint_token_pool::TypeProof", strings.Replace(state.BnMTokenPools["LINK"].PackageID, "0x", "", 1)),
						LockOrBurnParams:    []string{"0x0000000000000000000000000000000000000000000000000000000000000006", state.BnMTokenPools["LINK"].StateObjectId},
						ReleaseOrMintParams: []string{"0x0000000000000000000000000000000000000000000000000000000000000006", state.BnMTokenPools["LINK"].StateObjectId},
					},
				},
			},
			NonceManager: view.NonceManagerView{
				ContractMetaData: view.ContractMetaData{
					Address:        state.CCIPAddress,
					Owner:          owner,
					TypeAndVersion: "NonceManager 1.6.0",
				},
			},
			ReceiverRegistry: view.ReceiverRegistryView{
				ContractMetaData: view.ContractMetaData{
					Address:        state.CCIPAddress,
					Owner:          owner,
					TypeAndVersion: "ReceiverRegistry 1.6.0",
				},
			},
		},
		OnRamp: view.OnRampView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.OnRampAddress,
				Owner:          owner,
				TypeAndVersion: "OnRamp 1.6.0",
			},
			StaticConfig: view.OnRampStaticConfig{
				ChainSelector: SuiChainSelector,
			},
			DynamicConfig: view.OnRampDynamicConfig{
				FeeAggregator:  owner,
				AllowlistAdmin: owner,
			},
			DestChainSpecificData: map[uint64]view.DestChainSpecificData{}, // TODO: find out why this is not being populated
		},
		OffRamp: view.OffRampView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.OffRampAddress,
				Owner:          owner,
				TypeAndVersion: "OffRamp 1.6.0",
			},
			StaticConfig: view.OffRampStaticConfig{
				ChainSelector:      SuiChainSelector,
				RMNRemote:          state.CCIPAddress,
				TokenAdminRegistry: state.CCIPAddress,
				NonceManager:       state.CCIPAddress,
			},
			DynamicConfig: view.OffRampDynamicConfig{
				FeeQuoter:                               state.CCIPAddress,
				PermissionlessExecutionThresholdSeconds: 28800, // TODO: why does it defaulting to this value?
			},
			SourceChainConfigs: map[uint64]view.OffRampSourceChainConfig{
				EVMChainSelector: {
					Router:                    state.CCIPAddress,
					IsEnabled:                 true,
					MinSeqNr:                  1,
					IsRMNVerificationDisabled: true,
					OnRamp:                    fmt.Sprintf("0x%s", DestChainOnRampAddress),
				},
			},
		},
		Router: view.RouterView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.CCIPRouterAddress,
				Owner:          owner,
				TypeAndVersion: "Router 1.6.0",
			},
			IsTestRouter: false,
			OnRamps:      nil,
			OffRamps:     nil,
		},
		Tokens: nil,
		TokenPools: map[string]map[string]view.TokenPoolView{
			"LINK": {
				state.BnMTokenPools["LINK"].PackageID: {
					ContractMetaData: view.ContractMetaData{
						Address:        state.BnMTokenPools["LINK"].PackageID,
						Owner:          owner,
						TypeAndVersion: "BurnMintTokenPool 1.6.0",
					},
					Token: s.linkTokenMetadataID,
					RemoteChainConfigs: map[uint64]view.RemoteChainConfig{
						EVMChainSelector: {
							RemoteTokenAddress:  fmt.Sprintf("0x000000000000000000000000%s", EVMTokenAddress),
							RemotePoolAddresses: []string{EVMPoolAddress},
							InboundRateLimiterConfig: view.RateLimiterConfig{
								IsEnabled: false,
								Capacity:  RateLimiterCapacity,
								Rate:      RateLimiterRate,
							},
							OutboundRateLimiterConfig: view.RateLimiterConfig{
								IsEnabled: false,
								Capacity:  RateLimiterCapacity,
								Rate:      RateLimiterRate,
							},
						},
					},
					AllowList:        []string{},
					AllowListEnabled: false,
				},
			},
		},
	}

	s.Require().Equal(expectedView, actualView)
}

func (s *DeployTestSuite) DeployMCMS() {
	s.T().Log("Phase 1: Deploying MCMS...")

	out, err := changesets.DeployMCMS{}.Apply(s.env, mcmsops.DeployMCMSSeqInput{
		ChainSelector: SuiChainSelector,
		Bypasser: &types.Config{
			Quorum:       1,
			Signers:      MCMSMockSigners,
			GroupSigners: []types.Config{},
		},
		Proposer: &types.Config{
			Quorum:       1,
			Signers:      MCMSMockSigners,
			GroupSigners: []types.Config{},
		},
		Canceller: &types.Config{
			Quorum:       2,
			Signers:      MCMSMockSigners,
			GroupSigners: []types.Config{},
		},
	})
	s.Require().NoError(err, "failed to deploy MCMS")

	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge MCMS addresses")
}

func (s *DeployTestSuite) DeployLink() {
	s.T().Log("Phase 2: Deploying LINK Token...")

	out, err := changesets.DeployLinkToken{}.Apply(s.env, changesets.DeployLinkTokenConfig{
		ChainSelector: SuiChainSelector,
	})
	s.Require().NoError(err, "failed to deploy LINK token")

	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge LINK token addresses")
}

func (s *DeployTestSuite) DeployCCIPCore() {
	s.T().Log("Phase 3: Deploying Core CCIP Infrastructure...")

	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")

	var linkCoinMetadataObjectID string
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiLinkTokenObjectMetadataID {
			linkCoinMetadataObjectID = addr
			break
		}
	}
	s.Require().NotEmpty(linkCoinMetadataObjectID, "LINK CoinMetadata not found")

	out, err := changesets.DeploySuiChain{}.Apply(s.env, changesets.DeploySuiChainConfig{
		SuiChainSelector:              SuiChainSelector,
		DestChainSelector:             EVMChainSelector,
		DestChainOnRampAddressBytes:   DestChainOnRampAddressBytes,
		LinkTokenCoinMetadataObjectId: linkCoinMetadataObjectID,
	})
	s.Require().NoError(err, "failed to deploy CCIP")

	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge CCIP addresses")
}

func (s *DeployTestSuite) ConnectLanes() {
	s.T().Log("Phase 4: Connecting Lanes...")

	ccipConfig := deployment.DefaultCCIPSeqConfig
	offrampConfig := deployment.DefaultOffRampSeqConfig

	_, err := changesets.ConnectSuiToEVM{}.Apply(s.env, changesets.ConnectSuiToEVMConfig{
		SuiChainSelector: SuiChainSelector,
		FeeQuoterApplyTokenTransferFeeConfigUpdatesInput: ccipops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
			DestChainSelector:    EVMChainSelector,
			AddTokens:            []string{},
			AddMinFeeUsdCents:    []uint32{},
			AddMaxFeeUsdCents:    []uint32{},
			AddDeciBps:           []uint16{},
			AddDestGasOverhead:   []uint32{},
			AddDestBytesOverhead: []uint32{},
			AddIsEnabled:         []bool{},
			RemoveTokens:         []string{},
		},
		FeeQuoterApplyDestChainConfigUpdatesInput: ccipops.FeeQuoterApplyDestChainConfigUpdatesInput{
			DestChainSelector:                 EVMChainSelector,
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
			Tokens:                     []string{},
			PremiumMultiplierWeiPerEth: []uint64{},
		},
		ApplyDestChainConfigureOnRampInput: onrampops.ApplyDestChainConfigureOnRampInput{
			DestChainSelector:         []uint64{EVMChainSelector},
			DestChainAllowListEnabled: []bool{false},
			DestChainRouters:          []string{DestChainRouterAddress},
		},
		ApplySourceChainConfigUpdateInput: offrampops.ApplySourceChainConfigUpdateInput{
			SourceChainsSelectors:                 []uint64{EVMChainSelector},
			SourceChainsOnRamp:                    [][]byte{DestChainOnRampAddressBytes},
			SourceChainsIsEnabled:                 offrampConfig.InitializeOffRampInput.SourceChainsIsEnabled,
			SourceChainsIsRMNVerificationDisabled: offrampConfig.InitializeOffRampInput.SourceChainsIsRMNVerificationDisabled,
		},
	})
	s.Require().NoError(err, "failed to connect lanes")
}

func (s *DeployTestSuite) loadDeploymentAddresses() {
	s.T().Log("Loading deployment addresses...")

	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")

	for addr, typeAndVersion := range addresses {
		switch typeAndVersion.Type {
		case deployment.SuiLinkTokenType:
			s.linkTokenPackageID = addr
		case deployment.SuiLinkTokenObjectMetadataID:
			s.linkTokenMetadataID = addr
		case deployment.SuiLinkTokenTreasuryCapID:
			s.linkTokenTreasuryCapID = addr
		case deployment.SuiCCIPType:
			s.ccipPackageID = addr
		case deployment.SuiCCIPObjectRefType:
			s.ccipObjectRef = addr
		case deployment.SuiMcmsPackageIDType:
			s.mcmsPackageID = addr
		}
	}

	s.deployerAddr, err = s.signer.GetAddress()
	s.Require().NoError(err, "failed to get deployer address")

	s.Require().NotEmpty(s.linkTokenPackageID, "LINK token package ID not found")
	s.Require().NotEmpty(s.linkTokenMetadataID, "LINK token metadata ID not found")
	s.Require().NotEmpty(s.linkTokenTreasuryCapID, "LINK token treasury cap ID not found")
	s.Require().NotEmpty(s.ccipPackageID, "CCIP package ID not found")
	s.Require().NotEmpty(s.ccipObjectRef, "CCIP object ref not found")
	s.Require().NotEmpty(s.mcmsPackageID, "MCMS package ID not found")
}

func (s *DeployTestSuite) DeployTokenPools() {
	s.T().Log("Phase 5: Deploying Token Pools...")

	s.loadDeploymentAddresses()

	coinTypeArg := fmt.Sprintf("%s::link::LINK", s.linkTokenPackageID)

	tokenPoolOut, err := changesets.DeployTPAndConfigure{}.Apply(s.env, changesets.DeployTPAndConfigureConfig{
		SuiChainSelector: SuiChainSelector,
		TokenPoolTypes:   []string{"bnm"},
		BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
			BurnMintTokenPoolDeployInput: burnminttokenpoolops.BurnMintTokenPoolDeployInput{
				CCIPPackageId:    s.ccipPackageID,
				MCMSAddress:      s.mcmsPackageID,
				MCMSOwnerAddress: s.deployerAddr,
			},
			CoinObjectTypeArg:      coinTypeArg,
			CCIPObjectRefObjectId:  s.ccipObjectRef,
			CoinMetadataObjectId:   s.linkTokenMetadataID,
			TreasuryCapObjectId:    s.linkTokenTreasuryCapID,
			TokenPoolAdministrator: s.deployerAddr,
			// Remote chain configuration
			RemoteChainSelectorsToRemove: []uint64{},
			RemoteChainSelectorsToAdd:    []uint64{EVMChainSelector},
			RemotePoolAddressesToAdd:     [][]string{{EVMPoolAddress}},
			RemoteTokenAddressesToAdd:    []string{fmt.Sprintf("0x%s", EVMTokenAddress)},
			// Rate limiter configs
			RemoteChainSelectors: []uint64{EVMChainSelector},
			OutboundIsEnableds:   []bool{false},
			OutboundCapacities:   []uint64{RateLimiterCapacity},
			OutboundRates:        []uint64{RateLimiterRate},
			InboundIsEnableds:    []bool{false},
			InboundCapacities:    []uint64{RateLimiterCapacity},
			InboundRates:         []uint64{RateLimiterRate},
		},
	})
	s.Require().NoError(err, "failed to deploy LINK token pool")

	err = s.env.ExistingAddresses.Merge(tokenPoolOut.AddressBook)
	s.Require().NoError(err, "failed to merge LINK token pool addresses")
}
