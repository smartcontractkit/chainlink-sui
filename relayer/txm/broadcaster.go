// broadcaster.go provides functionality for broadcasting transactions to the Sui blockchain.
// It implements a non-blocking broadcast mechanism that can handle multiple transactions
// in a batch while maintaining proper ordering and state management.
package txm

import (
	"context"
	"sort"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suierrors"
)

// broadcastLoop is the main goroutine responsible for processing transactions from the broadcast channel
// and submitting them to the Sui blockchain.
//
// The function continuously monitors a channel for transaction IDs that need to be broadcast.
// When a transaction ID is received, it:
// 1. Gathers any additional transaction IDs waiting in the channel (without blocking)
// 2. Retrieves the corresponding transaction objects from the repository
// 3. Submits the transactions to the blockchain in order of their timestamp
// 4. Updates their state to reflect the submission status
//
// The loop also handles graceful shutdown through the stopChannel and properly
// cleans up resources using the done WaitGroup.
//
// The function never returns until the broadcast channel is closed or the stop signal is received.
func (txm *SuiTxm) broadcastLoop() {
	defer txm.done.Done()
	txm.lggr.Infow("Starting broadcast loop")

	loopCtx, cancel := services.StopRChan(txm.stopChannel).NewCtx()
	defer cancel()

	for {
		select {
		case <-txm.stopChannel:
			txm.lggr.Infow("Broadcast loop stopped")
			return
		case <-loopCtx.Done():
			txm.lggr.Infow("Loop context cancelled. Broadcast loop stopped")
			return
		case initialId, ok := <-txm.broadcastChannel:
			// Check if the channel is closed
			if !ok {
				txm.lggr.Infow("Broadcast channel closed")
				return
			}
			broadcastIds := getAllBroadcastIds(initialId, txm.broadcastChannel)

			txm.lggr.Infow("Broadcasting transactions", "ids", broadcastIds)
			transactions := getInflightTransactions(txm, broadcastIds)
			broadcastTransactions(loopCtx, txm, transactions)
		}
	}
}

func broadcastTransactions(loopCtx context.Context, txm *SuiTxm, transactions []SuiTx) {
	for _, tx := range transactions {
		// Process the transaction for broadcasting
		payload := client.TransactionBlockRequest{
			TxBytes:    tx.Payload,
			Signatures: tx.Signatures,
			Options: client.TransactionBlockOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
				ShowEvents:         true,
			},
			RequestType: tx.RequestType,
		}

		txm.lggr.Infow("Broadcasting transaction", "txID", tx.TransactionID, "payload", tx)

		resp, err := txm.suiGateway.SendTransaction(loopCtx, payload)

		// We increment the attempts here regardless of the error
		// This is because we want to keep track of how many times we tried to broadcast the transaction
		// Even in the case the transaction is malformed (e.g wrong function name)
		attemptErr := txm.transactionRepository.IncrementAttempts(tx.TransactionID)
		if attemptErr != nil {
			txm.lggr.Errorw("Failed to increment transaction attempts", "txID", tx.TransactionID, "error", attemptErr)
			continue
		}
		if err != nil {
			// In the case there is an error submitting
			txm.lggr.Errorw("Failed to broadcast transaction", "txID", tx.TransactionID, "function inputs", tx.Functions, "error", err)

			// special-case equivocation / locked coin
			if handled := handleLockCoinError(txm, tx, err.Error()); handled {
				// All state/TxError/lockedCoins already handled by helper.
				continue
			}

			// Default to retrying the transaction
			newState := StateRetriable

			if resp.Effects.Status.Status != "" && resp.TxDigest == "" {
				// Update the transaction state to Failed if the digest is empty
				// An empty digest indicates a total failure of the transaction
				txm.lggr.Errorw("Transaction failed without a digest", "txID", tx.TransactionID, "function inputs", tx.Functions)
				newState = StateFailed
			}

			err = txm.transactionRepository.ChangeState(tx.TransactionID, newState)
			if err != nil {
				txm.lggr.Errorw("Failed to change transaction state", "txID", tx.TransactionID, "error", err)
			}

			continue
		}
		txm.lggr.Infow("Transaction broadcasted", "response", resp, "txID", tx.TransactionID)

		err = txm.transactionRepository.UpdateTransactionDigest(tx.TransactionID, resp.TxDigest)
		if err != nil {
			txm.lggr.Errorw("Failed to update transaction digest", "txID", tx.TransactionID, "error", err)
			continue
		}

		err = txm.transactionRepository.ChangeState(tx.TransactionID, StateSubmitted)
		if err != nil {
			txm.lggr.Errorw("Failed to change transaction state to Submitted", "txID", tx.TransactionID, "error", err)
			continue
		}

		txm.lggr.Infow("Transaction state updated to Submitted", "txID", tx.TransactionID)
	}
}

func getInflightTransactions(txm *SuiTxm, broadcastIds []string) []SuiTx {
	transactions := []SuiTx{}
	for _, id := range broadcastIds {
		tx, err := txm.transactionRepository.GetTransaction(id)
		if err != nil {
			txm.lggr.Errorw("Failed to get transaction", "txID", id, "error", err)
			continue
		}
		transactions = append(transactions, tx)
	}

	// ensure older transactions are broadcast first
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Timestamp < transactions[j].Timestamp
	})

	return transactions
}

func getAllBroadcastIds(initalId string, channel chan string) []string {
	broadcastIds := []string{initalId}
	// read all available ids on broadcastChan without blocking, and broadcast in order of which they were
	// queued. this means that retries would take priority over newly submitted transactions.

DrainChannel:
	for {
		select {
		case nextId := <-channel:
			broadcastIds = append(broadcastIds, nextId)
		default:
			break DrainChannel
		}
	}

	// Get all the broadcast IDs
	return broadcastIds
}

// handleLockCoinError inspects the error message, and if it is a
// locked-coin / equivocation error returns true, false otherwise.
func handleLockCoinError(txm *SuiTxm, tx SuiTx, msg string) bool {
	txErr := suierrors.ParseSuiErrorMessage(msg)
	if txErr == nil || txErr.Category != suierrors.LockCoinErrors {
		return false
	}

	if objID, ver, ok := suierrors.ExtractLockedObjectRef(msg); ok {
		txm.lggr.Infow("Detected locked coin at broadcast time",
			"txID", tx.TransactionID,
			"objectID", objID,
			"version", ver,
		)
		txm.markLockedCoin(objID, ver)
	}

	// From Pending -> Failed (allowed)
	if err2 := txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed); err2 != nil {
		txm.lggr.Errorw("Failed to change transaction state", "txID", tx.TransactionID, "error", err2)
	}
	if err2 := txm.transactionRepository.UpdateTransactionError(tx.TransactionID, txErr); err2 != nil {
		txm.lggr.Errorw("Failed to update transaction error", "txID", tx.TransactionID, "error", err2)
	}

	txm.lggr.Errorw("Non-retriable locked-coin error at broadcast",
		"txID", tx.TransactionID,
		"error", msg,
	)

	return true
}
