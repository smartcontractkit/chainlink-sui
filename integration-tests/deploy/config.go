//go:build integration

package deploy

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	"github.com/smartcontractkit/chainlink-sui/deployment/view"

	"github.com/ethereum/go-ethereum/common"

	cselectors "github.com/smartcontractkit/chain-selectors"
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

// Configs for changesets
func GetMCMSConfig(quorum uint8) *types.Config {
	return &types.Config{
		Quorum:       quorum,
		Signers:      MCMSMockSigners,
		GroupSigners: []types.Config{},
	}
}

var (
	ccipConfig             = deployment.DefaultCCIPSeqConfig
	offrampConfig          = deployment.DefaultOffRampSeqConfig
	TokenTransferFeeConfig = ccipops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
		DestChainSelector:    EVMChainSelector,
		AddTokens:            []string{},
		AddMinFeeUsdCents:    []uint32{},
		AddMaxFeeUsdCents:    []uint32{},
		AddDeciBps:           []uint16{},
		AddDestGasOverhead:   []uint32{},
		AddDestBytesOverhead: []uint32{},
		AddIsEnabled:         []bool{},
		RemoveTokens:         []string{},
	}
	DestChainConfigUpdatesInput = ccipops.FeeQuoterApplyDestChainConfigUpdatesInput{
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
	}
	PremiumMultiplierWeiPerEth = ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
		Tokens:                     []string{},
		PremiumMultiplierWeiPerEth: []uint64{},
	}
	DestChainConfigureOnRamp = onrampops.ApplyDestChainConfigureOnRampInput{
		DestChainSelector:         []uint64{EVMChainSelector},
		DestChainAllowListEnabled: []bool{false},
		DestChainRouters:          []string{DestChainRouterAddress},
	}
	SourceChainConfigUpdate = offrampops.ApplySourceChainConfigUpdateInput{
		SourceChainsSelectors:                 []uint64{EVMChainSelector},
		SourceChainsOnRamp:                    [][]byte{DestChainOnRampAddressBytes},
		SourceChainsIsEnabled:                 offrampConfig.InitializeOffRampInput.SourceChainsIsEnabled,
		SourceChainsIsRMNVerificationDisabled: offrampConfig.InitializeOffRampInput.SourceChainsIsRMNVerificationDisabled,
	}
)

// buildExpectedSuiChainView returns the expected SuiChainView for the given state and owner
func buildExpectedSuiChainView(s *DeployTestSuite, state deployment.CCIPChainState, owner string) deployment.SuiChainView {
	return deployment.SuiChainView{
		ChainSelector: SuiChainSelector,
		ChainID:       "",
		MCMSWithTimelock: view.MCMSWithTimelockView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.MCMSPackageID,
				Owner:          owner,
				TypeAndVersion: "MCMS 1.6.0",
				StateObjectID:  state.MCMSStateObjectID,
			},
			Bypasser:                 *GetMCMSConfig(1),
			Proposer:                 *GetMCMSConfig(1),
			Canceller:                *GetMCMSConfig(2),
			TimelockMinDelay:         0,
			TimelockBlockedFunctions: []view.TimelockBlockedFunction{},
		},
		CCIP: view.CCIPView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.CCIPAddress,
				Owner:          owner,
				TypeAndVersion: "",
				StateObjectID:  state.CCIPObjectRef,
			},
			FeeQuoter: view.FeeQuoterView{
				ContractMetaData: view.ContractMetaData{
					Address:        state.CCIPAddress,
					Owner:          owner,
					TypeAndVersion: "FeeQuoter 1.6.0",
					StateObjectID:  state.CCIPObjectRef,
				},
				FeeTokens: []string{state.LinkTokenCoinMetadataId},
				StaticConfig: view.FeeQuoterStaticConfig{
					MaxFeeJuelsPerMsg:            deployment.DefaultCCIPSeqConfig.MaxFeeJuelsPerMsg,
					LinkToken:                    state.LinkTokenCoinMetadataId,
					TokenPriceStalenessThreshold: deployment.DefaultCCIPSeqConfig.TokenPriceStalenessThreshold,
				},
				DestinationChainConfigs: map[uint64]view.FeeQuoterDestChainConfig{
					EVMChainSelector: {
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
						ChainFamilySelector:               fmt.Sprintf("%x", ccipConfig.ChainFamilySelector),
						EnforceOutOfOrder:                 ccipConfig.EnforceOutOfOrder,
						DefaultTokenFeeUsdCents:           ccipConfig.DefaultTokenFeeUsdCents,
						DefaultTokenDestGasOverhead:       ccipConfig.DefaultTokenDestGasOverhead,
						DefaultTxGasLimit:                 ccipConfig.DefaultTxGasLimit,
						GasMultiplierWeiPerEth:            ccipConfig.GasMultiplierWeiPerEth,
						GasPriceStalenessThreshold:        ccipConfig.GasPriceStalenessThreshold,
						NetworkFeeUsdCents:                ccipConfig.NetworkFeeUsdCents,
					},
				},
			},
			RMNRemote: view.RMNRemoteView{
				ContractMetaData: view.ContractMetaData{
					Address:        state.CCIPAddress,
					Owner:          owner,
					TypeAndVersion: "RMNRemote 1.6.0",
					StateObjectID:  state.CCIPObjectRef,
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
					StateObjectID:  state.CCIPObjectRef,
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
					state.ManagedTokens[changesets.CCIPBnMSymbol].TokenCoinMetadataID: {
						TokenPoolPackageId:  state.ManagedTokenPools[changesets.CCIPBnMSymbol].PackageID,
						TokenPoolModule:     "managed_token_pool",
						TokenType:           fmt.Sprintf("%s::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", strings.Replace(state.ManagedTokens[changesets.CCIPBnMSymbol].TokenPackageID, "0x", "", 1)),
						Administrator:       owner,
						TokenPoolTypeProof:  fmt.Sprintf("%s::managed_token_pool::TypeProof", strings.Replace(state.ManagedTokenPools[changesets.CCIPBnMSymbol].PackageID, "0x", "", 1)),
						LockOrBurnParams:    []string{"0x0000000000000000000000000000000000000000000000000000000000000006", "0x0000000000000000000000000000000000000000000000000000000000000403", state.ManagedTokens[changesets.CCIPBnMSymbol].StateObjectId, state.ManagedTokenPools[changesets.CCIPBnMSymbol].StateObjectId},
						ReleaseOrMintParams: []string{"0x0000000000000000000000000000000000000000000000000000000000000006", "0x0000000000000000000000000000000000000000000000000000000000000403", state.ManagedTokens[changesets.CCIPBnMSymbol].StateObjectId, state.ManagedTokenPools[changesets.CCIPBnMSymbol].StateObjectId},
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
		OnRamp: map[string]view.OnRampView{
			state.OnRampAddress: {
				ContractMetaData: view.ContractMetaData{
					Address:        state.OnRampAddress,
					Owner:          owner,
					TypeAndVersion: "OnRamp 1.6.0",
					StateObjectID:  state.OnRampStateObjectId,
				},
				StaticConfig: view.OnRampStaticConfig{
					ChainSelector: SuiChainSelector,
				},
				DynamicConfig: view.OnRampDynamicConfig{
					FeeAggregator:  owner,
					AllowlistAdmin: owner,
				},
				DestChainSpecificData: map[uint64]view.DestChainSpecificData{
					EVMChainSelector: {
						AllowedSendersList: []string{},
						DestChainConfig: view.OnRampDestChainConfig{
							SequenceNumber:   0,
							AllowlistEnabled: DestChainConfigureOnRamp.DestChainAllowListEnabled[0],
							Router:           DestChainConfigureOnRamp.DestChainRouters[0],
						},
						ExpectedNextSeqNum: 1,
					},
				},
			},
		},
		OffRamp: map[string]view.OffRampView{
			state.OffRampAddress: {
				ContractMetaData: view.ContractMetaData{
					Address:        state.OffRampAddress,
					Owner:          owner,
					TypeAndVersion: "OffRamp 1.6.0",
					StateObjectID:  state.OffRampStateObjectId,
				},
				StaticConfig: view.OffRampStaticConfig{
					ChainSelector:      SuiChainSelector,
					RMNRemote:          state.CCIPAddress,
					TokenAdminRegistry: state.CCIPAddress,
					NonceManager:       state.CCIPAddress,
				},
				DynamicConfig: view.OffRampDynamicConfig{
					FeeQuoter:                               state.CCIPAddress,
					PermissionlessExecutionThresholdSeconds: 28800,
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
		},
		Router: view.RouterView{
			ContractMetaData: view.ContractMetaData{
				Address:        state.CCIPRouterAddress,
				Owner:          owner,
				TypeAndVersion: "Router 1.6.0",
				StateObjectID:  state.CCIPRouterStateObjectID,
			},
			IsTestRouter: false,
			OnRamps:      map[uint64]string{EVMChainSelector: state.OnRampAddress},
			OffRamps:     nil,
		},
		TokenPools: map[string]map[string]view.TokenPoolView{
			"LINK": {
				state.BnMTokenPools["LINK"].PackageID: {
					ContractMetaData: view.ContractMetaData{
						Address:        state.BnMTokenPools["LINK"].PackageID,
						Owner:          owner,
						TypeAndVersion: "BurnMintTokenPool 1.6.0",
						StateObjectID:  state.BnMTokenPools["LINK"].StateObjectId,
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
			changesets.CCIPBnMSymbol: {
				state.ManagedTokenPools[changesets.CCIPBnMSymbol].PackageID: {
					ContractMetaData: view.ContractMetaData{
						Address:        state.ManagedTokenPools[changesets.CCIPBnMSymbol].PackageID,
						Owner:          owner,
						TypeAndVersion: "ManagedTokenPool 1.6.0",
						StateObjectID:  state.ManagedTokenPools[changesets.CCIPBnMSymbol].StateObjectId,
					},
					Token: state.ManagedTokens[changesets.CCIPBnMSymbol].TokenCoinMetadataID,
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
}

func (s *DeployTestSuite) GetDeployTPAndConfigureConfig() changesets.DeployTPAndConfigureConfig {
	coinTypeArg := fmt.Sprintf("%s::link::LINK", s.linkTokenPackageID)

	return changesets.DeployTPAndConfigureConfig{
		SuiChainSelector: SuiChainSelector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeBurnMint},
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
	}
}
