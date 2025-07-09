package codec

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"testing"

	aptosBCS "github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/utils"

	"github.com/smartcontractkit/chainlink-sui/shared"

	"github.com/stretchr/testify/require"
)

// Test structs for single-field struct cases
type SingleFieldUint64 struct {
	Value uint64 `json:"value"`
}

type SingleFieldBytes struct {
	Data []byte `json:"data"`
}

type SingleFieldString struct {
	Text string `json:"text"`
}

type SingleFieldBigInt struct {
	Amount *big.Int `json:"amount"`
}

// Complex test structs
type ComplexStruct struct {
	ID       uint64   `json:"id"`
	Name     string   `json:"name"`
	Data     []byte   `json:"data"`
	Values   []uint32 `json:"values"`
	Enabled  bool     `json:"enabled"`
	Optional *string  `json:"optional,omitempty"`
}

type NestedStruct struct {
	Outer ComplexStruct `json:"outer"`
	Inner struct {
		Count uint32 `json:"count"`
	} `json:"inner"`
}

func TestDecodeSuiJsonValue_NilTarget(t *testing.T) {
	t.Parallel()

	err := DecodeSuiJsonValue("test", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "target cannot be nil")
}

func TestDecodeSuiJsonValue_DirectAssignment(t *testing.T) {
	t.Parallel()

	var target string
	data := "test_value"

	err := DecodeSuiJsonValue(data, &target)
	require.NoError(t, err)
	require.Equal(t, "test_value", target)
}

func TestDecodeSuiJsonValue_NumericTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     any
		target   any
		expected any
	}{
		// Float64 to various uint types
		{
			name:     "float64 to uint64",
			data:     float64(12345),
			target:   new(uint64),
			expected: uint64(12345),
		},
		{
			name:     "float64 to uint32",
			data:     float64(12345),
			target:   new(uint32),
			expected: uint32(12345),
		},
		{
			name:     "float64 to uint16",
			data:     float64(12345),
			target:   new(uint16),
			expected: uint16(12345),
		},
		{
			name:     "float64 to uint8",
			data:     float64(123),
			target:   new(uint8),
			expected: uint8(123),
		},
		// String to uint types
		{
			name:     "string to uint64",
			data:     "12345",
			target:   new(uint64),
			expected: uint64(12345),
		},
		{
			name:     "string to uint32",
			data:     "12345",
			target:   new(uint32),
			expected: uint32(12345),
		},
		// JSON Number to uint
		{
			name:     "json.Number to uint64",
			data:     json.Number("12345"),
			target:   new(uint64),
			expected: uint64(12345),
		},
		// Byte slice to uint (little-endian)
		{
			name:     "byte slice to uint64",
			data:     []byte{0x39, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // 12345 in little-endian
			target:   new(uint64),
			expected: uint64(12345),
		},
		{
			name:     "single byte to uint64",
			data:     []byte{0xFF},
			target:   new(uint64),
			expected: uint64(255),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := DecodeSuiJsonValue(tt.data, tt.target)
			require.NoError(t, err)

			targetValue := reflect.ValueOf(tt.target).Elem().Interface()
			require.Equal(t, tt.expected, targetValue)
		})
	}
}

func TestDecodeSuiJsonValue_NumericErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   any
		target any
	}{
		{
			name:   "invalid string to uint",
			data:   "not_a_number",
			target: new(uint64),
		},
		{
			name:   "negative json.Number to uint",
			data:   json.Number("-123"),
			target: new(uint64),
		},
		{
			name:   "empty byte slice",
			data:   []byte{},
			target: new(uint64),
		},
		{
			name:   "unsupported data type",
			data:   map[string]any{"key": "value"},
			target: new(uint64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := DecodeSuiJsonValue(tt.data, tt.target)
			require.Error(t, err)
		})
	}
}

func TestDecodeSuiJsonValue_StringType(t *testing.T) {
	t.Parallel()

	var target string
	err := DecodeSuiJsonValue("test_string", &target)
	require.NoError(t, err)
	require.Equal(t, "test_string", target)

	// Test error case
	err = DecodeSuiJsonValue(123, &target)
	require.Error(t, err)
	require.Contains(t, err.Error(), "expected string")
}

func TestDecodeSuiJsonValue_SliceTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     any
		target   any
		expected any
	}{
		// String to []byte conversions
		{
			name:     "numeric string to bytes",
			data:     "12345",
			target:   new([]byte),
			expected: []byte{0x39, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}[:2], // little-endian, trimmed
		},
		{
			name:     "hex string to bytes",
			data:     "0x48656c6c6f",
			target:   new([]byte),
			expected: []byte("Hello"),
		},
		{
			name:     "hex string odd length",
			data:     "0x123",
			target:   new([]byte),
			expected: []byte{0x01, 0x23},
		},
		{
			name:     "base64 string to bytes",
			data:     base64.StdEncoding.EncodeToString([]byte("Hello")),
			target:   new([]byte),
			expected: []byte("Hello"),
		},
		{
			name:     "regular string to bytes",
			data:     "Hello",
			target:   new([]byte),
			expected: []byte("Hello"),
		},
		// Array to slice
		{
			name:     "array to uint slice",
			data:     []any{float64(1), float64(2), float64(3)},
			target:   new([]uint32),
			expected: []uint32{1, 2, 3},
		},
		{
			name:     "array to string slice",
			data:     []any{"a", "b", "c"},
			target:   new([]string),
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := DecodeSuiJsonValue(tt.data, tt.target)
			require.NoError(t, err)

			targetValue := reflect.ValueOf(tt.target).Elem().Interface()
			require.Equal(t, tt.expected, targetValue)
		})
	}
}

func TestDecodeSuiJsonValue_SliceErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   any
		target any
	}{
		{
			name:   "invalid hex string",
			data:   "0xGGGG",
			target: new([]byte),
		},
		{
			name:   "non-slice data to slice target",
			data:   "not_a_slice",
			target: new([]uint32),
		},
		{
			name:   "slice element decode error",
			data:   []any{"not_a_number"},
			target: new([]uint32),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := DecodeSuiJsonValue(tt.data, tt.target)
			require.Error(t, err)
		})
	}
}

func TestDecodeSuiJsonValue_StructTypes(t *testing.T) {
	t.Parallel()

	// Test complex struct
	data := map[string]any{
		"id":      float64(123),
		"name":    "test_name",
		"data":    "0x48656c6c6f", // "Hello" in hex
		"values":  []any{float64(1), float64(2), float64(3)},
		"enabled": true,
	}

	var target ComplexStruct
	err := DecodeSuiJsonValue(data, &target)
	require.NoError(t, err)

	require.Equal(t, uint64(123), target.ID)
	require.Equal(t, "test_name", target.Name)
	require.Equal(t, []byte("Hello"), target.Data)
	require.Equal(t, []uint32{1, 2, 3}, target.Values)
	require.True(t, target.Enabled)

	// Test nested struct
	nestedData := map[string]any{
		"outer": data,
		"inner": map[string]any{
			"count": float64(42),
		},
	}

	var nestedTarget NestedStruct
	err = DecodeSuiJsonValue(nestedData, &nestedTarget)
	require.NoError(t, err)

	require.Equal(t, uint64(123), nestedTarget.Outer.ID)
	require.Equal(t, uint32(42), nestedTarget.Inner.Count)
}

func TestHexStringHook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     any
		target   any
		expected any
	}{
		// Non-string input should pass through
		{
			name:     "non-string input",
			data:     123,
			target:   new(string),
			expected: 123,
		},
		// Non-hex string should pass through
		{
			name:     "non-hex string",
			data:     "regular_string",
			target:   new(string),
			expected: "regular_string",
		},
		// Hex string to string should pass through
		{
			name:     "hex string to string",
			data:     "0x123",
			target:   new(string),
			expected: "0x123",
		},
		// Hex string to byte slice
		{
			name:     "hex to byte slice",
			data:     "0x48656c6c6f",
			target:   new([]byte),
			expected: []byte("Hello"),
		},
		{
			name:     "empty hex to byte slice",
			data:     "0x",
			target:   new([]byte),
			expected: []uint8{},
		},
		{
			name:     "odd length hex to byte slice",
			data:     "0x123",
			target:   new([]byte),
			expected: []byte{0x01, 0x23},
		},
		// Hex string to uint types
		{
			name:     "hex to uint64",
			data:     "0xFF",
			target:   new(uint64),
			expected: uint64(255),
		},
		{
			name:     "hex to uint32",
			data:     "0xFF",
			target:   new(uint32),
			expected: uint64(255),
		},
		// Hex string to int types
		{
			name:     "hex to int64",
			data:     "0xFF",
			target:   new(int64),
			expected: int64(255),
		},
		// Hex string to big.Int
		{
			name:     "hex to big.Int",
			data:     "0x123456789ABCDEF",
			target:   new(*big.Int),
			expected: func() *big.Int { bi := new(big.Int); bi.SetString("123456789ABCDEF", 16); return bi }(),
		},
		// Hex string to byte array
		{
			name:     "hex to byte array",
			data:     "0x123456",
			target:   new([4]uint8),
			expected: []uint8{0x12, 0x34, 0x56, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetType := reflect.TypeOf(tt.target).Elem()
			result, err := hexStringHook(reflect.TypeOf(tt.data), targetType, tt.data)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestHexStringHook_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   string
		target any
	}{
		{
			name:   "unsupported slice element type",
			data:   "0x123",
			target: new([]string),
		},
		{
			name:   "invalid hex for int",
			data:   "0xGGG",
			target: new(int64),
		},
		{
			name:   "unsupported array element type",
			data:   "0x123",
			target: new([4]string),
		},
		{
			name:   "unsupported target type",
			data:   "0x123",
			target: new(float64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetType := reflect.TypeOf(tt.target).Elem()
			_, err := hexStringHook(reflect.TypeOf(tt.data), targetType, tt.data)
			require.Error(t, err)
		})
	}
}

func TestBase64StringHook(t *testing.T) {
	t.Parallel()

	testData := []byte("Hello, World!")
	encoded := base64.StdEncoding.EncodeToString(testData)

	tests := []struct {
		name     string
		data     any
		target   any
		expected any
	}{
		{
			name:     "non-string input",
			data:     123,
			target:   new([]byte),
			expected: 123,
		},
		{
			name:     "base64 to byte slice",
			data:     encoded,
			target:   new([]byte),
			expected: testData,
		},
		{
			name:     "invalid base64 to byte slice",
			data:     "invalid_base64",
			target:   new([]byte),
			expected: "invalid_base64",
		},
		{
			name:     "base64 to non-byte-slice",
			data:     encoded,
			target:   new(string),
			expected: encoded,
		},
		{
			name:     "single-field struct",
			data:     encoded,
			target:   new(SingleFieldBytes),
			expected: SingleFieldBytes{Data: testData},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetType := reflect.TypeOf(tt.target).Elem()
			result, err := base64StringHook(reflect.TypeOf(tt.data), targetType, tt.data)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestNumericStringHook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     any
		target   any
		expected any
	}{
		{
			name:     "non-string input",
			data:     123,
			target:   new(int64),
			expected: 123,
		},
		{
			name:     "string to string",
			data:     "123",
			target:   new(string),
			expected: "123",
		},
		{
			name:     "string to int64",
			data:     "123",
			target:   new(int64),
			expected: int64(123),
		},
		{
			name:     "string to uint64",
			data:     "123",
			target:   new(uint64),
			expected: uint64(123),
		},
		{
			name:     "string to float64",
			data:     "123.45",
			target:   new(float64),
			expected: float64(123.45),
		},
		{
			name:     "string to big.Int",
			data:     "123456789012345678901234567890",
			target:   new(*big.Int),
			expected: func() *big.Int { bi := new(big.Int); bi.SetString("123456789012345678901234567890", 10); return bi }(),
		},
		{
			name:     "numeric string to byte slice",
			data:     "12345",
			target:   new([]byte),
			expected: []byte{0x39, 0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}[:2], // little-endian, trimmed
		},
		{
			name:     "single-field struct",
			data:     "123",
			target:   new(SingleFieldUint64),
			expected: SingleFieldUint64{Value: 123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetType := reflect.TypeOf(tt.target).Elem()
			result, err := numericStringHook(reflect.TypeOf(tt.data), targetType, tt.data)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestNumericStringHook_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   string
		target any
	}{
		{
			name:   "invalid int string",
			data:   "not_a_number",
			target: new(int64),
		},
		{
			name:   "invalid uint string",
			data:   "not_a_number",
			target: new(uint64),
		},
		{
			name:   "invalid float string",
			data:   "not_a_float",
			target: new(float64),
		},
		{
			name:   "invalid big.Int string",
			data:   "not_a_big_int",
			target: new(*big.Int),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetType := reflect.TypeOf(tt.target).Elem()
			_, err := numericStringHook(reflect.TypeOf(tt.data), targetType, tt.data)
			require.Error(t, err)
		})
	}
}

func TestBooleanHook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     any
		target   any
		expected any
	}{
		{
			name:     "non-bool input",
			data:     "not_bool",
			target:   new(bool),
			expected: "not_bool",
		},
		{
			name:     "bool to bool true",
			data:     true,
			target:   new(bool),
			expected: true,
		},
		{
			name:     "bool to bool false",
			data:     false,
			target:   new(bool),
			expected: false,
		},
		{
			name:     "bool true to int64",
			data:     true,
			target:   new(int64),
			expected: int64(1),
		},
		{
			name:     "bool false to int64",
			data:     false,
			target:   new(int64),
			expected: int64(0),
		},
		{
			name:     "bool true to uint64",
			data:     true,
			target:   new(uint64),
			expected: uint64(1),
		},
		{
			name:     "bool false to uint64",
			data:     false,
			target:   new(uint64),
			expected: uint64(0),
		},
		{
			name:     "bool true to big.Int",
			data:     true,
			target:   new(*big.Int),
			expected: big.NewInt(1),
		},
		{
			name:     "bool false to big.Int",
			data:     false,
			target:   new(*big.Int),
			expected: big.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetType := reflect.TypeOf(tt.target).Elem()
			result, err := booleanHook(reflect.TypeOf(tt.data), targetType, tt.data)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestBooleanHook_Errors(t *testing.T) {
	t.Parallel()

	_, err := booleanHook(reflect.TypeOf(true), reflect.TypeOf(float64(0)), true)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported target type for boolean conversion")
}

func TestArrayHook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     any
		target   any
		expected any
	}{
		{
			name:     "non-array input",
			data:     "not_array",
			target:   new([]int),
			expected: "not_array",
		},
		{
			name:     "array to non-slice target",
			data:     []any{1, 2, 3},
			target:   new(string),
			expected: []any{1, 2, 3},
		},
		{
			name:     "array to slice",
			data:     []any{float64(1), float64(2), float64(3)},
			target:   new([]uint32),
			expected: []uint32{1, 2, 3},
		},
		{
			name: "single-field struct",
			data: []any{float64(1), float64(2)},
			target: new(struct {
				Values []uint32 `json:"values"`
			}),
			expected: struct {
				Values []uint32 `json:"values"`
			}{Values: []uint32{1, 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targetType := reflect.TypeOf(tt.target).Elem()
			result, err := arrayHook(reflect.TypeOf(tt.data), targetType, tt.data)
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestArrayHook_Error(t *testing.T) {
	t.Parallel()

	// Test decode error in array element
	data := []any{"not_a_number"}
	targetType := reflect.TypeOf([]uint32{})

	_, err := arrayHook(reflect.TypeOf(data), targetType, data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to decode array element")
}

func TestOverflowFunctions(t *testing.T) {
	t.Parallel()

	// Test overflowFloat
	require.False(t, overflowFloat(reflect.TypeOf(float64(0)), 1.0))
	require.True(t, overflowFloat(reflect.TypeOf(float32(0)), math.MaxFloat64))

	// Test overflowFloat32
	require.False(t, overflowFloat32(1.0))
	require.True(t, overflowFloat32(math.MaxFloat64))

	// Test overflowInt
	require.False(t, overflowInt(reflect.TypeOf(int64(0)), 123))
	require.True(t, overflowInt(reflect.TypeOf(int8(0)), 1000))

	// Test overflowUint
	require.False(t, overflowUint(reflect.TypeOf(uint64(0)), 123))
	require.True(t, overflowUint(reflect.TypeOf(uint8(0)), 1000))
}

func TestOverflowFunctions_Panics(t *testing.T) {
	t.Parallel()

	// Test panic cases
	require.Panics(t, func() {
		overflowFloat(reflect.TypeOf(int(0)), 1.0)
	})

	require.Panics(t, func() {
		overflowInt(reflect.TypeOf(float64(0)), 1)
	})

	require.Panics(t, func() {
		overflowUint(reflect.TypeOf(float64(0)), 1)
	})
}

func TestDecodeBase64(t *testing.T) {
	t.Parallel()

	testData := []byte("Hello, World!")
	encoded := base64.StdEncoding.EncodeToString(testData)

	decoded, err := shared.DecodeBase64(encoded)
	require.NoError(t, err)
	require.Equal(t, testData, decoded)

	// Test invalid base64
	_, err = shared.DecodeBase64("invalid_base64!!!")
	require.Error(t, err)
}

func TestDecodeSuiJsonValue_EdgeCases(t *testing.T) {
	t.Parallel()

	// Test with very large numbers
	t.Run("large number string", func(t *testing.T) {
		t.Parallel()

		var target uint64
		err := DecodeSuiJsonValue("18446744073709551615", &target) // max uint64
		require.NoError(t, err)
		require.Equal(t, uint64(18446744073709551615), target)
	})

	// Test with zero values
	t.Run("zero values", func(t *testing.T) {
		t.Parallel()

		var target uint64
		err := DecodeSuiJsonValue(float64(0), &target)
		require.NoError(t, err)
		require.Equal(t, uint64(0), target)
	})

	// Test with empty structures
	t.Run("empty struct", func(t *testing.T) {
		t.Parallel()

		type EmptyStruct struct{}
		var target EmptyStruct
		err := DecodeSuiJsonValue(map[string]any{}, &target)
		require.NoError(t, err)
	})

	// Test with nil pointer fields
	t.Run("nil pointer fields", func(t *testing.T) {
		t.Parallel()

		type StructWithPointer struct {
			Value *string `json:"value"`
		}
		var target StructWithPointer
		err := DecodeSuiJsonValue(map[string]any{}, &target)
		require.NoError(t, err)
		require.Nil(t, target.Value)
	})
}

func TestDecodeSuiJsonValue_SuiSpecificCases(t *testing.T) {
	t.Parallel()

	// TODO: add test for Sui object ID format (hex string)

	// Test Sui balance format (string number)
	t.Run("sui balance", func(t *testing.T) {
		t.Parallel()

		var target uint64
		balance := "1000000000" // 1 SUI in MIST
		err := DecodeSuiJsonValue(balance, &target)
		require.NoError(t, err)
		require.Equal(t, uint64(1000000000), target)
	})

	// Test Sui transaction digest (base64)
	t.Run("sui transaction digest", func(t *testing.T) {
		t.Parallel()

		var target []byte
		digest := base64.StdEncoding.EncodeToString([]byte("transaction_digest"))
		err := DecodeSuiJsonValue(digest, &target)
		require.NoError(t, err)
		require.Equal(t, []byte("transaction_digest"), target)
	})

	// Test Sui event data structure
	t.Run("sui event data", func(t *testing.T) {
		t.Parallel()

		type SuiEvent struct {
			ID                string         `json:"id"`
			PackageID         string         `json:"packageId"`
			TransactionModule string         `json:"transactionModule"`
			Sender            string         `json:"sender"`
			Type              string         `json:"type"`
			ParsedJSON        map[string]any `json:"parsedJson"`
			Bcs               string         `json:"bcs"`
		}

		amount := float64(1000)
		eventData := map[string]any{
			"id":                "0xabcd1234",
			"packageId":         "0x45678901",
			"transactionModule": "test_module",
			"sender":            "0x12345678",
			"type":              "0x12345678::test::TestEvent",
			"parsedJson": map[string]any{
				"amount":    amount,
				"recipient": "0x12345678",
			},
			"bcs": base64.StdEncoding.EncodeToString([]byte("bcs_data")),
		}

		var target SuiEvent
		err := DecodeSuiJsonValue(eventData, &target)
		require.NoError(t, err)

		require.Equal(t, "0xabcd1234", target.ID)
		require.Equal(t, "0x45678901", target.PackageID)
		require.Equal(t, "test_module", target.TransactionModule)
		require.Equal(t, "0x12345678", target.Sender)
		require.Equal(t, "0x12345678::test::TestEvent", target.Type)
		require.NotNil(t, target.ParsedJSON)
		require.InDelta(t, amount, target.ParsedJSON["amount"], 0.0)
	})

	t.Run("JSON Struct Decoder", func(t *testing.T) {
		t.Parallel()

		jsonModule := `
		{
			"fileFormatVersion": 6,
			"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
			"name": "counter",
			"friends": [],
			"structs": {
				"AddressList": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop"
					]
				},
				"fields": [
					{
					"name": "addresses",
					"type": {
						"Vector": "Address"
					}
					},
					{
					"name": "count",
					"type": "U64"
					}
				],
				"typeParameters": []
				},
				"AdminCap": {
				"abilities": {
					"abilities": [
					"Store",
					"Key"
					]
				},
				"fields": [
					{
					"name": "id",
					"type": {
						"Struct": {
						"address": "0x2",
						"module": "object",
						"name": "UID",
						"typeArguments": []
						}
					}
					}
				],
				"typeParameters": []
				},
				"COUNTER": {
				"abilities": {
					"abilities": [
					"Drop"
					]
				},
				"fields": [
					{
					"name": "dummy_field",
					"type": "Bool"
					}
				],
				"typeParameters": []
				},
				"ComplexResult": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop"
					]
				},
				"fields": [
					{
					"name": "count",
					"type": "U64"
					},
					{
					"name": "addr",
					"type": "Address"
					},
					{
					"name": "is_complex",
					"type": "Bool"
					},
					{
					"name": "bytes",
					"type": {
						"Vector": "U8"
					}
					}
				],
				"typeParameters": []
				},
				"ConfigInfo": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop",
					"Store"
					]
				},
				"fields": [
					{
					"name": "config_digest",
					"type": {
						"Vector": "U8"
					}
					},
					{
					"name": "big_f",
					"type": "U8"
					},
					{
					"name": "n",
					"type": "U8"
					},
					{
					"name": "is_signature_verification_enabled",
					"type": "Bool"
					}
				],
				"typeParameters": []
				},
				"Counter": {
				"abilities": {
					"abilities": [
					"Store",
					"Key"
					]
				},
				"fields": [
					{
					"name": "id",
					"type": {
						"Struct": {
						"address": "0x2",
						"module": "object",
						"name": "UID",
						"typeArguments": []
						}
					}
					},
					{
					"name": "value",
					"type": "U64"
					}
				],
				"typeParameters": []
				},
				"CounterIncremented": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop"
					]
				},
				"fields": [
					{
					"name": "counter_id",
					"type": {
						"Struct": {
						"address": "0x2",
						"module": "object",
						"name": "ID",
						"typeArguments": []
						}
					}
					},
					{
					"name": "new_value",
					"type": "U64"
					}
				],
				"typeParameters": []
				},
				"CounterPointer": {
				"abilities": {
					"abilities": [
					"Store",
					"Key"
					]
				},
				"fields": [
					{
					"name": "id",
					"type": {
						"Struct": {
						"address": "0x2",
						"module": "object",
						"name": "UID",
						"typeArguments": []
						}
					}
					},
					{
					"name": "counter_id",
					"type": "Address"
					},
					{
					"name": "admin_cap_id",
					"type": "Address"
					}
				],
				"typeParameters": []
				},
				"MultiNestedStruct": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop"
					]
				},
				"fields": [
					{
					"name": "is_multi_nested",
					"type": "Bool"
					},
					{
					"name": "double_count",
					"type": "U64"
					},
					{
					"name": "nested_struct",
					"type": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "NestedStruct",
						"typeArguments": []
						}
					}
					},
					{
					"name": "nested_simple_struct",
					"type": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "SimpleResult",
						"typeArguments": []
						}
					}
					}
				],
				"typeParameters": []
				},
				"NestedStruct": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop"
					]
				},
				"fields": [
					{
					"name": "is_nested",
					"type": "Bool"
					},
					{
					"name": "double_count",
					"type": "U64"
					},
					{
					"name": "nested_struct",
					"type": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "ComplexResult",
						"typeArguments": []
						}
					}
					},
					{
					"name": "nested_simple_struct",
					"type": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "SimpleResult",
						"typeArguments": []
						}
					}
					}
				],
				"typeParameters": []
				},
				"OCRConfig": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop",
					"Store"
					]
				},
				"fields": [
					{
					"name": "config_info",
					"type": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "ConfigInfo",
						"typeArguments": []
						}
					}
					},
					{
					"name": "signers",
					"type": {
						"Vector": {
						"Vector": "U8"
						}
					}
					},
					{
					"name": "transmitters",
					"type": {
						"Vector": "Address"
					}
					}
				],
				"typeParameters": []
				},
				"SimpleResult": {
				"abilities": {
					"abilities": [
					"Copy",
					"Drop"
					]
				},
				"fields": [
					{
					"name": "value",
					"type": "U64"
					}
				],
				"typeParameters": []
				}
			},
			"exposedFunctions": {
				"array_size": {
				"isEntry": false,
				"parameters": [
					{
					"Vector": {
						"TypeParameter": 0
					}
					}
				],
				"return": [
					"U64"
				],
				"typeParameters": [
					{
					"abilities": [
						"Drop"
					]
					}
				],
				"visibility": "Public"
				},
				"create": {
				"isEntry": false,
				"parameters": [
					{
					"MutableReference": {
						"Struct": {
						"address": "0x2",
						"module": "tx_context",
						"name": "TxContext",
						"typeArguments": []
						}
					}
					}
				],
				"return": [
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_address_list": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "AddressList",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_count": {
				"isEntry": true,
				"parameters": [
					{
					"Reference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					}
				],
				"return": [
					"U64"
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_count_no_entry": {
				"isEntry": false,
				"parameters": [
					{
					"Reference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					}
				],
				"return": [
					"U64"
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_count_using_pointer": {
				"isEntry": true,
				"parameters": [
					{
					"Reference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					}
				],
				"return": [
					"U64"
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_multi_nested_result_struct": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "MultiNestedStruct",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_nested_result_struct": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "NestedStruct",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_ocr_config": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "OCRConfig",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_result_struct": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "ComplexResult",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_simple_result": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "SimpleResult",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_tuple_struct": {
				"isEntry": false,
				"parameters": [],
				"return": [
					"U64",
					"Address",
					"Bool",
					{
					"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "MultiNestedStruct",
						"typeArguments": []
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_vector_of_addresses": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Vector": "Address"
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_vector_of_u8": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Vector": "U8"
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"get_vector_of_vectors_of_u8": {
				"isEntry": false,
				"parameters": [],
				"return": [
					{
					"Vector": {
						"Vector": "U8"
					}
					}
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"increment": {
				"isEntry": true,
				"parameters": [
					{
					"MutableReference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					}
				],
				"return": [],
				"typeParameters": [],
				"visibility": "Public"
				},
				"increment_by": {
				"isEntry": true,
				"parameters": [
					{
					"MutableReference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					},
					"U64"
				],
				"return": [],
				"typeParameters": [],
				"visibility": "Public"
				},
				"increment_by_one": {
				"isEntry": false,
				"parameters": [
					{
					"MutableReference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					},
					{
					"MutableReference": {
						"Struct": {
						"address": "0x2",
						"module": "tx_context",
						"name": "TxContext",
						"typeArguments": []
						}
					}
					}
				],
				"return": [
					"U64"
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"increment_by_one_no_context": {
				"isEntry": false,
				"parameters": [
					{
					"MutableReference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					}
				],
				"return": [
					"U64"
				],
				"typeParameters": [],
				"visibility": "Public"
				},
				"increment_by_two": {
				"isEntry": false,
				"parameters": [
					{
					"Reference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "AdminCap",
						"typeArguments": []
						}
					}
					},
					{
					"MutableReference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					},
					{
					"MutableReference": {
						"Struct": {
						"address": "0x2",
						"module": "tx_context",
						"name": "TxContext",
						"typeArguments": []
						}
					}
					}
				],
				"return": [],
				"typeParameters": [],
				"visibility": "Public"
				},
				"increment_by_two_no_context": {
				"isEntry": true,
				"parameters": [
					{
					"Reference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "AdminCap",
						"typeArguments": []
						}
					}
					},
					{
					"MutableReference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					}
				],
				"return": [],
				"typeParameters": [],
				"visibility": "Public"
				},
				"increment_mult": {
				"isEntry": true,
				"parameters": [
					{
					"MutableReference": {
						"Struct": {
						"address": "0x66b827fe66f3bc50c9deef27624c7b705a1d0af4a8b0883280d729c728559b71",
						"module": "counter",
						"name": "Counter",
						"typeArguments": []
						}
					}
					},
					"U64",
					"U64",
					{
					"MutableReference": {
						"Struct": {
						"address": "0x2",
						"module": "tx_context",
						"name": "TxContext",
						"typeArguments": []
						}
					}
					}
				],
				"return": [],
				"typeParameters": [],
				"visibility": "Public"
				},
				"initialize": {
				"isEntry": true,
				"parameters": [
					{
					"MutableReference": {
						"Struct": {
						"address": "0x2",
						"module": "tx_context",
						"name": "TxContext",
						"typeArguments": []
						}
					}
					}
				],
				"return": [],
				"typeParameters": [],
				"visibility": "Public"
				}
			}
		}
		`

		var schema map[string]any
		err := json.Unmarshal([]byte(jsonModule), &schema)
		require.NoError(t, err)

		structs, ok := schema["structs"].(map[string]any)
		require.True(t, ok)

		bcsBytes := []byte{
			32, 0, 10, 163, 58, 124, 44, 0, 205, 25, 175, 172, 143, 227, 22, 8, 175, 42, 52, 252, 74, 32, 10, 107, 236, 80, 1, 177, 162, 131, 82, 115, 71, 1, 4, 1, 4, 32, 153, 199, 47, 13, 162, 190, 48, 139, 149, 84, 92, 112, 93, 249, 186, 231, 136, 123, 47, 47, 228, 6, 126, 60, 15, 225, 137, 169, 88, 36, 111, 223, 32, 185, 6, 58, 13, 126, 86, 237, 190, 192, 150, 150, 19, 74, 225, 21, 7, 83, 19, 164, 225, 70, 37, 68, 140, 138, 155, 195, 18, 14, 201, 54, 184, 32, 237, 81, 16, 104, 37, 16, 243, 198, 124, 89, 11, 86, 195, 24, 18, 132, 120, 108, 13, 25, 116, 159, 64, 190, 1, 184, 175, 103, 72, 18, 122, 255, 32, 213, 192, 56, 2, 175, 151, 186, 105, 250, 60, 206, 8, 54, 91, 208, 80, 45, 64, 142, 15, 45, 182, 101, 87, 125, 144, 114, 146, 189, 165, 130, 187, 4, 51, 226, 204, 208, 210, 225, 178, 251, 215, 161, 23, 19, 167, 250, 208, 102, 88, 245, 8, 211, 230, 10, 7, 91, 68, 202, 111, 169, 46, 217, 9, 137, 88, 2, 84, 235, 187, 66, 243, 57, 245, 194, 18, 92, 179, 94, 242, 121, 119, 226, 188, 125, 133, 223, 136, 196, 186, 122, 104, 225, 215, 140, 230, 170, 118, 124, 155, 157, 80, 221, 219, 194, 82, 0, 141, 107, 133, 15, 228, 127, 248, 169, 211, 92, 82, 137, 101, 86, 107, 86, 17, 193, 42, 182, 100, 61, 63, 88, 117, 124, 145, 87, 42, 74, 17, 60, 67, 61, 23, 200, 219, 8, 212, 84, 233, 97, 22, 211, 228, 125, 79, 118, 102, 0, 252, 175, 97, 116,
		}

		deserializer := aptosBCS.NewDeserializer(bcsBytes)

		jsonMap, err := DecodeSuiStructToJSON(
			structs,
			"OCRConfig",
			deserializer,
		)

		require.NoError(t, err)
		utils.PrettyPrint(jsonMap)
	})

	t.Run("JSON Struct Decoder (String case)", func(t *testing.T) {
		t.Parallel()

		jsonModule := `
		{
			"fileFormatVersion": 6,
			"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
			"name": "token_admin_registry",
			"friends": [],
			"structs": {
				"AdministratorTransferRequested": {
					"abilities": {
						"abilities": [
							"Copy",
							"Drop"
						]
					},
					"fields": [
						{
							"name": "coin_metadata_address",
							"type": "Address"
						},
						{
							"name": "current_admin",
							"type": "Address"
						},
						{
							"name": "new_admin",
							"type": "Address"
						}
					],
					"typeParameters": []
				},
				"AdministratorTransferred": {
					"abilities": {
						"abilities": [
							"Copy",
							"Drop"
						]
					},
					"fields": [
						{
							"name": "coin_metadata_address",
							"type": "Address"
						},
						{
							"name": "new_admin",
							"type": "Address"
						}
					],
					"typeParameters": []
				},
				"PoolInfos": {
					"abilities": {
						"abilities": [
							"Copy",
							"Drop"
						]
					},
					"fields": [
						{
							"name": "token_pool_package_ids",
							"type": {
								"Vector": "Address"
							}
						},
						{
							"name": "token_pool_state_addresses",
							"type": {
								"Vector": "Address"
							}
						},
						{
							"name": "token_pool_modules",
							"type": {
								"Vector": {
									"Struct": {
										"address": "0x1",
										"module": "string",
										"name": "String",
										"typeArguments": []
									}
								}
							}
						},
						{
							"name": "token_types",
							"type": {
								"Vector": {
									"Struct": {
										"address": "0x1",
										"module": "ascii",
										"name": "String",
										"typeArguments": []
									}
								}
							}
						}
					],
					"typeParameters": []
				},
				"PoolRegistered": {
					"abilities": {
						"abilities": [
							"Copy",
							"Drop"
						]
					},
					"fields": [
						{
							"name": "coin_metadata_address",
							"type": "Address"
						},
						{
							"name": "token_pool_package_id",
							"type": "Address"
						},
						{
							"name": "administrator",
							"type": "Address"
						},
						{
							"name": "type_proof",
							"type": {
								"Struct": {
									"address": "0x1",
									"module": "ascii",
									"name": "String",
									"typeArguments": []
								}
							}
						}
					],
					"typeParameters": []
				},
				"PoolSet": {
					"abilities": {
						"abilities": [
							"Copy",
							"Drop"
						]
					},
					"fields": [
						{
							"name": "coin_metadata_address",
							"type": "Address"
						},
						{
							"name": "previous_pool_package_id",
							"type": "Address"
						},
						{
							"name": "new_pool_package_id",
							"type": "Address"
						},
						{
							"name": "type_proof",
							"type": {
								"Struct": {
									"address": "0x1",
									"module": "ascii",
									"name": "String",
									"typeArguments": []
								}
							}
						}
					],
					"typeParameters": []
				},
				"PoolUnregistered": {
					"abilities": {
						"abilities": [
							"Copy",
							"Drop"
						]
					},
					"fields": [
						{
							"name": "coin_metadata_address",
							"type": "Address"
						},
						{
							"name": "previous_pool_address",
							"type": "Address"
						}
					],
					"typeParameters": []
				},
				"TokenAdminRegistryState": {
					"abilities": {
						"abilities": [
							"Store",
							"Key"
						]
					},
					"fields": [
						{
							"name": "id",
							"type": {
								"Struct": {
									"address": "0x2",
									"module": "object",
									"name": "UID",
									"typeArguments": []
								}
							}
						},
						{
							"name": "token_configs",
							"type": {
								"Struct": {
									"address": "0x2",
									"module": "linked_table",
									"name": "LinkedTable",
									"typeArguments": [
										"Address",
										{
											"Struct": {
												"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
												"module": "token_admin_registry",
												"name": "TokenConfig",
												"typeArguments": []
											}
										}
									]
								}
							}
						}
					],
					"typeParameters": []
				},
				"TokenConfig": {
					"abilities": {
						"abilities": [
							"Copy",
							"Drop",
							"Store"
						]
					},
					"fields": [
						{
							"name": "token_pool_package_id",
							"type": "Address"
						},
						{
							"name": "token_pool_state_address",
							"type": "Address"
						},
						{
							"name": "token_pool_module",
							"type": {
								"Struct": {
									"address": "0x1",
									"module": "string",
									"name": "String",
									"typeArguments": []
								}
							}
						},
						{
							"name": "token_type",
							"type": {
								"Struct": {
									"address": "0x1",
									"module": "ascii",
									"name": "String",
									"typeArguments": []
								}
							}
						},
						{
							"name": "administrator",
							"type": "Address"
						},
						{
							"name": "pending_administrator",
							"type": "Address"
						},
						{
							"name": "type_proof",
							"type": {
								"Struct": {
									"address": "0x1",
									"module": "ascii",
									"name": "String",
									"typeArguments": []
								}
							}
						}
					],
					"typeParameters": []
				}
			},
			"exposedFunctions": {
				"accept_admin_role": {
					"isEntry": false,
					"parameters": [
						{
							"MutableReference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address",
						{
							"MutableReference": {
								"Struct": {
									"address": "0x2",
									"module": "tx_context",
									"name": "TxContext",
									"typeArguments": []
								}
							}
						}
					],
					"return": [],
					"typeParameters": [],
					"visibility": "Public"
				},
				"get_all_configured_tokens": {
					"isEntry": false,
					"parameters": [
						{
							"Reference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address",
						"U64"
					],
					"return": [
						{
							"Vector": "Address"
						},
						"Address",
						"Bool"
					],
					"typeParameters": [],
					"visibility": "Public"
				},
				"get_pool": {
					"isEntry": false,
					"parameters": [
						{
							"Reference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address"
					],
					"return": [
						"Address"
					],
					"typeParameters": [],
					"visibility": "Public"
				},
				"get_pool_infos": {
					"isEntry": false,
					"parameters": [
						{
							"Reference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						{
							"Vector": "Address"
						}
					],
					"return": [
						{
							"Struct": {
								"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
								"module": "token_admin_registry",
								"name": "PoolInfos",
								"typeArguments": []
							}
						}
					],
					"typeParameters": [],
					"visibility": "Public"
				},
				"get_pools": {
					"isEntry": false,
					"parameters": [
						{
							"Reference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						{
							"Vector": "Address"
						}
					],
					"return": [
						{
							"Vector": "Address"
						}
					],
					"typeParameters": [],
					"visibility": "Public"
				},
				"get_token_config": {
					"isEntry": false,
					"parameters": [
						{
							"Reference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address"
					],
					"return": [
						"Address",
						"Address",
						{
							"Struct": {
								"address": "0x1",
								"module": "string",
								"name": "String",
								"typeArguments": []
							}
						},
						{
							"Struct": {
								"address": "0x1",
								"module": "ascii",
								"name": "String",
								"typeArguments": []
							}
						},
						"Address",
						"Address",
						{
							"Struct": {
								"address": "0x1",
								"module": "ascii",
								"name": "String",
								"typeArguments": []
							}
						}
					],
					"typeParameters": [],
					"visibility": "Public"
				},
				"initialize": {
					"isEntry": false,
					"parameters": [
						{
							"MutableReference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						{
							"Reference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "OwnerCap",
									"typeArguments": []
								}
							}
						},
						{
							"MutableReference": {
								"Struct": {
									"address": "0x2",
									"module": "tx_context",
									"name": "TxContext",
									"typeArguments": []
								}
							}
						}
					],
					"return": [],
					"typeParameters": [],
					"visibility": "Public"
				},
				"is_administrator": {
					"isEntry": false,
					"parameters": [
						{
							"Reference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address",
						"Address"
					],
					"return": [
						"Bool"
					],
					"typeParameters": [],
					"visibility": "Public"
				},
				"register_pool": {
					"isEntry": false,
					"parameters": [
						{
							"MutableReference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						{
							"Reference": {
								"Struct": {
									"address": "0x2",
									"module": "coin",
									"name": "TreasuryCap",
									"typeArguments": [
										{
											"TypeParameter": 0
										}
									]
								}
							}
						},
						{
							"Reference": {
								"Struct": {
									"address": "0x2",
									"module": "coin",
									"name": "CoinMetadata",
									"typeArguments": [
										{
											"TypeParameter": 0
										}
									]
								}
							}
						},
						"Address",
						"Address",
						{
							"Struct": {
								"address": "0x1",
								"module": "string",
								"name": "String",
								"typeArguments": []
							}
						},
						"Address",
						{
							"TypeParameter": 1
						}
					],
					"return": [],
					"typeParameters": [
						{
							"abilities": []
						},
						{
							"abilities": [
								"Drop"
							]
						}
					],
					"visibility": "Public"
				},
				"register_pool_by_admin": {
					"isEntry": false,
					"parameters": [
						{
							"MutableReference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address",
						"Address",
						"Address",
						{
							"Struct": {
								"address": "0x1",
								"module": "string",
								"name": "String",
								"typeArguments": []
							}
						},
						{
							"Struct": {
								"address": "0x1",
								"module": "ascii",
								"name": "String",
								"typeArguments": []
							}
						},
						"Address",
						{
							"Struct": {
								"address": "0x1",
								"module": "ascii",
								"name": "String",
								"typeArguments": []
							}
						},
						{
							"MutableReference": {
								"Struct": {
									"address": "0x2",
									"module": "tx_context",
									"name": "TxContext",
									"typeArguments": []
								}
							}
						}
					],
					"return": [],
					"typeParameters": [],
					"visibility": "Public"
				},
				"set_pool": {
					"isEntry": false,
					"parameters": [
						{
							"MutableReference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address",
						"Address",
						"Address",
						{
							"Struct": {
								"address": "0x1",
								"module": "string",
								"name": "String",
								"typeArguments": []
							}
						},
						{
							"TypeParameter": 0
						},
						{
							"MutableReference": {
								"Struct": {
									"address": "0x2",
									"module": "tx_context",
									"name": "TxContext",
									"typeArguments": []
								}
							}
						}
					],
					"return": [],
					"typeParameters": [
						{
							"abilities": [
								"Drop"
							]
						}
					],
					"visibility": "Public"
				},
				"transfer_admin_role": {
					"isEntry": false,
					"parameters": [
						{
							"MutableReference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address",
						"Address",
						{
							"MutableReference": {
								"Struct": {
									"address": "0x2",
									"module": "tx_context",
									"name": "TxContext",
									"typeArguments": []
								}
							}
						}
					],
					"return": [],
					"typeParameters": [],
					"visibility": "Public"
				},
				"type_and_version": {
					"isEntry": false,
					"parameters": [],
					"return": [
						{
							"Struct": {
								"address": "0x1",
								"module": "string",
								"name": "String",
								"typeArguments": []
							}
						}
					],
					"typeParameters": [],
					"visibility": "Public"
				},
				"unregister_pool": {
					"isEntry": false,
					"parameters": [
						{
							"MutableReference": {
								"Struct": {
									"address": "0x93e6aa7d7efd881081c4523eb1225d7345da0c2e3188067a43e09ac58ade00f2",
									"module": "state_object",
									"name": "CCIPObjectRef",
									"typeArguments": []
								}
							}
						},
						"Address",
						{
							"MutableReference": {
								"Struct": {
									"address": "0x2",
									"module": "tx_context",
									"name": "TxContext",
									"typeArguments": []
								}
							}
						}
					],
					"return": [],
					"typeParameters": [],
					"visibility": "Public"
				}
			}
		}
		`

		var schema map[string]any
		err := json.Unmarshal([]byte(jsonModule), &schema)
		require.NoError(t, err)

		structs, ok := schema["structs"].(map[string]any)
		require.True(t, ok)

		bcsBytes := []byte{1, 203, 131, 104, 67, 247, 212, 213, 63, 234, 0, 172, 253, 104, 80, 167, 198, 167, 122, 227, 165, 211, 158, 27, 218, 248, 148, 160, 212, 92, 144, 255, 107, 1, 5, 248, 236, 185, 248, 254, 11, 147, 254, 142, 242, 226, 6, 109, 14, 31, 63, 215, 241, 104, 164, 118, 6, 179, 25, 147, 35, 171, 213, 189, 127, 240, 1, 23, 108, 111, 99, 107, 95, 114, 101, 108, 101, 97, 115, 101, 95, 116, 111, 107, 101, 110, 95, 112, 111, 111, 108, 1, 88, 101, 57, 48, 98, 55, 52, 57, 101, 97, 53, 53, 50, 57, 100, 101, 49, 52, 51, 52, 49, 102, 55, 51, 54, 57, 98, 50, 53, 52, 98, 50, 56, 50, 54, 56, 101, 54, 99, 48, 50, 57, 97, 55, 50, 51, 57, 53, 51, 56, 99, 48, 54, 49, 98, 101, 97, 53, 56, 97, 54, 52, 50, 100, 50, 58, 58, 108, 105, 110, 107, 95, 116, 111, 107, 101, 110, 58, 58, 76, 73, 78, 75, 95, 84, 79, 75, 69, 78}

		deserializer := aptosBCS.NewDeserializer(bcsBytes)

		jsonMap, err := DecodeSuiStructToJSON(
			structs,
			"PoolInfos",
			deserializer,
		)

		require.NoError(t, err)
		utils.PrettyPrint(jsonMap)
	})
}
