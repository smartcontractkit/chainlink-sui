package txm

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
)

const (
	success = "success"
	failure = "failure"
)

// confirmerLoop is the main goroutine responsible for monitoring and confirming transactions
// that have been submitted to the Sui blockchain.
//
// The function runs on a periodic ticker (with jitter) and:
// 1. Retrieves all in-flight transactions from the repository
// 2. For each transaction in the submitted state, checks its status on-chain
// 3. Updates the transaction state based on the confirmation status
// 4. Handles retries and failures according to configured policies
//
// The loop continues until either:
// - The stop channel is closed
// - The context is cancelled
// - The service is shut down
//
// Parameters:
// - Uses the txm.configuration.ConfirmPollSecs for the base ticker period
// - Uses txm.stopChannel for shutdown signaling
// - Uses txm.done WaitGroup for cleanup
//
// The function never returns until a shutdown signal is received.

// checkConfirmations processes a batch of in-flight transactions and updates their
// confirmation status.
//
// For each transaction, it:
// 1. Retrieves the current status from the Sui blockchain
// 2. Updates the transaction state in the repository based on the response
// 3. Handles any errors or retries needed
//
// Parameters:
// - loopCtx: Context for cancellation and timeouts
// - txm: The transaction manager instance containing configuration and dependencies
//
// The function logs errors but continues processing remaining transactions if one fails.

func (txm *SuiTxm) confirmerLoop() {
	defer txm.done.Done()
	txm.lggr.Infow("Starting confimer loop")

	loopCtx, cancel := services.StopRChan(txm.stopChannel).NewCtx()
	defer cancel()

	basePeriod := txm.configuration.ConfirmPollSecs
	// Create initial ticker with jitter
	ticker, jitteredDuration := GetTicker(basePeriod)

	defer ticker.Stop()

	txm.lggr.Infow("Created confirmer ticker",
		"basePeriod", basePeriod,
		"jitteredDuration", jitteredDuration.String())

	// Loop to check for confirmations
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
		switch tx.State {
		case StateSubmitted:
			txm.lggr.Debugw("Transaction is in submitted state", "transactionID", tx.TransactionID)
			resp, err := txm.suiGateway.GetTransactionStatus(loopCtx, tx.Digest)
			if err != nil {
				txm.lggr.Errorw("Error getting transaction status", "transactionID", tx.TransactionID, "error", err)
				continue
			}

			switch resp.Status {
			case success:
				err := handleSuccess(txm, tx)
				if err != nil {
					txm.lggr.Errorw("Error handling successful transaction", "transactionID", tx.TransactionID, "error", err)
					continue
				}
			case failure:
				err := handleTransactionError(loopCtx, txm, tx, &resp)
				if err != nil {
					continue
				}
			default:
				txm.lggr.Infow("Unknown transaction status", "transactionID", tx.TransactionID, "status", resp.Status)
			}
		case StateRetriable:
			// TODO check if the transaction is still retriable
			txm.lggr.Debugw("Transaction is still retriable")
		case StatePending, StateFinalized, StateFailed:
			// Do nothing for pending, finalized and failed transactions
		}
	}
}

func handleSuccess(txm *SuiTxm, tx SuiTx) error {
	err := txm.transactionRepository.ChangeState(tx.TransactionID, StateFinalized)
	if err != nil {
		txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
		return err
	}
	txm.lggr.Infow("Transaction finalized", "transactionID", tx.TransactionID)

	return nil
}

func handleTransactionError(ctx context.Context, txm *SuiTxm, tx SuiTx, result *client.TransactionResult) error {
	txm.lggr.Debugw("Handling transaction error", "transactionID", tx.TransactionID, "error", result.Error)
	isRetryable, strategy := txm.retryManager.IsRetryable(&tx, result.Error)
	txError := suierrors.ParseSuiErrorMessage(result.Error)

	if txError == nil {
		txm.lggr.Errorw("Failed to parse transaction error", "transactionID", tx.TransactionID, "error", result.Error)
		return errors.New("failed to parse transaction error")
	}

	if isRetryable {
		txm.lggr.Infow("Transaction is retriable", "transactionID", tx.TransactionID, "strategy", strategy)
		switch strategy {
		case ExponentialBackoff:
			// TODO: for another PR implement exponential backoff
			txm.lggr.Infow("Exponential backoff strategy not implemented")
		case GasBump:
			updatedGas, err := txm.gasManager.GasBump(ctx, &tx)
			if err != nil {
				txm.lggr.Errorw("Failed to bump gas", "transactionID", tx.TransactionID, "error", err)
				err = txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed)
				if err != nil {
					txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
				}
				err = txm.transactionRepository.UpdateTransactionError(tx.TransactionID, txError)
				if err != nil {
					txm.lggr.Errorw("Failed to update transaction error", "transactionID", tx.TransactionID, "error", err)
				}

				return nil
			}

			err = txm.transactionRepository.UpdateTransactionGas(tx.TransactionID, &updatedGas)
			if err != nil {
				txm.lggr.Errorw("Failed to update transaction gas", "transactionID", tx.TransactionID, "error", err)
				return err
			}
			err = txm.transactionRepository.IncrementAttempts(tx.TransactionID)
			if err != nil {
				txm.lggr.Errorw("Failed to increment transaction attempts", "transactionID", tx.TransactionID, "error", err)
				return nil
			}

			// Change state to retriable
			err = txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable)
			if err != nil {
				txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
				return err
			}

			// Re-enqueue
			txm.broadcastChannel <- tx.TransactionID
		case NoRetry:
			txm.lggr.Infow("Transaction is not retriable", "transactionID", tx.TransactionID, "error", result.Error)
			err := txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed)
			if err != nil {
				txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
				return err
			}
		}
	} else {
		txm.lggr.Infow("Transaction is not retriable", "transactionID", tx.TransactionID, "result", result)
		err := txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed)
		if err != nil {
			txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
			return err
		}
		err = txm.transactionRepository.UpdateTransactionError(tx.TransactionID, txError)
		if err != nil {
			txm.lggr.Errorw("Failed to update transaction error", "transactionID", tx.TransactionID, "error", err)
			return err
		}
		txm.lggr.Infow("Transaction failed", "transactionID", tx.TransactionID)
	}

	return nil
}
