package codec

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

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
		// Handle structs by marshaling to JSON and then unmarshaling
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to marshal data: %w", err)
		}

		return json.Unmarshal(jsonBytes, target)
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
