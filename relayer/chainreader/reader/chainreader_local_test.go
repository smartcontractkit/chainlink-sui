//go:build integration

package reader

import (
	"context"
	"crypto/ed25519"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
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
	ctx := context.Background()

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

	pointerTag := "_::counter::CounterPointer::counter_id"

	// Set up the ChainReader
	chainReaderConfig := config.ChainReaderConfig{
		IsLoopPlugin: false,
		EventsIndexer: config.EventsIndexerConfig{
			PollingInterval: 10 * time.Second,
			SyncTimeout:     10 * time.Second,
		},
		TransactionsIndexer: config.TransactionsIndexerConfig{
			PollingInterval: 10 * time.Second,
			SyncTimeout:     10 * time.Second,
		},
		Modules: map[string]*config.ChainReaderModule{
			"counter": {
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
				},
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "counter",
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

	chainReader, err := NewChainReader(ctx, log, relayerClient, chainReaderConfig, db)
	require.NoError(t, err)

	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	log.Debugw("ChainReader setup complete")

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

	t.Run("GetLatestValue_SimpleStruct", func(t *testing.T) {
		var retSimpleResult SimpleResult

		log.Debugw("Testing get_simple_result function for BCS struct decoding",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "counter", "get_simple_result"}, "-"),
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
			strings.Join([]string{packageId, "counter", "get_address_list"}, "-"),
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
			strings.Join([]string{packageId, "counter", "get_tuple_struct"}, "-"),
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
		t.Parallel()
		ctx := context.Background()

		go func() {
			_ = chainReader.Start(ctx)
		}()

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

	t.Run("GetLatestValue_PointerTag", func(t *testing.T) {
		expectedUint64 := uint64(0)
		var retUint64 uint64

		log.Debugw("Testing get_simple_result function for BCS struct decoding",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "counter", "get_count_using_pointer"}, "-"),
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
