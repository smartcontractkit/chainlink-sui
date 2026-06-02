//go:build integration

package txm_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	modulecounter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	"github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

func newTestCoin(t *testing.T, coinType string, balance uint64) *suirpcv2.Object {
	t.Helper()

	objectID := fmt.Sprintf("0xcoin-%d", balance)
	digest := fmt.Sprintf("digest-%d", balance)
	version := uint64(1)

	return &suirpcv2.Object{
		ObjectId:   &objectID,
		ObjectType: &coinType,
		Balance:    &balance,
		Version:    &version,
		Digest:     &digest,
	}
}

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
	ptbClient, _, _, accountAddress, keystore, publicKeyBytes, packageId, counterObjectId := testutils.SetupTestEnv(t, ctx, lggr, gasLimit)

	// Fund the account multiple times to ensure sufficient balance as separate objects
	for i := 0; i < 5; i++ {
		err := testutils.FundWithFaucet(lggr, "localnet", accountAddress)
		require.NoError(t, err)
	}

	gasBudget := uint64(200000000000)

	publicKey := fmt.Sprintf("%064x", publicKeyBytes)
	suiSigner := utils.NewTestPrivateKeySigner(keystore.GetSuiSigner(ctx, publicKey).PriKey)

	opts := &bind.CallOpts{
		Signer:           suiSigner,
		WaitForExecution: true,
		GasBudget:        &gasBudget,
	}

	lggr.Debugw("Published Contract", "packageId", packageId)
	lggr.Debugw("Account Address", "accountAddress", accountAddress)
	lggr.Debugw("Counter object created", "counterObjectId", counterObjectId)

	counterInterface, err := modulecounter.NewCounter(packageId, ptbClient)
	require.NoError(t, err)
	counter, ok := counterInterface.(*modulecounter.CounterContract)
	require.True(t, ok, "Failed to cast to CounterContract")

	counterObj := bind.Object{
		Id: counterObjectId,
		// InitialSharedVersion will be resolved automatically by the object resolver
	}

	gasManager := txm.NewSuiGasManager(lggr, ptbClient, *big.NewInt(int64(gasBudget)), 0)
	txID := "1"

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
		require.NoError(t, err)
		require.NotNil(t, tx)

		finalGasBudget := tx.GasBudget
		lggr.Debugw("Final gas budget", "finalGasBudget", finalGasBudget)
		lggr.Debugw("PTB transaction generated", "tx", tx)

		resp, err := ptbClient.SignAndSendTransaction(ctx, tx.Payload, publicKeyBytes)
		require.NoError(t, err)

		gasUsed := resp.GetTransaction().GetEffects().GetGasUsed()
		computationCost := gasUsed.GetComputationCost()
		storageCost := gasUsed.GetStorageCost()
		storageRebate := gasUsed.GetStorageRebate()

		totalGasUsed := computationCost + storageCost - storageRebate
		require.Greater(t, totalGasUsed, uint64(0))
		require.Equal(t, totalGasUsed, finalGasBudget)

		require.True(t, resp.Transaction.GetEffects().GetStatus().GetSuccess())
	})
}

// TestCoinSelectionEdgeCases tests edge cases in coin selection logic
//
//nolint:paralleltest
func TestCoinSelectionEdgeCases(t *testing.T) {
	lggr := logger.Test(t)

	// Test case 1: Empty coin list
	t.Run("EmptyCoinList", func(t *testing.T) {
		_, err := txm.SelectCoinsForGasBudget(1000000, []*suirpcv2.Object{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "no coins available")
	})

	// Test case 2: No SUI coins available
	t.Run("NoSUICoins", func(t *testing.T) {
		t.Skip("Invalid test case")
		nonSuiCoins := []*suirpcv2.Object{
			newTestCoin(t, "0x123::other::TOKEN", 1000000000),
		}
		_, err := txm.SelectCoinsForGasBudget(1000000, nonSuiCoins)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no SUI coins available")
	})

	// Test case 3: Insufficient balance
	t.Run("InsufficientBalance", func(t *testing.T) {
		insufficientCoins := []*suirpcv2.Object{
			newTestCoin(t, "0x2::coin::Coin<0x2::sui::SUI>", 500000), // Less than required
		}
		_, err := txm.SelectCoinsForGasBudget(1000000, insufficientCoins)
		require.Error(t, err)
		require.Contains(t, err.Error(), "insufficient funds")
	})

	// Test case 4: Exact balance match
	t.Run("ExactBalanceMatch", func(t *testing.T) {
		exactCoins := []*suirpcv2.Object{
			newTestCoin(t, "0x2::coin::Coin<0x2::sui::SUI>", 1000000), // Exactly what's needed
		}
		selected, err := txm.SelectCoinsForGasBudget(1000000, exactCoins)
		require.NoError(t, err)
		require.Len(t, selected, 1)
	})

	// Test case 5: Multiple coins needed
	t.Run("MultipleCoinsCombined", func(t *testing.T) {
		multipleCoins := []*suirpcv2.Object{
			newTestCoin(t, "0x2::coin::Coin<0x2::sui::SUI>", 600000),
			newTestCoin(t, "0x2::coin::Coin<0x2::sui::SUI>", 500000),
		}
		selected, err := txm.SelectCoinsForGasBudget(1000000, multipleCoins)
		require.NoError(t, err)
		require.Len(t, selected, 2) // Should select both coins

		var totalBalance uint64
		for _, coin := range selected {
			totalBalance += coin.GetBalance()
		}
		require.GreaterOrEqual(t, totalBalance, uint64(1000000))
	})

	lggr.Debugw("Coin selection edge cases test completed")
}
