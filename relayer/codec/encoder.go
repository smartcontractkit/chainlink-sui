package codec

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func EncodeFunctionParams(argMap map[string]any, params []SuiFunctionParam) ([]string, []any, error) {
	types := make([]string, len(params))
	values := make([]any, len(params))

	for i, paramConfig := range params {
		argValue, ok := argMap[paramConfig.Name]
		if !ok {
			if paramConfig.Required {
				return nil, nil, fmt.Errorf("missing argument: %s", paramConfig.Name)
			}
			argValue = paramConfig.DefaultValue
		}

		types[i] = paramConfig.Type
		values[i] = argValue
	}

	return types, values, nil
}

func EncodeToSuiValue(typeName string, value any) (any, error) {
	switch typeName {
	case "address":
		return encodeAddress(value)
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

func encodeUint(typeName string, value any) (uint64, error) {
	switch v := value.(type) {
	case int:
		if v < 0 {
			return 0, fmt.Errorf("cannot convert negative int %d to %s", v, typeName)
		}

		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("cannot convert negative int %d to %s", v, typeName)
		}

		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint64:
		return v, nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("cannot convert negative int %f to %s", v, typeName)
		}

		return uint64(v), nil
	case string:
		i, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string %s to %s: %w", v, typeName, err)
		}

		return i, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to %s", value, typeName)
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

func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
