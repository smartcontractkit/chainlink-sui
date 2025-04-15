package testutils

import (
	"context"
	"errors"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	suiAltClient "github.com/pattonkan/sui-go/suiclient"
)

// FakeSuiClient fakes the SuiClient interface; only GetTransactionStatus is used
type FakeSuiClient struct {
	// status controls the simulated response.
	Status client.TransactionResult
}

func (c *FakeSuiClient) GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error) {
	return c.Status, nil
}
func (c *FakeSuiClient) MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error) {
	return models.TxnMetaData{}, nil
}
func (c *FakeSuiClient) SendTransaction(ctx context.Context, payload client.TransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
	return models.SuiTransactionBlockResponse{}, nil
}
func (c *FakeSuiClient) ReadObjectId(ctx context.Context, objectId string) (map[string]any, error) {
	return map[string]any{}, nil
}
func (c *FakeSuiClient) ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string) (*suiAltClient.ExecutionResultType, error) {
	return nil, errors.New("invalid value")
}
func (c *FakeSuiClient) SignAndSendTransaction(ctx context.Context, txBytes string, signerOverride *signer.SuiSigner, executionRequestType client.TransactionRequestType) (models.SuiTransactionBlockResponse, error) {
	return models.SuiTransactionBlockResponse{}, nil
}
func (c *FakeSuiClient) QueryEvents(ctx context.Context, filter models.EventFilterByMoveEventModule, limit uint64, cursor *models.EventId, descending bool) (models.PaginatedEventsResponse, error) {
	return models.PaginatedEventsResponse{}, nil
}
func (c *FakeSuiClient) WithRateLimit(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}
func (c *FakeSuiClient) GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error) {
	return []models.CoinData{}, nil
}
