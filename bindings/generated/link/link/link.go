// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_link

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

type ILink interface {
	MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error)
	Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ILinkDevInspect
	Encoder() LinkEncoder
	Bound() bind.IBoundContract
}

type ILinkDevInspect interface {
	Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (bind.Object, error)
}

type LinkEncoder interface {
	MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error)
	MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error)
	MintWithArgs(args ...any) (*bind.EncodedCall, error)
}

type LinkContract struct {
	*bind.BoundContract
	linkEncoder
	devInspect *LinkDevInspect
}

type LinkDevInspect struct {
	contract *LinkContract
}

var _ ILink = (*LinkContract)(nil)
var _ ILinkDevInspect = (*LinkDevInspect)(nil)

func NewLink(packageID string, client sui.ISuiAPI) (ILink, error) {
	contract, err := bind.NewBoundContract(packageID, "link", "link", client)
	if err != nil {
		return nil, err
	}

	c := &LinkContract{
		BoundContract: contract,
		linkEncoder:   linkEncoder{BoundContract: contract},
	}
	c.devInspect = &LinkDevInspect{contract: c}
	return c, nil
}

func (c *LinkContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *LinkContract) Encoder() LinkEncoder {
	return c.linkEncoder
}

func (c *LinkContract) DevInspect() ILinkDevInspect {
	return c.devInspect
}

type LINK struct {
}

func init() {
	bind.RegisterStructDecoder("link::link::LINK", func(data []byte) (interface{}, error) {
		var result LINK
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// MintAndTransfer executes the mint_and_transfer Move function.
func (c *LinkContract) MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.linkEncoder.MintAndTransfer(treasuryCap, amount, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Mint executes the mint Move function.
func (c *LinkContract) Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.linkEncoder.Mint(treasuryCap, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Mint executes the mint Move function using DevInspect to get return values.
//
// Returns: Coin<LINK>
func (d *LinkDevInspect) Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (bind.Object, error) {
	encoded, err := d.contract.linkEncoder.Mint(treasuryCap, amount)
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

type linkEncoder struct {
	*bind.BoundContract
}

// MintAndTransfer encodes a call to the mint_and_transfer Move function.
func (c linkEncoder) MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<LINK>",
		"u64",
		"address",
	}, []any{
		treasuryCap,
		amount,
		recipient,
	}, nil)
}

// MintAndTransferWithArgs encodes a call to the mint_and_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c linkEncoder) MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<LINK>",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Mint encodes a call to the mint Move function.
func (c linkEncoder) Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<LINK>",
		"u64",
	}, []any{
		treasuryCap,
		amount,
	}, []string{
		"Coin<LINK>",
	})
}

// MintWithArgs encodes a call to the mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c linkEncoder) MintWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<LINK>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Coin<LINK>",
	})
}
