//go:build integration

package mcmsops

import (
	"context"
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"

	"github.com/stretchr/testify/require"
)

func setupSuiTest(t *testing.T) (rel.SuiSigner, *suiclient.ClientImpl) {
	t.Helper()

	log := logger.Test(t)

	// Start the node.
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			if perr := cmd.Process.Kill(); perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	// Generate key pair and create a signer.
	pk, _, _, err := testutils.GenerateAccountKeyPair(t, log)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)

	// Create the client.
	client := suiclient.NewClient("http://localhost:9000")

	// Fund the account.
	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)
	err = testutils.FundWithFaucet(log, "localnet", signerAddress)
	require.NoError(t, err)

	return signer, client
}

func TestDeployMCMSSeq(t *testing.T) {
	t.Parallel()
	signer, client := setupSuiTest(t)

	deps := sui_ops.OpTxDeps{
		Client: *client,
		Signer: signer,
		GetTxOpts: func() bind.TxOpts {
			b := uint64(300_000_000)
			return bind.TxOpts{
				GasBudget: &b,
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
