package lockreleasetokenpoolops

import (
	"context"
	"testing"

	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	ccip_tokenpoolops "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
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

	// Deploy CCIP
	inputCCIP := ccip_ops.DeployCCIPInput{
		McmsPackageId: reportMCMs.Output.PackageId,
		McmsOwner:     "0x2",
	}

	reportCCIP, err := cld_ops.ExecuteOperation(bundle, ccip_ops.DeployCCIPOp, deps, inputCCIP)
	require.NoError(t, err, "failed to deploy CCIP Package")

	// deploy CCIP Token Pool
	inputTokenPool := ccip_tokenpoolops.TokenPoolDeployInput{
		CCIPPackageId:    reportCCIP.Output.PackageId,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: "0x2",
	}

	reportCCIPTokenPool, err := cld_ops.ExecuteOperation(bundle, ccip_tokenpoolops.DeployCCIPTokenPoolOp, deps, inputTokenPool)
	require.NoError(t, err, "failed to deploy CCIP TokenPool Package")

	// Deploy LINK
	linkReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	// Initialize TokenAdminRegistry
	inputTAR := ccip_ops.InitTARInput{
		CCIPPackageId:      reportCCIP.Output.PackageId,
		StateObjectId:      reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		OwnerCapObjectId:   reportCCIP.Output.Objects.OwnerCapObjectId,
		LocalChainSelector: 10,
	}

	_, err = cld_ops.ExecuteOperation(bundle, ccip_ops.TokenAdminRegistryInitializeOp, deps, inputTAR)
	require.NoError(t, err, "failed to deploy LINK token")

	// Run BurnMintTokenPool deploy + configure sequence
	LRTokenPoolInput := DeployAndInitLockReleaseTokenPoolInput{
		LockReleaseTokenPoolDeployInput: LockReleaseTokenPoolDeployInput{
			CCIPPackageId:          reportCCIP.Output.PackageId,
			CCIPTokenPoolPackageId: reportCCIPTokenPool.Output.PackageId,
			MCMSAddress:            reportMCMs.Output.PackageId,
			MCMSOwnerAddress:       "0x2",
		},

		// initialize
		CoinObjectTypeArg:      linkReport.Output.PackageId + "::link_token::LINK_TOKEN",
		CCIPObjectRefObjectId:  reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		CoinMetadataObjectId:   linkReport.Output.Objects.CoinMetadataObjectId,
		TreasuryCapObjectId:    linkReport.Output.Objects.TreasuryCapObjectId,
		TokenPoolPackageId:     reportCCIPTokenPool.Output.PackageId,
		TokenPoolAdministrator: signerAddress,
		Rebalancer:             "",

		// apply dest chain updates
		RemoteChainSelectorsToRemove: []uint64{},
		RemoteChainSelectorsToAdd:    []uint64{10},
		RemotePoolAddressesToAdd: [][][]byte{
			{
				[]byte{255, 42, 71, 253, 186, 134, 235, 132, 31, 239, 183, 238, 227, 238, 107, 189, 171, 84, 45, 136},
			},
		},
		RemoteTokenAddressesToAdd: [][]byte{
			{103, 150, 111, 194, 2, 150, 82, 27, 22, 140, 233, 220, 73, 235, 236, 155, 221, 2, 161, 243},
		},

		// set chain rate limiter configs
		RemoteChainSelectors: []uint64{10},
		OutboundIsEnableds:   []bool{true},
		OutboundCapacities:   []uint64{10},
		OutboundRates:        []uint64{10},
		InboundIsEnableds:    []bool{true},
		InboundCapacities:    []uint64{10},
		InboundRates:         []uint64{10},
	}

	_, err = cld_ops.ExecuteSequence(bundle, DeployAndInitLockReleaseTokenPoolSequence, deps, LRTokenPoolInput)
	require.NoError(t, err, "failed to deploy CCIP LockRelease token pool Sequence")
}
