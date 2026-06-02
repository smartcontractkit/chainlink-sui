//nolint:revive // var-naming: public client API mirrors Sui RPC parameter names (objectId, packageId).
package testutils

import (
	"context"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/patrickmn/go-cache"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"

	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// FakeSuiPTBClient implements the SuiPTBClient interface for testing
type FakeSuiPTBClient struct {
	// Status controls the simulated response for GetTransactionStatus
	Status client.TransactionResult
	// CoinsData controls the simulated response for GetCoinsByAddress
	CoinsData []models.CoinData
}

var _ client.SuiPTBClient = (*FakeSuiPTBClient)(nil)

func (c *FakeSuiPTBClient) MoveCall(ctx context.Context, req client.MoveCallRequest) (client.TxnMetaData, error) {
	return client.TxnMetaData{}, nil
}

func (c *FakeSuiPTBClient) SendTransaction(ctx context.Context, execRequest *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error) {
	return nil, nil
}

func (c *FakeSuiPTBClient) ReadObjectId(ctx context.Context, objectId string) (*suirpcv2.Object, error) {
	return nil, nil
}

func (c *FakeSuiPTBClient) ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor []byte) ([]*suirpcv2.Object, error) {
	return nil, nil
}

func (c *FakeSuiPTBClient) ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor []byte) ([]*suirpcv2.Object, error) {
	return nil, nil
}

func (c *FakeSuiPTBClient) ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	return []any{}, nil
}

func (c *FakeSuiPTBClient) SimulatePTB(ctx context.Context, bcsBytes []byte) ([]any, error) {
	return []any{}, nil
}

func (c *FakeSuiPTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte) (*suirpcv2.ExecuteTransactionResponse, error) {
	return nil, nil
}

func (c *FakeSuiPTBClient) QueryEvents(ctx context.Context, filter client.EventFilterByMoveEventModule, limit *uint, cursor *client.EventId, sortOptions *client.QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	return &models.PaginatedEventsResponse{}, nil
}

func (c *FakeSuiPTBClient) GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error) {
	return c.Status, nil
}

func (c *FakeSuiPTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]*suirpcv2.Object, error) {
	return nil, nil
}

func (c *FakeSuiPTBClient) QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]*suirpcv2.Object, error) {
	return nil, nil
}

// WithRateLimit is provided to maintain compatibility with previous implementations
func (c *FakeSuiPTBClient) WithRateLimit(ctx context.Context, methodName string, f func(ctx context.Context) error) error {
	return f(ctx)
}

func (c *FakeSuiPTBClient) ToPTBArg(ctx context.Context, builder any, argValue any, isMutable bool) (any, error) {
	return argValue, nil
}

func (c *FakeSuiPTBClient) EstimateGas(ctx context.Context, tx *transaction.Transaction) (uint64, error) {
	return 0, nil
}

func (c *FakeSuiPTBClient) BlockByDigest(ctx context.Context, txDigest string) (*suirpcv2.Checkpoint, error) {
	return &suirpcv2.Checkpoint{}, nil
}

func (c *FakeSuiPTBClient) FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType client.TransactionRequestType) (*suirpcv2.ExecuteTransactionResponse, error) {
	return nil, nil
}

func (c *FakeSuiPTBClient) GetSUIBalance(ctx context.Context, address string) (*suirpcv2.Balance, error) {
	return &suirpcv2.Balance{}, nil
}

func (c *FakeSuiPTBClient) GetNormalizedModule(ctx context.Context, packageId string, module string) (models.GetNormalizedMoveModuleResponse, error) {
	return models.GetNormalizedMoveModuleResponse{}, nil
}

func (c *FakeSuiPTBClient) GetBlockById(ctx context.Context, checkpointId string) (*suirpcv2.Checkpoint, error) {
	return &suirpcv2.Checkpoint{}, nil
}

func (c *FakeSuiPTBClient) GetLatestEpoch(ctx context.Context) (*suirpcv2.Epoch, error) {
	return &suirpcv2.Epoch{}, nil
}

func (c *FakeSuiPTBClient) GetLatestCheckpoint(ctx context.Context) (*suirpcv2.Checkpoint, error) {
	return &suirpcv2.Checkpoint{}, nil
}

func (c *FakeSuiPTBClient) GetCheckpointData(ctx context.Context, checkpointSequenceNumber uint64) (*client.CheckpointData, error) {
	return &client.CheckpointData{}, nil
}

func (c *FakeSuiPTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *suirpcv2.Checkpoint, limit *uint64) ([]*suirpcv2.ExecutedTransaction, error) {
	return nil, nil
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

func (c *FakeSuiPTBClient) GetCoinMetadata(ctx context.Context, coinType string) (models.CoinMetadataResponse, error) {
	return models.CoinMetadataResponse{}, nil
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

func (c *FakeSuiPTBClient) GetCCIPPackageID(ctx context.Context, offRampPackageID string) (string, error) {
	return "", nil
}

func (c *FakeSuiPTBClient) GetReferenceGasPrice(ctx context.Context) (*big.Int, error) {
	return big.NewInt(1000), nil
}

func (c *FakeSuiPTBClient) GetTokenPoolConfigByPackageAddress(ctx context.Context, accountAddress string, tokenPoolAddress string, ccipPackageAddress string) (module_token_admin_registry.TokenConfig, error) {
	return module_token_admin_registry.TokenConfig{}, nil
}
