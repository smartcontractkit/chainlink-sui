//go:build integration

package client_test

import (
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

const (
	grpcDefaultMaxConcurrent = int64(3)
	grpcTestTimeout          = 120 * time.Second
)

//nolint:paralleltest
func TestGrpcClient(t *testing.T) {
	log := logger.Test(t)

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	testCfg := client.PTBClientConfig{
		GrpcTarget:            "127.0.0.1:9000",
		GrpcToken:             "test",
		TransactionTimeout:    grpcTestTimeout,
		MaxConcurrentRequests: grpcDefaultMaxConcurrent,
		KeystoreService:       keystoreInstance,
	}

	var suiNodeCmd *exec.Cmd
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	suiNodeCmd = cmd

	testutils.CleanupTestContracts()

	t.Cleanup(func() {
		testutils.CleanupTestContracts()

		if suiNodeCmd != nil && suiNodeCmd.Process != nil {
			if perr := suiNodeCmd.Process.Kill(); perr != nil {
				t.Logf("Failed to kill local Sui node: %v", perr)
			}
		}
	})

	relayerClient, err := client.NewPTBClient(log, testCfg)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, relayerClient.Close())
	})

	t.Run("GrpcConnection", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		require.NoError(t, relayerClient.VerifyGrpcServices(ctx))

		chainID, err := relayerClient.HealthCheckGrpc(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, chainID, "expected chain ID from gRPC GetServiceInfo")
		log.Infow("gRPC connection verified", "grpcTarget", testCfg.GrpcTarget, "chainID", chainID)
	})

	err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	chainID, err := relayerClient.HealthCheckGrpc(context.Background())
	require.NoError(t, err)
	testutils.PatchEnvironmentTOML("contracts/test", "local", chainID)
	testutils.PatchEnvironmentTOML("contracts/test_secondary", "local", chainID)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	gasBudget := int(2000000000)
	packageId, tx, err := testutils.PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotNil(t, tx)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(tx.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// Test GetLatestValue for different data types
	//nolint:paralleltest
	t.Run("FunctionRead", func(t *testing.T) {
		args := []any{counterObjectId}
		argTypes := []string{"objectId"}

		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_count",
			args,
			argTypes,
			[]string{},
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

				err := relayerClient.WithRateLimit(ctx, "TestMethod", func(ctx context.Context) error {
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
		require.True(t, completeCount <= int(grpcDefaultMaxConcurrent),
			"Too many requests (%d) completed, limit is %d",
			completeCount, grpcDefaultMaxConcurrent)
	})

	//nolint:paralleltest
	t.Run("MoveCall_IncrementByValue", func(t *testing.T) {
		txnMetadata, err := relayerClient.MoveCall(context.Background(), client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "counter",
			Function:        "increment_by",
			Arguments:       []any{counterObjectId, uint64(100)},
			TypeArguments:   []any{},
			GasBudget:       1000000000,
		})
		require.NoError(t, err)

		// Verify we can execute the transaction
		resp, err := relayerClient.SignAndSendTransaction(context.Background(), txnMetadata.TxBytes, publicKeyBytes)
		require.NoError(t, err)
		log.Debugw("transaction response", "transaction", resp.Transaction.GetEffects().GetStatus())

		require.Equal(t, true, resp.Transaction.GetEffects().GetStatus().GetSuccess(), "Expected move call to succeed")
	})

	//nolint:paralleltest
	t.Run("GetCoinsByAddress", func(t *testing.T) {
		// Get coins owned by the account
		coins, err := relayerClient.GetCoinsByAddress(context.Background(), accountAddress)
		require.NoError(t, err)
		require.NotNil(t, coins)

		for _, coin := range coins {
			log.Debugw("coin", "coin", coin)
		}

		// Account should have at least one coin after faucet funding
		require.True(t, len(coins) > 0, "Expected at least one coin in account")

		// Verify coin data structure
		for _, coin := range coins {
			require.NotEmpty(t, coin.GetObjectId())
		}

		utils.PrettyPrint(coins)
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

	t.Run("ReadFilterOwnedObjectIds_(no_cursor)", func(t *testing.T) {
		objects, err := relayerClient.ReadFilterOwnedObjectIds(
			context.Background(),
			accountAddress,
			fmt.Sprintf("%s::counter::AdminCap", packageId),
			nil,
		)
		require.NoError(t, err)
		require.NotNil(t, objects)
		require.NotZero(t, len(objects))
		require.Equal(t, fmt.Sprintf("%s::counter::AdminCap", packageId), objects[0].GetObjectType())
	})

	t.Run("ReadFilterOwnedObjectIds_(many_pages)", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			err := testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
			require.NoError(t, err)
		}

		CreateManyObjects(t, relayerClient, packageId, accountAddress, publicKeyBytes)

		assert.Eventually(t, func() bool {
			objects, err := relayerClient.ReadFilterOwnedObjectIds(
				context.Background(),
				accountAddress,
				fmt.Sprintf("%s::counter::SomeObject", packageId),
				nil,
			)
			require.NoError(t, err)
			log.Debugw("objects", "objects", objects)
			return len(objects) > 99
		}, 10*time.Second, 2*time.Second)
	})

	//nolint:paralleltest
	t.Run("GetTransactionStatus", func(t *testing.T) {
		tx := IncrementCounterWithMoveCall(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		txDigest := tx.GetDigest()

		// Now check its status
		assert.Eventually(t, func() bool {
			txStatus, err := relayerClient.GetTransactionStatus(context.Background(), txDigest)
			return txStatus.Status == "success" && err == nil
		}, 10*time.Second, 2*time.Second, "Expected transaction status to be 'success' but condition not met")
	})

	t.Run("ReadFunction_JSONResponseParsing", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_result_struct",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_StringAsStruct", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"type_and_version",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		require.Equal(t, "Counter 1.6.0", response[0], "Expected type and version to match 'Counter 1.6.0'")
	})

	t.Run("ReadFunction_NestedStruct", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_nested_result_struct",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_MultiNestedStruct", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_multi_nested_result_struct",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_Tuple", func(t *testing.T) {
		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_tuple_struct",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(response)
	})

	t.Run("ReadFunction_OCRConfig", func(t *testing.T) {
		values, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_ocr_config",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(values)
	})

	t.Run("ReadFunction_VectorOfU8", func(t *testing.T) {
		values, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_vector_of_u8",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(values)
	})

	t.Run("ReadFunction_VectorOfAddresses", func(t *testing.T) {
		values, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_vector_of_addresses",
			[]any{},
			[]string{},
			[]string{},
		)
		require.NoError(t, err)
		utils.PrettyPrint(values)
	})

	t.Run("GetLatestPackageId", func(t *testing.T) {
		latestPackageId, err := relayerClient.GetLatestPackageId(
			context.Background(),
			packageId,
			"counter",
		)

		// The latest package ID should be the same as the provided package ID
		require.NoError(t, err)
		require.Equal(t, latestPackageId, packageId)
	})

	t.Run("QueryCoinsByAddress", func(t *testing.T) {
		coins, err := relayerClient.QueryCoinsByAddress(context.Background(), accountAddress, "0x2::coin::Coin<0x2::sui::SUI>")
		require.NoError(t, err)
		require.NotNil(t, coins)

		require.True(t, len(coins) > 0)
		for _, coin := range coins {
			require.Equal(t, accountAddress, coin.GetOwner().GetAddress())
			expectedCoinType := "0x0000000000000000000000000000000000000000000000000000000000000002::coin::Coin<0x0000000000000000000000000000000000000000000000000000000000000002::sui::SUI>"
			require.Equal(t, expectedCoinType, coin.GetObjectType())
		}

		utils.PrettyPrint(coins)
	})

	t.Run("GetLatestCheckpoint", func(t *testing.T) {
		checkpoint, err := relayerClient.GetLatestCheckpoint(context.Background())
		require.NoError(t, err)
		require.NotNil(t, checkpoint)
		require.NotZero(t, checkpoint.GetSequenceNumber())
	})

	t.Run("GetCheckpointData", func(t *testing.T) {
		// create a random tranaction to be included in the next checkpoint
		tx := IncrementCounterWithMoveCall(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		txDigest := tx.GetDigest()
		require.NotEmpty(t, txDigest)

		var checkpoint uint64

		// assert the transaction is included in the next checkpoint
		assert.Eventually(t, func() bool {
			transaction, err := relayerClient.GetTransactionStatus(context.Background(), txDigest)
			checkpoint = transaction.Checkpoint
			return transaction.Status == "success" && err == nil
		}, 10*time.Second, 1*time.Second)

		checkpointData, err := relayerClient.GetCheckpointData(context.Background(), checkpoint)
		require.NoError(t, err)
		require.NotNil(t, checkpoint)
		require.Equal(t, checkpoint, checkpointData.Checkpoint.GetSequenceNumber())

		// assert the transaction is included in the checkpoint data
		var found bool
		for _, tx := range checkpointData.Transactions {
			if tx.GetDigest() == txDigest {
				require.Equal(t, tx.GetCheckpoint(), checkpoint)
				require.Contains(t, checkpointData.Transactions, tx)
				found = true
				break
			}
		}

		require.True(t, found, "Transaction not found in checkpoint data")
	})
}

func IncrementCounterWithMoveCall(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) *suirpcv2.ExecutedTransaction {
	t.Helper()
	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "counter",
		Function:        "increment", // Assuming this function exists in the contract
		Arguments:       []any{counterObjectId},
		TypeArguments:   []any{},
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
	)
	require.NoError(t, err)
	require.Equal(t, true, resp.Transaction.GetEffects().GetStatus().GetSuccess(), "Expected move call to succeed")

	return resp.GetTransaction()
}

func CreateManyObjects(t *testing.T, relayerClient *client.PTBClient, packageId string, accountAddress string, signerPublicKey []byte) {
	t.Helper()
	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "counter",
		Function:        "create_many_objects",
		Arguments:       []any{uint64(100)},
		TypeArguments:   []any{},
		GasBudget:       2000000000,
	}

	txnMetadata, err := relayerClient.MoveCall(context.Background(), moveCallReq)
	require.NoError(t, err)
	require.NotEmpty(t, txnMetadata.TxBytes, "Expected non-empty transaction bytes")

	// Verify we can execute the transaction
	resp, err := relayerClient.SignAndSendTransaction(
		context.Background(),
		txnMetadata.TxBytes,
		signerPublicKey,
	)
	require.NoError(t, err)
	require.Equal(t, true, resp.Transaction.GetEffects().GetStatus().GetSuccess(), "Expected move call to succeed")
}

func CreateFailedTransaction(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) {
	t.Helper()
	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "counter",
		Function:        "increment_by",
		Arguments:       []any{counterObjectId, uint64(1000)},
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
	)
	require.NoError(t, err)
	require.Equal(t, false, resp.Transaction.GetEffects().GetStatus().GetSuccess(), "Expected move call to fail")
}
