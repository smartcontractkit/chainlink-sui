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
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/mitchellh/mapstructure"
)

const (
	// Bit and byte constants
	byteSize   = 8
	uint8Bits  = 8
	uint64Bits = 64
	bits128    = 128
	bits256    = 256

	// Number bases
	base10 = 10
	base16 = 16
	base2  = 2

	// Response parsing constants
	maxByteValue = 255
)

// DecodeSuiJsonValue decodes Sui JSON-RPC response data into the provided target
func DecodeSuiJsonValue(data any, target any) error {
	if target == nil {
		return fmt.Errorf("target cannot be nil")
	}

	// unwrap raw JSON bytes / RawMessage
	if raw, ok := data.(json.RawMessage); ok {
		var intermediate any
		if err := json.Unmarshal(raw, &intermediate); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", err)
		}

		return DecodeSuiJsonValue(intermediate, target)
	}
	// direct type‐match optimization
	if reflect.TypeOf(data) == reflect.TypeOf(target).Elem() {
		reflect.ValueOf(target).Elem().Set(reflect.ValueOf(data))
		return nil
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()

	// handle both big.Int and *big.Int before falling into Struct logic
	bigPtrT := reflect.TypeOf((*big.Int)(nil)) // *big.Int
	bigValT := bigPtrT.Elem()                  // big.Int
	if targetType == bigValT || targetType == bigPtrT {
		// expect a JSON string
		str, ok := data.(string)
		if !ok {
			return fmt.Errorf("big.Int decode: expected string, got %T", data)
		}
		bi, success := new(big.Int).SetString(str, 10)
		if !success {
			return fmt.Errorf("big.Int decode: invalid number %q", str)
		}
		if targetType == bigValT {
			// value form: big.Int
			targetValue.Set(reflect.ValueOf(*bi))
		} else {
			// pointer form: *big.Int
			targetValue.Set(reflect.ValueOf(bi))
		}

		return nil
	}

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
func DecodeSuiStructToJSON(normalizedStructs map[string]any, identifier string, bcsDecoder *aptosBCS.Deserializer) (map[string]any, error) {
	jsonResult := make(map[string]any)

	normalizedStruct, ok := normalizedStructs[identifier].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("struct with identifier '%s' not found in normalized structs", identifier)
	}

	fields, ok := normalizedStruct["fields"].([]any)
	if !ok {
		return nil, fmt.Errorf("fields not found for struct '%s'", identifier)
	}

	for _, field := range fields {
		fieldMap, ok := field.(map[string]any)
		if !ok {
			continue
		}

		fieldName, ok := fieldMap["name"].(string)
		if !ok {
			continue
		}

		fieldType := fieldMap["type"]

		// Handle different field types based on the new format
		switch v := fieldType.(type) {
		case string:
			// Primitive types like "U64", "Bool", "Address"
			value, err := getDefaultBCSConverter().DecodePrimitive(bcsDecoder, v)
			if err != nil {
				return nil, fmt.Errorf("failed to decode primitive field %s: %w", fieldName, err)
			}
			jsonResult[fieldName] = value

		case map[string]any:
			if vectorType, exists := v["Vector"]; exists {
				// Vector type
				decodedVector, err := decodeVectorField(bcsDecoder, vectorType, normalizedStructs)
				if err != nil {
					return nil, fmt.Errorf("failed to decode vector field %s: %w", fieldName, err)
				}
				jsonResult[fieldName] = decodedVector
			} else if structType, exists := v["Struct"]; exists {
				// Struct type
				structMap, ok := structType.(map[string]any)
				if !ok {
					return nil, fmt.Errorf("invalid struct type for field %s", fieldName)
				}
				structName, ok := structMap["name"].(string)
				if !ok {
					return nil, fmt.Errorf("struct name not found for field %s", fieldName)
				}

				// Special case for String struct - it's a primitive type in Sui
				if structName == "String" {
					jsonResult[fieldName] = bcsDecoder.ReadString()
				} else {
					inner, err := DecodeSuiStructToJSON(normalizedStructs, structName, bcsDecoder)
					if err != nil {
						return nil, fmt.Errorf("failed to decode struct field %s: %w", fieldName, err)
					}
					jsonResult[fieldName] = inner
				}
			}
		}
	}

	return jsonResult, nil
}

func decodeVectorField(bcsDecoder *aptosBCS.Deserializer, vectorType any, normalizedStructs map[string]any) (any, error) {
	// Read the length of the vector first
	vectorLength := bcsDecoder.Uleb128()

	switch v := vectorType.(type) {
	case string:
		// Try to use the BCS converter for registered vector types
		if getDefaultBCSConverter().HasVectorHandler(v) {
			// Use the registered vector handler which will handle the length internally
			// We need to "rewind" by putting the length back since the handler expects to read it
			// Actually, the handler receives the length as a parameter, so we're good
			handler, _ := getDefaultBCSConverter().vectorHandlers[v]
			return handler(bcsDecoder, uint64(vectorLength))
		}

		// Fall back to generic primitive vector handling
		if getDefaultBCSConverter().HasPrimitiveHandler(v) {
			primitiveVector := make([]any, vectorLength)
			for i := range vectorLength {
				value, err := getDefaultBCSConverter().DecodePrimitive(bcsDecoder, v)
				if err != nil {
					return nil, fmt.Errorf("failed to decode primitive vector element at index %d: %w", i, err)
				}
				primitiveVector[i] = value
			}
			return primitiveVector, nil
		}

		return nil, fmt.Errorf("unsupported vector element type: %s", v)

	case map[string]any:
		if innerVectorType, exists := v["Vector"]; exists {
			// This is vector<vector<T>> - recursively decode each inner vector
			outerVector := make([]any, vectorLength)
			for i := range vectorLength {
				innerResult, err := decodeVectorField(bcsDecoder, innerVectorType, normalizedStructs)
				if err != nil {
					return nil, fmt.Errorf("failed to decode inner vector at index %d: %w", i, err)
				}
				outerVector[i] = innerResult
			}

			return outerVector, nil
		} else if structType, exists := v["Struct"]; exists {
			// This is vector<SomeStruct> - decode each struct
			structVector := make([]any, vectorLength)
			structMap, ok := structType.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("invalid struct type in vector")
			}
			structName, ok := structMap["name"].(string)
			if !ok {
				return nil, fmt.Errorf("struct name not found in vector element")
			}

			// this is a special case where strings are defined as a struct in Sui normalized module structs definition
			if structName == "String" {
				vecOfStrings := make([]any, vectorLength)
				for i := range vectorLength {
					vecOfStrings[i] = bcsDecoder.ReadString()
				}

				return vecOfStrings, nil
			}

			for i := range vectorLength {
				structResult, err := DecodeSuiStructToJSON(normalizedStructs, structName, bcsDecoder)
				if err != nil {
					return nil, fmt.Errorf("failed to decode struct at index %d: %w", i, err)
				}
				structVector[i] = structResult
			}

			return structVector, nil
		}
	}

	return nil, fmt.Errorf("unsupported vector type: %v", vectorType)
}

func DecodeSuiPrimative(bcsDecoder *aptosBCS.Deserializer, primativeType string) (any, error) {
	// Try to decode as primitive using the BCS converter registry
	converter := getDefaultBCSConverter()
	if converter.HasPrimitiveHandler(primativeType) {
		return converter.DecodePrimitive(bcsDecoder, primativeType)
	}

	// Handle vector types
	if strings.HasPrefix(primativeType, "vector<") && strings.HasSuffix(primativeType, ">") {
		innerType := strings.TrimSuffix(strings.TrimPrefix(primativeType, "vector<"), ">")

		// Handle simple vector types
		if converter.HasVectorHandler(innerType) || converter.HasPrimitiveHandler(innerType) {
			return decodeVectorField(bcsDecoder, innerType, nil)
		}

		// Handle nested vector types (e.g., vector<vector<U8>>)
		if innerType == "vector<U8>" || innerType == "vector<u8>" {
			return decodeVectorField(bcsDecoder, map[string]any{"Vector": "U8"}, nil)
		}
	}

	return nil, fmt.Errorf("unsupported BCS primitive type: %s", primativeType)
}

// DecodeVectorOfStructs decodes a vector of structs from BCS bytes
// vectorType should be in format "vector<0xpackage::module::StructName>"
func DecodeVectorOfStructs(bcsDecoder *aptosBCS.Deserializer, vectorType string, normalizedStructs map[string]any) (any, error) {
	// Check if it's actually a vector type
	if !strings.HasPrefix(vectorType, "vector<") || !strings.HasSuffix(vectorType, ">") {
		return nil, fmt.Errorf("not a vector type: %s", vectorType)
	}

	// Extract inner type
	innerType := strings.TrimSuffix(strings.TrimPrefix(vectorType, "vector<"), ">")

	// Check if inner type is a struct (has 3 parts when split by ::)
	structParts := strings.Split(innerType, "::")
	if len(structParts) != 3 {
		return nil, fmt.Errorf("inner type is not a struct: %s", innerType)
	}

	structName := structParts[2]

	// Create vector type definition compatible with decodeVectorField
	vectorTypedef := map[string]any{
		"Struct": map[string]any{
			"name": structName,
		},
	}

	return decodeVectorField(bcsDecoder, vectorTypedef, normalizedStructs)
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
			UnifiedTypeConverterHook,
			mapstructure.StringToTimeDurationHookFunc(),
		),
		Result:           target,
		WeaklyTypedInput: true,
		TagName:          "json",
		MatchName: func(mapKey, fieldName string) bool {
			mk := strings.ReplaceAll(mapKey, "_", "")
			fn := strings.ReplaceAll(fieldName, "_", "")
			ok := strings.EqualFold(mk, fn)
			return ok
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	return decoder.Decode(data)
}

// temp fix for uint64 and int64 to string when marshaling to JSON
func preprocessForJSONSafeInteger(data any) any {
	switch v := data.(type) {
	case uint64:
		return strconv.FormatUint(v, 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case []uint64:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = strconv.FormatUint(item, 10)
		}
		return result
	case []int64:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = strconv.FormatInt(item, 10)
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = preprocessForJSONSafeInteger(item)
		}
		return result
	case map[string]any:
		result := make(map[string]any, len(v))
		for key, val := range v {
			result[key] = preprocessForJSONSafeInteger(val)
		}
		return result
	default:
		return data
	}
}

// decodeGeneric handles other types via JSON marshaling/unmarshaling
func decodeGeneric(data any, target any) error {
	preprocessedData := preprocessForJSONSafeInteger(data)
	jsonBytes, err := json.Marshal(preprocessedData)
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

// handleSingleFieldStruct processes structs with single fields
// This is kept here for backward compatibility but the main implementation is in type_converters.go
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

func DeserializeExecutionReport(data []byte) (*ExecutionReport, error) {
	deserializer := aptosBCS.NewDeserializer(data)

	// 1. Read source_chain_selector (u64)
	sourceChainSelector := deserializer.U64()

	// 2. Read message header
	messageID := make([]byte, 32)
	deserializer.ReadFixedBytesInto(messageID)

	headerSourceChain := deserializer.U64()
	destChainSelector := deserializer.U64()
	sequenceNumber := deserializer.U64()
	nonce := deserializer.U64()

	if sourceChainSelector != headerSourceChain {
		return nil, fmt.Errorf("source chain selector mismatch: %d != %d", sourceChainSelector, headerSourceChain)
	}

	header := RampMessageHeader{
		MessageID:           messageID,
		SourceChainSelector: headerSourceChain,
		DestChainSelector:   destChainSelector,
		SequenceNumber:      sequenceNumber,
		Nonce:               nonce,
	}

	// 3. Read sender (vector<u8>)
	sender := deserializer.ReadBytes()

	// 4. Read data (vector<u8>)
	msgData := deserializer.ReadBytes()

	// 5. Read receiver (address)
	receiver := deserializer.ReadFixedBytes(32)

	// 6. Read gas_limit (u256)
	gasLimit := deserializer.U256()

	tokenReceiver := [32]byte{}
	deserializer.ReadFixedBytesInto(tokenReceiver[:])

	// 7. Read token_amounts vector
	tokenAmountsLen := deserializer.Uleb128()
	tokenAmounts := make([]Any2SuiTokenTransfer, tokenAmountsLen)

	for i := range tokenAmountsLen {
		sourcePoolAddr := deserializer.ReadBytes()

		destToken := deserializer.ReadFixedBytes(32)

		destGas := deserializer.U32()
		extraData := deserializer.ReadBytes()
		amount := deserializer.U256()

		tokenAmounts[i] = Any2SuiTokenTransfer{
			SourcePoolAddress: sourcePoolAddr,
			DestTokenAddress:  models.SuiAddress(hex.EncodeToString(destToken)),
			DestGasAmount:     destGas,
			ExtraData:         extraData,
			Amount:            &amount,
		}
	}

	message := Any2SuiRampMessage{
		Header:        header,
		Sender:        sender,
		Data:          msgData,
		Receiver:      models.SuiAddress(hex.EncodeToString(receiver)),
		GasLimit:      &gasLimit,
		TokenReceiver: models.SuiAddressBytes(tokenReceiver),
		TokenAmounts:  tokenAmounts,
	}

	// 8. Read offchain_token_data (vector<vector<u8>>)
	offchainDataLen := deserializer.Uleb128()
	offchainData := make([][]byte, offchainDataLen)

	for i := range offchainDataLen {
		offchainData[i] = deserializer.ReadBytes()
	}

	// 9. Read proofs (vector<vector<u8>>)
	proofsLen := deserializer.Uleb128()
	proofs := make([][]byte, proofsLen)

	for i := range proofsLen {
		proofs[i] = deserializer.ReadFixedBytes(32)
	}

	if err := deserializer.Error(); err != nil {
		return nil, fmt.Errorf("failed to deserialize execution report: %w", err)
	}

	return &ExecutionReport{
		SourceChainSelector: sourceChainSelector,
		Message:             message,
		OffchainTokenData:   offchainData,
		Proofs:              proofs,
	}, nil
}
