package testutils

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func PrettyPrintDebug(log logger.Logger, data any, label string) {
	resultJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Errorw("Failed to marshal data to JSON", "error", err)
	} else {
		log.Debugf("%s:\n%s", label, string(resultJSON))
	}
}

//go:fix inline
func StringPointer(s string) *string {
	return new(s)
}

//go:fix inline
func BoolPointer(b bool) *bool {
	return new(b)
}

//go:fix inline
func IntPointer(i int) *int {
	return new(i)
}

//go:fix inline
func Uint64Pointer(i uint64) *uint64 {
	return new(i)
}

//go:fix inline
func Uint32Pointer(i uint32) *uint32 {
	return new(i)
}

//go:fix inline
func Uint16Pointer(i uint16) *uint16 {
	return new(i)
}

//go:fix inline
func Uint8Pointer(i uint8) *uint8 {
	return new(i)
}
