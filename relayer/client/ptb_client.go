package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"golang.org/x/sync/semaphore"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"
	"github.com/smartcontractkit/chainlink-sui/shared"
)

const maxCoinsPageSize uint = 50
const Base10 = 10
const DefaultGasPrice = 10_000
const DefaultGasBudget = 1_000_000_000

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
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error)
	EstimateGas(ctx context.Context, txBytes string) (uint64, error)
	FinishPTBAndSend(ctx context.Context, txnSigner *signer.Signer, tx *transaction.Transaction, requestType TransactionRequestType) (SuiTransactionBlockResponse, error)
	BlockByDigest(ctx context.Context, txDigest string) (*SuiTransactionBlockResponse, error)
	GetSUIBalance(ctx context.Context, address string) (*big.Int, error)
	GetClient() *sui.ISuiAPI
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
	normalizedModules  map[string]map[string]models.GetNormalizedMoveModuleResponse
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

	client := sui.NewSuiClient(rpcUrl)

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
	}, nil
}

func (c *PTBClient) WithRateLimit(ctx context.Context, f func(ctx context.Context) error) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	if c.rateLimiter == nil {
		return f(timeoutCtx)
	}

	if err := c.rateLimiter.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire rate limit: %w", err)
	}
	defer c.rateLimiter.Release(1)

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
		if err != nil {
			return fmt.Errorf("failed to read function: %w", err)
		}

		c.log.Debugw("ReadFunction RPC response", "RPC response", response, "functionTag", fmt.Sprintf("%s::%s::%s", packageId, module, function))

		if len(response.Results) == 0 {
			return fmt.Errorf("no results from function call")
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
			structParts := strings.Split(structTag, "::")

			// create a bcs decoder from the return value
			bcsBytes, err := codec.AnySliceToBytes(returnedValue[0].([]any))
			if err != nil {
				return fmt.Errorf("failed to convert return value to bytes: %w", err)
			}
			bcsDecoder := bcs.NewDeserializer(bcsBytes)

			// if the response type is not a struct (primitive type), skip the result (keep it as is)
			structPartsLen := 3
			if len(structParts) != structPartsLen {
				primitive, err := codec.DecodeSuiPrimative(bcsDecoder, structTag)
				if err != nil {
					return fmt.Errorf("failed to decode primitive: %w", err)
				}
				results[i] = primitive
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

				results[i] = jsonResult
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

	signerAddress, err := GetAddressFromPublicKey(signerPublicKey)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to get signer address: %w", err)
	}

	txBytes, err := shared.DecodeBase64(txBytesRaw)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to decode tx bytes: %w", err)
	}

	signatures, err := c.keystoreService.Sign(ctx, signerAddress, txBytes)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to sign tx: %w", err)
	}

	signaturesString := SerializeSuiSignature(signatures, signerPublicKey)

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

		eventFilter := models.EventFilterByMoveEventModule{
			MoveEventModule: models.MoveEventModule{
				Package: filter.Package,
				Module:  filter.Module,
				Event:   filter.Event,
			},
		}

		queryReq := models.SuiXQueryEventsRequest{
			SuiEventFilter:  eventFilter,
			Limit:           limitVal,
			DescendingOrder: sortOptions != nil && sortOptions.Descending,
		}

		if cursor != nil {
			queryReq.Cursor = cursor
		}

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

func (c *PTBClient) GetClient() *sui.ISuiAPI {
	return &c.client
}
