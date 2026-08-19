package deployment

import (
	"errors"
	"fmt"
	"strings"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// addressesForSuiChain loads address/type metadata for one Sui chain.
// It prefers env.DataStore address refs (address_refs.json) when present and
// falls back to env.ExistingAddresses (addresses.json).
func addressesForSuiChain(env cldf.Environment, chainSelector uint64) (map[string][]cldf.TypeAndVersion, error) {
	if addresses, ok, err := addressesForSuiChainFromDatastore(env, chainSelector); err != nil {
		return nil, err
	} else if ok {
		return addresses, nil
	}
	return addressesForSuiChainFromAddressBook(env, chainSelector)
}

func addressesForSuiChainFromDatastore(env cldf.Environment, chainSelector uint64) (map[string][]cldf.TypeAndVersion, bool, error) {
	if env.DataStore == nil {
		return nil, false, nil
	}
	refs := env.DataStore.Addresses().Filter(fdatastore.AddressRefByChainSelector(chainSelector))
	if len(refs) == 0 {
		return nil, false, nil
	}
	// One Sui object address can carry several typed refs — the MCMS state object id
	// is reused by the generic Proposer/Canceller/Bypasser role refs — so keep every ref
	// per address instead of last-write-wins, which would drop the MCMS state entry.
	addresses := make(map[string][]cldf.TypeAndVersion, len(refs))
	for _, ref := range refs {
		tv, err := typeAndVersionFromDatastoreRef(ref)
		if err != nil {
			return nil, false, fmt.Errorf("datastore ref %s chain %d: %w", ref.Address, chainSelector, err)
		}
		addresses[ref.Address] = append(addresses[ref.Address], tv)
	}
	return addresses, true, nil
}

func addressesForSuiChainFromAddressBook(env cldf.Environment, chainSelector uint64) (map[string][]cldf.TypeAndVersion, error) {
	abAddresses, err := env.ExistingAddresses.AddressesForChain(chainSelector)
	if err != nil {
		if errors.Is(err, cldf.ErrChainNotFound) {
			return make(map[string][]cldf.TypeAndVersion), nil
		}
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chainSelector, err)
	}
	addresses := make(map[string][]cldf.TypeAndVersion, len(abAddresses))
	for addr, tv := range abAddresses {
		addresses[addr] = []cldf.TypeAndVersion{tv}
	}
	return addresses, nil
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

func datastoreRefLabels(ref fdatastore.AddressRef) []string {
	if !ref.Labels.IsEmpty() {
		return ref.Labels.List()
	}
	qualifier := strings.TrimSpace(ref.Qualifier)
	if qualifier == MCMSFastCurseLabel {
		return []string{MCMSFastCurseLabel}
	}
	return nil
}
