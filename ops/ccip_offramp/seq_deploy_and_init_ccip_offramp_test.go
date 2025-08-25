//go:build integration

package offrampops

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	linkops "github.com/smartcontractkit/chainlink-sui/ops/link"
	mcms_ops "github.com/smartcontractkit/chainlink-sui/ops/mcms"

	"github.com/stretchr/testify/require"
)

func TestDeployAndInitCCIPOfframpSeq(t *testing.T) {
	t.Parallel()

	signer, client := testenv.SetupEnvironment(t)

	// create 4 additional signers + transmitters
	signerAddresses := make([]string, 0, 4) // Preallocate slice with capacity
	signerAddrBytes := make([][]byte, 0, 4)

	for range 4 {
		additionalSigner, _ := testenv.CreateTestAccount(t)

		signerAddress, err := additionalSigner.GetAddress()
		require.NoError(t, err)
		signerAddresses = append(signerAddresses, signerAddress)

		addrHex := strings.TrimPrefix(signerAddress, "0x")

		addrBytes, err := hex.DecodeString(addrHex)
		require.NoError(t, err)
		signerAddrBytes = append(signerAddrBytes, addrBytes)
	}

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
		McmsPackageID: reportMCMs.Output.PackageId,
		McmsOwner:     signerAddress,
	}

	reportCCIP, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	// Initialize feeQuoter
	feeQuoterInit := ccip_ops.InitFeeQuoterInput{
		CCIPPackageID:                 reportCCIP.Output.PackageId,
		StateObjectID:                 reportCCIP.Output.Objects.CCIPObjectRefObjectID,
		OwnerCapObjectID:              reportCCIP.Output.Objects.OwnerCapObjectID,
		MaxFeeJuelsPerMsg:             "100000000",
		LinkTokenCoinMetadataObjectID: linkReport.Output.Objects.CoinMetadataObjectId,
		TokenPriceStalenessThreshold:  60,
	}

	_, err = cld_ops.ExecuteOperation(bundle, ccip_ops.FeeQuoterInitializeOp, deps, feeQuoterInit)
	require.NoError(t, err, "failed to initialize Fee Quoter Package")

	// Issue fee quoter cap
	issueFeeQuoterCapInput := ccip_ops.IssueFeeQuoterCapInput{
		CCIPPackageID:    reportCCIP.Output.PackageId,
		OwnerCapObjectID: reportCCIP.Output.Objects.OwnerCapObjectID,
	}

	reportIssueFeeQuoterCap, err := cld_ops.ExecuteOperation(bundle, ccip_ops.FeeQuoterIssueFeeQuoterCapOp, deps, issueFeeQuoterCapInput)
	require.NoError(t, err, "failed to issue Fee Quoter Cap")

	// Run OffRamp Sequence
	seqOffRampInput := DeployAndInitCCIPOffRampSeqInput{
		DeployCCIPOffRampInput: DeployCCIPOffRampInput{
			CCIPPackageID: reportCCIP.Output.PackageId,
			MCMSPackageID: reportMCMs.Output.PackageId,
		},
		InitializeOffRampInput: InitializeOffRampInput{
			DestTransferCapID:                     reportCCIP.Output.Objects.DestTransferCapObjectID,
			FeeQuoterCapID:                        reportIssueFeeQuoterCap.Output.Objects.FeeQuoterCapObjectID,
			ChainSelector:                         2,
			PremissionExecThresholdSeconds:        10,
			SourceChainSelectors:                  []uint64{1},
			SourceChainsIsEnabled:                 []bool{true},
			SourceChainsIsRMNVerificationDisabled: []bool{true},
			SourceChainsOnRamp:                    [][]byte{{0x01}},
		},
		CommitOCR3Config: SetOCR3ConfigInput{
			// Sample config digest
			ConfigDigest: []byte{
				0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
				0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
				0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
				0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
			},
			OCRPluginType:                  byte(0),
			BigF:                           byte(1),
			IsSignatureVerificationEnabled: true,
			Signers:                        signerAddrBytes,
			Transmitters:                   signerAddresses,
		},
		ExecutionOCR3Config: SetOCR3ConfigInput{
			ConfigDigest: []byte{
				0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
				0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
				0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
				0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
			},
			OCRPluginType:                  byte(1),
			BigF:                           byte(1),
			IsSignatureVerificationEnabled: false,
			Signers:                        signerAddrBytes,
			Transmitters:                   signerAddresses,
		},
	}

	_, err = cld_ops.ExecuteSequence(bundle, DeployAndInitCCIPOffRampSequence, deps, seqOffRampInput)
	require.NoError(t, err, "failed to deploy CCIP OffRamp Package")
}
