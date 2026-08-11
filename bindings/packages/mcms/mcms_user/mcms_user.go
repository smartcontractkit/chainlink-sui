package mcmsuser

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mcms_user "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_user"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type MCMSUser interface {
	Address() string
	MCMSUser() module_mcms_user.IMcmsUser
}

var _ MCMSUser = MCMSUserPackage{}

type MCMSUserPackage struct {
	address string

	mcmsUser module_mcms_user.IMcmsUser
}

func (p MCMSUserPackage) Address() string {
	return p.address
}

func (p MCMSUserPackage) MCMSUser() module_mcms_user.IMcmsUser {
	return p.mcmsUser
}

func NewMCMSUser(address string, chainClient client.BindingsClient) (MCMSUser, error) {
	mcmsUserContract, err := module_mcms_user.NewMcmsUser(address, chainClient)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return MCMSUserPackage{
		address:  packageId,
		mcmsUser: mcmsUserContract,
	}, nil
}

func PublishMCMSUser(ctx context.Context, opts *bind.CallOpts, chainClient client.BindingsClient, mcmsAddress, mcmsOwnerAddress, suiRPC string) (MCMSUser, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.MCMSUser, map[string]string{
		"mcms_test":  "0x0",
		"mcms":       mcmsAddress,
		"mcms_owner": mcmsOwnerAddress,
		"signer":     signerAddr,
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

	contract, err := NewMCMSUser(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
