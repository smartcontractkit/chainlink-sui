package deployment

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"

	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

func TestLoadOnchainStatesui_DualMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector

	slow := map[string]cldf.ContractType{
		"0xslow_package":  SuiMcmsPackageIDType,
		"0xslow_state":    SuiMcmsObjectIDType,
		"0xslow_registry": SuiMcmsRegistryObjectIDType,
		"0xslow_account":  SuiMcmsAccountStateObjectIDType,
		"0xslow_timelock": SuiMcmsTimelockObjectIDType,
		"0xslow_deployer": SuiMcmsDeployerObjectIDType,
	}
	fast := map[string]cldf.ContractType{
		"0xfast_package":  SuiMcmsPackageIDType,
		"0xfast_state":    SuiMcmsObjectIDType,
		"0xfast_registry": SuiMcmsRegistryObjectIDType,
		"0xfast_account":  SuiMcmsAccountStateObjectIDType,
		"0xfast_timelock": SuiMcmsTimelockObjectIDType,
		"0xfast_deployer": SuiMcmsDeployerObjectIDType,
	}

	ab := cldf.NewMemoryAddressBook()
	for addr, typ := range slow {
		require.NoError(t, ab.Save(selector, addr, cldf.NewTypeAndVersion(typ, Version1_0_0)))
	}
	for addr, typ := range fast {
		tv := cldf.NewTypeAndVersion(typ, Version1_0_0)
		tv.Labels.Add(MCMSFastCurseLabel)
		require.NoError(t, ab.Save(selector, addr, tv))
	}

	got, err := LoadOnchainStatesui(cldf.Environment{
		ExistingAddresses: ab,
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)

	chainState := got[selector]
	require.True(t, chainState.HasMCMSInstance(MCMSInstanceSlow))
	require.True(t, chainState.HasMCMSInstance(MCMSInstanceFastCurse))

	require.Equal(t, "0xslow_registry", chainState.MCMSStateByInstance(MCMSInstanceSlow).RegistryObjectID)
	require.Equal(t, "0xfast_registry", chainState.MCMSStateByInstance(MCMSInstanceFastCurse).RegistryObjectID)
	require.NotEqual(t,
		chainState.MCMSStateByInstance(MCMSInstanceSlow).StateObjectID,
		chainState.MCMSStateByInstance(MCMSInstanceFastCurse).StateObjectID,
	)
}

func TestMCMSInstanceFromFastCurseFlag(t *testing.T) {
	t.Parallel()

	require.Equal(t, MCMSInstanceSlow, MCMSInstanceFromFastCurseFlag(false))
	require.Equal(t, MCMSInstanceFastCurse, MCMSInstanceFromFastCurseFlag(true))
	require.Equal(t, MCMSFastCurseLabel, MCMSInstanceFastCurse.AddressBookLabel())
	require.Empty(t, MCMSInstanceSlow.AddressBookLabel())
}

// DatastoreQualifier and MCMSInstanceFromQualifier must stay inverses: the loader recovers
// the instance from the qualifier when a ref arrives without labels, so a drift between the
// two would file a fastcurse deployment as the slow one.
func TestMCMSInstanceQualifierRoundTrip(t *testing.T) {
	t.Parallel()

	for _, instance := range []MCMSInstance{MCMSInstanceSlow, MCMSInstanceFastCurse} {
		got, ok := MCMSInstanceFromQualifier(instance.DatastoreQualifier())
		require.True(t, ok)
		require.Equal(t, instance, got)
	}

	// Rows written before the purpose qualifiers existed.
	got, ok := MCMSInstanceFromQualifier(MCMSFastCurseLabel)
	require.True(t, ok)
	require.Equal(t, MCMSInstanceFastCurse, got)

	for _, qualifier := range []string{"", "USDC", "0xabc-SuiMcmsPackageID"} {
		_, ok := MCMSInstanceFromQualifier(qualifier)
		require.False(t, ok, "qualifier %q names no MCMS instance", qualifier)
	}
}

// StoreMCMSInAddressBook and PlannedMCMSRefs share one type mapping (mcmsObjectTypes): the
// pre-deploy conflict check must cover exactly the keys the write will occupy. This pins the
// correspondence so a future MCMS object type cannot be added to one without the other.
func TestPlannedMCMSRefs_matchesStoreMCMSInAddressBook(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector

	for _, instance := range []MCMSInstance{MCMSInstanceSlow, MCMSInstanceFastCurse} {
		ab := cldf.NewMemoryAddressBook()
		ds := fdatastore.NewMemoryDataStore()

		report := mcmsops.DeployMCMSSeqOutput{
			PackageId: "0xmcms-pkg",
			Objects: mcmsops.DeployMCMSObjects{
				McmsMultisigStateObjectId:   "0xmcms-state",
				TimelockObjectId:            "0xmcms-timelock",
				McmsDeployerStateObjectId:   "0xmcms-deployer",
				McmsRegistryObjectId:        "0xmcms-registry",
				McmsAccountStateObjectId:    "0xmcms-account",
				McmsAccountOwnerCapObjectId: "0xmcms-ownercap",
			},
		}
		require.NoError(t, StoreMCMSInAddressBook(ab, ds.Addresses(), selector, report, instance))

		planned := PlannedMCMSRefs(instance)
		refs := ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector))
		require.Len(t, refs, len(planned))

		plannedKeys := make(map[[2]string]struct{}, len(planned))
		for _, p := range planned {
			plannedKeys[[2]string{string(p.Type), p.Qualifier}] = struct{}{}
		}
		for _, ref := range refs {
			_, ok := plannedKeys[[2]string{string(ref.Type), ref.Qualifier}]
			require.True(t, ok, "wrote %s qualified %q with no planned key", ref.Type, ref.Qualifier)
			require.Equal(t, instance.DatastoreQualifier(), ref.Qualifier)
		}
	}
}
