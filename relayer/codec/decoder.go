package codec

import (
	aptosBCS "github.com/aptos-labs/aptos-go-sdk/bcs"

	"github.com/smartcontractkit/chainlink-sui/codec"
)


// Deprecated: use codec.DecodeSuiJsonValue
//
//go:fix inline
func DecodeSuiJsonValue(data any, target any) error {
	return codec.DecodeSuiJsonValue(data, target)
}

// Deprecated: use codec.ConvertBase64StringsToHex
//
//go:fix inline
func ConvertBase64StringsToHex(data any) any {
	return codec.ConvertBase64StringsToHex(data)
}

// Deprecated: use codec.DecodeSuiStructToJSON
//
//go:fix inline
func DecodeSuiStructToJSON(normalizedStructs map[string]any, identifier string, bcsDecoder *aptosBCS.Deserializer) (map[string]any, error) {
	return codec.DecodeSuiStructToJSON(normalizedStructs, identifier, bcsDecoder)
}

// Deprecated: use codec.DecodeSuiPrimative
//
//go:fix inline
func DecodeSuiPrimative(bcsDecoder *aptosBCS.Deserializer, primativeType string) (any, error) {
	return codec.DecodeSuiPrimative(bcsDecoder, primativeType)
}

// Deprecated: use codec.DecodeVectorOfStructs
//
//go:fix inline
func DecodeVectorOfStructs(bcsDecoder *aptosBCS.Deserializer, vectorType string, normalizedStructs map[string]any) (any, error) {
	return codec.DecodeVectorOfStructs(bcsDecoder, vectorType, normalizedStructs)
}

// Deprecated: use codec.AnySliceToBytes
//
//go:fix inline
func AnySliceToBytes(src []any) ([]byte, error) {
	return codec.AnySliceToBytes(src)
}

// Deprecated: use codec.DeserializeExecutionReport
//
//go:fix inline
func DeserializeExecutionReport(data []byte) (*ExecutionReport, error) {
	return codec.DeserializeExecutionReport(data)
}

// Deprecated: use codec.UnwrapBCSPureBytes
//
//go:fix inline
func UnwrapBCSPureBytes(pure []byte) ([]byte, error) {
	return codec.UnwrapBCSPureBytes(pure)
}

// Deprecated: use DeserializeExecutionReportFromPure
//
//go:fix inline
func DeserializeExecutionReportFromPure(pure []byte) (*ExecutionReport, error) {
	return codec.DeserializeExecutionReportFromPure(pure)
}
