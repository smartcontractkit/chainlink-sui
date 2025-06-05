package tokenpoolops

import (
	"context"
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	ccip_routerops "github.com/smartcontractkit/chainlink-sui/ops/ccip_router"

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

func TestDeployAndInitSeq(t *testing.T) {
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

	inputCCIP := ccip_ops.DeployCCIPInput{
		McmsPackageId: "0x2",
	}

	// deploy CCIP package
	reportCCIP, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// deploy Router package
	reportCCIPRouter, err := cld_ops.ExecuteOperation(bundle, ccip_routerops.DeployCCIPRouterOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy CCIP Package")

	// deploy CCIP Token Pool
	inputTokenPool := TokenPoolDeployInput{
		CCIPPackageId:     reportCCIP.Output.PackageId,
		CCIPRouterAddress: reportCCIPRouter.Output.PackageId,
	}

	_, err = cld_ops.ExecuteOperation(bundle, DeployCCIPTokenPoolOp, deps, inputTokenPool)
	require.NoError(t, err, "failed to deploy CCIP Package")
}
