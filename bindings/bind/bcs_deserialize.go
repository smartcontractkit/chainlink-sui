package bind

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/transaction"
)

type bcsStructDecoder func([]byte) (any, error)

var bcsStructDecoders = make(map[string]bcsStructDecoder)

// RegisterStructDecoder registers a custom BCS decoder for MCMS and legacy calldata paths.
func RegisterStructDecoder(moveType string, decoder bcsStructDecoder) {
	bcsStructDecoders[moveType] = decoder
}

// DeserializeBCS decodes BCS-encoded transaction calldata (used by MCMS proposal decoding).
func DeserializeBCS(data []byte, moveTypes []string) ([]any, error) {
	reader := bytes.NewReader(data)
	deserializer := mystenbcs.NewDecoder(reader)
	ret := make([]any, 0, len(moveTypes))
	for _, moveType := range moveTypes {
		decoded, _, err := bcsDeserializeType(reader, deserializer, moveType)
		if err != nil {
			return ret, err
		}
		ret = append(ret, decoded)
	}
	if reader.Len() != 0 {
		return ret, errors.New("failed to deserialize, not all data consumed")
	}

	return ret, nil
}

func bcsDeserializeType(reader io.Reader, deserializer *mystenbcs.Decoder, moveType string) (any, reflect.Type, error) {
	switch {
	case moveType == "bool":
		var res bool
		typ, err := bcsDecode(deserializer, &res)
		return res, typ, err
	case moveType == "u8":
		var res uint8
		typ, err := bcsDecode(deserializer, &res)
		return res, typ, err
	case moveType == "u16":
		var res uint16
		typ, err := bcsDecode(deserializer, &res)
		return res, typ, err
	case moveType == "u32":
		var res uint32
		typ, err := bcsDecode(deserializer, &res)
		return res, typ, err
	case moveType == "u64":
		var res uint64
		typ, err := bcsDecode(deserializer, &res)
		return res, typ, err
	case moveType == "0x1::string::String":
		var res string
		typ, err := bcsDecode(deserializer, &res)
		return res, typ, err
	case strings.HasPrefix(moveType, "vector<") && strings.HasSuffix(moveType, ">"):
		return bcsDeserializeSlice(reader, deserializer, moveType)
	case moveType == "address":
		return bcsDeserializeAddress(deserializer)
	case moveType == "u128":
		return bcsDeserializeBigInt(deserializer, moveType, 16)
	case moveType == "u256":
		return bcsDeserializeBigInt(deserializer, moveType, 32)
	default:
		// custom move structs are decoded by their IDs
		return bcsDeserializeAddress(deserializer)
	}
}

func bcsDecode(deserializer *mystenbcs.Decoder, target any) (reflect.Type, error) {
	if _, err := deserializer.Decode(target); err != nil {
		return nil, err
	}
	return reflect.TypeOf(target).Elem(), nil
}

func bcsDeserializeSlice(reader io.Reader, deserializer *mystenbcs.Decoder, moveType string) (any, reflect.Type, error) {
	innerType := moveType[len("vector<") : len(moveType)-1]

	// vector<u8> uses ULEB128 length + raw bytes; mystenbcs handles this natively.
	if innerType == "u8" {
		var res []byte
		typ, err := bcsDecode(deserializer, &res)
		return res, typ, err
	}

	length, _, err := mystenbcs.ULEB128Decode[uint64](reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode vector length: %w", err)
	}
	if length > uint64(^uint(0)>>1) {
		return nil, nil, fmt.Errorf("vector length %d out of range", length)
	}

	// Resolve element reflect.Type from the move type string up front so length == 0
	// still yields a well-typed empty slice. Deriving the type only from decoded
	// elements would make reflect.SliceOf(nil) panic on an empty vector.
	elemType, err := reflectTypeForMoveType(innerType)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve element type %q: %w", innerType, err)
	}

	// length has been range-checked against maxInt above, so the narrowing is safe.
	n := int(length)
	sliceType := reflect.SliceOf(elemType)
	slice := reflect.MakeSlice(sliceType, n, n)
	for i := range n {
		dec, _, err := bcsDeserializeType(reader, deserializer, innerType)
		if err != nil {
			return nil, nil, err
		}
		slice.Index(i).Set(reflect.ValueOf(dec))
	}

	return slice.Interface(), sliceType, nil
}

// reflectTypeForMoveType returns the Go reflect.Type produced by bcsDeserializeType
// for the given moveType string, without decoding any bytes. Kept in lockstep with
// bcsDeserializeType so slices can be typed up front (needed for empty vectors).
func reflectTypeForMoveType(moveType string) (reflect.Type, error) {
	switch {
	case moveType == "bool":
		return reflect.TypeFor[bool](), nil
	case moveType == "u8":
		return reflect.TypeFor[uint8](), nil
	case moveType == "u16":
		return reflect.TypeFor[uint16](), nil
	case moveType == "u32":
		return reflect.TypeFor[uint32](), nil
	case moveType == "u64":
		return reflect.TypeFor[uint64](), nil
	case moveType == "0x1::string::String":
		return reflect.TypeFor[string](), nil
	case strings.HasPrefix(moveType, "vector<") && strings.HasSuffix(moveType, ">"):
		inner := moveType[len("vector<") : len(moveType)-1]
		if inner == "u8" {
			return reflect.TypeFor[[]byte](), nil
		}
		elem, err := reflectTypeForMoveType(inner)
		if err != nil {
			return nil, err
		}

		return reflect.SliceOf(elem), nil
	case moveType == "address":
		return reflect.TypeFor[models.SuiAddress](), nil
	case moveType == "u128", moveType == "u256":
		return reflect.TypeFor[*big.Int](), nil
	default:
		// Custom Move structs fall through to bcsDeserializeAddress in decoding.
		return reflect.TypeFor[models.SuiAddress](), nil
	}
}

func bcsDeserializeAddress(deserializer *mystenbcs.Decoder) (models.SuiAddress, reflect.Type, error) {
	var res [32]byte
	_, err := bcsDecode(deserializer, &res)
	if err != nil {
		return "", nil, err
	}
	addrStr := transaction.ConvertSuiAddressBytesToString(res)
	return addrStr, reflect.TypeFor[models.SuiAddress](), nil
}

func bcsDeserializeBigInt(deserializer *mystenbcs.Decoder, moveType string, size int) (*big.Int, reflect.Type, error) {
	switch size {
	case 16:
		var bytes [16]byte
		if _, err := deserializer.Decode(&bytes); err != nil {
			return nil, nil, fmt.Errorf("failed to decode %s: %w", moveType, err)
		}
		dec := new(big.Int).SetBytes(reverseBytes(bytes[:]))
		return dec, reflect.TypeFor[*big.Int](), nil
	case 32:
		var bytes [32]byte
		if _, err := deserializer.Decode(&bytes); err != nil {
			return nil, nil, fmt.Errorf("failed to decode %s: %w", moveType, err)
		}
		dec := new(big.Int).SetBytes(reverseBytes(bytes[:]))
		return dec, reflect.TypeFor[*big.Int](), nil
	default:
		return nil, nil, fmt.Errorf("unsupported big int size %d for type %s", size, moveType)
	}
}
