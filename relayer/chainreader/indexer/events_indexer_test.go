//go:build integration

package indexer_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/mr-tron/base58"

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

	relayerClient, err := client.NewPTBClient(log, testutils.LocalUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, err)

	chainID, err := testutils.GetChainIdentifier(testutils.LocalUrl)
	require.NoError(t, err)
	testutils.PatchEnvironmentTOML("contracts/test", "local", chainID)

	// testutils.PatchContractAddressTOML(t, "contracts/test", "test_secondary", "_")

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
	pollingInterval := time.Second
	syncTimeout := 10 * time.Second

	indexer := indexer2.NewEventIndexer(
		db,
		log,
		relayerClient,
		[]*client.EventSelector{eventSelector},
		pollingInterval,
		syncTimeout,
	)

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

		txnResult, sendErr := relayerClient.SignAndSendTransaction(ctx, txMetadata.TxBytes, publicKeyBytes, "WaitForLocalExecution")
		require.NoError(t, sendErr)

		log.Debugw("Event created successfully", "eventNumber", eventNum, "txDigest", txnResult.TxDigest)
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
		}, timeout, 500*time.Millisecond, "Should find %d events", expectedCount)

		return events
	}

	// Helper function to wait for events to be indexed
	waitForEventCountFromDB := func(expectedCount int, timeout time.Duration) []database.EventRecord {
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
		}, timeout, 500*time.Millisecond, "Should find %d events", expectedCount)

		return events
	}

	t.Run("TestCursorAndOffsetBasicFunctionality", func(t *testing.T) {
		log.Infow("Starting basic cursor and offset functionality test")

		// create initial events and test basic indexing
		t.Run("InitialSync", func(t *testing.T) {
			log.Infow("Creating initial events")

			// Create 3 events
			for i := 1; i <= 3; i++ {
				createEvent(i)
			}

			// Run sync to index events
			err := indexer.SyncEvent(ctx, eventSelector)
			require.NoError(t, err)

			// Wait for events to be indexed
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

			// Verify the cursor is set correctly
			eventHandle := packageId + "::" + eventSelector.Module + "::" + eventSelector.Event
			cursor, totalCount, err := dbStore.GetLatestOffset(ctx, packageId, eventHandle)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, uint64(3), totalCount, "Should have 3 events total")
		})

		// Test GetLatestOffset functionality
		t.Run("GetLatestOffset", func(t *testing.T) {
			log.Infow("Testing GetLatestOffset")

			eventHandle := packageId + "::" + eventSelector.Module + "::" + eventSelector.Event

			// Get the latest offset from database
			cursor, totalCount, err := dbStore.GetLatestOffset(ctx, packageId, eventHandle)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, uint64(3), totalCount, "Should have 3 events total")

			log.Debugw("Latest offset details",
				"cursor", cursor,
				"totalCount", totalCount)

			// The cursor should reflect the latest event
			require.NotEmpty(t, cursor.TxDigest, "Cursor should have TxDigest")
			require.NotEmpty(t, cursor.EventSeq, "Cursor should have EventSeq")
			require.Equal(t, uint64(3), totalCount, "Should have 3 events total")

			// Create more events
			for i := 4; i <= 6; i++ {
				createEvent(i)
			}

			// Run sync to index events
			err = indexer.SyncEvent(ctx, eventSelector)
			require.NoError(t, err)

			// Get the latest offset from database
			cursor, totalCount, err = dbStore.GetLatestOffset(ctx, packageId, eventHandle)
			require.NoError(t, err)
			require.NotNil(t, cursor)
			require.Equal(t, uint64(6), totalCount, "Should have 6 events total")
		})

		// Test multiple sync operations
		t.Run("MultipleSyncOperations", func(t *testing.T) {
			log.Infow("Testing multiple sync operations")

			// Create more events
			for i := 6; i <= 8; i++ {
				createEvent(i)
			}

			// Run sync multiple times to test idempotency
			for i := range 3 {
				err := indexer.SyncEvent(ctx, eventSelector)
				require.NoError(t, err)
				log.Debugw("Sync operation completed", "iteration", i+1)
			}

			// Wait for all events to be indexed
			allEvents := waitForEventCountFromDB(7, 60*time.Second)

			log.Infow("Fetched all events", "events", allEvents)
		})
	})

	t.Run("TestCursorAdvancementValidation", func(t *testing.T) {
		log.Infow("Testing cursor advancement validation")

		// This test validates that cursors advance properly between sync operations
		// by creating events in batches and checking cursor progression

		// Create a fresh event selector for isolation
		freshEventSelector := &client.EventSelector{
			Package: packageId,
			Module:  "counter",
			Event:   "CounterIncremented",
		}

		freshIndexer := indexer2.NewEventIndexer(
			db,
			log,
			relayerClient,
			[]*client.EventSelector{freshEventSelector},
			pollingInterval,
			syncTimeout,
		)

		// Create first batch of events
		log.Infow("Creating first batch of events")
		for i := 1; i <= 2; i++ {
			createEvent(i)
		}

		// Run first sync
		err := freshIndexer.SyncEvent(ctx, freshEventSelector)
		require.NoError(t, err)

		// Get cursor after first sync
		eventHandle := packageId + "::" + freshEventSelector.Module + "::" + freshEventSelector.Event
		cursor1, count1, err := dbStore.GetLatestOffset(ctx, packageId, eventHandle)
		require.NoError(t, err)
		log.Debugw("First sync cursor", "cursor", cursor1, "count", count1)

		// Create second batch of events
		log.Infow("Creating second batch of events")
		for i := 3; i <= 4; i++ {
			createEvent(i)
		}

		// Run second sync
		err = freshIndexer.SyncEvent(ctx, freshEventSelector)
		require.NoError(t, err)

		// Get cursor after second sync
		cursor2, count2, err := dbStore.GetLatestOffset(ctx, packageId, eventHandle)
		require.NoError(t, err)
		log.Debugw("Second sync cursor", "cursor", cursor2, "count", count2)

		// Verify cursor advancement
		require.Greater(t, count2, count1, "Event count should increase after second sync")

		// If cursors are the same, it might indicate the cursor update bug
		if cursor1 != nil && cursor2 != nil {
			// The cursors should be different if new events were processed
			// This helps identify the cursor update bug
			log.Infow("Cursor comparison",
				"cursor1", cursor1,
				"cursor2", cursor2,
				"same", cursor1.TxDigest == cursor2.TxDigest && cursor1.EventSeq == cursor2.EventSeq)
		}
	})

	t.Run("TestConcurrentEventIndexerAccess", func(t *testing.T) {
		log.Infow("Starting concurrent access test")

		// Start the background poller in a goroutine
		pollerCtx, cancelPoller := context.WithCancel(ctx)
		defer cancelPoller()

		var pollerWg sync.WaitGroup
		pollerWg.Add(1)
		go func() {
			defer pollerWg.Done()
			err := indexer.Start(pollerCtx)
			if err != nil && err != context.Canceled {
				log.Errorw("Background poller error", "error", err)
			}
		}()

		// Give the poller time to start
		time.Sleep(200 * time.Millisecond)

		t.Run("ConcurrentSyncEventCalls", func(t *testing.T) {
			log.Infow("Testing concurrent SyncEvent calls")

			var wg sync.WaitGroup
			numConcurrentCallers := 10
			numIterations := 5

			// Track errors from concurrent operations
			errChan := make(chan error, numConcurrentCallers*numIterations)

			// Call SyncEvent concurrently
			for i := 0; i < numConcurrentCallers; i++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()

					for j := 0; j < numIterations; j++ {
						// Each goroutine tries to sync the same event selector
						err := indexer.SyncEvent(ctx, eventSelector)
						if err != nil {
							log.Errorw("SyncEvent error",
								"goroutineID", goroutineID,
								"iteration", j,
								"error", err)
							errChan <- err
						}

						// Small sleep to vary timing
						time.Sleep(time.Duration(goroutineID*5) * time.Millisecond)
					}
				}(i)
			}

			// Wait for all goroutines to complete
			wg.Wait()
			close(errChan)

			// Check for any errors
			var errors []error
			for err := range errChan {
				errors = append(errors, err)
			}

			require.Empty(t, errors, "Should not have errors during concurrent operations")
			log.Infow("Concurrent SyncEvent calls completed", "totalErrors", len(errors))
		})

		t.Run("ConcurrentNewEventSelectors", func(t *testing.T) {
			log.Infow("Testing concurrent new event selectors")

			var wg sync.WaitGroup
			numConcurrentCallers := 20

			// Track errors from concurrent operations
			errChan := make(chan error, numConcurrentCallers)

			// Create different event selectors for each goroutine
			for i := 0; i < numConcurrentCallers; i++ {
				wg.Add(1)
				go func(goroutineID int) {
					defer wg.Done()

					// Create a new selector (simulating ad-hoc event registration)
					newSelector := &client.EventSelector{
						Package: packageId,
						Module:  "counter",
						Event:   "CounterIncremented",
					}

					// Try to sync with this new selector
					err := indexer.SyncEvent(ctx, newSelector)
					if err != nil {
						log.Errorw("SyncEvent error with new selector",
							"goroutineID", goroutineID,
							"error", err)
						errChan <- err
					}
				}(i)
			}

			// Wait for all goroutines to complete
			wg.Wait()
			close(errChan)

			// Check for any errors
			var errors []error
			for err := range errChan {
				errors = append(errors, err)
			}

			require.Empty(t, errors, "Should not have errors during concurrent new selector operations")
			log.Infow("Concurrent new event selectors completed", "totalErrors", len(errors))
		})

		t.Run("ConcurrentReadsAndWrites", func(t *testing.T) {
			log.Infow("Testing concurrent reads and writes stress test")

			var wg sync.WaitGroup
			testDuration := 2 * time.Second
			deadline := time.Now().Add(testDuration)

			// Track errors
			errChan := make(chan error, 100)

			// Multiple goroutines hammering SyncEvent
			for i := 0; i < 5; i++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					for time.Now().Before(deadline) {
						if err := indexer.SyncEvent(ctx, eventSelector); err != nil {
							select {
							case errChan <- err:
							default:
							}
						}
						time.Sleep(10 * time.Millisecond)
					}
				}(i)
			}

			// Wait for stress test to complete
			wg.Wait()
			close(errChan)

			// Collect errors
			var errors []error
			for err := range errChan {
				errors = append(errors, err)
			}

			require.Empty(t, errors, "Should not have errors during stress test")
			log.Infow("Stress test completed", "duration", testDuration, "errors", len(errors))
		})

		// Stop the background poller
		cancelPoller()
		pollerWg.Wait()

		log.Infow("All concurrent access tests completed successfully")
	})

	t.Run("TestWithTimestamps", func(t *testing.T) {
		log.Infow("Testing with timestamps")

		// Trigger some events
		for i := 1; i <= 3; i++ {
			createEvent(i)
		}

		// Create a new event selector for timestamps
		timestampEventSelector := &client.EventSelector{
			Package: packageId,
			Module:  "counter",
			Event:   "CounterIncremented",
		}

		// Run sync to index events
		err := indexer.SyncEvent(ctx, timestampEventSelector)
		require.NoError(t, err)

		// Wait for events to be indexed
		events := waitForEventCount(3, 60*time.Second)
		require.GreaterOrEqual(t, len(events), 3)

		// Check that events are recorded with timestamps in seconds
		for _, event := range events[:3] {
			require.Greater(t, event.BlockTimestamp, uint64(0), "Event should have a timestamp")
			require.Less(t, event.BlockTimestamp, uint64(time.Now().Unix()+1), "Event timestamp should be in the past")
		}
	})

	t.Run("TestRaceDetection", func(t *testing.T) {
		// Run with: go test -race -run TestEventsIndexer/TestRaceDetection
		log.Infow("Starting race detection test")

		// Create a fresh indexer with very short polling interval
		raceIndexer := indexer2.NewEventIndexer(
			db,
			log,
			relayerClient,
			[]*client.EventSelector{eventSelector},
			50*time.Millisecond,
			5*time.Second,
		)

		// Start background poller with short timeout
		pollerCtx, cancelPoller := context.WithTimeout(ctx, 1*time.Second)
		defer cancelPoller()

		var wg sync.WaitGroup

		// Start background poller
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = raceIndexer.Start(pollerCtx)
		}()

		// Hammer with concurrent operations
		for i := 0; i < 30; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				_ = raceIndexer.SyncEvent(ctx, eventSelector)
			}(i)
		}

		// Wait for all operations
		wg.Wait()

		log.Infow("Race detection test completed - no races detected!")
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
}
