package deployment

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// SaveSuiAddress writes a single Sui contract address to both the legacy address
// book and the datastore, deriving identical metadata from one TypeAndVersion so
// the two stores cannot drift. The datastore AddressRef uses the Sui convention
// qualifier "<address>-<type>" and copies the address-book labels (token symbol,
// "fastcurse") so LoadOnchainStatesui reconstructs the same TypeAndVersion.
func SaveSuiAddress(
	ab *cldf.AddressBookMap,
	ds fdatastore.MutableAddressRefStore,
	chainSelector uint64,
	address string,
	tv cldf.TypeAndVersion,
) error {
	if err := ab.Save(chainSelector, address, tv); err != nil {
		return fmt.Errorf("save to address book: %w", err)
	}
	version := tv.Version
	ref := fdatastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       address,
		Type:          fdatastore.ContractType(tv.Type),
		Version:       &version,
		Qualifier:     fmt.Sprintf("%s-%s", address, tv.Type),
		Labels:        fdatastore.NewLabelSet(tv.Labels.List()...),
	}
	if err := ds.Upsert(ref); err != nil {
		return fmt.Errorf("save to datastore: %w", err)
	}
	return nil
}
