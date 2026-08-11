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

func TestCurseUncurseChains_VerifyPreconditions_FastCurseUsesRegisteredCurserCap(t *testing.T) {
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

	cs := CurseUncurseChains{}
	err := cs.VerifyPreconditions(env, CurseUncurseChainsConfig{
		SuiChainSelector: selector,
		OperationType:    string(CurseOperationType),
		IsGlobalCurse:    true,
		IsFastCurse:      true,
	})
	require.NoError(t, err)

	err = cs.VerifyPreconditions(env, CurseUncurseChainsConfig{
		SuiChainSelector:  selector,
		OperationType:     string(CurseOperationType),
		IsGlobalCurse:     true,
		IsFastCurse:       true,
		CurserCapObjectId: registeredCap,
	})
	require.NoError(t, err)

	err = cs.VerifyPreconditions(env, CurseUncurseChainsConfig{
		SuiChainSelector:  selector,
		OperationType:     string(CurseOperationType),
		IsGlobalCurse:     true,
		IsFastCurse:       true,
		CurserCapObjectId: "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "does not match registered CurserCap")
}
