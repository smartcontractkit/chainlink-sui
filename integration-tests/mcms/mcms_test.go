//go:build integration

package mcms

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type CCIPMCMSTestSuite struct {
	MCMSTestSuite
}

func (s *CCIPMCMSTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()
}

func (s *CCIPMCMSTestSuite) Test_CCIP_MCMS() {
	s.T().Run("Transfer Ownership of CCIP to MCMS", func(t *testing.T) {
		s.RunOwnershipCCIPTransfer()
	})

	s.T().Run("Execute config proposal against CCIP from MCMS", func(t *testing.T) {
		RunTestCCIPFeeQuoterProposal(s)
		RunCCIPRampsProposal(s)
		RunTestRouterProposal(s)
	})
}

// TODO: For prod env, the initial deployment sequence should start the ownership transfer flow of every deployed contract
func RunTestCCIPFeeQuoterProposal(s *CCIPMCMSTestSuite) {
	// 1. Build configs
	expectedTTFC := module_fee_quoter.TokenTransferFeeConfig{
		MinFeeUsdCents:    3007,
		MaxFeeUsdCents:    30007,
		DeciBps:           1007,
		DestGasOverhead:   1000007,
		DestBytesOverhead: 1007,
		IsEnabled:         true,
	}

	expectedDestChainConfig := module_fee_quoter.DestChainConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           2,
		MaxDataBytes:                      2007,
		MaxPerMsgGasLimit:                 5000007,
		DestGasOverhead:                   1000007,
		DestGasPerPayloadByteBase:         byte(7),
		DestGasPerPayloadByteHigh:         byte(7),
		DestGasPerPayloadByteThreshold:    uint16(17),
		DestDataAvailabilityOverheadGas:   300007,
		DestGasPerDataAvailabilityByte:    7,
		DestDataAvailabilityMultiplierBps: 7,
		ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c},
		EnforceOutOfOrder:                 false,
		DefaultTokenFeeUsdCents:           7,
		DefaultTokenDestGasOverhead:       100007,
		DefaultTxGasLimit:                 500007,
		GasMultiplierWeiPerEth:            107,
		GasPriceStalenessThreshold:        307,
		NetworkFeeUsdCents:                7,
	}

	expectedPremiumMultiplier := uint64(77)
	destChainSelector := uint64(16015286601757825753)

	// 2. Run ops to generate proposal
	input := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			ccipops.FeeQuoterApplyFeeTokenUpdatesOp.Def(),
			ccipops.FeeQuoterApplyTokenTransferFeeConfigUpdatesOp.Def(),
			ccipops.FeeQuoterApplyDestChainConfigUpdatesOp.Def(),
			ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp.Def(),
		},
		Inputs: []any{
			ccipops.FeeQuoterApplyFeeTokenUpdatesInput{
				CCIPPackageId:     s.ccipPackageId,
				StateObjectId:     s.ccipObjects.CCIPObjectRefObjectId,
				OwnerCapObjectId:  s.ccipObjects.OwnerCapObjectId,
				FeeTokensToRemove: []string{},
				FeeTokensToAdd:    []string{s.linkObjects.CoinMetadataObjectId},
			},
			ccipops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
				CCIPPackageId:        s.ccipPackageId,
				StateObjectId:        s.ccipObjects.CCIPObjectRefObjectId,
				OwnerCapObjectId:     s.ccipObjects.OwnerCapObjectId,
				DestChainSelector:    destChainSelector,
				AddTokens:            []string{s.linkObjects.CoinMetadataObjectId},
				AddMinFeeUsdCents:    []uint32{expectedTTFC.MinFeeUsdCents},
				AddMaxFeeUsdCents:    []uint32{expectedTTFC.MaxFeeUsdCents},
				AddDeciBps:           []uint16{expectedTTFC.DeciBps},
				AddDestGasOverhead:   []uint32{expectedTTFC.DestGasOverhead},
				AddDestBytesOverhead: []uint32{expectedTTFC.DestBytesOverhead},
				AddIsEnabled:         []bool{expectedTTFC.IsEnabled},
				RemoveTokens:         []string{},
			},
			ccipops.FeeQuoterApplyDestChainConfigUpdatesInput{
				CCIPPackageId:                     s.ccipPackageId,
				StateObjectId:                     s.ccipObjects.CCIPObjectRefObjectId,
				OwnerCapObjectId:                  s.ccipObjects.OwnerCapObjectId,
				DestChainSelector:                 destChainSelector,
				IsEnabled:                         expectedDestChainConfig.IsEnabled,
				MaxNumberOfTokensPerMsg:           expectedDestChainConfig.MaxNumberOfTokensPerMsg,
				MaxDataBytes:                      expectedDestChainConfig.MaxDataBytes,
				MaxPerMsgGasLimit:                 expectedDestChainConfig.MaxPerMsgGasLimit,
				DestGasOverhead:                   expectedDestChainConfig.DestGasOverhead,
				DestGasPerPayloadByteBase:         expectedDestChainConfig.DestGasPerPayloadByteBase,
				DestGasPerPayloadByteHigh:         expectedDestChainConfig.DestGasPerPayloadByteHigh,
				DestGasPerPayloadByteThreshold:    expectedDestChainConfig.DestGasPerPayloadByteThreshold,
				DestDataAvailabilityOverheadGas:   expectedDestChainConfig.DestDataAvailabilityOverheadGas,
				DestGasPerDataAvailabilityByte:    expectedDestChainConfig.DestGasPerDataAvailabilityByte,
				DestDataAvailabilityMultiplierBps: expectedDestChainConfig.DestDataAvailabilityMultiplierBps,
				ChainFamilySelector:               expectedDestChainConfig.ChainFamilySelector,
				EnforceOutOfOrder:                 expectedDestChainConfig.EnforceOutOfOrder,
				DefaultTokenFeeUsdCents:           expectedDestChainConfig.DefaultTokenFeeUsdCents,
				DefaultTokenDestGasOverhead:       expectedDestChainConfig.DefaultTokenDestGasOverhead,
				DefaultTxGasLimit:                 expectedDestChainConfig.DefaultTxGasLimit,
				GasMultiplierWeiPerEth:            expectedDestChainConfig.GasMultiplierWeiPerEth,
				GasPriceStalenessThreshold:        expectedDestChainConfig.GasPriceStalenessThreshold,
				NetworkFeeUsdCents:                expectedDestChainConfig.NetworkFeeUsdCents,
			},
			ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
				CCIPPackageId:              s.ccipPackageId,
				StateObjectId:              s.ccipObjects.CCIPObjectRefObjectId,
				OwnerCapObjectId:           s.ccipObjects.OwnerCapObjectId,
				Tokens:                     []string{s.linkObjects.CoinMetadataObjectId},
				PremiumMultiplierWeiPerEth: []uint64{expectedPremiumMultiplier},
			},
		},
		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		ChainSelector:      uint64(s.chainSelector),
		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionSchedule,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}
	feeQuoterReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, input)
	s.Require().NoError(err, "executing fee quoter proposal sequence")

	timelockProposal := feeQuoterReport.Output

	// 3. Execute proposal
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// 4. Verify the changes in CCIP state object
	fqContract, err := module_fee_quoter.NewFeeQuoter(s.ccipPackageId, s.client)
	require.NoError(s.T(), err)

	ccipObjRef := bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}
	linkTokenID := s.linkObjects.CoinMetadataObjectId

	// Verify fee tokens
	feeTokens, err := fqContract.DevInspect().GetFeeTokens(s.T().Context(), s.deps.GetCallOpts(), ccipObjRef)
	require.NoError(s.T(), err)
	require.Contains(s.T(), feeTokens, linkTokenID)

	// Verify token transfer fee config matches input
	actualTTFC, err := fqContract.DevInspect().GetTokenTransferFeeConfig(s.T().Context(), s.deps.GetCallOpts(), ccipObjRef, destChainSelector, linkTokenID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedTTFC, actualTTFC)

	// Verify destination chain config matches input
	actualDestChainConfig, err := fqContract.DevInspect().GetDestChainConfig(s.T().Context(), s.deps.GetCallOpts(), ccipObjRef, destChainSelector)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedDestChainConfig, actualDestChainConfig)

	// Verify premium multiplier matches input
	actualPremiumMultiplier, err := fqContract.DevInspect().GetPremiumMultiplierWeiPerEth(s.T().Context(), s.deps.GetCallOpts(), ccipObjRef, linkTokenID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedPremiumMultiplier, actualPremiumMultiplier)
}

func RunCCIPRampsProposal(s *CCIPMCMSTestSuite) {
	// 1. Build configs
	mock32Bytes := []byte{
		0x33, 0x17, 0xaa, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
	}
	configDigest := []byte{
		0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
		0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
		0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
		0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
	}
	expectedSCC := module_offramp.SourceChainConfig{
		IsEnabled:                 false,
		IsRmnVerificationDisabled: true,
		OnRamp:                    mock32Bytes,
	}
	expectedDCC := module_onramp.DestChainConfig{
		AllowlistEnabled: true,
		Router:           "0x304121906bf93b21f915a04cffea4df21090432e3c2fd60e51ebe68f79c90a41",
		AllowedSenders: []string{
			"0x1cf00ee891001df44fc0736e56f469ab85dcf9b78511ac9268f292716fc04447",
			"0x2d011ff9a2112e0550d1847f67057abc96ed09b78511ac9268f292716fc04447",
		},
	}

	// 2. Run ops to generate proposal
	input := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			offrampops.ApplySourceChainConfigUpdatesOp.Def(),
			offrampops.SetOCR3ConfigOp.Def(),
			onrampops.ApplyDestChainConfigUpdateOp.Def(),
			onrampops.ApplyAllowListUpdateOp.Def(),
		},
		Inputs: []any{
			offrampops.ApplySourceChainConfigUpdateInput{
				CCIPObjectRef:                         s.ccipObjects.CCIPObjectRefObjectId,
				OffRampPackageId:                      s.ccipOfframpPackageId,
				OffRampStateId:                        s.ccipOfframpObjects.StateObjectId,
				OwnerCapObjectId:                      s.ccipOfframpObjects.OwnerCapId,
				SourceChainsSelectors:                 []uint64{cselectors.ETHEREUM_MAINNET.Selector},
				SourceChainsIsEnabled:                 []bool{expectedSCC.IsEnabled},
				SourceChainsIsRMNVerificationDisabled: []bool{expectedSCC.IsRmnVerificationDisabled},
				SourceChainsOnRamp:                    [][]byte{expectedSCC.OnRamp},
			},
			offrampops.SetOCR3ConfigInput{
				OffRampPackageId:               s.ccipOfframpPackageId,
				OffRampStateId:                 s.ccipOfframpObjects.StateObjectId,
				CCIPObjectRefId:                s.ccipObjects.CCIPObjectRefObjectId,
				OwnerCapObjectId:               s.ccipOfframpObjects.OwnerCapId,
				ConfigDigest:                   configDigest,
				OCRPluginType:                  byte(1),
				BigF:                           byte(1),
				IsSignatureVerificationEnabled: false,
				Signers:                        [][]byte{mock32Bytes},
				Transmitters:                   []string{"0x11223344556677889900aabbccddeeff00112233"},
			},
			onrampops.ApplyDestChainConfigureOnRampInput{
				OnRampPackageId:           s.ccipOnrampPackageId,
				CCIPObjectRefId:           s.ccipObjects.CCIPObjectRefObjectId,
				OwnerCapObjectId:          s.ccipOnrampObjects.OwnerCapObjectId,
				StateObjectId:             s.ccipOnrampObjects.StateObjectId,
				DestChainSelector:         []uint64{cselectors.ETHEREUM_MAINNET.Selector},
				DestChainAllowListEnabled: []bool{false},
				DestChainRouters:          []string{expectedDCC.Router},
			},
			onrampops.ApplyAllowListUpdatesInput{
				OnRampPackageId:               s.ccipOnrampPackageId,
				CCIPObjectRefId:               s.ccipObjects.CCIPObjectRefObjectId,
				OwnerCapObjectId:              s.ccipOnrampObjects.OwnerCapObjectId,
				StateObjectId:                 s.ccipOnrampObjects.StateObjectId,
				DestChainSelector:             []uint64{cselectors.ETHEREUM_MAINNET.Selector},
				DestChainAllowListEnabled:     []bool{expectedDCC.AllowlistEnabled},
				DestChainAddAllowedSenders:    [][]string{expectedDCC.AllowedSenders},
				DestChainRemoveAllowedSenders: [][]string{{}},
			},
		},
		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		ChainSelector:      uint64(s.chainSelector),
		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}
	rampsReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, input)
	s.Require().NoError(err, "executing ramps proposal sequence")

	timelockProposal := rampsReport.Output

	// 3. Execute proposal
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// 4. Assert changes in contracts

	// Create contract instances
	offrampContract, err := module_offramp.NewOfframp(s.ccipOfframpPackageId, s.client)
	require.NoError(s.T(), err, "creating offramp contract")

	onrampContract, err := module_onramp.NewOnramp(s.ccipOnrampPackageId, s.client)
	require.NoError(s.T(), err, "creating onramp contract")

	ccipObjRef := bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}
	offRampStateObj := bind.Object{Id: s.ccipOfframpObjects.StateObjectId}
	onRampStateObj := bind.Object{Id: s.ccipOnrampObjects.StateObjectId}

	// Verify SourceChainConfig changes in OffRamp
	actualSCC, err := offrampContract.DevInspect().GetSourceChainConfig(s.T().Context(), s.deps.GetCallOpts(), ccipObjRef, offRampStateObj, cselectors.ETHEREUM_MAINNET.Selector)
	require.NoError(s.T(), err, "getting source chain config")
	require.Equal(s.T(), expectedSCC.IsEnabled, actualSCC.IsEnabled, "source chain config IsEnabled should match")
	require.Equal(s.T(), expectedSCC.IsRmnVerificationDisabled, actualSCC.IsRmnVerificationDisabled, "source chain config IsRmnVerificationDisabled should match")
	require.Equal(s.T(), expectedSCC.OnRamp, actualSCC.OnRamp, "source chain config OnRamp should match")

	// Verify DestChainConfig changes in OnRamp
	actualDCCResults, err := onrampContract.DevInspect().GetDestChainConfig(s.T().Context(), s.deps.GetCallOpts(), onRampStateObj, cselectors.ETHEREUM_MAINNET.Selector)
	require.NoError(s.T(), err, "getting dest chain config")
	// GetDestChainConfig returns multiple values: [0]: u64, [1]: bool, [2]: address
	// Based on the bindings, it seems to return sequence number, allowlist enabled, and router
	require.Len(s.T(), actualDCCResults, 3, "dest chain config should return 3 values")
	actualAllowlistEnabled, ok := actualDCCResults[1].(bool)
	require.True(s.T(), ok, "second value should be bool for allowlist enabled")
	actualRouter, ok := actualDCCResults[2].(string)
	require.True(s.T(), ok, "third value should be string for router")

	// Note: The allowlist enabled in the ApplyDestChainConfigureOnRampInput was set to false,
	// but in ApplyAllowListUpdatesInput it was set to true. The final state should be true.
	require.Equal(s.T(), expectedDCC.AllowlistEnabled, actualAllowlistEnabled, "dest chain config AllowlistEnabled should match")
	require.Equal(s.T(), expectedDCC.Router, actualRouter, "dest chain config Router should match")

	// Verify AllowList changes in OnRamp
	actualAllowListResults, err := onrampContract.DevInspect().GetAllowedSendersList(s.T().Context(), s.deps.GetCallOpts(), onRampStateObj, cselectors.ETHEREUM_MAINNET.Selector)
	require.NoError(s.T(), err, "getting allowed senders list")
	// GetAllowedSendersList returns [0]: bool, [1]: vector<address>
	require.Len(s.T(), actualAllowListResults, 2, "allowed senders list should return 2 values")
	actualAllowlistEnabledFromList, ok := actualAllowListResults[0].(bool)
	require.True(s.T(), ok, "first value should be bool for allowlist enabled")
	actualAllowedSenders, ok := actualAllowListResults[1].([]string)
	require.True(s.T(), ok, "second value should be []string for allowed senders")

	require.Equal(s.T(), expectedDCC.AllowlistEnabled, actualAllowlistEnabledFromList, "allowlist enabled should match")
	for _, expectedSender := range expectedDCC.AllowedSenders {
		require.Contains(s.T(), actualAllowedSenders, expectedSender, "allowed senders should contain expected sender")
	}
}

func RunTestRouterProposal(s *CCIPMCMSTestSuite) {
	// 1. Build configs
	expectedDestChainSelectors := []uint64{
		cselectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
		cselectors.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.Selector,
	}
	expectedOnRampAddresses := []string{
		"0x1111111111111111111111111111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222222222222222222222222222",
	}

	// 2. Run ops to generate proposal
	input := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			routerops.SetOnRampsOp.Def(),
		},
		Inputs: []any{
			routerops.SetOnRampsInput{
				RouterPackageId:     s.ccipRouterPackageId,
				RouterStateObjectId: s.ccipRouterObjects.RouterStateObjectId,
				OwnerCapObjectId:    s.ccipRouterObjects.OwnerCapObjectId,
				DestChainSelectors:  expectedDestChainSelectors,
				OnRampAddresses:     expectedOnRampAddresses,
			},
		},
		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		ChainSelector:      uint64(s.chainSelector),
		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}
	routerReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, input)
	s.Require().NoError(err, "executing router proposal sequence")

	timelockProposal := routerReport.Output

	// 3. Execute proposal
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// 4. Assert changes in contracts

	// Create router contract instance
	routerContract, err := module_router.NewRouter(s.ccipRouterPackageId, s.client)
	require.NoError(s.T(), err, "creating router contract")

	routerStateObj := bind.Object{Id: s.ccipRouterObjects.RouterStateObjectId}

	// Verify OnRamp addresses are set correctly
	for i, destChainSelector := range expectedDestChainSelectors {
		// Verify chain is supported
		isSupported, err := routerContract.DevInspect().IsChainSupported(s.T().Context(), s.deps.GetCallOpts(), routerStateObj, destChainSelector)
		require.NoError(s.T(), err, "checking if chain is supported")
		require.True(s.T(), isSupported, "chain %d should be supported", destChainSelector)

		// Verify OnRamp address matches expected
		actualOnRampAddress, err := routerContract.DevInspect().GetOnRamp(s.T().Context(), s.deps.GetCallOpts(), routerStateObj, destChainSelector)
		require.NoError(s.T(), err, "getting on-ramp address for chain %d", destChainSelector)
		require.Equal(s.T(), expectedOnRampAddresses[i], actualOnRampAddress, "on-ramp address for chain %d should match", destChainSelector)
	}
}
