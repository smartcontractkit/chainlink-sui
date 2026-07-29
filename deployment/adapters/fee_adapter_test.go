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
	// Delegates to the Sui lane adapter (single source of truth).
	require.Equal(t, (&suilanes.SuiAdapter{}).GetFeeQuoterDestChainConfig(), got)
	require.True(t, got.IsEnabled)
	require.Equal(t, uint32(16_000), got.MaxDataBytes)
}

func TestSuiFeeAdapter_GetFeeContractRef(t *testing.T) {
	t.Parallel()
	a := &SuiFeeAdapter{}
	const selector uint64 = 123

	t.Run("empty onRamp errors", func(t *testing.T) {
		t.Parallel()
		_, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.NewMemoryDataStore().Seal(), datastore.AddressRef{}, selector, 456)
		require.Error(t, err)
	})

	t.Run("returns CCIP package ref when present", func(t *testing.T) {
		t.Parallel()
		ds := datastore.NewMemoryDataStore()
		require.NoError(t, ds.Addresses().Add(datastore.AddressRef{
			ChainSelector: selector,
			Type:          datastore.ContractType(suideploy.SuiCCIPType),
			Address:       "0xccippkg",
			Version:       semver.MustParse("1.0.0"),
		}))
		onRamp := datastore.AddressRef{ChainSelector: selector, Address: "0xonramp"}
		got, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, ds.Seal(), onRamp, selector, 456)
		require.NoError(t, err)
		require.Equal(t, "0xccippkg", got.Address)
		require.Equal(t, datastore.ContractType(suideploy.SuiCCIPType), got.Type)
	})

	t.Run("errors when CCIP package missing", func(t *testing.T) {
		t.Parallel()
		onRamp := datastore.AddressRef{ChainSelector: selector, Address: "0xonramp"}
		_, err := a.GetFeeContractRef(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.NewMemoryDataStore().Seal(), onRamp, selector, 456)
		require.Error(t, err)
	})
}

func TestSuiFeeAdapter_StubReturns(t *testing.T) {
	t.Parallel()
	a := &SuiFeeAdapter{}
	require.Nil(t, a.SetTokenTransferFee(nil, datastore.AddressRef{}))
	require.Nil(t, a.ApplyDestChainConfigUpdates(nil, datastore.AddressRef{}))

	_, err := a.GetOnchainTokenTransferFeeConfig(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.AddressRef{}, 1, 2, "0xtoken")
	require.Error(t, err)

	var _ lanes.FeeQuoterDestChainConfig
	_, err = a.GetOnchainDestChainConfig(cldf_ops.Bundle{}, cldf_chain.BlockChains{}, datastore.AddressRef{}, 1, 2)
	require.Error(t, err)
}
