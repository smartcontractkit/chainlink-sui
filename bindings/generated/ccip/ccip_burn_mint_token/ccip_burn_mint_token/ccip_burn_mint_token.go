// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_ccip_burn_mint_token

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

const FunctionInfo = `[{"package":"ccip_burn_mint_token","module":"ccip_burn_mint_token","name":"mint","parameters":[{"name":"treasury_cap","type":"TreasuryCap<CCIP_BURN_MINT_TOKEN>"},{"name":"amount","type":"u64"}]},{"package":"ccip_burn_mint_token","module":"ccip_burn_mint_token","name":"mint_and_transfer","parameters":[{"name":"treasury_cap","type":"TreasuryCap<CCIP_BURN_MINT_TOKEN>"},{"name":"amount","type":"u64"},{"name":"recipient","type":"address"}]}]`

type ICcipBurnMintToken interface {
	MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error)
	Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ICcipBurnMintTokenDevInspect
	Encoder() CcipBurnMintTokenEncoder
	Bound() bind.IBoundContract
}

type ICcipBurnMintTokenDevInspect interface {
}

type CcipBurnMintTokenEncoder interface {
	MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error)
	MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error)
	MintWithArgs(args ...any) (*bind.EncodedCall, error)
}

type CcipBurnMintTokenContract struct {
	*bind.BoundContract
	ccipBurnMintTokenEncoder
	devInspect *CcipBurnMintTokenDevInspect
}

type CcipBurnMintTokenDevInspect struct {
	contract *CcipBurnMintTokenContract
}

var _ ICcipBurnMintToken = (*CcipBurnMintTokenContract)(nil)
var _ ICcipBurnMintTokenDevInspect = (*CcipBurnMintTokenDevInspect)(nil)

func NewCcipBurnMintToken(packageID string, client sui.ISuiAPI) (ICcipBurnMintToken, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_burn_mint_token", "ccip_burn_mint_token", client)
	if err != nil {
		return nil, err
	}

	c := &CcipBurnMintTokenContract{
		BoundContract:            contract,
		ccipBurnMintTokenEncoder: ccipBurnMintTokenEncoder{BoundContract: contract},
	}
	c.devInspect = &CcipBurnMintTokenDevInspect{contract: c}
	return c, nil
}

func (c *CcipBurnMintTokenContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *CcipBurnMintTokenContract) Encoder() CcipBurnMintTokenEncoder {
	return c.ccipBurnMintTokenEncoder
}

func (c *CcipBurnMintTokenContract) DevInspect() ICcipBurnMintTokenDevInspect {
	return c.devInspect
}

type CCIP_BURN_MINT_TOKEN struct {
}

func init() {
	bind.RegisterStructDecoder("ccip_burn_mint_token::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", func(data []byte) (interface{}, error) {
		var result CCIP_BURN_MINT_TOKEN
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for CCIP_BURN_MINT_TOKEN
	bind.RegisterStructDecoder("vector<ccip_burn_mint_token::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN>", func(data []byte) (interface{}, error) {
		var results []CCIP_BURN_MINT_TOKEN
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// MintAndTransfer executes the mint_and_transfer Move function.
func (c *CcipBurnMintTokenContract) MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ccipBurnMintTokenEncoder.MintAndTransfer(treasuryCap, amount, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Mint executes the mint Move function.
func (c *CcipBurnMintTokenContract) Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ccipBurnMintTokenEncoder.Mint(treasuryCap, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type ccipBurnMintTokenEncoder struct {
	*bind.BoundContract
}

// MintAndTransfer encodes a call to the mint_and_transfer Move function.
func (c ccipBurnMintTokenEncoder) MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<CCIP_BURN_MINT_TOKEN>",
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
func (c ccipBurnMintTokenEncoder) MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<CCIP_BURN_MINT_TOKEN>",
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
func (c ccipBurnMintTokenEncoder) Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<CCIP_BURN_MINT_TOKEN>",
		"u64",
	}, []any{
		treasuryCap,
		amount,
	}, nil)
}

// MintWithArgs encodes a call to the mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ccipBurnMintTokenEncoder) MintWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<CCIP_BURN_MINT_TOKEN>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, expectedParams, args, nil)
}
