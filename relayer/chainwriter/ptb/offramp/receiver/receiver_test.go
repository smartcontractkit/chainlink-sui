package receiver_test

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/holiman/uint256"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	lockreleaseops "github.com/smartcontractkit/chainlink-sui/ops/ccip_lock_release_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_offramp"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
	mocklinktokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_link_token"
	receiver_module "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp/receiver"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

const SUI_CHAIN_SELECTOR = 2
const ETHEREUM_CHAIN_SELECTOR = 1

type EnvironmentSettings struct {
	// Deployment reports
	MockLinkReport      cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]]
	CCIPReport          cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput]
	OffRampReport       cld_ops.SequenceReport[offrampops.DeployAndInitCCIPOffRampSeqInput, offrampops.DeployCCIPOffRampSeqOutput]
	TokenPoolReport     cld_ops.SequenceReport[lockreleaseops.DeployAndInitLockReleaseTokenPoolInput, lockreleaseops.DeployLockReleaseTokenPoolOutput]
	DummyReceiverReport *cld_ops.SequenceReport[ccipops.DeployAndInitDummyReceiverInput, ccipops.DeployDummyReceiverSeqOutput]

	EthereumPoolAddress []byte

	// Signers
	SignersAddrBytes [][]byte
	Signer           rel.SuiSigner

	// Public keys
	PublicKeys [][]byte

	// Private keys
	PrivateKeys []ed25519.PrivateKey

	// Client
	Client sui.ISuiAPI
}

func setupClients(t *testing.T, lggr logger.Logger) (rel.SuiSigner, sui.ISuiAPI) {
	t.Helper()

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
	client := sui.NewSuiClient(testutils.LocalUrl)

	// Generate key pair and create a signer.
	pk, _, _, err := testutils.GenerateAccountKeyPair(t)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)

	// Fund the account.
	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)
	for range 5 {
		err = testutils.FundWithFaucet(lggr, "localnet", signerAddress)
		require.NoError(t, err)
	}

	return signer, client
}

func SetupTestEnvironment(t *testing.T) *EnvironmentSettings {
	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")

	signer, client := setupClients(t, lggr)

	// Declare all arrays
	signerAddrBytes := make([][]byte, 0, 4)
	signerPrivateKeys := make([]ed25519.PrivateKey, 0, 4)

	// Get the main account's public key first
	keystoreInstance := testutils.NewTestKeystore(t)
	_, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	// add 3 generated signers
	for range 3 {
		pk, _, _, err := testutils.GenerateAccountKeyPair(t)
		require.NoError(t, err)

		_signer := rel.NewPrivateKeySigner(pk)

		signerAddress, err := _signer.GetAddress()
		require.NoError(t, err)

		addrHex := strings.TrimPrefix(signerAddress, "0x")
		addrBytes, err := hex.DecodeString(addrHex)
		require.NoError(t, err)
		signerAddrBytes = append(signerAddrBytes, addrBytes)

		signerPrivateKeys = append(signerPrivateKeys, pk)
	}

	// the 4th signer is the account that will call the OffRamp
	// signerAddresses = append(signerAddresses, accountAddress)
	// signerAddrBytes = append(signerAddrBytes, accountAddressBytes)
	// signerPublicKeys = append(signerPublicKeys, publicKeyBytes)
	// signerPrivateKeys = append(signerPrivateKeys, privateKey)

	// Create 20-byte Ethereum addresses for RMN Remote signers
	ethAddr1, err := hex.DecodeString("8a1b2c3d4e5f60718293a4b5c6d7e8f901234567")
	require.NoError(t, err, "failed to decode eth address 1")
	ethAddr2, err := hex.DecodeString("7b8c9dab0c1d2e3f405162738495a6b7c8d9e0f1")
	require.NoError(t, err, "failed to decode eth address 2")
	ethAddr3, err := hex.DecodeString("1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err, "failed to decode eth address 3")
	// For the 4th address, derive a 20-byte address from the account's public key
	accountEthAddr := make([]byte, 20)
	copy(accountEthAddr, publicKeyBytes[:20]) // Take first 20 bytes of the Ed25519 public key

	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(500_000_000)
			return &bind.CallOpts{
				Signer:           signer,
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

	// Deploy LINK
	mockLinkReport, err := cld_ops.ExecuteOperation(bundle, mocklinktokenops.DeployMockLinkTokenOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	configDigest, err := uint256.FromHex("0xe3b1c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	require.NoError(t, err, "failed to convert config digest to uint256")

	// Deploy MCMs
	reportMCMs, err := cld_ops.ExecuteOperation(bundle, mcmsops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Package")
	lggr.Debugw("MCMS deployment report", "output", reportMCMs.Output)

	signerAddr, err := signer.GetAddress()
	require.NoError(t, err)

	lggr.Debugw("LINK report", "output", mockLinkReport.Output)

	report, err := cld_ops.ExecuteSequence(bundle, ccipops.DeployAndInitCCIPSequence, deps, ccipops.DeployAndInitCCIPSeqInput{
		LinkTokenCoinMetadataObjectId: mockLinkReport.Output.Objects.CoinMetadataObjectId,
		LocalChainSelector:            1,
		DestChainSelector:             2,
		DeployCCIPInput: ccipops.DeployCCIPInput{
			McmsPackageId: reportMCMs.Output.PackageId,
			McmsOwner:     signerAddr,
		},
		MaxFeeJuelsPerMsg:            "100000000",
		TokenPriceStalenessThreshold: 60,
		// Fee Quoter configuration
		AddMinFeeUsdCents:    []uint32{3000},
		AddMaxFeeUsdCents:    []uint32{30000},
		AddDeciBps:           []uint16{1000},
		AddDestGasOverhead:   []uint32{1000000},
		AddDestBytesOverhead: []uint32{1000},
		AddIsEnabled:         []bool{true},
		RemoveTokens:         []string{},
		// Fee Quoter destination chain configuration
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           2,
		MaxDataBytes:                      2000,
		MaxPerMsgGasLimit:                 5000000,
		DestGasOverhead:                   1000000,
		DestGasPerPayloadByteBase:         byte(2),
		DestGasPerPayloadByteHigh:         byte(5),
		DestGasPerPayloadByteThreshold:    uint16(10),
		DestDataAvailabilityOverheadGas:   300000,
		DestGasPerDataAvailabilityByte:    4,
		DestDataAvailabilityMultiplierBps: 1,
		ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c},
		EnforceOutOfOrder:                 false,
		DefaultTokenFeeUsdCents:           3,
		DefaultTokenDestGasOverhead:       100000,
		DefaultTxGasLimit:                 500000,
		GasMultiplierWeiPerEth:            100,
		GasPriceStalenessThreshold:        1000000000,
		NetworkFeeUsdCents:                10,
		// Premium multiplier updates
		PremiumMultiplierWeiPerEth: []uint64{10},

		RmnHomeContractConfigDigest: configDigest.Bytes(),
		SignerOnchainPublicKeys:     [][]byte{ethAddr1, ethAddr2, ethAddr3, accountEthAddr},
		NodeIndexes:                 []uint64{0, 1, 2, 3},
		FSign:                       uint64(1),
	})
	require.NoError(t, err, "failed to execute CCIP deploy sequence")
	require.NotEmpty(t, report.Output.CCIPPackageId, "CCIP package ID should not be empty")

	// Deploy and initialize the dummy receiver
	dummyReceiverReport, err := cld_ops.ExecuteSequence(bundle, ccipops.DeployAndInitDummyReceiverSequence, deps, ccipops.DeployAndInitDummyReceiverInput{
		DeployDummyReceiverInput: ccipops.DeployDummyReceiverInput{
			CCIPPackageId: report.Output.CCIPPackageId,
			McmsPackageId: reportMCMs.Output.PackageId,
			McmsOwner:     signerAddr,
		},
		CCIPObjectRefObjectId: report.Output.Objects.CCIPObjectRefObjectId,
	})

	require.NoError(t, err, "failed to deploy and initialize dummy receiver")
	lggr.Debugw("Dummy receiver deployment report", "output", dummyReceiverReport.Output)

	return &EnvironmentSettings{
		MockLinkReport:      mockLinkReport,
		CCIPReport:          report,
		DummyReceiverReport: &dummyReceiverReport,
		SignersAddrBytes:    signerAddrBytes,
		Signer:              signer,
		PublicKeys:          [][]byte{ethAddr1, ethAddr2, ethAddr3, accountEthAddr},
		PrivateKeys:         signerPrivateKeys,
		Client:              client,
	}
}

func TestReceiver(t *testing.T) {
	env := SetupTestEnvironment(t)
	lggr := logger.Test(t)
	ctx := context.Background()

	ccipObjectRef := env.CCIPReport.Output.Objects.CCIPObjectRefObjectId
	ccipPackageId := env.CCIPReport.Output.CCIPPackageId

	t.Run("TestFilterRegisteredReceivers", func(t *testing.T) {
		t.Skip()
		// Use the dummy receiver that was actually registered
		receiverPackageId := env.DummyReceiverReport.Output.DummyReceiverPackageId
		receiverModule := "ccip_dummy_receiver"
		receiver := fmt.Sprintf("%s::%s::ccip_receive", receiverPackageId, receiverModule)

		msg := ccipocr3.Message{
			Receiver: []byte(receiver),
			Data:     []byte("Hello World"),
		}

		signerAddress, err := env.Signer.GetAddress()
		require.NoError(t, err)

		ptbClient, err := client.NewPTBClient(lggr, testutils.LocalUrl, nil, 10*time.Second, nil, 5, "WaitForLocalExecution")
		require.NoError(t, err, "Failed to create PTB client for event querying")

		registeredReceivers, err := receiver_module.FilterRegisteredReceivers(
			ctx,
			lggr,
			[]ccipocr3.Message{msg},
			signerAddress,
			ptbClient,
			ccipObjectRef,
			ccipPackageId)
		require.NoError(t, err)
		require.Equal(t, 1, len(registeredReceivers))
		require.Equal(t, msg, registeredReceivers[0])
	})

	t.Run("TestAddReceiverCallCommands", func(t *testing.T) {
		receiverPackageId := env.DummyReceiverReport.Output.DummyReceiverPackageId
		receiverModule := "ccip_dummy_receiver"
		receiver := fmt.Sprintf("%s::%s::echo", receiverPackageId, receiverModule)

		msg := ccipocr3.Message{
			Receiver: []byte(receiver),
			Data:     []byte("Hello World"),
		}

		// Create a new transaction builder
		ptb := transaction.NewTransaction()
		ptb.SetSuiClient(env.Client.(*sui.Client))

		signerAddress, err := env.Signer.GetAddress()
		require.NoError(t, err)

		ptbClient, err := client.NewPTBClient(lggr, testutils.LocalUrl, nil, 10*time.Second, nil, 5, "WaitForLocalExecution")
		require.NoError(t, err, "Failed to create PTB client for event querying")

		receiverCommands, err := receiver_module.AddReceiverCallCommands(ctx, lggr, ptb, signerAddress, []ccipocr3.Message{msg}, 0, ccipObjectRef, ccipPackageId, ptbClient)
		require.NoError(t, err)
		lggr.Info("receiver commands", "commands", receiverCommands)

		opts := &bind.CallOpts{
			Signer:           env.Signer,
			WaitForExecution: true,
		}

		tx, err := bind.ExecutePTB(ctx, opts, env.Client, ptb)
		require.NoError(t, err)
		lggr.Infow("tx", "tx", tx)

		//require.Equal(t, 1, len(receiverCommands))
	})
}
