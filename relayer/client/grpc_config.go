package client

import (
	"os"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
)

const (
	envGrpcTarget = EnvGrpcTarget
	envGrpcToken  = EnvGrpcToken
)

// ResolveGrpcConfig returns gRPC target and token, preferring explicit node config
// and falling back to GRPC_TARGET / GRPC_TOKEN environment variables.
func ResolveGrpcConfig(nodeGrpcTarget, nodeGrpcToken *string) (target, token string) {
	if nodeGrpcTarget != nil {
		target = *nodeGrpcTarget
	}
	if nodeGrpcToken != nil {
		token = *nodeGrpcToken
	}

	if target == "" {
		target = os.Getenv(envGrpcTarget)
	}
	if token == "" {
		token = os.Getenv(envGrpcToken)
	}

	return target, token
}

// PTBClientConfigFromNode builds a PTBClientConfig from node connection settings.
func PTBClientConfigFromNode(
	rpcURL string,
	nodeGrpcTarget, nodeGrpcToken *string,
	maxRetries *int,
	transactionTimeout time.Duration,
	keystoreService loop.Keystore,
	maxConcurrentRequests int64,
	defaultRequestType TransactionRequestType,
) PTBClientConfig {
	grpcTarget, grpcToken := ResolveGrpcConfig(nodeGrpcTarget, nodeGrpcToken)

	return PTBClientConfig{
		RpcURL:                rpcURL,
		GrpcTarget:            grpcTarget,
		GrpcToken:             grpcToken,
		MaxRetries:            maxRetries,
		TransactionTimeout:    transactionTimeout,
		KeystoreService:       keystoreService,
		MaxConcurrentRequests: maxConcurrentRequests,
		DefaultRequestType:    defaultRequestType,
	}
}
