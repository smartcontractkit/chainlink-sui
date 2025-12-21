package txm

var (
	HandleLockCoinError   = handleLockCoinError
	PreparePTBTransaction = preparePTBTransaction
)

func (txm *SuiTxm) SnapshotLockedCoins() map[string]struct{} {
	return txm.snapshotLockedCoins()
}
