package utils

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// FindCurserCapObjectIDFromTx extracts the minted CurserCap object ID from an executed
// mint_and_register_curser_cap MCMS proposal transaction.
func FindCurserCapObjectIDFromTx(ctx context.Context, client sui.ISuiAPI, txDigest string) (string, error) {
	tx, err := client.SuiGetTransactionBlock(ctx, models.SuiGetTransactionBlockRequest{
		Digest: txDigest,
		Options: models.SuiTransactionBlockOptions{
			ShowObjectChanges: true,
		},
	})
	if err != nil {
		return "", fmt.Errorf("fetch transaction %s: %w", txDigest, err)
	}
	curserCapID, err := bind.FindObjectIdFromPublishTx(tx, "rmn_remote", "CurserCap")
	if err != nil {
		return "", fmt.Errorf("find CurserCap in transaction %s: %w", txDigest, err)
	}
	return curserCapID, nil
}
