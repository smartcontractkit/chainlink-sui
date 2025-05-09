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
			// Update the transaction state to Failed if the digest is empty
			// An empty digest indicates a total failure of the transaction
			if resp.TxDigest == "" {
				txm.lggr.Errorw("Transaction failed without a digest", "txID", tx.TransactionID, "function inputs", tx.Functions)
				err = txm.transactionRepository.ChangeState(tx.TransactionID, StateFailed)
				if err != nil {
					txm.lggr.Errorw("Failed to change transaction state to Failed", "txID", tx.TransactionID, "error", err)
				}
			}

			continue
		}
		txm.lggr.Infow("Transaction broadcasted", resp)
		err = txm.transactionRepository.UpdateTransactionDigest(tx.TransactionID, resp.TxDigest)
		if err != nil {
			txm.lggr.Errorw("Failed to update transaction digest", "txID", tx.TransactionID, "error", err)
			continue
		}
		_ = txm.transactionRepository.ChangeState(tx.TransactionID, StateSubmitted)
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
