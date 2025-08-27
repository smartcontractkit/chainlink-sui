//go:build integration

package mcmsops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"

	"github.com/stretchr/testify/require"
)

func TestDeployMCMSSeq(t *testing.T) {
	t.Parallel()

	signer, client := testenv.SetupEnvironment(t)

	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(300_000_000)
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

	report, err := cld_ops.ExecuteSequence(bundle, DeployMCMSSequence, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to execute MCMS deploy sequence")

	objects := report.Output.Objects
	require.NotEmpty(t, objects.McmsMultisigStateObjectId, "MCMS Multisig State Object ID should not be empty")
	require.NotEmpty(t, objects.TimelockObjectId, "MCMS Timelock Object ID should not be empty")
	require.NotEmpty(t, objects.McmsDeployerObjectId, "MCMS Deployer Object ID should not be empty")
	require.NotEmpty(t, objects.McmsRegistryObjectId, "MCMS Registry Object ID should not be empty")
	require.NotEmpty(t, objects.McmsAccountStateObjectId, "MCMS Account State Object ID should not be empty")
	require.NotEmpty(t, objects.McmsAccountOwnerCapObjectId, "MCMS Account Owner Cap Object ID should not be empty")
	require.NotEmpty(t, report.Output.Digest, "Transaction digest should not be empty")
	require.NotEmpty(t, report.Output.PackageId, "Package ID should not be empty")
}
