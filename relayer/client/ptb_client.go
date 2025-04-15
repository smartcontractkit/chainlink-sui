package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"golang.org/x/sync/semaphore"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

const maxCoinsPageSize = 50

type SuiPTBClient interface {
	MoveCall(ctx context.Context, req MoveCallRequest) (TxnMetaData, error)
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (SuiTransactionBlockResponse, error)
	ReadObjectId(ctx context.Context, objectId string) (map[string]any, error)
	ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string) (*suiclient.ExecutionResultType, error)
	SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerOverride *signer.SuiSigner, executionRequestType TransactionRequestType) (SuiTransactionBlockResponse, error)
	QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *suiclient.EventId, descending bool) (*suiclient.EventPage, error)
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]CoinData, error)
}

// PTBClient implements SuiClient interface using the bindings/bind package
type PTBClient struct {
	log                logger.Logger
	client             *suiclient.ClientImpl
	maxRetries         *int
	transactionTimeout time.Duration
	signer             *signer.SuiSigner
	rateLimiter        *semaphore.Weighted
}

var _ SuiPTBClient = (*PTBClient)(nil)

func NewPTBClient(log logger.Logger, rpcUrl string, maxRetries *int, transactionTimeout time.Duration, defaultSigner *signer.SuiSigner, maxConcurrentRequests int64) (*PTBClient, error) {
	client := suiclient.NewClient(rpcUrl)

	if maxConcurrentRequests <= 0 {
		maxConcurrentRequests = 100 // Default value
	}

	return &PTBClient{
		log:                log,
		client:             client,
		maxRetries:         maxRetries,
		transactionTimeout: transactionTimeout,
		signer:             defaultSigner,
		rateLimiter:        semaphore.NewWeighted(maxConcurrentRequests),
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

		// Get signer address
		signerAddr, err := (*c.signer).GetAddress()
		if err != nil {
			return fmt.Errorf("failed to get signer address: %w", err)
		}

		txBytes, err := bind.FinishTransactionFromBuilder(ctx, ptb, bind.TxOpts{}, signerAddr, *c.client)
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
		requestType := "WaitForLocalExecution"
		if payload.RequestType != "" {
			requestType = payload.RequestType
		}

		blockReq := &suiclient.ExecuteTransactionBlockRequest{
			TxDataBytes: *b64Tx,
			Signatures:  suiSignatures,
			Options:     options,
			RequestType: suiclient.ExecuteTransactionRequestType(requestType),
		}

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

func (c *PTBClient) ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string) (*suiclient.ExecutionResultType, error) {
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

		// Get signer address
		signerAddr, err := (*c.signer).GetAddress()
		if err != nil {
			return fmt.Errorf("failed to get signer address: %w", err)
		}

		txBytes, err := bind.FinishDevInspectTransactionFromBuilder(ctx, ptb, bind.TxOpts{}, signerAddr, *c.client)
		if err != nil {
			return fmt.Errorf("failed to finish transaction: %w", err)
		}

		response, err := bind.DevInspectTx(ctx, *c.signer, *c.client, txBytes)
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

func (c *PTBClient) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerOverride *signer.SuiSigner, executionRequestType TransactionRequestType) (SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	if signerOverride == nil {
		signerOverride = c.signer
	}

	txBytes, err := base64.StdEncoding.DecodeString(txBytesRaw)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to decode tx bytes: %w", err)
	}

	response, err := bind.SignAndSendTx(ctx, *signerOverride, *c.client, txBytes)
	if err != nil {
		return SuiTransactionBlockResponse{}, fmt.Errorf("failed to sign and send transaction: %w", err)
	}

	return convertSuiResponse(response), nil
}

func (c *PTBClient) QueryEvents(ctx context.Context, filter EventFilterByMoveEventModule, limit *uint, cursor *suiclient.EventId, descending bool) (*suiclient.EventPage, error) {
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

		// Create the query request
		req := &suiclient.QueryEventsRequest{
			Query:           queryFilter,
			Cursor:          cursor,
			Limit:           limit,
			DescendingOrder: descending,
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

		result.Status = response.Effects.Data.V1.Status.Status
		result.Error = response.Effects.Data.V1.Status.Error

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
	// Implementation of conversion logic
	// This is a simplified version
	return SuiTransactionBlockResponse{
		TxDigest: resp.Digest.String(),
		Status: SuiExecutionStatus{
			Status: resp.Effects.Data.V1.Status.Status,
			Error:  resp.Effects.Data.V1.Status.Error,
		},
		Effects: *resp.Effects.Data.V1,
	}
}
