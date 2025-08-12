//go:build integration

package indexer_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	indexer2 "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

//nolint:paralleltest
func TestEventsIndexer(t *testing.T) {
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

	// Setup Sui node and account
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

	// Deploy contract with proper account consistency
	contractPath := testutils.BuildSetup(t, "contracts/test")
	testutils.BuildContract(t, contractPath)

	packageId, publishOutput, err := testutils.PublishContract(t, "TestContract", contractPath, accountAddress, nil)
	require.NoError(t, err)

	log.Debugw("Published Contract", "packageId", packageId)

	counterObjectId, err := testutils.QueryCreatedObjectID(publishOutput.ObjectChanges, packageId, "counter", "Counter")
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

				// The newValue should be i+1 (counter starts from 0, so increment makes it 1, 2, 3)
				expectedValue := strconv.Itoa(i + 1)
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
}
