//go:build invisible

package ccip_test

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/holiman/uint256"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/require"

	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mocklinktoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_link_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	lockreleaseops "github.com/smartcontractkit/chainlink-sui/ops/ccip_lock_release_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_offramp"
	cciptokenpoolop "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
	mocklinktokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_link_token"
	chainwriter "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/expander"
	offramp "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"golang.org/x/crypto/blake2b"
)

const SUI_CHAIN_SELECTOR = 2
const ETHEREUM_CHAIN_SELECTOR = 1

// PoolInfos mimics the Move struct from token_admin_registry.move
type PoolInfos struct {
	TokenPoolPackageIds     []string `json:"token_pool_package_ids"`
	TokenPoolStateAddresses []string `json:"token_pool_state_addresses"`
	TokenPoolModules        []string `json:"token_pool_modules"`
	TokenTypes              []string `json:"token_types"`
}

// ReceivedMessageEvent represents the ReceivedMessage event from the dummy receiver contract
type ReceivedMessageEvent struct {
	MessageID []byte `json:"message_id"`
	Sender    []byte `json:"sender"`
	Data      []byte `json:"data"`
}

func AnyPointer[T any](v T) *T {
	return &v
}

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

// TestMessage contains all the parameters needed for both commit and execute phases
type TestMessage struct {
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
	MessageID           [32]byte
	Receiver            []byte
	GasLimit            *big.Int
	Sender              []byte
	Data                []byte
	OnRampAddress       []byte
	TokenAmounts        []ccipocr3.RampTokenAmount
	LinkMetadataBytes   []byte
}

// createTestMessage creates a consistent message that can be reused for both commit and execute phases
func createTestMessage(envSettings *EnvironmentSettings) *TestMessage {
	// Core message parameters
	sourceChainSelector := uint64(ETHEREUM_CHAIN_SELECTOR) // 1
	destChainSelector := uint64(SUI_CHAIN_SELECTOR)        // 2
	sequenceNumber := uint64(1)
	nonce := uint64(0)          // Always 0 for out-of-order execution
	messageID := [32]byte{}     // All zeros
	gasLimit := big.NewInt(200) // Must be consistent
	sender := []byte{}          // Empty sender for now
	receiverMessage := "Hello, world!"
	data := []byte(receiverMessage) // Test data - must be consistent

	// OnRamp address (32 bytes)
	onRampAddress := make([]byte, 32)
	onRampAddress[31] = 20 // Same as used elsewhere

	// Process link token address
	linkTokenAddress := envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId
	linkMetadataAddressBytes := make([]byte, 32)
	if len(linkTokenAddress) > 2 && linkTokenAddress[:2] == "0x" {
		if decoded, err := hex.DecodeString(linkTokenAddress[2:]); err == nil {
			copy(linkMetadataAddressBytes[32-len(decoded):], decoded)
		}
	}

	// Create token amounts
	tokenAmounts := []ccipocr3.RampTokenAmount{
		{
			SourcePoolAddress: envSettings.EthereumPoolAddress,
			DestTokenAddress:  linkMetadataAddressBytes,
			ExtraData:         []byte{},
			Amount:            ccipocr3.NewBigInt(big.NewInt(300)), // Must be consistent
		},
	}

	return &TestMessage{
		SourceChainSelector: sourceChainSelector,
		DestChainSelector:   destChainSelector,
		SequenceNumber:      sequenceNumber,
		Nonce:               nonce,
		MessageID:           messageID,
		Receiver:            nil,
		GasLimit:            gasLimit,
		Sender:              sender,
		Data:                data,
		OnRampAddress:       onRampAddress,
		TokenAmounts:        tokenAmounts,
		LinkMetadataBytes:   linkMetadataAddressBytes,
	}
}

// createTestMessageWithReceiver creates a test message that includes a receiver address
func createTestMessageWithReceiver(envSettings *EnvironmentSettings, receiverPackageId string, receiverModule string) *TestMessage {
	msg := createTestMessage(envSettings)

	receiverPayload := fmt.Sprintf("%s::%s::ccip_receive", receiverPackageId, receiverModule)
	msg.Receiver = []byte(receiverPayload)
	return msg
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
	pk, _, _, err := testutils.GenerateAccountKeyPair(t, lggr)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)

	// Fund the account.
	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)
	for range 10 {
		err = testutils.FundWithFaucet(lggr, "localnet", signerAddress)
		require.NoError(t, err)
	}

	return signer, client
}

func SetupTestEnvironment(t *testing.T) *EnvironmentSettings {
	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	accountAddress := testutils.GetAccountAndKeyFromSui(t, lggr)
	accountAddressBytes := []byte(accountAddress)

	signer, client := setupClients(t, lggr)

	// Declare all arrays
	signerAddresses := make([]string, 0, 4)
	signerAddrBytes := make([][]byte, 0, 4)
	signerPublicKeys := make([][]byte, 0, 4)
	signerPrivateKeys := make([]ed25519.PrivateKey, 0, 4)

	// Get the main account's public key first
	keystoreInstance, err := keystore.NewSuiKeystore(lggr, "")
	require.NoError(t, err)
	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes := []byte(publicKey)

	// add 3 generated signers
	for range 3 {
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

		// Extract the public key (32 bytes) for OCR3
		publicKey := pk.Public().(ed25519.PublicKey)
		signerPublicKeys = append(signerPublicKeys, []byte(publicKey))

		signerPrivateKeys = append(signerPrivateKeys, pk)
	}

	// the 4th signer is the account that will call the OffRamp
	signerAddresses = append(signerAddresses, accountAddress)
	signerAddrBytes = append(signerAddrBytes, accountAddressBytes)
	signerPublicKeys = append(signerPublicKeys, publicKeyBytes)
	signerPrivateKeys = append(signerPrivateKeys, privateKey)

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

	// Create a dummy OnRamp address
	OnRampAddress := make([]byte, 32)
	OnRampAddress[31] = 20

	seqOffRampInput := offrampops.DeployAndInitCCIPOffRampSeqInput{
		DeployCCIPOffRampInput: offrampops.DeployCCIPOffRampInput{
			CCIPPackageId: report.Output.CCIPPackageId,
			MCMSPackageId: reportMCMs.Output.PackageId,
		},
		InitializeOffRampInput: offrampops.InitializeOffRampInput{
			DestTransferCapId:                     report.Output.Objects.DestTransferCapObjectId,
			FeeQuoterCapId:                        report.Output.Objects.FeeQuoterCapObjectId,
			ChainSelector:                         SUI_CHAIN_SELECTOR,
			PremissionExecThresholdSeconds:        10,
			SourceChainSelectors:                  []uint64{ETHEREUM_CHAIN_SELECTOR},
			SourceChainsIsEnabled:                 []bool{true},
			SourceChainsIsRMNVerificationDisabled: []bool{true},
			SourceChainsOnRamp:                    [][]byte{OnRampAddress},
		},
		CommitOCR3Config: offrampops.SetOCR3ConfigInput{
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
			Signers:                        signerPublicKeys,
			Transmitters:                   signerAddresses,
		},
		ExecutionOCR3Config: offrampops.SetOCR3ConfigInput{
			ConfigDigest: []byte{
				0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
				0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
				0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
				0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
			},
			OCRPluginType:                  byte(1),
			BigF:                           byte(1),
			IsSignatureVerificationEnabled: false,
			Signers:                        signerPublicKeys,
			Transmitters:                   signerAddresses,
		},
	}

	offrampReport, err := cld_ops.ExecuteSequence(bundle, offrampops.DeployAndInitCCIPOffRampSequence, deps, seqOffRampInput)
	require.NoError(t, err, "failed to deploy CCIP Package")

	lggr.Debugw("Offramp deployment report", "output", offrampReport.Output)

	// Deploy CCIP token pool
	ccipTokenPoolReport, err := cld_ops.ExecuteOperation(bundle, cciptokenpoolop.DeployCCIPTokenPoolOp, deps, cciptokenpoolop.TokenPoolDeployInput{
		CCIPPackageId:    report.Output.CCIPPackageId,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: accountAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Token Pool")

	lggr.Debugw("CCIP Token Pool deployment report", "output", ccipTokenPoolReport.Output)

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", mockLinkReport.Output.PackageId)

	ethereumPoolAddress := []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}
	ethereumPoolAddressString := hex.EncodeToString(ethereumPoolAddress)
	remoteTokenAddress := []byte{0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f, 0x40}
	remoteTokenAddressString := hex.EncodeToString(remoteTokenAddress)

	// Deploy and initialize the lock release token pool
	seqLockReleaseDeployInput := lockreleaseops.DeployAndInitLockReleaseTokenPoolInput{
		LockReleaseTokenPoolDeployInput: lockreleaseops.LockReleaseTokenPoolDeployInput{
			CCIPPackageId:          report.Output.CCIPPackageId,
			CCIPTokenPoolPackageId: ccipTokenPoolReport.Output.PackageId,
			MCMSAddress:            reportMCMs.Output.PackageId,
			MCMSOwnerAddress:       accountAddress,
		},
		// Initialization parameters
		CoinObjectTypeArg:      linkTokenType,
		CCIPObjectRefObjectId:  report.Output.Objects.CCIPObjectRefObjectId,
		CoinMetadataObjectId:   mockLinkReport.Output.Objects.CoinMetadataObjectId,
		TreasuryCapObjectId:    mockLinkReport.Output.Objects.TreasuryCapObjectId,
		TokenPoolAdministrator: accountAddress,
		Rebalancer:             signerAddr,

		// Chain updates - adding the destination chain
		RemoteChainSelectorsToRemove: []uint64{},
		RemoteChainSelectorsToAdd:    []uint64{ETHEREUM_CHAIN_SELECTOR},       // Destination chain selector
		RemotePoolAddressesToAdd:     [][]string{{ethereumPoolAddressString}}, // 32-byte remote pool address
		RemoteTokenAddressesToAdd:    []string{remoteTokenAddressString},      // 32-byte remote token address
		// Rate limiter configurations
		RemoteChainSelectors: []uint64{ETHEREUM_CHAIN_SELECTOR}, // Destination chain selector
		OutboundIsEnableds:   []bool{false},
		OutboundCapacities:   []uint64{1000000}, // 1M tokens capacity
		OutboundRates:        []uint64{100000},  // 100K tokens per time window
		InboundIsEnableds:    []bool{false},
		InboundCapacities:    []uint64{1000000}, // 1M tokens capacity
		InboundRates:         []uint64{100000},  // 100K tokens per time window
	}

	tokenPoolLockReleaseReport, err := cld_ops.ExecuteSequence(bundle, lockreleaseops.DeployAndInitLockReleaseTokenPoolSequence, deps, seqLockReleaseDeployInput)
	require.NoError(t, err, "failed to deploy and initialize Lock Release Token Pool")

	lggr.Debugw("Token Pool Lock Release deployment report", "output", tokenPoolLockReleaseReport.Output)

	// Provide liquidity to the lock release token pool
	// First, mint some LINK tokens using the LINK token contract
	liquidityAmount := uint64(1000000) // 1M tokens for liquidity

	// Create LINK token contract instance
	linkContract, err := mocklinktoken.NewMockLinkToken(mockLinkReport.Output.PackageId, client)
	require.NoError(t, err, "failed to create LINK token contract")

	// Mint LINK tokens to the signer's address
	mintTx, err := linkContract.MockLinkToken().Mint(
		ctx,
		deps.GetCallOpts(),
		bind.Object{Id: mockLinkReport.Output.Objects.TreasuryCapObjectId},
		liquidityAmount,
	)
	require.NoError(t, err, "failed to mint LINK tokens for liquidity")

	lggr.Debugw("Minted LINK tokens for liquidity", "amount", liquidityAmount, "txDigest", mintTx.Digest)

	// Find the minted coin object ID from the transaction
	mintedCoinId, err := bind.FindCoinObjectIdFromTx(*mintTx, linkTokenType)
	require.NoError(t, err, "failed to find minted coin object ID")

	lggr.Debugw("Minted coin ID", "mintedCoinId", mintedCoinId)

	// Provide the minted tokens as liquidity to the pool
	provideLiquidityInput := lockreleaseops.LockReleaseTokenPoolProviderLiquidityInput{
		LockReleaseTokenPoolPackageId: tokenPoolLockReleaseReport.Output.LockReleaseTPPackageID,
		StateObjectId:                 tokenPoolLockReleaseReport.Output.Objects.StateObjectId,
		Coin:                          mintedCoinId,
		CoinObjectTypeArg:             linkTokenType,
	}

	_, err = cld_ops.ExecuteOperation(bundle, lockreleaseops.LockReleaseTokenPoolProviderLiquidityOp, deps, provideLiquidityInput)
	require.NoError(t, err, "failed to provide liquidity to Lock Release Token Pool")

	lggr.Debugw("Provided liquidity to Lock Release Token Pool", "amount", liquidityAmount)

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
		OffRampReport:       offrampReport,
		TokenPoolReport:     tokenPoolLockReleaseReport,
		DummyReceiverReport: &dummyReceiverReport,
		EthereumPoolAddress: []byte(ethereumPoolAddressString), // Use hex string as bytes to match configuration
		SignersAddrBytes:    signerAddrBytes,
		Signer:              signer,
		PublicKeys:          [][]byte{ethAddr1, ethAddr2, ethAddr3, accountEthAddr},
		PrivateKeys:         signerPrivateKeys,
		Client:              client,
	}
}

func setupChainWriter(t *testing.T, envSettings *EnvironmentSettings) (*chainwriter.SuiChainWriter, string) {
	lggr := logger.Test(t)

	keystoreInstance, err := keystore.NewSuiKeystore(lggr, "")
	require.NoError(t, err)

	accountAddress := testutils.GetAccountAndKeyFromSui(t, lggr)
	lggr.Infow("Using account", "address", accountAddress)

	err = testutils.FundWithFaucet(lggr, "localnet", accountAddress)
	require.NoError(t, err)

	// Get private key for signing
	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes := []byte(publicKey)

	offRampPackageId := envSettings.OffRampReport.Output.CCIPOffRampPackageId

	chainWriterConfig := cwConfig.ChainWriterConfig{
		Modules: map[string]*cwConfig.ChainWriterModule{
			cwConfig.PTBChainWriterModuleName: {
				Name:     cwConfig.PTBChainWriterModuleName,
				ModuleID: "0x123",
				Functions: map[string]*cwConfig.ChainWriterFunction{
					"register_token_pool": {
						Name:      "register_token_pool",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []cwConfig.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: AnyPointer(envSettings.CCIPReport.Output.CCIPPackageId),
								ModuleId:  AnyPointer("token_admin_registry"),
								Function:  AnyPointer("register_pool_by_admin"),
								Params: []codec.SuiFunctionParam{
									{
										Name:      "ccip_object_ref",
										Type:      "object_id",
										Required:  true,
										IsMutable: testutils.BoolPointer(true),
									},
									{
										Name:     "coin_metadata_address",
										Type:     "address",
										Required: true,
									},
									{
										Name:     "token_pool_package_id",
										Type:     "address",
										Required: true,
									},
									{
										Name:     "token_pool_state_address",
										Type:     "address",
										Required: true,
									},
									{
										Name:     "token_pool_module",
										Type:     "string",
										Required: true,
									},
									{
										Name:     "token_type",
										Type:     "string",
										Required: true,
									},
									{
										Name:     "initial_administrator",
										Type:     "address",
										Required: true,
									},
									{
										Name:     "proof",
										Type:     "string",
										Required: true,
									},
								},
							},
						},
					},
					"ccip_commit":                          offramp.GenerateCommitPTB(lggr, offRampPackageId, publicKeyBytes),
					cwConfig.CCIPExecuteReportFunctionName: offramp.GenerateExecutePTB(lggr, offRampPackageId, publicKeyBytes),
				},
			},
		},
	}

	_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr)

	chainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, chainWriterConfig, false)
	require.NoError(t, err)

	return chainWriter, accountAddress
}

func createCommitReport(t *testing.T, envSettings *EnvironmentSettings) ([]byte, testutils.CommitReport, [][]byte, *TestMessage) {
	// Create the centralized test message
	testMessage := createTestMessage(envSettings)

	// dummy report context - each element must be 32 bytes for Sui Move
	reportContext := [][]byte{
		make([]byte, 32), // config digest - 32 bytes
		make([]byte, 32), // epoch and round - 32 bytes
	}
	// Add some distinguishing data
	reportContext[0] = []byte{
		0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
		0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
		0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
		0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
	}
	reportContext[1][0] = 0x022

	linkTokenAddress := envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId
	chainSelector := testMessage.SourceChainSelector // Use from test message
	seqNumStart := testMessage.SequenceNumber        // Use from test message
	seqNumEnd := uint64(10)
	gasPrice := big.NewInt(1000000000000000000)
	r := make([]byte, 32)
	for i := 0; i < 32; i++ {
		r[i] = 0x01
	}
	s := make([]byte, 32)
	for i := 0; i < 32; i++ {
		s[i] = 0x02
	}

	// First create a temporary report to get the parameters
	tempReport := testutils.GetCommitReport(
		testMessage.OnRampAddress, // Use from test message
		make([]byte, 32),          // temporary merkle root
		linkTokenAddress,
		big.NewInt(1000000000000000000), // price in wei
		chainSelector,
		seqNumStart,
		seqNumEnd,
		gasPrice,
		r,
		s,
	)

	// Calculate the actual message hash that should be the merkle root
	calculatedMerkleRoot := calculateMessageHashForCommit(testMessage, tempReport)

	// Create the final report with the calculated merkle root
	report := testutils.GetCommitReport(
		testMessage.OnRampAddress, // Use from test message
		calculatedMerkleRoot,
		linkTokenAddress,
		big.NewInt(1000000000000000000), // price in wei
		chainSelector,
		seqNumStart,
		seqNumEnd,
		gasPrice,
		r,
		s,
	)

	// Use the helper function to properly serialize the commit report using BCS format
	bcsBytes, err := testutils.SerializeCommitReport(report)
	require.NoError(t, err, "failed to serialize commit report")

	return bcsBytes, report, reportContext, testMessage
}

// createCommitReportForMessage creates a commit report for a specific message
func createCommitReportForMessage(t *testing.T, envSettings *EnvironmentSettings, testMessage *TestMessage) ([]byte, testutils.CommitReport, [][]byte) {
	// dummy report context - each element must be 32 bytes for Sui Move
	reportContext := [][]byte{
		make([]byte, 32), // config digest - 32 bytes
		make([]byte, 32), // epoch and round - 32 bytes
	}
	// Add some distinguishing data
	reportContext[0] = []byte{
		0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
		0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
		0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
		0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
	}
	reportContext[1][0] = 0x022

	linkTokenAddress := envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId
	chainSelector := testMessage.SourceChainSelector // Use from test message
	seqNumStart := testMessage.SequenceNumber        // Use from test message
	seqNumEnd := uint64(10)
	gasPrice := big.NewInt(1000000000000000000)
	r := make([]byte, 32)
	for i := 0; i < 32; i++ {
		r[i] = 0x01
	}
	s := make([]byte, 32)
	for i := 0; i < 32; i++ {
		s[i] = 0x02
	}

	// First create a temporary report to get the parameters
	tempReport := testutils.GetCommitReport(
		testMessage.OnRampAddress, // Use from test message
		make([]byte, 32),          // temporary merkle root
		linkTokenAddress,
		big.NewInt(1000000000000000000), // price in wei
		chainSelector,
		seqNumStart,
		seqNumEnd,
		gasPrice,
		r,
		s,
	)

	// Calculate the actual message hash that should be the merkle root
	calculatedMerkleRoot := calculateMessageHashForCommit(testMessage, tempReport)

	// Create the final report with the calculated merkle root
	report := testutils.GetCommitReport(
		testMessage.OnRampAddress, // Use from test message
		calculatedMerkleRoot,
		linkTokenAddress,
		big.NewInt(1000000000000000000), // price in wei
		chainSelector,
		seqNumStart,
		seqNumEnd,
		gasPrice,
		r,
		s,
	)

	// Use the helper function to properly serialize the commit report using BCS format
	bcsBytes, err := testutils.SerializeCommitReport(report)
	require.NoError(t, err, "failed to serialize commit report")

	return bcsBytes, report, reportContext
}

// calculateMessageHashForCommit creates a hash that matches what the contract will calculate
// It takes the commit report as input to ensure parameter consistency
func calculateMessageHashForCommit(testMessage *TestMessage, commitReport testutils.CommitReport) []byte {
	// Extract parameters from the commit report to ensure consistency
	if len(commitReport.UnblessedMerkleRoots) == 0 {
		// Fallback to predictable hash if no unblessed roots
		return []byte{
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
		}
	}

	// Get parameters from the commit report's first unblessed merkle root
	merkleRootInfo := commitReport.UnblessedMerkleRoots[0]
	sourceChainSelector := merkleRootInfo.SourceChainSelector // Should be ETHEREUM_CHAIN_SELECTOR (1)
	destChainSelector := testMessage.DestChainSelector        // Use from test message
	sequenceNumber := merkleRootInfo.MinSeqNr                 // Should be 1
	nonce := testMessage.Nonce                                // Use from test message
	messageID := testMessage.MessageID                        // Use from test message
	receiver := testMessage.Receiver                          // Use from test message
	gasLimit := testMessage.GasLimit                          // Use from test message
	sender := testMessage.Sender                              // Use from test message
	data := testMessage.Data                                  // Use from test message

	// OnRamp address from commit report
	onRampAddress := merkleRootInfo.OnRampAddress

	// Use token amounts directly from test message (already properly formatted)
	tokenAmounts := testMessage.TokenAmounts

	// Calculate metadata hash using commit report parameters
	metadataHash, err := testutils.ComputeMetadataHash(
		sourceChainSelector,
		destChainSelector,
		onRampAddress,
	)
	if err != nil {
		// Fallback to a predictable hash if calculation fails
		return []byte{
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
		}
	}

	// Calculate message hash using the same parameters
	messageHash, err := testutils.ComputeMessageDataHash(
		metadataHash,
		messageID,
		receiver,
		sequenceNumber,
		gasLimit,
		nonce,
		sender,
		data,
		tokenAmounts,
		uint32(200),
	)
	if err != nil {
		// Fallback to a predictable hash if calculation fails
		return []byte{
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
			0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
		}
	}

	return messageHash[:]
}

// TestCCIPOffRamp tests the CCIP offramp functionality
func TestCCIPOffRamp(t *testing.T) {
	t.Parallel()
	t.Skip("Skipping CCIP Offramp test")

	lggr := logger.Test(t)

	envSettings := SetupTestEnvironment(t)

	lggr.Infow("Link token address", "linkTokenAddress", envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()
	chainWriter, accountAddress := setupChainWriter(t, envSettings)

	err := chainWriter.Start(ctx)
	require.NoError(t, err)

	clockObject := "0x6"

	bcsBytes, report, reportContext, testMessage := createCommitReport(t, envSettings)

	ccipObjectRef := envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId
	require.NotEmpty(t, ccipObjectRef, "CCIP object reference should not be empty")
	offrampStateObjectId := envSettings.OffRampReport.Output.Objects.StateObjectId
	require.NotEmpty(t, offrampStateObjectId, "Offramp state object ID should not be empty")
	linkTokenAddress := envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId

	// Add debugging: print the serialized bytes and the report structure
	lggr.Debugf("=== DEBUG INFO ===")
	lggr.Debugf("Report structure:")
	lggr.Debugf("  TokenPriceUpdates: %d items", len(report.PriceUpdates.TokenPriceUpdates))
	for i, tpu := range report.PriceUpdates.TokenPriceUpdates {
		lggr.Debugf("    [%d] SourceToken: %x (len=%d)", i, tpu.SourceToken, len(tpu.SourceToken))
		lggr.Debugf("    [%d] UsdPerToken: %s", i, tpu.UsdPerToken.String())
	}
	lggr.Debugf("  GasPriceUpdates: %d items", len(report.PriceUpdates.GasPriceUpdates))
	lggr.Debugf("  BlessedMerkleRoots: %d items", len(report.BlessedMerkleRoots))
	lggr.Debugf("  UnblessedMerkleRoots: %d items", len(report.UnblessedMerkleRoots))
	lggr.Debugf("  RMNSignatures: %d items", len(report.RMNSignatures))
	for i, sig := range report.RMNSignatures {
		lggr.Debugf("    [%d] Signature: %x (len=%d)", i, sig, len(sig))
	}
	lggr.Debugf("Serialized bytes: %x", bcsBytes)
	lggr.Debugf("Serialized length: %d", len(bcsBytes))
	lggr.Debugf("Committed Merkle Root: %x", report.UnblessedMerkleRoots[0].MerkleRoot)
	lggr.Debug("==================")

	// Create proper signatures for OCR3 verification
	configDigest := reportContext[0]
	sequenceBytes := reportContext[1]

	// Hash the report the same way the contract does: report + config_digest + sequence_bytes
	var reportForHashing []byte
	reportForHashing = append(reportForHashing, bcsBytes...)
	reportForHashing = append(reportForHashing, configDigest...)
	reportForHashing = append(reportForHashing, sequenceBytes...)
	hashedReport := blake2b.Sum256(reportForHashing)

	// Create signatures using the first big_f + 1 = 2 private keys
	signatures := make([][]byte, 2)
	for i := 0; i < 2; i++ {
		privateKey := envSettings.PrivateKeys[i]
		publicKey := privateKey.Public().(ed25519.PublicKey)

		// Sign the hashed report with raw Ed25519 (no Sui intent bytes)
		signature := ed25519.Sign(privateKey, hashedReport[:])

		require.Equal(t, len(signature), 64)

		// Create 96-byte signature: 32 bytes pubkey + 64 bytes signature
		sig96 := make([]byte, 96)
		copy(sig96[0:32], publicKey)  // First 32 bytes: public key
		copy(sig96[32:96], signature) // Next 64 bytes: signature
		signatures[i] = sig96
	}
	lggr.Debugf("  GasPriceUpdates: %d items", len(report.PriceUpdates.GasPriceUpdates))
	lggr.Debugf("  BlessedMerkleRoots: %d items", len(report.BlessedMerkleRoots))
	lggr.Debugf("  UnblessedMerkleRoots: %d items", len(report.UnblessedMerkleRoots))
	lggr.Debugf("  RMNSignatures: %d items", len(report.RMNSignatures))
	for i, sig := range report.RMNSignatures {
		lggr.Debugf("    [%d] Signature: %x (len=%d)", i, sig, len(sig))
	}
	lggr.Debugf("Serialized bytes: %x", bcsBytes)
	lggr.Debugf("Serialized length: %d", len(bcsBytes))
	lggr.Debug("==================")

	lggr.Debugw("Signatures", "signatures", signatures)

	ptbArgs := cwConfig.Arguments{
		Args: map[string]any{
			"ccip_object_ref": ccipObjectRef,
			"state":           offrampStateObjectId,
			"clock":           clockObject,
			"report_context":  reportContext,
			"report":          bcsBytes,
			"signatures":      signatures,
		},
	}

	lggr.Debugw("PTB args", "args", ptbArgs)

	// call commit
	txID := "ccip-commit-tx"
	err = chainWriter.SubmitTransaction(ctx,
		cwConfig.PTBChainWriterModuleName,
		"ccip_commit", // execute PTB
		&ptbArgs,
		txID,
		accountAddress,
		&commonTypes.TxMeta{GasLimit: big.NewInt(500000000)},
		nil,
	)
	require.NoError(t, err)
	lggr.Infow("Submitted transaction", "txID", txID)

	require.Eventually(t, func() bool {
		status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
		if statusErr != nil {
			return false
		}
		return status == commonTypes.Finalized
	}, 5*time.Second, 1*time.Second, "Transaction final state not reached")

	// Second transaction: execute flow (with expanded PTB)
	// Convert reportContext to the format expected by SuiOffRampExecCallArgs
	var reportContextArray [2][32]byte
	copy(reportContextArray[0][:], reportContext[0])
	copy(reportContextArray[1][:], reportContext[1])

	lggr.Debugw("linkTokenAddress", "linkTokenAddress", linkTokenAddress)

	// Use parameters from testMessage to ensure consistency with commit phase
	lggr.Debugw("Using testMessage parameters", "testMessage", testMessage)

	// Convert testMessage.MessageID to ccipocr3.Bytes32
	var messageIDBytes32 ccipocr3.Bytes32
	copy(messageIDBytes32[:], testMessage.MessageID[:])

	reportInfo := ccipocr3.ExecuteReportInfo{
		AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{
			{
				Messages: []ccipocr3.Message{
					{
						Header: ccipocr3.RampMessageHeader{
							MessageID:           messageIDBytes32,
							SourceChainSelector: ccipocr3.ChainSelector(testMessage.SourceChainSelector),
							DestChainSelector:   ccipocr3.ChainSelector(testMessage.DestChainSelector),
							SequenceNumber:      ccipocr3.SeqNum(testMessage.SequenceNumber),
							Nonce:               testMessage.Nonce,
						},
						TokenAmounts: testMessage.TokenAmounts,
						Receiver:     nil, //testMessage.Receiver[:], receiver is nil for now
						Data:         testMessage.Data,
					},
				},
			},
		},
	}

	// Add some distinguishing data
	// dummy report context - each element must be 32 bytes for Sui Move
	execReportContext := [][]byte{
		make([]byte, 32), // config digest - 32 bytes
		make([]byte, 32), // epoch and round - 32 bytes
	}
	execReportContext[0] = []byte{
		0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
		0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
		0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
		0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
	}
	execReportContext[1][0] = 0x022

	offChainTokenData := [][]byte{
		make([]byte, 32), // config digest - 32 bytes
		//make([]byte, 32), // epoch and round - 32 bytes
	}
	offChainTokenData[0] = []byte{
		0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
		0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
		0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
		0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
	}
	//offChainTokenData[1][0] = 0x022

	// Fix: Use empty proofs for single message tree
	// When proofs is empty, merkle_root(leaf, []) returns the leaf itself
	// This makes the message hash become the merkle root, matching our commit
	proofs := [][]byte{}

	lggr.Debugf("=== EXECUTION PHASE DEBUG ===")
	lggr.Debugf("Using empty proofs array (length: %d)", len(proofs))
	lggr.Debugf("Expected: message hash will become the merkle root")
	lggr.Debugf("Committed root was: %x", report.UnblessedMerkleRoots[0].MerkleRoot)
	lggr.Debugf("With empty proofs, contract will calculate: merkle_root(message_hash, []) = message_hash")
	lggr.Debugf("This should make committed_root == calculated_root")
	lggr.Debugf("=== PARAMETER CONSISTENCY ===")
	lggr.Debugf("Source Chain: %d", testMessage.SourceChainSelector)
	lggr.Debugf("Dest Chain: %d", testMessage.DestChainSelector)
	lggr.Debugf("Sequence Number: %d", testMessage.SequenceNumber)
	lggr.Debugf("Nonce: %d", testMessage.Nonce)
	lggr.Debugf("Gas Limit: %s", testMessage.GasLimit.String())
	lggr.Debugf("Token Amounts: %d tokens", len(testMessage.TokenAmounts))
	if len(testMessage.TokenAmounts) > 0 {
		lggr.Debugf("First token amount: %s", testMessage.TokenAmounts[0].Amount.Int.String())
	}
	lggr.Debugf("===============================")

	executeReport := testutils.GetExecutionReportFromCCIP(
		testMessage.SourceChainSelector, // Use consistent source chain selector
		reportInfo.AbstractReports[0].Messages[0],
		offChainTokenData,
		proofs,
		uint32(testMessage.GasLimit.Uint64()), // Use consistent gas limit
	)

	execReportBCSBytes, err := testutils.SerializeExecutionReport(executeReport)
	require.NoError(t, err, "failed to serialize execution report")

	// Convert execReportContext [][]byte to [2][32]byte
	var execReportContextArray [2][32]byte
	copy(execReportContextArray[0][:], execReportContext[0])
	copy(execReportContextArray[1][:], execReportContext[1])

	execReportArgs := expander.SuiOffRampExecCallArgs{
		ReportContext: execReportContextArray,
		Report:        execReportBCSBytes,
		Info:          reportInfo,
	}

	ptbArgsExecute := cwConfig.Arguments{
		Args: map[string]any{
			"ccip_object_ref": ccipObjectRef,
			"state":           offrampStateObjectId,
			"clock":           clockObject,
			"report_context":  execReportContext,
			"report":          execReportBCSBytes,
			"signatures":      signatures,
			"expanded_report": execReportArgs,
		},
	}

	lggr.Debugw("ptbArgsExecute", "ptbArgsExecute", ptbArgsExecute)

	txID = "ccip-execute-tx"
	err = chainWriter.SubmitTransaction(ctx,
		cwConfig.PTBChainWriterModuleName,
		cwConfig.CCIPExecuteReportFunctionName, // this will cause the PTB to be expanded
		&ptbArgsExecute,
		txID,
		envSettings.OffRampReport.Output.CCIPOffRampPackageId,
		&commonTypes.TxMeta{GasLimit: big.NewInt(500000000)},
		nil,
	)
	require.NoError(t, err)
	lggr.Infow("Submitted transaction", "txID", txID)

	require.Eventually(t, func() bool {
		status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
		if statusErr != nil {
			lggr.Errorw("Failed to get transaction status", "error", statusErr)
			return false
		}
		lggr.Debugw("Transaction status", "status", status)
		return status == commonTypes.Finalized
	}, 10*time.Second, 1*time.Second, "Transaction final state not reached")

	chainWriter.Close()
}

// TestCCIPOffRampWithReceiver tests the CCIP offramp functionality with a dummy receiver
func TestCCIPOffRampWithReceiver(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)

	envSettings := SetupTestEnvironment(t)

	// Ensure dummy receiver was deployed and registered
	require.NotNil(t, envSettings.DummyReceiverReport, "Dummy receiver should be deployed")
	require.NotEmpty(t, envSettings.DummyReceiverReport.Output.DummyReceiverPackageId, "Dummy receiver package ID should not be empty")
	require.NotEmpty(t, envSettings.DummyReceiverReport.Output.Objects.CCIPReceiverStateObjectId, "Dummy receiver state object ID should not be empty")

	lggr.Infow("Dummy receiver deployed",
		"packageId", envSettings.DummyReceiverReport.Output.DummyReceiverPackageId,
		"receiverStateId", envSettings.DummyReceiverReport.Output.Objects.CCIPReceiverStateObjectId,
	)

	// Create receiver address from the dummy receiver state object ID
	receiverPackageId := envSettings.DummyReceiverReport.Output.DummyReceiverPackageId

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()
	chainWriter, accountAddress := setupChainWriter(t, envSettings)

	err := chainWriter.Start(ctx)
	require.NoError(t, err)

	clockObject := "0x6"

	// Create test message with receiver
	receiverModule := "dummy_receiver"
	testMessage := createTestMessageWithReceiver(envSettings, receiverPackageId, receiverModule)

	// Create commit report for the message with receiver
	bcsBytes, report, reportContext := createCommitReportForMessage(t, envSettings, testMessage)

	ccipObjectRef := envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId
	require.NotEmpty(t, ccipObjectRef, "CCIP object reference should not be empty")
	offrampStateObjectId := envSettings.OffRampReport.Output.Objects.StateObjectId
	require.NotEmpty(t, offrampStateObjectId, "Offramp state object ID should not be empty")
	linkTokenAddress := envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId

	lggr.Infow("Test message with receiver",
		"sourceChain", testMessage.SourceChainSelector,
		"destChain", testMessage.DestChainSelector,
		"sequenceNumber", testMessage.SequenceNumber,
		"receiver", hex.EncodeToString(testMessage.Receiver),
		"data", hex.EncodeToString(testMessage.Data),
	)

	// Create proper signatures for OCR3 verification
	configDigest := reportContext[0]
	sequenceBytes := reportContext[1]

	// Hash the report the same way the contract does: report + config_digest + sequence_bytes
	var reportForHashing []byte
	reportForHashing = append(reportForHashing, bcsBytes...)
	reportForHashing = append(reportForHashing, configDigest...)
	reportForHashing = append(reportForHashing, sequenceBytes...)
	hashedReport := blake2b.Sum256(reportForHashing)

	// Create signatures using the first big_f + 1 = 2 private keys
	signatures := make([][]byte, 2)
	for i := 0; i < 2; i++ {
		privateKey := envSettings.PrivateKeys[i]
		publicKey := privateKey.Public().(ed25519.PublicKey)

		// Sign the hashed report with raw Ed25519 (no Sui intent bytes)
		signature := ed25519.Sign(privateKey, hashedReport[:])

		require.Equal(t, len(signature), 64)

		// Create 96-byte signature: 32 bytes pubkey + 64 bytes signature
		sig96 := make([]byte, 96)
		copy(sig96[0:32], publicKey)  // First 32 bytes: public key
		copy(sig96[32:96], signature) // Next 64 bytes: signature
		signatures[i] = sig96
	}

	lggr.Debugw("Signatures", "signatures", signatures)

	ptbArgs := cwConfig.Arguments{
		Args: map[string]any{
			"ccip_object_ref": ccipObjectRef,
			"state":           offrampStateObjectId,
			"clock":           clockObject,
			"report_context":  reportContext,
			"report":          bcsBytes,
			"signatures":      signatures,
		},
	}

	lggr.Debugw("PTB args", "args", ptbArgs)

	// call commit
	txID := "ccip-commit-with-receiver-tx"
	err = chainWriter.SubmitTransaction(ctx,
		cwConfig.PTBChainWriterModuleName,
		"ccip_commit", // execute PTB
		&ptbArgs,
		txID,
		accountAddress,
		&commonTypes.TxMeta{GasLimit: big.NewInt(500000000)},
		nil,
	)
	require.NoError(t, err)
	lggr.Infow("Submitted commit transaction", "txID", txID)

	require.Eventually(t, func() bool {
		status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
		if statusErr != nil {
			return false
		}
		return status == commonTypes.Finalized
	}, 5*time.Second, 1*time.Second, "Commit transaction final state not reached")

	// Second transaction: execute flow (with expanded PTB)
	// Convert reportContext to the format expected by SuiOffRampExecCallArgs
	var reportContextArray [2][32]byte
	copy(reportContextArray[0][:], reportContext[0])
	copy(reportContextArray[1][:], reportContext[1])

	lggr.Debugw("linkTokenAddress", "linkTokenAddress", linkTokenAddress)

	// Use parameters from testMessage to ensure consistency with commit phase
	lggr.Debugw("Using testMessage parameters", "testMessage", testMessage)

	// Convert testMessage.MessageID to ccipocr3.Bytes32
	var messageIDBytes32 ccipocr3.Bytes32
	copy(messageIDBytes32[:], testMessage.MessageID[:])

	reportInfo := ccipocr3.ExecuteReportInfo{
		AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{
			{
				Messages: []ccipocr3.Message{
					{
						Header: ccipocr3.RampMessageHeader{
							MessageID:           messageIDBytes32,
							SourceChainSelector: ccipocr3.ChainSelector(testMessage.SourceChainSelector),
							DestChainSelector:   ccipocr3.ChainSelector(testMessage.DestChainSelector),
							SequenceNumber:      ccipocr3.SeqNum(testMessage.SequenceNumber),
							Nonce:               testMessage.Nonce,
						},
						TokenAmounts: testMessage.TokenAmounts,
						Receiver:     testMessage.Receiver, // Use the receiver address
						Data:         testMessage.Data,
					},
				},
			},
		},
	}

	// Add some distinguishing data
	// dummy report context - each element must be 32 bytes for Sui Move
	execReportContext := [][]byte{
		make([]byte, 32), // config digest - 32 bytes
		make([]byte, 32), // epoch and round - 32 bytes
	}
	execReportContext[0] = []byte{
		0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
		0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
		0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
		0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
	}
	execReportContext[1][0] = 0x022

	offChainTokenData := [][]byte{
		make([]byte, 32), // config digest - 32 bytes
	}
	offChainTokenData[0] = []byte{
		0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
		0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
		0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
		0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
	}

	// Fix: Use empty proofs for single message tree
	// When proofs is empty, merkle_root(leaf, []) returns the leaf itself
	// This makes the message hash become the merkle root, matching our commit
	proofs := [][]byte{}

	lggr.Debugf("=== EXECUTION PHASE DEBUG ===")
	lggr.Debugf("Using empty proofs array (length: %d)", len(proofs))
	lggr.Debugf("Expected: message hash will become the merkle root")
	lggr.Debugf("Committed root was: %x", report.UnblessedMerkleRoots[0].MerkleRoot)
	lggr.Debugf("With empty proofs, contract will calculate: merkle_root(message_hash, []) = message_hash")
	lggr.Debugf("This should make committed_root == calculated_root")
	lggr.Debugf("=== PARAMETER CONSISTENCY ===")
	lggr.Debugf("Source Chain: %d", testMessage.SourceChainSelector)
	lggr.Debugf("Dest Chain: %d", testMessage.DestChainSelector)
	lggr.Debugf("Sequence Number: %d", testMessage.SequenceNumber)
	lggr.Debugf("Nonce: %d", testMessage.Nonce)
	lggr.Debugf("Gas Limit: %s", testMessage.GasLimit.String())
	lggr.Debugf("Receiver: %x", testMessage.Receiver)
	lggr.Debugf("Token Amounts: %d tokens", len(testMessage.TokenAmounts))
	if len(testMessage.TokenAmounts) > 0 {
		lggr.Debugf("First token amount: %s", testMessage.TokenAmounts[0].Amount.Int.String())
	}
	lggr.Debugf("===============================")

	executeReport := testutils.GetExecutionReportFromCCIP(
		testMessage.SourceChainSelector, // Use consistent source chain selector
		reportInfo.AbstractReports[0].Messages[0],
		offChainTokenData,
		proofs,
		uint32(testMessage.GasLimit.Uint64()), // Use consistent gas limit
	)

	execReportBCSBytes, err := testutils.SerializeExecutionReport(executeReport)
	require.NoError(t, err, "failed to serialize execution report")

	// Convert execReportContext [][]byte to [2][32]byte
	var execReportContextArray [2][32]byte
	copy(execReportContextArray[0][:], execReportContext[0])
	copy(execReportContextArray[1][:], execReportContext[1])

	execReportArgs := expander.SuiOffRampExecCallArgs{
		ReportContext: execReportContextArray,
		Report:        execReportBCSBytes,
		Info:          reportInfo,
	}

	ptbArgsExecute := cwConfig.Arguments{
		Args: map[string]any{
			"ccip_object_ref": ccipObjectRef,
			"state":           offrampStateObjectId,
			"clock":           clockObject,
			"report_context":  execReportContext,
			"report":          execReportBCSBytes,
			"signatures":      signatures,
			"expanded_report": execReportArgs,
		},
	}

	lggr.Debugw("ptbArgsExecute", "ptbArgsExecute", ptbArgsExecute)

	txID = "ccip-execute-with-receiver-tx"
	err = chainWriter.SubmitTransaction(ctx,
		cwConfig.PTBChainWriterModuleName,
		cwConfig.CCIPExecuteReportFunctionName, // this will cause the PTB to be expanded
		&ptbArgsExecute,
		txID,
		envSettings.OffRampReport.Output.CCIPOffRampPackageId,
		&commonTypes.TxMeta{GasLimit: big.NewInt(500000000)},
		nil,
	)
	require.NoError(t, err)
	lggr.Infow("Submitted execute transaction", "txID", txID)

	require.Eventually(t, func() bool {
		status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
		if statusErr != nil {
			lggr.Errorw("Failed to get transaction status", "error", statusErr)
			return false
		}
		lggr.Debugw("Transaction status", "status", status)
		return status == commonTypes.Finalized
	}, 10*time.Second, 1*time.Second, "Execute transaction final state not reached")

	// Read the emitted event from the dummy receiver
	lggr.Infow("Reading events from the execute transaction")

	// Create a PTB client to query events (the basic sui.ISuiAPI doesn't have QueryEvents)
	ptbClient, err := client.NewPTBClient(lggr, testutils.LocalUrl, nil, 10*time.Second, nil, 5, "WaitForLocalExecution")
	require.NoError(t, err, "Failed to create PTB client for event querying")

	receiverPackageId = envSettings.DummyReceiverReport.Output.DummyReceiverPackageId

	// Query for ReceivedMessage events emitted by the dummy receiver
	eventFilter := client.EventFilterByMoveEventModule{
		Package: receiverPackageId,
		Module:  "dummy_receiver",
		Event:   "ReceivedMessage",
	}

	// Query events with a small limit since we expect only one event
	limit := uint(10)
	eventsResponse, err := ptbClient.QueryEvents(ctx, eventFilter, &limit, nil, nil)
	lggr.Debugw("eventsResponse", "eventsResponse", eventsResponse)
	require.NoError(t, err, "Failed to query events")
	require.NotEmpty(t, eventsResponse.Data, "Expected at least one ReceivedMessage event")

	// Find the most recent event (should be ours)
	mostRecentEvent := eventsResponse.Data[0] // Events are typically returned in descending order

	lggr.Infow("Found ReceivedMessage event",
		"eventId", mostRecentEvent.Id,
		"packageId", mostRecentEvent.PackageId,
		"transactionModule", mostRecentEvent.TransactionModule,
		"sender", mostRecentEvent.Sender,
		"type", mostRecentEvent.Type,
		"parsedJson", mostRecentEvent.ParsedJson,
	)

	// Deserialize the event data into our struct
	var receivedMessage ReceivedMessageEvent
	parsedJson, err := json.Marshal(mostRecentEvent.ParsedJson)
	require.NoError(t, err, "Failed to marshal parsed JSON")
	err = json.Unmarshal(parsedJson, &receivedMessage)
	require.NoError(t, err, "Failed to deserialize ReceivedMessage event")

	lggr.Infow("Deserialized ReceivedMessage event",
		"messageId", hex.EncodeToString(receivedMessage.MessageID),
		"sender", hex.EncodeToString(receivedMessage.Sender),
		"data", string(receivedMessage.Data),
	)

	// Verify the event data matches our test message
	require.Equal(t, testMessage.Data, receivedMessage.Data, "Message data should match")

	chainWriter.Close()
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}
