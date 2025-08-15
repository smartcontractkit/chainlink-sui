package offramp_test

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/holiman/uint256"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/stretchr/testify/require"
)

const (
	evmReceiverAddress      = "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"
	SUI_CHAIN_SELECTOR      = 2
	ETHEREUM_CHAIN_SELECTOR = 1
)

var ConfigDigest = []byte{
	0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
	0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
	0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
	0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
}

func setupClients(t *testing.T, lggr logger.Logger) (rel.SuiSigner, sui.ISuiAPI, ed25519.PrivateKey) {
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
	for range 3 {
		err = testutils.FundWithFaucet(lggr, "localnet", signerAddress)
		require.NoError(t, err)
	}

	return signer, client, pk
}

func normalizeTo32Bytes(address string) []byte {
	addressHex := address
	if strings.HasPrefix(address, "0x") {
		addressHex = address[2:]
	}
	addressBytesFull, _ := hex.DecodeString(addressHex)
	addressBytes := addressBytesFull
	if len(addressBytesFull) > 32 {
		addressBytes = addressBytesFull[len(addressBytesFull)-32:]
	} else if len(addressBytesFull) < 32 {
		// pad left with zeros
		padding := make([]byte, 32-len(addressBytesFull))
		addressBytes = append(padding, addressBytesFull...)
	}
	return addressBytes
}

type EnvironmentSettings struct {
	AccountAddress string

	// Deployment reports
	MockLinkReport      cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]]
	CCIPReport          cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput]
	OffRampReport       cld_ops.SequenceReport[offrampops.DeployAndInitCCIPOffRampSeqInput, offrampops.DeployCCIPOffRampSeqOutput]
	TokenPoolReport     cld_ops.SequenceReport[lockreleaseops.DeployAndInitLockReleaseTokenPoolInput, lockreleaseops.DeployLockReleaseTokenPoolOutput]
	DummyReceiverReport cld_ops.SequenceReport[ccipops.DeployAndInitDummyReceiverInput, ccipops.DeployDummyReceiverSeqOutput]

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

func SetupOffRamp(t *testing.T,
	report cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	deps sui_ops.OpTxDeps,
	reportMCMs cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
	mockLinkReport cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	signerAddr string,
	accountAddress string,
	bundle cld_ops.Bundle,
	lggr logger.Logger,
	client sui.ISuiAPI,
	privateKey ed25519.PrivateKey,
) cld_ops.SequenceReport[offrampops.DeployAndInitCCIPOffRampSeqInput, offrampops.DeployCCIPOffRampSeqOutput] {
	t.Helper()
	lggr.Debugw("Setting up off ramp")

	// Get the main account's public key first
	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)
	accountAddressBytes, err := hex.DecodeString(strings.TrimPrefix(accountAddress, "0x"))
	require.NoError(t, err)

	// Create a dummy OnRamp address
	OnRampAddress := make([]byte, 32)
	OnRampAddress[31] = 20

	// Declare all arrays
	signerAddresses := make([]string, 0, 4)
	signerAddrBytes := make([][]byte, 0, 4)
	signerPublicKeys := make([][]byte, 0, 4)
	signerPrivateKeys := make([]ed25519.PrivateKey, 0, 4)

	// add 3 generated signers
	for range 3 {
		pk, _, _, err := testutils.GenerateAccountKeyPair(t)
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

	lggr.Infow("signer addresses", "signerAddresses", signerAddresses)

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
			ConfigDigest:                   ConfigDigest,
			OCRPluginType:                  byte(0),
			BigF:                           byte(1),
			IsSignatureVerificationEnabled: true,
			Signers:                        signerPublicKeys,
			Transmitters:                   signerAddresses,
		},
		ExecutionOCR3Config: offrampops.SetOCR3ConfigInput{
			ConfigDigest:                   ConfigDigest,
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

	return offrampReport
}

func SetupTokenPool(t *testing.T,
	report cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	deps sui_ops.OpTxDeps,
	reportMCMs cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
	mockLinkReport cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	signerAddr string,
	accountAddress string,
	linkTokenType string,
	ethereumPoolAddressString string,
	remoteTokenAddressString string,
	destChainSelector uint64,
	bundle cld_ops.Bundle,
	lggr logger.Logger,
	client sui.ISuiAPI,
) cld_ops.SequenceReport[lockreleaseops.DeployAndInitLockReleaseTokenPoolInput, lockreleaseops.DeployLockReleaseTokenPoolOutput] {
	t.Helper()

	lggr.Debugw("Setting up token pool")
	// Create a context for the operation
	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	// Deploy CCIP token pool
	ccipTokenPoolReport, err := cld_ops.ExecuteOperation(bundle, cciptokenpoolop.DeployCCIPTokenPoolOp, deps, cciptokenpoolop.TokenPoolDeployInput{
		CCIPPackageId:    report.Output.CCIPPackageId,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: accountAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Token Pool")

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

	return tokenPoolLockReleaseReport
}

func SetupDummyReceiver(t *testing.T,
	report cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	deps sui_ops.OpTxDeps,
	reportMCMs cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
	mockLinkReport cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	signerAddr string,
	accountAddress string,
	bundle cld_ops.Bundle,
	lggr logger.Logger,
) cld_ops.SequenceReport[ccipops.DeployAndInitDummyReceiverInput, ccipops.DeployDummyReceiverSeqOutput] {
	t.Helper()

	// Deploy and initialize the dummy receiver
	dummyReceiverReport, err := cld_ops.ExecuteSequence(bundle, ccipops.DeployAndInitDummyReceiverSequence, deps, ccipops.DeployAndInitDummyReceiverInput{
		DeployDummyReceiverInput: ccipops.DeployDummyReceiverInput{
			CCIPPackageId: report.Output.CCIPPackageId,
			McmsPackageId: reportMCMs.Output.PackageId,
			McmsOwner:     signerAddr,
		},
		CCIPObjectRefObjectId: report.Output.Objects.CCIPObjectRefObjectId,
		ReceiverStateParams:   []string{"0x6"}, // the clock object id
	})

	require.NoError(t, err, "failed to deploy and initialize dummy receiver")
	lggr.Debugw("Dummy receiver deployment report", "output", dummyReceiverReport.Output)

	return dummyReceiverReport
}

func SetupFeeQuoterPrices(t *testing.T,
	report cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	deps sui_ops.OpTxDeps,
	reportMCMs cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
	mockLinkReport cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	signerAddr string,
	destChainSelector uint64,
	bundle cld_ops.Bundle,
	lggr logger.Logger,
) {
	t.Helper()
	// **CRITICAL**: Set token prices in the fee quoter
	// The fee quoter needs to know USD prices to calculate fees
	// Set LINK token price to $5.00 USD (5 * 1e18 = 5e18)
	linkTokenPrice := big.NewInt(0)
	linkTokenPrice.SetString("5000000000000000000", 10) // $5.00 in 1e18 format

	// Set gas price for destination chain to 20 gwei (20 * 1e9 = 2e10)
	gasPrice := big.NewInt(20000000000) // 20 gwei in wei

	updatePricesInput := ccipops.FeeQuoterUpdateTokenPricesInput{
		CCIPPackageId:         report.Output.CCIPPackageId,
		CCIPObjectRef:         report.Output.Objects.CCIPObjectRefObjectId,
		FeeQuoterCapId:        report.Output.Objects.FeeQuoterCapObjectId,
		SourceTokens:          []string{mockLinkReport.Output.Objects.CoinMetadataObjectId},
		SourceUsdPerToken:     []*big.Int{linkTokenPrice},
		GasDestChainSelectors: []uint64{destChainSelector},
		GasUsdPerUnitGas:      []*big.Int{gasPrice},
	}

	_, err := cld_ops.ExecuteOperation(bundle, ccipops.FeeQuoterUpdateTokenPricesOp, deps, updatePricesInput)
	require.NoError(t, err, "failed to update token prices in fee quoter")

	lggr.Infow("Updated token prices in fee quoter", "linkPrice", linkTokenPrice.String(), "gasPrice", gasPrice.String())
}

func SetupTestEnvironment(t *testing.T, localChainSelector uint64, destChainSelector uint64, keystoreInstance *testutils.TestKeystore) *EnvironmentSettings {
	t.Helper()

	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")

	accountAddress, _ := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	signer, client, privateKey := setupClients(t, lggr)

	// Create 20-byte Ethereum addresses for RMN Remote signers
	ethAddr1, err := hex.DecodeString("8a1b2c3d4e5f60718293a4b5c6d7e8f901234567")
	require.NoError(t, err, "failed to decode eth address 1")
	ethAddr2, err := hex.DecodeString("7b8c9dab0c1d2e3f405162738495a6b7c8d9e0f1")
	require.NoError(t, err, "failed to decode eth address 2")
	ethAddr3, err := hex.DecodeString("1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err, "failed to decode eth address 3")

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
		LocalChainSelector:            localChainSelector,
		DestChainSelector:             destChainSelector,
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
		SignerOnchainPublicKeys:     [][]byte{ethAddr1, ethAddr2, ethAddr3},
		NodeIndexes:                 []uint64{0, 1, 2},
		FSign:                       uint64(1),
	})
	require.NoError(t, err, "failed to execute CCIP deploy sequence")
	require.NotEmpty(t, report.Output.CCIPPackageId, "CCIP package ID should not be empty")

	offrampReport := SetupOffRamp(t, report, deps, reportMCMs, mockLinkReport, signerAddr, accountAddress, bundle, lggr, client, privateKey)

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", mockLinkReport.Output.PackageId)

	ethereumPoolAddressString := string(normalizeTo32Bytes(evmReceiverAddress))
	remoteTokenAddressString := string(normalizeTo32Bytes(evmReceiverAddress))

	tokenPoolReport := SetupTokenPool(t, report, deps, reportMCMs, mockLinkReport,
		signerAddr, accountAddress, linkTokenType, ethereumPoolAddressString, remoteTokenAddressString,
		destChainSelector, bundle, lggr, client,
	)

	//SetupFeeQuoterPrices(t, report, deps, reportMCMs, mockLinkReport, signerAddr, destChainSelector, bundle, lggr)

	dummyReceiverReport := SetupDummyReceiver(t, report, deps, reportMCMs, mockLinkReport, signerAddr, accountAddress, bundle, lggr)

	return &EnvironmentSettings{
		AccountAddress:      accountAddress,
		MockLinkReport:      mockLinkReport,
		CCIPReport:          report,
		OffRampReport:       offrampReport,
		TokenPoolReport:     tokenPoolReport,
		DummyReceiverReport: dummyReceiverReport,
		Signer:              signer,
		Client:              client,
	}
}

func TestExecuteOffRamp(t *testing.T) {
	lggr := logger.Test(t)
	env := SetupTestEnvironment(t, ETHEREUM_CHAIN_SELECTOR, SUI_CHAIN_SELECTOR, testutils.NewTestKeystore(t))

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	lggr.Infow("Environment settings", "env", env)

	offrampPackageId := env.OffRampReport.Output.CCIPOffRampPackageId
	linkTokenPackageId := env.MockLinkReport.Output.Objects.CoinMetadataObjectId

	t.Run("TestTokenTransferWithArbitraryMessaging", func(t *testing.T) {
		lggr.Infow("Testing Token Transfer with Arbitrary Messaging")

		ptb := transaction.NewTransaction()
		ptb.SetSuiClient(env.Client.(*sui.Client))

		receiverPackageId := env.DummyReceiverReport.Output.DummyReceiverPackageId
		receiverModule := "dummy_receiver"
		receiver := fmt.Sprintf("%s::%s::ccip_receive", receiverPackageId, receiverModule)

		tokenAmount := ccipocr3.BigInt{Int: big.NewInt(300)}

		rawContent := "Do or do not, there is no try."
		msg := ccipocr3.Message{
			Receiver: []byte(receiver),
			Data:     []byte(rawContent),
		}

		lggr.Infow("Message", "msg", msg)

		ptbClient, err := client.NewPTBClient(lggr, testutils.LocalUrl, nil, 10*time.Second, nil, 5, "WaitForLocalExecution")
		require.NoError(t, err, "Failed to create PTB client for event querying")

		ctx := context.Background()

		addressMappings, err := offramp.GetOfframpAddressMappings(
			ctx,
			lggr,
			ptbClient,
			offrampPackageId,
			publicKeyBytes,
		)
		require.NoError(t, err, "failed to get offramp address mappings")
		lggr.Infow("Offramp address mappings", "addressMappings", addressMappings)

		hexEncodedLinkPackageId, err := hex.DecodeString(strings.Replace(linkTokenPackageId, "0x", "", 1))
		require.NoError(t, err, "failed to decode link token package id")

		var messageIDBytes32 ccipocr3.Bytes32

		execReport := ccipocr3.ExecuteReportInfo{
			AbstractReports: []ccipocr3.ExecutePluginReportSingleChain{
				{
					SourceChainSelector: ETHEREUM_CHAIN_SELECTOR,
					Messages: []ccipocr3.Message{
						{
							Header: ccipocr3.RampMessageHeader{
								MessageID:           messageIDBytes32,
								SourceChainSelector: ccipocr3.ChainSelector(ETHEREUM_CHAIN_SELECTOR),
								DestChainSelector:   ccipocr3.ChainSelector(SUI_CHAIN_SELECTOR),
								SequenceNumber:      ccipocr3.SeqNum(uint64(1)),
								Nonce:               uint64(0),
							},
							Receiver: []byte{}, // []byte(receiver),
							Data:     []byte(rawContent),
							TokenAmounts: []ccipocr3.RampTokenAmount{
								{
									DestTokenAddress: hexEncodedLinkPackageId,
									Amount:           tokenAmount,
								},
							},
						},
					},
				},
			},
		}

		offChainTokenData := [][]byte{
			make([]byte, 32), // config digest - 32 bytes
		}
		offChainTokenData[0] = []byte{
			0x00, 0x0A, 0x2F, 0x1F, 0x37, 0xB0, 0x33, 0xCC,
			0xC4, 0x42, 0x8A, 0xB6, 0x5C, 0x35, 0x39, 0xC9,
			0x31, 0x5D, 0xBF, 0x88, 0x2D, 0x4B, 0xAB, 0x13,
			0xF1, 0xE7, 0xEF, 0xE7, 0xB3, 0xDD, 0xDC, 0x36,
		}
		proofs := [][]byte{}

		report := testutils.GetExecutionReportFromCCIP(
			uint64(execReport.AbstractReports[0].SourceChainSelector),
			execReport.AbstractReports[0].Messages[0],
			offChainTokenData,
			proofs,
			uint32(1_000_000),
		)

		execReportBCSBytes, err := testutils.SerializeExecutionReport(report)
		require.NoError(t, err, "failed to serialize execution report")

		reportContext := [][]byte{
			make([]byte, 32), // config digest - 32 bytes
			make([]byte, 32), // epoch and round - 32 bytes
		}
		reportContext[0] = ConfigDigest
		reportContext[1][0] = 0x022

		args := cwConfig.Arguments{
			Args: map[string]interface{}{
				"ReportContext": reportContext,
				"Report":        execReportBCSBytes,
				"Info":          execReport,
			},
		}

		err = offramp.BuildOffRampExecutePTB(
			ctx,
			lggr,
			ptbClient,
			ptb,
			args,
			accountAddress,
			addressMappings,
		)
		require.NoError(t, err, "failed to build offramp execute PTB")
		lggr.Infow("Offramp execute PTB", "ptb", ptb)

		// Fund the account
		for range 3 {
			err = testutils.FundWithFaucet(lggr, "localnet", accountAddress)
			require.NoError(t, err)
		}

		_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr)
		txManager.Start(ctx)

		txID := "execute-offramp-test"
		txMetadata := &commontypes.TxMeta{}

		txManager.EnqueuePTB(ctx, txID, txMetadata, publicKeyBytes, ptb, false)

		require.Eventually(t, func() bool {
			status, statusErr := txManager.GetTransactionStatus(ctx, txID)
			if statusErr != nil {
				lggr.Errorw("Failed to get transaction status", "error", statusErr)
				return false
			}
			lggr.Debugw("Transaction status", "status", status)
			return status == commontypes.Finalized
		}, 10*time.Second, 1*time.Second, "Execute transaction final state not reached")
	})
}
