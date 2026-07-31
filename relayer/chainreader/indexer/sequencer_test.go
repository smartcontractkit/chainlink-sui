package indexer

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func resultWithEvents(seq uint64) checkpointResult {
	meta := CheckpointMeta{SequenceNumber: seq}
	return checkpointResult{
		seq: seq,
		events: CheckpointEventsBatch{
			Checkpoint: meta,
			Events:     []CheckpointEventItem{{Event: &suirpcv2.Event{}, TxDigest: "digest", EventIndex: 0}},
		},
		txs: CheckpointTransactionsBatch{Checkpoint: meta},
	}
}

func resultWithoutEvents(seq uint64) checkpointResult {
	meta := CheckpointMeta{SequenceNumber: seq}
	return checkpointResult{
		seq:    seq,
		events: CheckpointEventsBatch{Checkpoint: meta},
		txs:    CheckpointTransactionsBatch{Checkpoint: meta},
	}
}

// drainChannels consumes both channels until they are closed or ctx expires, recording the
// order of checkpoint sequences observed on each.
func drainChannels(
	eventsCh <-chan CheckpointEventsBatch,
	txCh <-chan CheckpointTransactionsBatch,
) (eventSeqs func() []uint64, txSeqs func() []uint64, wait func()) {
	var mu sync.Mutex
	var evs, txs []uint64
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for b := range eventsCh {
			mu.Lock()
			evs = append(evs, b.Checkpoint.SequenceNumber)
			mu.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		for b := range txCh {
			mu.Lock()
			txs = append(txs, b.Checkpoint.SequenceNumber)
			mu.Unlock()
		}
	}()

	snapshot := func(src *[]uint64) func() []uint64 {
		return func() []uint64 {
			mu.Lock()
			defer mu.Unlock()
			out := make([]uint64, len(*src))
			copy(out, *src)
			return out
		}
	}

	return snapshot(&evs), snapshot(&txs), wg.Wait
}

func TestSequencerOrdersShuffledConcurrentSubmits(t *testing.T) {
	ctx := context.Background()
	eventsCh := make(chan CheckpointEventsBatch, 256)
	txCh := make(chan CheckpointTransactionsBatch, 256)
	// Buffer sized so no Submit blocks: chunks are picked up in shuffled (non-FIFO) order here,
	// unlike the poller's ascending dispatch, so the watermark chunk may be processed last.
	s := newCheckpointSequencer(100, 100, eventsCh, txCh, logger.Test(t))

	// Workers pull contiguous chunks and submit each chunk's sequences in ascending order —
	// the sequencer's per-worker ordering invariant.
	const chunkSize = 5
	chunks := make([]chunkRange, 0, 20)
	for from := uint64(100); from < 200; from += chunkSize {
		chunks = append(chunks, chunkRange{start: from, end: from + chunkSize - 1})
	}
	rand.Shuffle(len(chunks), func(i, j int) { chunks[i], chunks[j] = chunks[j], chunks[i] })

	chunkCh := make(chan chunkRange, len(chunks))
	for _, c := range chunks {
		chunkCh <- c
	}
	close(chunkCh)

	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for c := range chunkCh {
				for seq := c.start; seq <= c.end; seq++ {
					assert.NoError(t, s.Submit(ctx, resultWithEvents(seq)))
				}
			}
		})
	}
	wg.Wait()

	close(eventsCh)
	close(txCh)
	eventSeqs, txSeqs, wait := drainChannels(eventsCh, txCh)
	wait()

	expected := make([]uint64, 0, 100)
	for seq := uint64(100); seq < 200; seq++ {
		expected = append(expected, seq)
	}
	assert.Equal(t, expected, eventSeqs())
	assert.Equal(t, expected, txSeqs())
	assert.Equal(t, uint64(200), s.NextExpected())
}

func TestSequencerEventlessCheckpointAdvancesWatermark(t *testing.T) {
	ctx := context.Background()
	eventsCh := make(chan CheckpointEventsBatch, 16)
	txCh := make(chan CheckpointTransactionsBatch, 16)
	s := newCheckpointSequencer(10, 8, eventsCh, txCh, logger.Test(t))

	require.NoError(t, s.Submit(ctx, resultWithEvents(10)))
	require.NoError(t, s.Submit(ctx, resultWithoutEvents(11)))
	require.NoError(t, s.Submit(ctx, resultWithEvents(12)))

	close(eventsCh)
	close(txCh)
	eventSeqs, txSeqs, wait := drainChannels(eventsCh, txCh)
	wait()

	// Eventless checkpoint 11 emits no events batch but still advances ordering.
	assert.Equal(t, []uint64{10, 12}, eventSeqs())
	assert.Equal(t, []uint64{10, 11, 12}, txSeqs())
	assert.Equal(t, uint64(13), s.NextExpected())
}

func TestSequencerBelowWatermarkBypasses(t *testing.T) {
	ctx := context.Background()
	eventsCh := make(chan CheckpointEventsBatch, 16)
	txCh := make(chan CheckpointTransactionsBatch, 16)
	s := newCheckpointSequencer(50, 4, eventsCh, txCh, logger.Test(t))

	// Fill the pending buffer with out-of-order results so capacity is exhausted.
	for seq := uint64(52); seq < 56; seq++ {
		require.NoError(t, s.Submit(ctx, resultWithEvents(seq)))
	}
	require.Equal(t, uint64(50), s.NextExpected())

	// A rescan of an already-emitted checkpoint must not block even with a full buffer.
	done := make(chan struct{})
	go func() {
		defer close(done)
		assert.NoError(t, s.Submit(ctx, resultWithEvents(42)))
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("below-watermark Submit blocked")
	}

	// Bypass emitted directly; watermark unchanged.
	assert.Equal(t, uint64(50), s.NextExpected())
	assert.Equal(t, uint64(42), (<-eventsCh).Checkpoint.SequenceNumber)
	assert.Equal(t, uint64(42), (<-txCh).Checkpoint.SequenceNumber)
}

func TestSequencerDropsDuplicatePending(t *testing.T) {
	ctx := context.Background()
	eventsCh := make(chan CheckpointEventsBatch, 16)
	txCh := make(chan CheckpointTransactionsBatch, 16)
	s := newCheckpointSequencer(10, 8, eventsCh, txCh, logger.Test(t))

	require.NoError(t, s.Submit(ctx, resultWithEvents(12)))
	require.NoError(t, s.Submit(ctx, resultWithEvents(12))) // duplicate: dropped
	require.NoError(t, s.Submit(ctx, resultWithEvents(10)))
	require.NoError(t, s.Submit(ctx, resultWithEvents(11)))

	close(eventsCh)
	close(txCh)
	eventSeqs, txSeqs, wait := drainChannels(eventsCh, txCh)
	wait()

	assert.Equal(t, []uint64{10, 11, 12}, eventSeqs())
	assert.Equal(t, []uint64{10, 11, 12}, txSeqs())
}

func TestSequencerBlocksWhenFullAndUnblocksOnDrain(t *testing.T) {
	ctx := context.Background()
	eventsCh := make(chan CheckpointEventsBatch, 16)
	txCh := make(chan CheckpointTransactionsBatch, 16)
	s := newCheckpointSequencer(10, 2, eventsCh, txCh, logger.Test(t))

	// Two out-of-order results fill the buffer (watermark stays at 10).
	require.NoError(t, s.Submit(ctx, resultWithEvents(12)))
	require.NoError(t, s.Submit(ctx, resultWithEvents(13)))

	blocked := make(chan struct{})
	go func() {
		defer close(blocked)
		assert.NoError(t, s.Submit(ctx, resultWithEvents(14))) // full and not next-expected: blocks
	}()

	select {
	case <-blocked:
		t.Fatal("Submit should have blocked while buffer is full")
	case <-time.After(100 * time.Millisecond):
	}

	// Submitting the next-expected sequences drains 10..13 and frees capacity for 14.
	require.NoError(t, s.Submit(ctx, resultWithEvents(10)))
	require.NoError(t, s.Submit(ctx, resultWithEvents(11)))

	select {
	case <-blocked:
	case <-time.After(5 * time.Second):
		t.Fatal("Submit did not unblock after drain")
	}

	close(eventsCh)
	close(txCh)
	eventSeqs, _, wait := drainChannels(eventsCh, txCh)
	wait()
	assert.Equal(t, []uint64{10, 11, 12, 13, 14}, eventSeqs())
}

func TestSequencerSkipToFlushesAndAdvances(t *testing.T) {
	ctx := context.Background()
	eventsCh := make(chan CheckpointEventsBatch, 16)
	txCh := make(chan CheckpointTransactionsBatch, 16)
	s := newCheckpointSequencer(10, 8, eventsCh, txCh, logger.Test(t))

	// 10 is missing (permanently failed); 11 and 13 are buffered.
	require.NoError(t, s.Submit(ctx, resultWithEvents(11)))
	require.NoError(t, s.Submit(ctx, resultWithEvents(13)))
	require.Equal(t, uint64(10), s.NextExpected())

	s.SkipTo(ctx, 12)

	// Buffered 11 was flushed in order, missing 10 skipped, watermark waits on 12.
	assert.Equal(t, uint64(12), s.NextExpected())

	require.NoError(t, s.Submit(ctx, resultWithEvents(12)))
	assert.Equal(t, uint64(14), s.NextExpected())

	close(eventsCh)
	close(txCh)
	eventSeqs, txSeqs, wait := drainChannels(eventsCh, txCh)
	wait()
	assert.Equal(t, []uint64{11, 12, 13}, eventSeqs())
	assert.Equal(t, []uint64{11, 12, 13}, txSeqs())
}

func TestSequencerCtxCancelUnblocksSubmit(t *testing.T) {
	eventsCh := make(chan CheckpointEventsBatch, 16)
	txCh := make(chan CheckpointTransactionsBatch, 16)
	s := newCheckpointSequencer(10, 1, eventsCh, txCh, logger.Test(t))

	require.NoError(t, s.Submit(context.Background(), resultWithEvents(12)))

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Submit(ctx, resultWithEvents(13)) // buffer full: blocks
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(5 * time.Second):
		t.Fatal("Submit did not unblock on ctx cancel")
	}
}

func TestSequencerCloseReleasesWaiters(t *testing.T) {
	ctx := context.Background()
	eventsCh := make(chan CheckpointEventsBatch, 16)
	txCh := make(chan CheckpointTransactionsBatch, 16)
	s := newCheckpointSequencer(10, 1, eventsCh, txCh, logger.Test(t))

	require.NoError(t, s.Submit(ctx, resultWithEvents(12)))

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Submit(ctx, resultWithEvents(13)) // buffer full: blocks
	}()

	time.Sleep(50 * time.Millisecond)
	s.close()

	select {
	case err := <-errCh:
		require.ErrorIs(t, err, errSequencerClosed)
	case <-time.After(5 * time.Second):
		t.Fatal("Submit did not unblock on close")
	}

	// Closed sequencer rejects further submissions.
	require.ErrorIs(t, s.Submit(ctx, resultWithEvents(10)), errSequencerClosed)
}
