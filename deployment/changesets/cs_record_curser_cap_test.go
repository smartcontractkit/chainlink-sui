package changesets

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment"
)

func TestRecordCurserCap_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	cs := RecordCurserCap{}

	err := cs.VerifyPreconditions(cldf.Environment{}, RecordCurserCapConfig{
		SuiChainSelector: cselectors.SUI_TESTNET.Selector,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "txDigest or curserCapObjectId")

	err = cs.VerifyPreconditions(cldf.Environment{}, RecordCurserCapConfig{
		SuiChainSelector: cselectors.SUI_TESTNET.Selector,
		TxDigest:         "0xabc",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no Sui chain client")

	err = cs.VerifyPreconditions(cldf.Environment{
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			cselectors.SUI_TESTNET.Selector: sui.Chain{},
		}),
	}, RecordCurserCapConfig{
		SuiChainSelector:  cselectors.SUI_TESTNET.Selector,
		CurserCapObjectId: "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
	})
	require.NoError(t, err)
}

func TestRecordCurserCap_Apply_SavesCurserCapObjectID(t *testing.T) {
	t.Parallel()

	const capID = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	selector := cselectors.SUI_TESTNET.Selector

	env := cldf.Environment{
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	}

	out, err := RecordCurserCap{}.Apply(env, RecordCurserCapConfig{
		SuiChainSelector:  selector,
		CurserCapObjectId: capID,
	})
	require.NoError(t, err)
	require.NotNil(t, out.AddressBook)

	addrs, err := out.AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	require.Contains(t, addrs, capID)
	require.Equal(t, deployment.SuiCurserCapObjectIDType, addrs[capID].Type)
}

func TestRecordCurserCap_Apply_RejectsConflictingCurserCap(t *testing.T) {
	t.Parallel()

	const registeredCap = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	const otherCap = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	selector := cselectors.SUI_TESTNET.Selector

	env := cldf.Environment{
		ExistingAddresses: cldf.NewMemoryAddressBookFromMap(map[uint64]map[string]cldf.TypeAndVersion{
			selector: {
				registeredCap: cldf.NewTypeAndVersion(deployment.SuiCurserCapObjectIDType, deployment.Version1_0_0),
			},
		}),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	}

	_, err := RecordCurserCap{}.Apply(env, RecordCurserCapConfig{
		SuiChainSelector:  selector,
		CurserCapObjectId: otherCap,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "conflicts with registered CurserCap")
}

func TestRecordCurserCap_Apply_IdempotentWhenCapAlreadyRegistered(t *testing.T) {
	t.Parallel()

	const registeredCap = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	selector := cselectors.SUI_TESTNET.Selector

	env := cldf.Environment{
		ExistingAddresses: cldf.NewMemoryAddressBookFromMap(map[uint64]map[string]cldf.TypeAndVersion{
			selector: {
				registeredCap: cldf.NewTypeAndVersion(deployment.SuiCurserCapObjectIDType, deployment.Version1_0_0),
			},
		}),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	}

	out, err := RecordCurserCap{}.Apply(env, RecordCurserCapConfig{
		SuiChainSelector:  selector,
		CurserCapObjectId: registeredCap,
	})
	require.NoError(t, err)

	addrs, err := out.AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	require.Contains(t, addrs, registeredCap)
}
