package txm

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
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
	// mu guards the whole check-then-reserve sequence so that reservations are
	// atomic and all-or-nothing across the full coin set.
	mu sync.Mutex
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
	m.mu.Lock()
	defer m.mu.Unlock()

	// First pass: verify the entire set is available before writing anything.
	// Holding the lock across the check and the writes closes the check-then-set
	// race between concurrent callers, and validating upfront makes the
	// reservation all-or-nothing (no partial reservations leaked on failure).
	for _, coin := range coinIDs {
		if m.isCoinReserved(coin.ObjectId) {
			return fmt.Errorf("coin %s is already reserved", hex.EncodeToString(coin.ObjectId[:]))
		}
	}

	expiresAt := DefaultAllocationTimeout
	if expiry != nil {
		expiresAt = *expiry
	}

	// Second pass: reserve the full set now that every coin is known to be free.
	for _, coin := range coinIDs {
		coinID := hex.EncodeToString(coin.ObjectId[:])
		m.coinsCache.Set(coinID, true, expiresAt)
	}

	// Track the reserved set under the txID using the same TTL so ReleaseCoins can
	// always unwind the individual coin entries for the reservation's lifetime.
	m.coinsCache.Set(txID, coinIDs, expiresAt)

	return nil
}

// ReleaseCoins only releases reservations stored under a txID key (txID -> []SuiObjectRef).
// It does not work with coinID keys (coinID -> bool) and cannot unlock those entries directly.
func (m *SuiGasCoinManager) ReleaseCoins(txID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

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
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.isCoinReserved(coinID)
}

// isCoinReserved is the lock-free implementation of IsCoinReserved. Callers must
// already hold m.mu.
func (m *SuiGasCoinManager) isCoinReserved(coinID models.SuiAddressBytes) bool {
	coinIDStr := hex.EncodeToString(coinID[:])
	isReserved, found := m.coinsCache.Get(coinIDStr)
	return found && isReserved.(bool)
}
