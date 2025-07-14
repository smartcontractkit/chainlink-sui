package lockreleasetokenpool

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type LockReleaseTokenPool interface {
	Address() string
}

var _ LockReleaseTokenPool = CCIPLockReleaseTokenPool{}

type CCIPLockReleaseTokenPool struct {
	address string

	tokenPool module_lock_release_token_pool.ILockReleaseTokenPool
}

func (p CCIPLockReleaseTokenPool) Address() string {
	return p.address
}

func NewCCIPLockReleaseTokenPool(address string, client sui.ISuiAPI) (LockReleaseTokenPool, error) {
	tokenPoolContract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPLockReleaseTokenPool{
		address:   packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPLockReleaseTokenPool(
	ctx context.Context,
	opts *bind.CallOpts,
	client sui.ISuiAPI,
	ccipAddress string,
	ccipTokenPoolAddress string,
	mcmsAddress,
	mcmsOwnerAddress string) (LockReleaseTokenPool, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.LockReleaseTokenPool, map[string]string{
		"ccip":                    ccipAddress,
		"ccip_token_pool":         ccipTokenPoolAddress,
		"lock_release_token_pool": "0x0",
		"mcms":                    mcmsAddress,
		"mcms_owner":              mcmsOwnerAddress,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compile package: %w", err)
	}

	packageId, tx, err := bind.PublishPackage(ctx, opts, client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to publish package: %w", err)
	}

	contract, err := NewCCIPLockReleaseTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
