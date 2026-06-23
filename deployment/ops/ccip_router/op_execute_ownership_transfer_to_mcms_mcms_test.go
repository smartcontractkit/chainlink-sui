package routerops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestExecuteOwnershipTransferToMcmsRouterOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	input := ExecuteOwnershipTransferToMcmsRouterInput{
		RouterPackageId:     mcmstest.PackageID,
		OwnerCapObjectId:    mcmstest.OwnerCapID,
		RouterStateObjectId: mcmstest.StateObjectID,
		RegistryObjectId:    mcmstest.RegistryID,
		To:                  mcmstest.Recipient,
	}
	report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ExecuteOwnershipTransferToMcmsRouterOp, sui_ops.OpTxDeps{}, input)
	require.NoError(t, err)

	contract, err := module_router.NewRouter(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().ExecuteOwnershipTransferToMcms(bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.RouterStateObjectId}, bind.Object{Id: input.RegistryObjectId}, input.To)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.RouterStateObjectId, nil)
}
