//go:build integration

package offramp_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/holiman/uint256"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mocklinktoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_link_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	lockreleaseops "github.com/smartcontractkit/chainlink-sui/ops/ccip_lock_release_token_pool"
	onrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_onramp"
	cciptokenpoolop "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
	mocklinktokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_link_token"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	ptbClient "github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/stretchr/testify/require"
)

const (
	evmReceiverAddress = "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"
)

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
	for range 3 {
		err = testutils.FundWithFaucet(lggr, "localnet", signerAddress)
		require.NoError(t, err)
	}

	return signer, client
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
	MockLinkReport  cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]]
	CCIPReport      cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput]
	OnnRampReport   cld_ops.SequenceReport[onrampops.DeployAndInitCCIPOnRampSeqInput, onrampops.DeployCCIPOnRampSeqOutput]
	TokenPoolReport cld_ops.SequenceReport[lockreleaseops.DeployAndInitLockReleaseTokenPoolInput, lockreleaseops.DeployLockReleaseTokenPoolOutput]

	EthereumPoolAddress []byte

	// Signers
	Signer rel.SuiSigner

	// Client
	Client sui.ISuiAPI
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
		RemoteChainSelectorsToAdd:    []uint64{destChainSelector},             // Destination chain selector
		RemotePoolAddressesToAdd:     [][]string{{ethereumPoolAddressString}}, // 32-byte remote pool address
		RemoteTokenAddressesToAdd:    []string{remoteTokenAddressString},      // 32-byte remote token address
		// Rate limiter configurations
		RemoteChainSelectors: []uint64{destChainSelector}, // Destination chain selector
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

func SetupTestEnvironment(t *testing.T, localChainSelector uint64, destChainSelector uint64, keystoreInstance *testutils.TestKeystore) *EnvironmentSettings {
	t.Helper()

	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")

	accountAddress, _ := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	signer, client := setupClients(t, lggr)

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

	seqOnrampInput := onrampops.DeployAndInitCCIPOnRampSeqInput{
		DeployCCIPOnRampInput: onrampops.DeployCCIPOnRampInput{
			CCIPPackageId:      report.Output.CCIPPackageId,
			MCMSPackageId:      reportMCMs.Output.PackageId,
			MCMSOwnerPackageId: signerAddr,
		},
		OnRampInitializeInput: onrampops.OnRampInitializeInput{
			NonceManagerCapId:         report.Output.Objects.NonceManagerCapObjectId,   // this is from NonceManager init Op
			SourceTransferCapId:       report.Output.Objects.SourceTransferCapObjectId, // this is from CCIP package publish
			ChainSelector:             destChainSelector,
			FeeAggregator:             signerAddr,
			AllowListAdmin:            signerAddr,
			DestChainSelectors:        []uint64{destChainSelector},
			DestChainEnabled:          []bool{true},
			DestChainAllowListEnabled: []bool{true},
		},
		ApplyDestChainConfigureOnRampInput: onrampops.ApplyDestChainConfigureOnRampInput{
			DestChainSelector:         []uint64{destChainSelector},
			DestChainEnabled:          []bool{true},
			DestChainAllowListEnabled: []bool{false},
		},
		ApplyAllowListUpdatesInput: onrampops.ApplyAllowListUpdatesInput{
			DestChainSelector:             []uint64{destChainSelector},
			DestChainAllowListEnabled:     []bool{false},
			DestChainAddAllowedSenders:    [][]string{{}},
			DestChainRemoveAllowedSenders: [][]string{{}},
		},
	}
	// Run onRamp deploy & Apply dest chain update sequence
	reportOnRamp, err := cld_ops.ExecuteSequence(bundle, onrampops.DeployAndInitCCIPOnRampSequence, deps, seqOnrampInput)
	require.NoError(t, err, "failed to execute CCIP OnRamp deploy sequence")

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", mockLinkReport.Output.PackageId)

	ethereumPoolAddressString := string(normalizeTo32Bytes(evmReceiverAddress))
	remoteTokenAddressString := string(normalizeTo32Bytes(evmReceiverAddress))

	tokenPoolReport := SetupTokenPool(t, report, deps, reportMCMs, mockLinkReport,
		signerAddr, accountAddress, linkTokenType, ethereumPoolAddressString, remoteTokenAddressString,
		destChainSelector, bundle, lggr, client,
	)

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

	_, err = cld_ops.ExecuteOperation(bundle, ccipops.FeeQuoterUpdateTokenPricesOp, deps, updatePricesInput)
	require.NoError(t, err, "failed to update token prices in fee quoter")

	lggr.Debugw("Updated token prices in fee quoter", "linkPrice", linkTokenPrice.String(), "gasPrice", gasPrice.String())

	return &EnvironmentSettings{
		AccountAddress:  accountAddress,
		MockLinkReport:  mockLinkReport,
		CCIPReport:      report,
		OnnRampReport:   reportOnRamp,
		TokenPoolReport: tokenPoolReport,
		Signer:          signer,
		Client:          client,
	}
}

func TestGetTokenPoolPTBConfig(t *testing.T) {
	lggr := logger.Test(t)
	sourceChainSelector := uint64(1)
	destChainSelector := uint64(2)

	// Create keystore and get account
	keystoreInstance := testutils.NewTestKeystore(t)

	t.Run("parse_param_type", func(t *testing.T) {
		jsonStr := `{"MutableReference": {"Struct": {"address": "0x2", "module": "object", "name": "UID"}}}`
		var param interface{}
		_ = json.Unmarshal([]byte(jsonStr), &param)

		paramType := offramp.ParseParamType(lggr, param)
		require.Equal(t, paramType, "object_id")
	})

	t.Run("decode_parameters", func(t *testing.T) {
		jsonStr := `
	{"increment_mult": {
		"isEntry": true,
		"parameters": [
			{
			"MutableReference": {
				"Struct": {
				"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
				"module": "counter",
				"name": "Counter",
				"typeArguments": []
				}
			}
			},
			"U64",
			"U64",
			{
			"MutableReference": {
				"Struct": {
				"address": "0x2",
				"module": "tx_context",
				"name": "TxContext",
				"typeArguments": []
				}
			}
			}
		],
		"return": [],
		"typeParameters": [],
		"visibility": "Public"
		}
	}
	`

		var function map[string]any
		err := json.Unmarshal([]byte(jsonStr), &function)
		require.NoError(t, err, "Failed to unmarshal JSON")
		require.NotNil(t, function, "Function map should not be nil")
		lggr.Debugw("Parsed function", "function", function)

		decodedParameters, err := offramp.DecodeParameters(lggr, function["increment_mult"].(map[string]any))
		require.NoError(t, err)
		require.Equal(t, len(decodedParameters), 4)

		for _, param := range decodedParameters {
			lggr.Debugw("Decoded parameter", "param", param)
		}
	})

	t.Run("lock_or_burn", func(t *testing.T) {
		envSettings := SetupTestEnvironment(t, sourceChainSelector, destChainSelector, keystoreInstance)

		accountAddress, _ := testutils.GetAccountAndKeyFromSui(keystoreInstance)
		lggr.Infow("Using account", "address", accountAddress)

		tokenPoolInfo := offramp.TokenPool{
			PackageId: envSettings.TokenPoolReport.Output.LockReleaseTPPackageID,
			ModuleId:  "lock_release_token_pool",
			Function:  "lock_or_burn",
		}

		client, err := ptbClient.NewPTBClient(lggr, testutils.LocalUrl, nil, 10*time.Second, nil, 100, "WaitForLocalExecution")
		require.NoError(t, err, "failed to create PTB client")

		ptbConfig, err := offramp.GetTokenPoolPTBConfig(context.Background(), lggr, client, tokenPoolInfo)
		require.NoError(t, err, "failed to get token pool PTB config")
		lggr.Debugw("PTB config", "ptbConfig", *ptbConfig)

		require.Equal(t, ptbConfig.Type, codec.SuiPTBCommandMoveCall)
		require.Equal(t, ptbConfig.PackageId, offramp.AnyPointer(tokenPoolInfo.PackageId))
		require.Equal(t, ptbConfig.ModuleId, offramp.AnyPointer(tokenPoolInfo.ModuleId))
		require.Equal(t, ptbConfig.Function, offramp.AnyPointer(tokenPoolInfo.Function))
		require.Equal(t, len(ptbConfig.Params), 5)
		params := ptbConfig.Params

		ccipObjectRef := params[0]
		coin := params[1]
		remoteChainSelector := params[2]
		clock := params[3]
		state := params[4]

		require.Equal(t, ccipObjectRef.Name, "token_pool_0_state_object_CCIPObjectRef")
		require.Equal(t, ccipObjectRef.Type, "object_id")
		require.Equal(t, coin.Name, "token_pool_0_coin_Coin")
		require.Equal(t, coin.Type, "object_id")
		require.Equal(t, remoteChainSelector.Name, "token_pool_0_U64")
		require.Equal(t, remoteChainSelector.Type, "int64")
		require.Equal(t, clock.Name, "token_pool_0_clock_Clock")
		require.Equal(t, clock.Type, "object_id")
		require.Equal(t, state.Name, "token_pool_0_lock_release_token_pool_LockReleaseTokenPoolState")
		require.Equal(t, state.Type, "object_id")

		for _, param := range params {
			// check that hot potato does not exist
			require.NotEqual(t, param.Name, "hot_potato")
			require.NotEqual(t, param.Name, "tx_context")
		}
	})
}
