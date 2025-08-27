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

func StringPointer(s string) *string {
	return &s
}

func BoolPointer(b bool) *bool {
	return &b
}

func IntPointer(i int) *int {
	return &i
}

func Uint64Pointer(i uint64) *uint64 {
	return &i
}

func Uint32Pointer(i uint32) *uint32 {
	return &i
}

func Uint16Pointer(i uint16) *uint16 {
	return &i
}

func Uint8Pointer(i uint8) *uint8 {
	return &i
}
