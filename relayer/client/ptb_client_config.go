package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/block-vision/sui-go-sdk/common/grpcconn"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/semaphore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
)

// PTBClientConfig configures a PTBClient with JSON-RPC and optional gRPC endpoints.
//
// RpcURL is always required for the hybrid migration period: methods not yet migrated
// to gRPC (events, transaction queries, DevInspect) continue to use JSON-RPC.
//
// When GrpcTarget and GrpcToken are both set, a gRPC client is also initialized for
// upcoming method migrations in Phase 3+.
type PTBClientConfig struct {
	RpcURL                string
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
	if cfg.RpcURL == "" {
		return nil, fmt.Errorf("rpc URL is required")
	}

	log.Infof("Creating new SUI PTBClient")

	maxConcurrentRequests := cfg.MaxConcurrentRequests
	if maxConcurrentRequests <= 0 {
		log.Warnw("maxConcurrentRequests is less than 0, setting to default value", "maxConcurrentRequests", maxConcurrentRequests)
		maxConcurrentRequests = 500
	}

	httpClient := &http.Client{
		Timeout: DefaultHTTPTimeout,
		Transport: &http.Transport{
			MaxConnsPerHost:     int(maxConcurrentRequests) * 2,
			MaxIdleConns:        int(maxConcurrentRequests) * 2,
			MaxIdleConnsPerHost: int(maxConcurrentRequests) * 2,
		},
	}
	jsonRpcClient := sui.NewSuiClientWithCustomClient(cfg.RpcURL, httpClient)

	var grpcClient *grpcconn.SuiGrpcClient
	if cfg.grpcEnabled() {
		log.Infow("Initializing Sui gRPC client", "target", cfg.GrpcTarget)
		grpcConfig := DefaultGrpcConfig(cfg.GrpcTarget, cfg.GrpcToken)
		grpcConfig.UseTLS = false
		grpcClient = NewSuiGrpcClient(grpcConfig)
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
		client:             jsonRpcClient,
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
