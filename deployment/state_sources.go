package deployment

import (
	"fmt"
	"strings"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type suiAddressRef struct {
	address   string
	tv        cldf.TypeAndVersion
	qualifier string
}

// addressesForSuiChain loads address/type metadata for one Sui chain from
// env.DataStore (address_refs.json). A chain with no refs yields empty state.
func addressesForSuiChain(env cldf.Environment, chainSelector uint64) (map[string][]suiAddressRef, error) {
	addresses, ok, err := addressesForSuiChainFromDatastore(env, chainSelector)
	if err != nil {
		return nil, err
	}
	if !ok {
		return make(map[string][]suiAddressRef), nil
	}
	return addresses, nil
}

func addressesForSuiChainFromDatastore(env cldf.Environment, chainSelector uint64) (map[string][]suiAddressRef, bool, error) {
	if env.DataStore == nil {
		return nil, false, nil
	}
	refs := env.DataStore.Addresses().Filter(fdatastore.AddressRefByChainSelector(chainSelector))
	if len(refs) == 0 {
		return nil, false, nil
	}
	addresses := make(map[string][]suiAddressRef, len(refs))
	for _, ref := range refs {
		if ref.Labels.Contains(SupersededLabel) {
			continue
		}
		tv, err := typeAndVersionFromDatastoreRef(ref)
		if err != nil {
			return nil, false, fmt.Errorf("datastore ref %s chain %d: %w", ref.Address, chainSelector, err)
		}
		addresses[ref.Address] = append(addresses[ref.Address], suiAddressRef{
			address:   ref.Address,
			tv:        tv,
			qualifier: ref.Qualifier,
		})
	}
	return addresses, true, nil
}

func typeAndVersionFromDatastoreRef(ref fdatastore.AddressRef) (cldf.TypeAndVersion, error) {
	if strings.TrimSpace(string(ref.Type)) == "" {
		return cldf.TypeAndVersion{}, fmt.Errorf("contract type is empty")
	}
	version := Version1_0_0
	if ref.Version != nil {
		version = *ref.Version
	}
	tv := cldf.NewTypeAndVersion(cldf.ContractType(ref.Type), version)
	for _, label := range datastoreRefLabels(ref) {
		tv.Labels.Add(label)
	}
	return tv, nil
}

// datastoreRefLabels returns the labels needed to rebuild a TypeAndVersion.
func datastoreRefLabels(ref fdatastore.AddressRef) []string {
	if !ref.Labels.IsEmpty() {
		return ref.Labels.List()
	}
	return nil
}
