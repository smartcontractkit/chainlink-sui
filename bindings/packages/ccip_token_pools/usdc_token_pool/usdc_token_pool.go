package usdctokenpool

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_usdc_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/usdc_token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type USDCTokenPool interface {
	Address() string
}

type CCIPUSDCTokenPoolPackage struct {
	address string

	tokenPool module_usdc_token_pool.IUsdcTokenPool
}

func (p CCIPUSDCTokenPoolPackage) Address() string {
	return p.address
}

func NewCCIPUSDCTokenPool(address string, client sui.ISuiAPI) (USDCTokenPool, error) {
	tokenPoolContract, err := module_usdc_token_pool.NewUsdcTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPUSDCTokenPoolPackage{
		address:   packageId,
		tokenPool: tokenPoolContract,
	}, nil
}

func PublishCCIPUSDCTokenPool(
	ctx context.Context,
	opts *bind.CallOpts,
	client sui.ISuiAPI,
	ccipAddress,
	ccipTokenPoolAddress,
	usdcLocalTokenAddress,
	messageTransmitterAddress,
	tokenMessengerMinterAddress,
	mcmsAddress,
	mcmsOwnerAddress string) (USDCTokenPool, *models.SuiTransactionBlockResponse, error) {

	artifact, err := bind.CompilePackage(contracts.USDCTokenPool, map[string]string{
		"ccip":                   ccipAddress,
		"ccip_token_pool":        ccipTokenPoolAddress,
		"usdc_token_pool":        "0x0",
		"usdc_local_token":       usdcLocalTokenAddress,
		"token_messenger_minter": tokenMessengerMinterAddress,
		"message_transmitter":    messageTransmitterAddress,
		"mcms":                   mcmsAddress,
		"mcms_owner":             mcmsOwnerAddress,
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

	contract, err := NewCCIPUSDCTokenPool(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
