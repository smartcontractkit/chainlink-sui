package tokenpool

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type TokenPool interface {
	Address() sui.Address
}

var _ TokenPool = CCIPTokenPoolPackage{}

type CCIPTokenPoolPackage struct {
	address sui.Address

	tokenPool module_token_pool.ITokenPool
}

func (p CCIPTokenPoolPackage) Address() sui.Address {
	return p.address
}

func NewCCIPTokenPool(address string, client suiclient.ClientImpl) (TokenPool, error) {
	tokenPoolContract, err := module_token_pool.NewTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPTokenPoolPackage{
		address:   *packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPTokenPool(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl, ccipRouterAddress, ccipAddress, mcmsAddress, mcmsOwner string) (TokenPool, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPTokenPools, map[string]string{
		"ccip":            ccipAddress,
		"ccip_router":     ccipRouterAddress,
		"ccip_token_pool": "0x0",
		"mcms":            mcmsAddress,
		"mcms_owner":      mcmsOwner,
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

	contract, err := NewCCIPTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
