// broadcaster.go provides functionality for broadcasting transactions to the Sui blockchain.
// It implements a non-blocking broadcast mechanism that can handle multiple transactions
// in a batch while maintaining proper ordering and state management.
package txm

import (
	"context"
	"encoding/base64"
	"sort"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
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
// Parameters:
//   - loopCtx: Context for the broadcast operations, used for cancellation and timeouts
//
// The function never returns until the broadcast channel is closed or the stop signal is received.
func (txm *SuiTxm) broadcastLoop(loopCtx context.Context) {
	defer txm.done.Done()

	_, cancel := services.StopRChan(txm.stopChannel).NewCtx()
	defer cancel()

	for {
		select {
		case <-txm.stopChannel:
			txm.lggr.Infow("Broadcast loop stopped")
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
		txm.lggr.Infow("Submitting transaction", "txID", tx.TransactionID)
		// Process the transaction for broadcasting
		payload := client.TransactionBlockRequest{
			TxBytes:    base64.StdEncoding.EncodeToString(tx.Payload),
			Signatures: tx.Signatures,
			Options: client.TransactionBlockOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
			},
			RequestType: tx.RequestType,
		}

		resp, _ := txm.suiGateway.SendTransaction(loopCtx, payload)

		txm.lggr.Debugw("Transaction response:", resp)
		err := txm.transactionRepository.IncrementAttempts(tx.TransactionID)
		if err != nil {
			txm.lggr.Errorw("Failed to increment transaction attempts", "txID", tx.TransactionID, "error", err)
			continue
		}
		err = txm.transactionRepository.ChangeState(tx.TransactionID, StateSubmitted)
		if err != nil {
			// By default falls back to Retryable. The confirmer routine will make sure to retry
			// the transaction if it is in a retriable state.
			txm.lggr.Errorw("Failed to update transaction state", "txID", tx.TransactionID, "error", err)
			err = txm.transactionRepository.ChangeState(tx.TransactionID, StateRetriable)
			if err != nil {
				txm.lggr.Errorw("Failed to update transaction state to retriable", "txID", tx.TransactionID, "error", err)
				continue
			}

			continue
		}
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
