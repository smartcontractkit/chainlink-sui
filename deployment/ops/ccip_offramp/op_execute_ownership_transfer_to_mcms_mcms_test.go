package offrampops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestExecuteOwnershipTransferToMcmsOffRampOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	input := ExecuteOwnershipTransferToMcmsOffRampInput{
		OffRampPackageId:     mcmstest.PackageID,
		OffRampRefObjectId:   mcmstest.StateObjectID,
		OwnerCapObjectId:     mcmstest.OwnerCapID,
		OffRampStateObjectId: mcmstest.CoinMetadata,
		RegistryObjectId:     mcmstest.RegistryID,
		To:                   mcmstest.Recipient,
	}
	report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ExecuteOwnershipTransferToMcmsOffRampOp, sui_ops.OpTxDeps{}, input)
	require.NoError(t, err)

	contract, err := module_offramp.NewOfframp(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().ExecuteOwnershipTransferToMcms(bind.Object{Id: input.OffRampRefObjectId}, bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.OffRampStateObjectId}, bind.Object{Id: input.RegistryObjectId}, input.To)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.OffRampStateObjectId, nil)
}
