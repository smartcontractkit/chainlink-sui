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
	"github.com/block-vision/sui-go-sdk/transaction"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

type IDummyReceiver interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverStateId string, receiverStateParams []string) (*models.SuiTransactionBlockResponse, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	CcipReceive(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object) (*models.SuiTransactionBlockResponse, error)
	CcipReceive1(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	CcipReceive2(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object, param bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IDummyReceiverDevInspect
	Encoder() DummyReceiverEncoder
}

type IDummyReceiverDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (uint64, error)
	CcipReceive(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object) (bind.Object, error)
	CcipReceive1(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object) (bind.Object, error)
	CcipReceive2(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object, param bind.Object) (bind.Object, error)
}

type DummyReceiverEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterReceiver(ref bind.Object, receiverStateId string, receiverStateParams []string) (*bind.EncodedCall, error)
	RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCounter(state bind.Object) (*bind.EncodedCall, error)
	GetCounterWithArgs(args ...any) (*bind.EncodedCall, error)
	CcipReceive(ref bind.Object, packageId string, receiverParams bind.Object) (*bind.EncodedCall, error)
	CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error)
	CcipReceive1(ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object) (*bind.EncodedCall, error)
	CcipReceive1WithArgs(args ...any) (*bind.EncodedCall, error)
	CcipReceive2(ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object, param bind.Object) (*bind.EncodedCall, error)
	CcipReceive2WithArgs(args ...any) (*bind.EncodedCall, error)
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

func (c *DummyReceiverContract) BuildPTB(ctx context.Context, ptb *transaction.Transaction, encoded *bind.EncodedCall) (*transaction.Argument, error) {
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

type OwnerCap struct {
	Id              string `move:"sui::object::UID"`
	ReceiverAddress string `move:"address"`
}

type ReceivedMessage struct {
	MessageId               []byte `move:"vector<u8>"`
	SourceChainSelector     uint64 `move:"u64"`
	Sender                  []byte `move:"vector<u8>"`
	Data                    []byte `move:"vector<u8>"`
	DestTokenTransferLength uint64 `move:"u64"`
}

type CCIPReceiverState struct {
	Id                      string `move:"sui::object::UID"`
	Counter                 uint64 `move:"u64"`
	MessageId               []byte `move:"vector<u8>"`
	SourceChainSelector     uint64 `move:"u64"`
	Sender                  []byte `move:"vector<u8>"`
	Data                    []byte `move:"vector<u8>"`
	DestTokenTransferLength uint64 `move:"u64"`
}

type DummyReceiverProof struct {
}

type bcsOwnerCap struct {
	Id              string
	ReceiverAddress [32]byte
}

func convertOwnerCapFromBCS(bcs bcsOwnerCap) OwnerCap {
	return OwnerCap{
		Id:              bcs.Id,
		ReceiverAddress: fmt.Sprintf("0x%x", bcs.ReceiverAddress),
	}
}

func init() {
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::OwnerCap", func(data []byte) (interface{}, error) {
		var temp bcsOwnerCap
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertOwnerCapFromBCS(temp)
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
func (c *DummyReceiverContract) RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverStateId string, receiverStateParams []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.RegisterReceiver(ref, receiverStateId, receiverStateParams)
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

// CcipReceive executes the ccip_receive Move function.
func (c *DummyReceiverContract) CcipReceive(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.CcipReceive(ref, packageId, receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CcipReceive1 executes the ccip_receive_1 Move function.
func (c *DummyReceiverContract) CcipReceive1(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.CcipReceive1(ref, packageId, receiverParams, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CcipReceive2 executes the ccip_receive_2 Move function.
func (c *DummyReceiverContract) CcipReceive2(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object, param bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.CcipReceive2(ref, packageId, receiverParams, state, param)
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

// CcipReceive executes the ccip_receive Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *DummyReceiverDevInspect) CcipReceive(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object) (bind.Object, error) {
	encoded, err := d.contract.dummyReceiverEncoder.CcipReceive(ref, packageId, receiverParams)
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

// CcipReceive1 executes the ccip_receive_1 Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *DummyReceiverDevInspect) CcipReceive1(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object) (bind.Object, error) {
	encoded, err := d.contract.dummyReceiverEncoder.CcipReceive1(ref, packageId, receiverParams, state)
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

// CcipReceive2 executes the ccip_receive_2 Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *DummyReceiverDevInspect) CcipReceive2(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object, param bind.Object) (bind.Object, error) {
	encoded, err := d.contract.dummyReceiverEncoder.CcipReceive2(ref, packageId, receiverParams, state, param)
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
func (c dummyReceiverEncoder) RegisterReceiver(ref bind.Object, receiverStateId string, receiverStateParams []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_receiver", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"address",
		"vector<address>",
	}, []any{
		ref,
		receiverStateId,
		receiverStateParams,
	}, nil)
}

// RegisterReceiverWithArgs encodes a call to the register_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"address",
		"vector<address>",
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

// CcipReceive encodes a call to the ccip_receive Move function.
func (c dummyReceiverEncoder) CcipReceive(ref bind.Object, packageId string, receiverParams bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"osh::ReceiverParams",
	}, []any{
		ref,
		packageId,
		receiverParams,
	}, []string{
		"osh::ReceiverParams",
	})
}

// CcipReceiveWithArgs encodes a call to the ccip_receive Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"osh::ReceiverParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}

// CcipReceive1 encodes a call to the ccip_receive_1 Move function.
func (c dummyReceiverEncoder) CcipReceive1(ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive_1", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"osh::ReceiverParams",
		"&mut CCIPReceiverState",
	}, []any{
		ref,
		packageId,
		receiverParams,
		state,
	}, []string{
		"osh::ReceiverParams",
	})
}

// CcipReceive1WithArgs encodes a call to the ccip_receive_1 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) CcipReceive1WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"osh::ReceiverParams",
		"&mut CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive_1", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}

// CcipReceive2 encodes a call to the ccip_receive_2 Move function.
func (c dummyReceiverEncoder) CcipReceive2(ref bind.Object, packageId string, receiverParams bind.Object, state bind.Object, param bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive_2", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"osh::ReceiverParams",
		"&mut CCIPReceiverState",
		"&Clock",
	}, []any{
		ref,
		packageId,
		receiverParams,
		state,
		param,
	}, []string{
		"osh::ReceiverParams",
	})
}

// CcipReceive2WithArgs encodes a call to the ccip_receive_2 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) CcipReceive2WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"osh::ReceiverParams",
		"&mut CCIPReceiverState",
		"&Clock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive_2", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}
