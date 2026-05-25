package client

import (
	"os"
	"strings"
)

const (
	EnvSuiRPCURL  = "SUI_RPC_URL"
	EnvGrpcTarget = "GRPC_TARGET"
	EnvGrpcToken  = "GRPC_TOKEN"
)

// IntegrationTestConfig holds connection settings for integration tests.
type IntegrationTestConfig struct {
	RpcURL     string
	GrpcTarget string
	GrpcToken  string
	UseLocal   bool
}

// LoadIntegrationTestConfig reads integration test settings from the environment.
// GRPC_TARGET and GRPC_TOKEN are required. SUI_RPC_URL defaults to http://127.0.0.1:9000.
func LoadIntegrationTestConfig() IntegrationTestConfig {
	rpcURL := os.Getenv(EnvSuiRPCURL)
	if rpcURL == "" {
		rpcURL = "http://127.0.0.1:9000"
	}

	return IntegrationTestConfig{
		RpcURL:     rpcURL,
		GrpcTarget: os.Getenv(EnvGrpcTarget),
		GrpcToken:  os.Getenv(EnvGrpcToken),
		UseLocal:   isLocalRPCURL(rpcURL),
	}
}

func isLocalRPCURL(rpcURL string) bool {
	return strings.Contains(rpcURL, "127.0.0.1") || strings.Contains(rpcURL, "localhost")
}
