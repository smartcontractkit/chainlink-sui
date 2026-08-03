//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonConfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

func strPtr(s string) *string { return &s }

func validNode() *NodeConfig {
	return &NodeConfig{
		Name:       strPtr("primary"),
		URL:        commonConfig.MustParseURL("https://sui-testnet.example.com:443"),
		GrpcTarget: strPtr("sui-testnet.example.com:443"),
		GrpcToken:  strPtr("token"),
	}
}

func TestNodeConfigValidateConfig_Valid(t *testing.T) {
	assert.NoError(t, validNode().ValidateConfig())
}

func TestNodeConfigValidateConfig_MissingGrpcFields(t *testing.T) {
	n := validNode()
	n.GrpcTarget = nil
	n.GrpcToken = nil

	err := n.ValidateConfig()
	require.Error(t, err)
	msg := err.Error()
	assert.Contains(t, msg, "GrpcTarget: missing")
	assert.Contains(t, msg, "GrpcToken: missing")
	assert.Contains(t, msg, `(node "primary")`)
}

func TestNodeConfigValidateConfig_EmptyGrpcTarget(t *testing.T) {
	n := validNode()
	n.GrpcTarget = strPtr("")

	err := n.ValidateConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "GrpcTarget: empty")
}

func TestNodeConfigValidateConfig_EmptyGrpcToken(t *testing.T) {
	n := validNode()
	n.GrpcToken = strPtr("")

	err := n.ValidateConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "GrpcToken: empty")
}

// Existing Name and URL requirements must still be enforced alongside the new gRPC checks.
func TestNodeConfigValidateConfig_StillRequiresNameAndURL(t *testing.T) {
	n := validNode()
	n.Name = nil
	n.URL = nil

	err := n.ValidateConfig()
	require.Error(t, err)
	msg := err.Error()
	assert.Contains(t, msg, "Name: missing")
	assert.Contains(t, msg, "URL: missing")
}
