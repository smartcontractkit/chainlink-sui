package deployment

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

func TestLoadOnchainStatesui_prefersDatastoreOverAddressBook(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	dsCCIP := "0xdatastore-ccip-package"
	abCCIP := "0xaddressbook-ccip-package"
	objectRef := "0xccip-object-ref"

	ds := fdatastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       dsCCIP,
		Type:          fdatastore.ContractType(SuiCCIPType),
		Version:       &Version1_0_0,
	}))
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       objectRef,
		Type:          fdatastore.ContractType(SuiCCIPObjectRefType),
		Version:       &Version1_0_0,
	}))

	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, ab.Save(selector, abCCIP, cldf.NewTypeAndVersion(SuiCCIPType, Version1_0_0)))

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore:         ds.Seal(),
		ExistingAddresses: ab,
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.Equal(t, dsCCIP, got[selector].CCIPAddress)
	require.Equal(t, objectRef, got[selector].CCIPObjectRef)
}

func TestLoadOnchainStatesui_fallsBackToAddressBookWhenDatastoreEmpty(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	abCCIP := "0xaddressbook-ccip-package"

	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, ab.Save(selector, abCCIP, cldf.NewTypeAndVersion(SuiCCIPType, Version1_0_0)))

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore:         fdatastore.NewMemoryDataStore().Seal(),
		ExistingAddresses: ab,
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.Equal(t, abCCIP, got[selector].CCIPAddress)
}

func TestLoadOnchainStatesui_fromDatastore_fastcurseLabels(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector

	ds := fdatastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xslow_package",
		Type:          fdatastore.ContractType(SuiMcmsPackageIDType),
		Version:       &Version1_0_0,
	}))
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xfast_package",
		Type:          fdatastore.ContractType(SuiMcmsPackageIDType),
		Version:       &Version1_0_0,
		Qualifier:     MCMSFastCurseLabel,
		Labels:        fdatastore.NewLabelSet(MCMSFastCurseLabel),
	}))

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore: ds.Seal(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.Equal(t, "0xslow_package", got[selector].MCMSPackageID)
	require.Equal(t, "0xfast_package", got[selector].FastCurseMCMSPackageID)
}

func TestLoadOnchainStatesui_fromDatastore_fastcurseQualifierFallback(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector

	ds := fdatastore.NewMemoryDataStore()
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xfast_package",
		Type:          fdatastore.ContractType(SuiMcmsPackageIDType),
		Version:       &Version1_0_0,
		Qualifier:     MCMSFastCurseLabel,
	}))

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore: ds.Seal(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.Equal(t, "0xfast_package", got[selector].FastCurseMCMSPackageID)
}

func TestDatastoreRefLabels_ignoresAddressTypeQualifier(t *testing.T) {
	t.Parallel()

	labels := datastoreRefLabels(fdatastore.AddressRef{
		Qualifier: "0x5ef4b483da6644c84aa78eae4f51a9bfb1fb4554d5134ac98892e931fcbdd6bf-SuiCCIP",
	})
	require.Empty(t, labels)
}

// TestLoadOnchainStatesui_fromDatastore_mcmsRoleRefCollision pins that the generic MCMS
// role refs (Canceller/BypasserManyChainMultiSig), which reuse the SuiManyChainMultisigObjectID
// address because roles are internal to the Sui MCMS state object, do not clobber the MCMS
// state entry. The loader must keep every ref per address instead of last-write-wins.
func TestLoadOnchainStatesui_fromDatastore_mcmsRoleRefCollision(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector

	ds := fdatastore.NewMemoryDataStore()
	// Slow MCMS instance: the state object id is reused by two generic role refs.
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xslow_state",
		Type:          fdatastore.ContractType(SuiMcmsObjectIDType),
		Version:       &Version1_0_0,
	}))
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xslow_state",
		Type:          fdatastore.ContractType("CancellerManyChainMultiSig"),
		Version:       &Version1_0_0,
	}))
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xslow_state",
		Type:          fdatastore.ContractType("BypasserManyChainMultiSig"),
		Version:       &Version1_0_0,
	}))
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xslow_package",
		Type:          fdatastore.ContractType(SuiMcmsPackageIDType),
		Version:       &Version1_0_0,
	}))
	// Fastcurse MCMS instance: no role refs share its state object id.
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xfast_state",
		Type:          fdatastore.ContractType(SuiMcmsObjectIDType),
		Version:       &Version1_0_0,
		Qualifier:     MCMSFastCurseLabel,
		Labels:        fdatastore.NewLabelSet(MCMSFastCurseLabel),
	}))
	require.NoError(t, ds.Addresses().Upsert(fdatastore.AddressRef{
		ChainSelector: selector,
		Address:       "0xfast_package",
		Type:          fdatastore.ContractType(SuiMcmsPackageIDType),
		Version:       &Version1_0_0,
		Qualifier:     MCMSFastCurseLabel,
		Labels:        fdatastore.NewLabelSet(MCMSFastCurseLabel),
	}))

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore: ds.Seal(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)

	chainState := got[selector]
	slowFields := chainState.MCMSState(false)
	fastFields := chainState.MCMSState(true)
	require.Equal(t, "0xslow_state", slowFields.StateObjectID,
		"slow MCMS state object must survive the role-ref collision")
	require.Equal(t, "0xslow_package", slowFields.PackageID)
	require.Equal(t, "0xfast_state", fastFields.StateObjectID)
	require.Equal(t, "0xfast_package", fastFields.PackageID)
}
