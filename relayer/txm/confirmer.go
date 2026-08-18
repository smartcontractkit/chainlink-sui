package txm

import (
	"context"
	"math"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
)

const (
	success                               = "success"
	failure                               = "failure"
	defaultExponentialBackoffDelaySeconds = 2
)

func (txm *SuiTxm) confirmerLoop() {
	defer txm.done.Done()
	txm.lggr.Infow("Starting confirmer loop")

	loopCtx, cancel := services.StopRChan(txm.stopChannel).NewCtx()
	defer cancel()

	basePeriod := txm.configuration.ConfirmPollSecs
	ticker, jitteredDuration := GetTicker(basePeriod)
	defer ticker.Stop()

	txm.lggr.Infow("Created confirmer ticker",
		"basePeriod", basePeriod,
		"jitteredDuration", jitteredDuration.String())

	for {
		select {
		case <-txm.stopChannel:
			txm.lggr.Infow("Confirmer loop stopped")
			return
		case <-loopCtx.Done():
			txm.lggr.Infow("Loop context cancelled. Confirmer loop stopped")
			return
		case <-ticker.C:
			txm.lggr.Debugw("Ticker fired, checking transaction confirmations")
			checkConfirmations(loopCtx, txm)
		}
	}
}

func checkConfirmations(loopCtx context.Context, txm *SuiTxm) {
	inFlightTransactions, err := txm.transactionRepository.GetInflightTransactions()
	if err != nil {
		txm.lggr.Errorw("Error getting in-flight transactions", "error", err)
		return
	}

	for _, tx := range inFlightTransactions {
		txm.lggr.Debugw("Checking transaction confirmations", "transactionID", tx.TransactionID)

		var resp client.TransactionResult
		var err error

		if tx.State == StateSubmitted {
			txm.lggr.Debugw("Transaction is in submitted state", "transactionID", tx.TransactionID)
			resp, err = txm.suiGateway.GetTransactionStatus(loopCtx, tx.Digest)
			if err != nil {
				txm.lggr.Errorw("Error getting transaction status", "transactionID", tx.TransactionID, "error", err)
				continue
			}
		} else if tx.State == StateRetriable {
			txm.lggr.Debugw("Transaction is in retriable state", "transactionID", tx.TransactionID)
			// Check if it's a broadcast error (never made it onchain)
			if tx.BroadcastError == "" {
				continue
			}
			resp.Status = failure
			resp.Error = tx.BroadcastError
		} else {
			txm.lggr.Debugw("Transaction is in not submitted or retriable state", "transactionID", tx.TransactionID, "state", tx.State)
			continue
		}

		switch resp.Status {
		case success:
			if err := handleSuccess(loopCtx, txm, tx); err != nil {
				txm.lggr.Errorw("Error handling successful transaction", "transactionID", tx.TransactionID, "error", err)
			}
		case failure:
			if err := handleTransactionError(loopCtx, txm, tx, &resp); err != nil {
				txm.lggr.Errorw("Error handling failed transaction", "transactionID", tx.TransactionID, "error", err)
			}
		default:
			txm.lggr.Infow("Unknown transaction status", "transactionID", tx.TransactionID, "status", resp.Status)
		}
	}
}

func handleSuccess(ctx context.Context, txm *SuiTxm, tx SuiTx) error {
	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateFinalized); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}
	txm.lggr.Infow("Transaction finalized", "transactionID", tx.TransactionID)

	// Record successful transaction in health metrics
	txm.recordLastSuccess(ctx)

	if err := tx.CoinManager.ReleaseCoins(tx.TransactionID); err != nil {
		// This error is not critical, can be safely ignored as the coins will auto-release after the default TTL
		txm.lggr.Debugw("Failed to release coins", "transactionID", tx.TransactionID, "error", err)
	}

	return nil
}

func handleTransactionError(ctx context.Context, txm *SuiTxm, tx SuiTx, result *client.TransactionResult) error {
	txm.lggr.Debugw("Handling transaction error", "transactionID", tx.TransactionID, "error", result.Error)

	txError := suierrors.ParseSuiErrorMessage(result.Error)

	// Check if the error is a locked object error, and exclude the coin so other
	// transactions (and this transaction's own coin-refresh retry) do not re-select it.
	if objectID, version, ok := suierrors.ExtractLockedObjectRef(result.Error); ok {
		txm.lggr.Infow("Detected locked coin at confirmation time",
			"txID", tx.TransactionID,
			"objectID", objectID,
			"version", version,
		)

		coinID, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(objectID))
		if err != nil {
			txm.lggr.Debugw(
				"Failed to convert locked coin object ID",
				"transactionID", tx.TransactionID,
				"objectID", objectID,
				"error", err,
			)
		} else {
			// Use a standalone, long-lived exclusion keyed by the coin itself (not by
			// txID). The upcoming coin-refresh retry calls ReleaseCoins(txID), which
			// would otherwise wipe a txID-scoped reservation for this coin; the
			// standalone exclusion survives that release and keeps the locked coin out
			// of re-selection until its TTL expires.
			tx.CoinManager.ReserveLockedCoin(*coinID, DefaultLockedCoinTTL)
		}
	}

	isRetryable, strategy := txm.retryManager.IsRetryable(&tx, result.Error)
	if !isRetryable {
		return markTransactionFailed(txm, tx, txError)
	}

	txm.lggr.Infow("Transaction is retryable", "transactionID", tx.TransactionID, "strategy", strategy)

	switch strategy {
	case ExponentialBackoff:
		return handleExponentialBackoffRetry(txm, tx)
	case GasBump:
		return handleGasBumpRetry(ctx, txm, tx, txError)
	case CoinRefresh:
		return handleCoinRefreshRetry(ctx, txm, tx, txError)
	case NoRetry:
		return markTransactionFailed(txm, tx, txError)
	default:
		return markTransactionFailed(txm, tx, txError)
	}
}

func handleGasBumpRetry(ctx context.Context, txm *SuiTxm, tx SuiTx, txError *suierrors.SuiError) error {
	txm.lggr.Infow("Gas bump strategy", "transactionID", tx.TransactionID)

	updatedGas, err := txm.gasManager.GasBump(ctx, &tx)
	if err != nil {
		txm.lggr.Errorw("Failed to bump gas, marking transaction as failed", "transactionID", tx.TransactionID, "error", err)
		if stateErr := txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed); stateErr != nil {
			txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", stateErr)
		}
		if txErrErr := txm.transactionRepository.UpdateTransactionError(tx.TransactionID, txError); txErrErr != nil {
			txm.lggr.Errorw("Failed to update transaction error", "transactionID", tx.TransactionID, "error", txErrErr)
		}
		if releaseErr := tx.CoinManager.ReleaseCoins(tx.TransactionID); releaseErr != nil {
			// Not critical - coins auto-release after TTL.
			txm.lggr.Debugw("Failed to release coins after gas bump failure", "transactionID", tx.TransactionID, "error", releaseErr)
		}
		return err
	}

	// Release the current reservation before the transaction is rebuilt. UpdateTransactionGas
	// re-selects and re-reserves coins via preparePTBTransaction; without releasing first, the
	// rebuild would exclude this transaction's own still-reserved coins and either pick a
	// different set (orphaning the old reservation until its TTL) or fail if no other coins are
	// free. Releasing lets the rebuild re-select the same coins with the bumped budget.
	if err := tx.CoinManager.ReleaseCoins(tx.TransactionID); err != nil {
		// Not critical - coins auto-release after TTL.
		txm.lggr.Debugw("Failed to release coins before gas bump rebuild", "transactionID", tx.TransactionID, "error", err)
	}

	if err := txm.transactionRepository.UpdateTransactionGas(ctx, txm.keystoreService, txm.suiGateway, tx.TransactionID, &updatedGas); err != nil {
		txm.lggr.Errorw("Failed to update transaction gas", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	clearBroadcastError(txm, tx.TransactionID)
	txm.broadcastChannel <- tx.TransactionID
	return nil
}

func handleCoinRefreshRetry(ctx context.Context, txm *SuiTxm, tx SuiTx, txError *suierrors.SuiError) error {
	txm.lggr.Infow("Coin refresh strategy - refreshing coins for locked coin error", "transactionID", tx.TransactionID)

	// Release the old coins that are locked
	if err := tx.CoinManager.ReleaseCoins(tx.TransactionID); err != nil {
		// This is not critical - coins will auto-release after TTL
		txm.lggr.Debugw("Failed to release old coins", "transactionID", tx.TransactionID, "error", err)
	}

	// Get the current transaction to ensure we have the latest state
	currentTx, err := txm.transactionRepository.GetTransaction(tx.TransactionID)
	if err != nil {
		txm.lggr.Errorw("Failed to get current transaction", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	// Calling UpdateTransactionGas will also update the gas coins used as the transaction gets re-built
	// with new (unlocked) coins.
	// Call chain: UpdateTransactionGas -> UpdateBSCPayload -> preparePTBTransaction (this refreshes the coins).
	if err := txm.transactionRepository.UpdateTransactionGas(ctx, txm.keystoreService, txm.suiGateway, tx.TransactionID, currentTx.Metadata.GasLimit); err != nil {
		txm.lggr.Errorw("Failed to update transaction with refreshed coins", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	txm.lggr.Infow("Transaction refreshed with new coins", "transactionID", tx.TransactionID)
	clearBroadcastError(txm, tx.TransactionID)
	txm.broadcastChannel <- tx.TransactionID
	return nil
}

func handleExponentialBackoffRetry(txm *SuiTxm, tx SuiTx) error {
	delaySeconds := float64(defaultExponentialBackoffDelaySeconds) * math.Pow(2, float64(tx.Attempt))

	txm.lggr.Infow("Exponential backoff strategy", "transactionID", tx.TransactionID, "delay", delaySeconds, "state", tx.State)

	// Check if enough time has elapsed since the last update
	timeElapsed := time.Since(time.Unix(int64(tx.LastUpdatedAt), 0))
	if timeElapsed.Seconds() < delaySeconds {
		// Not enough time has elapsed for the next retry
		txm.lggr.Debugw("Not enough time elapsed, no need to retry", "transactionID", tx.TransactionID, "elapsed", timeElapsed, "required", delaySeconds)
		return nil
	}

	txm.lggr.Debugw("Sufficient time elapsed, retrying transaction", "transactionID", tx.TransactionID, "elapsed", timeElapsed)

	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	clearBroadcastError(txm, tx.TransactionID)
	txm.broadcastChannel <- tx.TransactionID
	return nil
}

// clearBroadcastError acknowledges a stored broadcast failure once a retry has been
// scheduled. While BroadcastError is empty the confirmer skips the transaction, so the
// same stored failure cannot be re-handled (and re-enqueued) on every tick while the
// re-broadcast is queued or in flight. The broadcaster sets a fresh BroadcastError if
// the retry fails again. A failure to clear is not critical: it only means one extra
// retry pass may be handled.
func clearBroadcastError(txm *SuiTxm, transactionID string) {
	if err := txm.transactionRepository.UpdateTransactionBroadcastError(transactionID, ""); err != nil {
		txm.lggr.Errorw("Failed to clear broadcast error", "transactionID", transactionID, "error", err)
	}
}

func markTransactionFailed(txm *SuiTxm, tx SuiTx, txError *suierrors.SuiError) error {
	txm.lggr.Infow("Transaction is not retriable, marking as failed", "transactionID", tx.TransactionID)

	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	if err := txm.transactionRepository.UpdateTransactionError(tx.TransactionID, txError); err != nil {
		txm.lggr.Errorw("Failed to update transaction error", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	txm.lggr.Infow("Transaction failed", "transactionID", tx.TransactionID)

	if err := tx.CoinManager.ReleaseCoins(tx.TransactionID); err != nil {
		// This error is not critical, can be safely ignored as the coins will auto-release after the default TTL
		txm.lggr.Debugw("Failed to release coins", "transactionID", tx.TransactionID, "error", err)
	}

	return nil
}
