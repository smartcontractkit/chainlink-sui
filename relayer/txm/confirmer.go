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
			continue
		}

		switch resp.Status {
		case success:
			if err := handleSuccess(txm, tx); err != nil {
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

func handleSuccess(txm *SuiTxm, tx SuiTx) error {
	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateFinalized); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}
	txm.lggr.Infow("Transaction finalized", "transactionID", tx.TransactionID)

	if err := txm.coinManager.ReleaseCoins(tx.TransactionID); err != nil {
		// This error is not critical, can be safely ignored as the coins will auto-release after the default TTL
		txm.lggr.Debugw("Failed to release coins", "transactionID", tx.TransactionID, "error", err)
	}

	return nil
}

func handleTransactionError(ctx context.Context, txm *SuiTxm, tx SuiTx, result *client.TransactionResult) error {
	txm.lggr.Debugw("Handling transaction error", "transactionID", tx.TransactionID, "error", result.Error)

	txError := suierrors.ParseSuiErrorMessage(result.Error)

	// Check if the error is a locked object error, mark the coin as reserved if it is not already
	// to avoid other transactions from using it
	if objectID, version, ok := suierrors.ExtractLockedObjectRef(result.Error); ok {
		txm.lggr.Infow("Detected locked coin at confirmation time",
			"txID", tx.TransactionID,
			"objectID", objectID,
			"version", version,
		)

		coinID, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(objectID))
		if err == nil && !txm.coinManager.IsCoinReserved(*coinID) {
			// Coin lock duration
			expiry := DefaultLockedCoinTTL

			// The coin is not recorded is not marked as reserved, mark it as reserved
			err = txm.coinManager.TryReserveCoins(ctx, tx.TransactionID, []transaction.SuiObjectRef{
				{
					ObjectId: *coinID,
					Version:  0,
					Digest:   nil,
				},
			}, &expiry)

			if err != nil {
				// This is not a critical error, so we continue
				txm.lggr.Debugw(
					"Failed to mark locked coin as reserved",
					"transactionID", tx.TransactionID,
					"objectID", objectID,
					"error", err,
				)
			}
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
		return err
	}

	if err := txm.transactionRepository.UpdateTransactionGas(ctx, txm.keystoreService, txm.suiGateway, tx.TransactionID, &updatedGas); err != nil {
		txm.lggr.Errorw("Failed to update transaction gas", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	txm.broadcastChannel <- tx.TransactionID
	return nil
}

func handleCoinRefreshRetry(ctx context.Context, txm *SuiTxm, tx SuiTx, txError *suierrors.SuiError) error {
	txm.lggr.Infow("Coin refresh strategy - refreshing coins for locked coin error", "transactionID", tx.TransactionID)

	// Release the old coins that are locked
	if err := txm.coinManager.ReleaseCoins(tx.TransactionID); err != nil {
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
	txm.broadcastChannel <- tx.TransactionID
	return nil
}

func handleExponentialBackoffRetry(txm *SuiTxm, tx SuiTx) error {
	delaySeconds := float64(defaultExponentialBackoffDelaySeconds) * math.Pow(2, float64(tx.Attempt))

	txm.lggr.Infow("Exponential backoff strategy", "transactionID", tx.TransactionID, "delay", delaySeconds, "state", tx.State)

	// Check if enough time has elapsed since the last update
	timeElapsed := time.Since(time.Unix(int64(tx.LastUpdatedAt), 0))
	if timeElapsed.Seconds() < delaySeconds {
		// Not enough time has elapsed for the next retry, mark the transaction as failed
		txm.lggr.Debugw("Not enough time elapsed, no need to retry", "transactionID", tx.TransactionID, "elapsed", timeElapsed, "required", delaySeconds)
		return nil
	}

	txm.lggr.Debugw("Sufficient time elapsed, retrying transaction", "transactionID", tx.TransactionID, "elapsed", timeElapsed)

	if err := txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable); err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}

	txm.broadcastChannel <- tx.TransactionID
	return nil
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

	if err := txm.coinManager.ReleaseCoins(tx.TransactionID); err != nil {
		// This error is not critical, can be safely ignored as the coins will auto-release after the default TTL
		txm.lggr.Debugw("Failed to release coins", "transactionID", tx.TransactionID, "error", err)
	}

	return nil
}
