// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_generics

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

type IGenerics interface {
	CreateBox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, value bind.Object) (*models.SuiTransactionBlockResponse, error)
	Unbox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, box bind.Object) (*models.SuiTransactionBlockResponse, error)
	Deposit(ctx context.Context, opts *bind.CallOpts, typeArgs []string, token bind.Object, coin bind.Object) (*models.SuiTransactionBlockResponse, error)
	Balance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, token bind.Object) (*models.SuiTransactionBlockResponse, error)
	CreatePair(ctx context.Context, opts *bind.CallOpts, typeArgs []string, first bind.Object, second bind.Object) (*models.SuiTransactionBlockResponse, error)
	CreateSuiToken(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	CreateAndTransferSuiToken(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	CreateAndTransferToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (*models.SuiTransactionBlockResponse, error)
	CreateAndTransferBox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, value bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IGenericsDevInspect
	Encoder() GenericsEncoder
}

type IGenericsDevInspect interface {
	CreateBox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, value bind.Object) (any, error)
	Unbox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, box bind.Object) (any, error)
	Balance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, token bind.Object) (uint64, error)
	CreatePair(ctx context.Context, opts *bind.CallOpts, typeArgs []string, first bind.Object, second bind.Object) (any, error)
	CreateSuiToken(ctx context.Context, opts *bind.CallOpts) (bind.Object, error)
}

type GenericsEncoder interface {
	CreateBox(typeArgs []string, value bind.Object) (*bind.EncodedCall, error)
	CreateBoxWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Unbox(typeArgs []string, box bind.Object) (*bind.EncodedCall, error)
	UnboxWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Deposit(typeArgs []string, token bind.Object, coin bind.Object) (*bind.EncodedCall, error)
	DepositWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Balance(typeArgs []string, token bind.Object) (*bind.EncodedCall, error)
	BalanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	CreatePair(typeArgs []string, first bind.Object, second bind.Object) (*bind.EncodedCall, error)
	CreatePairWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	CreateSuiToken() (*bind.EncodedCall, error)
	CreateSuiTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateAndTransferSuiToken() (*bind.EncodedCall, error)
	CreateAndTransferSuiTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateAndTransferToken(typeArgs []string) (*bind.EncodedCall, error)
	CreateAndTransferTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	CreateAndTransferBox(typeArgs []string, value bind.Object) (*bind.EncodedCall, error)
	CreateAndTransferBoxWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
}

type GenericsContract struct {
	*bind.BoundContract
	genericsEncoder
	devInspect *GenericsDevInspect
}

type GenericsDevInspect struct {
	contract *GenericsContract
}

var _ IGenerics = (*GenericsContract)(nil)
var _ IGenericsDevInspect = (*GenericsDevInspect)(nil)

func NewGenerics(packageID string, client sui.ISuiAPI) (*GenericsContract, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "generics", client, nil)
	if err != nil {
		return nil, err
	}

	c := &GenericsContract{
		BoundContract:   contract,
		genericsEncoder: genericsEncoder{BoundContract: contract},
	}
	c.devInspect = &GenericsDevInspect{contract: c}
	return c, nil
}

func (c *GenericsContract) Encoder() GenericsEncoder {
	return c.genericsEncoder
}

func (c *GenericsContract) DevInspect() IGenericsDevInspect {
	return c.devInspect
}

type Box struct {
	Id    string      `move:"sui::object::UID"`
	Value bind.Object `move:"T"`
}

type Token struct {
	Id      string `move:"sui::object::UID"`
	Balance uint64 `move:"u64"`
}

type Pair struct {
	First  bind.Object `move:"T"`
	Second bind.Object `move:"U"`
}

func init() {
	bind.RegisterStructDecoder("test::generics::Box", func(data []byte) (interface{}, error) {
		var result Box
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::generics::Token", func(data []byte) (interface{}, error) {
		var result Token
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::generics::Pair", func(data []byte) (interface{}, error) {
		var result Pair
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// CreateBox executes the create_box Move function.
func (c *GenericsContract) CreateBox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, value bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.CreateBox(typeArgs, value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Unbox executes the unbox Move function.
func (c *GenericsContract) Unbox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, box bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.Unbox(typeArgs, box)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Deposit executes the deposit Move function.
func (c *GenericsContract) Deposit(ctx context.Context, opts *bind.CallOpts, typeArgs []string, token bind.Object, coin bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.Deposit(typeArgs, token, coin)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Balance executes the balance Move function.
func (c *GenericsContract) Balance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, token bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.Balance(typeArgs, token)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePair executes the create_pair Move function.
func (c *GenericsContract) CreatePair(ctx context.Context, opts *bind.CallOpts, typeArgs []string, first bind.Object, second bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.CreatePair(typeArgs, first, second)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateSuiToken executes the create_sui_token Move function.
func (c *GenericsContract) CreateSuiToken(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.CreateSuiToken()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateAndTransferSuiToken executes the create_and_transfer_sui_token Move function.
func (c *GenericsContract) CreateAndTransferSuiToken(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.CreateAndTransferSuiToken()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateAndTransferToken executes the create_and_transfer_token Move function.
func (c *GenericsContract) CreateAndTransferToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.CreateAndTransferToken(typeArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateAndTransferBox executes the create_and_transfer_box Move function.
func (c *GenericsContract) CreateAndTransferBox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, value bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.genericsEncoder.CreateAndTransferBox(typeArgs, value)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateBox executes the create_box Move function using DevInspect to get return values.
//
// Returns: Box<T>
func (d *GenericsDevInspect) CreateBox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, value bind.Object) (any, error) {
	encoded, err := d.contract.genericsEncoder.CreateBox(typeArgs, value)
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

// Unbox executes the unbox Move function using DevInspect to get return values.
//
// Returns: T
func (d *GenericsDevInspect) Unbox(ctx context.Context, opts *bind.CallOpts, typeArgs []string, box bind.Object) (any, error) {
	encoded, err := d.contract.genericsEncoder.Unbox(typeArgs, box)
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

// Balance executes the balance Move function using DevInspect to get return values.
//
// Returns: u64
func (d *GenericsDevInspect) Balance(ctx context.Context, opts *bind.CallOpts, typeArgs []string, token bind.Object) (uint64, error) {
	encoded, err := d.contract.genericsEncoder.Balance(typeArgs, token)
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

// CreatePair executes the create_pair Move function using DevInspect to get return values.
//
// Returns: Pair<T, U>
func (d *GenericsDevInspect) CreatePair(ctx context.Context, opts *bind.CallOpts, typeArgs []string, first bind.Object, second bind.Object) (any, error) {
	encoded, err := d.contract.genericsEncoder.CreatePair(typeArgs, first, second)
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

// CreateSuiToken executes the create_sui_token Move function using DevInspect to get return values.
//
// Returns: Token<sui::sui::SUI>
func (d *GenericsDevInspect) CreateSuiToken(ctx context.Context, opts *bind.CallOpts) (bind.Object, error) {
	encoded, err := d.contract.genericsEncoder.CreateSuiToken()
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

type genericsEncoder struct {
	*bind.BoundContract
}

// CreateBox encodes a call to the create_box Move function.
func (c genericsEncoder) CreateBox(typeArgs []string, value bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_box", typeArgsList, typeParamsList, []string{
		"T",
	}, []any{
		value,
	}, []string{
		"Box<T>",
	})
}

// CreateBoxWithArgs encodes a call to the create_box Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) CreateBoxWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"T",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_box", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Box<T>",
	})
}

// Unbox encodes a call to the unbox Move function.
func (c genericsEncoder) Unbox(typeArgs []string, box bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("unbox", typeArgsList, typeParamsList, []string{
		"Box<T>",
	}, []any{
		box,
	}, []string{
		"T",
	})
}

// UnboxWithArgs encodes a call to the unbox Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) UnboxWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"Box<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("unbox", typeArgsList, typeParamsList, expectedParams, args, []string{
		"T",
	})
}

// Deposit encodes a call to the deposit Move function.
func (c genericsEncoder) Deposit(typeArgs []string, token bind.Object, coin bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("deposit", typeArgsList, typeParamsList, []string{
		"&mut Token<T>",
		"Token<T>",
	}, []any{
		token,
		coin,
	}, nil)
}

// DepositWithArgs encodes a call to the deposit Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) DepositWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Token<T>",
		"Token<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("deposit", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Balance encodes a call to the balance Move function.
func (c genericsEncoder) Balance(typeArgs []string, token bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("balance", typeArgsList, typeParamsList, []string{
		"&Token<T>",
	}, []any{
		token,
	}, []string{
		"u64",
	})
}

// BalanceWithArgs encodes a call to the balance Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) BalanceWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Token<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("balance", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// CreatePair encodes a call to the create_pair Move function.
func (c genericsEncoder) CreatePair(typeArgs []string, first bind.Object, second bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"U",
	}
	return c.EncodeCallArgsWithGenerics("create_pair", typeArgsList, typeParamsList, []string{
		"T",
		"U",
	}, []any{
		first,
		second,
	}, []string{
		"Pair<T, U>",
	})
}

// CreatePairWithArgs encodes a call to the create_pair Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) CreatePairWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"T",
		"U",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"U",
	}
	return c.EncodeCallArgsWithGenerics("create_pair", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Pair<T, U>",
	})
}

// CreateSuiToken encodes a call to the create_sui_token Move function.
func (c genericsEncoder) CreateSuiToken() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_sui_token", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"Token<sui::sui::SUI>",
	})
}

// CreateSuiTokenWithArgs encodes a call to the create_sui_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) CreateSuiTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_sui_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Token<sui::sui::SUI>",
	})
}

// CreateAndTransferSuiToken encodes a call to the create_and_transfer_sui_token Move function.
func (c genericsEncoder) CreateAndTransferSuiToken() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_and_transfer_sui_token", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// CreateAndTransferSuiTokenWithArgs encodes a call to the create_and_transfer_sui_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) CreateAndTransferSuiTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_and_transfer_sui_token", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CreateAndTransferToken encodes a call to the create_and_transfer_token Move function.
func (c genericsEncoder) CreateAndTransferToken(typeArgs []string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_and_transfer_token", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// CreateAndTransferTokenWithArgs encodes a call to the create_and_transfer_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) CreateAndTransferTokenWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_and_transfer_token", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CreateAndTransferBox encodes a call to the create_and_transfer_box Move function.
func (c genericsEncoder) CreateAndTransferBox(typeArgs []string, value bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_and_transfer_box", typeArgsList, typeParamsList, []string{
		"T",
	}, []any{
		value,
	}, nil)
}

// CreateAndTransferBoxWithArgs encodes a call to the create_and_transfer_box Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c genericsEncoder) CreateAndTransferBoxWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"T",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_and_transfer_box", typeArgsList, typeParamsList, expectedParams, args, nil)
}
