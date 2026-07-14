//go:build unit

package txm_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink-sui/relayer/txm"
)

// mustCoinRef builds a SuiObjectRef from a 0x-prefixed 32-byte address string.
func mustCoinRef(t *testing.T, addr string) transaction.SuiObjectRef {
	t.Helper()
	b, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(addr))
	require.NoError(t, err)
	return transaction.SuiObjectRef{ObjectId: *b, Version: 1, Digest: nil}
}

func TestSuiGasCoinManager_TryReserveCoins(t *testing.T) {
	lggr := logger.Test(t)
	mockClient := &testutils.FakeSuiPTBClient{}
	gcm := txm.NewGasCoinManager(lggr, mockClient)

	ctx := context.Background()
	txID := "test-tx-123"

	// Create test coin IDs
	coinID1 := models.SuiAddress("0x1234567890abcdef1234567890abcdef12345678")
	coinID2 := models.SuiAddress("0xabcdef1234567890abcdef1234567890abcdef12")

	coinID1Bytes, err := transaction.ConvertSuiAddressStringToBytes(coinID1)
	require.NoError(t, err)
	coinID2Bytes, err := transaction.ConvertSuiAddressStringToBytes(coinID2)
	require.NoError(t, err)

	coinIDs := []transaction.SuiObjectRef{
		{
			ObjectId: *coinID1Bytes,
			Version:  1,
			Digest:   nil,
		},
		{
			ObjectId: *coinID2Bytes,
			Version:  1,
			Digest:   nil,
		},
	}

	t.Run("successfully reserve coins", func(t *testing.T) {
		err := gcm.TryReserveCoins(ctx, txID, coinIDs, nil)
		assert.NoError(t, err)

		// Verify coins are reserved
		assert.True(t, gcm.IsCoinReserved(*coinID1Bytes))
		assert.True(t, gcm.IsCoinReserved(*coinID2Bytes))

		// Verify transaction is stored
		isReserved := gcm.IsCoinReserved(*coinID1Bytes)
		assert.True(t, isReserved)
		isReserved = gcm.IsCoinReserved(*coinID2Bytes)
		assert.True(t, isReserved)
	})

	t.Run("fail to reserve already reserved coin", func(t *testing.T) {
		// Try to reserve the same coins again
		err := gcm.TryReserveCoins(ctx, "tx-test-2", coinIDs, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is already reserved")
	})

	t.Run("reserved then released coins can be reserved again", func(t *testing.T) {
		coinID1 := models.SuiAddress("0x1234567890abcdef1234123890abcdef12345678")
		coinID1Bytes, err := transaction.ConvertSuiAddressStringToBytes(coinID1)
		require.NoError(t, err)

		coinID2 := models.SuiAddress("0x1234567890abcdef1234153890abcdef12345678")
		coinID2Bytes, err := transaction.ConvertSuiAddressStringToBytes(coinID2)
		require.NoError(t, err)

		coinIDs := []transaction.SuiObjectRef{
			{
				ObjectId: *coinID1Bytes,
				Version:  1,
				Digest:   nil,
			},
			{
				ObjectId: *coinID2Bytes,
				Version:  1,
				Digest:   nil,
			},
		}

		err = gcm.TryReserveCoins(ctx, "tx-test-3", coinIDs, nil)
		assert.NoError(t, err)

		// coins are reserved
		assert.True(t, gcm.IsCoinReserved(*coinID1Bytes))
		assert.True(t, gcm.IsCoinReserved(*coinID2Bytes))

		// release the coins
		err = gcm.ReleaseCoins("tx-test-3")
		assert.NoError(t, err)

		// coins are not reserved
		assert.False(t, gcm.IsCoinReserved(*coinID1Bytes))
		assert.False(t, gcm.IsCoinReserved(*coinID2Bytes))

		// try to reserve the coins again
		err = gcm.TryReserveCoins(ctx, "tx-test-3", coinIDs, nil)
		assert.NoError(t, err)
	})

	t.Run("coins should be released automatically after the default TTL", func(t *testing.T) {
		coinID1Bytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress("0x1234567890abcdef1234123890abcdef12345123"))
		require.NoError(t, err)
		coinIDs := []transaction.SuiObjectRef{
			{
				ObjectId: *coinID1Bytes,
				Version:  1,
				Digest:   nil,
			},
		}

		err = gcm.TryReserveCoins(ctx, "tx-test-4", coinIDs, nil)
		assert.NoError(t, err)

		// coins should be released automatically after the default TTL (30 seconds)
		require.Eventually(t, func() bool {
			isReserved := gcm.IsCoinReserved(*coinID1Bytes)
			if isReserved {
				fmt.Println("coin is still reserved")
			}

			return !isReserved
		}, 45*time.Second, 10*time.Second)
	})
}

// TestSuiGasCoinManager_AtomicReservation exercises the atomic, all-or-nothing
// behaviour of TryReserveCoins that guards against two transactions reserving the
// same gas coin.
func TestSuiGasCoinManager_AtomicReservation(t *testing.T) {
	lggr := logger.Test(t)
	mockClient := &testutils.FakeSuiPTBClient{}
	ctx := context.Background()

	t.Run("concurrent reservations of the same coin: exactly one wins", func(t *testing.T) {
		gcm := txm.NewGasCoinManager(lggr, mockClient)
		coin := mustCoinRef(t, fmt.Sprintf("0x%064x", 1))

		const numGoroutines = 50
		var wg sync.WaitGroup
		var successes atomic.Int32
		start := make(chan struct{})

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				<-start // release all goroutines at once to maximize contention
				txID := fmt.Sprintf("tx-%d", i)
				if err := gcm.TryReserveCoins(ctx, txID, []transaction.SuiObjectRef{coin}, nil); err == nil {
					successes.Add(1)
				}
			}(i)
		}

		close(start)
		wg.Wait()

		assert.Equal(t, int32(1), successes.Load(), "exactly one caller should reserve the coin")
		assert.True(t, gcm.IsCoinReserved(coin.ObjectId))
	})

	t.Run("reservation is all-or-nothing on partial overlap", func(t *testing.T) {
		gcm := txm.NewGasCoinManager(lggr, mockClient)
		coinA := mustCoinRef(t, fmt.Sprintf("0x%064x", 0xAA))
		coinB := mustCoinRef(t, fmt.Sprintf("0x%064x", 0xBB))

		// Reserve coinA under tx1.
		require.NoError(t, gcm.TryReserveCoins(ctx, "tx1", []transaction.SuiObjectRef{coinA}, nil))

		// tx2 requests [coinB, coinA]; coinA is already taken, so the whole call must
		// fail without leaving coinB reserved.
		err := gcm.TryReserveCoins(ctx, "tx2", []transaction.SuiObjectRef{coinB, coinA}, nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is already reserved")
		assert.False(t, gcm.IsCoinReserved(coinB.ObjectId), "coinB must not be leaked when the reservation fails")

		// The failed reservation left no txID mapping to release.
		assert.Error(t, gcm.ReleaseCoins("tx2"))

		// coinB is still free and can be claimed by another transaction.
		assert.NoError(t, gcm.TryReserveCoins(ctx, "tx3", []transaction.SuiObjectRef{coinB}, nil))
	})

	t.Run("concurrent disjoint reservations all succeed", func(t *testing.T) {
		gcm := txm.NewGasCoinManager(lggr, mockClient)

		const numGoroutines = 50
		// Pre-build coin refs so the goroutines don't call require/t inside them.
		coins := make([]transaction.SuiObjectRef, numGoroutines)
		for i := range coins {
			coins[i] = mustCoinRef(t, fmt.Sprintf("0x%064x", 0x1000+i))
		}

		var wg sync.WaitGroup
		var successes atomic.Int32
		start := make(chan struct{})

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				<-start
				if err := gcm.TryReserveCoins(ctx, fmt.Sprintf("tx-%d", i), []transaction.SuiObjectRef{coins[i]}, nil); err == nil {
					successes.Add(1)
				}
			}(i)
		}

		close(start)
		wg.Wait()

		assert.Equal(t, int32(numGoroutines), successes.Load(), "all disjoint reservations should succeed")
	})
}

// TestSuiGasCoinManager_ReserveLockedCoin verifies that a locked-coin exclusion is
// independent of any txID reservation and survives ReleaseCoins, which is what keeps
// a chain-locked coin out of re-selection during the coin-refresh retry.
func TestSuiGasCoinManager_ReserveLockedCoin(t *testing.T) {
	lggr := logger.Test(t)
	mockClient := &testutils.FakeSuiPTBClient{}
	ctx := context.Background()

	t.Run("exclusion survives ReleaseCoins of the payment reservation", func(t *testing.T) {
		gcm := txm.NewGasCoinManager(lggr, mockClient)
		coin := mustCoinRef(t, fmt.Sprintf("0x%064x", 0xCC))

		// The coin was reserved as a normal payment coin under the transaction's ID.
		require.NoError(t, gcm.TryReserveCoins(ctx, "tx-locked", []transaction.SuiObjectRef{coin}, nil))
		assert.True(t, gcm.IsCoinReserved(coin.ObjectId))

		// The chain reports it as locked; add a standalone exclusion (as handleTransactionError does).
		gcm.ReserveLockedCoin(coin.ObjectId, txm.DefaultLockedCoinTTL)

		// The coin-refresh retry releases the payment reservation for this txID.
		require.NoError(t, gcm.ReleaseCoins("tx-locked"))

		// The standalone exclusion must survive, so the coin is still excluded and
		// cannot be re-selected or re-reserved by any transaction.
		assert.True(t, gcm.IsCoinReserved(coin.ObjectId), "locked coin exclusion must survive ReleaseCoins")
		err := gcm.TryReserveCoins(ctx, "tx-other", []transaction.SuiObjectRef{coin}, nil)
		assert.Error(t, err, "a locked coin must not be reservable by another transaction")
	})

	t.Run("exclusion is created even when the coin was never a payment reservation", func(t *testing.T) {
		gcm := txm.NewGasCoinManager(lggr, mockClient)
		coin := mustCoinRef(t, fmt.Sprintf("0x%064x", 0xDD))

		assert.False(t, gcm.IsCoinReserved(coin.ObjectId))
		gcm.ReserveLockedCoin(coin.ObjectId, txm.DefaultLockedCoinTTL)
		assert.True(t, gcm.IsCoinReserved(coin.ObjectId))

		// There is no txID mapping for a locked-only exclusion, so ReleaseCoins keyed
		// by the coin's own hex finds nothing and the exclusion remains.
		assert.Error(t, gcm.ReleaseCoins(fmt.Sprintf("0x%064x", 0xDD)))
		assert.True(t, gcm.IsCoinReserved(coin.ObjectId))
	})
}
