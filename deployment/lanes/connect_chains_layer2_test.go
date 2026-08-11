package lanes_test

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/lanes"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

const (
	testOnRampPackageID      = "0xf87c6010be571a304f0d860857204bc66f037842156f0f6c9d80be265fd83752"
	testOnRampStateObjectID  = "0x75ec1e10b4302f7c69476eb196c88a0aa43a4d509bbbe5cc1feb213e4b6dd58b"
	testOnRampOwnerCapID     = "0x8101795ff02d4935a05fb519e1b21b83855a970639cf0c28fa7a51f7d2e689ae"
	testRouterPackageID      = "0xed4613bd35004954c07150c3e9b10230f5e23e3058bc2ca0e3e676cb43eb4dc1"
	testRouterStateObjectID  = "0xbb2486d233b0d358f82fb8c4c5c75881e65069ea8ebe5ab692a636c9e0eff7cd"
	testRouterOwnerCapID     = "0xb000000000000000000000000000000000000000000000000000000000000001"
	testOffRampStateObjectID = "0x4fdacea0d627df26a6f34ac62952ed8d3c32ea70aacb90a7f2134b39e36a79cd"
	testOffRampOwnerCapID    = "0x8f7eb5b7879449519b39db447b94c87b6bb88f0a1deff40b3e0ffef3d1058f69"
)

func TestConfigureLaneLegAsSource_MCMSBatchOp_EncoderParity(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := testSuiChains()
	evmRouter := common.HexToAddress(evmRouterAddress).Bytes()
	input := sourceLegInput(evmRouter)
	destRouterHex := leftPaddedAddressHex(evmRouter)

	report, err := runSourceLeg(t, env, chains, input)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 5)

	feeQuoter, err := module_fee_quoter.NewFeeQuoter(testCCIPPackageID, nil)
	require.NoError(t, err)
	onRamp, err := module_onramp.NewOnramp(testOnRampPackageID, nil)
	require.NoError(t, err)
	router, err := module_router.NewRouter(testRouterPackageID, nil)
	require.NoError(t, err)

	t.Run("apply_token_transfer_fee_config_updates", func(t *testing.T) {
		encoded, err := feeQuoter.Encoder().ApplyTokenTransferFeeConfigUpdates(
			bind.Object{Id: testCCIPObjectRef},
			bind.Object{Id: testCCIPOwnerCapID},
			evmSepoliaSelector,
			[]string{testLinkCoinMetadata},
			[]uint32{lanes.DefaultLinkTokenTransferMinFeeUsdCents},
			[]uint32{lanes.DefaultLinkTokenTransferMaxFeeUsdCents},
			[]uint16{lanes.DefaultLinkTokenTransferDeciBps},
			[]uint32{lanes.DefaultLinkTokenTransferDestGasOverhead},
			[]uint32{lanes.DefaultLinkTokenTransferDestBytesOverhead},
			[]bool{true},
			[]string{},
		)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.BatchOps[0].Transactions[0].Data, encoded, testCCIPObjectRef, nil)
	})

	t.Run("apply_dest_chain_config_updates", func(t *testing.T) {
		cfg := input.Dest.FeeQuoterDestChainConfig
		translated := lanes.TranslateDestChainConfig(cfg, input.Dest.Selector)
		encoded, err := feeQuoter.Encoder().ApplyDestChainConfigUpdates(
			bind.Object{Id: testCCIPObjectRef},
			bind.Object{Id: testCCIPOwnerCapID},
			translated.DestChainSelector,
			translated.IsEnabled,
			translated.MaxNumberOfTokensPerMsg,
			translated.MaxDataBytes,
			translated.MaxPerMsgGasLimit,
			translated.DestGasOverhead,
			translated.DestGasPerPayloadByteBase,
			translated.DestGasPerPayloadByteHigh,
			translated.DestGasPerPayloadByteThreshold,
			translated.DestDataAvailabilityOverheadGas,
			translated.DestGasPerDataAvailabilityByte,
			translated.DestDataAvailabilityMultiplierBps,
			translated.ChainFamilySelector,
			translated.EnforceOutOfOrder,
			translated.DefaultTokenFeeUsdCents,
			translated.DefaultTokenDestGasOverhead,
			translated.DefaultTxGasLimit,
			translated.GasMultiplierWeiPerEth,
			translated.GasPriceStalenessThreshold,
			translated.NetworkFeeUsdCents,
		)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.BatchOps[1].Transactions[0].Data, encoded, testCCIPObjectRef, nil)
	})

	t.Run("apply_premium_multiplier_wei_per_eth_updates", func(t *testing.T) {
		encoded, err := feeQuoter.Encoder().ApplyPremiumMultiplierWeiPerEthUpdates(
			bind.Object{Id: testCCIPObjectRef},
			bind.Object{Id: testCCIPOwnerCapID},
			[]string{testLinkCoinMetadata},
			[]uint64{lanes.DefaultLinkPremiumMultiplierWeiPerEth},
		)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.BatchOps[2].Transactions[0].Data, encoded, testCCIPObjectRef, nil)
	})

	t.Run("onramp_apply_dest_chain_config_updates", func(t *testing.T) {
		encoded, err := onRamp.Encoder().ApplyDestChainConfigUpdates(
			bind.Object{Id: testCCIPObjectRef},
			bind.Object{Id: testOnRampStateObjectID},
			bind.Object{Id: testOnRampOwnerCapID},
			[]uint64{evmSepoliaSelector},
			[]bool{input.Source.AllowListEnabled},
			[]string{destRouterHex},
		)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.BatchOps[3].Transactions[0].Data, encoded, testOnRampStateObjectID, nil)
	})

	t.Run("router_set_on_ramps", func(t *testing.T) {
		encoded, err := router.Encoder().SetOnRamps(
			bind.Object{Id: testRouterOwnerCapID},
			bind.Object{Id: testRouterStateObjectID},
			[]uint64{evmSepoliaSelector},
			[]string{testOnRampPackageID},
		)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.BatchOps[4].Transactions[0].Data, encoded, testRouterStateObjectID, nil)
	})
}

func TestConfigureLaneLegAsSource_OnRampEncoderParity_20And32ByteRouterEquivalent(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := testSuiChains()
	evmRouter := common.HexToAddress(evmRouterAddress).Bytes()

	report20, err := runSourceLeg(t, env, chains, sourceLegInput(evmRouter))
	require.NoError(t, err)

	padded := make([]byte, 32)
	copy(padded[12:], evmRouter)
	report32, err := runSourceLeg(t, env, chains, sourceLegInput(padded))
	require.NoError(t, err)

	require.Equal(
		t,
		report20.Output.BatchOps[3].Transactions[0].Data,
		report32.Output.BatchOps[3].Transactions[0].Data,
		"OnRamp dest router encoding should match for 20-byte and left-padded 32-byte EVM routers",
	)
}

func TestConfigureLaneLegAsDest_MCMSBatchOp_EncoderParity(t *testing.T) {
	env := testEnvWithAddressBook(t)
	input := destLegInput()

	var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
	err := lanes.WithConnectChainsEnvironment(env, func() error {
		var execErr error
		report, execErr = cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsDest,
			testSuiChains(),
			input,
		)
		return execErr
	})
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)

	offRamp, err := module_offramp.NewOfframp(testOffRampPackageID, nil)
	require.NoError(t, err)
	encoded, err := offRamp.Encoder().ApplySourceChainConfigUpdates(
		bind.Object{Id: testCCIPObjectRef},
		bind.Object{Id: testOffRampStateObjectID},
		bind.Object{Id: testOffRampOwnerCapID},
		[]uint64{evmSepoliaSelector},
		[]bool{true},
		[]bool{false},
		[][]byte{input.Source.OnRamp},
	)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(
		t,
		report.Output.BatchOps[0].Transactions[0].Data,
		encoded,
		testOffRampStateObjectID,
		nil,
	)
}

func TestConfigureLaneLegAsDest_Flags(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := testSuiChains()

	tests := []struct {
		name            string
		isDisabled      bool
		rmnEnabled      bool
		wantEnabled     bool
		wantRMNDisabled bool
	}{
		{name: "enabled with RMN verification", isDisabled: false, rmnEnabled: true, wantEnabled: true, wantRMNDisabled: false},
		{name: "enabled without RMN verification", isDisabled: false, rmnEnabled: false, wantEnabled: true, wantRMNDisabled: true},
		{name: "disabled lane", isDisabled: true, rmnEnabled: true, wantEnabled: false, wantRMNDisabled: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := destLegInput()
			input.IsDisabled = tc.isDisabled
			input.Source.RMNVerificationEnabled = tc.rmnEnabled

			var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
			err := lanes.WithConnectChainsEnvironment(env, func() error {
				var execErr error
				report, execErr = cldf_ops.ExecuteSequence(
					mcmstest.Bundle(t),
					lanes.ConfigureLaneLegAsDest,
					chains,
					input,
				)
				return execErr
			})
			require.NoError(t, err)

			offRamp, err := module_offramp.NewOfframp(testOffRampPackageID, nil)
			require.NoError(t, err)
			encoded, err := offRamp.Encoder().ApplySourceChainConfigUpdates(
				bind.Object{Id: testCCIPObjectRef},
				bind.Object{Id: testOffRampStateObjectID},
				bind.Object{Id: testOffRampOwnerCapID},
				[]uint64{evmSepoliaSelector},
				[]bool{tc.wantEnabled},
				[]bool{tc.wantRMNDisabled},
				[][]byte{input.Source.OnRamp},
			)
			require.NoError(t, err)
			mcmstest.AssertProposalDataMatches(
				t,
				report.Output.BatchOps[0].Transactions[0].Data,
				encoded,
				testOffRampStateObjectID,
				nil,
			)
		})
	}
}

func TestConfigureLaneLegAsDest_PartialLatestPackageIDs(t *testing.T) {
	env := testEnvWithAddressBook(t)
	const latestOffRampPackageID = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"

	var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
	err := lanes.RunConnectChainsWithSuiScopes(env, map[uint64]lanes.LatestPackageIDsConfig{
		suiTestnetSelector: {OffRamp: latestOffRampPackageID},
	}, func() error {
		var execErr error
		report, execErr = cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsDest,
			testSuiChains(),
			destLegInput(),
		)
		return execErr
	})
	require.NoError(t, err)
	tx := report.Output.BatchOps[0].Transactions[0]
	require.Equal(t, testOffRampPackageID, tx.To)
	latest, err := utils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Equal(t, latestOffRampPackageID, latest)
}

func TestRunConnectChainsWithSuiScopes(t *testing.T) {
	env := testEnvWithAddressBook(t)
	adapter := &lanes.SuiAdapter{}
	const latestOffRamp = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"

	err := lanes.RunConnectChainsWithSuiScopes(env, map[uint64]lanes.LatestPackageIDsConfig{
		suiTestnetSelector: {OffRamp: latestOffRamp},
	}, func() error {
		onRamp, err := adapter.GetOnRampAddress(nil, suiTestnetSelector)
		require.NoError(t, err)
		require.Equal(t, testOnRampPackageID, "0x"+hex.EncodeToString(onRamp))

		var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
		report, err = cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsDest,
			testSuiChains(),
			destLegInput(),
		)
		require.NoError(t, err)
		latest, err := utils.TransactionLatestPackageID(report.Output.BatchOps[0].Transactions[0])
		require.NoError(t, err)
		require.Equal(t, latestOffRamp, latest)
		return nil
	})
	require.NoError(t, err)
}

func TestConfigureLaneLegAsSource_MissingAddressBookFields(t *testing.T) {
	chains := testSuiChains()
	router := common.HexToAddress(evmRouterAddress).Bytes()

	t.Run("missing link token metadata", func(t *testing.T) {
		env := testEnvOmittingAddressTypes(t, string(suideploy.SuiLinkTokenObjectMetadataID))
		err := lanes.WithConnectChainsEnvironment(env, func() error {
			_, execErr := cldf_ops.ExecuteSequence(
				mcmstest.Bundle(t),
				lanes.ConfigureLaneLegAsSource,
				chains,
				sourceLegInput(router),
			)
			return execErr
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "SuiLinkTokenObjectMetadataID")
	})
}

func TestConfigureLaneLegAsDest_MissingAddressBookFields(t *testing.T) {
	chains := testSuiChains()

	t.Run("missing offramp owner cap", func(t *testing.T) {
		env := testEnvOmittingAddressTypes(t, string(suideploy.SuiOffRampOwnerCapObjectIDType))
		err := lanes.WithConnectChainsEnvironment(env, func() error {
			_, execErr := cldf_ops.ExecuteSequence(
				mcmstest.Bundle(t),
				lanes.ConfigureLaneLegAsDest,
				chains,
				destLegInput(),
			)
			return execErr
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "SuiOffRampOwnerCapObjectID")
	})
}

func TestDisableRemoteChain_NoOp(t *testing.T) {
	report, err := cldf_ops.ExecuteSequence(
		mcmstest.Bundle(t),
		lanes.DisableRemoteChain,
		testSuiChains(),
		laneapi.DisableRemoteChainInput{
			LocalChainSelector:  suiTestnetSelector,
			RemoteChainSelector: evmSepoliaSelector,
		},
	)
	require.NoError(t, err)
	require.Empty(t, report.Output.BatchOps)
}

func TestSuiLaneAdapter_RegisteredInRegistry(t *testing.T) {
	version := semver.MustParse("1.6.0")
	adapter, ok := laneapi.GetLaneAdapterRegistry().GetLaneAdapter(chain_selectors.FamilySui, version)
	require.True(t, ok)
	require.IsType(t, &lanes.SuiAdapter{}, adapter)

	selector := adapter.(laneapi.ChainMetadataProvider).GetChainFamilySelector()
	require.Equal(t, [4]byte{0xc4, 0xe0, 0x59, 0x53}, selector)
}

func TestConfigureLaneLegAsSource_AllowListEnabled(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := testSuiChains()
	input := sourceLegInput(common.HexToAddress(evmRouterAddress).Bytes())
	input.Source.AllowListEnabled = true

	report, err := runSourceLeg(t, env, chains, input)
	require.NoError(t, err)

	onRamp, err := module_onramp.NewOnramp(testOnRampPackageID, nil)
	require.NoError(t, err)
	encoded, err := onRamp.Encoder().ApplyDestChainConfigUpdates(
		bind.Object{Id: testCCIPObjectRef},
		bind.Object{Id: testOnRampStateObjectID},
		bind.Object{Id: testOnRampOwnerCapID},
		[]uint64{evmSepoliaSelector},
		[]bool{true},
		[]string{leftPaddedAddressHex(common.HexToAddress(evmRouterAddress).Bytes())},
	)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.BatchOps[3].Transactions[0].Data, encoded, testOnRampStateObjectID, nil)
}

func leftPaddedAddressHex(b []byte) string {
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return "0x" + hex.EncodeToString(padded)
}

func testEnvOmittingAddressTypes(t *testing.T, omitTypes ...string) cldf.Environment {
	t.Helper()

	b, err := os.ReadFile("../testdata/addresses.json")
	require.NoError(t, err)

	addrsByChain := make(map[uint64]map[string]cldf.TypeAndVersion)
	require.NoError(t, json.Unmarshal(b, &addrsByChain))

	omit := make(map[string]struct{}, len(omitTypes))
	for _, typ := range omitTypes {
		omit[typ] = struct{}{}
	}

	chainAddrs := addrsByChain[suiTestnetSelector]
	filtered := make(map[string]cldf.TypeAndVersion, len(chainAddrs))
	for addr, tv := range chainAddrs {
		if _, skip := omit[string(tv.Type)]; skip {
			continue
		}
		filtered[addr] = tv
	}
	addrsByChain[suiTestnetSelector] = filtered

	return cldf.Environment{
		Name:              "test",
		ExistingAddresses: cldf.NewMemoryAddressBookFromMap(addrsByChain),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			suiTestnetSelector: sui.Chain{},
		}),
	}
}
