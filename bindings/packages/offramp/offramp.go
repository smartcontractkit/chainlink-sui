package offramp

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type Offramp interface {
	Address() string
	Offramp() module_offramp.IOfframp
}

var _ Offramp = OfframpPackage{}

type OfframpPackage struct {
	address string

	offramp module_offramp.IOfframp
}

func (p OfframpPackage) Address() string {
	return p.address
}

func (p OfframpPackage) Offramp() module_offramp.IOfframp {
	return p.offramp
}

func NewOfframp(address string, client sui.ISuiAPI) (Offramp, error) {
	offrampContract, err := module_offramp.NewOfframp(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return OfframpPackage{
		address: packageId,
		offramp: offrampContract,
	}, nil
}

func PublishOfframp(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI, ccipAddress string, mcmsAddress string) (Offramp, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPOfframp, map[string]string{
		"mcms":                      mcmsAddress,
		"ccip":                      ccipAddress,
		"ccip_offramp":              "0x0",
		"mcms_owner":                "0x1",
		"mcms_register_entrypoints": "0x2",
	})
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

	contract, err := NewOfframp(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
