//go:build integration

package ccip_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/holiman/uint256"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	managedtoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/managed_token"
	mockethtoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_eth_token"
	mocklinktoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_link_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/ops/ccip"
	burnmintops "github.com/smartcontractkit/chainlink-sui/ops/ccip_burn_mint_token_pool"
	lockreleaseops "github.com/smartcontractkit/chainlink-sui/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/ops/ccip_managed_token_pool"
	onrampops "github.com/smartcontractkit/chainlink-sui/ops/ccip_onramp"
	cciptokenpoolop "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/ops/managed_token"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
	mockethtokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_eth_token"
	mocklinktokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_link_token"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/stretchr/testify/require"
)

type ContractAddresses struct {
	CCIPPackageID              string
	CCIPOnrampPackageID        string
	LinkLockReleaseTokenPool   string
	CCIPTokenPoolPackageID     string
	CCIPTokenPoolStateObjectId string
}

const (
	evmReceiverAddress = "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"
	ethereumAddress    = "0x80226fc0ee2b096224eeac085bb9a8cba1146f7d"
	clockObjectId      = "0x6"
	denyListObjectId   = "0x403"
)

func setupClients(t *testing.T, lggr logger.Logger) (rel.SuiSigner, sui.ISuiAPI) {
	t.Helper()

	client := sui.NewSuiClient(testutils.LocalUrl)

	// Generate key pair and create a signer.
	pk, _, _, err := testutils.GenerateAccountKeyPair(t)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)

	// Fund the signer for contract deployment
	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)
	for range 10 {
		err = testutils.FundWithFaucet(lggr, "localnet", signerAddress)
		require.NoError(t, err)
	}

	return signer, client
}

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

func SetupEthTokenPoolBurnMint(
	t *testing.T,
	deps sui_ops.OpTxDeps,
	reportCCIP *cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	reportMCMs *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
	reportMockEthToken *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mockethtokenops.DeployMockEthTokenObjects]],
	signerAddr string,
	accountAddress string,
	tokenType string,
	ethereumPoolAddressString string,
	remoteTokenAddressString string,
	destChainSelector uint64,
	bundle cld_ops.Bundle,
	client sui.ISuiAPI,
	lggr logger.Logger,
) *cld_ops.SequenceReport[burnmintops.DeployAndInitBurnMintTokenPoolInput, burnmintops.DeployBurnMintTokenPoolOutput] {
	t.Helper()

	lggr.Debugw("Setting up ETH burn-mint token pool")

	// Deploy CCIP token pool
	ccipTokenPoolReport, err := cld_ops.ExecuteOperation(bundle, cciptokenpoolop.DeployCCIPTokenPoolOp, deps, cciptokenpoolop.TokenPoolDeployInput{
		CCIPPackageId:    reportCCIP.Output.CCIPPackageId,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: accountAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Token Pool")

	// Deploy and initialize the burn mint token pool
	seqBurnMintDeployInput := burnmintops.DeployAndInitBurnMintTokenPoolInput{
		BurnMintTokenPoolDeployInput: burnmintops.BurnMintTokenPoolDeployInput{
			CCIPPackageId:          reportCCIP.Output.CCIPPackageId,
			CCIPTokenPoolPackageId: ccipTokenPoolReport.Output.PackageId,
			MCMSAddress:            reportMCMs.Output.PackageId,
			MCMSOwnerAddress:       accountAddress,
		},
		// Initialization parameters
		CoinObjectTypeArg:      tokenType,
		CCIPObjectRefObjectId:  reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		CoinMetadataObjectId:   reportMockEthToken.Output.Objects.CoinMetadataObjectId,
		TreasuryCapObjectId:    reportMockEthToken.Output.Objects.TreasuryCapObjectId,
		TokenPoolAdministrator: accountAddress,

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

	tokenPoolBurnMintReport, err := cld_ops.ExecuteSequence(bundle, burnmintops.DeployAndInitBurnMintTokenPoolSequence, deps, seqBurnMintDeployInput)
	require.NoError(t, err, "failed to deploy and initialize Burn Mint Token Pool")

	lggr.Debugw("ETH Token Pool Burn Mint deployment report", "output", tokenPoolBurnMintReport.Output)

	// Note: Burn mint pools don't need liquidity provision like lock-release pools
	// because they mint/burn tokens on demand rather than locking them

	return &tokenPoolBurnMintReport
}

func SetupManagedTokenPool(
	t *testing.T,
	deps sui_ops.OpTxDeps,
	reportCCIP *cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	reportMCMs *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
	reportMockEthToken *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mockethtokenops.DeployMockEthTokenObjects]],
	signerAddr string,
	accountAddress string,
	tokenType string,
	ethereumPoolAddressString string,
	remoteTokenAddressString string,
	destChainSelector uint64,
	bundle cld_ops.Bundle,
	client sui.ISuiAPI,
	lggr logger.Logger,
) (
	*cld_ops.SequenceReport[managedtokenpoolops.SeqDeployAndInitManagedTokenPoolInput, managedtokenpoolops.DeployManagedTokenPoolOutput],
	*cld_ops.SequenceReport[managedtokenops.DeployAndInitManagedTokenInput, managedtokenops.DeployManagedTokenOutput],
) {
	t.Helper()

	lggr.Debugw("Setting up managed token pool")

	// First, deploy CCIP token pool
	ccipTokenPoolReport, err := cld_ops.ExecuteOperation(bundle, cciptokenpoolop.DeployCCIPTokenPoolOp, deps, cciptokenpoolop.TokenPoolDeployInput{
		CCIPPackageId:    reportCCIP.Output.CCIPPackageId,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: accountAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Token Pool")

	// Deploy and initialize the managed token first
	seqManagedTokenDeployInput := managedtokenops.DeployAndInitManagedTokenInput{
		ManagedTokenDeployInput: managedtokenops.ManagedTokenDeployInput{
			MCMSAddress:      reportMCMs.Output.PackageId,
			MCMSOwnerAddress: accountAddress,
		},
		// Initialization parameters
		CoinObjectTypeArg:   tokenType,
		TreasuryCapObjectId: reportMockEthToken.Output.Objects.TreasuryCapObjectId,
		DenyCapObjectId:     "", // Optional - not using deny cap for this example
		// Configure a new minter
		MinterAddress: signerAddr,
		Allowance:     1000000, // 1M tokens allowance
		IsUnlimited:   false,
	}

	managedTokenReport, err := cld_ops.ExecuteSequence(bundle, managedtokenops.DeployAndInitManagedTokenSequence, deps, seqManagedTokenDeployInput)
	require.NoError(t, err, "failed to deploy and initialize Managed Token")

	lggr.Debugw("Managed Token deployment report", "output", managedTokenReport.Output)

	mintCapObjectId := configureNewMinter(
		t,
		client,
		deps.Signer,
		managedTokenReport.Output.ManagedTokenPackageId,
		tokenType,
		managedTokenReport.Output.Objects.StateObjectId,
		managedTokenReport.Output.Objects.OwnerCapObjectId,
		signerAddr,
		0,
		true,
		lggr,
	)

	// Now deploy and initialize the managed token pool
	seqManagedTokenPoolDeployInput := managedtokenpoolops.SeqDeployAndInitManagedTokenPoolInput{
		// Deploy inputs
		CCIPPackageId:          reportCCIP.Output.CCIPPackageId,
		CCIPTokenPoolPackageId: ccipTokenPoolReport.Output.PackageId,
		ManagedTokenPackageId:  managedTokenReport.Output.ManagedTokenPackageId,
		MCMSAddress:            reportMCMs.Output.PackageId,
		MCMSOwnerAddress:       accountAddress,
		// Initialize inputs
		CoinObjectTypeArg:         tokenType,
		CCIPObjectRefObjectId:     reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		ManagedTokenStateObjectId: managedTokenReport.Output.Objects.StateObjectId,
		ManagedTokenOwnerCapId:    managedTokenReport.Output.Objects.OwnerCapObjectId,
		CoinMetadataObjectId:      reportMockEthToken.Output.Objects.CoinMetadataObjectId,
		MintCapObjectId:           mintCapObjectId,
		TokenPoolAdministrator:    accountAddress,
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

	managedTokenPoolReport, err := cld_ops.ExecuteSequence(bundle, managedtokenpoolops.DeployAndInitManagedTokenPoolSequence, deps, seqManagedTokenPoolDeployInput)
	require.NoError(t, err, "failed to deploy and initialize Managed Token Pool")

	lggr.Debugw("Managed Token Pool deployment report", "output", managedTokenPoolReport.Output)

	return &managedTokenPoolReport, &managedTokenReport
}

// ConfigureManagedTokenMinter configures a new minter for a managed token contract.
// This function calls the configure_new_minter operation on the managed token,
// which allows the specified address to mint tokens up to the given allowance.
//
// Parameters:
// - managedTokenPackageId: The package ID of the deployed managed token
// - tokenType: The fully qualified token type (e.g., "package_id::token::TOKEN_TYPE")
// - stateObjectId: The state object ID of the managed token
// - ownerCapObjectId: The owner capability object ID for authorization
// - minterAddress: The address that will be granted minting permissions
// - allowance: Maximum number of tokens this minter can mint
// - isUnlimited: If true, the minter has unlimited minting capability
func ConfigureManagedTokenMinter(
	t *testing.T,
	deps sui_ops.OpTxDeps,
	managedTokenPackageId string,
	tokenType string,
	stateObjectId string,
	ownerCapObjectId string,
	minterAddress string,
	allowance uint64,
	isUnlimited bool,
	bundle cld_ops.Bundle,
	lggr logger.Logger,
) {
	t.Helper()

	lggr.Debugw("Configuring managed token minter",
		"packageId", managedTokenPackageId,
		"minterAddress", minterAddress,
		"allowance", allowance,
		"isUnlimited", isUnlimited)

	configureInput := managedtokenops.ManagedTokenConfigureNewMinterInput{
		ManagedTokenPackageId: managedTokenPackageId,
		CoinObjectTypeArg:     tokenType,
		StateObjectId:         stateObjectId,
		OwnerCapObjectId:      ownerCapObjectId,
		MinterAddress:         minterAddress,
		Allowance:             allowance,
		IsUnlimited:           isUnlimited,
	}

	_, err := cld_ops.ExecuteOperation(bundle, managedtokenops.ManagedTokenConfigureNewMinterOp, deps, configureInput)
	require.NoError(t, err, "failed to configure new minter for managed token")

	lggr.Debugw("Successfully configured managed token minter")
}

// configureNewMinter configures a new minter for a managed token.
// This function follows the pattern of getEthCoins/getLinkCoins by handling all the complexity internally.
// It configures minting permissions for the specified address and returns the mint cap object ID.
//
// Returns the mint cap object ID that was transferred to the minter address.
func configureNewMinter(
	t *testing.T,
	client sui.ISuiAPI,
	signer rel.SuiSigner,
	managedTokenPackageId string,
	tokenType string,
	stateObjectId string,
	ownerCapObjectId string,
	minterAddress string,
	allowance uint64,
	isUnlimited bool,
	lggr logger.Logger,
) string {
	t.Helper()

	lggr.Debugw("Configuring managed token minter",
		"packageId", managedTokenPackageId,
		"minterAddress", minterAddress,
		"allowance", allowance,
		"isUnlimited", isUnlimited)

	// Set up dependencies locally
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

	// Create managed token contract instance
	managedTokenContract, err := managedtoken.NewCCIPManagedToken(managedTokenPackageId, client)
	require.NoError(t, err, "failed to create managed token contract")

	// Call ConfigureNewMinter directly on the contract
	configureTx, err := managedTokenContract.ManagedToken().ConfigureNewMinter(
		context.Background(),
		deps.GetCallOpts(),
		[]string{tokenType},
		bind.Object{Id: stateObjectId},
		bind.Object{Id: ownerCapObjectId},
		minterAddress,
		allowance,
		isUnlimited,
	)
	require.NoError(t, err, "failed to configure new minter for managed token")

	// Find the mint cap object ID that was transferred to the minter
	mintCapObjectId, err := bind.FindObjectIdFromPublishTx(*configureTx, "managed_token", "MintCap")
	require.NoError(t, err, "failed to find mint cap object ID from configure minter transaction")

	lggr.Infow("Successfully configured managed token minter",
		"minter", minterAddress,
		"allowance", allowance,
		"isUnlimited", isUnlimited,
		"mintCapObjectId", mintCapObjectId)

	return mintCapObjectId
}

func SetupTokenPool(
	t *testing.T,
	deps sui_ops.OpTxDeps,
	reportCCIP *cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	reportMCMs *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mcmsops.DeployMCMSObjects]],
	mockLinkReport *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	signerAddr string,
	accountAddress string,
	linkTokenType string,
	ethereumPoolAddressString string,
	remoteTokenAddressString string,
	destChainSelector uint64,
	bundle cld_ops.Bundle,
	client sui.ISuiAPI,
	lggr logger.Logger,
) *cld_ops.SequenceReport[lockreleaseops.DeployAndInitLockReleaseTokenPoolInput, lockreleaseops.DeployLockReleaseTokenPoolOutput] {
	t.Helper()

	lggr.Debugw("Setting up token pool")
	// Create a context for the operation
	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	// Deploy CCIP token pool
	ccipTokenPoolReport, err := cld_ops.ExecuteOperation(bundle, cciptokenpoolop.DeployCCIPTokenPoolOp, deps, cciptokenpoolop.TokenPoolDeployInput{
		CCIPPackageId:    reportCCIP.Output.CCIPPackageId,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: accountAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Token Pool")

	// Deploy and initialize the lock release token pool
	seqLockReleaseDeployInput := lockreleaseops.DeployAndInitLockReleaseTokenPoolInput{
		LockReleaseTokenPoolDeployInput: lockreleaseops.LockReleaseTokenPoolDeployInput{
			CCIPPackageId:          reportCCIP.Output.CCIPPackageId,
			CCIPTokenPoolPackageId: ccipTokenPoolReport.Output.PackageId,
			MCMSAddress:            reportMCMs.Output.PackageId,
			MCMSOwnerAddress:       accountAddress,
		},
		// Initialization parameters
		CoinObjectTypeArg:      linkTokenType,
		CCIPObjectRefObjectId:  reportCCIP.Output.Objects.CCIPObjectRefObjectId,
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

	return &tokenPoolLockReleaseReport
}

func SetupTestEnvironment(t *testing.T, localChainSelector uint64, destChainSelector uint64, keystoreInstance *testutils.TestKeystore) *EnvironmentSettings {
	t.Helper()

	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")

	accountAddress, _, signer, client, deps, bundle := basicSetUp(t, lggr, keystoreInstance)
	signerAddr, err := signer.GetAddress()
	require.NoError(t, err)

	reportCCIP, reportOnRamp, reportMockLinkToken, reportMockEthToken, reportMCMs := deployCCIPAndOnrampAndTokens(t, localChainSelector, destChainSelector, keystoreInstance, signerAddr, bundle, deps, lggr)

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", reportMockLinkToken.Output.PackageId)
	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", reportMockEthToken.Output.PackageId)

	ethereumPoolAddressString := string(normalizeTo32Bytes(evmReceiverAddress))
	remoteTokenAddressString := string(normalizeTo32Bytes(evmReceiverAddress))

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

	ethCoins := getEthCoins(t, client, signer, reportMockEthToken.Output.PackageId, reportMockEthToken.Output.Objects.TreasuryCapObjectId, ethTokenType, accountAddress, lggr, 1000000, 1000000)

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

	updatePrices(t, reportCCIP, reportMockLinkToken, deps, bundle, destChainSelector, lggr)

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

func basicSetUp(t *testing.T, lggr logger.Logger, keystoreInstance *testutils.TestKeystore) (string, []byte, rel.SuiSigner, sui.ISuiAPI, sui_ops.OpTxDeps, cld_ops.Bundle) {
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	signer, client := setupClients(t, lggr)

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

func updatePrices(
	t *testing.T,
	reportCCIP *cld_ops.SequenceReport[ccipops.DeployAndInitCCIPSeqInput, ccipops.DeployCCIPSeqOutput],
	reportMockLink *cld_ops.Report[cld_ops.EmptyInput, sui_ops.OpTxResult[mocklinktokenops.DeployMockLinkTokenObjects]],
	deps sui_ops.OpTxDeps,
	bundle cld_ops.Bundle,
	destChainSelector uint64,
	lggr logger.Logger,
) {
	// **CRITICAL**: Set token prices in the fee quoter
	// The fee quoter needs to know USD prices to calculate fees
	// Set LINK token price to $5.00 USD (5 * 1e18 = 5e18)
	linkTokenPrice := big.NewInt(0)
	linkTokenPrice.SetString("5000000000000000000", 10) // $5.00 in 1e18 format

	// Set gas price for destination chain to 20 gwei (20 * 1e9 = 2e10)
	gasPrice := big.NewInt(20000000000) // 20 gwei in wei

	updatePricesInput := ccipops.FeeQuoterUpdateTokenPricesInput{
		CCIPPackageId:         reportCCIP.Output.CCIPPackageId,
		CCIPObjectRef:         reportCCIP.Output.Objects.CCIPObjectRefObjectId,
		FeeQuoterCapId:        reportCCIP.Output.Objects.FeeQuoterCapObjectId,
		SourceTokens:          []string{reportMockLink.Output.Objects.CoinMetadataObjectId},
		SourceUsdPerToken:     []*big.Int{linkTokenPrice},
		GasDestChainSelectors: []uint64{destChainSelector},
		GasUsdPerUnitGas:      []*big.Int{gasPrice},
	}

	_, err := cld_ops.ExecuteOperation(bundle, ccipops.FeeQuoterUpdateTokenPricesOp, deps, updatePricesInput)
	require.NoError(t, err, "failed to update token prices in fee quoter")

	lggr.Debugw("Updated token prices in fee quoter", "linkPrice", linkTokenPrice.String(), "gasPrice", gasPrice.String())
}

func deployCCIPAndOnrampAndTokens(
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
	require.NotEmpty(t, reportCCIP.Output.CCIPPackageId, "CCIP package ID should not be empty")

	seqOnrampInput := onrampops.DeployAndInitCCIPOnRampSeqInput{
		DeployCCIPOnRampInput: onrampops.DeployCCIPOnRampInput{
			CCIPPackageId:      reportCCIP.Output.CCIPPackageId,
			MCMSPackageId:      reportMCMs.Output.PackageId,
			MCMSOwnerPackageId: signerAddr,
		},
		OnRampInitializeInput: onrampops.OnRampInitializeInput{
			NonceManagerCapId:         reportCCIP.Output.Objects.NonceManagerCapObjectId,   // this is from NonceManager init Op
			SourceTransferCapId:       reportCCIP.Output.Objects.SourceTransferCapObjectId, // this is from CCIP package publish
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

func SetupTestEnvironmentForManagedTokenPool(t *testing.T, client sui.ISuiAPI, signer rel.SuiSigner, accountAddress string, bundle cld_ops.Bundle, deps sui_ops.OpTxDeps, localChainSelector uint64, destChainSelector uint64, keystoreInstance *testutils.TestKeystore) *EnvironmentSettings {
	t.Helper()

	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")
	signerAddr, err := signer.GetAddress()
	require.NoError(t, err)

	reportCCIP, reportOnRamp, reportMockLinkToken, reportMockEthToken, reportMCMs := deployCCIPAndOnrampAndTokens(t, localChainSelector, destChainSelector, keystoreInstance, signerAddr, bundle, deps, lggr)

	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", reportMockEthToken.Output.PackageId)
	ethCoins := getEthCoins(t, client, signer, reportMockEthToken.Output.PackageId, reportMockEthToken.Output.Objects.TreasuryCapObjectId, ethTokenType, accountAddress, lggr, 1000000, 1000000)

	ethereumPoolAddressString := string(normalizeTo32Bytes(evmReceiverAddress))
	remoteTokenAddressString := string(normalizeTo32Bytes(evmReceiverAddress))

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

	updatePrices(t, reportCCIP, reportMockLinkToken, deps, bundle, destChainSelector, lggr)

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

func getLinkCoins(t *testing.T, envSettings *EnvironmentSettings, linkTokenType string, accountAddress string, lggr logger.Logger, tokenAmount uint64, feeAmount uint64) (string, string) {
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

func getEthCoins(t *testing.T, client sui.ISuiAPI, signer rel.SuiSigner, ethTokenPackageId string, treasuryCapObjectId string, ethTokenType string, accountAddress string, lggr logger.Logger, tokenAmount uint64, feeAmount uint64) []string {
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

// TestCCIPSuiOnRamp tests the CCIP onramp send functionality
func TestCCIPSuiOnRamp(t *testing.T) {
	lggr := logger.Test(t)

	localChainSelector := uint64(1)
	destChainSelector := uint64(2)

	// Create keystore and get account
	keystoreInstance := testutils.NewTestKeystore(t)

	// Start dedicated Sui node for this test
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			if perr := cmd.Process.Kill(); perr != nil {
				t.Logf("Failed to kill Sui node process: %v", perr)
			}
		}
	})

	// Wait for the node to be fully ready
	time.Sleep(3 * time.Second)

	envSettings := SetupTestEnvironment(t, localChainSelector, destChainSelector, keystoreInstance)

	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)
	lggr.Infow("Using account", "address", accountAddress)

	// Fund the account for gas payments
	for range 10 {
		err := testutils.FundWithFaucet(lggr, "localnet", accountAddress)
		require.NoError(t, err)
	}

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", envSettings.MockLinkReport.Output.PackageId)
	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", envSettings.MockEthTokenReport.Output.PackageId)

	_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr)

	tokenPoolDetails := testutils.TokenToolDetails{
		TokenPoolPackageId: envSettings.LockReleaseTokenPoolReport.Output.LockReleaseTPPackageID,
		TokenPoolType:      testutils.TokenPoolTypeLockRelease,
	}
	ethTokenPoolDetails := testutils.TokenToolDetails{
		TokenPoolPackageId: envSettings.BurnMintTokenPoolReport.Output.BurnMintTPPackageID,
		TokenPoolType:      testutils.TokenPoolTypeBurnMint,
	}

	chainWriterConfig, err := testutils.ConfigureOnRampChainWriter(envSettings.CCIPReport.Output.CCIPPackageId, envSettings.OnRampReport.Output.CCIPOnRampPackageId, []testutils.TokenToolDetails{tokenPoolDetails, ethTokenPoolDetails}, publicKeyBytes)
	require.NoError(t, err)
	lggr.Infow("chainWriterConfig", "chainWriterConfig", chainWriterConfig)
	chainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, chainWriterConfig, false)
	require.NoError(t, err)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	err = chainWriter.Start(ctx)
	require.NoError(t, err)

	t.Run("CCIP SUI messaging", func(t *testing.T) {
		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		mintedCoinId1, mintedCoinId2 := getLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

		// Create array with both coins for the PTB arguments
		linkCoins := []string{mintedCoinId1, mintedCoinId2}

		// Set up arguments for the PTB
		ptbArgs := createCCIPSendPTBArgsForBMAndLRTokenPools(
			lggr,
			destChainSelector,
			linkTokenType,
			ethTokenType,
			envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId,
			envSettings.MockEthTokenReport.Output.Objects.CoinMetadataObjectId,
			linkCoins,
			envSettings.EthCoins,
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId,
			clockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectId,
			envSettings.LockReleaseTokenPoolReport.Output.Objects.StateObjectId,
			envSettings.BurnMintTokenPoolReport.Output.Objects.StateObjectId,
			ethereumAddress,
		)
		txID := "ccip_send_test_message"

		lggr.Infow("Submitting transaction",
			"txID", txID,
			"accountAddress", accountAddress,
			"ptbArgs", ptbArgs,
			"chainWriterConfig", chainWriterConfig)

		err = chainWriter.SubmitTransaction(ctx,
			cwConfig.PTBChainWriterModuleName,
			"message_passing",
			&ptbArgs,
			txID,
			accountAddress,
			&commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			nil,
		)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
			if statusErr != nil {
				return false
			}

			return status == commonTypes.Finalized
		}, 5*time.Second, 1*time.Second, "Transaction final state not reached")

		// Create a PTB client to query events (the basic sui.ISuiAPI doesn't have QueryEvents)
		ptbClient, err := client.NewPTBClient(lggr, testutils.LocalUrl, nil, 10*time.Second, nil, 5, "WaitForLocalExecution")
		require.NoError(t, err, "Failed to create PTB client for event querying")

		// Query for ReceivedMessage events emitted by the dummy receiver
		eventFilter := client.EventFilterByMoveEventModule{
			Package: envSettings.OnRampReport.Output.CCIPOnRampPackageId,
			Module:  "onramp",
			Event:   "CCIPMessageSent",
		}

		// Query events with a small limit since we expect only one event
		limit := uint(10)
		eventsResponse, err := ptbClient.QueryEvents(ctx, eventFilter, &limit, nil, nil)
		lggr.Debugw("eventsResponse", "eventsResponse", eventsResponse)
		require.NoError(t, err, "Failed to query events")
		require.NotEmpty(t, eventsResponse.Data, "Expected at least one ReceivedMessage event")
		lggr.Infow("mostRecentEvent", "mostRecentEvent", eventsResponse.Data)
	})

	t.Run("CCIP SUI messaging with 1 LR TP and 1 BM TP", func(t *testing.T) {
		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		mintedCoinId1, mintedCoinId2 := getLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

		// Create array with both coins for the PTB arguments
		linkCoins := []string{mintedCoinId1, mintedCoinId2}

		// Set up arguments for the PTB
		ptbArgs := createCCIPSendPTBArgsForBMAndLRTokenPools(
			lggr,
			destChainSelector,
			linkTokenType,
			ethTokenType,
			envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId,
			envSettings.MockEthTokenReport.Output.Objects.CoinMetadataObjectId,
			linkCoins,
			envSettings.EthCoins,
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId,
			clockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectId,
			envSettings.LockReleaseTokenPoolReport.Output.Objects.StateObjectId,
			envSettings.BurnMintTokenPoolReport.Output.Objects.StateObjectId,
			ethereumAddress,
		)
		txID := "ccip_send_test_token"

		err = chainWriter.SubmitTransaction(ctx,
			cwConfig.PTBChainWriterModuleName,
			"token_transfer_with_messaging",
			&ptbArgs,
			txID,
			accountAddress,
			&commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			nil,
		)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
			if statusErr != nil {
				return false
			}

			return status == commonTypes.Finalized
		}, 5*time.Second, 1*time.Second, "Transaction final state not reached")
	})

	chainWriter.Close()
}

func TestCCIPSuiOnRampWithManagedTokenPool(t *testing.T) {
	lggr := logger.Test(t)

	localChainSelector := uint64(1)
	destChainSelector := uint64(2)

	// Create keystore and get account
	keystoreInstance := testutils.NewTestKeystore(t)

	// Wait a bit to ensure previous test's node is fully shut down
	time.Sleep(2 * time.Second)

	// Start dedicated Sui node for this test
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			if perr := cmd.Process.Kill(); perr != nil {
				t.Logf("Failed to kill Sui node process: %v", perr)
			}
		}
	})

	// Wait for the node to be fully ready
	time.Sleep(3 * time.Second)

	accountAddress, publicKeyBytes, signer, client, deps, bundle := basicSetUp(t, lggr, keystoreInstance)

	// Fund the account for gas payments
	for range 10 {
		err := testutils.FundWithFaucet(lggr, "localnet", accountAddress)
		require.NoError(t, err)
	}

	envSettings := SetupTestEnvironmentForManagedTokenPool(t, client, signer, accountAddress, bundle, deps, localChainSelector, destChainSelector, keystoreInstance)

	lggr.Infow("Using account", "address", accountAddress)

	_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr)

	ethManagedTokenPoolDetails := testutils.TokenToolDetails{
		TokenPoolPackageId: envSettings.ManagedTokenPoolReport.Output.ManagedTPPackageId,
		TokenPoolType:      testutils.TokenPoolTypeManaged,
	}

	chainWriterConfig, err := testutils.ConfigureOnRampChainWriter(envSettings.CCIPReport.Output.CCIPPackageId, envSettings.OnRampReport.Output.CCIPOnRampPackageId, []testutils.TokenToolDetails{ethManagedTokenPoolDetails}, publicKeyBytes)
	require.NoError(t, err)
	lggr.Infow("chainWriterConfig", "chainWriterConfig", chainWriterConfig)
	chainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, chainWriterConfig, false)
	require.NoError(t, err)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	err = chainWriter.Start(ctx)
	require.NoError(t, err)

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", envSettings.MockLinkReport.Output.PackageId)
	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", envSettings.MockEthTokenReport.Output.PackageId)

	t.Run("CCIP SUI messaging with 1 managed TP", func(t *testing.T) {
		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		mintedCoinId0, _ := getLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

		// Set up arguments for the PTB
		ptbArgs := createCCIPSendPTBArgsForManagedTokenPool(
			lggr,
			destChainSelector,
			linkTokenType,
			ethTokenType,
			envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId,
			envSettings.MockEthTokenReport.Output.Objects.CoinMetadataObjectId,
			mintedCoinId0,
			envSettings.EthCoins[0],
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId,
			clockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectId,
			envSettings.ManagedTokenReport.Output.Objects.StateObjectId,
			envSettings.ManagedTokenPoolReport.Output.Objects.StateObjectId,
			ethereumAddress,
		)
		txID := "ccip_send_test_token"

		err = chainWriter.SubmitTransaction(ctx,
			cwConfig.PTBChainWriterModuleName,
			"token_transfer_with_messaging",
			&ptbArgs,
			txID,
			accountAddress,
			&commonTypes.TxMeta{GasLimit: big.NewInt(10000000)},
			nil,
		)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			status, statusErr := chainWriter.GetTransactionStatus(ctx, txID)
			if statusErr != nil {
				return false
			}

			return status == commonTypes.Finalized
		}, 5*time.Second, 1*time.Second, "Transaction final state not reached")
	})

	chainWriter.Close()
}

// createCCIPSendPTBArgsForBMAndLRTokenPools creates PTBArgMapping for a CCIP send operation
func createCCIPSendPTBArgsForBMAndLRTokenPools(
	lggr logger.Logger,
	destChainSelector uint64,
	linkTokenType string,
	ethTokenType string,
	linkTokenMetadata string,
	ethTokenMetadata string,
	linkTokenCoinObjects []string,
	ethTokenCoinObjects []string,
	ccipObjectRef string,
	clockObject string,
	ccipOnrampState string,
	tokenPoolState string,
	ethTokenPoolState string,
	ethereumAddress string,
) cwConfig.Arguments {
	lggr.Infow("createCCIPSendPTBArgsForBMAndLRTokenPools", "destChainSelector", destChainSelector, "linkTokenType", linkTokenType, "linkTokenMetadata", linkTokenMetadata, "linkTokenCoinObjects", linkTokenCoinObjects, "ccipObjectRef", ccipObjectRef, "clockObject", clockObject, "ccipOnrampState", ccipOnrampState, "tokenPoolState", tokenPoolState)

	// Remove 0x prefix if present
	evmAddressBytes := normalizeTo32Bytes(ethereumAddress)

	lggr.Infow("evmAddressBytes", "evmAddressBytes", evmAddressBytes)

	return cwConfig.Arguments{
		Args: map[string]any{
			"ccip_object_ref":                    ccipObjectRef,
			"ccip_object_ref_mutable":            ccipObjectRef, // Same object, different parameter name
			"clock":                              clockObject,
			"destination_chain_selector":         destChainSelector,
			"link_lock_release_token_pool_state": tokenPoolState,
			"eth_burn_mint_token_pool_state":     ethTokenPoolState,
			"c_link":                             linkTokenCoinObjects[0],
			"c_eth":                              ethTokenCoinObjects[0],
			"onramp_state":                       ccipOnrampState,
			"receiver":                           evmAddressBytes,
			"data":                               []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"fee_token_metadata":                 linkTokenMetadata,
			"fee_token":                          linkTokenCoinObjects[1],
			"extra_args":                         []byte{}, // Empty array to use default gas limit
		},
		ArgTypes: map[string]string{
			"c_link":    linkTokenType,
			"c_eth":     ethTokenType,
			"fee_token": linkTokenType,
		},
	}
}

func createCCIPSendPTBArgsForManagedTokenPool(
	lggr logger.Logger,
	destChainSelector uint64,
	linkTokenType string,
	ethTokenType string,
	linkTokenMetadata string,
	ethTokenMetadata string,
	linkTokenCoinObject string,
	ethTokenCoinObject string,
	ccipObjectRef string,
	clockObject string,
	ccipOnrampState string,
	ethManagedTokenState string,
	ethManagedTokenPoolState string,
	ethereumAddress string,
) cwConfig.Arguments {
	lggr.Infow("createCCIPSendPTBArgsForManagedTokenPool", "destChainSelector", destChainSelector, "linkTokenType", linkTokenType, "linkTokenMetadata", linkTokenMetadata, "linkTokenCoinObject", linkTokenCoinObject, "ethTokenCoinObject", ethTokenCoinObject, "ccipObjectRef", ccipObjectRef, "clockObject", clockObject, "ccipOnrampState", ccipOnrampState, "ethManagedTokenState", ethManagedTokenState, "ethManagedTokenPoolState", ethManagedTokenPoolState)

	// Remove 0x prefix if present
	evmAddressBytes := normalizeTo32Bytes(ethereumAddress)

	lggr.Infow("evmAddressBytes", "evmAddressBytes", evmAddressBytes)

	return cwConfig.Arguments{
		Args: map[string]any{
			"ccip_object_ref":              ccipObjectRef,
			"ccip_object_ref_mutable":      ccipObjectRef, // Same object, different parameter name
			"clock":                        clockObject,
			"destination_chain_selector":   destChainSelector,
			"eth_managed_token_state":      ethManagedTokenState,
			"eth_managed_token_pool_state": ethManagedTokenPoolState,
			"c_managed_eth":                ethTokenCoinObject,
			"onramp_state":                 ccipOnrampState,
			"receiver":                     evmAddressBytes,
			"data":                         []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"fee_token_metadata":           linkTokenMetadata,
			"fee_token":                    linkTokenCoinObject,
			"extra_args":                   []byte{}, // Empty array to use default gas limit
			"deny_list":                    denyListObjectId,
		},
		ArgTypes: map[string]string{
			"c_managed_eth": ethTokenType,
			"fee_token":     linkTokenType,
		},
	}
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
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
