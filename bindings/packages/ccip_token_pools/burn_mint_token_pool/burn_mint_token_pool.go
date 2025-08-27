package burnminttokenpool

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type BurnMintTokenPool interface {
	Address() string
}

var _ BurnMintTokenPool = CCIPBurnMintTokenPoolPackage{}

type CCIPBurnMintTokenPoolPackage struct {
	address string

	tokenPool module_burn_mint_token_pool.IBurnMintTokenPool
}

func (p CCIPBurnMintTokenPoolPackage) Address() string {
	return p.address
}

func NewCCIPBurnMintTokenPool(address string, client sui.ISuiAPI) (BurnMintTokenPool, error) {
	tokenPoolContract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPBurnMintTokenPoolPackage{
		address:   packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPBurnMintTokenPool(
	ctx context.Context,
	opts *bind.CallOpts,
	client sui.ISuiAPI,
	ccipAddress,
	ccipTokenPoolAddress,
	mcmsAddress,
	mcmsOwnerAddress string) (BurnMintTokenPool, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.BurnMintTokenPool, map[string]string{
		"ccip":                 ccipAddress,
		"ccip_token_pool":      ccipTokenPoolAddress,
		"burn_mint_token_pool": "0x0",
		"mcms":                 mcmsAddress,
		"mcms_owner":           mcmsOwnerAddress,
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

	contract, err := NewCCIPBurnMintTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
