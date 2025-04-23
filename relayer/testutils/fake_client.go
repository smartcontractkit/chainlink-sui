package testutils

import (
	"context"
	"errors"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	"github.com/pattonkan/sui-go/suiclient"
)

// FakeSuiPTBClient implements the SuiPTBClient interface for testing
type FakeSuiPTBClient struct {
	// Status controls the simulated response for GetTransactionStatus
	Status client.TransactionResult
}

var _ client.SuiPTBClient = (*FakeSuiPTBClient)(nil)

func (c *FakeSuiPTBClient) MoveCall(ctx context.Context, req client.MoveCallRequest) (client.TxnMetaData, error) {
	return client.TxnMetaData{}, nil
}

func (c *FakeSuiPTBClient) SendTransaction(ctx context.Context, payload client.TransactionBlockRequest) (client.SuiTransactionBlockResponse, error) {
	return client.SuiTransactionBlockResponse{}, nil
}

func (c *FakeSuiPTBClient) ReadObjectId(ctx context.Context, objectId string) (map[string]any, error) {
	return map[string]any{}, nil
}

func (c *FakeSuiPTBClient) ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string) (*suiclient.ExecutionResultType, error) {
	return nil, errors.New("invalid value")
}

func (c *FakeSuiPTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerOverride *signer.SuiSigner, executionRequestType client.TransactionRequestType) (client.SuiTransactionBlockResponse, error) {
	return client.SuiTransactionBlockResponse{}, nil
}

func (c *FakeSuiPTBClient) QueryEvents(ctx context.Context, filter client.EventFilterByMoveEventModule, limit *uint, cursor *client.EventId, sortOptions *client.QuerySortOptions) (*suiclient.EventPage, error) {
	return &suiclient.EventPage{}, nil
}

func (c *FakeSuiPTBClient) GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error) {
	return c.Status, nil
}

func (c *FakeSuiPTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]client.CoinData, error) {
	return []client.CoinData{}, nil
}

// WithRateLimit is provided to maintain compatibility with previous implementations
func (c *FakeSuiPTBClient) WithRateLimit(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

func (c *FakeSuiPTBClient) BlockByDigest(ctx context.Context, txDigest string) (*client.SuiTransactionBlockResponse, error) {
	return &client.SuiTransactionBlockResponse{}, nil
}
