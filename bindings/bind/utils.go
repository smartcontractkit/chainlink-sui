package bind

import (
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/utils"
)

// ToSuiAddress normalizes and validates a Sui address
func ToSuiAddress(address string) (string, error) {
	normalized := utils.NormalizeSuiAddress(address)
	if !utils.IsValidSuiAddress(normalized) {
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
