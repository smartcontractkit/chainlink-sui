package codec

import (
	"encoding/base64"
	"encoding/json"
	"math"
	"math/big"
	"reflect"
	"testing"

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

	decoded, err := DecodeBase64(encoded)
	require.NoError(t, err)
	require.Equal(t, testData, decoded)

	// Test invalid base64
	_, err = DecodeBase64("invalid_base64!!!")
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
}
