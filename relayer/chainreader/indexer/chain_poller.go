package indexer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// ChainPollerAPI defines the interface for the ChainPoller.
// It fetches checkpoint data and fans it out to EventsIndexer and TransactionsIndexer via channels.
type ChainPollerAPI interface {
	Start(ctx context.Context) error
	EventsChannel() <-chan CheckpointEventsBatch
	TransactionsChannel() <-chan CheckpointTransactionsBatch
	Close() error
}

// SelectorProvider is a function that returns the current list of event selectors.
// This allows dynamic selector registration via AddEventSelector.
type SelectorProvider func() []*client.EventSelector

// ChainPoller implements checkpoint-based polling for the Sui chain.
// It fetches checkpoints, extracts events and transactions, and sends them over channels
// to the indexers.
type ChainPoller struct {
	client           client.SuiPTBClient
	logger           logger.Logger
	config           config.ChainPollerConfig
	eventsCh         chan CheckpointEventsBatch
	transactionsCh   chan CheckpointTransactionsBatch
	selectorProvider SelectorProvider

	starter       services.StateMachine
	lastProcessed uint64 // in-memory only
	mu            sync.RWMutex
	wg            sync.WaitGroup
	cancel        context.CancelFunc
}

// NewChainPoller creates a new ChainPoller instance.
func NewChainPoller(
	client client.SuiPTBClient,
	log logger.Logger,
	cfg config.ChainPollerConfig,
	selectorProvider SelectorProvider,
) *ChainPoller {
	bufferSize := cfg.ChannelBufferSize
	if bufferSize <= 0 {
		bufferSize = 16 // default
	}

	return &ChainPoller{
		client:           client,
		logger:           logger.Named(log, "ChainPoller"),
		config:           cfg,
		eventsCh:         make(chan CheckpointEventsBatch, bufferSize),
		transactionsCh:   make(chan CheckpointTransactionsBatch, bufferSize),
		selectorProvider: selectorProvider,
	}
}

// Start begins the polling loop.
func (cp *ChainPoller) Start(ctx context.Context) error {
	return cp.starter.StartOnce("ChainPoller", func() error {
		//nolint:gosec // G118: cancel is invoked from Close()
		pollerCtx, cancel := context.WithCancel(ctx)
		cp.cancel = cancel

		cp.wg.Add(1)
		go func() {
			defer cp.wg.Done()
			defer cp.closeChannels()
			cp.run(pollerCtx)
		}()

		return nil
	})
}

// EventsChannel returns the channel for checkpoint event batches.
func (cp *ChainPoller) EventsChannel() <-chan CheckpointEventsBatch {
	return cp.eventsCh
}

// TransactionsChannel returns the channel for checkpoint transaction batches.
func (cp *ChainPoller) TransactionsChannel() <-chan CheckpointTransactionsBatch {
	return cp.transactionsCh
}

// Close stops the poller and closes channels.
func (cp *ChainPoller) Close() error {
	return cp.starter.StopOnce("ChainPoller", func() error {
		if cp.cancel != nil {
			cp.cancel()
		}
		cp.wg.Wait()
		return nil
	})
}

// closeChannels closes both output channels.
func (cp *ChainPoller) closeChannels() {
	close(cp.eventsCh)
	close(cp.transactionsCh)
}

// run is the main polling loop.
func (cp *ChainPoller) run(ctx context.Context) {
	cp.logger.Info("ChainPoller starting")

	// Compute start sequence
	startSeq, err := cp.computeStartSequence(ctx)
	if err != nil {
		cp.logger.Errorw("Failed to comnpute start sequence", "error", err)
		// Continue with startSeq = 0 as fallback
		startSeq = 0
	}

	cp.logger.Infow("Starting checkpoint polling",
		"startSequence", startSeq,
		"backfillCount", cp.config.BackfillCheckpointCount,
		"startCheckpoint", cp.config.StartCheckpointSequence,
	)

	// Initial catch-up loop
	latestSeq, err := cp.getLatestCheckpointSequence(ctx)
	if err != nil {
		cp.logger.Errorw("Failed to get latest checkpoint", "error", err)
		// Continue anyway, will retry in polling loop
	} else {
		cp.logger.Infow("Catch-up phase", "from", startSeq, "to", latestSeq)
		cp.catchUp(ctx, startSeq, latestSeq)
	}

	// Live polling
	ticker := time.NewTicker(cp.config.PollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			cp.logger.Info("ChainPoller stopping")
			return
		case <-ticker.C:
			cp.logger.Debugw("Polling for latest checkpoint")
			latestSeq, err := cp.getLatestCheckpointSequence(ctx)
			if err != nil {
				cp.logger.Warnw("Failed to get latest checkpoint", "error", err)
				continue
			}

			cp.mu.RLock()
			lastProcessed := cp.lastProcessed
			cp.mu.RUnlock()

			if latestSeq > lastProcessed {
				cp.catchUp(ctx, lastProcessed+1, latestSeq)
			}
		}
	}
}

// computeStartSequence calculates the starting checkpoint sequence number.
func (cp *ChainPoller) computeStartSequence(ctx context.Context) (uint64, error) {
	// If StartCheckpointSequence is configured, use it directly
	if cp.config.StartCheckpointSequence != nil {
		return *cp.config.StartCheckpointSequence, nil
	}

	// Get the latest checkpoint
	latestSeq, err := cp.getLatestCheckpointSequence(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest checkpoint: %w", err)
	}

	cp.logger.Infow("Latest checkpoint fetched in chain poller", "sequence", latestSeq)

	// If BackfillCheckpointCount is configured, start from latest - N
	if cp.config.BackfillCheckpointCount != nil {
		count := *cp.config.BackfillCheckpointCount
		if latestSeq > count {
			return latestSeq - count, nil
		}
		// If latest < count, start from 0
		return 0, nil
	}

	// Default: start from latest (no backfill)
	return latestSeq, nil
}

// getLatestCheckpointSequence fetches the latest checkpoint and returns its sequence number.
func (cp *ChainPoller) getLatestCheckpointSequence(ctx context.Context) (uint64, error) {
	checkpoint, err := cp.client.GetLatestCheckpoint(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest checkpoint: %w", err)
	}

	seq := checkpoint.GetSequenceNumber()
	return seq, nil
}

// defaultRescanRewindCheckpoints is how far RescanRecent rewinds when no explicit backfill window is set.
const defaultRescanRewindCheckpoints uint64 = 100

// RescanRecent rewinds lastProcessed so the most recent checkpoints are re-processed on the next poll. It
// is wired to fire when a new event selector is registered after the poller has advanced: a selector for
// an event (e.g. the OnRamp CCIPMessageSent) only registers once its contract is discovered/bound, by
// which point the poller may have already processed — and, lacking the selector, discarded — the
// checkpoints carrying that event. Event inserts are idempotent (ON CONFLICT DO NOTHING), so re-scanning
// is safe. The rewind span matches the configured backfill window.
func (cp *ChainPoller) RescanRecent() {
	rewind := defaultRescanRewindCheckpoints
	if cp.config.BackfillCheckpointCount != nil && *cp.config.BackfillCheckpointCount > 0 {
		rewind = *cp.config.BackfillCheckpointCount
	}

	cp.mu.Lock()
	prev := cp.lastProcessed
	if cp.lastProcessed > rewind {
		cp.lastProcessed -= rewind
	} else {
		cp.lastProcessed = 0
	}
	curr := cp.lastProcessed
	cp.mu.Unlock()

	if curr != prev {
		cp.logger.Infow("Rewinding poller to re-scan for newly registered event selector",
			"from", prev, "to", curr, "rewind", rewind)
	}
}

// catchUp processes checkpoints from startSeq to endSeq (inclusive).
func (cp *ChainPoller) catchUp(ctx context.Context, startSeq, endSeq uint64) {
	for seq := startSeq; seq <= endSeq; seq++ {
		select {
		case <-ctx.Done():
			cp.logger.Infow("Catch-up interrupted", "atSequence", seq)
			return
		default:
			if err := cp.processCheckpoint(ctx, seq); err != nil {
				if isCheckpointNotFound(err) {
					if seq == endSeq {
						cp.logger.Debugw("Latest checkpoint not yet available",
							"sequence", seq,
							"error", err,
						)
						return
					}
					cp.logger.Warnw("Checkpoint not found during catch-up, skipping",
						"sequence", seq,
						"error", err,
					)
					continue
				}
				cp.logger.Errorw("Failed to process checkpoint, will retry on next poll",
					"sequence", seq,
					"error", err,
				)
				// Don't advance lastProcessed - will retry on next poll
				return
			}
		}
	}
}

func isCheckpointNotFound(err error) bool {
	for err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return true
		}
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

// processCheckpoint fetches and processes a single checkpoint.
func (cp *ChainPoller) processCheckpoint(ctx context.Context, seq uint64) error {
	ctx, cancel := context.WithTimeout(ctx, cp.config.SyncTimeout)
	defer cancel()

	cp.logger.Debugw("Processing checkpoint", "sequence", seq)

	start := time.Now()

	// Fetch checkpoint data
	data, err := cp.client.GetCheckpointData(ctx, seq)
	if err != nil {
		return fmt.Errorf("failed to fetch checkpoint %d: %w", seq, err)
	}

	if data.Checkpoint == nil {
		return fmt.Errorf("checkpoint %d is nil", seq)
	}

	checkpoint := data.Checkpoint
	// Convert protobuf timestamp (seconds) to milliseconds
	timestampMs := uint64(0)
	if ts := checkpoint.GetSummary().GetTimestamp(); ts != nil {
		secs := ts.GetSeconds()
		if secs < 0 {
			return fmt.Errorf("checkpoint %d has negative timestamp", seq)
		}

		timestampMs = uint64(secs) * 1000
	}

	meta := CheckpointMeta{
		SequenceNumber: checkpoint.GetSequenceNumber(),
		Digest:         checkpoint.GetDigest(),
		TimestampMs:    timestampMs,
	}

	cp.logger.Debugw("Processing checkpoint",
		"sequence", seq,
		"digest", meta.Digest,
		"transactions", len(data.Transactions),
	)

	// Build events batch by filtering against selectors
	selectors := cp.selectorProvider()
	eventsBatch := cp.filterEvents(meta, data.Transactions, selectors)

	// Send events before transactions so event indexing is not blocked when the
	// transactions channel is full (e.g. TransactionsIndexer still starting up).
	if len(eventsBatch.Events) > 0 {
		select {
		case cp.eventsCh <- eventsBatch:
			cp.logger.Debugw("Sent events batch",
				"sequence", seq,
				"eventCount", len(eventsBatch.Events),
			)
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Build transaction batch
	txBatch := CheckpointTransactionsBatch{
		Checkpoint:   meta,
		Transactions: data.Transactions,
	}

	select {
	case cp.transactionsCh <- txBatch:
	case <-ctx.Done():
		return ctx.Err()
	}

	// Update last processed
	cp.mu.Lock()
	cp.lastProcessed = seq
	cp.mu.Unlock()

	elapsed := time.Since(start)
	cp.logger.Debugw("Checkpoint processed",
		"sequence", seq,
		"duration", elapsed,
	)

	return nil
}

// filterEvents extracts events matching the registered selectors from all transactions.
func (cp *ChainPoller) filterEvents(meta CheckpointMeta, transactions []*suirpcv2.ExecutedTransaction, selectors []*client.EventSelector) CheckpointEventsBatch {
	batch := CheckpointEventsBatch{
		Checkpoint: meta,
		Events:     make([]CheckpointEventItem, 0),
	}

	if len(selectors) == 0 {
		return batch
	}

	// Build a map for faster lookup (key: "package::module::event")
	selectorMap := make(map[string]*client.EventSelector)
	for _, sel := range selectors {
		key := fmt.Sprintf("%s::%s::%s", sel.Package, sel.Module, sel.Event)
		selectorMap[key] = sel
	}

	for _, tx := range transactions {
		txDigest := tx.GetDigest()
		if txDigest == "" {
			continue
		}

		// Get events from transaction
		txEvents := tx.GetEvents()
		if txEvents == nil {
			continue
		}

		for eventIdx, event := range txEvents.GetEvents() {
			if event == nil {
				continue
			}

			// Check if event matches any selector
			for _, sel := range selectors {
				if eventMatchesSelector(event, sel) {
					item := CheckpointEventItem{
						Event:      event,
						TxDigest:   txDigest,
						EventIndex: uint32(eventIdx),
					}
					batch.Events = append(batch.Events, item)
					break // Matched this selector, move to next event
				}
			}
		}
	}

	return batch
}

// eventMatchesSelector checks if an event matches a given selector.
func eventMatchesSelector(event *suirpcv2.Event, sel *client.EventSelector) bool {
	if event == nil || sel == nil {
		return false
	}

	// Check package ID (handle potential 0x prefix differences)
	eventPackage := strings.TrimPrefix(event.GetPackageId(), "0x")
	selectorPackage := strings.TrimPrefix(sel.Package, "0x")
	if eventPackage != selectorPackage {
		return false
	}

	// Check module
	if event.GetModule() != sel.Module {
		return false
	}

	// Check event type - EventType in v2 is fully qualified like "0x123::module::EventName"
	// We extract just the event name for comparison
	eventTypeName := event.GetEventType()
	if parts := strings.Split(eventTypeName, "::"); len(parts) >= 3 {
		eventTypeName = parts[2]
	}
	if eventTypeName != sel.Event {
		return false
	}

	return true
}

// Ready returns nil if the poller has started successfully.
func (cp *ChainPoller) Ready() error {
	return cp.starter.Ready()
}

// HealthReport returns the health status of the poller.
func (cp *ChainPoller) HealthReport() map[string]error {
	return map[string]error{
		"ChainPoller": cp.starter.Healthy(),
	}
}
