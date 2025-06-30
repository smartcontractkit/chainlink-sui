package burnminttokenpool

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type BurnMintTokenPool interface {
	Address() sui.Address
}

var _ BurnMintTokenPool = CCIPBurnMintTokenPoolPackage{}

type CCIPBurnMintTokenPoolPackage struct {
	address sui.Address

	tokenPool module_burn_mint_token_pool.IBurnMintTokenPool
}

func (p CCIPBurnMintTokenPoolPackage) Address() sui.Address {
	return p.address
}

func NewCCIPBurnMintTokenPool(address string, client suiclient.ClientImpl) (BurnMintTokenPool, error) {
	tokenPoolContract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPBurnMintTokenPoolPackage{
		address:   *packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPBurnMintTokenPool(
	ctx context.Context,
	opts bind.TxOpts,
	signer rel.SuiSigner,
	client suiclient.ClientImpl,
	ccipAddress string,
	ccipTokenPoolAddress string) (BurnMintTokenPool, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIPTokenPools, map[string]string{
		"ccip":                  ccipAddress,
		"ccip_token_pool":       ccipTokenPoolAddress,
		"burn_mint_token_pool":  "0x0",
		"burn_mint_local_token": "0x1",
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

	contract, err := NewCCIPBurnMintTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
