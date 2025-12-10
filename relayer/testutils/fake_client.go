package testutils

import (
	"context"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/patrickmn/go-cache"

	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// FakeSuiPTBClient implements the SuiPTBClient interface for testing
type FakeSuiPTBClient struct {
	// Status controls the simulated response for GetTransactionStatus
	Status client.TransactionResult
	// CoinsData controls the simulated response for GetCoinsByAddress
	CoinsData  []models.CoinData
	MockClient sui.ISuiAPI
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

func (c *FakeSuiPTBClient) ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor string) ([]models.SuiObjectData, error) {
	return []models.SuiObjectData{}, nil
}

func (c *FakeSuiPTBClient) ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor *models.ObjectId) ([]models.SuiObjectResponse, error) {
	return []models.SuiObjectResponse{}, nil
}

func (c *FakeSuiPTBClient) ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
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
	// If CoinsData is set, return it; otherwise return default coins with sufficient balance
	if len(c.CoinsData) > 0 {
		return c.CoinsData, nil
	}

	// Return default coins with sufficient balance for gas
	return []models.CoinData{
		{
			CoinType:     "0x2::sui::SUI",
			Balance:      "100000000", // 100M units, should be enough for most gas budgets
			CoinObjectId: "0x1234567890abcdef1234567890abcdef12345678",
			Version:      "1",
			Digest:       "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT", // Valid base58 digest
		},
	}, nil
}

// WithRateLimit is provided to maintain compatibility with previous implementations
func (c *FakeSuiPTBClient) WithRateLimit(ctx context.Context, methodName string, f func(ctx context.Context) error) error {
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

func (c *FakeSuiPTBClient) GetClient() sui.ISuiAPI {
	// Return the mock client if set, otherwise nil.
	// Tests that need specific ISuiAPI behavior should set MockClient explicitly.
	// Note: Returning nil is acceptable since most tests don't use GetClient(),
	// and those that do should provide their own mock implementation.
	return c.MockClient
}

func (c *FakeSuiPTBClient) GetBlockById(ctx context.Context, checkpointId string) (models.CheckpointResponse, error) {
	return models.CheckpointResponse{}, nil
}

func (c *FakeSuiPTBClient) GetLatestEpoch(ctx context.Context) (string, error) {
	return "1000", nil
}

func (c *FakeSuiPTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *string, limit *uint64) (models.SuiXQueryTransactionBlocksResponse, error) {
	return models.SuiXQueryTransactionBlocksResponse{}, nil
}

func (c *FakeSuiPTBClient) HashTxBytes(txBytes []byte) []byte {
	return []byte{}
}

func (c *FakeSuiPTBClient) SuiXGetReferenceGasPrice(ctx context.Context) (string, error) {
	// Return a default gas price for testing
	return "1000", nil
}

func (c *FakeSuiPTBClient) GetLatestPackageId(ctx context.Context, packageId string, module string) (string, error) {
	// Return the provided package ID as the latest for testing
	return packageId, nil
}

func (c *FakeSuiPTBClient) LoadModulePackageIds(ctx context.Context, packageId string, module string) ([]string, error) {
	// Return a single package ID for testing
	return []string{packageId}, nil
}

func (c *FakeSuiPTBClient) GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (c *FakeSuiPTBClient) GetParentObjectID(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error) {
	return "", nil
}

func (c *FakeSuiPTBClient) GetCache() *cache.Cache {
	return nil
}

func (c *FakeSuiPTBClient) GetCachedValue(key string) (any, bool) {
	return nil, false
}

func (c *FakeSuiPTBClient) GetCachedValues(keys []string) (map[string]any, bool) {
	return nil, false
}

func (c *FakeSuiPTBClient) SetCachedValue(key string, value any) {
}

func (c *FakeSuiPTBClient) SetCachedValues(keyValues map[string]any) {
}

func (c *FakeSuiPTBClient) GetCCIPPackageID(ctx context.Context, offRampPackageID string, signerAddress string) (string, error) {
	return "", nil
}

func (c *FakeSuiPTBClient) GetReferenceGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1000), nil
}

func (c *FakeSuiPTBClient) QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]models.CoinData, error) {
	return []models.CoinData{}, nil
}

func (c *FakeSuiPTBClient) GetTokenPoolConfigByPackageAddress(ctx context.Context, accountAddress string, tokenPoolAddress string, ccipPackageAddress string) (module_token_admin_registry.TokenConfig, error) {
	return module_token_admin_registry.TokenConfig{}, nil
}

// StatefulFakeSuiPTBClient is a more sophisticated fake client that can change behavior
// based on gas budget and track call counts for testing gas bump scenarios
type StatefulFakeSuiPTBClient struct {
	CoinsData                    []models.CoinData
	GasBudgetThreshold           uint64 // Minimum gas budget required for success
	CallCount                    int    // Track number of calls to GetTransactionStatus
	ForcedTransactionStatusError string // Force the transaction status to be this value
	CurrentGasBudget             uint64 // Track the current gas budget being tested
	// MockClient allows tests to inject custom ISuiAPI behavior when GetClient() is called.
	// If nil, GetClient() will return nil, which is fine for most tests that don't use it.
	// Tests that need specific ISuiAPI behavior should set this field explicitly.
	MockClient sui.ISuiAPI
}

var _ client.SuiPTBClient = (*StatefulFakeSuiPTBClient)(nil)

func (c *StatefulFakeSuiPTBClient) MoveCall(ctx context.Context, req client.MoveCallRequest) (client.TxnMetaData, error) {
	return client.TxnMetaData{}, nil
}

func (c *StatefulFakeSuiPTBClient) SendTransaction(ctx context.Context, payload client.TransactionBlockRequest) (client.SuiTransactionBlockResponse, error) {
	return client.SuiTransactionBlockResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) ReadObjectId(ctx context.Context, objectId string) (models.SuiObjectData, error) {
	return models.SuiObjectData{}, nil
}

func (c *StatefulFakeSuiPTBClient) ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor string) ([]models.SuiObjectData, error) {
	return []models.SuiObjectData{}, nil
}

func (c *StatefulFakeSuiPTBClient) ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor *models.ObjectId) ([]models.SuiObjectResponse, error) {
	return []models.SuiObjectResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	return []any{}, nil
}

func (c *StatefulFakeSuiPTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte, executionRequestType client.TransactionRequestType) (client.SuiTransactionBlockResponse, error) {
	return client.SuiTransactionBlockResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) QueryEvents(ctx context.Context, filter client.EventFilterByMoveEventModule, limit *uint, cursor *client.EventId, sortOptions *client.QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	return &models.PaginatedEventsResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error) {
	c.CallCount++

	if c.ForcedTransactionStatusError != "" {
		return client.TransactionResult{
			Status: "failure",
			Error:  c.ForcedTransactionStatusError,
		}, nil
	}

	// Simulate behavior: fail with gas budget too low for first 2 attempts, then succeed
	if c.CallCount <= 2 {
		return client.TransactionResult{
			Status: "failure",
			Error:  "GasBudgetTooLow",
		}, nil
	}

	// After 2 attempts, return success
	return client.TransactionResult{
		Status: "success",
		Error:  "",
	}, nil
}

func (c *StatefulFakeSuiPTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error) {
	if len(c.CoinsData) > 0 {
		return c.CoinsData, nil
	}

	// Return default coins with sufficient balance for gas
	return []models.CoinData{
		{
			CoinType:     "0x2::sui::SUI",
			Balance:      "100000000",
			CoinObjectId: "0x1234567890abcdef1234567890abcdef12345678",
			Version:      "1",
			Digest:       "9WzSXdwbky8tNbH7juvyaui4QzMUYEjdCEKMrMgLhXHT",
		},
	}, nil
}

func (c *StatefulFakeSuiPTBClient) WithRateLimit(ctx context.Context, methodName string, f func(ctx context.Context) error) error {
	return f(ctx)
}

func (c *StatefulFakeSuiPTBClient) ToPTBArg(ctx context.Context, builder any, argValue any, isMutable bool) (any, error) {
	return argValue, nil
}

func (c *StatefulFakeSuiPTBClient) EstimateGas(ctx context.Context, txBytes string) (uint64, error) {
	return 0, nil
}

func (c *StatefulFakeSuiPTBClient) BlockByDigest(ctx context.Context, txDigest string) (*client.SuiTransactionBlockResponse, error) {
	return &client.SuiTransactionBlockResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType client.TransactionRequestType) (client.SuiTransactionBlockResponse, error) {
	return client.SuiTransactionBlockResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) GetLatestPackageId(ctx context.Context, packageId string, module string) (string, error) {
	return "", nil
}

func (c *StatefulFakeSuiPTBClient) LoadModulePackageIds(ctx context.Context, packageId string, module string) ([]string, error) {
	return []string{}, nil
}

func (c *StatefulFakeSuiPTBClient) GetSUIBalance(ctx context.Context, address string) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (c *StatefulFakeSuiPTBClient) GetNormalizedModule(ctx context.Context, packageId string, module string) (models.GetNormalizedMoveModuleResponse, error) {
	return models.GetNormalizedMoveModuleResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) GetClient() sui.ISuiAPI {
	// Return the mock client if set, otherwise nil.
	// Tests that need specific ISuiAPI behavior should set MockClient explicitly.
	// Note: Returning nil is acceptable since most tests don't use GetClient(),
	// and those that do should provide their own mock implementation.
	return c.MockClient
}

func (c *StatefulFakeSuiPTBClient) GetBlockById(ctx context.Context, checkpointId string) (models.CheckpointResponse, error) {
	return models.CheckpointResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) GetLatestEpoch(ctx context.Context) (string, error) {
	return "1000", nil
}

func (c *StatefulFakeSuiPTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *string, limit *uint64) (models.SuiXQueryTransactionBlocksResponse, error) {
	return models.SuiXQueryTransactionBlocksResponse{}, nil
}

func (c *StatefulFakeSuiPTBClient) HashTxBytes(txBytes []byte) []byte {
	return []byte{}
}

func (c *StatefulFakeSuiPTBClient) SuiXGetReferenceGasPrice(ctx context.Context) (string, error) {
	return "1000", nil
}

func (c *StatefulFakeSuiPTBClient) GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (c *StatefulFakeSuiPTBClient) GetParentObjectID(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error) {
	return "", nil
}

func (c *StatefulFakeSuiPTBClient) GetCache() *cache.Cache {
	return nil
}

func (c *StatefulFakeSuiPTBClient) GetCachedValue(key string) (any, bool) {
	return nil, false
}

func (c *StatefulFakeSuiPTBClient) GetCachedValues(keys []string) (map[string]any, bool) {
	return nil, false
}

func (c *StatefulFakeSuiPTBClient) SetCachedValue(key string, value any) {
}

func (c *StatefulFakeSuiPTBClient) SetCachedValues(keyValues map[string]any) {
}

func (c *StatefulFakeSuiPTBClient) GetCCIPPackageID(ctx context.Context, offRampPackageID string, signerAddress string) (string, error) {
	return "", nil
}

func (c *StatefulFakeSuiPTBClient) GetReferenceGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1000), nil
}

func (c *StatefulFakeSuiPTBClient) QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]models.CoinData, error) {
	return []models.CoinData{}, nil
}

func (c *StatefulFakeSuiPTBClient) GetTokenPoolConfigByPackageAddress(ctx context.Context, accountAddress string, tokenPoolAddress string, ccipPackageAddress string) (module_token_admin_registry.TokenConfig, error) {
	return module_token_admin_registry.TokenConfig{}, nil
}
