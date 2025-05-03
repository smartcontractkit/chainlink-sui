//go:build integration

package chainreader

import (
	"context"
	"crypto/ed25519"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestChainReaderLocal(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	var err error
	accountAddress := testutils.GetAccountAndKeyFromSui(t, log)
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	// Ensure the process is killed when the test completes.
	t.Cleanup(func() {
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	log.Debugw("Started Sui node")

	err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	runChainReaderCounterTest(t, log, testutils.LocalUrl)
}

func runChainReaderCounterTest(t *testing.T, log logger.Logger, rpcUrl string) {
	t.Helper()

	accountAddress := testutils.GetAccountAndKeyFromSui(t, log)
	keystoreInstance, keystoreErr := keystore.NewSuiKeystore(log, "")
	require.NoError(t, keystoreErr)

	privateKey, err := keystoreInstance.GetPrivateKeyByAddress(accountAddress)
	require.NoError(t, err)
	publicKey := privateKey.Public().(ed25519.PublicKey)
	publicKeyBytes := []byte(publicKey)

	relayerClient, clientErr := client.NewPTBClient(log, rpcUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, clientErr)

	faucetFundErr := testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, faucetFundErr)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, publishOutput, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// Set up the ChainReader
	chainReaderConfig := ChainReaderConfig{
		Modules: map[string]*ChainReaderModule{
			"counter": {
				Name: "counter",
				Functions: map[string]*ChainReaderFunction{
					"get_count": {
						Name:          "get_count",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:         "address",
								Name:         "counter_id",
								DefaultValue: counterObjectId,
								Required:     true,
							},
						},
					},
				},
				Events: map[string]*ChainReaderEvent{
					"counter_incremented": {
						Name:      "counter_incremented",
						EventType: "CounterIncremented",
					},
				},
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "counter",
		Address: packageId, // Package ID of the deployed counter contract
	}

	chainReader := NewChainReader(log, *relayerClient, chainReaderConfig)
	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	log.Debugw("ChainReader setup complete")

	// Test GetLatestValue for different data types

	t.Run("GetLatestValue_Uint64", func(t *testing.T) {
		expectedUint64 := uint64(0)
		var retUint64 uint64
		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, counterBinding.Name, counterObjectId}, "-"),
			primitives.Finalized,
			struct {
				Value uint64
			}{Value: expectedUint64},
			&retUint64,
		)
		require.NoError(t, err)
		require.Equal(t, expectedUint64, retUint64)
	})

	t.Run("GetLatestValue_FunctionRead", func(t *testing.T) {
		expectedUint64 := uint64(0)
		var retUint64 uint64

		log.Debugw("Testing get_count",
			"counterObjectId", counterObjectId,
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "counter", "get_count"}, "-"),
			primitives.Finalized,
			map[string]any{
				"counter_id": counterObjectId,
			},
			&retUint64,
		)
		require.NoError(t, err)
		require.Equal(t, expectedUint64, retUint64)
	})

	t.Run("QueryKey_Events", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		// Increment the counter to emit an event
		log.Debugw("Incrementing counter to emit event", "counterObjectId", counterObjectId)

		// Use relayerClient to call increment instead of using CLI
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "counter",
			Function:        "increment",
			TypeArguments:   []any{},
			Arguments:       []any{counterObjectId},
			GasBudget:       "2000000",
		}

		log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

		txMetadata, err := relayerClient.MoveCall(context.Background(), moveCallReq)
		require.NoError(t, err)

		txnResult, err := relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, err)

		log.Debugw("Transaction result", "result", txnResult)

		// Query for counter increment events
		type CounterEvent struct {
			CounterID string `json:"counter_id"`
			NewValue  uint64 `json:"new_value"`
		}

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "counter_incremented",
		}

		// Setup limit and sort
		limitAndSort := query.LimitAndSort{
			Limit: query.Limit{
				Count:  50,
				Cursor: "",
			},
		}

		log.Debugw("Querying for counter events",
			"filter", filter.Key,
			"limit", limitAndSort.Limit.Count,
			"packageId", packageId,
			"contract", counterBinding.Name,
			"eventType", "CounterIncremented")

		sequences := []types.Sequence{}
		require.Eventually(t, func() bool {
			// Query for events
			var counterEvent CounterEvent
			sequences, err = chainReader.QueryKey(
				ctx,
				counterBinding,
				filter,
				limitAndSort,
				&counterEvent,
			)
			if err != nil {
				log.Errorw("Failed to query events", "error", err)
				require.NoError(t, err)
			}

			return len(sequences) > 0
		}, 60*time.Second, 1*time.Second, "Event should eventually be indexed and found")

		log.Debugw("Query results", "sequences", sequences)

		// Verify we got at least one event
		require.NotEmpty(t, sequences, "Expected at least one event")

		// Verify the event data
		event := sequences[0].Data.(*CounterEvent)
		require.NotNil(t, event)
		log.Debugw("Event data", "counterId", event.CounterID, "newValue", event.NewValue)
		require.Equal(t, uint64(1), event.NewValue, "Expected counter value to be 1")
	})
}
