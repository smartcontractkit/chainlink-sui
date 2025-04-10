// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package modulecounter

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
const CounterJSON = `
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

func PublishCounter(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*CounterContract, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.ToArtifact(CounterJSON)
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
	contract, err := NewCounter(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}

type ICounter interface {
	Initialize() bind.IMethod
	Increment(counter string) bind.IMethod
	IncrementByOne(counter string) bind.IMethod
	IncrementByOneNoContext(counter string) bind.IMethod
	IncrementByTwo(admin string, counter string) bind.IMethod
	IncrementByTwoNoContext(admin string, counter string) bind.IMethod
	IncrementMult(counter string, a uint64, b uint64) bind.IMethod
	GetCount(counter string) bind.IMethod
	GetCountNoEntry(counter string) bind.IMethod
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

type AdminCap struct {
	Id string `move:"sui::object::UID"`
}

type Counter struct {
	Id    string `move:"sui::object::UID"`
	Value uint64 `move:"u64"`
}

// Functions

func (c *CounterContract) Initialize() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "initialize", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) Increment(counter string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment", false, "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByOne(counter string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_one", false, "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_one", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByOneNoContext(counter string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_one_no_context", false, "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_one_no_context", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByTwo(admin string, counter string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_two", false, "", admin, counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_two", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementByTwoNoContext(admin string, counter string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_by_two_no_context", false, "", admin, counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_by_two_no_context", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) IncrementMult(counter string, a uint64, b uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_mult", false, "", counter, a, b)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_mult", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetCount(counter string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_count", false, "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_count", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *CounterContract) GetCountNoEntry(counter string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_count_no_entry", false, "", counter)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_count_no_entry", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
