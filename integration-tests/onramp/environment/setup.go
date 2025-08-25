package environment

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/holiman/uint256"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mockethtoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_eth_token"
	mocklinktoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_link_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	burnmintops "github.com/smartcontractkit/chainlink-sui/ops/ccip_burn_mint_token_pool"
	lockreleaseops "github.com/smartcontractkit/chainlink-sui/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/ops/ccip_managed_token_pool"
	onrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_onramp"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/ops/managed_token"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
	mockethtokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_eth_token"
	mocklinktokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_link_token"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/stretchr/testify/require"
)

// Constants used across the environment setup
const (
	EvmReceiverAddress = "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"
	EthereumAddress    = "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"
	ClockObjectId      = "0x6"
	DenyListObjectId   = "0x403"
)

// EnvironmentSettings holds all the deployed contract information and client settings
// needed for running CCIP onramp integration tests.
type EnvironmentSettings struct {
	AccountAddress string
	// Deployment reports
	MockLinkReport             *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]]
	MockEthTokenReport         *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mockethtokenops.DeployMockEthTokenObjects]]
	ManagedTokenReport         *cld_ops.SequenceReport[managedtokenops.DeployAndInitManagedTokenInput, managedtokenops.DeployManagedTokenOutput]
	CCIPReport                 *cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput]
	OnRampReport               *cld_ops.SequenceReport[onrampops.DeployAndInitCCIPOnRampSeqInput, onrampops.DeployCCIPOnRampSeqOutput]
	LockReleaseTokenPoolReport *cld_ops.SequenceReport[lockreleaseops.DeployAndInitLockReleaseTokenPoolInput, lockreleaseops.DeployLockReleaseTokenPoolOutput]
	BurnMintTokenPoolReport    *cld_ops.SequenceReport[burnmintops.DeployAndInitBurnMintTokenPoolInput, burnmintops.DeployBurnMintTokenPoolOutput]
	ManagedTokenPoolReport     *cld_ops.SequenceReport[managedtokenpoolops.SeqDeployAndInitManagedTokenPoolInput, managedtokenpoolops.DeployManagedTokenPoolOutput]

	EthereumPoolAddress []byte
	EthCoins            []string

	// Signers
	Signer rel.SuiSigner

	// Client
	Client sui.ISuiAPI
}

// SetupClients creates and configures Sui client and signer for testing.
// It generates a new key pair, creates a signer, and funds the signer address.
func SetupClients(t *testing.T, lggr logger.Logger) (rel.SuiSigner, sui.ISuiAPI) {
	t.Helper()

	client := sui.NewSuiClient(testutils.LocalUrl)

	// Generate key pair and create a signer.
	pk, _, _, err := testutils.GenerateAccountKeyPair(t)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)

	// Fund the signer for contract deployment
	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)
	for range 3 {
		err = testutils.FundWithFaucet(lggr, "localnet", signerAddress)
		require.NoError(t, err)
	}

	return signer, client
}

// BasicSetUp performs basic environment setup including account creation, client setup,
// and bundle initialization. This is the foundation for all test environments.
func BasicSetUp(t *testing.T, lggr logger.Logger, keystoreInstance *testutils.TestKeystore) (string, []byte, rel.SuiSigner, sui.ISuiAPI, sui_ops.OpTxDeps, cld_ops.Bundle) {
	t.Helper()

	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	signer, client := SetupClients(t, lggr)

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

	return accountAddress, publicKeyBytes, signer, client, deps, bundle
}

// UpdatePrices sets token prices in the fee quoter contract.
// This is critical for fee calculations in CCIP operations.
func UpdatePrices(
	t *testing.T,
	reportCCIP *cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	reportMockLink *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	deps sui_ops.OpTxDeps,
	bundle cld_ops.Bundle,
	destChainSelector uint64,
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
		CCIPPackageID:         reportCCIP.Output.CCIPPackageID,
		CCIPObjectRef:         reportCCIP.Output.Objects.CCIPObjectRefObjectID,
		FeeQuoterCapID:        reportCCIP.Output.Objects.FeeQuoterCapObjectID,
		SourceTokens:          []string{reportMockLink.Output.Objects.CoinMetadataObjectId},
		SourceUsdPerToken:     []*big.Int{linkTokenPrice},
		GasDestChainSelectors: []uint64{destChainSelector},
		GasUsdPerUnitGas:      []*big.Int{gasPrice},
	}

	_, err := cld_ops.ExecuteOperation(bundle, ccipops.FeeQuoterUpdateTokenPricesOp, deps, updatePricesInput)
	require.NoError(t, err, "failed to update token prices in fee quoter")

	lggr.Debugw("Updated token prices in fee quoter", "linkPrice", linkTokenPrice.String(), "gasPrice", gasPrice.String())
}

// DeployCCIPAndOnrampAndTokens deploys all the core CCIP infrastructure including
// mock tokens, MCMS contracts, CCIP core contracts, and onramp contracts.
func DeployCCIPAndOnrampAndTokens(
	t *testing.T,
	localChainSelector uint64,
	destChainSelector uint64,
	keystoreInstance *testutils.TestKeystore,
	signerAddr string,
	bundle cld_ops.Bundle,
	deps sui_ops.OpTxDeps,
	lggr logger.Logger,
) (
	*cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	*cld_ops.SequenceReport[onrampops.DeployAndInitCCIPOnRampSeqInput, onrampops.DeployCCIPOnRampSeqOutput],
	*cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	*cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mockethtokenops.DeployMockEthTokenObjects]],
	*cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
) {
	t.Helper()

	// Deploy LINK
	mockLinkReport, err := cld_ops.ExecuteOperation(bundle, mocklinktokenops.DeployMockLinkTokenOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy LINK token")

	// Deploy Mock ETH Token
	mockEthTokenReport, err := cld_ops.ExecuteOperation(bundle, mockethtokenops.DeployMockEthTokenOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy Mock ETH token")

	configDigest, err := uint256.FromHex("0xe3b1c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	require.NoError(t, err, "failed to convert config digest to uint256")

	// Deploy MCMs
	reportMCMs, err := cld_ops.ExecuteOperation(bundle, mcmsops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	require.NoError(t, err, "failed to deploy MCMS Package")
	lggr.Debugw("MCMS deployment report", "output", reportMCMs.Output)

	lggr.Debugw("LINK report", "output", mockLinkReport.Output)
	lggr.Debugw("Mock ETH token report", "output", mockEthTokenReport.Output)

	// Create 20-byte Ethereum addresses for RMN Remote signers
	ethAddr1, err := hex.DecodeString("8a1b2c3d4e5f60718293a4b5c6d7e8f901234567")
	require.NoError(t, err, "failed to decode eth address 1")
	ethAddr2, err := hex.DecodeString("7b8c9dab0c1d2e3f405162738495a6b7c8d9e0f1")
	require.NoError(t, err, "failed to decode eth address 2")
	ethAddr3, err := hex.DecodeString("1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err, "failed to decode eth address 3")

	reportCCIP, err := cld_ops.ExecuteSequence(bundle, ccipops.DeployAndInitCCIPSequence, deps, ccipops.DeployAndInitCCIPSeqInput{
		LinkTokenCoinMetadataObjectID: mockLinkReport.Output.Objects.CoinMetadataObjectId,
		LocalChainSelector:            localChainSelector,
		DestChainSelector:             destChainSelector,
		DeployCCIPInput: ccipops.DeployCCIPInput{
			McmsPackageID: reportMCMs.Output.PackageId,
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
	require.NotEmpty(t, reportCCIP.Output.CCIPPackageID, "CCIP package ID should not be empty")

	seqOnrampInput := onrampops.DeployAndInitCCIPOnRampSeqInput{
		DeployCCIPOnRampInput: onrampops.DeployCCIPOnRampInput{
			CCIPPackageID:      reportCCIP.Output.CCIPPackageID,
			MCMSPackageID:      reportMCMs.Output.PackageId,
			MCMSOwnerPackageID: signerAddr,
		},
		OnRampInitializeInput: onrampops.OnRampInitializeInput{
			NonceManagerCapID:         reportCCIP.Output.Objects.NonceManagerCapObjectID,   // this is from NonceManager init Op
			SourceTransferCapID:       reportCCIP.Output.Objects.SourceTransferCapObjectID, // this is from CCIP package publish
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

	return &reportCCIP, &reportOnRamp, &mockLinkReport, &mockEthTokenReport, &reportMCMs
}

// SetupTestEnvironment sets up a complete test environment with CCIP infrastructure
// and both lock/release and burn/mint token pools.
func SetupTestEnvironment(t *testing.T, localChainSelector uint64, destChainSelector uint64, keystoreInstance *testutils.TestKeystore) *EnvironmentSettings {
	t.Helper()

	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")

	accountAddress, _, signer, client, deps, bundle := BasicSetUp(t, lggr, keystoreInstance)
	signerAddr, err := signer.GetAddress()
	require.NoError(t, err)

	reportCCIP, reportOnRamp, reportMockLinkToken, reportMockEthToken, reportMCMs := DeployCCIPAndOnrampAndTokens(t, localChainSelector, destChainSelector, keystoreInstance, signerAddr, bundle, deps, lggr)

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", reportMockLinkToken.Output.PackageId)
	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", reportMockEthToken.Output.PackageId)

	ethereumPoolAddressString := string(NormalizeTo32Bytes(EvmReceiverAddress))
	remoteTokenAddressString := string(NormalizeTo32Bytes(EvmReceiverAddress))

	linkTokenPoolReport := SetupTokenPool(
		t,
		deps,
		reportCCIP,
		reportMCMs,
		reportMockLinkToken,
		signerAddr,
		accountAddress,
		linkTokenType,
		ethereumPoolAddressString,
		remoteTokenAddressString,
		destChainSelector,
		bundle,
		client,
		lggr,
	)

	ethCoins := GetEthCoins(t, client, signer, reportMockEthToken.Output.PackageId, reportMockEthToken.Output.Objects.TreasuryCapObjectId, ethTokenType, accountAddress, lggr, 1000000, 1000000)

	ethTokenPoolReport := SetupEthTokenPoolBurnMint(
		t,
		deps,
		reportCCIP,
		reportMCMs,
		reportMockEthToken,
		signerAddr,
		accountAddress,
		ethTokenType,
		ethereumPoolAddressString,
		remoteTokenAddressString,
		destChainSelector,
		bundle,
		client,
		lggr,
	)

	UpdatePrices(t, reportCCIP, reportMockLinkToken, deps, bundle, destChainSelector, lggr)

	return &EnvironmentSettings{
		AccountAddress:             accountAddress,
		MockLinkReport:             reportMockLinkToken,
		MockEthTokenReport:         reportMockEthToken,
		CCIPReport:                 reportCCIP,
		OnRampReport:               reportOnRamp,
		LockReleaseTokenPoolReport: linkTokenPoolReport,
		BurnMintTokenPoolReport:    ethTokenPoolReport,
		EthCoins:                   ethCoins,
		Signer:                     signer,
		Client:                     client,
	}
}

// SetupTestEnvironmentForManagedTokenPool sets up a test environment specifically
// for managed token pool testing.
func SetupTestEnvironmentForManagedTokenPool(t *testing.T, client sui.ISuiAPI, signer rel.SuiSigner, accountAddress string, bundle cld_ops.Bundle, deps sui_ops.OpTxDeps, localChainSelector uint64, destChainSelector uint64, keystoreInstance *testutils.TestKeystore) *EnvironmentSettings {
	t.Helper()

	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")
	signerAddr, err := signer.GetAddress()
	require.NoError(t, err)

	reportCCIP, reportOnRamp, reportMockLinkToken, reportMockEthToken, reportMCMs := DeployCCIPAndOnrampAndTokens(t, localChainSelector, destChainSelector, keystoreInstance, signerAddr, bundle, deps, lggr)

	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", reportMockEthToken.Output.PackageId)
	ethCoins := GetEthCoins(t, client, signer, reportMockEthToken.Output.PackageId, reportMockEthToken.Output.Objects.TreasuryCapObjectId, ethTokenType, accountAddress, lggr, 1000000, 1000000)

	ethereumPoolAddressString := string(NormalizeTo32Bytes(EvmReceiverAddress))
	remoteTokenAddressString := string(NormalizeTo32Bytes(EvmReceiverAddress))

	// Setup managed token pool for ETH token
	managedTokenPoolReport, managedTokenReport := SetupManagedTokenPool(
		t,
		deps,
		reportCCIP,
		reportMCMs,
		reportMockEthToken,
		signerAddr,
		accountAddress,
		ethTokenType,
		ethereumPoolAddressString,
		remoteTokenAddressString,
		destChainSelector,
		bundle,
		client,
		lggr,
	)

	UpdatePrices(t, reportCCIP, reportMockLinkToken, deps, bundle, destChainSelector, lggr)

	return &EnvironmentSettings{
		AccountAddress:         accountAddress,
		MockLinkReport:         reportMockLinkToken,
		MockEthTokenReport:     reportMockEthToken,
		CCIPReport:             reportCCIP,
		OnRampReport:           reportOnRamp,
		ManagedTokenPoolReport: managedTokenPoolReport,
		ManagedTokenReport:     managedTokenReport,
		EthCoins:               ethCoins,
		Signer:                 signer,
		Client:                 client,
	}
}

// GetLinkCoins mints LINK tokens for testing CCIP operations.
// Returns two coin IDs: one for token transfer and one for fee payment.
func GetLinkCoins(t *testing.T, envSettings *EnvironmentSettings, linkTokenType string, accountAddress string, lggr logger.Logger, tokenAmount uint64, feeAmount uint64) (string, string) {
	t.Helper()

	// Mint LINK tokens for the CCIP send operation
	// We need two separate coins: one for the token transfer and one for the fee payment

	// Use the setup account to mint tokens (since it owns the TreasuryCapObjectId)
	// but then transfer them to the transaction account
	deps := sui_ops.OpTxDeps{
		Client: envSettings.Client,
		Signer: envSettings.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(500_000_000)
			return &bind.CallOpts{
				Signer:           envSettings.Signer,
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	// Create LINK token contract instance
	linkContract, err := mocklinktoken.NewMockLinkToken(envSettings.MockLinkReport.Output.PackageId, envSettings.Client)
	require.NoError(t, err, "failed to create LINK token contract")

	// Use MintAndTransfer to mint directly to the transaction account
	// This avoids the ownership issue by minting directly to the account that will use the coins

	// Mint first coin for token transfer directly to transaction account
	mintTx1, err := linkContract.MockLinkToken().MintAndTransfer(
		context.Background(),
		deps.GetCallOpts(),
		bind.Object{Id: envSettings.MockLinkReport.Output.Objects.TreasuryCapObjectId},
		tokenAmount,
		accountAddress, // Mint directly to transaction account
	)
	require.NoError(t, err, "failed to mint and transfer LINK tokens for transfer")

	lggr.Debugw("Minted and transferred LINK tokens for transfer", "amount", tokenAmount, "txDigest", mintTx1.Digest, "recipient", accountAddress)

	// Find the first minted coin object ID from the transaction
	mintedCoinId1, err := bind.FindCoinObjectIdFromTx(*mintTx1, linkTokenType)
	require.NoError(t, err, "failed to find first minted coin object ID")
	lggr.Infow("First mintedCoinId", "coin", mintedCoinId1)

	// Mint second coin for fee payment directly to transaction account
	mintTx2, err := linkContract.MockLinkToken().MintAndTransfer(
		context.Background(),
		deps.GetCallOpts(),
		bind.Object{Id: envSettings.MockLinkReport.Output.Objects.TreasuryCapObjectId},
		feeAmount,
		accountAddress, // Mint directly to transaction account
	)
	require.NoError(t, err, "failed to mint and transfer LINK tokens for fee")

	lggr.Debugw("Minted and transferred LINK tokens for fee", "amount", feeAmount, "txDigest", mintTx2.Digest, "recipient", accountAddress)

	// Find the second minted coin object ID from the transaction
	mintedCoinId2, err := bind.FindCoinObjectIdFromTx(*mintTx2, linkTokenType)
	require.NoError(t, err, "failed to find second minted coin object ID")
	lggr.Infow("Second mintedCoinId", "coin", mintedCoinId2)

	return mintedCoinId1, mintedCoinId2
}

// GetEthCoins mints ETH tokens for testing CCIP operations.
// Returns an array of coin IDs for use in testing.
func GetEthCoins(t *testing.T, client sui.ISuiAPI, signer rel.SuiSigner, ethTokenPackageId string, treasuryCapObjectId string, ethTokenType string, accountAddress string, lggr logger.Logger, tokenAmount uint64, feeAmount uint64) []string {
	t.Helper()

	// Mint ETH tokens for the CCIP send operation
	// We need two separate coins: one for the token transfer and one for the fee payment

	// Use the setup account to mint tokens (since it owns the TreasuryCapObjectId)
	// but then transfer them to the transaction account
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

	// Create ETH token contract instance
	ethContract, err := mockethtoken.NewMockEthToken(ethTokenPackageId, client)
	require.NoError(t, err, "failed to create ETH token contract")

	// Use MintAndTransfer to mint directly to the transaction account
	// This avoids the ownership issue by minting directly to the account that will use the coins

	// Mint first coin for token transfer directly to transaction account
	mintTx1, err := ethContract.MockEthToken().MintAndTransfer(
		context.Background(),
		deps.GetCallOpts(),
		bind.Object{Id: treasuryCapObjectId},
		tokenAmount,
		accountAddress, // Mint directly to transaction account
	)
	require.NoError(t, err, "failed to mint and transfer ETH tokens for transfer")

	lggr.Debugw("Minted and transferred ETH tokens for transfer", "amount", tokenAmount, "txDigest", mintTx1.Digest, "recipient", accountAddress)

	// Find the first minted coin object ID from the transaction
	mintedCoinId1, err := bind.FindCoinObjectIdFromTx(*mintTx1, ethTokenType)
	require.NoError(t, err, "failed to find first minted coin object ID")
	lggr.Infow("First ETH mintedCoinId", "coin", mintedCoinId1)

	// Mint second coin for fee payment directly to transaction account
	mintTx2, err := ethContract.MockEthToken().MintAndTransfer(
		context.Background(),
		deps.GetCallOpts(),
		bind.Object{Id: treasuryCapObjectId},
		feeAmount,
		accountAddress, // Mint directly to transaction account
	)
	require.NoError(t, err, "failed to mint and transfer ETH tokens for fee")

	lggr.Debugw("Minted and transferred ETH tokens for fee", "amount", feeAmount, "txDigest", mintTx2.Digest, "recipient", accountAddress)

	// Find the second minted coin object ID from the transaction
	mintedCoinId2, err := bind.FindCoinObjectIdFromTx(*mintTx2, ethTokenType)
	require.NoError(t, err, "failed to find second minted coin object ID")
	lggr.Infow("Second ETH mintedCoinId", "coin", mintedCoinId2)

	return []string{mintedCoinId1, mintedCoinId2}
}

// NormalizeTo32Bytes converts an address string to a 32-byte representation.
// This is used for converting Ethereum addresses to the format expected by Sui contracts.
func NormalizeTo32Bytes(address string) []byte {
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
