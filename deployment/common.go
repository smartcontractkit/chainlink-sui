package deployment

import (
	"encoding/hex"
	"fmt"
	"strings"
)

func StrToBytes(s string) ([]byte, error) {
	raw, err := hex.DecodeString(strings.TrimSpace(strings.TrimPrefix(s, "0x")))
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func StrTo32(s string) ([]byte, error) {
	raw, err := hex.DecodeString(strings.TrimSpace(strings.TrimPrefix(s, "0x")))
	if err != nil {
		return nil, err
	}
	if len(raw) > 32 {
		return nil, fmt.Errorf("address longer than 32 bytes: %d", len(raw))
	}
	out := make([]byte, 32)
	copy(out[32-len(raw):], raw) // left-pad with zeros
	return out, nil
}
