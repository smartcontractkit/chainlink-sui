package adapters

import (
	"testing"

	semver "github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	fees "github.com/smartcontractkit/chainlink-ccip/deployment/fees"
	lanes "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"

	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	suilanes "github.com/smartcontractkit/chainlink-sui/deployment/lanes"
	"github.com/stretchr/testify/require"
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

func TestSuiFeeAdapter_ImplementedAndStubs(t *testing.T) {
	t.Parallel()
	a := &SuiFeeAdapter{}
	// Both sequence methods return non-nil sequences.
	require.NotNil(t, a.SetTokenTransferFee(nil, datastore.AddressRef{}))
	require.NotNil(t, a.ApplyDestChainConfigUpdates(nil, datastore.AddressRef{}))

	// GetOnchainTokenTransferFeeConfig errors without a Sui chain / valid fee ref.
	_, err := a.GetOnchainTokenTransferFeeConfig(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.AddressRef{}, 1, 2, "0xtoken")
	require.Error(t, err)

	// Still-stubbed read errors.
	var _ lanes.FeeQuoterDestChainConfig
	_, err = a.GetOnchainDestChainConfig(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.AddressRef{}, 1, 2)
	require.Error(t, err)
}
