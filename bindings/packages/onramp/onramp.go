package onramp

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type Onramp interface {
	Address() string
	Onramp() module_onramp.IOnramp
}

var _ Onramp = OnrampPackage{}

type OnrampPackage struct {
	address string

	onramp module_onramp.IOnramp
}

func (p OnrampPackage) Address() string {
	return p.address
}

func (p OnrampPackage) Onramp() module_onramp.IOnramp {
	return p.onramp
}

func NewOnramp(address string, client sui.ISuiAPI) (Onramp, error) {
	onrampContract, err := module_onramp.NewOnramp(address, client)
	if err != nil {
		return nil, err
	}

	return OnrampPackage{
		address: address,
		onramp:  onrampContract,
	}, nil
}

func PublishOnramp(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI, ccipAddress, mcmsAddress, mcmsOwnerAddress string) (Onramp, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPOnramp, map[string]string{
		"ccip":        ccipAddress,
		"ccip_onramp": "0x0",
		"mcms":        mcmsAddress,
		"mcms_owner":  mcmsOwnerAddress,
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

	contract, err := NewOnramp(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
