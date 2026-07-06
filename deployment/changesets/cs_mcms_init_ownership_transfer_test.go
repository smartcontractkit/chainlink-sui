package changesets

import (
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

func TestInitMCMSOwnershipTransfer_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	cs := InitMCMSOwnershipTransfer{}

	err := cs.VerifyPreconditions(cldf.Environment{}, InitMCMSOwnershipTransferConfig{
		ChainSelector: 17529533435026248318,
	})
	require.NoError(t, err)

	err = cs.VerifyPreconditions(cldf.Environment{}, InitMCMSOwnershipTransferConfig{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "chainSelector is required")
}
