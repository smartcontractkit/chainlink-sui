package client

import (
	"time"

	"github.com/block-vision/sui-go-sdk/common/grpcconn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// gRPC support requires sui-go-sdk with grpcconn (main branch or v2+):
//
//	GOPROXY=direct go get github.com/block-vision/sui-go-sdk@main
//	go mod tidy

const (
	DefaultGrpcTimeout    = 30 * time.Second
	DefaultGrpcRetryCount = 3
	DefaultGrpcMaxMsgSize = 20 * 1024 * 1024
)

// GrpcClientConfig holds configuration for a Sui gRPC client connection.
type GrpcClientConfig struct {
	Target     string
	Token      string
	Timeout    time.Duration
	RetryCount int
	UseTLS     bool
}

// DefaultGrpcConfig returns sensible defaults for a gRPC endpoint and auth token.
func DefaultGrpcConfig(target, token string) GrpcClientConfig {
	return GrpcClientConfig{
		Target:     target,
		Token:      token,
		Timeout:    DefaultGrpcTimeout,
		RetryCount: DefaultGrpcRetryCount,
		UseTLS:     true,
	}
}

// NewSuiGrpcClient creates an authenticated Sui gRPC client.
func NewSuiGrpcClient(config GrpcClientConfig) *grpcconn.SuiGrpcClient {
	opts := []grpcconn.GrpcConnOption{
		grpcconn.WithTimeout(config.Timeout),
		grpcconn.WithRetryCount(config.RetryCount),
	}

	if config.UseTLS {
		opts = append(opts, grpcconn.WithDialOptions(
			grpc.WithTransportCredentials(credentials.NewTLS(nil)),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(DefaultGrpcMaxMsgSize)),
		))
	}

	return grpcconn.NewSuiGrpcClientWithAuth(config.Target, config.Token, opts...)
}
