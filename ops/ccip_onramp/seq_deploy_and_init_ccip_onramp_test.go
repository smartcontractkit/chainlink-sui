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
		McmsPackageID: reportMCMs.Output.PackageID,
		McmsOwner:     signerAddress,
	}

	report, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// report from CCIP
	nonceManagerInput := ccip_ops.InitNMInput{
		CCIPPackageID:    report.Output.PackageID,
		StateObjectID:    report.Output.Objects.CCIPObjectRefObjectID,
		OwnerCapObjectID: report.Output.Objects.OwnerCapObjectID,
	}

	reportNonceManagerInit, err := cld_ops.ExecuteOperation(bundle, ccip_ops.NonceManagerInitializeOp, deps, nonceManagerInput)
	require.NoError(t, err, "failed to initialize Nonce Manager Package")

	inputOnRamp := DeployAndInitCCIPOnRampSeqInput{
		DeployCCIPOnRampInput: DeployCCIPOnRampInput{
			CCIPPackageID:      report.Output.PackageID,
			MCMSPackageID:      reportMCMs.Output.PackageID,
			MCMSOwnerPackageID: signerAddress,
		},
		OnRampInitializeInput: OnRampInitializeInput{
			NonceManagerCapID:         reportNonceManagerInit.Output.Objects.NonceManagerCapObjectID, // this is from NonceManager init Op
			SourceTransferCapID:       report.Output.Objects.SourceTransferCapObjectID,               // this is from CCIP package publish
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
		OnRampPackageID:   reportOnRamp.Output.CCIPOnRampPackageID,
		StateObjectID:     reportOnRamp.Output.Objects.StateObjectID,
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
