package onramp

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type Onramp interface {
	Address() sui.Address
	Onramp() module_onramp.IOnramp
}

var _ Onramp = OnrampPackage{}

type OnrampPackage struct {
	address sui.Address

	onramp module_onramp.IOnramp
}

func (p OnrampPackage) Address() sui.Address {
	return p.address
}

func (p OnrampPackage) Onramp() module_onramp.IOnramp {
	return p.onramp
}

func NewOnramp(address string, client suiclient.ClientImpl) (Onramp, error) {
	onrampContract, err := module_onramp.NewOnramp(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return OnrampPackage{
		address: *packageId,
		onramp:  onrampContract,
	}, nil
}

func PublishOnramp(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl, ccipAddress string, mcmsAddress string) (Onramp, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPOnramp, map[string]string{
		"mcms":   mcmsAddress,
		"ccip":   ccipAddress,
		"onramp": "0x0",
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

	contract, err := NewOnramp(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
