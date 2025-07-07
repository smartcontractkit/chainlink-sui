// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_counter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
	_ = uint256.NewInt
)

type ICounter interface {
	Initialize() bind.IMethod
	Increment(counter bind.Object) bind.IMethod
	Create() bind.IMethod
	IncrementByOne(counter bind.Object) bind.IMethod
	IncrementByOneNoContext(counter bind.Object) bind.IMethod
	IncrementByTwo(admin bind.Object, counter bind.Object) bind.IMethod
	IncrementByTwoNoContext(admin bind.Object, counter bind.Object) bind.IMethod
	IncrementBy(counter bind.Object, by uint64) bind.IMethod
	IncrementMult(counter bind.Object, a uint64, b uint64) bind.IMethod
	GetCount(counter bind.Object) bind.IMethod
	GetCountUsingPointer(counter bind.Object) bind.IMethod
	GetCountNoEntry(counter bind.Object) bind.IMethod
	GetCoinValue(typeArgs string, coin bind.Object) bind.IMethod
	GetAddressList() bind.IMethod
	GetSimpleResult() bind.IMethod
	GetResultStruct() bind.IMethod
	GetNestedResultStruct() bind.IMethod
	GetMultiNestedResultStruct() bind.IMethod
	GetTupleStruct() bind.IMethod
	GetOcrConfig() bind.IMethod
	GetVectorOfU8() bind.IMethod
	GetVectorOfAddresses() bind.IMethod
	GetVectorOfVectorsOfU8() bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type CounterContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ ICounter = (*CounterContract)(nil)

func NewCounter(packageID string, client suiclient.ClientImpl) (*CounterContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &CounterContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *CounterContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type COUNTER struct {
}

type CounterIncremented struct {
	CounterId bind.Object `move:"ID"`
	NewValue  uint64      `move:"u64"`
}

type AdminCap struct {
	Id string `move:"sui::object::UID"`
}

type Counter struct {
	Id    string `move:"sui::object::UID"`
	Value uint64 `move:"u64"`
}

type CounterPointer struct {
	Id         string `move:"sui::object::UID"`
	CounterId  string `move:"address"`
	AdminCapId string `move:"address"`
}

type AddressList struct {
	Addresses []string `move:"vector<address>"`
	Count     uint64   `move:"u64"`
}

type SimpleResult struct {
	Value uint64 `move:"u64"`
}

type ComplexResult struct {
	Count     uint64 `move:"u64"`
	Addr      string `move:"address"`
	IsComplex bool   `move:"bool"`
	Bytes     []byte `move:"vector<u8>"`
}

type NestedStruct struct {
	IsNested           bool          `move:"bool"`
	DoubleCount        uint64        `move:"u64"`
	NestedStruct       ComplexResult `move:"ComplexResult"`
	NestedSimpleStruct SimpleResult  `move:"SimpleResult"`
}

type MultiNestedStruct struct {
	IsMultiNested      bool         `move:"bool"`
	DoubleCount        uint64       `move:"u64"`
	NestedStruct       NestedStruct `move:"NestedStruct"`
	NestedSimpleStruct SimpleResult `move:"SimpleResult"`
}

type ConfigInfo struct {
	ConfigDigest                   []byte `move:"vector<u8>"`
	BigF                           byte   `move:"u8"`
	N                              byte   `move:"u8"`
	IsSignatureVerificationEnabled bool   `move:"bool"`
}

type OCRConfig struct {
	ConfigInfo   ConfigInfo `move:"ConfigInfo"`
	Signers      [][]byte   `move:"vector<vector<u8>>"`
	Transmitters []string   `move:"vector<address>"`
}

// Functions

func (c *CounterContract) Initialize() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "initialize", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) Increment(counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment", false, "", "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) Create() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "create", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "create", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByOne(counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_one", false, "", "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_one", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByOneNoContext(counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_one_no_context", false, "", "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_one_no_context", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByTwo(admin bind.Object, counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_two", false, "", "", admin, counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_two", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByTwoNoContext(admin bind.Object, counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_two_no_context", false, "", "", admin, counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_two_no_context", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementBy(counter bind.Object, by uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by", false, "", "", counter, by)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementMult(counter bind.Object, a uint64, b uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_mult", false, "", "", counter, a, b)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_mult", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetCount(counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_count", false, "", "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_count", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetCountUsingPointer(counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_count_using_pointer", false, "", "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_count_using_pointer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetCountNoEntry(counter bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_count_no_entry", false, "", "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_count_no_entry", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetCoinValue(typeArgs string, coin bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_coin_value", false, "", typeArgs, coin)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_coin_value", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetAddressList() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_address_list", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_address_list", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetSimpleResult() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_simple_result", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_simple_result", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetResultStruct() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_result_struct", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_result_struct", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetNestedResultStruct() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_nested_result_struct", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_nested_result_struct", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetMultiNestedResultStruct() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_multi_nested_result_struct", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_multi_nested_result_struct", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetTupleStruct() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_tuple_struct", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_tuple_struct", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetOcrConfig() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_ocr_config", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_ocr_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetVectorOfU8() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_vector_of_u8", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_vector_of_u8", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetVectorOfAddresses() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_vector_of_addresses", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_vector_of_addresses", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetVectorOfVectorsOfU8() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_vector_of_vectors_of_u8", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_vector_of_vectors_of_u8", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
