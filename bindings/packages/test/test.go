package test

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	modulecomplex "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/complex"
	modulecounter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	module_generics "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/generics"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type Test interface {
	Address() string
	Counter() modulecounter.ICounter
	Complex() modulecomplex.IComplex
	Generics() module_generics.IGenerics
}

var _ Test = TestPackage{}

type TestPackage struct {
	address string

	counter   modulecounter.ICounter
	complex   modulecomplex.IComplex
	generics  module_generics.IGenerics
	PackageId string
}

func (p TestPackage) Address() string {
	return p.address
}

func (p TestPackage) Counter() modulecounter.ICounter {
	return p.counter
}

func (p TestPackage) Complex() modulecomplex.IComplex {
	return p.complex
}

func (p TestPackage) Generics() module_generics.IGenerics {
	return p.generics
}

func NewTest(address string, client sui.ISuiAPI) (Test, error) {
	counterContract, err := modulecounter.NewCounter(address, client)
	if err != nil {
		return nil, err
	}

	complexContract, err := modulecomplex.NewComplex(address, client)
	if err != nil {
		return nil, err
	}

	genericsContract, err := module_generics.NewGenerics(address, client)
	if err != nil {
		return nil, err
	}

	return TestPackage{
		address:   address,
		counter:   counterContract,
		complex:   complexContract,
		generics:  genericsContract,
		PackageId: address,
	}, nil
}

func PublishTest(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI) (Test, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.Test, map[string]string{
		"test": "0x0",
	})
	if err != nil {
		return nil, nil, err
	}

	if err != nil {
		return nil, nil, err
	}

	packageId, tx, err := bind.PublishPackage(ctx, opts, client, bind.PublishRequest{
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
