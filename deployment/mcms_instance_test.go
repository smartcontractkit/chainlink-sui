package deployment

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
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
