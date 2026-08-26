package deployment

import (
	"errors"
	"fmt"
	"strings"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// SaveSuiAddress writes a single Sui contract address to both the legacy address
// book and the datastore, deriving identical metadata from one TypeAndVersion so
// the two stores cannot drift. The datastore AddressRef copies the address-book
// labels (token symbol, "fastcurse") so LoadOnchainStatesui reconstructs the same
// TypeAndVersion.
//
// The caller that knows the semantic qualifier provides it explicitly: it is never derived
// from the contract's own address and never derived from the label list or its order.
// Build it with the constructors in qualifiers.go TokenQualifier for token-scoped types,
// MinterCapQualifier for holder-scoped capability objects, MCMSInstance.DatastoreQualifier
// for the MCMS instances, ChainSingletonQualifier for the rest.
// The address-book labels are still copied onto the ref so the Sui loader
// reconstructs the same TypeAndVersion.
func SaveSuiAddress(
	ab *cldf.AddressBookMap,
	ds fdatastore.MutableAddressRefStore,
	chainSelector uint64,
	address string,
	tv cldf.TypeAndVersion,
	qualifier string,
) error {
	if address != "" && strings.Contains(strings.ToLower(qualifier), strings.ToLower(address)) {
		return fmt.Errorf(
			"qualifier %q for %s contains the address being written: a qualifier must identify the instance in domain terms, not by its own address",
			qualifier, tv.Type,
		)
	}
	if err := ab.Save(chainSelector, address, tv); err != nil {
		return fmt.Errorf("save to address book: %w", err)
	}
	version := tv.Version
	ref := fdatastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       address,
		Type:          fdatastore.ContractType(tv.Type),
		Version:       &version,
		Qualifier:     qualifier,
		Labels:        fdatastore.NewLabelSet(tv.Labels.List()...),
	}
	// Add, not Upsert: within one changeset run, writing the same key twice is always a bug —
	// there is no redeploy interpretation available, because both writes belong to the same
	// deployment. Upsert would drop the first ref and report success, leaving an object on
	// chain that the registry never mentions. Conflicts with rows written by *previous* runs
	// are a different question and cannot be seen from here (this store is fresh, and the
	// merge into the environment datastore happens after the changeset returns) — those are
	// checked upfront by ValidateNoDatastoreConflicts.
	if err := ds.Add(ref); err != nil {
		if errors.Is(err, fdatastore.ErrAddressRefExists) {
			return fmt.Errorf(
				"this changeset already wrote a %s: two objects cannot share one (chain, type, version, qualifier), so one of them needs a distinct qualifier: %w",
				describe(PlannedRef{Type: tv.Type, Qualifier: qualifier, Version: &version}), err,
			)
		}
		return fmt.Errorf("save to datastore: %w", err)
	}
	return nil
}
