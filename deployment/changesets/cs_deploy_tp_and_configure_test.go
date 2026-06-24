package changesets

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

type stubSigner struct{}

func (stubSigner) Sign([]byte) ([]string, error) { return nil, nil }
func (stubSigner) GetAddress() (string, error)   { return "0x1", nil }

func deployTPEnv(t *testing.T, ab *cldf.AddressBookMap, selector uint64) cldf.Environment {
	t.Helper()
	return cldf.Environment{
		ExistingAddresses: ab,
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{Signer: stubSigner{}},
		}),
	}
}

func storeSlowMCMS(t *testing.T, ab *cldf.AddressBookMap, selector uint64) {
	t.Helper()
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xslow_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xslow_state",
			McmsRegistryObjectId:        "0xslow_registry",
			McmsAccountStateObjectId:    "0xslow_account",
			McmsAccountOwnerCapObjectId: "0xslow_owner_cap",
			TimelockObjectId:            "0xslow_timelock",
			McmsDeployerStateObjectId:   "0xslow_deployer",
		},
	}, deployment.MCMSInstanceSlow))
}

func storeFastMCMS(t *testing.T, ab *cldf.AddressBookMap, selector uint64) {
	t.Helper()
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xfast_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xfast_state",
			McmsRegistryObjectId:        "0xfast_registry",
			McmsAccountStateObjectId:    "0xfast_account",
			McmsAccountOwnerCapObjectId: "0xfast_owner_cap",
			TimelockObjectId:            "0xfast_timelock",
			McmsDeployerStateObjectId:   "0xfast_deployer",
		},
	}, deployment.MCMSInstanceFastCurse))
}

func TestDeployTPAndConfigure_VerifyPreconditions_RequiresFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	storeSlowMCMS(t, ab, selector)

	cs := DeployTPAndConfigure{}
	err := cs.VerifyPreconditions(deployTPEnv(t, ab, selector), DeployTPAndConfigureConfig{
		SuiChainSelector: selector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeManaged},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fast MCMS package not deployed")
}

func TestDeployTPAndConfigure_VerifyPreconditions_SucceedsWithFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	storeSlowMCMS(t, ab, selector)
	storeFastMCMS(t, ab, selector)

	cs := DeployTPAndConfigure{}
	err := cs.VerifyPreconditions(deployTPEnv(t, ab, selector), DeployTPAndConfigureConfig{
		SuiChainSelector: selector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeBurnMint},
	})
	require.NoError(t, err)
}

func TestDeployTPAndConfigure_Apply_RequiresFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	storeSlowMCMS(t, ab, selector)

	_, err := DeployTPAndConfigure{}.Apply(deployTPEnv(t, ab, selector), DeployTPAndConfigureConfig{
		SuiChainSelector: selector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeLockRelease},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fast MCMS package not deployed")
}
