package changesets

import (
	"testing"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

// seedDatastore converts legacy test fixtures into datastore refs. The address is a
// test-only qualifier so refs with the same type/version remain distinct.
func seedDatastore(t *testing.T, addrsByChain map[uint64]map[string]cldf.TypeAndVersion) *fdatastore.MemoryDataStore {
	t.Helper()

	ds := fdatastore.NewMemoryDataStore()
	for chainSelector, addrs := range addrsByChain {
		for addr, tv := range addrs {
			version := tv.Version
			labels := tv.Labels.List()
			require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       addr,
				Type:          fdatastore.ContractType(tv.Type),
				Version:       &version,
				Qualifier:     addr,
				Labels:        fdatastore.NewLabelSet(labels...),
			}))
		}
	}
	return ds
}
