package counter

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// This should be auto-generated when compiling, same as the bindings
// Currently built with `sui move build --dump-bytecode-as-base64 --with-unpublished-dependencies`
const CounterJSON = `{
  "modules": [
    "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"
  ],
  "dependencies": [
    "0x0000000000000000000000000000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000000000000000000000000000002"
  ],
  "digest": [
    205, 51, 124, 253, 124, 136, 134, 192, 240, 141, 61, 232, 201, 214, 35, 229,
    116, 23, 208, 55, 0, 255, 80, 13, 43, 253, 114, 188, 176, 168, 115, 210
  ]
}`

func PublishCounter(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (*Counter, *suiclient.SuiTransactionBlockResponse, error) {
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
	counter, err := NewCounter(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return counter, tx, nil
}

type ICounter interface {
	// We require ctx even for building
	Increment(objectId string) bind.IMethod
	GetCount(objectId string) bind.IMethod
	// TODO: Add rest of methods

	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
	// Gets the object information. This case, only the object value
	Inspect(ctx context.Context, objectId string) (uint64, error)
}

type Counter struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ ICounter = (*Counter)(nil)

func NewCounter(packageID string, client suiclient.ClientImpl) (*Counter, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &Counter{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *Counter) Connect(client suiclient.ClientImpl) {
	c.client = client
}

func (c *Counter) Increment(counterObjectId string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment", counterObjectId)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *Counter) IncrementMult(counterObjectId string, a, b uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "increment_mult", counterObjectId, a, b)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "increment_mult", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *Counter) GetCount(counterObjectId string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "counter", "get_count", counterObjectId)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "counter", "get_count", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *Counter) Inspect(ctx context.Context, counterObjectId string) (uint64, error) {
	obj, err := bind.ReadObject(ctx, counterObjectId, c.client)
	if err != nil {
		return 0, err
	}
	var count string
	err = bind.GetCustomValueFromObjectData(*obj.Data, &count)
	if err != nil {
		return 0, err
	}
	// Convert count string to uint64
	countUint, err := strconv.ParseUint(count, 10, 64)
	if err != nil {
		return 0, err
	}

	return countUint, nil
}
