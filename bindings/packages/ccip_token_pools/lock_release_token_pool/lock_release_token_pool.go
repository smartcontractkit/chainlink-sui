package lockreleasetokenpool

import (
	"context"
	"fmt"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

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

func NewCCIPLockReleaseTokenPool(address string, chainClient client.BindingsClient) (LockReleaseTokenPool, error) {
	tokenPoolContract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(address, chainClient)
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
	chainClient client.BindingsClient,
	ccipAddress string,
	mcmsAddress,
	fastMcmsAddress,
	mcmsOwnerAddress, suiRPC string) (LockReleaseTokenPool, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.LockReleaseTokenPool, map[string]string{
		"ccip":                    ccipAddress,
		"lock_release_token_pool": "0x0",
		"mcms":                    mcmsAddress,
		"fast_mcms":               fastMcmsAddress,
		"mcms_owner":              mcmsOwnerAddress,
		"signer":                  signerAddr,
	}, false, suiRPC)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to compile package: %w", err)
	}

	//nolint:revive // var-naming: generated bindings keep packageId naming
	packageId, tx, err := bind.PublishPackage(ctx, opts, chainClient, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to publish package: %w", err)
	}

	contract, err := NewCCIPLockReleaseTokenPool(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
