package adapters

import (
	"fmt"
	"testing"

	semver "github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	fees "github.com/smartcontractkit/chainlink-ccip/deployment/fees"
	lanes "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	suilanes "github.com/smartcontractkit/chainlink-sui/deployment/lanes"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	suideployutils "github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func TestSuiFeeAdapter_Registered(t *testing.T) {
	t.Parallel()
	adapter, ok := fees.GetRegistry().GetFeeAdapter(chainsel.FamilySui, semver.MustParse("1.6.0"))
	require.True(t, ok, "sui fee adapter must be registered for v1.6.0")
	require.IsType(t, &SuiFeeAdapter{}, adapter)
}

func TestSuiFeeResolver_Registered(t *testing.T) {
	t.Parallel()
	resolver, ok := fees.GetRegistry().GetFeeResolver(chainsel.FamilySui)
	require.True(t, ok, "sui fee resolver must be registered")
	require.IsType(t, &SuiFeeResolver{}, resolver)
}

func TestSuiFeeResolver_GetOnRampRef(t *testing.T) {
	t.Parallel()
	r := &SuiFeeResolver{}
	const selector uint64 = 123

	t.Run("returns OnRamp ref versioned at 1.6.0", func(t *testing.T) {
		t.Parallel()
		ds := datastore.NewMemoryDataStore()
		require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: selector,
			Type:          datastore.ContractType(suideploy.SuiOnRampType),
			Address:       "0xonramp",
			Version:       semver.MustParse("1.0.0"),
		}))
		got, err := r.GetOnRampRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, ds.Seal(), selector, 456)
		require.NoError(t, err)
		require.Equal(t, "0xonramp", got.Address)
		require.Equal(t, datastore.ContractType(suideploy.SuiOnRampType), got.Type)
		require.True(t, got.Version.Equal(semver.MustParse("1.6.0")), "onRamp ref must be versioned at 1.6.0 so the generic flow selects the adapter")
	})

	t.Run("errors when OnRamp missing", func(t *testing.T) {
		t.Parallel()
		_, err := r.GetOnRampRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.NewMemoryDataStore().Seal(), selector, 456)
		require.Error(t, err)
	})
}

func TestSuiFeeAdapter_GetDefaultTokenTransferFeeConfig(t *testing.T) {
	t.Parallel()
	a := &SuiFeeAdapter{}
	got := a.GetDefaultTokenTransferFeeConfig(123, 456)
	require.Equal(t, fees.GetDefaultChainAgnosticTokenTransferFeeConfig(123, 456), got)
}

func TestSuiFeeAdapter_GetDefaultDestChainConfig(t *testing.T) {
	t.Parallel()
	a := &SuiFeeAdapter{}
	got := a.GetDefaultDestChainConfig(123, 456)
	require.Equal(t, (&suilanes.SuiAdapter{}).GetFeeQuoterDestChainConfig(), got)
	require.True(t, got.IsEnabled)
	require.Equal(t, uint32(16_000), got.MaxDataBytes)
}

func TestSuiFeeAdapter_GetFeeContractRef(t *testing.T) {
	t.Parallel()
	a := &SuiFeeAdapter{}
	const selector uint64 = 123
	onRamp := datastore.AddressRef{ChainSelector: selector, Address: "0xonramp"}

	t.Run("empty onRamp errors", func(t *testing.T) {
		t.Parallel()
		_, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.NewMemoryDataStore().Seal(), datastore.AddressRef{}, selector, 456)
		require.Error(t, err)
	})

	t.Run("returns enriched CCIP package ref versioned at 1.6.0", func(t *testing.T) {
		t.Parallel()
		ds := datastore.NewMemoryDataStore()
		v := semver.MustParse("1.0.0")
		for _, r := range []datastore.AddressRef{
			{ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiCCIPType), Address: "0xccippkg", Version: v},
			{ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiCCIPObjectRefType), Address: "0xobjref", Version: v},
			{ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiCCIPOwnerCapObjectIDType), Address: "0xownercap", Version: v},
		} {
			require.NoError(t, ds.Addresses().Add(r))
		}
		got, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, ds.Seal(), onRamp, selector, 456)
		require.NoError(t, err)
		require.Equal(t, "0xccippkg", got.Address)
		require.Equal(t, datastore.ContractType(suideploy.SuiCCIPType), got.Type)
		require.True(t, got.Version.Equal(semver.MustParse("1.6.0")), "fee ref must be versioned at 1.6.0 so the generic flow selects the adapter")
		require.Equal(t, "0xobjref", feeRefLabelValue(got, suiFeeCCIPObjectRefLabel))
		require.Equal(t, "0xownercap", feeRefLabelValue(got, suiFeeCCIPOwnerCapLabel))
		require.Empty(t, feeRefLabelValue(got, suiFeeLatestCCIPPkgLabel), "no latest package ref present")
	})

	t.Run("includes latest package label when upgraded", func(t *testing.T) {
		t.Parallel()
		ds := datastore.NewMemoryDataStore()
		v := semver.MustParse("1.0.0")
		for _, r := range []datastore.AddressRef{
			{ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiCCIPType), Address: "0xccippkg", Version: v},
			{ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiCCIPObjectRefType), Address: "0xobjref", Version: v},
			{ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiCCIPOwnerCapObjectIDType), Address: "0xownercap", Version: v},
			{ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiLatestCCIPPackageIDType), Address: "0xlatestpkg", Version: v},
		} {
			require.NoError(t, ds.Addresses().Add(r))
		}
		got, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, ds.Seal(), onRamp, selector, 456)
		require.NoError(t, err)
		require.Equal(t, "0xlatestpkg", feeRefLabelValue(got, suiFeeLatestCCIPPkgLabel))
	})

	t.Run("errors when CCIP object ref missing", func(t *testing.T) {
		t.Parallel()
		ds := datastore.NewMemoryDataStore()
		require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: selector, Type: datastore.ContractType(suideploy.SuiCCIPType), Address: "0xccippkg", Version: semver.MustParse("1.0.0"),
		}))
		_, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, ds.Seal(), onRamp, selector, 456)
		require.Error(t, err)
	})

	t.Run("errors when CCIP package missing", func(t *testing.T) {
		t.Parallel()
		_, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.NewMemoryDataStore().Seal(), onRamp, selector, 456)
		require.Error(t, err)
	})
}

func TestFeeRefLabelValue(t *testing.T) {
	t.Parallel()
	ref := datastore.AddressRef{Labels: datastore.NewLabelSet("ccip-object-ref:0xobj", "other", "ccip-owner-cap:0xcap")}
	require.Equal(t, "0xobj", feeRefLabelValue(ref, suiFeeCCIPObjectRefLabel))
	require.Equal(t, "0xcap", feeRefLabelValue(ref, suiFeeCCIPOwnerCapLabel))
	require.Empty(t, feeRefLabelValue(ref, suiFeeLatestCCIPPkgLabel))
	require.Empty(t, feeRefLabelValue(datastore.AddressRef{}, suiFeeCCIPObjectRefLabel))
}

func TestBuildSuiTokenTransferFeeUpdate(t *testing.T) {
	t.Parallel()
	enabled := &fees.TokenTransferFeeArgs{IsEnabled: true, MinFeeUSDCents: 10, MaxFeeUSDCents: 100, DeciBps: 5, DestGasOverhead: 200, DestBytesOverhead: 32}
	disabled := &fees.TokenTransferFeeArgs{IsEnabled: false}
	u, err := buildSuiTokenTransferFeeUpdate(map[string]*fees.TokenTransferFeeArgs{
		"0xtokenA": enabled,
		"0xtokenB": disabled,
		"0xtokenC": nil,
	})
	require.NoError(t, err)
	require.Len(t, u.AddTokens, 1)
	require.Equal(t, "0xtokenA", u.AddTokens[0])
	require.Equal(t, uint32(10), u.AddMinFeeUsdCents[0])
	require.Equal(t, uint32(100), u.AddMaxFeeUsdCents[0])
	require.Equal(t, uint16(5), u.AddDeciBps[0])
	require.Equal(t, uint32(200), u.AddDestGasOverhead[0])
	require.Equal(t, uint32(32), u.AddDestBytesOverhead[0])
	require.True(t, u.AddIsEnabled[0])
	require.ElementsMatch(t, []string{"0xtokenB", "0xtokenC"}, u.RemoveTokens)

	_, err = buildSuiTokenTransferFeeUpdate(map[string]*fees.TokenTransferFeeArgs{"": enabled})
	require.Error(t, err)
}

func TestFqConfigToArgs(t *testing.T) {
	t.Parallel()
	cfg := module_fee_quoter.TokenTransferFeeConfig{
		MinFeeUsdCents:    7,
		MaxFeeUsdCents:    70,
		DeciBps:           3,
		DestGasOverhead:   400,
		DestBytesOverhead: 64,
		IsEnabled:         true,
	}
	got := fqConfigToArgs(cfg)
	require.Equal(t, fees.TokenTransferFeeArgs{
		MinFeeUSDCents:    7,
		MaxFeeUSDCents:    70,
		DeciBps:           3,
		DestGasOverhead:   400,
		DestBytesOverhead: 64,
		IsEnabled:         true,
	}, got)
}

func TestSuiFeeAdapter_SequencesAndReadGuards(t *testing.T) {
	t.Parallel()
	a := &SuiFeeAdapter{}
	// Both sequence methods return non-nil sequences.
	require.NotNil(t, a.SetTokenTransferFee(nil, datastore.AddressRef{}))
	require.NotNil(t, a.ApplyDestChainConfigUpdates(nil, datastore.AddressRef{}))

	// GetOnchainTokenTransferFeeConfig errors without a Sui chain / valid fee ref.
	_, err := a.GetOnchainTokenTransferFeeConfig(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.AddressRef{}, 1, 2, "0xtoken")
	require.Error(t, err)

	// GetOnchainDestChainConfig errors without a Sui chain / valid fee ref.
	var _ lanes.FeeQuoterDestChainConfig
	_, err = a.GetOnchainDestChainConfig(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.AddressRef{}, 1, 2)
	require.Error(t, err)
}

const (
	fqAdapterTestPkgID     = "0x1111111111111111111111111111111111111111111111111111111111111111"
	fqAdapterTestLatestPkg = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	fqAdapterTestStateObj  = "0x2222222222222222222222222222222222222222222222222222222222222222"
	fqAdapterTestOwnerCap  = "0x3333333333333333333333333333333333333333333333333333333333333333"
)

// fqAdapterTestFeeRef builds a fee ref shaped like the one GetFeeContractRef returns.
func fqAdapterTestFeeRef(selector uint64) datastore.AddressRef {
	return datastore.AddressRef{
		ChainSelector: selector,
		Address:       fqAdapterTestPkgID,
		Labels: datastore.NewLabelSet(
			suiFeeCCIPObjectRefLabel+":"+fqAdapterTestStateObj,
			suiFeeCCIPOwnerCapLabel+":"+fqAdapterTestOwnerCap,
			suiFeeLatestCCIPPkgLabel+":"+fqAdapterTestLatestPkg,
		),
	}
}

func fqAdapterTestBundle(t *testing.T) cldf_ops.Bundle {
	t.Helper()
	registry := cldf_ops.NewOperationRegistry(ccipops.FeeQuoterApplyDestChainConfigUpdatesOp.AsUntyped())
	return cldf_ops.NewBundle(
		t.Context,
		logger.Test(t),
		cldf_ops.NewMemoryReporter(),
		cldf_ops.WithOperationRegistry(registry),
	)
}

// fqAdapterTestLaneCfg returns a full lane config with distinct values per dest so
// encoded calls are distinguishable in the batch ops.
func fqAdapterTestLaneCfg(dst uint64) lanes.FeeQuoterDestChainConfig {
	return lanes.FeeQuoterDestChainConfig{
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             uint32(300_000 + dst%10),
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         0x2812d52c,
		DefaultTokenFeeUSDCents:     25,
		DefaultTokenDestGasOverhead: uint32(90_000 + dst%10),
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          10,
		V1Params: &lanes.FeeQuoterV1Params{
			MaxNumberOfTokensPerMsg:           1,
			DestGasPerPayloadByteHigh:         40,
			DestGasPerPayloadByteThreshold:    3_000,
			DestDataAvailabilityOverheadGas:   100,
			DestGasPerDataAvailabilityByte:    16,
			DestDataAvailabilityMultiplierBps: 1,
			EnforceOutOfOrder:                 true,
			GasMultiplierWeiPerEth:            1e18,
			GasPriceStalenessThreshold:        1_000_000,
		},
	}
}

// expectedDestChainConfigTx encodes the expected MCMS transaction for one dst config,
// mirroring the op's encode path: encode against the latest package head, then stamp
// the original package id as the MCMS registry identity.
func expectedDestChainConfigTx(t *testing.T, dst uint64, cfg lanes.FeeQuoterDestChainConfig) mcmstypes.Transaction {
	t.Helper()
	contract, err := module_fee_quoter.NewFeeQuoter(fqAdapterTestLatestPkg, nil)
	require.NoError(t, err)
	in := suilanes.TranslateDestChainConfig(cfg, dst)
	encodedCall, err := contract.Encoder().ApplyDestChainConfigUpdates(
		bind.Object{Id: fqAdapterTestStateObj},
		bind.Object{Id: fqAdapterTestOwnerCap},
		dst,
		in.IsEnabled,
		in.MaxNumberOfTokensPerMsg,
		in.MaxDataBytes,
		in.MaxPerMsgGasLimit,
		in.DestGasOverhead,
		in.DestGasPerPayloadByteBase,
		in.DestGasPerPayloadByteHigh,
		in.DestGasPerPayloadByteThreshold,
		in.DestDataAvailabilityOverheadGas,
		in.DestGasPerDataAvailabilityByte,
		in.DestDataAvailabilityMultiplierBps,
		in.ChainFamilySelector,
		in.EnforceOutOfOrder,
		in.DefaultTokenFeeUsdCents,
		in.DefaultTokenDestGasOverhead,
		in.DefaultTxGasLimit,
		in.GasMultiplierWeiPerEth,
		in.GasPriceStalenessThreshold,
		in.NetworkFeeUsdCents,
	)
	require.NoError(t, err)
	call, err := sui_ops.ToTransactionCall(encodedCall, fqAdapterTestStateObj)
	require.NoError(t, err)
	call.LatestPackageID = call.PackageID
	call.PackageID = fqAdapterTestPkgID
	tx, err := suideployutils.TransactionCallToMCMSTransaction(call)
	require.NoError(t, err)
	return tx
}

func TestSuiFeeAdapter_ApplyDestChainConfigUpdates(t *testing.T) {
	t.Parallel()

	selector := chainsel.SUI_TESTNET.Selector
	a := &SuiFeeAdapter{}
	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})

	// Deliberately unsorted: sepolia (larger) before jovay (smaller). The sequence
	// must emit batch ops in ascending dst-selector order for byte-deterministic
	// proposals.
	settings := map[uint64]lanes.FeeQuoterDestChainConfig{
		16015286601757825753: fqAdapterTestLaneCfg(16015286601757825753),
		945045181441419236:   fqAdapterTestLaneCfg(945045181441419236),
	}

	report, err := cldf_ops.ExecuteSequence(
		fqAdapterTestBundle(t),
		a.ApplyDestChainConfigUpdates(nil, fqAdapterTestFeeRef(selector)),
		chains,
		fees.ApplyDestChainConfigSequenceInput{Selector: selector, Settings: settings},
	)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 2)

	// One batch op per dst, each carrying exactly one transaction that matches the
	// binding-encoded call, and emitted in ascending dst selector order.
	expectedJovay := expectedDestChainConfigTx(t, 945045181441419236, settings[945045181441419236])
	expectedSepolia := expectedDestChainConfigTx(t, 16015286601757825753, settings[16015286601757825753])
	require.Equal(t, mcmstypes.ChainSelector(selector), report.Output.BatchOps[0].ChainSelector)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)
	require.Equal(t, expectedJovay, report.Output.BatchOps[0].Transactions[0])
	require.Equal(t, mcmstypes.ChainSelector(selector), report.Output.BatchOps[1].ChainSelector)
	require.Len(t, report.Output.BatchOps[1].Transactions, 1)
	require.Equal(t, expectedSepolia, report.Output.BatchOps[1].Transactions[0])

	latest, err := suideployutils.TransactionLatestPackageID(report.Output.BatchOps[1].Transactions[0])
	require.NoError(t, err)
	require.Equal(t, fqAdapterTestLatestPkg, latest, "upgraded package head must ride as the latest package id")
}

func TestSuiFeeAdapter_ApplyDestChainConfigUpdates_Guards(t *testing.T) {
	t.Parallel()

	selector := chainsel.SUI_TESTNET.Selector
	a := &SuiFeeAdapter{}
	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})
	input := fees.ApplyDestChainConfigSequenceInput{
		Selector: selector,
		Settings: map[uint64]lanes.FeeQuoterDestChainConfig{456: fqAdapterTestLaneCfg(456)},
	}

	t.Run("unknown source chain errors", func(t *testing.T) {
		t.Parallel()
		_, err := cldf_ops.ExecuteSequence(
			fqAdapterTestBundle(t),
			a.ApplyDestChainConfigUpdates(nil, fqAdapterTestFeeRef(selector)),
			chains,
			fees.ApplyDestChainConfigSequenceInput{Selector: 999, Settings: input.Settings},
		)
		require.ErrorContains(t, err, "sui chain with selector 999 not found")
	})

	t.Run("empty fee ref address errors", func(t *testing.T) {
		t.Parallel()
		_, err := cldf_ops.ExecuteSequence(
			fqAdapterTestBundle(t),
			a.ApplyDestChainConfigUpdates(nil, datastore.AddressRef{}),
			chains,
			input,
		)
		require.ErrorContains(t, err, "empty CCIP package address")
	})

	t.Run("missing labels error", func(t *testing.T) {
		t.Parallel()
		_, err := cldf_ops.ExecuteSequence(
			fqAdapterTestBundle(t),
			a.ApplyDestChainConfigUpdates(nil, datastore.AddressRef{Address: fqAdapterTestPkgID}),
			chains,
			input,
		)
		require.ErrorContains(t, err, "missing ccip-object-ref/ccip-owner-cap labels")
	})
}

func TestIsUnknownDestChainAbort(t *testing.T) {
	t.Parallel()

	// Abort descriptions embed the module and code in varying shapes depending on
	// the node; code 3 is unique to EUnknownDestChainSelector in the fee_quoter module.
	for _, desc := range []string{
		"simulate failed: Move abort in 0x1::fee_quoter: location 0x1 (line: 1556) code: 3",
		"simulate failed: Move abort in 0x1::fee_quoter code 3",
		"simulate failed: 0x1::fee_quoter aborted with code=3",
	} {
		require.True(t, isUnknownDestChainAbort(fmt.Errorf("%s", desc)), desc)
	}
	for _, desc := range []string{
		"simulate failed: Move abort in 0x1::fee_quoter: code: 2",
		"simulate failed: Move abort in 0x1::onramp: code: 3",
		"simulate failed: gas budget exceeded",
	} {
		require.False(t, isUnknownDestChainAbort(fmt.Errorf("%s", desc)), desc)
	}
	require.False(t, isUnknownDestChainAbort(nil))
}
