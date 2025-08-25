package environment

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
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
	cciptokenpoolop "github.com/smartcontractkit/chainlink-sui/ops/ccip_token_pool"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/ops/managed_token"
	mcmsops "github.com/smartcontractkit/chainlink-sui/ops/mcms"
	mockethtokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_eth_token"
	mocklinktokenops "github.com/smartcontractkit/chainlink-sui/ops/mock_link_token"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/stretchr/testify/require"
)

// SetupEthTokenPoolBurnMint sets up a burn/mint token pool for ETH tokens.
// This type of pool mints/burns tokens on demand rather than locking them.
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
		CCIPPackageId:    reportCCIP.Output.CCIPPackageID,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: accountAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Token Pool")

	// Deploy and initialize the burn mint token pool
	seqBurnMintDeployInput := burnmintops.DeployAndInitBurnMintTokenPoolInput{
		BurnMintTokenPoolDeployInput: burnmintops.BurnMintTokenPoolDeployInput{
			CCIPPackageId:          reportCCIP.Output.CCIPPackageID,
			CCIPTokenPoolPackageId: ccipTokenPoolReport.Output.PackageId,
			MCMSAddress:            reportMCMs.Output.PackageId,
			MCMSOwnerAddress:       accountAddress,
		},
		// Initialization parameters
		CoinObjectTypeArg:      tokenType,
		CCIPObjectRefObjectId:  reportCCIP.Output.Objects.CCIPObjectRefObjectID,
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

// SetupManagedTokenPool sets up a managed token pool with associated managed token.
// This includes deploying both the managed token and the managed token pool.
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
		CCIPPackageId:    reportCCIP.Output.CCIPPackageID,
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
		CCIPPackageId:          reportCCIP.Output.CCIPPackageID,
		CCIPTokenPoolPackageId: ccipTokenPoolReport.Output.PackageId,
		ManagedTokenPackageId:  managedTokenReport.Output.ManagedTokenPackageId,
		MCMSAddress:            reportMCMs.Output.PackageId,
		MCMSOwnerAddress:       accountAddress,
		// Initialize inputs
		CoinObjectTypeArg:         tokenType,
		CCIPObjectRefObjectId:     reportCCIP.Output.Objects.CCIPObjectRefObjectID,
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

// SetupTokenPool sets up a lock/release token pool for LINK tokens.
// This type of pool locks tokens when sending and releases them when receiving.
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
		CCIPPackageId:    reportCCIP.Output.CCIPPackageID,
		MCMSAddress:      reportMCMs.Output.PackageId,
		MCMSOwnerAddress: accountAddress,
	})
	require.NoError(t, err, "failed to deploy CCIP Token Pool")

	// Deploy and initialize the lock release token pool
	seqLockReleaseDeployInput := lockreleaseops.DeployAndInitLockReleaseTokenPoolInput{
		LockReleaseTokenPoolDeployInput: lockreleaseops.LockReleaseTokenPoolDeployInput{
			CCIPPackageID:          reportCCIP.Output.CCIPPackageID,
			CCIPTokenPoolPackageID: ccipTokenPoolReport.Output.PackageId,
			MCMSAddress:            reportMCMs.Output.PackageId,
			MCMSOwnerAddress:       accountAddress,
		},
		// Initialization parameters
		CoinObjectTypeArg:      linkTokenType,
		CCIPObjectRefObjectID:  reportCCIP.Output.Objects.CCIPObjectRefObjectID,
		CoinMetadataObjectID:   mockLinkReport.Output.Objects.CoinMetadataObjectId,
		TreasuryCapObjectID:    mockLinkReport.Output.Objects.TreasuryCapObjectId,
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
		LockReleaseTokenPoolPackageID: tokenPoolLockReleaseReport.Output.LockReleaseTPPackageID,
		StateObjectID:                 tokenPoolLockReleaseReport.Output.Objects.StateObjectID,
		Coin:                          mintedCoinId,
		CoinObjectTypeArg:             linkTokenType,
	}

	_, err = cld_ops.ExecuteOperation(bundle, lockreleaseops.LockReleaseTokenPoolProviderLiquidityOp, deps, provideLiquidityInput)
	require.NoError(t, err, "failed to provide liquidity to Lock Release Token Pool")

	lggr.Debugw("Provided liquidity to Lock Release Token Pool", "amount", liquidityAmount)

	return &tokenPoolLockReleaseReport
}

// MintTestTokens mints tokens for testing purposes.
// This is a helper function to mint both transfer and fee tokens.
func MintTestTokens(
	t *testing.T,
	client sui.ISuiAPI,
	signer rel.SuiSigner,
	packageId, treasuryCapId, tokenType, recipient string,
	transferAmount, feeAmount uint64,
	lggr logger.Logger,
) (transferCoin, feeCoin string) {
	t.Helper()

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

	// Determine contract type and create appropriate instance
	var mintTransferTx, mintFeeTx *models.SuiTransactionBlockResponse
	var err error

	// Check if this is a LINK token or ETH token based on the package structure
	if packageId != "" {
		// Try LINK token first
		if linkContract, linkErr := mocklinktoken.NewMockLinkToken(packageId, client); linkErr == nil {
			// Mint transfer token
			mintTransferTx, err = linkContract.MockLinkToken().MintAndTransfer(
				context.Background(),
				deps.GetCallOpts(),
				bind.Object{Id: treasuryCapId},
				transferAmount,
				recipient,
			)
			require.NoError(t, err, "failed to mint and transfer LINK tokens for transfer")

			// Mint fee token
			mintFeeTx, err = linkContract.MockLinkToken().MintAndTransfer(
				context.Background(),
				deps.GetCallOpts(),
				bind.Object{Id: treasuryCapId},
				feeAmount,
				recipient,
			)
			require.NoError(t, err, "failed to mint and transfer LINK tokens for fee")
		} else {
			// Try ETH token
			ethContract, ethErr := mockethtoken.NewMockEthToken(packageId, client)
			require.NoError(t, ethErr, "failed to create token contract")

			// Mint transfer token
			mintTransferTx, err = ethContract.MockEthToken().MintAndTransfer(
				context.Background(),
				deps.GetCallOpts(),
				bind.Object{Id: treasuryCapId},
				transferAmount,
				recipient,
			)
			require.NoError(t, err, "failed to mint and transfer ETH tokens for transfer")

			// Mint fee token
			mintFeeTx, err = ethContract.MockEthToken().MintAndTransfer(
				context.Background(),
				deps.GetCallOpts(),
				bind.Object{Id: treasuryCapId},
				feeAmount,
				recipient,
			)
			require.NoError(t, err, "failed to mint and transfer ETH tokens for fee")
		}
	}

	// Find coin object IDs from transactions
	transferCoinId, err := bind.FindCoinObjectIdFromTx(*mintTransferTx, tokenType)
	require.NoError(t, err, "failed to find transfer coin object ID")

	feeCoinId, err := bind.FindCoinObjectIdFromTx(*mintFeeTx, tokenType)
	require.NoError(t, err, "failed to find fee coin object ID")

	lggr.Infow("Successfully minted test tokens",
		"transferCoin", transferCoinId,
		"feeCoin", feeCoinId,
		"transferAmount", transferAmount,
		"feeAmount", feeAmount,
		"recipient", recipient)

	return transferCoinId, feeCoinId
}
