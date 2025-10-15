//go:build unit

package codec

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeConverter_HexConversions(t *testing.T) {
	tc := NewTypeConverter()

	tests := []struct {
		name     string
		data     string
		toType   reflect.Type
		expected any
		wantErr  bool
	}{
		{
			name:     "hex to string",
			data:     "0xdeadbeef",
			toType:   reflect.TypeOf(""),
			expected: "deadbeef",
		},
		{
			name:     "hex to bytes",
			data:     "0xdeadbeef",
			toType:   reflect.TypeOf([]byte{}),
			expected: []byte{0xde, 0xad, 0xbe, 0xef},
		},
		{
			name:     "hex to uint64",
			data:     "0xff",
			toType:   reflect.TypeOf(uint64(0)),
			expected: uint64(255),
		},
		{
			name:   "hex to big.Int",
			data:   "0x123456789abcdef",
			toType: reflect.TypeOf((*big.Int)(nil)),
			expected: func() *big.Int {
				bi := new(big.Int)
				bi.SetString("123456789abcdef", 16)
				return bi
			}(),
		},
		{
			name:     "empty hex to bytes",
			data:     "0x",
			toType:   reflect.TypeOf([]byte{}),
			expected: []byte{},
		},
		{
			name:     "odd length hex to bytes (should pad)",
			data:     "0xabc",
			toType:   reflect.TypeOf([]byte{}),
			expected: []byte{0x0a, 0xbc},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tc.Convert(reflect.TypeOf(""), tt.toType, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTypeConverter_NumericStringConversions(t *testing.T) {
	tc := NewTypeConverter()

	tests := []struct {
		name     string
		data     string
		toType   reflect.Type
		expected any
		wantErr  bool
	}{
		{
			name:     "string to int64",
			data:     "12345",
			toType:   reflect.TypeOf(int64(0)),
			expected: int64(12345),
		},
		{
			name:     "string to uint64",
			data:     "67890",
			toType:   reflect.TypeOf(uint64(0)),
			expected: uint64(67890),
		},
		{
			name:     "string to float64",
			data:     "123.45",
			toType:   reflect.TypeOf(float64(0)),
			expected: float64(123.45),
		},
		{
			name:   "string to big.Int",
			data:   "999999999999999999999",
			toType: reflect.TypeOf((*big.Int)(nil)),
			expected: func() *big.Int {
				bi := new(big.Int)
				bi.SetString("999999999999999999999", 10)
				return bi
			}(),
		},
		{
			name:    "invalid string to int",
			data:    "not a number",
			toType:  reflect.TypeOf(int64(0)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tc.Convert(reflect.TypeOf(""), tt.toType, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTypeConverter_BooleanConversions(t *testing.T) {
	tc := NewTypeConverter()

	tests := []struct {
		name     string
		data     bool
		toType   reflect.Type
		expected any
	}{
		{
			name:     "bool true to int",
			data:     true,
			toType:   reflect.TypeOf(int(0)),
			expected: int(1),
		},
		{
			name:     "bool false to int",
			data:     false,
			toType:   reflect.TypeOf(int(0)),
			expected: int(0),
		},
		{
			name:     "bool true to uint",
			data:     true,
			toType:   reflect.TypeOf(uint(0)),
			expected: uint(1),
		},
		{
			name:     "bool false to uint",
			data:     false,
			toType:   reflect.TypeOf(uint(0)),
			expected: uint(0),
		},
		{
			name:     "bool true to big.Int",
			data:     true,
			toType:   reflect.TypeOf((*big.Int)(nil)),
			expected: big.NewInt(1),
		},
		{
			name:     "bool false to big.Int",
			data:     false,
			toType:   reflect.TypeOf((*big.Int)(nil)),
			expected: big.NewInt(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tc.Convert(reflect.TypeOf(true), tt.toType, tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTypeConverter_Base64Conversions(t *testing.T) {
	tc := NewTypeConverter()

	tests := []struct {
		name     string
		data     string
		toType   reflect.Type
		expected any
	}{
		{
			name:     "base64 to bytes",
			data:     "SGVsbG8gV29ybGQ=", // "Hello World"
			toType:   reflect.TypeOf([]byte{}),
			expected: []byte("Hello World"),
		},
		{
			name:     "empty base64",
			data:     "",
			toType:   reflect.TypeOf([]byte{}),
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tc.Convert(reflect.TypeOf(""), tt.toType, tt.data)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTypeConverter_CustomRegistration(t *testing.T) {
	tc := NewTypeConverter()

	// Register a custom converter
	tc.RegisterConverter("custom_test", func(from, to reflect.Type, data any) (any, error) {
		return "custom_result", nil
	})

	// Verify it's registered
	fn, ok := tc.converters["custom_test"]
	assert.True(t, ok)
	assert.NotNil(t, fn)

	// Test custom converter
	result, err := fn(reflect.TypeOf(""), reflect.TypeOf(""), "test_data")
	require.NoError(t, err)
	assert.Equal(t, "custom_result", result)
}

func TestUnifiedTypeConverterHook_HexString(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		fromType reflect.Type
		toType   reflect.Type
		expected any
		wantErr  bool
	}{
		{
			name:     "hex string to bytes",
			data:     "0xdeadbeef",
			fromType: reflect.TypeOf(""),
			toType:   reflect.TypeOf([]byte{}),
			expected: []byte{0xde, 0xad, 0xbe, 0xef},
		},
		{
			name:     "hex string to uint64",
			data:     "0xff",
			fromType: reflect.TypeOf(""),
			toType:   reflect.TypeOf(uint64(0)),
			expected: uint64(255),
		},
		{
			name:     "non-hex string (passthrough)",
			data:     "regular string",
			fromType: reflect.TypeOf(""),
			toType:   reflect.TypeOf(""),
			expected: "regular string",
		},
		{
			name:     "hex string to interface",
			data:     "0xabc",
			fromType: reflect.TypeOf(""),
			toType:   reflect.TypeOf((*any)(nil)).Elem(),
			expected: "0xabc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnifiedTypeConverterHook(tt.fromType, tt.toType, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnifiedTypeConverterHook_NumericString(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		toType   reflect.Type
		expected any
		wantErr  bool
	}{
		{
			name:     "numeric string to int",
			data:     "42",
			toType:   reflect.TypeOf(int(0)),
			expected: int(42),
		},
		{
			name:     "numeric string to uint",
			data:     "100",
			toType:   reflect.TypeOf(uint64(0)),
			expected: uint64(100),
		},
		{
			name:     "numeric string to float",
			data:     "3.14",
			toType:   reflect.TypeOf(float64(0)),
			expected: float64(3.14),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnifiedTypeConverterHook(reflect.TypeOf(""), tt.toType, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnifiedTypeConverterHook_SameType(t *testing.T) {
	// When types are the same, should return data as-is
	data := "test"
	result, err := UnifiedTypeConverterHook(reflect.TypeOf(""), reflect.TypeOf(""), data)
	require.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestTypeConverter_OverflowChecks(t *testing.T) {
	tc := NewTypeConverter()

	tests := []struct {
		name    string
		data    string
		toType  reflect.Type
		wantErr bool
	}{
		{
			name:    "uint8 overflow",
			data:    "256", // max uint8 is 255
			toType:  reflect.TypeOf(uint8(0)),
			wantErr: true,
		},
		{
			name:    "uint8 no overflow",
			data:    "255",
			toType:  reflect.TypeOf(uint8(0)),
			wantErr: false,
		},
		{
			name:    "int8 overflow",
			data:    "128", // max int8 is 127
			toType:  reflect.TypeOf(int8(0)),
			wantErr: true,
		},
		{
			name:    "int8 no overflow",
			data:    "127",
			toType:  reflect.TypeOf(int8(0)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tc.Convert(reflect.TypeOf(""), tt.toType, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
