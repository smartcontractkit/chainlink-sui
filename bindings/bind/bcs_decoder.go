package bind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/transaction"
)

const (
	AddressType = "address"

	// BCS data length constants
	U8Len      = 1
	U16Len     = 2
	U32Len     = 4
	U64Len     = 8
	U128Len    = 16
	U256Len    = 32
	AddressLen = 32
)

type StructDecoder func([]byte) (any, error)

var structDecoders = make(map[string]StructDecoder)

func RegisterStructDecoder(moveType string, decoder StructDecoder) {
	structDecoders[moveType] = decoder
}

type DevInspectResult struct {
	ReturnValues [][]any `json:"returnValues"`
}

func DecodeDevInspectResults(rawResults json.RawMessage, returnTypes []string, resolver *TypeResolver) ([]any, error) {
	var executionResults []DevInspectResult
	if err := json.Unmarshal(rawResults, &executionResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DevInspect results: %w", err)
	}

	if len(executionResults) == 0 {
		return nil, fmt.Errorf("no execution results found")
	}

	result := executionResults[0]

	if len(result.ReturnValues) != len(returnTypes) {
		return nil, fmt.Errorf("expected %d return values, got %d", len(returnTypes), len(result.ReturnValues))
	}

	decodedValues := make([]any, len(result.ReturnValues))
	for i, returnValue := range result.ReturnValues {
		const minReturnValueLen = 2
		if len(returnValue) < minReturnValueLen {
			return nil, fmt.Errorf("invalid return value format at index %d", i)
		}

		bcsBytes, err := extractBCSBytes(returnValue[0])
		if err != nil {
			return nil, fmt.Errorf("failed to extract BCS bytes at index %d: %w", i, err)
		}

		moveType := returnTypes[i]
		if resolver != nil {
			moveType = resolver.ResolveType(moveType)
		}

		decoded, err := decodeBCSValue(bcsBytes, moveType)
		if err != nil {
			return nil, fmt.Errorf("failed to decode value at index %d (type %s): %w", i, moveType, err)
		}

		decodedValues[i] = decoded
	}

	return decodedValues, nil
}

func extractBCSBytes(value any) ([]byte, error) {
	bcsArray, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("BCS bytes should be an array, got %T", value)
	}

	bytes := make([]byte, len(bcsArray))
	for i, b := range bcsArray {
		num, ok := b.(float64)
		if !ok {
			return nil, fmt.Errorf("BCS byte at index %d is not a number", i)
		}
		bytes[i] = byte(num)
	}

	return bytes, nil
}

func decodeBCSValue(data []byte, moveType string) (any, error) {
	if decoder, ok := structDecoders[moveType]; ok {
		return decoder(data)
	}

	switch moveType {
	case "bool":
		if len(data) != 1 {
			return nil, fmt.Errorf("invalid bool BCS data length %d", len(data))
		}
		// check that it's 0 or 1 as expected since sui-go-sdk has a more relaxed check
		// https://github.com/block-vision/sui-go-sdk/blob/5434626f683dcd308be2c418d5611072fee52484/mystenbcs/decode.go#L117
		if data[0] != 0 && data[0] != 1 {
			return nil, fmt.Errorf("invalid bool BCS value %v", data)
		}
		var result bool
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	case "u8":
		if len(data) != U8Len {
			return nil, fmt.Errorf("invalid u8 BCS data length %d", len(data))
		}
		var result uint8
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	case "u16":
		if len(data) != U16Len {
			return nil, fmt.Errorf("invalid u16 BCS data length %d", len(data))
		}
		var result uint16
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	case "u32":
		if len(data) != U32Len {
			return nil, fmt.Errorf("invalid u32 BCS data length %d", len(data))
		}
		var result uint32
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	case "u64":
		if len(data) != U64Len {
			return nil, fmt.Errorf("invalid u64 BCS data length %d", len(data))
		}
		var result uint64
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	case "u128":
		// mystenbcs.Unmarshal doesn't support u128, u256, or address handling.
		// https://github.com/block-vision/sui-go-sdk/blob/5434626f683dcd308be2c418d5611072fee52484/mystenbcs/decode.go#L131
		if len(data) != U128Len {
			return nil, fmt.Errorf("invalid u128 data length: %d", len(data))
		}
		var bytes [16]byte
		copy(bytes[:], data)
		return DecodeU128Value(bytes)

	case "u256":
		if len(data) != U256Len {
			return nil, fmt.Errorf("invalid u256 data length: %d", len(data))
		}
		var bytes [32]byte
		copy(bytes[:], data)
		return DecodeU256Value(bytes)

	case AddressType:
		if len(data) != AddressLen {
			return nil, fmt.Errorf("invalid address BCS data length %d", len(data))
		}
		var result [32]byte
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return fmt.Sprintf("0x%x", result), nil

	case "vector<u8>":
		var result []byte
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	case "vector<address>":
		var result [][32]byte
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		addresses := make([]string, len(result))
		for i, addr := range result {
			addresses[i] = fmt.Sprintf("0x%x", addr)
		}

		return addresses, nil

	case "vector<vector<address>>":
		var result [][][32]byte
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		addresses := make([][]string, len(result))
		for i, a := range result {
			subAddresses := make([]string, len(a))
			for j, addr := range a {
				subAddresses[j] = fmt.Sprintf("0x%x", addr)
			}
			addresses[i] = subAddresses
		}

		return addresses, nil

	case "vector<vector<u8>>":
		var result [][]byte
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	case "0x1::string::String":
		var result string
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	// TODO: handle vectors recursively
	case "vector<0x1::string::String>":
		var result []string
		if _, err := mystenbcs.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		return result, nil

	default:
		return data, fmt.Errorf("unsupported type for automatic decoding: %s", moveType)
	}
}

func reverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i := range data {
		result[i] = data[len(data)-1-i]
	}

	return result
}

// DecodeU256Value decodes a 32-byte array to *big.Int for u256 values
func DecodeU256Value(bcsBytes [32]byte) (*big.Int, error) {
	result := new(big.Int)
	result.SetBytes(reverseBytes(bcsBytes[:]))
	return result, nil
}

// DecodeU128Value decodes a 16-byte array to *big.Int for u128 values
func DecodeU128Value(bcsBytes [16]byte) (*big.Int, error) {
	result := new(big.Int)
	result.SetBytes(reverseBytes(bcsBytes[:]))
	return result, nil
}

// DeserializeBCS consumes a slice of bytes containing multiple BCS-encoded values
// and decodes them according to the provided Move types.
// TODO: this function should also serve extractBCSBytes, but currently
// getElementType handles objects as their ID (32-byte address)
func DeserializeBCS(data []byte, moveTypes []string) ([]any, error) {
	deserializer := mystenbcs.NewDecoder(bytes.NewReader(data))
	ret := make([]any, 0, len(moveTypes))
	for _, moveType := range moveTypes {
		decoded, _, err := decodeType(deserializer, moveType)
		if err != nil {
			return ret, err
		}
		ret = append(ret, decoded)
	}
	return ret, nil
}

func decodeType(deserializer *mystenbcs.Decoder, moveType string) (any, reflect.Type, error) {
	// Handle special cases first
	switch {
	case moveType == "bool":
		var res bool
		typ, err := decode(deserializer, &res)
		return res, typ, err
	case moveType == "u8":
		var res uint8
		typ, err := decode(deserializer, &res)
		return res, typ, err
	case moveType == "u16":
		var res uint16
		typ, err := decode(deserializer, &res)
		return res, typ, err
	case moveType == "u32":
		var res uint32
		typ, err := decode(deserializer, &res)
		return res, typ, err
	case moveType == "u64":
		var res uint64
		typ, err := decode(deserializer, &res)
		return res, typ, err
	case moveType == "0x1::string::String":
		var res string
		typ, err := decode(deserializer, &res)
		return res, typ, err
	case strings.HasPrefix(moveType, "vector"):
		// decodeSlice calls decodeType recursively
		return decodeSlice(deserializer, moveType)
	case moveType == "address":
		return decodeAddress(deserializer)
	case moveType == "u128":
		return decodeBigInt(deserializer, moveType, 16)
	case moveType == "u256":
		return decodeBigInt(deserializer, moveType, 32)
	default:
		// Custom move structs are deserialized as their ID (32-byte address)
		return decodeAddress(deserializer)
	}
}

func decodeSlice(deserializer *mystenbcs.Decoder, moveType string) (any, reflect.Type, error) {
	var length uint8

	// Decode length prefix
	deserializer.Decode(&length)
	innerType := moveType[7 : len(moveType)-1]

	elements := make([]any, length)
	var elemType reflect.Type
	for i := range elements {
		dec, refT, err := decodeType(deserializer, innerType)
		if err != nil {
			return nil, nil, err
		}
		elements[i] = dec
		elemType = refT
	}

	// Create properly typed slice using reflection
	sliceType := reflect.SliceOf(elemType)
	slice := reflect.MakeSlice(sliceType, int(length), int(length))

	// Set each element in the slice
	for i, elem := range elements {
		slice.Index(i).Set(reflect.ValueOf(elem))
	}

	return slice.Interface(), sliceType, nil
}

func decodeAddress(deserializer *mystenbcs.Decoder) (models.SuiAddress, reflect.Type, error) {
	var res [32]byte
	_, err := decode(deserializer, &res)
	if err != nil {
		return "", nil, err
	}
	addrStr := transaction.ConvertSuiAddressBytesToString(res)
	return addrStr, reflect.TypeOf(addrStr), nil
}

func decodeBigInt(deserializer *mystenbcs.Decoder, moveType string, size int) (*big.Int, reflect.Type, error) {
	// Direct decoding based on size
	switch size {
	case 16:
		var bytes [16]byte
		if _, err := deserializer.Decode(&bytes); err != nil {
			return nil, nil, fmt.Errorf("failed to decode %s: %w", moveType, err)
		}
		dec, err := DecodeU128Value(bytes)
		return dec, reflect.TypeOf(dec), err
	case 32:
		var bytes [32]byte
		if _, err := deserializer.Decode(&bytes); err != nil {
			return nil, nil, fmt.Errorf("failed to decode %s: %w", moveType, err)
		}
		dec, err := DecodeU256Value(bytes)
		return dec, reflect.TypeOf(dec), err
	default:
		return nil, nil, fmt.Errorf("unsupported big int size %d for type %s", size, moveType)
	}
}

func decode[T any](deserializer *mystenbcs.Decoder, target *T) (reflect.Type, error) {
	_, err := deserializer.Decode(target)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}
	return reflect.TypeOf(*target), nil
}
