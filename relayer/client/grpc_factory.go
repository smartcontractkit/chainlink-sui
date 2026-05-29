package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/block-vision/sui-go-sdk/common/grpcconn"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	cache "github.com/patrickmn/go-cache"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"golang.org/x/sync/semaphore"
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

// PTBClientConfig configures a PTBClient with gRPC endpoints.
type PTBClientConfig struct {
	GrpcTarget            string
	GrpcToken             string
	MaxRetries            *int
	TransactionTimeout    time.Duration
	KeystoreService       loop.Keystore
	MaxConcurrentRequests int64
	DefaultRequestType    TransactionRequestType
}

func (cfg PTBClientConfig) grpcEnabled() bool {
	return cfg.GrpcTarget != "" && cfg.GrpcToken != ""
}

// NewPTBClientFromConfig creates a PTBClient from a full configuration.
func NewPTBClientFromConfig(log logger.Logger, cfg PTBClientConfig) (*PTBClient, error) {
	log.Infof("Creating new SUI PTBClient")

	maxConcurrentRequests := cfg.MaxConcurrentRequests
	if maxConcurrentRequests <= 0 {
		log.Warnw("maxConcurrentRequests is less than 0, setting to default value", "maxConcurrentRequests", maxConcurrentRequests)
		maxConcurrentRequests = 500
	}

	var grpcClient *grpcconn.SuiGrpcClient
	var moveModuleClient sui.ISuiAPI
	if cfg.grpcEnabled() {
		log.Infow("Initializing Sui gRPC client", "target", cfg.GrpcTarget)
		grpcConfig := DefaultGrpcConfig(cfg.GrpcTarget, cfg.GrpcToken)
		grpcConfig.UseTLS = false
		grpcClient = NewSuiGrpcClient(grpcConfig)
		moveModuleClient = sui.NewSuiClientWithCustomClient(
			fmt.Sprintf("http://%s", cfg.GrpcTarget),
			&http.Client{Timeout: cfg.TransactionTimeout},
		)
	} else {
		log.Info("gRPC client not configured; using JSON-RPC only")
	}

	log.Infof(
		"PTBClient config transactionTimeout: %s, maxConcurrentRequests: %d, grpcEnabled: %t",
		cfg.TransactionTimeout,
		maxConcurrentRequests,
		cfg.grpcEnabled(),
	)

	return &PTBClient{
		log:                log,
		moveModuleClient:   moveModuleClient,
		grpcClient:         grpcClient,
		maxRetries:         cfg.MaxRetries,
		transactionTimeout: cfg.TransactionTimeout,
		keystoreService:    cfg.KeystoreService,
		rateLimiter:        semaphore.NewWeighted(maxConcurrentRequests),
		defaultRequestType: cfg.DefaultRequestType,
		normalizedModules:  make(map[string]map[string]models.GetNormalizedMoveModuleResponse),
		cache:              cache.New(DefaultCacheExpiration, DefaultCacheCleanupInterval),
	}, nil
}
