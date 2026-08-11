package client

import (
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/semaphore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-sui/relayer/client/suigrpcconn"
)

// gRPC support uses the local suigrpcconn helper, which fixes TLS dial options for public Sui fullnodes.

const (
	DefaultGrpcMaxMsgSize = 20 * 1024 * 1024

	// DefaultMaxGrpcConnections is the number of independent gRPC connections each client opens to the
	// node. Spreading RPCs round-robin across several connections multiplies the available concurrent
	// HTTP/2 streams so a single connection's stream limit does not throttle bursty reads.
	DefaultMaxGrpcConnections = 128
)

// GrpcClientConfig holds configuration for a Sui gRPC client connection.
type GrpcClientConfig struct {
	Target string
	Token  string
	UseTLS bool
}

// DefaultGrpcConfig returns sensible defaults for a gRPC endpoint and auth token.
func DefaultGrpcConfig(target, token string) GrpcClientConfig {
	return GrpcClientConfig{
		Target: target,
		Token:  token,
		UseTLS: true,
	}
}

// NewSuiGrpcClient creates an authenticated Sui gRPC connection.
func NewSuiGrpcClient(config GrpcClientConfig) *suigrpcconn.Connection {
	return suigrpcconn.NewConnectionWithAuth(config.Target, config.Token, config.UseTLS)
}

// PTBClientConfig configures a PTBClient with gRPC endpoints.
type PTBClientConfig struct {
	GrpcTarget            string
	GrpcToken             string
	MaxRetries            *int
	TransactionTimeout    time.Duration
	KeystoreService       loop.Keystore
	MaxConcurrentRequests int64
	// MaxGrpcConnections is the size of the round-robin gRPC connection pool. Zero means use
	// DefaultMaxGrpcConnections.
	MaxGrpcConnections int
	// ObjectCache, when set, caches version-stable object reference metadata to avoid redundant GetObject
	// RPCs on the read hot path. Nil disables object-metadata caching.
	ObjectCache        ObjectMetadataCache
	DefaultRequestType TransactionRequestType
}

func (cfg PTBClientConfig) grpcEnabled() bool {
	return cfg.GrpcTarget != "" && cfg.GrpcToken != ""
}

// grpcTargetUsesTLS reports whether the gRPC dial should use TLS for the given host:port target.
// Public Sui full nodes serve gRPC over TLS on port 443; local nodes use plain gRPC on 9000.
func grpcTargetUsesTLS(target string) bool {
	_, port, err := net.SplitHostPort(target)
	if err != nil {
		if target == "" || isLocalGrpcHost(target) {
			return false
		}
		// Bare public hostnames default to port 443 when derived from https RPC URLs.
		return true
	}

	if port == "" || port == "443" {
		return true
	}

	// Local and other non-443 ports use plain gRPC.
	return false
}

func isLocalGrpcHost(host string) bool {
	switch host {
	case "127.0.0.1", "localhost", "::1":
		return true
	default:
		return strings.HasPrefix(host, "[::1]")
	}
}

// NewPTBClientFromConfig creates a PTBClient from a full configuration.
func NewPTBClientFromConfig(log logger.Logger, cfg PTBClientConfig) (*PTBClient, error) {
	log.Infof("Creating new SUI PTBClient")

	maxConcurrentRequests := cfg.MaxConcurrentRequests
	if maxConcurrentRequests <= 0 {
		log.Warnw("maxConcurrentRequests is less than 0, setting to default value", "maxConcurrentRequests", maxConcurrentRequests)
		maxConcurrentRequests = 500
	}

	maxGrpcConnections := cfg.MaxGrpcConnections
	if maxGrpcConnections <= 0 {
		maxGrpcConnections = DefaultMaxGrpcConnections
	}

	var connPool *grpcConnPool
	var moveModuleClient sui.ISuiAPI
	if cfg.grpcEnabled() {
		useTLS := grpcTargetUsesTLS(cfg.GrpcTarget)

		log.Infow("Initializing Sui gRPC client", "target", cfg.GrpcTarget, "connections", maxGrpcConnections, "useTLS", useTLS)

		grpcConfig := DefaultGrpcConfig(cfg.GrpcTarget, cfg.GrpcToken)
		grpcConfig.UseTLS = useTLS
		connPool = newGrpcConnPool(maxGrpcConnections, func() *suigrpcconn.Connection {
			return NewSuiGrpcClient(grpcConfig)
		})

		// TODO: remove this to use the gRPC client only
		jsonRPCScheme := "http"
		if useTLS {
			jsonRPCScheme = "https"
		}
		moveModuleClient = sui.NewSuiClientWithCustomClient(
			jsonRPCScheme+"://"+cfg.GrpcTarget,
			&http.Client{Timeout: cfg.TransactionTimeout},
		)
	} else {
		log.Info("gRPC client not configured; using JSON-RPC only")
	}

	log.Infof(
		"PTBClient config transactionTimeout: %s, maxConcurrentRequests: %d, grpcConnections: %d, grpcEnabled: %t",
		cfg.TransactionTimeout,
		maxConcurrentRequests,
		maxGrpcConnections,
		cfg.grpcEnabled(),
	)

	// TODO: use a "read-only" signer for this operation instead of creating a new one
	// each time, it can be empty since gas selection is disabled
	seed := make([]byte, 32)
	if _, seedErr := rand.Read(seed); seedErr != nil {
		return nil, fmt.Errorf("failed to generate random seed: %w", seedErr)
	}
	devInspectSigner := signer.NewSigner(seed)

	return &PTBClient{
		log:                log,
		moveModuleClient:   moveModuleClient,
		connPool:           connPool,
		maxRetries:         cfg.MaxRetries,
		transactionTimeout: cfg.TransactionTimeout,
		keystoreService:    cfg.KeystoreService,
		devInspectSigner:   devInspectSigner,
		rateLimiter:        semaphore.NewWeighted(maxConcurrentRequests),
		defaultRequestType: cfg.DefaultRequestType,
		normalizedModules:  make(map[string]map[string]models.GetNormalizedMoveModuleResponse),
		cache:              cache.New(DefaultCacheExpiration, DefaultCacheCleanupInterval),
		objectCache:        cfg.ObjectCache,
	}, nil
}
