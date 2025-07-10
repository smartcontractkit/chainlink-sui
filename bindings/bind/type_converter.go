package bind

import (
	"bytes"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

// Constants for magic numbers
const (
	MaxUint8Value     = 255
	MaxUint16Value    = 65535
	U128ByteSize      = 16
	U256ByteSize      = 32
	DecimalBase       = 10
	HexCharPairLength = 2
)

func convertToAddressString(value any) (string, error) {
	switch v := value.(type) {
	case string:
		// normalize address to 64 chars
		addr, err := bindutils.ConvertAddressToString(v)
		if err != nil {
			return "", fmt.Errorf("invalid address %v: %w", v, err)
		}

		return addr, nil
	case []byte:
		addr, err := bindutils.ConvertBytesToAddress(v)
		if err != nil {
			return "", fmt.Errorf("invalid address bytes: %w", err)
		}

		return addr, nil
	default:
		return "", fmt.Errorf("cannot convert %T to address", value)
	}
}

func convertToByteArray(value any) ([]uint8, error) {
	switch v := value.(type) {
	case []uint8:
		return v, nil
	case string:
		v = strings.TrimPrefix(v, "0x")
		byteSlice := []uint8{}
		for i := 0; i < len(v); i += HexCharPairLength {
			if i+1 < len(v) {
				b, err := strconv.ParseUint(v[i:i+HexCharPairLength], 16, 8)
				if err != nil {
					return nil, err
				}
				byteSlice = append(byteSlice, uint8(b))
			}
		}

		return byteSlice, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to vector<u8>", value)
	}
}

// ConvertToCallArg converts a Go value to a CallArg for use in PTB
func ConvertToCallArg(typeName string, value any) (*transaction.CallArg, error) {
	if value == nil {
		return nil, fmt.Errorf("nil value for type %s", typeName)
	}

	typeName = strings.TrimSpace(typeName)
	isMutableRef := strings.HasPrefix(typeName, "&mut ")
	isImmutableRef := strings.HasPrefix(typeName, "&") && !isMutableRef

	if obj, ok := value.(Object); ok {
		arg, err := convertObjectStructToCallArg(obj, isMutableRef)
		if err != nil {
			return nil, err
		}

		return arg, nil
	}

	baseType := typeName
	if isMutableRef {
		baseType = strings.TrimPrefix(baseType, "&mut ")
		baseType = strings.TrimSpace(baseType)
	} else if isImmutableRef {
		baseType = strings.TrimPrefix(baseType, "&")
		baseType = strings.TrimSpace(baseType)
	}

	if isMutableRef || isImmutableRef {
		// check if it's a primitive reference
		primitives := []string{"u8", "u16", "u32", "u64", "u128", "u256", "bool", "address", "vector"}
		isPrimitive := false
		for _, prim := range primitives {
			if baseType == prim || strings.HasPrefix(baseType, prim+"<") {
				isPrimitive = true
				break
			}
		}

		if !isPrimitive {
			// if this is not a primitive, then this should be an object id.
			// TODO: consider supporting [32]byte and []byte len 32 object ids
			objId, ok := value.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported object type for call arg: %T", value)
			}
			arg, err := convertObjectIdToCallArg(objId, isMutableRef)
			if err != nil {
				return nil, fmt.Errorf("failed to convert object ID string (%s) to CallArg: %w", objId, err)
			}

			return arg, nil
		}
	}

	// BCS encode pure values
	return convertPureValueToCallArg(baseType, value)
}

func convertObjectStructToCallArg(obj Object, isMutable bool) (*transaction.CallArg, error) {
	objIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(obj.Id))
	if err != nil {
		return nil, fmt.Errorf("failed to convert object ID to bytes: %w", err)
	}

	if obj.InitialSharedVersion != nil {
		return &transaction.CallArg{
			Object: &transaction.ObjectArg{
				SharedObject: &transaction.SharedObjectRef{
					ObjectId:             *objIdBytes,
					InitialSharedVersion: *obj.InitialSharedVersion,
					Mutable:              isMutable,
				},
			},
		}, nil
	}

	return &transaction.CallArg{
		UnresolvedObject: &transaction.UnresolvedObject{
			ObjectId: *objIdBytes,
		},
	}, nil
}

func convertObjectIdToCallArg(objId string, _ bool) (*transaction.CallArg, error) {
	if !strings.HasPrefix(objId, "0x") {
		return nil, fmt.Errorf("object ID should start with 0x: %s", objId)
	}
	normalizedId, err := ToSuiAddress(objId)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize object ID: %w", err)
	}
	objIdBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(normalizedId))
	if err != nil {
		return nil, fmt.Errorf("failed to convert object ID to bytes: %w", err)
	}

	return &transaction.CallArg{
		UnresolvedObject: &transaction.UnresolvedObject{
			ObjectId: *objIdBytes,
		},
	}, nil
}

func convertPureValueToCallArg(typeName string, value any) (*transaction.CallArg, error) {
	var valueToEncode any
	var err error

	switch typeName {
	case "bool":
		v, ok := value.(bool)
		if !ok {
			return nil, fmt.Errorf("expected bool, got %T", value)
		}
		valueToEncode = v

	case "u8":
		valueToEncode, err = convertToUint8(value)
	case "u16":
		valueToEncode, err = convertToUint16(value)
	case "u32":
		valueToEncode, err = convertToUint32(value)
	case "u64":
		valueToEncode, err = convertToUint64(value)
	case "u128":
		bigInt, u128Err := convertToUint128(value)
		if u128Err != nil {
			return nil, u128Err
		}

		bigIntBytes := bigInt.Bytes()
		if len(bigIntBytes) > U128ByteSize {
			return nil, fmt.Errorf("u128 value too large")
		}

		result := make([]byte, U128ByteSize)
		for i, b := range bigIntBytes {
			result[U128ByteSize-len(bigIntBytes)+i] = b
		}

		return &transaction.CallArg{
			Pure: &transaction.Pure{
				Bytes: reverseBytes(result),
			},
		}, nil

	case "u256":
		bigInt, u256Err := convertToUint256(value)
		if u256Err != nil {
			return nil, u256Err
		}
		bigIntBytes := bigInt.Bytes()
		if len(bigIntBytes) > U256ByteSize {
			return nil, fmt.Errorf("u256 value too large")
		}

		result := make([]byte, U256ByteSize)
		for i, b := range bigIntBytes {
			result[U256ByteSize-len(bigIntBytes)+i] = b
		}

		return &transaction.CallArg{
			Pure: &transaction.Pure{
				Bytes: reverseBytes(result),
			},
		}, nil

	case "address":
		addrStr, addrErr := convertToAddressString(value)
		if addrErr != nil {
			return nil, addrErr
		}
		addrBytes, addrBytesErr := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(addrStr))
		if addrBytesErr != nil {
			return nil, fmt.Errorf("failed to convert address to bytes: %w", addrBytesErr)
		}

		return &transaction.CallArg{
			Pure: &transaction.Pure{
				Bytes: addrBytes[:],
			},
		}, nil

	case "vector<u8>":
		valueToEncode, err = convertToByteArray(value)

	default:
		if !strings.HasPrefix(typeName, "vector<") {
			return nil, fmt.Errorf("unsupported type for CallArg: %s", typeName)
		}
		innerType := typeName[7 : len(typeName)-1] // Remove "vector<" and ">"
		valueToEncode, err = convertVectorToBCS(innerType, value)
	}

	if err != nil {
		return nil, err
	}

	bcsBytes, err := bcsEncode(valueToEncode)
	if err != nil {
		return nil, fmt.Errorf("failed to BCS encode value: %w", err)
	}

	return &transaction.CallArg{
		Pure: &transaction.Pure{
			Bytes: bcsBytes,
		},
	}, nil
}

func convertToUint8(value any) (uint8, error) {
	switch v := value.(type) {
	case uint8:
		return v, nil
	case int:
		if v < 0 || v > MaxUint8Value {
			return 0, fmt.Errorf("value %d out of range for u8", v)
		}

		return uint8(v), nil
	case uint64:
		if v > MaxUint8Value {
			return 0, fmt.Errorf("value %d out of range for u8", v)
		}

		return uint8(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to u8", value)
	}
}

func convertToUint16(value any) (uint16, error) {
	switch v := value.(type) {
	case uint16:
		return v, nil
	case int:
		if v < 0 || v > MaxUint16Value {
			return 0, fmt.Errorf("value %d out of range for u16", v)
		}

		return uint16(v), nil
	case uint64:
		if v > MaxUint16Value {
			return 0, fmt.Errorf("value %d out of range for u16", v)
		}

		return uint16(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to u16", value)
	}
}

func convertToUint32(value any) (uint32, error) {
	switch v := value.(type) {
	case uint32:
		return v, nil
	case int:
		if v < 0 || v > int(^uint32(0)) {
			return 0, fmt.Errorf("value %d out of range for u32", v)
		}

		return uint32(v), nil
	case uint64:
		if v > uint64(^uint32(0)) {
			return 0, fmt.Errorf("value %d out of range for u32", v)
		}

		return uint32(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to u32", value)
	}
}

func convertToUint64(value any) (uint64, error) {
	switch v := value.(type) {
	case uint64:
		return v, nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("negative value %d for u64", v)
		}

		return uint64(v), nil
	case string:
		val, err := strconv.ParseUint(v, DecimalBase, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot parse %q as u64: %w", v, err)
		}

		return val, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to u64", value)
	}
}

func convertToUint128(value any) (*big.Int, error) {
	switch v := value.(type) {
	case *big.Int:
		return v, nil
	case string:
		bi := new(big.Int)
		if _, ok := bi.SetString(v, DecimalBase); !ok {
			return nil, fmt.Errorf("invalid big int string %s", v)
		}

		return bi, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to u128", value)
	}
}

func convertToUint256(value any) (*big.Int, error) {
	switch v := value.(type) {
	case *big.Int:
		return v, nil
	case string:
		bi := new(big.Int)
		if _, ok := bi.SetString(v, DecimalBase); !ok {
			return nil, fmt.Errorf("invalid big int string %s", v)
		}

		return bi, nil
	default:
		return nil, fmt.Errorf("cannot convert %T to u256", value)
	}
}

func convertVectorToBCS(innerType string, value any) (any, error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("expected slice or array for vector, got %T", value)
	}

	switch innerType {
	case "u8":
		result := make([]uint8, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			val, err := convertToUint8(elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert vector element %d: %w", i, err)
			}
			result[i] = val
		}

		return result, nil

	case "u16":
		result := make([]uint16, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			val, err := convertToUint16(elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert vector element %d: %w", i, err)
			}
			result[i] = val
		}

		return result, nil

	case "u32":
		result := make([]uint32, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			val, err := convertToUint32(elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert vector element %d: %w", i, err)
			}
			result[i] = val
		}

		return result, nil

	case "u64":
		result := make([]uint64, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			val, err := convertToUint64(elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert vector element %d: %w", i, err)
			}
			result[i] = val
		}

		return result, nil

	case "u128":
		// For u128, we need to encode as a slice of 16-byte arrays in little-endian
		result := make([][16]byte, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			bigInt, err := convertToUint128(elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert vector element %d: %w", i, err)
			}
			bigIntBytes := bigInt.Bytes()
			if len(bigIntBytes) > U128ByteSize {
				return nil, fmt.Errorf("u128 value too large at index %d", i)
			}
			// Create 16-byte array in big-endian, then reverse for little-endian
			var arr [16]byte
			for j, b := range bigIntBytes {
				arr[16-len(bigIntBytes)+j] = b
			}
			// Reverse for little-endian
			for j := range 8 {
				arr[j], arr[15-j] = arr[15-j], arr[j]
			}
			result[i] = arr
		}

		return result, nil

	case "u256":
		// For u256, we need to encode as a slice of 32-byte arrays in little-endian
		result := make([][32]byte, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			bigInt, err := convertToUint256(elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert vector element %d: %w", i, err)
			}
			bigIntBytes := bigInt.Bytes()
			if len(bigIntBytes) > U256ByteSize {
				return nil, fmt.Errorf("u256 value too large at index %d", i)
			}
			// Create 32-byte array in big-endian, then reverse for little-endian
			var arr [32]byte
			for j, b := range bigIntBytes {
				arr[32-len(bigIntBytes)+j] = b
			}
			// Reverse for little-endian
			for j := range 16 {
				arr[j], arr[31-j] = arr[31-j], arr[j]
			}
			result[i] = arr
		}

		return result, nil

	case "bool":
		result := make([]bool, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			val, ok := elem.(bool)
			if !ok {
				return nil, fmt.Errorf("expected bool at index %d, got %T", i, elem)
			}
			result[i] = val
		}

		return result, nil

	case "address":
		result := make([][32]byte, rv.Len())
		for i := range rv.Len() {
			elem := rv.Index(i).Interface()
			addrStr, err := convertToAddressString(elem)
			if err != nil {
				return nil, fmt.Errorf("failed to convert address at index %d: %w", i, err)
			}
			addrBytes, err := transaction.ConvertSuiAddressStringToBytes(models.SuiAddress(addrStr))
			if err != nil {
				return nil, fmt.Errorf("failed to convert address to bytes at index %d: %w", i, err)
			}
			result[i] = *addrBytes
		}

		return result, nil

	default:
		if strings.HasPrefix(innerType, "vector<") {
			nextInnerType := innerType[7 : len(innerType)-1]
			result := make([]any, rv.Len())
			for i := range rv.Len() {
				elem := rv.Index(i).Interface()
				converted, err := convertVectorToBCS(nextInnerType, elem)
				if err != nil {
					return nil, fmt.Errorf("failed to convert nested vector element %d: %w", i, err)
				}
				result[i] = converted
			}

			return result, nil
		}

		return nil, fmt.Errorf("unsupported vector inner type: %s", innerType)
	}
}

type SuiAddressBytes [32]byte

func bcsEncode(value any) ([]byte, error) {
	if addrs, ok := value.([][32]byte); ok {
		suiAddrs := make([]SuiAddressBytes, len(addrs))
		for i, addr := range addrs {
			suiAddrs[i] = SuiAddressBytes(addr)
		}
		value = suiAddrs
	}

	bcsEncodedMsg := bytes.Buffer{}
	bcsEncoder := mystenbcs.NewEncoder(&bcsEncodedMsg)
	err := bcsEncoder.Encode(value)
	if err != nil {
		return nil, err
	}

	return bcsEncodedMsg.Bytes(), nil
}
