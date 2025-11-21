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

var _ USDCTokenPool = CCIPUSDCTokenPoolPackage{}

type CCIPUSDCTokenPoolPackage struct {
	address string

	usdcTokenPool module_usdc_token_pool.IUsdcTokenPool
}

func (p CCIPUSDCTokenPoolPackage) Address() string {
	return p.address
}

func NewCCIPUSDCTokenPool(address string, client sui.ISuiAPI) (USDCTokenPool, error) {
	usdcTokenPoolContract, err := module_usdc_token_pool.NewUsdcTokenPool(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPUSDCTokenPoolPackage{
		address:       packageId,
		usdcTokenPool: usdcTokenPoolContract,
	}, nil
}

func PublishCCIPUSDCTokenPool(
	ctx context.Context,
	opts *bind.CallOpts,
	client sui.ISuiAPI,
	ccipAddress,
	usdcCoinMetadataObjectId,
	tokenMessengerMinterPackageId,
	tokenMessengerMinterStateObjectId,
	messageTransmitterPackageId,
	messageTransmitterStateObjectId,
	treasuryObjectId,
	mcmsAddress,
	mcmsOwnerAddress, suiRPC string) (USDCTokenPool, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.USDCTokenPool, map[string]string{
		"ccip":                              ccipAddress,
		"usdc_token_pool":                   "0x0",
		"usdc_coin_metadata_object_id":      usdcCoinMetadataObjectId,
		"token_messenger_minter_package_id": tokenMessengerMinterPackageId,
		"token_messenger_minter_state":      tokenMessengerMinterStateObjectId,
		"message_transmitter_package_id":    messageTransmitterPackageId,
		"message_transmitter_state":         messageTransmitterStateObjectId,
		"treasury":                          treasuryObjectId,
		"mcms":                              mcmsAddress,
		"mcms_owner":                        mcmsOwnerAddress,

		"signer": signerAddr,
	}, false, suiRPC)
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
