//go:build integration

package offrampops

import (
	"context"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	linkops "github.com/smartcontractkit/chainlink-sui/ops/link"
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

	// Create the client.
	client := suiclient.NewClient("http://localhost:9000")

	// Generate key pair and create a signer.
	pk, _, _, err := testutils.GenerateAccountKeyPair(t, log)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)

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

	// create 4 additional signers + transmitters
	signerAddresses := make([]string, 0, 4) // Preallocate slice with capacity
	signerAddrBytes := make([][]byte, 0, 4)

	for range 4 {
		pk, _, _, err := testutils.GenerateAccountKeyPair(t, logger.Test(t))
		require.NoError(t, err)

		_signer := rel.NewPrivateKeySigner(pk)

		signerAddress, err := _signer.GetAddress()
		require.NoError(t, err)
		signerAddresses = append(signerAddresses, signerAddress)

		addrHex := strings.TrimPrefix(signerAddress, "0x")

		addrBytes, err := hex.DecodeString(addrHex)
		require.NoError(t, err)
		signerAddrBytes = append(signerAddrBytes, addrBytes)
	}

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

	reportMCMs, err := cld_ops.ExecuteOperation(bundle, mcms_ops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Package")

	// Deploy CCIP
	inputCCIP := ccip_ops.DeployCCIPInput{
		McmsPackageId: reportMCMs.Output.PackageId,
		McmsOwner:     "0x2",
	}

	reportCCIP, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	// Initialize feeQuoter
	feeQuoterInit := ccip_ops.InitFeeQuoterInput{
		CCIPPackageId:                 reportCCIP.Output.PackageId,
		StateObjectId:                 reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		OwnerCapObjectId:              reportCCIP.Output.Objects.OwnerCapObjectId,
		MaxFeeJuelsPerMsg:             "100000000",
		LinkTokenCoinMetadataObjectId: linkReport.Output.Objects.CoinMetadataObjectId,
		TokenPriceStalenessThreshold:  60,
	}

	reportFeeQuoterInit, err := cld_ops.ExecuteOperation(bundle, ccip_ops.FeeQuoterInitializeOp, deps, feeQuoterInit)
	require.NoError(t, err, "failed to initialize Fee Quoter Package")

	// Run OffRamp Sequence
	seqOffRampInput := DeployAndInitCCIPOffRampSeqInput{
		DeployCCIPOffRampInput: DeployCCIPOffRampInput{
			CCIPPackageId: reportCCIP.Output.PackageId,
			MCMSPackageId: reportMCMs.Output.PackageId,
		},
		InitializeOffRampInput: InitializeOffRampInput{
			DestTransferCapId:                     reportCCIP.Output.Objects.DestTransferCapObjectId,
			FeeQuoterCapId:                        reportFeeQuoterInit.Output.Objects.FeeQuoterCapObjectId,
			ChainSelector:                         2,
			PremissionExecThresholdSeconds:        10,
			SourceChainSelectors:                  []uint64{1},
			SourceChainsIsEnabled:                 []bool{true},
			SourceChainsIsRMNVerificationDisabled: []bool{true},
			SourceChainsOnRamp:                    [][]byte{{0x01}},
		},
		SetOCR3ConfigInput: SetOCR3ConfigInput{
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
	}

	_, err = cld_ops.ExecuteSequence(bundle, DeployAndInitCCIPOffRampSequence, deps, seqOffRampInput)
	require.NoError(t, err, "failed to deploy CCIP OffRamp Package")
}
