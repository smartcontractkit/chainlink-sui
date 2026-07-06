package codec

import (
	aptosBCS "github.com/aptos-labs/aptos-go-sdk/bcs"

	"github.com/smartcontractkit/chainlink-sui/codec"
)

// Deprecated: use codec.BCSPrimitiveHandler
//
//go:fix inline
type BCSPrimitiveHandler = codec.BCSPrimitiveHandler

// Deprecated: use codec.BCSVectorHandler
//
//go:fix inline
type BCSVectorHandler = codec.BCSVectorHandler

// Deprecated: use codec.BCSTypeConverter
//
//go:fix inline
type BCSTypeConverter = codec.BCSTypeConverter

// Deprecated: use codec.NewBCSTypeConverter
//
//go:fix inline
func NewBCSTypeConverter() *BCSTypeConverter {
	return codec.NewBCSTypeConverter()
}

// Deprecated: codec.DecodeBCSPrimitive
//
//go:fix inline
func DecodeBCSPrimitive(deserializer *aptosBCS.Deserializer, primitiveType string) (any, error) {
	return codec.DecodeBCSPrimitive(deserializer, primitiveType)
}

// Deprecated: codec.DecodeBCSVector
//
//go:fix inline
func DecodeBCSVector(deserializer *aptosBCS.Deserializer, elementType string) (any, error) {
	return codec.DecodeBCSVector(deserializer, elementType)
}
