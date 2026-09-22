package lanes_test

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_datastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-sui/deployment/lanes"
)

const suiTestnetSelector = uint64(9762610643973837292)

func testEnvWithDatastore(t *testing.T) cldf.Environment {
	t.Helper()

	return cldf.Environment{
		Name:      "test",
		DataStore: loadTestDatastore(t),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			suiTestnetSelector: sui.Chain{},
		}),
	}
}

// loadTestDatastore loads the shared test fixture rows as datastore address refs.
func loadTestDatastore(t *testing.T) cldf_datastore.DataStore {
	t.Helper()

	b, err := os.ReadFile("../testdata/addresses.json")
	require.NoError(t, err)

	addrsByChain := make(map[uint64]map[string]cldf.TypeAndVersion)
	require.NoError(t, json.Unmarshal(b, &addrsByChain))

	return seedTestDatastore(t, addrsByChain)
}

// seedTestDatastore converts the legacy fixture into datastore refs. The address is a
// test-only qualifier so refs with the same type/version remain distinct.
func seedTestDatastore(t *testing.T, addrsByChain map[uint64]map[string]cldf.TypeAndVersion) cldf_datastore.DataStore {
	t.Helper()

	ds := cldf_datastore.NewMemoryDataStore()
	for chainSelector, addrs := range addrsByChain {
		for addr, tv := range addrs {
			version := tv.Version
			labels := tv.Labels.List()
			require.NoError(t, ds.Addresses().Upsert(cldf_datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       addr,
				Type:          cldf_datastore.ContractType(tv.Type),
				Version:       &version,
				Qualifier:     addr,
				Labels:        cldf_datastore.NewLabelSet(labels...),
			}))
		}
	}
	return ds.Seal()
}

func TestSuiAdapter_AddressGettersFromDatastore(t *testing.T) {
	env := testEnvWithDatastore(t)
	adapter := &lanes.SuiAdapter{}

	onRampPkg := "0xf87c6010be571a304f0d860857204bc66f037842156f0f6c9d80be265fd83752"
	offRampPkg := "0x9438693fb18f5660aff9277240a2282be44dc01cdd7eed4e1d8de0591ad52c03"
	ccipPkg := "0xece742a763bddf1e36629fa06b605497e413241afd14f05e558e80eef4f64e95"
	routerPkg := "0xed4613bd35004954c07150c3e9b10230f5e23e3058bc2ca0e3e676cb43eb4dc1"

	err := lanes.WithConnectChainsEnvironment(env, func() error {
		onRamp, err := adapter.GetOnRampAddress(nil, suiTestnetSelector)
		require.NoError(t, err)
		require.Equal(t, onRampPkg, "0x"+hex.EncodeToString(onRamp))

		offRamp, err := adapter.GetOffRampAddress(nil, suiTestnetSelector)
		require.NoError(t, err)
		require.Equal(t, offRampPkg, "0x"+hex.EncodeToString(offRamp))

		fq, err := adapter.GetFQAddress(nil, suiTestnetSelector)
		require.NoError(t, err)
		require.Equal(t, ccipPkg, "0x"+hex.EncodeToString(fq))

		router, err := adapter.GetRouterAddress(nil, suiTestnetSelector)
		require.NoError(t, err)
		require.Equal(t, routerPkg, "0x"+hex.EncodeToString(router))

		return nil
	})
	require.NoError(t, err)
}

func TestSuiAdapter_AddressGetters_MissingRef(t *testing.T) {
	adapter := &lanes.SuiAdapter{}
	env := cldf.Environment{
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			suiTestnetSelector: sui.Chain{},
		}),
	}

	err := lanes.WithConnectChainsEnvironment(env, func() error {
		_, err := adapter.GetOnRampAddress(nil, suiTestnetSelector)
		require.Error(t, err)
		require.Contains(t, err.Error(), "SuiOnRamp")
		return nil
	})
	require.NoError(t, err)
}

func TestSuiAdapter_AddressGetters_RequireConnectChainsScope(t *testing.T) {
	adapter := &lanes.SuiAdapter{}
	_, err := adapter.GetOnRampAddress(nil, suiTestnetSelector)
	require.Error(t, err)
	require.Contains(t, err.Error(), "WithConnectChainsEnvironment")
}
