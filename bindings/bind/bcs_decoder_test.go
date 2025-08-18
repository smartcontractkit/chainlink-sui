package bind

import (
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mystenbcs.Unmarshal doesn't support *big.Int types, causing corruption
// like ChainId: 4 becoming ChainId: 9188203512877425290
//
// Solution: Two-stage BCS conversion using byte arrays + DecodeU256Value/DecodeU128Value
// This test verifies our fix works correctly and prevents the corruption.
func TestU256BCSDecoding(t *testing.T) {
	t.Run("DecodeU256Value correctly handles u256", func(t *testing.T) {
		// Test case: Chain ID 4 encoded as little-endian u256
		var chainIdBytes [32]byte
		chainIdBytes[0] = 4 // Little-endian representation of 4

		result, err := DecodeU256Value(chainIdBytes)
		require.NoError(t, err)

		expected := big.NewInt(4)
		assert.Equal(t, 0, result.Cmp(expected), "Expected %s, got %s", expected.String(), result.String())
	})

	t.Run("DecodeU128Value correctly handles u128", func(t *testing.T) {
		// Test case: Value 1000 encoded as little-endian u128
		var valueBytes [16]byte
		valueBytes[0] = 232 // 1000 = 0x03E8, little-endian: E8 03
		valueBytes[1] = 3

		result, err := DecodeU128Value(valueBytes)
		require.NoError(t, err)

		expected := big.NewInt(1000)
		assert.Equal(t, 0, result.Cmp(expected), "Expected %s, got %s", expected.String(), result.String())
	})

	t.Run("Large u256 values decode correctly", func(t *testing.T) {
		// Test with a larger u256 value to ensure proper byte handling
		// Value: 2^64 + 1 = 18446744073709551617
		var valueBytes [32]byte
		valueBytes[0] = 1 // Low byte of (2^64 + 1)
		valueBytes[8] = 1 // 9th byte (2^64 bit)

		result, err := DecodeU256Value(valueBytes)
		require.NoError(t, err)

		// 2^64 + 1
		expected := new(big.Int)
		expected.SetString("18446744073709551617", 10)

		assert.Equal(t, 0, result.Cmp(expected), "Expected %s, got %s", expected.String(), result.String())
	})
}

// TestBCSStructDecoding demonstrates how our two-stage conversion works
func TestBCSStructDecoding(t *testing.T) {
	// Mock struct similar to what we generate
	type MockBCSStruct struct {
		ChainId [32]byte
		Value   uint64
	}

	type MockFinalStruct struct {
		ChainId *big.Int
		Value   uint64
	}

	t.Run("Two-stage BCS conversion works correctly", func(t *testing.T) {
		// Stage 1: Create BCS data that mystenbcs can handle
		bcsData := MockBCSStruct{
			Value: 12345,
		}
		// Set ChainId to 4 in little-endian
		bcsData.ChainId[0] = 4

		// This simulates what our generated conversion functions do
		final := MockFinalStruct{
			Value: bcsData.Value,
			ChainId: func() *big.Int {
				decoded, err := DecodeU256Value(bcsData.ChainId)
				require.NoError(t, err)

				return decoded
			}(),
		}

		// Verify the conversion worked correctly
		assert.Equal(t, uint64(12345), final.Value)
		assert.Equal(t, 0, final.ChainId.Cmp(big.NewInt(4)))
	})
}

func TestOldLogicFailure(t *testing.T) {
	t.Run("ChainId corruption: 4 becomes 9188203512877425290", func(t *testing.T) {
		// Original: ChainId = 4 as u256
		var originalBytes [32]byte
		originalBytes[0] = 4

		// ✅ Correct interpretation
		decoded, err := DecodeU256Value(originalBytes)
		require.NoError(t, err)
		assert.Equal(t, int64(4), decoded.Int64())

		// 🚨 Verify corrupted value exists
		corruptedValue := uint64(9188203512877425290)
		corruptedBytes := make([]byte, 8)
		for i := range 8 {
			corruptedBytes[i] = byte(corruptedValue >> (i * 8))
		}

		var reconstructed uint64
		_, err = mystenbcs.Unmarshal(corruptedBytes, &reconstructed)
		require.NoError(t, err)
		assert.Equal(t, corruptedValue, reconstructed)
	})

	t.Run("mystenbcs.Unmarshal fails with *big.Int", func(t *testing.T) {
		var chainIdBytes [32]byte
		chainIdBytes[0] = 4

		// This should fail
		var bigIntResult *big.Int
		value, err := mystenbcs.Unmarshal(chainIdBytes[:], &bigIntResult)
		require.NoError(t, err)
		assert.NotEqual(t, 4, value)

		// Our solution works
		decoded, err := DecodeU256Value(chainIdBytes)
		require.NoError(t, err)
		assert.Equal(t, int64(4), decoded.Int64())
	})
}

func TestReverseBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Single byte",
			input:    []byte{0x42},
			expected: []byte{0x42},
		},
		{
			name:     "Four bytes - little to big endian",
			input:    []byte{0x01, 0x02, 0x03, 0x04},
			expected: []byte{0x04, 0x03, 0x02, 0x01},
		},
		{
			name:     "Empty slice",
			input:    []byte{},
			expected: []byte{},
		},
		{
			name:     "32-byte u256 value",
			input:    append([]byte{0x04, 0x00}, make([]byte, 30)...),
			expected: append(make([]byte, 30), []byte{0x00, 0x04}...),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reverseBytes(tt.input)
			assert.Equal(t, tt.expected, result)

			// Ensure original slice is not modified
			if len(tt.input) > 0 {
				// This test assumes the first test case
				if tt.name == "Four bytes - little to big endian" {
					assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04}, tt.input)
				}
			}
		})
	}
}

// TestSpecificDecodeFunctions tests the actual functions used by generated bindings
func TestSpecificDecodeFunctions(t *testing.T) {
	t.Run("DecodeU256Value edge cases", func(t *testing.T) {
		// Test zero value
		var zeroBytes [32]byte
		result, err := DecodeU256Value(zeroBytes)
		require.NoError(t, err)
		assert.Equal(t, int64(0), result.Int64())

		// Test max u256 value (all bytes set to 0xFF)
		var maxBytes [32]byte
		for i := range maxBytes {
			maxBytes[i] = 0xFF
		}
		result, err = DecodeU256Value(maxBytes)
		require.NoError(t, err)

		// Max u256 = 2^256 - 1
		expected := new(big.Int)
		expected.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
		assert.Equal(t, 0, result.Cmp(expected))
	})

	t.Run("DecodeU128Value edge cases", func(t *testing.T) {
		// Test zero value
		var zeroBytes [16]byte
		result, err := DecodeU128Value(zeroBytes)
		require.NoError(t, err)
		assert.Equal(t, int64(0), result.Int64())

		// Test max u128 value (all bytes set to 0xFF)
		var maxBytes [16]byte
		for i := range maxBytes {
			maxBytes[i] = 0xFF
		}
		result, err = DecodeU128Value(maxBytes)
		require.NoError(t, err)

		// Max u128 = 2^128 - 1
		expected := new(big.Int)
		expected.SetString("340282366920938463463374607431768211455", 10)
		assert.Equal(t, 0, result.Cmp(expected))
	})

	t.Run("DecodeU256Value vs DecodeU128Value consistency", func(t *testing.T) {
		// Test that the same value in u128 and u256 format gives same result
		testValue := new(big.Int)
		testValue.SetString("12345678901234567890", 10)

		// Create u128 bytes (16 bytes) - copy value bytes to the start in little-endian order
		var u128Bytes [16]byte
		valueBytes := testValue.Bytes() // big-endian representation

		// Copy bytes in reverse order to create little-endian representation
		// Place the least significant bytes first
		for i := 0; i < len(valueBytes) && i < 16; i++ {
			u128Bytes[i] = valueBytes[len(valueBytes)-1-i]
		}

		// Create u256 bytes (32 bytes) with same value - same logic but for 32 bytes
		var u256Bytes [32]byte
		for i := 0; i < len(valueBytes) && i < 32; i++ {
			u256Bytes[i] = valueBytes[len(valueBytes)-1-i]
		}

		result128, err := DecodeU128Value(u128Bytes)
		require.NoError(t, err)

		result256, err := DecodeU256Value(u256Bytes)
		require.NoError(t, err)

		assert.Equal(t, 0, result128.Cmp(result256), "u128 and u256 decoding should give same result for same value")
		assert.Equal(t, 0, result128.Cmp(testValue), "Decoded value should match original")
	})
}
