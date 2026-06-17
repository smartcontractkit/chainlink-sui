package usdctokenpoolops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_usdc_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/usdc_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestUsdcTokenPoolDualModeOps_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	typeArgs := []string{mcmstest.CoinTypeArg}
	contract, err := module_usdc_token_pool.NewUsdcTokenPool(mcmstest.PackageID, nil)
	require.NoError(t, err)

	t.Run("accept_ownership", func(t *testing.T) {
		t.Parallel()
		input := AcceptOwnershipUsdcTokenPoolInput{
			UsdcTokenPoolPackageId: mcmstest.PackageID,
			TypeArgs:               typeArgs,
			StateObjectId:          mcmstest.StateObjectID,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), AcceptOwnershipUsdcTokenPoolOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().AcceptOwnership(typeArgs, bind.Object{Id: input.StateObjectId})
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})

	t.Run("execute_ownership_transfer_to_mcms", func(t *testing.T) {
		t.Parallel()
		input := ExecuteOwnershipTransferToMcmsUsdcTokenPoolInput{
			UsdcTokenPoolPackageId: mcmstest.PackageID,
			TypeArgs:               typeArgs,
			OwnerCapObjectId:       mcmstest.OwnerCapID,
			StateObjectId:          mcmstest.StateObjectID,
			RegistryObjectId:       mcmstest.RegistryID,
			To:                     mcmstest.Recipient,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ExecuteOwnershipTransferToMcmsUsdcTokenPoolOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().ExecuteOwnershipTransferToMcms(typeArgs, bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.RegistryObjectId}, input.To)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})
}
