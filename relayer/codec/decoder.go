package codec

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	aptosBCS "github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
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

// DecodeSuiJsonValue decodes Sui JSON response data into the provided target.
func DecodeSuiJsonValue(data any, target any) error {
	return bind.DecodeJSONReturn(data, target)
}

// ConvertBase64StringsToHex walks arbitrary JSON-like structures and converts any
// base64-encoded strings into 0x-prefixed hex strings, preserving []byte slices.
func ConvertBase64StringsToHex(data any) any {
	switch v := data.(type) {
	case nil:
		return nil
	case string:
		// check if the string is entirely numeric
		if _, err := strconv.ParseUint(v, 10, 64); err == nil {
			return v
		}

		decoded, err := base64.StdEncoding.DecodeString(v)
		if err == nil && len(decoded) > 0 {
			return "0x" + hex.EncodeToString(decoded)
		}
		return v
	case json.RawMessage:
		var inner any
		if err := json.Unmarshal(v, &inner); err != nil {
			return v
		}
		return ConvertBase64StringsToHex(inner)
	case map[string]any:
		result := make(map[string]any, len(v))
		for key, value := range v {
			result[key] = ConvertBase64StringsToHex(value)
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, value := range v {
			result[i] = ConvertBase64StringsToHex(value)
		}
		return result
	default:
		rv := reflect.ValueOf(data)
		if rv.Kind() == reflect.Slice {
			// Preserve []byte and other byte slices as-is
			if rv.Type().Elem().Kind() == reflect.Uint8 {
				return data
			}

			result := make([]any, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				result[i] = ConvertBase64StringsToHex(rv.Index(i).Interface())
			}
			return result
		}

		if rv.Kind() == reflect.Map && rv.Type().Key().Kind() == reflect.String {
			result := make(map[string]any, rv.Len())
			iter := rv.MapRange()
			for iter.Next() {
				result[iter.Key().String()] = ConvertBase64StringsToHex(iter.Value().Interface())
			}
			return result
		}
	}

	return data
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

// numericToBytes converts a number to byte slice (little-endian)
// Used by type_converters.go
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

	sourceChainSelector, fields, err := deserializeExecutionReportMessageFields(deserializer)
	if err != nil {
		return nil, err
	}

	tokenAmounts, err := deserializeAny2SuiTokenTransfers(deserializer)
	if err != nil {
		return nil, err
	}

	offchainData, proofs, err := deserializeOffchainTokenDataAndProofs(deserializer)
	if err != nil {
		return nil, err
	}

	if err := finalizeExecutionReportDecode(deserializer, "execution report"); err != nil {
		return nil, err
	}

	message := Any2SuiRampMessage{
		Header:        fields.Header,
		Sender:        fields.Sender,
		Data:          fields.Data,
		Receiver:      fields.Receiver,
		GasLimit:      fields.GasLimit,
		TokenReceiver: fields.TokenReceiver,
		TokenAmounts:  tokenAmounts,
	}

	return &ExecutionReport{
		SourceChainSelector: sourceChainSelector,
		Message:             message,
		OffchainTokenData:   offchainData,
		Proofs:              proofs,
	}, nil
}

func DeserializeExecutionReportV2(data []byte) (*ExecutionReportV2, error) {
	deserializer := aptosBCS.NewDeserializer(data)

	sourceChainSelector, fields, err := deserializeExecutionReportMessageFields(deserializer)
	if err != nil {
		return nil, err
	}

	receiverObjectIDs, err := deserializeReceiverObjectIDs(deserializer)
	if err != nil {
		return nil, err
	}

	tokenAmounts, err := deserializeAny2SuiTokenTransfers(deserializer)
	if err != nil {
		return nil, err
	}

	offchainData, proofs, err := deserializeOffchainTokenDataAndProofs(deserializer)
	if err != nil {
		return nil, err
	}

	if err := finalizeExecutionReportDecode(deserializer, "V2 execution report"); err != nil {
		return nil, err
	}

	message := Any2SuiRampMessageV2{
		Header:            fields.Header,
		Sender:            fields.Sender,
		Data:              fields.Data,
		Receiver:          fields.Receiver,
		GasLimit:          fields.GasLimit,
		TokenReceiver:     fields.TokenReceiver,
		ReceiverObjectIDs: receiverObjectIDs,
		TokenAmounts:      tokenAmounts,
	}

	return &ExecutionReportV2{
		SourceChainSelector: sourceChainSelector,
		Message:             message,
		OffchainTokenData:   offchainData,
		Proofs:              proofs,
	}, nil
}

type executionReportMessageFields struct {
	Header        RampMessageHeader
	Sender        []byte
	Data          []byte
	Receiver      models.SuiAddress
	GasLimit      *big.Int
	TokenReceiver models.SuiAddressBytes
}

func deserializeExecutionReportMessageFields(deserializer *aptosBCS.Deserializer) (uint64, executionReportMessageFields, error) {
	sourceChainSelector := deserializer.U64()

	header, err := deserializeRampMessageHeader(deserializer, sourceChainSelector)
	if err != nil {
		return 0, executionReportMessageFields{}, err
	}

	sender := deserializer.ReadBytes()
	msgData := deserializer.ReadBytes()
	receiver := deserializer.ReadFixedBytes(32)
	gasLimit := deserializer.U256()

	tokenReceiver := [32]byte{}
	deserializer.ReadFixedBytesInto(tokenReceiver[:])
	if err := deserializer.Error(); err != nil {
		return 0, executionReportMessageFields{}, fmt.Errorf("failed to deserialize ramp message fields: %w", err)
	}

	return sourceChainSelector, executionReportMessageFields{
		Header:        header,
		Sender:        sender,
		Data:          msgData,
		Receiver:      models.SuiAddress(hex.EncodeToString(receiver)),
		GasLimit:      &gasLimit,
		TokenReceiver: models.SuiAddressBytes(tokenReceiver),
	}, nil
}

func deserializeRampMessageHeader(deserializer *aptosBCS.Deserializer, sourceChainSelector uint64) (RampMessageHeader, error) {
	messageID := make([]byte, 32)
	deserializer.ReadFixedBytesInto(messageID)

	headerSourceChain := deserializer.U64()
	destChainSelector := deserializer.U64()
	sequenceNumber := deserializer.U64()
	nonce := deserializer.U64()
	if err := deserializer.Error(); err != nil {
		return RampMessageHeader{}, fmt.Errorf("failed to deserialize message header: %w", err)
	}

	if sourceChainSelector != headerSourceChain {
		return RampMessageHeader{}, fmt.Errorf("source chain selector mismatch: %d != %d", sourceChainSelector, headerSourceChain)
	}

	return RampMessageHeader{
		MessageID:           messageID,
		SourceChainSelector: headerSourceChain,
		DestChainSelector:   destChainSelector,
		SequenceNumber:      sequenceNumber,
		Nonce:               nonce,
	}, nil
}

func finalizeExecutionReportDecode(deserializer *aptosBCS.Deserializer, context string) error {
	if deserializer.Remaining() > 0 {
		return fmt.Errorf("unexpected remaining bytes after decoding %s: %d", context, deserializer.Remaining())
	}
	if err := deserializer.Error(); err != nil {
		return fmt.Errorf("failed to deserialize %s: %w", context, err)
	}
	return nil
}

func deserializeReceiverObjectIDs(deserializer *aptosBCS.Deserializer) ([]models.SuiAddressBytes, error) {
	receiverObjectIDsLen := deserializer.Uleb128()
	if deserializer.Error() != nil {
		return nil, fmt.Errorf("failed to deserialize receiver_object_ids length: %w", deserializer.Error())
	}

	receiverObjectIDs := make([]models.SuiAddressBytes, receiverObjectIDsLen)
	for i := range receiverObjectIDsLen {
		var objectID [32]byte
		deserializer.ReadFixedBytesInto(objectID[:])
		if deserializer.Error() != nil {
			return nil, fmt.Errorf("failed to deserialize receiver_object_ids[%d]: %w", i, deserializer.Error())
		}
		receiverObjectIDs[i] = models.SuiAddressBytes(objectID)
	}

	return receiverObjectIDs, nil
}

func deserializeAny2SuiTokenTransfers(deserializer *aptosBCS.Deserializer) ([]Any2SuiTokenTransfer, error) {
	tokenAmountsLen := deserializer.Uleb128()
	if deserializer.Error() != nil {
		return nil, fmt.Errorf("failed to deserialize token_amounts length: %w", deserializer.Error())
	}

	tokenAmounts := make([]Any2SuiTokenTransfer, tokenAmountsLen)
	for i := range tokenAmountsLen {
		sourcePoolAddr := deserializer.ReadBytes()

		destToken := deserializer.ReadFixedBytes(32)

		destGas := deserializer.U32()
		extraData := deserializer.ReadBytes()
		amount := deserializer.U256()

		if deserializer.Error() != nil {
			return nil, fmt.Errorf("failed to deserialize token_amounts[%d]: %w", i, deserializer.Error())
		}

		tokenAmounts[i] = Any2SuiTokenTransfer{
			SourcePoolAddress: sourcePoolAddr,
			DestTokenAddress:  models.SuiAddress(hex.EncodeToString(destToken)),
			DestGasAmount:     destGas,
			ExtraData:         extraData,
			Amount:            &amount,
		}
	}

	return tokenAmounts, nil
}

func deserializeOffchainTokenDataAndProofs(deserializer *aptosBCS.Deserializer) ([][]byte, [][]byte, error) {
	offchainDataLen := deserializer.Uleb128()
	if deserializer.Error() != nil {
		return nil, nil, fmt.Errorf("failed to deserialize offchain_token_data length: %w", deserializer.Error())
	}

	offchainData := make([][]byte, offchainDataLen)
	for i := range offchainDataLen {
		offchainData[i] = deserializer.ReadBytes()
		if deserializer.Error() != nil {
			return nil, nil, fmt.Errorf("failed to deserialize offchain_token_data[%d]: %w", i, deserializer.Error())
		}
	}

	proofsLen := deserializer.Uleb128()
	if deserializer.Error() != nil {
		return nil, nil, fmt.Errorf("failed to deserialize proofs length: %w", deserializer.Error())
	}

	proofs := make([][]byte, proofsLen)
	for i := range proofsLen {
		proofs[i] = deserializer.ReadFixedBytes(32)
		if deserializer.Error() != nil {
			return nil, nil, fmt.Errorf("failed to deserialize proofs[%d]: %w", i, deserializer.Error())
		}
	}

	return offchainData, proofs, nil
}

// FormatReceiverObjectIDStrings converts report-bound receiver object IDs to 0x-prefixed hex strings.
func FormatReceiverObjectIDStrings(ids []models.SuiAddressBytes) []string {
	if len(ids) == 0 {
		return []string{}
	}

	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = "0x" + hex.EncodeToString(id[:])
	}
	return result
}

// UnwrapBCSPureBytes decodes a BCS-encoded pure input value stored on-chain.
// Pure vector<u8> arguments are stored with a ULEB128 length prefix.
func UnwrapBCSPureBytes(pure []byte) ([]byte, error) {
	if len(pure) == 0 {
		return nil, errors.New("pure bytes are empty")
	}

	deserializer := aptosBCS.NewDeserializer(pure)
	unwrapped := deserializer.ReadBytes()
	if err := deserializer.Error(); err != nil {
		return nil, fmt.Errorf("failed to unwrap pure bytes: %w", err)
	}

	return unwrapped, nil
}

// DeserializeExecutionReportFromPure deserializes an execution report from a
// BCS-encoded pure input containing vector<u8>.
func DeserializeExecutionReportFromPure(pure []byte) (*ExecutionReport, error) {
	reportBytes, err := UnwrapBCSPureBytes(pure)
	if err != nil {
		return nil, err
	}

	return DeserializeExecutionReport(reportBytes)
}
