package offramp

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

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

func NewOfframp(address string, chainClient client.BindingsClient) (Offramp, error) {
	offrampContract, err := module_offramp.NewOfframp(address, chainClient)
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

func PublishOfframp(ctx context.Context, opts *bind.CallOpts, chainClient client.BindingsClient, ccipAddress string, mcmsAddress, suiRPC string) (Offramp, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.CCIPOfframp, map[string]string{
		"mcms":                      mcmsAddress,
		"ccip":                      ccipAddress,
		"ccip_offramp":              "0x0",
		"mcms_owner":                "0x1",
		"mcms_register_entrypoints": "0x2",
		"signer":                    signerAddr,
	}, false, suiRPC)
	if err != nil {
		return nil, nil, err
	}

	//nolint:revive // var-naming: generated bindings keep packageId naming
	packageId, tx, err := bind.PublishPackage(ctx, opts, chainClient, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return nil, nil, err
	}

	contract, err := NewOfframp(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
