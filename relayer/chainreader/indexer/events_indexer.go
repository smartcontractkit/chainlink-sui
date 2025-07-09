package indexer

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

type EventsIndexer struct {
	db                  *database.DBStore
	client              client.SuiPTBClient
	logger              logger.Logger
	pollingInterval     time.Duration
	syncTimeout         time.Duration
	eventConfigurations []*client.EventSelector
	// a map of event handles to the last processed cursor
	lastProcessedCursors map[string]*models.EventId
}

type EventsIndexerApi interface {
	Start(ctx context.Context) error
	SyncAllEvents(ctx context.Context) error
	SyncEvent(ctx context.Context, selector *client.EventSelector) error
}

const batchSizeRecords = 50

func NewEventIndexer(
	db *database.DBStore,
	log logger.Logger,
	ptbClient client.SuiPTBClient,
	eventConfigurations []*client.EventSelector,
	pollingInterval time.Duration,
	syncTimeout time.Duration,
) EventsIndexerApi {
	return &EventsIndexer{
		db:                   db,
		client:               ptbClient,
		logger:               log,
		pollingInterval:      pollingInterval,
		syncTimeout:          syncTimeout,
		eventConfigurations:  eventConfigurations,
		lastProcessedCursors: make(map[string]*models.EventId),
	}
}

func (eIndexer *EventsIndexer) Start(ctx context.Context) error {
	ticker := time.NewTicker(eIndexer.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			syncCtx, cancel := context.WithTimeout(ctx, eIndexer.syncTimeout)
			start := time.Now()

			err := eIndexer.SyncAllEvents(syncCtx)
			elapsed := time.Since(start)

			if err != nil && !errors.Is(err, context.DeadlineExceeded) {
				eIndexer.logger.Warnw("EventSync completed with errors", "error", err, "duration", elapsed)
			} else if err != nil {
				eIndexer.logger.Warnw("EventSync timed out", "duration", elapsed)
			} else {
				eIndexer.logger.Debugw("Event sync completed successfully", "duration", elapsed)
			}

			cancel()
		case <-ctx.Done():
			eIndexer.logger.Infow("Event polling stopped")
			return nil
		}
	}
}

func (eIndexer *EventsIndexer) SyncAllEvents(ctx context.Context) error {
	eIndexer.logger.Debug("SyncAllEvents: starting")

	if eIndexer.db == nil {
		return fmt.Errorf("SyncAllEvents only operates with database store")
	}

	if err := eIndexer.db.EnsureSchema(ctx); err != nil {
		return fmt.Errorf("SyncAllEvents: failed to ensure schema: %w", err)
	}

	successCount := 0
	errorCount := 0
	var lastErr error

	// Iterate through all configured modules and their events
	for _, selector := range eIndexer.eventConfigurations {
		packageAddress, moduleName, eventName := selector.Package, selector.Module, selector.Event

		select {
		case <-ctx.Done():
			if successCount > 0 {
				eIndexer.logger.Infow("SyncAllEvents: interrupted, some events synced", "successCount", successCount, "errorCount", errorCount)
			}

			return ctx.Err()
		default:
			err := eIndexer.SyncEvent(ctx, selector)
			if err != nil {
				errorCount++
				lastErr = fmt.Errorf("SyncAllEvents: module %s event %s: %w", moduleName, eventName, err)
				eIndexer.logger.Errorw("SyncAllEvents: error syncing event",
					"package", packageAddress,
					"module", moduleName, "event",
					eventName, "error", err)
			} else {
				successCount++
			}
		}
	}

	if errorCount > 0 {
		eIndexer.logger.Errorw("SyncAllEvents: completed with errors", "successCount", successCount, "errorCount", errorCount, "lastError", lastErr)
		return lastErr
	}

	eIndexer.logger.Infow("SyncAllEvents: successfully synced all events", "count", successCount)

	return nil
}

func (eIndexer *EventsIndexer) SyncEvent(ctx context.Context, selector *client.EventSelector) error {
	if selector == nil {
		return fmt.Errorf("unspecified selector for SyncEvent call")
	}

	eventHandle := fmt.Sprintf("%s::%s::%s", selector.Package, selector.Module, selector.Event)

	// check if the event selector is already tracked, if not add it to the list
	if !eIndexer.isEventSelectorAdded(*selector) {
		eIndexer.eventConfigurations = append(eIndexer.eventConfigurations, selector)
	}

	eIndexer.logger.Debugw("syncEvent: searching for event", "handle", eventHandle)

	// Get the cursor for pagination - either from memory or start fresh
	cursor := eIndexer.lastProcessedCursors[eventHandle]
	var totalCount uint64
	var err error
	if cursor == nil {
		// attempt to get the latest event sync of the given type and use its data to construct a cursor
		cursor, totalCount, err = eIndexer.db.GetLatestOffset(ctx, selector.Package, eventHandle)
		if err != nil {
			return err
		}

		eIndexer.logger.Debugw("syncEvent: starting fresh sync", "handle", eventHandle, "cursor", cursor)
	}

	batchSize := uint(batchSizeRecords)
	var totalProcessed int

	sortOptions := &client.QuerySortOptions{
		Descending: false, // Process events in chronological order
	}

	// Convert cursor to client format if we have one
	var clientCursor *client.EventId
	if cursor != nil {
		clientCursor = &client.EventId{
			TxDigest: cursor.TxDigest,
			EventSeq: cursor.EventSeq,
		}
	}

eventLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Query events from the Sui blockchain
			eventsPage, err := eIndexer.client.QueryEvents(ctx, *selector, &batchSize, clientCursor, sortOptions)
			if err != nil {
				eIndexer.logger.Errorw("syncEvent: failed to fetch events",
					"error", err, "handle", eventHandle)

				return fmt.Errorf("syncEvent: failed to fetch events: %w", err)
			}

			eIndexer.logger.Debugw("syncEvent: fetched events",
				"count", len(eventsPage.Data),
				"handle", eventHandle,
				"cursor", clientCursor)

			if len(eventsPage.Data) == 0 {
				break eventLoop
			}

			// Convert events to database records
			var batchRecords []database.EventRecord
			for i, event := range eventsPage.Data {
				// Get block information
				block, err := eIndexer.client.BlockByDigest(ctx, event.Id.TxDigest)
				if err != nil {
					eIndexer.logger.Errorw("syncEvent: failed to fetch block metadata",
						"txDigest", event.Id.TxDigest, "error", err)

					continue
				}

				offset, err := strconv.ParseUint(event.Id.EventSeq, 10, 64)
				if err != nil {
					eIndexer.logger.Errorw("syncEvent: failed to parse event offset",
						"eventSeq", event.Id.EventSeq, "error", err)

					continue
				}

				// offset is the event sequence number, we need to add the total number of events processed so far
				// and the index of the event in the current batch
				//nolint:gosec
				offset += uint64(i) + totalCount

				// Convert event to database record
				record := database.EventRecord{
					EventAccountAddress: selector.Package,
					EventHandle:         eventHandle,
					// TODO: event offset is a string and should be stored in the DB as a string
					EventOffset:    offset,
					TxDigest:       event.Id.TxDigest,
					BlockVersion:   0,
					BlockHeight:    fmt.Sprintf("%d", block.Height),
					BlockHash:      []byte(block.TxDigest),
					BlockTimestamp: block.Timestamp,
					Data:           event.ParsedJson,
				}
				batchRecords = append(batchRecords, record)
			}

			// Insert batch of events into database
			if len(batchRecords) > 0 {
				if err := eIndexer.db.InsertEvents(ctx, batchRecords); err != nil {
					return fmt.Errorf("syncEvent: failed to insert batch of events: %w", err)
				}

				totalProcessed += len(batchRecords)
				eIndexer.logger.Debugw("syncEvent: saved batch of events",
					"batch_count", len(batchRecords),
					"total_processed", totalProcessed,
					"handle", eventHandle)
			}

			// Update cursor for next iteration
			if eventsPage.HasNextPage && eventsPage.NextCursor.TxDigest != "" && eventsPage.NextCursor.EventSeq != "" {
				cursor = &models.EventId{
					TxDigest: eventsPage.NextCursor.TxDigest,
					EventSeq: eventsPage.NextCursor.EventSeq,
				}
				clientCursor = &client.EventId{
					TxDigest: eventsPage.NextCursor.TxDigest,
					EventSeq: eventsPage.NextCursor.EventSeq,
				}
				eIndexer.lastProcessedCursors[eventHandle] = cursor
			} else {
				// No more events to process
				break eventLoop
			}

			// If we received fewer events than the batch size, we're caught up
			if uint(len(eventsPage.Data)) < batchSize {
				break eventLoop
			}
		}
	}

	return nil
}

// IsEventSelectorAdded checks if a specific event selector has already been included in the list of events
// to sync
func (eIndexer *EventsIndexer) isEventSelectorAdded(eConfig client.EventSelector) bool {
	for _, selector := range eIndexer.eventConfigurations {
		if selector.Package == eConfig.Package && selector.Module == eConfig.Module && selector.Event == eConfig.Event {
			return true
		}
	}

	return false
}
