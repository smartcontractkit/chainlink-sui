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

func TestSaveSuiAddress_dualWrite(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()

	tv := cldf.NewTypeAndVersion(SuiCCIPType, Version1_0_0)
	require.NoError(t, SaveSuiAddress(ab, ds.Addresses(), selector, "0xccip-pkg", tv, ChainSingletonQualifier))

	abAddrs, err := ab.AddressesForChain(selector)
	require.NoError(t, err)
	require.Contains(t, abAddrs, "0xccip-pkg")
	require.Equal(t, SuiCCIPType, abAddrs["0xccip-pkg"].Type)

	refs := ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector))
	require.Len(t, refs, 1)
	require.Equal(t, "0xccip-pkg", refs[0].Address)
	require.Equal(t, fdatastore.ContractType(SuiCCIPType), refs[0].Type)
	require.Equal(t, "1.0.0", refs[0].Version.String())
	require.Empty(t, refs[0].Qualifier)
	require.True(t, refs[0].Labels.IsEmpty())
}

func TestSaveSuiAddress_symbolLabelRoundTrip(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()

	tv := cldf.NewTypeAndVersion(SuiManagedTokenPackageIDType, Version1_0_0)
	tv.AddLabel("USDC")
	require.NoError(t, SaveSuiAddress(ab, ds.Addresses(), selector, "0xusdc-pkg", tv, "USDC"))

	refs := ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector))
	require.Len(t, refs, 1)
	require.True(t, refs[0].Labels.Contains("USDC"))

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore: ds.Seal(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.Contains(t, got[selector].ManagedTokens, "USDC")
	require.Equal(t, "0xusdc-pkg", got[selector].ManagedTokens["USDC"].PackageID)
}

func TestSaveSuiAddress_fastcurseLabelRoundTrip(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()

	tv := cldf.NewTypeAndVersion(SuiMcmsPackageIDType, Version1_0_0)
	tv.Labels.Add(MCMSFastCurseLabel)
	require.NoError(t, SaveSuiAddress(ab, ds.Addresses(), selector, "0xfast-pkg", tv, RMNMCMSQualifier))

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore: ds.Seal(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.Equal(t, "0xfast-pkg", got[selector].FastCurseMCMSPackageID)
}

func TestStoreMCMSInAddressBook_dualWrite(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
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
	require.NoError(t, StoreMCMSInAddressBook(ab, ds.Addresses(), selector, report, MCMSInstanceFastCurse))

	refs := ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector))
	require.Len(t, refs, 7)
	for _, ref := range refs {
		require.True(t, ref.Labels.Contains(MCMSFastCurseLabel))
	}

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore: ds.Seal(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.Equal(t, report.PackageId, got[selector].FastCurseMCMSPackageID)
	require.Equal(t, report.Objects.McmsRegistryObjectId, got[selector].FastCurseMCMSRegistryObjectID)
}
