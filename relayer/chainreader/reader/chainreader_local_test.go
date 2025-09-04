//go:build integration

package reader

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-aptos/relayer/chainreader"
	aptosCRConfig "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type AddressList struct {
	Addresses [][]byte `json:"addresses"`
	Count     uint64   `json:"count"`
}

// Go struct that matches the Move SimpleResult struct
type SimpleResult struct {
	Value uint64 `json:"value"`
}

func TestChainReaderLocal(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

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

	runChainReaderCounterTest(t, log, testutils.LocalUrl)
}

func runChainReaderCounterTest(t *testing.T, log logger.Logger, rpcUrl string) {
	t.Helper()
	ctx := context.Background()

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	relayerClient, clientErr := client.NewPTBClient(log, rpcUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, clientErr)

	faucetFundErr := testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, faucetFundErr)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	gasBudget := int(2000000000)
	packageId, tx, err := testutils.PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotNil(t, packageId)
	require.NotNil(t, tx)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(tx.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	pointerTag := "_::counter::CounterPointer::counter_id"

	pollingInterval := 10 * time.Second
	syncTimeout := 10 * time.Second

	// Set up the ChainReader
	chainReaderConfig := config.ChainReaderConfig{
		IsLoopPlugin: false,
		Modules: map[string]*config.ChainReaderModule{
			"Counter": {
				Name: "counter",
				Functions: map[string]*config.ChainReaderFunction{
					"get_count": {
						Name:          "get_count",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:         "object_id",
								Name:         "counter_id",
								DefaultValue: counterObjectId,
								Required:     true,
							},
						},
					},
					"get_address_list": {
						Name:          "get_address_list",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
					},
					"get_simple_result": {
						Name:          "get_simple_result",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
					},
					"get_tuple_struct": {
						Name:                "get_tuple_struct",
						SignerAddress:       accountAddress,
						Params:              []codec.SuiFunctionParam{}, // No parameters needed
						ResultTupleToStruct: []string{"value", "address", "bool", "struct_tag"},
					},
					"get_count_using_pointer": {
						Name:          "get_count_using_pointer",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:       "object_id",
								Name:       "counter_id",
								PointerTag: &pointerTag,
								Required:   true,
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{
					"counter_incremented": {
						Name:      "counter_incremented",
						EventType: "CounterIncremented",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "counter",
							Event:   "CounterIncremented",
						},
					},
					"counter_decremented": {
						Name:      "counter_decremented",
						EventType: "CounterDecremented",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "counter",
							Event:   "CounterDecremented",
						},
					},
				},
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "Counter",
		Address: packageId, // Package ID of the deployed counter contract
	}

	datastoreUrl := os.Getenv("TEST_DB_URL")
	if datastoreUrl == "" {
		t.Skip("Skipping persistent tests as TEST_DB_URL is not set in CI")
	}
	db := sqltest.NewDB(t, datastoreUrl)

	// attempt to connect
	_, err = db.Connx(ctx)
	require.NoError(t, err)

	// Create the indexers
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		log,
		relayerClient,
		pollingInterval,
		syncTimeout,
		// start without any configs, they will be set when ChainReader is initialized and gets a reference
		// to the transaction indexer to avoid having to reading ChainReader configs here as well
		map[string]*config.ChainReaderEvent{},
	)
	evIndexer := indexer.NewEventIndexer(
		db,
		log,
		relayerClient,
		// start without any selectors, they will be added during .Bind() calls on ChainReader
		[]*client.EventSelector{},
		pollingInterval,
		syncTimeout,
	)
	indexerInstance := indexer.NewIndexer(
		log,
		evIndexer,
		txnIndexer,
	)

	chainReader, err := NewChainReader(ctx, log, relayerClient, chainReaderConfig, db, indexerInstance)
	require.NoError(t, err)

	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	log.Debugw("ChainReader setup complete")

	go func() {
		err = chainReader.Start(ctx)
		require.NoError(t, err)
		log.Debugw("ChainReader started")
	}()
	go func() {
		err = indexerInstance.Start(ctx)
		require.NoError(t, err)
		log.Debugw("Indexers started")
	}()

	t.Run("GetLatestValue_FunctionRead", func(t *testing.T) {
		expectedUint64 := uint64(0)
		var retUint64 uint64

		log.Debugw("Testing get_count",
			"counterObjectId", counterObjectId,
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_count"}, "-"),
			primitives.Finalized,
			map[string]any{
				"counter_id": counterObjectId,
			},
			&retUint64,
		)
		require.NoError(t, err)
		require.Equal(t, expectedUint64, retUint64)
	})

	t.Run("GetLatestValue_SimpleStruct", func(t *testing.T) {
		var retSimpleResult SimpleResult

		log.Debugw("Testing get_simple_result function for BCS struct decoding",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_simple_result"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed
			&retSimpleResult,
		)
		require.NoError(t, err)

		// Verify the returned struct
		require.NotNil(t, retSimpleResult)
		require.Equal(t, uint64(42), retSimpleResult.Value, "Expected value to be 42")

		log.Debugw("SimpleResult test completed successfully",
			"value", retSimpleResult.Value)
	})

	t.Run("GetLatestValue_AddressList", func(t *testing.T) {
		var retAddressList AddressList

		log.Debugw("Testing get_address_list function",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_address_list"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed
			&retAddressList,
		)
		require.NoError(t, err)

		// Verify the returned struct
		require.NotNil(t, retAddressList)

		log.Debugw("retAddressList", "retAddressList", retAddressList)

		require.Equal(t, uint64(4), retAddressList.Count, "Expected 4 addresses")
		require.Len(t, retAddressList.Addresses, 4, "Expected 4 addresses in the list")

		// Verify the expected addresses match what we defined in the Move function
		expectedAddresses := [][32]byte{
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3},
			{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4},
		}

		for i, addr := range retAddressList.Addresses {
			log.Debugw("Address comparison", "index", i, "expected", expectedAddresses[i], "actual", addr)
		}

		log.Debugw("AddressList test completed successfully",
			"count", retAddressList.Count,
			"addresses", retAddressList.Addresses)
	})

	t.Run("GetLatestValue_TupleToStruct", func(t *testing.T) {
		var retTupleStruct map[string]any

		log.Debugw("Testing get_tuple_struct function for BCS struct decoding",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_tuple_struct"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed
			&retTupleStruct,
		)
		require.NoError(t, err)

		// Verify the returned struct
		require.NotNil(t, retTupleStruct)
		require.Equal(t, uint64(42), retTupleStruct["value"], "Expected value to be 42")
		require.Len(t, retTupleStruct["address"].([]byte), 32, "Expected address to be 0x1")
		require.Equal(t, true, retTupleStruct["bool"], "Expected bool to be true")

		log.Debugw("TupleStruct test completed successfully",
			"value", retTupleStruct["value"],
			"address", retTupleStruct["address"],
			"bool", retTupleStruct["bool"],
			"struct_tag", retTupleStruct["struct_tag"])
	})

	t.Run("QueryKey_Events", func(t *testing.T) {
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
			GasBudget:       2000000,
		}

		log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

		txMetadata, testErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, testErr)

		txnResult, testErr := relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, testErr)

		log.Debugw("Transaction result", "result", txnResult)

		// Query for counter increment events
		type CounterEvent struct {
			CounterID string `json:"counterId"`
			NewValue  uint64 `json:"newValue"`
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

	t.Run("QueryKey_WithFilter", func(t *testing.T) {
		// Decrement the counter to emit an event (different from what has been previously emitted)
		log.Debugw("Decrementing counter to emit event", "counterObjectId", counterObjectId)
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "counter",
			Function:        "decrement",
			TypeArguments:   []any{},
			Arguments:       []any{counterObjectId},
			GasBudget:       2000000,
		}

		txMetadata, testErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, testErr)

		_, testErr = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, testErr)

		// Query for counter increment events
		type CounterDecrementEvent struct {
			EventType string `json:"eventType"`
			CounterID string `json:"counterId"`
			NewValue  uint64 `json:"newValue"`
		}

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "counter_decremented",
		}

		// Setup limit and sort
		limitAndSort := query.LimitAndSort{
			Limit: query.Limit{
				Count:  50,
				Cursor: "",
			},
		}

		sequences := []types.Sequence{}
		require.Eventually(t, func() bool {
			// Query for events
			var counterEvent CounterDecrementEvent
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
		require.NotEmpty(t, sequences, "Expected at least one event")
	})

	t.Run("QueryKey_WithMetadata", func(t *testing.T) {
		type CounterDecrementEvent struct {
			EventType string `json:"eventType"`
			CounterID string `json:"counterId"`
			NewValue  uint64 `json:"newValue"`
		}

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "counter_decremented",
		}

		// Setup limit and sort
		limitAndSort := query.LimitAndSort{
			Limit: query.Limit{
				Count:  50,
				Cursor: "",
			},
		}

		sequences := []aptosCRConfig.SequenceWithMetadata{}
		require.Eventually(t, func() bool {
			// Query for events
			var counterEvent CounterDecrementEvent
			sequences, err = chainReader.(chainreader.ExtendedContractReader).QueryKeyWithMetadata(
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
		require.NotEmpty(t, sequences, "Expected at least one event")
	})

	t.Run("GetLatestValue_PointerTag", func(t *testing.T) {
		expectedUint64 := uint64(0)
		var retUint64 uint64

		log.Debugw("Testing get_simple_result function for BCS struct decoding",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_count_using_pointer"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed, the counter_id object should be populated from the pointer tag
			&retUint64,
		)
		require.NoError(t, err)

		// Verify the returned struct
		require.NotNil(t, retUint64)
		require.Equal(t, expectedUint64, retUint64, "Expected value to be 0")
	})
}
