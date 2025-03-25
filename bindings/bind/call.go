package bind

import (
	"context"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
)

type TxOpts struct {
	GasObject string
	GasBudget string
}

// TODO: We should stanardize on calling contracts. We could use the same in CW and here
func CallMethod(ctx context.Context, signer signer.Signer, client sui.ISuiAPI, req models.MoveCallRequest) (*models.SuiTransactionBlockResponse, error) {
	unsignedTx, err := client.MoveCall(ctx, req)
	if err != nil {
		msg := fmt.Sprintf("failed to generate move call transaction: %v", err)
		return nil, errors.New(msg)
	}

	signedTx := unsignedTx.SignSerializedSigWith(signer.PriKey)
	blockReq := &models.SuiExecuteTransactionBlockRequest{
		TxBytes:   signedTx.TxBytes,
		Signature: []string{signedTx.Signature},
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		// RequestType:
	}

	tx, err := client.SuiExecuteTransactionBlock(ctx, *blockReq)
	if err != nil {
		// TODO: include more details about the function and arguments
		msg := fmt.Sprintf("tx failed calling move method: %v", err)
		return nil, errors.New(msg)
	}

	return &tx, nil
}
