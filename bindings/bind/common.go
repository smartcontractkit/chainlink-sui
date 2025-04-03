package bind

import (
	"context"
	"fmt"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

const FailureResultType = "failure"

type TxOpts struct {
	GasObject string
	// Optional. GasLimit
	GasBudget *uint64
	GasPrice  *uint64
}

func SignAndSendTx(ctx context.Context, signer rel.SuiSigner, client suiclient.ClientImpl, txBytes []byte) (*suiclient.SuiTransactionBlockResponse, error) {
	signatures, err := signer.Sign(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	b64bytes := codec.EncodeBase64(txBytes)
	b64Tx, err := sui.NewBase64Data(b64bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert tx to base64: %w", err)
	}

	suiSignatures, err := ToSuiSignatures(signatures)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signatures to Sui format: %w", err)
	}
	blockReq := &suiclient.ExecuteTransactionBlockRequest{
		TxDataBytes: *b64Tx,
		Signatures:  suiSignatures,
		Options: &suiclient.SuiTransactionBlockResponseOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: "WaitForLocalExecution",
	}

	tx, err := client.ExecuteTransactionBlock(ctx, blockReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling move method: %w", err)
		return nil, msg
	}

	return tx, nil
}

func DevInspectTx(ctx context.Context, signer rel.SuiSigner, client suiclient.ClientImpl, txBytes []byte) (*suiclient.DevInspectTransactionBlockResponse, error) {
	_address, err := signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	address, err := ToSuiAddress(_address)
	if err != nil {
		return nil, fmt.Errorf("failed to convert address: %w", err)
	}

	b64bytes := codec.EncodeBase64(txBytes)
	b64Tx, err := sui.NewBase64Data(b64bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert tx to base64: %w", err)
	}
	blockReq := &suiclient.DevInspectTransactionBlockRequest{
		TxKindBytes:   *b64Tx,
		SenderAddress: address,
	}

	tx, err := client.DevInspectTransactionBlock(ctx, blockReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling dev inspect method: %w", err)
		return nil, msg
	}

	return tx, nil
}
