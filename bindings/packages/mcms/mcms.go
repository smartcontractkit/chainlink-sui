package mcms

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mcms "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type MCMS interface {
	Address() string
	MCMS() module_mcms.IMcms
}

var _ MCMS = MCMSPackage{}

type MCMSPackage struct {
	address string

	mcms module_mcms.IMcms
}

func (p MCMSPackage) Address() string {
	return p.address
}

func (p MCMSPackage) MCMS() module_mcms.IMcms {
	return p.mcms
}

func NewMCMS(address string, chainClient client.BindingsClient) (MCMS, error) {
	mcmsContract, err := module_mcms.NewMcms(address, chainClient)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return MCMSPackage{
		address: packageId,
		mcms:    mcmsContract,
	}, nil
}

func PublishMCMS(ctx context.Context, opts *bind.CallOpts, chainClient client.BindingsClient, suiRPC string) (MCMS, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.MCMS, map[string]string{
		"mcms":       "0x0",
		"mcms_owner": "0x2",
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

	contract, err := NewMCMS(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
