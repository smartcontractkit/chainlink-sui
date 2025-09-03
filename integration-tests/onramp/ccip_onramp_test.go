//go:build integration

package ccip_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-sui/integration-tests/onramp/environment"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
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

// TestCCIPSuiOnRamp tests the CCIP onramp send functionality
func TestCCIPSuiOnRamp(t *testing.T) {
	lggr := logger.Test(t)

	localChainSelector := uint64(1)
	destChainSelector := uint64(2)

	gasBudget := int64(500_000_000)

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

	c := context.Background()
	ctx, cancel := context.WithCancel(c)
	defer cancel()

	t.Run("CCIP SUI messaging", func(t *testing.T) {
		_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr, gasBudget)
		tokenPoolDetails := testutils.TokenToolDetails{
			TokenPoolPackageId: envSettings.LockReleaseTokenPoolReport.Output.LockReleaseTPPackageID,
			TokenPoolType:      testutils.TokenPoolTypeLockRelease,
		}
		ethTokenPoolDetails := testutils.TokenToolDetails{
			TokenPoolPackageId: envSettings.BurnMintTokenPoolReport.Output.BurnMintTPPackageID,
			TokenPoolType:      testutils.TokenPoolTypeBurnMint,
		}

		err = txManager.Start(ctx)
		require.NoError(t, err)

		chainWriterConfig, err := testutils.ConfigureOnRampChainWriter(
			lggr,
			envSettings.CCIPReport.Output.CCIPPackageId,
			envSettings.OnRampReport.Output.CCIPOnRampPackageId,
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

		err = chainWriter.Start(ctx)
		require.NoError(t, err)

		t.Cleanup(func() {
			txManager.Close()
			chainWriter.Close()
		})

		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		gasBudget := int64(500_000_000)

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
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId,
			environment.ClockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectId,
			envSettings.LockReleaseTokenPoolReport.Output.Objects.StateObjectId,
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

		offrampPackageId := envSettings.OnRampReport.Output.CCIPOnRampPackageId

		err = chainWriter.SubmitTransaction(ctx,
			cwConfig.PTBChainWriterModuleName,
			"message_passing",
			&ptbArgs,
			txID,
			offrampPackageId,
			&commonTypes.TxMeta{GasLimit: big.NewInt(gasBudget)},
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

	t.Run("CCIP SUI messaging with Lock Release Token Pool", func(t *testing.T) {
		_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr, gasBudget)

		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		mintedCoinId1, mintedCoinId2 := environment.GetLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

		// Create array with both coins for the PTB arguments
		linkCoins := []string{mintedCoinId1, mintedCoinId2}

		// Create chain writer config for Lock Release Token Pool only
		lrTokenPoolDetails := testutils.TokenToolDetails{
			TokenPoolPackageId: envSettings.LockReleaseTokenPoolReport.Output.LockReleaseTPPackageID,
			TokenPoolType:      testutils.TokenPoolTypeLockRelease,
		}

		lrChainWriterConfig, err := testutils.ConfigureOnRampChainWriter(
			lggr,
			envSettings.CCIPReport.Output.CCIPPackageId,
			envSettings.OnRampReport.Output.CCIPOnRampPackageId,
			[]testutils.TokenToolDetails{lrTokenPoolDetails},
			publicKeyBytes,
			linkTokenType,
			linkTokenType,
			ethTokenType,
		)
		require.NoError(t, err)

		lrChainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, lrChainWriterConfig, false)
		require.NoError(t, err)

		err = txManager.Start(ctx)
		require.NoError(t, err)

		err = lrChainWriter.Start(ctx)
		require.NoError(t, err)

		t.Cleanup(func() {
			txManager.Close()
			lrChainWriter.Close()
		})

		// Set up arguments for the PTB - only Lock Release Token Pool
		ptbArgs := createCCIPSendPTBArgsForLRTokenPool(
			lggr,
			destChainSelector,
			linkTokenType,
			envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId,
			linkCoins,
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId,
			environment.ClockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectId,
			envSettings.LockReleaseTokenPoolReport.Output.Objects.StateObjectId,
			environment.EthereumAddress,
		)
		txID := "ccip_send_lock_release_token_pool"

		err = lrChainWriter.SubmitTransaction(ctx,
			cwConfig.PTBChainWriterModuleName,
			"token_transfer_with_messaging",
			&ptbArgs,
			txID,
			accountAddress,
			&commonTypes.TxMeta{GasLimit: big.NewInt(gasBudget)},
			nil,
		)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			status, statusErr := lrChainWriter.GetTransactionStatus(ctx, txID)
			if statusErr != nil {
				return false
			}

			return status == commonTypes.Finalized
		}, 5*time.Second, 1*time.Second, "Transaction final state not reached")
	})

	t.Run("CCIP SUI messaging with Burn Mint Token Pool", func(t *testing.T) {
		_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr, gasBudget)

		tokenAmount := uint64(500000) // 500K tokens for transfer
		feeAmount := uint64(100000)   // 100K tokens for fee payment

		gasBudget := int64(500_000_000)

		_, mintedCoinId2 := environment.GetLinkCoins(t, envSettings, linkTokenType, accountAddress, lggr, tokenAmount, feeAmount)

		// Create chain writer config for Burn Mint Token Pool only
		bmTokenPoolDetails := testutils.TokenToolDetails{
			TokenPoolPackageId: envSettings.BurnMintTokenPoolReport.Output.BurnMintTPPackageID,
			TokenPoolType:      testutils.TokenPoolTypeBurnMint,
		}

		bmChainWriterConfig, err := testutils.ConfigureOnRampChainWriter(
			lggr,
			envSettings.CCIPReport.Output.CCIPPackageId,
			envSettings.OnRampReport.Output.CCIPOnRampPackageId,
			[]testutils.TokenToolDetails{bmTokenPoolDetails},
			publicKeyBytes,
			linkTokenType,
			linkTokenType,
			ethTokenType,
		)
		require.NoError(t, err)

		lggr.Debugw("bmChainWriterConfig", "bmChainWriterConfig", bmChainWriterConfig)

		bmChainWriter, err := chainwriter.NewSuiChainWriter(lggr, txManager, bmChainWriterConfig, false)
		require.NoError(t, err)

		err = txManager.Start(ctx)
		require.NoError(t, err)

		err = bmChainWriter.Start(ctx)
		require.NoError(t, err)

		t.Cleanup(func() {
			txManager.Close()
			bmChainWriter.Close()
		})

		// Set up arguments for the PTB - only Burn Mint Token Pool
		ptbArgs := createCCIPSendPTBArgsForBMTokenPool(
			lggr,
			destChainSelector,
			linkTokenType,
			ethTokenType,
			envSettings.MockLinkReport.Output.Objects.CoinMetadataObjectId,
			envSettings.MockEthTokenReport.Output.Objects.CoinMetadataObjectId,
			mintedCoinId2,           // fee token
			envSettings.EthCoins[0], // token to transfer
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId,
			environment.ClockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectId,
			envSettings.BurnMintTokenPoolReport.Output.Objects.StateObjectId,
			environment.EthereumAddress,
		)
		txID := "ccip_send_burn_mint_token_pool"

		offrampPackageId := envSettings.OnRampReport.Output.CCIPOnRampPackageId

		err = bmChainWriter.SubmitTransaction(ctx,
			cwConfig.PTBChainWriterModuleName,
			"token_transfer_with_messaging",
			&ptbArgs,
			txID,
			offrampPackageId,
			&commonTypes.TxMeta{GasLimit: big.NewInt(gasBudget)},
			nil,
		)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			status, statusErr := bmChainWriter.GetTransactionStatus(ctx, txID)
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

	gasBudget := int64(500_000_000)

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

	_, txManager, _ := testutils.SetupClients(t, testutils.LocalUrl, keystoreInstance, lggr, gasBudget)

	ethManagedTokenPoolDetails := testutils.TokenToolDetails{
		TokenPoolPackageId: envSettings.ManagedTokenPoolReport.Output.ManagedTPPackageId,
		TokenPoolType:      testutils.TokenPoolTypeManaged,
	}

	linkTokenType := fmt.Sprintf("%s::mock_link_token::MOCK_LINK_TOKEN", envSettings.MockLinkReport.Output.PackageId)
	ethTokenType := fmt.Sprintf("%s::mock_eth_token::MOCK_ETH_TOKEN", envSettings.MockEthTokenReport.Output.PackageId)

	chainWriterConfig, err := testutils.ConfigureOnRampChainWriter(
		lggr,
		envSettings.CCIPReport.Output.CCIPPackageId,
		envSettings.OnRampReport.Output.CCIPOnRampPackageId,
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
			envSettings.CCIPReport.Output.Objects.CCIPObjectRefObjectId,
			environment.ClockObjectId,
			envSettings.OnRampReport.Output.Objects.StateObjectId,
			envSettings.ManagedTokenReport.Output.Objects.StateObjectId,
			envSettings.ManagedTokenPoolReport.Output.Objects.StateObjectId,
			environment.EthereumAddress,
		)
		txID := "ccip_send_test_token"

		offrampPackageId := envSettings.OnRampReport.Output.CCIPOnRampPackageId

		err = chainWriter.SubmitTransaction(ctx,
			cwConfig.PTBChainWriterModuleName,
			"token_transfer_with_messaging",
			&ptbArgs,
			txID,
			offrampPackageId,
			&commonTypes.TxMeta{GasLimit: big.NewInt(gasBudget)},
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

// createCCIPSendPTBArgsForLRTokenPool creates PTBArgMapping for a CCIP send operation with Lock Release Token Pool only
func createCCIPSendPTBArgsForLRTokenPool(
	lggr logger.Logger,
	destChainSelector uint64,
	linkTokenType string,
	linkTokenMetadata string,
	linkTokenCoinObjects []string,
	ccipObjectRef string,
	clockObject string,
	ccipOnrampState string,
	tokenPoolState string,
	ethereumAddress string,
) map[string]any {
	lggr.Infow("createCCIPSendPTBArgsForLRTokenPool", "destChainSelector", destChainSelector, "linkTokenType", linkTokenType, "linkTokenMetadata", linkTokenMetadata, "linkTokenCoinObjects", linkTokenCoinObjects, "ccipObjectRef", ccipObjectRef, "clockObject", clockObject, "ccipOnrampState", ccipOnrampState, "tokenPoolState", tokenPoolState)

	// Remove 0x prefix if present
	evmAddressBytes := environment.NormalizeTo32Bytes(ethereumAddress)

	lggr.Infow("evmAddressBytes", "evmAddressBytes", evmAddressBytes)

	return map[string]any{
		"ccip_object_ref":                    ccipObjectRef,
		"ccip_object_ref_mutable":            ccipObjectRef, // Same object, different parameter name
		"clock":                              environment.ClockObjectId,
		"destination_chain_selector":         destChainSelector,
		"link_lock_release_token_pool_state": tokenPoolState,
		"c_link":                             linkTokenCoinObjects[0],
		"onramp_state":                       ccipOnrampState,
		"receiver":                           evmAddressBytes,
		"data":                               []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"fee_token_metadata":                 linkTokenMetadata,
		"fee_token":                          linkTokenCoinObjects[1],
		"extra_args":                         []byte{}, // Empty array to use default gas limit
		"token_receiver":                     testutils.ZeroAddress,
	}
}

// createCCIPSendPTBArgsForBMTokenPool creates PTBArgMapping for a CCIP send operation with Burn Mint Token Pool only
func createCCIPSendPTBArgsForBMTokenPool(
	lggr logger.Logger,
	destChainSelector uint64,
	linkTokenType string,
	ethTokenType string,
	linkTokenMetadata string,
	ethTokenMetadata string,
	feeTokenCoinObject string,
	ethTokenCoinObject string,
	ccipObjectRef string,
	clockObject string,
	ccipOnrampState string,
	ethTokenPoolState string,
	ethereumAddress string,
) map[string]any {
	lggr.Infow("createCCIPSendPTBArgsForBMTokenPool", "destChainSelector", destChainSelector, "ethTokenType", ethTokenType, "ethTokenMetadata", ethTokenMetadata, "ethTokenCoinObject", ethTokenCoinObject, "ccipObjectRef", ccipObjectRef, "clockObject", clockObject, "ccipOnrampState", ccipOnrampState, "ethTokenPoolState", ethTokenPoolState)

	// Remove 0x prefix if present
	evmAddressBytes := environment.NormalizeTo32Bytes(ethereumAddress)

	lggr.Infow("evmAddressBytes", "evmAddressBytes", evmAddressBytes)

	return map[string]any{
		"ccip_object_ref":                ccipObjectRef,
		"ccip_object_ref_mutable":        ccipObjectRef, // Same object, different parameter name
		"clock":                          environment.ClockObjectId,
		"destination_chain_selector":     destChainSelector,
		"eth_burn_mint_token_pool_state": ethTokenPoolState,
		"c_eth":                          ethTokenCoinObject,
		"onramp_state":                   ccipOnrampState,
		"receiver":                       evmAddressBytes,
		"data":                           []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"fee_token_metadata":             linkTokenMetadata,
		"fee_token":                      feeTokenCoinObject,
		"extra_args":                     []byte{}, // Empty array to use default gas limit
		"token_receiver":                 testutils.ZeroAddress,
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
