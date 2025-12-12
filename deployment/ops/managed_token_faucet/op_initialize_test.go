//go:build integration

package managedtokenfaucetops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
	mcms_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

func TestInitializeManagedTokenFaucet(t *testing.T) {
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
	managedTokenInput := managedtokenops.DeployAndInitManagedTokenInput{
		ManagedTokenDeployInput: managedtokenops.ManagedTokenDeployInput{
			MCMSAddress:      reportMCMs.Output.PackageId,
			MCMSOwnerAddress: signerAddress,
		},
		CoinObjectTypeArg:   reportLinkToken.Output.PackageId + "::link::LINK",
		TreasuryCapObjectId: reportLinkToken.Output.Objects.TreasuryCapObjectId,
		DenyCapObjectId:     "", // Empty for basic initialization
		MinterAddress:       signerAddress,
		Allowance:           0,
		IsUnlimited:         true,
	}

	managedTokenReport, err := cld_ops.ExecuteSequence(bundle, managedtokenops.DeployAndInitManagedTokenSequence, deps, managedTokenInput)
	require.NoError(t, err, "failed to deploy and initialize managed token")

	// Deploy managed token faucet package
	deployFaucetInput := DeployManagedTokenFaucetInput{
		ManagedTokenPackageId: managedTokenReport.Output.ManagedTokenPackageId,
		MCMSAddress:           reportMCMs.Output.PackageId,
		MCMSOwnerAddress:      signerAddress,
	}

	deployFaucetReport, err := cld_ops.ExecuteOperation(bundle, DeployManagedTokenFaucetOp, deps, deployFaucetInput)
	require.NoError(t, err, "failed to deploy managed token faucet package")

	// Initialize managed token faucet
	initInput := InitializeManagedTokenFaucetInput{
		ManagedTokenFaucetPackageId: deployFaucetReport.Output.PackageId,
		CoinObjectTypeArg:           managedTokenInput.CoinObjectTypeArg,
		MintCapObjectId:             managedTokenReport.Output.Objects.MinterCapObjectId,
	}

	initReport, err := cld_ops.ExecuteOperation(bundle, InitializeManagedTokenFaucetOp, deps, initInput)
	require.NoError(t, err, "failed to initialize managed token faucet")

	// Verify the initialization was successful
	require.NotEmpty(t, initReport.Output.Digest, "transaction digest should not be empty")
	require.Equal(t, deployFaucetReport.Output.PackageId, initReport.Output.PackageId, "package ID should match deployed faucet package")
	require.NotEmpty(t, initReport.Output.Objects.FaucetStateObjectId, "faucet state object ID should not be empty")

	// Verify the object ID is a valid address format
	require.Len(t, initReport.Output.Objects.FaucetStateObjectId, 66, "faucet state object ID should be 66 characters (0x + 64 hex chars)")
	require.Contains(t, initReport.Output.Objects.FaucetStateObjectId, "0x", "faucet state object ID should start with 0x")

	t.Logf("Successfully initialized ManagedTokenFaucet")
	t.Logf("Package ID: %s", initReport.Output.PackageId)
	t.Logf("Faucet State Object ID: %s", initReport.Output.Objects.FaucetStateObjectId)
}
