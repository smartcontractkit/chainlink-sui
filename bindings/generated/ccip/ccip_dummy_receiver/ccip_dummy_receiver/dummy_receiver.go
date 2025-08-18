// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_dummy_receiver

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

type IDummyReceiver interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error)
	GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error)
	Echo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, message []byte) (*models.SuiTransactionBlockResponse, error)
	CcipReceive(ctx context.Context, opts *bind.CallOpts, expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IDummyReceiverDevInspect
	Encoder() DummyReceiverEncoder
}

type IDummyReceiverDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (uint64, error)
	GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]TokenAmount, error)
	GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (string, error)
	GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (uint64, error)
	Echo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, message []byte) ([]byte, error)
}

type DummyReceiverEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterReceiver(ref bind.Object) (*bind.EncodedCall, error)
	RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCounter(state bind.Object) (*bind.EncodedCall, error)
	GetCounterWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestTokenAmounts(state bind.Object) (*bind.EncodedCall, error)
	GetDestTokenAmountsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenAmountToken(tokenAmount TokenAmount) (*bind.EncodedCall, error)
	GetTokenAmountTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenAmountAmount(tokenAmount TokenAmount) (*bind.EncodedCall, error)
	GetTokenAmountAmountWithArgs(args ...any) (*bind.EncodedCall, error)
	Echo(ref bind.Object, message []byte) (*bind.EncodedCall, error)
	EchoWithArgs(args ...any) (*bind.EncodedCall, error)
	CcipReceive(expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*bind.EncodedCall, error)
	CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error)
}

type DummyReceiverContract struct {
	*bind.BoundContract
	dummyReceiverEncoder
	devInspect *DummyReceiverDevInspect
}

type DummyReceiverDevInspect struct {
	contract *DummyReceiverContract
}

var _ IDummyReceiver = (*DummyReceiverContract)(nil)
var _ IDummyReceiverDevInspect = (*DummyReceiverDevInspect)(nil)

func NewDummyReceiver(packageID string, client sui.ISuiAPI) (*DummyReceiverContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_dummy_receiver", "dummy_receiver", client)
	if err != nil {
		return nil, err
	}

	c := &DummyReceiverContract{
		BoundContract:        contract,
		dummyReceiverEncoder: dummyReceiverEncoder{BoundContract: contract},
	}
	c.devInspect = &DummyReceiverDevInspect{contract: c}
	return c, nil
}

func (c *DummyReceiverContract) Encoder() DummyReceiverEncoder {
	return c.dummyReceiverEncoder
}

func (c *DummyReceiverContract) DevInspect() IDummyReceiverDevInspect {
	return c.devInspect
}

type OwnerCap struct {
	Id              string `move:"sui::object::UID"`
	ReceiverAddress string `move:"address"`
}

type ReceivedMessage struct {
	MessageId               []byte        `move:"vector<u8>"`
	SourceChainSelector     uint64        `move:"u64"`
	Sender                  []byte        `move:"vector<u8>"`
	Data                    []byte        `move:"vector<u8>"`
	DestTokenTransferLength uint64        `move:"u64"`
	DestTokenAmounts        []TokenAmount `move:"vector<TokenAmount>"`
}

type CCIPReceiverState struct {
	Id                      string        `move:"sui::object::UID"`
	Counter                 uint64        `move:"u64"`
	MessageId               []byte        `move:"vector<u8>"`
	SourceChainSelector     uint64        `move:"u64"`
	Sender                  []byte        `move:"vector<u8>"`
	Data                    []byte        `move:"vector<u8>"`
	DestTokenTransferLength uint64        `move:"u64"`
	DestTokenAmounts        []TokenAmount `move:"vector<TokenAmount>"`
}

type DummyReceiverProof struct {
}

type TokenAmount struct {
	Token  string `move:"address"`
	Amount uint64 `move:"u64"`
}

type bcsOwnerCap struct {
	Id              string
	ReceiverAddress [32]byte
}

func convertOwnerCapFromBCS(bcs bcsOwnerCap) (OwnerCap, error) {

	return OwnerCap{
		Id:              bcs.Id,
		ReceiverAddress: fmt.Sprintf("0x%x", bcs.ReceiverAddress),
	}, nil
}

type bcsTokenAmount struct {
	Token  [32]byte
	Amount uint64
}

func convertTokenAmountFromBCS(bcs bcsTokenAmount) TokenAmount {
	return TokenAmount{
		Token:  fmt.Sprintf("0x%x", bcs.Token),
		Amount: bcs.Amount,
	}
}

func init() {
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::OwnerCap", func(data []byte) (interface{}, error) {
		var temp bcsOwnerCap
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOwnerCapFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::ReceivedMessage", func(data []byte) (interface{}, error) {
		var result ReceivedMessage
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::CCIPReceiverState", func(data []byte) (interface{}, error) {
		var result CCIPReceiverState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::DummyReceiverProof", func(data []byte) (interface{}, error) {
		var result DummyReceiverProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::TokenAmount", func(data []byte) (interface{}, error) {
		var temp bcsTokenAmount
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertTokenAmountFromBCS(temp)
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *DummyReceiverContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterReceiver executes the register_receiver Move function.
func (c *DummyReceiverContract) RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.RegisterReceiver(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCounter executes the get_counter Move function.
func (c *DummyReceiverContract) GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetCounter(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestTokenAmounts executes the get_dest_token_amounts Move function.
func (c *DummyReceiverContract) GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetDestTokenAmounts(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenAmountToken executes the get_token_amount_token Move function.
func (c *DummyReceiverContract) GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetTokenAmountToken(tokenAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenAmountAmount executes the get_token_amount_amount Move function.
func (c *DummyReceiverContract) GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetTokenAmountAmount(tokenAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Echo executes the echo Move function.
func (c *DummyReceiverContract) Echo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, message []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.Echo(ref, message)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CcipReceive executes the ccip_receive Move function.
func (c *DummyReceiverContract) CcipReceive(ctx context.Context, opts *bind.CallOpts, expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.CcipReceive(expectedMessageId, ref, message, param, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *DummyReceiverDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.dummyReceiverEncoder.TypeAndVersion()
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

// GetCounter executes the get_counter Move function using DevInspect to get return values.
//
// Returns: u64
func (d *DummyReceiverDevInspect) GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (uint64, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetCounter(state)
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

// GetDestTokenAmounts executes the get_dest_token_amounts Move function using DevInspect to get return values.
//
// Returns: vector<TokenAmount>
func (d *DummyReceiverDevInspect) GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]TokenAmount, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetDestTokenAmounts(state)
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
	result, ok := results[0].([]TokenAmount)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []TokenAmount, got %T", results[0])
	}
	return result, nil
}

// GetTokenAmountToken executes the get_token_amount_token Move function using DevInspect to get return values.
//
// Returns: address
func (d *DummyReceiverDevInspect) GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (string, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetTokenAmountToken(tokenAmount)
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

// GetTokenAmountAmount executes the get_token_amount_amount Move function using DevInspect to get return values.
//
// Returns: u64
func (d *DummyReceiverDevInspect) GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (uint64, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetTokenAmountAmount(tokenAmount)
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

// Echo executes the echo Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *DummyReceiverDevInspect) Echo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, message []byte) ([]byte, error) {
	encoded, err := d.contract.dummyReceiverEncoder.Echo(ref, message)
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
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

type dummyReceiverEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c dummyReceiverEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// RegisterReceiver encodes a call to the register_receiver Move function.
func (c dummyReceiverEncoder) RegisterReceiver(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_receiver", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
	}, []any{
		ref,
	}, nil)
}

// RegisterReceiverWithArgs encodes a call to the register_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_receiver", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCounter encodes a call to the get_counter Move function.
func (c dummyReceiverEncoder) GetCounter(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_counter", typeArgsList, typeParamsList, []string{
		"&CCIPReceiverState",
	}, []any{
		state,
	}, []string{
		"u64",
	})
}

// GetCounterWithArgs encodes a call to the get_counter Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetCounterWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_counter", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetDestTokenAmounts encodes a call to the get_dest_token_amounts Move function.
func (c dummyReceiverEncoder) GetDestTokenAmounts(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_token_amounts", typeArgsList, typeParamsList, []string{
		"&CCIPReceiverState",
	}, []any{
		state,
	}, []string{
		"vector<ccip_dummy_receiver::dummy_receiver::TokenAmount>",
	})
}

// GetDestTokenAmountsWithArgs encodes a call to the get_dest_token_amounts Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetDestTokenAmountsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_token_amounts", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ccip_dummy_receiver::dummy_receiver::TokenAmount>",
	})
}

// GetTokenAmountToken encodes a call to the get_token_amount_token Move function.
func (c dummyReceiverEncoder) GetTokenAmountToken(tokenAmount TokenAmount) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_token", typeArgsList, typeParamsList, []string{
		"&TokenAmount",
	}, []any{
		tokenAmount,
	}, []string{
		"address",
	})
}

// GetTokenAmountTokenWithArgs encodes a call to the get_token_amount_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetTokenAmountTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenAmount",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetTokenAmountAmount encodes a call to the get_token_amount_amount Move function.
func (c dummyReceiverEncoder) GetTokenAmountAmount(tokenAmount TokenAmount) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_amount", typeArgsList, typeParamsList, []string{
		"&TokenAmount",
	}, []any{
		tokenAmount,
	}, []string{
		"u64",
	})
}

// GetTokenAmountAmountWithArgs encodes a call to the get_token_amount_amount Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetTokenAmountAmountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenAmount",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_amount", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// Echo encodes a call to the echo Move function.
func (c dummyReceiverEncoder) Echo(ref bind.Object, message []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("echo", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"vector<u8>",
	}, []any{
		ref,
		message,
	}, []string{
		"vector<u8>",
	})
}

// EchoWithArgs encodes a call to the echo Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) EchoWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("echo", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// CcipReceive encodes a call to the ccip_receive Move function.
func (c dummyReceiverEncoder) CcipReceive(expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"&CCIPObjectRef",
		"client::Any2SuiMessage",
		"&Clock",
		"&mut CCIPReceiverState",
	}, []any{
		expectedMessageId,
		ref,
		message,
		param,
		state,
	}, nil)
}

// CcipReceiveWithArgs encodes a call to the ccip_receive Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"&CCIPObjectRef",
		"client::Any2SuiMessage",
		"&Clock",
		"&mut CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, expectedParams, args, nil)
}
