package txm

import (
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

func (txm *SuiTxm) reaperLoop() {
	defer txm.done.Done()
	txm.lggr.Infow("Starting reaper loop")

	loopCtx, cancel := services.StopRChan(txm.stopChannel).NewCtx()
	defer cancel()

	basePeriod := txm.configuration.ReaperPollSecs
	ticker, jitteredDuration := GetTicker(uint(basePeriod))
	defer ticker.Stop()

	txm.lggr.Infow("Created reaper ticker",
		"basePeriod", basePeriod,
		"jitteredDuration", jitteredDuration.String())

	for {
		select {
		case <-txm.stopChannel:
			txm.lggr.Infow("Reaper loop stopped")
			return
		case <-loopCtx.Done():
			txm.lggr.Infow("Loop context cancelled. Reaper loop stopped")
			return
		case <-ticker.C:
			txm.lggr.Debugw("Ticker fired, cleaning up transactions")
			cleanupTransactions(txm)
		}
	}
}

func cleanupTransactions(txm *SuiTxm) {
	txm.lggr.Debugw("Cleaning up transactions")

	// Get all finalized and failed transactions, never in-flight transactions
	states := []TransactionState{StateFinalized, StateFailed}

	finalizedTransactions, err := txm.transactionRepository.GetTransactionsByStates(states)
	if err != nil {
		txm.lggr.Errorw("Error getting finalized transactions", "error", err)
		return
	}

	currentTimestamp := GetCurrentUnixTimestamp()

	if currentTimestamp == 0 {
		txm.lggr.Errorw("Found a 0 timestamp, skipping cleanup")
		return
	}

	for _, tx := range finalizedTransactions {
		txm.lggr.Debugw("Cleaning up finalized transaction", "transactionID", tx.TransactionID)
		timeDiff := currentTimestamp - tx.LastUpdatedAt
		txm.lggr.Debugw("Time difference", "timeDiff", timeDiff)
		if timeDiff > txm.configuration.TransactionRetentionSecs {
			txm.lggr.Debugw("Cleaning up finalized transaction", "transactionID", tx.TransactionID)
			err := txm.transactionRepository.DeleteTransaction(tx.TransactionID)
			if err != nil {
				txm.lggr.Errorw("Error deleting transaction", "transactionID", tx.TransactionID, "error", err)
			}
		}
	}
}
