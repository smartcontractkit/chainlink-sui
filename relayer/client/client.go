package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/fardream/go-bcs/bcs"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/signer"

	suiAlt "github.com/pattonkan/sui-go/sui"
	suiPtb "github.com/pattonkan/sui-go/sui/suiptb"
	suiAltClient "github.com/pattonkan/sui-go/suiclient"
	"golang.org/x/sync/semaphore"
)

const maxCoinsPageSize = 50

type TransactionResult struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type SuiClient interface {
	MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error)
	SendTransaction(ctx context.Context, payload TransactionBlockRequest) (models.SuiTransactionBlockResponse, error)
	ReadObjectId(ctx context.Context, objectId string) (map[string]any, error)
	ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string) (*suiAltClient.ExecutionResultType, error)
	SignAndSendTransaction(ctx context.Context, txBytes string, signerOverride *signer.SuiSigner, executionRequestType TransactionRequestType) (models.SuiTransactionBlockResponse, error)
	QueryEvents(ctx context.Context, filter models.EventFilterByMoveEventModule, limit uint64, cursor *models.EventId, descending bool) (models.PaginatedEventsResponse, error)
	WithRateLimit(ctx context.Context, f func(ctx context.Context) error) error
	GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error)
	GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error)
}

type Client struct {
	log                logger.Logger
	client             sui.ISuiAPI
	ptbClient          *suiAltClient.ClientImpl
	maxRetries         *int
	transactionTimeout time.Duration
	signer             *signer.SuiSigner
	// while this is a weighted semaphore, we currently only use a weight of 1 for each request
	// that the client makes to the code
	rateLimiter *semaphore.Weighted
}

var _ SuiClient = (*Client)(nil)

func NewClient(log logger.Logger, rpcUrl string, maxRetries *int, transactionTimeout time.Duration, defaultSigner *signer.SuiSigner, maxConcurrentRequests int64) (*Client, error) {
	baseClient := sui.NewSuiClient(rpcUrl)
	ptbClient := suiAltClient.NewClient(rpcUrl)

	if maxConcurrentRequests <= 0 {
		maxConcurrentRequests = 100 // Default value
	}

	return &Client{
		log:                log,
		client:             baseClient,
		ptbClient:          ptbClient,
		maxRetries:         maxRetries,
		transactionTimeout: transactionTimeout,
		signer:             defaultSigner,
		rateLimiter:        semaphore.NewWeighted(maxConcurrentRequests),
	}, nil
}

func (c *Client) WithRateLimit(ctx context.Context, f func(ctx context.Context) error) error {
	// First create a timeout context that will be used for the actual operation
	timeoutCtx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	// If no rate limiter, just run the function with the timeout context
	if c.rateLimiter == nil {
		return f(timeoutCtx)
	}

	// Use the ORIGINAL context for rate limiting acquisition
	// This ensures parent cancellation can still cancel before we even try to acquire
	if err := c.rateLimiter.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("failed to acquire rate limit: %w", err)
	}
	defer c.rateLimiter.Release(1)

	// Execute the function with the timeout context
	return f(timeoutCtx)
}

func (c *Client) MoveCall(ctx context.Context, req models.MoveCallRequest) (models.TxnMetaData, error) {
	var result models.TxnMetaData
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		var err error
		result, err = c.client.MoveCall(ctx, req)

		return err
	})

	return result, err
}

func (c *Client) SendTransaction(ctx context.Context, payload TransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
	var result models.SuiTransactionBlockResponse
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		clientPayload := models.SuiExecuteTransactionBlockRequest{
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

		var err error
		result, err = c.client.SuiExecuteTransactionBlock(ctx, clientPayload)

		return err
	})

	return result, err
}

// ReadObjectId reads an object from the Sui blockchain
func (c *Client) ReadObjectId(ctx context.Context, objectId string) (map[string]any, error) {
	var result map[string]any
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		object, err := c.client.SuiGetObject(ctx, models.SuiGetObjectRequest{
			ObjectId: objectId,
			Options: models.SuiObjectDataOptions{
				ShowContent: true,
			},
		})

		if err != nil {
			return fmt.Errorf("failed to get object by ID: %w", err)
		}

		result = object.Data.Content.Fields

		return nil
	})

	return result, err
}

// ReadFunction calls a Move contract function and returns the value.
// The implementation internally signs the transactions with the signer attached to the client.
// This method also calls the Move contract in "devInspect" execution mode since it is only reading values.
func (c *Client) ReadFunction(ctx context.Context, packageId string, module string, function string, args []any, argTypes []string) (*suiAltClient.ExecutionResultType, error) {
	var result *suiAltClient.ExecutionResultType
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		addressHex, _ := (*c.signer).GetAddress()
		address, err := suiAlt.AddressFromHex(addressHex)
		if err != nil {
			return err
		}

		packageIdObj, err := suiAlt.PackageIdFromHex(packageId)
		if err != nil {
			return err
		}

		// Create a new client instance pointing to your Sui node's RPC endpoint.
		ptb := suiPtb.NewTransactionDataTransactionBuilder()

		// Convert each string type into a "TypeArg"
		typeTagArgs := make([]suiAlt.TypeTag, len(argTypes))
		for i, argType := range argTypes {
			typeTag, tagErr := suiAlt.NewTypeTag(argType)
			if tagErr != nil {
				return fmt.Errorf("failed to create type tag: %w", err)
			}
			typeTagArgs[i] = *typeTag
		}

		// Convert each arg into a "CallArg" type
		callArgs := make([]suiPtb.CallArg, len(args))
		for i, arg := range args {
			encodedArg, encodedArgErr := codec.EncodePtbFunctionParam(argTypes[i], arg)
			if encodedArgErr != nil {
				return fmt.Errorf("failed to encode argument: %w", err)
			}
			callArgs[i] = encodedArg
		}

		err = ptb.MoveCall(packageIdObj, module, function, []suiAlt.TypeTag{}, callArgs)
		if err != nil {
			return fmt.Errorf("failed to move call: %w", err)
		}

		pt := ptb.Finish()
		tx := suiPtb.NewTransactionData(
			address,
			pt,
			nil,
			suiAltClient.DefaultGasBudget,
			suiAltClient.DefaultGasPrice,
		)

		txBytes, err := bcs.Marshal(tx.V1.Kind)
		if err != nil {
			return fmt.Errorf("failed to marshal transaction: %w", err)
		}

		resp, err := c.ptbClient.DevInspectTransactionBlock(ctx, &suiAltClient.DevInspectTransactionBlockRequest{
			SenderAddress: address,
			TxKindBytes:   txBytes,
		})
		if err != nil {
			return fmt.Errorf("failed to call read function: %w", err)
		}
		if len(resp.Results) == 0 {
			return fmt.Errorf("failed to call read function: no results")
		}

		c.log.Debugw("Dev inspect results", "results", resp)
		result = &resp.Results[0]

		return nil
	})

	return result, err
}

// SignAndSendTransaction given a plain (non-encoded) transaction, signs it and sends it to the node.
// The implementation uses the signer attached (default) to the client or the signer provided in the argument if specified.
// The transaction bytes should be in base64 encoded format.
// The executionRequestType parameter determines how the transaction is executed (e.g., "WaitForLocalExecution").
// Returns a SuiTransactionBlockResponse containing the transaction results, including inputs, effects, and changes.
// If signing or sending fails, an error is returned with context about the failure.
// This method doesn't use rate limiting because it only signs the request and internally calls SendTransaction which already has rate limiting applied.
func (c *Client) SignAndSendTransaction(ctx context.Context, txBytesRaw string, signerOverride *signer.SuiSigner, executionRequestType TransactionRequestType) (models.SuiTransactionBlockResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.transactionTimeout)
	defer cancel()

	if signerOverride == nil {
		// fallback to the default signer if no override is provided
		signerOverride = c.signer
	}

	txBytes, err := base64.StdEncoding.DecodeString(txBytesRaw)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("failed to decode tx bytes: %w", err)
	}

	signatures, err := (*signerOverride).Sign(txBytes)
	if err != nil {
		return models.SuiTransactionBlockResponse{}, fmt.Errorf("error signing transaction: %w", err)
	}

	// Call SendTransaction which already has rate limiting applied
	return c.SendTransaction(ctx, TransactionBlockRequest{
		TxBytes:    txBytesRaw,
		Signatures: signatures,
		Options: TransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowObjectChanges:  true,
			ShowBalanceChanges: true,
		},
		RequestType: string(executionRequestType),
	})
}

// QueryEvents queries events from the Sui network with flexible filtering options.
// When used with package/module/eventType parameters, it constructs a MoveEventType filter.
// Parameters:
// - ctx: Context for the request
// - filter: The event filter object which specifies the package, the module and the event type
// - cursor: (optional) the EventId cursor to offset the result by
// - descending: Whether to sort in descending order
// Returns events matching the criteria or an error.
func (c *Client) QueryEvents(
	ctx context.Context,
	filter models.EventFilterByMoveEventModule,
	limit uint64,
	cursor *models.EventId,
	descending bool,
) (models.PaginatedEventsResponse, error) {
	var result models.PaginatedEventsResponse
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		// Execute the query to get paginated events
		response, err := c.client.SuiXQueryEvents(ctx, models.SuiXQueryEventsRequest{
			SuiEventFilter:  filter,
			Cursor:          cursor,
			Limit:           limit,
			DescendingOrder: descending,
		})
		if err != nil {
			return fmt.Errorf("failed to query events: %w", err)
		}

		c.log.Debugw("Query events",
			"filter", fmt.Sprintf("%+v", filter),
			"limit", limit,
			"cursor", cursor,
			"response", fmt.Sprintf("%+v", response),
		)

		result = response

		return nil
	})

	return result, err
}

// GetTransactionStatus implements SuiClient.
func (c *Client) GetTransactionStatus(ctx context.Context, digest string) (TransactionResult, error) {
	var result models.SuiTransactionBlockResponse
	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		req := models.SuiGetTransactionBlockRequest{
			Digest: digest,
			Options: models.SuiTransactionBlockOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
			},
		}

		statusResponse, err := c.client.SuiGetTransactionBlock(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get transaction status: %w", err)
		}

		result = statusResponse

		return nil
	})

	return TransactionResult{
		Status: result.Effects.Status.Status,
		Error:  result.Effects.Status.Error,
	}, err
}

func (c *Client) GetCoinsByAddress(ctx context.Context, address string) ([]models.CoinData, error) {
	result := []models.CoinData{}
	pageLimit := uint64(maxCoinsPageSize) // Set the maximum page size

	err := c.WithRateLimit(ctx, func(ctx context.Context) error {
		var cursor *string // Start with nil cursor for first page

		// Loop until we've fetched all pages
		for {
			// Create request with pagination parameters
			request := models.SuiXGetAllCoinsRequest{
				Owner:  address,
				Limit:  pageLimit,
				Cursor: cursor,
			}

			// Fetch this page of coins
			resp, err := c.client.SuiXGetAllCoins(ctx, request)
			if err != nil {
				return fmt.Errorf("failed to get coins by address: %w", err)
			}

			// Add coins from this page to our result set
			result = append(result, resp.Data...)

			// Log how many coins we've collected so far
			c.log.Debugw("Fetched coins page",
				"address", address,
				"page_size", len(resp.Data),
				"total_so_far", len(result))

			if !resp.HasNextPage {
				// No more pages, exit the loop
				break
			}
			// Update cursor for next request
			cursor = &resp.NextCursor

			// Check if context is cancelled before making next request
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		return nil
	})

	return result, err
}
