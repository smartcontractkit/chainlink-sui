package codec

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
)

// TypeConversionFunc defines a function that converts data from one type to another
type TypeConversionFunc func(from reflect.Type, to reflect.Type, data any) (any, error)

// TypeConverter provides a registry approach to type conversions for mapstructure
type TypeConverter struct {
	converters map[string]TypeConversionFunc
}

// NewTypeConverter creates a new type converter with all standard conversions registered
func NewTypeConverter() *TypeConverter {
	tc := &TypeConverter{
		converters: make(map[string]TypeConversionFunc),
	}

	tc.registerStandardConverters()
	return tc
}

// registerStandardConverters registers all standard type conversion handlers
func (tc *TypeConverter) registerStandardConverters() {
	// Hex string conversions
	tc.RegisterConverter("hex_string_to_string", tc.hexToString)
	tc.RegisterConverter("hex_string_to_bytes", tc.hexToBytes)
	tc.RegisterConverter("hex_string_to_uint", tc.hexToUint)
	tc.RegisterConverter("hex_string_to_int", tc.hexToInt)
	tc.RegisterConverter("hex_string_to_bigint", tc.hexToBigInt)
	tc.RegisterConverter("hex_string_to_array", tc.hexToArray)

	// Base64 string conversions
	tc.RegisterConverter("base64_string_to_bytes", tc.base64ToBytes)

	// Numeric string conversions
	tc.RegisterConverter("string_to_int", tc.stringToInt)
	tc.RegisterConverter("string_to_uint", tc.stringToUint)
	tc.RegisterConverter("string_to_float", tc.stringToFloat)
	tc.RegisterConverter("string_to_bytes", tc.stringToBytes)
	tc.RegisterConverter("string_to_bigint", tc.stringToBigInt)

	// Boolean conversions
	tc.RegisterConverter("bool_to_int", tc.boolToInt)
	tc.RegisterConverter("bool_to_uint", tc.boolToUint)
	tc.RegisterConverter("bool_to_bigint", tc.boolToBigInt)

	// Array/Slice conversions
	tc.RegisterConverter("slice_to_slice", tc.sliceToSlice)
}

// RegisterConverter registers a conversion function with a unique key
func (tc *TypeConverter) RegisterConverter(key string, fn TypeConversionFunc) {
	tc.converters[key] = fn
}

// Convert attempts to convert data using registered converters
func (tc *TypeConverter) Convert(from reflect.Type, to reflect.Type, data any) (any, error) {
	// Try hex string conversions
	if from.Kind() == reflect.String {
		if str, ok := data.(string); ok && strings.HasPrefix(str, "0x") {
			result, err, handled := tc.handleHexString(str, to, data)
			if handled {
				return result, err
			}
		}
	}

	// Try base64 conversions
	if from.Kind() == reflect.String && (to.Kind() == reflect.Slice || to.Kind() == reflect.Array) && to.Elem().Kind() == reflect.Uint8 {
		result, err, handled := tc.handleBase64String(data.(string), to)
		if handled {
			return result, err
		}
	}

	// Try numeric string conversions
	if from.Kind() == reflect.String {
		result, err, handled := tc.handleNumericString(data.(string), to)
		if handled {
			return result, err
		}
	}

	// Try boolean conversions
	if from.Kind() == reflect.Bool {
		result, err, handled := tc.handleBoolean(data.(bool), to)
		if handled {
			return result, err
		}
	}

	// Try array/slice conversions
	if (from.Kind() == reflect.Slice || from.Kind() == reflect.Array) && to.Kind() == reflect.Slice {
		result, err, handled := tc.handleSlice(from, to, data)
		if handled {
			return result, err
		}
	}

	// No conversion found, return data as-is
	return data, nil
}

// handleHexString handles hex string conversions
func (tc *TypeConverter) handleHexString(str string, to reflect.Type, data any) (result any, err error, handled bool) {
	hexStr := strings.TrimPrefix(str, "0x")

	switch to.Kind() {
	case reflect.String:
		if fn, ok := tc.converters["hex_string_to_string"]; ok {
			result, err = fn(reflect.TypeOf(str), to, hexStr)
			return result, err, true
		}
	case reflect.Slice:
		if to.Elem().Kind() == reflect.Uint8 {
			if fn, ok := tc.converters["hex_string_to_bytes"]; ok {
				result, err = fn(reflect.TypeOf(str), to, hexStr)
				return result, err, true
			}
		}
	case reflect.Array:
		if to.Elem().Kind() == reflect.Uint8 {
			if fn, ok := tc.converters["hex_string_to_array"]; ok {
				result, err = fn(reflect.TypeOf(str), to, hexStr)
				return result, err, true
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fn, ok := tc.converters["hex_string_to_uint"]; ok {
			result, err = fn(reflect.TypeOf(str), to, hexStr)
			return result, err, true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fn, ok := tc.converters["hex_string_to_int"]; ok {
			result, err = fn(reflect.TypeOf(str), to, hexStr)
			return result, err, true
		}
	case reflect.Ptr:
		if to == reflect.TypeOf((*big.Int)(nil)) {
			if fn, ok := tc.converters["hex_string_to_bigint"]; ok {
				result, err = fn(reflect.TypeOf(str), to, hexStr)
				return result, err, true
			}
		}
	case reflect.Interface:
		return "0x" + hexStr, nil, true
	}

	return nil, nil, false
}

// handleBase64String handles base64 string conversions
func (tc *TypeConverter) handleBase64String(str string, to reflect.Type) (result any, err error, handled bool) {
	if fn, ok := tc.converters["base64_string_to_bytes"]; ok {
		result, err = fn(reflect.TypeOf(str), to, str)
		return result, err, true
	}
	return nil, nil, false
}

// handleNumericString handles numeric string conversions
func (tc *TypeConverter) handleNumericString(str string, to reflect.Type) (result any, err error, handled bool) {
	switch to.Kind() {
	case reflect.String:
		return str, nil, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fn, ok := tc.converters["string_to_int"]; ok {
			result, err = fn(reflect.TypeOf(str), to, str)
			return result, err, true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fn, ok := tc.converters["string_to_uint"]; ok {
			result, err = fn(reflect.TypeOf(str), to, str)
			return result, err, true
		}
	case reflect.Float32, reflect.Float64:
		if fn, ok := tc.converters["string_to_float"]; ok {
			result, err = fn(reflect.TypeOf(str), to, str)
			return result, err, true
		}
	case reflect.Slice:
		if to.Elem().Kind() == reflect.Uint8 {
			if fn, ok := tc.converters["string_to_bytes"]; ok {
				result, err = fn(reflect.TypeOf(str), to, str)
				return result, err, true
			}
		}
	case reflect.Ptr:
		if to == reflect.TypeOf((*big.Int)(nil)) {
			if fn, ok := tc.converters["string_to_bigint"]; ok {
				result, err = fn(reflect.TypeOf(str), to, str)
				return result, err, true
			}
		}
	}

	return nil, nil, false
}

// handleBoolean handles boolean conversions
func (tc *TypeConverter) handleBoolean(boolValue bool, to reflect.Type) (result any, err error, handled bool) {
	switch to.Kind() {
	case reflect.Bool:
		return boolValue, nil, true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fn, ok := tc.converters["bool_to_int"]; ok {
			result, err = fn(reflect.TypeOf(boolValue), to, boolValue)
			return result, err, true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if fn, ok := tc.converters["bool_to_uint"]; ok {
			result, err = fn(reflect.TypeOf(boolValue), to, boolValue)
			return result, err, true
		}
	case reflect.Ptr:
		if to == reflect.TypeOf((*big.Int)(nil)) {
			if fn, ok := tc.converters["bool_to_bigint"]; ok {
				result, err = fn(reflect.TypeOf(boolValue), to, boolValue)
				return result, err, true
			}
		}
	}

	return nil, nil, false
}

// handleSlice handles array/slice conversions
func (tc *TypeConverter) handleSlice(from, to reflect.Type, data any) (result any, err error, handled bool) {
	if fn, ok := tc.converters["slice_to_slice"]; ok {
		result, err = fn(from, to, data)
		return result, err, true
	}
	return nil, nil, false
}

// Conversion function implementations

func (tc *TypeConverter) hexToString(from, to reflect.Type, data any) (any, error) {
	hexStr, ok := data.(string)
	if !ok {
		return data, nil
	}
	return hexStr, nil
}

func (tc *TypeConverter) hexToBytes(from, to reflect.Type, data any) (any, error) {
	hexStr, ok := data.(string)
	if !ok {
		return data, nil
	}

	if hexStr == "" {
		return []uint8{}, nil
	}

	if len(hexStr)%2 == 1 {
		hexStr = "0" + hexStr
	}

	return hex.DecodeString(hexStr)
}

func (tc *TypeConverter) hexToUint(from, to reflect.Type, data any) (any, error) {
	hexStr, ok := data.(string)
	if !ok {
		return data, nil
	}

	return strconv.ParseUint(hexStr, base16, uint64Bits)
}

func (tc *TypeConverter) hexToInt(from, to reflect.Type, data any) (any, error) {
	hexStr, ok := data.(string)
	if !ok {
		return data, nil
	}

	val, err := strconv.ParseInt(hexStr, base16, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hex to int: %w", err)
	}

	return reflect.ValueOf(val).Convert(to).Interface(), nil
}

func (tc *TypeConverter) hexToBigInt(from, to reflect.Type, data any) (any, error) {
	hexStr, ok := data.(string)
	if !ok {
		return data, nil
	}

	bi := new(big.Int)
	bi.SetString(hexStr, base16)
	return bi, nil
}

func (tc *TypeConverter) hexToArray(from, to reflect.Type, data any) (any, error) {
	hexStr, ok := data.(string)
	if !ok {
		return data, nil
	}

	bytes, err := tc.hexToBytes(from, reflect.SliceOf(to.Elem()), hexStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string %q: %w", hexStr, err)
	}

	byteSlice := bytes.([]byte)
	out := make([]uint8, to.Len())
	copy(out, byteSlice)

	return out, nil
}

func (tc *TypeConverter) base64ToBytes(from, to reflect.Type, data any) (any, error) {
	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	bytes, err := base64.StdEncoding.DecodeString(str)
	if err == nil {
		return bytes, nil
	}

	return data, nil
}

func (tc *TypeConverter) stringToInt(from, to reflect.Type, data any) (any, error) {
	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	val, err := strconv.ParseInt(str, base10, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string to int: %w", err)
	}

	if overflowInt(to, val) {
		return nil, fmt.Errorf("value %d overflows %v", val, to)
	}

	return reflect.ValueOf(val).Convert(to).Interface(), nil
}

func (tc *TypeConverter) stringToUint(from, to reflect.Type, data any) (any, error) {
	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	val, err := strconv.ParseUint(str, base10, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string to uint: %w", err)
	}

	if overflowUint(to, val) {
		return nil, fmt.Errorf("value %d overflows %v", val, to)
	}

	return reflect.ValueOf(val).Convert(to).Interface(), nil
}

func (tc *TypeConverter) stringToFloat(from, to reflect.Type, data any) (any, error) {
	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	val, err := strconv.ParseFloat(str, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string to float: %w", err)
	}

	if overflowFloat(to, val) {
		return nil, fmt.Errorf("value %f overflows %v", val, to)
	}

	return reflect.ValueOf(val).Convert(to).Interface(), nil
}

func (tc *TypeConverter) stringToBytes(from, to reflect.Type, data any) (any, error) {
	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	if num, err := strconv.ParseUint(str, base10, uint64Bits); err == nil {
		return numericToBytes(num), nil
	}

	return str, nil
}

func (tc *TypeConverter) stringToBigInt(from, to reflect.Type, data any) (any, error) {
	str, ok := data.(string)
	if !ok {
		return data, nil
	}

	result := new(big.Int)
	_, success := result.SetString(str, base10)
	if !success {
		return nil, fmt.Errorf("cannot parse string %s as big.Int", str)
	}

	return result, nil
}

func (tc *TypeConverter) boolToInt(from, to reflect.Type, data any) (any, error) {
	boolValue, ok := data.(bool)
	if !ok {
		return data, nil
	}

	if boolValue {
		return reflect.ValueOf(1).Convert(to).Interface(), nil
	}

	return reflect.ValueOf(0).Convert(to).Interface(), nil
}

func (tc *TypeConverter) boolToUint(from, to reflect.Type, data any) (any, error) {
	boolValue, ok := data.(bool)
	if !ok {
		return data, nil
	}

	if boolValue {
		return reflect.ValueOf(1).Convert(to).Interface(), nil
	}

	return reflect.ValueOf(0).Convert(to).Interface(), nil
}

func (tc *TypeConverter) boolToBigInt(from, to reflect.Type, data any) (any, error) {
	boolValue, ok := data.(bool)
	if !ok {
		return data, nil
	}

	if boolValue {
		return big.NewInt(1), nil
	}

	return big.NewInt(0), nil
}

func (tc *TypeConverter) sliceToSlice(from, to reflect.Type, data any) (any, error) {
	sourceSlice := reflect.ValueOf(data)
	targetSlice := reflect.MakeSlice(to, sourceSlice.Len(), sourceSlice.Cap())

	for i := range sourceSlice.Len() {
		sourceElem := sourceSlice.Index(i).Interface()
		targetElem := reflect.New(to.Elem()).Interface()

		if err := DecodeSuiJsonValue(sourceElem, targetElem); err != nil {
			return nil, fmt.Errorf("failed to decode array element at index %d: %w", i, err)
		}

		targetSlice.Index(i).Set(reflect.ValueOf(targetElem).Elem())
	}

	return targetSlice.Interface(), nil
}

// Global type converter instance for package-wide use (lazy initialization)
var defaultTypeConverter *TypeConverter

// getDefaultTypeConverter returns the global type converter, initializing it if necessary
func getDefaultTypeConverter() *TypeConverter {
	if defaultTypeConverter == nil {
		defaultTypeConverter = NewTypeConverter()
	}
	return defaultTypeConverter
}

// UnifiedTypeConverterHook is a mapstructure hook that uses the type converter registry
func UnifiedTypeConverterHook(from, to reflect.Type, data any) (any, error) {
	// Skip if types are the same
	if from == to {
		return data, nil
	}

	// Handle single-field struct case only for primitive source types
	// (string, bool, numeric types) to avoid interfering with map-to-struct decoding
	if to.Kind() == reflect.Struct && to.NumField() == 1 {
		switch from.Kind() {
		case reflect.String, reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Slice, reflect.Array:
			return handleSingleFieldStruct(to, data, DecodeSuiJsonValue)
		}
	}

	// Use the global converter
	return getDefaultTypeConverter().Convert(from, to, data)
}
