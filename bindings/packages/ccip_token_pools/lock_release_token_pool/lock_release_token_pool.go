package lockreleasetokenpool

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type LockReleaseTokenPool interface {
	Address() sui.Address
}

var _ LockReleaseTokenPool = CCIPLockReleaseTokenPool{}

type CCIPLockReleaseTokenPool struct {
	address sui.Address

	tokenPool module_lock_release_token_pool.ILockReleaseTokenPool
}

func (p CCIPLockReleaseTokenPool) Address() sui.Address {
	return p.address
}

func NewCCIPLockReleaseTokenPool(address string, client suiclient.ClientImpl) (LockReleaseTokenPool, error) {
	tokenPoolContract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPLockReleaseTokenPool{
		address:   *packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPLockReleaseTokenPool(
	ctx context.Context,
	opts bind.TxOpts,
	signer rel.SuiSigner,
	client suiclient.ClientImpl,
	ccipAddress,
	ccipTokenPoolAddress,
	mcmsAddress,
	mcmsOwnerAddres string) (LockReleaseTokenPool, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.LockReleaseTokenPool, map[string]string{
		"ccip":                    ccipAddress,
		"ccip_token_pool":         ccipTokenPoolAddress,
		"lock_release_token_pool": "0x0",
		"mcms":                    mcmsAddress,
		"mcms_owner":              mcmsOwnerAddres,
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

	contract, err := NewCCIPLockReleaseTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
