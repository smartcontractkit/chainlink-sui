//go:build integration

package loop

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestLoopChainReaderLocal(t *testing.T) {
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

	runLoopChainReaderEchoTest(t, log, testutils.LocalUrl)
}

func runLoopChainReaderEchoTest(t *testing.T, log logger.Logger, rpcUrl string) {
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
	log.Debugw("Publish output", "output", publishOutput)

	// Set up the base ChainReader with echo function configurations
	chainReaderConfig := chainreader.ChainReaderConfig{
		IsLoopPlugin: true,
		Modules: map[string]*chainreader.ChainReaderModule{
			"echo": {
				Name: "echo",
				Functions: map[string]*chainreader.ChainReaderFunction{
					"echo_u64": {
						Name:          "echo_u64",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "u64",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_u256": {
						Name:          "echo_u256",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "u256",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_u32_u64_tuple": {
						Name:          "echo_u32_u64_tuple",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "u32",
								Name:     "val1",
								Required: true,
							},
							{
								Type:     "u64",
								Name:     "val2",
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"first", "second"},
					},
					"echo_string": {
						Name:          "echo_string",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "0x1::string::String",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_byte_vector": {
						Name:          "echo_byte_vector",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "vector<u8>",
								Name:     "val",
								Required: true,
							},
						},
					},
					"echo_byte_vector_vector": {
						Name:          "echo_byte_vector_vector",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "vector<vector<u8>>",
								Name:     "val",
								Required: true,
							},
						},
					},
					"simple_event_echo": {
						Name:          "simple_event_echo",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:     "u64",
								Name:     "number",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*chainreader.ChainReaderEvent{
					"single_value_event": {
						Name:      "single_value_event",
						EventType: "SingleValueEvent",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "echo",
							Event:   "SingleValueEvent",
						},
					},
					"double_value_event": {
						Name:      "double_value_event",
						EventType: "DoubleValueEvent",
					},
					"triple_value_event": {
						Name:      "triple_value_event",
						EventType: "TripleValueEvent",
					},
				},
			},
			"counter": {
				Name: "counter",
				Functions: map[string]*chainreader.ChainReaderFunction{
					"get_tuple_struct": {
						Name:                "get_tuple_struct",
						SignerAddress:       accountAddress,
						Params:              []codec.SuiFunctionParam{},
						ResultTupleToStruct: []string{"value", "address", "bool", "struct_tag"},
					},
					"get_ocr_config": {
						Name:          "get_ocr_config",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{},
						// used to wrap entire result
						ResultTupleToStruct: []string{"OCRConfig"},
					},
				},
			},
		},
	}

	echoBinding := types.BoundContract{
		Name:    "echo",
		Address: packageId, // Package ID of the deployed echo contract
	}

	counterBinding := types.BoundContract{
		Name:    "counter",
		Address: packageId, // Package ID of the deployed echo contract
	}

	// Set up DB
	datastoreUrl := os.Getenv("TEST_DB_URL")
	if datastoreUrl == "" {
		t.Skip("Skipping persistent tests as TEST_DB_URL is not set in CI")
	}
	db := sqltest.NewDB(t, datastoreUrl)

	// Create the base chain reader
	chainReader, err := chainreader.NewChainReader(ctx, log, relayerClient, chainReaderConfig, db)
	require.NoError(t, err)

	// Wrap the base chain reader with loop chain reader
	loopReader := NewLoopChainReader(log, chainReader)

	// Bind the contracts to the loop reader
	err = loopReader.Bind(context.Background(), []types.BoundContract{echoBinding})
	require.NoError(t, err)

	err = loopReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	log.Debugw("LoopChainReader setup complete")

	// Test 1: echo_u64 function call
	t.Run("LoopReader_GetLatestValue_EchoU64", func(t *testing.T) {
		testValue := uint64(42)
		var retUint64 uint64

		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u64"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testValue,
			},
			&retUint64,
		)
		require.NoError(t, err)
		require.Equal(t, testValue, retUint64)
	})

	// Test 2: echo_u64 with different values
	t.Run("LoopReader_GetLatestValue_EchoU64_VariousValues", func(t *testing.T) {
		testCases := []uint64{0, 1, 100, 1000, 1000000000}

		for _, testValue := range testCases {
			t.Run(fmt.Sprintf("Value_%d", testValue), func(t *testing.T) {
				var retUint64 uint64
				err = loopReader.GetLatestValue(
					context.Background(),
					strings.Join([]string{packageId, echoBinding.Name, "echo_u64"}, "-"),
					primitives.Finalized,
					map[string]any{
						"val": testValue,
					},
					&retUint64,
				)
				require.NoError(t, err)
				require.Equal(t, testValue, retUint64)
			})
		}
	})

	// Test 3: echo_u256 function call
	t.Run("LoopReader_GetLatestValue_EchoU256", func(t *testing.T) {
		t.Skip("Skipping u256 test")

		testValue := big.NewInt(123456789)
		var retBigInt *big.Int
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u256"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testValue,
			},
			&retBigInt,
		)
		require.NoError(t, err)
		require.Equal(t, testValue, retBigInt)
	})

	// Test 4: echo_u256 with large values
	t.Run("LoopReader_GetLatestValue_EchoU256_LargeValue", func(t *testing.T) {
		t.Skip("Skipping large value test")

		// Test with a very large number
		testValue := new(big.Int)
		testValue.SetString("123456789012345678901234567890", 10)
		var retBigInt *big.Int
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u256"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testValue,
			},
			&retBigInt,
		)
		require.NoError(t, err)
		require.Equal(t, testValue, retBigInt)
	})

	// Test 5: echo_u32_u64_tuple function call
	t.Run("LoopReader_GetLatestValue_EchoTuple", func(t *testing.T) {
		testVal1 := uint32(100)
		testVal2 := uint64(200)

		type TupleResult struct {
			First  uint32 `json:"first"`
			Second uint64 `json:"second"`
		}

		var retTuple TupleResult
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_u32_u64_tuple"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val1": testVal1,
				"val2": testVal2,
			},
			&retTuple,
		)
		require.NoError(t, err)
		require.Equal(t, testVal1, retTuple.First)
		require.Equal(t, testVal2, retTuple.Second)
	})

	// Test 6: echo_string function call
	t.Run("LoopReader_GetLatestValue_EchoString", func(t *testing.T) {
		t.Skip("Skipping string test")

		testString := "Hello, Sui!"
		var retString string
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, echoBinding.Name, "echo_string"}, "-"),
			primitives.Finalized,
			map[string]any{
				"val": testString,
			},
			&retString,
		)
		require.NoError(t, err)
		require.Equal(t, testString, retString)
	})

	// Test 7: echo_with_events function call and event querying
	t.Run("LoopReader_EchoWithEvents_AndQueryEvents", func(t *testing.T) {
		// Test data
		testNumber := uint64(12345)
		testText := "Hello Events!"
		testBytes := []byte("test bytes data")

		// First, call the function that emits events
		var retUint64 uint64
		err = loopReader.GetLatestValue(
			ctx,
			strings.Join([]string{packageId, echoBinding.Name, "simple_event_echo"}, "-"),
			primitives.Finalized,
			map[string]any{
				"number": testNumber,
			},
			&retUint64,
		)
		require.NoError(t, err)

		// Define event structures to match the Move contract
		type SingleValueEvent struct {
			Value uint64 `json:"value"`
		}

		type NoConfigSingleValueEvent struct {
			Value uint64 `json:"value"`
		}

		type DoubleValueEvent struct {
			Number uint64 `json:"number"`
			Text   string `json:"text"`
		}

		type TripleValueEvent struct {
			Values [][]byte `json:"values"`
		}

		// Query for SingleValueEvent
		t.Run("QuerySingleValueEvent", func(t *testing.T) {
			singleValueEvent := &SingleValueEvent{}
			var sequences []types.Sequence
			//nolint:govet
			var err error

			// Use relayerClient to call increment instead of using CLI
			moveCallReq := client.MoveCallRequest{
				Signer:          accountAddress,
				PackageObjectId: packageId,
				Module:          "echo",
				Function:        "simple_event_echo",
				TypeArguments:   []any{},
				Arguments: []any{
					testNumber,
				},
				GasBudget: "2000000",
			}

			log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

			txMetadata, err := relayerClient.MoveCall(ctx, moveCallReq)
			require.NoError(t, err)

			_, err = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				sequences, err = loopReader.QueryKey(
					ctx,
					echoBinding,
					query.KeyFilter{
						Key: "single_value_event",
					},
					query.LimitAndSort{
						SortBy: []query.SortBy{},
						Limit:  query.CountLimit(10),
					},
					singleValueEvent,
				)

				if err != nil {
					log.Errorw("Error querying for SingleValueEvent", "err", err)
				}

				return err == nil && len(sequences) > 0
			}, 30*time.Second, 1*time.Second)

			require.NoError(t, err)
			require.NotEmpty(t, sequences, "Expected to find SingleValueEvent")
			log.Debugw("Sequences found", "sequences", sequences)
		})

		// Query for SingleValueEvent
		t.Run("QuerySingleValueEvent_WithoutConfig", func(t *testing.T) {
			singleValueEvent := &NoConfigSingleValueEvent{}
			var sequences []types.Sequence
			//nolint:govet
			var err error

			// Use relayerClient to call increment instead of using CLI
			moveCallReq := client.MoveCallRequest{
				Signer:          accountAddress,
				PackageObjectId: packageId,
				Module:          "echo",
				Function:        "no_config_event_echo",
				TypeArguments:   []any{},
				Arguments: []any{
					testNumber,
				},
				GasBudget: "2000000",
			}

			log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

			txMetadata, err := relayerClient.MoveCall(ctx, moveCallReq)
			require.NoError(t, err)

			_, err = relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				sequences, err = loopReader.QueryKey(
					ctx,
					echoBinding,
					query.KeyFilter{
						Key: "NoConfigSingleValueEvent",
					},
					query.LimitAndSort{
						SortBy: []query.SortBy{},
						Limit:  query.CountLimit(10),
					},
					singleValueEvent,
				)

				if err != nil {
					log.Errorw("Error querying for NoValueSingleValueEvent", "err", err)
				}

				return err == nil && len(sequences) > 0
			}, 30*time.Second, 1*time.Second)

			require.NoError(t, err)
			require.NotEmpty(t, sequences, "Expected to find SingleValueEvent")
			log.Debugw("Sequences found", "sequences", sequences)
		})

		// Query for DoubleValueEvent
		t.Run("QueryDoubleValueEvent", func(t *testing.T) {
			t.Skip("Skipping double value test")
			doubleValueEvent := &DoubleValueEvent{}
			//nolint:govet
			sequences, err := loopReader.QueryKey(
				ctx,
				echoBinding,
				query.KeyFilter{
					Key: "double_value_event",
				},
				query.LimitAndSort{
					Limit: query.CountLimit(10),
				},
				doubleValueEvent,
			)
			require.NoError(t, err)
			require.NotEmpty(t, sequences, "Expected to find DoubleValueEvent")

			// Check the latest event
			latestEvent := sequences[0].Data.(*DoubleValueEvent)
			require.Equal(t, testNumber, latestEvent.Number)
			require.Equal(t, testText, latestEvent.Text)
		})

		// Query for TripleValueEvent
		t.Run("QueryTripleValueEvent", func(t *testing.T) {
			t.Skip("Skipping triple value event test")
			tripleValueEvent := &TripleValueEvent{}
			//nolint:govet
			sequences, err := loopReader.QueryKey(
				ctx,
				echoBinding,
				query.KeyFilter{
					Key: "triple_value_event",
				},
				query.LimitAndSort{
					Limit: query.CountLimit(10),
				},
				tripleValueEvent,
			)
			require.NoError(t, err)
			require.NotEmpty(t, sequences, "Expected to find TripleValueEvent")

			// Check the latest event
			latestEvent := sequences[0].Data.(*TripleValueEvent)
			require.NotEmpty(t, latestEvent.Values, "Expected non-empty values array")
			require.Equal(t, testBytes, latestEvent.Values[0])
		})
	})

	t.Run("LoopReader_GetLatestValue_GetTupleStruct", func(t *testing.T) {
		var retTupleStruct map[string]any
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, counterBinding.Name, "get_tuple_struct"}, "-"),
			primitives.Finalized,
			map[string]any{},
			&retTupleStruct,
		)
		require.NoError(t, err)
		require.NotEmpty(t, retTupleStruct, "Expected to find TupleStruct")
		log.Debugw("retTupleStruct", "retTupleStruct", retTupleStruct)

		// require.Equal(t, uint64(42), retTupleStruct["value"], "Expected value to be 42")
		// require.Equal(t, "0x1", retTupleStruct["address"], "Expected address to be 0x1")
		// require.Equal(t, true, retTupleStruct["bool"], "Expected bool to be true")
		// require.Equal(t, "0x1", retTupleStruct["struct_tag"], "Expected struct_tag to be 0x1")
	})

	t.Run("LoopReader_GetLatestValue_GetOCRConfig", func(t *testing.T) {
		type ConfigInfo struct {
			ConfigDigest                   []byte `json:"config_digest"`
			BigF                           uint64 `json:"big_f"`
			N                              uint64 `json:"n"`
			IsSignatureVerificationEnabled bool   `json:"is_signature_verification_enabled"`
		}

		type OCRConfig struct {
			ConfigInfo   ConfigInfo `json:"config_info"`
			Signers      [][]byte   `json:"signers"`
			Transmitters [][]byte   `json:"transmitters"`
		}

		type OCRConfigWrapped struct {
			OCRConfig OCRConfig `json:"OCRConfig"`
		}

		var retOCRConfig OCRConfigWrapped
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, counterBinding.Name, "get_ocr_config"}, "-"),
			primitives.Finalized,
			map[string]any{},
			&retOCRConfig,
		)

		require.NoError(t, err)
		require.NotEmpty(t, retOCRConfig, "Expected to find OCRConfig")
		log.Debugw("retOCRConfig", "retOCRConfig", retOCRConfig)
	})
}
