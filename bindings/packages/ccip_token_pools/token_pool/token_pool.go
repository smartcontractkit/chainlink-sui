package tokenpool

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type TokenPool interface {
	Address() string
}

var _ TokenPool = CCIPTokenPoolPackage{}

type CCIPTokenPoolPackage struct {
	address string

	tokenPool module_token_pool.ITokenPool
}

func (p CCIPTokenPoolPackage) Address() string {
	return p.address
}

func NewCCIPTokenPool(address string, client sui.ISuiAPI) (TokenPool, error) {
	tokenPoolContract, err := module_token_pool.NewTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPTokenPoolPackage{
		address:   packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPTokenPool(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI, ccipAddress, mcmsAddress, mcmsOwner string) (TokenPool, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPTokenPool, map[string]string{
		"ccip":            ccipAddress,
		"ccip_token_pool": "0x0",
		"mcms":            mcmsAddress,
		"mcms_owner":      mcmsOwner,
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

	contract, err := NewCCIPTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
