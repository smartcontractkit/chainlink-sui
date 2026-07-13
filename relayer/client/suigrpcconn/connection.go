package suigrpcconn

import (
	"context"
	"fmt"
	"sync"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

const (
	defaultKeepAlive        = 60 * time.Second
	defaultKeepaliveTimeout = 20 * time.Second
	defaultMaxRecvMsgSize   = 20 * 1024 * 1024
)

// Connection owns a single gRPC connection to a Sui node and the typed service stubs built on
// top of it. It supports both TLS (public fullnodes on :443) and plaintext (local nodes on :9000)
// transports, selected via the useTLS flag at construction.
//
// Each typed service getter lazily connects on first use; in practice grpcConnPool drives that
// initialization exactly once via sync.Once, after which the cached stubs are read-only and the
// getters are lock-free on the hot path.
type Connection struct {
	mu     sync.RWMutex
	target string
	conn   *grpc.ClientConn
	useTLS bool
	// interceptors are composed once at construction (e.g. auth) and always appended after the
	// transport options, so dial-option ordering can never be inverted by a caller.
	interceptors []grpc.DialOption

	ledgerService               suirpcv2.LedgerServiceClient
	stateService                suirpcv2.StateServiceClient
	movePackageService          suirpcv2.MovePackageServiceClient
	transactionExecutionService suirpcv2.TransactionExecutionServiceClient
}

// NewConnectionWithAuth builds a client that injects the given auth token on every RPC. useTLS
// selects TLS (public fullnodes) vs plaintext (local nodes) transport credentials.
func NewConnectionWithAuth(target, token string, useTLS bool) *Connection {
	return &Connection{
		target:       target,
		useTLS:       useTLS,
		interceptors: []grpc.DialOption{grpc.WithUnaryInterceptor(authUnaryInterceptor(token))},
	}
}

// authMetadata mirrors the token across the header names different Sui gRPC gateways expect.
func authMetadata(token string) metadata.MD {
	return metadata.Pairs(
		"authorization", "Bearer "+token,
		"x-api-key", token,
		"x-token", token,
	)
}

func authUnaryInterceptor(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
		ctx = metadata.NewOutgoingContext(ctx, authMetadata(token))
		return invoker(ctx, method, req, reply, cc, callOpts...)
	}
}

// dialOptions assembles the dial options in a fixed, order-independent way: transport credentials
// first, then connection-level options, then any interceptors. Unlike the previous option-mutator
// model, there is no way for a caller to reorder transport credentials behind the interceptors.
func (c *Connection) dialOptions() []grpc.DialOption {
	var creds credentials.TransportCredentials
	if c.useTLS {
		creds = credentials.NewTLS(nil)
	} else {
		creds = insecure.NewCredentials()
	}

	opts := make([]grpc.DialOption, 0, 3)
	opts = append(opts, grpc.WithTransportCredentials(creds))
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                defaultKeepAlive,
		Timeout:             defaultKeepaliveTimeout,
		PermitWithoutStream: !c.useTLS,
	}))
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(defaultMaxRecvMsgSize)))

	return append(opts, c.interceptors...)
}

// Connect dials the node (once) and constructs the typed service stubs. It is idempotent: a second
// call while already connected is a no-op, which is what grpcConnPool's sync.Once relies on.
func (c *Connection) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return nil
	}

	conn, err := grpc.NewClient(c.target, c.dialOptions()...)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.conn = conn
	c.ledgerService = suirpcv2.NewLedgerServiceClient(conn)
	c.stateService = suirpcv2.NewStateServiceClient(conn)
	c.movePackageService = suirpcv2.NewMovePackageServiceClient(conn)
	c.transactionExecutionService = suirpcv2.NewTransactionExecutionServiceClient(conn)
	return nil
}

func (c *Connection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

func (c *Connection) LedgerService(ctx context.Context) (suirpcv2.LedgerServiceClient, error) {
	if c.ledgerService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.ledgerService, nil
}

func (c *Connection) StateService(ctx context.Context) (suirpcv2.StateServiceClient, error) {
	if c.stateService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.stateService, nil
}

func (c *Connection) MovePackageService(ctx context.Context) (suirpcv2.MovePackageServiceClient, error) {
	if c.movePackageService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.movePackageService, nil
}

func (c *Connection) TransactionExecutionService(ctx context.Context) (suirpcv2.TransactionExecutionServiceClient, error) {
	if c.transactionExecutionService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.transactionExecutionService, nil
}
