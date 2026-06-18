package loop_test

// This load test reproduces, in isolation, the read-path load scenario that makes the Sui CCIP E2E lanes
// fail: a large burst of concurrent read RPCs against a node whose per-connection concurrency is limited.
//
// It stands up an in-process mock Sui LedgerService whose GetObject has a configurable artificial latency
// and whose gRPC server advertises a configurable MaxConcurrentStreams. It then drives the real
// client.PTBClient (the same code path used in production) with a configurable burst of concurrent
// ReadObjectId calls, for a configurable gRPC connection-pool size, and reports latency/throughput.
//
// The goal is to validate the hypothesis that a SINGLE shared gRPC connection becomes a head-of-line
// bottleneck under bursty reads, and that spreading calls across a pool of connections relieves it.
//
// Everything is configurable via loadScenario so the assumptions can be probed independently:
//   - PoolSize:       number of gRPC connections (round-robin)
//   - StreamLimit:    server-advertised MaxConcurrentStreams per connection
//   - CallLatency:    artificial per-RPC server latency
//   - Concurrency:    number of reads issued simultaneously (the burst size)

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// mockLedgerServer is a minimal LedgerService whose GetObject sleeps for a fixed latency and tracks the
// maximum number of concurrently in-flight GetObject calls it observed (i.e. the realized concurrency at
// the "node", which is bounded by poolSize * StreamLimit).
type mockLedgerServer struct {
	suirpcv2.UnimplementedLedgerServiceServer

	latency     time.Duration
	inFlight    int64
	maxInFlight int64
	totalCalls  int64
}

func (m *mockLedgerServer) GetObject(ctx context.Context, _ *suirpcv2.GetObjectRequest) (*suirpcv2.GetObjectResponse, error) {
	atomic.AddInt64(&m.totalCalls, 1)
	cur := atomic.AddInt64(&m.inFlight, 1)
	defer atomic.AddInt64(&m.inFlight, -1)

	for {
		old := atomic.LoadInt64(&m.maxInFlight)
		if cur <= old || atomic.CompareAndSwapInt64(&m.maxInFlight, old, cur) {
			break
		}
	}

	select {
	case <-time.After(m.latency):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	return &suirpcv2.GetObjectResponse{Object: &suirpcv2.Object{}}, nil
}

type loadScenario struct {
	Name        string
	PoolSize    int
	StreamLimit uint32
	CallLatency time.Duration
	Concurrency int
}

func (s loadScenario) withDefaults() loadScenario {
	if s.PoolSize <= 0 {
		s.PoolSize = 1
	}
	if s.StreamLimit == 0 {
		s.StreamLimit = 4
	}
	if s.CallLatency == 0 {
		s.CallLatency = 150 * time.Millisecond
	}
	if s.Concurrency <= 0 {
		s.Concurrency = 64
	}
	return s
}

type loadResult struct {
	WallTime   time.Duration
	P50        time.Duration
	P99        time.Duration
	ServerPeak int64
	TotalCalls int64
	Errors     int
	PoolSize   int
}

// startMockLedger starts an in-process gRPC LedgerService with the given per-call latency and
// MaxConcurrentStreams, returning the listen address and a stop func.
func startMockLedger(t *testing.T, latency time.Duration, streamLimit uint32) (string, *mockLedgerServer, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv := grpc.NewServer(grpc.MaxConcurrentStreams(streamLimit))
	mock := &mockLedgerServer{latency: latency}
	suirpcv2.RegisterLedgerServiceServer(srv, mock)

	go func() { _ = srv.Serve(lis) }()

	return lis.Addr().String(), mock, func() { srv.Stop() }
}

// runLoadScenario drives a burst of concurrent ReadObjectId calls through a PTBClient configured with the
// scenario's pool size and returns aggregate timing.
func runLoadScenario(t *testing.T, s loadScenario) loadResult {
	t.Helper()
	s = s.withDefaults()

	addr, mock, stop := startMockLedger(t, s.CallLatency, s.StreamLimit)
	defer stop()

	c, err := client.NewPTBClient(logger.Test(t), client.PTBClientConfig{
		GrpcTarget:            addr,
		GrpcToken:             "load-test-token",
		TransactionTimeout:    2 * time.Minute,
		MaxConcurrentRequests: int64(s.Concurrency * 4),
		MaxGrpcConnections:    s.PoolSize,
	})
	require.NoError(t, err)
	defer func() { _ = c.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Warm up all pooled connections so connection establishment is not counted in the timed burst.
	for i := 0; i < s.PoolSize; i++ {
		_, _ = c.ReadObjectId(ctx, "0x1")
	}

	latencies := make([]time.Duration, s.Concurrency)
	var errCount int64

	var wg sync.WaitGroup
	wg.Add(s.Concurrency)
	start := time.Now()
	for i := 0; i < s.Concurrency; i++ {
		go func(idx int) {
			defer wg.Done()
			callStart := time.Now()
			if _, err := c.ReadObjectId(ctx, "0x1"); err != nil {
				atomic.AddInt64(&errCount, 1)
			}
			latencies[idx] = time.Since(callStart)
		}(i)
	}
	wg.Wait()
	wall := time.Since(start)

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	return loadResult{
		WallTime:   wall,
		P50:        latencies[len(latencies)/2],
		P99:        latencies[(len(latencies)*99)/100],
		ServerPeak: atomic.LoadInt64(&mock.maxInFlight),
		TotalCalls: atomic.LoadInt64(&mock.totalCalls),
		Errors:     int(errCount),
		PoolSize:   s.PoolSize,
	}
}

func logResult(t *testing.T, label string, r loadResult) {
	t.Helper()
	t.Logf("%-18s pool=%d wall=%-10s p50=%-10s p99=%-10s serverPeakConcurrency=%d calls=%d errors=%d",
		label, r.PoolSize, r.WallTime.Round(time.Millisecond), r.P50.Round(time.Millisecond),
		r.P99.Round(time.Millisecond), r.ServerPeak, r.TotalCalls, r.Errors)
}

// TestConnectionPoolRelievesStreamBottleneck is the conclusive comparison: under a burst of concurrent
// reads against a connection with a small MaxConcurrentStreams, a single connection serializes the work
// into waves while a pool of connections runs them in parallel. This mirrors the E2E read-path stall.
func TestConnectionPoolRelievesStreamBottleneck(t *testing.T) {
	t.Parallel()

	const (
		streamLimit = uint32(4)
		latency     = 150 * time.Millisecond
		concurrency = 64
		poolSize    = 8
	)

	single := runLoadScenario(t, loadScenario{
		Name:        "single-connection",
		PoolSize:    1,
		StreamLimit: streamLimit,
		CallLatency: latency,
		Concurrency: concurrency,
	})
	logResult(t, "single-connection", single)

	pooled := runLoadScenario(t, loadScenario{
		Name:        "pooled-connections",
		PoolSize:    poolSize,
		StreamLimit: streamLimit,
		CallLatency: latency,
		Concurrency: concurrency,
	})
	logResult(t, "pooled-connections", pooled)

	require.Zero(t, single.Errors, "single-connection scenario should not error")
	require.Zero(t, pooled.Errors, "pooled scenario should not error")

	// The single connection is capped at streamLimit concurrent calls at the node; the pool reaches
	// roughly poolSize*streamLimit. The realized server-side concurrency is the smoking gun.
	require.LessOrEqual(t, single.ServerPeak, int64(streamLimit)+1,
		"single connection should be limited to ~StreamLimit concurrent calls at the node")
	require.Greater(t, pooled.ServerPeak, single.ServerPeak,
		"pool should achieve higher node-side concurrency than a single connection")

	// And that translates into materially better wall-clock throughput for the burst.
	ratio := float64(single.WallTime) / float64(pooled.WallTime)
	t.Logf("speedup (single/pooled) = %.2fx", ratio)
	require.Greater(t, ratio, 1.8,
		"connection pool should significantly outperform a single connection under stream contention")
}

// TestConnectionPoolScaling sweeps pool sizes so the relationship between pool size and burst latency is
// visible. It is informational (no hard assertions beyond success) and configurable for ad-hoc probing.
func TestConnectionPoolScaling(t *testing.T) {
	t.Parallel()

	for _, poolSize := range []int{1, 2, 4, 8, 16} {
		r := runLoadScenario(t, loadScenario{
			Name:        fmt.Sprintf("pool-%d", poolSize),
			PoolSize:    poolSize,
			StreamLimit: 4,
			CallLatency: 150 * time.Millisecond,
			Concurrency: 64,
		})
		logResult(t, fmt.Sprintf("pool-%d", poolSize), r)
		require.Zero(t, r.Errors)
	}
}
