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
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
)

func TestDeployAndInitSeq(t *testing.T) {
	t.Parallel()
	signer, client := testenv.SetupEnvironment(t)

	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(1_000_000_000)
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

	inputCCIP, err := ccip_ops.DeployCCIPDependencyPackages(bundle, deps)
	require.NoError(t, err, "failed to deploy CCIP dependency packages")

	reportCCIP, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// deploy managed token
	reportManagedToken, err := cld_ops.ExecuteOperation(bundle, managedtokenops.DeployCCIPManagedTokenOp, deps, managedtokenops.ManagedTokenDeployInput{
		MCMSAddress:      inputCCIP.McmsPackageId,
		MCMSOwnerAddress: signerAddress,
	})
	require.NoError(t, err, "failed to deploy ManagedToken Package")

	// Initialize TokenAdminRegistry
	inputTAR := ccip_ops.InitTARInput{
		CCIPPackageId:      reportCCIP.Output.PackageId,
		StateObjectId:      reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		OwnerCapObjectId:   reportCCIP.Output.Objects.OwnerCapObjectId,
		LocalChainSelector: 10,
	}

	_, err = cld_ops.ExecuteOperation(bundle, ccip_ops.TokenAdminRegistryInitializeOp, deps, inputTAR)
	require.NoError(t, err, "failed to initialize TokenAdminRegistry")

	_, err = cld_ops.ExecuteOperation(bundle, ccip_ops.TokenAdminRegistryInitializeLocalDecimalsOp, deps, ccip_ops.InitLocalDecimalsInput{
		CCIPPackageId:    reportCCIP.Output.PackageId,
		StateObjectId:    reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		OwnerCapObjectId: reportCCIP.Output.Objects.OwnerCapObjectId,
	})
	require.NoError(t, err, "failed to initialize local decimals state")

	// Test just the package deployment for now
	managedTokenPoolInput := ManagedTokenPoolDeployInput{
		CCIPPackageId:         reportCCIP.Output.PackageId,
		ManagedTokenPackageId: reportManagedToken.Output.PackageId,
		MCMSAddress:           inputCCIP.McmsPackageId,
		FastMcmsAddress:       inputCCIP.FastMcmsPackageId,
		MCMSOwnerAddress:      signerAddress,
	}

	_, err = cld_ops.ExecuteOperation(bundle, DeployCCIPManagedTokenPoolOp, deps, managedTokenPoolInput)
	require.NoError(t, err, "failed to deploy CCIP ManagedTokenPool")
}
