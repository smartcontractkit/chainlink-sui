//go:build unit

package codec

import (
	"math/big"
	"testing"

	aptosBCS "github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBCSTypeConverter_DecodePrimitive(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *aptosBCS.Deserializer
		typeStr  string
		expected any
		wantErr  bool
	}{
		{
			name: "U8",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.U8(42)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "U8",
			expected: uint8(42),
		},
		{
			name: "u8 lowercase",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.U8(42)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "u8",
			expected: uint8(42),
		},
		{
			name: "U16",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.U16(1234)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "U16",
			expected: uint16(1234),
		},
		{
			name: "U32",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.U32(123456)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "U32",
			expected: uint32(123456),
		},
		{
			name: "U64 as string",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.U64(12345678901234)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "U64",
			expected: "12345678901234",
		},
		{
			name: "U128 as string",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				val := big.NewInt(123456789012345678)
				s.U128(*val)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "U128",
			expected: "123456789012345678",
		},
		{
			name: "U256 as string",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				val := big.NewInt(0)
				val.SetString("123456789012345678901234567890", 10)
				s.U256(*val)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "U256",
			expected: "123456789012345678901234567890",
		},
		{
			name: "Bool true",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Bool(true)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "Bool",
			expected: true,
		},
		{
			name: "Bool false",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Bool(false)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr:  "bool",
			expected: false,
		},
		{
			name: "Address",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				addr := make([]byte, 32)
				for i := range addr {
					addr[i] = byte(i)
				}
				s.FixedBytes(addr)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			typeStr: "Address",
			expected: func() []byte {
				addr := make([]byte, 32)
				for i := range addr {
					addr[i] = byte(i)
				}
				return addr
			}(),
		},
		{
			name: "Unsupported type",
			setup: func() *aptosBCS.Deserializer {
				return aptosBCS.NewDeserializer([]byte{})
			},
			typeStr: "InvalidType",
			wantErr: true,
		},
	}

	converter := NewBCSTypeConverter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deserializer := tt.setup()
			result, err := converter.DecodePrimitive(deserializer, tt.typeStr)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBCSTypeConverter_DecodeVector(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *aptosBCS.Deserializer
		elemType string
		expected any
		wantErr  bool
	}{
		{
			name: "Vector<U8> as bytes",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Uleb128(3) // length
				s.U8(1)
				s.U8(2)
				s.U8(3)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			elemType: "U8",
			expected: []byte{1, 2, 3},
		},
		{
			name: "Vector<U64> as uint64 slice",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Uleb128(2) // length
				s.U64(100)
				s.U64(200)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			elemType: "U64",
			expected: []uint64{100, 200},
		},
		{
			name: "Vector<Address>",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Uleb128(2) // length
				addr1 := make([]byte, 32)
				addr1[0] = 1
				s.FixedBytes(addr1)
				addr2 := make([]byte, 32)
				addr2[0] = 2
				s.FixedBytes(addr2)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			elemType: "Address",
			expected: []any{
				func() []byte {
					addr := make([]byte, 32)
					addr[0] = 1
					return addr
				}(),
				func() []byte {
					addr := make([]byte, 32)
					addr[0] = 2
					return addr
				}(),
			},
		},
		{
			name: "Vector of primitives (fallback to generic)",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Uleb128(2) // length
				s.U16(100)
				s.U16(200)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			elemType: "U16",
			expected: []any{uint16(100), uint16(200)},
		},
		{
			name: "Empty vector",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Uleb128(0) // length
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			elemType: "U8",
			expected: []byte{},
		},
		{
			name: "Unsupported element type",
			setup: func() *aptosBCS.Deserializer {
				s := &aptosBCS.Serializer{}
				s.Uleb128(1)
				return aptosBCS.NewDeserializer(s.ToBytes())
			},
			elemType: "InvalidType",
			wantErr:  true,
		},
	}

	converter := NewBCSTypeConverter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deserializer := tt.setup()
			result, err := converter.DecodeVector(deserializer, tt.elemType)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBCSTypeConverter_CustomRegistration(t *testing.T) {
	converter := NewBCSTypeConverter()

	// Register a custom type
	converter.RegisterPrimitive("CustomType", func(d *aptosBCS.Deserializer) (any, error) {
		return "custom_value", nil
	})

	// Test the custom type
	s := &aptosBCS.Serializer{}
	// Doesn't matter what we serialize for this test
	s.U8(0)
	deserializer := aptosBCS.NewDeserializer(s.ToBytes())

	result, err := converter.DecodePrimitive(deserializer, "CustomType")
	require.NoError(t, err)
	assert.Equal(t, "custom_value", result)
}

func TestBCSTypeConverter_HasHandlers(t *testing.T) {
	converter := NewBCSTypeConverter()

	// Test primitive handlers
	assert.True(t, converter.HasPrimitiveHandler("U8"))
	assert.True(t, converter.HasPrimitiveHandler("u64"))
	assert.True(t, converter.HasPrimitiveHandler("Bool"))
	assert.False(t, converter.HasPrimitiveHandler("NonExistent"))

	// Test vector handlers
	assert.True(t, converter.HasVectorHandler("U8"))
	assert.True(t, converter.HasVectorHandler("U64"))
	assert.False(t, converter.HasVectorHandler("NonExistent"))
}

func TestGlobalBCSConverter(t *testing.T) {
	// Test that the global converter works
	s := &aptosBCS.Serializer{}
	s.U8(42)
	deserializer := aptosBCS.NewDeserializer(s.ToBytes())

	result, err := DecodeBCSPrimitive(deserializer, "U8")
	require.NoError(t, err)
	assert.Equal(t, uint8(42), result)
}
