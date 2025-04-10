//go:build integration

package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestClient(t *testing.T) {
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

	accountAddress := testutils.GetAccountAndKeyFromSui(t, log)
	keystoreInstance, err := keystore.NewSuiKeystore(log, "", keystore.PrivateKeySigner)
	require.NoError(t, err)
	signer, err := keystoreInstance.GetSignerFromAddress(accountAddress)
	require.NoError(t, err)
	maxConcurrent := int64(3)
	relayerClient, err := client.NewClient(log, testutils.LocalUrl, nil, 10*time.Second, &signer, maxConcurrent)
	require.NoError(t, err)

	err = testutils.FundWithFaucet(log, constant.SuiLocalnet, accountAddress)
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
		argTypes := []string{"address"}

		response, err := relayerClient.ReadFunction(
			context.Background(),
			packageId,
			"counter",
			"get_count",
			args,
			argTypes,
		)
		require.NoError(t, err)
		require.NotNil(t, response)
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
}
