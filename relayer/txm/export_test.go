package txm

import "time"

var (
	HandleLockCoinError   = handleLockCoinError
	PreparePTBTransaction = preparePTBTransaction
	CoinKey               = coinKey
)

func (txm *SuiTxm) SnapshotLockedCoins() map[string]struct{} {
	return txm.snapshotLockedCoins()
}

func (txm *SuiTxm) MarkLockedCoin(objectID string, version uint64, currTime time.Time) {
	txm.markLockedCoin(objectID, version, currTime)
}
