// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package modulecomplex

import (
	"context"
	"fmt"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// Built with `sui move build --dump-bytecode-as-base64 --with-unpublished-dependencies`
const ComplexJSON = `
{
  "modules": [
    "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"
  ],
  "dependencies": [
    "0x0000000000000000000000000000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000000000000000000000000000002"
  ]
}
`

func PublishComplex(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*ComplexContract, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.ToArtifact(ComplexJSON)
	if err != nil {
		return nil, nil, err
	}

	packageId, tx, err := bind.PublishPackage(ctx, opts, signer, client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return nil, nil, err
	}
	contract, err := NewComplex(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}

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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "complex", "flatten_u8", false, "", input)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "complex", "flatten_u8", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
