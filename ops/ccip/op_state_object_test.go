//go:build integration

package ccipops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/ops/link"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"

	"github.com/stretchr/testify/require"
)

func TestStateObjectPackageIdOperations(t *testing.T) {
	t.Parallel()

	signer, client := testenv.SetupEnvironment(t)

	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)

	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	// Deploy MCMS
	mcmsReport, err := cld_ops.ExecuteOperation(bundle, mcmsops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Contract")

	signerAddress, err := signer.GetAddress()
	require.NoError(t, err, "failed to get signer address")

	// Deploy CCIP
	ccipReport, err := cld_ops.ExecuteOperation(bundle, DeployCCIPOp, deps, DeployCCIPInput{
		McmsPackageId: mcmsReport.Output.PackageId,
		McmsOwner:     signerAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Package")

	t.Run("Test Get Initial Package ID", func(t *testing.T) {
		// Test getting initial package ID
		getInitialReport, err := cld_ops.ExecuteOperation(bundle, GetInitialPackageIdStateObjectOp, deps, GetInitialPackageIdStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
		})
		require.NoError(t, err, "failed to get initial package ID")
		require.NotEmpty(t, getInitialReport.Output.Objects.InitialPackageId, "initial package ID should not be empty")
		require.Equal(t, ccipReport.Output.PackageId, getInitialReport.Output.Objects.InitialPackageId, "initial package ID should match deployed package ID")
	})

	t.Run("Test Get Package IDs", func(t *testing.T) {
		// Test getting all package IDs
		getPackageIdsReport, err := cld_ops.ExecuteOperation(bundle, GetPackageIdsStateObjectOp, deps, GetPackageIdsStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
		})
		require.NoError(t, err, "failed to get package IDs")
		require.NotEmpty(t, getPackageIdsReport.Output.Objects.PackageIds, "package IDs should not be empty")
		require.Contains(t, getPackageIdsReport.Output.Objects.PackageIds, ccipReport.Output.PackageId, "package IDs should contain deployed package ID")
	})

	t.Run("Test Add Package ID", func(t *testing.T) {
		// Add a new package ID
		newPackageId := "0x1234567890abcdef1234567890abcdef12345678"
		_, err := cld_ops.ExecuteOperation(bundle, AddPackageIdStateObjectOp, deps, AddPackageIdStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId:      ccipReport.Output.Objects.OwnerCapObjectId,
			PackageId:             newPackageId,
		})
		require.NoError(t, err, "failed to add package ID")

		// Verify the package ID was added
		getPackageIdsReport, err := cld_ops.ExecuteOperation(bundle, GetPackageIdsStateObjectOp, deps, GetPackageIdsStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
		})
		require.NoError(t, err, "failed to get package IDs after adding")
		require.Contains(t, getPackageIdsReport.Output.Objects.PackageIds, newPackageId, "package IDs should contain the newly added package ID")
		require.Len(t, getPackageIdsReport.Output.Objects.PackageIds, 2, "should have 2 package IDs now")
	})
}
