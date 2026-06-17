package chainaccessor

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// suiAddressFromBytes renders raw address bytes as a 0x-prefixed hex string.
// ccipocr3 carries Sui addresses/object IDs as raw bytes; on-chain calls and the
// event store use the 0x-prefixed string form.
func suiAddressFromBytes(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// suiAddressToBytes parses a 0x-prefixed (or bare) hex Sui address into its raw
// 32-byte representation, left-padding shorter values as Sui does.
func suiAddressToBytes(addr string) ([]byte, error) {
	h := strings.TrimPrefix(strings.TrimPrefix(addr, "0x"), "0X")
	if len(h)%2 == 1 {
		h = "0" + h
	}
	b, err := hex.DecodeString(h)
	if err != nil {
		return nil, fmt.Errorf("invalid sui address %q: %w", addr, err)
	}
	if len(b) > 32 {
		return nil, fmt.Errorf("sui address %q is longer than 32 bytes", addr)
	}
	// left-pad to 32 bytes
	if len(b) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(b):], b)
		b = padded
	}
	return b, nil
}
