package indexer

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/mr-tron/base58"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"
)

// EventsIndexer consumes checkpoint event batches from a channel and persists them to the database.
// It no longer polls the RPC directly; instead, the ChainPoller feeds it checkpoint data.
type EventsIndexer struct {
	db     *database.DBStore
	logger logger.Logger
	wg     sync.WaitGroup

	// Protected by configMutex
	eventConfigurations []*sui.EventFilterByMoveEventModule
	configMutex         sync.RWMutex

	starter services.StateMachine
}

// EventsIndexerApi defines the interface for the events indexer.
type EventsIndexerApi interface {
	// Start begins consuming from the events channel.
	// The channel is closed when the poller stops.
	Start(ctx context.Context, eventsCh <-chan CheckpointEventsBatch) error
	// ProcessCheckpointEvents processes a batch of events from a single checkpoint.
	ProcessCheckpointEvents(ctx context.Context, batch CheckpointEventsBatch) error
	// AddEventSelector registers a new event selector to be matched against incoming events.
	AddEventSelector(ctx context.Context, selector *sui.EventFilterByMoveEventModule) error
	// GetEventSelectors returns a snapshot of all registered event selectors.
	// Used by ChainPoller to filter events.
	GetEventSelectors() []*sui.EventFilterByMoveEventModule
	// SetEventOffsetOverrides is deprecated; kept for backward compatibility.
	// Events are now checkpoint-ordered; this method logs a warning.
	SetEventOffsetOverrides(ctx context.Context, overrides map[string]sui.EventId) error
	Ready() error
	Close() error
}

// NewEventIndexer creates a new EventsIndexer instance.
func NewEventIndexer(
	db sqlutil.DataSource,
	log logger.Logger,
	eventConfigurations []*sui.EventFilterByMoveEventModule,
) EventsIndexerApi {
	dataStore := database.NewDBStore(db, log)
	namedLogger := logger.Named(log, "EventsIndexer")

	return &EventsIndexer{
		db:                  dataStore,
		logger:              namedLogger,
		eventConfigurations: eventConfigurations,
	}
}

// Start begins consuming events from the provided channel.
func (eIndexer *EventsIndexer) Start(ctx context.Context, eventsCh <-chan CheckpointEventsBatch) error {
	return eIndexer.starter.StartOnce("EventsIndexer", func() error {
		// Ensure database schema is set up
		syncCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		if err := eIndexer.db.EnsureSchema(syncCtx); err != nil {
			return fmt.Errorf("failed to ensure schema: %w", err)
		}

		eIndexer.wg.Go(func() {
			eIndexer.run(ctx, eventsCh)
		})

		return nil
	})
}

// run is the main consumer loop.
func (eIndexer *EventsIndexer) run(ctx context.Context, eventsCh <-chan CheckpointEventsBatch) {
	eIndexer.logger.Info("EventsIndexer starting")

	for {
		select {
		case <-ctx.Done():
			eIndexer.logger.Info("EventsIndexer stopping")
			return
		case batch, ok := <-eventsCh:
			if !ok {
				eIndexer.logger.Info("EventsIndexer channel closed, exiting")
				return
			}
			if err := eIndexer.ProcessCheckpointEvents(ctx, batch); err != nil {
				eIndexer.logger.Errorw("Failed to process checkpoint events",
					"sequence", batch.Checkpoint.SequenceNumber,
					"error", err)
			}
		}
	}
}

// ProcessCheckpointEvents processes a batch of events from a single checkpoint.
func (eIndexer *EventsIndexer) ProcessCheckpointEvents(ctx context.Context, batch CheckpointEventsBatch) error {
	if len(batch.Events) == 0 {
		return nil
	}

	// Group events by their selector for efficient batch processing
	eventsByHandle := make(map[string][]CheckpointEventItem)

	for _, item := range batch.Events {
		// Find matching selector
		selector := eIndexer.findMatchingSelector(item.Event)
		if selector == nil {
			eIndexer.logger.Debugw("Event does not match any selector, skipping",
				"package", item.Event.GetPackageId(),
				"module", item.Event.GetModule(),
				"type", item.Event.GetEventType())
			continue
		}

		handle := fmt.Sprintf("%s::%s::%s", selector.Package, selector.Module, selector.Event)
		eventsByHandle[handle] = append(eventsByHandle[handle], item)
	}

	// Process each event handle group
	for handle, items := range eventsByHandle {
		if err := eIndexer.processEventsForHandle(ctx, handle, items, batch.Checkpoint); err != nil {
			eIndexer.logger.Errorw("Failed to process events for handle",
				"handle", handle,
				"error", err)
			// Continue processing other handles
		}
	}

	return nil
}

// findMatchingSelector returns the first selector that matches the given event.
func (eIndexer *EventsIndexer) findMatchingSelector(event *suirpcv2.Event) *sui.EventFilterByMoveEventModule {
	if event == nil {
		eIndexer.logger.Warnw("Event is nil, skipping", "event", event)
		return nil
	}

	// Normalize the event type once so the address segment is comparable to the
	// selector types. See normalizeEventType for why only the address is lowercased.
	eventType := normalizeEventType(event.GetEventType())

	eIndexer.configMutex.RLock()
	selectors := make([]*sui.EventFilterByMoveEventModule, len(eIndexer.eventConfigurations))
	copy(selectors, eIndexer.eventConfigurations)
	eIndexer.configMutex.RUnlock()

	for _, selector := range selectors {
		expectedEventType := normalizeEventType(fmt.Sprintf("%s::%s::%s", selector.Package, selector.Module, selector.Event))
		if eventType == expectedEventType {
			return selector
		}
	}

	return nil
}

// processEventsForHandle processes all events for a specific event handle.
func (eIndexer *EventsIndexer) processEventsForHandle(
	ctx context.Context,
	handle string,
	items []CheckpointEventItem,
	meta CheckpointMeta,
) error {
	// Parse event handle: package::module::event
	parts := strings.Split(handle, "::")
	if len(parts) != 3 {
		return fmt.Errorf("invalid event handle format: %s", handle)
	}
	packageID := parts[0]

	// Build event records
	records := make([]database.EventRecord, 0, len(items))
	for _, item := range items {
		// Convert event JSON data using SDK's AsInterface()
		data := make(map[string]any)
		if jsonVal := item.Event.GetJson(); jsonVal != nil {
			if iface := jsonVal.AsInterface(); iface != nil {
				if mapData, ok := iface.(map[string]any); ok {
					data = mapData
				}
			}
		}

		// Normalize keys to camelCase
		if normalized := common.ConvertMapKeysToCamelCase(data); normalized != nil {
			if mapData, ok := normalized.(map[string]any); ok {
				data = mapData
			}
		}

		// Convert tx digest to hex format
		txDigestHex := item.TxDigest
		if !strings.HasPrefix(item.TxDigest, "0x") {
			if base58Bytes, err := base58.Decode(item.TxDigest); err == nil {
				txDigestHex = "0x" + hex.EncodeToString(base58Bytes)
			}
		}

		// Convert block hash to bytes
		blockHashBytes := []byte(meta.Digest)
		if decoded, err := base58.Decode(meta.Digest); err == nil {
			blockHashBytes = decoded
		}

		record := database.EventRecord{
			EventAccountAddress: packageID,
			EventHandle:         handle,
			EventOffset:         uint64(item.EventIndex),
			TxDigest:            txDigestHex,
			BlockVersion:        0,
			BlockHeight:         strconv.FormatUint(meta.SequenceNumber, 10),
			BlockHash:           blockHashBytes,
			BlockTimestamp:      meta.TimestampMs / 1000, // Convert ms to seconds
			Data:                data,
			IsSynthetic:         false,
		}
		records = append(records, record)

		eIndexer.logger.Debugw("Prepared event record to insert", "event", record)
	}

	// Batch insert with fallback
	if err := eIndexer.db.InsertEvents(ctx, records); err != nil {
		eIndexer.logger.Errorw("Batch insert failed, falling back to per-event insert",
			"error", err,
			"handle", handle)

		// Fallback: insert each record individually
		for _, record := range records {
			if err := eIndexer.db.InsertEvents(ctx, []database.EventRecord{record}); err != nil {
				eIndexer.logger.Errorw("Failed to insert single event, skipping",
					"error", err,
					"txDigest", record.TxDigest,
					"offset", record.EventOffset)
				continue
			}
		}
	}

	eIndexer.logger.Debugw("Inserted checkpoint events",
		"handle", handle,
		"count", len(records),
		"checkpoint", meta.SequenceNumber)

	return nil
}

// AddEventSelector registers a new event selector to be matched against incoming events.
func (eIndexer *EventsIndexer) AddEventSelector(ctx context.Context, selector *sui.EventFilterByMoveEventModule) error {
	if selector == nil {
		return fmt.Errorf("unspecified selector for AddEventSelector call")
	}

	if eIndexer.isEventSelectorAdded(*selector) {
		eIndexer.logger.Debugw("Event selector already registered, skipping",
			"package", selector.Package, "module", selector.Module, "event", selector.Event)
		return nil
	}

	eIndexer.configMutex.Lock()
	added := false
	if !eIndexer.isEventSelectorAddedLocked(*selector) {
		eIndexer.eventConfigurations = append(eIndexer.eventConfigurations, selector)
		added = true
	}
	total := len(eIndexer.eventConfigurations)
	eIndexer.configMutex.Unlock()

	if added {
		// This is the signal to watch: until at least one selector is registered, the ChainPoller's
		// filterEvents returns empty batches and the EventsIndexer never indexes anything (e.g. the
		// OnRamp CCIPMessageSent events the commit plugin needs for Sui-source merkle roots).
		eIndexer.logger.Infow("Registered event selector",
			"package", selector.Package,
			"module", selector.Module,
			"event", selector.Event,
			"totalSelectors", total)
	}

	return nil
}

// GetEventSelectors returns a snapshot of all registered event selectors.
func (eIndexer *EventsIndexer) GetEventSelectors() []*sui.EventFilterByMoveEventModule {
	eIndexer.configMutex.RLock()
	defer eIndexer.configMutex.RUnlock()

	selectors := make([]*sui.EventFilterByMoveEventModule, len(eIndexer.eventConfigurations))
	copy(selectors, eIndexer.eventConfigurations)
	return selectors
}

// SetEventOffsetOverrides is deprecated; events are now checkpoint-ordered.
func (eIndexer *EventsIndexer) SetEventOffsetOverrides(ctx context.Context, overrides map[string]sui.EventId) error {
	eIndexer.logger.Warn("SetEventOffsetOverrides is deprecated; events are now checkpoint-ordered")
	return nil
}

// isEventSelectorAdded checks if a specific event selector has already been added.
func (eIndexer *EventsIndexer) isEventSelectorAdded(eConfig sui.EventFilterByMoveEventModule) bool {
	eIndexer.configMutex.RLock()
	defer eIndexer.configMutex.RUnlock()
	return eIndexer.isEventSelectorAddedLocked(eConfig)
}

// isEventSelectorAddedLocked assumes the lock is already held.
func (eIndexer *EventsIndexer) isEventSelectorAddedLocked(eConfig sui.EventFilterByMoveEventModule) bool {
	for _, selector := range eIndexer.eventConfigurations {
		if selector.Package == eConfig.Package && selector.Module == eConfig.Module && selector.Event == eConfig.Event {
			return true
		}
	}
	return false
}

// Ready returns nil if the indexer has started successfully.
func (eIndexer *EventsIndexer) Ready() error {
	return eIndexer.starter.Ready()
}

// Close stops the events indexer.
func (eIndexer *EventsIndexer) Close() error {
	return eIndexer.starter.StopOnce("EventsIndexer", func() error {
		eIndexer.wg.Wait()
		return nil
	})
}
