// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mock_link_token

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

type IMockLinkToken interface {
	MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error)
	Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMockLinkTokenDevInspect
	Encoder() MockLinkTokenEncoder
	Bound() bind.IBoundContract
}

type IMockLinkTokenDevInspect interface {
}

type MockLinkTokenEncoder interface {
	MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error)
	MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error)
	MintWithArgs(args ...any) (*bind.EncodedCall, error)
}

type MockLinkTokenContract struct {
	*bind.BoundContract
	mockLinkTokenEncoder
	devInspect *MockLinkTokenDevInspect
}

type MockLinkTokenDevInspect struct {
	contract *MockLinkTokenContract
}

var _ IMockLinkToken = (*MockLinkTokenContract)(nil)
var _ IMockLinkTokenDevInspect = (*MockLinkTokenDevInspect)(nil)

func NewMockLinkToken(packageID string, client sui.ISuiAPI) (IMockLinkToken, error) {
	contract, err := bind.NewBoundContract(packageID, "mock_link_token", "mock_link_token", client)
	if err != nil {
		return nil, err
	}

	c := &MockLinkTokenContract{
		BoundContract:        contract,
		mockLinkTokenEncoder: mockLinkTokenEncoder{BoundContract: contract},
	}
	c.devInspect = &MockLinkTokenDevInspect{contract: c}
	return c, nil
}

func (c *MockLinkTokenContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *MockLinkTokenContract) Encoder() MockLinkTokenEncoder {
	return c.mockLinkTokenEncoder
}

func (c *MockLinkTokenContract) DevInspect() IMockLinkTokenDevInspect {
	return c.devInspect
}

type MOCK_LINK_TOKEN struct {
}

func init() {
	bind.RegisterStructDecoder("mock_link_token::mock_link_token::MOCK_LINK_TOKEN", func(data []byte) (interface{}, error) {
		var result MOCK_LINK_TOKEN
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// MintAndTransfer executes the mint_and_transfer Move function.
func (c *MockLinkTokenContract) MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mockLinkTokenEncoder.MintAndTransfer(treasuryCap, amount, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Mint executes the mint Move function.
func (c *MockLinkTokenContract) Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mockLinkTokenEncoder.Mint(treasuryCap, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type mockLinkTokenEncoder struct {
	*bind.BoundContract
}

// MintAndTransfer encodes a call to the mint_and_transfer Move function.
func (c mockLinkTokenEncoder) MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<MOCK_LINK_TOKEN>",
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
func (c mockLinkTokenEncoder) MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<MOCK_LINK_TOKEN>",
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
func (c mockLinkTokenEncoder) Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<MOCK_LINK_TOKEN>",
		"u64",
	}, []any{
		treasuryCap,
		amount,
	}, nil)
}

// MintWithArgs encodes a call to the mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mockLinkTokenEncoder) MintWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<MOCK_LINK_TOKEN>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, expectedParams, args, nil)
}
