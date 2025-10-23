//go:build integration

package mcms

import (
	"testing"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
)

type CCIPMCMSTestSuite struct {
	MCMSTestSuite
}

func (s *CCIPMCMSTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()
}

func (s *CCIPMCMSTestSuite) Test_CCIP_MCMS() {
	s.T().Run("Transfer Ownership of CCIP to MCMS", func(t *testing.T) {
		RunTestCCIPOwnershipTransfer(s)
	})

	s.T().Run("Execute config proposal against CCIP from MCMS", func(t *testing.T) {
		RunTestCCIPFeeQuoterProposal(s)
		RunCCIPOffRampProposal(s)
	})
}

// TODO: For prod env, the initial deployment sequence should start the ownership transfer flow of every deployed contract
func RunTestCCIPOwnershipTransfer(s *CCIPMCMSTestSuite) {
	// 1a. Transfer OwnerCap of CCIP to MCMS (this should be done in the initial deployment sequence)
	ccipContract, err := module_state_object.NewStateObject(s.ccipPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip state object contract")

	tx, err := ccipContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId},
		bind.Object{Id: s.ccipObjects.OwnerCapObjectId},
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	s.T().Logf("✅ Transferred ownership of CCIP to MCMS in tx: %s", tx.Digest)

	// 1b. Transfer ownership of the CCIP Router to MCMS
	ccipRouterContract, err := module_router.NewRouter(s.ccipRouterPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip router contract")

	tx, err = ccipRouterContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipRouterObjects.RouterStateObjectId},
		bind.Object{Id: s.ccipRouterObjects.OwnerCapObjectId},
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP Router to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	s.T().Logf("✅ Transferred ownership of CCIP Router to MCMS in tx: %s", tx.Digest)

	// 1c. Transfer ownership of the CCIP OnRamp to MCMS
	ccipOnRampContract, err := module_onramp.NewOnramp(s.ccipOnrampPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip onramp contract")

	tx, err = ccipOnRampContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId},
		bind.Object{Id: s.ccipOnrampObjects.StateObjectId},
		bind.Object{Id: s.ccipOnrampObjects.OwnerCapObjectId},
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP OnRamp to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	s.T().Logf("✅ Transferred ownership of CCIP OnRamp to MCMS in tx: %s", tx.Digest)

	// 1d. Transfer ownership of the CCIP OffRamp to MCMS
	ccipOffRampContract, err := module_offramp.NewOfframp(s.ccipOfframpPackageId, s.client)
	require.NoError(s.T(), err, "creating ccip offramp contract")

	tx, err = ccipOffRampContract.TransferOwnership(
		s.T().Context(),
		&bind.CallOpts{
			Signer:           s.signer,
			WaitForExecution: true,
		},
		bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId},
		bind.Object{Id: s.ccipOfframpObjects.StateObjectId},
		bind.Object{Id: s.ccipOfframpObjects.OwnerCapId},
		s.mcmsPackageID,
	)
	require.NoError(s.T(), err, "transferring ownership of CCIP OffRamp to MCMS")
	require.NotEmpty(s.T(), tx, "Transaction should not be empty")

	// 2. Proposal execution with acceptance from MCMS (through bypasser)
	input := ownershipops.AcceptCCIPOwnershipInput{
		// MCMS related
		MCMSPackageId:     s.mcmsPackageID,
		MCMSStateObjId:    s.mcmsObj,
		MCMSTimelockObjId: s.timelockObj,
		MCMSAccountObjId:  s.accountObj,
		MCMSRegistryObjId: s.registryObj,

		CCIPPackageId: s.ccipPackageId,
		CCIPObjectRef: s.ccipObjects.CCIPObjectRefObjectId,

		RouterPackageId:     s.ccipRouterPackageId,
		RouterStateObjectId: s.ccipRouterObjects.RouterStateObjectId,

		OnRampPackageId:     s.ccipOnrampPackageId,
		OnRampStateObjectId: s.ccipOnrampObjects.StateObjectId,

		OffRampPackageId:     s.ccipOfframpPackageId,
		OffRampStateObjectId: s.ccipOfframpObjects.StateObjectId,

		// Proposal
		Role: suisdk.TimelockRoleBypasser,

		ChainSelector: uint64(s.chainSelector),
	}
	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(s.bundle, ownershipops.AcceptCCIPOwnershipSeq, s.deps, input)
	s.Require().NoError(err, "executing ownership acceptance proposal sequence")

	timelockProposal := acceptOwnershipProposalReport.Output

	// 3. Execute transfer ownership from original owner
	// 3.1. Execute the proposal
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// 3.2. Finish the ownership transfer with the original owner signer
	_, err = ccipContract.ExecuteOwnershipTransferToMcms(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}, bind.Object{Id: s.ccipObjects.OwnerCapObjectId}, bind.Object{Id: s.registryObj}, s.mcmsPackageID)
	s.Require().NoError(err, "executing ownership transfer of CCIP to MCMS")
	_, err = ccipOnRampContract.ExecuteOwnershipTransferToMcms(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}, bind.Object{Id: s.ccipOnrampObjects.OwnerCapObjectId}, bind.Object{Id: s.ccipOnrampObjects.StateObjectId}, bind.Object{Id: s.registryObj}, s.mcmsPackageID)
	s.Require().NoError(err, "executing ownership transfer of OnRamp to MCMS")
	_, err = ccipOffRampContract.ExecuteOwnershipTransferToMcms(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}, bind.Object{Id: s.ccipOfframpObjects.OwnerCapId}, bind.Object{Id: s.ccipOfframpObjects.StateObjectId}, bind.Object{Id: s.registryObj}, s.mcmsPackageID)
	s.Require().NoError(err, "executing ownership transfer of OffRamp to MCMS")

	// 4. Verify the new owner is MCMS
	newOwner, err := ccipContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId})
	s.Require().NoError(err, "getting new owner of CCIP state object")
	s.Require().Equal(s.mcmsPackageID, newOwner, "new owner of CCIP should be MCMS")

	newOwnerOnRamp, err := ccipOnRampContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipOnrampObjects.StateObjectId})
	s.Require().NoError(err, "getting new owner of OnRamp state object")
	s.Require().Equal(s.mcmsPackageID, newOwnerOnRamp, "new owner of OnRamp should be MCMS")

	newOwnerOffRamp, err := ccipOffRampContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipOfframpObjects.StateObjectId})
	s.Require().NoError(err, "getting new owner of OffRamp state object")
	s.Require().Equal(s.mcmsPackageID, newOwnerOffRamp, "new owner of OffRamp should be MCMS")
}

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
		MmcsPackageID:  s.mcmsPackageID,
		McmsStateObjID: s.mcmsObj,
		TimelockObjID:  s.timelockObj,
		AccountObjID:   s.accountObj,
		RegistryObjID:  s.registryObj,
		// Proposal
		Role:          suisdk.TimelockRoleBypasser,
		ChainSelector: uint64(s.chainSelector),
		Delay:         0,
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

func RunCCIPOffRampProposal(s *CCIPMCMSTestSuite) {
	// 1. Build configs
	// onRampBytes := []byte{
	// 	0x33, 0x17, 0xaa, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
	// 	0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
	// }
	// expectedSCC := module_offramp.SourceChainConfig{
	// 	IsEnabled:                 false,
	// 	IsRmnVerificationDisabled: true,
	// 	OnRamp:                    onRampBytes,
	// }
	expectedDCC := module_onramp.DestChainConfig{
		AllowlistEnabled: true,
		Router:           "0x304121906bf93b21f915a04cffea4df21090432e3c2fd60e51ebe68f79c90a41",
		AllowedSenders: []string{
			"0x1cf00ee891001df44fc0736e56f469ab85dcf9b78511ac9268f292716fc04447",
			"0x1cf00ee891001df44fc0736e56f469ab85dcf9b78511ac9268f292716fc04447",
		},
	}

	// 2. Run ops to generate proposal
	input := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			// offrampops.ApplySourceChainConfigUpdatesOp.Def(),
			// offrampops.SetOCR3ConfigOp.Def(),
			onrampops.ApplyDestChainConfigUpdateOp.Def(),
			onrampops.ApplyAllowListUpdateOp.Def(),
		},
		Inputs: []any{
			// offrampops.ApplySourceChainConfigUpdateInput{
			// 	CCIPObjectRef:                         s.ccipObjects.CCIPObjectRefObjectId,
			// 	OffRampPackageId:                      s.ccipOfframpPackageId,
			// 	OffRampStateId:                        s.ccipOfframpObjects.StateObjectId,
			// 	OwnerCapObjectId:                      s.ccipOfframpObjects.OwnerCapId,
			// 	SourceChainsSelectors:                 []uint64{cselectors.ETHEREUM_MAINNET.Selector},
			// 	SourceChainsIsEnabled:                 []bool{expectedSCC.IsEnabled},
			// 	SourceChainsIsRMNVerificationDisabled: []bool{expectedSCC.IsRmnVerificationDisabled},
			// 	SourceChainsOnRamp:                    [][]byte{expectedSCC.OnRamp},
			// },
			// offrampops.SetOCR3ConfigInput{
			// 	OffRampPackageId:               s.ccipOfframpPackageId,
			// 	OffRampStateId:                 s.ccipOfframpObjects.StateObjectId,
			// 	CCIPObjectRefId:                s.ccipObjects.CCIPObjectRefObjectId,
			// 	OwnerCapObjectId:               s.ccipOfframpObjects.OwnerCapId,
			// 	ConfigDigest:                   []byte{0x01, 0x02, 0x03},
			// 	OCRPluginType:                  1,
			// 	BigF:                           0,
			// 	IsSignatureVerificationEnabled: false,
			// 	Signers:                        [][]byte{onRampBytes},
			// 	Transmitters:                   []string{"0x11223344556677889900aabbccddeeff00112233"},
			// },
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
		MmcsPackageID:  s.mcmsPackageID,
		McmsStateObjID: s.mcmsObj,
		TimelockObjID:  s.timelockObj,
		AccountObjID:   s.accountObj,
		RegistryObjID:  s.registryObj,
		// Proposal
		Role:          suisdk.TimelockRoleBypasser,
		ChainSelector: uint64(s.chainSelector),
		Delay:         0,
	}
	rampsReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, input)
	s.Require().NoError(err, "executing ramps proposal sequence")

	timelockProposal := rampsReport.Output

	// 3. Execute proposal
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// TODO: 4. Assert changes in contracts
}
