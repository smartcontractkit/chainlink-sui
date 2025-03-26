package client

import (
	"context"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

const LocalEndpoint = "http://localhost:9000"

type SuiClient interface {
	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (models.SuiTransactionBlockResponse, error)
}

type Client struct {
	log                logger.Logger
	client             sui.ISuiAPI
	maxRetries         *int
	transactionTimeout time.Duration
}

func NewClient(log logger.Logger, client sui.ISuiAPI, maxRetries *int, transactionTimeout time.Duration) (*Client, error) {
	return &Client{
		log:                log,
		client:             client,
		maxRetries:         maxRetries,
		transactionTimeout: transactionTimeout,
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
