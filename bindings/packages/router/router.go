package router

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

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

func NewCCIPRouter(address string, client sui.ISuiAPI) (CCIPRouter, error) {
	routerContract, err := module_router.NewRouter(address, client)
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

func PublishCCIPRouter(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI, mcmsAddress string, mcmsOwner string) (CCIPRouter, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPRouter, map[string]string{
		"ccip_router": "0x0",
		"mcms":        mcmsAddress,
		"mcms_owner":  mcmsOwner,
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

	contract, err := NewCCIPRouter(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
