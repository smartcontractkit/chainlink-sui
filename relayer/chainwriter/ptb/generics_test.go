//go:build unit

package ptb_test

import (
	"testing"

	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestResolveGenericTypeTags(t *testing.T) {
	t.Parallel()

	// Create a dummy config for testing
	writerConfig := cwConfig.ChainWriterConfig{
		Modules: map[string]*cwConfig.ChainWriterModule{},
	}

	// Create a mock client
	mockClient := &testutils.FakeSuiPTBClient{
		Status: client.TransactionResult{
			Status: "success",
			Error:  "",
		},
	}

	// Create a test logger
	log := logger.Test(t)

	// Create PTBConstructor using NewPTBConstructor
	ptbService := ptb.NewPTBConstructor(writerConfig, mockClient, log)

	tests := []struct {
		name        string
		params      []codec.SuiFunctionParam
		arguments   cwConfig.Arguments
		expectError bool
		errorMsg    string
		expectedLen int
		validate    func(t *testing.T, result []transaction.TypeTag)
	}{
		{
			name:        "no parameters",
			params:      []codec.SuiFunctionParam{},
			arguments:   cwConfig.Arguments{},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "no generic parameters",
			params: []codec.SuiFunctionParam{
				{Name: "value", Type: "u64", IsGeneric: false},
			},
			arguments:   cwConfig.Arguments{},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "single generic parameter",
			params: []codec.SuiFunctionParam{
				{Name: "coin", Type: "Coin<T>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"coin": "0x2::sui::SUI",
				},
			},
			expectError: false,
			expectedLen: 1,
			validate: func(t *testing.T, result []transaction.TypeTag) {
				t.Helper()
				assert.Len(t, result, 1)
				assert.NotNil(t, result[0].Struct)
				assert.Equal(t, "sui", result[0].Struct.Module)
				assert.Equal(t, "SUI", result[0].Struct.Name)
			},
		},
		{
			name: "multiple generic parameters with same type",
			params: []codec.SuiFunctionParam{
				{Name: "coin1", Type: "Coin<T>", IsGeneric: true},
				{Name: "coin2", Type: "Coin<T>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"coin1": "0x2::sui::SUI",
					"coin2": "0x2::sui::SUI",
				},
			},
			expectError: false,
			expectedLen: 1, // Should deduplicate
		},
		{
			name: "multiple generic parameters with different types",
			params: []codec.SuiFunctionParam{
				{Name: "coin1", Type: "Coin<T>", IsGeneric: true},
				{Name: "coin2", Type: "Coin<U>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"coin1": "0x2::sui::SUI",
					"coin2": "0x2::coin::Coin",
				},
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name: "generic parameter with empty name",
			params: []codec.SuiFunctionParam{
				{Name: "", Type: "Coin<T>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{},
			},
			expectError: true,
			errorMsg:    "generic parameter missing name",
		},
		{
			name: "missing type in ArgTypes",
			params: []codec.SuiFunctionParam{
				{Name: "coin", Type: "Coin<T>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{},
			},
			expectError: true,
			errorMsg:    "generic parameter \"coin\" not found in ArgTypes",
		},
		{
			name: "invalid type format",
			params: []codec.SuiFunctionParam{
				{Name: "coin", Type: "Coin<T>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"coin": "invalid_type",
				},
			},
			expectError: true,
			errorMsg:    "failed to create type tag",
		},
		{
			name: "vector type parameter",
			params: []codec.SuiFunctionParam{
				{Name: "coins", Type: "vector<Coin<T>>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"coins": "vector<0x2::sui::SUI>",
				},
			},
			expectError: false,
			expectedLen: 1,
			validate: func(t *testing.T, result []transaction.TypeTag) {
				t.Helper()
				assert.Len(t, result, 1)
				// Should extract the inner type from the vector
				assert.NotNil(t, result[0].Struct)
				assert.Equal(t, "sui", result[0].Struct.Module)
				assert.Equal(t, "SUI", result[0].Struct.Name)
			},
		},
		{
			name: "mixed generic and non-generic parameters",
			params: []codec.SuiFunctionParam{
				{Name: "value", Type: "u64", IsGeneric: false},
				{Name: "coin", Type: "Coin<T>", IsGeneric: true},
				{Name: "amount", Type: "u64", IsGeneric: false},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"coin": "0x2::sui::SUI",
				},
			},
			expectError: false,
			expectedLen: 1,
		},
		{
			name: "empty vector type parameter",
			params: []codec.SuiFunctionParam{
				{Name: "items", Type: "vector<T>", IsGeneric: true},
			},
			arguments: config.Arguments{
				ArgTypes: map[string]string{
					"items": "vector<>",
				},
			},
			expectError: true,
			errorMsg:    "failed to extract vector inner type",
		},
		{
			name: "nested generic types",
			params: []codec.SuiFunctionParam{
				{Name: "nested_coins", Type: "vector<Coin<T>>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"nested_coins": "vector<0x2::coin::Coin>",
				},
			},
			expectError: false,
			expectedLen: 1,
			validate: func(t *testing.T, result []transaction.TypeTag) {
				t.Helper()
				assert.Len(t, result, 1)
				assert.NotNil(t, result[0].Struct)
				assert.Equal(t, "coin", result[0].Struct.Module)
				assert.Equal(t, "Coin", result[0].Struct.Name)
			},
		},
		{
			name: "multiple different generic types complex",
			params: []codec.SuiFunctionParam{
				{Name: "sui_coin", Type: "Coin<T>", IsGeneric: true},
				{Name: "link_coin", Type: "Coin<U>", IsGeneric: true},
				{Name: "coin_vector", Type: "vector<Coin<V>>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"sui_coin":    "0x2::sui::SUI",
					"link_coin":   "0x2::link::LINK",
					"coin_vector": "vector<0x2::coin::Coin>",
				},
			},
			expectError: false,
			expectedLen: 3,
			validate: func(t *testing.T, result []transaction.TypeTag) {
				t.Helper()
				assert.Len(t, result, 3)

				// Check SUI coin type
				assert.NotNil(t, result[0].Struct)
				assert.Equal(t, "sui", result[0].Struct.Module)
				assert.Equal(t, "SUI", result[0].Struct.Name)

				// Check LINK coin type
				assert.NotNil(t, result[1].Struct)
				assert.Equal(t, "link", result[1].Struct.Module)
				assert.Equal(t, "LINK", result[1].Struct.Name)

				// Check vector inner type (Coin)
				assert.NotNil(t, result[2].Struct)
				assert.Equal(t, "coin", result[2].Struct.Module)
				assert.Equal(t, "Coin", result[2].Struct.Name)
			},
		},
		{
			name: "invalid address format",
			params: []codec.SuiFunctionParam{
				{Name: "token", Type: "Token<T>", IsGeneric: true},
			},
			arguments: cwConfig.Arguments{
				ArgTypes: map[string]string{
					"token": "0xZZZ::token::Token",
				},
			},
			expectError: true,
			errorMsg:    "failed to convert package address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := ptbService.ResolveGenericTypeTags(tt.params, tt.arguments)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)

				if tt.validate != nil {
					tt.validate(t, result)
				}
			}
		})
	}
}

func TestResolveGenericTypeTags_OrderingAndDeduplication(t *testing.T) {
	t.Parallel()

	// Create a dummy config for testing
	writerConfig := cwConfig.ChainWriterConfig{
		Modules: map[string]*cwConfig.ChainWriterModule{},
	}

	// Create a mock client
	mockClient := &testutils.FakeSuiPTBClient{
		Status: client.TransactionResult{
			Status: "success",
			Error:  "",
		},
	}

	// Create a test logger
	log := logger.Test(t)

	// Create PTBConstructor using NewPTBConstructor
	ptbService := ptb.NewPTBConstructor(writerConfig, mockClient, log)

	// Test that type tags are deduplicated but order is preserved
	params := []codec.SuiFunctionParam{
		{Name: "coin1", Type: "Coin<T>", IsGeneric: true},
		{Name: "coin2", Type: "Coin<U>", IsGeneric: true},
		{Name: "coin3", Type: "Coin<T>", IsGeneric: true}, // Same as coin1
		{Name: "coin4", Type: "Coin<V>", IsGeneric: true},
		{Name: "coin5", Type: "Coin<U>", IsGeneric: true}, // Same as coin2
	}

	arguments := cwConfig.Arguments{
		ArgTypes: map[string]string{
			"coin1": "0x2::sui::SUI",
			"coin2": "0x2::coin::Coin",
			"coin3": "0x2::sui::SUI", // Duplicate
			"coin4": "0x2::link::LINK",
			"coin5": "0x2::coin::Coin", // Duplicate
		},
	}

	result, err := ptbService.ResolveGenericTypeTags(params, arguments)
	require.NoError(t, err)

	// Should have 3 unique types in the order they first appeared
	assert.Len(t, result, 3)

	// Check first type (0x2::sui::SUI)
	assert.NotNil(t, result[0].Struct)
	assert.Equal(t, "sui", result[0].Struct.Module)
	assert.Equal(t, "SUI", result[0].Struct.Name)

	// Check second type (0x2::coin::Coin)
	assert.NotNil(t, result[1].Struct)
	assert.Equal(t, "coin", result[1].Struct.Module)
	assert.Equal(t, "Coin", result[1].Struct.Name)

	// Check third type (0x2::link::LINK)
	assert.NotNil(t, result[2].Struct)
	assert.Equal(t, "link", result[2].Struct.Module)
	assert.Equal(t, "LINK", result[2].Struct.Name)
}

func TestResolveGenericTypeTags_EmptyArgTypes(t *testing.T) {
	t.Parallel()

	// Create a dummy config for testing
	writerConfig := cwConfig.ChainWriterConfig{
		Modules: map[string]*cwConfig.ChainWriterModule{},
	}

	// Create a mock client
	mockClient := &testutils.FakeSuiPTBClient{
		Status: client.TransactionResult{
			Status: "success",
			Error:  "",
		},
	}

	// Create a test logger
	log := logger.Test(t)

	// Create PTBConstructor using NewPTBConstructor
	ptbService := ptb.NewPTBConstructor(writerConfig, mockClient, log)

	params := []codec.SuiFunctionParam{
		{Name: "coin", Type: "Coin<T>", IsGeneric: true},
	}

	// Test with nil ArgTypes
	arguments := cwConfig.Arguments{
		ArgTypes: nil,
	}

	result, err := ptbService.ResolveGenericTypeTags(params, arguments)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "generic parameter \"coin\" not found in ArgTypes")
	assert.Nil(t, result)
}
