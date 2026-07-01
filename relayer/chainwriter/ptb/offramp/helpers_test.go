package offramp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
)

func TestExposedFunctionSignature_InvalidShape(t *testing.T) {
	tests := []struct {
		name string
		raw  any
	}{
		{name: "string", raw: "not-a-map"},
		{name: "nil", raw: nil},
		{name: "float64", raw: float64(1)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				m, err := exposedFunctionSignature("ccip_receive", tc.raw)
				require.Error(t, err)
				assert.Nil(t, m)
				assert.ErrorIs(t, err, ErrUnsupportedReceiverABI)
			})
		})
	}
}

func TestExposedFunctionSignature_ValidMap(t *testing.T) {
	m, err := exposedFunctionSignature("ccip_receive", map[string]any{"parameters": []any{}})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"parameters": []any{}}, m)
}

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
			require.ErrorIs(t, err, ErrUnsupportedReceiverABI,
				"expected ErrUnsupportedReceiverABI, got: %v", err)
			assert.Contains(t, err.Error(), "TypeParameter")
			assert.Equal(t, SuiArgumentMetadata{}, result)
		})
	}
}

func TestDecodeParam_UnsupportedABI_DefaultBranch(t *testing.T) {
	lggr := logger.Test(t)

	tests := []struct {
		name  string
		param any
	}{
		{
			name:  "unknown key with non-map value",
			param: map[string]any{"SomeUnknownKey": float64(99)},
		},
		{
			name:  "unknown key with map missing Struct",
			param: map[string]any{"SomeKey": map[string]any{"NotStruct": "x"}},
		},
		{
			name:  "unknown key with non-map Struct value",
			param: map[string]any{"SomeKey": map[string]any{"Struct": "not-a-map"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := decodeParam(lggr, tc.param, "Reference")
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrUnsupportedReceiverABI,
				"expected ErrUnsupportedReceiverABI, got: %v", err)
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

func TestDecodeParam_MalformedInput_NotUnsupportedABI(t *testing.T) {
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := decodeParam(lggr, tc.param, "Reference")
				require.Error(t, err, "expected error for input: %v", tc.param)
				assert.NotErrorIs(t, err, ErrUnsupportedReceiverABI,
					"should NOT be ErrUnsupportedReceiverABI for malformed input: %v", err)
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
		{
			name: "0x1::string::String struct",
			param: map[string]any{"Struct": map[string]any{
				"address":       "0x1",
				"module":        "string",
				"name":          "String",
				"typeArguments": []any{},
			}},
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Address:       "0x1",
				Module:        "string",
				Name:          "String",
				Reference:     "Reference",
				TypeArguments: []TypeParameter{},
				Type:          moveStringType,
			},
		},
		{
			name: "Vector of 0x1::string::String",
			param: map[string]any{"Vector": map[string]any{"Struct": map[string]any{
				"address":       "0x1",
				"module":        "string",
				"name":          "String",
				"typeArguments": []any{},
			}}},
			reference: "Reference",
			expected: SuiArgumentMetadata{
				Address:       "0x1",
				Module:        "string",
				Name:          "String",
				Reference:     "Vector",
				TypeArguments: []TypeParameter{},
				Type:          moveStringType,
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
		assert.ErrorIs(t, err, ErrUnsupportedReceiverABI,
			"expected ErrUnsupportedReceiverABI to propagate through DecodeParameters, got: %v", err)
	})
}

func TestDecodeParameters_StringParam(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{
		"parameters": []any{
			map[string]any{"Struct": map[string]any{
				"address": "0x1", "module": "string", "name": "String", "typeArguments": []any{},
			}},
		},
	}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.NoError(t, err)
	assert.Equal(t, []string{moveStringType}, result)
}

func TestDecodeParameters_AsciiStringParam(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{
		"parameters": []any{
			map[string]any{"Struct": map[string]any{
				"address": "0x1", "module": "ascii", "name": "String", "typeArguments": []any{},
			}},
		},
	}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.NoError(t, err)
	assert.Equal(t, []string{moveASCIIStringType}, result)
}

func TestDecodeParameters_VectorAsciiString(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{
		"parameters": []any{
			map[string]any{"Vector": map[string]any{"Struct": map[string]any{
				"address": "0x1", "module": "ascii", "name": "String", "typeArguments": []any{},
			}}},
		},
	}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.NoError(t, err)
	assert.Equal(t, []string{"vector<" + moveASCIIStringType + ">"}, result)
}

func TestDecodeParameters_VectorString(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{
		"parameters": []any{
			map[string]any{"Vector": map[string]any{"Struct": map[string]any{
				"address": "0x1", "module": "string", "name": "String", "typeArguments": []any{},
			}}},
		},
	}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.NoError(t, err)
	assert.Equal(t, []string{"vector<" + moveStringType + ">"}, result)
}

func TestParseParamType_String(t *testing.T) {
	lggr := logger.Test(t)

	inner := map[string]any{
		"address": "0x1", "module": "string", "name": "String", "typeArguments": []any{},
	}
	wrapper := map[string]any{"Struct": inner}

	assert.Equal(t, moveStringType, ParseParamType(lggr, inner))
	assert.Equal(t, moveStringType, ParseParamType(lggr, wrapper))
	assert.Equal(t, "vector<"+moveStringType+">", ParseParamType(lggr, map[string]any{"Vector": wrapper}))
}

func suiStringStructABI() map[string]any {
	return map[string]any{
		"address": "0x1", "module": "string", "name": "String", "typeArguments": []any{},
	}
}

func suiObjectStructABI(address, module, name string) map[string]any {
	return map[string]any{
		"address": address, "module": module, "name": name, "typeArguments": []any{},
	}
}

// TestDecodeParameters_ComplexSignature exercises a mixed parameter list resembling a rich
// Move entry function: object ref, address, vectors, std::string, and primitives.
func TestDecodeParameters_ComplexSignature(t *testing.T) {
	lggr := logger.Test(t)

	function := map[string]any{
		"parameters": []any{
			// &CCIPObjectRef (object reference)
			map[string]any{"Reference": map[string]any{"Struct": suiObjectStructABI(
				"0xccip", "state_object", "CCIPObjectRef",
			)}},
			// address primitive (Move Address)
			"Address",
			// vector<u8>
			map[string]any{"Vector": "U8"},
			// vector<SomeStruct> (vector of object-like structs)
			map[string]any{"Vector": map[string]any{"Struct": suiObjectStructABI(
				"0xpool", "burn_mint", "BurnMintTokenPoolState",
			)}},
			// 0x1::string::String by value
			map[string]any{"Struct": suiStringStructABI()},
			// vector<0x1::string::String>
			map[string]any{"Vector": map[string]any{"Struct": suiStringStructABI()}},
			// bool
			"Bool",
		},
	}

	result, err := DecodeParameters(lggr, function, "parameters")
	require.NoError(t, err)

	// Address uses default Reference label and object_id typing → &object today (not bindings' "address").
	expected := []string{
		"&object",
		"&object",
		"vector<u8>",
		"vector<object_id>",
		moveStringType,
		"vector<" + moveStringType + ">",
		"bool",
	}
	assert.Equal(t, expected, result)

	decoded := make([]SuiArgumentMetadata, len(function["parameters"].([]any)))
	for i, parameter := range function["parameters"].([]any) {
		meta, decodeErr := decodeParam(lggr, parameter, "Reference")
		require.NoError(t, decodeErr)
		decoded[i] = meta
	}

	assert.Equal(t, "object_id", decoded[0].Type)
	assert.Equal(t, "CCIPObjectRef", decoded[0].Name)
	assert.Equal(t, "object_id", decoded[1].Type)
	assert.Equal(t, "Address", decoded[1].Name)
	assert.Equal(t, "u8", decoded[2].Type)
	assert.Equal(t, "object_id", decoded[3].Type)
	assert.Equal(t, "BurnMintTokenPoolState", decoded[3].Name)
	assert.Equal(t, moveStringType, decoded[4].Type)
	assert.Equal(t, moveStringType, decoded[5].Type)
	assert.Equal(t, "bool", decoded[6].Type)
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

// The tests below exercise DecodeParametersFromFunctionDescriptor, the gRPC-based counterpart to
// DecodeParameters, using suirpcv2.FunctionDescriptor fixtures instead of JSON-RPC-shaped maps.
// Each case mirrors an equivalent DecodeParameters test above to confirm behavioral parity.

func openSig(reference *suirpcv2.OpenSignature_Reference, body *suirpcv2.OpenSignatureBody) *suirpcv2.OpenSignature {
	return &suirpcv2.OpenSignature{Reference: reference, Body: body}
}

func primitiveBody(t suirpcv2.OpenSignatureBody_Type) *suirpcv2.OpenSignatureBody {
	return &suirpcv2.OpenSignatureBody{Type: AnyPointer(t)}
}

func datatypeBody(typeName string) *suirpcv2.OpenSignatureBody {
	return &suirpcv2.OpenSignatureBody{Type: AnyPointer(suirpcv2.OpenSignatureBody_DATATYPE), TypeName: AnyPointer(typeName)}
}

func vectorBody(element *suirpcv2.OpenSignatureBody) *suirpcv2.OpenSignatureBody {
	return &suirpcv2.OpenSignatureBody{
		Type:                       AnyPointer(suirpcv2.OpenSignatureBody_VECTOR),
		TypeParameterInstantiation: []*suirpcv2.OpenSignatureBody{element},
	}
}

func typeParameterBody(position uint32) *suirpcv2.OpenSignatureBody {
	return &suirpcv2.OpenSignatureBody{Type: AnyPointer(suirpcv2.OpenSignatureBody_TYPE_PARAMETER), TypeParameter: AnyPointer(position)}
}

func TestDecodeParametersFromFunctionDescriptor_NilDescriptor(t *testing.T) {
	lggr := logger.Test(t)

	result, err := DecodeParametersFromFunctionDescriptor(lggr, nil)
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDecodeParametersFromFunctionDescriptor_EmptyParameters(t *testing.T) {
	lggr := logger.Test(t)

	result, err := DecodeParametersFromFunctionDescriptor(lggr, &suirpcv2.FunctionDescriptor{})
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestDecodeParametersFromFunctionDescriptor_StringParam(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			openSig(nil, datatypeBody("0x1::string::String")),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.NoError(t, err)
	assert.Equal(t, []string{moveStringType}, result)
}

func TestDecodeParametersFromFunctionDescriptor_VectorAsciiString(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			openSig(nil, vectorBody(datatypeBody("0x1::ascii::String"))),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.NoError(t, err)
	assert.Equal(t, []string{"vector<" + moveASCIIStringType + ">"}, result)
}

func TestDecodeParametersFromFunctionDescriptor_NestedVectorU8(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			openSig(nil, vectorBody(vectorBody(primitiveBody(suirpcv2.OpenSignatureBody_U8)))),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.NoError(t, err)
	assert.Equal(t, []string{"vector<vector<u8>>"}, result)
}

// TestDecodeParametersFromFunctionDescriptor_ComplexSignature mirrors
// TestDecodeParameters_ComplexSignature using the gRPC FunctionDescriptor shape.
func TestDecodeParametersFromFunctionDescriptor_ComplexSignature(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			// &CCIPObjectRef (object reference)
			openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), datatypeBody("0xccip::state_object::CCIPObjectRef")),
			// address primitive (Move Address)
			openSig(nil, primitiveBody(suirpcv2.OpenSignatureBody_ADDRESS)),
			// vector<u8>
			openSig(nil, vectorBody(primitiveBody(suirpcv2.OpenSignatureBody_U8))),
			// vector<SomeStruct> (vector of object-like structs)
			openSig(nil, vectorBody(datatypeBody("0xpool::burn_mint::BurnMintTokenPoolState"))),
			// 0x1::string::String by value
			openSig(nil, datatypeBody("0x1::string::String")),
			// vector<0x1::string::String>
			openSig(nil, vectorBody(datatypeBody("0x1::string::String"))),
			// bool
			openSig(nil, primitiveBody(suirpcv2.OpenSignatureBody_BOOL)),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.NoError(t, err)

	expected := []string{
		"&object",
		"&object",
		"vector<u8>",
		"vector<object_id>",
		moveStringType,
		"vector<" + moveStringType + ">",
		"bool",
	}
	assert.Equal(t, expected, result)
}

// TestDecodeParametersFromFunctionDescriptor_ValidDummyReceiverSignature mirrors
// TestDecodeParameters_ValidDummyReceiverSignature using the gRPC FunctionDescriptor shape.
func TestDecodeParametersFromFunctionDescriptor_ValidDummyReceiverSignature(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			openSig(nil, vectorBody(primitiveBody(suirpcv2.OpenSignatureBody_U8))),
			openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), datatypeBody("0xccip::state_object::CCIPObjectRef")),
			openSig(nil, datatypeBody("0xccip::client::Any2SuiMessage")),
			openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), datatypeBody("0x2::clock::Clock")),
			openSig(AnyPointer(suirpcv2.OpenSignature_MUTABLE), datatypeBody("0xreceiver::dummy_receiver::CCIPReceiverState")),
			openSig(AnyPointer(suirpcv2.OpenSignature_MUTABLE), datatypeBody("0x2::tx_context::TxContext")),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.NoError(t, err)

	// TxContext is skipped.
	expected := []string{
		"vector<u8>",
		"&object",
		"&object",
		"&object",
		"&mut object",
	}
	assert.Equal(t, expected, result)
}

func TestDecodeParametersFromFunctionDescriptor_TopLevelTypeParameter_Unsupported(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			openSig(nil, typeParameterBody(0)),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrUnsupportedReceiverABI)
}

func TestDecodeParametersFromFunctionDescriptor_VectorWrappingTypeParameter_Unsupported(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			openSig(AnyPointer(suirpcv2.OpenSignature_IMMUTABLE), vectorBody(typeParameterBody(0))),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrUnsupportedReceiverABI)
}

func TestDecodeParametersFromFunctionDescriptor_NilBody(t *testing.T) {
	lggr := logger.Test(t)

	fd := &suirpcv2.FunctionDescriptor{
		Parameters: []*suirpcv2.OpenSignature{
			openSig(nil, nil),
		},
	}

	result, err := DecodeParametersFromFunctionDescriptor(lggr, fd)
	require.Error(t, err)
	assert.Nil(t, result)
}
