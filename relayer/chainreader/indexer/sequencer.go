package indexer

import (
	"context"
	"errors"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var errSequencerClosed = errors.New("checkpoint sequencer is closed")

// checkpointResult is the outcome of processing a single checkpoint, handed to the
// sequencer by concurrent catch-up workers.
type checkpointResult struct {
	seq    uint64
	events CheckpointEventsBatch       // may be empty; only emitted when it has events
	txs    CheckpointTransactionsBatch // always emitted
}

// checkpointSequencer is a reorder buffer between concurrent checkpoint workers and the
// poller's output channels. Workers submit results in any interleaving; the sequencer emits
// them to the events and transactions channels strictly in ascending sequence order, so the
// downstream indexers (single consumers) insert into the database sequentially by checkpoint.
//
// Liveness invariant: each worker must submit its own range in ascending order (catchUp does).
// The next-expected sequence is then always submitted by a worker that is not blocked on the
// full buffer — a worker whose chunk contains the watermark has already submitted everything
// below it — so draining always progresses. Arbitrary per-sequence interleaving across workers
// without that per-worker ordering could deadlock a full buffer.
type checkpointSequencer struct {
	mu   sync.Mutex
	cond *sync.Cond

	// nextSeq is the watermark: the next sequence expected to be emitted. Everything below
	// it has already been emitted in order.
	nextSeq uint64
	// pending holds out-of-order results waiting for the watermark to reach them.
	pending map[uint64]*checkpointResult
	// maxPending bounds the reorder buffer; Submit blocks when full (except for the
	// next-expected sequence, which must always be accepted so draining can progress).
	maxPending int
	closed     bool

	eventsCh chan<- CheckpointEventsBatch
	txCh     chan<- CheckpointTransactionsBatch
	lggr     logger.Logger
}

func newCheckpointSequencer(
	startSeq uint64,
	maxPending int,
	eventsCh chan<- CheckpointEventsBatch,
	txCh chan<- CheckpointTransactionsBatch,
	lggr logger.Logger,
) *checkpointSequencer {
	if maxPending <= 0 {
		maxPending = 1
	}
	s := &checkpointSequencer{
		nextSeq:    startSeq,
		pending:    make(map[uint64]*checkpointResult),
		maxPending: maxPending,
		eventsCh:   eventsCh,
		txCh:       txCh,
		lggr:       logger.Named(lggr, "CheckpointSequencer"),
	}
	s.cond = sync.NewCond(&s.mu)

	return s
}

// Submit hands a processed checkpoint to the sequencer. It blocks while the reorder buffer
// is full, and while draining blocks on full output channels (intentional backpressure on
// the workers). Results below the watermark (re-scans, duplicates of already-emitted
// checkpoints) bypass ordering and are emitted immediately — inserts are idempotent.
func (s *checkpointSequencer) Submit(ctx context.Context, r checkpointResult) error {
	// Wake any cond-waiters when the caller's context is canceled so they can observe it.
	stop := context.AfterFunc(ctx, func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.cond.Broadcast()
	})
	defer stop()

	s.mu.Lock()
	for {
		if s.closed {
			s.mu.Unlock()
			return errSequencerClosed
		}
		if err := ctx.Err(); err != nil {
			s.mu.Unlock()
			return err
		}
		if r.seq < s.nextSeq {
			// Already emitted once in order; this is a rescan or retry duplicate. Emit
			// directly without consuming pending capacity (cannot deadlock the buffer).
			s.mu.Unlock()
			return s.emit(ctx, &r)
		}
		if _, dup := s.pending[r.seq]; dup {
			// A retry chunk overlapped an in-flight result; drop the duplicate.
			s.mu.Unlock()
			return nil
		}
		if len(s.pending) < s.maxPending || r.seq == s.nextSeq {
			break
		}
		s.cond.Wait()
	}

	s.pending[r.seq] = &r
	err := s.drainLocked(ctx)
	s.mu.Unlock()

	return err
}

// NextExpected returns the current watermark: the lowest sequence not yet emitted.
func (s *checkpointSequencer) NextExpected() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.nextSeq
}

// setStart positions the watermark. Must only be called before any Submit (the poller sets
// it once at startup after computing the start sequence).
func (s *checkpointSequencer) setStart(seq uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextSeq = seq
}

// isPending reports whether seq has already been fetched and is buffered awaiting emission.
func (s *checkpointSequencer) isPending(seq uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.pending[seq]
	return ok
}

// SkipTo advances the watermark to seq, emitting any buffered results below it in order and
// logging sequences that are skipped without ever having been submitted. Used when a gap can
// never be filled (pruned provider history, or a checkpoint that permanently fails).
func (s *checkpointSequencer) SkipTo(ctx context.Context, seq uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.nextSeq < seq {
		if r, ok := s.pending[s.nextSeq]; ok {
			if err := s.emit(ctx, r); err != nil {
				return
			}
			delete(s.pending, s.nextSeq)
		} else {
			s.lggr.Errorw("Skipping checkpoint without emitting it", "sequence", s.nextSeq, "skipTo", seq)
		}
		s.nextSeq++
	}

	// The skip may have made buffered successors contiguous.
	if err := s.drainLocked(ctx); err != nil {
		return
	}
	s.cond.Broadcast()
}

// close marks the sequencer closed and wakes all blocked submitters. Buffered results are
// discarded; the poller re-fetches them on restart from the persisted cursor.
func (s *checkpointSequencer) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
	s.cond.Broadcast()
}

// drainLocked emits contiguous results starting at the watermark. Caller must hold s.mu.
// Emitting under the lock is intentional: a full output channel stalls all submitters,
// propagating backpressure to the catch-up workers.
func (s *checkpointSequencer) drainLocked(ctx context.Context) error {
	drained := false
	for {
		r, ok := s.pending[s.nextSeq]
		if !ok {
			break
		}
		if err := s.emit(ctx, r); err != nil {
			// Shutting down; leave the result pending — it is re-fetched after restart.
			return err
		}
		delete(s.pending, s.nextSeq)
		s.nextSeq++
		drained = true
	}
	if drained {
		s.cond.Broadcast()
	}

	return nil
}

// emit sends a result to the output channels: events first (only when non-empty, preserving
// the channel contract), then the transactions batch.
func (s *checkpointSequencer) emit(ctx context.Context, r *checkpointResult) error {
	if len(r.events.Events) > 0 {
		select {
		case s.eventsCh <- r.events:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	select {
	case s.txCh <- r.txs:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
