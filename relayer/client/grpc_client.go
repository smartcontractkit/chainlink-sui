//nolint:revive // var-naming: public client API mirrors Sui RPC parameter names (objectId, packageId).
package client

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"

	"github.com/smartcontractkit/chainlink-sui/relayer/common"
)

const (
	maxCoinsPageSize            uint          = 50
	Base10                      int           = 10
	DefaultGasPrice             uint64        = 10_000
	DefaultGasBudget            uint64        = 1_000_000_000
	DefaultMinGasBudget         uint64        = 1_000_000
	DefaultCacheExpiration      time.Duration = 120 * time.Minute
	DefaultCacheCleanupInterval time.Duration = 240 * time.Minute
	DefaultHTTPTimeout          time.Duration = 30 * time.Second
	// DefaultReadOpTimeout caps a single read chain (transform + simulate). Prevents gRPC calls from
	// outliving the CCIP config poller's 30s deadline when the node is overloaded.
	DefaultReadOpTimeout time.Duration = 30 * time.Second
	// DefaultPackageIdCacheTTL caches resolved latest package IDs to avoid repeated heavy
	// GetFunction/GetPackage/ListOwnedObjects chains under bursty config polling.
	DefaultPackageIdCacheTTL time.Duration = 5 * time.Minute
)

var RateLimitWeights = map[string]int64{
	"MoveCall":                             0,
	"SendTransaction":                      0,
	"ReadFunction":                         0,
	"SignAndSendTransaction":               0,
	"QueryEvents":                          0,
	"QueryTransactions":                    0,
	"GetCoinsByAddress":                    0,
	"QueryCoinsByAddress":                  0,
	"EstimateGas":                          0,
	"GetTransactionStatus":                 0,
	"GetBlockById":                         0,
	"GetNormalizedModule":                  0,
	"GetSUIBalance":                        0,
	"GetCoinBalanceByAddress":              0,
	"GetValuesFromPackageOwnedObjectField": 0,
	"GetReferenceGasPrice":                 0,
	"FinishPTBAndSend":                     0,
	"ReadFilterOwnedObjectIds":             0,
	"ReadOwnedObjects":                     0,
	"ReadObjectId":                         0,
	"GetLatestPackageId":                   0,
	"LoadModulePackageIds":                 0,
	"GetParentObjectID":                    0,
	"GetCCIPPackageID":                     0,
	"GetTokenPoolConfigByPackageAddress":   0,
	"GetLatestEpoch":                       0,
	"GetTransactionsByCheckpoint":          0,
	"GetLatestCheckpoint":                  0,
	"GetCheckpointData":                    0,
	"GetCheckpointAvailability":            0,
	"SimulatePTB":                          0,
	"GetCoinMetadata":                      0,
}

// BindingsClient is the subset of SuiPTBClient used by the bindings module.
type BindingsClient interface {
	ReadObjectId(ctx context.Context, objectId string) (*suirpcv2.Object, error)
	QueryCoinsByAddress(ctx context.Context, address, coinType string) ([]*suirpcv2.Object, error)
	GetReferenceGasPrice(ctx context.Context) (*big.Int, error)
	SimulatePTB(ctx context.Context, bcsBytes []byte) ([]any, error)
	SendTransaction(ctx context.Context, req *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error)
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	FinishPTBAndSend(ctx context.Context, signer *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (*suirpcv2.ExecuteTransactionResponse, error)
}

var _ BindingsClient = (*PTBClient)(nil)

type SuiPTBClient interface {
	MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error)
	SendTransaction(ctx context.Context, execRequest *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error)
	ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor []byte) ([]*suirpcv2.Object, error)
	ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor []byte) ([]*suirpcv2.Object, error)
	ReadObjectId(ctx context.Context, objectId string) (*suirpcv2.Object, error)
	ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error)
	SimulatePTB(ctx context.Context, bcsBytes []byte) ([]any, error)
	SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte) (*suirpcv2.ExecuteTransactionResponse, error)
	QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*models.PaginatedEventsResponse, error)
	QueryTransactions(ctx context.Context, fromAddress string, cursor *suirpcv2.Checkpoint, limit *uint64) ([]*suirpcv2.ExecutedTransaction, error)
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]*suirpcv2.Object, error)
	QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]*suirpcv2.Object, error)
	EstimateGas(ctx context.Context, tx *transaction.Transaction) (uint64, error)
	GetReferenceGasPrice(ctx context.Context) (*big.Int, error)
	FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (*suirpcv2.ExecuteTransactionResponse, error)
	GetBlockById(ctx context.Context, checkpointDigest string) (*suirpcv2.Checkpoint, error)
	GetLatestEpoch(ctx context.Context) (*suirpcv2.Epoch, error)
	GetLatestCheckpoint(ctx context.Context) (*suirpcv2.Checkpoint, error)
	GetCheckpointData(ctx context.Context, checkpointSequenceNumber uint64) (*CheckpointData, error)
	GetNormalizedModule(ctx context.Context, packageId string, moduleId string) (models.GetNormalizedMoveModuleResponse, error)
	GetMoveModuleFunction(ctx context.Context, packageId string, moduleId string, functionName string) (*suirpcv2.FunctionDescriptor, error)
	GetSUIBalance(ctx context.Context, address string) (*suirpcv2.Balance, error)
	LoadModulePackageIds(ctx context.Context, packageId string, module string) ([]string, error)
	GetLatestPackageId(ctx context.Context, packageId string, module string) (string, error)
	GetCoinMetadata(ctx context.Context, coinType string) (models.CoinMetadataResponse, error)
	GetCache() *cache.Cache
	GetCachedValue(key string) (any, bool)
	SetCachedValue(key string, value any)
	GetCachedValues(keys []string) (map[string]any, bool)
	SetCachedValues(keyValues map[string]any)
	HashTxBytes(txBytes []byte) []byte
	GetCCIPPackageID(ctx context.Context, offRampPackageID string) (string, error)
	GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error)
	GetParentObjectID(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error)
}

// PTBClient implements SuiClient interface using the blockvision SDK.
// During the gRPC migration, JSON-RPC (client) and gRPC (grpcClient) coexist:
// migrated methods use gRPC service accessors; others continue via JSON-RPC.
type PTBClient struct {
	log              logger.Logger
	moveModuleClient sui.ISuiAPI // internal JSON-RPC client for unmigrated read APIs
	// connPool holds one or more independent gRPC connections to the node, handed out round-robin by the
	// service getters. Using several connections multiplies the available concurrent HTTP/2 streams so a
	// single connection's stream limit does not become a head-of-line bottleneck under bursty reads.
	connPool           *grpcConnPool
	maxRetries         *int
	transactionTimeout time.Duration
	keystoreService    loop.Keystore
	rateLimiter        *semaphore.Weighted
	defaultRequestType TransactionRequestType
	devInspectSigner   *signer.Signer

	// map of module name to normalized module definition (similar to an ABI)
	normalizedModules map[string]map[string]models.GetNormalizedMoveModuleResponse

	cache *cache.Cache // used for caching object IDs (e.g. offramp state object ID or state pointers)

	// objectCache, when set, de-duplicates and caches version-stable object reference metadata so the
	// per-read GetObject fan-out does not hit the node on every read. It is injected via PTBClientConfig
	// (the implementation lives in chainreader/reader, which imports this package, so it is referenced
	// here only through the ObjectMetadataCache interface to avoid an import cycle). May be nil.
	objectCache ObjectMetadataCache
}

// ObjectMetadataCache caches version-stable object reference metadata (owner/version/digest) to avoid
// redundant GetObject RPCs on the read hot path. It is satisfied by chainreader/reader.Cache and
// supplied via PTBClientConfig.ObjectCache.
type ObjectMetadataCache interface {
	GetObjectMetadata(ctx context.Context, objectID string, loader func(context.Context) (*suirpcv2.Object, error)) (*suirpcv2.Object, error)
}

var _ SuiPTBClient = (*PTBClient)(nil)

// ExtendedPTBClient augments SuiPTBClient with additional methods that are not part
// of the stable SuiPTBClient contract. Consumers that need transaction details by
// digest can depend on this interface without forcing every SuiPTBClient
// implementer to change.
type ExtendedPTBClient interface {
	SuiPTBClient
	GetTransaction(ctx context.Context, digest string) (TransactionDetails, error)
	GetCheckpointAvailability(ctx context.Context) (*suirpcv2.GetServiceInfoResponse, error)
}

var _ ExtendedPTBClient = (*PTBClient)(nil)

func NewPTBClient(log logger.Logger, cfg PTBClientConfig) (*PTBClient, error) {
	return NewPTBClientFromConfig(log, cfg)
}

func (c *PTBClient) WithRateLimit(ctx context.Context, methodName string, f func(ctx context.Context) error) error {
	weight := int64(1)
	if weightValue, ok := RateLimitWeights[methodName]; ok {
		weight = weightValue
	}
	// start a timer to track the duration of the function
	timer := time.Now()

	workCtx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	// If rate limiter is disabled or weight is 0, skip semaphore entirely.
	// This will skip adding to the semaphore queue and prevent unnecessary queuing.
	if c.rateLimiter == nil || weight == 0 {
		err := f(workCtx)
		duration := time.Since(timer)
		c.log.Debugw("WithRateLimit completed", "methodName", methodName, "duration", duration)
		return err
	}

	// acquire with the timeout context so it can't hang forever
	if err := c.rateLimiter.Acquire(workCtx, weight); err != nil {
		duration := time.Since(timer)
		c.log.Debugw("WithRateLimit failed to acquire rate limit", "methodName", methodName, "duration", duration)
		return fmt.Errorf("failed to acquire rate limit for %s: %w", methodName, err)
	}
	acquireWait := time.Since(timer)
	execStart := time.Now()

	// ensure cleanup on exit
	defer func() {
		c.rateLimiter.Release(weight)
	}()

	// run the user function with the timeout context
	// if the function respects the context, it will return and lock will be released in defer
	err := f(workCtx)
	duration := time.Since(timer)
	c.log.Debugw("WithRateLimit completed", "methodName", methodName, "duration", duration)

	// distinguish time spent waiting on the semaphore vs running the RPC.
	if duration > 2*time.Second || acquireWait > 500*time.Millisecond {
		c.log.Debugw(
			"WithRateLimit slow rate-limited call",
			"methodName", methodName,
			"acquireWaitSec", acquireWait.Seconds(),
			"execSec", time.Since(execStart).Seconds(),
			"totalSec", duration.Seconds(),
			"ctxErr", ctxErrString(ctx),
		)
	}
	return err
}

func ctxErrString(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if err := ctx.Err(); err != nil {
		return err.Error()
	}
	return ""
}

// MoveCall is a method that's used primarily in tests for adhoc contract calls.
// It simply builds the transaction bytes and returns them.
func (c *PTBClient) MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error) {
	movePkgService, err := c.getMovePackageService(ctx)
	if err != nil {
		return TxnMetaData{}, fmt.Errorf("failed to get move package service: %w", err)
	}

	txn := transaction.NewTransaction()

	var args []transaction.Argument
	var typeArgs []transaction.TypeTag

	functionDefinition, err := movePkgService.GetFunction(ctx, &suirpcv2.GetFunctionRequest{
		PackageId:  &req.PackageObjectId,
		ModuleName: &req.Module,
		Name:       &req.Function,
	})
	if err != nil {
		return TxnMetaData{}, fmt.Errorf("failed to get normalized module: %w", err)
	}

	for i, arg := range req.Arguments {
		paramBody := functionDefinition.GetFunction().GetParameters()[i].GetBody()
		txArg, transformErr := c.transformMoveCallArgFromSignature(ctx, txn, arg, paramBody, true)
		if transformErr != nil {
			return TxnMetaData{}, fmt.Errorf("failed to transform transaction arg: %w", transformErr)
		}
		args = append(args, *txArg)
	}

	txn.MoveCall(models.SuiAddress(req.PackageObjectId), req.Module, req.Function, typeArgs, args)

	gasPrice, err := c.GetReferenceGasPrice(ctx)
	if err != nil {
		return TxnMetaData{}, fmt.Errorf("failed to get reference gas price: %w", err)
	}

	txn.SetGasBudget(req.GasBudget)
	txn.SetGasPrice(gasPrice.Uint64())
	txn.SetSender(models.SuiAddress(req.Signer))

	// Note: this is only a placeholder to ensure `buildBCSBytes` doesn't fail.
	// The actual signing is handled externally.
	txn.SetSigner(&signer.Signer{
		Address: req.Signer,
	})

	paymentCoinBytes, paymentCoinVersion, paymentCoinDigest, err := c.GetTransactionPaymentCoinForAddress(context.Background(), req.Signer)
	if err != nil {
		return TxnMetaData{}, fmt.Errorf("failed to get transaction payment coin for address: %w", err)
	}

	txn.SetGasPayment([]transaction.SuiObjectRef{
		{
			ObjectId: paymentCoinBytes,
			Version:  paymentCoinVersion,
			Digest:   paymentCoinDigest,
		},
	})

	bcsBytes, err := txn.BuildBCSBytes(context.Background())
	if err != nil {
		return TxnMetaData{}, fmt.Errorf("failed to build bcs bytes: %w", err)
	}

	return TxnMetaData{
		TxBytes: base64.StdEncoding.EncodeToString(bcsBytes),
	}, nil
}

// SendTransaction executes an already signed transaction, using the execution service, given the BCS bytes.
func (c *PTBClient) SendTransaction(ctx context.Context, execRequest *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error) {
	txExecService, err := c.getTransactionExecutionService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction execution service: %w", err)
	}

	response, err := txExecService.ExecuteTransaction(ctx, execRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to execute transaction: %w", err)
	}
	return response, nil
}

func (c *PTBClient) ReadObjectId(ctx context.Context, objectId string) (*suirpcv2.Object, error) {
	var result *suirpcv2.Object
	err := c.WithRateLimit(ctx, "ReadObjectId", func(ctx context.Context) error {
		var err error
		result, err = c.readObjectIdInternal(ctx, objectId)
		return err
	})
	return result, err
}

// readObjectIdInternal is the internal implementation without rate limiting
func (c *PTBClient) readObjectIdInternal(ctx context.Context, objectId string) (*suirpcv2.Object, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	objectReq := suirpcv2.GetObjectRequest{
		ObjectId: &objectId,
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"contents",
				"object_type",
				"owner",
				"version",
				"digest",
				"json",
			},
		},
	}

	response, err := ledgerService.GetObject(ctx, &objectReq)
	if err != nil || response.Object == nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return response.Object, nil
}

// readObjectMetadataInternal fetches only the lightweight reference metadata of an object
// (object_type, owner, version, digest) and intentionally omits the potentially large `contents`/`json`
// fields. This is all that is needed to build a transaction ObjectArg, and avoids transferring/serializing
// the full (and growing) object state on every read. When an objectCache is configured it is consulted
// first (and concurrent reads of the same object are collapsed) so version-stable shared/immutable objects
// are not re-fetched on every read.
func (c *PTBClient) readObjectMetadataInternal(ctx context.Context, objectId string) (*suirpcv2.Object, error) {
	if c.objectCache != nil {
		return c.objectCache.GetObjectMetadata(ctx, objectId, func(ctx context.Context) (*suirpcv2.Object, error) {
			return c.fetchObjectMetadata(ctx, objectId)
		})
	}
	return c.fetchObjectMetadata(ctx, objectId)
}

// fetchObjectMetadata performs the actual GetObject RPC for an object's reference metadata.
func (c *PTBClient) fetchObjectMetadata(ctx context.Context, objectId string) (*suirpcv2.Object, error) {
	metaStart := time.Now()
	ledgerService, err := c.getLedgerService(ctx)
	svcDur := time.Since(metaStart)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	objectReq := suirpcv2.GetObjectRequest{
		ObjectId: &objectId,
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"object_type",
				"owner",
				"version",
				"digest",
			},
		},
	}

	rpcStart := time.Now()
	response, err := ledgerService.GetObject(ctx, &objectReq)
	rpcDur := time.Since(rpcStart)
	if time.Since(metaStart) > 2*time.Second {
		objType := ""
		if response != nil && response.Object != nil {
			objType = response.Object.GetObjectType()
		}
		common.DebugLog("grpc_client.go:readObjectMetadataInternal", "metadata read breakdown", map[string]any{
			"hypothesisId": "H12",
			"objectId":     objectId,
			"objectType":   objType,
			"svcSec":       svcDur.Seconds(),
			"rpcSec":       rpcDur.Seconds(),
			"totalSec":     time.Since(metaStart).Seconds(),
			"err":          err != nil,
		})
	}
	if err != nil || response.Object == nil {
		return nil, fmt.Errorf("failed to read object metadata: %w", err)
	}

	return response.Object, nil
}

func (c *PTBClient) ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor []byte) ([]*suirpcv2.Object, error) {
	var result []*suirpcv2.Object

	err := c.WithRateLimit(ctx, "ReadFilterOwnedObjectIds", func(ctx context.Context) error {
		response, err := c.readFilterOwnedObjectIdsInternal(ctx, ownerAddress, &structType, cursor)
		if err != nil {
			return fmt.Errorf("failed to read filter owned object ids: %w", err)
		}

		result = response.Objects

		return nil
	})

	return result, err
}

func (c *PTBClient) readFilterOwnedObjectIdsInternal(ctx context.Context, ownerAddress string, structType *string, cursor []byte) (*suirpcv2.ListOwnedObjectsResponse, error) {
	stateService, err := c.getStateService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get state service: %w", err)
	}

	response, err := stateService.ListOwnedObjects(ctx, &suirpcv2.ListOwnedObjectsRequest{
		Owner:      &ownerAddress,
		ObjectType: structType,
		PageToken:  cursor,
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"contents",
				"object_type",
				"owner",
				"version",
				"digest",
				"json",
				"balance",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read owned objects: %w", err)
	}

	// NOTE: this will recursively call itself until there are no more pages
	if response.NextPageToken != nil {
		nextPage, err := c.readFilterOwnedObjectIdsInternal(ctx, ownerAddress, structType, response.NextPageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to read next page of owned objects: %w", err)
		}
		response.Objects = append(response.Objects, nextPage.Objects...)
	}

	return response, nil
}

func (c *PTBClient) ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor []byte) ([]*suirpcv2.Object, error) {
	var result []*suirpcv2.Object
	err := c.WithRateLimit(ctx, "ReadOwnedObjects", func(ctx context.Context) error {
		var err error
		response, err := c.readFilterOwnedObjectIdsInternal(ctx, ownerAddress, nil, cursor)
		if err != nil {
			return fmt.Errorf("failed to read owned objects: %w", err)
		}
		result = response.Objects
		return err
	})
	return result, err
}

func (c *PTBClient) EstimateGas(ctx context.Context, tx *transaction.Transaction) (uint64, error) {
	var result uint64
	txExecService, err := c.getTransactionExecutionService(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get transaction execution service: %w", err)
	}

	err = c.WithRateLimit(ctx, "EstimateGas", func(ctx context.Context) error {
		bcsBytes, buildErr := tx.BuildBCSBytes(ctx)
		if buildErr != nil {
			return fmt.Errorf("failed to build bcs bytes: %w", buildErr)
		}

		doGasSelection := true

		response, simErr := txExecService.SimulateTransaction(ctx, &suirpcv2.SimulateTransactionRequest{
			Transaction:    &suirpcv2.Transaction{Bcs: &suirpcv2.Bcs{Value: bcsBytes}},
			DoGasSelection: &doGasSelection,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"transaction",
					"transaction.effects.status",
					"transaction.effects.gas_used",
					"transaction.effects.gas_used.computation_cost",
					"transaction.effects.gas_used.storage_cost",
					"transaction.effects.gas_used.storage_rebate",
					"transaction.effects.gas_used.non_refundable_storage_fee",
				},
			},
		})
		if simErr != nil {
			return fmt.Errorf("failed to simulate transaction: %w", simErr)
		}

		gasUsed := response.GetTransaction().GetEffects().GetGasUsed()

		computationCost := gasUsed.GetComputationCost()
		storageCost := gasUsed.GetStorageCost()
		storageRebate := gasUsed.GetStorageRebate()
		nonRefundableStorageFee := gasUsed.GetNonRefundableStorageFee()

		result = computationCost + storageCost + nonRefundableStorageFee

		// Avoid uint64 underflow in case of a node bug since this is a simulated transaction
		if storageRebate < result {
			result -= storageRebate
		}

		c.log.Debugw("Estimated gas", "computationCost", computationCost, "storageCost", storageCost, "storageRebate", storageRebate, "nonRefundableStorageFee", nonRefundableStorageFee, "result", result)

		return nil
	})

	return result, err
}

func (c *PTBClient) GetReferenceGasPrice(ctx context.Context) (*big.Int, error) {
	var result *big.Int
	err := c.WithRateLimit(ctx, "GetReferenceGasPrice", func(ctx context.Context) error {
		ledgerService, err := c.getLedgerService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get ledger service: %w", err)
		}

		resp, err := ledgerService.GetEpoch(ctx, &suirpcv2.GetEpochRequest{
			ReadMask: &fieldmaskpb.FieldMask{Paths: []string{"reference_gas_price"}},
		})
		if err != nil {
			return fmt.Errorf("GetEpoch: %w", err)
		}
		var gasPrice uint64
		if resp.Epoch.ReferenceGasPrice != nil {
			gasPrice = *resp.Epoch.ReferenceGasPrice
		}

		result = new(big.Int).SetUint64(gasPrice)
		return nil
	})
	return result, err
}

func (c *PTBClient) ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	var results []any
	err := c.WithRateLimit(ctx, "ReadFunction", func(ctx context.Context) error {
		var err error
		results, err = c.readFunctionInternal(ctx, packageId, module, function, args, argTypes, nil, typeArgs)
		return err
	})
	return results, err
}

// ReadFunctionWithMutability behaves like ReadFunction but lets callers control, per argument,
// whether an object argument is passed as a mutable shared-object reference. This is required for
// Move view functions that take immutable references (&T): passing such an object as mutable causes
// the node to reject the simulation ("Mutable parameter provided, immutable parameter expected").
// mutabilities is index-aligned with args; a nil or short slice defaults an argument to mutable.
func (c *PTBClient) ReadFunctionWithMutability(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, mutabilities []bool, typeArgs []string) ([]any, error) {
	var results []any
	err := c.WithRateLimit(ctx, "ReadFunction", func(ctx context.Context) error {
		var err error
		results, err = c.readFunctionInternal(ctx, packageId, module, function, args, argTypes, mutabilities, typeArgs)
		return err
	})
	return results, err
}

// readFunctionInternal is the internal implementation without rate limiting
func (c *PTBClient) readFunctionInternal(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, mutabilities []bool, typeArgs []string) ([]any, error) {
	readCtx := ctx
	cancel := func() {}
	timeout := DefaultReadOpTimeout
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			// Parent deadline already passed; ctx is already cancelled.
		} else {
			if remaining < timeout {
				timeout = remaining
			}
			readCtx, cancel = context.WithTimeout(ctx, timeout)
		}
	} else {
		readCtx, cancel = context.WithTimeout(ctx, timeout)
	}
	defer cancel()

	rfStart := time.Now()
	var svcDur, typeArgDur, transformDur, buildDur, simDur time.Duration
	txExecService, err := c.getTransactionExecutionService(readCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction execution service: %w", err)
	}
	svcDur = time.Since(rfStart)

	var results []any
	txn := transaction.NewTransaction()

	var txnArgs []transaction.Argument
	var txnTypeArgs []transaction.TypeTag

	// Process type arguments
	typeArgStart := time.Now()
	for _, typeArg := range typeArgs {
		typeTag, typeTagErr := c.CreateTypeTag(typeArg)
		if typeTagErr != nil {
			return nil, fmt.Errorf("failed to create type tag for %s: %w", typeArg, typeTagErr)
		}
		txnTypeArgs = append(txnTypeArgs, typeTag)
	}
	typeArgDur = time.Since(typeArgStart)
	transformStart := time.Now()

	for i, arg := range args {
		argType, ok := common.ValueAt(argTypes, i)
		if !ok {
			argType = common.InferArgumentType(arg)
		}

		// Default object arguments to mutable to preserve prior behavior; honor an explicit
		// per-argument mutability when provided (needed for immutable-reference view functions).
		mutable := true
		if i < len(mutabilities) {
			mutable = mutabilities[i]
		}

		txArg, transformErr := c.TransformTransactionArg(readCtx, txn, arg, argType, mutable)
		if transformErr != nil {
			return nil, fmt.Errorf("failed to transform transaction arg: %w", transformErr)
		}
		txnArgs = append(txnArgs, *txArg)
	}
	transformDur = time.Since(transformStart)

	// set a read-only (no funds) signer for read (simulate) calls
	txn.SetSigner(c.devInspectSigner)
	txn.SetSender(models.SuiAddress(c.devInspectSigner.Address))

	txn.SetGasBudget(DefaultGasBudget)
	txn.SetGasPrice(DefaultGasPrice)
	txn.MoveCall(models.SuiAddress(packageId), module, function, txnTypeArgs, txnArgs)

	// Get transaction bytes
	buildStart := time.Now()
	bcsBytes, err := txn.BuildBCSBytes(readCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to build bcs bytes: %w", err)
	}
	buildDur = time.Since(buildStart)
	simStart := time.Now()

	results, err = c.simulatePTBInternal(readCtx, txExecService, bcsBytes, false)
	simDur = time.Since(simStart)
	totalDur := time.Since(rfStart)
	if totalDur > 3*time.Second {
		c.log.Debugw("ReadFunction slow read breakdown", "module", module, "function", function, "numArgs", len(args), "numTypeArgs", len(typeArgs), "getSvcSec", svcDur.Seconds(), "typeArgSec", typeArgDur.Seconds(), "transformSec", transformDur.Seconds(), "buildSec", buildDur.Seconds(), "simulateSec", simDur.Seconds(), "totalSec", totalDur.Seconds(), "ctxErr", ctxErrString(readCtx))
	}
	if err != nil {
		return nil, err
	}

	c.log.Debugw("ReadFunction results", "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function), "results", results)

	return results, nil
}

// SimulatePTB simulates a pre-built PTB and returns JSON-decoded Move return values.
func (c *PTBClient) SimulatePTB(ctx context.Context, bcsBytes []byte) ([]any, error) {
	var results []any
	err := c.WithRateLimit(ctx, "SimulatePTB", func(ctx context.Context) error {
		txExecService, err := c.getTransactionExecutionService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get transaction execution service: %w", err)
		}

		var simErr error
		results, simErr = c.simulatePTBInternal(ctx, txExecService, bcsBytes, true)
		return simErr
	})
	return results, err
}

func (c *PTBClient) simulatePTBInternal(ctx context.Context, txExecService suirpcv2.TransactionExecutionServiceClient, bcsBytes []byte, checks bool) ([]any, error) {
	doGasSelection := false
	checksEnum := suirpcv2.SimulateTransactionRequest_DISABLED.Enum()
	if checks {
		checksEnum = suirpcv2.SimulateTransactionRequest_ENABLED.Enum()
	}

	// measure the raw SimulateTransaction RPC latency and the number of simulate calls
	// concurrently hitting the single gRPC connection / local Sui node.
	simStart := time.Now()
	response, err := txExecService.SimulateTransaction(ctx, &suirpcv2.SimulateTransactionRequest{
		Transaction:    &suirpcv2.Transaction{Bcs: &suirpcv2.Bcs{Value: bcsBytes}},
		DoGasSelection: &doGasSelection,
		Checks:         checksEnum,
	})
	simDur := time.Since(simStart)
	if simDur > time.Second {
		c.log.Debugw("SimulateTransaction slow SimulateTransaction RPC", "simSec", simDur.Seconds(), "simErr", err != nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to simulate transaction: %w", err)
	}

	if response.Transaction != nil && response.Transaction.Effects != nil && response.Transaction.Effects.Status != nil {
		if !response.Transaction.Effects.Status.GetSuccess() {
			errMsg := response.Transaction.Effects.Status.GetError()
			if errMsg != nil {
				return nil, fmt.Errorf("simulate failed: %s", errMsg.GetDescription())
			}
			return nil, errors.New("simulate failed")
		}
	}

	if len(response.CommandOutputs) == 0 {
		return []any{}, nil
	}

	returnedValues := response.CommandOutputs[0].GetReturnValues()
	if len(returnedValues) == 0 {
		return []any{}, nil
	}

	results := make([]any, len(returnedValues))
	for i, returnedValue := range returnedValues {
		results[i] = returnedValue.Json.AsInterface()
	}

	return results, nil
}

func (c *PTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte) (*suirpcv2.ExecuteTransactionResponse, error) {
	signerPublicKeyId := fmt.Sprintf("%064x", signerPublicKey)
	bcsBytes, err := base64.StdEncoding.DecodeString(txBytesRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx bytes: %w", err)
	}

	hashedTxBytes := c.HashTxBytes(bcsBytes)
	signatureBytes, err := c.keystoreService.Sign(ctx, signerPublicKeyId, hashedTxBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Sui signature format: [scheme_byte (0x00 for ED25519) + signature_bytes + public_key_bytes]
	// The scheme byte must be prepended to the signature
	schemeByte := byte(0x00) // ED25519 scheme identifier
	suiSignature := append([]byte{schemeByte}, signatureBytes...)
	suiSignature = append(suiSignature, signerPublicKey...)

	// Verify we can execute the transaction
	resp, err := c.SendTransaction(ctx, &suirpcv2.ExecuteTransactionRequest{
		Transaction: &suirpcv2.Transaction{Bcs: &suirpcv2.Bcs{Value: bcsBytes}},
		Signatures: []*suirpcv2.UserSignature{
			{
				Bcs: &suirpcv2.Bcs{Value: suiSignature},
			},
		},
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"transaction",
				"transaction.digest",
				"digest",
				"effects.digest",
				"effects.status",
				"effects.gas_used",
				"events",
				"events.events.package_id",
				"events.events.module",
				"events.events.event_type",
				"events.events.json",
			},
		},
	})

	return resp, err
}

func (c *PTBClient) QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	return nil, errors.New("method implementation pending gRPC migration")
}

// GetEventsByCheckpoint returns all the events for a given checkpoint sequence number.
// @param checkpointSequenceNumber - the checkpoint sequence number to get events for
// @param eventTypes - the types of events to get (must be fully qualified `packageId::moduleId::EventName`)
// @return the events and an error if any
func (c *PTBClient) GetEventsByCheckpoint(ctx context.Context, checkpointSequenceNumber uint64, eventTypes []string) ([]*suirpcv2.Event, error) {
	transactions, err := c.GetTransactionsByCheckpoint(ctx, checkpointSequenceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	var events []*suirpcv2.Event
	for _, transaction := range transactions {
		for _, event := range transaction.GetEvents().Events {
			qualifiedEventHandle := strings.Join([]string{event.GetPackageId(), event.GetModule(), event.GetEventType()}, "::")
			if slices.Contains(eventTypes, qualifiedEventHandle) {
				events = append(events, event)
			}
		}
	}

	return events, nil
}

func (c *PTBClient) GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return TransactionResult{}, fmt.Errorf("failed to get ledger service: %w", err)
	}

	var result TransactionResult

	err = c.WithRateLimit(ctx, "GetTransactionStatus", func(ctx context.Context) error {
		response, txErr := ledgerService.GetTransaction(ctx, &suirpcv2.GetTransactionRequest{
			Digest: &digest,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"effects.status",
					"effects.status.success",
					"effects.status.error",
					"effects.gas_used",
					"checkpoint",
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		tx := response.GetTransaction()
		if tx == nil || tx.GetEffects() == nil || tx.GetEffects().GetStatus() == nil {
			return fmt.Errorf("transaction %s missing status", digest)
		}

		var status string
		success := tx.GetEffects().GetStatus().GetSuccess()
		if success {
			status = "success"
		} else {
			status = "failure"
		}

		result = TransactionResult{
			Status:     status,
			Error:      tx.GetEffects().GetStatus().GetError().String(),
			Checkpoint: tx.GetCheckpoint(),
		}

		return nil
	})

	return result, err
}

// GetTransaction fetches transaction details by digest, including the sender
// address, in a single gRPC call. It avoids the checkpoint scan required to
// recover the sender when only GetTransactionStatus is available.
func (c *PTBClient) GetTransaction(ctx context.Context, digest string) (TransactionDetails, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return TransactionDetails{}, fmt.Errorf("failed to get ledger service: %w", err)
	}

	var result TransactionDetails

	err = c.WithRateLimit(ctx, "GetTransaction", func(ctx context.Context) error {
		response, txErr := ledgerService.GetTransaction(ctx, &suirpcv2.GetTransactionRequest{
			Digest: &digest,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"digest",
					"transaction.sender",
					"effects.status",
					"effects.status.success",
					"effects.status.error",
					"checkpoint",
					"timestamp",
				},
			},
		})
		if txErr != nil {
			return txErr
		}

		tx := response.GetTransaction()
		if tx == nil || tx.GetEffects() == nil || tx.GetEffects().GetStatus() == nil {
			return fmt.Errorf("transaction %s missing status", digest)
		}

		status := "failure"
		if tx.GetEffects().GetStatus().GetSuccess() {
			status = "success"
		}

		var timestampSeconds uint64
		if ts := tx.GetTimestamp(); ts != nil {
			//nolint:gosec // timestampSeconds is a positive value that is safe to convert to uint64
			timestampSeconds = uint64(ts.GetSeconds())
		}

		result = TransactionDetails{
			Digest:     tx.GetDigest(),
			Status:     status,
			Error:      tx.GetEffects().GetStatus().GetError().String(),
			Checkpoint: tx.GetCheckpoint(),
			Sender:     tx.GetTransaction().GetSender(),
			Timestamp:  timestampSeconds,
		}

		return nil
	})

	return result, err
}

func (c *PTBClient) GetTransactionChangedObjects(ctx context.Context, digest string) ([]*suirpcv2.ChangedObject, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	var changed []*suirpcv2.ChangedObject
	err = c.WithRateLimit(ctx, "GetTransactionChangedObjects", func(ctx context.Context) error {
		response, txErr := ledgerService.GetTransaction(ctx, &suirpcv2.GetTransactionRequest{
			Digest: &digest,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"effects.changed_objects",
					"effects.changed_objects.object_id",
					"effects.changed_objects.object_type",
					"effects.changed_objects.id_operation",
					"effects.changed_objects.output_state",
					"effects.changed_objects.output_owner",
					"effects.changed_objects.output_version",
					"effects.changed_objects.output_digest",
					"effects.changed_objects.input_version",
				},
			},
		})
		if txErr != nil {
			return txErr
		}
		if response.GetTransaction() == nil || response.GetTransaction().GetEffects() == nil {
			return fmt.Errorf("transaction %s missing effects", digest)
		}
		changed = response.GetTransaction().GetEffects().GetChangedObjects()
		return nil
	})
	if err != nil {
		return nil, err
	}

	return changed, nil
}

// QueryTransactions queries the transactions for a given address.
// @param fromAddress - the address to query transactions for
// @param cursor - a checkpoint ID to start from, if nil the latest checkpoint is used
// @param limit - the limit of transactions to return
// @return the transactions and an error if any
func (c *PTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *suirpcv2.Checkpoint, limit *uint64) ([]*suirpcv2.ExecutedTransaction, error) {
	return nil, errors.New("method implementation pending gRPC migration")
}

// GetTransactionsByCheckpoint returns all the transactions for a given checkpoint sequence number.
func (c *PTBClient) GetTransactionsByCheckpoint(ctx context.Context, checkpointSequenceNumber uint64) ([]*suirpcv2.ExecutedTransaction, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	response, err := ledgerService.GetCheckpoint(ctx, &suirpcv2.GetCheckpointRequest{
		CheckpointId: &suirpcv2.GetCheckpointRequest_SequenceNumber{
			SequenceNumber: checkpointSequenceNumber,
		},
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: checkpointTransactionReadMaskPaths,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoint: %w", err)
	}

	return response.GetCheckpoint().GetTransactions(), nil
}

// CheckpointData combines checkpoint metadata with its transactions.
type CheckpointData struct {
	Checkpoint   *suirpcv2.Checkpoint
	Transactions []*suirpcv2.ExecutedTransaction
}

// Read mask paths for checkpoint indexing. Nested event fields must be listed
// explicitly; requesting only "transactions.events" can return TransactionEvents
// containers without Event metadata, causing the indexer to skip real events.
var checkpointTransactionReadMaskPaths = []string{
	"sequence_number",
	"digest",
	"summary.timestamp",
	"transactions",
	"transactions.digest",
	"transactions.transaction.sender",
	"transactions.transaction.kind",
	"transactions.transaction.kind.programmable_transaction.commands",
	"transactions.transaction.kind.programmable_transaction.commands.move_call.package",
	"transactions.transaction.kind.programmable_transaction.commands.move_call.module",
	"transactions.transaction.kind.programmable_transaction.commands.move_call.function",
	"transactions.transaction.kind.programmable_transaction.commands.move_call.arguments",
	"transactions.transaction.kind.programmable_transaction.commands.move_call.arguments.kind",
	"transactions.transaction.kind.programmable_transaction.commands.move_call.arguments.input",
	"transactions.transaction.kind.programmable_transaction.inputs",
	"transactions.transaction.kind.programmable_transaction.inputs.pure",
	"transactions.effects.status",
	"transactions.effects.status.success",
	"transactions.effects.status.error",
	"transactions.effects.status.error.command",
	"transactions.effects.status.error.abort",
	"transactions.effects.status.error.abort.abort_code",
	"transactions.effects.status.error.abort.location",
	"transactions.effects.events_digest",
	"transactions.events.events.package_id",
	"transactions.events.events.module",
	"transactions.events.events.event_type",
	"transactions.events.events.json",
}

var transactionEventsReadMaskPaths = []string{
	"digest",
	"effects.events_digest",
	"events.events.package_id",
	"events.events.module",
	"events.events.event_type",
	"events.events.json",
}

// TODO: this should be the responsibility of the indexer, not the client
// HydrateTransactionEvents fetches full event payloads when a checkpoint
// transaction reports an events_digest but omits inline TransactionEvents.
func (c *PTBClient) HydrateTransactionEvents(ctx context.Context, tx *suirpcv2.ExecutedTransaction) {
	if tx == nil {
		return
	}
	if tx.GetEvents() != nil && len(tx.GetEvents().GetEvents()) > 0 {
		return
	}
	if tx.GetEffects() == nil || tx.GetEffects().GetEventsDigest() == "" {
		return
	}

	digest := tx.GetDigest()
	if digest == "" {
		return
	}

	err := c.WithRateLimit(ctx, "GetTransactionStatus", func(ctx context.Context) error {
		ledgerService, err := c.getLedgerService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get ledger service: %w", err)
		}

		response, err := ledgerService.GetTransaction(ctx, &suirpcv2.GetTransactionRequest{
			Digest: &digest,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: transactionEventsReadMaskPaths,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get transaction %s: %w", digest, err)
		}

		events := response.GetTransaction().GetEvents()
		if events != nil && len(events.GetEvents()) > 0 {
			tx.Events = events
		}

		return nil
	})
	if err != nil {
		c.log.Warnw("Failed to hydrate transaction events", "digest", digest, "error", err)
	}
}

// GetCheckpointData returns checkpoint metadata plus all transactions for a given sequence number.
func (c *PTBClient) GetCheckpointData(ctx context.Context, checkpointSequenceNumber uint64) (*CheckpointData, error) {
	var result *CheckpointData

	err := c.WithRateLimit(ctx, "GetCheckpointData", func(ctx context.Context) error {
		ledgerService, err := c.getLedgerService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get ledger service: %w", err)
		}

		response, err := ledgerService.GetCheckpoint(ctx, &suirpcv2.GetCheckpointRequest{
			CheckpointId: &suirpcv2.GetCheckpointRequest_SequenceNumber{
				SequenceNumber: checkpointSequenceNumber,
			},
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: checkpointTransactionReadMaskPaths,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get checkpoint: %w", err)
		}

		transactions := response.GetCheckpoint().GetTransactions()
		for _, tx := range transactions {
			c.HydrateTransactionEvents(ctx, tx)
		}

		result = &CheckpointData{
			Checkpoint:   response.GetCheckpoint(),
			Transactions: transactions,
		}
		return nil
	})

	return result, err
}

// GetLatestCheckpoint returns the latest checkpoint from the chain.
// Uses GetCheckpointRequest with empty CheckpointId which returns the latest.
func (c *PTBClient) GetLatestCheckpoint(ctx context.Context) (*suirpcv2.Checkpoint, error) {
	var result *suirpcv2.Checkpoint

	err := c.WithRateLimit(ctx, "GetLatestCheckpoint", func(ctx context.Context) error {
		ledgerService, err := c.getLedgerService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get ledger service: %w", err)
		}

		response, err := ledgerService.GetCheckpoint(ctx, &suirpcv2.GetCheckpointRequest{
			// Empty CheckpointId returns the latest checkpoint
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"sequence_number",
					"digest",
					"summary.timestamp",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get latest checkpoint: %w", err)
		}

		result = response.GetCheckpoint()
		return nil
	})

	return result, err
}

// GetCoinsByAddress returns all coin objects for a given address.
func (c *PTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]*suirpcv2.Object, error) {
	var result []*suirpcv2.Object
	err := c.WithRateLimit(ctx, "GetCoinsByAddress", func(ctx context.Context) error {
		coinTag := "0x2::coin::Coin"
		response, err := c.readFilterOwnedObjectIdsInternal(ctx, address, &coinTag, nil)
		if err != nil {
			return fmt.Errorf("failed to get coins: %w", err)
		}
		result = response.GetObjects()

		return nil
	})

	return result, err
}

// QueryCoinsByAddress is the same as GetCoinsByAddress, but it allows you to filter by coin type.
// The `coinType` parameter should be the full coin type, e.g. "0x2::coin::Coin<0x2::sui::SUI>".
func (c *PTBClient) QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]*suirpcv2.Object, error) {
	var result []*suirpcv2.Object
	err := c.WithRateLimit(ctx, "GetCoinsByAddress", func(ctx context.Context) error {
		response, err := c.readFilterOwnedObjectIdsInternal(ctx, address, &coinType, nil)
		if err != nil {
			return fmt.Errorf("failed to get coins: %w", err)
		}
		result = response.GetObjects()

		return nil
	})

	return result, err
}

// FinishPTBAndSend finishes the PTB transaction and sends it to the network.
// IMPORTANT: This method is only used for testing purposes.
func (c *PTBClient) FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (*suirpcv2.ExecuteTransactionResponse, error) {
	// This method should only be used in test environments
	if !testing.Testing() {
		return nil, errors.New("FinishPTBAndSend is only available in test environments")
	}

	gasPrice, err := c.GetReferenceGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get reference gas price: %w", err)
	}
	tx.SetGasPrice(gasPrice.Uint64())
	tx.SetSigner(txnSigner)
	tx.SetGasBudget(DefaultGasBudget)

	paymentCoinBytes, paymentCoinVersion, paymentCoinDigest, err := c.GetTransactionPaymentCoinForAddress(ctx, txnSigner.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction payment coin for address: %w", err)
	}

	tx.SetGasPayment([]transaction.SuiObjectRef{
		{
			ObjectId: paymentCoinBytes,
			Version:  paymentCoinVersion,
			Digest:   paymentCoinDigest,
		},
	})

	bcsBytes, err := tx.BuildBCSBytes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build bcs bytes: %w", err)
	}

	encodedBcsBytes := base64.StdEncoding.EncodeToString(bcsBytes)

	return c.SignAndSendTransaction(ctx, encodedBcsBytes, txnSigner.PubKey)
}

// GetBlockById (i.e. get checkpoint by id) returns the checkpoint details given its ID
func (c *PTBClient) GetBlockById(ctx context.Context, checkpointDigest string) (*suirpcv2.Checkpoint, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	var result *suirpcv2.Checkpoint
	err = c.WithRateLimit(ctx, "GetBlockById", func(ctx context.Context) error {
		response, cpErr := ledgerService.GetCheckpoint(ctx, &suirpcv2.GetCheckpointRequest{
			CheckpointId: &suirpcv2.GetCheckpointRequest_Digest{
				Digest: checkpointDigest,
			},
		})
		if cpErr != nil {
			return fmt.Errorf("failed to get checkpoint: %w", cpErr)
		}

		result = response.GetCheckpoint()
		return nil
	})

	return result, err
}

func (c *PTBClient) GetLatestEpoch(ctx context.Context) (*suirpcv2.Epoch, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get state service: %w", err)
	}

	var result *suirpcv2.Epoch
	err = c.WithRateLimit(ctx, "GetLatestEpoch", func(ctx context.Context) error {
		response, epochErr := ledgerService.GetEpoch(ctx, &suirpcv2.GetEpochRequest{
			Epoch: nil,
		})
		if epochErr != nil {
			return fmt.Errorf("failed to get latest epoch: %w", epochErr)
		}
		result = response.Epoch
		return nil
	})

	return result, err
}

func (c *PTBClient) GetCoinBalanceByAddress(ctx context.Context, address string, coinType string) (*suirpcv2.Balance, error) {
	stateService, err := c.getStateService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get state service: %w", err)
	}

	var result *suirpcv2.Balance
	err = c.WithRateLimit(ctx, "GetCoinBalanceByAddress", func(ctx context.Context) error {
		response, balanceErr := stateService.GetBalance(ctx, &suirpcv2.GetBalanceRequest{
			Owner:    &address,
			CoinType: &coinType,
		})
		if balanceErr != nil {
			return fmt.Errorf("failed to get coin balance: %w", balanceErr)
		}

		result = response.Balance

		return nil
	})

	return result, err
}

func (c *PTBClient) GetSUIBalance(ctx context.Context, address string) (*suirpcv2.Balance, error) {
	suiCoinType := "0x2::sui::SUI"
	return c.GetCoinBalanceByAddress(ctx, address, suiCoinType)
}

func (c *PTBClient) GetNormalizedModule(ctx context.Context, packageId string, module string) (models.GetNormalizedMoveModuleResponse, error) {
	var result models.GetNormalizedMoveModuleResponse
	err := c.WithRateLimit(ctx, "GetNormalizedModule", func(ctx context.Context) error {
		var err error
		result, err = c.getNormalizedModuleInternal(ctx, packageId, module)
		return err
	})
	return result, err
}

// getNormalizedModuleInternal is the internal implementation without rate limiting
func (c *PTBClient) getNormalizedModuleInternal(ctx context.Context, packageId string, module string) (models.GetNormalizedMoveModuleResponse, error) {
	// check if the normalized module is already cached
	normalizedModule, ok := c.normalizedModules[packageId][module]
	if ok {
		return normalizedModule, nil
	}

	if c.moveModuleClient == nil {
		return models.GetNormalizedMoveModuleResponse{}, errors.New("move module client not configured")
	}

	normalizedModule, err := c.moveModuleClient.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
		Package:    packageId,
		ModuleName: module,
	})
	if err != nil {
		return models.GetNormalizedMoveModuleResponse{}, fmt.Errorf("failed to get normalized module: %w", err)
	}

	if _, ok := c.normalizedModules[packageId]; !ok {
		c.normalizedModules[packageId] = make(map[string]models.GetNormalizedMoveModuleResponse)
	}

	// cache the normalized module
	c.normalizedModules[packageId][module] = normalizedModule

	return normalizedModule, nil
}

func (c *PTBClient) GetMoveModuleFunction(ctx context.Context, packageId string, moduleId string, functionName string) (*suirpcv2.FunctionDescriptor, error) {
	var result *suirpcv2.FunctionDescriptor
	err := c.WithRateLimit(ctx, "GetMoveModuleFunction", func(ctx context.Context) error {
		var err error
		movePkgService, err := c.getMovePackageService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get move package service: %w", err)
		}
		response, err := movePkgService.GetFunction(ctx, &suirpcv2.GetFunctionRequest{
			PackageId:  &packageId,
			ModuleName: &moduleId,
			Name:       &functionName,
		})
		result = response.GetFunction()
		return err
	})
	return result, err
}

func (c *PTBClient) GetCoinMetadata(ctx context.Context, coinType string) (models.CoinMetadataResponse, error) {
	var result models.CoinMetadataResponse
	err := c.WithRateLimit(ctx, "GetCoinMetadata", func(ctx context.Context) error {
		stateService, err := c.getStateService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get state service: %w", err)
		}

		coinInfoResponse, err := stateService.GetCoinInfo(ctx, &suirpcv2.GetCoinInfoRequest{
			CoinType: &coinType,
		})
		if err != nil {
			return fmt.Errorf("failed to get coin info: %w", err)
		}

		metadata := coinInfoResponse.GetMetadata()
		if metadata == nil {
			return fmt.Errorf("no metadata found for coin type %s", coinType)
		}

		result = models.CoinMetadataResponse{
			Id:          metadata.GetId(),
			Decimals:    int(metadata.GetDecimals()),
			Name:        metadata.GetName(),
			Symbol:      metadata.GetSymbol(),
			IconUrl:     metadata.GetIconUrl(),
			Description: metadata.GetDescription(),
		}
		return nil
	})

	return result, err
}

// LoadModulePackages returns the set of package IDs for a given module using its original package ID
// This method assumes that module names are unique across all packages
func (c *PTBClient) LoadModulePackageIds(ctx context.Context, packageId string, module string) ([]string, error) {
	var result []string
	err := c.WithRateLimit(ctx, "LoadModulePackageIds", func(ctx context.Context) error {
		var err error
		result, err = c.loadModulePackageIdsInternal(ctx, packageId, module)
		return err
	})
	return result, err
}

// loadModulePackageIdsInternal is the internal implementation without rate limiting
func (c *PTBClient) loadModulePackageIdsInternal(ctx context.Context, packageId string, module string) ([]string, error) {
	// TODO: this can be simplified significantly by getting the pointer config directly using the module's name.
	// If the pointer config does not exist for the module, we can return the provided package ID as it's the only package ID.
	// If the pointer config exists, we can use the parent field name and query the pointer object's ID directly.

	movePkgService, err := c.getMovePackageService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get move package service: %w", err)
	}

	addPackageIdFunctionName := "add_package_id"
	response, err := movePkgService.GetFunction(ctx, &suirpcv2.GetFunctionRequest{
		PackageId:  &packageId,
		ModuleName: &module,
		Name:       &addPackageIdFunctionName,
	})
	if err != nil || response.Function == nil {
		c.log.Warnw("module does not have the `add_package_id` function", "module", module)
		// fallback to using the provided package ID as it's the only package ID
		return []string{packageId}, nil
	}

	packageDetails, err := movePkgService.GetPackage(ctx, &suirpcv2.GetPackageRequest{
		PackageId: &packageId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get package details for package %s: %w", packageId, err)
	}

	// Iterate through the structs to find the pointer object
	var pointerStructTypeName string
	for _, module := range packageDetails.Package.Modules {
		for _, datatype := range module.GetDatatypes() {
			structTypeName := datatype.GetTypeName()
			if strings.Contains(structTypeName, "Pointer") {
				// Must be the full type name: <defining_id>::<module>::<name>
				// so we can use it to filter safely when querying owned objects
				pointerStructTypeName = structTypeName
				break
			}
		}
	}

	if pointerStructTypeName == "" {
		return nil, fmt.Errorf("pointer struct name not found for package %s and module %s", packageId, module)
	}

	// Read the owned objects to get the pointer object's ID (we use the internal method to avoid nested semaphore acquisition)
	pointerObjects, err := c.readFilterOwnedObjectIdsInternal(ctx, packageId, &pointerStructTypeName, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned objects: %w", err)
	}

	if len(pointerObjects.Objects) == 0 {
		return nil, fmt.Errorf("no pointer objects found for package %s and module %s", packageId, module)
	}

	// We assume a single pointer object per package and module
	structName := strings.Split(pointerStructTypeName, "::")[2]
	pointerObject := pointerObjects.Objects[0]
	parentObjectID := pointerObject.Json.GetStructValue().Fields[common.GetParentFieldName(structName)].GetStringValue()

	// The derivation key is the state object name within that module
	derivationKey := common.GetStateObjectNameByModule(module)

	stateObjectID, err := DeriveObjectIDWithVectorU8Key(parentObjectID, []byte(derivationKey))
	if err != nil {
		return nil, fmt.Errorf("failed to derive state object ID in LoadModulePackageIds: %w", err)
	}

	c.log.Debugw("stateObjectId", "stateObjectId", stateObjectID, "derivationKey", derivationKey)

	// Read the state object
	stateObject, err := c.readObjectIdInternal(ctx, stateObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state object: %w", err)
	}

	c.log.Debugw("stateObject", "stateObject", stateObject.GetJson())

	packageIdsValues := stateObject.Json.GetStructValue().GetFields()["package_ids"].GetListValue()
	packageIds := make([]string, len(packageIdsValues.Values))
	for i, packageIdValue := range packageIdsValues.Values {
		packageIds[i] = packageIdValue.GetStringValue()
	}

	return packageIds, nil
}

func (c *PTBClient) GetLatestPackageId(ctx context.Context, packageId string, module string) (string, error) {
	// attempt reading the value from cache first. Only non-empty entries are ever cached, so a
	// cache hit is always a valid package ID.
	cacheKey := "latest_pkg_id_fetch:" + packageId + ":" + module
	if cachedID, found := c.cache.Get(cacheKey); found {
		if id, ok := cachedID.(string); ok && id != "" {
			return id, nil
		}
	}

	missStart := time.Now()
	var result string
	err := c.WithRateLimit(ctx, "GetLatestPackageId", func(ctx context.Context) error {
		var err error
		result, err = c.getLatestPackageIdInternal(ctx, packageId, module)

		// Only cache successful, non-empty resolutions. Caching an empty result on error would
		// poison the cache: subsequent calls would return ("", nil) for the TTL window, causing
		// reads to target the zero package address (0x0). Package IDs are stable for the duration
		// of a CCIP deployment, so a longer TTL on success avoids repeated
		// GetFunction/GetPackage/ListOwnedObjects storms during config polling.
		if err == nil && result != "" {
			c.cache.Set(cacheKey, result, DefaultPackageIdCacheTTL)
		}

		return err
	})
	c.log.Debugw("GetLatestPackageId package id cache MISS", "module", module, "missSec", time.Since(missStart).Seconds(), "err", err != nil)
	return result, err
}

// getLatestPackageIdInternal is the internal implementation without rate limiting
func (c *PTBClient) getLatestPackageIdInternal(ctx context.Context, packageId string, module string) (string, error) {
	// Use internal method to avoid nested semaphore acquisition
	packageIds, err := c.loadModulePackageIdsInternal(ctx, packageId, module)
	if err != nil {
		return "", fmt.Errorf("failed to load module package ids: %w", err)
	}

	if len(packageIds) == 0 {
		return "", fmt.Errorf("nil or empty package ids found for package %s and module %s", packageId, module)
	}

	return packageIds[len(packageIds)-1], nil
}

func (c *PTBClient) GetCache() *cache.Cache {
	return c.cache
}

func (c *PTBClient) GetCachedValue(key string) (any, bool) {
	return c.cache.Get(key)
}

func (c *PTBClient) GetCachedValues(keys []string) (map[string]any, bool) {
	result := make(map[string]any)
	for _, key := range keys {
		value, found := c.cache.Get(key)
		if !found {
			return nil, false
		}
		result[key] = value
	}
	return result, true
}

func (c *PTBClient) SetCachedValue(key string, value any) {
	c.cache.Set(key, value, cache.NoExpiration)
}

func (c *PTBClient) SetCachedValues(keyValues map[string]any) {
	for key, value := range keyValues {
		c.cache.Set(key, value, cache.NoExpiration)
	}
}

// GetCCIPPackageId gets the CCIP package ID from the offramp package ID.
// IMPORTANT: This function expects to call the original (un-upgraded / first version) offramp package ID.
func (c *PTBClient) GetCCIPPackageID(ctx context.Context, offRampPackageID string) (string, error) {
	cacheKey := "ccip_package_id_fetch:" + offRampPackageID
	if cachedID, found := c.GetCachedValue(cacheKey); found {
		if id, ok := cachedID.(string); ok && id != "" {
			return id, nil
		}
	}

	response, err := c.ReadFunction(
		ctx,
		offRampPackageID,
		"offramp",
		"get_ccip_package_id",
		[]any{},
		[]string{},
		[]string{},
	)
	if err != nil {
		return "", err
	}

	ccipPackageID := response[0].(string)
	if ccipPackageID == "" {
		return "", fmt.Errorf("no CCIP package ID found for offramp package %s", offRampPackageID)
	}

	// Since the original CCIP package ID is the value we require here, we can cache it
	// without a specific TTL. Upgrades that change the latest package ID will be resolved
	// using `GetLatestPackageId` which will re-read the value from the chain and applies
	// a TTL if it uses a cache.
	c.SetCachedValue(cacheKey, ccipPackageID)
	return ccipPackageID, nil
}

// GetValueFromPackageOwnedObjectField gets the value of a field from a package owned object.
// This is used to get addresses stored within pointer objects on-chain. For example, the state object ID of a package is stored in the pointer object,
// so we need to get the value of the pointer object's field to get the state object ID.
func (c *PTBClient) GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error) {
	var result map[string]string
	err := c.WithRateLimit(ctx, "GetValuesFromPackageOwnedObjectField", func(ctx context.Context) error {
		var err error
		result, err = c.getValuesFromPackageOwnedObjectFieldInternal(ctx, packageID, moduleID, objectName, fieldKeys)
		return err
	})
	return result, err
}

// getValuesFromPackageOwnedObjectFieldInternal is the internal implementation without rate limiting
func (c *PTBClient) getValuesFromPackageOwnedObjectFieldInternal(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error) {
	structType := fmt.Sprintf("%s::%s::%s", packageID, moduleID, objectName)
	ownedObjects, err := c.readFilterOwnedObjectIdsInternal(ctx, packageID, &structType, nil)
	if err != nil {
		c.log.Errorw("Error reading owned objects", "error", err)
		return nil, err
	}

	if len(ownedObjects.Objects) == 0 {
		return nil, fmt.Errorf("no owned objects found for structType %s", structType)
	}

	foundValues := make(map[string]string)
	objectFields := ownedObjects.Objects[0].Json.GetStructValue().GetFields()
	for fieldKey, fieldValue := range objectFields {
		foundValues[fieldKey] = fieldValue.GetStringValue()
	}

	return foundValues, nil
}

// GetParentObjectID gets the parent object ID from a pointer object's field.
// With derived objects, pointers now store a reference to the parent "Object" struct (e.g., OffRampObject, CCIPObject).
// e.g. OffRampStatePointer contains "off_ramp_object_id" field pointing to OffRampObject.
func (c *PTBClient) GetParentObjectID(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error) {
	cacheKey := "parent_object_id_fetch:" + packageID + ":" + moduleID + ":" + pointerObjectName
	if cachedID, found := c.GetCachedValue(cacheKey); found {
		if id, ok := cachedID.(string); ok && id != "" {
			return id, nil
		}
	}

	var result string
	err := c.WithRateLimit(ctx, "GetParentObjectID", func(ctx context.Context) error {
		var err error
		result, err = c.getParentObjectIDInternal(ctx, packageID, moduleID, pointerObjectName)

		if err == nil && result != "" {
			c.SetCachedValue(cacheKey, result)
		}

		return err
	})
	return result, err
}

// getParentObjectIDInternal is the internal implementation without rate limiting
func (c *PTBClient) getParentObjectIDInternal(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error) {
	// Use internal method to avoid nested semaphore acquisition
	structType := fmt.Sprintf("%s::%s::%s", packageID, moduleID, pointerObjectName)
	ownedObjects, err := c.readFilterOwnedObjectIdsInternal(ctx, packageID, &structType, nil)
	if err != nil {
		c.log.Errorw("Error reading owned objects", "error", err)
		return "", err
	}

	if len(ownedObjects.Objects) == 0 {
		return "", fmt.Errorf("no pointer objects found for package %s and module %s", packageID, moduleID)
	}

	// We assume there is only one pointer object per package and module
	ownedObject := ownedObjects.Objects[0]

	// Get the parent field name from shared configuration
	fieldName := common.GetParentFieldName(pointerObjectName)
	if fieldName == "" {
		return "", fmt.Errorf("unknown pointer object type: %s", pointerObjectName)
	}

	ownedObjectStructValue := ownedObject.GetJson().GetStructValue()
	if ownedObjectStructValue == nil {
		return "", fmt.Errorf("pointer object %s does not have a struct value", pointerObjectName)
	}
	if ownedObjectStructValue.Fields[fieldName] == nil {
		return "", fmt.Errorf("pointer object %s does not have the field %s", pointerObjectName, fieldName)
	}

	c.log.Debugw("ownedObjectStructValue", "pointerFieldValue", ownedObjectStructValue.Fields[fieldName])

	parentObjectID := ownedObjectStructValue.Fields[fieldName].GetStringValue()
	if parentObjectID == "" {
		return "", fmt.Errorf("pointer object %s not found in package %s", structType, packageID)
	}

	return parentObjectID, nil
}

// GetCheckpointAvailability returns the provider's checkpoint history bounds from GetServiceInfo.
func (c *PTBClient) GetCheckpointAvailability(ctx context.Context) (*suirpcv2.GetServiceInfoResponse, error) {
	var result *suirpcv2.GetServiceInfoResponse

	err := c.WithRateLimit(ctx, "GetCheckpointAvailability", func(ctx context.Context) error {
		service, err := c.getLedgerService(ctx)
		if err != nil {
			return fmt.Errorf("failed to get ledger service: %w", err)
		}

		resp, err := service.GetServiceInfo(ctx, &suirpcv2.GetServiceInfoRequest{})
		if err != nil {
			return fmt.Errorf("GetServiceInfo failed: %w", err)
		}

		result = resp
		return nil
	})

	return result, err
}
