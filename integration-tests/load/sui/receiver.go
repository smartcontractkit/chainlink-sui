package sui

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// ParseEVMReceiver32 parses a hex EVM receiver and returns a 32-byte left-padded value.
// Accepts either 20-byte (EVM address) or 32-byte already padded input.
func ParseEVMReceiver32(receiver string) ([]byte, error) {
	raw := strings.TrimPrefix(strings.TrimSpace(receiver), "0x")
	if raw == "" {
		return nil, fmt.Errorf("receiver is empty")
	}

	decoded, err := hex.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("receiver is not valid hex: %w", err)
	}

	switch len(decoded) {
	case 20:
		padded := make([]byte, 32)
		copy(padded[12:], decoded)
		if isAllZero(decoded) {
			return nil, fmt.Errorf("receiver EVM address cannot be zero")
		}
		return padded, nil
	case 32:
		if isAllZero(decoded[12:]) {
			return nil, fmt.Errorf("receiver EVM address cannot be zero")
		}
		return decoded, nil
	default:
		return nil, fmt.Errorf("receiver must decode to 20 or 32 bytes, got %d", len(decoded))
	}
}

func isAllZero(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
