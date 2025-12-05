package indexer

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

type EventsIndexer struct {
	db              *database.DBStore
	client          client.SuiPTBClient
	logger          logger.Logger
	pollingInterval time.Duration
	syncTimeout     time.Duration

	// Protected by configMutex
	eventConfigurations []*client.EventSelector
	configMutex         sync.RWMutex

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

// Converts snake_case to camelCase
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if i > 0 && len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}

	return strings.Join(parts, "")
}

// Recursively convert all keys in the map to camelCase,
// with a special case for message.header.sequence_number → seqNum
func convertMapKeysToCamelCase(input any) any {
	return convertMapKeysToCamelCaseWithPath(input, "")
}

func convertMapKeysToCamelCaseWithPath(input any, path string) any {
	switch typed := input.(type) {
	case map[string]any:
		result := make(map[string]any)
		for k, v := range typed {
			camelKey := snakeToCamel(k)
			fullPath := path
			if fullPath != "" {
				fullPath += "." + camelKey
			} else {
				fullPath = camelKey
			}

			if fullPath == "message.header.sequenceNumber" {
				camelKey = "seqNum"
			}

			result[camelKey] = convertMapKeysToCamelCaseWithPath(v, fullPath)
		}

		return result

	case []any:
		for i, v := range typed {
			typed[i] = convertMapKeysToCamelCaseWithPath(v, path)
		}
	}

	return input
}

func (eIndexer *EventsIndexer) SyncEvent(ctx context.Context, selector *client.EventSelector) error {
	if selector == nil {
		return fmt.Errorf("unspecified selector for SyncEvent call")
	}

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
			eIndexer.logger.Errorw("syncEvent: failed to get latest offset", "error", offsetErr)
			return offsetErr
		}

		if dbOffsetCursor != nil {
			txDigestBytes, err := hex.DecodeString(strings.TrimPrefix(dbOffsetCursor.TxDigest, "0x"))
			if err != nil {
				eIndexer.logger.Errorw("syncEvent: failed to decode tx digest", "error", err)
				eIndexer.cursorMutex.RUnlock()
				return err
			}
			// convert the db offset cursor digest from hex (the format stored in the DB) to base58 (the format expected by the client)
			cursor = &models.EventId{
				TxDigest: base58.Encode(txDigestBytes),
				EventSeq: dbOffsetCursor.EventSeq,
			}

			totalCount = dbTotalCount
		} else {
			eIndexer.logger.Debugw("syncEvent: starting fresh sync", "handle", eventHandle)

			// hardcoded cursor starting point
			cursor = &models.EventId{
				TxDigest: "CpFQ8JsaHwTEuNLCfeJQopu3eM3ipViowkWmg23k4fNk",
				EventSeq: "0",
			}
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

				// normalize the data, convert snake case to camel case
				normalizedData := convertMapKeysToCamelCase(event.ParsedJson)

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
