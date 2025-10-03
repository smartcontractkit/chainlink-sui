package common

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"slices"
	"strings"
)

func ValueAt[T any](slice []T, idx int) (T, bool) {
	var zero T
	if idx < 0 || idx >= len(slice) {
		return zero, false
	}

	return slice[idx], true
}

// InferArgumentType attempts to determine the argument type from the value
// NOTE: this method shouldn't be needed, it's a fallback for when the argument type is not known
func InferArgumentType(arg any) string {
	switch arg := arg.(type) {
	case string:
		if strings.HasPrefix(arg, "0x") {
			return "objectId"
		}

		return "address"
	case []byte:
		return "vector<u8>"
	case uint64, int64:
		return "u64"
	case int:
		return "u64"
	case int32, uint32:
		return "u32"
	case int16, uint16:
		return "u16"
	case int8, uint8:
		return "u8"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}

func SerializeUBigInt(size uint, v *big.Int) []byte {
	ub := make([]byte, size)
	v.FillBytes(ub)
	// Reverse, since big.Int outputs bytes in BigEndian
	slices.Reverse(ub)

	return ub
}

// ConvertBytesToHex recursively walks through any value and hex-encodes all []byte values.
func ConvertBytesToHex(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for k, val := range v {
			v[k] = ConvertBytesToHex(val) // recursive
		}
		return v

	case []any:
		for i, val := range v {
			v[i] = ConvertBytesToHex(val) // recursive
		}
		return v

	case []uint8:
		// Confirm it's a real []byte and not some other []uint8 misuse
		return "0x" + hex.EncodeToString(v)

	case string:
		// length prevents any random string from being encoded
		if b, err := base64.StdEncoding.DecodeString(v); err == nil && len(b) == 32 {
			return "0x" + hex.EncodeToString(b)
		}
		return v

	default:
		return value
	}
}

func IsZeroAddress(address []byte) bool {
	return len(address) == 32 && bytes.Equal(address, make([]byte, 32))
}
