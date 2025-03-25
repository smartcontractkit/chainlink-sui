package bind

import (
	"context"
	"errors"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
)

type ObjectID = string

func PublishPackage(
	// TODO: Replace by a Signer common interface
	signer signer.Signer,
	client sui.ISuiAPI,
	req models.PublishRequest,
) (ObjectID, *models.SuiTransactionBlockResponse, error) {
	unsignedTx, err := client.Publish(context.Background(), req)
	if err != nil {
		return "", nil, err
	}

	// TODO: We need to be able to sign without passing the private key...
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
	tx, err := client.SuiExecuteTransactionBlock(context.Background(), *blockReq)
	if err != nil {
		return "", nil, err
	}

	objectId, err := findObjectIdFromPublishTx(tx)
	if err != nil {
		return "", nil, err
	}

	return objectId, &tx, err
}

func findObjectIdFromPublishTx(tx models.SuiTransactionBlockResponse) (string, error) {
	for _, change := range tx.ObjectChanges {
		if change.Type == "published" {
			return change.ObjectId, nil
		}
	}

	return "", errors.New("object ID not found in transaction")
}
