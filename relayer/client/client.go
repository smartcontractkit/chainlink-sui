package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fardream/go-bcs/bcs"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
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
	ReadFunction(ctx context.Context, packageId string, module string, function string, args []interface{}, argTypes []interface{}, signer *signer.SuiSigner) (models.SuiTransactionBlockResponse, error)
	SignAndSendTransaction(ctx context.Context, txBytes string, signer *signer.SuiSigner) (models.SuiTransactionBlockResponse, error)
}

type Client struct {
	log                logger.Logger
	client             sui.ISuiAPI
	maxRetries         *int
	transactionTimeout time.Duration
	signer             *signer.SuiSigner
}

func NewClient(log logger.Logger, client sui.ISuiAPI, maxRetries *int, transactionTimeout time.Duration, signer *signer.SuiSigner) (*Client, error) {
	return &Client{
		log:                log,
		client:             client,
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
func (c *Client) ReadFunction(ctx context.Context, packageId string, module string, function string, args []interface{}, argTypes []interface{}) (models.SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	// get the account address from the signer attached to the client
	sender, err := (*c.signer).GetAddress()
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to get address: %v", err)
	}

	c.log.Debugw("Preparing to call move function", "packageId", packageId, "module", module, "function", function, "args", args, "argTypes", argTypes, "sender", sender)

	txn, err := c.client.MoveCall(ctx, models.MoveCallRequest{
		PackageObjectId: packageId,
		Module:          module,
		Function:        function,
		TypeArguments:   argTypes,
		Arguments:       args,
		Signer:          sender,
		GasBudget:       "200000",
		Gas:             nil,
		ExecutionMode:   models.TransactionExecutionDevInspect,
	})
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to move call: %v", err)
	}

	// use the default signer
	// results, err := c.SignAndSendTransaction(ctx, txn.TxBytes, nil)
	// if err != nil {
	// 	return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to dev inspect transaction: %v", err)
	// }

	results, err := c.client.SuiDevInspectTransactionBlock(ctx, models.SuiDevInspectTransactionBlockRequest{
		TxBytes: txn.TxBytes,
		Sender:  sender,
	})
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to dev inspect transaction: %v", err)
	}

	c.log.Debugw("Dev inspect results", "results", results)

	return results, nil
}

// DevInspectRequest defines the JSON-RPC request structure.
type DevInspectRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func (c *Client) DevInspectCall(ctx context.Context, packageId string, module string, function string, args []interface{}, argTypes []interface{}) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	// get the account address from the signer attached to the client
	sender, err := (*c.signer).GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %v", err)
	}

	encodedTxn, err := bcs.Marshal([]interface{}{
		sender,
		packageId,
		module,
		function,
		argTypes,
		args,
		"-", // GasObject - leaving empty
		"100000",
		"DevInspect",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal move call request: %v", err)
	}

	c.log.Debugw("Encoded transaction", "encodedTxn", encodedTxn)

	resp, err := c.client.SuiDevInspectTransactionBlock(ctx, models.SuiDevInspectTransactionBlockRequest{
		TxBytes: string(encodedTxn),
		Sender:  sender,
	})

	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to dev inspect transaction: %v", err)
	}

	return resp, nil
}

func (c *Client) DevInspectAlt(ctx context.Context, packageId string, module string, function string, args []string, argTypes []string) (any, error) {
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
	client := suiAltClient.NewClient(testutils.LocalUrl)
	ptb := suiPtb.NewTransactionDataTransactionBuilder()

	typeTagArgs := make([]suiAlt.TypeTag, len(argTypes))
	for i, argType := range argTypes {
		typeTag, err := suiAlt.NewTypeTag(argType)
		if err != nil {
			return nil, fmt.Errorf("failed to create type tag: %v", err)
		}
		typeTagArgs[i] = *typeTag
	}

	callArgs := make([]suiPtb.CallArg, len(args))
	for i, arg := range args {
		objectId, err := suiAlt.ObjectIdFromHex(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to create object ID: %v", err)
		}
		// Convert string argument to CallArg
		// Assuming the string is an object ID/address
		callArgs[i] = suiPtb.CallArg{
			Object: &suiPtb.ObjectArg{
				SharedObject: &suiPtb.SharedObjectArg{
					Id: objectId,
				},
			},
		}
	}

	c.log.Debugw("Encoded transaction", "callArgs", callArgs, "typeArgs", typeTagArgs)

	err = ptb.MoveCall(packageIdObj, module, function, []suiAlt.TypeTag{}, callArgs)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to move call: %v", err)
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
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to marshal transaction: %v", err)
	}

	resp, err := client.DevInspectTransactionBlock(ctx, &suiAltClient.DevInspectTransactionBlockRequest{
		SenderAddress: address,
		TxKindBytes:   txBytes,
	})
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to move call: %v", err)
	}

	c.log.Debugw("Dev inspect results", "results", resp)

	return resp, nil
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
