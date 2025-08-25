//go:build integration

package managedtokenops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/ops/link"
	mcms_ops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
)

func TestDeployAndInitManagedToken(t *testing.T) {
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

	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
	)

	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)

	// Deploy MCMS
	reportMCMs, err := cld_ops.ExecuteOperation(bundle, mcms_ops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Package")

	// Deploy LINK Token
	reportLinkToken, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK Token")

	// Deploy and initialize managed token
	input := DeployAndInitManagedTokenInput{
		ManagedTokenDeployInput: ManagedTokenDeployInput{
			MCMSAddress:      reportMCMs.Output.PackageID,
			MCMSOwnerAddress: signerAddress,
		},
		CoinObjectTypeArg:   reportLinkToken.Output.PackageID + "::link::LINK",
		TreasuryCapObjectId: reportLinkToken.Output.Objects.TreasuryCapObjectId,
		DenyCapObjectId:     "", // Empty for basic initialization
		MinterAddress:       signerAddress,
		Allowance:           0,
		IsUnlimited:         true,
	}

	output, err := cld_ops.ExecuteSequence(bundle, DeployAndInitManagedTokenSequence, deps, input)
	require.NoError(t, err)

	// Verify the deployment was successful
	require.NotEmpty(t, output.Output.ManagedTokenPackageId)
	require.NotEmpty(t, output.Output.Objects.OwnerCapObjectId)
	require.NotEmpty(t, output.Output.Objects.StateObjectId)

	// Verify the package ID is a valid address
	require.Len(t, output.Output.ManagedTokenPackageId, 66) // 0x + 64 hex chars
	require.Contains(t, output.Output.ManagedTokenPackageId, "0x")

	// Verify the object IDs are valid
	require.Len(t, output.Output.Objects.OwnerCapObjectId, 66)
	require.Contains(t, output.Output.Objects.OwnerCapObjectId, "0x")
	require.Len(t, output.Output.Objects.StateObjectId, 66)
	require.Contains(t, output.Output.Objects.StateObjectId, "0x")

	t.Logf("Successfully deployed and initialized ManagedToken")
	t.Logf("Package ID: %s", output.Output.ManagedTokenPackageId)
	t.Logf("OwnerCap Object ID: %s", output.Output.Objects.OwnerCapObjectId)
	t.Logf("State Object ID: %s", output.Output.Objects.StateObjectId)
}
