//go:build integration

package loop

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	chainreaderConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestLoopChainReaderLocal(t *testing.T) {
	log := logger.Test(t)

	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	testutils.CleanupTestContracts()

	// Ensure the process is killed when the test completes.
	t.Cleanup(func() {
		testutils.CleanupTestContracts()

		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
	})

	log.Debugw("Started Sui node")

	runLoopChainReaderEchoTest(t, log, testutils.LocalGrpcURL)
}

func runLoopChainReaderEchoTest(t *testing.T, log logger.Logger, rpcUrl string) {
	t.Helper()
	ctx := context.Background()

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	relayerClient, clientErr := client.NewPTBClient(log, client.PTBClientConfig{
		GrpcTarget:            rpcUrl,
		GrpcToken:             "test",
		TransactionTimeout:    10 * time.Second,
		MaxConcurrentRequests: 5,
		KeystoreService:       keystoreInstance,
		DefaultRequestType:    client.TransactionRequestType("WaitForLocalExecution"),
	})
	require.NoError(t, clientErr)

	faucetFundErr := testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
	require.NoError(t, faucetFundErr)

	chainID, err := testutils.GetChainIdentifier(rpcUrl)
	require.NoError(t, err)
	testutils.PatchEnvironmentTOML("contracts/test", "local", chainID)
	testutils.PatchEnvironmentTOML("contracts/test_secondary", "local", chainID)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	gasBudget := int(2000000000)
	packageId, tx, err := testutils.PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotNil(t, packageId)
	require.NotNil(t, tx)

	// Set up the base ChainReader with echo function configurations
	chainReaderConfigs := chainreaderConfig.ChainReaderConfig{
		IsLoopPlugin: true,
		Modules: map[string]*chainreaderConfig.ChainReaderModule{
			"echo": {
				Name: "echo",
				Functions: map[string]*chainreaderConfig.ChainReaderFunction{
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
				Events: map[string]*chainreaderConfig.ChainReaderEvent{
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
				Functions: map[string]*chainreaderConfig.ChainReaderFunction{
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
		EventsIndexer: config.EventsIndexerConfig{
			PollingInterval: 2 * time.Second,
			SyncTimeout:     60 * time.Second,
		},
		TransactionsIndexer: config.TransactionsIndexerConfig{
			PollingInterval: 2 * time.Second,
			SyncTimeout:     60 * time.Second,
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

	// Create the indexers
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		log,
		map[string]*config.ChainReaderEvent{},
	)
	evIndexer := indexer.NewEventIndexer(
		db,
		log,
		[]*client.EventSelector{},
	)

	chainPoller := indexer.NewChainPoller(
		relayerClient,
		log,
		config.ChainPollerConfig{
			PollingInterval:         1 * time.Second,
			SyncTimeout:             60 * time.Second,
			BackfillCheckpointCount: testutils.Uint64Pointer(uint64(1)),
		},
		evIndexer.GetEventSelectors,
	)

	indexerInstance := indexer.NewIndexerFromComponents(
		log,
		chainPoller,
		evIndexer,
		txnIndexer,
	)

	chainReader, err := reader.NewChainReader(ctx, log, relayerClient, chainReaderConfigs, db, indexerInstance, nil)
	require.NoError(t, err)

	// Wrap the base chain reader with loop chain reader
	loopReader := NewLoopChainReader(log, chainReader)

	// Bind the contracts
	err = loopReader.Bind(context.Background(), []types.BoundContract{echoBinding, counterBinding})
	require.NoError(t, err)

	// Start the indexers
	err = indexerInstance.Start(ctx)
	require.NoError(t, err)

	log.Debugw("LoopChainReader setup complete")

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

	t.Run("LoopReader_GetLatestValue_EchoU256", func(t *testing.T) {
		t.Skip("Skipping, the entire test suite will be removed in favor of DefaultAccessor")
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

	t.Run("LoopReader_GetLatestValue_EchoU256_LargeValue", func(t *testing.T) {
		t.Skip("Skipping, the entire test suite will be removed in favor of DefaultAccessor")
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

	t.Run("LoopReader_GetLatestValue_EchoString", func(t *testing.T) {
		t.Skip("Skipping, the entire test suite will be removed in favor of DefaultAccessor")
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

	t.Run("LoopReader_EchoWithEvents_AndQueryEvents", func(t *testing.T) {
		// Test data
		testNumber := uint64(12345)

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

		// Query for SingleValueEvent
		t.Run("QuerySingleValueEvent", func(t *testing.T) {
			singleValueEvent := &SingleValueEvent{}
			var sequences []types.Sequence
			//nolint:govet
			var err error

			evIndexer.AddEventSelector(ctx, &client.EventSelector{
				Package: packageId,
				Module:  "echo",
				Event:   "SingleValueEvent",
			})

			// Use relayerClient to call increment instead of using CLI
			moveCallReq := client.MoveCallRequest{
				Signer:          accountAddress,
				PackageObjectId: packageId,
				Module:          "echo",
				Function:        "simple_event_echo",
				TypeArguments:   []any{},
				Arguments: []any{
					uint64(testNumber),
				},
				GasBudget: 2000000,
			}

			log.Debugw("Calling moveCall", "moveCallReq", moveCallReq)

			txMetadata, err := relayerClient.MoveCall(ctx, moveCallReq)
			require.NoError(t, err)

			txResponse, err := relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes)
			require.NoError(t, err)

			// Make sure the transaction succeeded to make sure the event is indexed
			require.Eventually(t, func() bool {
				status, err := relayerClient.GetTransactionStatus(ctx, txResponse.Transaction.GetDigest())
				log.Debugw("Transaction status for SingleValueEvent", "status", status, "error", err)
				return err == nil && status.Status == "success"
			}, 30*time.Second, 1*time.Second)

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
			}, 60*time.Second, 1*time.Second)

			require.NoError(t, err)
			require.NotEmpty(t, sequences, "Expected to find SingleValueEvent")
			log.Debugw("Sequences found", "sequences", sequences)
		})
	})

	t.Run("LoopReader_GetLatestValue_GetTupleStruct", func(t *testing.T) {
		t.Skip("Skipping, the entire test suite will be removed in favor of DefaultAccessor")
		var retTupleStruct map[string]any
		err = loopReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{packageId, counterBinding.Name, "get_tuple_struct"}, "-"),
			primitives.Finalized,
			map[string]any{},
			&retTupleStruct,
		)
		require.NoError(t, err)

		log.Debugw("retTupleStruct", "retTupleStruct", retTupleStruct)

		require.NotEmpty(t, retTupleStruct, "Expected to find TupleStruct")
		// Accept either float64 (from generic JSON map) or uint64
		if v, ok := retTupleStruct["value"].(float64); ok {
			require.Equal(t, float64(42), v, "Expected value to be 42")
		} else {
			require.Equal(t, uint64(42), retTupleStruct["value"], "Expected value to be 42")
		}
		require.Equal(t, "0x1", retTupleStruct["address"], "Expected address to be 0x1")
		require.Equal(t, true, retTupleStruct["bool"], "Expected bool to be true")
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

	t.Run("BatchGetLatestValues", func(t *testing.T) {
		type TupleResult struct {
			First  uint32 `json:"first"`
			Second uint64 `json:"second"`
		}
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

		var (
			echoU64Values = make([]uint64, 6)
			echoU64Ptrs   = make([]*uint64, len(echoU64Values))
			tupleResults  = make([]TupleResult, 2)
			tuplePtrs     = make([]*TupleResult, len(tupleResults))
			ocrConfigs    = make([]OCRConfigWrapped, 4)
			ocrConfigPtrs = make([]*OCRConfigWrapped, len(ocrConfigs))
		)
		for i := range echoU64Values {
			echoU64Values[i] = uint64((i + 1) * 10)
			echoU64Ptrs[i] = &echoU64Values[i]
		}
		for i := range tupleResults {
			tupleResults[i] = TupleResult{First: uint32(100 + i), Second: uint64(200 + i)}
			tuplePtrs[i] = &tupleResults[i]
		}
		for i := range ocrConfigs {
			ocrConfigPtrs[i] = &ocrConfigs[i]
		}

		request := types.BatchGetLatestValuesRequest{
			echoBinding: {
				{ReadName: "echo_u64", Params: map[string]any{"val": echoU64Values[0]}, ReturnVal: echoU64Ptrs[0]},
				{ReadName: "echo_u64", Params: map[string]any{"val": echoU64Values[1]}, ReturnVal: echoU64Ptrs[1]},
				{ReadName: "echo_u64", Params: map[string]any{"val": echoU64Values[2]}, ReturnVal: echoU64Ptrs[2]},
				{ReadName: "echo_u64", Params: map[string]any{"val": echoU64Values[3]}, ReturnVal: echoU64Ptrs[3]},
				{ReadName: "echo_u64", Params: map[string]any{"val": echoU64Values[4]}, ReturnVal: echoU64Ptrs[4]},
				{ReadName: "echo_u64", Params: map[string]any{"val": echoU64Values[5]}, ReturnVal: echoU64Ptrs[5]},
				{ReadName: "echo_u32_u64_tuple", Params: map[string]any{"val1": tupleResults[0].First, "val2": tupleResults[0].Second}, ReturnVal: tuplePtrs[0]},
				{ReadName: "echo_u32_u64_tuple", Params: map[string]any{"val1": tupleResults[1].First, "val2": tupleResults[1].Second}, ReturnVal: tuplePtrs[1]},
			},
			counterBinding: {
				{ReadName: "get_ocr_config", Params: map[string]any{}, ReturnVal: ocrConfigPtrs[0]},
				{ReadName: "get_ocr_config", Params: map[string]any{}, ReturnVal: ocrConfigPtrs[1]},
				{ReadName: "get_ocr_config", Params: map[string]any{}, ReturnVal: ocrConfigPtrs[2]},
				{ReadName: "get_ocr_config", Params: map[string]any{}, ReturnVal: ocrConfigPtrs[3]},
			},
		}

		totalReads := 0
		for _, batch := range request {
			totalReads += len(batch)
		}
		require.Equal(t, 12, totalReads, "expected a batch of 12 reads")

		result, err := loopReader.BatchGetLatestValues(ctx, request)
		require.NoError(t, err)
		require.Len(t, result, 2, "expected results for echo and counter contracts")

		echoResults := result[echoBinding]
		require.Len(t, echoResults, 8)
		for i, readResult := range echoResults {
			_, readErr := readResult.GetResult()
			require.NoError(t, readErr, "echo batch read %d failed", i)
		}

		counterResults := result[counterBinding]
		require.Len(t, counterResults, 4)
		for i, readResult := range counterResults {
			_, readErr := readResult.GetResult()
			require.NoError(t, readErr, "counter batch read %d failed", i)
		}

		for i, expected := range echoU64Values {
			require.Equal(t, expected, echoU64Values[i], "echo_u64 read %d", i)
		}
		for i, expected := range tupleResults {
			require.Equal(t, expected.First, tupleResults[i].First, "echo_u32_u64_tuple read %d first", i)
			require.Equal(t, expected.Second, tupleResults[i].Second, "echo_u32_u64_tuple read %d second", i)
		}
		require.NotEmpty(t, ocrConfigs[0], "expected OCR config from counter contract")
	})
}
