// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_complex

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
)

type IComplex interface {
	NewObjectWithTransfer(someId []byte, someNumber uint64, someAddress string, someAddresses []string) bind.IMethod
	NewObject(someId []byte, someNumber uint64, someAddress string, someAddresses []string) bind.IMethod
	FlattenAddress(someAddress string, someAddresses []string) bind.IMethod
	FlattenU8(input [][]byte) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type ComplexContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IComplex = (*ComplexContract)(nil)

func NewComplex(packageID string, client suiclient.ClientImpl) (*ComplexContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &ComplexContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *ComplexContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

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

// Functions

func (c *ComplexContract) NewObjectWithTransfer(someId []byte, someNumber uint64, someAddress string, someAddresses []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "complex", "new_object_with_transfer", false, "", someId, someNumber, someAddress, someAddresses)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "complex", "new_object_with_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ComplexContract) NewObject(someId []byte, someNumber uint64, someAddress string, someAddresses []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "complex", "new_object", false, "", someId, someNumber, someAddress, someAddresses)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "complex", "new_object", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ComplexContract) FlattenAddress(someAddress string, someAddresses []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "complex", "flatten_address", false, "", someAddress, someAddresses)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "complex", "flatten_address", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ComplexContract) FlattenU8(input [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "complex", "flatten_u8", false, "", input)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "complex", "flatten_u8", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
