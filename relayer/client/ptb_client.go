package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/semaphore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
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
	DefaultCacheExpiration      time.Duration = 120 * time.Minute
	DefaultCacheCleanupInterval time.Duration = 240 * time.Minute
	DefaultHTTPTimeout          time.Duration = 15 * time.Second
)

// var since it's passed via pointer
var maxPageSize uint = 50

type SuiPTBClient interface {
	MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error)
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (SuiTransactionBlockResponse, error)
	ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor *models.ObjectId) ([]models.SuiObjectResponse, error)
	ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, limit *uint) ([]models.SuiObjectData, error)
	ReadObjectId(ctx context.Context, objectId string) (models.SuiObjectData, error)
	ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string) ([]any, error)
	SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte, executionRequestType TransactionRequestType) (SuiTransactionBlockResponse, error)
	QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*models.PaginatedEventsResponse, error)
	QueryTransactions(ctx context.Context, fromAddress string, cursor *string, limit *uint64) (models.SuiXQueryTransactionBlocksResponse, error)
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error)
	EstimateGas(ctx context.Context, txBytes string) (uint64, error)
	FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (SuiTransactionBlockResponse, error)
	BlockByDigest(ctx context.Context, txDigest string) (*SuiTransactionBlockResponse, error)
	GetBlockById(ctx context.Context, checkpointId string) (models.CheckpointResponse, error)
	GetNormalizedModule(ctx context.Context, packageId string, moduleId string) (models.GetNormalizedMoveModuleResponse, error)
	GetSUIBalance(ctx context.Context, address string) (*big.Int, error)
	LoadModulePackageIds(ctx context.Context, packageId string, module string, signerAddress string) ([]string, error)
	GetLatestPackageId(ctx context.Context, packageId string, module string, signerAddress string) (string, error)
	GetClient() sui.ISuiAPI
	GetCache() *cache.Cache
	GetCachedValue(key string) (any, bool)
	SetCachedValue(key string, value any)
	GetCachedValues(keys []string) (map[string]any, bool)
	SetCachedValues(keyValues map[string]any)
	HashTxBytes(txBytes []byte) []byte
	GetCCIPPackageID(ctx context.Context, offRampPackageID string, signerAddress string) (string, error)
	GetValuesFromPackageOwnedObjectField(ctx context.Context, packageID string, moduleID string, objectName string, fieldKeys []string) (map[string]string, error)
}

// PTBClient implements SuiClient interface using the blockvision SDK
type PTBClient struct {
	log                logger.Logger
	client             sui.ISuiAPI
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

func NewPTBClient(
	log logger.Logger,
	rpcUrl string,
	maxRetries *int,
	transactionTimeout time.Duration,
	keystoreService loop.Keystore,
	maxConcurrentRequests int64,
	defaultRequestType TransactionRequestType,
) (*PTBClient, error) {
	log.Infof("Creating new SUI client with blockvision SDK")

	httpClient := &http.Client{
		Timeout: DefaultHTTPTimeout,
	}
	client := sui.NewSuiClientWithCustomClient(rpcUrl, httpClient)

	if maxConcurrentRequests <= 0 {
		maxConcurrentRequests = 100 // Default value
	}

	return &PTBClient{
		log:                log,
		client:             client,
		maxRetries:         maxRetries,
		transactionTimeout: transactionTimeout,
		keystoreService:    keystoreService,
		rateLimiter:        semaphore.NewWeighted(maxConcurrentRequests),
		defaultRequestType: defaultRequestType,
		normalizedModules:  make(map[string]map[string]models.GetNormalizedMoveModuleResponse),
		cache:              cache.New(DefaultCacheExpiration, DefaultCacheCleanupInterval),
	}, nil
}

func (c *PTBClient) WithRateLimit(ctx context.Context, f func(ctx context.Context) error) error {
	// the work itself should time out by transactionTimeout
	timeoutCtx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	if c.rateLimiter == nil {
		return f(timeoutCtx)
	}

	// acquire with the timeout context so it can't hang forever
	if err := c.rateLimiter.Acquire(timeoutCtx, 1); err != nil {
		return fmt.Errorf("failed to acquire rate limit: %w", err)
	}

	// force the release if the function never returns
	var released int32 // 0 = not released, 1 = released
	releaseOnce := func() {
		if atomic.CompareAndSwapInt32(&released, 0, 1) {
			c.rateLimiter.Release(1)
		}
	}

	// if we haven't returned by the deadline (+ small grace),
	// reclaim the permit to prevent global deadlock
	grace := 2 * time.Second
	_ = time.AfterFunc(c.transactionTimeout+grace, func() {
		c.log.Warnw("rate limit lease expired; forcing release",
			"timeout", c.transactionTimeout.String())
		releaseOnce()
	})

	// ensure cleanup on exit and panic recovery
	defer func() {
		// NOTE: we can check time.Stop() to see if the timer fired but this redundant
		// since we are always releasing

		// release if not already reclaimed by watchdog
		releaseOnce()

		// bubble up panic after releasing
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	// run the user function with the timeout context
	// if the function respects the context, it will return and lock will be released in defer
	return f(timeoutCtx)
}

func (c *PTBClient) MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error) {
	var result TxnMetaData
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		moveCallReq := models.MoveCallRequest{
			Signer:          req.Signer,
			PackageObjectId: req.PackageObjectId,
			Module:          req.Module,
			Function:        req.Function,
			// TODO: handle type arguments
			TypeArguments: []any{},
			Arguments:     req.Arguments,
			GasBudget:     strconv.FormatUint(req.GasBudget, 10),
			Gas:           nil,
			ExecutionMode: models.TransactionExecutionCommit,
		}

		c.log.Debugw("MoveCall request", "request", moveCallReq)

		response, err := c.client.MoveCall(ctx, moveCallReq)
		if err != nil {
			return fmt.Errorf("failed to create move call: %w", err)
		}

		result.TxBytes = response.TxBytes

		return nil
	})

	return result, err
}

func (c *PTBClient) SendTransaction(ctx context.Context, payload TransactionBlockRequest) (SuiTransactionBlockResponse, error) {
	var result SuiTransactionBlockResponse
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		// Use blockvision SDK's execute transaction
		executeReq := models.SuiExecuteTransactionBlockRequest{
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

		c.log.Debugw("Executing transaction", "request", executeReq)

		response, err := c.client.SuiExecuteTransactionBlock(ctx, executeReq)
		if err != nil {
			return fmt.Errorf("failed to execute transaction: %w", err)
		}

		// Convert blockvision response to models response
		result = c.convertBlockvisionResponse(&response)

		return nil
	})

	return result, err
}

func (c *PTBClient) ReadObjectId(ctx context.Context, objectId string) (models.SuiObjectData, error) {
	var result models.SuiObjectData
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		objectReq := models.SuiGetObjectRequest{
			ObjectId: objectId,
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
				ShowType:    true,
				ShowOwner:   true,
			},
		}

		response, err := c.client.SuiGetObject(ctx, objectReq)
		if err != nil {
			return fmt.Errorf("failed to read object: %w", err)
		}

		c.log.Infow("ReadObjectId response", "response", response)

		if response.Data == nil || response.Data.Content == nil {
			return fmt.Errorf("object has no content")
		}

		result = *response.Data

		return nil
	})

	return result, err
}

func (c *PTBClient) ReadFilterOwnedObjectIds(ctx context.Context, ownerAddress string, structType string, limit *uint) ([]models.SuiObjectData, error) {
	var result []models.SuiObjectData
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		limitVal := uint64(maxPageSize)
		if limit != nil {
			limitVal = uint64(*limit)
		}

		ownedObjectsReq := models.SuiXGetOwnedObjectsRequest{
			Address: ownerAddress,
			Query: models.SuiObjectResponseQuery{
				Filter: models.ObjectFilterByStructType{
					StructType: structType,
				},
				Options: models.SuiObjectDataOptions{
					ShowType: true,
				},
			},
			Limit: limitVal,
		}

		response, err := c.client.SuiXGetOwnedObjects(ctx, ownedObjectsReq)
		if err != nil {
			return fmt.Errorf("failed to read owned objects: %w", err)
		}

		for _, obj := range response.Data {
			if obj.Data != nil {
				result = append(result, *obj.Data)
			}
		}

		return nil
	})

	return result, err
}

func (c *PTBClient) ReadOwnedObjects(ctx context.Context, ownerAddress string, cursor *models.ObjectId) ([]models.SuiObjectResponse, error) {
	var result []models.SuiObjectResponse
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		ownedObjectsReq := models.SuiXGetOwnedObjectsRequest{
			Address: ownerAddress,
			Query: models.SuiObjectResponseQuery{
				Options: models.SuiObjectDataOptions{
					ShowContent: true,
					ShowType:    true,
					ShowOwner:   true,
				},
			},
			Limit: uint64(maxPageSize),
		}

		if cursor != nil {
			cursorHex := cursor
			ownedObjectsReq.Cursor = string(cursorHex.Data())
		}

		response, err := c.client.SuiXGetOwnedObjects(ctx, ownedObjectsReq)
		if err != nil {
			return fmt.Errorf("failed to read owned objects: %w", err)
		}

		result = response.Data

		return nil
	})

	return result, err
}

func (c *PTBClient) EstimateGas(ctx context.Context, txBytes string) (uint64, error) {
	var result uint64
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		// Use blockvision SDK's dry run transaction
		dryRunReq := models.SuiDryRunTransactionBlockRequest{
			TxBytes: txBytes,
		}

		response, err := c.client.SuiDryRunTransactionBlock(ctx, dryRunReq)
		if err != nil {
			return fmt.Errorf("failed to estimate gas: %w", err)
		}

		// Extract gas used from response
		if response.Effects.GasUsed.ComputationCost != "" {
			computationCost, err := strconv.ParseUint(response.Effects.GasUsed.ComputationCost, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse computation cost: %w", err)
			}
			storageCost, err := strconv.ParseUint(response.Effects.GasUsed.StorageCost, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse storage cost: %w", err)
			}
			result = computationCost + storageCost
		}

		return nil
	})

	return result, err
}

func (c *PTBClient) ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string) ([]any, error) {
	var results []any
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		txn := transaction.NewTransaction()

		var txnArgs []transaction.Argument
		var txnTypeArgs []transaction.TypeTag
		for i, arg := range args {
			argType, ok := common.ValueAt(argTypes, i)
			if !ok {
				argType = common.InferArgumentType(arg)
			}

			arg, err := c.TransformTransactionArg(ctx, txn, arg, argType, true)
			if err != nil {
				return fmt.Errorf("failed to transform transaction arg: %w", err)
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
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}
		txBytes := mystenbcs.ToBase64(bcsEncodedMsg)

		// Use dev inspect for read-only function calls
		devInspectReq := models.SuiDevInspectTransactionBlockRequest{
			Sender:  signerAddress,
			TxBytes: txBytes,
		}

		response, err := c.client.SuiDevInspectTransactionBlock(ctx, devInspectReq)
		if err != nil && response.Effects.Status.Status != "success" {
			return fmt.Errorf("failed to read function: %w", err)
		}

		c.log.Debugw("ReadFunction RPC response", "RPC response", response, "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function))

		if len(response.Results) == 0 {
			return fmt.Errorf("no results from function call: %+v", response)
		}

		resultsMarshalled, err := response.Results.MarshalJSON()
		if err != nil {
			return fmt.Errorf("failed to marshal results: %w", err)
		}
		var functionReadResponse []FunctionReadResponse
		err = json.Unmarshal(resultsMarshalled, &functionReadResponse)
		if err != nil {
			return fmt.Errorf("failed to unmarshal results: %w", err)
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
				return fmt.Errorf("failed to convert return value to bytes: %w", err)
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

					normalizedModule, err := c.GetNormalizedModule(ctx, innerPackageId, innerModuleName)
					if err != nil {
						return fmt.Errorf("failed to get normalized module for vector struct: %w", err)
					}

					c.log.Debugw("normalizedModule for vector", "normalizedModule", normalizedModule)

					// Use the new DecodeVectorOfStructs function
					jsonResult, err := codec.DecodeVectorOfStructs(bcsDecoder, structTag, normalizedModule.Structs)
					if err != nil {
						return fmt.Errorf("failed to decode vector of structs: %w", err)
					}

					results[i] = jsonResult
				} else {
					// This is vector<primitive> - use existing primitive vector handling
					primitive, err := codec.DecodeSuiPrimative(bcsDecoder, structTag)
					if err != nil {
						return fmt.Errorf("failed to decode primitive vector: %w", err)
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
					return fmt.Errorf("failed to decode primitive: %w", err)
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
				normalizedModule, err := c.GetNormalizedModule(ctx, packageId, structParts[1])
				c.log.Debugw("normalizedModule", "normalizedModule", normalizedModule)
				if err != nil {
					return fmt.Errorf("failed to get normalized struct: %w", err)
				}

				jsonResult, err := codec.DecodeSuiStructToJSON(normalizedModule.Structs, structParts[2], bcsDecoder)
				if err != nil {
					return fmt.Errorf("failed to parse struct into JSON: %w", err)
				}

				// convert any []uint8 fields to hex strings
				hexified := common.ConvertBytesToHex(jsonResult)
				results[i] = hexified
			}
		}

		c.log.Debugw("ReadFunction results", "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function), "results", results)

		return nil
	})

	return results, err
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
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
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
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
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

	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
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
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		coinsReq := models.SuiXGetAllCoinsRequest{
			Owner: address,
			Limit: uint64(maxCoinsPageSize),
		}

		response, err := c.client.SuiXGetAllCoins(ctx, coinsReq)
		if err != nil {
			return fmt.Errorf("failed to get coins: %w", err)
		}

		result = response.Data

		return nil
	})

	return result, err
}

func (c *PTBClient) FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (SuiTransactionBlockResponse, error) {
	tx.SetSigner(txnSigner)
	// TODO: get gas price and budget from the txn
	tx.SetGasPrice(DefaultGasPrice)
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

	response, err := tx.Execute(ctx, models.SuiTransactionBlockOptions{
		ShowInput:          true,
		ShowRawInput:       true,
		ShowEffects:        true,
		ShowEvents:         true,
		ShowObjectChanges:  true,
		ShowBalanceChanges: true,
	}, string(requestType))
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to execute transaction: %w", err)
	}

	return c.convertBlockvisionResponse(response), nil
}

func (c *PTBClient) BlockByDigest(ctx context.Context, txDigest string) (*SuiTransactionBlockResponse, error) {
	var result *SuiTransactionBlockResponse
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
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
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
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

func (c *PTBClient) GetSUIBalance(ctx context.Context, address string) (*big.Int, error) {
	var result *big.Int
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
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
func (c *PTBClient) LoadModulePackageIds(ctx context.Context, packageId string, module string, signerAddress string) ([]string, error) {
	// Ensure that the module keeps track of its package IDs by checking that it has `add_package_id` function
	normalizedModule, err := c.GetNormalizedModule(ctx, packageId, module)
	if err != nil {
		return nil, fmt.Errorf("failed to get normalized module: %w", err)
	}

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

	// Read the owned objects to get the pointer object's ID
	ownedObjects, err := c.ReadOwnedObjects(ctx, packageId, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned objects: %w", err)
	}

	// Use the normalized module to determine the state ref object
	var pointerObject models.SuiObjectData
	pointerObjectFound := false

	for _, ownedObject := range ownedObjects {
		c.log.Debugw("ownedObject", "ownedObject", ownedObject.Data.Type)
		if ownedObject.Data.Type == fmt.Sprintf("%s::%s::%s", packageId, module, pointerStructName) {
			pointerObject = *ownedObject.Data
			pointerObjectFound = true
			break
		}
	}

	if !pointerObjectFound {
		return nil, fmt.Errorf("pointer object not found for package %s and module %s", packageId, module)
	}

	c.log.Debugw("pointer ref object", "pointerObject", pointerObject)

	stateObjectId := ""
	switch module {
	case "offramp":
		stateObjectId = pointerObject.Content.SuiMoveObject.Fields["off_ramp_state_id"].(string)
	case "onramp":
		stateObjectId = pointerObject.Content.SuiMoveObject.Fields["on_ramp_state_id"].(string)
	case "ccip":
	case "state_object":
		stateObjectId = pointerObject.Content.SuiMoveObject.Fields["object_ref_id"].(string)
	}

	if stateObjectId == "" {
		return nil, fmt.Errorf("state object id not found for package %s and module %s", packageId, module)
	}

	// Read the state object
	stateObject, err := c.ReadObjectId(ctx, stateObjectId)
	if err != nil {
		return nil, fmt.Errorf("failed to get state object: %w", err)
	}

	// Read the package IDs from the state object
	packageIds := []string{}
	for _, packageId := range stateObject.Content.SuiMoveObject.Fields["package_ids"].([]any) {
		packageIds = append(packageIds, packageId.(string))
	}

	return packageIds, nil
}

func (c *PTBClient) GetLatestPackageId(ctx context.Context, packageId string, module string, signerAddress string) (string, error) {
	packageIds, err := c.LoadModulePackageIds(ctx, packageId, module, signerAddress)
	if err != nil {
		return "", fmt.Errorf("failed to load module package ids: %w", err)
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
	ownedObjects, err := c.ReadOwnedObjects(ctx, packageID, nil)
	if err != nil {
		c.log.Errorw("Error reading owned objects", "error", err)
		return nil, err
	}

	foundValues := make(map[string]string)
	for _, ownedObject := range ownedObjects {
		qualifiedName := fmt.Sprintf("%s::%s::%s", packageID, moduleID, objectName)
		if ownedObject.Data.Type != "" && ownedObject.Data.Type == qualifiedName {
			// parse the object into a map
			parsedObject := ownedObject.Data.Content.Fields
			for _, fieldKey := range fieldKeys {
				fieldValue, ok := parsedObject[fieldKey].(string)
				if !ok {
					return nil, fmt.Errorf("field %s not found in object %s", fieldKey, qualifiedName)
				}

				foundValues[fieldKey] = fieldValue
			}
		}
	}

	return foundValues, nil
}
