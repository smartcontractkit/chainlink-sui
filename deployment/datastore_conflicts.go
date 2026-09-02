package deployment

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// PlannedRef is one datastore row a changeset intends to write, named by the part of the key
// that is knowable before anything is deployed.
//
// The datastore key is (chainSelector, type, version, qualifier)
type PlannedRef struct {
	Type      cldf.ContractType
	Qualifier string
	// Version defaults to Version1_0_0 when nil, which is what every Sui write site uses.
	Version *semver.Version
	// Address is optional: deploy changesets plan their keys before anything is on chain,
	// when the address is not yet knowable. Import-style changesets (e.g. RecordCurserCap)
	// do know it. When set and equal to the address already recorded under the key, the write
	// is a re-record of the same object a no-op, not a conflict.
	Address string
}

func (p PlannedRef) version() *semver.Version {
	if p.Version != nil {
		return p.Version
	}
	v := Version1_0_0
	return &v
}

// ValidateNoDatastoreConflicts reports whether any of the rows a changeset is about to write
// would take a datastore key that is already occupied.
//
// Call it from VerifyPreconditions, which the changeset runner invokes before Apply. A
// conflict caught there costs a re-run; the same conflict discovered during Apply means the
// objects are already deployed on chain, and since a changeset that fails mid-way returns an
// empty ChangesetOutput, every address it had recorded up to that point is discarded too.
//
// An occupied key is an error unless replaceExisting says the caller means to take it, a
// redeploy legitimately replaces the row a previous deployment wrote. That intent has to be
// stated rather than inferred: the two cases are indistinguishable from the datastore's side,
// and the framework resolves them identically and silently (both MemoryAddressRefStore.Upsert
// and MemoryDataStore.Merge replace on a key hit, with no error and no record of what was
// displaced). When the caller does opt in, the displaced address is logged.
//
// Two things this cannot see, both structural. It compares against the environment datastore
// as it stands before the changeset runs, so it will not catch two changesets in the same
// pipeline claiming one key that collision only materialises when their outputs are merged.
// And it does not mark the row it displaces as superseded: keeping the old row alongside the
// new one means re-qualifying it, which is a key change and belongs to the re-qualification
// tooling, not to a deploy changeset.
//
// This check reads only the datastore, not the address book, and legacy address-derived rows
// occupy different keys from their semantic replacements. The state loader gives semantic rows
// deterministic precedence while the re-qualification tooling retires the legacy rows.
func ValidateNoDatastoreConflicts(
	e cldf.Environment,
	chainSelector uint64,
	replaceExisting bool,
	plan func() ([]PlannedRef, error),
) error {
	if e.DataStore == nil {
		return nil
	}

	planned, err := plan()
	if err != nil {
		return err
	}
	if len(planned) == 0 {
		return nil
	}

	// A duplicate key inside the plan itself means Apply would write two objects under one
	// (chain, type, version, qualifier): SaveSuiAddress rejects the second write only after
	// the first object is already deployed. Catch it here instead this is always a bug, so
	// replaceExisting does not excuse it.
	type planKey struct {
		typ       cldf.ContractType
		version   string
		qualifier string
	}
	seen := make(map[planKey]struct{}, len(planned))
	for _, ref := range planned {
		k := planKey{typ: ref.Type, version: ref.version().String(), qualifier: ref.Qualifier}
		if _, dup := seen[k]; dup {
			return fmt.Errorf(
				"changeset plan writes %s twice: two objects cannot share one (chain, type, version, qualifier), so one of them needs a distinct qualifier",
				describe(ref),
			)
		}
		seen[k] = struct{}{}
	}

	type conflict struct {
		planned  PlannedRef
		existing string
	}
	var conflicts []conflict
	for _, ref := range planned {
		key := fdatastore.NewAddressRefKey(chainSelector, fdatastore.ContractType(ref.Type), ref.version(), ref.Qualifier)
		existing, err := e.DataStore.Addresses().Get(key)
		if errors.Is(err, fdatastore.ErrAddressRefNotFound) {
			continue
		}
		if err != nil {
			return fmt.Errorf("checking datastore for %s: %w", describe(ref), err)
		}
		if ref.Address != "" && ref.Address == existing.Address {
			// Re-recording the object the key already points at: a no-op, not a takeover.
			continue
		}
		conflicts = append(conflicts, conflict{planned: ref, existing: existing.Address})
	}
	if len(conflicts) == 0 {
		return nil
	}
	sort.Slice(conflicts, func(i, j int) bool { return describe(conflicts[i].planned) < describe(conflicts[j].planned) })

	if replaceExisting {
		for _, c := range conflicts {
			if e.Logger != nil {
				e.Logger.Warnw("replacing an existing datastore ref",
					"chainSelector", chainSelector,
					"ref", describe(c.planned),
					"replacedAddress", c.existing,
				)
			}
		}
		return nil
	}

	var lines []string
	for _, c := range conflicts {
		lines = append(lines, fmt.Sprintf("  %s is already held by %s", describe(c.planned), c.existing))
	}

	return fmt.Errorf(
		"datastore conflict on chain %d: %d ref(s) this changeset would write are already recorded:\n%s\n"+
			"nothing has been deployed. If this is a redeploy and taking these keys is intended, set the changeset's "+
			"ReplaceExisting flag; otherwise the qualifiers are wrong — two different objects cannot share one "+
			"(chain, type, version, qualifier)",
		chainSelector, len(conflicts), strings.Join(lines, "\n"),
	)
}

func describe(ref PlannedRef) string {
	if ref.Qualifier == ChainSingletonQualifier {
		return fmt.Sprintf("%s %s (no qualifier)", ref.Type, ref.version())
	}
	return fmt.Sprintf("%s %s qualified %q", ref.Type, ref.version(), ref.Qualifier)
}
