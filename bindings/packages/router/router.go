package router

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type CCIPRouter interface {
	Address() sui.Address
}

var _ CCIPRouter = CCIPRouterPackage{}

type CCIPRouterPackage struct {
	address sui.Address

	router module_router.IRouter
}

func (p CCIPRouterPackage) Address() sui.Address {
	return p.address
}

func NewCCIPRouter(address string, client suiclient.ClientImpl) (CCIPRouter, error) {
	routerContract, err := module_router.NewRouter(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPRouterPackage{
		address: *packageId,
		router:  routerContract,
	}, nil
}

func PublishCCIPRouter(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl, mcmsAddress string, mcmsOwner string) (CCIPRouter, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPRouter, map[string]string{
		"ccip_router": "0x0",
		"mcms":        mcmsAddress,
		"mcms_owner":  mcmsOwner,
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

	contract, err := NewCCIPRouter(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
