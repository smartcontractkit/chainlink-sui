package txm

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/patrickmn/go-cache"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const (
	DefaultLockedCoinTTL     = 24 * time.Hour
	DefaultAllocationTimeout = 30 * time.Second
)

type GasCoinManager interface {
	TryReserveCoins(ctx context.Context, txID string, paymentCoins []transaction.SuiObjectRef, expiry *time.Duration) error
	ReleaseCoins(txID string) error
	IsCoinReserved(coinID models.SuiAddressBytes) bool
}

// SuiGasCoinManager is the concrete implementation of GasCoinManager.
type SuiGasCoinManager struct {
	lggr       logger.Logger
	client     client.SuiPTBClient
	coinsCache *cache.Cache
}

// NewGasCoinManager creates a new SuiGasCoinManager.
func NewGasCoinManager(lggr logger.Logger, suiClient client.SuiPTBClient) *SuiGasCoinManager {
	gcm := &SuiGasCoinManager{
		lggr:       logger.Named(lggr, "SuiGasCoinManager"),
		client:     suiClient,
		coinsCache: cache.New(DefaultAllocationTimeout, DefaultLockedCoinTTL),
	}
	return gcm
}

func (m *SuiGasCoinManager) TryReserveCoins(
	ctx context.Context,
	txID string,
	coinIDs []transaction.SuiObjectRef,
	expiry *time.Duration,
) error {
	for _, coin := range coinIDs {
		if m.IsCoinReserved(coin.ObjectId) {
			return fmt.Errorf("coin %s is already reserved", hex.EncodeToString(coin.ObjectId[:]))
		}

		coinID := hex.EncodeToString(coin.ObjectId[:])
		expiresAt := DefaultAllocationTimeout

		if expiry != nil {
			expiresAt = *expiry
		}

		m.coinsCache.Set(coinID, true, expiresAt)
	}

	m.coinsCache.Set(txID, coinIDs, DefaultAllocationTimeout)

	return nil
}

// ReleaseCoins only releases reservations stored under a txID key (txID -> []SuiObjectRef).
// It does not work with coinID keys (coinID -> bool) and cannot unlock those entries directly.
func (m *SuiGasCoinManager) ReleaseCoins(txID string) error {
	coinIDs, ok := m.coinsCache.Get(txID)
	if !ok {
		return fmt.Errorf("no coins reserved for transaction %s", txID)
	}

	for _, coin := range coinIDs.([]transaction.SuiObjectRef) {
		coinID := hex.EncodeToString(coin.ObjectId[:])
		m.coinsCache.Delete(coinID)
	}

	m.coinsCache.Delete(txID)
	return nil
}

func (m *SuiGasCoinManager) IsCoinReserved(coinID models.SuiAddressBytes) bool {
	coinIDStr := hex.EncodeToString(coinID[:])
	isReserved, found := m.coinsCache.Get(coinIDStr)
	return found && isReserved.(bool)
}
