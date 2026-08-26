//go:build integration

package indexer_test

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	v2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/codec"
	chainreaderutil "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/chainreader_util"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

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

		decodedReport, err := codec.DeserializeExecutionReport(reportBytes)
		require.NoError(t, err)
		require.NotNil(t, decodedReport)

		// The message ID is expected to be encoded as a hex string due to the use of `ExpectedEventType`
		// in the ChainReader config for the relevant event.
		require.Equal(t, "0x"+hex.EncodeToString(decodedReport.Message.Header.MessageID), executionStateChanged.MessageId)
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

/***
 * The following test attempts to achieve the same verification done
 * by the previous `TestTransactionIndexer` test while running against a
 * known checkpoint from testnet.
 *
 * It is expected that this test will eventually fail as either the testnet
 * RPC hardcoded here is removed or the checkpoint used for verification is
 * pruned.
 *
 * A second requirement of this version of the test is to avoid
 * adding aux. events (ex: SetOCRConfig) manually and allow those events
 * to be fetched from the contracts to resemble a production deployment.
 *
 * This test can be safely skipped if necessary.
 */
func TestTransactionsIndexerTestnet(t *testing.T) {
	grpcTarget := os.Getenv("GRPC_TARGET")
	if grpcTarget == "" {
		t.Fatal("environment value for GRPC_TARGET is not set")
	}
	grpcToken := os.Getenv("GRPC_TOKEN")
	if grpcToken == "" {
		t.Fatal("environment value for GRPC_TOKEN is not set")
	}
	// Known failurd transaction at checkpoint #370592649
	// Explorer URL: https://suiscan.xyz/testnet/tx/5syaidZp4YxG2tMRyDUKHAsNSTMiqRUVEq3EAKZo3JzY
	const CHECKPOINT = 370_000_000
	// The checkpoint the known failed transaction landed in. Polling starts at CHECKPOINT and
	// catches up through this one, which is also the window in which the OffRamp's ConfigSet and
	// SourceChainConfigSet events were emitted, so no config event has to be created by hand.
	const FAILED_TX_CHECKPOINT = 370592649
	const FAILED_TX_DIGEST = "5syaidZp4YxG2tMRyDUKHAsNSTMiqRUVEq3EAKZo3JzY"

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
	t.Cleanup(func() { dbConnection.Close() })

	dbStore := database.NewDBStore(db, log)
	require.NoError(t, dbStore.EnsureSchema(ctx))

	relayerClient, err := client.NewPTBClient(log, client.PTBClientConfig{
		GrpcTarget:            grpcTarget,
		GrpcToken:             grpcToken,
		TransactionTimeout:    30 * time.Second,
		MaxConcurrentRequests: 50,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		if cerr := relayerClient.Close(); cerr != nil {
			t.Logf("Failed to close client: %v", cerr)
		}
	})

	// 1. Read the known failed transaction back out of its checkpoint. The transmitter and the
	// execution report are taken straight from the on-chain transaction; the OffRamp package is the
	// pinned original package ID (see OFFRAMP_PACKAGE_ID).
	checkpointTxns, err := relayerClient.GetTransactionsByCheckpoint(ctx, FAILED_TX_CHECKPOINT)
	require.NoError(t, err, "failed to read checkpoint %d, it may have been pruned by the RPC", uint64(FAILED_TX_CHECKPOINT))

	var failedTxn *v2.ExecutedTransaction
	for _, txn := range checkpointTxns {
		if txn.GetDigest() == FAILED_TX_DIGEST {
			failedTxn = txn
			break
		}
	}
	require.NotNil(t, failedTxn, "transaction %s not found in checkpoint %d", FAILED_TX_DIGEST, uint64(FAILED_TX_CHECKPOINT))
	require.False(t, failedTxn.GetEffects().GetStatus().GetSuccess(), "expected %s to be a failed transaction", FAILED_TX_DIGEST)
	require.NotEmpty(t, failedTxn.GetTransaction().GetSender(), "expected the failed transaction to have a sender")

	offrampPackageId, executeCommandIndex, reportBytes := FindOfframpExecuteCall(t, failedTxn)

	// The indexer only records a synthetic FAILURE when the abort happened strictly after
	// init_execute, i.e. when the report itself passed full onchain validation. Assert the
	// transaction still has that shape, so a change to it is reported as such instead of
	// surfacing later as an unexplained missing event.
	execError := failedTxn.GetEffects().GetStatus().GetError()
	moveAbort := execError.GetAbort()
	require.NotNil(t, moveAbort, "expected the failed transaction to carry MoveAbort details")
	require.NotEmpty(t, moveAbort.GetLocation().GetFunctionName(), "expected the abort location to name a function")
	require.NotEqual(t, "init_execute", moveAbort.GetLocation().GetFunctionName())
	require.Greater(t, execError.GetCommand(), executeCommandIndex, "expected the abort to happen after init_execute")

	log.Infow("Found known failed OffRamp execution",
		"digest", FAILED_TX_DIGEST,
		"offrampPackageId", offrampPackageId,
		"transmitter", failedTxn.GetTransaction().GetSender(),
	)

	executionReport, err := codec.DeserializeExecutionReportFromPure(reportBytes)
	require.NoError(t, err)
	require.NotNil(t, executionReport)

	sourceChainSelector := executionReport.Message.Header.SourceChainSelector

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
							Package: offrampPackageId,
							Module:  "offramp",
							Event:   "ExecutionStateChanged",
						},
						ExpectedEventType: &OfframpExecutionStateChanged{},
					},
					"SourceChainConfigSet": {
						Name:      "offramp",
						EventType: "SourceChainConfigSet",
						EventSelector: client.EventSelector{
							Package: offrampPackageId,
							Module:  "offramp",
							Event:   "SourceChainConfigSet",
						},
					},
					"ConfigSet": {
						Name:      "ocr3_base",
						EventType: "ConfigSet",
						EventSelector: client.EventSelector{
							Package: offrampPackageId,
							Module:  "ocr3_base",
							Event:   "ConfigSet",
						},
					},
				},
			},
		},
	}

	// 2. Wire up the same components the localnet test uses, pointed at testnet and started from
	// the known checkpoint rather than the one the contracts were published in.
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		log,
		// start without any configs, they will be set when ChainReader is initialized
		map[string]*config.ChainReaderEvent{},
	)

	evIndexer := indexer.NewEventIndexer(
		db,
		log,
		// the OCR and source chain configs are indexed from the chain, exactly as they would be in
		// a production deployment, so their selectors have to be in place before the first poll
		[]*client.EventSelector{
			{
				Package: offrampPackageId,
				Module:  "ocr3_base",
				Event:   "ConfigSet",
			},
			{
				Package: offrampPackageId,
				Module:  "offramp",
				Event:   "SourceChainConfigSet",
			},
		},
	)

	startCheckpointSequence := uint64(CHECKPOINT)
	chainPoller := indexer.NewChainPoller(
		relayerClient,
		log,
		config.ChainPollerConfig{
			PollingInterval:         2 * time.Second,
			SyncTimeout:             60 * time.Second,
			ChannelBufferSize:       16,
			StartCheckpointSequence: &startCheckpointSequence,
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
			Address: offrampPackageId,
		},
	}

	err = cReader.Bind(ctx, boundContracts)
	require.NoError(t, err)

	err = indexerInstance.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		if cerr := indexerInstance.Close(); cerr != nil {
			t.Logf("Failed to close indexer: %v", cerr)
		}
		if cerr := cReader.Close(); cerr != nil {
			t.Logf("Failed to close chain reader: %v", cerr)
		}
	})

	// helper: checks the database directly for an event, without going through the ChainReader
	hasEventDBOnlyCheck := func(module string, key string) bool {
		events, qerr := dbStore.QueryEvents(ctx, offrampPackageId, fmt.Sprintf("%s::%s::%s", offrampPackageId, module, key), []query.Expression{}, query.LimitAndSort{})
		if qerr != nil {
			log.Errorw("Error querying events", "module", module, "key", key, "error", qerr)
			return false
		}

		return len(events) > 0
	}

	// 3. Wait for the OCR and source chain configs to be indexed from the chain. The transactions
	// indexer stays idle until it sees a ConfigSet event, and needs the SourceChainConfigSet to
	// hash the message, so both must land before any synthetic event can be written.
	require.Eventually(t, func() bool {
		okConfig := hasEventDBOnlyCheck("ocr3_base", "ConfigSet")
		okSrcCfg := hasEventDBOnlyCheck("offramp", "SourceChainConfigSet")

		log.Debugw("config event wait progress",
			"ConfigSet", okConfig,
			"SourceChainConfigSet", okSrcCfg,
		)

		return okConfig && okSrcCfg
	}, 10*time.Minute, 5*time.Second, "OffRamp config events were not indexed from checkpoints at or after %d", uint64(CHECKPOINT))

	// 4. Wait for the poller to reach the failed transaction's checkpoint and for the transactions
	// indexer to write the synthetic ExecutionStateChanged event for it. The polled window also
	// covers genuine executions by the same OffRamp, so the message ID is used to narrow the query
	// down to the message the known failed transaction carried. It is stored base64-encoded, which
	// is the representation both real and synthetic events are written in.
	messageIdFilter := []query.Expression{
		query.Comparator("messageId",
			primitives.ValueComparator{
				Value:    base64.StdEncoding.EncodeToString(executionReport.Message.Header.MessageID),
				Operator: primitives.Eq,
			},
		),
	}

	require.Eventually(t, func() bool {
		events, qerr := dbStore.QueryEvents(ctx, offrampPackageId, offrampPackageId+"::offramp::ExecutionStateChanged", messageIdFilter, query.LimitAndSort{})
		if qerr != nil {
			log.Errorw("Error querying execution state changed events", "error", qerr)
			return false
		}

		log.Debugw("synthetic event wait progress", "ExecutionStateChanged", len(events))

		return len(events) > 0
	}, 10*time.Minute, 5*time.Second, "no ExecutionStateChanged event was indexed for the message executed in checkpoint %d", uint64(FAILED_TX_CHECKPOINT))

	// 5. Read the event back through the ChainReader, which is how CCIP consumes it.
	events, err := cReader.QueryKey(ctx, boundContracts[0], query.KeyFilter{Key: "ExecutionStateChanged", Expressions: messageIdFilter}, query.LimitAndSort{}, &OfframpExecutionStateChanged{})
	require.NoError(t, err)
	require.NotEmpty(t, events)

	// The same message may also have been executed successfully elsewhere in the polled window; the
	// event written by the transactions indexer is the one in the FAILURE state.
	var executionStateChanged *OfframpExecutionStateChanged
	for _, event := range events {
		if data, ok := event.Data.(*OfframpExecutionStateChanged); ok && data.State == 3 { // 3 = FAILURE
			executionStateChanged = data
			break
		}
	}
	require.NotNil(t, executionStateChanged, "no synthetic ExecutionStateChanged event in the FAILURE state was indexed")

	require.Equal(t, sourceChainSelector, executionStateChanged.SourceChainSelector)
	require.Equal(t, executionReport.Message.Header.SequenceNumber, executionStateChanged.SequenceNumber)
	// The message ID and hash are stored base64-encoded and rendered as hex strings by the response
	// transform driven by `ExpectedEventType` in the ChainReader config.
	require.Equal(t, "0x"+hex.EncodeToString(executionReport.Message.Header.MessageID), executionStateChanged.MessageId)

	// Re-hash the report against the onRamp from the indexed SourceChainConfigSet event: the same
	// inputs the indexer had, so the stored hash must match exactly.
	onRamp := QueryIndexedOnRamp(ctx, t, dbStore, offrampPackageId, sourceChainSelector)
	expectedMessageHash, err := chainreaderutil.NewMessageHasherV1(log).Hash(ctx, executionReport, onRamp)
	require.NoError(t, err)
	require.Equal(t, "0x"+hex.EncodeToString(expectedMessageHash[:]), executionStateChanged.MessageHash)
}

// OFFRAMP_PACKAGE_ID is the original (defining) package ID of the testnet OffRamp deployment the
// known failed transaction belongs to.
//
// A PTB calls into whichever package version was latest when it was submitted, so the package on
// the move call is an upgraded package ID. Move events, however, are always emitted under the
// event type's defining package, which stays at the original ID across upgrades. Binding the
// OffRamp to the move call's package would therefore build event handles no indexed event can ever
// match, which is why the original ID is pinned here instead of being read off the transaction.
const OFFRAMP_PACKAGE_ID = "0x5ef4b483da6644c84aa78eae4f51a9bfb1fb4554d5134ac98892e931fcbdd6bf"

// FindOfframpExecuteCall locates the offramp::init_execute command in a PTB and returns the OffRamp
// package to bind (the original package ID, not the upgraded one the call targets), the command's
// index, and the BCS-encoded execution report passed to it. It mirrors how the TransactionsIndexer
// recovers the report from a failed transmitter transaction.
func FindOfframpExecuteCall(t *testing.T, txn *v2.ExecutedTransaction) (packageId string, commandIndex uint64, reportBytes []byte) {
	t.Helper()

	programmableTxn := txn.GetTransaction().GetKind().GetProgrammableTransaction()
	require.NotNil(t, programmableTxn, "expected %s to be a programmable transaction", txn.GetDigest())

	inputs := programmableTxn.GetInputs()

	for i, cmd := range programmableTxn.GetCommands() {
		moveCall := cmd.GetMoveCall()
		if moveCall == nil || moveCall.GetModule() != "offramp" || moveCall.GetFunction() != "init_execute" {
			continue
		}

		var callArgs []*v2.Input
		for _, arg := range moveCall.GetArguments() {
			if arg.GetKind() != v2.Argument_INPUT {
				continue
			}
			inputIdx := int(arg.GetInput())
			require.Less(t, inputIdx, len(inputs), "input index out of range")
			callArgs = append(callArgs, inputs[inputIdx])
		}

		require.GreaterOrEqual(t, len(callArgs), 5, "expected the execution report in argument position 4")
		report := callArgs[4].GetPure()
		require.NotEmpty(t, report, "expected the execution report argument to be a Pure input")

		t.Logf("offramp::init_execute called on package %s, binding original package %s", moveCall.GetPackage(), OFFRAMP_PACKAGE_ID)

		return OFFRAMP_PACKAGE_ID, uint64(i), report
	}

	t.Fatalf("no offramp::init_execute command found in transaction %s", txn.GetDigest())

	return "", 0, nil
}

// QueryIndexedOnRamp reads the indexed SourceChainConfigSet event for the given source chain and
// returns the onRamp address it carries, which is the value the TransactionsIndexer hashes with.
func QueryIndexedOnRamp(ctx context.Context, t *testing.T, dbStore *database.DBStore, offrampPackageId string, sourceChainSelector uint64) []byte {
	t.Helper()

	events, err := dbStore.QueryEvents(
		ctx,
		offrampPackageId,
		offrampPackageId+"::offramp::SourceChainConfigSet",
		[]query.Expression{
			query.Comparator("sourceChainSelector",
				primitives.ValueComparator{Value: strconv.FormatUint(sourceChainSelector, 10), Operator: primitives.Eq},
			),
		},
		query.LimitAndSort{
			Limit:  query.CountLimit(1),
			SortBy: []query.SortBy{query.NewSortBySequence(query.Desc)},
		},
	)
	require.NoError(t, err)
	require.NotEmpty(t, events, "no SourceChainConfigSet event indexed for source chain selector %d", sourceChainSelector)

	var configEvent codec.SourceChainConfigSet
	require.NoError(t, codec.DecodeSuiJsonValue(events[0].Data, &configEvent))
	require.NotEmpty(t, configEvent.SourceChainConfig.OnRamp, "indexed SourceChainConfigSet has no onRamp")

	return configEvent.SourceChainConfig.OnRamp
}
