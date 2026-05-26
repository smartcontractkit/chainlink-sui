package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/common/grpcconn"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	v2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/semaphore"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"

	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	suiSigner "github.com/smartcontractkit/chainlink-sui/relayer/signer"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"
	"github.com/smartcontractkit/chainlink-sui/shared"
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
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (SuiTransactionBlockResponse, error)
	ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor []byte) ([]*suirpcv2.Object, error)
	ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, cursor []byte) ([]*suirpcv2.Object, error)
	ReadObjectId(ctx context.Context, objectId string) (*suirpcv2.Object, error)
	ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error)
	SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte, executionRequestType TransactionRequestType) (SuiTransactionBlockResponse, error)
	QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*models.PaginatedEventsResponse, error)
	QueryTransactions(ctx context.Context, fromAddress string, cursor *string, limit *uint64) (models.SuiXQueryTransactionBlocksResponse, error)
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error)
	QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]models.CoinData, error)
	EstimateGas(ctx context.Context, tx *transaction.Transaction) (uint64, error)
	GetReferenceGasPrice(ctx context.Context) (*big.Int, error)
	FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (SuiTransactionBlockResponse, error)
	BlockByDigest(ctx context.Context, txDigest string) (*SuiTransactionBlockResponse, error)
	GetBlockById(ctx context.Context, checkpointId string) (models.CheckpointResponse, error)
	GetLatestEpoch(ctx context.Context) (string, error)
	GetNormalizedModule(ctx context.Context, packageId string, moduleId string) (models.GetNormalizedMoveModuleResponse, error)
	GetSUIBalance(ctx context.Context, address string) (*big.Int, error)
	LoadModulePackageIds(ctx context.Context, packageId string, module string) ([]string, error)
	GetLatestPackageId(ctx context.Context, packageId string, module string) (string, error)
	GetClient() sui.ISuiAPI
	GetCache() *cache.Cache
	GetCachedValue(key string) (any, bool)
	SetCachedValue(key string, value any)
	GetCachedValues(keys []string) (map[string]any, bool)
	SetCachedValues(keyValues map[string]any)
	HashTxBytes(txBytes []byte) []byte
	GetCCIPPackageID(ctx context.Context, offRampPackageID string, signerAddress string) (string, error)
	GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error)
	GetParentObjectID(ctx context.Context, packageID string, moduleID string, pointerObjectName string) (string, error)
	GetTokenPoolConfigByPackageAddress(ctx context.Context, accountAddress string, tokenPoolAddress string, ccipPackageAddress string) (module_token_admin_registry.TokenConfig, error)
}

// PTBClient implements SuiClient interface using the blockvision SDK.
// During the gRPC migration, JSON-RPC (client) and gRPC (grpcClient) coexist:
// migrated methods use gRPC service accessors; others continue via JSON-RPC.
type PTBClient struct {
	log                logger.Logger
	grpcClient         *grpcconn.SuiGrpcClient
	grpcServicesMu     sync.Mutex
	ledgerService      suirpcv2.LedgerServiceClient
	stateService       suirpcv2.StateServiceClient
	txExecService      suirpcv2.TransactionExecutionServiceClient
	movePkgService     suirpcv2.MovePackageServiceClient
	maxRetries         *int
	transactionTimeout time.Duration
	keystoreService    loop.Keystore
	rateLimiter        *semaphore.Weighted
	defaultRequestType TransactionRequestType

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

func (c *PTBClient) MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error) {
	return TxnMetaData{}, errors.New("method removed in gRPC migration")
}

func (c *PTBClient) SendTransaction(ctx context.Context, payload TransactionBlockRequest) (SuiTransactionBlockResponse, error) {
	var result SuiTransactionBlockResponse
	err := c.WithRateLimit(ctx, "SendTransaction", func(ctx context.Context) error {
		return errors.New("method pending gRPC migration")
	})

	return result, err
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
	objectReq := suirpcv2.GetObjectRequest{
		ObjectId: &objectId,
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"object.contents",
				"object.object_type",
				"object.owner",
				"object.version",
				"object.digest",
			},
		},
	}

	response, err := c.ledgerService.GetObject(ctx, &objectReq)
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
	response, err := c.stateService.ListOwnedObjects(ctx, &suirpcv2.ListOwnedObjectsRequest{
		Owner:      &ownerAddress,
		ObjectType: structType,
		PageToken:  cursor,
		ReadMask: &fieldmaskpb.FieldMask{
			Paths: []string{
				"object.contents",
				"object.object_type",
				"object.owner",
				"object.version",
				"object.digest",
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

func (c *PTBClient) ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	var results []any
	err := c.WithRateLimit(ctx, "ReadFunction", func(ctx context.Context) error {
		var err error
		results, err = c.readFunctionInternal(ctx, signerAddress, packageId, module, function, args, argTypes, typeArgs)
		return err
	})
	return results, err
}

// readFunctionInternal is the internal implementation without rate limiting
func (c *PTBClient) readFunctionInternal(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string, typeArgs []string) ([]any, error) {
	var results []any
	txn := transaction.NewTransaction()

	var txnArgs []transaction.Argument
	var txnTypeArgs []transaction.TypeTag

	// Process type arguments
	for _, typeArg := range typeArgs {
		typeTag, err := c.createTypeTag(typeArg)
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

	txn.SetSuiClient(c.client.(*sui.Client))
	txn.SetSender(models.SuiAddress(signerAddress))
	txn.SetGasBudget(DefaultGasBudget)
	txn.SetGasPrice(DefaultGasPrice)
	txn.MoveCall(models.SuiAddress(packageId), module, function, txnTypeArgs, txnArgs)

	// Get transaction bytes
	bcsEncodedMsg, err := txn.Data.V1.Kind.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transaction: %w", err)
	}
	txBytes := mystenbcs.ToBase64(bcsEncodedMsg)

	// Use dev inspect for read-only function calls
	devInspectReq := models.SuiDevInspectTransactionBlockRequest{
		Sender:  signerAddress,
		TxBytes: txBytes,
	}

	response, err := c.client.SuiDevInspectTransactionBlock(ctx, devInspectReq)
	if err != nil || response.Effects.Status.Status != "success" {
		return nil, fmt.Errorf("failed to read function: %w (%s)", err, response.Effects.Status.Error)
	}

	c.log.Debugw("ReadFunction RPC response", "RPC response", response, "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function))

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results from function call: %+v", response)
	}

	resultsMarshalled, err := response.Results.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal results: %w", err)
	}
	var functionReadResponse []FunctionReadResponse
	err = json.Unmarshal(resultsMarshalled, &functionReadResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal results: %w", err)
	}

	results = make([]any, len(functionReadResponse[0].ReturnValues))

	// parse one or more results
	for i, returnedValue := range functionReadResponse[0].ReturnValues {
		returnedValue := returnedValue.([]any)
		structTag := returnedValue[1].(string)
		structPartsLen := 3

		// create a bcs decoder from the return value
		bcsBytes, err := codec.AnySliceToBytes(returnedValue[0].([]any))
		if err != nil {
			return nil, fmt.Errorf("failed to convert return value to bytes: %w", err)
		}
		bcsDecoder := bcs.NewDeserializer(bcsBytes)

		// This is a special case for Sui strings as they are represented as a struct with tag "0x1::string::String"
		// Since we use the tag to fetch the normalized module, it causes a failure since a module "string" does not exist.
		// We parse the value separately here. The BCS decoder is not actually needed for this case but we are already initializing it for the complex structs
		// so we can use it to read the string value.
		if structTag == "0x1::string::String" {
			strValue := bcsDecoder.ReadString()
			results[i] = strValue
			continue
		}

		// Handle vector<T> types. We need to distinguish between vectors of structs
		// (e.g. vector<0x123::module::MyStruct>) and vectors of primitive values
		// (e.g. vector<u8>, vector<u64>).
		//
		// - For vector<Struct>, fetch the normalized module metadata and decode each
		//   element into a JSON object using DecodeVectorOfStructs.
		// - For vector<primitive>, fall back to DecodeSuiPrimative to decode the raw
		//   values.
		// In both cases, run the result through ConvertBytesToHex so that any []byte
		// fields (vector<u8> especially) are hex-encoded instead of base64 when
		// marshalled to JSON.
		if strings.HasPrefix(structTag, "vector<") && strings.HasSuffix(structTag, ">") {
			// Extract inner type to determine if it's a vector of structs or primitives
			innerType := strings.TrimSuffix(strings.TrimPrefix(structTag, "vector<"), ">")
			innerStructParts := strings.Split(innerType, "::")

			// Check if inner type is a struct (3 parts: package::module::struct)
			if len(innerStructParts) == structPartsLen {
				// This is vector<Struct> - get normalized module for the inner struct
				innerPackageId := innerStructParts[0]
				innerModuleName := innerStructParts[1]

				normalizedModule, err := c.getNormalizedModuleInternal(ctx, innerPackageId, innerModuleName)
				if err != nil {
					return nil, fmt.Errorf("failed to get normalized module for vector struct: %w", err)
				}

				// Use the new DecodeVectorOfStructs function
				jsonResult, err := codec.DecodeVectorOfStructs(bcsDecoder, structTag, normalizedModule.Structs)
				if err != nil {
					return nil, fmt.Errorf("failed to decode vector of structs: %w", err)
				}

				results[i] = jsonResult
			} else {
				// This is vector<primitive> - use existing primitive vector handling
				primitive, err := codec.DecodeSuiPrimative(bcsDecoder, structTag)
				if err != nil {
					return nil, fmt.Errorf("failed to decode primitive vector: %w", err)
				}
				results[i] = primitive
			}
			continue
		}

		// Handle non-vector types (existing logic)
		structParts := strings.Split(structTag, "::")

		// if the response type is not a struct (primitive type), skip the result (keep it as is)
		if len(structParts) != structPartsLen {
			primitive, err := codec.DecodeSuiPrimative(bcsDecoder, structTag)
			if err != nil {
				return nil, fmt.Errorf("failed to decode primitive: %w", err)
			}
			// Normalize large ints to decimal strings to be LOOP/JSON friendly
			if structTag == "u128" || structTag == "u256" {
				switch v := primitive.(type) {
				case *big.Int:
					results[i] = v.String()
				case big.Int:
					vv := v
					results[i] = vv.String()
				default:
					results[i] = fmt.Sprint(v)
				}
			} else {
				results[i] = primitive
			}
		} else {
			// otherwise, get the normalized struct and attempt turning the result into JSON
			normalizedModule, err := c.getNormalizedModuleInternal(ctx, packageId, structParts[1])
			if err != nil {
				return nil, fmt.Errorf("failed to get normalized struct: %w", err)
			}

			jsonResult, err := codec.DecodeSuiStructToJSON(normalizedModule.Structs, structParts[2], bcsDecoder)
			if err != nil {
				return nil, fmt.Errorf("failed to parse struct into JSON: %w", err)
			}

			// convert any []uint8 fields to hex strings
			hexified := common.ConvertBytesToHex(jsonResult)
			results[i] = hexified
		}
	}

	c.log.Debugw("ReadFunction results", "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function), "results", results)

	return results, nil
}

func (c *PTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte, executionRequestType TransactionRequestType) (SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	signerId := fmt.Sprintf("%064x", signerPublicKey)

	txBytes, err := shared.DecodeBase64(txBytesRaw)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to decode tx bytes: %w", err)
	}

	// Hash the transaction bytes to include intent messages for Sui signing protocol
	txBytesToSign := c.HashTxBytes(txBytes)
	signature, err := c.keystoreService.Sign(ctx, signerId, txBytesToSign)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to sign tx: %w", err)
	}

	signaturesString := SerializeSuiSignature(signature, signerPublicKey)

	return c.SendTransaction(ctx, TransactionBlockRequest{
		TxBytes:     txBytesRaw,
		Signatures:  []string{signaturesString},
		RequestType: string(executionRequestType),
		Options: TransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
	})
}

func (c *PTBClient) QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*models.PaginatedEventsResponse, error) {
	var result *models.PaginatedEventsResponse
	err := c.WithRateLimit(ctx, "QueryEvents", func(ctx context.Context) error {
		limitVal := uint64(maxPageSize)
		if limit != nil {
			limitVal = uint64(*limit)
		}

		eventFilter := models.EventFilterByMoveEventType{
			MoveEventType: fmt.Sprintf("%s::%s::%s", filter.Package, filter.Module, filter.Event),
		}

		queryReq := models.SuiXQueryEventsRequest{
			SuiEventFilter:  eventFilter,
			Limit:           limitVal,
			DescendingOrder: sortOptions != nil && sortOptions.Descending,
		}

		if cursor != nil {
			queryReq.Cursor = cursor
		}

		c.log.Infow("querying events",
			"filter", queryReq.SuiEventFilter,
			"limit", queryReq.Limit,
			"descending", queryReq.DescendingOrder,
			"cursor", cursor,
		)

		response, err := c.client.SuiXQueryEvents(ctx, queryReq)
		if err != nil {
			return fmt.Errorf("failed to query events: %w", err)
		}

		result = &response

		return nil
	})

	return result, err
}

func (c *PTBClient) GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error) {
	var result TransactionResult
	err := c.WithRateLimit(ctx, "GetTransactionStatus", func(ctx context.Context) error {
		txReq := models.SuiGetTransactionBlockRequest{
			Digest: digest,
			Options: models.SuiTransactionBlockOptions{
				ShowEffects: true,
			},
		}

		response, err := c.client.SuiGetTransactionBlock(ctx, txReq)
		if err != nil {
			return err
		}

		result = TransactionResult{
			Status: response.Effects.Status.Status,
			Error:  response.Effects.Status.Error,
		}

		return nil
	})

	return result, err
}

func (c *PTBClient) QueryTransactions(ctx context.Context, fromAddress string, cursor *string, limit *uint64) (models.SuiXQueryTransactionBlocksResponse, error) {
	var result models.SuiXQueryTransactionBlocksResponse

	limitVal := uint64(maxPageSize)
	if limit != nil {
		limitVal = *limit
	}

	// if the cursor is empty, set it to nil to avoid RPC errors
	if cursor != nil && *cursor == "" {
		cursor = nil
	}

	err := c.WithRateLimit(ctx, "QueryTransactions", func(ctx context.Context) error {
		c.log.Debugw("Querying transactions", "fromAddress", fromAddress, "cursor", cursor, "limit", limitVal)

		txns, err := c.client.SuiXQueryTransactionBlocks(ctx, models.SuiXQueryTransactionBlocksRequest{
			SuiTransactionBlockResponseQuery: models.SuiTransactionBlockResponseQuery{
				TransactionFilter: map[string]any{
					"FromAddress": fromAddress,
				},
				Options: models.SuiTransactionBlockOptions{
					ShowInput:          true,
					ShowEffects:        true,
					ShowEvents:         false,
					ShowObjectChanges:  false,
					ShowBalanceChanges: false,
				},
			},
			Limit:  limitVal,
			Cursor: cursor,
		})
		if err != nil {
			return fmt.Errorf("failed to get account transactions: %w", err)
		}

		result = txns

		return nil
	})

	return result, err
}

func (c *PTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error) {
	var result []models.CoinData
	err := c.WithRateLimit(ctx, "GetCoinsByAddress", func(ctx context.Context) error {
		coinsReq := models.SuiXGetAllCoinsRequest{
			Owner: address,
			Limit: uint64(maxCoinsPageSize),
		}
		hasNextPage := true

		for hasNextPage {
			response, err := c.client.SuiXGetAllCoins(ctx, coinsReq)
			if err != nil {
				return fmt.Errorf("failed to get coins: %w", err)
			}

			result = append(result, response.Data...)

			hasNextPage = response.HasNextPage
			coinsReq.Cursor = response.NextCursor
		}

		return nil
	})

	return result, err
}

func (c *PTBClient) QueryCoinsByAddress(ctx context.Context, address string, coinType string) ([]models.CoinData, error) {
	var result []models.CoinData
	err := c.WithRateLimit(ctx, "QueryCoinsByAddress", func(ctx context.Context) error {
		coinsReq := models.SuiXGetCoinsRequest{
			Owner:    address,
			CoinType: coinType,
			Limit:    uint64(maxCoinsPageSize),
		}
		hasNextPage := true

		for hasNextPage {
			response, err := c.client.SuiXGetCoins(ctx, coinsReq)
			if err != nil {
				return fmt.Errorf("failed to get coins: %w", err)
			}

			result = append(result, response.Data...)

			hasNextPage = response.HasNextPage
			coinsReq.Cursor = response.NextCursor
		}

		return nil
	})

	return result, err
}

// FinishPTBAndSend finishes the PTB transaction and sends it to the network.
// IMPORTANT: This method is only used for testing purposes.
func (c *PTBClient) FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (SuiTransactionBlockResponse, error) {
	// This method should only be used in test environments
	if !testing.Testing() {
		return SuiTransactionBlockResponse{}, fmt.Errorf("FinishPTBAndSend is only available in test environments")
	}

	gasPrice, err := c.GetReferenceGasPrice(ctx)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to get reference gas price: %w", err)
	}
	tx.SetGasPrice(gasPrice.Uint64())

	tx.SetSigner(txnSigner)
	tx.SetGasBudget(DefaultGasBudget)

	// Set gas payment - use the first coin available for the signer
	coins, err := c.GetCoinsByAddress(ctx, txnSigner.Address)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to get coins for gas payment: %w", err)
	}
	if len(coins) == 0 {
		return SuiTransactionBlockResponse{}, fmt.Errorf("no coins available for gas payment")
	}
	// Use the first coin as gas payment
	paymentCoin, version, digest, err := c.GetTransactionPaymentCoinForAddress(ctx, txnSigner.Address)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to create coin object id: %w", err)
	}
	tx.SetGasPayment([]transaction.SuiObjectRef{
		{
			ObjectId: paymentCoin,
			Version:  version,
			Digest:   digest,
		},
	})

	c.log.Debugw("Executing transaction in PTB Client", "tx", tx)

	waitForLocal := requestType == WaitForLocalExecution
	execRequestType := string(requestType)
	if waitForLocal {
		execRequestType = string(WaitForEffectsCert)
	}

	response, err := tx.Execute(ctx, models.SuiTransactionBlockOptions{
		ShowInput:          true,
		ShowRawInput:       true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, execRequestType)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to execute transaction: %w", err)
	}

	// Same as SendTransaction: preserve err == nil for executed txs; callers use Status.

	if waitForLocal && response.Digest != "" {
		if waitErr := bind.WaitForTransactionIndexed(ctx, c.client, response.Digest); waitErr != nil {
			return c.convertBlockvisionResponse(response), waitErr
		}
	}

	return c.convertBlockvisionResponse(response), nil
}

func (c *PTBClient) BlockByDigest(ctx context.Context, txDigest string) (*SuiTransactionBlockResponse, error) {
	var result *SuiTransactionBlockResponse
	err := c.WithRateLimit(ctx, "BlockByDigest", func(ctx context.Context) error {
		txReq := models.SuiGetTransactionBlockRequest{
			Digest: txDigest,
			Options: models.SuiTransactionBlockOptions{
				ShowInput:          true,
				ShowEffects:        true,
				ShowEvents:         true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
			},
		}

		response, err := c.client.SuiGetTransactionBlock(ctx, txReq)
		if err != nil {
			return fmt.Errorf("failed to get transaction block: %w", err)
		}

		converted := c.convertBlockvisionResponse(&response)
		result = &converted

		return nil
	})

	return result, err
}

// GetBlockById (i.e. get checkpoint by id) returns the checkpoint details given its ID
func (c *PTBClient) GetBlockById(ctx context.Context, checkpointId string) (models.CheckpointResponse, error) {
	var result models.CheckpointResponse
	err := c.WithRateLimit(ctx, "GetBlockById", func(ctx context.Context) error {
		response, err := c.client.SuiGetCheckpoint(ctx, models.SuiGetCheckpointRequest{
			CheckpointID: checkpointId,
		})
		if err != nil {
			return fmt.Errorf("failed to get checkpoint: %w", err)
		}

		result = response

		return nil
	})

	return result, err
}

func (c *PTBClient) GetLatestEpoch(ctx context.Context) (string, error) {
	var result string
	err := c.WithRateLimit(ctx, "GetLatestEpoch", func(ctx context.Context) error {
		response, err := c.client.SuiXGetLatestSuiSystemState(ctx)
		if err != nil {
			return fmt.Errorf("failed to get latest epoch: %w", err)
		}
		result = response.Epoch
		return nil
	})
	return result, err
}

func (c *PTBClient) GetSUIBalance(ctx context.Context, address string) (*big.Int, error) {
	var result *big.Int
	err := c.WithRateLimit(ctx, "GetSUIBalance", func(ctx context.Context) error {
		balanceReq := models.SuiXGetBalanceRequest{
			Owner:    address,
			CoinType: "0x2::sui::SUI", // Default SUI coin type
		}

		response, err := c.client.SuiXGetBalance(ctx, balanceReq)
		if err != nil {
			return fmt.Errorf("failed to get SUI balance: %w", err)
		}

		balance, ok := new(big.Int).SetString(response.TotalBalance, Base10)
		if !ok {
			return fmt.Errorf("failed to parse balance: %s", response.TotalBalance)
		}
		result = balance

		return nil
	})

	return result, err
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

// TODO: replace w/ GraphQL query
// loadModulePackageIdsInternal is the internal implementation without rate limiting
func (c *PTBClient) loadModulePackageIdsInternal(ctx context.Context, packageId string, module string) ([]string, error) {
	// Ensure that the module keeps track of its package IDs by checking that it has `add_package_id` function
	normalizedModule, err := c.getNormalizedModuleInternal(ctx, packageId, module)
	if err != nil {
		return nil, fmt.Errorf("failed to get normalized module: %w", err)
	}

	// TODO: this is not a clean way to determine the if the package holds a history of its IDs,
	// we need to replace this with a GraphQL query to come up with a more robust solution.
	// Check that the module has the `add_package_id` function
	if _, ok := normalizedModule.ExposedFunctions["add_package_id"]; !ok {
		c.log.Warnw("module does not have the `add_package_id` function", "module", module)
		// fallback to using the provided package ID as it's the only package ID
		return []string{packageId}, nil
	}

	// Iterate through the structs to find the pointer object
	pointerStructName := ""
	pointerStructNameFound := false
	for structName := range normalizedModule.Structs {
		if strings.Contains(structName, "Pointer") {
			pointerStructName = structName
			pointerStructNameFound = true
			break
		}
	}

	if !pointerStructNameFound {
		return nil, fmt.Errorf("pointer struct name not found for package %s and module %s", packageId, module)
	}

	// Read the owned objects to get the pointer object's ID (we use the internal method to avoid nested semaphore acquisition)
	objectType := fmt.Sprintf("%s::%s::%s", packageId, module, pointerStructName)
	pointerObjects, err := c.readFilterOwnedObjectIdsInternal(ctx, packageId, &objectType, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned objects: %w", err)
	}

	if len(pointerObjects.Objects) == 0 {
		return nil, fmt.Errorf("no pointer objects found for package %s and module %s", packageId, module)
	}

	// Use the normalized module to determine the state ref object
	pointerObject := pointerObjects.Objects[0]
	c.log.Debugw("pointer ref object", "pointerObject", pointerObject)
	parentObjectID := pointerObject.Json.GetStructValue().Fields[common.GetParentFieldName(pointerStructName)].GetStringValue()

	c.log.Debugw("parentObjectID", "parentObjectID", parentObjectID)

	// TODO: put this in the config instead of having a match statement here
	derivationKey := ""
	switch module {
	case "offramp":
		derivationKey = "OffRampState"
	case "onramp":
		derivationKey = "OnRampState"
	case "ccip":
	case "state_object":
		derivationKey = "CCIPObjectRef"
	case "router":
		derivationKey = "RouterState"
	case "burn_mint_token_pool":
		derivationKey = "BurnMintTokenPoolState"
	case "lock_release_token_pool":
		derivationKey = "LockReleaseTokenPoolState"
	case "managed_token_pool":
		derivationKey = "ManagedTokenPoolState"
	case "usdc_token_pool":
		derivationKey = "USDCTokenPoolState"
	case "counter":
		derivationKey = "Counter"
	}

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
func (c *PTBClient) GetCCIPPackageID(ctx context.Context, offRampPackageID string, signerAddress string) (string, error) {
	offRamp, err := module_offramp.NewOfframp(offRampPackageID, c.GetClient())
	if err != nil {
		return "", err
	}

	devInspectSigner := suiSigner.NewDevInspectSigner(signerAddress)

	ccipPkgID, err := offRamp.DevInspect().GetCcipPackageId(ctx, &bind.CallOpts{
		Signer:           devInspectSigner,
		WaitForExecution: true,
	})
	if err != nil {
		return "", err
	}

	return ccipPkgID, nil
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

	parentObjectID := ownedObject.Json.GetStructValue().Fields[fieldName].GetStringValue()
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

// Add helper method to create type tags
func (c *PTBClient) createTypeTag(typeStr string) (transaction.TypeTag, error) {
	if typeStr == "" {
		return transaction.TypeTag{}, fmt.Errorf("type string cannot be empty")
	}

	// Handle struct types (package::module::name)
	if strings.Contains(typeStr, "::") {
		parts := strings.Split(typeStr, "::")
		if len(parts) != 3 {
			return transaction.TypeTag{}, fmt.Errorf("invalid struct type format %q, expected package::module::name", typeStr)
		}

		packageID, module, name := parts[0], parts[1], parts[2]

		// Convert package ID to address bytes
		packageAddr := models.SuiAddress(packageID)
		addressBytes, err := transaction.ConvertSuiAddressStringToBytes(packageAddr)
		if err != nil {
			return transaction.TypeTag{}, fmt.Errorf("failed to convert package address %q: %w", packageID, err)
		}

		return transaction.TypeTag{
			Struct: &transaction.StructTag{
				Address:    *addressBytes,
				Module:     module,
				Name:       name,
				TypeParams: []*transaction.TypeTag{},
			},
		}, nil
	}

	// TODO: Handle primitive types if needed
	return transaction.TypeTag{}, fmt.Errorf("unsupported type format: %s", typeStr)
}
