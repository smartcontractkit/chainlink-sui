//go:build integration

package managedtokenpoolops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	ccip_tokenpoolops "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/ops/managed_token"
	mcms_ops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
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

	signerAddress, err := signer.GetAddress()
	require.NoError(t, err, "failed to get signer address")

	reportMCMs, err := cld_ops.ExecuteOperation(bundle, mcms_ops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Package")

	// Deploy CCIP
	inputCCIP := ccip_ops.DeployCCIPInput{
		McmsPackageID: reportMCMs.Output.PackageID,
		McmsOwner:     signerAddress,
	}

	reportCCIP, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// deploy CCIP Token Pool
	inputTokenPool := ccip_tokenpoolops.TokenPoolDeployInput{
		CCIPPackageID:    reportCCIP.Output.PackageID,
		MCMSAddress:      reportMCMs.Output.PackageID,
		MCMSOwnerAddress: signerAddress,
	}

	reportCCIPTokenPool, err := cld_ops.ExecuteOperation(bundle, ccip_tokenpoolops.DeployCCIPTokenPoolOp, deps, inputTokenPool)
	require.NoError(t, err, "failed to deploy CCIP TokenPool Package")

	// deploy managed token
	reportManagedToken, err := cld_ops.ExecuteOperation(bundle, managedtokenops.DeployCCIPManagedTokenOp, deps, managedtokenops.ManagedTokenDeployInput{
		MCMSAddress:      reportMCMs.Output.PackageID,
		MCMSOwnerAddress: signerAddress,
	})
	require.NoError(t, err, "failed to deploy ManagedToken Package")

	// Initialize TokenAdminRegistry
	inputTAR := ccip_ops.InitTARInput{
		CCIPPackageID:      reportCCIP.Output.PackageID,
		StateObjectID:      reportCCIP.Output.Objects.CCIPObjectRefObjectID,
		OwnerCapObjectID:   reportCCIP.Output.Objects.OwnerCapObjectID,
		LocalChainSelector: 10,
	}

	_, err = cld_ops.ExecuteOperation(bundle, ccip_ops.TokenAdminRegistryInitializeOp, deps, inputTAR)
	require.NoError(t, err, "failed to initialize TokenAdminRegistry")

	// Test just the package deployment for now
	managedTokenPoolInput := ManagedTokenPoolDeployInput{
		CCIPPackageID:          reportCCIP.Output.PackageID,
		CCIPTokenPoolPackageID: reportCCIPTokenPool.Output.PackageID,
		ManagedTokenPackageID:  reportManagedToken.Output.PackageID,
		MCMSAddress:            reportMCMs.Output.PackageID,
		MCMSOwnerAddress:       signerAddress,
	}

	_, err = cld_ops.ExecuteOperation(bundle, DeployCCIPManagedTokenPoolOp, deps, managedTokenPoolInput)
	require.NoError(t, err, "failed to deploy CCIP ManagedTokenPool")
}
