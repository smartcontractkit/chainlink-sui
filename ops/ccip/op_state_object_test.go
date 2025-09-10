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

func TestStateObjectOperations(t *testing.T) {
	t.Parallel()

	signer, client := testenv.SetupEnvironment(t)

	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(500_000_000)
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
	_, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
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

	t.Run("Test Get Owner", func(t *testing.T) {
		// Test getting owner
		getOwnerReport, err := cld_ops.ExecuteOperation(bundle, GetOwnerStateObjectOp, deps, GetOwnerStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
		})
		require.NoError(t, err, "failed to get owner")
		require.NotEmpty(t, getOwnerReport.Output.Objects.Owner, "owner should not be empty")
		require.Equal(t, signerAddress, getOwnerReport.Output.Objects.Owner, "owner should match signer address")
	})

	t.Run("Test Get Owner Cap ID", func(t *testing.T) {
		// Test getting owner cap ID
		getOwnerCapIdReport, err := cld_ops.ExecuteOperation(bundle, GetOwnerCapIdStateObjectOp, deps, GetOwnerCapIdStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
		})
		require.NoError(t, err, "failed to get owner cap ID")
		require.NotEmpty(t, getOwnerCapIdReport.Output.Objects.OwnerCapId, "owner cap ID should not be empty")
		require.Equal(t, ccipReport.Output.Objects.OwnerCapObjectId, getOwnerCapIdReport.Output.Objects.OwnerCapId, "owner cap ID should match deployed owner cap ID")
	})

	t.Run("Test Get Pending Transfer", func(t *testing.T) {
		// Test getting pending transfer info (should be empty initially)
		getPendingTransferReport, err := cld_ops.ExecuteOperation(bundle, GetPendingTransferStateObjectOp, deps, GetPendingTransferStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
		})
		require.NoError(t, err, "failed to get pending transfer info")
		require.False(t, getPendingTransferReport.Output.Objects.HasPendingTransfer, "should not have pending transfer initially")
		require.Nil(t, getPendingTransferReport.Output.Objects.PendingTransferFrom, "pending transfer from should be nil")
		require.Nil(t, getPendingTransferReport.Output.Objects.PendingTransferTo, "pending transfer to should be nil")
		require.Nil(t, getPendingTransferReport.Output.Objects.PendingTransferAccepted, "pending transfer accepted should be nil")
	})

	t.Run("Test Add Package ID", func(t *testing.T) {
		// Add a new package ID
		newPackageId := "0x1234567890abcdef1234567890abcdef12345678"
		addReport, err := cld_ops.ExecuteOperation(bundle, AddPackageIdStateObjectOp, deps, AddPackageIdStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId:      ccipReport.Output.Objects.OwnerCapObjectId,
			PackageId:             newPackageId,
		})
		require.NoError(t, err, "failed to add package ID")
		require.NotEmpty(t, addReport.Output.Digest, "add package ID transaction should have a digest")
	})

	t.Run("Test Remove Package ID", func(t *testing.T) {
		// First add a package ID to remove
		newPackageId := "0xabcdef1234567890abcdef1234567890abcdef12"
		_, err := cld_ops.ExecuteOperation(bundle, AddPackageIdStateObjectOp, deps, AddPackageIdStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId:      ccipReport.Output.Objects.OwnerCapObjectId,
			PackageId:             newPackageId,
		})
		require.NoError(t, err, "failed to add package ID for removal test")

		// Now remove the package ID
		removeReport, err := cld_ops.ExecuteOperation(bundle, RemovePackageIdStateObjectOp, deps, RemovePackageIdStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId:      ccipReport.Output.Objects.OwnerCapObjectId,
			PackageId:             newPackageId,
		})
		require.NoError(t, err, "failed to remove package ID")
		require.NotEmpty(t, removeReport.Output.Digest, "remove package ID transaction should have a digest")
	})

	t.Run("Test Ownership Transfer", func(t *testing.T) {
		// Create a new signer for the transfer
		newSigner, _ := testenv.SetupEnvironment(t)
		newSignerAddress, err := newSigner.GetAddress()
		require.NoError(t, err, "failed to get new signer address")

		// Transfer ownership
		transferReport, err := cld_ops.ExecuteOperation(bundle, TransferOwnershipStateObjectOp, deps, TransferOwnershipStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
			OwnerCapObjectId:      ccipReport.Output.Objects.OwnerCapObjectId,
			To:                    newSignerAddress,
		})
		require.NoError(t, err, "failed to transfer ownership")
		require.NotEmpty(t, transferReport.Output.Digest, "transfer ownership transaction should have a digest")

		// Verify pending transfer exists
		getPendingTransferReport, err := cld_ops.ExecuteOperation(bundle, GetPendingTransferStateObjectOp, deps, GetPendingTransferStateObjectInput{
			CCIPPackageId:         ccipReport.Output.PackageId,
			CCIPObjectRefObjectId: ccipReport.Output.Objects.CCIPObjectRefObjectId,
		})
		require.NoError(t, err, "failed to get pending transfer info after transfer")
		require.True(t, getPendingTransferReport.Output.Objects.HasPendingTransfer, "should have pending transfer")
		require.NotNil(t, getPendingTransferReport.Output.Objects.PendingTransferFrom, "pending transfer from should not be nil")
		require.NotNil(t, getPendingTransferReport.Output.Objects.PendingTransferTo, "pending transfer to should not be nil")
		require.Equal(t, signerAddress, *getPendingTransferReport.Output.Objects.PendingTransferFrom, "pending transfer from should match original owner")
		require.Equal(t, newSignerAddress, *getPendingTransferReport.Output.Objects.PendingTransferTo, "pending transfer to should match new owner")
		require.NotNil(t, getPendingTransferReport.Output.Objects.PendingTransferAccepted, "pending transfer accepted should not be nil")
		require.False(t, *getPendingTransferReport.Output.Objects.PendingTransferAccepted, "pending transfer should not be accepted yet")
	})
}
