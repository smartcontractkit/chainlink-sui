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

	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// ChainPollerAPI defines the interface for the ChainPoller.
// It fetches checkpoint data and fans it out to EventsIndexer and TransactionsIndexer via channels.
type ChainPollerAPI interface {
	Start(ctx context.Context) error
	EventsChannel() <-chan CheckpointEventsBatch
	TransactionsChannel() <-chan CheckpointTransactionsBatch
	// RescanRecent rewinds the poller so recently-processed checkpoints are re-scanned (used when a new
	// event selector is registered after the poller has already advanced past the events it matches).
	RescanRecent()
	Close() error
}

// SelectorProvider is a function that returns the current list of event selectors.
// This allows dynamic selector registration via AddEventSelector.
type SelectorProvider func() []*sui.EventFilterByMoveEventModule

// ChainPoller implements checkpoint-based polling for the Sui chain.
// It fetches checkpoints, extracts events and transactions, and sends them over channels
// to the indexers.
type ChainPoller struct {
	client           client.SuiPTBClient
	extendedClient   client.ExtendedPTBClient
	logger           logger.Logger
	config           sui.ChainPollerConfig
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
	cfg sui.ChainPollerConfig,
	selectorProvider SelectorProvider,
) *ChainPoller {
	bufferSize := cfg.ChannelBufferSize
	if bufferSize <= 0 {
		bufferSize = 16 // default
	}

	return &ChainPoller{
		client:           client,
		extendedClient:   asExtendedPTBClient(client),
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

		cp.wg.Go(func() {
			defer cp.closeChannels()
			cp.run(pollerCtx)
		})

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

	cp.logger.Infow(
		"Starting checkpoint polling",
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
	var startSeq uint64

	// If StartCheckpointSequence is configured, use it directly
	if cp.config.StartCheckpointSequence != nil {
		startSeq = *cp.config.StartCheckpointSequence
	} else {
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
				startSeq = latestSeq - count
			}
		} else {
			// Default: start from latest (no backfill)
			startSeq = latestSeq
		}
	}

	return cp.clampToProviderFloor(ctx, startSeq), nil
}

// clampToProviderFloor raises startSeq to the provider's lowest available checkpoint when history
// has been pruned below the configured backfill/start point.
func (cp *ChainPoller) clampToProviderFloor(ctx context.Context, startSeq uint64) uint64 {
	if cp.extendedClient == nil {
		cp.logger.Warnw(
			"Failed to get provider checkpoint availability, using requested start sequence",
			"startSequence", startSeq,
			"error", errors.New("client does not implement ExtendedPTBClient"),
		)
		return startSeq
	}

	info, err := cp.extendedClient.GetCheckpointAvailability(ctx)
	if err != nil {
		cp.logger.Warnw(
			"Failed to get provider checkpoint availability, using requested start sequence",
			"startSequence", startSeq,
			"error", err,
		)
		return startSeq
	}

	lowest := info.GetLowestAvailableCheckpoint()
	if lowest > 0 && startSeq < lowest {
		cp.logger.Warnw(
			"Start sequence below provider history floor, clamping",
			"requested", startSeq,
			"lowestAvailable", lowest,
		)
		return lowest
	}

	return startSeq
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

// catchUpConcurrency bounds the number of checkpoints fetched in parallel during
// catch-up. The per-checkpoint gRPC fetch is the bottleneck, so issuing several at
// once closes the gap to the tip faster. Commits stay in sequence order, so
// lastProcessed never advances past an unprocessed checkpoint. Conservative to
// avoid overloading the RPC node.
const catchUpConcurrency = 8

// catchUp processes checkpoints from startSeq to endSeq (inclusive). Checkpoints
// are fetched concurrently up to catchUpConcurrency at a time, then committed in
// sequence order. In-order commits preserve the guarantee that lastProcessed only
// advances past fully-processed checkpoints, so a failed checkpoint is retried on
// the next poll rather than skipped.
func (cp *ChainPoller) catchUp(ctx context.Context, startSeq, endSeq uint64) {
	startSeq = cp.clampToProviderFloor(ctx, startSeq)
	if startSeq > endSeq {
		return
	}

	// Canceling catchUpCtx on return stops the fetch stage so fetch goroutines
	// never outlive this call.
	catchUpCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	type fetchResult struct {
		seq  uint64
		data *client.CheckpointData
		err  error
	}

	results := make(chan fetchResult, catchUpConcurrency)

	// Fetch stage: issues GetCheckpointData calls in increasing sequence order,
	// up to catchUpConcurrency in flight. gRPC clients are safe for concurrent use.
	go func() {
		sem := make(chan struct{}, catchUpConcurrency)
		var wg sync.WaitGroup

		for seq := startSeq; seq <= endSeq; seq++ {
			select {
			case <-catchUpCtx.Done():
				wg.Wait()
				close(results)
				return
			case sem <- struct{}{}:
			}

			wg.Add(1)
			go func(seq uint64) {
				defer wg.Done()
				defer func() { <-sem }()

				data, err := cp.fetchCheckpoint(catchUpCtx, seq)
				select {
				case results <- fetchResult{seq: seq, data: data, err: err}:
				case <-catchUpCtx.Done():
				}
			}(seq)
		}

		wg.Wait()
		close(results)
	}()

	// Commit stage: drains fetch results and commits in sequence order, buffering
	// out-of-order results in pending until the next expected sequence arrives.
	pending := make(map[uint64]fetchResult, catchUpConcurrency)
	nextSeq := startSeq

	for r := range results {
		pending[r.seq] = r

		for {
			pr, ok := pending[nextSeq]
			if !ok {
				break
			}
			delete(pending, nextSeq)

			if pr.err != nil {
				if isCheckpointNotFound(pr.err) {
					if nextSeq == endSeq {
						cp.logger.Debugw(
							"Latest checkpoint not yet available",
							"sequence", nextSeq,
							"error", pr.err,
						)
						return
					}

					if lowest := cp.providerLowestAvailable(ctx); lowest > 0 && nextSeq < lowest {
						cp.logger.Errorw(
							"Checkpoint below provider history floor during catch-up, stopping; next poll resumes from the floor",
							"sequence", nextSeq,
							"lowestAvailable", lowest,
							"error", pr.err,
						)
						// Stop this catch-up rather than fetching and buffering the now-pruned
						// range [nextSeq, lowest). The fetch stage is canceled on return. The next
						// poll calls catchUp(lastProcessed+1, latest), and clampToProviderFloor
						// raises the start to the current floor, so no available checkpoint is
						// missed and pending cannot grow with the pruned range.
						return
					}

					cp.logger.Warnw(
						"Checkpoint not found during catch-up, will retry on next poll",
						"sequence", nextSeq,
						"error", pr.err,
					)
					// Don't advance lastProcessed - will retry on next poll
					return
				}

				cp.logger.Errorw(
					"Failed to process checkpoint, will retry on next poll",
					"sequence", nextSeq,
					"error", pr.err,
				)
				// Don't advance lastProcessed - will retry on next poll
				return
			}

			if err := cp.commitCheckpoint(ctx, nextSeq, pr.data); err != nil {
				cp.logger.Errorw(
					"Failed to process checkpoint, will retry on next poll",
					"sequence", nextSeq,
					"error", err,
				)
				// Don't advance lastProcessed - will retry on next poll
				return
			}
			nextSeq++
		}
	}
}

func (cp *ChainPoller) providerLowestAvailable(ctx context.Context) uint64 {
	if cp.extendedClient == nil {
		cp.logger.Warnw(
			"Failed to get provider checkpoint availability",
			"error", errors.New("client does not implement ExtendedPTBClient"),
		)
		return 0
	}

	info, err := cp.extendedClient.GetCheckpointAvailability(ctx)
	if err != nil {
		cp.logger.Warnw("Failed to get provider checkpoint availability", "error", err)
		return 0
	}
	return info.GetLowestAvailableCheckpoint()
}

func asExtendedPTBClient(suiClient client.SuiPTBClient) client.ExtendedPTBClient {
	ext, ok := suiClient.(client.ExtendedPTBClient)
	if !ok {
		return nil
	}
	return ext
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

// fetchCheckpoint fetches a single checkpoint's data. Safe to call concurrently:
// only the gRPC client is touched, which supports parallel calls.
func (cp *ChainPoller) fetchCheckpoint(ctx context.Context, seq uint64) (*client.CheckpointData, error) {
	ctx, cancel := context.WithTimeout(ctx, cp.config.SyncTimeout)
	defer cancel()

	cp.logger.Debugw("Fetching checkpoint", "sequence", seq)

	data, err := cp.client.GetCheckpointData(ctx, seq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch checkpoint %d: %w", seq, err)
	}

	if data.Checkpoint == nil {
		return nil, fmt.Errorf("checkpoint %d is nil", seq)
	}

	return data, nil
}

// commitCheckpoint builds event and transaction batches for an already-fetched
// checkpoint, sends them to the indexers, and advances lastProcessed. Called in
// sequence order from catchUp so lastProcessed only advances past committed
// checkpoints.
func (cp *ChainPoller) commitCheckpoint(ctx context.Context, seq uint64, data *client.CheckpointData) error {
	ctx, cancel := context.WithTimeout(ctx, cp.config.SyncTimeout)
	defer cancel()

	start := time.Now()

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

	cp.logger.Debugw(
		"Processing checkpoint",
		"sequence", seq,
		"digest", meta.Digest,
		"transactions", len(data.Transactions),
	)

	// Build events batch by filtering against selectors
	selectors := cp.selectorProvider()
	cp.logger.Debugw("current event selectors", "selectors", selectors)

	eventsBatch := cp.filterEvents(meta, data.Transactions, selectors)

	// Send events before transactions so event indexing is not blocked when the
	// transactions channel is full (e.g. TransactionsIndexer still starting up).
	if len(eventsBatch.Events) > 0 {
		select {
		case cp.eventsCh <- eventsBatch:
			cp.logger.Debugw(
				"Sent events batch",
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
	cp.logger.Debugw(
		"Checkpoint processed",
		"sequence", seq,
		"duration", elapsed,
	)

	return nil
}

// filterEvents extracts events matching the registered selectors from all transactions.
func (cp *ChainPoller) filterEvents(meta CheckpointMeta, transactions []*suirpcv2.ExecutedTransaction, selectors []*sui.EventFilterByMoveEventModule) CheckpointEventsBatch {
	batch := CheckpointEventsBatch{
		Checkpoint: meta,
		Events:     make([]CheckpointEventItem, 0),
	}

	if len(selectors) == 0 {
		return batch
	}

	// Precompute each selector's normalized type string once. Event types from the
	// chain are 0x-prefixed, so compare on the trimmed form instead of rebuilding
	// and trimming the selector string for every event/selector pair.
	selectorTypes := make([]string, len(selectors))
	for i, sel := range selectors {
		if sel == nil {
			continue
		}
		selectorTypes[i] = strings.TrimPrefix(fmt.Sprintf("%s::%s::%s", sel.Package, sel.Module, sel.Event), "0x")
	}

	for _, tx := range transactions {
		cp.logger.Debugw("Processing transaction for event filtering", "transaction", tx)

		txDigest := tx.GetDigest()
		if txDigest == "" {
			cp.logger.Warnf("Transaction digest is empty, skipping")
			continue
		}

		// Get events from transaction
		txEvents := tx.GetEvents()
		if txEvents == nil {
			cp.logger.Debugw("Transaction has no events, skipping", "transaction", txDigest)
			continue
		}

		for eventIdx, event := range txEvents.GetEvents() {
			if event == nil {
				cp.logger.Warnf("Event is nil in transaction %s, skipping", txDigest)
				continue
			}

			// Check if event matches any selector
			eventType := strings.TrimPrefix(event.GetEventType(), "0x")
			for i, sel := range selectors {
				if sel == nil {
					continue
				}

				cp.logger.Debugw(
					"Checking if event matches selector",
					"event", event.GetEventType(),
					"selector", selectorTypes[i],
				)

				if eventType == selectorTypes[i] {
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
func eventMatchesSelector(event *suirpcv2.Event, sel *sui.EventFilterByMoveEventModule) bool {
	if event == nil || sel == nil {
		return false
	}

	expectedEventType := fmt.Sprintf("%s::%s::%s", sel.Package, sel.Module, sel.Event)
	expectedEventType = strings.TrimPrefix(expectedEventType, "0x")
	eventType := strings.TrimPrefix(event.GetEventType(), "0x")

	return expectedEventType == eventType
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
