package changesets

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func TestRegisterCurserCap_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	cs := RegisterCurserCap{}

	err := cs.VerifyPreconditions(cldf.Environment{}, RegisterCurserCapConfig{
		SuiChainSelector: cselectors.SUI_TESTNET.Selector,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "timelockConfig is required")

	err = cs.VerifyPreconditions(cldf.Environment{}, RegisterCurserCapConfig{
		SuiChainSelector: cselectors.SUI_TESTNET.Selector,
		TimelockConfig:   &utils.TimelockConfig{MCMSAction: types.TimelockActionBypass},
	})
	require.NoError(t, err)
}
