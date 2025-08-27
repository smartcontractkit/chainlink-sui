package bind

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/mystenbcs"
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
