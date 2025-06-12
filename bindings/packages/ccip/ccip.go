package ccip

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	module_nonce_manager "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/nonce_manager"
	module_receiver_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/receiver_registry"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type CCIP interface {
	Address() sui.Address
	FeeQuoter() module_fee_quoter.IFeeQuoter
	NonceManager() module_nonce_manager.INonceManager
	ReceiverRegistry() module_receiver_registry.IReceiverRegistry
	RmnRemote() module_rmn_remote.IRmnRemote
	TokenAdminRegistry() module_token_admin_registry.ITokenAdminRegistry
}

var _ CCIP = CCIPPackage{}

type CCIPPackage struct {
	address sui.Address

	feeQuoter          module_fee_quoter.IFeeQuoter
	nonceManager       module_nonce_manager.INonceManager
	receiverRegistry   module_receiver_registry.IReceiverRegistry
	rmnRemove          module_rmn_remote.IRmnRemote
	tokenAdminRegistry module_token_admin_registry.ITokenAdminRegistry
}

func (p CCIPPackage) Address() sui.Address {
	return p.address
}

func (p CCIPPackage) FeeQuoter() module_fee_quoter.IFeeQuoter {
	return p.feeQuoter
}

func (p CCIPPackage) NonceManager() module_nonce_manager.INonceManager {
	return p.nonceManager
}

func (p CCIPPackage) ReceiverRegistry() module_receiver_registry.IReceiverRegistry {
	return p.receiverRegistry
}

func (p CCIPPackage) RmnRemote() module_rmn_remote.IRmnRemote {
	return p.rmnRemove
}

func (p CCIPPackage) TokenAdminRegistry() module_token_admin_registry.ITokenAdminRegistry {
	return p.tokenAdminRegistry
}

func NewCCIP(address string, client suiclient.ClientImpl) (CCIP, error) {
	feeQuoterContract, err := module_fee_quoter.NewFeeQuoter(address, client)
	if err != nil {
		return nil, err
	}

	nonceManagerContract, err := module_nonce_manager.NewNonceManager(address, client)
	if err != nil {
		return nil, err
	}

	receiverRegistryContract, err := module_receiver_registry.NewReceiverRegistry(address, client)
	if err != nil {
		return nil, err
	}

	rmnRemoteContract, err := module_rmn_remote.NewRmnRemote(address, client)
	if err != nil {
		return nil, err
	}

	tokenAdminRegistryContract, err := module_token_admin_registry.NewTokenAdminRegistry(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPPackage{
		address:            *packageId,
		feeQuoter:          feeQuoterContract,
		nonceManager:       nonceManagerContract,
		receiverRegistry:   receiverRegistryContract,
		rmnRemove:          rmnRemoteContract,
		tokenAdminRegistry: tokenAdminRegistryContract,
	}, nil
}

func PublishCCIP(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl, mcmsAddress string, mcmsOwner string) (CCIP, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.CCIP, map[string]string{
		"mcms":       mcmsAddress,
		"mcms_owner": mcmsOwner,
		"ccip":       "0x0",
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

	contract, err := NewCCIP(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
