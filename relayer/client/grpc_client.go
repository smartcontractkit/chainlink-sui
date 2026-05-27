package client

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/common/grpcconn"
	"github.com/block-vision/sui-go-sdk/models"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	v2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/hasura/go-graphql-client"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"

	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	suiSigner "github.com/smartcontractkit/chainlink-sui/relayer/signer"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
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
)

var RateLimitWeights = map[string]int64{
	"MoveCall":                             1,
	"SendTransaction":                      1,
	"ReadFunction":                         1,
	"SignAndSendTransaction":               1,
	"QueryEvents":                          1,
	"QueryTransactions":                    1,
	"GetCoinsByAddress":                    1,
	"QueryCoinsByAddress":                  1,
	"EstimateGas":                          1,
	"GetTransactionStatus":                 1,
	"GetBlockById":                         1,
	"GetNormalizedModule":                  1,
	"GetSUIBalance":                        1,
	"GetCoinBalanceByAddress":              1,
	"GetValuesFromPackageOwnedObjectField": 1,
	"GetReferenceGasPrice":                 1,
	"FinishPTBAndSend":                     1,
	"BlockByDigest":                        1,
	// Keep 0, these methods are often called at the same time as ReadFunction
	// from ChainReader, high load of GetLatestValue calls could cause a deadlock.
	"ReadFilterOwnedObjectIds":           0,
	"ReadOwnedObjects":                   0,
	"ReadObjectId":                       0,
	"GetLatestPackageId":                 0,
	"LoadModulePackageIds":               0,
	"GetParentObjectID":                  0,
	"GetCCIPPackageID":                   0,
	"GetTokenPoolConfigByPackageAddress": 0,
	"GetLatestEpoch":                     0,
}

// var since it's passed via pointer
var maxPageSize uint = 50

type SuiPTBClient interface {
	MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error)
	SendTransaction(ctx context.Context, execRequest *suirpcv2.ExecuteTransactionRequest) (*suirpcv2.ExecuteTransactionResponse, error)
	ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor []byte) ([]*suirpcv2.Object, error)
	ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor []byte) ([]*suirpcv2.Object, error)
	ReadObjectId(ctx context.Context, objectId string) (*suirpcv2.Object, error)
	ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error)
	SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte) (*suirpcv2.ExecuteTransactionResponse, error)
	QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*models.PaginatedEventsResponse, error)
	QueryTransactions(ctx context.Context, fromAddress string, cursor *suirpcv2.Checkpoint, limit *uint64) ([]*suirpcv2.ExecutedTransaction, error)
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]*suirpcv2.Object, error)
	QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]*suirpcv2.Object, error)
	EstimateGas(ctx context.Context, tx *transaction.Transaction) (uint64, error)
	GetReferenceGasPrice(ctx context.Context) (*big.Int, error)
	FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (*suirpcv2.ExecuteTransactionResponse, error)
	BlockByDigest(ctx context.Context, txDigest string) (*suirpcv2.Checkpoint, error)
	GetBlockById(ctx context.Context, checkpointDigest string) (*suirpcv2.Checkpoint, error)
	GetLatestEpoch(ctx context.Context) (*suirpcv2.Epoch, error)
	GetNormalizedModule(ctx context.Context, packageId string, moduleId string) (models.GetNormalizedMoveModuleResponse, error)
	GetSUIBalance(ctx context.Context, address string) (*suirpcv2.Balance, error)
	LoadModulePackageIds(ctx context.Context, packageId string, module string) ([]string, error)
	GetLatestPackageId(ctx context.Context, packageId string, module string) (string, error)
	GetClient() sui.ISuiAPI
	GetCache() *cache.Cache
	GetCachedValue(key string) (any, bool)
	SetCachedValue(key string, value any)
	GetCachedValues(keys []string) (map[string]any, bool)
	SetCachedValues(keyValues map[string]any)
	HashTxBytes(txBytes []byte) []byte
	GetCCIPPackageID(ctx context.Context, offRampPackageID string) (string, error)
	GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error)
	GetParentObjectID(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error)
	GetTokenPoolConfigByPackageAddress(ctx context.Context, accountAddress string, tokenPoolAddress string, ccipPackageAddress string) (module_token_admin_registry.TokenConfig, error)
}

// PTBClient implements SuiClient interface using the blockvision SDK.
// During the gRPC migration, JSON-RPC (client) and gRPC (grpcClient) coexist:
// migrated methods use gRPC service accessors; others continue via JSON-RPC.
type PTBClient struct {
	log                 logger.Logger
	client              sui.ISuiAPI // TODO: remove this once the gRPC migration is complete
	graphqlClient       *graphql.Client
	grpcClient          *grpcconn.SuiGrpcClient
	grpcServicesMu      sync.Mutex
	ledgerService       suirpcv2.LedgerServiceClient
	stateService        suirpcv2.StateServiceClient
	txExecService       suirpcv2.TransactionExecutionServiceClient
	movePkgService      suirpcv2.MovePackageServiceClient
	subscriptionService suirpcv2.SubscriptionServiceClient
	maxRetries          *int
	transactionTimeout  time.Duration
	keystoreService     loop.Keystore
	rateLimiter         *semaphore.Weighted
	defaultRequestType  TransactionRequestType

	// map of module name to normalized module definition (similar to an ABI)
	normalizedModules map[string]map[string]models.GetNormalizedMoveModuleResponse

	cache *cache.Cache // used for caching object IDs (e.g. offramp state object ID or state pointers)
}

var _ SuiPTBClient = (*PTBClient)(nil)

func NewPTBClient(log logger.Logger, cfg PTBClientConfig) (*PTBClient, error) {
	return NewPTBClientFromConfig(log, cfg)
}

func (c *PTBClient) WithRateLimit(ctx context.Context, methodName string, f func(ctx context.Context) error) error {
	start := time.Now()

	weight := int64(1)
	if weightValue, ok := RateLimitWeights[methodName]; ok {
		weight = weightValue
	}

	workCtx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	// If rate limiter is disabled or weight is 0, skip semaphore entirely.
	// This will skip adding to the semaphore queue and prevent unnecessary queuing.
	if c.rateLimiter == nil || weight == 0 {
		return f(workCtx)
	}

	// acquire with the timeout context so it can't hang forever
	if err := c.rateLimiter.Acquire(ctx, weight); err != nil {
		return fmt.Errorf("failed to acquire rate limit for %s: %w", methodName, err)
	}

	// ensure cleanup on exit
	defer func() {
		c.rateLimiter.Release(weight)
		c.log.Debugw("WithRateLimit released", "methodName", methodName, "duration", time.Since(start))
	}()

	// run the user function with the timeout context
	// if the function respects the context, it will return and lock will be released in defer
	return f(workCtx)
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

	functionDefinition, err := movePkgService.GetFunction(ctx, &v2.GetFunctionRequest{
		PackageId:  &req.PackageObjectId,
		ModuleName: &req.Module,
		Name:       &req.Function,
	})
	if err != nil {
		return TxnMetaData{}, fmt.Errorf("failed to get normalized module: %w", err)
	}

	for i, arg := range req.Arguments {
		argType := functionDefinition.GetFunction().GetParameters()[i].Body.GetType()
		argTypeString := v2.OpenSignatureBody_Type_name[int32(argType)]
		arg, err := c.TransformTransactionArg(ctx, txn, arg, argTypeString, true)
		if err != nil {
			return TxnMetaData{}, fmt.Errorf("failed to transform transaction arg: %w", err)
		}
		args = append(args, *arg)
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
	err := c.WithRateLimit(ctx, "EstimateGas", func(ctx context.Context) error {
		bcsBytes, err := tx.BuildBCSBytes(ctx)
		if err != nil {
			return fmt.Errorf("failed to build bcs bytes: %w", err)
		}

		doGasSelection := true
		response, err := c.txExecService.SimulateTransaction(ctx, &suirpcv2.SimulateTransactionRequest{
			Transaction:    &v2.Transaction{Bcs: &v2.Bcs{Value: bcsBytes}},
			DoGasSelection: &doGasSelection,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{
					"transaction.effects.status",
					"transaction.effects.gas_used",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to simulate transaction: %w", err)
		}

		gasUsed := response.Transaction.Effects.GasUsed

		// Extract gas used from response
		var computationCost, storageCost, storageRebate uint64
		if gasUsed.ComputationCost != nil {
			computationCost = *gasUsed.ComputationCost
		}
		if gasUsed.StorageCost != nil {
			storageCost = *gasUsed.StorageCost
		}
		if gasUsed.StorageRebate != nil {
			storageRebate = *gasUsed.StorageRebate
		}

		// Override the estimate with a minimum threshold
		result = max(computationCost+storageCost-storageRebate, DefaultMinGasBudget)

		return nil
	})

	return result, err
}

func (c *PTBClient) GetReferenceGasPrice(ctx context.Context) (*big.Int, error) {
	var result *big.Int
	err := c.WithRateLimit(ctx, "GetReferenceGasPrice", func(ctx context.Context) error {
		resp, err := c.ledgerService.GetEpoch(ctx, &v2.GetEpochRequest{
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
		results, err = c.readFunctionInternal(ctx, packageId, module, function, args, argTypes, typeArgs)
		return err
	})
	return results, err
}

// readFunctionInternal is the internal implementation without rate limiting
func (c *PTBClient) readFunctionInternal(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	txExecService, err := c.getTransactionExecutionService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction execution service: %w", err)
	}

	var results []any
	txn := transaction.NewTransaction()

	var txnArgs []transaction.Argument
	var txnTypeArgs []transaction.TypeTag

	// Process type arguments
	for _, typeArg := range typeArgs {
		typeTag, err := c.CreateTypeTag(typeArg)
		if err != nil {
			return nil, fmt.Errorf("failed to create type tag for %s: %w", typeArg, err)
		}
		txnTypeArgs = append(txnTypeArgs, typeTag)
	}

	for i, arg := range args {
		argType, ok := common.ValueAt(argTypes, i)
		if !ok {
			argType = common.InferArgumentType(arg)
		}

		arg, err := c.TransformTransactionArg(ctx, txn, arg, argType, true)
		if err != nil {
			return nil, fmt.Errorf("failed to transform transaction arg: %w", err)
		}
		txnArgs = append(txnArgs, *arg)
	}

	// TODO: use a "read-only" signer for this operation instead of creating a new one
	// each time, it can be empty since gas selection is disabled
	seed := make([]byte, 32)
	if _, err := rand.Read(seed); err != nil {
		return nil, fmt.Errorf("failed to generate random seed: %w", err)
	}
	devInspectSigner := signer.NewSigner(seed)
	txn.SetSigner(devInspectSigner)

	txn.SetSender(models.SuiAddress(devInspectSigner.Address))
	txn.SetGasBudget(DefaultGasBudget)
	txn.SetGasPrice(DefaultGasPrice)
	txn.MoveCall(models.SuiAddress(packageId), module, function, txnTypeArgs, txnArgs)

	// Get transaction bytes
	bcsBytes, err := txn.BuildBCSBytes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build bcs bytes: %w", err)
	}

	doGasSelection := false
	response, err := txExecService.SimulateTransaction(ctx, &suirpcv2.SimulateTransactionRequest{
		Transaction:    &v2.Transaction{Bcs: &v2.Bcs{Value: bcsBytes}},
		DoGasSelection: &doGasSelection,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to simulate transaction: %w", err)
	}

	c.log.Debugw("ReadFunction RPC response", "RPC response", response, "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function))

	if len(response.CommandOutputs) == 0 {
		return nil, fmt.Errorf("no results from function call: %+v", response)
	}

	returnedValues := response.CommandOutputs[0].GetReturnValues()
	if len(returnedValues) == 0 {
		return nil, fmt.Errorf("no return values from function call: %+v", response)
	}

	results = make([]any, len(returnedValues))
	for i, returnedValue := range returnedValues {
		results[i] = returnedValue.Json.AsInterface()
	}

	c.log.Debugw("ReadFunction results", "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function), "results", results)

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
		Transaction: &v2.Transaction{Bcs: &v2.Bcs{Value: bcsBytes}},
		Signatures: []*v2.UserSignature{
			{
				Bcs: &v2.Bcs{Value: suiSignature},
			},
		},
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: []string{"digest", "effects.status", "effects.gas_used"},
		},
	})

	return resp, err
}

func (c *PTBClient) QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	return nil, fmt.Errorf("method implementation pending gRPC migration")
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
		response, err := ledgerService.GetTransaction(ctx, &suirpcv2.GetTransactionRequest{
			Digest: &digest,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{"effects.status", "effects.gas_used"},
			},
		})
		if err != nil {
			return err
		}

		var status string
		success := response.Transaction.GetEffects().GetStatus().GetSuccess()
		if success {
			status = "success"
		} else {
			status = "failure"
		}

		result = TransactionResult{
			Status: status,
			Error:  response.Transaction.GetEffects().GetStatus().GetError().String(),
		}

		return nil
	})

	return result, err
}

// QueryTransactions queries the transactions for a given address.
// @param fromAddress - the address to query transactions for
// @param cursor - a checkpoint ID to start from, if nil the latest checkpoint is used
// @param limit - the limit of transactions to return
// @return the transactions and an error if any
func (c *PTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *suirpcv2.Checkpoint, limit *uint64) ([]*suirpcv2.ExecutedTransaction, error) {
	return nil, fmt.Errorf("method implementation pending gRPC migration")
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
			Paths: []string{
				"checkpoint",
				"checkpoint.transactions",
				"checkpoint.transactions.effects",
				"checkpoint.transactions.events",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoint: %w", err)
	}

	return response.GetCheckpoint().GetTransactions(), nil
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
		return nil, fmt.Errorf("FinishPTBAndSend is only available in test environments")
	}

	gasPrice, err := c.GetReferenceGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get reference gas price: %w", err)
	}
	tx.SetGasPrice(gasPrice.Uint64())
	tx.SetSigner(txnSigner)
	tx.SetGasBudget(DefaultGasBudget)

	bcsBytes, err := tx.BuildBCSBytes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build bcs bytes: %w", err)
	}

	encodedBcsBytes := base64.StdEncoding.EncodeToString(bcsBytes)

	return c.SignAndSendTransaction(ctx, encodedBcsBytes, txnSigner.PubKey)
}

// BlockByDigest returns the transaction block using the transaction digest.
func (c *PTBClient) BlockByDigest(ctx context.Context, txDigest string) (*suirpcv2.Checkpoint, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	var result *suirpcv2.Checkpoint
	err = c.WithRateLimit(ctx, "BlockByDigest", func(ctx context.Context) error {
		response, err := ledgerService.GetTransaction(ctx, &suirpcv2.GetTransactionRequest{
			Digest: &txDigest,
			ReadMask: &fieldmaskpb.FieldMask{
				Paths: []string{"transaction.checkpoint"},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get transaction block: %w", err)
		}

		checkpointSeqNumber := response.Transaction.Checkpoint
		if checkpointSeqNumber == nil {
			return fmt.Errorf("checkpoint sequence number not found")
		}

		checkpoint, err := ledgerService.GetCheckpoint(ctx, &suirpcv2.GetCheckpointRequest{
			CheckpointId: &suirpcv2.GetCheckpointRequest_SequenceNumber{
				SequenceNumber: *checkpointSeqNumber,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get checkpoint: %w", err)
		}

		result = checkpoint.GetCheckpoint()

		return nil
	})

	return result, err
}

// GetBlockById (i.e. get checkpoint by id) returns the checkpoint details given its ID
func (c *PTBClient) GetBlockById(ctx context.Context, checkpointDigest string) (*suirpcv2.Checkpoint, error) {
	ledgerService, err := c.getLedgerService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger service: %w", err)
	}

	var result *suirpcv2.Checkpoint
	err = c.WithRateLimit(ctx, "GetBlockById", func(ctx context.Context) error {
		response, err := ledgerService.GetCheckpoint(ctx, &suirpcv2.GetCheckpointRequest{
			CheckpointId: &suirpcv2.GetCheckpointRequest_Digest{
				Digest: checkpointDigest,
			},
		})
		if err != nil {
			return fmt.Errorf("failed to get checkpoint: %w", err)
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
		response, err := ledgerService.GetEpoch(ctx, &suirpcv2.GetEpochRequest{
			Epoch: nil,
		})
		if err != nil {
			return fmt.Errorf("failed to get latest epoch: %w", err)
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
		response, err := stateService.GetBalance(ctx, &suirpcv2.GetBalanceRequest{
			Owner:    &address,
			CoinType: &coinType,
		})
		if err != nil {
			return fmt.Errorf("failed to get coin balance: %w", err)
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

	normalizedModule, err := c.client.SuiGetNormalizedMoveModule(ctx, models.GetNormalizedMoveModuleRequest{
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

	stateObjectID, err := bind.DeriveObjectIDWithVectorU8Key(parentObjectID, []byte(derivationKey))
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
	var result string
	err := c.WithRateLimit(ctx, "GetLatestPackageId", func(ctx context.Context) error {
		var err error
		result, err = c.getLatestPackageIdInternal(ctx, packageId, module)
		return err
	})
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

func (c *PTBClient) GetClient() sui.ISuiAPI {
	return c.client
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

	return response[0].(string), nil
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
	var result string
	err := c.WithRateLimit(ctx, "GetParentObjectID", func(ctx context.Context) error {
		var err error
		result, err = c.getParentObjectIDInternal(ctx, packageID, moduleID, pointerObjectName)
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

// A helper to abstract away having to provide the generic type of a token pool state. Requires a CCIP / StateObject package binding.
func (c *PTBClient) GetTokenPoolConfigByPackageAddress(ctx context.Context, accountAddress string, tokenPoolAddress string, ccipPackageAddress string) (module_token_admin_registry.TokenConfig, error) {
	devInspectSigner := suiSigner.NewDevInspectSigner(accountAddress)
	tokenAdminRegistry, err := module_token_admin_registry.NewTokenAdminRegistry(ccipPackageAddress, c.GetClient())
	if err != nil {
		return module_token_admin_registry.TokenConfig{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	// Obtain the CCIPObjectRef ID from the CCIP package
	ccipPointerConfigs := common.GetPointerConfigsByContract("ccip")
	if len(ccipPointerConfigs) == 0 {
		return module_token_admin_registry.TokenConfig{}, fmt.Errorf("ccip pointer config not found")
	}

	ccipPointerConfig := ccipPointerConfigs[0]

	var ccipObjectRefID string
	if cached, ok := c.GetCachedValue(ccipPointerConfig.ParentFieldName); ok {
		ccipObjectRefID = cached.(string)
	} else {
		// Use internal method to avoid nested semaphore acquisition
		ccipObjectID, err := c.getParentObjectIDInternal(ctx, ccipPackageAddress, "state_object", ccipPointerConfig.Pointer)
		if err != nil {
			return module_token_admin_registry.TokenConfig{}, fmt.Errorf("failed to get ccip parent object ID: %w", err)
		}

		ccipObjectRefID, err = bind.DeriveObjectIDWithVectorU8Key(ccipObjectID, []byte("CCIPObjectRef"))
		if err != nil {
			return module_token_admin_registry.TokenConfig{}, fmt.Errorf("failed to derive ccip object ref ID: %w", err)
		}

		c.SetCachedValue(ccipPointerConfig.ParentFieldName, ccipObjectRefID)
	}

	// Obtain the pool token metadata using the token pool package ID by calling into TokenAdminRegistry
	poolTokenMetadataAddress, err := tokenAdminRegistry.DevInspect().GetPoolLocalToken(ctx, &bind.CallOpts{
		WaitForExecution: true,
		Signer:           devInspectSigner,
	}, bind.Object{
		Id: ccipObjectRefID,
	}, tokenPoolAddress)

	if err != nil {
		return module_token_admin_registry.TokenConfig{}, fmt.Errorf("failed to get pool local token: %w", err)
	} else if poolTokenMetadataAddress == "" {
		return module_token_admin_registry.TokenConfig{}, fmt.Errorf("pool token metadata address not found")
	}

	// Obtain the token pool config using the token pool metadata address by calling into TokenAdminRegistry
	tokenPoolConfig, err := tokenAdminRegistry.DevInspect().GetTokenConfigStruct(ctx, &bind.CallOpts{
		WaitForExecution: true,
		Signer:           devInspectSigner,
	}, bind.Object{
		Id: ccipObjectRefID,
	}, poolTokenMetadataAddress)

	if err != nil {
		return module_token_admin_registry.TokenConfig{}, fmt.Errorf("failed to get token pool config: %w", err)
	} else if tokenPoolConfig.TokenType == "" || tokenPoolConfig.TokenPoolPackageId == "" {
		return module_token_admin_registry.TokenConfig{}, fmt.Errorf("failed to get token pool config: empty response fields")
	}

	return tokenPoolConfig, nil
}
