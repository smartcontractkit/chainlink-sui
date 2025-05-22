package offramp

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type Offramp interface {
	Address() sui.Address
	Offramp() module_offramp.IOfframp
}

var _ Offramp = OfframpPackage{}

type OfframpPackage struct {
	address sui.Address

	offramp module_offramp.IOfframp
}

func (p OfframpPackage) Address() sui.Address {
	return p.address
}

func (p OfframpPackage) Offramp() module_offramp.IOfframp {
	return p.offramp
}

func NewOfframp(address string, client suiclient.ClientImpl) (Offramp, error) {
	offrampContract, err := module_offramp.NewOfframp(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return OfframpPackage{
		address: *packageId,
		offramp: offrampContract,
	}, nil
}

func PublishOfframp(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl, ccipAddress string, mcmsAddress string) (Offramp, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPOfframp, map[string]string{
		"mcms":    mcmsAddress,
		"ccip":    ccipAddress,
		"offramp": "0x0",
	})
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

	contract, err := NewOfframp(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
