//go:build integration

package onrampops

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	mcms_ops "github.com/smartcontractkit/chainlink-sui/ops/mcms"

	"github.com/stretchr/testify/require"
)

func TestDeployAndInitCCIPOnrampSeq(t *testing.T) {
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

	inputCCIP := ccip_ops.DeployCCIPInput{
		McmsPackageId: reportMCMs.Output.PackageId,
		McmsOwner:     signerAddress,
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
			MCMSOwnerPackageId: signerAddress,
		},
		OnRampInitializeInput: OnRampInitializeInput{
			NonceManagerCapId:         reportNonceManagerInit.Output.Objects.NonceManagerCapObjectId, // this is from NonceManager init Op
			SourceTransferCapId:       report.Output.Objects.SourceTransferCapObjectId,               // this is from CCIP package publish
			ChainSelector:             909606746561742123,
			FeeAggregator:             signerAddress,
			AllowListAdmin:            signerAddress,
			DestChainSelectors:        []uint64{909606746561742123},
			DestChainEnabled:          []bool{true},
			DestChainAllowListEnabled: []bool{true},
		},
		ApplyDestChainConfigureOnRampInput: ApplyDestChainConfigureOnRampInput{
			DestChainSelector:         []uint64{909606746561742123},
			DestChainEnabled:          []bool{true},
			DestChainAllowListEnabled: []bool{false},
		},
		ApplyAllowListUpdatesInput: ApplyAllowListUpdatesInput{
			DestChainSelector:             []uint64{909606746561742123},
			DestChainAllowListEnabled:     []bool{false},
			DestChainAddAllowedSenders:    [][]string{{}},
			DestChainRemoveAllowedSenders: [][]string{{}},
		},
	}

	// Run onRamp deploy & Apply dest chain update sequence
	reportOnRamp, err := cld_ops.ExecuteSequence(bundle, DeployAndInitCCIPOnRampSequence, deps, inputOnRamp)
	require.NoError(t, err, "failed to execute CCIP OnRamp deploy sequence")

	// success case
	isChainSupportedInput := IsChainSupportedInput{
		OnRampPackageId:   reportOnRamp.Output.CCIPOnRampPackageId,
		StateObjectId:     reportOnRamp.Output.Objects.StateObjectId,
		DestChainSelector: 909606746561742123,
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
