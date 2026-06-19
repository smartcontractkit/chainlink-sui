package usdctokenpool

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

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

func NewCCIPUSDCTokenPool(address string, chainClient client.BindingsClient) (USDCTokenPool, error) {
	usdcTokenPoolContract, err := module_usdc_token_pool.NewUsdcTokenPool(address, chainClient)
	if err != nil {
		return nil, err
	}

	packageID, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPUSDCTokenPoolPackage{
		address:       packageID,
		usdcTokenPool: usdcTokenPoolContract,
	}, nil
}

func PublishCCIPUSDCTokenPool(
	ctx context.Context,
	opts *bind.CallOpts,
	chainClient client.BindingsClient,
	ccipAddress,
	usdcCoinMetadataObjectID,
	tokenMessengerMinterPackageID,
	tokenMessengerMinterStateObjectID,
	messageTransmitterPackageID,
	messageTransmitterStateObjectID,
	treasuryObjectID,
	mcmsAddress,
	fastMcmsAddress,
	mcmsOwnerAddress, suiRPC string) (USDCTokenPool, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.USDCTokenPool, map[string]string{
		"ccip":                              ccipAddress,
		"usdc_token_pool":                   "0x0",
		"usdc_coin_metadata_object_id":      usdcCoinMetadataObjectID,
		"token_messenger_minter_package_id": tokenMessengerMinterPackageID,
		"token_messenger_minter_state":      tokenMessengerMinterStateObjectID,
		"message_transmitter_package_id":    messageTransmitterPackageID,
		"message_transmitter_state":         messageTransmitterStateObjectID,
		"treasury":                          treasuryObjectID,
		"mcms":                              mcmsAddress,
		"fast_mcms":                         fastMcmsAddress,
		"mcms_owner":                        mcmsOwnerAddress,

		"signer": signerAddr,
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

	contract, err := NewCCIPUSDCTokenPool(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
