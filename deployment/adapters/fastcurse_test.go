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
	suideployutils "github.com/smartcontractkit/chainlink-sui/deployment/utils"
	rmnops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
)

const (
	testCCIPPackageID        = "0x1111111111111111111111111111111111111111111111111111111111111111"
	testLatestCCIPPackageID  = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testCCIPObjectRef        = "0x2222222222222222222222222222222222222222222222222222222222222222"
	testCurserCapID          = "0x4444444444444444444444444444444444444444444444444444444444444444"
)

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
	adapter.CCIPAddress = testCCIPPackageID
	adapter.CCIPObjectRef = testCCIPObjectRef
	adapter.CurserCapObjectID = testCurserCapID

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
	adapter.CCIPAddress = testCCIPPackageID
	adapter.CCIPObjectRef = testCCIPObjectRef

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
	adapter.CCIPAddress = testCCIPPackageID
	adapter.LatestCCIPPackageID = testLatestCCIPPackageID
	adapter.CCIPObjectRef = testCCIPObjectRef
	adapter.CurserCapObjectID = testCurserCapID

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
