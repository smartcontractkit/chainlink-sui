//go:build integration

package lnrops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"

	"github.com/stretchr/testify/require"
)

func TestDeployAndInitSeq(t *testing.T) {
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

	// Deploy LnR Token
	lnrTokenReport, err := cld_ops.ExecuteOperation(bundle, DeployLnROp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LnR token")

	// Mint LnR Token
	_, err = cld_ops.ExecuteOperation(bundle, MintLnROp, deps, MintLnRTokenInput{
		LnRTokenPackageId: lnrTokenReport.Output.PackageId,
		TreasuryCapId:     lnrTokenReport.Output.Objects.TreasuryCapObjectId,
		Amount:            10,
		ToAddress:         "0x40d438a47eafc6bee64a7f0addeb468d2939920f5661462f90cd8dbae2cdd9cb",
	})
	require.NoError(t, err, "failed to mint LnR token")
}
