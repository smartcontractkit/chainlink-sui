//go:build unit

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTOMLConfigValidateConfig_RequiresAtLeastOneNode(t *testing.T) {
	c := &TOMLConfig{ChainID: strPtr("2"), Nodes: nil}

	err := c.ValidateConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Nodes: missing: must have at least one node")
}

func TestTOMLConfigValidateConfig_RequiresChainID(t *testing.T) {
	c := &TOMLConfig{ChainID: nil, Nodes: NodeConfigs{validNode()}}

	err := c.ValidateConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ChainID: missing")
}

func TestTOMLConfigValidateConfig_ValidPasses(t *testing.T) {
	c := &TOMLConfig{ChainID: strPtr("2"), Nodes: NodeConfigs{validNode()}}

	assert.NoError(t, c.ValidateConfig())
}

// Reproduces the layering failure mode at the chain level: a node missing GrpcTarget must be
// rejected by ValidateConfig, which is the fail-fast path that prevents the nil deref in NewRelayer.
func TestTOMLConfigValidateConfig_RejectsNodeWithoutGrpc(t *testing.T) {
	n := validNode()
	n.Name = strPtr("2_primary_rpcProxy_0")
	n.GrpcTarget = nil
	c := &TOMLConfig{ChainID: strPtr("2"), Nodes: NodeConfigs{n}}

	err := c.ValidateConfig()
	require.Error(t, err)
	msg := err.Error()
	assert.Contains(t, msg, "GrpcTarget: missing")
	assert.Contains(t, msg, `(node "2_primary_rpcProxy_0")`)
}
