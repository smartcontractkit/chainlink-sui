package txm

import (
	"context"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
)

func (txm *SuiTxm) confimerLoop(loopCtx context.Context) {
	txm.lggr.Infow("Starting confimer loop")

	_, cancel := services.StopRChan(txm.stopChannel).NewCtx()
	defer cancel()

	basePeriod := txm.configuration.ConfirmerPoolPeriodSeconds
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
			txm.lggr.Infow("Loop context cancelled")
			return
		case <-ticker.C:
			// Check for confirmations
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
			resp, err := txm.suiGateway.GetTransactionStatus(loopCtx, tx.Digest)
			if err != nil {
				txm.lggr.Errorw("Error getting transaction status", "transactionID", tx.TransactionID, "error", err)
				continue
			}

			switch resp.Status {
			case "success":
				err := handleSuccess(txm, tx)
				if err != nil {
					txm.lggr.Errorw("Error handling successful transaction", "transactionID", tx.TransactionID, "error", err)
					continue
				}
			case "failure":
				err := handleTransactionError(txm, tx, &resp)
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

func handleTransactionError(txm *SuiTxm, tx SuiTx, result *client.TransactionResult) error {
	errorMessage := suierrors.ParseSuiErrorMessage(result.Error)

	if suierrors.IsRetryable(errorMessage) {
		txm.lggr.Infow("Transaction retriable, retrying", "transactionID", tx.TransactionID)
		err := txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable)
		if err != nil {
			txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
			return err
		}
	} else {
		txm.lggr.Infow("Transaction is not retryable, marking as failed", "transactionID", tx.TransactionID)
		err := txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed)
		if err != nil {
			txm.lggr.Errorw("Failed to update transaction state", "transactionID", tx.TransactionID, "error", err)
			return err
		}
	}

	return nil
}
