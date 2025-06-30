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

	aptosBCS "github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/mitchellh/mapstructure"
	"github.com/pattonkan/sui-go/sui"
)

const (
	// Bit and byte constants
	byteSize     = 8
	uint8Bits    = 8
	uint8Bytes   = 1
	uint16Bits   = 16
	uint16Bytes  = 2
	uint32Bits   = 32
	uint32Bytes  = 4
	uint64Bits   = 64
	uint64Bytes  = 8
	bits128      = 128
	bits128Bytes = 16
	bits256      = 256
	bits256Bytes = 32

	// Number bases
	base10 = 10
	base16 = 16
	base2  = 2

	// Response parsing constants
	maxByteValue        = 255
	minResponseArrayLen = 2
	bitShift            = 8
)

// DecodeSuiJsonValue decodes Sui JSON-RPC response data into the provided target
func DecodeSuiJsonValue(data any, target any) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	// Direct type match optimization
	if reflect.TypeOf(data) == reflect.TypeOf(target).Elem() {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(data))
		return nil
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()

	//nolint:exhaustive
	switch targetType.Kind() {
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return decodeNumeric(data, targetValue)
	case reflect.String:
		return decodeString(data, targetValue)
	case reflect.Slice:
		return decodeSlice(data, targetValue)
	case reflect.Struct:
		return decodeStruct(data, target)
	default:
		return decodeGeneric(data, target)
	}
}

// DecodeSuiStructToJSON decodes a Sui struct into a JSON object
// using the normalized struct and the result
func DecodeSuiStructToJSON(normalizedStructs map[string]*sui.MoveNormalizedStruct, identifier string, bcsDecoder *aptosBCS.Deserializer) (map[string]any, error) {
	jsonResult := make(map[string]any)

	// Add nil check for the normalizedStructs map itself
	if normalizedStructs == nil {
		return nil, fmt.Errorf("normalizedStructs map is nil")
	}

	normalizedStruct := normalizedStructs[identifier]

	// Add nil check for the normalized struct
	if normalizedStruct == nil {
		return nil, fmt.Errorf("struct with identifier '%s' not found in normalized structs", identifier)
	}

	for _, field := range normalizedStruct.Fields {
		fieldType := field.Type
		fieldName := field.Name

		if fieldType.Address != nil {
			addressBytesLen := 32
			jsonResult[fieldName] = bcsDecoder.ReadFixedBytes(addressBytesLen)
		} else if fieldType.Bool != nil {
			jsonResult[fieldName] = bcsDecoder.Bool()
		} else if fieldType.U8 != nil {
			jsonResult[fieldName] = bcsDecoder.U8()
		} else if fieldType.U16 != nil {
			jsonResult[fieldName] = bcsDecoder.U16()
		} else if fieldType.U32 != nil {
			jsonResult[fieldName] = bcsDecoder.U32()
		} else if fieldType.U64 != nil {
			jsonResult[fieldName] = bcsDecoder.U64()
		} else if fieldType.U128 != nil {
			jsonResult[fieldName] = bcsDecoder.U128()
		} else if fieldType.U256 != nil {
			jsonResult[fieldName] = bcsDecoder.U256()
		} else if fieldType.Vector != nil {
			decodedVector, err := decodeVectorField(bcsDecoder, fieldType.Vector, normalizedStructs)
			if err != nil {
				return nil, fmt.Errorf("failed to decode vector field: %w", err)
			}
			jsonResult[fieldName] = decodedVector
		} else if fieldType.Struct != nil {
			inner, err := DecodeSuiStructToJSON(normalizedStructs, fieldType.Struct.Name, bcsDecoder)
			if err != nil {
				return nil, err
			}
			jsonResult[fieldName] = inner
		}
	}

	return jsonResult, nil
}

// func decodeVectorField(bcsDecoder *aptosBCS.Deserializer, vectorType *sui.MoveNormalizedType, _normalizedStructs map[string]*sui.MoveNormalizedStruct) (any, error) {
// 	// Read the length of the vector first
// 	vectorLength := bcsDecoder.Uleb128()

// 	// Check the inner type of the vector
// 	if vectorType.U8 != nil {
// 		// This is vector<u8> - read as bytes
// 		:= make([]byte, vectorLength)
// 		for i := range vectorLength {
// 			bytes[i] = bcsDecoder.U8()
// 		}

// 		return bytes, nil
// 	} else if vectorType.Vector != nil {
// 		vectorItem, err := decodeVectorField(bcsDecoder, vectorType.Vector, _normalizedStructs)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to decode inner vector: %w", err)
// 		}

// 		return vectorItem, nil
// 	} else if vectorType.Struct != nil {
// 		// This is vector<SomeStruct> - decode each struct
// 		structVector := make([]any, vectorLength)
// 		for range vectorLength {
// 			// For struct vectors, we'd need to calculate remaining bytes for each struct
// 			// This is more complex and would require careful BCS parsing
// 			return nil, fmt.Errorf("vector<struct> not yet implemented")
// 		}

// 		return structVector, nil
// 	}

// 	// Handle other primitive types in vectors
// 	primitiveVector := make([]any, vectorLength)
// 	for i := range vectorLength {
// 		if vectorType.Bool != nil {
// 			primitiveVector[i] = bcsDecoder.Bool()
// 		} else if vectorType.U16 != nil {
// 			primitiveVector[i] = bcsDecoder.U16()
// 		} else if vectorType.U32 != nil {
// 			primitiveVector[i] = bcsDecoder.U32()
// 		} else if vectorType.U64 != nil {
// 			primitiveVector[i] = bcsDecoder.U64()
// 		} else if vectorType.U128 != nil {
// 			primitiveVector[i] = bcsDecoder.U128()
// 		} else if vectorType.U256 != nil {
// 			primitiveVector[i] = bcsDecoder.U256()
// 		} else if vectorType.Address != nil {
// 			addressBytesLen := 32
// 			primitiveVector[i] = bcsDecoder.ReadFixedBytes(addressBytesLen)
// 		} else {
// 			return nil, fmt.Errorf("unsupported vector inner type")
// 		}
// 	}

// 	return primitiveVector, nil
// }

func decodeVectorField(bcsDecoder *aptosBCS.Deserializer, vectorType *sui.MoveNormalizedType, _normalizedStructs map[string]*sui.MoveNormalizedStruct) (any, error) {
	// Read the length of the vector first
	vectorLength := bcsDecoder.Uleb128()

	// Check the inner type of the vector
	if vectorType.U8 != nil {
		// This is vector<u8> - read as bytes
		bytes := make([]byte, vectorLength)
		for i := range vectorLength {
			bytes[i] = bcsDecoder.U8()
		}

		return bytes, nil
	} else if vectorType.Vector != nil {
		// This is vector<vector<T>> - recursively decode each inner vector
		outerVector := make([]any, vectorLength)
		for i := range vectorLength {
			innerResult, err := decodeVectorField(bcsDecoder, vectorType.Vector, _normalizedStructs)
			if err != nil {
				return nil, fmt.Errorf("failed to decode inner vector at index %d: %w", i, err)
			}
			outerVector[i] = innerResult
		}

		return outerVector, nil
	} else if vectorType.Struct != nil {
		// This is vector<SomeStruct> - decode each struct
		structVector := make([]any, vectorLength)
		for range vectorLength {
			// For struct vectors, we'd need to calculate remaining bytes for each struct
			// This is more complex and would require careful BCS parsing
			return nil, fmt.Errorf("vector<struct> not yet implemented")
		}

		return structVector, nil
	}

	primitiveVector := make([]any, vectorLength)
	for i := range vectorLength {
		if vectorType.Bool != nil {
			primitiveVector[i] = bcsDecoder.Bool()
		} else if vectorType.U16 != nil {
			primitiveVector[i] = bcsDecoder.U16()
		} else if vectorType.U32 != nil {
			primitiveVector[i] = bcsDecoder.U32()
		} else if vectorType.U64 != nil {
			primitiveVector[i] = bcsDecoder.U64()
		} else if vectorType.U128 != nil {
			primitiveVector[i] = bcsDecoder.U128()
		} else if vectorType.U256 != nil {
			primitiveVector[i] = bcsDecoder.U256()
		} else if vectorType.Address != nil {
			addressBytesLen := 32
			primitiveVector[i] = bcsDecoder.ReadFixedBytes(addressBytesLen)
		} else {
			return nil, fmt.Errorf("unsupported vector inner type")
		}
	}

	return primitiveVector, nil
}

func DecodeSuiPrimative(bcsDecoder *aptosBCS.Deserializer, primativeType string) (any, error) {
	switch primativeType {
	case "u8":
		return bcsDecoder.U8(), nil
	case "u16":
		return bcsDecoder.U16(), nil
	case "u32":
		return bcsDecoder.U32(), nil
	case "u64":
		return bcsDecoder.U64(), nil
	case "u128":
		return bcsDecoder.U128(), nil
	case "u256":
		return bcsDecoder.U256(), nil
	case "bool":
		return bcsDecoder.Bool(), nil
	//nolint:goconst
	case "address":
		addressBytesLen := 32
		return bcsDecoder.ReadFixedBytes(addressBytesLen), nil
	case "vector<address>":
		return decodeVectorField(bcsDecoder, &sui.MoveNormalizedType{Address: &sui.EmptyEnum{}}, nil)
	case "vector<u8>":
		return decodeVectorField(bcsDecoder, &sui.MoveNormalizedType{U8: &sui.EmptyEnum{}}, nil)
	case "vector<vector<u8>>":
		return decodeVectorField(bcsDecoder, &sui.MoveNormalizedType{Vector: &sui.MoveNormalizedType{U8: &sui.EmptyEnum{}}}, nil)
	}

	return nil, fmt.Errorf("unsupported BCS primitive type: %s", primativeType)
}

// decodeString handles string type decoding
func decodeString(data any, targetValue reflect.Value) error {
	str, ok := data.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", data)
	}
	targetValue.SetString(str)

	return nil
}

// decodeStruct handles struct decoding with mapstructure hooks
func decodeStruct(data any, target any) error {
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
}

// decodeGeneric handles other types via JSON marshaling/unmarshaling
func decodeGeneric(data any, target any) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	return json.Unmarshal(jsonBytes, target)
}

// decodeNumeric handles numeric types (u64, u32, etc.)
func decodeNumeric(data any, targetValue reflect.Value) error {
	//nolint:exhaustive
	switch v := data.(type) {
	case float64:
		return setNumericValue(targetValue, uint64(v))
	case string:
		n, err := strconv.ParseUint(v, base10, uint64Bits)
		if err != nil {
			return fmt.Errorf("failed to parse string as number: %w", err)
		}

		return setNumericValue(targetValue, n)
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return fmt.Errorf("failed to parse JSON number: %w", err)
		}
		if n < 0 {
			return fmt.Errorf("cannot convert negative value %d to uint", n)
		}

		return setNumericValue(targetValue, uint64(n))
	case []byte:
		return decodeNumericFromBytes(v, targetValue)
	case []any:
		bytes, err := AnySliceToBytes(v)
		if err != nil {
			return fmt.Errorf("failed to convert slice to bytes: %w", err)
		}

		return decodeNumericFromBytes(bytes, targetValue)
	default:
		return fmt.Errorf("unsupported data type for numeric target: %T", data)
	}
}

// setNumericValue sets a numeric value on the target based on its kind
func setNumericValue(targetValue reflect.Value, value uint64) error {
	//nolint:exhaustive
	switch targetValue.Kind() {
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		targetValue.SetUint(value)
		return nil
	default:
		return fmt.Errorf("unsupported target type for numeric value: %s", targetValue.Type())
	}
}

// decodeNumericFromBytes converts a byte array to a numeric value (little-endian)
func decodeNumericFromBytes(bytes []byte, targetValue reflect.Value) error {
	if len(bytes) == 0 {
		return fmt.Errorf("empty byte array cannot be converted to numeric value")
	}

	var result uint64
	// Process bytes in little-endian order
	for i := 0; i < len(bytes) && i < byteSize; i++ {
		result |= uint64(bytes[i]) << (byteSize * i)
	}

	return setNumericValue(targetValue, result)
}

// decodeSlice handles slice types
func decodeSlice(data any, targetValue reflect.Value) error {
	// Handle string to []byte conversion
	if str, ok := data.(string); ok && targetValue.Type().Elem().Kind() == reflect.Uint8 {
		return decodeStringToBytes(str, targetValue)
	}

	sourceSlice, ok := data.([]any)
	if !ok {
		return fmt.Errorf("expected slice, got %T", data)
	}

	return decodeSliceElements(sourceSlice, targetValue)
}

// decodeStringToBytes converts various string formats to byte slices
func decodeStringToBytes(str string, targetValue reflect.Value) error {
	// Try numeric string first
	if num, err := strconv.ParseUint(str, base10, uint64Bits); err == nil {
		bytes := numericToBytes(num)
		targetValue.Set(reflect.ValueOf(bytes))

		return nil
	}

	// Try hex decoding
	if strings.HasPrefix(str, "0x") {
		return decodeHexToBytes(str, targetValue)
	}

	// Try base64 decoding
	if bytes, err := base64.StdEncoding.DecodeString(str); err == nil {
		targetValue.Set(reflect.ValueOf(bytes))
		return nil
	}

	// Default: convert string directly to bytes
	targetValue.Set(reflect.ValueOf([]byte(str)))

	return nil
}

// numericToBytes converts a number to byte slice (little-endian)
func numericToBytes(num uint64) []byte {
	bytes := make([]byte, uint64Bits/uint8Bits)
	for i := range uint8Bits {
		bytes[i] = byte(num >> (i * uint8Bits))
	}
	// Remove trailing zeros
	for len(bytes) > 1 && bytes[len(bytes)-1] == 0 {
		bytes = bytes[:len(bytes)-1]
	}

	return bytes
}

// decodeHexToBytes decodes hex string to bytes
func decodeHexToBytes(str string, targetValue reflect.Value) error {
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

// decodeSliceElements decodes individual slice elements
func decodeSliceElements(sourceSlice []any, targetValue reflect.Value) error {
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

// AnySliceToBytes converts slice of interface{} to byte slice
func AnySliceToBytes(src []any) ([]byte, error) {
	dst := make([]byte, len(src))
	for i, v := range src {
		//nolint:exhaustive
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
		case float64:
			if x > maxByteValue {
				return nil, fmt.Errorf("element %d: float64 %f out of byte range", i, x)
			}
			dst[i] = byte(x)
		default:
			return nil, fmt.Errorf("element %d: unsupported type %T", i, v)
		}
	}

	return dst, nil
}

// ParseSuiResponseValue extracts the actual value from Sui's response format
func ParseSuiResponseValue(rawResponse any) (any, error) {
	responseArray, ok := rawResponse.([]any)
	if !ok {
		return nil, fmt.Errorf("expected Sui response to be an array, got %T", rawResponse)
	}

	if len(responseArray) < minResponseArrayLen {
		return nil, fmt.Errorf("expected Sui response array to have at least 2 elements, got %d", len(responseArray))
	}

	responseValue := responseArray[0]
	responseType, ok := responseArray[1].(string)
	if !ok {
		return nil, fmt.Errorf("expected second response element to be type string, got %T", responseArray[1])
	}

	return parseValueByType(responseValue, responseType)
}

// parseValueByType parses response value based on its Sui type
func parseValueByType(responseValue any, responseType string) (any, error) {
	//nolint:exhaustive
	switch {
	case isUintType(responseType):
		return parseUintValue(responseValue, responseType)
	case isBigUintType(responseType):
		return parseBigUintValue(responseValue)
	case responseType == "bool":
		return responseValue, nil
	case isStringType(responseType):
		return parseStringValue(responseValue)
	case isVectorType(responseType):
		return responseValue, nil
	case isTupleType(responseType):
		return parseTupleValue(responseValue, responseType)
	case isStructType(responseType):
		return parseStructValue(responseValue)
	default:
		return responseValue, nil
	}
}

// Type checking helper functions
func isUintType(t string) bool {
	return t == "u8" || t == "u16" || t == "u32" || t == "u64"
}

func isBigUintType(t string) bool {
	return t == "u128" || t == "u256"
}

func isStringType(t string) bool {
	return strings.Contains(t, "string")
}

func isVectorType(t string) bool {
	return strings.HasPrefix(t, "vector<")
}

func isTupleType(t string) bool {
	return strings.Contains(t, ",")
}

func isStructType(t string) bool {
	parts := strings.Split(t, "::")
	return len(parts) == 3 && strings.HasPrefix(parts[0], "0x")
}

// parseUintValue handles parsing of uint types
func parseUintValue(responseValue any, responseType string) (any, error) {
	byteArray, ok := responseValue.([]any)
	if !ok {
		return responseValue, nil
	}

	expectedBytes := getExpectedBytesForUintType(responseType)
	if len(byteArray) != expectedBytes {
		return nil, fmt.Errorf("expected %d bytes for %s, got %d", expectedBytes, responseType, len(byteArray))
	}

	return convertBytesToUint64(byteArray)
}

// getExpectedBytesForUintType returns expected byte length for uint types
func getExpectedBytesForUintType(responseType string) int {
	switch responseType {
	case "u8":
		return uint8Bytes
	case "u16":
		return uint16Bytes
	case "u32":
		return uint32Bytes
	case "u64":
		return uint64Bytes
	case "u128":
		return bits128Bytes
	case "u256":
		return bits256Bytes
	default:
		return uint64Bytes
	}
}

// convertBytesToUint64 converts byte array to uint64 (little-endian)
func convertBytesToUint64(byteArray []any) (uint64, error) {
	var result uint64
	for i, v := range byteArray {
		num, ok := v.(float64)
		if !ok {
			return 0, fmt.Errorf("expected byte value at index %d, got %T", i, v)
		}
		result |= uint64(byte(num)) << (i * bitShift)
	}

	return result, nil
}

// parseBigUintValue handles parsing of large uint types (u128, u256)
func parseBigUintValue(responseValue any) (any, error) {
	if byteArray, ok := responseValue.([]any); ok {
		return convertBytesToBigInt(byteArray)
	}

	// Handle direct values
	switch v := responseValue.(type) {
	case float64:
		return big.NewInt(int64(v)), nil
	case string:
		return parseBigIntFromString(v)
	default:
		return responseValue, nil
	}
}

// convertBytesToBigInt converts byte array to big.Int
func convertBytesToBigInt(byteArray []any) (*big.Int, error) {
	bytesArray, err := AnySliceToBytes(byteArray)
	if err != nil {
		return nil, err
	}

	result := new(big.Int)
	result.SetBytes(bytesArray)

	return result, nil
}

// parseBigIntFromString parses big.Int from string
func parseBigIntFromString(str string) (*big.Int, error) {
	result := new(big.Int)
	_, ok := result.SetString(str, base10)
	if !ok {
		return nil, fmt.Errorf("cannot parse string %s as big.Int", str)
	}

	return result, nil
}

// parseStringValue handles string type parsing
func parseStringValue(responseValue any) (any, error) {
	if byteArray, ok := responseValue.([]any); ok {
		return convertBytesToString(byteArray)
	}

	return responseValue, nil
}

// convertBytesToString converts byte array to string
func convertBytesToString(byteArray []any) (string, error) {
	bytes := make([]byte, len(byteArray))
	for i, v := range byteArray {
		num, ok := v.(float64)
		if !ok {
			return "", fmt.Errorf("expected byte value at index %d, got %T", i, v)
		}
		bytes[i] = byte(num)
	}

	return string(bytes), nil
}

// parseTupleValue handles tuple type parsing
func parseTupleValue(responseValue any, responseType string) (any, error) {
	tupleArray, ok := responseValue.([]any)
	if !ok {
		return responseValue, nil
	}

	types := extractTupleTypes(responseType)
	if len(tupleArray) != len(types) {
		return nil, fmt.Errorf("tuple length mismatch: expected %d elements, got %d", len(types), len(tupleArray))
	}

	return convertTupleToMap(tupleArray, types)
}

// extractTupleTypes extracts individual types from tuple type string
func extractTupleTypes(responseType string) []string {
	typeStr := strings.Trim(responseType, "()")
	return strings.Split(typeStr, ", ")
}

// convertTupleToMap converts tuple array to map with string keys
func convertTupleToMap(tupleArray []any, types []string) (map[string]any, error) {
	result := make(map[string]any)

	for i, item := range tupleArray {
		key := fmt.Sprintf("%d", i)

		if i < len(types) {
			elemType := strings.TrimSpace(types[i])
			parsedValue, err := ParseSuiResponseValue([]any{item, elemType})
			if err != nil {
				return nil, fmt.Errorf("failed to parse tuple element %d as %s: %w", i, elemType, err)
			}
			result[key] = parsedValue
		} else {
			result[key] = item
		}
	}

	return result, nil
}

// parseStructValue handles Move struct type parsing
func parseStructValue(responseValue any) (any, error) {
	byteArray, ok := responseValue.([]any)
	if !ok {
		return nil, fmt.Errorf("expected byte array for struct type, got %T", responseValue)
	}

	return convertToBcsBytes(byteArray)
}

// convertToBcsBytes converts interface slice to byte slice
func convertToBcsBytes(byteArray []any) ([]byte, error) {
	bcsBytes := make([]byte, len(byteArray))
	for i, v := range byteArray {
		if num, ok := v.(float64); ok {
			bcsBytes[i] = byte(num)
			continue
		}

		return nil, fmt.Errorf("expected float64 for BCS byte at index %d, got %T", i, v)
	}

	return bcsBytes, nil
}

// Mapstructure hook functions

// hexStringHook handles hex string conversions
func hexStringHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	// if the source type is not a string or the target type is the same as the source type, we simply return the data as is
	if f.Kind() != reflect.String || t.Kind() == f.Kind() {
		return data, nil
	}

	str, ok := data.(string)
	if !ok || !strings.HasPrefix(str, "0x") {
		return data, nil
	}

	str = strings.TrimPrefix(str, "0x")

	// Handle single-field struct case
	if t.Kind() == reflect.Struct && t.NumField() == 1 {
		return handleSingleFieldStruct(t, data, DecodeSuiJsonValue)
	}

	return processHexConversion(str, t)
}

// handleSingleFieldStruct processes structs with single fields
func handleSingleFieldStruct(t reflect.Type, data any, decodeFn func(any, any) error) (any, error) {
	field := t.Field(0)
	newStructVal := reflect.New(t).Elem()
	fieldPtr := newStructVal.Field(0).Addr().Interface()

	if err := decodeFn(data, fieldPtr); err != nil {
		return nil, fmt.Errorf("failed decoding for single-field struct %v field %s (%v): %w",
			t, field.Name, field.Type, err)
	}

	return newStructVal.Interface(), nil
}

// processHexConversion handles hex string conversion based on target type
func processHexConversion(hexStr string, t reflect.Type) (any, error) {
	//nolint:exhaustive
	switch t.Kind() {
	case reflect.String:
		return hexStr, nil
	case reflect.Slice:
		return processHexToSlice(hexStr, t)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.ParseUint(hexStr, base16, uint64Bits)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return processHexToInt(hexStr, t)
	case reflect.Ptr:
		return processHexToPointer(hexStr, t)
	case reflect.Array:
		return processHexToArray(hexStr, t)
	case reflect.Interface:
		return "0x" + hexStr, nil
	default:
		return nil, fmt.Errorf("unsupported target type for hex string conversion: %v", t.Kind())
	}
}

// processHexToSlice converts hex string to byte slice
func processHexToSlice(hexStr string, t reflect.Type) (any, error) {
	if t.Elem().Kind() != reflect.Uint8 {
		return nil, fmt.Errorf("unsupported target slice element type for hex string conversion: %v", t.Elem().Kind())
	}

	if hexStr == "" {
		return []uint8{}, nil
	}

	if len(hexStr)%2 == 1 {
		hexStr = "0" + hexStr
	}

	return hex.DecodeString(hexStr)
}

// processHexToInt converts hex string to integer types
func processHexToInt(hexStr string, t reflect.Type) (any, error) {
	val, err := strconv.ParseInt(hexStr, base16, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hex to int: %w", err)
	}

	return reflect.ValueOf(val).Convert(t).Interface(), nil
}

// processHexToPointer converts hex string for pointer types
func processHexToPointer(hexStr string, t reflect.Type) (any, error) {
	if t == reflect.TypeOf((*big.Int)(nil)) {
		bi := new(big.Int)
		bi.SetString(hexStr, base16)

		return bi, nil
	}

	return nil, fmt.Errorf("unsupported pointer type for hex conversion: %v", t)
}

// processHexToArray converts hex string to array types
func processHexToArray(hexStr string, t reflect.Type) (any, error) {
	if t.Elem().Kind() != reflect.Uint8 {
		return nil, fmt.Errorf("unsupported target array element type for hex string conversion: %v", t.Elem().Kind())
	}

	bytes, err := processHexToSlice(hexStr, reflect.SliceOf(t.Elem()))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string %q: %w", hexStr, err)
	}

	byteSlice := bytes.([]byte)
	out := make([]uint8, t.Len())
	copy(out, byteSlice)

	return out, nil
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
		return handleSingleFieldStruct(t, data, DecodeSuiJsonValue)
	}

	// Only try base64 decoding for byte slices
	if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
		if bytes, err := base64.StdEncoding.DecodeString(str); err == nil {
			return bytes, nil
		}
	}

	return data, nil
}

// numericStringHook handles numeric string conversions
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
		return handleSingleFieldStruct(t, data, DecodeSuiJsonValue)
	}

	return processNumericString(str, t)
}

// processNumericString handles numeric string conversion based on target type
func processNumericString(str string, t reflect.Type) (any, error) {
	//nolint:exhaustive
	switch t.Kind() {
	case reflect.String:
		return str, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return processStringToInt(str, t)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return processStringToUint(str, t)
	case reflect.Float32, reflect.Float64:
		return processStringToFloat(str, t)
	case reflect.Slice:
		return processStringToSlice(str, t)
	case reflect.Ptr:
		return processStringToPointer(str, t)
	default:
		return str, nil
	}
}

// processStringToInt converts string to integer types
func processStringToInt(str string, t reflect.Type) (any, error) {
	val, err := strconv.ParseInt(str, base10, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string to int: %w", err)
	}
	if overflowInt(t, val) {
		return nil, fmt.Errorf("value %d overflows %v", val, t)
	}

	return reflect.ValueOf(val).Convert(t).Interface(), nil
}

// processStringToUint converts string to unsigned integer types
func processStringToUint(str string, t reflect.Type) (any, error) {
	val, err := strconv.ParseUint(str, base10, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string to uint: %w", err)
	}
	if overflowUint(t, val) {
		return nil, fmt.Errorf("value %d overflows %v", val, t)
	}

	return reflect.ValueOf(val).Convert(t).Interface(), nil
}

// processStringToFloat converts string to float types
func processStringToFloat(str string, t reflect.Type) (any, error) {
	val, err := strconv.ParseFloat(str, uint64Bits)
	if err != nil {
		return nil, fmt.Errorf("failed to parse string to float: %w", err)
	}
	if overflowFloat(t, val) {
		return nil, fmt.Errorf("value %f overflows %v", val, t)
	}

	return reflect.ValueOf(val).Convert(t).Interface(), nil
}

// processStringToSlice handles string to byte slice conversion for numeric strings
func processStringToSlice(str string, t reflect.Type) (any, error) {
	if t.Elem().Kind() == reflect.Uint8 {
		if num, err := strconv.ParseUint(str, base10, uint64Bits); err == nil {
			return numericToBytes(num), nil
		}
	}

	return str, nil
}

// processStringToPointer handles string to pointer conversion
func processStringToPointer(str string, t reflect.Type) (any, error) {
	if t == reflect.TypeOf((*big.Int)(nil)) {
		return parseBigIntFromString(str)
	}

	return str, nil
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
		return handleSingleFieldStruct(t, data, DecodeSuiJsonValue)
	}

	return processBooleanConversion(boolValue, t)
}

// processBooleanConversion handles boolean conversion based on target type
func processBooleanConversion(boolValue bool, t reflect.Type) (any, error) {
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
		return handleSingleFieldStruct(t, data, DecodeSuiJsonValue)
	}

	if t.Kind() != reflect.Slice {
		return data, nil
	}

	return processArrayConversion(data, t)
}

// processArrayConversion handles array to slice conversion
func processArrayConversion(data any, t reflect.Type) (any, error) {
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

// Overflow checking functions
func overflowFloat(t reflect.Type, x float64) bool {
	//nolint:exhaustive
	switch t.Kind() {
	case reflect.Float32:
		return overflowFloat32(x)
	case reflect.Float64:
		return false
	default:
		panic("reflect: OverflowFloat of non-float type " + t.String())
	}
}

func overflowFloat32(x float64) bool {
	if x < 0 {
		x = -x
	}

	return math.MaxFloat32 < x && x <= math.MaxFloat64
}

func overflowInt(t reflect.Type, x int64) bool {
	//nolint:exhaustive
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := t.Size() * uint8Bits
		trunc := (x << (uint64Bits - bitSize)) >> (uint64Bits - bitSize)

		return x != trunc
	default:
		panic("reflect: OverflowInt of non-int type " + t.String())
	}
}

func overflowUint(t reflect.Type, x uint64) bool {
	//nolint:exhaustive
	switch t.Kind() {
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bitSize := t.Size() * uint8Bits
		trunc := (x << (uint64Bits - bitSize)) >> (uint64Bits - bitSize)

		return x != trunc
	default:
		panic("reflect: OverflowUint of non-uint type " + t.String())
	}
}
