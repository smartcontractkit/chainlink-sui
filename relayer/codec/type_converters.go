package codec

import (
	"reflect"

	"github.com/smartcontractkit/chainlink-sui/codec"
)

// Deprecated: use codec.TypeConversionFunc
//
//go:fix inline
type TypeConversionFunc = codec.TypeConversionFunc

// Deprecated: use codec.TypeConverter
//
//go:fix inline
type TypeConverter = codec.TypeConverter

// Deprecated: use codec.NewTypeConverter
//
//go:fix inline
func NewTypeConverter() *codec.TypeConverter {
	return codec.NewTypeConverter()
}

// Deprecated: use codec.UnifiedTypeConverterHook
//
//go:fix inline
func UnifiedTypeConverterHook(from, to reflect.Type, data any) (any, error) {
	return codec.UnifiedTypeConverterHook(from, to, data)
}

// Deprecated: use codec.BytesToAnySlice
//
//go:fix inline
func BytesToAnySlice(b []byte) []any {
	return codec.BytesToAnySlice(b)
}
