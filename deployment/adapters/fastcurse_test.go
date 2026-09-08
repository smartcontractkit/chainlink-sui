package adapters_test

import (
	"encoding/json"
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment/adapters"
	rmnops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
	suideployutils "github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

const (
	testCCIPPackageID       = "0x1111111111111111111111111111111111111111111111111111111111111111"
	testLatestCCIPPackageID = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testCCIPObjectRef       = "0x2222222222222222222222222222222222222222222222222222222222222222"
	testCurserCapID         = "0x4444444444444444444444444444444444444444444444444444444444444444"
)

// stubSuiSigner is a minimal signer stub used to prove the adapter never routes
// operations through the direct-execution path (which would try to submit a real
// transaction when a signer is present).
type stubSuiSigner struct{}

func (stubSuiSigner) Sign(_ []byte) ([]string, error) { return nil, nil }

func (stubSuiSigner) GetAddress() (string, error) { return "0x1", nil }

func fastCurseTestBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()

	registry := cld_ops.NewOperationRegistry(rmnops.CurseWithCurserCapOp.AsUntyped())
	return cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
		cld_ops.WithOperationRegistry(registry),
	)
}

func uncurseTestBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()

	registry := cld_ops.NewOperationRegistry(rmnops.UncurseChainOp.AsUntyped())
	return cld_ops.NewBundle(
		t.Context,
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
		cld_ops.WithOperationRegistry(registry),
	)
}

func txAdditionalFields(t *testing.T, tx mcmstypes.Transaction) suisdk.AdditionalFields {
	t.Helper()

	var fields suisdk.AdditionalFields
	require.NoError(t, json.Unmarshal(tx.AdditionalFields, &fields))
	return fields
}

func TestCurseAdapter_Curse_UsesFastCurseSequence(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:       testCCIPPackageID,
		CCIPObjectRef:     testCCIPObjectRef,
		CurserCapObjectID: testCurserCapID,
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})

	report, err := cld_ops.ExecuteSequence(
		fastCurseTestBundle(t),
		adapter.Curse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)

	tx := report.Output.BatchOps[0].Transactions[0]
	require.Equal(t, testCCIPPackageID, tx.To)

	fields := txAdditionalFields(t, tx)
	require.Equal(t, "rmn_remote", fields.ModuleName)
	require.Equal(t, "curse_multiple_with_curser_cap", fields.Function)
}

func TestCurseAdapter_Curse_RequiresRegisteredCurserCap(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:   testCCIPPackageID,
		CCIPObjectRef: testCCIPObjectRef,
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})

	_, err := cld_ops.ExecuteSequence(
		fastCurseTestBundle(t),
		adapter.Curse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "registered CurserCap is required")
}

func TestCurseAdapter_Curse_WithUpgradedPackage_PropagatesLatestPackageID(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:         testCCIPPackageID,
		LatestCCIPPackageID: testLatestCCIPPackageID,
		CCIPObjectRef:       testCCIPObjectRef,
		CurserCapObjectID:   testCurserCapID,
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})

	report, err := cld_ops.ExecuteSequence(
		fastCurseTestBundle(t),
		adapter.Curse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)

	tx := report.Output.BatchOps[0].Transactions[0]
	require.Equal(t, testCCIPPackageID, tx.To)

	latestPackageID, err := suideployutils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Equal(t, testLatestCCIPPackageID, latestPackageID)
}

func TestCurseAdapter_Uncurse_UsesOwnerCapSequence(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:          testCCIPPackageID,
		CCIPObjectRef:        testCCIPObjectRef,
		CCIPOwnerCapObjectID: "0x3333333333333333333333333333333333333333333333333333333333333333",
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})

	report, err := cld_ops.ExecuteSequence(
		uncurseTestBundle(t),
		adapter.Uncurse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)

	tx := report.Output.BatchOps[0].Transactions[0]
	require.Equal(t, testCCIPPackageID, tx.To)

	fields := txAdditionalFields(t, tx)
	require.Equal(t, "rmn_remote", fields.ModuleName)
	require.Equal(t, "uncurse_multiple", fields.Function)
}

func TestCurseAdapter_Uncurse_WithUpgradedPackage_PropagatesLatestPackageID(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:          testCCIPPackageID,
		LatestCCIPPackageID:  testLatestCCIPPackageID,
		CCIPObjectRef:        testCCIPObjectRef,
		CCIPOwnerCapObjectID: "0x3333333333333333333333333333333333333333333333333333333333333333",
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})

	report, err := cld_ops.ExecuteSequence(
		uncurseTestBundle(t),
		adapter.Uncurse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      []fastcurse.Subject{fastcurse.GenericSelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector)},
		},
	)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)

	tx := report.Output.BatchOps[0].Transactions[0]
	require.Equal(t, testCCIPPackageID, tx.To)

	latestPackageID, err := suideployutils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Equal(t, testLatestCCIPPackageID, latestPackageID)
}

func TestCurseAdapter_Uncurse_MultipleSubjects(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:          testCCIPPackageID,
		CCIPObjectRef:        testCCIPObjectRef,
		CCIPOwnerCapObjectID: "0x3333333333333333333333333333333333333333333333333333333333333333",
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selector}},
	})

	subjects := []fastcurse.Subject{
		fastcurse.GenericSelectorToSubject(cselectors.ETHEREUM_MAINNET.Selector),
		fastcurse.GenericSelectorToSubject(cselectors.POLYGON_MAINNET.Selector),
	}

	report, err := cld_ops.ExecuteSequence(
		uncurseTestBundle(t),
		adapter.Uncurse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      subjects,
		},
	)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	require.Equal(t, mcmstypes.ChainSelector(selector), report.Output.BatchOps[0].ChainSelector)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)
}

// TestCurseAdapter_Curse_ForcesProposalOnly_WhenChainHasSigner is a regression test for
// the "Object not found" failure that occurs when a Sui deployer key is loaded (chain.Signer
// non-nil) and the CurserCap is stored inside the fast MCMS Registry.
//
// The adapter is only used by the generic fastcurse framework, which always builds an MCMS
// proposal. It must therefore force ProposalOnly on the sequence input so the underlying
// operation encodes a leaf instead of attempting direct PTB assembly + submission against
// a CurserCap that has no top-level object owner.
func TestCurseAdapter_Curse_ForcesProposalOnly_WhenChainHasSigner(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:       testCCIPPackageID,
		CCIPObjectRef:     testCCIPObjectRef,
		CurserCapObjectID: testCurserCapID,
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{
			ChainMetadata: cldfsui.ChainMetadata{Selector: selector},
			Signer:        stubSuiSigner{},
		},
	})

	report, err := cld_ops.ExecuteSequence(
		fastCurseTestBundle(t),
		adapter.Curse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.NoError(t, err, "adapter.Curse() must produce a proposal leaf even when the chain has a signer")
	require.Len(t, report.Output.BatchOps, 1)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)

	tx := report.Output.BatchOps[0].Transactions[0]
	fields := txAdditionalFields(t, tx)
	require.Equal(t, "rmn_remote", fields.ModuleName)
	require.Equal(t, "curse_multiple_with_curser_cap", fields.Function)
}

// TestCurseAdapter_Uncurse_ForcesProposalOnly_WhenChainHasSigner is the uncurse counterpart
// of TestCurseAdapter_Curse_ForcesProposalOnly_WhenChainHasSigner. In production the OwnerCap
// is owned by the slow MCMS timelock; the loaded Sui deployer signer cannot authorize direct
// execution, so the adapter must always produce an MCMS proposal leaf.
func TestCurseAdapter_Uncurse_ForcesProposalOnly_WhenChainHasSigner(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selector, adapters.SuiChainState{
		CCIPAddress:          testCCIPPackageID,
		CCIPObjectRef:        testCCIPObjectRef,
		CCIPOwnerCapObjectID: "0x3333333333333333333333333333333333333333333333333333333333333333",
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selector: cldfsui.Chain{
			ChainMetadata: cldfsui.ChainMetadata{Selector: selector},
			Signer:        stubSuiSigner{},
		},
	})

	report, err := cld_ops.ExecuteSequence(
		uncurseTestBundle(t),
		adapter.Uncurse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selector,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.NoError(t, err, "adapter.Uncurse() must produce a proposal leaf even when the chain has a signer")
	require.Len(t, report.Output.BatchOps, 1)
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)

	tx := report.Output.BatchOps[0].Transactions[0]
	fields := txAdditionalFields(t, tx)
	require.Equal(t, "rmn_remote", fields.ModuleName)
	require.Equal(t, "uncurse_multiple", fields.Function)
}

// TestCurseAdapter_Curse_MultiSelector_NoClobber is the regression test for the whole point of the
// selector-keyed refactor: the fastcurse registry stores one shared CurseAdapter instance per
// family+version, so initializing/cursing a second Sui selector must not clobber the first. With
// the old flat-field design, populating selB overwrote selA's CCIPAddress and a later Curse for
// selA would route against selB's package. Here we populate two selectors with distinct
// CCIPAddresses and assert each Curse still targets its own.
func TestCurseAdapter_Curse_MultiSelector_NoClobber(t *testing.T) {
	t.Parallel()

	selA := cselectors.SUI_TESTNET.Selector
	selB := cselectors.SUI_MAINNET.Selector
	require.NotEqual(t, selA, selB)

	addrA := testCCIPPackageID
	addrB := "0x5555555555555555555555555555555555555555555555555555555555555555"

	adapter := adapters.NewCurseAdapter()
	adapter.SetChainState(selA, adapters.SuiChainState{
		CCIPAddress:       addrA,
		CCIPObjectRef:     testCCIPObjectRef,
		CurserCapObjectID: testCurserCapID,
	})
	adapter.SetChainState(selB, adapters.SuiChainState{
		CCIPAddress:       addrB,
		CCIPObjectRef:     testCCIPObjectRef,
		CurserCapObjectID: testCurserCapID,
	})

	chains := cldf_chain.NewBlockChains(map[uint64]cldf_chain.BlockChain{
		selA: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selA}},
		selB: cldfsui.Chain{ChainMetadata: cldfsui.ChainMetadata{Selector: selB}},
	})

	// Cursing selA after selB was populated must still target addrA (no clobber).
	reportA, err := cld_ops.ExecuteSequence(
		fastCurseTestBundle(t),
		adapter.Curse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selA,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.NoError(t, err)
	require.Len(t, reportA.Output.BatchOps, 1)
	require.Len(t, reportA.Output.BatchOps[0].Transactions, 1)
	require.Equal(t, addrA, reportA.Output.BatchOps[0].Transactions[0].To)

	// Cursing selB must target addrB.
	reportB, err := cld_ops.ExecuteSequence(
		fastCurseTestBundle(t),
		adapter.Curse(),
		chains,
		fastcurse.CurseInput{
			ChainSelector: selB,
			Subjects:      []fastcurse.Subject{fastcurse.GlobalCurseSubject()},
		},
	)
	require.NoError(t, err)
	require.Len(t, reportB.Output.BatchOps, 1)
	require.Len(t, reportB.Output.BatchOps[0].Transactions, 1)
	require.Equal(t, addrB, reportB.Output.BatchOps[0].Transactions[0].To)
}
