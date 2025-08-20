package managedtoken

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managed_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/managed_token/managed_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type ManagedToken interface {
	Address() string
	ManagedToken() module_managed_token.IManagedToken
}

var _ ManagedToken = CCIPManagedTokenPackage{}

type CCIPManagedTokenPackage struct {
	address string

	managedToken module_managed_token.IManagedToken
}

func (p CCIPManagedTokenPackage) Address() string {
	return p.address
}

func (p CCIPManagedTokenPackage) ManagedToken() module_managed_token.IManagedToken {
	return p.managedToken
}

func NewCCIPManagedToken(address string, client sui.ISuiAPI) (ManagedToken, error) {
	managedTokenContract, err := module_managed_token.NewManagedToken(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return CCIPManagedTokenPackage{
		address:      packageId,
		managedToken: managedTokenContract,
	}, nil
}

func PublishCCIPManagedToken(
	ctx context.Context,
	opts *bind.CallOpts,
	client sui.ISuiAPI,
	mcmsAddress,
	mcmsOwnerAddress string) (ManagedToken, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.ManagedToken, map[string]string{
		"managed_token": "0x0",
		"mcms":          mcmsAddress,
		"mcms_owner":    mcmsOwnerAddress,
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

	contract, err := NewCCIPManagedToken(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
