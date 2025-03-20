//go:build unit

package chainreader

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

func TestChainReaderBindUnbind(t *testing.T) {
	mockLogger := logger.Test(t)
	config := ChainReaderConfig{
		Modules: map[string]*ChainReaderModule{
			"testContract": {
				Name: "counter",
				Functions: map[string]*ChainReaderFunction{
					"getValue": {
						Name: "get_value",
					},
				},
			},
		},
	}

	// Create chain reader with nil client (we're just testing binding/unbinding)
	reader := NewChainReader(mockLogger, nil, config)

	// Test binding a contract
	err := reader.Bind(context.Background(), []types.BoundContract{
		{
			Name:    "testContract",
			Address: "0x1234567890abcdef1234567890abcdef",
		},
	})
	assert.NoError(t, err, "Binding should succeed")

	// Test unbinding the contract
	err = reader.Unbind(context.Background(), []types.BoundContract{
		{
			Name:    "testContract",
			Address: "0x1234567890abcdef1234567890abcdef",
		},
	})
	assert.NoError(t, err, "Unbinding should succeed")

	// Test binding a contract with invalid address
	err = reader.Bind(context.Background(), []types.BoundContract{
		{
			Name:    "testContract",
			Address: "invalid-address",
		},
	})
	assert.Error(t, err, "Binding with invalid address should fail")
}
