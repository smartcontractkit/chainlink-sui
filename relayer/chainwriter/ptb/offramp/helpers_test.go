package offramp

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeParam_PrimitiveString(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	tests := []struct {
		input    string
		wantName string
		wantType string
	}{
		{"U8", "U8", "u8"},
		{"U64", "U64", "u64"},
		{"Bool", "Bool", "bool"},
		{"Address", "Address", "object_id"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()
			meta, err := decodeParam(lggr, tc.input, "Reference")
			require.NoError(t, err)
			assert.Equal(t, tc.wantName, meta.Name)
			assert.Equal(t, tc.wantType, meta.Type)
			assert.Equal(t, "Reference", meta.Reference)
		})
	}
}

func TestDecodeParam_StructDirect(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"Struct": map[string]any{
			"address":       "0xcccc",
			"module":        "state_object",
			"name":          "CCIPObjectRef",
			"typeArguments": []any{},
		},
	}
	meta, err := decodeParam(lggr, param, "Reference")
	require.NoError(t, err)
	assert.Equal(t, "0xcccc", meta.Address)
	assert.Equal(t, "state_object", meta.Module)
	assert.Equal(t, "CCIPObjectRef", meta.Name)
	assert.Equal(t, "Reference", meta.Reference)
	assert.Empty(t, meta.TypeArguments)
}

func TestDecodeParam_StructWithTypeArguments(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"Struct": map[string]any{
			"address": "0xaabb",
			"module":  "publisher_wrapper",
			"name":    "PublisherWrapper",
			"typeArguments": []any{
				map[string]any{"TypeParameter": float64(0)},
			},
		},
	}
	meta, err := decodeParam(lggr, param, "Reference")
	require.NoError(t, err)
	assert.Equal(t, "PublisherWrapper", meta.Name)
	assert.Len(t, meta.TypeArguments, 1)
	assert.Equal(t, float64(0), meta.TypeArguments[0].TypeParameter)
}

func TestDecodeParam_Reference(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"Reference": map[string]any{
			"Struct": map[string]any{
				"address":       "0x2",
				"module":        "clock",
				"name":          "Clock",
				"typeArguments": []any{},
			},
		},
	}
	meta, err := decodeParam(lggr, param, "Reference")
	require.NoError(t, err)
	assert.Equal(t, "Clock", meta.Name)
	assert.Equal(t, "Reference", meta.Reference)
}

func TestDecodeParam_MutableReference(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"MutableReference": map[string]any{
			"Struct": map[string]any{
				"address":       "0xdead",
				"module":        "my_receiver",
				"name":          "ReceiverState",
				"typeArguments": []any{},
			},
		},
	}
	meta, err := decodeParam(lggr, param, "Reference")
	require.NoError(t, err)
	assert.Equal(t, "ReceiverState", meta.Name)
	assert.Equal(t, "MutableReference", meta.Reference)
}

func TestDecodeParam_VectorPrimitive(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{"Vector": "U8"}
	meta, err := decodeParam(lggr, param, "Reference")
	require.NoError(t, err)
	assert.Equal(t, "U8", meta.Name)
	assert.Equal(t, "Vector", meta.Reference)
	assert.Equal(t, "u8", meta.Type)
}

func TestDecodeParam_TypeParameter_ReturnsError(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{"TypeParameter": float64(0)}
	_, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported TypeParameter")
}

func TestDecodeParam_VectorTypeParameter_ReturnsError(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"Vector": map[string]any{"TypeParameter": float64(0)},
	}
	_, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported TypeParameter")
}

func TestDecodeParam_NonMapNonString_ReturnsError(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	_, err := decodeParam(lggr, float64(42), "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected string or map")
}

func TestDecodeParam_StructMissingAddress_ReturnsError(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"Struct": map[string]any{
			"module":        "foo",
			"name":          "Bar",
			"typeArguments": []any{},
		},
	}
	_, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing 'address'")
}

func TestDecodeParam_StructNonMapValue_ReturnsError(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"Struct": "not-a-map",
	}
	_, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Struct value is")
}

func TestDecodeParam_DefaultBranch_MissingNestedStruct(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"SomeOtherKey": map[string]any{"NotStruct": "value"},
	}
	_, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing nested Struct")
}

func TestDecodeParam_DefaultBranch_NonMapValue(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"SomeOtherKey": "not-a-map",
	}
	_, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected map value")
}

func TestDecodeParam_TypeArgumentsBadShape(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	param := map[string]any{
		"Struct": map[string]any{
			"address":       "0x1",
			"module":        "foo",
			"name":          "Bar",
			"typeArguments": []any{"not-a-map"},
		},
	}
	_, err := decodeParam(lggr, param, "Reference")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "typeArguments[0]")
}

func TestDecodeParameters_RejectsTypeParameter(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	funcSig := map[string]any{
		"parameters": []any{
			map[string]any{"Vector": "U8"},
			map[string]any{"TypeParameter": float64(0)},
		},
	}
	_, err := DecodeParameters(lggr, funcSig, "parameters")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode parameter 1")
	assert.Contains(t, err.Error(), "unsupported TypeParameter")
}

func TestDecodeParameters_ValidStandardReceiver(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	funcSig := map[string]any{
		"parameters": []any{
			map[string]any{"Vector": "U8"},
			map[string]any{
				"Reference": map[string]any{
					"Struct": map[string]any{
						"address":       "0xcccc",
						"module":        "state_object",
						"name":          "CCIPObjectRef",
						"typeArguments": []any{},
					},
				},
			},
			map[string]any{
				"Struct": map[string]any{
					"address":       "0xcccc",
					"module":        "client",
					"name":          "Any2SuiMessage",
					"typeArguments": []any{},
				},
			},
			map[string]any{
				"MutableReference": map[string]any{
					"Struct": map[string]any{
						"address":       "0x2",
						"module":        "tx_context",
						"name":          "TxContext",
						"typeArguments": []any{},
					},
				},
			},
		},
	}
	paramTypes, err := DecodeParameters(lggr, funcSig, "parameters")
	require.NoError(t, err)
	assert.Equal(t, []string{"vector<u8>", "&object", "&object"}, paramTypes)
}

func TestDecodeParameters_MissingKey(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	_, err := DecodeParameters(lggr, map[string]any{}, "parameters")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing or nil")
}

func TestDecodeParameters_NotArray(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	_, err := DecodeParameters(lggr, map[string]any{"parameters": "oops"}, "parameters")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not an array")
}

func TestExtractTypeArguments_Empty(t *testing.T) {
	t.Parallel()

	s := map[string]any{"typeArguments": []any{}}
	ta, err := extractTypeArguments(s)
	require.NoError(t, err)
	assert.Empty(t, ta)
}

func TestExtractTypeArguments_Missing(t *testing.T) {
	t.Parallel()

	s := map[string]any{}
	ta, err := extractTypeArguments(s)
	require.NoError(t, err)
	assert.Empty(t, ta)
}

func TestExtractTypeArguments_WrongType(t *testing.T) {
	t.Parallel()

	s := map[string]any{"typeArguments": "not-an-array"}
	_, err := extractTypeArguments(s)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected array")
}

func TestExtractStructIdentity_Valid(t *testing.T) {
	t.Parallel()

	s := map[string]any{
		"address": "0x1",
		"module":  "foo",
		"name":    "Bar",
	}
	addr, mod, name, err := extractStructIdentity(s)
	require.NoError(t, err)
	assert.Equal(t, "0x1", addr)
	assert.Equal(t, "foo", mod)
	assert.Equal(t, "Bar", name)
}

func TestExtractStructIdentity_MissingFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input map[string]any
		want  string
	}{
		{"missing address", map[string]any{"module": "a", "name": "b"}, "missing 'address'"},
		{"missing module", map[string]any{"address": "a", "name": "b"}, "missing 'module'"},
		{"missing name", map[string]any{"address": "a", "module": "b"}, "missing 'name'"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, _, _, err := extractStructIdentity(tc.input)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}
