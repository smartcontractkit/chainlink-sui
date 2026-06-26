package changesets

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func dualMCMSEnv(t *testing.T, ab cldf.AddressBook, selector uint64) cldf.Environment {
	t.Helper()
	return cldf.Environment{
		ExistingAddresses: ab,
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	}
}

func TestDeployMCMS_VerifyPreconditions_RejectsDuplicateInstance(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
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

	cs := DeployMCMS{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), DeployMCMSConfig{
		DeployMCMSSeqInput: mcmsops.DeployMCMSSeqInput{ChainSelector: selector},
		IsFastCurse:        false,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "slow MCMS is already recorded")
}

func TestDeployMCMS_VerifyPreconditions_RejectsDuplicateFastInstance(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
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

	cs := DeployMCMS{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), DeployMCMSConfig{
		DeployMCMSSeqInput: mcmsops.DeployMCMSSeqInput{ChainSelector: selector},
		IsFastCurse:        true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fastcurse MCMS is already recorded")
}

func TestDeployMCMS_VerifyPreconditions_AllowsFastWhenOnlySlowExists(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
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

	cs := DeployMCMS{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), DeployMCMSConfig{
		DeployMCMSSeqInput: mcmsops.DeployMCMSSeqInput{ChainSelector: selector},
		IsFastCurse:        true,
	})
	require.NoError(t, err)
}

func TestRegisterCurserCap_VerifyPreconditions_RequiresBothMCMSInstances(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
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

	cs := RegisterCurserCap{}
	cfg := RegisterCurserCapConfig{
		SuiChainSelector: selector,
		TimelockConfig:   &utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	}

	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "fastcurse MCMS must be deployed")

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

	require.NoError(t, cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), cfg))
}
