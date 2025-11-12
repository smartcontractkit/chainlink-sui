package bind

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/utils"
)

// IsSuiAddress returns true if addr is a valid Sui address/ObjectID.
// It is an improvement over the sui-go-sdk's IsValidSuiAddress function.
func IsSuiAddress(addr string) bool {
	if !(strings.HasPrefix(addr, "0x") || strings.HasPrefix(addr, "0X")) {
		return false
	}

	h := addr[2:]

	// 1..64 hex chars
	if len(h) == 0 || len(h) > 64 {
		return false
	}

	// hex only
	for _, r := range h {
		if !isHexRune(r) {
			return false
		}
	}

	// hex.DecodeString requires even length; allow odd by left-padding one '0'
	if len(h)%2 == 1 {
		h = "0" + h
	}

	b, err := hex.DecodeString(h)

	if err != nil {
		return false
	}

	// must be ≤32 bytes (32 bytes == 64 hex chars)
	return len(b) <= 32
}

// ToSuiAddress normalizes and validates a Sui address
func ToSuiAddress(address string) (string, error) {
	normalized := utils.NormalizeSuiAddress(address)
	if !IsSuiAddress(string(normalized)) {
		return "", fmt.Errorf("invalid sui address: %s", address)
	}

	return string(normalized), nil
}

func GetFailedTxError(tx *models.SuiTransactionBlockResponse) error {
	if tx.Effects.Status.Status != "failure" {
		return nil
	}

	return fmt.Errorf("transaction failed with error: %s", tx.Effects.Status.Error)
}

func isHexRune(r rune) bool {
	return unicode.IsDigit(r) ||
		('a' <= r && r <= 'f') ||
		('A' <= r && r <= 'F')
}
