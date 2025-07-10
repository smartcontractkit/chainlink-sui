package utils

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/block-vision/sui-go-sdk/utils"
)

const (
	// SuiAddressLength is the expected length of a Sui address in bytes
	SuiAddressLength = 32

	// ObjectRefPartsCount is the expected number of parts in an object reference string
	ObjectRefPartsCount = 3
)

// ConvertAddressToString normalizes a string address to the standard format
func ConvertAddressToString(addr string) (string, error) {
	if addr == "" {
		return "", fmt.Errorf("empty address string")
	}
	normalized := string(utils.NormalizeSuiAddress(addr))

	return normalized, nil
}

// ConvertBytesToAddress converts 32-byte address to hex string
func ConvertBytesToAddress(addrBytes []byte) (string, error) {
	if len(addrBytes) != SuiAddressLength {
		return "", fmt.Errorf("invalid address bytes length: expected %d bytes, got %d", SuiAddressLength, len(addrBytes))
	}

	return "0x" + hex.EncodeToString(addrBytes), nil
}

// ConvertStringToAddressBytes converts a string address to SuiAddressBytes
func ConvertStringToAddressBytes(addr string) (*models.SuiAddressBytes, error) {
	normalized := utils.NormalizeSuiAddress(addr)
	bytes, err := transaction.ConvertSuiAddressStringToBytes(normalized)
	if err != nil {
		return nil, fmt.Errorf("failed to convert address %s: %w", addr, err)
	}

	return bytes, nil
}

// ConvertDigestToString converts digest bytes to string
func ConvertDigestToString(digest []byte) models.ObjectDigest {
	return models.ObjectDigest(EncodeBase64(digest))
}

// ConvertStringToDigestBytes converts a string digest to ObjectDigestBytes
func ConvertStringToDigestBytes(digest string) (*models.ObjectDigestBytes, error) {
	// use the SDK's built-in conversion function
	return transaction.ConvertObjectDigestStringToBytes(models.ObjectDigest(digest))
}

// ParseTypeString parses a Move type string into components
func ParseTypeString(typeStr string) (packageAddr, module, structName string, typeParams []string, err error) {
	typeStr = strings.ReplaceAll(typeStr, " ", "")

	// check for generic parameters
	genericStart := strings.Index(typeStr, "<")
	if genericStart != -1 {
		genericEnd := strings.LastIndex(typeStr, ">")
		if genericEnd == -1 || genericEnd <= genericStart {
			err = fmt.Errorf("invalid type string: %s", typeStr)
			return
		}

		// extract type parameters
		typeParamsStr := typeStr[genericStart+1 : genericEnd]
		typeParams = SplitTypeParams(typeParamsStr)

		// get base type
		typeStr = typeStr[:genericStart]
	}

	// split by ::
	parts := strings.Split(typeStr, "::")
	if len(parts) != ObjectRefPartsCount {
		err = fmt.Errorf("invalid type format, expected address::module::struct, got: %s", typeStr)
		return
	}

	packageAddr = parts[0]
	module = parts[1]
	structName = parts[2]

	return
}

func SplitTypeParams(params string) []string {
	var result []string
	var current strings.Builder
	depth := 0

	for _, ch := range params {
		switch ch {
		case '<':
			depth++
			current.WriteRune(ch)
		case '>':
			depth--
			current.WriteRune(ch)
		case ',':
			if depth == 0 {
				result = append(result, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

func ConvertTypeStringToTypeTag(typeStr string) (*transaction.TypeTag, error) {
	typeStr = strings.TrimSpace(typeStr)

	switch typeStr {
	case "bool":
		return &transaction.TypeTag{Bool: &[]bool{true}[0]}, nil
	case "u8":
		return &transaction.TypeTag{U8: &[]bool{true}[0]}, nil
	case "u16":
		return &transaction.TypeTag{U16: &[]bool{true}[0]}, nil
	case "u32":
		return &transaction.TypeTag{U32: &[]bool{true}[0]}, nil
	case "u64":
		return &transaction.TypeTag{U64: &[]bool{true}[0]}, nil
	case "u128":
		return &transaction.TypeTag{U128: &[]bool{true}[0]}, nil
	case "u256":
		return &transaction.TypeTag{U256: &[]bool{true}[0]}, nil
	case "address":
		return &transaction.TypeTag{Address: &[]bool{true}[0]}, nil
	}

	// check for vector
	if strings.HasPrefix(typeStr, "vector<") && strings.HasSuffix(typeStr, ">") {
		innerType := typeStr[7 : len(typeStr)-1]
		innerTag, err := ConvertTypeStringToTypeTag(innerType)
		if err != nil {
			return nil, fmt.Errorf("failed to parse vector inner type: %w", err)
		}

		return &transaction.TypeTag{Vector: innerTag}, nil
	}

	// must be a struct
	packageAddr, module, structName, typeParams, err := ParseTypeString(typeStr)
	if err != nil {
		return nil, err
	}

	// convert address
	addrBytes, err := ConvertStringToAddressBytes(packageAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert package address: %w", err)
	}

	// convert type parameters
	typeParamTags := make([]*transaction.TypeTag, 0, len(typeParams))
	for _, param := range typeParams {
		tag, err := ConvertTypeStringToTypeTag(param)
		if err != nil {
			return nil, fmt.Errorf("failed to convert type parameter %s: %w", param, err)
		}
		typeParamTags = append(typeParamTags, tag)
	}

	return &transaction.TypeTag{
		Struct: &transaction.StructTag{
			Address:    *addrBytes,
			Module:     module,
			Name:       structName,
			TypeParams: typeParamTags,
		},
	}, nil
}

// EncodeBase64 encodes data to base64 string
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeBase64 decodes a base64 encoded string
func DecodeBase64(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
