//go:build integration

package ccip_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mockethtoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_eth_token"
	mocklinktoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/mock_link_token"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/onramp/environment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
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

func getLinkCoins_OLD_DELETE_ME(t *testing.T, envSettings *environment.EnvironmentSettings, linkTokenType string, accountAddress string, lggr logger.Logger, tokenAmount uint64, feeAmount uint64) (string, string) {
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

func getEthCoins_OLD_DELETE_ME(t *testing.T, client sui.ISuiAPI, signer rel.SuiSigner, ethTokenPackageId string, treasuryCapObjectId string, ethTokenType string, accountAddress string, lggr logger.Logger, tokenAmount uint64, feeAmount uint64) []string {
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

	envSettings := environment.SetupTestEnvironment(t, localChainSelector, destChainSelector, keystoreInstance)

	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)
	lggr.Infow("Using account", "address", accountAddress)

	// Fund the account for gas payments
	for range 3 {
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

	chainWriterConfig, err := testutils.ConfigureOnRampChainWriter(
		envSettings.CCIPReport.Output.CCIPPackageID,
		envSettings.OnRampReport.Output.CCIPOnRampPackageID,
		[]testutils.TokenToolDetails{tokenPoolDetails, ethTokenPoolDetails},
		publicKeyBytes,
		linkTokenType,
		linkTokenType,
		ethTokenType,
	)
	require.NoError(t, err)

	lggr.Infow("chainWriterConfig", "chainWriterConfig", chainWriterConfig)
	chainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, chainWriterConfig, false)
	require.NoError(t, err)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	err = chainWriter.Start(ctx)
	require.NoError(t, err)

	err = txManager.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		txManager.Close()
		chainWriter.Close()
	})

	t.Run("CCIP SUI messaging", func(t *testing.T) {
		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		mintedCoinId1, mintedCoinId2 := environment.GetLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

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
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectID,
			environment.ClockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectID,
			envSettings.LockReleaseTokenPoolReport.Output.Objects.StateObjectID,
			envSettings.BurnMintTokenPoolReport.Output.Objects.StateObjectId,
			environment.EthereumAddress,
		)
		txID := "ccip_send_test_message"

		lggr.Infow("ptbArgs", "ptbArgs", ptbArgs)

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
			Package: envSettings.OnRampReport.Output.CCIPOnRampPackageID,
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

		mintedCoinId1, mintedCoinId2 := environment.GetLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

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
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectID,
			environment.ClockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectID,
			envSettings.LockReleaseTokenPoolReport.Output.Objects.StateObjectID,
			envSettings.BurnMintTokenPoolReport.Output.Objects.StateObjectId,
			environment.EthereumAddress,
		)
		txID := "ccip_send_two_token_pool_calls"

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

	accountAddress, publicKeyBytes, signer, client, deps, bundle := environment.BasicSetUp(t, lggr, keystoreInstance)

	// Fund the account for gas payments
	for range 3 {
		err := testutils.FundWithFaucet(lggr, "localnet", accountAddress)
		require.NoError(t, err)
	}

	envSettings := environment.SetupTestEnvironmentForManagedTokenPool(t, client, signer, accountAddress, bundle, deps, localChainSelector, destChainSelector, keystoreInstance)

	lggr.Infow("Using account", "address", accountAddress)

	_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr)

	ethManagedTokenPoolDetails := testutils.TokenToolDetails{
		TokenPoolPackageId: envSettings.ManagedTokenPoolReport.Output.ManagedTPPackageId,
		TokenPoolType:      testutils.TokenPoolTypeManaged,
	}

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", envSettings.MockLinkReport.Output.PackageId)
	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", envSettings.MockEthTokenReport.Output.PackageId)

	chainWriterConfig, err := testutils.ConfigureOnRampChainWriter(
		envSettings.CCIPReport.Output.CCIPPackageID,
		envSettings.OnRampReport.Output.CCIPOnRampPackageID,
		[]testutils.TokenToolDetails{ethManagedTokenPoolDetails},
		publicKeyBytes,
		linkTokenType,
		linkTokenType,
		ethTokenType,
	)
	require.NoError(t, err)

	lggr.Infow("chainWriterConfig", "chainWriterConfig", chainWriterConfig)
	chainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, chainWriterConfig, false)
	require.NoError(t, err)

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	// Start the TxManager first (required for broadcasting transactions)
	err = txManager.Start(ctx)
	require.NoError(t, err)

	err = chainWriter.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		txManager.Close()
		chainWriter.Close()
	})

	t.Run("CCIP SUI messaging with 1 managed TP", func(t *testing.T) {
		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		mintedCoinId0, _ := environment.GetLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

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
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectID,
			environment.ClockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectID,
			envSettings.ManagedTokenReport.Output.Objects.StateObjectId,
			envSettings.ManagedTokenPoolReport.Output.Objects.StateObjectId,
			environment.EthereumAddress,
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
) map[string]any {
	lggr.Infow("createCCIPSendPTBArgsForBMAndLRTokenPools", "destChainSelector", destChainSelector, "linkTokenType", linkTokenType, "linkTokenMetadata", linkTokenMetadata, "linkTokenCoinObjects", linkTokenCoinObjects, "ccipObjectRef", ccipObjectRef, "clockObject", clockObject, "ccipOnrampState", ccipOnrampState, "tokenPoolState", tokenPoolState)

	// Remove 0x prefix if present
	evmAddressBytes := environment.NormalizeTo32Bytes(ethereumAddress)

	lggr.Infow("evmAddressBytes", "evmAddressBytes", evmAddressBytes)

	return map[string]any{
		"ccip_object_ref":                    ccipObjectRef,
		"ccip_object_ref_mutable":            ccipObjectRef, // Same object, different parameter name
		"clock":                              environment.ClockObjectId,
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
		"token_receiver":                     testutils.ZeroAddress,
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
) map[string]any {
	lggr.Infow("createCCIPSendPTBArgsForManagedTokenPool", "destChainSelector", destChainSelector, "linkTokenType", linkTokenType, "linkTokenMetadata", linkTokenMetadata, "linkTokenCoinObject", linkTokenCoinObject, "ethTokenCoinObject", ethTokenCoinObject, "ccipObjectRef", ccipObjectRef, "clockObject", clockObject, "ccipOnrampState", ccipOnrampState, "ethManagedTokenState", ethManagedTokenState, "ethManagedTokenPoolState", ethManagedTokenPoolState)

	// Remove 0x prefix if present
	evmAddressBytes := environment.NormalizeTo32Bytes(ethereumAddress)

	lggr.Infow("evmAddressBytes", "evmAddressBytes", evmAddressBytes)

	return map[string]any{
		"ccip_object_ref":              ccipObjectRef,
		"ccip_object_ref_mutable":      ccipObjectRef, // Same object, different parameter name
		"clock":                        environment.ClockObjectId,
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
		"deny_list":                    environment.DenyListObjectId,
		"token_receiver":               testutils.ZeroAddress,
	}
}

// Helper function to convert a string to a string pointer
func strPtr(s string) *string {
	return &s
}
