package common

import (
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

// inferArgumentType attempts to determine the argument type from the value
// TODO: this method shouldn't be needed, it's a fallback for when the argument type is not known
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
