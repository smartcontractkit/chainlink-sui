//go:build integration

package txm_test

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	modulecounter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

// TestTransactionGeneration tests the complete flow of generating and executing a Sui transaction
// using PTBs. This integration test verifies:
//
// 1. PTB client setup and account funding
// 2. Smart contract interaction (counter increment operation)
// 3. Gas management and estimation
// 4. Transaction generation with proper gas budget calculation
// 5. Transaction execution and verification of results
// 6. Gas consumption validation and coin usage optimization (gas smashing)
//
// The test ensures we can properly:
// - Generate transactions with accurate gas estimates
// - Execute smart contract calls through PTB
// - Optimize gas usage by consolidating multiple coins into a single payment
// - Validate that the final gas budget matches the actual gas consumed
//
// This test requires a running local Sui network and performs actual on-chain operations.
//
//nolint:paralleltest
func TestTransactionGeneration(t *testing.T) {
	ctx := context.Background()
	lggr := logger.Test(t)
	lggr.Debugw("Starting Sui node")

	gasLimit := int64(200000000000)
	ptbClient, _, _, accountAddress, _, publicKeyBytes, packageId, counterObjectId := testutils.SetupTestEnv(t, ctx, lggr, gasLimit)

	// Generate key pair and create a signer - use the same key for both signer and keystore
	pk, _, _, err := testutils.GenerateAccountKeyPair(t)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)
	accountAddress, err = signer.GetAddress()
	require.NoError(t, err)

	err = testutils.FundWithFaucet(lggr, "localnet", accountAddress)
	require.NoError(t, err)

	coins, err := ptbClient.GetCoinsByAddress(ctx, accountAddress)
	require.NoError(t, err)
	lggr.Debugw("Coins", "coins", coins)

	gasBudget := uint64(200000000000)

	opts := &bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
		GasBudget:        &gasBudget,
	}

	suiClient := ptbClient.GetClient()

	lggr.Debugw("Published Contract", "packageId", packageId)
	lggr.Debugw("Account Address", "accountAddress", accountAddress)
	lggr.Debugw("Counter object created", "counterObjectId", counterObjectId)

	counterInterface, err := modulecounter.NewCounter(packageId, suiClient)
	require.NoError(t, err)
	counter, ok := counterInterface.(*modulecounter.CounterContract)
	require.True(t, ok, "Failed to cast to CounterContract")

	counterObj := bind.Object{
		Id: counterObjectId,
		// InitialSharedVersion will be resolved automatically by the object resolver
	}

	gasManager := txm.NewSuiGasManager(lggr, ptbClient, *big.NewInt(int64(gasBudget)), 0)
	txID := "1"

	// Create a test keystore and add the signer's key (use the same key as the publicKeyBytes)
	keystore := testutils.NewTestKeystore(t)
	keystore.AddKey(pk)

	// Get the public key bytes from the private key for the transaction
	publicKeyBytes = pk.Public().(ed25519.PublicKey)

	t.Run("GeneratePTBTransactionWithGasEstimation", func(t *testing.T) {
		ptb := transaction.NewTransaction()
		inc, err := counter.Encoder().IncrementBy(counterObj, 10)
		require.NoError(t, err)

		_, err = counter.AppendPTB(ctx, opts, ptb, inc)
		require.NoError(t, err)

		ptb.SetGasPrice(10000000)

		txMeta := &commontypes.TxMeta{
			GasLimit: big.NewInt(int64(gasBudget)),
		}

		coinManager := txm.NewGasCoinManager(lggr, ptbClient)

		tx, err := txm.GeneratePTBTransactionWithGasEstimation(
			ctx,
			publicKeyBytes,
			lggr,
			keystore,
			ptbClient,
			"WaitForEffectsCert",
			txID,
			txMeta,
			ptb,
			true,
			gasManager,
			coinManager,
		)

		finalGasBudget := tx.GasBudget
		lggr.Debugw("Final gas budget", "finalGasBudget", finalGasBudget)

		require.NoError(t, err)
		lggr.Debugw("PTB transaction generated", "tx", tx)

		payload := client.TransactionBlockRequest{
			TxBytes:    tx.Payload,
			Signatures: tx.Signatures,
			Options: client.TransactionBlockOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
				ShowEvents:         true,
			},
			RequestType: tx.RequestType,
		}

		resp, err := ptbClient.SendTransaction(ctx, payload)
		require.NoError(t, err)

		gasUsed := resp.Effects.GasUsed
		lggr.Debugw("Gas used", "gasUsed", gasUsed)
		computationCost, err := strconv.ParseInt(gasUsed.ComputationCost, 10, 64)
		require.NoError(t, err)
		storageCost, err := strconv.ParseInt(gasUsed.StorageCost, 10, 64)
		require.NoError(t, err)
		storageRebate, _ := strconv.ParseInt(gasUsed.StorageRebate, 10, 64)
		if storageRebate != 0 {
			storageCost = storageCost - storageRebate
		}

		totalGasUsed := computationCost + storageCost
		require.Greater(t, totalGasUsed, int64(0))
		require.Equal(t, totalGasUsed, int64(finalGasBudget))

		objectChanges := resp.ObjectChanges
		usedCoins := []transaction.SuiObjectRef{}
		for _, objectChange := range objectChanges {
			if objectChange.Type == "mutated" && objectChange.ObjectType == "0x2::coin::Coin<0x2::sui::SUI>" {
				version, err := strconv.ParseUint(objectChange.PreviousVersion, 10, 64)
				require.NoError(t, err)

				objectIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(objectChange.ObjectId))

				usedCoins = append(usedCoins, transaction.SuiObjectRef{
					ObjectId: *objectIdBytes,
					Version:  version,
					Digest:   nil,
				})
			}
		}

		lggr.Debugw("Transaction broadcasted", "resp", resp)
		lggr.Debugw("Used coins", "usedCoins", usedCoins)

		// Test that the used coins in a single element, to confirm that gas smashing was used
		require.Len(t, usedCoins, 1)
	})
}

// TestCoinSelectionEdgeCases tests edge cases in coin selection logic
//
//nolint:paralleltest
func TestCoinSelectionEdgeCases(t *testing.T) {
	lggr := logger.Test(t)

	// Test case 1: Empty coin list
	t.Run("EmptyCoinList", func(t *testing.T) {
		_, err := txm.SelectCoinsForGasBudget(1000000, []models.CoinData{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "no coins available")
	})

	// Test case 2: No SUI coins available
	t.Run("NoSUICoins", func(t *testing.T) {
		nonSuiCoins := []models.CoinData{
			{
				CoinType: "0x123::other::TOKEN",
				Balance:  "1000000000",
			},
		}
		_, err := txm.SelectCoinsForGasBudget(1000000, nonSuiCoins)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no SUI coins available")
	})

	// Test case 3: Insufficient balance
	t.Run("InsufficientBalance", func(t *testing.T) {
		insufficientCoins := []models.CoinData{
			{
				CoinType: "0x2::sui::SUI",
				Balance:  "500000", // Less than required
			},
		}
		_, err := txm.SelectCoinsForGasBudget(1000000, insufficientCoins)
		require.Error(t, err)
		require.Contains(t, err.Error(), "insufficient funds")
	})

	// Test case 4: Exact balance match
	t.Run("ExactBalanceMatch", func(t *testing.T) {
		exactCoins := []models.CoinData{
			{
				CoinType: "0x2::sui::SUI",
				Balance:  "1000000", // Exactly what's needed
			},
		}
		selected, err := txm.SelectCoinsForGasBudget(1000000, exactCoins)
		require.NoError(t, err)
		require.Len(t, selected, 1)
	})

	// Test case 5: Multiple coins needed
	t.Run("MultipleCoinsCombined", func(t *testing.T) {
		multipleCoins := []models.CoinData{
			{
				CoinType: "0x2::sui::SUI",
				Balance:  "600000",
			},
			{
				CoinType: "0x2::sui::SUI",
				Balance:  "500000",
			},
		}
		selected, err := txm.SelectCoinsForGasBudget(1000000, multipleCoins)
		require.NoError(t, err)
		require.Len(t, selected, 2) // Should select both coins

		var totalBalance uint64
		for _, coin := range selected {
			var balance uint64
			_, parseErr := fmt.Sscanf(coin.Balance, "%d", &balance)
			require.NoError(t, parseErr)
			totalBalance += balance
		}
		require.GreaterOrEqual(t, totalBalance, uint64(1000000))
	})

	lggr.Debugw("Coin selection edge cases test completed")
}
