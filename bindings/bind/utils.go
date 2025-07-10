package bind

import (
	"fmt"

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
