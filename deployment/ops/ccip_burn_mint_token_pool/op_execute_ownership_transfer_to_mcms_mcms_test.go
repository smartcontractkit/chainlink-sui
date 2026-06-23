package burnminttokenpoolops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestExecuteOwnershipTransferToMcmsBurnMintTokenPoolOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	typeArgs := []string{mcmstest.CoinTypeArg}
	input := ExecuteOwnershipTransferToMcmsBurnMintTokenPoolInput{
		BurnMintTokenPoolPackageId: mcmstest.PackageID,
		TypeArgs:                   typeArgs,
		OwnerCapObjectId:           mcmstest.OwnerCapID,
		StateObjectId:              mcmstest.StateObjectID,
		RegistryObjectId:           mcmstest.RegistryID,
		To:                         mcmstest.Recipient,
	}
	report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ExecuteOwnershipTransferToMcmsBurnMintTokenPoolOp, sui_ops.OpTxDeps{}, input)
	require.NoError(t, err)

	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().ExecuteOwnershipTransferToMcms(typeArgs, bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.RegistryObjectId}, input.To)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
}
