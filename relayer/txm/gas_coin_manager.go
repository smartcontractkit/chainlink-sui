package txm

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/patrickmn/go-cache"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const (
	DefaultLockedCoinTTL = 24 * time.Hour
	DefaultAllocationTimeout = 30 * time.Second
)

type GasCoinManager interface {
	TryReserveCoins(ctx context.Context, txID string, coinIDs []models.ObjectId) error
	ReleaseCoins(txID string) error
	IsCoinReserved(coinID models.ObjectId) bool
}

// SuiGasCoinManager is the concrete implementation of GasCoinManager.
type SuiGasCoinManager struct {
	lggr   logger.Logger
	client client.SuiPTBClient
	coinsCache *cache.Cache
}

// NewGasCoinManager creates a new SuiGasCoinManager.
func NewGasCoinManager(lggr logger.Logger, suiClient client.SuiPTBClient) *SuiGasCoinManager {
	gcm := &SuiGasCoinManager{
		lggr:      logger.Named(lggr, "SuiGasCoinManager"),
		client:    suiClient,
		coinsCache: cache.New(DefaultAllocationTimeout, DefaultLockedCoinTTL),
	}
	return gcm
}

func (m *SuiGasCoinManager) TryReserveCoins(ctx context.Context, txID string, coinIDs []models.ObjectId) error {
	for _, coinIDBytes := range coinIDs {
		if m.IsCoinReserved(coinIDBytes) {
			return fmt.Errorf("coin %s is already reserved", coinIDBytes)
		}
		
		coinID := hex.EncodeToString(coinIDBytes.Data())
		m.coinsCache.Set(coinID, true, DefaultAllocationTimeout)
	}

	m.coinsCache.Set(txID, coinIDs, DefaultAllocationTimeout)

	return nil
}

func (m *SuiGasCoinManager) ReleaseCoins(txID string) error {
	coinIDs, ok := m.coinsCache.Get(txID)
	if !ok {
		return fmt.Errorf("no coins reserved for transaction %s", txID)
	}

	for _, coinIDBytes := range coinIDs.([]models.ObjectId) {
		coinID := hex.EncodeToString(coinIDBytes.Data())
		m.coinsCache.Delete(coinID)
	}

	m.coinsCache.Delete(txID)
	return nil
}

func (m *SuiGasCoinManager) IsCoinReserved(coinID models.ObjectId) bool {
	coinIDStr := hex.EncodeToString(coinID.Data())
	isReserved, found := m.coinsCache.Get(coinIDStr)
	return found && isReserved.(bool)
}