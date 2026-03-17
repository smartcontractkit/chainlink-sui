package indexer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"
)

type EventsIndexer struct {
	db              *database.DBStore
	client          client.SuiPTBClient
	logger          logger.Logger
	pollingInterval time.Duration
	syncTimeout     time.Duration

	// Protected by configMutex
	eventConfigurations  []*client.EventSelector
	eventOffsetOverrides map[string]client.EventId
	configMutex          sync.RWMutex

	// Protected by cursorMutex
	// a map of event handles to the last processed cursor
	lastProcessedCursors map[string]*models.EventId
	cursorMutex          sync.RWMutex
}

type EventsIndexerApi interface {
	Start(ctx context.Context) error
	SyncAllEvents(ctx context.Context) error
	SyncEvent(ctx context.Context, selector *client.EventSelector) error
	AddEventSelector(ctx context.Context, selector *client.EventSelector) error
	SetEventOffsetOverrides(ctx context.Context, offsetOverrides map[string]client.EventId) error
	Ready() error
	Close() error
}

const batchSizeRecords = 50

func NewEventIndexer(
	db sqlutil.DataSource,
	log logger.Logger,
	ptbClient client.SuiPTBClient,
	eventConfigurations []*client.EventSelector,
	pollingInterval time.Duration,
	syncTimeout time.Duration,
) EventsIndexerApi {
	dataStore := database.NewDBStore(db, log)
	namedLogger := logger.Named(log, "EventsIndexer")

	return &EventsIndexer{
		db:                   dataStore,
		client:               ptbClient,
		logger:               namedLogger,
		pollingInterval:      pollingInterval,
		syncTimeout:          syncTimeout,
		eventConfigurations:  eventConfigurations,
		lastProcessedCursors: make(map[string]*models.EventId),
	}
}

func (eIndexer *EventsIndexer) Start(ctx context.Context) error {

	syncCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := eIndexer.db.EnsureSchema(syncCtx); err != nil {
		return fmt.Errorf("Start: failed to ensure schema: %w", err)
	}

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

	successCount := 0
	errorCount := 0
	var lastErr error

	// Avoid holding lock during iteration by making a copy of the selectors
	eIndexer.configMutex.RLock()
	selectors := make([]*client.EventSelector, len(eIndexer.eventConfigurations))
	copy(selectors, eIndexer.eventConfigurations)
	eIndexer.configMutex.RUnlock()

	// Iterate through all configured modules and their events
	for _, selector := range selectors {
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

	eventKey := fmt.Sprintf("%s::%s", selector.Module, selector.Event)
	eventHandle := fmt.Sprintf("%s::%s::%s", selector.Package, selector.Module, selector.Event)

	// check if the event selector is already tracked, if not add it to the list
	if !eIndexer.isEventSelectorAdded(*selector) {
		eIndexer.configMutex.Lock()
		// Double-check after acquiring write lock (avoid race with concurrent adds)
		if !eIndexer.isEventSelectorAddedLocked(*selector) {
			eIndexer.eventConfigurations = append(eIndexer.eventConfigurations, selector)
		}
		eIndexer.configMutex.Unlock()
	}

	eIndexer.logger.Debugw("syncEvent: searching for event", "handle", eventHandle)

	// Get the cursor for pagination - either from memory or start fresh
	eIndexer.cursorMutex.RLock()
	cursor := eIndexer.lastProcessedCursors[eventHandle]

	var totalCount uint64

	if cursor == nil {
		// attempt to get the latest event sync of the given type and use its data to construct a cursor
		dbOffsetCursor, dbTotalCount, offsetErr := eIndexer.db.GetLatestOffset(ctx, selector.Package, eventHandle)
		if offsetErr != nil {
			eIndexer.cursorMutex.RUnlock()
			eIndexer.logger.Errorw("syncEvent: failed to get latest offset", "error", offsetErr)
			return offsetErr
		}

		if dbOffsetCursor != nil {
			// Some DB records have hex formatted txDigest while newer entries have base58 formatted txDigest.
			// We check if the txDigest is hex formatted and decode it if needed for backwards compatibility.
			if strings.ToLower(dbOffsetCursor.TxDigest[:2]) == "0x" {
				txDigestBytes, err := hex.DecodeString(dbOffsetCursor.TxDigest[2:])
				if err != nil {
					eIndexer.cursorMutex.RUnlock()
					eIndexer.logger.Errorw("syncEvent: failed to decode tx digest", "error", err, "txDigest", dbOffsetCursor.TxDigest)
					return fmt.Errorf("syncEvent: failed to decode tx digest: %w", err)
				}
				// convert the db offset cursor digest from hex (the format stored in the DB) to base58 (the format expected by the client)
				cursor = &models.EventId{
					TxDigest: base58.Encode(txDigestBytes),
					EventSeq: dbOffsetCursor.EventSeq,
				}
			} else {
				// DB already has base58
				cursor = &models.EventId{
					TxDigest: dbOffsetCursor.TxDigest,
					EventSeq: dbOffsetCursor.EventSeq,
				}
			}

			totalCount = dbTotalCount
		} else {
			eIndexer.configMutex.RLock()
			if override, ok := eIndexer.eventOffsetOverrides[eventKey]; ok {
				cursor = &models.EventId{
					TxDigest: override.TxDigest,
					EventSeq: override.EventSeq,
				}
			} else {
				eIndexer.logger.Debugw("syncEvent: starting fresh sync", "handle", eventHandle)
			}
			eIndexer.configMutex.RUnlock()
		}
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
	eIndexer.cursorMutex.RUnlock()

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

				// We simply increment the inserted event count to get the next offset.
				// Each new loop interation will increment the totalCount, so we add 1 to the offset.
				// This is done on LoC 375 if there are more events pages to fetch and process.
				//nolint:gosec
				offset := uint64(i) + totalCount + 1

				// normalize the data, convert snake case to camel case
				normalizedData := common.ConvertMapKeysToCamelCase(event.ParsedJson)

				// Convert the txDigest to hex
				txDigestHex := event.Id.TxDigest
				if base58Bytes, err := base58.Decode(txDigestHex); err == nil {
					hexTxId := hex.EncodeToString(base58Bytes)
					txDigestHex = "0x" + hexTxId
				}

				blockHashBytes, err := base58.Decode(block.TxDigest)
				if err != nil {
					eIndexer.logger.Errorw("Failed to decode block hash", "error", err)
					// fallback
					blockHashBytes = []byte(block.TxDigest)
				}

				// Convert event to database record
				record := database.EventRecord{
					EventAccountAddress: selector.Package,
					EventHandle:         eventHandle,
					EventOffset:         offset,
					TxDigest:            txDigestHex,
					BlockVersion:        0,
					BlockHeight:         fmt.Sprintf("%d", block.Height),
					BlockHash:           blockHashBytes,
					// Sui returns block.Timestamp in ms; convert to seconds for consistency with CCIP readers.
					BlockTimestamp: block.Timestamp / 1000,
					Data:           normalizedData.(map[string]any),
					IsSynthetic:    false,
				}
				batchRecords = append(batchRecords, record)
			}

			// Insert batch of events into database
			if len(batchRecords) > 0 {
				if err := eIndexer.db.InsertEvents(ctx, batchRecords); err != nil {
					eIndexer.logger.Errorw("syncEvent: failed to insert batch of events, falling back to per-event insert", "error", err)

					// Fallback: insert each record individually, skip bad ones
					totalProcessedFallback := 0
					for _, record := range batchRecords {
						if err := eIndexer.db.InsertEvents(ctx, []database.EventRecord{record}); err != nil {
							eIndexer.logger.Errorw("Failed to insert single event, skipping...",
								"error", err,
								"handle", eventHandle,
								"txDigest", record.TxDigest,
								"offset", record.EventOffset,
							)

							continue
						}

						totalProcessedFallback++
					}
					eIndexer.logger.Debugw("syncEvent: inserted batch of events", "count", totalProcessedFallback, "handle", eventHandle)
					totalProcessed += totalProcessedFallback
				} else {
					totalProcessed += len(batchRecords)
				}

				eIndexer.logger.Debugw("syncEvent: saved batch of events",
					"batch_count", len(batchRecords),
					"total_processed", totalProcessed,
					"handle", eventHandle)
			}

			// Update cursor for next iteration and the total count of events processed so far
			if eventsPage.HasNextPage && eventsPage.NextCursor.TxDigest != "" && eventsPage.NextCursor.EventSeq != "" {
				cursor = &models.EventId{
					TxDigest: eventsPage.NextCursor.TxDigest,
					EventSeq: eventsPage.NextCursor.EventSeq,
				}
				clientCursor = &client.EventId{
					TxDigest: eventsPage.NextCursor.TxDigest,
					EventSeq: eventsPage.NextCursor.EventSeq,
				}

				eIndexer.cursorMutex.Lock()
				eIndexer.lastProcessedCursors[eventHandle] = cursor
				eIndexer.cursorMutex.Unlock()

				totalCount, err = eIndexer.db.GetTotalCount(ctx, selector.Package, eventHandle)
				if err != nil {
					return fmt.Errorf("syncEvent: failed to get total count: %w", err)
				}
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

func (eIndexer *EventsIndexer) AddEventSelector(ctx context.Context, selector *client.EventSelector) error {
	if selector == nil {
		return fmt.Errorf("unspecified selector for AddEventSelector call")
	}

	// check if the event selector is already tracked, if not add it to the list
	if !eIndexer.isEventSelectorAdded(*selector) {
		eIndexer.configMutex.Lock()
		// Double-check after acquiring write lock (avoid race with concurrent adds)
		if !eIndexer.isEventSelectorAddedLocked(*selector) {
			eIndexer.eventConfigurations = append(eIndexer.eventConfigurations, selector)
		}
		eIndexer.configMutex.Unlock()
	}

	return nil
}

func (eIndexer *EventsIndexer) SetEventOffsetOverrides(ctx context.Context, offsetOverrides map[string]client.EventId) error {
	eIndexer.configMutex.Lock()
	defer eIndexer.configMutex.Unlock()
	eIndexer.eventOffsetOverrides = offsetOverrides
	return nil
}

// IsEventSelectorAdded checks if a specific event selector has already been included in the list of events to sync
func (eIndexer *EventsIndexer) isEventSelectorAdded(eConfig client.EventSelector) bool {
	eIndexer.configMutex.RLock()
	defer eIndexer.configMutex.RUnlock()
	return eIndexer.isEventSelectorAddedLocked(eConfig)
}

// isEventSelectorAddedLocked assumes the lock is already held
func (eIndexer *EventsIndexer) isEventSelectorAddedLocked(eConfig client.EventSelector) bool {
	for _, selector := range eIndexer.eventConfigurations {
		if selector.Package == eConfig.Package && selector.Module == eConfig.Module && selector.Event == eConfig.Event {
			return true
		}
	}

	return false
}

func (eIndexer *EventsIndexer) Ready() error {
	// TODO: implement
	return nil
}

func (eIndexer *EventsIndexer) Close() error {
	// TODO: implement
	return nil
}
