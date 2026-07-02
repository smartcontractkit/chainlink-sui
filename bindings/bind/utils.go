package bind

import (
	"fmt"
	"unicode"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/codec"
)

// IsSuiAddress returns true if addr is a valid Sui address/ObjectID.
// Deprecated: use codec.IsSuiAddress directly.
func IsSuiAddress(addr string) bool {
	return codec.IsSuiAddress(addr)
}

// ToSuiAddress normalizes and validates a Sui address.
// Deprecated: use codec.ToSuiAddress directly.
func ToSuiAddress(address string) (string, error) {
	return codec.ToSuiAddress(address)
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
