package burnminttokenpool

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

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

func NewCCIPBurnMintTokenPool(address string, chainClient client.BindingsClient) (BurnMintTokenPool, error) {
	tokenPoolContract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(address, chainClient)
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
	chainClient client.BindingsClient,
	ccipAddress,
	mcmsAddress,
	fastMcmsAddress,
	mcmsOwnerAddress, suiRPC string) (BurnMintTokenPool, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.BurnMintTokenPool, map[string]string{
		"ccip":                 ccipAddress,
		"burn_mint_token_pool": "0x0",
		"mcms":                 mcmsAddress,
		"fast_mcms":            fastMcmsAddress,
		"mcms_owner":           mcmsOwnerAddress,
		"signer":               signerAddr,
	}, false, suiRPC)
	if err != nil {
		return nil, nil, err
	}

	//nolint:revive // var-naming: generated bindings keep packageId naming
	packageId, tx, err := bind.PublishPackage(ctx, opts, chainClient, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return nil, nil, err
	}

	contract, err := NewCCIPBurnMintTokenPool(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
