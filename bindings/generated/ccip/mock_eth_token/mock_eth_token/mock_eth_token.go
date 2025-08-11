// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mock_eth_token

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

type IMockEthToken interface {
	MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error)
	Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMockEthTokenDevInspect
	Encoder() MockEthTokenEncoder
}

type IMockEthTokenDevInspect interface {
}

type MockEthTokenEncoder interface {
	MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error)
	MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error)
	MintWithArgs(args ...any) (*bind.EncodedCall, error)
}

type MockEthTokenContract struct {
	*bind.BoundContract
	mockEthTokenEncoder
	devInspect *MockEthTokenDevInspect
}

type MockEthTokenDevInspect struct {
	contract *MockEthTokenContract
}

var _ IMockEthToken = (*MockEthTokenContract)(nil)
var _ IMockEthTokenDevInspect = (*MockEthTokenDevInspect)(nil)

func NewMockEthToken(packageID string, client sui.ISuiAPI) (*MockEthTokenContract, error) {
	contract, err := bind.NewBoundContract(packageID, "mock_eth_token", "mock_eth_token", client)
	if err != nil {
		return nil, err
	}

	c := &MockEthTokenContract{
		BoundContract:       contract,
		mockEthTokenEncoder: mockEthTokenEncoder{BoundContract: contract},
	}
	c.devInspect = &MockEthTokenDevInspect{contract: c}
	return c, nil
}

func (c *MockEthTokenContract) Encoder() MockEthTokenEncoder {
	return c.mockEthTokenEncoder
}

func (c *MockEthTokenContract) DevInspect() IMockEthTokenDevInspect {
	return c.devInspect
}

func (c *MockEthTokenContract) BuildPTB(ctx context.Context, ptb *transaction.Transaction, encoded *bind.EncodedCall) (*transaction.Argument, error) {
	var callArgManager *bind.CallArgManager
	if ptb.Data.V1 != nil && ptb.Data.V1.Kind.ProgrammableTransaction != nil &&
		ptb.Data.V1.Kind.ProgrammableTransaction.Inputs != nil {
		callArgManager = bind.NewCallArgManagerWithExisting(ptb.Data.V1.Kind.ProgrammableTransaction.Inputs)
	} else {
		callArgManager = bind.NewCallArgManager()
	}

	arguments, err := callArgManager.ConvertEncodedCallArgsToArguments(encoded.CallArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to convert EncodedCallArguments to Arguments: %w", err)
	}

	ptb.Data.V1.Kind.ProgrammableTransaction.Inputs = callArgManager.GetInputs()

	typeTagValues := make([]transaction.TypeTag, len(encoded.TypeArgs))
	for i, tag := range encoded.TypeArgs {
		if tag != nil {
			typeTagValues[i] = *tag
		}
	}

	argumentValues := make([]transaction.Argument, len(arguments))
	for i, arg := range arguments {
		if arg != nil {
			argumentValues[i] = *arg
		}
	}

	result := ptb.MoveCall(
		models.SuiAddress(encoded.Module.PackageID),
		encoded.Module.ModuleName,
		encoded.Function,
		typeTagValues,
		argumentValues,
	)

	return &result, nil
}

type MOCK_ETH_TOKEN struct {
}

func init() {
	bind.RegisterStructDecoder("mock_eth_token::mock_eth_token::MOCK_ETH_TOKEN", func(data []byte) (interface{}, error) {
		var result MOCK_ETH_TOKEN
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// MintAndTransfer executes the mint_and_transfer Move function.
func (c *MockEthTokenContract) MintAndTransfer(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mockEthTokenEncoder.MintAndTransfer(treasuryCap, amount, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Mint executes the mint Move function.
func (c *MockEthTokenContract) Mint(ctx context.Context, opts *bind.CallOpts, treasuryCap bind.Object, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mockEthTokenEncoder.Mint(treasuryCap, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type mockEthTokenEncoder struct {
	*bind.BoundContract
}

// MintAndTransfer encodes a call to the mint_and_transfer Move function.
func (c mockEthTokenEncoder) MintAndTransfer(treasuryCap bind.Object, amount uint64, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint_and_transfer", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<MOCK_ETH_TOKEN>",
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
func (c mockEthTokenEncoder) MintAndTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<MOCK_ETH_TOKEN>",
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
func (c mockEthTokenEncoder) Mint(treasuryCap bind.Object, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, []string{
		"&mut TreasuryCap<MOCK_ETH_TOKEN>",
		"u64",
	}, []any{
		treasuryCap,
		amount,
	}, nil)
}

// MintWithArgs encodes a call to the mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mockEthTokenEncoder) MintWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TreasuryCap<MOCK_ETH_TOKEN>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mint", typeArgsList, typeParamsList, expectedParams, args, nil)
}
