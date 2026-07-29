package adapters

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"

	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
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
	require.Nil(t, a.ConfigureTokenForTransfersSequence())
	require.Nil(t, a.ManualRegistration())
	require.Nil(t, a.DeployToken())
	require.Nil(t, a.DeployTokenPoolForToken())
	require.NoError(t, a.DeployTokenVerify(cldf.Environment{}, tokensapi.DeployTokenInput{}))

	pool, err := a.DeriveTokenPoolCounterpart(cldf.Environment{}, 0, []byte{1, 2, 3}, []byte{4})
	require.NoError(t, err)
	require.Equal(t, []byte{1, 2, 3}, pool)

	_, err = a.DeriveTokenAddress(cldf.Environment{}, 0, datastore.AddressRef{})
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
