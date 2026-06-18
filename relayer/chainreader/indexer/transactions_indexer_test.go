//go:build integration

package indexer_test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	v2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"

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

	log.Debugw("Started Sui node")

	// Setup keystore and client
	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	// Fund the account multiple times to ensure sufficient balance
	require.Eventually(t, func() bool {
		failed := false

		for i := 0; i < 5; i++ {
			err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
			if err != nil {
				failed = true
				break
			}
		}

		return !failed
	}, 15*time.Second, 1*time.Second, "Failed to fund account with sufficient SUI balance")

	txnSigner := keystoreInstance.GetSuiSigner(ctx, hex.EncodeToString(publicKeyBytes))

	ptbClientConfig := client.PTBClientConfig{
		GrpcTarget:            testutils.LocalGrpcURL,
		GrpcToken:             "test",
		TransactionTimeout:    10 * time.Second,
		MaxConcurrentRequests: 5,
		KeystoreService:       keystoreInstance,
		DefaultRequestType:    client.WaitForLocalExecution,
	}
	relayerClient, err := client.NewPTBClient(log, ptbClientConfig)
	require.NoError(t, err)

	chainID, chainIDErr := testutils.GetChainIdentifier(testutils.LocalURL)
	require.NoError(t, chainIDErr)
	testutils.PatchEnvironmentTOML("contracts/test", "local", chainID)
	testutils.PatchEnvironmentTOML("contracts/test_secondary", "local", chainID)

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

	publishCheckpoint, err := relayerClient.GetLatestCheckpoint(ctx)
	require.NoError(t, err)
	publishCheckpointSeq := publishCheckpoint.GetSequenceNumber()

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
	}

	// Create the indexers with new channel-based API
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		log,
		// start without any configs, they will be set when ChainReader is initialized
		map[string]*config.ChainReaderEvent{},
	)

	evIndexer := indexer.NewEventIndexer(
		db,
		log,
		// start without any selectors, they will be added during .Bind() calls on ChainReader
		[]*client.EventSelector{
			{
				Package: packageId,
				Module:  "ocr3_base",
				Event:   "ConfigSet",
			},
			{
				Package: packageId,
				Module:  "offramp",
				Event:   "SourceChainConfigSet",
			},
		},
	)

	chainPoller := indexer.NewChainPoller(
		relayerClient,
		log,
		config.ChainPollerConfig{
			PollingInterval:         2 * time.Second,
			SyncTimeout:             60 * time.Second,
			ChannelBufferSize:       16,
			StartCheckpointSequence: &publishCheckpointSeq,
		},
		evIndexer.GetEventSelectors,
	)

	indexerInstance := indexer.NewIndexerFromComponents(
		log,
		chainPoller,
		evIndexer,
		txnIndexer,
	)

	// Create ChainReader
	cReader, err := reader.NewChainReader(
		ctx,
		log,
		relayerClient,
		readerConfig,
		db,
		indexerInstance,
		nil,
	)
	require.NoError(t, err)

	boundContracts := []types.BoundContract{
		{
			Name:    "OffRamp",
			Address: packageId,
		},
	}

	err = cReader.Bind(ctx, boundContracts)
	require.NoError(t, err)

	err = indexerInstance.Start(ctx)
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
		indexerInstance.Close()
		cReader.Close()
	})

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
				log.Debugw("Event not found")
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
				log.Debugw("Event not found")
			}

			return found
		}

		// 1. Create a few transactions
		for range 3 {
			CreateFailedTransaction(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)
		}

		// 2. Query the transactions and ensure that they are findable from the RPC (using checkpoint-based API)
		// Note: QueryTransactions is deprecated, but we can verify via checkpoint data
		latestCheckpoint, err := relayerClient.GetLatestCheckpoint(ctx)
		require.NoError(t, err)
		require.NotNil(t, latestCheckpoint)
		log.Debugw("Latest checkpoint found", "sequence", latestCheckpoint.GetSequenceNumber())

		// 3. Create a successful transaction (exercises the checkpoint pipeline)
		CreateSuccessfulTransaction(t, relayerClient, packageId, counterObjectId, accountAddress, publicKeyBytes)

		// 4. Create the initial OCR event to initiate transaction indexing
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

		ptbTx := BuildFailedOfframpExecutionPTB(ctx, t, relayerClient, packageId, accountAddress, ccipObjectRefId, offrampStateObjectId, reportBytes)
		ptbTx.SetSigner(txnSigner)

		bcsBytes, err := ptbTx.BuildBCSBytes(ctx)
		require.NoError(t, err)

		response, err := relayerClient.SignAndSendTransaction(ctx, base64.StdEncoding.EncodeToString(bcsBytes), publicKeyBytes)
		require.NoError(t, err)
		require.False(t, response.Transaction.GetEffects().GetStatus().GetSuccess())

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

	// Cleanup
	err = indexerInstance.Close()
	require.NoError(t, err)
}

func BuildFailedOfframpExecutionPTB(
	ctx context.Context,
	t *testing.T,
	relayerClient *client.PTBClient,
	packageId string,
	accountAddress string,
	ccipObjectRefId string,
	offrampStateObjectId string,
	reportBytes []byte,
) *transaction.Transaction {
	t.Helper()

	txn := transaction.NewTransaction()
	txn.SetSender(models.SuiAddress(accountAddress))

	refArg, err := relayerClient.TransformTransactionArg(ctx, txn, ccipObjectRefId, "object_id", false)
	require.NoError(t, err)

	stateArg, err := relayerClient.TransformTransactionArg(ctx, txn, offrampStateObjectId, "object_id", true)
	require.NoError(t, err)

	clockArg, err := relayerClient.TransformTransactionArg(ctx, txn, "0x6", "object_id", false)
	require.NoError(t, err)

	reportContextArg := txn.Pure([][]byte{})
	reportArg := txn.Pure(reportBytes)

	txn.MoveCall(
		models.SuiAddress(packageId),
		"offramp",
		"init_execute",
		nil,
		[]transaction.Argument{*refArg, *stateArg, *clockArg, reportContextArg, reportArg},
	)
	txn.MoveCall(
		models.SuiAddress(packageId),
		"offramp",
		"finish_execute",
		nil,
		nil,
	)

	referenceGasPrice, err := relayerClient.GetReferenceGasPrice(ctx)
	require.NoError(t, err)
	txn.SetGasPrice(referenceGasPrice.Uint64())
	txn.SetGasBudget(client.DefaultGasBudget)

	paymentCoinBytes, paymentCoinVersion, paymentCoinDigest, err := relayerClient.GetTransactionPaymentCoinForAddress(ctx, accountAddress)
	require.NoError(t, err)
	txn.SetGasPayment([]transaction.SuiObjectRef{
		{
			ObjectId: paymentCoinBytes,
			Version:  paymentCoinVersion,
			Digest:   paymentCoinDigest,
		},
	})

	return txn
}

func CreateFailedTransaction(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) {
	t.Helper()
	_, err := BasicIncrementBy(t, relayerClient, packageId, counterObjectId, accountAddress, signerPublicKey, 1000, false)
	require.NoError(t, err)
}

func CreateSuccessfulTransaction(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) {
	t.Helper()
	_, err := BasicIncrementBy(t, relayerClient, packageId, counterObjectId, accountAddress, signerPublicKey, 10, true)
	require.NoError(t, err)
}

func BasicIncrementBy(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte, val uint64, expectSuccess bool) (*v2.ExecuteTransactionResponse, error) {
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
	)
	require.NoError(t, err)

	// SignAndSendTransaction may return before effects are finalized; poll for the
	// terminal on-chain status and assert against that instead of the initial response.
	require.Eventually(t, func() bool {
		status, err := relayerClient.GetTransactionStatus(context.Background(), resp.Transaction.GetDigest())
		if err != nil {
			return false
		}
		if expectSuccess {
			return status.Status == "success"
		}
		return status.Status == "failure"
	}, 10*time.Second, 200*time.Millisecond)

	return resp, err
}

func SetOCRConfig(t *testing.T, relayerClient *client.PTBClient, packageId string, counterObjectId string, accountAddress string, signerPublicKey []byte) (*v2.ExecuteTransactionResponse, error) {
	t.Helper()

	// Prepare arguments for a move call
	moveCallReq := client.MoveCallRequest{
		Signer:          accountAddress,
		PackageObjectId: packageId,
		Module:          "ocr3_base",
		Function:        "set_ocr3_config",
		Arguments: []any{
			[]byte{1, 2, 3, 4, 5},
			uint8(0),
			uint8(1),
			[][]byte{signerPublicKey},
			[]string{accountAddress},
		},
		GasBudget: 1000000000,
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
	)

	require.Eventually(t, func() bool {
		status, err := relayerClient.GetTransactionStatus(context.Background(), resp.Transaction.GetDigest())
		return err == nil && status.Status == "success"
	}, 10*time.Second, 200*time.Millisecond)

	return resp, err
}

func ptr[T any](v T) *T {
	return &v
}
