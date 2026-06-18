package client

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/block-vision/sui-go-sdk/common/grpcconn"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
)

// pooledConn wraps a single Sui gRPC connection with a one-time, race-free initialization. The Sui SDK's
// SuiGrpcClient lazily (and non-atomically) creates its service stubs on first use, so we drive that
// initialization exactly once via sync.Once. After initialization the cached stubs are read-only, which
// makes the per-call service getters lock-free on the hot path.
type pooledConn struct {
	client  *grpcconn.SuiGrpcClient
	once    sync.Once
	initErr error
}

func (pc *pooledConn) ensure(ctx context.Context) error {
	pc.once.Do(func() {
		pc.initErr = pc.client.Connect(ctx)
	})
	return pc.initErr
}

// grpcConnPool holds several independent Sui gRPC connections and hands them out round-robin. Each
// connection is a separate HTTP/2 connection to the node, so spreading calls across the pool multiplies
// the number of concurrent in-flight streams available to the node and avoids a single connection's
// MaxConcurrentStreams becoming a head-of-line bottleneck under bursty CCIP read load.
type grpcConnPool struct {
	conns []*pooledConn
	idx   uint64
}

// newGrpcConnPool builds a pool of `size` connections using the provided factory. size < 1 is treated as 1.
func newGrpcConnPool(size int, factory func() *grpcconn.SuiGrpcClient) *grpcConnPool {
	if size < 1 {
		size = 1
	}
	conns := make([]*pooledConn, size)
	for i := range conns {
		conns[i] = &pooledConn{client: factory()}
	}
	return &grpcConnPool{conns: conns}
}

// Size returns the number of connections in the pool.
func (p *grpcConnPool) Size() int {
	if p == nil {
		return 0
	}
	return len(p.conns)
}

// next returns the next connection in round-robin order.
func (p *grpcConnPool) next() (*pooledConn, error) {
	if p == nil || len(p.conns) == 0 {
		return nil, errors.New("gRPC connection pool is not configured")
	}
	if len(p.conns) == 1 {
		return p.conns[0], nil
	}
	n := atomic.AddUint64(&p.idx, 1)
	//nolint:gosec // G115: modulo keeps the index within bounds
	return p.conns[n%uint64(len(p.conns))], nil
}

func (p *grpcConnPool) ledgerService(ctx context.Context) (suirpcv2.LedgerServiceClient, error) {
	pc, err := p.next()
	if err != nil {
		return nil, err
	}
	if err := pc.ensure(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC connection: %w", err)
	}
	return pc.client.LedgerService(ctx)
}

func (p *grpcConnPool) stateService(ctx context.Context) (suirpcv2.StateServiceClient, error) {
	pc, err := p.next()
	if err != nil {
		return nil, err
	}
	if err := pc.ensure(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC connection: %w", err)
	}
	return pc.client.StateService(ctx)
}

func (p *grpcConnPool) transactionExecutionService(ctx context.Context) (suirpcv2.TransactionExecutionServiceClient, error) {
	pc, err := p.next()
	if err != nil {
		return nil, err
	}
	if err := pc.ensure(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC connection: %w", err)
	}
	return pc.client.TransactionExecutionService(ctx)
}

func (p *grpcConnPool) movePackageService(ctx context.Context) (suirpcv2.MovePackageServiceClient, error) {
	pc, err := p.next()
	if err != nil {
		return nil, err
	}
	if err := pc.ensure(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize gRPC connection: %w", err)
	}
	return pc.client.MovePackageService(ctx)
}

func (p *grpcConnPool) Close() error {
	if p == nil {
		return nil
	}
	var closeErr error
	for _, pc := range p.conns {
		if pc == nil || pc.client == nil {
			continue
		}
		if err := pc.client.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
	}
	return closeErr
}

// Close closes all underlying gRPC connections in the pool.
func (c *PTBClient) Close() error {
	if c.connPool == nil {
		return nil
	}
	return c.connPool.Close()
}

// HealthCheckGrpc verifies gRPC connectivity by calling LedgerService.GetServiceInfo on a pooled connection.
func (c *PTBClient) HealthCheckGrpc(ctx context.Context) (chainID string, err error) {
	service, err := c.getLedgerService(ctx)
	if err != nil {
		return "", err
	}

	resp, err := service.GetServiceInfo(ctx, &suirpcv2.GetServiceInfoRequest{})
	if err != nil {
		return "", fmt.Errorf("GetServiceInfo failed: %w", err)
	}

	if resp.ChainId != nil {
		chainID = *resp.ChainId
	}

	return chainID, nil
}

// VerifyGrpcServices initializes all gRPC service stubs to verify connectivity.
func (c *PTBClient) VerifyGrpcServices(ctx context.Context) error {
	if _, err := c.getLedgerService(ctx); err != nil {
		return fmt.Errorf("LedgerService: %w", err)
	}
	if _, err := c.getStateService(ctx); err != nil {
		return fmt.Errorf("StateService: %w", err)
	}
	if _, err := c.getTransactionExecutionService(ctx); err != nil {
		return fmt.Errorf("TransactionExecutionService: %w", err)
	}
	if _, err := c.getMovePackageService(ctx); err != nil {
		return fmt.Errorf("MovePackageService: %w", err)
	}

	return nil
}

// getLedgerService returns a LedgerService client from the next pooled connection.
func (c *PTBClient) getLedgerService(ctx context.Context) (suirpcv2.LedgerServiceClient, error) {
	if c.connPool == nil {
		return nil, errors.New("gRPC client is not configured")
	}
	return c.connPool.ledgerService(ctx)
}

// getStateService returns a StateService client from the next pooled connection.
func (c *PTBClient) getStateService(ctx context.Context) (suirpcv2.StateServiceClient, error) {
	if c.connPool == nil {
		return nil, errors.New("gRPC client is not configured")
	}
	return c.connPool.stateService(ctx)
}

// getTransactionExecutionService returns a TransactionExecutionService client from the next pooled connection.
func (c *PTBClient) getTransactionExecutionService(ctx context.Context) (suirpcv2.TransactionExecutionServiceClient, error) {
	if c.connPool == nil {
		return nil, errors.New("gRPC client is not configured")
	}
	return c.connPool.transactionExecutionService(ctx)
}

// getMovePackageService returns a MovePackageService client from the next pooled connection.
func (c *PTBClient) getMovePackageService(ctx context.Context) (suirpcv2.MovePackageServiceClient, error) {
	if c.connPool == nil {
		return nil, errors.New("gRPC client is not configured")
	}
	return c.connPool.movePackageService(ctx)
}
