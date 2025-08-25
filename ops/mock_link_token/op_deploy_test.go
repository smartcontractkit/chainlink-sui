//go:build integration

package mocklinktokenops

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

	// Deploy Mock LINK Token
	mockLinkReport, err := cld_ops.ExecuteOperation(bundle, DeployMockLinkTokenOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy Mock LINK token")

	// Mint Mock Link Token
	_, err = cld_ops.ExecuteOperation(bundle, MintMockLinkTokenOp, deps, MintMockLinkTokenInput{
		MockLinkTokenPackageID: mockLinkReport.Output.PackageID,
		TreasuryCapID:          mockLinkReport.Output.Objects.TreasuryCapObjectID,
		Amount:                 10,
	})
	require.NoError(t, err, "failed to mint Mock LINK token")
}
