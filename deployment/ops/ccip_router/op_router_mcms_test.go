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

func TestRouterDualModeOps_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	router, err := module_router.NewRouter(mcmstest.PackageID, nil)
	require.NoError(t, err)

	t.Run("set_on_ramps", func(t *testing.T) {
		t.Parallel()
		input := SetOnRampsInput{
			RouterPackageId:     mcmstest.PackageID,
			RouterStateObjectId: mcmstest.StateObjectID,
			OwnerCapObjectId:    mcmstest.OwnerCapID,
			DestChainSelectors:  []uint64{mcmstest.DestChainSel},
			OnRampAddresses:     []string{mcmstest.Recipient},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), SetOnRampsOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := router.Encoder().SetOnRamps(bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.RouterStateObjectId}, input.DestChainSelectors, input.OnRampAddresses)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.RouterStateObjectId, nil)
	})

	t.Run("accept_ownership", func(t *testing.T) {
		t.Parallel()
		input := AcceptOwnershipInput{
			RouterPackageId:     mcmstest.PackageID,
			RouterStateObjectId: mcmstest.StateObjectID,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), AcceptOwnershipOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := router.Encoder().AcceptOwnership(bind.Object{Id: input.RouterStateObjectId})
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.RouterStateObjectId, nil)
	})
}
