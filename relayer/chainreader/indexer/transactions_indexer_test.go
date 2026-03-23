//go:build integration

package indexer_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	cwPTB "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestTransactionsIndexer(t *testing.T) {
	ctx := context.Background()
	log := logger.Test(t)
	testutils.CleanupTestContracts()

	// Setup database
	datastoreUrl := os.Getenv("TEST_DB_URL")
	if datastoreUrl == "" {
		t.Skip("Skipping persistent tests as TEST_DB_URL is not set in CI")
	}
	db := sqltest.NewDB(t, datastoreUrl)

	// Verify database connection
	dbConnection, err := db.Connx(ctx)
	require.NoError(t, err)

	dbStore := database.NewDBStore(db, log)
	require.NoError(t, dbStore.EnsureSchema(ctx))

	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)

	t.Cleanup(func() {
		testutils.CleanupTestContracts()
		if cmd.Process != nil {
			perr := cmd.Process.Kill()
			if perr != nil {
				t.Logf("Failed to kill process: %v", perr)
			}
		}
		dbConnection.Close()
	})

	log.Debugw("Started Sui node")

	// Setup keystore and client
	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	// Fund the account multiple times to ensure sufficient balance
	require.Eventually(t, func() bool {
		failed := false

		for i := 0; i < 3; i++ {
			err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
			if err != nil {
				failed = true
				break
			}
		}

		return !failed
	}, 15*time.Second, 1*time.Second, "Failed to fund account with sufficient SUI balance")

	txnSigner := keystoreInstance.GetSuiSigner(ctx, hex.EncodeToString(publicKeyBytes))

	relayerClient, err := client.NewPTBClient(log, testutils.LocalUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, err)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	gasBudget := int(2000000000)
	packageId, tx, err := testutils.PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotEmpty(t, packageId)
	require.NotEmpty(t, tx)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(tx.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	ccipObjectRefId, err := testutils.QueryCreatedObjectID(tx.ObjectChanges, packageId, "offramp", "CCIPObjectRef")
	require.NoError(t, err)

	offrampStateObjectId, err := testutils.QueryCreatedObjectID(tx.ObjectChanges, packageId, "offramp", "OffRampState")
	require.NoError(t, err)

	chainWriterConfig := cwConfig.ChainWriterConfig{
		Modules: map[string]*cwConfig.ChainWriterModule{
			"counter": {
				Name:     "counter",
				ModuleID: packageId,
				Functions: map[string]*cwConfig.ChainWriterFunction{
					"increment_by": {
						Name:      "increment_by",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []cwConfig.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  testutils.StringPointer("counter"),
								Function:  testutils.StringPointer("increment_by"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "by",
										Type:     "u64",
										Required: true,
									},
								},
							},
						},
					},
					"increment_by_bytes_length": {
						Name:      "increment_by_bytes_length",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []cwConfig.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  testutils.StringPointer("counter"),
								Function:  testutils.StringPointer("increment_by_bytes_length"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "counter",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "bytes",
										Type:     "vector<u8>",
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"offramp": {
				Name:     "offramp",
				ModuleID: packageId,
				Functions: map[string]*cwConfig.ChainWriterFunction{
					"offramp_execution_with_error": {
						Name:      "offramp_execution_with_error",
						PublicKey: publicKeyBytes,
						Params:    []codec.SuiFunctionParam{},
						PTBCommands: []cwConfig.ChainWriterPTBCommand{
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  testutils.StringPointer("offramp"),
								Function:  testutils.StringPointer("init_execute"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "ref",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:      "clock",
										Type:      "object_id",
										Required:  true,
										IsMutable: testutils.BoolPointer(false),
									},
									{
										Name:     "report_context",
										Type:     "vector<vector<u8>>",
										Required: true,
									},
									{
										Name:     "report",
										Type:     "vector<u8>",
										Required: true,
									},
								},
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: &packageId,
								ModuleId:  testutils.StringPointer("offramp"),
								Function:  testutils.StringPointer("finish_execute"),
								Params:    []codec.SuiFunctionParam{},
							},
						},
					},
				},
			},
		},
	}

	// Create transactions indexer
	pollingInterval := 4 * time.Second
	syncTimeout := 3 * time.Second

	type OfframpExecutionStateChanged struct {
		SourceChainSelector uint64 `json:"sourceChainSelector"`
		SequenceNumber      uint64 `json:"sequenceNumber"`
		MessageId           string `json:"messageId"`
		MessageHash         string `json:"messageHash"`
		State               int    `json:"state"`
	}

	readerConfig := config.ChainReaderConfig{
		IsLoopPlugin: false,
		Modules: map[string]*config.ChainReaderModule{
			"OffRamp": {
				Name:      "offramp",
				Functions: map[string]*config.ChainReaderFunction{},
				Events: map[string]*config.ChainReaderEvent{
					"ExecutionStateChanged": {
						Name:      "offramp",
						EventType: "ExecutionStateChanged",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "offramp",
							Event:   "ExecutionStateChanged",
						},
						ExpectedEventType: &OfframpExecutionStateChanged{},
					},
					"SourceChainConfigSet": {
						Name:      "offramp",
						EventType: "SourceChainConfigSet",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "offramp",
							Event:   "SourceChainConfigSet",
						},
					},
					"ConfigSet": {
						Name:      "ocr3_base",
						EventType: "ConfigSet",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "ocr3_base",
							Event:   "ConfigSet",
						},
					},
				},
			},
			"counter": {
				Name:      "counter",
				Functions: map[string]*config.ChainReaderFunction{},
				Events: map[string]*config.ChainReaderEvent{
					"CounterIncremented": {
						Name:      "counter",
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
		EventsIndexer: config.EventsIndexerConfig{
			PollingInterval: pollingInterval,
			SyncTimeout:     syncTimeout,
		},
		TransactionsIndexer: config.TransactionsIndexerConfig{
			PollingInterval: pollingInterval,
			SyncTimeout:     syncTimeout,
		},
	}

	// Create the indexers
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		log,
		relayerClient,
		readerConfig.TransactionsIndexer.PollingInterval,
		readerConfig.TransactionsIndexer.SyncTimeout,
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
		readerConfig.EventsIndexer.PollingInterval,
		readerConfig.EventsIndexer.SyncTimeout,
	)
	indexerInstance := indexer.NewIndexer(
		log,
		evIndexer,
		txnIndexer,
	)

	// Create ChainReader (remove the schema creation comment since it's already done)
	cReader, err := reader.NewChainReader(
		ctx,
		log,
		relayerClient,
		readerConfig,
		db,
		indexerInstance,
	)
	require.NoError(t, err)

	err = indexerInstance.Start(ctx)
	require.NoError(t, err)

	// Clean the events table again with a temporary connection
	func() {
		dbConn, dbConnErr := db.Connx(ctx)
		require.NoError(t, dbConnErr)
		defer dbConn.Close() // Explicitly close this connection

		_, deleteEventsErr := dbConn.ExecContext(ctx, `DELETE FROM sui.events WHERE TRUE`)
		require.NoError(t, deleteEventsErr)
	}()

	boundContracts := []types.BoundContract{
		{
			Name:    "OffRamp",
			Address: packageId,
		},
	}

	err = cReader.Bind(ctx, boundContracts)
	require.NoError(t, err)

	t.Run("TestBasicFailedTransactionIndexing", func(t *testing.T) {
		ctx := context.Background()

		// helper: returns true if at least one event with the given key exists for the contract
		hasEvent := func(contract types.BoundContract, key string) bool {
			dataType := map[string]any{}
			events, err := cReader.QueryKey(ctx, contract, query.KeyFilter{Key: key}, query.LimitAndSort{}, &dataType)
			if err != nil {
				log.Errorw("Error querying events", "contract", contract.Name, "key", key, "error", err)
				return false
			}
			found := len(events) > 0

			if found {
				log.Debugw("Event found (hasEvent)", "events", events)
			} else {
				log.Debugw("Event not found", events)
			}

			return found
		}

		// helper: same as hasEvent but only checks the database for the event without using the ChainReader
		hasEventDBOnlyCheck := func(packageId string, module string, key string) bool {
			events, err := dbStore.QueryEvents(ctx, packageId, fmt.Sprintf("%s::%s::%s", packageId, module, key), []query.Expression{}, query.LimitAndSort{})
			if err != nil {
				log.Errorw("Error querying events", "packageId", packageId, "module", module, "key", key, "error", err)
				return false
			}
			found := len(events) > 0

			if found {
				log.Debugw("Event found (hasEventDBOnlyCheck)", "events", events)
			} else {
				log.Debugw("Event not found", events)
			}

			return found
		}

		// 1. Create a few transactions
		for range 3 {
			CreateFailedTransaction(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		}

		// 2. Query the transactions and ensure that they are findable from the RPC
		txs_1, err := relayerClient.QueryTransactions(ctx, accountAddress, nil, nil)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(txs_1.Data), 3, "Expected at least 3 transactions")

		// 3. Start the indexers and ensure that the events / transactions are indexed
		go func() {
			_ = cReader.Start(ctx)
			_ = txnIndexer.Start(ctx)
		}()

		// 4. Create a successful transaction to trigger the transactions indexer
		CreateSuccessfulTransaction(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		time.Sleep(15 * time.Second)

		// 5. Create the initial OCR event to initiate transaction indexing
		setConfigResponse, setConfigErr := SetOCRConfig(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		require.NoError(t, setConfigErr)
		testutils.PrettyPrintDebug(log, setConfigResponse, "setConfigResponse")

		// 4.a. Wait for the configs to be set
		require.Eventually(t, func() bool {
			okConfig := hasEventDBOnlyCheck(packageId, "ocr3_base", "ConfigSet")
			okSrcCfg := hasEventDBOnlyCheck(packageId, "offramp", "SourceChainConfigSet")

			log.Debugw("event wait progress",
				"ConfigSet", okConfig,
				"SourceChainConfigSet", okSrcCfg,
			)

			return okConfig && okSrcCfg
		}, 90*time.Second, 5*time.Second)

		// 5. Create a failed PTB transaction
		reportStr := "9b3c1f221aa3f0cc579b9518768ead0a57cc3d9d782049b702fab91dd723c757f287d20217d8e69b9b3c1f221aa3f0ccec1182faa7c27b87a40200000000000000000000000000001407775923481a094e41d51449b0b0f979c126a3b003486579b4dcbf61d5f5f447ae448e3c1503a811d83bdc074a8712ebeb241fd649b372e040420f00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
		reportBytes, err := hex.DecodeString(reportStr)
		require.NoError(t, err)

		ptb := cwPTB.NewPTBConstructor(chainWriterConfig, relayerClient, log)
		ptbTx, err := ptb.BuildPTBCommands(context.Background(), "offramp", "offramp_execution_with_error", cwConfig.Arguments{
			Args: map[string]any{
				"ref":            ccipObjectRefId,
				"state":          offrampStateObjectId,
				"clock":          "0x06",
				"report_context": [][]byte{},
				"report":         reportBytes,
			},
		}, packageId, chainWriterConfig.Modules["offramp"].Functions["offramp_execution_with_error"])
		require.NoError(t, err)

		response, _ := relayerClient.FinishPTBAndSend(ctx, txnSigner, ptbTx, client.WaitForLocalExecution)
		require.Equal(t, "failure", response.Status.Status)

		// 5.b. Wait for the execution state changed event to be indexed
		require.Eventually(t, func() bool {
			return hasEvent(boundContracts[0], "ExecutionStateChanged")
		}, 90*time.Second, 5*time.Second)

		events, err := cReader.QueryKey(ctx, boundContracts[0], query.KeyFilter{Key: "ExecutionStateChanged"}, query.LimitAndSort{}, &OfframpExecutionStateChanged{})
		require.NoError(t, err)
		require.NotEmpty(t, events)

		executionStateChanged := events[0].Data.(*OfframpExecutionStateChanged)

		require.True(t, strings.HasPrefix(executionStateChanged.MessageId, "0x"))
		require.True(t, strings.HasPrefix(executionStateChanged.MessageHash, "0x"))
	})
}

func CreateFailedTransaction(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) {
	t.Helper()
	// Verify we can execute the transaction
	resp, err := BasicIncrementBy(t, relayerClient, packageId, counterObjectId, accountAddress, signerPublicKey, "1000")
	require.NoError(t, err)
	require.Equal(t, "failure", resp.Status.Status, "Expected move call to fail")
}

func CreateSuccessfulTransaction(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) {
	t.Helper()
	// Verify we can execute the transaction
	resp, err := BasicIncrementBy(t, relayerClient, packageId, counterObjectId, accountAddress, signerPublicKey, "10")
	require.NoError(t, err)
	require.Equal(t, "success", resp.Status.Status, "Expected move call to succeed")
}

func BasicIncrementBy(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte, val string) (client.SuiTransactionBlockResponse, error) {
	t.Helper()
	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "counter",
		Function:        "increment_by",
		Arguments:       []any{counterObjectId, val},
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

	return resp, err
}

func SetOCRConfig(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) (client.SuiTransactionBlockResponse, error) {
	t.Helper()
	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "ocr3_base",
		Function:        "set_ocr3_config",
		Arguments:       []any{[]byte{1, 2, 3, 4, 5}, uint8(0), uint8(1), [][]byte{signerPublicKey}, []string{accountAddress}},
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

	return resp, err
}
