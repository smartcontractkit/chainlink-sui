package deployment

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

func TestTokenQualifier_usesSymbolFormNotDisplayName(t *testing.T) {
	t.Parallel()

	// Sui labels carry display names with spaces; the datastore key uses the symbol form so
	// Sui rows match the EVM and Solana rows for the same token.
	require.Equal(t, "CCIP-BnM", TokenQualifier("CCIP BnM"))
	require.Equal(t, "CCIP-LnR", TokenQualifier("CCIP LnR"))
	require.Equal(t, "LINK", TokenQualifier("LINK"))
	require.Equal(t, "USDC", TokenQualifier("USDC"))
}

func TestMinterCapQualifier_keepsHolderAddressAsPassed(t *testing.T) {
	t.Parallel()
	require.Equal(t, "CCIP-BnM-0xab", MinterCapQualifier("CCIP BnM", "0xab"))
	require.Equal(t, "CCIP-BnM-0xAB", MinterCapQualifier("CCIP BnM", "0xAB"))
}

func TestSaveSuiAddress_rejectsAddressDerivedQualifier(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	address := "0xAbC123"
	tv := cldf.NewTypeAndVersion(SuiCCIPType, Version1_0_0)

	// The shim the convention exists to eliminate: a qualifier you can only supply if you
	// already know the address you are looking up.
	for _, qualifier := range []string{
		address + "-" + string(SuiCCIPType),
		"0xabc123-SuiCCIP", // the detector is case-insensitive
		address,
	} {
		ab := cldf.NewMemoryAddressBook()
		ds := fdatastore.NewMemoryDataStore()

		err := SaveSuiAddress(ab, ds.Addresses(), selector, address, tv, qualifier)
		require.ErrorContains(t, err, "contains the address being written")

		// Rejected before either store is touched, so the two cannot drift.
		require.Empty(t, ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector)))
		_, abErr := ab.AddressesForChain(selector)
		require.Error(t, abErr)
	}
}

// TestMinterCapQualifiers_noCollisionAcrossWriters reproduces the shape that exists in
// testnet today: one chain, one token, several minter caps written by different changesets.
// Under a bare token-symbol qualifier all of these share a datastore key and Upsert keeps
// only the last, so the loader reports fewer caps than were actually granted.
func TestMinterCapQualifiers_noCollisionAcrossWriters(t *testing.T) {
	t.Parallel()

	const (
		symbol   = "CCIP BnM"
		deployer = "0xdeployer"
		minter   = "0xthird-party-minter"
	)
	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()

	save := func(t *testing.T, contractType cldf.ContractType, address, qualifier string) {
		t.Helper()
		tv := cldf.NewTypeAndVersion(contractType, Version1_0_0)
		tv.AddLabel(symbol)
		require.NoError(t, SaveSuiAddress(ab, ds.Addresses(), selector, address, tv, qualifier))
	}

	// DeployManagedToken: the cap granted to the token's initial minter.
	save(t, SuiManagedTokenMinterCapID, "0xcap-deployer", MinterCapQualifier(symbol, deployer))
	// DeployManagedTokenFaucet: the cap the faucet mints with. Held by the deployer too, so
	// it is separated by type rather than by holder.
	save(t, SuiManagedTokenFaucetMinterCapIDType, "0xcap-faucet", TokenQualifier(symbol))
	// ManagedTokenConfigureNewMinter: a cap per additional authorised minter.
	save(t, SuiManagedTokenMinterCapID, "0xcap-third-party", MinterCapQualifier(symbol, minter))

	refs := ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector))
	require.Len(t, refs, 3, "each cap must occupy its own datastore key")

	got, err := LoadOnchainStatesui(cldf.Environment{
		DataStore: ds.Seal(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	})
	require.NoError(t, err)
	require.ElementsMatch(t,
		[]string{"0xcap-deployer", "0xcap-third-party"},
		got[selector].ManagedTokens[symbol].MinterCapObjectIds,
	)
	require.Equal(t, "0xcap-faucet", got[selector].ManagedTokenFaucets[symbol].MinterCapObjectID)
}

// Writing one key twice inside a single changeset run has no legitimate reading — both writes
// belong to the same deployment, so one of the two objects would go unrecorded. It must fail
// rather than let the second write displace the first.
func TestSaveSuiAddress_rejectsDuplicateKeyWithinOneChangeset(t *testing.T) {
	t.Parallel()

	const symbol = "USDC"
	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()

	tv := cldf.NewTypeAndVersion(SuiManagedTokenMinterCapID, Version1_0_0)
	tv.AddLabel(symbol)
	qualifier := MinterCapQualifier(symbol, "0xminter")

	require.NoError(t, SaveSuiAddress(ab, ds.Addresses(), selector, "0xcap-v1", tv, qualifier))

	err := SaveSuiAddress(ab, ds.Addresses(), selector, "0xcap-v2", tv, qualifier)
	require.ErrorContains(t, err, "already wrote a SuiManagedTokenMinterCapID")
	require.ErrorContains(t, err, "needs a distinct qualifier")

	// The first ref survives: a rejected write does not disturb what is already recorded.
	refs := ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector))
	require.Len(t, refs, 1)
	require.Equal(t, "0xcap-v1", refs[0].Address)
}
