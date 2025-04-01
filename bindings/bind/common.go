package bind

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	sui_signer "github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/fardream/go-bcs/bcs"
	sui_pattokan "github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type TxOpts struct {
	GasObject string
	// Optional. GasLimit
	GasBudget *uint64
	GasPrice  *uint64
}

func SignAndSendTx(ctx context.Context, signer sui_signer.Signer, client sui.ISuiAPI, txBytes []byte) (*models.SuiTransactionBlockResponse, error) {
	relayerSigner := rel.NewPrivateKeySigner(signer.PriKey)
	signatures, err := relayerSigner.Sign(txBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign tx: %w", err)
	}

	blockReq := &models.SuiExecuteTransactionBlockRequest{
		TxBytes:   codec.EncodeBase64(txBytes),
		Signature: signatures,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: "WaitForLocalExecution",
	}

	tx, err := client.SuiExecuteTransactionBlock(ctx, *blockReq)
	if err != nil {
		msg := fmt.Errorf("tx failed calling move method: %w", err)
		return nil, msg
	}

	return &tx, nil
}

const defaultGasBudget = 200000000

func FinishTransactionFromBuilder(ctx context.Context, ptb *suiptb.ProgrammableTransactionBuilder, opts TxOpts, signer string, client sui.ISuiAPI) ([]byte, error) {
	pt := ptb.Finish()

	address, err := ToSuiAddress(signer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert signer address")
	}

	var coinData *sui_pattokan.ObjectRef
	if opts.GasObject != "" {
		coinData, err = ToSuiObjectRef(ctx, client, opts.GasObject, signer)
	} else {
		coinData, err = FetchDefaultGasCoinRef(ctx, client, signer)
	}
	if err != nil {
		return nil, err
	}

	gasBudget := uint64(defaultGasBudget)
	if opts.GasBudget != nil {
		gasBudget = *opts.GasBudget
	}
	gasPrice := suiclient.DefaultGasPrice
	if opts.GasPrice != nil {
		gasPrice = *opts.GasPrice
	}
	txData := suiptb.NewTransactionData(
		address,
		pt,
		[]*sui_pattokan.ObjectRef{coinData},
		gasBudget,
		gasPrice,
	)

	txBytes, err := bcs.Marshal(txData)
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}
