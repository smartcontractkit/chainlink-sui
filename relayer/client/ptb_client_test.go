//go:build integration

package client_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/utils"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestPTBClient(t *testing.T) {
	log := logger.Test(t)

	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	maxConcurrent := int64(3)
	relayerClient, err := client.NewPTBClient(log, testutils.LocalUrl, nil, 120*time.Second, keystoreInstance, maxConcurrent, "WaitForLocalExecution")
	require.NoError(t, err)

	err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, publishOutput, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// Test GetLatestValue for different data types
	//nolint:paralleltest
	t.Run("FunctionRead", func(t *testing.T) {
		args := []any{counterObjectId}
		argTypes := []string{"objectId"}

		response, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_count",
			args,
			argTypes,
		)
		require.NoError(t, err)
		require.NotNil(t, response)
		utils.PrettyPrint(response)
	})

	//nolint:paralleltest
	t.Run("WithRateLimit", func(t *testing.T) {
		// Block operations with channel to observe concurrency
		completionCh := make(chan int, 50) // Buffer large enough for all completions

		// Block until manual release to measure concurrency precisely
		// This ensures we can observe exactly how many goroutines acquired the semaphore
		blockingOperation := func(id int) {
			// Make request that will block
			ctx := context.Background()
			go func() {
				defer func() {
					completionCh <- id // Signal this request completed
				}()

				err := relayerClient.WithRateLimit(ctx, func(ctx context.Context) error {
					time.Sleep(1 * time.Second)
					return nil
				})
				require.NoError(t, err)
			}()
		}

		// Start more requests than our concurrency limit
		numRequests := 100
		for i := range numRequests {
			blockingOperation(i)
		}

		// Wait a moment to ensure requests have time to acquire semaphore
		time.Sleep(500 * time.Millisecond)

		// Count how many completed without unblocking
		completeCount := 0
	countLoop:
		for {
			select {
			case <-completionCh:
				completeCount++
			case <-time.After(100 * time.Millisecond):
				break countLoop
			}
		}

		// Verify only maxConcurrent requests completed
		require.True(t, completeCount <= int(maxConcurrent),
			"Too many requests (%d) completed, limit is %d",
			completeCount, maxConcurrent)
	})

	//nolint:paralleltest
	t.Run("MoveCall", func(t *testing.T) {
		// Prepare arguments for a move call
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "counter",
			Function:        "increment", // Assuming this function exists in the contract
			Arguments:       []any{counterObjectId},
			TypeArguments:   []any{"objectId"},
			Gas:             1000000000,
			GasBudget:       1000000000,
		}

		// Call MoveCall to prepare the transaction
		txnMetadata, err := relayerClient.MoveCall(context.Background(), moveCallReq)
		require.NoError(t, err)
		require.NotEmpty(t, txnMetadata.TxBytes, "Expected non-empty transaction bytes")

		// Verify we can execute the transaction
		resp, err := relayerClient.SignAndSendTransaction(
			context.Background(),
			txnMetadata.TxBytes,
			publicKeyBytes,
			"WaitForLocalExecution",
		)
		require.NoError(t, err)
		require.Equal(t, "success", resp.Status.Status, "Expected move call to succeed")
	})

	//nolint:paralleltest
	t.Run("MoveCall_IncrementByValue", func(t *testing.T) {
		// Prepare arguments for a move call
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "counter",
			Function:        "increment_by",
			Arguments:       []any{counterObjectId, "10"},
			TypeArguments:   []any{},
			Gas:             1000000000,
			GasBudget:       1000000000,
		}

		// Call MoveCall to prepare the transaction
		txnMetadata, err := relayerClient.MoveCall(context.Background(), moveCallReq)
		require.NoError(t, err)
		require.NotEmpty(t, txnMetadata.TxBytes, "Expected non-empty transaction bytes")

		// Verify we can execute the transaction
		resp, err := relayerClient.SignAndSendTransaction(
			context.Background(),
			txnMetadata.TxBytes,
			publicKeyBytes,
			"WaitForLocalExecution",
		)
		require.NoError(t, err)
		require.Equal(t, "success", resp.Status.Status, "Expected move call to succeed")
	})

	//nolint:paralleltest
	t.Run("QueryEvents", func(t *testing.T) {
		// Increment the counter 3 times to create multiple events
		IncrementCounterWithMoveCall(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		IncrementCounterWithMoveCall(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		IncrementCounterWithMoveCall(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)

		// Create event filter for the counter module
		filter := client.EventFilterByMoveEventModule{
			Package: packageId,
			Module:  "counter",
			Event:   "CounterIncremented",
		}

		limit := uint(1)
		descending := true

		// Query events
		events, err := relayerClient.QueryEvents(context.Background(), filter, &limit, nil, &client.QuerySortOptions{
			Descending: descending,
		})
		require.NoError(t, err)
		require.NotNil(t, events)
		require.Equal(t, 1, len(events.Data))

		// Query events again with the cursor of the previous query
		cursor := client.EventId{
			TxDigest: events.Data[0].Id.TxDigest,
			EventSeq: events.Data[0].Id.EventSeq,
		}
		eventsWithCursor, errWithCursor := relayerClient.QueryEvents(context.Background(), filter, &limit, &cursor, &client.QuerySortOptions{
			Descending: descending,
		})
		require.NoError(t, errWithCursor)
		require.NotNil(t, eventsWithCursor)
		require.True(t, len(eventsWithCursor.Data) > 0)
	})

	//nolint:paralleltest
	t.Run("QueryEvents_(high_limit)", func(t *testing.T) {
		IncrementCounterWithMoveCall(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)

		// Create event filter for the counter module
		filter := client.EventFilterByMoveEventModule{
			Package: packageId,
			Module:  "counter",
			Event:   "CounterIncremented",
		}

		limit := uint(50)
		descending := true

		// Query events
		events, err := relayerClient.QueryEvents(context.Background(), filter, &limit, nil, &client.QuerySortOptions{
			Descending: descending,
		})
		require.NoError(t, err)
		require.NotNil(t, events)
		require.True(t, len(events.Data) > 0)
	})

	//nolint:paralleltest
	t.Run("GetCoinsByAddress", func(t *testing.T) {
		// Get coins owned by the account
		coins, err := relayerClient.GetCoinsByAddress(context.Background(), accountAddress)
		require.NoError(t, err)
		require.NotNil(t, coins)

		// Account should have at least one coin after faucet funding
		require.True(t, len(coins) > 0, "Expected at least one coin in account")

		// Verify coin data structure
		for _, coin := range coins {
			require.NotEmpty(t, coin.CoinObjectId)
			require.NotEmpty(t, coin.CoinType)
			require.NotEmpty(t, coin.Balance)
		}
	})

	//nolint:paralleltest
	t.Run("ReadObjectId", func(t *testing.T) {
		// Read the counter object
		objectData, err := relayerClient.ReadObjectId(context.Background(), counterObjectId)
		require.NoError(t, err)
		require.NotNil(t, objectData)
	})

	//nolint:paralleltest
	t.Run("ReadOwnedObjects", func(t *testing.T) {
		// Read owned objects for account
		objects, err := relayerClient.ReadOwnedObjects(
			context.Background(),
			accountAddress,
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, objects)
		require.True(t, len(objects) > 0)
	})

	t.Run("ReadFilterOwnedObjectIds", func(t *testing.T) {
		objects, err := relayerClient.ReadFilterOwnedObjectIds(
			context.Background(),
			accountAddress,
			fmt.Sprintf("%s::counter::AdminCap", packageId),
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, objects)
		require.Equal(t, 1, len(objects))
		require.Equal(t, fmt.Sprintf("%s::counter::AdminCap", packageId), objects[0].Type)
	})

	//nolint:paralleltest
	t.Run("GetTransactionStatus", func(t *testing.T) {
		txDigest := IncrementCounterWithMoveCall(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)

		// Now check its status
		txStatus, err := relayerClient.GetTransactionStatus(context.Background(), txDigest)
		require.NoError(t, err)
		require.Equal(t, "success", txStatus.Status, "Expected transaction status to be 'success', got: %s with error: %s",
			txStatus.Status, txStatus.Error)
	})

	t.Run("ReadFunction_JSONResponseParsing", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_result_struct",
			[]any{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_NestedStruct", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_nested_result_struct",
			[]any{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_MultiNestedStruct", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_multi_nested_result_struct",
			[]any{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_Tuple", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_tuple_struct",
			[]any{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_OCRConfig", func(t *testing.T) {
		values, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_ocr_config",
			[]any{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(values)
	})

	t.Run("ReadFunction_VectorOfU8", func(t *testing.T) {
		values, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_vector_of_u8",
			[]any{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(values)
	})

	t.Run("ReadFunction_VectorOfAddresses", func(t *testing.T) {
		values, err := relayerClient.ReadFunction(
			context.Background(),
			accountAddress,
			packageId,
			"counter",
			"get_vector_of_addresses",
			[]any{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(values)
	})

	t.Run("QueryTransactions", func(t *testing.T) {
		values, err := relayerClient.QueryTransactions(
			context.Background(),
			accountAddress,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, values)
		require.True(t, len(values.Data) > 0)

		utils.PrettyPrint(values)
	})

	t.Run("QueryFailedTransactions", func(t *testing.T) {
		CreateFailedTransaction(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)

		values, err := relayerClient.QueryTransactions(
			context.Background(),
			accountAddress,
			nil,
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, values)
		require.True(t, len(values.Data) > 0)

		failuresCount := 0
		for _, tx := range values.Data {
			if tx.Effects.Status.Status == "failure" {
				failuresCount++
			}
		}

		// expect to find the failed transaction in the list
		require.True(t, failuresCount > 0, "Expected at least one failure")

		// create another failed transaction and use the cursor to ignore the previously fetched ones
		CreateFailedTransaction(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		cursor := &values.NextCursor
		values, err = relayerClient.QueryTransactions(
			context.Background(),
			accountAddress,
			cursor,
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, values)
		require.Equal(t, 1, len(values.Data))
	})
}

func IncrementCounterWithMoveCall(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) string {
	t.Helper()
	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "counter",
		Function:        "increment", // Assuming this function exists in the contract
		Arguments:       []any{counterObjectId},
		GasBudget:       1000000000,
	}

	// Call MoveCall to prepare the transaction
	txnMetadata, err := relayerClient.MoveCall(context.Background(), moveCallReq)
	require.NoError(t, err)
	require.NotEmpty(t, txnMetadata.TxBytes, "Expected non-empty transaction bytes")

	// Verify we can execute the transaction
	resp, err := relayerClient.SignAndSendTransaction(
		context.Background(),
		txnMetadata.TxBytes,
		signerPublicKey,
		"WaitForLocalExecution",
	)
	require.NoError(t, err)
	require.Equal(t, "success", resp.Status.Status, "Expected move call to succeed")

	return resp.TxDigest
}

func CreateFailedTransaction(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) {
	t.Helper()
	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "counter",
		Function:        "increment_by",
		Arguments:       []any{counterObjectId, "1000"},
		GasBudget:       1000000000,
	}

	// Call MoveCall to prepare the transaction
	txnMetadata, err := relayerClient.MoveCall(context.Background(), moveCallReq)
	require.NoError(t, err)
	require.NotEmpty(t, txnMetadata.TxBytes, "Expected non-empty transaction bytes")

	// Verify we can execute the transaction
	resp, err := relayerClient.SignAndSendTransaction(
		context.Background(),
		txnMetadata.TxBytes,
		signerPublicKey,
		"WaitForLocalExecution",
	)
	require.NoError(t, err)
	require.Equal(t, "failure", resp.Status.Status, "Expected move call to fail")
}
