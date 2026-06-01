package onramp

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

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

func NewOnramp(address string, chainClient client.BindingsClient) (Onramp, error) {
	onrampContract, err := module_onramp.NewOnramp(address, chainClient)
	if err != nil {
		return nil, err
	}

	return OnrampPackage{
		address: address,
		onramp:  onrampContract,
	}, nil
}

func PublishOnramp(ctx context.Context, opts *bind.CallOpts, chainClient client.BindingsClient, ccipAddress, mcmsAddress, mcmsOwnerAddress, suiRPC string) (Onramp, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.CCIPOnramp, map[string]string{
		"ccip":        ccipAddress,
		"ccip_onramp": "0x0",
		"mcms":        mcmsAddress,
		"mcms_owner":  mcmsOwnerAddress,
		"signer":      signerAddr,
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

	contract, err := NewOnramp(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
