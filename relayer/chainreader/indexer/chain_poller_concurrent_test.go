package indexer

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// concurrentCheckpointClient is a thread-safe fake used by tests that run the full poller
// loop (concurrent workers). Failures are injected per sequence and decremented per call.
type concurrentCheckpointClient struct {
	testutils.FakeSuiPTBClient

	mu             sync.Mutex
	latest         uint64
	lowest         uint64
	failures       map[uint64]int // remaining generic RPC errors per sequence
	notFound       map[uint64]int // remaining NotFound responses per sequence
	processedCount map[uint64]int
}

func newConcurrentCheckpointClient(latest, lowest uint64) *concurrentCheckpointClient {
	return &concurrentCheckpointClient{
		latest:         latest,
		lowest:         lowest,
		failures:       make(map[uint64]int),
		notFound:       make(map[uint64]int),
		processedCount: make(map[uint64]int),
	}
}

func (c *concurrentCheckpointClient) GetLatestCheckpoint(ctx context.Context) (*suirpcv2.Checkpoint, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	seq := c.latest
	return &suirpcv2.Checkpoint{SequenceNumber: &seq}, nil
}

func (c *concurrentCheckpointClient) GetCheckpointAvailability(ctx context.Context) (*suirpcv2.GetServiceInfoResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	lowest := c.lowest
	return &suirpcv2.GetServiceInfoResponse{LowestAvailableCheckpoint: &lowest}, nil
}

func (c *concurrentCheckpointClient) GetTransaction(ctx context.Context, digest string) (client.TransactionDetails, error) {
	return client.TransactionDetails{}, nil
}

func (c *concurrentCheckpointClient) GetCheckpointData(ctx context.Context, seq uint64) (*client.CheckpointData, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.notFound[seq] > 0 {
		c.notFound[seq]--
		return nil, status.Error(codes.NotFound, "checkpoint not found")
	}
	if c.failures[seq] > 0 {
		c.failures[seq]--
		return nil, errors.New("simulated rpc failure")
	}

	c.processedCount[seq]++
	s := seq
	return &client.CheckpointData{
		Checkpoint: &suirpcv2.Checkpoint{SequenceNumber: &s},
	}, nil
}

func (c *concurrentCheckpointClient) processed(seq uint64) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.processedCount[seq]
}

func (c *concurrentCheckpointClient) totalProcessed() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := 0
	for _, count := range c.processedCount {
		n += count
	}
	return n
}

// fakeCursorStore is an in-memory CheckpointCursorStore recording all upserts.
type fakeCursorStore struct {
	mu      sync.Mutex
	seq     uint64
	found   bool
	getErr  error
	upserts []uint64
}

func (f *fakeCursorStore) GetCheckpointCursor(ctx context.Context, id string) (uint64, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.getErr != nil {
		return 0, false, f.getErr
	}
	return f.seq, f.found, nil
}

func (f *fakeCursorStore) UpsertCheckpointCursor(ctx context.Context, id string, seq uint64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.upserts = append(f.upserts, seq)
	if !f.found || seq > f.seq {
		f.seq = seq
		f.found = true
	}
	return nil
}

func (f *fakeCursorStore) latest() (uint64, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.seq, f.found
}

func noSelectors() []*sui.EventFilterByMoveEventModule { return nil }

// startPollerCollectingTxSeqs starts the poller and returns a snapshot function for the
// checkpoint sequences observed on the transactions channel (events channel is drained).
func startPollerCollectingTxSeqs(t *testing.T, cp *ChainPoller) func() []uint64 {
	t.Helper()

	var mu sync.Mutex
	var seqs []uint64

	go func() {
		for b := range cp.EventsChannel() {
			_ = b
		}
	}()
	go func() {
		for b := range cp.TransactionsChannel() {
			mu.Lock()
			seqs = append(seqs, b.Checkpoint.SequenceNumber)
			mu.Unlock()
		}
	}()

	require.NoError(t, cp.Start(context.Background()))
	t.Cleanup(func() { require.NoError(t, cp.Close()) })

	return func() []uint64 {
		mu.Lock()
		defer mu.Unlock()
		out := make([]uint64, len(seqs))
		copy(out, seqs)
		return out
	}
}

func seqRange(from, to uint64) []uint64 {
	out := make([]uint64, 0, to-from+1)
	for seq := from; seq <= to; seq++ {
		out = append(out, seq)
	}
	return out
}

// firstOccurrences returns seqs with duplicates removed, keeping first-occurrence order. The
// channel contract is "first delivery of each checkpoint is in ascending order, with possible
// idempotent re-deliveries" (gap retries and rescans re-emit already-delivered checkpoints), so
// ordering assertions are made on first occurrences.
func firstOccurrences(seqs []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(seqs))
	out := make([]uint64, 0, len(seqs))
	for _, s := range seqs {
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func uint64Ptr(v uint64) *uint64 { return &v }

func testPollerConfig(start uint64) sui.ChainPollerConfig {
	return sui.ChainPollerConfig{
		PollingInterval:         10 * time.Millisecond,
		SyncTimeout:             time.Minute,
		ChannelBufferSize:       16,
		StartCheckpointSequence: uint64Ptr(start),
	}
}

func TestPollerChunkedDispatchEmitsInOrder(t *testing.T) {
	t.Parallel()

	mockClient := newConcurrentCheckpointClient(160, 1)
	cursors := &fakeCursorStore{}
	cp := NewChainPoller(
		mockClient,
		logger.Test(t),
		testPollerConfig(100),
		noSelectors,
		WithWorkerPool(4, 5),
		WithCursorStore(cursors, "test"),
	)

	collected := startPollerCollectingTxSeqs(t, cp)

	require.Eventually(t, func() bool {
		return len(firstOccurrences(collected())) == 160-100+1
	}, 10*time.Second, 10*time.Millisecond, "expected all checkpoints emitted")

	// Despite 4 concurrent workers, first delivery is strictly ordered by checkpoint sequence.
	require.Equal(t, seqRange(100, 160), firstOccurrences(collected()))

	// The cursor converges on the last emitted checkpoint.
	require.Eventually(t, func() bool {
		seq, found := cursors.latest()
		return found && seq == 160
	}, 10*time.Second, 10*time.Millisecond)

	// With all chunks completed, the in-flight catch-up gauge returns to zero.
	require.Eventually(t, func() bool {
		return cp.goroutineCount.Load() == 0
	}, 10*time.Second, 10*time.Millisecond)
}

func TestPollerDispatchBoundedByReorderWindow(t *testing.T) {
	t.Parallel()

	// Huge backlog, but NO consumers on the output channels: the sequencer and channel
	// buffers fill, the watermark freezes, and dispatch must stop within the reorder-buffer
	// window instead of racing ahead fetching checkpoints nobody can emit yet.
	mockClient := newConcurrentCheckpointClient(10_000, 1)
	cfg := testPollerConfig(0)
	cfg.ChannelBufferSize = 4
	cp := NewChainPoller(
		mockClient,
		logger.Test(t),
		cfg,
		noSelectors,
		WithWorkerPool(2, 4), // reorder window = 8
	)

	require.NoError(t, cp.Start(context.Background()))
	t.Cleanup(func() { require.NoError(t, cp.Close()) })

	time.Sleep(500 * time.Millisecond) // ~50 poll ticks

	// Bound ≈ emitted-into-channels + reorder buffer + in-flight fetches; far below the
	// queue-capacity lookahead (32 chunks = 128) an unbounded dispatcher would consume.
	require.Less(t, mockClient.totalProcessed(), 40)
}

func TestPollerRetriesTransientFailureAndKeepsOrder(t *testing.T) {
	t.Parallel()

	mockClient := newConcurrentCheckpointClient(130, 1)
	mockClient.failures[115] = 3 // fails 3 times, then succeeds

	cursors := &fakeCursorStore{}
	cp := NewChainPoller(
		mockClient,
		logger.Test(t),
		testPollerConfig(100),
		noSelectors,
		WithWorkerPool(4, 5),
		WithCursorStore(cursors, "test"),
	)

	collected := startPollerCollectingTxSeqs(t, cp)

	require.Eventually(t, func() bool {
		return len(firstOccurrences(collected())) == 130-100+1
	}, 10*time.Second, 10*time.Millisecond)

	require.Equal(t, seqRange(100, 130), firstOccurrences(collected()))
}

func TestPollerNotFoundAtTipRetriesNextTick(t *testing.T) {
	t.Parallel()

	mockClient := newConcurrentCheckpointClient(110, 1)
	mockClient.notFound[110] = 3 // tip not yet available for the first few fetches

	cp := NewChainPoller(
		mockClient,
		logger.Test(t),
		testPollerConfig(100),
		noSelectors,
		WithWorkerPool(2, 4),
	)

	collected := startPollerCollectingTxSeqs(t, cp)

	require.Eventually(t, func() bool {
		return len(firstOccurrences(collected())) == 110-100+1
	}, 10*time.Second, 10*time.Millisecond)

	require.Equal(t, seqRange(100, 110), firstOccurrences(collected()))
}

func TestPollerSkipsPermanentlyFailingCheckpoint(t *testing.T) {
	t.Parallel()

	mockClient := newConcurrentCheckpointClient(110, 1)
	mockClient.failures[105] = 1 << 30 // effectively permanent

	cursors := &fakeCursorStore{}
	cp := NewChainPoller(
		mockClient,
		logger.Test(t),
		testPollerConfig(100),
		noSelectors,
		WithWorkerPool(2, 4),
		WithCursorStore(cursors, "test"),
	)
	cp.maxStalledTicks = 3 // keep the test fast

	collected := startPollerCollectingTxSeqs(t, cp)

	expected := append(seqRange(100, 104), seqRange(106, 110)...)
	require.Eventually(t, func() bool {
		return len(firstOccurrences(collected())) == len(expected)
	}, 10*time.Second, 10*time.Millisecond)

	// 105 was skipped after repeated stalled ticks; everything else is emitted in order.
	require.Equal(t, expected, firstOccurrences(collected()))

	require.Eventually(t, func() bool {
		seq, found := cursors.latest()
		return found && seq == 110
	}, 10*time.Second, 10*time.Millisecond)
}

func TestPollerRescanRecentReEmitsBelowWatermark(t *testing.T) {
	t.Parallel()

	mockClient := newConcurrentCheckpointClient(160, 1)
	// Chunk queue capacity (workers × 4 = 32) comfortably fits the rescan range
	// (100 checkpoints / 15 per chunk = 7 chunks) so nothing is dropped.
	cp := NewChainPoller(
		mockClient,
		logger.Test(t),
		testPollerConfig(150),
		noSelectors,
		WithWorkerPool(8, 15),
	)

	collected := startPollerCollectingTxSeqs(t, cp)

	require.Eventually(t, func() bool {
		return len(firstOccurrences(collected())) == 160-150+1
	}, 10*time.Second, 10*time.Millisecond)

	// Rescan re-enqueues [60, 159] (default window of 100 checkpoints from the requested start);
	// results below the watermark bypass ordering and are re-emitted directly.
	require.NoError(t, cp.RescanFrom(context.Background(), 60))

	require.Eventually(t, func() bool {
		return mockClient.processed(60) >= 1 && mockClient.processed(155) >= 2
	}, 10*time.Second, 10*time.Millisecond, "rescan should re-fetch old checkpoints")

	require.Eventually(t, func() bool {
		seen := make(map[uint64]int)
		for _, seq := range collected() {
			seen[seq]++
		}
		return seen[60] >= 1 && seen[155] >= 2
	}, 10*time.Second, 10*time.Millisecond, "rescan results should be re-emitted")
}

func TestComputeStartSequenceWithCursor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	// overlap = workers*chunkSize + bufferSize = 2*3 + 4 = 10
	newPoller := func(cfg sui.ChainPollerConfig, store CheckpointCursorStore) *ChainPoller {
		cfg.ChannelBufferSize = 4
		mockClient := newConcurrentCheckpointClient(5000, 1)
		return NewChainPoller(
			mockClient,
			logger.Test(t),
			cfg,
			noSelectors,
			WithWorkerPool(2, 3),
			WithCursorStore(store, "test"),
		)
	}

	t.Run("resumes from cursor minus overlap", func(t *testing.T) {
		t.Parallel()
		cp := newPoller(sui.ChainPollerConfig{SyncTimeout: time.Minute}, &fakeCursorStore{seq: 1000, found: true})
		start, err := cp.computeStartSequence(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(991), start) // 1000 + 1 - 10
	})

	t.Run("explicit start skips ahead of cursor", func(t *testing.T) {
		t.Parallel()
		cp := newPoller(sui.ChainPollerConfig{
			SyncTimeout:             time.Minute,
			StartCheckpointSequence: uint64Ptr(1500),
		}, &fakeCursorStore{seq: 1000, found: true})
		start, err := cp.computeStartSequence(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1500), start)
	})

	t.Run("explicit start does not rewind below cursor", func(t *testing.T) {
		t.Parallel()
		cp := newPoller(sui.ChainPollerConfig{
			SyncTimeout:             time.Minute,
			StartCheckpointSequence: uint64Ptr(500),
		}, &fakeCursorStore{seq: 1000, found: true})
		start, err := cp.computeStartSequence(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(991), start)
	})

	t.Run("small cursor floors resume at zero then clamps to provider floor", func(t *testing.T) {
		t.Parallel()
		cp := newPoller(sui.ChainPollerConfig{SyncTimeout: time.Minute}, &fakeCursorStore{seq: 5, found: true})
		start, err := cp.computeStartSequence(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), start) // resume 0, clamped to provider floor 1
	})

	t.Run("cursor load error falls back to configured start", func(t *testing.T) {
		t.Parallel()
		cp := newPoller(sui.ChainPollerConfig{
			SyncTimeout:             time.Minute,
			StartCheckpointSequence: uint64Ptr(700),
		}, &fakeCursorStore{getErr: errors.New("db down")})
		start, err := cp.computeStartSequence(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(700), start)
	})

	t.Run("no cursor defaults to latest", func(t *testing.T) {
		t.Parallel()
		cp := newPoller(sui.ChainPollerConfig{SyncTimeout: time.Minute}, &fakeCursorStore{})
		start, err := cp.computeStartSequence(ctx)
		require.NoError(t, err)
		require.Equal(t, uint64(5000), start)
	})
}
