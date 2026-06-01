// broadcaster.go provides functionality for broadcasting transactions to the Sui blockchain.
// It implements a non-blocking broadcast mechanism that can handle multiple transactions
// in a batch while maintaining proper ordering and state management.
package txm

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	v2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
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
		signatures := []*v2.UserSignature{}
		for _, signature := range tx.Signatures {
			signatureBytes, err := base64.StdEncoding.DecodeString(signature)
			if err != nil {
				txm.lggr.Errorw("Failed to decode signature", "txID", tx.TransactionID, "error", err)
				continue
			}

			signatures = append(signatures, &v2.UserSignature{
				Bcs: &v2.Bcs{Value: signatureBytes},
			})
		}

		payloadBytes, err := base64.StdEncoding.DecodeString(tx.Payload)
		if err != nil {
			txm.lggr.Errorw("Failed to decode payload", "txID", tx.TransactionID, "error", err)
			continue
		}

		payload := &suirpcv2.ExecuteTransactionRequest{
			Transaction: &v2.Transaction{Bcs: &v2.Bcs{Value: payloadBytes}},
			Signatures:  signatures,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{"transaction", "digest", "effects.digest", "effects.status", "effects.gas_used"},
			},
		}

		txm.lggr.Infow("Broadcasting transaction", "txID", tx.TransactionID)

		resp, err := txm.suiGateway.SendTransaction(loopCtx, payload)

		// We increment the attempts here regardless of the error
		// This is because we want to keep track of how many times we tried to broadcast the transaction
		// Even in the case the transaction is malformed (e.g wrong function name)
		attemptErr := txm.transactionRepository.IncrementAttempts(tx.TransactionID)
		if attemptErr != nil {
			txm.lggr.Errorw("Failed to increment transaction attempts", "txID", tx.TransactionID, "error", attemptErr)
			continue
		}

		// If the error is at the tranaction level (submitted but rejected by the node), set it
		// as the error message
		if err == nil && resp.GetTransaction().GetEffects().GetStatus().GetError() != nil {
			err = fmt.Errorf("transaction failed with error: %s", resp.GetTransaction().GetEffects().GetStatus().GetError().GetDescription())
		}

		if err != nil {
			// In the case there is an error submitting
			txm.lggr.Errorw("Failed to broadcast transaction", "txID", tx.TransactionID, "function inputs", tx.Functions, "error", err)

			// Default to retrying the transaction
			newState := StateRetriable

			// Attach the transaction submission error to the transaction object. If it remains marked as retriable,
			// the "confirmer" loop will pick it up and potentially retry it.
			err = txm.transactionRepository.UpdateTransactionBroadcastError(tx.TransactionID, err.Error())
			if err != nil {
				txm.lggr.Errorw("Failed to update transaction broadcast error", "txID", tx.TransactionID, "error", err)
			}

			// Attempt updating the state if it has changed
			if tx.State != newState {
				err = txm.transactionRepository.ChangeState(tx.TransactionID, newState)
				if err != nil {
					txm.lggr.Errorw("Failed to change transaction state", "txID", tx.TransactionID, "error", err)
				}
			}

			continue
		}
		txm.lggr.Infow("Transaction broadcasted", "response", resp, "txID", tx.TransactionID)

		err = txm.transactionRepository.UpdateTransactionDigest(tx.TransactionID, resp.GetTransaction().GetDigest())
		if err != nil {
			txm.lggr.Errorw("Failed to update transaction digest", "txID", tx.TransactionID, "error", err)
			continue
		}

		// Update the transaction state to submitted as we have not yet confirmed its status.
		// The "confirmer" loop checks the transactions statuses and possibly marks them as finalized.
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
