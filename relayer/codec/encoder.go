package codec

import "github.com/smartcontractkit/chainlink-sui/codec"

// Deprecated: use codec.EncodeToSuiValue
//
//go:fix inline
func EncodeToSuiValue(typeName string, value any) (any, error) {
	return codec.EncodeToSuiValue(typeName, value)
}
