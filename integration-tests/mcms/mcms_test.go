//go:build integration

package mcms

import (
	"encoding/hex"
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
)

type CCIPMCMSTestSuite struct {
	MCMSTestSuite

	// CCIP
	ccipPackageId string
	ccipObjects   ccipops.DeployCCIPSeqObjects
	linkObjects   linkops.DeployLinkObjects

	// Router
	ccipRouterPackageId string
	ccipRouterObjects   routerops.DeployCCIPRouterObjects

	// Onramp
	ccipOnrampPackageId string
	ccipOnrampObjects   onrampops.DeployCCIPOnRampSeqObjects

	// offramp
	ccipOfframpPackageId string
	ccipOfframpObjects   offrampops.DeployCCIPOffRampSeqObjects
}

func (s *CCIPMCMSTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()

	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(s.bundle, linkops.DeployLINKOp, s.deps, cld_ops.EmptyInput{})
	require.NoError(s.T(), err, "failed to deploy LINK token")

	configDigestHex := "e3b1c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	configDigest, err := hex.DecodeString(configDigestHex)
	require.NoError(s.T(), err, "failed to decode config digest")

	publicKey1Hex := "8a1b2c3d4e5f60718293a4b5c6d7e8f901234567"
	publicKey1, err := hex.DecodeString(publicKey1Hex)
	require.NoError(s.T(), err, "failed to decode public key 1")

	publicKey2Hex := "7b8c9dab0c1d2e3f405162738495a6b7c8d9e0f1"
	publicKey2, err := hex.DecodeString(publicKey2Hex)
	require.NoError(s.T(), err, "failed to decode public key 2")

	publicKey3Hex := "1234567890abcdef1234567890abcdef12345678"
	publicKey3, err := hex.DecodeString(publicKey3Hex)
	require.NoError(s.T(), err, "failed to decode public key 3")

	publicKey4Hex := "90abcdef1234567890abcdef1234567890abcdef"
	publicKey4, err := hex.DecodeString(publicKey4Hex)
	require.NoError(s.T(), err, "failed to decode public key 4")

	// Use the same seq as in production deployment
	ccipReport, err := cld_ops.ExecuteSequence(s.bundle, ccipops.DeployAndInitCCIPSequence, s.deps, ccipops.DeployAndInitCCIPSeqInput{
		LinkTokenCoinMetadataObjectId: linkReport.Output.Objects.CoinMetadataObjectId,
		LocalChainSelector:            1,
		DestChainSelector:             2,
		DeployCCIPInput: ccipops.DeployCCIPInput{
			McmsPackageId: s.mcmsPackageID,
			McmsOwner:     s.mcmsOwnerAddress,
		},
		MaxFeeJuelsPerMsg:            "100000000",
		TokenPriceStalenessThreshold: 60,
		// Fee Quoter configuration
		AddMinFeeUsdCents:    []uint32{3000},
		AddMaxFeeUsdCents:    []uint32{30000},
		AddDeciBps:           []uint16{1000},
		AddDestGasOverhead:   []uint32{1000000},
		AddDestBytesOverhead: []uint32{1000},
		AddIsEnabled:         []bool{true},
		RemoveTokens:         []string{},
		// Fee Quoter destination chain configuration
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           2,
		MaxDataBytes:                      2000,
		MaxPerMsgGasLimit:                 5000000,
		DestGasOverhead:                   1000000,
		DestGasPerPayloadByteBase:         byte(2),
		DestGasPerPayloadByteHigh:         byte(5),
		DestGasPerPayloadByteThreshold:    uint16(10),
		DestDataAvailabilityOverheadGas:   300000,
		DestGasPerDataAvailabilityByte:    4,
		DestDataAvailabilityMultiplierBps: 1,
		ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c},
		EnforceOutOfOrder:                 false,
		DefaultTokenFeeUsdCents:           3,
		DefaultTokenDestGasOverhead:       100000,
		DefaultTxGasLimit:                 500000,
		GasMultiplierWeiPerEth:            100,
		GasPriceStalenessThreshold:        1000000000,
		NetworkFeeUsdCents:                10,
		// Premium multiplier updates
		PremiumMultiplierWeiPerEth: []uint64{10},

		RmnHomeContractConfigDigest: configDigest,
		SignerOnchainPublicKeys:     [][]byte{publicKey1, publicKey2, publicKey3, publicKey4},
		NodeIndexes:                 []uint64{0, 1, 2, 3},
		FSign:                       uint64(1),
	})
	require.NoError(s.T(), err, "failed to execute CCIP deploy sequence")
	require.NotEmpty(s.T(), ccipReport.Output.CCIPPackageId, "CCIP package ID should not be empty")

	s.linkObjects = linkReport.Output.Objects
	s.ccipPackageId = ccipReport.Output.CCIPPackageId
	s.ccipObjects = ccipReport.Output.Objects

	// Deploy Router
	routerReport, err := cld_ops.ExecuteOperation(s.bundle, routerops.DeployCCIPRouterOp, s.deps, routerops.DeployCCIPRouterInput{
		McmsPackageId: s.mcmsPackageID,
		McmsOwner:     s.mcmsOwnerAddress,
	})
	require.NoError(s.T(), err, "failed to execute CCIP deploy sequence")

	s.ccipRouterPackageId = routerReport.Output.PackageId
	s.ccipRouterObjects = routerReport.Output.Objects

	// Deploy Onramp
	ccipOnRampSeqInput := deployment.DefaultOnRampSeqConfig
	ccipOnRampSeqInput.DeployCCIPOnRampInput.CCIPPackageId = ccipReport.Output.CCIPPackageId
	ccipOnRampSeqInput.DeployCCIPOnRampInput.MCMSPackageId = s.mcmsPackageID
	ccipOnRampSeqInput.DeployCCIPOnRampInput.MCMSOwnerPackageId = s.mcmsOwnerAddress
	ccipOnRampSeqInput.OnRampInitializeInput.NonceManagerCapId = ccipReport.Output.Objects.NonceManagerCapObjectId
	ccipOnRampSeqInput.OnRampInitializeInput.SourceTransferCapId = ccipReport.Output.Objects.SourceTransferCapObjectId
	ccipOnRampSeqInput.OnRampInitializeInput.ChainSelector = uint64(s.chainSelector)
	ccipOnRampSeqInput.OnRampInitializeInput.FeeAggregator = s.mcmsOwnerAddress
	ccipOnRampSeqInput.OnRampInitializeInput.AllowListAdmin = s.mcmsOwnerAddress
	ccipOnRampSeqInput.OnRampInitializeInput.DestChainSelectors = []uint64{cselectors.ETHEREUM_MAINNET.Selector}
	ccipOnRampSeqInput.OnRampInitializeInput.DestChainRouters = []string{routerReport.Output.PackageId}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.DestChainSelector = []uint64{cselectors.ETHEREUM_MAINNET.Selector}
	ccipOnRampSeqInput.ApplyAllowListUpdatesInput.DestChainSelector = []uint64{cselectors.ETHEREUM_MAINNET.Selector}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.DestChainRouters = []string{routerReport.Output.PackageId}
	ccipOnRampSeqInput.ApplyDestChainConfigureOnRampInput.CCIPObjectRefId = ccipReport.Output.Objects.CCIPObjectRefObjectId

	ccipOnRampSeqReport, err := operations.ExecuteSequence(s.bundle, onrampops.DeployAndInitCCIPOnRampSequence, s.deps, ccipOnRampSeqInput)
	require.NoError(s.T(), err, "failed to execute CCIP OnRamp deploy sequence")

	s.ccipOnrampPackageId = ccipOnRampSeqReport.Output.CCIPOnRampPackageId
	s.ccipOnrampObjects = ccipOnRampSeqReport.Output.Objects

	// Deploy offramp
	ccipOffRampSeqInput := deployment.DefaultOffRampSeqConfig
	// note: this is a regression, can't acess other chains state very cleanly
	onRampBytes := [][]byte{
		{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
	}

	// Inject dynamic values for deployment
	ccipOffRampSeqInput.CCIPObjectRefId = ccipReport.Output.Objects.CCIPObjectRefObjectId
	ccipOffRampSeqInput.DeployCCIPOffRampInput.CCIPPackageId = ccipReport.Output.CCIPPackageId
	ccipOffRampSeqInput.DeployCCIPOffRampInput.MCMSPackageId = s.mcmsPackageID

	ccipOffRampSeqInput.InitializeOffRampInput.DestTransferCapId = ccipReport.Output.Objects.DestTransferCapObjectId
	ccipOffRampSeqInput.InitializeOffRampInput.FeeQuoterCapId = ccipReport.Output.Objects.FeeQuoterCapObjectId
	ccipOffRampSeqInput.InitializeOffRampInput.ChainSelector = uint64(s.chainSelector)
	ccipOffRampSeqInput.InitializeOffRampInput.SourceChainSelectors = []uint64{
		cselectors.ETHEREUM_MAINNET.Selector,
	}
	ccipOffRampSeqInput.InitializeOffRampInput.SourceChainsOnRamp = onRampBytes

	ccipOffRampSeqReport, err := operations.ExecuteSequence(s.bundle, offrampops.DeployAndInitCCIPOffRampSequence, s.deps, ccipOffRampSeqInput)
	require.NoError(s.T(), err, "failed to execute CCIP OffRamp deploy sequence")

	s.ccipOfframpPackageId = ccipOffRampSeqReport.Output.CCIPOffRampPackageId
	s.ccipOfframpObjects = ccipOffRampSeqReport.Output.Objects
}

func (s *CCIPMCMSTestSuite) Test_CCIP_MCMS() {
	s.T().Run("Transfer Ownership of CCIP to MCMS", func(t *testing.T) {
		RunTestCCIPOwnershipTransfer(s)
	})

	s.T().Run("Execute config proposal against CCIP from MCMS", func(t *testing.T) {
		RunTestCCIPFeeQuoterProposal(s)
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

	// 4. Verify the new owner is MCMS
	newOwner, err := ccipContract.DevInspect().Owner(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId})
	s.Require().NoError(err, "getting new owner of CCIP state object")
	s.Require().Equal(s.mcmsPackageID, newOwner, "new owner of CCIP should be MCMS")
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
