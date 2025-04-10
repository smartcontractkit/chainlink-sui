package codec

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/mitchellh/mapstructure"
)

const BYTE_SIZE = 8

// DecodeSuiJsonValue takes Sui JSON-RPC response data and decodes it into the provided target
func DecodeSuiJsonValue(data any, target any) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	// If data is already in the right format, just assign it
	if reflect.TypeOf(data) == reflect.TypeOf(target).Elem() {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(data))
		return nil
	}

	// Handle different types of data
	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()

	//nolint:exhaustive // default case handles remaining kinds
	switch targetType.Kind() {
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		// Handle numeric types
		return decodeNumeric(data, targetValue)
	case reflect.String:
		// Handle string type
		if str, ok := data.(string); ok {
			targetValue.SetString(str)
			return nil
		}

		return fmt.Errorf("expected string, got %T", data)
	case reflect.Slice:
		// Handle slices
		return decodeSlice(data, targetValue)
	case reflect.Struct:
		// Use mapstructure for struct types
		config := &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				numericStringHook,
				mapstructure.StringToTimeDurationHookFunc(),
			),
			Result:           target,
			WeaklyTypedInput: true,
			TagName:          "json",
		}

		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			return fmt.Errorf("failed to create decoder: %w", err)
		}

		return decoder.Decode(data)
	default:
		// Attempt direct JSON unmarshaling for other types
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}

		return json.Unmarshal(jsonBytes, target)
	}
}

// decodeNumeric handles numeric types (u64, u32, etc.)
func decodeNumeric(data any, targetValue reflect.Value) error {
	switch v := data.(type) {
	case float64:
		//nolint:exhaustive // default case handles remaining kinds
		switch targetValue.Kind() {
		case reflect.Uint64:
			targetValue.SetUint(uint64(v))
		case reflect.Uint32:
			targetValue.SetUint(uint64(v))
		case reflect.Uint16:
			targetValue.SetUint(uint64(v))
		case reflect.Uint8:
			targetValue.SetUint(uint64(v))
		default:
			return fmt.Errorf("unsupported target type for numeric value: %s", targetValue.Type())
		}

		return nil
	case string:
		// Numeric values can be returned as strings in JSON
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse string as number: %w", err)
		}
		targetValue.SetUint(n)

		return nil
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return fmt.Errorf("failed to parse JSON number: %w", err)
		}
		if n < 0 {
			return fmt.Errorf("cannot convert negative value %d to uint", n)
		}
		targetValue.SetUint(uint64(n))

		return nil
	case []byte:
		if len(v) > 0 {
			var result uint64
			// Process bytes in little-endian order (least significant byte first)
			for i := 0; i < len(v) && i < BYTE_SIZE; i++ {
				result |= uint64(v[i]) << (BYTE_SIZE * i)
			}
			targetValue.SetUint(result)

			return nil
		}

		return fmt.Errorf("empty byte array cannot be converted to numeric value")
	default:
		return fmt.Errorf("unsupported data type for numeric target: %T", data)
	}
}

// decodeSlice handles slice types
func decodeSlice(data any, targetValue reflect.Value) error {
	sourceSlice, ok := data.([]any)
	if !ok {
		return fmt.Errorf("expected slice, got %T", data)
	}

	elemType := targetValue.Type().Elem()
	slice := reflect.MakeSlice(targetValue.Type(), len(sourceSlice), len(sourceSlice))

	for i, item := range sourceSlice {
		elemValue := reflect.New(elemType)
		if err := DecodeSuiJsonValue(item, elemValue.Interface()); err != nil {
			return fmt.Errorf("failed to decode slice element at index %d: %w", i, err)
		}
		slice.Index(i).Set(elemValue.Elem())
	}

	targetValue.Set(slice)

	return nil
}

func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// Add new numericStringHook for mapstructure
func numericStringHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	switch t.Kind() {
	case reflect.String:
		return data, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string to uint: %w", err)
		}

		return reflect.ValueOf(val).Convert(t).Interface(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string to int: %w", err)
		}

		return reflect.ValueOf(val).Convert(t).Interface(), nil
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string to float: %w", err)
		}

		return reflect.ValueOf(val).Convert(t).Interface(), nil
	case reflect.Invalid, reflect.Bool, reflect.Complex64, reflect.Complex128,
		reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer, reflect.Uintptr:
		return data, nil
	default:
		return data, nil
	}
}
