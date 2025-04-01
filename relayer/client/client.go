package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/fardream/go-bcs/bcs"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	suiAlt "github.com/pattonkan/sui-go/sui"
	suiPtb "github.com/pattonkan/sui-go/sui/suiptb"
	suiAltClient "github.com/pattonkan/sui-go/suiclient"
)

type SuiClient interface {
	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (models.SuiTransactionBlockResponse, error)
	ReadObjectId(ctx context.Context, objectId string) (map[string]interface{}, error)
	ReadFunction(ctx context.Context, packageId string, module string, function string, args []interface{}, argTypes []string) (*suiAltClient.ExecutionResultType, error)
	SignAndSendTransaction(ctx context.Context, txBytes string, signer *signer.SuiSigner) (models.SuiTransactionBlockResponse, error)
}

type Client struct {
	log                logger.Logger
	client             sui.ISuiAPI
	ptbClient          *suiAltClient.ClientImpl
	maxRetries         *int
	transactionTimeout time.Duration
	signer             *signer.SuiSigner
}

func NewClient(log logger.Logger, rpcUrl string, maxRetries *int, transactionTimeout time.Duration, signer *signer.SuiSigner) (*Client, error) {
	baseClient := sui.NewSuiClient(rpcUrl)
	ptbClient := suiAltClient.NewClient(rpcUrl)

	return &Client{
		log:                log,
		client:             baseClient,
		ptbClient:          ptbClient,
		maxRetries:         maxRetries,
		transactionTimeout: transactionTimeout,
		signer:             signer,
	}, nil
}

func (c *Client) SendTransaction(ctx context.Context, payload TransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	clientPayload := models.SuiExecuteTransactionBlockRequest{
		TxBytes:   payload.TxBytes,
		Signature: payload.Signatures,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          payload.Options.ShowInput,
			ShowRawInput:       payload.Options.ShowRawInput,
			ShowEffects:        payload.Options.ShowEffects,
			ShowEvents:         payload.Options.ShowEvents,
			ShowObjectChanges:  payload.Options.ShowObjectChanges,
			ShowBalanceChanges: payload.Options.ShowBalanceChanges,
		},
		RequestType: payload.RequestType,
	}

	return c.client.SuiExecuteTransactionBlock(ctx, clientPayload)
}

func (c *Client) MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error) {
	return c.client.MoveCall(ctx, req)
}

// ReadObjectId reads an object from the Sui blockchain
func (c *Client) ReadObjectId(ctx context.Context, objectId string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	object, err := c.client.SuiGetObject(ctx, models.SuiGetObjectRequest{
		ObjectId: objectId,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object by ID: %v", err)
	}

	return object.Data.Content.Fields, nil
}

// ReadFunction calls a Move contract function and returns the value.
// The implementation internally signs the transactions with the signer attached to the client.
// This method also calls the Move contract in "devInspect" execution mode since it is only reading values.
func (c *Client) ReadFunction(ctx context.Context, packageId string, module string, function string, args []interface{}, argTypes []string) (*suiAltClient.ExecutionResultType, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	addressHex, _ := (*c.signer).GetAddress()
	address, err := suiAlt.AddressFromHex(addressHex)
	if err != nil {
	}

	packageIdObj, err := suiAlt.PackageIdFromHex(packageId)
	if err != nil {
	}

	// Create a new client instance pointing to your Sui node's RPC endpoint.
	ptb := suiPtb.NewTransactionDataTransactionBuilder()

	// Convert each string type into a "TypeArg"
	typeTagArgs := make([]suiAlt.TypeTag, len(argTypes))
	for i, argType := range argTypes {
		typeTag, err := suiAlt.NewTypeTag(argType)
		if err != nil {
			return nil, fmt.Errorf("failed to create type tag: %v", err)
		}
		typeTagArgs[i] = *typeTag
	}

	// Convert each arg into a "CallArg" type
	callArgs := make([]suiPtb.CallArg, len(args))
	for i, arg := range args {
		encodedArg, err := codec.EncodePtbFunctionParam(argTypes[i], arg)
		if err != nil {
			return nil, fmt.Errorf("failed to encode argument: %v", err)
		}
		callArgs[i] = encodedArg
	}

	err = ptb.MoveCall(packageIdObj, module, function, []suiAlt.TypeTag{}, callArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to move call: %v", err)
	}

	pt := ptb.Finish()
	tx := suiPtb.NewTransactionData(
		address,
		pt,
		nil,
		suiAltClient.DefaultGasBudget,
		suiAltClient.DefaultGasPrice,
	)

	txBytes, err := bcs.Marshal(tx.V1.Kind)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	resp, err := c.ptbClient.DevInspectTransactionBlock(ctx, &suiAltClient.DevInspectTransactionBlockRequest{
		SenderAddress: address,
		TxKindBytes:   txBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call read function: %v", err)
	}
	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("failed to call read function: no results")
	}

	c.log.Debugw("Dev inspect results", "results", resp)

	return &resp.Results[0], nil
}

// SignAndSendTransaction given a plain (non-encoded) transaction, signs it and sends it to the node.
// The implementation uses the signer attached (default) to the client or the signer provided in the argument if specified.
func (c *Client) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signer *signer.SuiSigner) (models.SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	if signer == nil {
		// fallback to the default signer if no override is provided
		signer = c.signer
	}

	txBytes, err := base64.StdEncoding.DecodeString(txBytesRaw)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to decode tx bytes: %v", err)
	}

	signatures, err := (*signer).Sign(txBytes)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("error signing transaction: %v", err)
	}

	return c.SendTransaction(ctx, TransactionBlockRequest{
		TxBytes:    txBytesRaw,
		Signatures: signatures,
		Options: TransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: "WaitForLocalExecution",
	})
}
