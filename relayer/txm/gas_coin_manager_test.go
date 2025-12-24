//go:build unit

package txm_test

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
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
	coinID1, err := models.NewHexData("0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)
	coinID2, err := models.NewHexData("0xabcdef1234567890abcdef1234567890abcdef12")
	require.NoError(t, err)
	
	coinIDs := []models.ObjectId{*coinID1, *coinID2}
	
	t.Run("successfully reserve coins", func(t *testing.T) {
		err := gcm.TryReserveCoins(ctx, txID, coinIDs)
		assert.NoError(t, err)
		
		// Verify coins are reserved
		assert.True(t, gcm.IsCoinReserved(*coinID1))
		assert.True(t, gcm.IsCoinReserved(*coinID2))
		
		// Verify transaction is stored
		isReserved := gcm.IsCoinReserved(*coinID1)
		assert.True(t, isReserved)
		isReserved = gcm.IsCoinReserved(*coinID2)
		assert.True(t, isReserved)
	})
	
	t.Run("fail to reserve already reserved coin", func(t *testing.T) {
		// Try to reserve the same coins again
		err := gcm.TryReserveCoins(ctx, "tx-test-2", coinIDs)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is already reserved")
	})

	t.Run("reserved then released coins can be reserved again", func(t *testing.T) {
		coinID1, err := models.NewHexData("0x1234567890abcdef1234123890abcdef12345678")
		require.NoError(t, err)
		coinID2, err := models.NewHexData("0x1234567890abcdef1234153890abcdef12345678")
		require.NoError(t, err)
		
		coinIDs := []models.ObjectId{*coinID1, *coinID2}

		err = gcm.TryReserveCoins(ctx, "tx-test-3", coinIDs)
		assert.NoError(t, err)

		// coins are reserved
		assert.True(t, gcm.IsCoinReserved(*coinID1))
		assert.True(t, gcm.IsCoinReserved(*coinID2))

		// release the coins
		err = gcm.ReleaseCoins("tx-test-3")
		assert.NoError(t, err)

		// coins are not reserved
		assert.False(t, gcm.IsCoinReserved(*coinID1))
		assert.False(t, gcm.IsCoinReserved(*coinID2))

		// try to reserve the coins again
		err = gcm.TryReserveCoins(ctx, "tx-test-3", coinIDs)
		assert.NoError(t, err)
	})
}
