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
	ReserveLockedCoin(coinID models.SuiAddressBytes, ttl time.Duration)
	ReleaseCoins(txID string) error
	IsCoinReserved(coinID models.SuiAddressBytes) bool
}

// lockedCoinKeyPrefix namespaces the cache entries used to exclude coins that the
// chain reported as locked. These entries are deliberately kept separate from the
// txID-associated reservations so that ReleaseCoins (which only unwinds a txID's
// reservation) can never remove a locked-coin exclusion before its TTL expires.
const lockedCoinKeyPrefix = "locked:"

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

// ReserveLockedCoin marks a single coin object as excluded for the given TTL,
// independently of any transaction ID. Unlike TryReserveCoins, the exclusion is
// not associated with a txID and therefore is never removed by ReleaseCoins; it
// persists until the TTL expires. This is used to exclude a coin that the chain
// reported as locked so it is not re-selected while the on-chain lock persists,
// even if that same coin was previously reserved as a payment coin (whose
// reservation is released during the coin-refresh retry).
func (m *SuiGasCoinManager) ReserveLockedCoin(coinID models.SuiAddressBytes, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.coinsCache.Set(lockedCoinCacheKey(coinID), true, ttl)
}

// lockedCoinCacheKey returns the namespaced cache key for a locked-coin exclusion.
func lockedCoinCacheKey(coinID models.SuiAddressBytes) string {
	return lockedCoinKeyPrefix + hex.EncodeToString(coinID[:])
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
// already hold m.mu. A coin is considered reserved if it has either a txID-scoped
// reservation (TryReserveCoins) or a standalone locked-coin exclusion
// (ReserveLockedCoin).
func (m *SuiGasCoinManager) isCoinReserved(coinID models.SuiAddressBytes) bool {
	if isReserved, found := m.coinsCache.Get(hex.EncodeToString(coinID[:])); found {
		if reserved, ok := isReserved.(bool); ok && reserved {
			return true
		}
	}

	if locked, found := m.coinsCache.Get(lockedCoinCacheKey(coinID)); found {
		if isLocked, ok := locked.(bool); ok && isLocked {
			return true
		}
	}

	return false
}
