//go:build integration

package chainreader

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestChainReaderLocal(t *testing.T) {
	logger := logger.Test(t)

	privateKey, publicKey, accountAddress := testutils.LoadAccountFromEnv(t, logger)
	// if the env does not contain a private key to be loaded, create one
	if privateKey == nil {
		privateKey, publicKey, accountAddress = testutils.GenerateAccountKeyPair(t, logger)
	}

	err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	logger.Debugw("Started Sui node")

	err = testutils.FundWithFaucet(logger, constant.SuiLocalnet, accountAddress)
	require.NoError(t, err)

	runChainReaderCounterTest(t, logger, testutils.LocalUrl)
}

func runChainReaderCounterTest(t *testing.T, logger logger.Logger, rpcUrl string) {
	client := sui.NewSuiClient(rpcUrl)

	// start by deploying the counter contract to local net
	packageId, counterObjectId, err := testutils.DeployCounterContract(t)
	require.NoError(t, err)

	// Set up the ChainReader
	chainReaderConfig := ChainReaderConfig{
		Modules: map[string]*ChainReaderModule{
			"counter": {
				Name: "counter",
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "counter",
		Address: packageId, // Package ID of the deployed counter contract
	}

	chainReader := NewChainReader(logger, client, chainReaderConfig)
	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	logger.Debugw("ChainReader setup complete")

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
}

// func runChainReaderTest(t *testing.T, logger logger.Logger, rpcUrl string, accountAddress string, publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) {
// 	keystore := testutils.NewTestKeystore(t)
// 	keystore.AddKey(privateKey)

// 	client := sui.NewSuiClient(rpcUrl)

// 	// Compile and publish test module
// 	packageObjectId := "0xabf6d4f0df648f9305bb32f8de07486b438d597779d9611e64c3add4044e71fc"

// 	// Initialize modules
// 	err := testutils.CallMoveFunction(
// 		client,
// 		privateKey,
// 		accountAddress,
// 		packageObjectId,
// 		"echo",
// 		"init",
// 		[]string{},      // type args
// 		[]any{}, // args
// 	)
// 	require.NoError(t, err)

// 	// Wait for the transaction to be confirmed
// 	time.Sleep(2 * time.Second)

// 	// Configure the ChainReader
// 	config := ChainReaderConfig{
// 		Modules: map[string]*ChainReaderModule{
// 			"testContract": {
// 				Name: "echo",
// 				Functions: map[string]*ChainReaderFunction{
// 					"echo_u64": {
// 						Params: []codec.SuiFunctionParam{
// 							{
// 								Name:     "value",
// 								Type:     "u64",
// 								Required: true,
// 							},
// 						},
// 					},
// 					"echo_u32_u64_tuple": {
// 						Params: []codec.SuiFunctionParam{
// 							{
// 								Name:     "value1",
// 								Type:     "u32",
// 								Required: true,
// 							},
// 							{
// 								Name:     "value2",
// 								Type:     "u64",
// 								Required: true,
// 							},
// 						},
// 					},
// 					"echo_string": {
// 						Params: []codec.SuiFunctionParam{
// 							{
// 								Name:     "value",
// 								Type:     "string",
// 								Required: true,
// 							},
// 						},
// 					},
// 					"echo_byte_vector": {
// 						Params: []codec.SuiFunctionParam{
// 							{
// 								Name:     "value",
// 								Type:     "vector<u8>",
// 								Required: true,
// 							},
// 						},
// 					},
// 					"echo_u256": {
// 						Params: []codec.SuiFunctionParam{
// 							{
// 								Name:     "value",
// 								Type:     "u256",
// 								Required: true,
// 							},
// 						},
// 					},
// 				},
// 				Events: map[string]*ChainReaderEvent{
// 					"SingleValueEvent": {
// 						EventType: "SingleValueEvent",
// 					},
// 					"DoubleValueEvent": {
// 						EventType: "DoubleValueEvent",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	binding := types.BoundContract{
// 		Name:    "testContract",
// 		Address: packageObjectId,
// 	}

// 	chainReader := NewChainReader(logger, client, config)
// 	err = chainReader.Bind(context.Background(), []types.BoundContract{binding})
// 	require.NoError(t, err)

// 	// Test GetLatestValue for different data types
// 	t.Run("GetLatestValue_Uint64", func(t *testing.T) {
// 		expectedUint64 := uint64(42)
// 		var retUint64 uint64
// 		err = chainReader.GetLatestValue(
// 			context.Background(),
// 			binding.ReadIdentifier("echo_u64"),
// 			primitives.Finalized,
// 			struct {
// 				Value uint64
// 			}{Value: expectedUint64},
// 			&retUint64,
// 		)
// 		require.NoError(t, err)
// 		require.Equal(t, expectedUint64, retUint64)
// 	})

// 	t.Run("GetLatestValue_String", func(t *testing.T) {
// 		expectedString := "hello world"
// 		var retString string
// 		err = chainReader.GetLatestValue(
// 			context.Background(),
// 			binding.ReadIdentifier("echo_string"),
// 			primitives.Finalized,
// 			struct {
// 				Value string
// 			}{Value: expectedString},
// 			&retString,
// 		)
// 		require.NoError(t, err)
// 		require.Equal(t, expectedString, retString)
// 	})

// 	// Test BatchGetLatestValues
// 	t.Run("BatchGetLatestValues", func(t *testing.T) {
// 		var retUint64 uint64
// 		var retString string

// 		batch := []types.BatchRead{
// 			{
// 				ReadName: "echo_u64",
// 				Params: struct {
// 					Value uint64
// 				}{Value: 123},
// 				ReturnVal: &retUint64,
// 			},
// 			{
// 				ReadName: "echo_string",
// 				Params: struct {
// 					Value string
// 				}{Value: "batch test"},
// 				ReturnVal: &retString,
// 			},
// 		}

// 		request := types.BatchGetLatestValuesRequest{
// 			binding: batch,
// 		}

// 		result, err := chainReader.BatchGetLatestValues(context.Background(), request)
// 		require.NoError(t, err)
// 		require.Len(t, result[binding], 2)

// 		// Check the results
// 		require.Equal(t, "echo_u64", result[binding][0].ReadName)
// 		require.Equal(t, "echo_string", result[binding][1].ReadName)
// 		require.Equal(t, uint64(123), retUint64)
// 		require.Equal(t, "batch test", retString)
// 	})

// 	// Test events by emitting some test events
// 	t.Run("QueryKey_Events", func(t *testing.T) {
// 		// First, emit some test events
// 		for i := 0; i < 3; i++ {
// 			err := testutils.CallMoveFunction(
// 				client,
// 				privateKey,
// 				accountAddress,
// 				packageObjectId,
// 				"echo",
// 				"echo_with_events",
// 				[]string{}, // type args
// 				[]any{
// 					uint64(i + 1),                  // number
// 					"Event " + uuid.New().String(), // text
// 					[]byte("test data"),            // bytes
// 				},
// 			)
// 			require.NoError(t, err)
// 			time.Sleep(1 * time.Second) // Give time for events to be processed
// 		}

// 		// Now query the events
// 		type SingleValueEventData struct {
// 			Value uint64 `json:"value"`
// 		}

// 		filter := query.KeyFilter{
// 			Key: "SingleValueEvent",
// 		}

// 		limitAndSort := query.LimitAndSort{
// 			Limit: query.Limit{
// 				Count: 10,
// 			},
// 			SortBy: []query.SortType{
// 				query.SortBySequence{Order: query.Asc},
// 			},
// 		}

// 		sequences, err := chainReader.QueryKey(
// 			context.Background(),
// 			binding,
// 			filter,
// 			limitAndSort,
// 			&SingleValueEventData{},
// 		)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, sequences)

// 		// Verify the events
// 		for _, seq := range sequences {
// 			event, ok := seq.Data.(*SingleValueEventData)
// 			require.True(t, ok)
// 			require.Greater(t, event.Value, uint64(0))
// 		}
// 	})

// 	// Test unbinding
// 	t.Run("Unbind", func(t *testing.T) {
// 		err := chainReader.Unbind(context.Background(), []types.BoundContract{binding})
// 		require.NoError(t, err)

// 		// Trying to read after unbinding should fail
// 		var retUint64 uint64
// 		err = chainReader.GetLatestValue(
// 			context.Background(),
// 			binding.ReadIdentifier("echo_u64"),
// 			primitives.Finalized,
// 			struct {
// 				Value uint64
// 			}{Value: 42},
// 			&retUint64,
// 		)
// 		require.Error(t, err)
// 	})
// }

func getSampleTxMetadata() *types.TxMeta {
	workflowID := "sample-workflow-id"
	return &types.TxMeta{
		WorkflowExecutionID: &workflowID,
		GasLimit:            big.NewInt(21000),
	}
}
