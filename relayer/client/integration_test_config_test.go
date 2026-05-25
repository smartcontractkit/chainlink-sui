package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadIntegrationTestConfig(t *testing.T) {
	t.Setenv(EnvSuiRPCURL, "http://127.0.0.1:9000")
	t.Setenv(EnvGrpcTarget, "grpc.example.com:443")
	t.Setenv(EnvGrpcToken, "token")

	cfg := LoadIntegrationTestConfig()
	require.Equal(t, "http://127.0.0.1:9000", cfg.RpcURL)
	require.Equal(t, "grpc.example.com:443", cfg.GrpcTarget)
	require.Equal(t, "token", cfg.GrpcToken)
	require.True(t, cfg.UseLocal)
}

func TestLoadIntegrationTestConfig_defaultsLocalRPC(t *testing.T) {
	t.Setenv(EnvSuiRPCURL, "")

	cfg := LoadIntegrationTestConfig()
	require.Equal(t, "http://127.0.0.1:9000", cfg.RpcURL)
	require.True(t, cfg.UseLocal)
}

func TestIsLocalRPCURL(t *testing.T) {
	require.True(t, isLocalRPCURL("http://127.0.0.1:9000"))
	require.True(t, isLocalRPCURL("http://localhost:9000"))
	require.False(t, isLocalRPCURL("https://fullnode.testnet.sui.io:443"))
}
