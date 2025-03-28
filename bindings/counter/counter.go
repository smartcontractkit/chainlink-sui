package counter

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
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

func PublishCounter(ctx context.Context, opts bind.TxOpts, signer signer.Signer, client sui.ISuiAPI) (*Counter, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.ToArtifact(CounterJSON)
	if err != nil {
		return nil, nil, err
	}

	packageid, tx, err := bind.PublishPackage(ctx, opts, signer, client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return nil, nil, err
	}

	return NewCounter(packageid), tx, nil
}

type ICounter interface {
	Increment(objectId string) bind.IMethod
	// TODO: Add rest of methods
}

type Counter struct {
	packageID bind.PackageID
}

var _ ICounter = (*Counter)(nil)

func NewCounter(packageID string) *Counter {
	return &Counter{
		packageID: packageID,
	}
}

func (c *Counter) EncodeIncrement(counterObjectId string) (encodedArgs []any, err error) {
	return bind.Encode(
		[]string{
			"address",
		},
		[]any{counterObjectId},
	)
}

func (c *Counter) EncodeIncrementMult(counterObjectId string, a, b uint64) (encodedArgs []any, err error) {
	return bind.Encode(
		[]string{
			"address",
			"u64",
			"u64",
		},
		[]any{counterObjectId, a, b},
	)
}

func (c *Counter) Increment(counterObjectId string) bind.IMethod {
	build := func(opts bind.TxOpts, signer string) (*suiptb.ProgrammableTransactionBuilder, error) {
		payload, err := c.EncodeIncrement(counterObjectId)
		if err != nil {
			return nil, err
		}
		ptb, err := bind.BuildCallTransaction(opts, c.packageID, "counter", "increment", payload)
		if err != nil {
			return nil, err
		}
		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build))
}
