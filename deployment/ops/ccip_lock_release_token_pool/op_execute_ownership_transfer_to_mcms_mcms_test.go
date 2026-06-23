package lockreleasetokenpoolops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestExecuteOwnershipTransferToMcmsLockReleaseTokenPoolOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	typeArgs := []string{mcmstest.CoinTypeArg}
	input := ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolInput{
		LockReleaseTokenPoolPackageId: mcmstest.PackageID,
		TypeArgs:                      typeArgs,
		OwnerCapObjectId:              mcmstest.OwnerCapID,
		StateObjectId:                 mcmstest.StateObjectID,
		RegistryObjectId:              mcmstest.RegistryID,
		To:                            mcmstest.Recipient,
	}
	report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolOp, sui_ops.OpTxDeps{}, input)
	require.NoError(t, err)

	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().ExecuteOwnershipTransferToMcms(typeArgs, bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.RegistryObjectId}, input.To)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
}
