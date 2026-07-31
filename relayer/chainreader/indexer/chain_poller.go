package indexer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

const (
	// defaultChannelBufferSize is used when the config does not set ChannelBufferSize.
	defaultChannelBufferSize = 64
	// defaultMaxConcurrentWorkers is the default number of goroutines processing checkpoint chunks.
	defaultMaxConcurrentWorkers = 64
	// defaultCatchupChunkSize is the default number of checkpoints per catch-up chunk.
	defaultCatchupChunkSize = 12
	// defaultMaxStalledTicksBeforeSkip is how many consecutive stalled poll ticks (no watermark
	// progress while work is outstanding) are tolerated before the stuck checkpoint is skipped.
	defaultMaxStalledTicksBeforeSkip = 30
	// minStalledTicksBeforeRetry gives in-flight workers at least one full tick before the run
	// loop re-fetches the gap at the watermark; retrying too eagerly re-emits checkpoints a
	// worker had just delivered (harmless duplicates, but wasted RPC calls).
	minStalledTicksBeforeRetry = 2
	// chunkQueueSizeMultiplier sizes the chunk queue relative to the worker count.
	chunkQueueSizeMultiplier = 4
	// defaultRescanRewindCheckpoints is how far RescanRecent rewinds when no explicit backfill window is set.
	defaultRescanRewindCheckpoints uint64 = 100
	// minStartSequence is a value used as a fallback to avoid using 0 and doing a full chain rescan
	// when the required configs are missing and the RPC node used is a full archive node
	minStartSequence uint64 = 299_000_000
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

// CheckpointCursorStore persists the last checkpoint sequence that was emitted in order, so the
// poller can resume where it left off after a restart. Implemented by *database.DBStore.
type CheckpointCursorStore interface {
	GetCheckpointCursor(ctx context.Context, id string) (seq uint64, found bool, err error)
	UpsertCheckpointCursor(ctx context.Context, id string, seq uint64) error
}

// chunkRange is an inclusive range of checkpoint sequences handed to a catch-up worker.
type chunkRange struct {
	start, end uint64
}

// PollerOption configures optional ChainPoller behavior.
type PollerOption func(*ChainPoller)

// WithWorkerPool sets the number of concurrent catch-up workers and the checkpoints-per-chunk
// size. Non-positive values keep the defaults.
func WithWorkerPool(workers, chunkSize int) PollerOption {
	return func(cp *ChainPoller) {
		if workers > 0 {
			cp.workers = workers
		}
		if chunkSize > 0 {
			cp.chunkSize = chunkSize
		}
	}
}

// WithCursorStore enables checkpoint-cursor persistence so the poller resumes from where it left
// off after a restart. id namespaces the cursor row (e.g. "chain_poller").
func WithCursorStore(store CheckpointCursorStore, id string) PollerOption {
	return func(cp *ChainPoller) {
		cp.cursorStore = store
		cp.cursorID = id
	}
}

// ChainPoller implements checkpoint-based polling for the Sui chain.
// It fetches checkpoints concurrently in bounded chunks, reorders the results via a sequencer,
// and sends events and transactions strictly in checkpoint order over channels to the indexers.
type ChainPoller struct {
	client           client.SuiPTBClient
	extendedClient   client.ExtendedPTBClient
	logger           logger.Logger
	config           sui.ChainPollerConfig
	eventsCh         chan CheckpointEventsBatch
	transactionsCh   chan CheckpointTransactionsBatch
	selectorProvider SelectorProvider

	bufferSize      int
	workers         int
	chunkSize       int
	maxPending      int // reorder-buffer capacity: workers * chunkSize
	maxStalledTicks int

	sequencer   *checkpointSequencer
	chunkCh     chan chunkRange
	cursorStore CheckpointCursorStore
	cursorID    string

	starter services.StateMachine
	wg      sync.WaitGroup
	cancel  context.CancelFunc

	// goroutineCount gauges catch-ups currently in flight (workers + the stall retry).
	goroutineCount atomic.Int64
	// retryInFlight ensures at most one background stall retry runs at a time.
	retryInFlight atomic.Bool
}

// NewChainPoller creates a new ChainPoller instance.
func NewChainPoller(
	client client.SuiPTBClient,
	log logger.Logger,
	cfg sui.ChainPollerConfig,
	selectorProvider SelectorProvider,
	opts ...PollerOption,
) *ChainPoller {
	bufferSize := cfg.ChannelBufferSize
	if bufferSize <= 0 {
		bufferSize = defaultChannelBufferSize
	}

	cp := &ChainPoller{
		client:           client,
		extendedClient:   asExtendedPTBClient(client),
		logger:           logger.Named(log, "ChainPoller"),
		config:           cfg,
		eventsCh:         make(chan CheckpointEventsBatch, bufferSize),
		transactionsCh:   make(chan CheckpointTransactionsBatch, bufferSize),
		selectorProvider: selectorProvider,
		bufferSize:       bufferSize,
		workers:          defaultMaxConcurrentWorkers,
		chunkSize:        defaultCatchupChunkSize,
		maxStalledTicks:  defaultMaxStalledTicksBeforeSkip,
	}

	for _, opt := range opts {
		opt(cp)
	}

	cp.maxPending = cp.workers * cp.chunkSize
	cp.chunkCh = make(chan chunkRange, cp.workers*chunkQueueSizeMultiplier)
	// The watermark is positioned in run() once the start sequence is known.
	cp.sequencer = newCheckpointSequencer(0, cp.maxPending, cp.eventsCh, cp.transactionsCh, cp.logger)

	return cp
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

// run is the main polling loop: it dispatches checkpoint chunks to the worker pool, retries
// stalled gaps, and persists the checkpoint cursor as the sequencer watermark advances.
func (cp *ChainPoller) run(ctx context.Context) {
	cp.logger.Info("ChainPoller starting")

	// Compute start sequence
	startSeq, err := cp.computeStartSequence(ctx)
	if err != nil {
		cp.logger.Errorw("Failed to compute start sequence", "error", err)
		// use a reasonable fallback value
		startSeq = minStartSequence
	}

	cp.sequencer.setStart(startSeq)

	cp.logger.Infow(
		"Starting checkpoint polling",
		"startSequence", startSeq,
		"backfillCount", cp.config.BackfillCheckpointCount,
		"startCheckpoint", cp.config.StartCheckpointSequence,
		"workers", cp.workers,
		"chunkSize", cp.chunkSize,
	)

	// Worker pool consuming checkpoint chunks. Workers must be done before the output channels
	// are closed (Start's defer), so wait for them on the way out.
	var workerWg sync.WaitGroup
	for range cp.workers {
		workerWg.Go(func() {
			cp.runWorker(ctx)
		})
	}
	defer func() {
		workerWg.Wait()
		cp.sequencer.close()
	}()

	// Live polling
	ticker := time.NewTicker(cp.config.PollingInterval)
	defer ticker.Stop()

	nextToDispatch := startSeq
	lastWatermark := startSeq
	stalledTicks := 0
	var persistedCursor uint64
	persistedCursorValid := false

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

			watermark := cp.sequencer.NextExpected()

			// 1. Dispatch new checkpoints in fixed-size chunks, bounded to the reorder-buffer
			// window ahead of the watermark. Without the bound, workers commit to chunks far
			// past what the sequencer can buffer, all park in Submit, and no idle capacity is
			// ever available for the sequences the watermark actually needs (head-of-line
			// blocking that collapses catch-up throughput to a single sequential fetcher).
			windowEnd := watermark + cp.pendingWindow() - 1
			nextToDispatch = cp.dispatchChunks(nextToDispatch, min(latestSeq, windowEnd))

			// retry stalled gaps: the watermark not moving while dispatched work is still
			// outstanding means a chunk failed (RPC error) or the tip checkpoint was not yet
			// available. Re-process the missing range in the background (at most one retry in
			// flight) — submitting the next-expected sequence is always accepted by the
			// sequencer, so this also unblocks workers waiting on a full reorder buffer.
			// A watermark whose result is already buffered is NOT a gap — it is mid-emission
			// behind downstream backpressure — so it must neither be retried nor skipped.
			if watermark < nextToDispatch && watermark == lastWatermark && !cp.sequencer.isPending(watermark) {
				stalledTicks++
				if stalledTicks >= minStalledTicksBeforeRetry && cp.retryInFlight.CompareAndSwap(false, true) {
					retryEnd := min(watermark+cp.chunkSpan()-1, nextToDispatch-1)
					workerWg.Go(func() {
						defer cp.retryInFlight.Store(false)
						cp.goroutineCount.Add(1)
						defer cp.goroutineCount.Add(-1)
						cp.catchUp(ctx, watermark, retryEnd)
					})
				}

				if stalledTicks >= cp.maxStalledTicks && cp.sequencer.NextExpected() == watermark {
					if lowest := cp.providerLowestAvailable(ctx); lowest > watermark {
						cp.logger.Errorw(
							"Watermark below provider history floor, skipping pruned checkpoints",
							"watermark", watermark,
							"lowestAvailable", lowest,
						)
						cp.sequencer.SkipTo(ctx, lowest)
					} else {
						cp.logger.Errorw(
							"Checkpoint failed to process after repeated retries, skipping",
							"sequence", watermark,
							"stalledTicks", stalledTicks,
						)
						cp.sequencer.SkipTo(ctx, watermark+1)
					}
					stalledTicks = 0
				}
			} else {
				stalledTicks = 0
			}
			lastWatermark = cp.sequencer.NextExpected()

			// persist the cursor
			// everything below the watermark has been emitted in order
			if cp.cursorStore != nil && lastWatermark > startSeq {
				cursorSeq := lastWatermark - 1
				if !persistedCursorValid || cursorSeq > persistedCursor {
					if err := cp.cursorStore.UpsertCheckpointCursor(ctx, cp.cursorID, cursorSeq); err != nil {
						cp.logger.Warnw("Failed to persist checkpoint cursor", "sequence", cursorSeq, "error", err)
					} else {
						persistedCursor = cursorSeq
						persistedCursorValid = true
					}
				}
			}

			cp.logger.Debugw("Catch-up state",
				"latest", latestSeq,
				"nextToDispatch", nextToDispatch,
				"watermark", lastWatermark,
				"activeCatchups", cp.goroutineCount.Load(),
			)
		}
	}
}

// chunkSpan returns the chunk size as uint64 for sequence arithmetic.
func (cp *ChainPoller) chunkSpan() uint64 {
	//nolint:gosec // G115: chunkSize is a small positive int
	return uint64(cp.chunkSize)
}

// pendingWindow returns the reorder-buffer capacity as uint64 for sequence arithmetic.
func (cp *ChainPoller) pendingWindow() uint64 {
	//nolint:gosec // G115: maxPending is a small positive int
	return uint64(cp.maxPending)
}

// runWorker consumes checkpoint chunks until the poller context is canceled.
func (cp *ChainPoller) runWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-cp.chunkCh:
			cp.goroutineCount.Add(1)
			cp.catchUp(ctx, r.start, r.end)
			cp.goroutineCount.Add(-1)
		}
	}
}

// dispatchChunks enqueues [from, latest] as chunkSize-sized ranges onto the chunk queue without
// blocking. It returns the new dispatch frontier: the first sequence not yet enqueued.
func (cp *ChainPoller) dispatchChunks(from, latest uint64) uint64 {
	for from <= latest {
		end := min(from+cp.chunkSpan()-1, latest)
		select {
		case cp.chunkCh <- chunkRange{start: from, end: end}:
			from = end + 1
		default:
			// Queue full; the remaining range is dispatched on a later tick.
			return from
		}
	}

	return from
}

// computeStartSequence calculates the starting checkpoint sequence number.
func (cp *ChainPoller) computeStartSequence(ctx context.Context) (uint64, error) {
	var startSeq uint64

	if resume, ok := cp.loadResumeSequence(ctx); ok {
		startSeq = resume
		// An explicit StartCheckpointSequence may skip ahead of the persisted cursor, but never
		// rewinds below it; delete the cursor row to force a rewind.
		if cp.config.StartCheckpointSequence != nil && *cp.config.StartCheckpointSequence > startSeq {
			startSeq = *cp.config.StartCheckpointSequence
		}
	} else if cp.config.StartCheckpointSequence != nil {
		// If StartCheckpointSequence is configured, use it directly
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

// loadResumeSequence reads the persisted checkpoint cursor and converts it into a resume start
// sequence, rewound by the in-flight window (reorder buffer + channel buffers) to cover batches
// that were emitted but possibly not yet inserted when the process stopped. Inserts are
// idempotent, so the overlap only costs a little re-scanning.
func (cp *ChainPoller) loadResumeSequence(ctx context.Context) (uint64, bool) {
	if cp.cursorStore == nil {
		return 0, false
	}

	cursorSeq, found, err := cp.cursorStore.GetCheckpointCursor(ctx, cp.cursorID)
	if err != nil {
		cp.logger.Warnw("Failed to load checkpoint cursor, falling back to configured start", "error", err)
		return 0, false
	}
	if !found {
		return 0, false
	}

	//nolint:gosec // G115: maxPending and bufferSize are small positive ints
	overlap := uint64(cp.maxPending + cp.bufferSize)
	resume := cursorSeq + 1
	if resume > overlap {
		resume -= overlap
	} else {
		resume = 0
	}

	cp.logger.Infow(
		"Resuming from persisted checkpoint cursor",
		"cursor", cursorSeq,
		"resumeStart", resume,
		"overlap", overlap,
	)

	return resume, true
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

// RescanRecent enqueues the most recent N checkpoints for re-scanning. Results below the
// sequencer watermark bypass ordering and are re-emitted directly; inserts are idempotent.
func (cp *ChainPoller) RescanRecent() {
	latestCheckpoint, err := cp.getLatestCheckpointSequence(context.Background())
	if err != nil {
		cp.logger.Errorw("Failed to get latest checkpoint", "error", err)
		return
	}

	startSeq := uint64(0)
	if latestCheckpoint > defaultRescanRewindCheckpoints {
		startSeq = latestCheckpoint - defaultRescanRewindCheckpoints
	}

	for from := startSeq; from <= latestCheckpoint; {
		end := min(from+cp.chunkSpan()-1, latestCheckpoint)
		select {
		case cp.chunkCh <- chunkRange{start: from, end: end}:
			from = end + 1
		default:
			cp.logger.Warnw(
				"Chunk queue full, dropping remainder of rescan range",
				"droppedFrom", from,
				"droppedTo", latestCheckpoint,
			)
			return
		}
	}
}

// catchUp processes checkpoints from startSeq to endSeq (inclusive).
func (cp *ChainPoller) catchUp(ctx context.Context, startSeq, endSeq uint64) {
	startSeq = cp.clampToProviderFloor(ctx, startSeq)
	if startSeq > endSeq {
		return
	}

	cp.logger.Infow("Starting checkpoint catchup goroutine", "startSequence", startSeq, "endSequence", endSeq)

	for seq := startSeq; seq <= endSeq; seq++ {
		select {
		case <-ctx.Done():
			cp.logger.Infow("Catch-up interrupted", "atSequence", seq)
			return
		default:
			if err := cp.processCheckpoint(ctx, seq); err != nil {
				if isCheckpointNotFound(err) {
					if seq == endSeq {
						cp.logger.Debugw(
							"Latest checkpoint not yet available",
							"sequence", seq,
							"error", err,
						)
						return
					}

					if lowest := cp.providerLowestAvailable(ctx); lowest > 0 && seq < lowest {
						cp.logger.Errorw(
							"Checkpoint below provider history floor during catch-up, jumping to lowest available",
							"sequence", seq,
							"lowestAvailable", lowest,
							"error", err,
						)
						seq = lowest - 1
						continue
					}

					cp.logger.Warnw(
						"Checkpoint not found during catch-up, will retry on next poll",
						"sequence", seq,
						"error", err,
					)
					return
				}
				cp.logger.Warnw(
					"Failed to process checkpoint, must rescan",
					"sequence", seq,
					"error", err,
				)
				// TODO: add to queue of retriable checkpoints
				return
			}
		}
	}

	cp.logger.Debugw("Checkpoint catchup goroutine completed",
		"start", startSeq,
		"end", endSeq,
	)
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

// processCheckpoint fetches a single checkpoint and submits the result to the sequencer, which
// emits batches to the output channels strictly in checkpoint order.
func (cp *ChainPoller) processCheckpoint(ctx context.Context, seq uint64) error {
	// Skip the fetch when the result is already buffered awaiting emission — a duplicate
	// submit would be dropped anyway, so re-fetching (e.g. from a stall retry overlapping an
	// in-flight chunk) is wasted RPC. Below-watermark sequences are still fetched so
	// RescanRecent can re-emit them.
	if cp.sequencer.isPending(seq) {
		cp.logger.Debugw("Checkpoint already buffered, skipping fetch", "sequence", seq)
		return nil
	}

	cp.logger.Debugw("Processing checkpoint", "sequence", seq)

	start := time.Now()

	// Only the fetch is bounded by SyncTimeout. The sequencer submission below blocks on
	// ordering/backpressure and is bounded by the poller context instead.
	fetchCtx, cancel := context.WithTimeout(ctx, cp.config.SyncTimeout)
	data, err := cp.client.GetCheckpointData(fetchCtx, seq)
	cancel()
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

	if err := cp.sequencer.Submit(ctx, checkpointResult{
		seq:    seq,
		events: eventsBatch,
		txs: CheckpointTransactionsBatch{
			Checkpoint:   meta,
			Transactions: data.Transactions,
		},
	}); err != nil {
		return fmt.Errorf("failed to submit checkpoint %d: %w", seq, err)
	}

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

	for txIdx, tx := range transactions {
		cp.logger.Debugw("Processing transaction for event filtering", "transaction", tx.GetDigest())

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
			for _, sel := range selectors {
				cp.logger.Debugw(
					"Checking if event matches selector",
					"event", event.GetEventType(),
					"selector", fmt.Sprintf("%s::%s::%s", sel.Package, sel.Module, sel.Event),
				)

				if eventMatchesSelector(event, sel) {
					cp.logger.Debugw(
						"Event does match selector",
						"event", event.GetEventType(),
						"rpcPackageId", event.GetPackageId(),
						"selector", fmt.Sprintf("%s::%s::%s", sel.Package, sel.Module, sel.Event),
					)

					item := CheckpointEventItem{
						Event:      event,
						TxDigest:   txDigest,
						TxIndex:    uint32(txIdx),
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
