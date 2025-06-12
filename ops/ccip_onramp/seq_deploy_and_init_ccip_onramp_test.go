//go:build integration

package onrampops

import (
	"context"
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	mcms_ops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
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
			b := uint64(400_000_000)
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

	signerAddress, err := signer.GetAddress()
	require.NoError(t, err, "failed to get signer address")

	reportMCMs, err := cld_ops.ExecuteOperation(bundle, mcms_ops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Package")

	inputCCIP := ccip_ops.DeployCCIPInput{
		McmsPackageId: reportMCMs.Output.PackageId,
		McmsOwner:     "0x2",
	}

	report, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// report from CCIP
	nonceManagerInput := ccip_ops.InitNMInput{
		CCIPPackageId:    report.Output.PackageId,
		StateObjectId:    report.Output.Objects.CCIPObjectRefObjectId,
		OwnerCapObjectId: report.Output.Objects.OwnerCapObjectId,
	}

	reportNonceManagerInit, err := cld_ops.ExecuteOperation(bundle, ccip_ops.NonceManagerInitializeOp, deps, nonceManagerInput)
	require.NoError(t, err, "failed to initialize Nonce Manager Package")

	inputOnRamp := DeployAndInitCCIPOnRampSeqInput{
		DeployCCIPOnRampInput: DeployCCIPOnRampInput{
			CCIPPackageId:      report.Output.PackageId,
			MCMSPackageId:      reportMCMs.Output.PackageId,
			MCMSOwnerPackageId: "0x2",
		},
		OnRampInitializeInput: OnRampInitializeInput{
			NonceManagerCapId:         reportNonceManagerInit.Output.Objects.NonceManagerCapObjectId, // this is from NonceManager init Op
			SourceTransferCapId:       report.Output.Objects.SourceTransferCapObjectId,               // this is from CCIP package publish
			ChainSelector:             1,
			FeeAggregator:             signerAddress,
			AllowListAdmin:            signerAddress,
			DestChainSelectors:        []uint64{2},
			DestChainEnabled:          []bool{true},
			DestChainAllowListEnabled: []bool{true},
		},
	}

	// Run onRamp deploy & Apply dest chain update sequence
	reportOnRamp, err := cld_ops.ExecuteSequence(bundle, DeployAndInitCCIPOnRampSequence, deps, inputOnRamp)
	require.NoError(t, err, "failed to execute CCIP OnRamp deploy sequence")

	// success case
	isChainSupportedInput := IsChainSupportedInput{
		OnRampPackageId:   reportOnRamp.Output.CCIPOnRampPackageId,
		StateObjectId:     reportOnRamp.Output.Objects.StateObjectId,
		DestChainSelector: 2,
	}

	reportIsChainSupported, err := cld_ops.ExecuteOperation(bundle, IsChainSupportedOp, deps, isChainSupportedInput)
	require.NoError(t, err, "failed to execute isChainSupported operation")
	require.True(t, reportIsChainSupported.Output.Objects.IsChainSupported)

	reportIsChainEnabled, err := cld_ops.ExecuteOperation(bundle, GetDestChainConfigOp, deps, isChainSupportedInput)
	require.NoError(t, err, "failed to execute GetDestChainConfigHandler operation")
	require.True(t, reportIsChainEnabled.Output.Objects.IsChainSupported)

	// failure case
	isChainSupportedInput.DestChainSelector = 3
	reportIsChainSupportedError, err := cld_ops.ExecuteOperation(bundle, IsChainSupportedOp, deps, isChainSupportedInput)
	require.NoError(t, err, "failed to execute isChainSupported operation")

	require.False(t, reportIsChainSupportedError.Output.Objects.IsChainSupported)
}
