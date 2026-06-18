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

// mockRPCNode tracks in-flight RPCs and per-method call counts for a mock Sui read node.
type mockRPCNode struct {
	latency     time.Duration
	inFlight    int64
	maxInFlight int64
	totalCalls  int64

	getObjectCalls          int64
	simulateCalls           int64
	listOwnedObjectsCalls   int64
	getEpochCalls           int64
}

func (m *mockRPCNode) track(ctx context.Context) error {
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
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// mockLedgerServer is a minimal LedgerService whose GetObject sleeps for a fixed latency and tracks the
// maximum number of concurrently in-flight GetObject calls it observed (i.e. the realized concurrency at
// the "node", which is bounded by poolSize * StreamLimit).
type mockLedgerServer struct {
	suirpcv2.UnimplementedLedgerServiceServer
	node *mockRPCNode
}

func (m *mockLedgerServer) GetObject(ctx context.Context, _ *suirpcv2.GetObjectRequest) (*suirpcv2.GetObjectResponse, error) {
	atomic.AddInt64(&m.node.getObjectCalls, 1)
	if err := m.node.track(ctx); err != nil {
		return nil, err
	}
	return &suirpcv2.GetObjectResponse{Object: &suirpcv2.Object{}}, nil
}

// mockSuiReadNode implements the gRPC services exercised by CCIP config polling and background
// chain-reader/indexer traffic against a single read node.
type mockSuiReadNode struct {
	suirpcv2.UnimplementedLedgerServiceServer
	suirpcv2.UnimplementedStateServiceServer
	suirpcv2.UnimplementedTransactionExecutionServiceServer

	node *mockRPCNode
}

func (m *mockSuiReadNode) GetObject(ctx context.Context, _ *suirpcv2.GetObjectRequest) (*suirpcv2.GetObjectResponse, error) {
	atomic.AddInt64(&m.node.getObjectCalls, 1)
	if err := m.node.track(ctx); err != nil {
		return nil, err
	}
	return &suirpcv2.GetObjectResponse{Object: &suirpcv2.Object{}}, nil
}

func (m *mockSuiReadNode) GetEpoch(ctx context.Context, _ *suirpcv2.GetEpochRequest) (*suirpcv2.GetEpochResponse, error) {
	atomic.AddInt64(&m.node.getEpochCalls, 1)
	if err := m.node.track(ctx); err != nil {
		return nil, err
	}
	gasPrice := uint64(1000)
	return &suirpcv2.GetEpochResponse{
		Epoch: &suirpcv2.Epoch{ReferenceGasPrice: &gasPrice},
	}, nil
}

func (m *mockSuiReadNode) ListOwnedObjects(ctx context.Context, _ *suirpcv2.ListOwnedObjectsRequest) (*suirpcv2.ListOwnedObjectsResponse, error) {
	atomic.AddInt64(&m.node.listOwnedObjectsCalls, 1)
	if err := m.node.track(ctx); err != nil {
		return nil, err
	}
	return &suirpcv2.ListOwnedObjectsResponse{Objects: []*suirpcv2.Object{{}}}, nil
}

func (m *mockSuiReadNode) SimulateTransaction(ctx context.Context, _ *suirpcv2.SimulateTransactionRequest) (*suirpcv2.SimulateTransactionResponse, error) {
	atomic.AddInt64(&m.node.simulateCalls, 1)
	if err := m.node.track(ctx); err != nil {
		return nil, err
	}
	success := true
	return &suirpcv2.SimulateTransactionResponse{
		Transaction: &suirpcv2.ExecutedTransaction{
			Effects: &suirpcv2.TransactionEffects{
				Status: &suirpcv2.ExecutionStatus{Success: &success},
			},
		},
	}, nil
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

type ccipLoadScenario struct {
	loadScenario
	// BatchCount is the number of concurrent BatchGetLatestValues callers (CCIP oracle nodes).
	BatchCount int
	// ReadsPerBatch is the number of reads inside each batch (CCIP config poller issues 6).
	ReadsPerBatch int
	// IntraBatchConcurrency limits concurrent reads within one batch. Production runs all reads in a
	// batch concurrently (BatchGetLatestValues spawns a goroutine per read with no intra-batch cap), so
	// the default equals ReadsPerBatch; the connection pool — not intra-batch serialization — is what
	// keeps the node from saturating. A lower value can be set to probe a throttled variant.
	IntraBatchConcurrency int
	// ObjectRefsPerRead is the number of shared-object metadata GetObject RPCs per read (TransformTransactionArg).
	ObjectRefsPerRead int
	// BackgroundWorkers is the number of concurrent poller/indexer goroutines interleaved with batches.
	BackgroundWorkers int
	// BatchDeadline, when > 0, wraps each batch in its own context deadline — modeling the CCIP config
	// poller's bgRefreshTimeout. A batch that cannot finish its reads before the deadline fails, which is
	// exactly the production symptom (config refresh times out => no merkle root committed).
	BatchDeadline time.Duration
}

func (s ccipLoadScenario) withDefaults() ccipLoadScenario {
	s.loadScenario = s.loadScenario.withDefaults()
	if s.BatchCount <= 0 {
		s.BatchCount = 4
	}
	if s.ReadsPerBatch <= 0 {
		s.ReadsPerBatch = 6
	}
	if s.IntraBatchConcurrency <= 0 {
		s.IntraBatchConcurrency = s.ReadsPerBatch
	}
	if s.ObjectRefsPerRead <= 0 {
		s.ObjectRefsPerRead = 2
	}
	if s.BackgroundWorkers <= 0 {
		s.BackgroundWorkers = 8
	}
	return s
}

type ccipLoadResult struct {
	loadResult
	GetObjectCalls        int64
	SimulateCalls         int64
	ListOwnedObjectsCalls int64
	GetEpochCalls         int64
	// BatchErrors counts BatchGetLatestValues-shaped batches that failed (e.g. a read hit the
	// per-batch deadline). In production each such failure is a config refresh that times out and
	// therefore a merkle root that never gets committed.
	BatchErrors int
}

// startMockLedger starts an in-process gRPC LedgerService with the given per-call latency and
// MaxConcurrentStreams, returning the listen address and a stop func.
func startMockLedger(t *testing.T, latency time.Duration, streamLimit uint32) (string, *mockRPCNode, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	node := &mockRPCNode{latency: latency}
	mock := &mockLedgerServer{node: node}

	srv := grpc.NewServer(grpc.MaxConcurrentStreams(streamLimit))
	suirpcv2.RegisterLedgerServiceServer(srv, mock)

	go func() { _ = srv.Serve(lis) }()

	return lis.Addr().String(), node, func() { srv.Stop() }
}

// startMockSuiReadNode starts an in-process mock Sui read node exposing Ledger, State, and
// TransactionExecution services with configurable per-RPC latency and MaxConcurrentStreams.
func startMockSuiReadNode(t *testing.T, latency time.Duration, streamLimit uint32) (string, *mockRPCNode, func()) {
	t.Helper()

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	node := &mockRPCNode{latency: latency}
	mock := &mockSuiReadNode{node: node}

	srv := grpc.NewServer(grpc.MaxConcurrentStreams(streamLimit))
	suirpcv2.RegisterLedgerServiceServer(srv, mock)
	suirpcv2.RegisterStateServiceServer(srv, mock)
	suirpcv2.RegisterTransactionExecutionServiceServer(srv, mock)

	go func() { _ = srv.Serve(lis) }()

	return lis.Addr().String(), node, func() { srv.Stop() }
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

// simulateConfigRead mirrors one GetLatestValue inside BatchGetLatestValues: object-metadata
// GetObject calls for shared config objects followed by a SimulateTransaction read.
func simulateConfigRead(ctx context.Context, c *client.PTBClient, objectIDs []string, simBCS []byte) error {
	for _, objectID := range objectIDs {
		if _, err := c.ReadObjectId(ctx, objectID); err != nil {
			return err
		}
	}
	_, err := c.SimulatePTB(ctx, simBCS)
	return err
}

// runBatchGetLatestValuesWorkload runs one BatchGetLatestValues-shaped burst. Reads within the batch
// are limited by intraBatchConcurrency (6 in E2E before serialization, 1 in production today).
func runBatchGetLatestValuesWorkload(
	ctx context.Context,
	c *client.PTBClient,
	readsPerBatch int,
	intraBatchConcurrency int,
	objectRefsPerRead int,
	simBCS []byte,
) error {
	objectIDs := make([]string, objectRefsPerRead)
	for i := range objectIDs {
		objectIDs[i] = fmt.Sprintf("0xconfig_object_%d", i)
	}

	readSem := make(chan struct{}, intraBatchConcurrency)
	var wg sync.WaitGroup
	var errCount int64

	for i := 0; i < readsPerBatch; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			readSem <- struct{}{}
			defer func() { <-readSem }()

			if err := simulateConfigRead(ctx, c, objectIDs, simBCS); err != nil {
				atomic.AddInt64(&errCount, 1)
			}
		}()
	}
	wg.Wait()

	if errCount > 0 {
		return fmt.Errorf("batch workload had %d read errors", errCount)
	}
	return nil
}

// runBackgroundReadWorkload interleaves poller/indexer RPCs that share the read-node connection.
func runBackgroundReadWorkload(ctx context.Context, c *client.PTBClient, workerID int) error {
	switch workerID % 3 {
	case 0:
		_, err := c.ReadFilterOwnedObjectIds(ctx, "0xowner", "0x2::coin::Coin<0x2::sui::SUI>", nil)
		return err
	case 1:
		_, err := c.GetReferenceGasPrice(ctx)
		return err
	default:
		_, err := c.ReadObjectId(ctx, "0xpoller_checkpoint")
		return err
	}
}

// runCCIPMixedLoadScenario drives concurrent BatchGetLatestValues batches interleaved with
// background poller/indexer reads through a PTBClient against a single read-node connection.
func runCCIPMixedLoadScenario(t *testing.T, s ccipLoadScenario) ccipLoadResult {
	t.Helper()
	s = s.withDefaults()

	addr, mock, stop := startMockSuiReadNode(t, s.CallLatency, s.StreamLimit)
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

	// Warm up all pooled connections across the services used by the mixed workload.
	for i := 0; i < s.PoolSize; i++ {
		_, _ = c.ReadObjectId(ctx, "0x1")
		_, _ = c.GetReferenceGasPrice(ctx)
		_, _ = c.ReadFilterOwnedObjectIds(ctx, "0xowner", "0x2::coin::Coin<0x2::sui::SUI>", nil)
		_, _ = c.SimulatePTB(ctx, []byte{0x01})
	}

	simBCS := []byte{0x01}
	totalWorkers := s.BatchCount + s.BackgroundWorkers
	latencies := make([]time.Duration, totalWorkers)
	var errCount int64
	var batchErrCount int64

	var wg sync.WaitGroup
	wg.Add(totalWorkers)
	start := time.Now()

	for batch := 0; batch < s.BatchCount; batch++ {
		go func(idx int) {
			defer wg.Done()

			// Each batch optionally runs under its own deadline, modeling the CCIP config poller's
			// bgRefreshTimeout that wraps every BatchGetLatestValues call.
			batchCtx := ctx
			if s.BatchDeadline > 0 {
				var cancel context.CancelFunc
				batchCtx, cancel = context.WithTimeout(ctx, s.BatchDeadline)
				defer cancel()
			}

			callStart := time.Now()
			if err := runBatchGetLatestValuesWorkload(batchCtx, c, s.ReadsPerBatch, s.IntraBatchConcurrency, s.ObjectRefsPerRead, simBCS); err != nil {
				atomic.AddInt64(&errCount, 1)
				atomic.AddInt64(&batchErrCount, 1)
			}
			latencies[idx] = time.Since(callStart)
		}(batch)
	}

	for worker := 0; worker < s.BackgroundWorkers; worker++ {
		go func(idx, workerID int) {
			defer wg.Done()
			callStart := time.Now()
			if err := runBackgroundReadWorkload(ctx, c, workerID); err != nil {
				atomic.AddInt64(&errCount, 1)
			}
			latencies[idx] = time.Since(callStart)
		}(s.BatchCount+worker, worker)
	}

	wg.Wait()
	wall := time.Since(start)

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	return ccipLoadResult{
		loadResult: loadResult{
			WallTime:   wall,
			P50:        latencies[len(latencies)/2],
			P99:        latencies[(len(latencies)*99)/100],
			ServerPeak: atomic.LoadInt64(&mock.maxInFlight),
			TotalCalls: atomic.LoadInt64(&mock.totalCalls),
			Errors:     int(errCount),
			PoolSize:   s.PoolSize,
		},
		GetObjectCalls:        atomic.LoadInt64(&mock.getObjectCalls),
		SimulateCalls:         atomic.LoadInt64(&mock.simulateCalls),
		ListOwnedObjectsCalls: atomic.LoadInt64(&mock.listOwnedObjectsCalls),
		GetEpochCalls:         atomic.LoadInt64(&mock.getEpochCalls),
		BatchErrors:           int(atomic.LoadInt64(&batchErrCount)),
	}
}

func logCCIPResult(t *testing.T, label string, r ccipLoadResult) {
	t.Helper()
	logResult(t, label, r.loadResult)
	t.Logf("%-18s getObject=%d simulate=%d listOwned=%d getEpoch=%d batchErrors=%d",
		label, r.GetObjectCalls, r.SimulateCalls, r.ListOwnedObjectsCalls, r.GetEpochCalls, r.BatchErrors)
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

// TestMixedCCIPReadLoadOnReadNodeConnection simulates production CCIP oracle load against a single
// read-node gRPC connection: multiple concurrent BatchGetLatestValues batches (4 nodes × 6 reads,
// each read doing object-metadata GetObject + SimulateTransaction) interleaved with background
// poller/indexer traffic (ReadFilterOwnedObjectIds, GetReferenceGasPrice, ReadObjectId). Validates
// that a connection pool relieves the stream bottleneck under this mixed workload.
func TestMixedCCIPReadLoadOnReadNodeConnection(t *testing.T) {
	t.Parallel()

	const (
		streamLimit = uint32(4)
		latency     = 150 * time.Millisecond
		poolSize    = 8
	)

	base := ccipLoadScenario{
		loadScenario: loadScenario{
			StreamLimit: streamLimit,
			CallLatency: latency,
		},
		BatchCount:            4,
		ReadsPerBatch:         6,
		IntraBatchConcurrency: 6,
		ObjectRefsPerRead:     2,
		BackgroundWorkers:     8,
	}

	single := runCCIPMixedLoadScenario(t, ccipLoadScenario{
		loadScenario: loadScenario{
			Name:        "single-read-node",
			PoolSize:    1,
			StreamLimit: base.StreamLimit,
			CallLatency: base.CallLatency,
		},
		BatchCount:            base.BatchCount,
		ReadsPerBatch:         base.ReadsPerBatch,
		IntraBatchConcurrency: base.IntraBatchConcurrency,
		ObjectRefsPerRead:     base.ObjectRefsPerRead,
		BackgroundWorkers:     base.BackgroundWorkers,
	})
	logCCIPResult(t, "single-read-node", single)

	pooled := runCCIPMixedLoadScenario(t, ccipLoadScenario{
		loadScenario: loadScenario{
			Name:        "pooled-read-node",
			PoolSize:    poolSize,
			StreamLimit: base.StreamLimit,
			CallLatency: base.CallLatency,
		},
		BatchCount:            base.BatchCount,
		ReadsPerBatch:         base.ReadsPerBatch,
		IntraBatchConcurrency: base.IntraBatchConcurrency,
		ObjectRefsPerRead:     base.ObjectRefsPerRead,
		BackgroundWorkers:     base.BackgroundWorkers,
	})
	logCCIPResult(t, "pooled-read-node", pooled)

	require.Zero(t, single.Errors, "single-connection mixed workload should not error")
	require.Zero(t, pooled.Errors, "pooled mixed workload should not error")

	countBackgroundByKind := func(workers, kind int) int {
		count := 0
		for i := 0; i < workers; i++ {
			if i%3 == kind {
				count++
			}
		}
		return count
	}

	expectedConfigReads := int64(base.BatchCount * base.ReadsPerBatch)
	expectedGetObject := expectedConfigReads*int64(base.ObjectRefsPerRead) +
		int64(countBackgroundByKind(base.BackgroundWorkers, 2)) +
		int64(single.PoolSize) // warmup ReadObjectId per connection
	expectedSimulate := expectedConfigReads + int64(single.PoolSize) // warmup SimulatePTB per connection
	expectedListOwned := int64(countBackgroundByKind(base.BackgroundWorkers, 0)) + int64(single.PoolSize)
	expectedGetEpoch := int64(countBackgroundByKind(base.BackgroundWorkers, 1)) + int64(single.PoolSize)

	require.Equal(t, expectedGetObject, single.GetObjectCalls)
	require.Equal(t, expectedSimulate, single.SimulateCalls)
	require.Equal(t, expectedListOwned, single.ListOwnedObjectsCalls)
	require.Equal(t, expectedGetEpoch, single.GetEpochCalls)

	require.LessOrEqual(t, single.ServerPeak, int64(streamLimit)+1,
		"single read-node connection should be limited to ~StreamLimit concurrent RPCs")
	require.Greater(t, pooled.ServerPeak, single.ServerPeak,
		"connection pool should achieve higher node-side concurrency under mixed CCIP load")

	ratio := float64(single.WallTime) / float64(pooled.WallTime)
	t.Logf("mixed-workload speedup (single/pooled) = %.2fx", ratio)
	require.Greater(t, ratio, 1.5,
		"connection pool should significantly outperform a single read-node connection under mixed CCIP load")
}

// TestConnectionPoolPreventsReadDeadlineExceeded is the faithful reproduction of the E2E failure, not just
// a throughput benchmark. It models the real symptom chain:
//
//	many concurrent CCIP config batches  ->  reads queue on the node's limited per-connection streams  ->
//	per-batch reads can't finish within the config poller's bgRefreshTimeout  ->  the batch returns
//	context.DeadlineExceeded  ->  in production that is a config refresh that fails and a merkle root that
//	never gets committed.
//
// With a SINGLE connection the reads serialize into waves and blow the deadline, so batches fail. With the
// connection POOL the same reads fan out across enough streams to finish well inside the deadline, so no
// batch fails. The assertion is therefore on the production-meaningful outcome (batch failures), not on a
// timing ratio: single-connection batches MUST time out and pooled batches MUST NOT.
func TestConnectionPoolPreventsReadDeadlineExceeded(t *testing.T) {
	t.Parallel()

	const (
		streamLimit = uint32(4)
		latency     = 150 * time.Millisecond
		poolSize    = 8
		// Sized so the single-connection burst serializes well past the deadline (~7s of work) while the
		// pooled burst finishes comfortably inside it (~0.9s of work). bgRefreshTimeout is 30s in
		// production; we scale both the load and the deadline down to keep the unit test fast and
		// deterministic while preserving the single-fails / pooled-passes separation.
		batchDeadline = 3 * time.Second
		batchCount    = 10
	)

	base := func(poolSize int, name string) ccipLoadScenario {
		return ccipLoadScenario{
			loadScenario: loadScenario{
				Name:        name,
				PoolSize:    poolSize,
				StreamLimit: streamLimit,
				CallLatency: latency,
			},
			BatchCount:        batchCount,
			ReadsPerBatch:     6,
			ObjectRefsPerRead: 2,
			BackgroundWorkers: 8,
			BatchDeadline:     batchDeadline,
		}
	}

	single := runCCIPMixedLoadScenario(t, base(1, "single-read-node-deadline"))
	logCCIPResult(t, "single-deadline", single)

	pooled := runCCIPMixedLoadScenario(t, base(poolSize, "pooled-read-node-deadline"))
	logCCIPResult(t, "pooled-deadline", pooled)

	// Reproduce the bug: a single connection cannot drain the config-read backlog within the per-batch
	// deadline, so batches time out. This is the E2E "merkle root never confirmed" failure in miniature.
	require.Positive(t, single.BatchErrors,
		"single connection should fail config batches by exceeding the per-batch deadline (the E2E symptom)")

	// Prove the fix: with the pool, every config batch completes within the deadline.
	require.Zero(t, pooled.BatchErrors,
		"connection pool should let every config batch finish within the per-batch deadline")

	// And the underlying mechanism: the single connection is pinned at ~StreamLimit node-side concurrency
	// while the pool reaches a multiple of it.
	require.LessOrEqual(t, single.ServerPeak, int64(streamLimit)+1,
		"single connection should be limited to ~StreamLimit concurrent RPCs at the node")
	require.Greater(t, pooled.ServerPeak, single.ServerPeak,
		"connection pool should achieve higher node-side concurrency")
}
