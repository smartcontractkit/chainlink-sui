// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_complex

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

type IComplex interface {
	NewObjectWithTransfer(ctx context.Context, opts *bind.CallOpts, someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*models.SuiTransactionBlockResponse, error)
	NewObject(ctx context.Context, opts *bind.CallOpts, someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*models.SuiTransactionBlockResponse, error)
	FlattenAddress(ctx context.Context, opts *bind.CallOpts, someAddress string, someAddresses []string) (*models.SuiTransactionBlockResponse, error)
	FlattenU8(ctx context.Context, opts *bind.CallOpts, input [][]byte) (*models.SuiTransactionBlockResponse, error)
	CheckU128(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*models.SuiTransactionBlockResponse, error)
	CheckU256(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*models.SuiTransactionBlockResponse, error)
	CheckWithObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object) (*models.SuiTransactionBlockResponse, error)
	CheckWithMutObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object, newNumber uint64) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IComplexDevInspect
	Encoder() ComplexEncoder
}

type IComplexDevInspect interface {
	NewObject(ctx context.Context, opts *bind.CallOpts, someId []byte, someNumber uint64, someAddress string, someAddresses []string) (DroppableObject, error)
	FlattenAddress(ctx context.Context, opts *bind.CallOpts, someAddress string, someAddresses []string) ([]string, error)
	FlattenU8(ctx context.Context, opts *bind.CallOpts, input [][]byte) ([]byte, error)
	CheckU128(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*big.Int, error)
	CheckU256(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*big.Int, error)
	CheckWithObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object) (uint64, error)
	CheckWithMutObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object, newNumber uint64) (uint64, error)
}

type ComplexEncoder interface {
	NewObjectWithTransfer(someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*bind.EncodedCall, error)
	NewObjectWithTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	NewObject(someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*bind.EncodedCall, error)
	NewObjectWithArgs(args ...any) (*bind.EncodedCall, error)
	FlattenAddress(someAddress string, someAddresses []string) (*bind.EncodedCall, error)
	FlattenAddressWithArgs(args ...any) (*bind.EncodedCall, error)
	FlattenU8(input [][]byte) (*bind.EncodedCall, error)
	FlattenU8WithArgs(args ...any) (*bind.EncodedCall, error)
	CheckU128(input *big.Int) (*bind.EncodedCall, error)
	CheckU128WithArgs(args ...any) (*bind.EncodedCall, error)
	CheckU256(input *big.Int) (*bind.EncodedCall, error)
	CheckU256WithArgs(args ...any) (*bind.EncodedCall, error)
	CheckWithObjectRef(obj bind.Object) (*bind.EncodedCall, error)
	CheckWithObjectRefWithArgs(args ...any) (*bind.EncodedCall, error)
	CheckWithMutObjectRef(obj bind.Object, newNumber uint64) (*bind.EncodedCall, error)
	CheckWithMutObjectRefWithArgs(args ...any) (*bind.EncodedCall, error)
}

type ComplexContract struct {
	*bind.BoundContract
	complexEncoder
	devInspect *ComplexDevInspect
}

type ComplexDevInspect struct {
	contract *ComplexContract
}

var _ IComplex = (*ComplexContract)(nil)
var _ IComplexDevInspect = (*ComplexDevInspect)(nil)

func NewComplex(packageID string, client sui.ISuiAPI) (*ComplexContract, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "complex", client)
	if err != nil {
		return nil, err
	}

	c := &ComplexContract{
		BoundContract:  contract,
		complexEncoder: complexEncoder{BoundContract: contract},
	}
	c.devInspect = &ComplexDevInspect{contract: c}
	return c, nil
}

func (c *ComplexContract) Encoder() ComplexEncoder {
	return c.complexEncoder
}

func (c *ComplexContract) DevInspect() IComplexDevInspect {
	return c.devInspect
}

func (c *ComplexContract) BuildPTB(ctx context.Context, ptb *transaction.Transaction, encoded *bind.EncodedCall) (*transaction.Argument, error) {
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

type SampleObject struct {
	Id            string   `move:"sui::object::UID"`
	SomeId        []byte   `move:"vector<u8>"`
	SomeNumber    uint64   `move:"u64"`
	SomeAddress   string   `move:"address"`
	SomeAddresses []string `move:"vector<address>"`
}

type DroppableObject struct {
	SomeId        []byte   `move:"vector<u8>"`
	SomeNumber    uint64   `move:"u64"`
	SomeAddress   string   `move:"address"`
	SomeAddresses []string `move:"vector<address>"`
}

type bcsSampleObject struct {
	Id            string
	SomeId        []byte
	SomeNumber    uint64
	SomeAddress   [32]byte
	SomeAddresses [][32]byte
}

func convertSampleObjectFromBCS(bcs bcsSampleObject) SampleObject {
	return SampleObject{
		Id:          bcs.Id,
		SomeId:      bcs.SomeId,
		SomeNumber:  bcs.SomeNumber,
		SomeAddress: fmt.Sprintf("0x%x", bcs.SomeAddress),
		SomeAddresses: func() []string {
			addrs := make([]string, len(bcs.SomeAddresses))
			for i, addr := range bcs.SomeAddresses {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}
}

type bcsDroppableObject struct {
	SomeId        []byte
	SomeNumber    uint64
	SomeAddress   [32]byte
	SomeAddresses [][32]byte
}

func convertDroppableObjectFromBCS(bcs bcsDroppableObject) DroppableObject {
	return DroppableObject{
		SomeId:      bcs.SomeId,
		SomeNumber:  bcs.SomeNumber,
		SomeAddress: fmt.Sprintf("0x%x", bcs.SomeAddress),
		SomeAddresses: func() []string {
			addrs := make([]string, len(bcs.SomeAddresses))
			for i, addr := range bcs.SomeAddresses {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}
}

func init() {
	bind.RegisterStructDecoder("test::complex::SampleObject", func(data []byte) (interface{}, error) {
		var temp bcsSampleObject
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertSampleObjectFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("test::complex::DroppableObject", func(data []byte) (interface{}, error) {
		var temp bcsDroppableObject
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertDroppableObjectFromBCS(temp)
		return result, nil
	})
}

// NewObjectWithTransfer executes the new_object_with_transfer Move function.
func (c *ComplexContract) NewObjectWithTransfer(ctx context.Context, opts *bind.CallOpts, someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.NewObjectWithTransfer(someId, someNumber, someAddress, someAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// NewObject executes the new_object Move function.
func (c *ComplexContract) NewObject(ctx context.Context, opts *bind.CallOpts, someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.NewObject(someId, someNumber, someAddress, someAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// FlattenAddress executes the flatten_address Move function.
func (c *ComplexContract) FlattenAddress(ctx context.Context, opts *bind.CallOpts, someAddress string, someAddresses []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.FlattenAddress(someAddress, someAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// FlattenU8 executes the flatten_u8 Move function.
func (c *ComplexContract) FlattenU8(ctx context.Context, opts *bind.CallOpts, input [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.FlattenU8(input)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CheckU128 executes the check_u128 Move function.
func (c *ComplexContract) CheckU128(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.CheckU128(input)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CheckU256 executes the check_u256 Move function.
func (c *ComplexContract) CheckU256(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.CheckU256(input)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CheckWithObjectRef executes the check_with_object_ref Move function.
func (c *ComplexContract) CheckWithObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.CheckWithObjectRef(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CheckWithMutObjectRef executes the check_with_mut_object_ref Move function.
func (c *ComplexContract) CheckWithMutObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object, newNumber uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.complexEncoder.CheckWithMutObjectRef(obj, newNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// NewObject executes the new_object Move function using DevInspect to get return values.
//
// Returns: DroppableObject
func (d *ComplexDevInspect) NewObject(ctx context.Context, opts *bind.CallOpts, someId []byte, someNumber uint64, someAddress string, someAddresses []string) (DroppableObject, error) {
	encoded, err := d.contract.complexEncoder.NewObject(someId, someNumber, someAddress, someAddresses)
	if err != nil {
		return DroppableObject{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return DroppableObject{}, err
	}
	if len(results) == 0 {
		return DroppableObject{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(DroppableObject)
	if !ok {
		return DroppableObject{}, fmt.Errorf("unexpected return type: expected DroppableObject, got %T", results[0])
	}
	return result, nil
}

// FlattenAddress executes the flatten_address Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *ComplexDevInspect) FlattenAddress(ctx context.Context, opts *bind.CallOpts, someAddress string, someAddresses []string) ([]string, error) {
	encoded, err := d.contract.complexEncoder.FlattenAddress(someAddress, someAddresses)
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
	result, ok := results[0].([]string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []string, got %T", results[0])
	}
	return result, nil
}

// FlattenU8 executes the flatten_u8 Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *ComplexDevInspect) FlattenU8(ctx context.Context, opts *bind.CallOpts, input [][]byte) ([]byte, error) {
	encoded, err := d.contract.complexEncoder.FlattenU8(input)
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

// CheckU128 executes the check_u128 Move function using DevInspect to get return values.
//
// Returns: u128
func (d *ComplexDevInspect) CheckU128(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*big.Int, error) {
	encoded, err := d.contract.complexEncoder.CheckU128(input)
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
	result, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *big.Int, got %T", results[0])
	}
	return result, nil
}

// CheckU256 executes the check_u256 Move function using DevInspect to get return values.
//
// Returns: u256
func (d *ComplexDevInspect) CheckU256(ctx context.Context, opts *bind.CallOpts, input *big.Int) (*big.Int, error) {
	encoded, err := d.contract.complexEncoder.CheckU256(input)
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
	result, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *big.Int, got %T", results[0])
	}
	return result, nil
}

// CheckWithObjectRef executes the check_with_object_ref Move function using DevInspect to get return values.
//
// Returns: u64
func (d *ComplexDevInspect) CheckWithObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object) (uint64, error) {
	encoded, err := d.contract.complexEncoder.CheckWithObjectRef(obj)
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

// CheckWithMutObjectRef executes the check_with_mut_object_ref Move function using DevInspect to get return values.
//
// Returns: u64
func (d *ComplexDevInspect) CheckWithMutObjectRef(ctx context.Context, opts *bind.CallOpts, obj bind.Object, newNumber uint64) (uint64, error) {
	encoded, err := d.contract.complexEncoder.CheckWithMutObjectRef(obj, newNumber)
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

type complexEncoder struct {
	*bind.BoundContract
}

// NewObjectWithTransfer encodes a call to the new_object_with_transfer Move function.
func (c complexEncoder) NewObjectWithTransfer(someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_object_with_transfer", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"u64",
		"address",
		"vector<address>",
	}, []any{
		someId,
		someNumber,
		someAddress,
		someAddresses,
	}, nil)
}

// NewObjectWithTransferWithArgs encodes a call to the new_object_with_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) NewObjectWithTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"u64",
		"address",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_object_with_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// NewObject encodes a call to the new_object Move function.
func (c complexEncoder) NewObject(someId []byte, someNumber uint64, someAddress string, someAddresses []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_object", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"u64",
		"address",
		"vector<address>",
	}, []any{
		someId,
		someNumber,
		someAddress,
		someAddresses,
	}, []string{
		"test::complex::DroppableObject",
	})
}

// NewObjectWithArgs encodes a call to the new_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) NewObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"u64",
		"address",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_object", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::complex::DroppableObject",
	})
}

// FlattenAddress encodes a call to the flatten_address Move function.
func (c complexEncoder) FlattenAddress(someAddress string, someAddresses []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("flatten_address", typeArgsList, typeParamsList, []string{
		"address",
		"vector<address>",
	}, []any{
		someAddress,
		someAddresses,
	}, []string{
		"vector<address>",
	})
}

// FlattenAddressWithArgs encodes a call to the flatten_address Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) FlattenAddressWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"address",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("flatten_address", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// FlattenU8 encodes a call to the flatten_u8 Move function.
func (c complexEncoder) FlattenU8(input [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("flatten_u8", typeArgsList, typeParamsList, []string{
		"vector<vector<u8>>",
	}, []any{
		input,
	}, []string{
		"vector<u8>",
	})
}

// FlattenU8WithArgs encodes a call to the flatten_u8 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) FlattenU8WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("flatten_u8", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// CheckU128 encodes a call to the check_u128 Move function.
func (c complexEncoder) CheckU128(input *big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_u128", typeArgsList, typeParamsList, []string{
		"u128",
	}, []any{
		input,
	}, []string{
		"u128",
	})
}

// CheckU128WithArgs encodes a call to the check_u128 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) CheckU128WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u128",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_u128", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u128",
	})
}

// CheckU256 encodes a call to the check_u256 Move function.
func (c complexEncoder) CheckU256(input *big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_u256", typeArgsList, typeParamsList, []string{
		"u256",
	}, []any{
		input,
	}, []string{
		"u256",
	})
}

// CheckU256WithArgs encodes a call to the check_u256 Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) CheckU256WithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u256",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_u256", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u256",
	})
}

// CheckWithObjectRef encodes a call to the check_with_object_ref Move function.
func (c complexEncoder) CheckWithObjectRef(obj bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_with_object_ref", typeArgsList, typeParamsList, []string{
		"&SampleObject",
	}, []any{
		obj,
	}, []string{
		"u64",
	})
}

// CheckWithObjectRefWithArgs encodes a call to the check_with_object_ref Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) CheckWithObjectRefWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&SampleObject",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_with_object_ref", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// CheckWithMutObjectRef encodes a call to the check_with_mut_object_ref Move function.
func (c complexEncoder) CheckWithMutObjectRef(obj bind.Object, newNumber uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_with_mut_object_ref", typeArgsList, typeParamsList, []string{
		"&mut SampleObject",
		"u64",
	}, []any{
		obj,
		newNumber,
	}, []string{
		"u64",
	})
}

// CheckWithMutObjectRefWithArgs encodes a call to the check_with_mut_object_ref Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c complexEncoder) CheckWithMutObjectRefWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut SampleObject",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("check_with_mut_object_ref", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}
