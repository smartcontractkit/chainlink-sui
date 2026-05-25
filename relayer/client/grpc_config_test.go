package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveGrpcConfig(t *testing.T) {
	t.Setenv(envGrpcTarget, "env-target:443")
	t.Setenv(envGrpcToken, "env-token")

	nodeTarget := "node-target:443"
	nodeToken := "node-token"

	target, token := ResolveGrpcConfig(&nodeTarget, &nodeToken)
	require.Equal(t, nodeTarget, target)
	require.Equal(t, nodeToken, token)

	target, token = ResolveGrpcConfig(nil, nil)
	require.Equal(t, "env-target:443", target)
	require.Equal(t, "env-token", token)
}

func TestPTBClientConfig_grpcEnabled(t *testing.T) {
	require.False(t, PTBClientConfig{GrpcTarget: "host:443"}.grpcEnabled())
	require.False(t, PTBClientConfig{GrpcToken: "token"}.grpcEnabled())
	require.True(t, PTBClientConfig{GrpcTarget: "host:443", GrpcToken: "token"}.grpcEnabled())
}
