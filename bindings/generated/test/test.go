package test

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	modulecomplex "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/complex"
	modulecounter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type Test interface {
	Address() sui.Address
	Counter() modulecounter.ICounter
	Complex() modulecomplex.IComplex
}

var _ Test = TestPackage{}

type TestPackage struct {
	address sui.Address

	counter modulecounter.ICounter
	complex modulecomplex.IComplex
}

func (p TestPackage) Address() sui.Address {
	return p.address
}

func (p TestPackage) Counter() modulecounter.ICounter {
	return p.counter
}

func (p TestPackage) Complex() modulecomplex.IComplex {
	return p.complex
}

func NewTest(address string, client suiclient.ClientImpl) (Test, error) {
	counterContract, err := modulecounter.NewCounter(address, client)
	if err != nil {
		return nil, err
	}

	complexContract, err := modulecomplex.NewComplex(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return TestPackage{
		address: *packageId,
		counter: counterContract,
		complex: complexContract,
	}, nil
}

func PublishTest(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (Test, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.Test, map[string]string{
		"test": "0x0",
	})
	if err != nil {
		return nil, nil, err
	}

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
	contract, err := NewTest(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
