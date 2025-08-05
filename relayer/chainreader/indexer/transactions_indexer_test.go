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

	txnSigner := keystoreInstance.GetSuiSigner(context.Background(), hex.EncodeToString(publicKeyBytes))

	relayerClient, err := client.NewPTBClient(log, testutils.LocalUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, err)

	// Deploy and configure offramp
	DeployAndConfigureOffRamp(t, relayerClient, accountAddress, publicKeyBytes)

	// Deploy contract
	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, publishOutput, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	ccipObjectRefId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "offramp", "CCIPObjectRef")
	require.NoError(t, err)

	offrampStateObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "offramp", "OffRampState")
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

	readerConfig := config.ChainReaderConfig{
		Modules: map[string]*config.ChainReaderModule{
			"offramp": {
				Name:      "offramp",
				Functions: map[string]*config.ChainReaderFunction{},
				Events: map[string]*config.ChainReaderEvent{
					"offramp::ExecutionStateChanged": {
						Name:      "offramp",
						EventType: "ExecutionStateChanged",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "offramp",
							Event:   "ExecutionStateChanged",
						},
					},
					"offramp::SourceChainConfigSet": {
						Name:      "offramp",
						EventType: "SourceChainConfigSet",
						EventSelector: client.EventSelector{
							Package: packageId,
							Module:  "offramp",
							Event:   "SourceChainConfigSet",
						},
					},
				},
			},
			"ocr3_base": {
				Name:      "ocr3_base",
				Functions: map[string]*config.ChainReaderFunction{},
				Events: map[string]*config.ChainReaderEvent{
					"ocr3_base::ConfigSet": {
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
					"counter::CounterIncremented": {
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
		IsLoopPlugin: false,
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
			Name:    "ocr3_base",
			Address: packageId,
		},
		{
			Name:    "offramp",
			Address: packageId,
		},
	}

	err = cReader.Bind(ctx, boundContracts)
	require.NoError(t, err)

	t.Run("TestBasicFailedTransactionIndexing", func(t *testing.T) {
		ctx := context.Background()

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
		}()

		// 4. Create a successful transaction to trigger the transactions indexer
		CreateSuccessfulTransaction(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		time.Sleep(15 * time.Second)

		// 5. Create the initial OCR event to initiate transaction indexing
		setConfigResponse, setConfigErr := SetOCRConfig(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		require.NoError(t, setConfigErr)
		testutils.PrettyPrintDebug(log, setConfigResponse, "setConfigResponse")

		// 5. Create a failed PTB transaction
		reportStr := "0xd91ad9c94fba41de8869e580deb6dbc08e84fb41431d41d04f8849ed00be4a070dca7c34e2f78ecdd91ad9c94fba41de15a9c133ee53500a0300000000000000000000000000000014e30b40bfb1baeed9e4c62f145be85eb3d19ae932184920616d206120746573742063636970206d6573736167654010af5717948371a0b649a59530f8e80e0e1247e015f05f1f3e09c715288dd040420f00000000000000000000000000000000000000000000000000000000000114bd10ffa3815c010d5cf7d38815a0eaabc959eb84a1b6cf2e878987deb2624f9a122297abf6332d45b48c4df6fc3ea705f810980fa08601002000000000000000000000000000000000000000000000000000000000000000120000c16ff2862300000000000000000000000000000000000000000000000000010000"
		reportBytes, err := hex.DecodeString(strings.TrimPrefix(reportStr, "0x"))
		require.NoError(t, err)

		ptb := cwPTB.NewPTBConstructor(chainWriterConfig, relayerClient, log)
		ptbTx, err := ptb.BuildPTBCommands(context.Background(), "counter", "offramp_execution_with_error", cwConfig.Arguments{
			Args: map[string]any{
				"ref":            ccipObjectRefId,
				"state":          offrampStateObjectId,
				"clock":          "0x06",
				"report_context": [][]byte{},
				"report":         reportBytes,
			},
		}, nil)
		require.NoError(t, err)

		// Execute the PTB command using the PTB client, we don't check errors because we expect a failure
		_, _ = relayerClient.FinishPTBAndSend(ctx, txnSigner, ptbTx, client.WaitForLocalExecution)

		// 5. Use eventually to check if the correct synth events have been inserted
		require.Eventually(t, func() bool {
			eventHandle := fmt.Sprintf("%s::%s", boundContracts[0].Name, "ConfigSet")
			events, err := cReader.QueryKey(ctx, boundContracts[0], query.KeyFilter{Key: eventHandle}, query.LimitAndSort{}, &database.EventRecord{})
			if err != nil {
				log.Errorw("Error querying ConfigSet events", "error", err)
				return false
			}
			testutils.PrettyPrintDebug(log, events, "ConfigSet events in DB")

			return len(events) > 0
		}, 120*time.Second, 10*time.Second)

		require.Eventually(t, func() bool {
			eventHandle := fmt.Sprintf("%s::%s", boundContracts[1].Name, "SourceChainConfigSet")
			events, err := cReader.QueryKey(ctx, boundContracts[1], query.KeyFilter{Key: eventHandle}, query.LimitAndSort{}, &database.EventRecord{})
			if err != nil {
				log.Errorw("Error querying SourceChainConfigSet events", "error", err)
				return false
			}

			return len(events) > 0
		}, 120*time.Second, 10*time.Second)

		// 5. Use eventually to check if the correct synth events have been inserted
		require.Eventually(t, func() bool {
			eventHandle := fmt.Sprintf("%s::%s", boundContracts[1].Name, "ExecutionStateChanged")
			events, err := cReader.QueryKey(ctx, boundContracts[1], query.KeyFilter{Key: eventHandle}, query.LimitAndSort{}, &database.EventRecord{})
			if err != nil {
				log.Errorw("Error querying ExecutionStateChanged events", "error", err)
				return false
			}

			return len(events) > 0
		}, 120*time.Second, 10*time.Second)
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

func DeployAndConfigureOffRamp(t *testing.T, relayerClient *client.PTBClient, accountAddress string, signerPublicKey []byte) {
	t.Helper()

	// Stretch: implement or re-use what's in the offramp integration test instead of sample contracts
}
