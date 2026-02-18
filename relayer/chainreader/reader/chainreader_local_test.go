//go:build integration

package reader

import (
	"context"
	"encoding/json"
	"math/big"
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

	testutils.CleanupTestContracts()

	t.Cleanup(func() {
		testutils.CleanupTestContracts()
	})

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	relayerClient, clientErr := client.NewPTBClient(log, rpcUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, clientErr)

	faucetFundErr := testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, faucetFundErr)

	// Publish test_secondary first (before counter, since counter depends on it)
	gasBudget := int(2000000000)
	contractPath := testutils.BuildSetup(t, "contracts/test_secondary")
	secondaryPackageId, tx, err := testutils.PublishContract(t, "test_secondary", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotNil(t, secondaryPackageId)
	require.NotNil(t, tx)

	log.Debugw("Published Secondary Contract", "packageId", secondaryPackageId)

	// Now publish counter with test_secondary package ID patched in Move.toml
	contractPath = testutils.BuildSetup(t, "contracts/test")
	packageId, tx, err := testutils.PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotNil(t, packageId)
	require.NotNil(t, tx)

	log.Debugw("Published Contract", "packageId", packageId)
	counterObjectId, err := testutils.QueryCreatedObjectID(tx.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	type RampMessageHeader struct {
		MessageId           string
		SourceChainSelector uint64
		DestChainSelector   uint64
		SequenceNumber      uint64
		Nonce               uint64
	}

	type Sui2AnyTokenTransfer struct {
		SourcePoolAddress string
		DestTokenAddress  string
		ExtraData         string
		Amount            uint64
		DestExecData      string
	}

	type Sui2AnyRampMessage struct {
		Header         RampMessageHeader
		Sender         string
		Data           string
		Receiver       string
		ExtraArgs      string
		FeeToken       string
		FeeTokenAmount uint64
		FeeValueJuels  *big.Int
		TokenAmounts   []Sui2AnyTokenTransfer
	}

	type CCIPMessageSent struct {
		DestChainSelector uint64
		SequenceNumber    uint64
		Message           Sui2AnyRampMessage
	}

	type OfframpExecutionStateChanged struct {
		SourceChainSelector uint64 `json:"sourceChainSelector"`
		SequenceNumber      uint64 `json:"sequenceNumber"`
		MessageId           string `json:"messageId"`
		MessageHash         string `json:"messageHash"`
		State               uint8  `json:"state"`
	}

	// Define pointer tag for counter object derivation
	pointerTag := &codec.PointerTag{
		Module:        "counter",
		PointerName:   "CounterPointer",
		FieldName:     "counter_object_id",
		DerivationKey: "Counter",
	}

	pointerTagSecondary := &codec.PointerTag{
		Module:        "state_object",
		PointerName:   "CCIPObjectRefPointer",
		FieldName:     "ccip_object_id",
		DerivationKey: "CCIPObjectRef",
		PackageID:     secondaryPackageId,
	}

	type CounterIncrementedEvent struct {
		CounterID string `json:"counterId"`
		NewValue  uint64 `json:"newValue"`
	}

	type CounterDecrementedEvent struct {
		EventType string `json:"eventType"`
		CounterID string `json:"counterId"`
		NewValue  uint64 `json:"newValue"`
	}

	type NestedCounterBytesEvent struct {
		Value uint64 `json:"value"`
		Bytes string `json:"bytes"`
	}

	type CounterBytesEvent struct {
		Bytes  string                  `json:"bytes"`
		Nested NestedCounterBytesEvent `json:"nested"`
		Values []uint64                `json:"values"`
	}

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
					"get_address_list_renamed": {
						Name:          "get_address_list",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
						ResultFieldRenames: map[string]aptosCRConfig.RenamedField{
							"addresses": {NewName: "wallets"},
							"count":     {NewName: "size"},
						},
					},
					"get_simple_result": {
						Name:          "get_simple_result",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
					},
					"get_simple_result_renamed": {
						Name:          "get_simple_result",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
						ResultFieldRenames: map[string]aptosCRConfig.RenamedField{
							"value": {NewName: "renamedValue"},
						},
					},
					"get_tuple_struct": {
						Name:                "get_tuple_struct",
						SignerAddress:       accountAddress,
						Params:              []codec.SuiFunctionParam{}, // No parameters needed
						ResultTupleToStruct: []string{"value", "address", "bool", "struct_tag"},
					},
					"get_tuple_struct_renamed": {
						Name:                "get_tuple_struct",
						SignerAddress:       accountAddress,
						Params:              []codec.SuiFunctionParam{}, // No parameters needed
						ResultTupleToStruct: []string{"value", "address", "bool", "struct_tag"},
						ResultFieldRenames: map[string]aptosCRConfig.RenamedField{
							"value":      {NewName: "answer"},
							"struct_tag": {NewName: "tag"},
						},
					},
					"get_count_using_pointer": {
						Name:          "get_count_using_pointer",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:       "object_id",
								Name:       "counter_id",
								PointerTag: pointerTag,
								Required:   true,
							},
						},
					},
					"get_value_with_pointer_dependency": {
						Name:          "get_value_with_pointer_dependency",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:       "object_id",
								Name:       "counter_id",
								PointerTag: pointerTag,
								Required:   true,
							},
							{
								Type:       "object_id",
								Name:       "ccip_object_ref_id",
								PointerTag: pointerTagSecondary,
								Required:   true,
							},
						},
					},
					"static_response": {
						Name:          "static_response",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:       "object_id",
								Name:       "counter_id",
								PointerTag: pointerTag,
								Required:   true,
							},
							{
								Type:       "object_id",
								Name:       "ccip_object_ref_id",
								PointerTag: pointerTagSecondary,
								Required:   true,
							},
						},
						StaticResponse:      []any{1, 2, 3},
						ResultTupleToStruct: []string{"a", "b", "c"},
					},
					"response_from_inputs": {
						Name:               "response_from_inputs",
						SignerAddress:      accountAddress,
						Params:             []codec.SuiFunctionParam{},
						ResponseFromInputs: []string{"package_id"},
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
						EventFieldRenames: map[string]aptosCRConfig.RenamedField{
							"counter_id": {NewName: "counterId"},
							"new_value":  {NewName: "newValue"},
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
					"counter_bytes": {
						Name:      "counter_bytes",
						EventType: "CounterBytes",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "counter",
							Event:   "CounterBytes",
						},
						ExpectedEventType: &CounterBytesEvent{},
					},
				},
			},
			"OffRamp": {
				Name: "offramp",
				Functions: map[string]*config.ChainReaderFunction{
					"get_all_source_chain_configs": {
						Name:          "get_all_source_chain_configs",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
					},
					"emit_execution_state_changed_event": {
						Name:          "emit_execution_state_changed_event",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type: "u64",
								Name: "source_chain_selector",
							},
							{
								Type: "u64",
								Name: "sequence_number",
							},
							{
								Type: "vector<u8>",
								Name: "message_id",
							},
							{
								Type: "vector<u8>",
								Name: "message_hash",
							},
							{
								Type: "u8",
								Name: "state",
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{
					"execution_state_changed": {
						Name:      "execution_state_changed",
						EventType: "ExecutionStateChanged",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "offramp",
							Event:   "ExecutionStateChanged",
						},
						EventFieldRenames: map[string]aptosCRConfig.RenamedField{},
						ExpectedEventType: &OfframpExecutionStateChanged{},
					},
				},
			},
			"OnRamp": {
				Name: "onramp",
				Functions: map[string]*config.ChainReaderFunction{
					"emit_sample_ccip_message_sent_event": {
						Name:          "emit_sample_ccip_message_sent_event",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
					},
				},
				Events: map[string]*config.ChainReaderEvent{
					"ccip_message_sent": {
						Name:      "ccip_message_sent",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "onramp",
							Event:   "CCIPMessageSent",
						},
						ExpectedEventType: &CCIPMessageSent{},
					},
				},
			},
			"Router": {
				Name: "router",
				Functions: map[string]*config.ChainReaderFunction{
					"get_mock_onramp_address": {
						Name:          "get_mock_onramp_address",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{}, // No parameters needed
					},
				},
				Events: map[string]*config.ChainReaderEvent{},
			},
			"FeeQuoter": {
				Name: "fee_quoter",
				Functions: map[string]*config.ChainReaderFunction{
					"emit_usd_per_token_updated_event": {
						Name:          "emit_usd_per_token_updated_event",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type: "address",
								Name: "token",
							},
							{
								Type: "u256",
								Name: "usd_per_token",
							},
							{
								Type: "u64",
								Name: "timestamp",
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{
					"usd_per_token_updated": {
						Name:      "usd_per_token_updated",
						EventType: "UsdPerTokenUpdated",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "fee_quoter",
							Event:   "UsdPerTokenUpdated",
						},
						EventFieldRenames: map[string]aptosCRConfig.RenamedField{
							"token":       {NewName: "myToken"},
							"usdPerToken": {NewName: "dollarPerToken"},
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

	offRampBinding := types.BoundContract{
		Name:    "OffRamp",
		Address: packageId, // Package ID of the deployed offramp contract
	}

	onRampBinding := types.BoundContract{
		Name:    "OnRamp",
		Address: packageId, // Package ID of the deployed onramp contract
	}

	routerBinding := types.BoundContract{
		Name:    "Router",
		Address: packageId, // Package ID of the deployed router contract
	}

	feeQuoterBinding := types.BoundContract{
		Name:    "FeeQuoter",
		Address: packageId, // Package ID of the deployed fee_quoter contract
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
		chainReaderConfig.TransactionsIndexer.PollingInterval,
		chainReaderConfig.TransactionsIndexer.SyncTimeout,
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
		chainReaderConfig.EventsIndexer.PollingInterval,
		chainReaderConfig.EventsIndexer.SyncTimeout,
	)
	indexerInstance := indexer.NewIndexer(
		log,
		evIndexer,
		txnIndexer,
	)

	// ChainReader in non-loop mode
	chainReader, err := NewChainReader(ctx, log, relayerClient, chainReaderConfig, db, indexerInstance)
	require.NoError(t, err)

	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding, offRampBinding, onRampBinding, routerBinding, feeQuoterBinding})
	require.NoError(t, err)

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

	// Verify renamed field on simple struct output
	t.Run("GetLatestValue_SimpleStruct_Renamed", func(t *testing.T) {
		var renamed map[string]any

		log.Debugw("Testing get_simple_result with field rename",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_simple_result_renamed"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed
			&renamed,
		)
		require.NoError(t, err)

		require.NotNil(t, renamed)
		// original key should not be present when renamed
		_, hasOriginal := renamed["value"]
		require.False(t, hasOriginal)
		require.Equal(t, "42", renamed["renamedValue"])
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

	t.Run("GetLatestValue_Address", func(t *testing.T) {
		var retAddress any

		log.Debugw("Testing get_mock_onramp_address function",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Router", "get_mock_onramp_address"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed
			&retAddress,
		)
		require.NoError(t, err)

		// Verify the returned struct
		require.NotNil(t, retAddress)
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", retAddress, "Expected address to be 0x0000000000000000000000000000000000000000000000000000000000000001")
	})

	// Verify renamed fields on address list output
	t.Run("GetLatestValue_AddressList_Renamed", func(t *testing.T) {
		var renamed map[string]any

		log.Debugw("Testing get_address_list with field rename",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_address_list_renamed"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed
			&renamed,
		)
		require.NoError(t, err)

		require.NotNil(t, renamed)
		// renamed keys should be present
		require.Contains(t, renamed, "wallets")
		require.Contains(t, renamed, "size")
		// original keys should not be present
		require.NotContains(t, renamed, "addresses")
		require.NotContains(t, renamed, "count")
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
		require.Equal(t, "42", retTupleStruct["value"], "Expected value to be 42")
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", retTupleStruct["address"], "Expected address to be 0x0000000000000000000000000000000000000000000000000000000000000001")
		require.Equal(t, true, retTupleStruct["bool"], "Expected bool to be true")

		log.Debugw("TupleStruct test completed successfully",
			"value", retTupleStruct["value"],
			"address", retTupleStruct["address"],
			"bool", retTupleStruct["bool"],
			"struct_tag", retTupleStruct["struct_tag"])
	})

	// Verify renamed fields on tuple-to-struct output
	t.Run("GetLatestValue_TupleToStruct_Renamed", func(t *testing.T) {
		var renamed map[string]any

		log.Debugw("Testing get_tuple_struct with field rename",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_tuple_struct_renamed"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed
			&renamed,
		)
		require.NoError(t, err)

		require.NotNil(t, renamed)
		// renamed keys should be present
		require.Contains(t, renamed, "answer")
		require.Contains(t, renamed, "tag")
		// original keys should not be present
		require.NotContains(t, renamed, "value")
		require.NotContains(t, renamed, "struct_tag")
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

	t.Run("QueryKeyWithMetadata_RenamedFields", func(t *testing.T) {
		// Increment the counter to emit an event
		log.Debugw("Emitting UsdPerTokenUpdated event")

		// Use relayerClient to call increment instead of using CLI
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "fee_quoter",
			Function:        "emit_usd_per_token_updated_event",
			TypeArguments: []any{
				"address",
				"u256",
				"u64",
			},
			Arguments: []any{
				"0x0000000000000000000000000000000000000000000000000000000000000001",
				"1000000000000000000",
				"1714953600",
			},
			GasBudget: 2000000,
		}

		log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

		txMetadata, testErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, testErr)

		_, testErr = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, testErr)

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "usd_per_token_updated",
		}

		// Setup limit and sort
		limitAndSort := query.LimitAndSort{
			Limit: query.Limit{
				Count:  50,
				Cursor: "",
			},
		}

		log.Debugw("Querying for UsdPerTokenUpdated events",
			"filter", filter.Key,
			"limit", limitAndSort.Limit.Count,
			"packageId", packageId,
			"contract", feeQuoterBinding.Name,
			"eventType", "UsdPerTokenUpdated")

		sequences := []aptosCRConfig.SequenceWithMetadata{}
		require.Eventually(t, func() bool {
			// Query for events
			var usdPerTokenUpdated any
			sequences, err = chainReader.(chainreader.ExtendedContractReader).QueryKeyWithMetadata(
				ctx,
				feeQuoterBinding,
				filter,
				limitAndSort,
				&usdPerTokenUpdated,
			)
			if err != nil {
				log.Errorw("Failed to query events", err)
				require.NoError(t, err)
			}

			return len(sequences) > 0
		}, 60*time.Second, 1*time.Second, "Event should eventually be indexed and found")

		log.Debugw("Query results", "sequences", sequences)

		// Verify we got at least one event
		require.NotEmpty(t, sequences, "Expected at least one event")

		// Verify the event data
		event := sequences[0].Sequence.Data
		require.NotNil(t, event)

		log.Debugw("Sequence data", "sequenceData", event)

		// Check the fields of the event
		eventMap, ok := (*event.(*any)).(map[string]any)
		require.True(t, ok, "Event data should be a map")

		require.Contains(t, eventMap, "myToken")
		require.Contains(t, eventMap, "dollarPerToken")
		require.Contains(t, eventMap, "timestamp")
	})

	t.Run("QueryKey_CCIPMessageSent_concrete_sequenceDataType", func(t *testing.T) {
		// Increment the counter to emit an event
		log.Debugw("Emitting CCIPMessageSent event")

		// Use relayerClient to call increment instead of using CLI
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "onramp",
			Function:        "emit_sample_ccip_message_sent_event",
			TypeArguments:   []any{},
			Arguments:       []any{},
			GasBudget:       2000000,
		}

		log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

		txMetadata, testErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, testErr)

		_, testErr = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, testErr)

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "ccip_message_sent",
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
			"contract", onRampBinding.Name,
			"eventType", "CCIPMessageSent")

		sequences := []types.Sequence{}
		require.Eventually(t, func() bool {
			// Query for events
			var ccipMessageSent CCIPMessageSent
			sequences, err = chainReader.QueryKey(
				ctx,
				onRampBinding,
				filter,
				limitAndSort,
				&ccipMessageSent,
			)
			if err != nil {
				log.Errorw("Failed to query events", err)
				require.NoError(t, err)
			}

			return len(sequences) > 0
		}, 60*time.Second, 1*time.Second, "Event should eventually be indexed and found")

		log.Debugw("Query results", "sequences", sequences)

		// Verify we got at least one event
		require.NotEmpty(t, sequences, "Expected at least one event")

		// Verify the event data, it should be castable since we specified a concrete type in the request
		event := sequences[0].Data.(*CCIPMessageSent)
		require.NotNil(t, event)

		require.Equal(t, uint64(1), event.SequenceNumber, "Expected sequence number to be 1")
		require.Equal(t, "0xabcdef123456", event.Message.Data, "Expected data to be 0xabcdef123456")
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000789", event.Message.Sender, "Expected sender to be 0x00...0000789")
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000abc", event.Message.FeeToken, "Expected feeToken to be 0x00...0000abc")
		require.Equal(t, uint64(500), event.Message.FeeTokenAmount, "Expected feeTokenAmount to be 500")
		require.Equal(t, "0xabcdef123456", event.Message.Receiver, "Receiver must be hex encoded")
	})

	t.Run("QueryKey_CCIPMessageSent_map[string]any_sequenceDataType", func(t *testing.T) {
		// Increment the counter to emit an event
		log.Debugw("Emitting CCIPMessageSent event")

		// Use relayerClient to call increment instead of using CLI
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "onramp",
			Function:        "emit_sample_ccip_message_sent_event",
			TypeArguments:   []any{},
			Arguments:       []any{},
			GasBudget:       2000000,
		}

		log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

		txMetadata, testErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, testErr)

		_, testErr = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, testErr)

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "ccip_message_sent",
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
			"contract", onRampBinding.Name,
			"eventType", "CCIPMessageSent")

		sequences := []types.Sequence{}
		require.Eventually(t, func() bool {
			// Query for events
			var ccipMessageSent map[string]any
			sequences, err = chainReader.QueryKey(
				ctx,
				onRampBinding,
				filter,
				limitAndSort,
				&ccipMessageSent,
			)
			if err != nil {
				log.Errorw("Failed to query events", err)
				require.NoError(t, err)
			}

			return len(sequences) > 0
		}, 60*time.Second, 1*time.Second, "Event should eventually be indexed and found")

		log.Debugw("Query results", "sequences", sequences)

		// Verify we got at least one event
		require.NotEmpty(t, sequences, "Expected at least one event")

		// Verify the event data
		event := sequences[0].Data
		require.NotNil(t, event)

		// Marshal and unmarshal the event data just to check values are correct in a simple way
		jsonData, err := json.Marshal(event)
		require.NoError(t, err)

		var typedData CCIPMessageSent
		err = json.Unmarshal(jsonData, &typedData)
		require.NoError(t, err)

		require.Equal(t, uint64(1), typedData.SequenceNumber, "Expected sequence number to be 1")
		require.Equal(t, "0xabcdef123456", typedData.Message.Data, "Expected data to be 0xabcdef123456")
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000789", typedData.Message.Sender, "Expected sender to be 0x00...0000789")
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000abc", typedData.Message.FeeToken, "Expected feeToken to be 0x00...0000abc")
		require.Equal(t, uint64(500), typedData.Message.FeeTokenAmount, "Expected feeTokenAmount to be 500")
		require.Equal(t, "0xabcdef123456", typedData.Message.Receiver, "Receiver must be hex encoded")
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

	t.Run("QueryKey_WithMetadata_CCIPMessageSent_untyped", func(t *testing.T) {
		// Increment the counter to emit an event
		log.Debugw("Emitting CCIPMessageSent event")

		// Use relayerClient to call increment instead of using CLI
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "onramp",
			Function:        "emit_sample_ccip_message_sent_event",
			TypeArguments:   []any{},
			Arguments:       []any{},
			GasBudget:       2000000,
		}

		log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

		txMetadata, testErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, testErr)

		_, testErr = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, testErr)

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "ccip_message_sent",
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
			"contract", onRampBinding.Name,
			"eventType", "CCIPMessageSent")

		sequences := []aptosCRConfig.SequenceWithMetadata{}
		require.Eventually(t, func() bool {
			// Query for events
			var ccipMessageSent map[string]any
			sequences, err = chainReader.(chainreader.ExtendedContractReader).QueryKeyWithMetadata(
				ctx,
				onRampBinding,
				filter,
				limitAndSort,
				&ccipMessageSent,
			)
			if err != nil {
				log.Errorw("Failed to query events", err)
				require.NoError(t, err)
			}

			return len(sequences) > 0
		}, 60*time.Second, 1*time.Second, "Event should eventually be indexed and found")

		log.Debugw("Query results", "sequences", sequences)

		// Verify we got at least one event
		require.NotEmpty(t, sequences, "Expected at least one event")

		// Verify the event data
		event := sequences[0].Sequence.Data
		require.NotNil(t, event)

		// Marshal and unmarshal the event data just to check values are correct in a simple way
		jsonData, err := json.Marshal(event)
		require.NoError(t, err)

		var typedData CCIPMessageSent
		err = json.Unmarshal(jsonData, &typedData)
		require.NoError(t, err)

		require.Equal(t, uint64(1), typedData.SequenceNumber, "Expected sequence number to be 1")
		require.Equal(t, "0xabcdef123456", typedData.Message.Data, "Expected data to be 0xabcdef123456")
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000789", typedData.Message.Sender, "Expected sender to be 0x00...0000789")
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000abc", typedData.Message.FeeToken, "Expected feeToken to be 0x00...0000abc")
		require.Equal(t, uint64(500), typedData.Message.FeeTokenAmount, "Expected feeTokenAmount to be 500")
		require.Equal(t, "0xabcdef123456", typedData.Message.Receiver, "Receiver must be hex encoded")
	})

	t.Run("QueryKey_WithMetadata_ExecutionStateChanged_untyped", func(t *testing.T) {
		// Increment the counter to emit an event
		log.Debugw("Emitting ExecutionStateChanged event")

		// Use relayerClient to call increment instead of using CLI
		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "offramp",
			Function:        "emit_sample_execution_state_changed_event",
			TypeArguments:   []any{},
			Arguments:       []any{},
			GasBudget:       2000000,
		}

		log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

		txMetadata, testErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, testErr)

		_, testErr = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, testErr)

		// Create a filter for events
		filter := query.KeyFilter{
			Key: "execution_state_changed",
		}

		// Setup limit and sort
		limitAndSort := query.LimitAndSort{
			Limit: query.Limit{
				Count:  50,
				Cursor: "",
			},
		}

		log.Debugw("Querying for contract events",
			"filter", filter.Key,
			"limit", limitAndSort.Limit.Count,
			"packageId", packageId,
			"contract", offRampBinding.Name,
			"eventType", "ExecutionStateChanged")

		sequences := []aptosCRConfig.SequenceWithMetadata{}
		require.Eventually(t, func() bool {
			// Query for events
			var executionStateChanged map[string]any
			sequences, err = chainReader.(chainreader.ExtendedContractReader).QueryKeyWithMetadata(
				ctx,
				offRampBinding,
				filter,
				limitAndSort,
				&executionStateChanged,
			)
			if err != nil {
				log.Errorw("Failed to query events", err)
				require.NoError(t, err)
			}

			return len(sequences) > 0
		}, 60*time.Second, 1*time.Second, "Event should eventually be indexed and found")

		log.Debugw("Query results", "sequences", sequences)

		// Verify we got at least one event
		require.NotEmpty(t, sequences, "Expected at least one event")

		log.Debugw("Event data", "event", sequences[0].Sequence.Data)

		// Verify the event data
		event := sequences[0].Sequence.Data
		require.NotNil(t, event)

		// Marshal and unmarshal the event data just to check values are correct in a simple way
		jsonData, err := json.Marshal(event)
		require.NoError(t, err)

		var typedData OfframpExecutionStateChanged
		err = json.Unmarshal(jsonData, &typedData)
		require.NoError(t, err)

		require.Equal(t, uint64(12), typedData.SequenceNumber, "Expected sequence number to be 12")
		require.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", typedData.MessageId, "Expected message id to be 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		require.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", typedData.MessageHash, "Expected message hash to be 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		require.Equal(t, uint8(1), typedData.State, "Expected state to be 1")
		require.Equal(t, uint64(24), typedData.SourceChainSelector, "Expected source chain selector to be 24")
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

	t.Run("GetLatestValue_WithSecondaryPointerTag", func(t *testing.T) {
		expectedUint64 := uint64(5)
		var retUint64 uint64

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "get_value_with_pointer_dependency"}, "-"),
			primitives.Finalized,
			map[string]any{}, // No parameters needed, pointer tags will take care of it
			&retUint64,
		)
		require.NoError(t, err)

		// Verify the returned struct
		require.NotNil(t, retUint64)
		require.Equal(t, expectedUint64, retUint64, "Expected value to be 5")
	})

	t.Run("GetLatestValue_GetAllSourceChainConfigs", func(t *testing.T) {
		var retAllSourceChainConfigs any
		params := map[string]any{}

		log.Debugw("Testing get_all_source_chain_configs function for BCS struct decoding",
			"packageId", packageId,
		)

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "OffRamp", "get_all_source_chain_configs"}, "-"),
			primitives.Finalized,
			&params, // no parameters needed
			&retAllSourceChainConfigs,
		)
		require.NoError(t, err)

		// Verify the returned data structure
		require.NotNil(t, retAllSourceChainConfigs)
		require.Len(t, retAllSourceChainConfigs, 2, "Expected 2 elements in the response")

		// Define JSON schema for the expected response format
		expectedSchema := `{
			"$schema": "https://json-schema.org/draft/2019-09/schema",
			"type": "array",
			"prefixItems": [
				{
					"type": "array",
					"items": {
						"type": "integer"
					}
				},
				{
					"type": "array",
					"items": {
						"type": "object",
						"properties": {
							"is_enabled": {
								"type": "boolean"
							},
							"is_rmn_verification_disabled": {
								"type": "boolean"
							},
							"min_seq_nr": {
								"type": "integer"
							},
							"on_ramp": {
								"type": "string",
								"pattern": "^0x[a-fA-F0-9]{64}$"
							},
							"router": {
								"type": "string",
								"pattern": "^0x[a-fA-F0-9]{64}$"
							}
						},
						"required": ["is_enabled", "is_rmn_verification_disabled", "min_seq_nr", "on_ramp", "router"],
						"additionalProperties": false
					}
				}
			],
			"minItems": 2,
			"maxItems": 2
		}`

		jsonResult, err := json.Marshal(retAllSourceChainConfigs)
		require.NoError(t, err)
		log.Debugw("jsonResult", "jsonResult", string(jsonResult))

		err = testutils.ValidateJSON(retAllSourceChainConfigs, expectedSchema)
		require.NoError(t, err)
	})

	t.Run("GetLatestValue_ResponseFromInputs", func(t *testing.T) {
		var retResponseFromInputs any
		params := map[string]any{}

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "response_from_inputs"}, "-"),
			primitives.Finalized,
			&params, // no parameters needed
			&retResponseFromInputs,
		)
		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, retResponseFromInputs, "retResponseFromInputs")
	})

	t.Run("GetLatestValue_StaticResponse", func(t *testing.T) {
		var retStaticResponse any
		params := map[string]any{}

		err = chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, "Counter", "static_response"}, "-"),
			primitives.Finalized,
			&params, // no parameters needed
			&retStaticResponse,
		)
		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, retStaticResponse, "retStaticResponse")
		require.Equal(t, map[string]any{"a": 1, "b": 2, "c": 3}, retStaticResponse, "Expected static response to be map[string]any with keys a, b, and c")
	})
}
