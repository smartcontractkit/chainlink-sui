//go:build unit

package ptb_test

import (
	"testing"

	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func strPtr(s string) *string {
	return &s
}

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
		expectError bool
		errorMsg    string
		expectedLen int
		validate    func(t *testing.T, result []transaction.TypeTag)
	}{
		{
			name:        "no parameters",
			params:      []codec.SuiFunctionParam{},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "no generic parameters",
			params: []codec.SuiFunctionParam{
				{Name: "value", Type: "u64", GenericType: nil},
			},
			expectError: false,
			expectedLen: 0,
		},
		{
			name: "single generic parameter",
			params: []codec.SuiFunctionParam{
				{Name: "coin", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")},
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
				{Name: "coin1", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")},
				{Name: "coin2", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")},
			},
			expectError: false,
			expectedLen: 1, // Should deduplicate
		},
		{
			name: "multiple generic parameters with different types",
			params: []codec.SuiFunctionParam{
				{Name: "coin1", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")},
				{Name: "coin2", Type: "Coin<U>", GenericType: strPtr("0x2::coin::Coin")},
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name: "invalid type format",
			params: []codec.SuiFunctionParam{
				{Name: "coin", Type: "Coin<T>", GenericType: strPtr("invalid_type")},
			},
			expectError: true,
			errorMsg:    "failed to create type tag",
		},
		{
			name: "vector type parameter",
			params: []codec.SuiFunctionParam{
				{Name: "coins", Type: "vector<Coin<T>>", GenericType: strPtr("vector<0x2::sui::SUI>")},
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
				{Name: "value", Type: "u64", GenericType: nil},
				{Name: "coin", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")},
				{Name: "amount", Type: "u64", GenericType: nil},
			},
			expectError: false,
			expectedLen: 1,
		},
		{
			name: "empty vector type parameter",
			params: []codec.SuiFunctionParam{
				{Name: "items", Type: "vector<T>", GenericType: strPtr("vector<>")},
			},
			expectError: true,
			errorMsg:    "failed to extract vector inner type",
		},
		{
			name: "nested generic types",
			params: []codec.SuiFunctionParam{
				{Name: "nested_coins", Type: "vector<Coin<T>>", GenericType: strPtr("vector<0x2::coin::Coin>")},
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
				{Name: "sui_coin", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")},
				{Name: "link_coin", Type: "Coin<U>", GenericType: strPtr("0x2::link::LINK")},
				{Name: "coin_vector", Type: "vector<Coin<V>>", GenericType: strPtr("vector<0x2::coin::Coin>")},
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
				{Name: "token", Type: "Token<T>", GenericType: strPtr("0xZZZ::token::Token")},
			},
			expectError: true,
			errorMsg:    "failed to convert package address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := ptbService.ResolveGenericTypeTags(tt.params)

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
		{Name: "coin1", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")},
		{Name: "coin2", Type: "Coin<U>", GenericType: strPtr("0x2::coin::Coin")},
		{Name: "coin3", Type: "Coin<T>", GenericType: strPtr("0x2::sui::SUI")}, // Same as coin1
		{Name: "coin4", Type: "Coin<V>", GenericType: strPtr("0x2::link::LINK")},
		{Name: "coin5", Type: "Coin<U>", GenericType: strPtr("0x2::coin::Coin")}, // Same as coin2
	}

	result, err := ptbService.ResolveGenericTypeTags(params)
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
