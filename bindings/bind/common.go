package bind

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
)

type Object struct {
	Id                   string
	InitialSharedVersion *uint64
}

type CallOpts struct {
	Signer           bindutils.SuiSigner
	GasObject        string
	GasBudget        *uint64
	GasPrice         *uint64
	WaitForExecution bool

	ObjectResolver *ObjectResolver
}

func SignAndSendTx(ctx context.Context, signer bindutils.SuiSigner, client sui.ISuiAPI, txBytes []byte, waitForExecution bool) (*models.SuiTransactionBlockResponse, error) {
	signatures, err := signer.Sign(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	b64bytes := bindutils.EncodeBase64(txBytes)

	// Convert signatures to string array
	signatureStrings := make([]string, 0, len(signatures))
	signatureStrings = append(signatureStrings, signatures...)

	blockReq := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   b64bytes,
		Signature: signatureStrings,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
			ShowEvents:         true,
		},
		RequestType: "WaitForEffectsCert",
	}

	if waitForExecution {
		// TODO: this has been noted to be deprecated and removed, but still works and is used in sui-go-sdk:
		// https://forums.sui.io/t/deprecating-waitforlocalexecution/45988
		// The typescript SDK has switched to polling: https://github.com/MystenLabs/ts-sdks/blob/502ad7f2803bf6443f7cb000c802d78110585b6f/packages/typescript/src/experimental/core.ts#L114
		blockReq.RequestType = "WaitForLocalExecution"
	}

	tx, err := client.SuiExecuteTransactionBlock(ctx, blockReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling move method: %w", err)
		return nil, msg
	}

	if err := GetFailedTxError(&tx); err != nil {
		return &tx, err
	}

	return &tx, nil
}

func DevInspectTx(ctx context.Context, signerAddress string, client sui.ISuiAPI, txBytes []byte) (*models.SuiTransactionBlockResponse, error) {
	b64bytes := bindutils.EncodeBase64(txBytes)

	devInspectReq := models.SuiDevInspectTransactionBlockRequest{
		Sender:  signerAddress,
		TxBytes: b64bytes,
	}

	tx, err := client.SuiDevInspectTransactionBlock(ctx, devInspectReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling dev inspect method: %w", err)
		return nil, msg
	}

	return &tx, nil
}

// DevInspectPTB executes a PTB using DevInspect
func DevInspectPTB(ctx context.Context, signerAddress string, client sui.ISuiAPI, ptb *transaction.Transaction) (*models.SuiTransactionBlockResponse, error) {
	// ensure the PTB has the required data
	if ptb.Data.V1 == nil || ptb.Data.V1.Kind == nil {
		return nil, fmt.Errorf("PTB is not properly initialized")
	}

	// at this stage, we do not have any type information, and all unresolved variants should be resolved.
	if ptb.Data.V1.Kind.ProgrammableTransaction != nil && len(ptb.Data.V1.Kind.ProgrammableTransaction.Inputs) > 0 {
		for _, input := range ptb.Data.V1.Kind.ProgrammableTransaction.Inputs {
			if input.UnresolvedPure != nil {
				return nil, fmt.Errorf("UnresolvedPure found in PTB inputs")
			}
			if input.UnresolvedObject != nil {
				return nil, fmt.Errorf("UnresolvedObject found in PTB inputs")
			}
		}
	}

	txBytes, err := ptb.Data.V1.Kind.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction kind: %w", err)
	}

	b64TxBytes := bindutils.EncodeBase64(txBytes)

	devInspectReq := models.SuiDevInspectTransactionBlockRequest{
		Sender:  signerAddress,
		TxBytes: b64TxBytes,
	}

	tx, err := client.SuiDevInspectTransactionBlock(ctx, devInspectReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling dev inspect method: %w", err)
		return nil, msg
	}

	return &tx, nil
}
