//go:build unit

package txm_test

import (
	"context"
	"fmt"
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
