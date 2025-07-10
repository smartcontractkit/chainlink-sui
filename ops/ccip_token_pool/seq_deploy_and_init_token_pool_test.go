//go:build integration

package tokenpoolops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"

	"github.com/stretchr/testify/require"
)

func TestDeployAndInitTokenPoolSeq(t *testing.T) {
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

	// Deploy MCMS
	mcmsReport, err := cld_ops.ExecuteOperation(bundle, mcmsops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Contract")

	inputCCIP := ccip_ops.DeployCCIPInput{
		McmsPackageId: mcmsReport.Output.PackageId,
		McmsOwner:     "0x2",
	}

	// deploy CCIP package
	reportCCIP, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// deploy CCIP Token Pool
	inputTokenPool := TokenPoolDeployInput{
		CCIPPackageId:    reportCCIP.Output.PackageId,
		MCMSAddress:      mcmsReport.Output.PackageId,
		MCMSOwnerAddress: "0x2",
	}

	_, err = cld_ops.ExecuteOperation(bundle, DeployCCIPTokenPoolOp, deps, inputTokenPool)
	require.NoError(t, err, "failed to deploy CCIP Token Pool Package")
}
