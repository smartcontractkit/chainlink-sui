package codec

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const BYTE_SIZE = 8
const UINT8_BITS = 8
const UINT16_BITS = 16
const UINT32_BITS = 32
const UINT64_BITS = 64
const BASE_10 = 10
const BASE_16 = 16

// Additional constants for decoder
const (
	maxByteValue        = 255
	minResponseArrayLen = 2
	bitShift            = 8
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
		// Use mapstructure for struct types with hooks
		config := &mapstructure.DecoderConfig{
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				hexStringHook,
				base64StringHook,
				numericStringHook,
				booleanHook,
				arrayHook,
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
	// Handle string to []byte conversion
	if str, ok := data.(string); ok && targetValue.Type().Elem().Kind() == reflect.Uint8 {
		// Try to parse as numeric string first and convert to bytes
		if num, err := strconv.ParseUint(str, 10, UINT64_BITS); err == nil {
			// Convert number to byte slice (little-endian)
			bytes := make([]byte, UINT64_BITS/UINT8_BITS) // Use 8 bytes for uint64
			for i := range UINT8_BITS {
				bytes[i] = byte(num >> (i * UINT8_BITS))
			}
			// Remove trailing zeros
			for len(bytes) > 1 && bytes[len(bytes)-1] == 0 {
				bytes = bytes[:len(bytes)-1]
			}
			targetValue.Set(reflect.ValueOf(bytes))

			return nil
		}

		// Try hex decoding if numeric parsing failed
		if strings.HasPrefix(str, "0x") {
			hexStr := strings.TrimPrefix(str, "0x")
			if len(hexStr)%2 == 1 {
				hexStr = "0" + hexStr
			}
			bytes, err := hex.DecodeString(hexStr)
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}
			targetValue.Set(reflect.ValueOf(bytes))

			return nil
		}

		// Try base64 decoding
		if bytes, err := base64.StdEncoding.DecodeString(str); err == nil {
			targetValue.Set(reflect.ValueOf(bytes))
			return nil
		}

		// Otherwise convert string directly to bytes
		targetValue.Set(reflect.ValueOf([]byte(str)))

		return nil
	}

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

// hexStringHook handles hex string conversions (ported from Aptos)
func hexStringHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	str, ok := data.(string)
	if !ok || !strings.HasPrefix(str, "0x") {
		return data, nil
	}

	str = strings.TrimPrefix(str, "0x")

	// Handle single-field struct case first by recursing via DecodeSuiJsonValue
	if t.Kind() == reflect.Struct && t.NumField() == 1 {
		field := t.Field(0)
		// Create a new zero value struct
		newStructVal := reflect.New(t).Elem()
		// Get a pointer to the field within the new struct
		fieldPtr := newStructVal.Field(0).Addr().Interface()

		// Recursively decode the original hex string data into the field pointer
		if err := DecodeSuiJsonValue(data, fieldPtr); err != nil {
			return nil, fmt.Errorf("failed decoding hex string for single-field struct %v field %s (%v): %w", t, field.Name, field.Type, err)
		}
		// Return the populated struct instance
		return newStructVal.Interface(), nil
	}

	//nolint:exhaustive
	switch t.Kind() {
	case reflect.String:
		return data, nil
	case reflect.Slice:
		if t.Elem().Kind() != reflect.Uint8 {
			return nil, fmt.Errorf("unsupported target slice element type for hex string conversion: %v", t.Elem().Kind())
		}
		if str == "" {
			return []uint8{}, nil
		} else if len(str)%2 == 1 {
			str = "0" + str
		}

		return hex.DecodeString(str)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.ParseUint(str, BASE_16, UINT64_BITS)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(str, BASE_16, UINT64_BITS)
		if err != nil {
			return nil, fmt.Errorf("failed to parse hex to int: %w", err)
		}

		return reflect.ValueOf(val).Convert(t).Interface(), nil
	case reflect.Ptr:
		if t == reflect.TypeOf((*big.Int)(nil)) {
			bi := new(big.Int)
			bi.SetString(str, BASE_16)

			return bi, nil
		}
	case reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			if str == "" {
				return []uint8{}, nil
			} else if len(str)%2 == 1 {
				str = "0" + str
			}
			bytes, err := hex.DecodeString(str)
			if err != nil {
				return nil, fmt.Errorf("failed to decode hex string %q: %w", str, err)
			}
			out := make([]uint8, t.Len())
			copy(out, bytes)

			return out, nil
		}

		return nil, fmt.Errorf("unsupported target array element type for hex string conversion: %v", t.Elem().Kind())
	case reflect.Interface:
		// return the original value for type "any" as target
		return data, nil
	default:
	}

	return nil, fmt.Errorf("unsupported target type for hex string conversion: %v", t.Kind())
}

// base64StringHook handles base64 string conversions
func base64StringHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	// Handle single-field struct case
	if t.Kind() == reflect.Struct && t.NumField() == 1 {
		field := t.Field(0)
		newStructVal := reflect.New(t).Elem()
		fieldPtr := newStructVal.Field(0).Addr().Interface()

		if err := DecodeSuiJsonValue(data, fieldPtr); err != nil {
			return nil, fmt.Errorf("failed decoding base64 string for single-field struct %v field %s (%v): %w", t, field.Name, field.Type, err)
		}

		return newStructVal.Interface(), nil
	}

	// Only try base64 decoding for byte slices
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		// Try base64 decoding
		if bytes, err := base64.StdEncoding.DecodeString(str); err == nil {
			return bytes, nil
		}
	}

	return data, nil
}

// numericStringHook handles numeric string conversions (enhanced from existing)
func numericStringHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	// Handle single-field struct case
	if t.Kind() == reflect.Struct && t.NumField() == 1 {
		field := t.Field(0)
		newStructVal := reflect.New(t).Elem()
		fieldPtr := newStructVal.Field(0).Addr().Interface()

		if err := DecodeSuiJsonValue(data, fieldPtr); err != nil {
			return nil, fmt.Errorf("failed decoding numeric string for single-field struct %v field %s (%v): %w", t, field.Name, field.Type, err)
		}

		return newStructVal.Interface(), nil
	}

	//nolint:exhaustive
	switch t.Kind() {
	case reflect.String:
		return data, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(str, 10, UINT64_BITS)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string to int: %w", err)
		}
		if overflowInt(t, val) {
			return nil, fmt.Errorf("value %d overflows %v", val, t)
		}

		return reflect.ValueOf(val).Convert(t).Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(str, 10, UINT64_BITS)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string to uint: %w", err)
		}
		if overflowUint(t, val) {
			return nil, fmt.Errorf("value %d overflows %v", val, t)
		}

		return reflect.ValueOf(val).Convert(t).Interface(), nil
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(str, UINT64_BITS)
		if err != nil {
			return nil, fmt.Errorf("failed to parse string to float: %w", err)
		}
		if overflowFloat(t, val) {
			return nil, fmt.Errorf("value %f overflows %v", val, t)
		}

		return reflect.ValueOf(val).Convert(t).Interface(), nil
	case reflect.Slice:
		// Handle string to byte slice conversion for numeric strings
		if t.Elem().Kind() == reflect.Uint8 {
			// Parse as number and convert to bytes
			if num, err := strconv.ParseUint(str, 10, 64); err == nil {
				// Convert number to byte slice (little-endian)
				bytes := make([]byte, UINT8_BITS)
				for i := range UINT8_BITS {
					bytes[i] = byte(num >> (i * UINT8_BITS))
				}
				// Remove trailing zeros
				for len(bytes) > 1 && bytes[len(bytes)-1] == 0 {
					bytes = bytes[:len(bytes)-1]
				}

				return bytes, nil
			}
		}

		return data, nil
	case reflect.Ptr:
		if t == reflect.TypeOf((*big.Int)(nil)) {
			bi := new(big.Int)
			_, ok := bi.SetString(str, BASE_10)
			if !ok {
				return nil, fmt.Errorf("failed to parse string as big.Int: %s", str)
			}

			return bi, nil
		}
	default:
		return data, nil
	}

	return data, nil
}

// booleanHook handles boolean conversions
func booleanHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() != reflect.Bool {
		return data, nil
	}

	boolValue, ok := data.(bool)
	if !ok {
		return data, nil
	}

	// Handle single-field struct case
	if t.Kind() == reflect.Struct && t.NumField() == 1 {
		field := t.Field(0)
		newStructVal := reflect.New(t).Elem()
		fieldPtr := newStructVal.Field(0).Addr().Interface()

		if err := DecodeSuiJsonValue(data, fieldPtr); err != nil {
			return nil, fmt.Errorf("failed decoding boolean for single-field struct %v field %s (%v): %w", t, field.Name, field.Type, err)
		}

		return newStructVal.Interface(), nil
	}

	//nolint:exhaustive
	switch t.Kind() {
	case reflect.Bool:
		return boolValue, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if boolValue {
			return reflect.ValueOf(1).Convert(t).Interface(), nil
		}

		return reflect.ValueOf(0).Convert(t).Interface(), nil
	case reflect.Ptr:
		if t == reflect.TypeOf((*big.Int)(nil)) {
			if boolValue {
				return big.NewInt(1), nil
			}

			return big.NewInt(0), nil
		}
	default:
	}

	return nil, fmt.Errorf("unsupported target type for boolean conversion: %v", t.Kind())
}

// arrayHook handles array conversions
func arrayHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	fKind := f.Kind()
	if fKind != reflect.Slice && fKind != reflect.Array {
		return data, nil
	}

	// Handle single-field struct case
	if t.Kind() == reflect.Struct && t.NumField() == 1 {
		field := t.Field(0)
		newStructVal := reflect.New(t).Elem()
		fieldPtr := newStructVal.Field(0).Addr().Interface()

		if err := DecodeSuiJsonValue(data, fieldPtr); err != nil {
			return nil, fmt.Errorf("failed decoding array for single-field struct %v field %s (%v): %w", t, field.Name, field.Type, err)
		}

		return newStructVal.Interface(), nil
	}

	if t.Kind() != reflect.Slice {
		return data, nil
	}

	sourceSlice := reflect.ValueOf(data)
	targetSlice := reflect.MakeSlice(t, sourceSlice.Len(), sourceSlice.Cap())

	for i := range sourceSlice.Len() {
		sourceElem := sourceSlice.Index(i).Interface()
		targetElem := reflect.New(t.Elem()).Interface()

		if err := DecodeSuiJsonValue(sourceElem, targetElem); err != nil {
			return nil, fmt.Errorf("failed to decode array element at index %d: %w", i, err)
		}

		targetSlice.Index(i).Set(reflect.ValueOf(targetElem).Elem())
	}

	return targetSlice.Interface(), nil
}

// Overflow checking functions (ported from Aptos)
func overflowFloat(t reflect.Type, x float64) bool {
	k := t.Kind()
	//nolint:exhaustive
	switch k {
	case reflect.Float32:
		return overflowFloat32(x)
	case reflect.Float64:
		return false
	default:
	}
	panic("reflect: OverflowFloat of non-float type " + t.String())
}

func overflowFloat32(x float64) bool {
	if x < 0 {
		x = -x
	}

	return math.MaxFloat32 < x && x <= math.MaxFloat64
}

func overflowInt(t reflect.Type, x int64) bool {
	k := t.Kind()
	//nolint:exhaustive
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := t.Size() * UINT8_BITS
		trunc := (x << (UINT64_BITS - bitSize)) >> (UINT64_BITS - bitSize)

		return x != trunc
	default:
	}
	panic("reflect: OverflowFloat of non-float type " + t.String())
}

func overflowUint(t reflect.Type, x uint64) bool {
	k := t.Kind()
	//nolint:exhaustive
	switch k {
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bitSize := t.Size() * UINT8_BITS
		trunc := (x << (UINT64_BITS - bitSize)) >> (UINT64_BITS - bitSize)

		return x != trunc
	default:
	}
	panic("reflect: OverflowUint of non-uint type " + t.String())
}

// AnySliceToBytes converts a slice of interface{} into a byte slice.
// It returns an error if any element isn't a valid byte value.
func AnySliceToBytes(src []any) ([]byte, error) {
	dst := make([]byte, len(src))
	for i, v := range src {
		switch x := v.(type) {
		case uint8:
			dst[i] = x
		case int:
			if x < 0 || x > maxByteValue {
				return nil, fmt.Errorf("element %d: int %d out of byte range", i, x)
			}
			dst[i] = byte(x)
		case uint:
			if x > maxByteValue {
				return nil, fmt.Errorf("element %d: uint %d out of byte range", i, x)
			}
			dst[i] = byte(x)
		default:
			return nil, fmt.Errorf("element %d: unsupported type %T, need byte/int", i, v)
		}
	}

	return dst, nil
}

// ParseSuiResponseValue extracts the actual value from Sui's response format
// Sui responses come as [value, typeString] tuples, this function extracts the value part
func ParseSuiResponseValue(rawResponse any) (any, error) {
	responseArray, ok := rawResponse.([]any)
	if !ok {
		return nil, fmt.Errorf("expected Sui response to be an array, got %T", rawResponse)
	}

	if len(responseArray) < minResponseArrayLen {
		return nil, fmt.Errorf("expected Sui response array to have at least 2 elements, got %d", len(responseArray))
	}

	// Extract the actual value (first element) and type (second element)
	responseValue := responseArray[0]
	responseType, ok := responseArray[1].(string)
	if !ok {
		return nil, fmt.Errorf("expected second response element to be type string, got %T", responseArray[1])
	}

	// Handle different response structures based on type
	switch {
	case responseType == "u64" || responseType == "u32" || responseType == "u16" || responseType == "u8":
		// For uint types, convert byte array to actual number
		if byteArray, ok := responseValue.([]any); ok {
			// Determine expected byte length based on type
			var expectedBytes int
			switch responseType {
			case "u8":
				expectedBytes = 1
			case "u16":
				expectedBytes = 2
			case "u32":
				expectedBytes = 4
			case "u64":
				expectedBytes = 8
			default:
				expectedBytes = 8 // fallback
			}

			if len(byteArray) != expectedBytes {
				return nil, fmt.Errorf("expected %d bytes for %s, got %d", expectedBytes, responseType, len(byteArray))
			}

			// Convert byte array to uint64 (little-endian)
			var result uint64
			for i, v := range byteArray {
				num, ok := v.(float64)
				if !ok {
					return nil, fmt.Errorf("expected byte value at index %d, got %T", i, v)
				}
				result |= uint64(byte(num)) << (i * bitShift)
			}

			return result, nil
		}

		return responseValue, nil

	case responseType == "u128" || responseType == "u256":
		// For large uints, return as *big.Int so JSON marshaling works correctly in LOOP mode
		if byteArray, ok := responseValue.([]any); ok {
			// Convert byte array to big.Int (little-endian)
			result := new(big.Int)
			bytesArray, err := AnySliceToBytes(byteArray)
			if err != nil {
				return nil, err
			}
			result.SetBytes(bytesArray)

			return result, nil
		}
		// If it's already a number, convert to *big.Int
		if num, ok := responseValue.(float64); ok {
			return big.NewInt(int64(num)), nil
		}
		// If it's a string, parse as *big.Int
		if str, ok := responseValue.(string); ok {
			result := new(big.Int)
			_, ok := result.SetString(str, BASE_10)
			if !ok {
				return nil, fmt.Errorf("cannot parse string %s as big.Int", str)
			}

			return result, nil
		}

		return responseValue, nil

	case responseType == "bool":
		return responseValue, nil

	case strings.Contains(responseType, "string"):
		// Handle string types - may come as byte array
		if byteArray, ok := responseValue.([]any); ok {
			bytes := make([]byte, len(byteArray))
			for i, v := range byteArray {
				num, ok := v.(float64)
				if !ok {
					return nil, fmt.Errorf("expected byte value at index %d, got %T", i, v)
				}
				bytes[i] = byte(num)
			}

			return string(bytes), nil
		}

		return responseValue, nil

	case strings.HasPrefix(responseType, "vector<u8>"):
		// Return byte arrays as-is for vector<u8>
		return responseValue, nil

	case strings.HasPrefix(responseType, "vector<"):
		// Return other vectors as-is
		return responseValue, nil

	case strings.Contains(responseType, ","):
		// Handle tuples - return as map with numeric string keys to match struct JSON tags
		if tupleArray, ok := responseValue.([]any); ok {
			// Extract individual types from the tuple type string
			// e.g., "(u32, u64)" -> ["u32", "u64"]
			typeStr := strings.Trim(responseType, "()")
			types := strings.Split(typeStr, ", ")

			if len(tupleArray) != len(types) {
				return nil, fmt.Errorf("tuple length mismatch: expected %d elements, got %d", len(types), len(tupleArray))
			}

			// Return as map with string keys "0", "1", etc. to match JSON struct tags
			result := make(map[string]any)
			for i, item := range tupleArray {
				if i < len(types) {
					// Parse each tuple element according to its type
					elemType := strings.TrimSpace(types[i])

					// Create a fake response structure for each element to reuse our parsing logic
					fakeResponse := []any{item, elemType}
					parsedValue, err := ParseSuiResponseValue(fakeResponse)
					if err != nil {
						return nil, fmt.Errorf("failed to parse tuple element %d as %s: %w", i, elemType, err)
					}
					result[fmt.Sprintf("%d", i)] = parsedValue
				} else {
					result[fmt.Sprintf("%d", i)] = item
				}
			}

			return result, nil
		}

		return responseValue, nil

	default:
		// For unknown types, return the value as-is
		return responseValue, nil
	}
}
