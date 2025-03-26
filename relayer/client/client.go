package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type SuiClient interface {
	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (models.SuiTransactionBlockResponse, error)
	ReadObjectId(ctx context.Context, objectId string) (any, error)
	ReadFunction(ctx context.Context, packageId string, module string, function string, args []interface{}) (models.TxnMetaData, error)
	SignAndSendTransaction(ctx context.Context, transaction models.SuiTransactionBlockData, signer *signer.SuiSigner) (models.SuiTransactionBlockResponse, error)
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
func (c *Client) ReadFunction(ctx context.Context, packageId string, module string, function string, args []interface{}, argTypes []interface{}, signer *signer.SuiSigner) (models.SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	txn, err := c.client.MoveCall(ctx, models.MoveCallRequest{
		PackageObjectId: packageId,
		Module:          module,
		Function:        function,
		TypeArguments:   argTypes,
		Arguments:       args,
		Signer:          packageId, // Using packageId as signer for read operations
	})

	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to move call: %v", err)
	}

	if signer == nil {
		// fallback to the default signer if no override is provided
		signer = c.signer
	}

	return c.SignAndSendTransaction(ctx, txn.TxBytes, signer)
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
