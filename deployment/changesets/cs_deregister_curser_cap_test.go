package changesets

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func TestDeregisterCurserCap_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	const capID = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"

	cs := DeregisterCurserCap{}

	err := cs.VerifyPreconditions(cldf.Environment{}, DeregisterCurserCapConfig{
		SuiChainSelector:   selector,
		CurserCapObjectIds: []string{capID},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "timelockConfig is required")

	err = cs.VerifyPreconditions(cldf.Environment{}, DeregisterCurserCapConfig{
		SuiChainSelector: selector,
		TimelockConfig:   &utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "curserCapObjectIds must not be empty")

	err = cs.VerifyPreconditions(cldf.Environment{}, DeregisterCurserCapConfig{
		SuiChainSelector:   selector,
		CurserCapObjectIds: []string{capID},
		TimelockConfig:     &utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no Sui chain state")
}

func TestDeregisterCurserCap_VerifyPreconditions_RequiresSlowMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	const capID = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"

	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, ab.Save(selector, "0xccip_pkg", cldf.NewTypeAndVersion(deployment.SuiCCIPType, deployment.Version1_0_0)))
	require.NoError(t, ab.Save(selector, "0xccip_ref", cldf.NewTypeAndVersion(deployment.SuiCCIPObjectRefType, deployment.Version1_0_0)))
	require.NoError(t, ab.Save(selector, "0xowner_cap", cldf.NewTypeAndVersion(deployment.SuiCCIPOwnerCapObjectIDType, deployment.Version1_0_0)))

	cs := DeregisterCurserCap{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), DeregisterCurserCapConfig{
		SuiChainSelector:   selector,
		CurserCapObjectIds: []string{capID},
		TimelockConfig:     &utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "slow MCMS must be deployed")
}

func TestDeregisterCurserCap_VerifyPreconditions_SucceedsWithSlowMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	const capID = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"

	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, ab.Save(selector, "0xccip_pkg", cldf.NewTypeAndVersion(deployment.SuiCCIPType, deployment.Version1_0_0)))
	require.NoError(t, ab.Save(selector, "0xccip_ref", cldf.NewTypeAndVersion(deployment.SuiCCIPObjectRefType, deployment.Version1_0_0)))
	require.NoError(t, ab.Save(selector, "0xowner_cap", cldf.NewTypeAndVersion(deployment.SuiCCIPOwnerCapObjectIDType, deployment.Version1_0_0)))
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, fdatastore.NewMemoryDataStore().Addresses(), selector, mcmsops.DeployMCMSSeqOutput{
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

	cs := DeregisterCurserCap{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), DeregisterCurserCapConfig{
		SuiChainSelector:   selector,
		CurserCapObjectIds: []string{capID},
		TimelockConfig:     &utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	require.NoError(t, err)
}
