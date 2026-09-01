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

	// The managed_token wrapper package (SuiManagedTokenPackageIDType) is a management framework
	// over TreasuryCap<T>, not a coin, so a ref stored under it must not derive a coin type.
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenPackageIDType),
		Address:       "0xmanagedpkg",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("MT"),
	}))
	// BnM coin package (stored as SuiManagedTokenType), carrying its module::STRUCT as a
	// coinType= label -> ::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN.
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xbnmpkg",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("CCIP BnM", "coinType="+suideploy.SuiCCIPBnMCoinTypeSuffix),
	}))

	sealed := ds.Seal()

	// "MT" only has a wrapper-package ref, which is not a coin source -> error.
	_, err := deriveSuiCoinType(sealed, selector, "MT")
	require.Error(t, err)

	got, err := deriveSuiCoinType(sealed, selector, "CCIP BnM")
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

// TestNormalizeSuiAddr pins the address normalization that the configure-before-own guard relies
// on to compare the on-chain pool owner against the deployer signer. The deployer still owns a
// freshly deployed pool, so owner == deployer must hold regardless of casing, leading 0x, or
// surrounding whitespace coming back from DevInspect vs GetAddress.
func TestNormalizeSuiAddr(t *testing.T) {
	t.Parallel()
	require.Equal(t, "0xabc", normalizeSuiAddr("0xABC"))
	require.Equal(t, "0xabc", normalizeSuiAddr("ABC"))
	require.Equal(t, "0xabc", normalizeSuiAddr("  0xAbC  "))
	require.Empty(t, normalizeSuiAddr(""))
	require.Empty(t, normalizeSuiAddr("   "))
}

// TestSuiAddrEqual pins the configure-before-own guard predicate: deployer-owned pools route
// EOA-direct, MCMS-owned pools keep the collect path. Equality must be insensitive to the same
// casing/0x/whitespace differences normalizeSuiAddr absorbs.
func TestSuiAddrEqual(t *testing.T) {
	t.Parallel()
	deployer := "0x40d438a47eafc6bee64a7f0addeb468d2939920f5661462f90cd8dbae2cdd9cb"

	// On-chain owner reported without a leading 0x, deployer with one -> still equal.
	require.True(t, suiAddrEqual("40d438a47eafc6bee64a7f0addeb468d2939920f5661462f90cd8dbae2cdd9cb", deployer))
	// Casing and whitespace differences collapse to the same address.
	require.True(t, suiAddrEqual("  0x40D438A47EAFC6BEE64A7F0ADDEB468D2939920F5661462F90CD8DBAE2CDD9CB  ", deployer))
	// A different owner, e.g. MCMS after accept_ownership, must not match the deployer.
	require.False(t, suiAddrEqual("0xdeadbeef", deployer))
	require.False(t, suiAddrEqual("", deployer))
}

func TestSuiTokenType(t *testing.T) {
	t.Parallel()
	require.Equal(t, datastore.ContractType(suideploy.SuiLinkTokenType),
		suiTokenType("0x1::link::LINK"))
	require.Equal(t, datastore.ContractType(suideploy.SuiManagedTokenType),
		suiTokenType("0x1::managed_token::MANAGED_TOKEN"))
	require.Equal(t, datastore.ContractType(suideploy.SuiManagedTokenType),
		suiTokenType("0x1::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN"))
	require.Equal(t, datastore.ContractType(suideploy.SuiLnRTokenType),
		suiTokenType("0x1::ccip_lock_release_token::CCIP_LOCK_RELEASE_TOKEN"))
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

func TestSuiPoolSymbol(t *testing.T) {
	t.Parallel()
	// The generic datastore-first resolver returns a raw ref whose Qualifier is the synthetic
	// "<addr>-<type>" form; the symbol is carried as the first label, so it must win.
	poolRef := datastore.AddressRef{
		Qualifier: "0xpool-SuiManagedTokenPool",
		Labels:    datastore.NewLabelSet("CCIP BnM"),
	}
	require.Equal(t, "CCIP BnM", suiPoolSymbol(poolRef))

	// With no label, fall back to the qualifier rather than returning empty.
	require.Equal(t, "0xpool-SuiManagedTokenPool",
		suiPoolSymbol(datastore.AddressRef{Qualifier: "0xpool-SuiManagedTokenPool"}))
	require.Empty(t, suiPoolSymbol(datastore.AddressRef{}))
}

// TestSuiTokenAdapter_DeriveTokenAddress_FromPoolLabels pins label-based coin-type derivation
// when the pool ref arrives from the generic resolver with a synthetic "<addr>-<type>" qualifier
// but the correct symbol label. Reading the symbol from the qualifier would fail to match any
// token package ref; reading it from the label derives the coin type.
func TestSuiTokenAdapter_DeriveTokenAddress_FromPoolLabels(t *testing.T) {
	t.Parallel()
	const selector uint64 = 123
	ds := datastore.NewMemoryDataStore()
	// The BnM coin package (ccip_burn_mint_token), stored under SuiManagedTokenType, carrying
	// its module::STRUCT as a coinType= label.
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xmanagedpkg",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("CCIP BnM", "coinType="+suideploy.SuiCCIPBnMCoinTypeSuffix),
	}))
	env := cldf.Environment{DataStore: ds.Seal()}
	a := &SuiTokenAdapter{}

	// Generic-resolver-style pool ref: synthetic qualifier, symbol as label.
	got, err := a.DeriveTokenAddress(env, selector, datastore.AddressRef{
		Type:      datastore.ContractType(suideploy.SuiManagedTokenPoolType),
		Address:   "0xpool",
		Qualifier: "0xpool-SuiManagedTokenPool",
		Labels:    datastore.NewLabelSet("CCIP BnM"),
	})
	require.NoError(t, err)
	require.Equal(t, "0xmanagedpkg::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", got)

	// Same ref without the label: falls back to the synthetic qualifier, which matches no
	// coin package ref, so derivation errors.
	_, err = a.DeriveTokenAddress(env, selector, datastore.AddressRef{
		Type:      datastore.ContractType(suideploy.SuiManagedTokenPoolType),
		Address:   "0xpool",
		Qualifier: "0xpool-SuiManagedTokenPool",
		Labels:    datastore.NewLabelSet(),
	})
	require.Error(t, err)
}

func TestCoinTypeSuffixFromLabels(t *testing.T) {
	t.Parallel()
	require.Equal(t, "usdc::USDC",
		coinTypeSuffixFromLabels(datastore.AddressRef{Labels: datastore.NewLabelSet("USDC", "coinType=usdc::USDC")}))
	require.Equal(t, "ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN",
		coinTypeSuffixFromLabels(datastore.AddressRef{Labels: datastore.NewLabelSet("coinType=ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN")}))

	// No coinType= label -> empty, even when other labels are present.
	require.Empty(t, coinTypeSuffixFromLabels(datastore.AddressRef{Labels: datastore.NewLabelSet("USDC")}))
	require.Empty(t, coinTypeSuffixFromLabels(datastore.AddressRef{Labels: datastore.NewLabelSet()}))
}

// TestDeriveSuiCoinType_CoinTypeLabel pins that a coin package ref carrying a coinType= label
// derives the coin type from that label, so BnM and a genuinely new coin (here a real USDC coin)
// derive identically through the same label-driven path. A ref with only the symbol label and no
// coinType= label does not derive.
func TestDeriveSuiCoinType_CoinTypeLabel(t *testing.T) {
	t.Parallel()
	ds := datastore.NewMemoryDataStore()
	const selector uint64 = 123

	// A genuinely new coin (USDC) filed under SuiManagedTokenType, carrying its coinType suffix.
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xusdcpkg",
		Qualifier:     "0xusdcpkg-SuiManagedToken",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("USDC", "coinType=usdc::USDC"),
	}))
	// BnM coin ref carrying its coinType label, the same shape as the USDC ref above.
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xbnmpkg",
		Qualifier:     "0xbnmpkg-SuiManagedToken",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("CCIP BnM", "coinType=ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN"),
	}))
	// LnR coin ref carrying its coinType label, the same shape as the BnM ref above.
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiLnRTokenType),
		Address:       "0xlnrpkg",
		Qualifier:     "0xlnrpkg-SuiLnRToken",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("CCIP LnR", "coinType="+suideploy.SuiCCIPLnRCoinTypeSuffix),
	}))
	// A ref with only the symbol label and no coinType= label must not derive.
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xunlabeledpkg",
		Qualifier:     "0xunlabeledpkg-SuiManagedToken",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("UNLABELED"),
	}))

	sealed := ds.Seal()

	got, err := deriveSuiCoinType(sealed, selector, "USDC")
	require.NoError(t, err)
	require.Equal(t, "0xusdcpkg::usdc::USDC", got)

	// BnM derives through the identical label-driven path as the new coin.
	got, err = deriveSuiCoinType(sealed, selector, "CCIP BnM")
	require.NoError(t, err)
	require.Equal(t, "0xbnmpkg::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", got)

	// LnR derives through the identical label-driven path as BnM and the new coin.
	got, err = deriveSuiCoinType(sealed, selector, "CCIP LnR")
	require.NoError(t, err)
	require.Equal(t, "0xlnrpkg::ccip_lock_release_token::CCIP_LOCK_RELEASE_TOKEN", got)

	// A coin package ref without a coinType= label does not derive.
	_, err = deriveSuiCoinType(sealed, selector, "UNLABELED")
	require.Error(t, err)
}

// TestSuiTokenAdapter_DeriveTokenAddress_NewCoinFromLabels pins end-to-end coin-type
// derivation for a new managed token whose pool ref arrives from the generic resolver (synthetic
// qualifier, symbol label) and whose coin package ref carries a coinType= label. No explicit
// tokenRef is needed; the correct coin type is derived purely from the datastore.
func TestSuiTokenAdapter_DeriveTokenAddress_NewCoinFromLabels(t *testing.T) {
	t.Parallel()
	const selector uint64 = 123
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xusdcpkg",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("USDC", "coinType=usdc::USDC"),
	}))
	env := cldf.Environment{DataStore: ds.Seal()}
	a := &SuiTokenAdapter{}

	got, err := a.DeriveTokenAddress(env, selector, datastore.AddressRef{
		Type:      datastore.ContractType(suideploy.SuiManagedTokenPoolType),
		Address:   "0xpool",
		Qualifier: "0xpool-SuiManagedTokenPool",
		Labels:    datastore.NewLabelSet("USDC"),
	})
	require.NoError(t, err)
	require.Equal(t, "0xusdcpkg::usdc::USDC", got)
}

// TestSuiTokenAdapter_DeriveTokenAddress_PoolFamilyAgnostic pins that coin-type derivation keys
// on the coin, not the pool family: the same USDC coin package ref serves a BurnMint pool and a
// Managed pool labelled USDC, and both derive the identical coin type.
func TestSuiTokenAdapter_DeriveTokenAddress_PoolFamilyAgnostic(t *testing.T) {
	t.Parallel()
	const selector uint64 = 123
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenType),
		Address:       "0xusdcpkg",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("USDC", "coinType=usdc::USDC"),
	}))
	env := cldf.Environment{DataStore: ds.Seal()}
	a := &SuiTokenAdapter{}

	burnMintPool := datastore.AddressRef{
		Type:      datastore.ContractType(suideploy.SuiBnMTokenPoolType),
		Address:   "0xbnmPool",
		Qualifier: "0xbnmPool-SuiBnMTokenPool",
		Labels:    datastore.NewLabelSet("USDC"),
	}
	managedPool := datastore.AddressRef{
		Type:      datastore.ContractType(suideploy.SuiManagedTokenPoolType),
		Address:   "0xmanagedPool",
		Qualifier: "0xmanagedPool-SuiManagedTokenPool",
		Labels:    datastore.NewLabelSet("USDC"),
	}

	gotBnM, err := a.DeriveTokenAddress(env, selector, burnMintPool)
	require.NoError(t, err)
	gotManaged, err := a.DeriveTokenAddress(env, selector, managedPool)
	require.NoError(t, err)

	require.Equal(t, "0xusdcpkg::usdc::USDC", gotBnM)
	require.Equal(t, "0xusdcpkg::usdc::USDC", gotManaged)
	require.Equal(t, gotBnM, gotManaged)
}

// TestResolveSuiPoolObjects_ManagedByLabel pins that pool state and owner-cap objects for a
// managed token pool are resolved by symbol label, independent of the underlying coin type.
func TestResolveSuiPoolObjects_ManagedByLabel(t *testing.T) {
	t.Parallel()
	const selector uint64 = 123
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenPoolStateType),
		Address:       "0xusdcstate",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("USDC"),
	}))
	require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: selector,
		Type:          datastore.ContractType(suideploy.SuiManagedTokenPoolOwnerIDType),
		Address:       "0xusdcowner",
		Version:       semver.MustParse("1.0.0"),
		Labels:        datastore.NewLabelSet("USDC"),
	}))

	stateObjID, ownerCapID, err := resolveSuiPoolObjects(ds.Seal(), selector,
		datastore.ContractType(suideploy.SuiManagedTokenPoolType), "USDC")
	require.NoError(t, err)
	require.Equal(t, "0xusdcstate", stateObjID)
	require.Equal(t, "0xusdcowner", ownerCapID)

	// Missing owner-cap ref for another symbol -> error.
	_, _, err = resolveSuiPoolObjects(ds.Seal(), selector,
		datastore.ContractType(suideploy.SuiManagedTokenPoolType), "nope")
	require.Error(t, err)
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

	lnr, err := suiPoolTypeFromStr("lnr")
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiLnRTokenPoolType), lnr)

	// Generic cross-family contract-type strings used by EVM/Solana YAMLs map to the Sui pool
	// types so a single token_expansion YAML can use one poolType across families.
	bnmGeneric, err := suiPoolTypeFromStr("BurnMintTokenPool")
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiBnMTokenPoolType), bnmGeneric)

	lnrGeneric, err := suiPoolTypeFromStr("LockReleaseTokenPool")
	require.NoError(t, err)
	require.Equal(t, datastore.ContractType(suideploy.SuiLnRTokenPoolType), lnrGeneric)

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
