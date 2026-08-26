package deployment

import (
	"errors"
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

func conflictEnv(t *testing.T, selector uint64, refs ...fdatastore.AddressRef) cldf.Environment {
	t.Helper()

	ds := fdatastore.NewMemoryDataStore()
	for _, ref := range refs {
		require.NoError(t, ds.Addresses().Upsert(ref))
	}

	return cldf.Environment{DataStore: ds.Seal(), Logger: logger.Test(t)}
}

func minterCapRef(selector uint64, address, qualifier string) fdatastore.AddressRef {
	version := Version1_0_0
	return fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       address,
		Type:          fdatastore.ContractType(SuiManagedTokenMinterCapID),
		Version:       &version,
		Qualifier:     qualifier,
	}
}

func TestValidateNoDatastoreConflicts_reportsOccupiedKey(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	qualifier := MinterCapQualifier("CCIP BnM", "0xminter")
	e := conflictEnv(t, selector, minterCapRef(selector, "0xexisting-cap", qualifier))

	err := ValidateNoDatastoreConflicts(e, selector, false, func() ([]PlannedRef, error) {
		return []PlannedRef{{Type: SuiManagedTokenMinterCapID, Qualifier: qualifier}}, nil
	})

	// The error has to be actionable without a datastore dump: which ref, and what already
	// holds it, so the operator can tell a wrong qualifier from an intended redeploy.
	require.ErrorContains(t, err, "SuiManagedTokenMinterCapID")
	require.ErrorContains(t, err, `qualified "CCIP-BnM-0xminter"`)
	require.ErrorContains(t, err, "0xexisting-cap")
	require.ErrorContains(t, err, "nothing has been deployed")
	require.ErrorContains(t, err, "ReplaceExisting")
}

func TestValidateNoDatastoreConflicts_allowsFreeKey(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	e := conflictEnv(t, selector, minterCapRef(selector, "0xexisting-cap", MinterCapQualifier("CCIP BnM", "0xminter")))

	// A different holder is a different key, which is the whole point of holder scoping.
	require.NoError(t, ValidateNoDatastoreConflicts(e, selector, false, func() ([]PlannedRef, error) {
		return []PlannedRef{{
			Type:      SuiManagedTokenMinterCapID,
			Qualifier: MinterCapQualifier("CCIP BnM", "0xanother-minter"),
		}}, nil
	}))
}

func TestValidateNoDatastoreConflicts_replaceExistingIsOptIn(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	qualifier := MinterCapQualifier("CCIP BnM", "0xminter")
	e := conflictEnv(t, selector, minterCapRef(selector, "0xexisting-cap", qualifier))

	plan := func() ([]PlannedRef, error) {
		return []PlannedRef{{Type: SuiManagedTokenMinterCapID, Qualifier: qualifier}}, nil
	}

	require.Error(t, ValidateNoDatastoreConflicts(e, selector, false, plan))
	require.NoError(t, ValidateNoDatastoreConflicts(e, selector, true, plan))
}

// The plan may need on-chain reads to resolve a qualifier. Nothing can conflict on a chain
// with no refs recorded, so the check must not make the caller pay for them.
func TestValidateNoDatastoreConflicts_skipsPlanWhenChainHasNoRefs(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	otherSelector := cselectors.SUI_MAINNET.Selector

	planned := false
	plan := func() ([]PlannedRef, error) {
		planned = true
		return nil, errors.New("plan should not have run")
	}

	// Empty datastore.
	require.NoError(t, ValidateNoDatastoreConflicts(conflictEnv(t, selector), selector, false, plan))
	require.False(t, planned)

	// Refs exist, but for a different chain.
	e := conflictEnv(t, otherSelector, minterCapRef(otherSelector, "0xcap", "CCIP-BnM-0xminter"))
	require.NoError(t, ValidateNoDatastoreConflicts(e, selector, false, plan))
	require.False(t, planned)

	// No datastore at all.
	require.NoError(t, ValidateNoDatastoreConflicts(cldf.Environment{}, selector, false, plan))
	require.False(t, planned)
}

func TestValidateNoDatastoreConflicts_propagatesPlanError(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	e := conflictEnv(t, selector, minterCapRef(selector, "0xcap", "CCIP-BnM-0xminter"))

	sentinel := errors.New("failed to read coin symbol")
	err := ValidateNoDatastoreConflicts(e, selector, false, func() ([]PlannedRef, error) {
		return nil, sentinel
	})
	require.ErrorIs(t, err, sentinel)
}

// Two planned refs on one key means Apply would deploy an object and then fail to record it
// (SaveSuiAddress rejects the second write). That must surface before anything is deployed,
// and replaceExisting must not excuse it: the collision is inside the plan, not against
// earlier runs.
func TestValidateNoDatastoreConflicts_rejectsDuplicateKeysWithinPlan(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	qualifier := MinterCapQualifier("CCIP BnM", "0xminter")

	// The chain needs one recorded ref for the plan to run at all: on an empty chain the
	// check is skipped (see TestValidateNoDatastoreConflicts_skipsPlanWhenChainHasNoRefs) and
	// a duplicate write is caught by SaveSuiAddress in Apply instead.
	e := conflictEnv(t, selector, minterCapRef(selector, "0xcap", MinterCapQualifier("CCIP BnM", "0xanother-minter")))

	plan := func() ([]PlannedRef, error) {
		return []PlannedRef{
			{Type: SuiManagedTokenMinterCapID, Qualifier: qualifier},
			{Type: SuiManagedTokenMinterCapID, Qualifier: qualifier},
		}, nil
	}

	for _, replaceExisting := range []bool{false, true} {
		err := ValidateNoDatastoreConflicts(e, selector, replaceExisting, plan)
		require.ErrorContains(t, err, "twice")
		require.ErrorContains(t, err, "SuiManagedTokenMinterCapID")
	}
}

// Import-style changesets know the address upfront. Re-recording the object a key already
// points at is a no-op; the key pointing at a *different* object is the conflict.
func TestValidateNoDatastoreConflicts_sameAddressRerecordIsNoOp(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	e := conflictEnv(t, selector, minterCapRef(selector, "0xexisting-cap", ChainSingletonQualifier))

	// Same address under the key: no conflict, no flag needed.
	require.NoError(t, ValidateNoDatastoreConflicts(e, selector, false, func() ([]PlannedRef, error) {
		return []PlannedRef{{Type: SuiManagedTokenMinterCapID, Qualifier: ChainSingletonQualifier, Address: "0xexisting-cap"}}, nil
	}))

	// A different address under the key: conflict, and Sui addresses are case-sensitive.
	for _, addr := range []string{"0xother-cap", "0xEXISTING-cap"} {
		err := ValidateNoDatastoreConflicts(e, selector, false, func() ([]PlannedRef, error) {
			return []PlannedRef{{Type: SuiManagedTokenMinterCapID, Qualifier: ChainSingletonQualifier, Address: addr}}, nil
		})
		require.ErrorContains(t, err, "0xexisting-cap", "address %q must conflict", addr)
	}
}
