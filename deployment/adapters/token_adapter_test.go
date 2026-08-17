package adapters

import (
	"math/big"
	"testing"

	semver "github.com/Masterminds/semver/v3"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"

	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmstest "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
	"github.com/stretchr/testify/require"
)

func TestSuiTokenAdapter_AddressRefToBytes(t *testing.T) {
	t.Parallel()
	a := &SuiTokenAdapter{}

	got, err := a.AddressRefToBytes(datastore.AddressRef{Address: "0xdeadbeef"})
	require.NoError(t, err)
	require.Equal(t, []byte{0xde, 0xad, 0xbe, 0xef}, got)

	coinType := "0x1::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN"
	got, err = a.AddressRefToBytes(datastore.AddressRef{Address: coinType})
	require.NoError(t, err)
	require.Equal(t, []byte(coinType), got)

	_, err = a.AddressRefToBytes(datastore.AddressRef{Address: ""})
	require.Error(t, err)
}

func TestSuiTokenAdapter_StubReturns(t *testing.T) {
	t.Parallel()
	a := &SuiTokenAdapter{}

	require.Nil(t, a.MigrateLockReleasePoolLiquiditySequence())
	require.NotNil(t, a.ManualRegistration())
	require.NotNil(t, a.DeployToken())
	require.NotNil(t, a.DeployTokenPoolForToken())
	require.NotNil(t, a.ConfigureTokenForTransfersSequence())
	require.NotNil(t, a.SetTokenPoolRateLimits())
	require.NotNil(t, a.UpdateAuthorities())
	require.NoError(t, a.DeployTokenVerify(cldf.Environment{}, tokensapi.DeployTokenInput{}))

	pool, err := a.DeriveTokenPoolCounterpart(cldf.Environment{}, 0, []byte{1, 2, 3}, []byte{4})
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3}, pool)

	// DeriveTokenAddress errors when no symbol/package is resolvable.
	_, err = a.DeriveTokenAddress(cldf.Environment{}, 0, datastore.AddressRef{})
	require.Error(t, err)
}

func TestSuiTokenAdapter_DeployToken_Errors(t *testing.T) {
	t.Parallel()
	a := &SuiTokenAdapter{}

	_, err := cldf_ops.ExecuteSequence(
		mcmstest.Bundle(t),
		a.DeployToken(),
		cldf_chain.BlockChains{},
		tokensapi.DeployTokenInput{},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not supported")
}

func TestDeriveSuiCoinType(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()
	const selector uint64 = 123

	// Managed token package -> 0x<pkg>::managed_token::MANAGED_TOKEN
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenPackageIDType),
		Address:       "0xmanagedpkg",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("MT"),
	}))
	// BnM token package (stored as SuiManagedTokenType) -> ::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xbnmpkg",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("CCIP BnM"),
	}))

	sealed := ds.Seal()

	got, err := deriveSuiCoinType(sealed, selector, "MT")
	require.NoError(t, err)
	require.Equal(t, "0xmanagedpkg::managed_token::MANAGED_TOKEN", got)

	got, err = deriveSuiCoinType(sealed, selector, "CCIP BnM")
	require.NoError(t, err)
	require.Equal(t, "0xbnmpkg::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", got)

	// Unknown symbol -> error.
	_, err = deriveSuiCoinType(sealed, selector, "nope")
	require.Error(t, err)

	// Empty symbol -> error.
	_, err = deriveSuiCoinType(sealed, selector, "")
	require.Error(t, err)
}

func TestBigToU64(t *testing.T) {
	t.Parallel()
	v, err := bigToU64(nil)
	require.NoError(t, err)
	require.Zero(t, v)

	v, err = bigToU64(new(big.Int).SetUint64(42))
	require.NoError(t, err)
	require.Equal(t, uint64(42), v)

	_, err = bigToU64(new(big.Int).Lsh(big.NewInt(1), 64))
	require.Error(t, err)

	_, err = bigToU64(big.NewInt(-1))
	require.Error(t, err)
}

func TestNormalizeCoinType(t *testing.T) {
	t.Parallel()
	require.Equal(t, "0x1::m::T", normalizeCoinType("0x1::m::T"))
	require.Equal(t, "0x1::m::T", normalizeCoinType("1::m::T"))
}

func TestSuiTokenType(t *testing.T) {
	t.Parallel()
	require.Equal(t, datastore.ContractType(suideploy.SuiLinkTokenType),
		suiTokenType("0x1::link::LINK"))
	require.Equal(t, datastore.ContractType(suideploy.SuiManagedTokenType),
		suiTokenType("0x1::managed_token::MANAGED_TOKEN"))
	require.Equal(t, datastore.ContractType(suideploy.SuiManagedTokenType),
		suiTokenType("0x1::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN"))
}

func TestPoolObjectTypes(t *testing.T) {
	t.Parallel()
	st, ot, err := poolObjectTypes(datastore.ContractType(suideploy.SuiBnMTokenPoolType))
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolStateType), st)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolOwnerIDType), ot)

	st, ot, err = poolObjectTypes(datastore.ContractType(suideploy.SuiManagedTokenPoolType))
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiManagedTokenPoolStateType), st)
	require.Equal(t, datastore.ContractType(suideploy.SuiManagedTokenPoolOwnerIDType), ot)

	st, ot, err = poolObjectTypes(datastore.ContractType(suideploy.SuiLnRTokenPoolType))
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiLnRTokenPoolStateType), st)
	require.Equal(t, datastore.ContractType(suideploy.SuiLnRTokenPoolOwnerIDType), ot)

	_, _, err = poolObjectTypes(datastore.ContractType("unknown"))
	require.Error(t, err)
}

func TestIsSuiPoolType(t *testing.T) {
	t.Parallel()
	require.True(t, isSuiPoolType(datastore.ContractType(suideploy.SuiBnMTokenPoolType)))
	require.True(t, isSuiPoolType(datastore.ContractType(suideploy.SuiManagedTokenPoolType)))
	require.True(t, isSuiPoolType(datastore.ContractType(suideploy.SuiLnRTokenPoolType)))
	require.False(t, isSuiPoolType(datastore.ContractType(suideploy.SuiCCIPType)))
}

func TestSymbolFromLabels(t *testing.T) {
	t.Parallel()
	labels := datastore.NewLabelSet("BnM")
	require.Equal(t, "BnM", symbolFromLabels(datastore.AddressRef{Labels: labels}))
	require.Empty(t, symbolFromLabels(datastore.AddressRef{Labels: datastore.NewLabelSet()}))
}

func TestBatchOpFromCall(t *testing.T) {
	t.Parallel()
	call := sui_ops.TransactionCall{
		PackageID:  "0x" + repeatHex(64),
		Module:     "burn_mint_token_pool",
		Function:   "set_chain_rate_limiter_configs",
		Data:       []byte{0x00},
		StateObjID: "0x" + repeatHex(64),
		TypeArgs:   []string{"0x1::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN"},
	}
	out, err := batchOpFromCall(123, call)
	require.NoError(t, err)
	require.Len(t, out.BatchOps, 1)
	require.NotEmpty(t, out.BatchOps[0].Transactions)
}

func repeatHex(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '0'
	}
	return string(b)
}

func TestSuiPoolTypeFromStr(t *testing.T) {
	t.Parallel()
	bnm, err := suiPoolTypeFromStr("bnm")
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolType), bnm)

	bnm2, err := suiPoolTypeFromStr(string(suideploy.SuiBnMTokenPoolType))
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolType), bnm2)

	mng, err := suiPoolTypeFromStr("managed")
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiManagedTokenPoolType), mng)

	_, err = suiPoolTypeFromStr("nonsense")
	require.Error(t, err)
}

func TestAppendSuiPoolAddresses(t *testing.T) {
	t.Parallel()
	refs := appendSuiPoolAddresses(nil, 7, datastore.ContractType(suideploy.SuiBnMTokenPoolType), "BnM", "0xpool", "0xstate", "0xownercap")
	require.Len(t, refs, 3)
	require.Equal(t, "0xpool", refs[0].Address)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolType), refs[0].Type)
	require.Equal(t, "0xstate", refs[1].Address)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolStateType), refs[1].Type)
	require.Equal(t, "0xownercap", refs[2].Address)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolOwnerIDType), refs[2].Type)
	for _, r := range refs {
		require.Equal(t, uint64(7), r.ChainSelector)
		require.True(t, r.Labels.Contains("BnM"))
		require.NotNil(t, r.Version)
	}
}

func TestRefAddressHelpers(t *testing.T) {
	t.Parallel()
	require.Empty(t, firstRefAddress(nil))
	require.Equal(t, "0xa", firstRefAddress([]datastore.AddressRef{{Address: "0xa"}}))
	require.Empty(t, refAddress(datastore.AddressRef{}, false))
	require.Equal(t, "0xb", refAddress(datastore.AddressRef{Address: "0xb"}, true))
}

func TestSuiTokenAdapter_ManualRegistration_Noop(t *testing.T) {
	t.Parallel()
	a := &SuiTokenAdapter{}
	report, err := cldf_ops.ExecuteSequence(
		mcmstest.Bundle(t),
		a.ManualRegistration(),
		cldf_chain.BlockChains{},
		tokensapi.ManualRegistrationSequenceInput{
			RegisterTokenConfig: tokensapi.RegisterTokenConfig{ChainSelector: 123},
		},
	)
	require.NoError(t, err)
	require.Empty(t, report.Output.BatchOps)
	require.Empty(t, report.Output.Addresses)
}
