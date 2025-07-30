package testutils

import (
	"context"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
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

func (c *FakeSuiPTBClient) ReadObjectId(ctx context.Context, objectId string) (models.SuiObjectData, error) {
	return models.SuiObjectData{}, nil
}

func (c *FakeSuiPTBClient) ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, limit *uint) ([]models.SuiObjectData, error) {
	return []models.SuiObjectData{}, nil
}

func (c *FakeSuiPTBClient) ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor *models.ObjectId) ([]models.SuiObjectResponse, error) {
	return []models.SuiObjectResponse{}, nil
}

func (c *FakeSuiPTBClient) ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string) ([]any, error) {
	return []any{}, nil
}

func (c *FakeSuiPTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte, executionRequestType client.TransactionRequestType) (client.SuiTransactionBlockResponse, error) {
	return client.SuiTransactionBlockResponse{}, nil
}

func (c *FakeSuiPTBClient) QueryEvents(ctx context.Context, filter client.EventFilterByMoveEventModule, limit *uint, cursor *client.EventId, sortOptions *client.QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	return &models.PaginatedEventsResponse{}, nil
}

func (c *FakeSuiPTBClient) GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error) {
	return c.Status, nil
}

func (c *FakeSuiPTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error) {
	return []models.CoinData{}, nil
}

// WithRateLimit is provided to maintain compatibility with previous implementations
func (c *FakeSuiPTBClient) WithRateLimit(ctx context.Context, f func(ctx context.Context) error) error {
	return f(ctx)
}

func (c *FakeSuiPTBClient) ToPTBArg(ctx context.Context, builder any, argValue any, isMutable bool) (any, error) {
	return argValue, nil
}

func (c *FakeSuiPTBClient) EstimateGas(ctx context.Context, txBytes string) (uint64, error) {
	return 0, nil
}

func (c *FakeSuiPTBClient) BlockByDigest(ctx context.Context, txDigest string) (*client.SuiTransactionBlockResponse, error) {
	return &client.SuiTransactionBlockResponse{}, nil
}

func (c *FakeSuiPTBClient) FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType client.TransactionRequestType) (client.SuiTransactionBlockResponse, error) {
	return client.SuiTransactionBlockResponse{}, nil
}

func (c *FakeSuiPTBClient) GetSUIBalance(ctx context.Context, address string) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (c *FakeSuiPTBClient) GetNormalizedModule(ctx context.Context, packageId string, module string) (models.GetNormalizedMoveModuleResponse, error) {
	return models.GetNormalizedMoveModuleResponse{}, nil
}

func (c *FakeSuiPTBClient) GetClient() *sui.ISuiAPI {
	return nil
}

func (c *FakeSuiPTBClient) GetBlockById(ctx context.Context, checkpointId string) (models.CheckpointResponse, error) {
	return models.CheckpointResponse{}, nil
}

func (c *FakeSuiPTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *string, limit *uint64) (models.SuiXQueryTransactionBlocksResponse, error) {
	return models.SuiXQueryTransactionBlocksResponse{}, nil
}
