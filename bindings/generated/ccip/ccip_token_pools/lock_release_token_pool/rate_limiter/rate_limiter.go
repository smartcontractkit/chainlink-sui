// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_rate_limiter

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

const FunctionInfo = `[{"package":"lock_release_token_pool","module":"rate_limiter","name":"consume","parameters":[{"name":"clock","type":"Clock"},{"name":"bucket","type":"TokenBucket"},{"name":"requested_tokens","type":"u64"}]},{"package":"lock_release_token_pool","module":"rate_limiter","name":"get_current_token_bucket_state","parameters":[{"name":"clock","type":"Clock"},{"name":"state","type":"TokenBucket"}]},{"package":"lock_release_token_pool","module":"rate_limiter","name":"get_rate","parameters":[{"name":"bucket","type":"TokenBucket"}]},{"package":"lock_release_token_pool","module":"rate_limiter","name":"get_token_bucket_fields","parameters":[{"name":"bucket","type":"TokenBucket"}]},{"package":"lock_release_token_pool","module":"rate_limiter","name":"is_enabled","parameters":[{"name":"bucket","type":"TokenBucket"}]},{"package":"lock_release_token_pool","module":"rate_limiter","name":"new","parameters":[{"name":"clock","type":"Clock"},{"name":"is_enabled","type":"bool"},{"name":"capacity","type":"u64"},{"name":"rate","type":"u64"}]},{"package":"lock_release_token_pool","module":"rate_limiter","name":"set_token_bucket_config","parameters":[{"name":"clock","type":"Clock"},{"name":"bucket","type":"TokenBucket"},{"name":"is_enabled","type":"bool"},{"name":"capacity","type":"u64"},{"name":"rate","type":"u64"}]}]`

type IRateLimiter interface {
	New(ctx context.Context, opts *bind.CallOpts, clock bind.Object, isEnabled bool, capacity uint64, rate uint64) (*models.SuiTransactionBlockResponse, error)
	GetCurrentTokenBucketState(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	Consume(ctx context.Context, opts *bind.CallOpts, clock bind.Object, bucket bind.Object, requestedTokens uint64) (*models.SuiTransactionBlockResponse, error)
	SetTokenBucketConfig(ctx context.Context, opts *bind.CallOpts, clock bind.Object, bucket bind.Object, isEnabled bool, capacity uint64, rate uint64) (*models.SuiTransactionBlockResponse, error)
	IsEnabled(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetRate(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenBucketFields(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IRateLimiterDevInspect
	Encoder() RateLimiterEncoder
	Bound() bind.IBoundContract
}

type IRateLimiterDevInspect interface {
	New(ctx context.Context, opts *bind.CallOpts, clock bind.Object, isEnabled bool, capacity uint64, rate uint64) (bind.Object, error)
	GetCurrentTokenBucketState(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state bind.Object) (bind.Object, error)
	IsEnabled(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (bool, error)
	GetRate(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (uint64, error)
	GetTokenBucketFields(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) ([]any, error)
}

type RateLimiterEncoder interface {
	New(clock bind.Object, isEnabled bool, capacity uint64, rate uint64) (*bind.EncodedCall, error)
	NewWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCurrentTokenBucketState(clock bind.Object, state bind.Object) (*bind.EncodedCall, error)
	GetCurrentTokenBucketStateWithArgs(args ...any) (*bind.EncodedCall, error)
	Consume(clock bind.Object, bucket bind.Object, requestedTokens uint64) (*bind.EncodedCall, error)
	ConsumeWithArgs(args ...any) (*bind.EncodedCall, error)
	SetTokenBucketConfig(clock bind.Object, bucket bind.Object, isEnabled bool, capacity uint64, rate uint64) (*bind.EncodedCall, error)
	SetTokenBucketConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	IsEnabled(bucket bind.Object) (*bind.EncodedCall, error)
	IsEnabledWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRate(bucket bind.Object) (*bind.EncodedCall, error)
	GetRateWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenBucketFields(bucket bind.Object) (*bind.EncodedCall, error)
	GetTokenBucketFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
}

type RateLimiterContract struct {
	*bind.BoundContract
	rateLimiterEncoder
	devInspect *RateLimiterDevInspect
}

type RateLimiterDevInspect struct {
	contract *RateLimiterContract
}

var _ IRateLimiter = (*RateLimiterContract)(nil)
var _ IRateLimiterDevInspect = (*RateLimiterDevInspect)(nil)

func NewRateLimiter(packageID string, client sui.ISuiAPI) (IRateLimiter, error) {
	contract, err := bind.NewBoundContract(packageID, "lock_release_token_pool", "rate_limiter", client)
	if err != nil {
		return nil, err
	}

	c := &RateLimiterContract{
		BoundContract:      contract,
		rateLimiterEncoder: rateLimiterEncoder{BoundContract: contract},
	}
	c.devInspect = &RateLimiterDevInspect{contract: c}
	return c, nil
}

func (c *RateLimiterContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *RateLimiterContract) Encoder() RateLimiterEncoder {
	return c.rateLimiterEncoder
}

func (c *RateLimiterContract) DevInspect() IRateLimiterDevInspect {
	return c.devInspect
}

func init() {
}

// New executes the new Move function.
func (c *RateLimiterContract) New(ctx context.Context, opts *bind.CallOpts, clock bind.Object, isEnabled bool, capacity uint64, rate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rateLimiterEncoder.New(clock, isEnabled, capacity, rate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCurrentTokenBucketState executes the get_current_token_bucket_state Move function.
func (c *RateLimiterContract) GetCurrentTokenBucketState(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rateLimiterEncoder.GetCurrentTokenBucketState(clock, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Consume executes the consume Move function.
func (c *RateLimiterContract) Consume(ctx context.Context, opts *bind.CallOpts, clock bind.Object, bucket bind.Object, requestedTokens uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rateLimiterEncoder.Consume(clock, bucket, requestedTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetTokenBucketConfig executes the set_token_bucket_config Move function.
func (c *RateLimiterContract) SetTokenBucketConfig(ctx context.Context, opts *bind.CallOpts, clock bind.Object, bucket bind.Object, isEnabled bool, capacity uint64, rate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rateLimiterEncoder.SetTokenBucketConfig(clock, bucket, isEnabled, capacity, rate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsEnabled executes the is_enabled Move function.
func (c *RateLimiterContract) IsEnabled(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rateLimiterEncoder.IsEnabled(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRate executes the get_rate Move function.
func (c *RateLimiterContract) GetRate(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rateLimiterEncoder.GetRate(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenBucketFields executes the get_token_bucket_fields Move function.
func (c *RateLimiterContract) GetTokenBucketFields(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rateLimiterEncoder.GetTokenBucketFields(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// New executes the new Move function using DevInspect to get return values.
//
// Returns: TokenBucket
func (d *RateLimiterDevInspect) New(ctx context.Context, opts *bind.CallOpts, clock bind.Object, isEnabled bool, capacity uint64, rate uint64) (bind.Object, error) {
	encoded, err := d.contract.rateLimiterEncoder.New(clock, isEnabled, capacity, rate)
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

// GetCurrentTokenBucketState executes the get_current_token_bucket_state Move function using DevInspect to get return values.
//
// Returns: TokenBucket
func (d *RateLimiterDevInspect) GetCurrentTokenBucketState(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state bind.Object) (bind.Object, error) {
	encoded, err := d.contract.rateLimiterEncoder.GetCurrentTokenBucketState(clock, state)
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

// IsEnabled executes the is_enabled Move function using DevInspect to get return values.
//
// Returns: bool
func (d *RateLimiterDevInspect) IsEnabled(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (bool, error) {
	encoded, err := d.contract.rateLimiterEncoder.IsEnabled(bucket)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// GetRate executes the get_rate Move function using DevInspect to get return values.
//
// Returns: u64
func (d *RateLimiterDevInspect) GetRate(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) (uint64, error) {
	encoded, err := d.contract.rateLimiterEncoder.GetRate(bucket)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// GetTokenBucketFields executes the get_token_bucket_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u64
//	[1]: u64
//	[2]: bool
//	[3]: u64
//	[4]: u64
func (d *RateLimiterDevInspect) GetTokenBucketFields(ctx context.Context, opts *bind.CallOpts, bucket bind.Object) ([]any, error) {
	encoded, err := d.contract.rateLimiterEncoder.GetTokenBucketFields(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

type rateLimiterEncoder struct {
	*bind.BoundContract
}

// New encodes a call to the new Move function.
func (c rateLimiterEncoder) New(clock bind.Object, isEnabled bool, capacity uint64, rate uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new", typeArgsList, typeParamsList, []string{
		"&Clock",
		"bool",
		"u64",
		"u64",
	}, []any{
		clock,
		isEnabled,
		capacity,
		rate,
	}, []string{
		"TokenBucket",
	})
}

// NewWithArgs encodes a call to the new Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rateLimiterEncoder) NewWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"bool",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new", typeArgsList, typeParamsList, expectedParams, args, []string{
		"TokenBucket",
	})
}

// GetCurrentTokenBucketState encodes a call to the get_current_token_bucket_state Move function.
func (c rateLimiterEncoder) GetCurrentTokenBucketState(clock bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_current_token_bucket_state", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&TokenBucket",
	}, []any{
		clock,
		state,
	}, []string{
		"TokenBucket",
	})
}

// GetCurrentTokenBucketStateWithArgs encodes a call to the get_current_token_bucket_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rateLimiterEncoder) GetCurrentTokenBucketStateWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&TokenBucket",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_current_token_bucket_state", typeArgsList, typeParamsList, expectedParams, args, []string{
		"TokenBucket",
	})
}

// Consume encodes a call to the consume Move function.
func (c rateLimiterEncoder) Consume(clock bind.Object, bucket bind.Object, requestedTokens uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("consume", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&mut TokenBucket",
		"u64",
	}, []any{
		clock,
		bucket,
		requestedTokens,
	}, nil)
}

// ConsumeWithArgs encodes a call to the consume Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rateLimiterEncoder) ConsumeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&mut TokenBucket",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("consume", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// SetTokenBucketConfig encodes a call to the set_token_bucket_config Move function.
func (c rateLimiterEncoder) SetTokenBucketConfig(clock bind.Object, bucket bind.Object, isEnabled bool, capacity uint64, rate uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_token_bucket_config", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&mut TokenBucket",
		"bool",
		"u64",
		"u64",
	}, []any{
		clock,
		bucket,
		isEnabled,
		capacity,
		rate,
	}, nil)
}

// SetTokenBucketConfigWithArgs encodes a call to the set_token_bucket_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rateLimiterEncoder) SetTokenBucketConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&mut TokenBucket",
		"bool",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_token_bucket_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IsEnabled encodes a call to the is_enabled Move function.
func (c rateLimiterEncoder) IsEnabled(bucket bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_enabled", typeArgsList, typeParamsList, []string{
		"&TokenBucket",
	}, []any{
		bucket,
	}, []string{
		"bool",
	})
}

// IsEnabledWithArgs encodes a call to the is_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rateLimiterEncoder) IsEnabledWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenBucket",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_enabled", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetRate encodes a call to the get_rate Move function.
func (c rateLimiterEncoder) GetRate(bucket bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_rate", typeArgsList, typeParamsList, []string{
		"&TokenBucket",
	}, []any{
		bucket,
	}, []string{
		"u64",
	})
}

// GetRateWithArgs encodes a call to the get_rate Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rateLimiterEncoder) GetRateWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenBucket",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_rate", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetTokenBucketFields encodes a call to the get_token_bucket_fields Move function.
func (c rateLimiterEncoder) GetTokenBucketFields(bucket bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_bucket_fields", typeArgsList, typeParamsList, []string{
		"&TokenBucket",
	}, []any{
		bucket,
	}, []string{
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	})
}

// GetTokenBucketFieldsWithArgs encodes a call to the get_token_bucket_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rateLimiterEncoder) GetTokenBucketFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenBucket",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_bucket_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	})
}
