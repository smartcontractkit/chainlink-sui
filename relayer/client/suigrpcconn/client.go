package suigrpcconn

import (
	"context"
	"fmt"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type SuiGrpcClient struct {
	conn                         *GrpcConn
	nameService                  suirpcv2.NameServiceClient
	ledgerService                suirpcv2.LedgerServiceClient
	stateService                 suirpcv2.StateServiceClient
	movePackageService           suirpcv2.MovePackageServiceClient
	subscriptionService          suirpcv2.SubscriptionServiceClient
	transactionExecutionService  suirpcv2.TransactionExecutionServiceClient
	signatureVerificationService suirpcv2.SignatureVerificationServiceClient
}

func NewSuiGrpcClient(target string, useTLS bool, opts ...GrpcConnOption) *SuiGrpcClient {
	allOpts := append([]GrpcConnOption{WithTransportSecurity(useTLS)}, opts...)
	conn := NewGrpcConn(target, allOpts...)
	return &SuiGrpcClient{conn: conn}
}

func NewSuiGrpcClientWithAuth(target, token string, useTLS bool, opts ...GrpcConnOption) *SuiGrpcClient {
	authOpts := []GrpcConnOption{
		WithTransportSecurity(useTLS),
		WithDialOptions(
			grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
				md := metadata.Pairs(
					"authorization", "Bearer "+token,
					"x-api-key", token,
					"x-token", token,
				)
				ctx = metadata.NewOutgoingContext(ctx, md)
				return invoker(ctx, method, req, reply, cc, callOpts...)
			}),
			grpc.WithChainStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				md := metadata.Pairs(
					"authorization", "Bearer "+token,
					"x-api-key", token,
					"x-token", token,
				)
				ctx = metadata.NewOutgoingContext(ctx, md)
				return streamer(ctx, desc, cc, method, opts...)
			}),
		),
	}

	allOpts := append(authOpts, opts...)
	conn := NewGrpcConn(target, allOpts...)
	return &SuiGrpcClient{conn: conn}
}

func (c *SuiGrpcClient) Connect(ctx context.Context) error {
	if err := c.conn.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}

	grpcConn, err := c.conn.GetConn(ctx)
	if err != nil {
		return fmt.Errorf("failed to get connection: %v", err)
	}

	c.nameService = suirpcv2.NewNameServiceClient(grpcConn)
	c.ledgerService = suirpcv2.NewLedgerServiceClient(grpcConn)
	c.stateService = suirpcv2.NewStateServiceClient(grpcConn)
	c.movePackageService = suirpcv2.NewMovePackageServiceClient(grpcConn)
	c.subscriptionService = suirpcv2.NewSubscriptionServiceClient(grpcConn)
	c.transactionExecutionService = suirpcv2.NewTransactionExecutionServiceClient(grpcConn)
	c.signatureVerificationService = suirpcv2.NewSignatureVerificationServiceClient(grpcConn)
	return nil
}

func (c *SuiGrpcClient) Close() error {
	return c.conn.Close()
}

func (c *SuiGrpcClient) NameService(ctx context.Context) (suirpcv2.NameServiceClient, error) {
	if c.nameService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.nameService, nil
}

func (c *SuiGrpcClient) LedgerService(ctx context.Context) (suirpcv2.LedgerServiceClient, error) {
	if c.ledgerService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.ledgerService, nil
}

func (c *SuiGrpcClient) StateService(ctx context.Context) (suirpcv2.StateServiceClient, error) {
	if c.stateService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.stateService, nil
}

func (c *SuiGrpcClient) MovePackageService(ctx context.Context) (suirpcv2.MovePackageServiceClient, error) {
	if c.movePackageService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.movePackageService, nil
}

func (c *SuiGrpcClient) SubscriptionService(ctx context.Context) (suirpcv2.SubscriptionServiceClient, error) {
	if c.subscriptionService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.subscriptionService, nil
}

func (c *SuiGrpcClient) TransactionExecutionService(ctx context.Context) (suirpcv2.TransactionExecutionServiceClient, error) {
	if c.transactionExecutionService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.transactionExecutionService, nil
}

func (c *SuiGrpcClient) SignatureVerificationService(ctx context.Context) (suirpcv2.SignatureVerificationServiceClient, error) {
	if c.signatureVerificationService == nil {
		if err := c.Connect(ctx); err != nil {
			return nil, err
		}
	}
	return c.signatureVerificationService, nil
}

func (c *SuiGrpcClient) CallWithRetry(ctx context.Context, method string, req interface{}, reply interface{}) error {
	return c.conn.Call(ctx, method, req, reply)
}

func (c *SuiGrpcClient) GetConnection(ctx context.Context) (*grpc.ClientConn, error) {
	return c.conn.GetConn(ctx)
}
