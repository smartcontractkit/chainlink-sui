package managedtokenpool

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type ManagedTokenPool interface {
	Address() string
}

var _ ManagedTokenPool = CCIPManagedTokenPoolPackage{}

type CCIPManagedTokenPoolPackage struct {
	address string

	tokenPool module_managed_token_pool.IManagedTokenPool
}

func (p CCIPManagedTokenPoolPackage) Address() string {
	return p.address
}

func NewCCIPManagedTokenPool(address string, client sui.ISuiAPI) (ManagedTokenPool, error) {
	tokenPoolContract, err := module_managed_token_pool.NewManagedTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPManagedTokenPoolPackage{
		address:   packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPManagedTokenPool(
	ctx context.Context,
	opts *bind.CallOpts,
	client sui.ISuiAPI,
	ccipAddress,
	ccipTokenPoolAddress,
	managedTokenAddress,
	mcmsAddress,
	mcmsOwnerAddress string) (ManagedTokenPool, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.ManagedTokenPool, map[string]string{
		"ccip":               ccipAddress,
		"ccip_token_pool":    ccipTokenPoolAddress,
		"managed_token_pool": "0x0",
		"managed_token":      managedTokenAddress,
		"mcms":               mcmsAddress,
		"mcms_owner":         mcmsOwnerAddress,
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

	contract, err := NewCCIPManagedTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
