// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_token_pool_rate_limiter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

const FunctionInfo = `[{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"consume_inbound","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"RateLimitState"},{"name":"dest_chain_selector","type":"u64"},{"name":"requested_tokens","type":"u64"}]},{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"consume_outbound","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"RateLimitState"},{"name":"dest_chain_selector","type":"u64"},{"name":"requested_tokens","type":"u64"}]},{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"destroy_rate_limiter","parameters":[{"name":"state","type":"RateLimitState"}]},{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"get_current_inbound_rate_limiter_state","parameters":[{"name":"state","type":"RateLimitState"},{"name":"clock","type":"Clock"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"get_current_outbound_rate_limiter_state","parameters":[{"name":"state","type":"RateLimitState"},{"name":"clock","type":"Clock"},{"name":"remote_chain_selector","type":"u64"}]},{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"new","parameters":null},{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"set_chain_rate_limiter_config","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"RateLimitState"},{"name":"remote_chain_selector","type":"u64"},{"name":"outbound_is_enabled","type":"bool"},{"name":"outbound_capacity","type":"u64"},{"name":"outbound_rate","type":"u64"},{"name":"inbound_is_enabled","type":"bool"},{"name":"inbound_capacity","type":"u64"},{"name":"inbound_rate","type":"u64"}]},{"package":"lock_release_token_pool","module":"token_pool_rate_limiter","name":"verify_zero_config","parameters":[{"name":"state","type":"RateLimitState"}]}]`

type ITokenPoolRateLimiter interface {
	New(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	ConsumeInbound(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*models.SuiTransactionBlockResponse, error)
	ConsumeOutbound(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state RateLimitState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error)
	GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	DestroyRateLimiter(ctx context.Context, opts *bind.CallOpts, state RateLimitState) (*models.SuiTransactionBlockResponse, error)
	VerifyZeroConfig(ctx context.Context, opts *bind.CallOpts, state RateLimitState) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ITokenPoolRateLimiterDevInspect
	Encoder() TokenPoolRateLimiterEncoder
	Bound() bind.IBoundContract
}

type ITokenPoolRateLimiterDevInspect interface {
	New(ctx context.Context, opts *bind.CallOpts) (RateLimitState, error)
	GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (bind.Object, error)
	GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (bind.Object, error)
}

type TokenPoolRateLimiterEncoder interface {
	New() (*bind.EncodedCall, error)
	NewWithArgs(args ...any) (*bind.EncodedCall, error)
	ConsumeInbound(clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*bind.EncodedCall, error)
	ConsumeInboundWithArgs(args ...any) (*bind.EncodedCall, error)
	ConsumeOutbound(clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*bind.EncodedCall, error)
	ConsumeOutboundWithArgs(args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfig(clock bind.Object, state RateLimitState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCurrentInboundRateLimiterState(state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetCurrentInboundRateLimiterStateWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCurrentOutboundRateLimiterState(state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetCurrentOutboundRateLimiterStateWithArgs(args ...any) (*bind.EncodedCall, error)
	DestroyRateLimiter(state RateLimitState) (*bind.EncodedCall, error)
	DestroyRateLimiterWithArgs(args ...any) (*bind.EncodedCall, error)
	VerifyZeroConfig(state RateLimitState) (*bind.EncodedCall, error)
	VerifyZeroConfigWithArgs(args ...any) (*bind.EncodedCall, error)
}

type TokenPoolRateLimiterContract struct {
	*bind.BoundContract
	tokenPoolRateLimiterEncoder
	devInspect *TokenPoolRateLimiterDevInspect
}

type TokenPoolRateLimiterDevInspect struct {
	contract *TokenPoolRateLimiterContract
}

var _ ITokenPoolRateLimiter = (*TokenPoolRateLimiterContract)(nil)
var _ ITokenPoolRateLimiterDevInspect = (*TokenPoolRateLimiterDevInspect)(nil)

func NewTokenPoolRateLimiter(packageID string, client sui.ISuiAPI) (ITokenPoolRateLimiter, error) {
	contract, err := bind.NewBoundContract(packageID, "lock_release_token_pool", "token_pool_rate_limiter", client)
	if err != nil {
		return nil, err
	}

	c := &TokenPoolRateLimiterContract{
		BoundContract:               contract,
		tokenPoolRateLimiterEncoder: tokenPoolRateLimiterEncoder{BoundContract: contract},
	}
	c.devInspect = &TokenPoolRateLimiterDevInspect{contract: c}
	return c, nil
}

func (c *TokenPoolRateLimiterContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *TokenPoolRateLimiterContract) Encoder() TokenPoolRateLimiterEncoder {
	return c.tokenPoolRateLimiterEncoder
}

func (c *TokenPoolRateLimiterContract) DevInspect() ITokenPoolRateLimiterDevInspect {
	return c.devInspect
}

type RateLimitState struct {
	OutboundRateLimiterConfig bind.Object `move:"Table<u64, TokenBucket>"`
	InboundRateLimiterConfig  bind.Object `move:"Table<u64, TokenBucket>"`
}

type TokensConsumed struct {
	RemoteChainSelector uint64 `move:"u64"`
	Tokens              uint64 `move:"u64"`
}

type ConfigChanged struct {
	RemoteChainSelector uint64 `move:"u64"`
	OutboundIsEnabled   bool   `move:"bool"`
	OutboundCapacity    uint64 `move:"u64"`
	OutboundRate        uint64 `move:"u64"`
	InboundIsEnabled    bool   `move:"bool"`
	InboundCapacity     uint64 `move:"u64"`
	InboundRate         uint64 `move:"u64"`
}

func init() {
	bind.RegisterStructDecoder("lock_release_token_pool::token_pool_rate_limiter::RateLimitState", func(data []byte) (interface{}, error) {
		var result RateLimitState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for RateLimitState
	bind.RegisterStructDecoder("vector<lock_release_token_pool::token_pool_rate_limiter::RateLimitState>", func(data []byte) (interface{}, error) {
		var results []RateLimitState
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("lock_release_token_pool::token_pool_rate_limiter::TokensConsumed", func(data []byte) (interface{}, error) {
		var result TokensConsumed
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TokensConsumed
	bind.RegisterStructDecoder("vector<lock_release_token_pool::token_pool_rate_limiter::TokensConsumed>", func(data []byte) (interface{}, error) {
		var results []TokensConsumed
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("lock_release_token_pool::token_pool_rate_limiter::ConfigChanged", func(data []byte) (interface{}, error) {
		var result ConfigChanged
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ConfigChanged
	bind.RegisterStructDecoder("vector<lock_release_token_pool::token_pool_rate_limiter::ConfigChanged>", func(data []byte) (interface{}, error) {
		var results []ConfigChanged
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// New executes the new Move function.
func (c *TokenPoolRateLimiterContract) New(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.New()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConsumeInbound executes the consume_inbound Move function.
func (c *TokenPoolRateLimiterContract) ConsumeInbound(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.ConsumeInbound(clock, state, destChainSelector, requestedTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConsumeOutbound executes the consume_outbound Move function.
func (c *TokenPoolRateLimiterContract) ConsumeOutbound(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.ConsumeOutbound(clock, state, destChainSelector, requestedTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfig executes the set_chain_rate_limiter_config Move function.
func (c *TokenPoolRateLimiterContract) SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state RateLimitState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.SetChainRateLimiterConfig(clock, state, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCurrentInboundRateLimiterState executes the get_current_inbound_rate_limiter_state Move function.
func (c *TokenPoolRateLimiterContract) GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.GetCurrentInboundRateLimiterState(state, clock, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCurrentOutboundRateLimiterState executes the get_current_outbound_rate_limiter_state Move function.
func (c *TokenPoolRateLimiterContract) GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.GetCurrentOutboundRateLimiterState(state, clock, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyRateLimiter executes the destroy_rate_limiter Move function.
func (c *TokenPoolRateLimiterContract) DestroyRateLimiter(ctx context.Context, opts *bind.CallOpts, state RateLimitState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.DestroyRateLimiter(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// VerifyZeroConfig executes the verify_zero_config Move function.
func (c *TokenPoolRateLimiterContract) VerifyZeroConfig(ctx context.Context, opts *bind.CallOpts, state RateLimitState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolRateLimiterEncoder.VerifyZeroConfig(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// New executes the new Move function using DevInspect to get return values.
//
// Returns: RateLimitState
func (d *TokenPoolRateLimiterDevInspect) New(ctx context.Context, opts *bind.CallOpts) (RateLimitState, error) {
	encoded, err := d.contract.tokenPoolRateLimiterEncoder.New()
	if err != nil {
		return RateLimitState{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return RateLimitState{}, err
	}
	if len(results) == 0 {
		return RateLimitState{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(RateLimitState)
	if !ok {
		return RateLimitState{}, fmt.Errorf("unexpected return type: expected RateLimitState, got %T", results[0])
	}
	return result, nil
}

// GetCurrentInboundRateLimiterState executes the get_current_inbound_rate_limiter_state Move function using DevInspect to get return values.
//
// Returns: rate_limiter::TokenBucket
func (d *TokenPoolRateLimiterDevInspect) GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (bind.Object, error) {
	encoded, err := d.contract.tokenPoolRateLimiterEncoder.GetCurrentInboundRateLimiterState(state, clock, remoteChainSelector)
	if err != nil {
		return bind.Object{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return bind.Object{}, err
	}
	if len(results) == 0 {
		return bind.Object{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bind.Object)
	if !ok {
		return bind.Object{}, fmt.Errorf("unexpected return type: expected bind.Object, got %T", results[0])
	}
	return result, nil
}

// GetCurrentOutboundRateLimiterState executes the get_current_outbound_rate_limiter_state Move function using DevInspect to get return values.
//
// Returns: rate_limiter::TokenBucket
func (d *TokenPoolRateLimiterDevInspect) GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, state RateLimitState, clock bind.Object, remoteChainSelector uint64) (bind.Object, error) {
	encoded, err := d.contract.tokenPoolRateLimiterEncoder.GetCurrentOutboundRateLimiterState(state, clock, remoteChainSelector)
	if err != nil {
		return bind.Object{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return bind.Object{}, err
	}
	if len(results) == 0 {
		return bind.Object{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bind.Object)
	if !ok {
		return bind.Object{}, fmt.Errorf("unexpected return type: expected bind.Object, got %T", results[0])
	}
	return result, nil
}

type tokenPoolRateLimiterEncoder struct {
	*bind.BoundContract
}

// New encodes a call to the new Move function.
func (c tokenPoolRateLimiterEncoder) New() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"lock_release_token_pool::token_pool_rate_limiter::RateLimitState",
	})
}

// NewWithArgs encodes a call to the new Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) NewWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new", typeArgsList, typeParamsList, expectedParams, args, []string{
		"lock_release_token_pool::token_pool_rate_limiter::RateLimitState",
	})
}

// ConsumeInbound encodes a call to the consume_inbound Move function.
func (c tokenPoolRateLimiterEncoder) ConsumeInbound(clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("consume_inbound", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&mut RateLimitState",
		"u64",
		"u64",
	}, []any{
		clock,
		state,
		destChainSelector,
		requestedTokens,
	}, nil)
}

// ConsumeInboundWithArgs encodes a call to the consume_inbound Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) ConsumeInboundWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&mut RateLimitState",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("consume_inbound", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ConsumeOutbound encodes a call to the consume_outbound Move function.
func (c tokenPoolRateLimiterEncoder) ConsumeOutbound(clock bind.Object, state RateLimitState, destChainSelector uint64, requestedTokens uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("consume_outbound", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&mut RateLimitState",
		"u64",
		"u64",
	}, []any{
		clock,
		state,
		destChainSelector,
		requestedTokens,
	}, nil)
}

// ConsumeOutboundWithArgs encodes a call to the consume_outbound Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) ConsumeOutboundWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&mut RateLimitState",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("consume_outbound", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetChainRateLimiterConfig encodes a call to the set_chain_rate_limiter_config Move function.
func (c tokenPoolRateLimiterEncoder) SetChainRateLimiterConfig(clock bind.Object, state RateLimitState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&mut RateLimitState",
		"u64",
		"bool",
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	}, []any{
		clock,
		state,
		remoteChainSelector,
		outboundIsEnabled,
		outboundCapacity,
		outboundRate,
		inboundIsEnabled,
		inboundCapacity,
		inboundRate,
	}, nil)
}

// SetChainRateLimiterConfigWithArgs encodes a call to the set_chain_rate_limiter_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) SetChainRateLimiterConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&mut RateLimitState",
		"u64",
		"bool",
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCurrentInboundRateLimiterState encodes a call to the get_current_inbound_rate_limiter_state Move function.
func (c tokenPoolRateLimiterEncoder) GetCurrentInboundRateLimiterState(state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_current_inbound_rate_limiter_state", typeArgsList, typeParamsList, []string{
		"&RateLimitState",
		"&Clock",
		"u64",
	}, []any{
		state,
		clock,
		remoteChainSelector,
	}, []string{
		"rate_limiter::TokenBucket",
	})
}

// GetCurrentInboundRateLimiterStateWithArgs encodes a call to the get_current_inbound_rate_limiter_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) GetCurrentInboundRateLimiterStateWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RateLimitState",
		"&Clock",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_current_inbound_rate_limiter_state", typeArgsList, typeParamsList, expectedParams, args, []string{
		"rate_limiter::TokenBucket",
	})
}

// GetCurrentOutboundRateLimiterState encodes a call to the get_current_outbound_rate_limiter_state Move function.
func (c tokenPoolRateLimiterEncoder) GetCurrentOutboundRateLimiterState(state RateLimitState, clock bind.Object, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_current_outbound_rate_limiter_state", typeArgsList, typeParamsList, []string{
		"&RateLimitState",
		"&Clock",
		"u64",
	}, []any{
		state,
		clock,
		remoteChainSelector,
	}, []string{
		"rate_limiter::TokenBucket",
	})
}

// GetCurrentOutboundRateLimiterStateWithArgs encodes a call to the get_current_outbound_rate_limiter_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) GetCurrentOutboundRateLimiterStateWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RateLimitState",
		"&Clock",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_current_outbound_rate_limiter_state", typeArgsList, typeParamsList, expectedParams, args, []string{
		"rate_limiter::TokenBucket",
	})
}

// DestroyRateLimiter encodes a call to the destroy_rate_limiter Move function.
func (c tokenPoolRateLimiterEncoder) DestroyRateLimiter(state RateLimitState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("destroy_rate_limiter", typeArgsList, typeParamsList, []string{
		"lock_release_token_pool::token_pool_rate_limiter::RateLimitState",
	}, []any{
		state,
	}, nil)
}

// DestroyRateLimiterWithArgs encodes a call to the destroy_rate_limiter Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) DestroyRateLimiterWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"lock_release_token_pool::token_pool_rate_limiter::RateLimitState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("destroy_rate_limiter", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// VerifyZeroConfig encodes a call to the verify_zero_config Move function.
func (c tokenPoolRateLimiterEncoder) VerifyZeroConfig(state RateLimitState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify_zero_config", typeArgsList, typeParamsList, []string{
		"&RateLimitState",
	}, []any{
		state,
	}, nil)
}

// VerifyZeroConfigWithArgs encodes a call to the verify_zero_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolRateLimiterEncoder) VerifyZeroConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RateLimitState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify_zero_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}
