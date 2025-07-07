package codec

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

const (
	// Type bounds
	maxUint8  = 255
	maxUint16 = 65535
	maxUint32 = 4294967295
)

// Safe conversion functions to avoid lint issues
func safeUint8(val uint64) (uint8, error) {
	if val > maxUint8 {
		return 0, fmt.Errorf("value %d exceeds uint8 maximum", val)
	}

	return uint8(val), nil
}

func safeUint16(val uint64) (uint16, error) {
	if val > maxUint16 {
		return 0, fmt.Errorf("value %d exceeds uint16 maximum", val)
	}

	return uint16(val), nil
}

func safeUint32(val uint64) (uint32, error) {
	if val > maxUint32 {
		return 0, fmt.Errorf("value %d exceeds uint32 maximum", val)
	}

	return uint32(val), nil
}

func EncodeToSuiValue(typeName string, value any) (any, error) {
	switch typeName {
	case "address":
		return encodeAddress(value)
	case "object_id":
		return encodeObject(value)
	case "u8", "u16", "u32", "u64", "u128", "u256":
		return encodeUint(typeName, value)
	case "bool":
		return encodeBool(value)
	case "string":
		return encodeString(value)
	default:
		if strings.HasPrefix(typeName, "0x") && strings.Contains(typeName, "::") {
			// TODO: need to use go-bsc to encode this. Reference https://github.com/fardream/go-bcs/blob/main/bcs/encode_test.go
			return nil, errors.New("struct types are not supported")
		}
		if strings.HasPrefix(typeName, "vector<") && strings.HasSuffix(typeName, ">") {
			return encodeVector(typeName, value)
		}

		return nil, fmt.Errorf("unsupported type: %s", typeName)
	}
}

func encodeObject(value any) (bind.Object, error) {
	switch v := value.(type) {
	case bind.Object:
		return bind.Object{Id: v.Id}, nil
	default:
		return bind.Object{}, fmt.Errorf("cannot convert %T to object", value)
	}
}

func encodeAddress(value any) (string, error) {
	switch v := value.(type) {
	case string:
		// Ensure it's a valid Sui address format
		if !strings.HasPrefix(v, "0x") {
			v = "0x" + v
		}

		return v, nil
	case []byte:
		return "0x" + hex.EncodeToString(v), nil
	default:
		return "", fmt.Errorf("cannot convert %T to address", value)
	}
}

func encodeUint(typeName string, value any) (any, error) {
	// First convert to a common intermediate type
	var baseValue uint64
	var bigIntValue *big.Int

	switch v := value.(type) {
	case int:
		if v < 0 {
			return nil, fmt.Errorf("cannot convert negative int %d to %s", v, typeName)
		}
		baseValue = uint64(v)
	case int64:
		if v < 0 {
			return nil, fmt.Errorf("cannot convert negative int %d to %s", v, typeName)
		}
		baseValue = uint64(v)
	case uint:
		baseValue = uint64(v)
	case uint64:
		baseValue = v
	case uint32:
		baseValue = uint64(v)
	case uint16:
		baseValue = uint64(v)
	case uint8:
		baseValue = uint64(v)
	case float64:
		if v < 0 {
			return nil, fmt.Errorf("cannot convert negative float %f to %s", v, typeName)
		}
		baseValue = uint64(v)
	case json.Number:
		// Handle JSON numbers properly
		if strings.Contains(string(v), ".") {
			f, err := v.Float64()
			if err != nil {
				return nil, fmt.Errorf("cannot convert json.Number %s to %s: %w", v, typeName, err)
			}
			baseValue = uint64(f)
		} else {
			i, err := v.Int64()
			if err != nil {
				return nil, fmt.Errorf("cannot convert json.Number %s to %s: %w", v, typeName, err)
			}
			//nolint:gosec
			// we assume safe conversion without negative numbers
			baseValue = uint64(i)
		}
	case string:
		// Handle big numbers that might come as strings
		if typeName == "u128" || typeName == "u256" {
			bigIntValue = new(big.Int)
			_, ok := bigIntValue.SetString(v, base10)
			if !ok {
				return nil, fmt.Errorf("cannot convert string %s to %s", v, typeName)
			}
		} else {
			i, err := strconv.ParseUint(v, base10, 64)
			if err != nil {
				return nil, fmt.Errorf("cannot convert string %s to %s: %w", v, typeName, err)
			}
			baseValue = i
		}
	case *big.Int:
		bigIntValue = v
	default:
		return nil, fmt.Errorf("cannot convert %T to %s", value, typeName)
	}

	// Now convert to the appropriate type based on typeName
	switch typeName {
	case "u8":
		if bigIntValue != nil {
			if !bigIntValue.IsUint64() || bigIntValue.Uint64() > maxUint8 {
				return nil, fmt.Errorf("value %s too large for u8", bigIntValue.String())
			}
			// Safe conversion after bounds check
			return safeUint8(bigIntValue.Uint64())
		}
		if baseValue > maxUint8 {
			return nil, fmt.Errorf("value %d too large for u8", baseValue)
		}

		return safeUint8(baseValue)
	case "u16":
		if bigIntValue != nil {
			if !bigIntValue.IsUint64() || bigIntValue.Uint64() > maxUint16 {
				return nil, fmt.Errorf("value %s too large for u16", bigIntValue.String())
			}
			// Safe conversion after bounds check
			return safeUint16(bigIntValue.Uint64())
		}
		if baseValue > maxUint16 {
			return nil, fmt.Errorf("value %d too large for u16", baseValue)
		}

		return safeUint16(baseValue)
	case "u32":
		if bigIntValue != nil {
			if !bigIntValue.IsUint64() || bigIntValue.Uint64() > maxUint32 {
				return nil, fmt.Errorf("value %s too large for u32", bigIntValue.String())
			}
			// Safe conversion after bounds check
			return safeUint32(bigIntValue.Uint64())
		}
		if baseValue > maxUint32 {
			return nil, fmt.Errorf("value %d too large for u32", baseValue)
		}

		return safeUint32(baseValue)
	case "u64":
		if bigIntValue != nil {
			if !bigIntValue.IsUint64() {
				return nil, fmt.Errorf("value %s too large for u64", bigIntValue.String())
			}

			return bigIntValue.Uint64(), nil
		}

		return baseValue, nil
	case "u128":
		if bigIntValue != nil {
			// Validate it fits in u128 (2^128 - 1)
			maxVal := new(big.Int)
			maxVal.Exp(big.NewInt(base2), big.NewInt(bits128), nil)
			maxVal.Sub(maxVal, big.NewInt(1))
			if bigIntValue.Cmp(maxVal) > 0 {
				return nil, fmt.Errorf("value %s too large for u128", bigIntValue.String())
			}

			return bigIntValue.String(), nil
		}

		return strconv.FormatUint(baseValue, base10), nil
	case "u256":
		if bigIntValue != nil {
			// Validate it fits in u256 (2^256 - 1)
			maxVal := new(big.Int)
			maxVal.Exp(big.NewInt(base2), big.NewInt(bits256), nil)
			maxVal.Sub(maxVal, big.NewInt(1))
			if bigIntValue.Cmp(maxVal) > 0 {
				return nil, fmt.Errorf("value %s too large for u256", bigIntValue.String())
			}

			return bigIntValue.String(), nil
		}

		return strconv.FormatUint(baseValue, base10), nil
	default:
		return nil, fmt.Errorf("unsupported uint type: %s", typeName)
	}
}

func encodeBool(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		b, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("cannot convert string %s to bool: %w", v, err)
		}

		return b, nil
	case int:
		return v != 0, nil
	case float64:
		return v != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

func encodeString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", value)
	}
}

func encodeVector(typeName string, value any) ([]any, error) {
	// Extract the inner type, e.g., "vector<string>" -> "string"
	if !strings.HasPrefix(typeName, "vector<") || !strings.HasSuffix(typeName, ">") {
		return nil, fmt.Errorf("invalid vector type: %s", typeName)
	}
	innerType := typeName[len("vector<") : len(typeName)-1]

	// Use reflection to ensure 'value' is a slice or array
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("expected a slice/array for vector type %s, got %T", typeName, value)
	}

	// Special handling for vector<u8> (byte arrays)
	if innerType == "u8" {
		// Handle []any from JSON unmarshaling
		if interfaceSlice, ok := value.([]any); ok {
			bytes := make([]byte, len(interfaceSlice))
			for i, item := range interfaceSlice {
				if num, ok := item.(float64); ok {
					if num < 0 || num > maxUint8 {
						return nil, fmt.Errorf("invalid byte value at index %d: %f", i, num)
					}
					bytes[i] = byte(num)
				} else if num, ok := item.(int); ok {
					if num < 0 || num > maxUint8 {
						return nil, fmt.Errorf("invalid byte value at index %d: %d", i, num)
					}
					bytes[i] = byte(num)
				} else {
					return nil, fmt.Errorf("invalid byte value at index %d: %T", i, item)
				}
			}

			return []any{bytes}, nil
		}
		// Handle []byte directly
		if bytes, ok := value.([]byte); ok {
			return []any{bytes}, nil
		}
	}

	encodedElements := make([]any, 0, rv.Len())
	for i := range rv.Len() {
		elem := rv.Index(i).Interface()
		encodedElem, err := EncodeToSuiValue(innerType, elem)
		if err != nil {
			return nil, fmt.Errorf("failed to encode element at index %d: %w", i, err)
		}
		encodedElements = append(encodedElements, encodedElem)
	}

	return encodedElements, nil
}
