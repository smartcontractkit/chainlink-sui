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

	elements := make([]any, length)
	var elemType reflect.Type
	for i := uint64(0); i < length; i++ {
		dec, refT, err := bcsDeserializeType(reader, deserializer, innerType)
		if err != nil {
			return nil, nil, err
		}
		elements[i] = dec
		elemType = refT
	}

	sliceType := reflect.SliceOf(elemType)
	slice := reflect.MakeSlice(sliceType, int(length), int(length))
	for i, elem := range elements {
		slice.Index(i).Set(reflect.ValueOf(elem))
	}

	return slice.Interface(), sliceType, nil
}

func bcsDeserializeAddress(deserializer *mystenbcs.Decoder) (models.SuiAddress, reflect.Type, error) {
	var res [32]byte
	_, err := bcsDecode(deserializer, &res)
	if err != nil {
		return "", nil, err
	}
	addrStr := transaction.ConvertSuiAddressBytesToString(res)
	return addrStr, reflect.TypeOf(addrStr), nil
}

func bcsDeserializeBigInt(deserializer *mystenbcs.Decoder, moveType string, size int) (*big.Int, reflect.Type, error) {
	switch size {
	case 16:
		var bytes [16]byte
		if _, err := deserializer.Decode(&bytes); err != nil {
			return nil, nil, fmt.Errorf("failed to decode %s: %w", moveType, err)
		}
		dec := new(big.Int).SetBytes(reverseBytes(bytes[:]))
		return dec, reflect.TypeOf(dec), nil
	case 32:
		var bytes [32]byte
		if _, err := deserializer.Decode(&bytes); err != nil {
			return nil, nil, fmt.Errorf("failed to decode %s: %w", moveType, err)
		}
		dec := new(big.Int).SetBytes(reverseBytes(bytes[:]))
		return dec, reflect.TypeOf(dec), nil
	default:
		return nil, nil, fmt.Errorf("unsupported big int size %d for type %s", size, moveType)
	}
}
