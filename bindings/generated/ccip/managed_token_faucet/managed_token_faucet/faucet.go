// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_faucet

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

const FunctionInfo = `[{"package":"managed_token_faucet","module":"faucet","name":"drip","parameters":[{"name":"metadata","type":"CoinMetadata<T>"},{"name":"state","type":"FaucetState<T>"},{"name":"token_state","type":"TokenState<T>"},{"name":"deny_list","type":"DenyList"}]},{"package":"managed_token_faucet","module":"faucet","name":"drip_and_send","parameters":[{"name":"metadata","type":"CoinMetadata<T>"},{"name":"state","type":"FaucetState<T>"},{"name":"token_state","type":"TokenState<T>"},{"name":"deny_list","type":"DenyList"},{"name":"recipient","type":"address"}]},{"package":"managed_token_faucet","module":"faucet","name":"initialize","parameters":[{"name":"mint_cap","type":"MintCap<T>"}]},{"package":"managed_token_faucet","module":"faucet","name":"type_and_version","parameters":null}]`

type IFaucet interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, mintCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	Drip(ctx context.Context, opts *bind.CallOpts, typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object) (*models.SuiTransactionBlockResponse, error)
	DripAndSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IFaucetDevInspect
	Encoder() FaucetEncoder
	Bound() bind.IBoundContract
}

type IFaucetDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	Drip(ctx context.Context, opts *bind.CallOpts, typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object) (any, error)
}

type FaucetEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(typeArgs []string, mintCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Drip(typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object) (*bind.EncodedCall, error)
	DripWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	DripAndSend(typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object, recipient string) (*bind.EncodedCall, error)
	DripAndSendWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
}

type FaucetContract struct {
	*bind.BoundContract
	faucetEncoder
	devInspect *FaucetDevInspect
}

type FaucetDevInspect struct {
	contract *FaucetContract
}

var _ IFaucet = (*FaucetContract)(nil)
var _ IFaucetDevInspect = (*FaucetDevInspect)(nil)

func NewFaucet(packageID string, client sui.ISuiAPI) (IFaucet, error) {
	contract, err := bind.NewBoundContract(packageID, "managed_token_faucet", "faucet", client)
	if err != nil {
		return nil, err
	}

	c := &FaucetContract{
		BoundContract: contract,
		faucetEncoder: faucetEncoder{BoundContract: contract},
	}
	c.devInspect = &FaucetDevInspect{contract: c}
	return c, nil
}

func (c *FaucetContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *FaucetContract) Encoder() FaucetEncoder {
	return c.faucetEncoder
}

func (c *FaucetContract) DevInspect() IFaucetDevInspect {
	return c.devInspect
}

type FaucetState struct {
	Id      string      `move:"sui::object::UID"`
	MintCap bind.Object `move:"MintCap<T>"`
}

func init() {
	bind.RegisterStructDecoder("managed_token_faucet::faucet::FaucetState", func(data []byte) (interface{}, error) {
		var result FaucetState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for FaucetState
	bind.RegisterStructDecoder("vector<managed_token_faucet::faucet::FaucetState>", func(data []byte) (interface{}, error) {
		var results []FaucetState
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *FaucetContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.faucetEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *FaucetContract) Initialize(ctx context.Context, opts *bind.CallOpts, typeArgs []string, mintCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.faucetEncoder.Initialize(typeArgs, mintCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Drip executes the drip Move function.
func (c *FaucetContract) Drip(ctx context.Context, opts *bind.CallOpts, typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.faucetEncoder.Drip(typeArgs, metadata, state, tokenState, denyList)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DripAndSend executes the drip_and_send Move function.
func (c *FaucetContract) DripAndSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.faucetEncoder.DripAndSend(typeArgs, metadata, state, tokenState, denyList, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *FaucetDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.faucetEncoder.TypeAndVersion()
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// Drip executes the drip Move function using DevInspect to get return values.
//
// Returns: Coin<T>
func (d *FaucetDevInspect) Drip(ctx context.Context, opts *bind.CallOpts, typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object) (any, error) {
	encoded, err := d.contract.faucetEncoder.Drip(typeArgs, metadata, state, tokenState, denyList)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	return results[0], nil
}

type faucetEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c faucetEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c faucetEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// Initialize encodes a call to the initialize Move function.
func (c faucetEncoder) Initialize(typeArgs []string, mintCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"MintCap<T>",
	}, []any{
		mintCap,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c faucetEncoder) InitializeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"MintCap<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Drip encodes a call to the drip Move function.
func (c faucetEncoder) Drip(typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("drip", typeArgsList, typeParamsList, []string{
		"&CoinMetadata<T>",
		"&FaucetState<T>",
		"&mut TokenState<T>",
		"&DenyList",
	}, []any{
		metadata,
		state,
		tokenState,
		denyList,
	}, []string{
		"Coin<T>",
	})
}

// DripWithArgs encodes a call to the drip Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c faucetEncoder) DripWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CoinMetadata<T>",
		"&FaucetState<T>",
		"&mut TokenState<T>",
		"&DenyList",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("drip", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Coin<T>",
	})
}

// DripAndSend encodes a call to the drip_and_send Move function.
func (c faucetEncoder) DripAndSend(typeArgs []string, metadata bind.Object, state bind.Object, tokenState bind.Object, denyList bind.Object, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("drip_and_send", typeArgsList, typeParamsList, []string{
		"&CoinMetadata<T>",
		"&FaucetState<T>",
		"&mut TokenState<T>",
		"&DenyList",
		"address",
	}, []any{
		metadata,
		state,
		tokenState,
		denyList,
		recipient,
	}, nil)
}

// DripAndSendWithArgs encodes a call to the drip_and_send Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c faucetEncoder) DripAndSendWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CoinMetadata<T>",
		"&FaucetState<T>",
		"&mut TokenState<T>",
		"&DenyList",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("drip_and_send", typeArgsList, typeParamsList, expectedParams, args, nil)
}
