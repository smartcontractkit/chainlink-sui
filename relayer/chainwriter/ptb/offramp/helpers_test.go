package offramp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestDecodeParam_PoisonABI_TypeParameter(t *testing.T) {
	lggr := logger.Test(t)

	tests := []struct {
		name  string
		param any
	}{
		{
			name:  "Vector wrapping TypeParameter",
			param: map[string]any{"Vector": map[string]any{"TypeParameter": float64(0)}},
		},
		{
			name:  "top-level TypeParameter",
			param: map[string]any{"TypeParameter": float64(0)},
		},
		{
			name:  "Reference wrapping Vector wrapping TypeParameter",
			param: map[string]any{"Reference": map[string]any{"Vector": map[string]any{"TypeParameter": float64(0)}}},
		},
		{
			name:  "MutableReference wrapping TypeParameter",
			param: map[string]any{"MutableReference": map[string]any{"TypeParameter": float64(1)}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := decodeParam(lggr, tc.param, "Reference")
			require.Error(t, err)
			assert.Contains(t, err.Error(), "TypeParameter")
			assert.Equal(t, SuiArgumentMetadata{}, result)
		})
	}
}

func TestDecodeParam_MultiKeyWrapperMap_RejectsDeterministically(t *testing.T) {
	lggr := logger.Test(t)

	poison := map[string]any{
		"Vector":        map[string]any{"TypeParameter": float64(0)},
		"TypeParameter": float64(0),
	}

	for i := 0; i < 50; i++ {
		result, err := decodeParam(lggr, poison, "Reference")
		require.Error(t, err, "iteration %d", i)
		assert.Contains(t, err.Error(), "exactly one ABI wrapper key")
		assert.Equal(t, SuiArgumentMetadata{}, result)
	}
}

func TestDecodeParam_MultiKeyWrapperMap_NestedRejects(t *testing.T) {
	lggr := logger.Test(t)

	param := map[string]any{
		"Vector": map[string]any{
			"Struct":        map[string]any{"address": "0x1", "module": "m", "name": "S", "typeArguments": []any{}},
			"TypeParameter": float64(0),
		},
	}

	result, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "exactly one ABI wrapper key")
	assert.Equal(t, SuiArgumentMetadata{}, result)
}

func TestDecodeParam_MalformedInput_NoPanic(t *testing.T) {
	lggr := logger.Test(t)

	tests := []struct {
		name  string
		param any
	}{
		{
			name:  "nil parameter",
			param: nil,
		},
		{
			name:  "float64 instead of map or string",
			param: float64(42),
		},
		{
			name:  "integer instead of map or string",
			param: int(7),
		},
		{
			name:  "Struct with nil value",
			param: map[string]any{"Struct": nil},
		},
		{
			name:  "Struct with wrong type value",
			param: map[string]any{"Struct": "not-a-map"},
		},
		{
			name:  "Struct missing address field",
			param: map[string]any{"Struct": map[string]any{"module": "m", "name": "S", "typeArguments": []any{}}},
		},
		{
			name:  "Struct with non-string address",
			param: map[string]any{"Struct": map[string]any{"address": 123, "module": "m", "name": "S", "typeArguments": []any{}}},
		},
		{
			name:  "Struct with bad typeArguments type",
			param: map[string]any{"Struct": map[string]any{"address": "0x1", "module": "m", "name": "S", "typeArguments": "not-array"}},
		},
		{
			name:  "Struct with typeArgument missing TypeParameter key",
			param: map[string]any{"Struct": map[string]any{"address": "0x1", "module": "m", "name": "S", "typeArguments": []any{map[string]any{"NotTypeParameter": float64(0)}}}},
		},
		{
			name:  "empty map",
			param: map[string]any{},
		},
		{
			name: "multiple top-level wrapper keys",
			param: map[string]any{
				"Vector":        map[string]any{"U8": nil},
				"TypeParameter": float64(0),
			},
		},
		{
			name:  "default key with non-map value",
			param: map[string]any{"SomeUnknownKey": float64(99)},
		},
		{
			name:  "default key with map missing Struct",
			param: map[string]any{"SomeKey": map[string]any{"NotStruct": "x"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := decodeParam(lggr, tc.param, "Reference")
				require.Error(t, err, "expected error for input: %v", tc.param)
			})
		})
	}
}

func TestDecodeParam_ValidABI(t *testing.T) {
	lggr := logger.Test(t)

	tests := []struct {
		name      string
		param     any
		reference string
		expected  SuiArgumentMetadata
	}{
		{
			name:      "primitive U64",
			param:     "U64",
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Name:          "U64",
				Reference:     "Reference",
				TypeArguments: []TypeParameter{},
				Type:          "u64",
			},
		},
		{
			name:      "primitive Bool",
			param:     "Bool",
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Name:          "Bool",
				Reference:     "Reference",
				TypeArguments: []TypeParameter{},
				Type:          "bool",
			},
		},
		{
			name:      "Vector of U8",
			param:     map[string]any{"Vector": "U8"},
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Name:          "U8",
				Reference:     "Vector",
				TypeArguments: []TypeParameter{},
				Type:          "u8",
			},
		},
		{
			name:      "nested Vector of U8",
			param:     map[string]any{"Vector": map[string]any{"Vector": "U8"}},
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Name:          "U8",
				Reference:     "Vector",
				TypeArguments: []TypeParameter{},
				Type:          "vector<u8>",
			},
		},
		{
			name: "Struct with no type arguments",
			param: map[string]any{"Struct": map[string]any{
				"address":       "0x1",
				"module":        "state_object",
				"name":          "CCIPObjectRef",
				"typeArguments": []any{},
			}},
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Address:       "0x1",
				Module:        "state_object",
				Name:          "CCIPObjectRef",
				Reference:     "Reference",
				TypeArguments: []TypeParameter{},
				Type:          "object_id",
			},
		},
		{
			name: "MutableReference wrapping Struct",
			param: map[string]any{"MutableReference": map[string]any{"Struct": map[string]any{
				"address":       "0xabc",
				"module":        "dummy_receiver",
				"name":          "CCIPReceiverState",
				"typeArguments": []any{},
			}}},
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Address:       "0xabc",
				Module:        "dummy_receiver",
				Name:          "CCIPReceiverState",
				Reference:     "MutableReference",
				TypeArguments: []TypeParameter{},
				Type:          "object_id",
			},
		},
		{
			name: "Struct with one type argument",
			param: map[string]any{"Struct": map[string]any{
				"address":       "0xpool",
				"module":        "burn_mint",
				"name":          "BurnMintTokenPoolState",
				"typeArguments": []any{map[string]any{"TypeParameter": float64(0)}},
			}},
			reference: "MutableReference",
			expected: SuiArgumentMetadata{
				Address:       "0xpool",
				Module:        "burn_mint",
				Name:          "BurnMintTokenPoolState",
				Reference:     "MutableReference",
				TypeArguments: []TypeParameter{{TypeParameter: 0}},
				Type:          "object_id",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := decodeParam(lggr, tc.param, tc.reference)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDecodeParameters_PoisonABI_ReturnsError(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{
		"parameters": []any{
			map[string]any{"Vector": map[string]any{"TypeParameter": float64(0)}},
			"U64",
		},
	}

	assert.NotPanics(t, func() {
		result, err := DecodeParameters(lggr, function, "parameters")
		require.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "TypeParameter")
	})
}

func TestDecodeParameters_NestedVectorU8(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{
		"parameters": []any{
			map[string]any{"Vector": map[string]any{"Vector": "U8"}},
		},
	}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.NoError(t, err)
	assert.Equal(t, []string{"vector<vector<u8>>"}, result)
}

func TestDecodeParameters_ValidDummyReceiverSignature(t *testing.T) {
	lggr := logger.Test(t)

	// Matches the normalized ABI of dummy_receiver::ccip_receive:
	// ccip_receive(expected_message_id: vector<u8>, ref: &CCIPObjectRef, message: Any2SuiMessage, _: &Clock, state: &mut CCIPReceiverState)
	function := map[string]any{
		"parameters": []any{
			map[string]any{"Vector": "U8"},
			map[string]any{"Reference": map[string]any{"Struct": map[string]any{
				"address": "0xccip", "module": "state_object", "name": "CCIPObjectRef", "typeArguments": []any{},
			}}},
			map[string]any{"Struct": map[string]any{
				"address": "0xccip", "module": "client", "name": "Any2SuiMessage", "typeArguments": []any{},
			}},
			map[string]any{"Reference": map[string]any{"Struct": map[string]any{
				"address": "0x2", "module": "clock", "name": "Clock", "typeArguments": []any{},
			}}},
			map[string]any{"MutableReference": map[string]any{"Struct": map[string]any{
				"address": "0xreceiver", "module": "dummy_receiver", "name": "CCIPReceiverState", "typeArguments": []any{},
			}}},
			map[string]any{"MutableReference": map[string]any{"Struct": map[string]any{
				"address": "0x2", "module": "tx_context", "name": "TxContext", "typeArguments": []any{},
			}}},
		},
	}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.NoError(t, err)

	// TxContext is skipped. Any2SuiMessage is a Struct with default reference "Reference" → "&object"
	expected := []string{
		"vector<u8>",
		"&object",
		"&object",
		"&object",
		"&mut object",
	}
	assert.Equal(t, expected, result)
}

func TestDecodeParameters_MissingKey(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{"return": []any{}}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDecodeParameters_NilValue(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{"parameters": nil}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.Error(t, err)
	assert.Nil(t, result)
}
