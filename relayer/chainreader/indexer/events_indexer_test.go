//go:build integration

package indexer_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	indexer2 "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestEventsIndexer(t *testing.T) {
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

	// Setup Sui node and account
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

	// Create keystore for PTB client and add the generated key
	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, publicKeyBytes := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	// Fund the account multiple times to ensure sufficient balance
	for i := 0; i < 3; i++ {
		err = testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress)
		require.NoError(t, err)
	}

	ptbClientConfig := client.PTBClientConfig{
		GrpcTarget:            testutils.LocalGrpcURL,
		GrpcToken:             "test",
		TransactionTimeout:    10 * time.Second,
		MaxConcurrentRequests: 10,
		KeystoreService:       keystoreInstance,
		DefaultRequestType:    client.WaitForLocalExecution,
	}

	relayerClient, err := client.NewPTBClient(log, ptbClientConfig)
	require.NoError(t, err)

	chainID, err := testutils.GetChainIdentifier(testutils.LocalURL)
	require.NoError(t, err)
	testutils.PatchEnvironmentTOML("contracts/test", "local", chainID)
	testutils.PatchEnvironmentTOML("contracts/test_secondary", "local", chainID)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	gasBudget := int(2000000000)
	packageId, tx, err := testutils.PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotNil(t, packageId)
	require.NotNil(t, tx)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(tx.ObjectChanges, packageId, "counter", "Counter")
	require.NoError(t, err)

	// Setup event selector
	eventSelector := &client.EventSelector{
		Package: packageId,
		Module:  "counter",
		Event:   "CounterIncremented",
	}

	// Create events indexer
	evIndexer := indexer2.NewEventIndexer(
		db,
		log,
		// start without any selectors, they will be added during .Bind() calls
		[]*client.EventSelector{},
	)

	// Get the latest checkpoint sequence number
	latestCheckpoint, err := relayerClient.GetLatestCheckpoint(ctx)
	require.NoError(t, err)
	require.NotNil(t, latestCheckpoint)
	latestCheckpointSequence := latestCheckpoint.GetSequenceNumber()

	// Create chain poller to feed events to the indexer
	chainPoller := indexer2.NewChainPoller(
		relayerClient,
		log,
		config.ChainPollerConfig{
			PollingInterval:         1 * time.Second,
			SyncTimeout:             120 * time.Second,
			ChannelBufferSize:       16,
			StartCheckpointSequence: &latestCheckpointSequence,
		},
		evIndexer.GetEventSelectors,
	)

	// Start the poller and events indexer directly. This test exercises event
	// indexing only; drain the transactions channel so the poller never blocks.
	err = chainPoller.Start(ctx)
	require.NoError(t, err)
	err = evIndexer.Start(ctx, chainPoller.EventsChannel())
	require.NoError(t, err)

	go func() {
		for range chainPoller.TransactionsChannel() {
		}
	}()

	// Helper function to create events by calling contract
	createEvent := func(eventNum int) {
		log.Debugw("Creating event by calling contract", "eventNumber", eventNum)

		moveCallReq := client.MoveCallRequest{
			Signer:          accountAddress,
			PackageObjectId: packageId,
			Module:          "counter",
			Function:        "increment",
			TypeArguments:   []any{},
			Arguments:       []any{counterObjectId},
			GasBudget:       2000000,
		}

		txMetadata, callErr := relayerClient.MoveCall(ctx, moveCallReq)
		require.NoError(t, callErr)

		txnResult, sendErr := relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes)
		require.NoError(t, sendErr)

		// wait for transaction confirmation
		require.Eventually(t, func() bool {
			status, err := relayerClient.GetTransactionStatus(ctx, txnResult.Transaction.GetDigest())
			return err == nil && status.Status == "success"
		}, 10*time.Second, 200*time.Millisecond)

		log.Debugw("Event created successfully", "eventNumber", eventNum, "txDigest", txnResult.Transaction.GetDigest())
	}

	// Helper function to wait for events to be indexed
	waitForEventCount := func(expectedCount int, timeout time.Duration) []database.EventRecord {
		log.Debugw("Waiting for events to be indexed", "expectedCount", expectedCount)

		var events []database.EventRecord
		eventHandle := packageId + "::" + eventSelector.Module + "::" + eventSelector.Event

		require.Eventually(t, func() bool {
			var err error
			events, err = dbStore.QueryEvents(ctx, packageId, eventHandle, nil, query.LimitAndSort{
				Limit: query.Limit{
					//nolint:gosec
					Count: uint64(expectedCount) + uint64(1),
				},
			})
			if err != nil {
				log.Errorw("Failed to query events", "error", err)
				return false
			}

			log.Debugw("Current event count", "count", len(events), "expected", expectedCount)

			return len(events) >= expectedCount
		}, timeout, 1*time.Second, "Should find %d events", expectedCount)

		return events
	}

	t.Run("TestBasicEventIndexingViaChannel", func(t *testing.T) {
		log.Infow("Starting basic channel-based event indexing test")

		// Add the event selector so the poller starts filtering for it
		err := evIndexer.AddEventSelector(ctx, eventSelector)
		require.NoError(t, err)

		// Create 3 events
		for i := 1; i <= 3; i++ {
			createEvent(i)
		}

		// Wait for events to be indexed via the channel
		events := waitForEventCount(3, 60*time.Second)

		log.Infow("Fetched all events", "eventsFound", len(events))

		// Verify events have correct sequential values
		for i, event := range events[:3] {
			log.Debugw("Event details",
				"index", i,
				"offset", event.EventOffset,
				"txDigest", event.TxDigest,
				"data", event.Data)

			// Verify event data
			require.NotNil(t, event.Data)
			newValue, ok := event.Data["newValue"]
			require.True(t, ok, "Event should have newValue field")

			expectedValue := strconv.Itoa(3 - i)
			require.Equal(t, expectedValue, newValue, "Event %d should have newValue %d", i, expectedValue)
		}

		// Verify the offset is tracked correctly
		eventHandle := packageId + "::" + eventSelector.Module + "::" + eventSelector.Event
		cursor, totalCount, err := dbStore.GetLatestOffset(ctx, packageId, eventHandle)
		require.NoError(t, err)
		require.NotNil(t, cursor)
		require.Equal(t, uint64(3), totalCount, "Should have 3 events total")
	})

	t.Run("TestMultipleSyncOperationsViaChannel", func(t *testing.T) {
		log.Infow("Testing multiple sync operations via channel")

		// Create more events
		for i := 4; i <= 7; i++ {
			createEvent(i)
		}

		// Wait for all events to be indexed via the channel pipeline
		allEvents := waitForEventCount(7, 60*time.Second)
		require.GreaterOrEqual(t, len(allEvents), 7, "Should have at least 7 events")

		log.Infow("Fetched all events", "events", len(allEvents))
	})

	t.Run("TestWithTimestamps", func(t *testing.T) {
		log.Infow("Testing with timestamps")

		// Trigger some events
		for i := 8; i <= 10; i++ {
			createEvent(i)
		}

		// Wait for events to be indexed
		events := waitForEventCount(10, 120*time.Second)
		require.GreaterOrEqual(t, len(events), 10)

		// Check that events are recorded with timestamps in seconds
		for _, event := range events[:3] {
			require.Greater(t, event.BlockTimestamp, uint64(0), "Event should have a timestamp")
			require.Less(t, event.BlockTimestamp, uint64(time.Now().Unix()+1), "Event timestamp should be in the past")
		}
	})

	t.Run("TestOrderedEventsQueryWithOutOfOrderEventOffset", func(t *testing.T) {
		t.Skip("Skipping test ordered events query with out of order event offset until the relevant index is re-added")
		// insert duplicate events with out of order event_offset for CCIPMessageSent

		packageId := "0x30e087460af8a8aacccbc218aa358cdcde8d43faf61ec0638d71108e276e2f1d"
		eventHandle := packageId + "::onramp::CCIPMessageSent"
		baseRecord := database.EventRecord{
			EventAccountAddress: accountAddress,
			EventHandle:         eventHandle,
			EventOffset:         0,
			TxDigest:            "5HueCGU5rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ",
			BlockVersion:        0,
			BlockHeight:         "100",
			BlockHash:           []byte("5HueCGU5rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"),
			BlockTimestamp:      1000000000,
			Data:                map[string]any{},
			IsSynthetic:         false,
		}

		// insert duplicate and incorrect event offsets
		for i := range 200_000 {
			recordA := baseRecord
			recordB := baseRecord

			// use different event_offset for both records
			recordA.EventOffset = uint64(i)
			recordB.EventOffset = uint64(i%2 + 1000)

			// use duplicate data for both records
			recordA.BlockHeight = strconv.Itoa(100 + i)
			recordB.BlockHeight = strconv.Itoa(100 + i)

			recordA.TxDigest = base58.Encode([]byte("record" + strconv.Itoa(i)))
			recordB.TxDigest = base58.Encode([]byte("record" + strconv.Itoa(i)))

			recordA.Data = map[string]any{
				"destChainSelector": 3478487238524512106,
				"sequenceNumber":    776 + uint64(i),
			}
			recordB.Data = map[string]any{
				"destChainSelector": 3478487238524512106,
				"sequenceNumber":    776 + uint64(i),
			}

			dbStore.InsertEvents(ctx, []database.EventRecord{recordA, recordB})
		}

		// insert some other unrelated events
		for i := range 10_000 {
			recordA := baseRecord

			// use different event_offset for both records
			recordA.EventOffset = uint64(i + 1)

			// use duplicate data for both records
			recordA.BlockHeight = "100"

			recordA.TxDigest = base58.Encode([]byte("record" + strconv.Itoa(i)))

			recordA.EventHandle = packageId + "::onramp::SomeOtherEvent"

			recordA.Data = map[string]any{
				"destChainSelector": 3478487238524512106,
				"sequenceNumber":    176 + uint64(i),
			}

			dbStore.InsertEvents(ctx, []database.EventRecord{recordA})
		}

		// query events with out of order event_offset
		events, err := dbStore.QueryEvents(ctx, accountAddress, eventHandle, []query.Expression{
			{
				BoolExpression: query.BoolExpression{
					BoolOperator: query.AND,
					Expressions: []query.Expression{
						{
							Primitive: &primitives.Comparator{
								Name: "sequenceNumber",
								ValueComparators: []primitives.ValueComparator{
									{Value: uint64(776), Operator: primitives.Gte},
									{Value: uint64(779), Operator: primitives.Lte},
								},
							},
						},
						{
							Primitive: &primitives.Comparator{
								Name: "destChainSelector",
								ValueComparators: []primitives.ValueComparator{
									{Value: "3478487238524512106", Operator: primitives.Eq},
								},
							},
						},
					},
				},
			},
		}, query.LimitAndSort{
			Limit: query.Limit{
				Count: 100,
			},
			SortBy: []query.SortBy{
				query.NewSortBySequence(query.Asc),
			},
		})

		// we should only get 10 events
		require.NoError(t, err)
		require.Equal(t, 4, len(events))

		for _, event := range events {
			fmt.Printf("eventHandle: %s\n", event.EventHandle)
			fmt.Printf("sequenceNumber: %f\n", event.Data["sequenceNumber"].(float64))
			fmt.Println("--------------------------------")
		}

		// events should have strictly increasing sequence numbers and be in order
		for i := range len(events) - 1 {
			require.Equal(t, events[i].Data["sequenceNumber"].(float64)+1, events[i+1].Data["sequenceNumber"].(float64))
		}

		// query another range for the same event handle
		events, err = dbStore.QueryEvents(ctx, accountAddress, eventHandle, []query.Expression{
			{
				BoolExpression: query.BoolExpression{
					BoolOperator: query.AND,
					Expressions: []query.Expression{
						{
							Primitive: &primitives.Comparator{
								Name: "sequenceNumber",
								ValueComparators: []primitives.ValueComparator{
									{Value: uint64(779), Operator: primitives.Gte},
									{Value: uint64(785), Operator: primitives.Lte},
								},
							},
						},
						{
							Primitive: &primitives.Comparator{
								Name: "destChainSelector",
								ValueComparators: []primitives.ValueComparator{
									{Value: "3478487238524512106", Operator: primitives.Eq},
								},
							},
						},
					},
				},
			},
		}, query.LimitAndSort{
			Limit: query.Limit{
				Count: 100,
			},
			SortBy: []query.SortBy{
				query.NewSortBySequence(query.Asc),
			},
		})

		require.NoError(t, err)
		require.Equal(t, 7, len(events))

		for _, event := range events {
			fmt.Printf("eventHandle: %s\n", event.EventHandle)
			fmt.Printf("sequenceNumber: %f\n", event.Data["sequenceNumber"].(float64))
			fmt.Println("--------------------------------")
		}

		// events should have strictly increasing sequence numbers and be in order
		for i := range len(events) - 1 {
			require.Equal(t, events[i].EventHandle, eventHandle)
			require.Equal(t, events[i].Data["sequenceNumber"].(float64)+1, events[i+1].Data["sequenceNumber"].(float64))
		}
	})

	t.Run("TestSyntheticEventsSkipForOffset", func(t *testing.T) {
		eventHandle := packageId + "::offramp::ExecutionStateChanged"
		record := database.EventRecord{
			EventAccountAddress: accountAddress,
			EventHandle:         eventHandle,
			EventOffset:         0,
			TxDigest:            "fake_digest",
			BlockVersion:        0,
			BlockHeight:         "100",
			BlockHash:           []byte("5HueCGU5rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"),
			BlockTimestamp:      1000000000,
			Data:                map[string]any{},
			IsSynthetic:         true,
		}

		recordB := database.EventRecord{
			EventAccountAddress: accountAddress,
			EventHandle:         eventHandle,
			EventOffset:         1,
			TxDigest:            "real_digest",
			BlockVersion:        0,
			BlockHeight:         "100",
			BlockHash:           []byte("5HueCGU5rMjxEXxiPuD5BDku4MkFqeZyd4dZ1jvhTVqvbTLvyTJ"),
			BlockTimestamp:      1000000000,
			Data:                map[string]any{},
			IsSynthetic:         false,
		}

		dbStore.InsertEvents(ctx, []database.EventRecord{record, recordB})

		// query events with out of order event_offset
		cursor, totalCount, err := dbStore.GetLatestOffset(ctx, accountAddress, eventHandle)
		require.NoError(t, err)
		require.Equal(t, recordB.TxDigest, cursor.TxDigest)
		require.Equal(t, uint64(2), totalCount)
	})

	// Cleanup
	require.NoError(t, chainPoller.Close())
	require.NoError(t, evIndexer.Close())
}
