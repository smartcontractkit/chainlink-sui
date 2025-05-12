package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/pattonkan/sui-go/sui/suiptb"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"golang.org/x/sync/semaphore"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

const maxCoinsPageSize = 50

type SuiPTBClient interface {
	MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error)
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (SuiTransactionBlockResponse, error)
	ReadObjectId(ctx context.Context, objectId string) (map[string]any, error)
	ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string) (*suiclient.ExecutionResultType, error)
	SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte, executionRequestType TransactionRequestType) (SuiTransactionBlockResponse, error)
	QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*suiclient.EventPage, error)
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]CoinData, error)
	ToPTBArg(ctx context.Context, builder *suiptb.ProgrammableTransactionBuilder, argValue any) (suiptb.Argument, error)
	EstimateGas(ctx context.Context, txBytes string) (uint64, error)
	FinishPTBAndSend(ctx context.Context, signerPublicKey []byte, builder *suiptb.ProgrammableTransactionBuilder) (SuiTransactionBlockResponse, error)
	BlockByDigest(ctx context.Context, txDigest string) (*SuiTransactionBlockResponse, error)
	GetSUIBalance(ctx context.Context, address string) (*big.Int, error)
}

// PTBClient implements SuiClient interface using the bindings/bind package
type PTBClient struct {
	log                logger.Logger
	client             *suiclient.ClientImpl
	maxRetries         *int
	transactionTimeout time.Duration
	keystoreService    loop.Keystore
	rateLimiter        *semaphore.Weighted
	defaultRequestType TransactionRequestType
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
	log.Infof("Creating new SUI client")

	client := suiclient.NewClient(rpcUrl)

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
		packageId, err := sui.AddressFromHex(req.PackageObjectId)
		if err != nil {
			return fmt.Errorf("invalid package ID: %w", err)
		}

		ptb, err := bind.BuildPTBFromArgs(ctx, *c.client, packageId, req.Module, req.Function, false, "", req.Arguments...)
		if err != nil {
			return fmt.Errorf("failed to build PTB: %w", err)
		}

		txBytes, err := bind.FinishTransactionFromBuilder(ctx, ptb, bind.TxOpts{}, req.Signer, *c.client)
		if err != nil {
			return fmt.Errorf("failed to finish transaction: %w", err)
		}

		result.TxBytes = base64.StdEncoding.EncodeToString(txBytes)

		return nil
	})

	return result, err
}

func (c *PTBClient) SendTransaction(ctx context.Context, payload TransactionBlockRequest) (SuiTransactionBlockResponse, error) {
	var result SuiTransactionBlockResponse
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		suiSignatures, err := bind.ToSuiSignatures(payload.Signatures)
		if err != nil {
			return fmt.Errorf("failed to convert signatures: %w", err)
		}

		b64Tx, err := sui.NewBase64Data(payload.TxBytes)
		if err != nil {
			return fmt.Errorf("invalid transaction bytes: %w", err)
		}

		options := &suiclient.SuiTransactionBlockResponseOptions{
			ShowInput:          payload.Options.ShowInput,
			ShowRawInput:       payload.Options.ShowRawInput,
			ShowEffects:        payload.Options.ShowEffects,
			ShowEvents:         payload.Options.ShowEvents,
			ShowObjectChanges:  payload.Options.ShowObjectChanges,
			ShowBalanceChanges: payload.Options.ShowBalanceChanges,
		}

		// Convert string request type to the appropriate type
		requestType := c.defaultRequestType
		if payload.RequestType != "" {
			requestType = TransactionRequestType(payload.RequestType)
		}

		blockReq := &suiclient.ExecuteTransactionBlockRequest{
			TxDataBytes: *b64Tx,
			Signatures:  suiSignatures,
			Options:     options,
			RequestType: suiclient.ExecuteTransactionRequestType(requestType),
		}

		c.log.Debugw("Executing transaction", "request", blockReq)

		response, err := c.client.ExecuteTransactionBlock(ctx, blockReq)
		if err != nil {
			return fmt.Errorf("failed to execute transaction: %w", err)
		}

		// Convert suiclient response to models response
		result = convertSuiResponse(response)

		return nil
	})

	return result, err
}

func (c *PTBClient) ReadObjectId(ctx context.Context, objectId string) (map[string]any, error) {
	var result map[string]any
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		response, err := bind.ReadObject(ctx, objectId, *c.client)
		if err != nil {
			return fmt.Errorf("failed to read object: %w", err)
		}

		if response.Data == nil || response.Data.Content == nil {
			return fmt.Errorf("object has no content")
		}

		if response.Data.Content.Data.MoveObject == nil ||
			response.Data.Content.Data.MoveObject.Fields == nil {
			return fmt.Errorf("object has no fields")
		}

		// Decode JSON fields to map
		err = json.Unmarshal(response.Data.Content.Data.MoveObject.Fields, &result)
		if err != nil {
			return fmt.Errorf("failed to decode object fields: %w", err)
		}

		return nil
	})

	return result, err
}

func (c *PTBClient) EstimateGas(ctx context.Context, txBytes string) (uint64, error) {
	response, err := c.client.DryRunTransaction(ctx, sui.Base64Data(txBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to dry run transaction: %w", err)
	}

	// Referenced from https://docs.sui.io/concepts/tokenomics/gas-in-sui
	fee := response.Effects.Data.V1.GasUsed.StorageCost.Uint64() -
		response.Effects.Data.V1.GasUsed.StorageRebate.Uint64() +
		response.Effects.Data.V1.GasUsed.ComputationCost.Uint64()

	return fee, nil
}

func (c *PTBClient) ReadFunction(ctx context.Context, signerAddress string, packageId string, module string, function string, args []any, argTypes []string) (*suiclient.ExecutionResultType, error) {
	var result *suiclient.ExecutionResultType
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		pkgId, err := sui.AddressFromHex(packageId)
		if err != nil {
			return fmt.Errorf("invalid package ID: %w", err)
		}

		ptb, err := bind.BuildPTBFromArgs(ctx, *c.client, pkgId, module, function, false, "", args...)
		if err != nil {
			return fmt.Errorf("failed to build PTB: %w", err)
		}

		txBytes, err := bind.FinishDevInspectTransactionFromBuilder(ctx, ptb, bind.TxOpts{}, signerAddress, *c.client)
		if err != nil {
			return fmt.Errorf("failed to finish transaction: %w", err)
		}

		response, err := bind.DevInspectTx(ctx, signerAddress, *c.client, txBytes)
		if err != nil {
			return fmt.Errorf("failed to inspect transaction: %w", err)
		}

		if len(response.Results) == 0 {
			return fmt.Errorf("no results from function call")
		}

		result = &response.Results[0]

		return nil
	})

	return result, err
}

func (c *PTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerPublicKey []byte, executionRequestType TransactionRequestType) (SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	txBytes, err := base64.StdEncoding.DecodeString(txBytesRaw)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to decode tx bytes: %w", err)
	}

	signerAddress, err := GetAddressFromPublicKey(signerPublicKey)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to get signer address: %w", err)
	}

	signatures, err := c.keystoreService.Sign(ctx, signerAddress, txBytes)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to sign tx: %w", err)
	}

	signaturesString := SerializeSuiSignature(signatures, signerPublicKey)

	b64bytes := codec.EncodeBase64(txBytes)

	return c.SendTransaction(ctx, TransactionBlockRequest{
		TxBytes:     b64bytes,
		Signatures:  []string{signaturesString},
		RequestType: string(c.defaultRequestType),
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

func (c *PTBClient) QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *EventId, sortOptions *QuerySortOptions) (*suiclient.EventPage, error) {
	var result *suiclient.EventPage
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		// Create package ID
		packageId, err := sui.AddressFromHex(filter.Package)
		if err != nil {
			return fmt.Errorf("invalid package ID: %w", err)
		}

		// Create the query filter structure
		queryFilter := &suiclient.EventFilter{
			MoveEventType: &sui.StructTag{
				Address: packageId,
				Module:  filter.Module,
				Name:    filter.Event,
			},
		}

		// Convert cursor to SDK client format
		var queryCursor *suiclient.EventId
		if cursor != nil {
			queryCursor = &suiclient.EventId{
				TxDigest: *sui.MustNewDigest(cursor.TxDigest),
				EventSeq: cursor.EventSeq,
			}
		} else {
			queryCursor = nil
		}

		// Create the query request
		req := &suiclient.QueryEventsRequest{
			Query:           queryFilter,
			Cursor:          queryCursor,
			Limit:           limit,
			DescendingOrder: sortOptions.Descending,
		}

		response, err := c.client.QueryEvents(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to query events: %w", err)
		}

		// Convert response to models format
		result = response

		return nil
	})

	return result, err
}

func (c *PTBClient) GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error) {
	var result TransactionResult
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		txDigest, err := sui.NewDigest(digest)
		if err != nil {
			return fmt.Errorf("invalid tx digest: %w", err)
		}

		req := &suiclient.GetTransactionBlockRequest{
			Digest: txDigest,
			Options: &suiclient.SuiTransactionBlockResponseOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
			},
		}

		response, err := c.client.GetTransactionBlock(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get transaction: %w", err)
		}

		// Set status if available
		if response.Effects != nil && response.Effects.Data.V1 != nil {
			result.Status = response.Effects.Data.V1.Status.Status
			result.Error = response.Effects.Data.V1.Status.Error
		}

		return nil
	})

	return result, err
}

func (c *PTBClient) GetCoinsByAddress(ctx context.Context, address string) ([]CoinData, error) {
	var result []CoinData
	pageLimit := uint(maxCoinsPageSize)

	// NOTE: the context with the timeout (from the WithRateLimit callback) is being ignored
	// because it is taking a significant amount of time in local testing (>30s) which is unusual.
	// Currently deferring to using a timeout context when calling the GetCoinsByAddress method.
	err := c.WithRateLimit(ctx, func(_ctx context.Context) error {
		suiAddress, err := sui.AddressFromHex(address)
		if err != nil {
			return fmt.Errorf("invalid address: %w", err)
		}

		var cursor *sui.ObjectId
		for {
			req := &suiclient.GetCoinsRequest{
				Owner:  suiAddress,
				Limit:  pageLimit,
				Cursor: cursor,
			}

			c.log.Debugw("About to fetch coin data", "request", req)

			resp, err := c.client.GetCoins(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to get coins: %w", err)
			}

			c.log.Debugw("Coins page", "response", resp)

			// Convert coin data to models format
			for _, coin := range resp.Data {
				coinObjId := coin.CoinObjectId
				result = append(result, CoinData{
					CoinType:     coin.CoinType,
					CoinObjectId: coinObjId.String(),
					Version:      fmt.Sprint(coin.Version),
					Digest:       coin.Digest.String(),
					Balance:      coin.Balance.String(),
					PreviousTx:   coin.PreviousTransaction.String(),
				})
			}

			if !resp.HasNextPage {
				break
			}

			cursor = resp.NextCursor

			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		return nil
	})

	return result, err
}

// Helper functions to convert between suiclient and models types
func convertSuiResponse(resp *suiclient.SuiTransactionBlockResponse) SuiTransactionBlockResponse {
	result := SuiTransactionBlockResponse{}

	// If response is nil, return empty result
	if resp == nil {
		return result
	}

	// Set digest if available
	if resp.Digest != nil {
		result.TxDigest = resp.Digest.String()
	}

	if resp.Events != nil {
		result.Events = resp.Events
	}

	// Check for nil at each level before accessing nested properties
	if resp.Effects != nil &&
		resp.Effects.Data.V1 != nil {
		// Copy effects
		result.Effects = *resp.Effects.Data.V1

		// Set status
		result.Status = SuiExecutionStatus{
			Status: resp.Effects.Data.V1.Status.Status,
			Error:  resp.Effects.Data.V1.Status.Error,
		}
	}

	if resp.ObjectChanges != nil {
		result.ObjectChanges = resp.ObjectChanges
	}

	return result
}

// ToPTBArg converts an argument into a format compatible with PTB based on the specified type.
func (c *PTBClient) ToPTBArg(ctx context.Context, builder *suiptb.ProgrammableTransactionBuilder, argValue any) (suiptb.Argument, error) {
	// TODO: this method should be improved in the following ways
	// 1. Given that we already know the expect type from the config, the actual conversion to a PTB arg type should be more strict and well defined
	// 		than what's currently available in the bindings
	// 2. There's no need to pass the builder (by value) around which incurs a lot of (extra) work on the underlying Go process
	//
	// NOTE: This is currently placed here simply to avoid leaking the SDK client outside
	return bind.ToPTBArg(ctx, builder, *c.client, argValue)
}

// FinishPTBAndSend receives a constructed PTB and proceeds to attach a gas token and finally signs and sends the request.
func (c *PTBClient) FinishPTBAndSend(ctx context.Context, signerPublicKey []byte, builder *suiptb.ProgrammableTransactionBuilder) (SuiTransactionBlockResponse, error) {
	// TODO: edit `bind.FinishTransactionFromBuilder()` to receive a reference to the client instead of passing by value

	signerAddress, err := GetAddressFromPublicKey(signerPublicKey)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to get signer address: %w", err)
	}

	txBytes, err := bind.FinishTransactionFromBuilder(ctx, builder, bind.TxOpts{}, signerAddress, *c.client)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to finish transaction: %w", err)
	}

	b64bytes := codec.EncodeBase64(txBytes)

	return c.SignAndSendTransaction(ctx, b64bytes, signerPublicKey, c.defaultRequestType)
}

func (c *PTBClient) BlockByDigest(ctx context.Context, txDigest string) (*SuiTransactionBlockResponse, error) {
	response, err := c.client.GetTransactionBlock(ctx, &suiclient.GetTransactionBlockRequest{
		Digest: sui.MustNewDigest(txDigest),
		Options: &suiclient.SuiTransactionBlockResponseOptions{
			ShowBalanceChanges: false,
			ShowEffects:        false,
			ShowObjectChanges:  false,
			ShowRawInput:       false,
			ShowInput:          false,
			ShowEvents:         false,
			ShowRawEffects:     false,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get block: %w", err)
	}

	return &SuiTransactionBlockResponse{
		TxDigest:  response.Digest.String(),
		Height:    response.Checkpoint.Uint64(),
		Timestamp: response.TimestampMs.Uint64(),
	}, nil
}

func (c *PTBClient) GetSUIBalance(ctx context.Context, address string) (*big.Int, error) {
	accountAddress, err := sui.AddressFromHex(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	balanceResponse, err := c.client.GetBalance(ctx, &suiclient.GetBalanceRequest{
		CoinType: "",
		Owner:    accountAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %w", err)
	}

	return balanceResponse.TotalBalance.Int, nil
}
