package changesets

import (
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

func TestMCMSExecuteTransferOwnership_VerifyPreconditions_FastCurseBlocksCCIPOwnership(t *testing.T) {
	t.Parallel()

	cs := MCMSExecuteTransferOwnership{}
	err := cs.VerifyPreconditions(cldf.Environment{}, MCMSExecuteTransferOwnershipInput{
		ChainSelector: 1,
		IsFastCurse:   true,
		StateObject:   true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fastcurse MCMS cannot receive CCIP ownership transfer")
}

func TestMCMSExecuteTransferOwnership_VerifyPreconditions_FastCurseAllowsMcmsSelfTransfer(t *testing.T) {
	t.Parallel()

	cs := MCMSExecuteTransferOwnership{}
	err := cs.VerifyPreconditions(cldf.Environment{}, MCMSExecuteTransferOwnershipInput{
		ChainSelector: 1,
		IsFastCurse:   true,
		MCMS:          true,
	})
	require.NoError(t, err)
}

func TestCurseUncurseChains_VerifyPreconditions_FastUncurseBlocked(t *testing.T) {
	t.Parallel()

	cs := CurseUncurseChains{}
	err := cs.VerifyPreconditions(cldf.Environment{}, CurseUncurseChainsConfig{
		SuiChainSelector: 1,
		OperationType:    string(UncurseOperationType),
		IsGlobalCurse:    true,
		IsFastCurse:        true,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "uncurse via fastcurse MCMS is not supported")
}

func TestCurseUncurseChains_VerifyPreconditions_FastCurseAllowed(t *testing.T) {
	t.Parallel()

	cs := CurseUncurseChains{}
	err := cs.VerifyPreconditions(cldf.Environment{}, CurseUncurseChainsConfig{
		SuiChainSelector:  1,
		OperationType:     string(CurseOperationType),
		IsGlobalCurse:     true,
		IsFastCurse:       true,
		CurserCapObjectId: "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
	})
	require.NoError(t, err)
}
