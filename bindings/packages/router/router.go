package router

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type CCIPRouter interface {
	Address() string
}

var _ CCIPRouter = CCIPRouterPackage{}

type CCIPRouterPackage struct {
	address string

	router module_router.IRouter
}

func (p CCIPRouterPackage) Address() string {
	return p.address
}

func NewCCIPRouter(address string, chainClient client.BindingsClient) (CCIPRouter, error) {
	routerContract, err := module_router.NewRouter(address, chainClient)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPRouterPackage{
		address: packageId,
		router:  routerContract,
	}, nil
}

func PublishCCIPRouter(ctx context.Context, opts *bind.CallOpts, chainClient client.BindingsClient, mcmsAddress string, mcmsOwner, suiRPC string) (CCIPRouter, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.CCIPRouter, map[string]string{
		"ccip_router":               "0x0",
		"mcms":                      mcmsAddress,
		"mcms_owner":                mcmsOwner,
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

	contract, err := NewCCIPRouter(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
