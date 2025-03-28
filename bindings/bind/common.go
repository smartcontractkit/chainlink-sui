package bind

import (
	"context"
	"errors"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
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

func SignAndSendTx(ctx context.Context, signer signer.Signer, client sui.ISuiAPI, txBytes []byte) (*models.SuiTransactionBlockResponse, error) {
	relayerSigner := rel.NewPrivateKeySigner(signer.PriKey)
	signatures, err := relayerSigner.Sign(txBytes)

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
		msg := fmt.Sprintf("tx failed calling move method: %v", err)
		return nil, errors.New(msg)
	}

	return &tx, nil
}

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

	gasBudget := uint64(200000000)
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
