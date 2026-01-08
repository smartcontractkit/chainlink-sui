package common

import (
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"slices"
	"strings"
)

func ValueAt[T any](slice []T, idx int) (T, bool) {
	var zero T
	if idx < 0 || idx >= len(slice) {
		return zero, false
	}

	return slice[idx], true
}

// InferArgumentType attempts to determine the argument type from the value
// NOTE: this method shouldn't be needed, it's a fallback for when the argument type is not known
func InferArgumentType(arg any) string {
	switch arg := arg.(type) {
	case string:
		if strings.HasPrefix(arg, "0x") {
			return "objectId"
		}

		return "address"
	case []byte:
		return "vector<u8>"
	case uint64, int64:
		return "u64"
	case int:
		return "u64"
	case int32, uint32:
		return "u32"
	case int16, uint16:
		return "u16"
	case int8, uint8:
		return "u8"
	case bool:
		return "bool"
	default:
		return "unknown"
	}
}

func SerializeUBigInt(size uint, v *big.Int) []byte {
	ub := make([]byte, size)
	v.FillBytes(ub)
	// Reverse, since big.Int outputs bytes in BigEndian
	slices.Reverse(ub)

	return ub
}

// ConvertBytesToHex recursively walks through any value and hex-encodes all []byte values.
func ConvertBytesToHex(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for k, val := range v {
			v[k] = ConvertBytesToHex(val) // recursive
		}
		return v

	case []any:
		for i, val := range v {
			v[i] = ConvertBytesToHex(val) // recursive
		}
		return v

	case []uint8:
		// Confirm it's a real []byte and not some other []uint8 misuse
		return "0x" + hex.EncodeToString(v)

	case string:
		// length prevents any random string from being encoded
		if b, err := base64.StdEncoding.DecodeString(v); err == nil && len(b) == 32 {
			return "0x" + hex.EncodeToString(b)
		}
		return v

	default:
		return value
	}
}

// Normalize a string by removing underscores and dashes and converting to lowercase
func NormalizeName(moduleName string) string {
	moduleName = strings.ReplaceAll(moduleName, "_", "")
	moduleName = strings.ReplaceAll(moduleName, "-", "")
	return strings.ToLower(moduleName)
}

func GetModuleForContract(contractName string) string {
	switch contractName {
	case "offramp", "onramp":
		return contractName
	case "Counter", "counter":
		return "counter"
	case "Router", "router":
		return "router"
	case "BurnMintTokenPool", "burn_mint_token_pool", "burnminttokenpool":
		return "burn_mint_token_pool"
	case "ManagedTokenPool", "managed_token_pool", "managedtokenpool":
		return "managed_token_pool"
	case "USDCTokenPool", "usdc_token_pool", "usdctokenpool":
		return "usdc_token_pool"
	case "LockReleaseTokenPool", "lock_release_token_pool", "lockreleasetokenpool":
		return "lock_release_token_pool"
	default:
		// anything under ccip module has to be state_object module
		return "state_object"
	}
}

// Converts snake_case to camelCase
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if i > 0 && len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(string(parts[i][0])) + parts[i][1:]
		}
	}

	return strings.Join(parts, "")
}

// Recursively convert all keys in the map to camelCase,
// with a special case for message.header.sequence_number → seqNum
func ConvertMapKeysToCamelCase(input any) any {
	return ConvertMapKeysToCamelCaseWithPath(input, "")
}

func ConvertMapKeysToCamelCaseWithPath(input any, path string) any {
	switch typed := input.(type) {
	case map[string]any:
		result := make(map[string]any)
		for k, v := range typed {
			camelKey := SnakeToCamel(k)
			fullPath := path
			if fullPath != "" {
				fullPath += "." + camelKey
			} else {
				fullPath = camelKey
			}

			if fullPath == "message.header.sequenceNumber" {
				camelKey = "seqNum"
			}

			result[camelKey] = ConvertMapKeysToCamelCaseWithPath(v, fullPath)
		}

		return result

	case []any:
		for i, v := range typed {
			typed[i] = ConvertMapKeysToCamelCaseWithPath(v, path)
		}
	}

	return input
}
