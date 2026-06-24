package onrampops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestExecuteOwnershipTransferToMcmsOnRampOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	input := ExecuteOwnershipTransferToMcmsOnRampInput{
		OnRampPackageId:     mcmstest.PackageID,
		OnRampRefObjectId:   mcmstest.StateObjectID,
		OwnerCapObjectId:    mcmstest.OwnerCapID,
		OnRampStateObjectId: mcmstest.CoinMetadata,
		RegistryObjectId:    mcmstest.RegistryID,
		To:                  mcmstest.Recipient,
	}
	report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ExecuteOwnershipTransferToMcmsOnRampOp, sui_ops.OpTxDeps{}, input)
	require.NoError(t, err)

	contract, err := module_onramp.NewOnramp(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().ExecuteOwnershipTransferToMcms(bind.Object{Id: input.OnRampRefObjectId}, bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.OnRampStateObjectId}, bind.Object{Id: input.RegistryObjectId}, input.To)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.OnRampStateObjectId, nil)
}
