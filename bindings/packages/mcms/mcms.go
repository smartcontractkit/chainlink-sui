package mcms

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mcms "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms"

	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type MCMS interface {
	Address() sui.Address
	MCMS() module_mcms.IMcms
}

var _ MCMS = MCMSPackage{}

type MCMSPackage struct {
	address sui.Address

	mcms module_mcms.IMcms
}

func (p MCMSPackage) Address() sui.Address {
	return p.address
}

func (p MCMSPackage) MCMS() module_mcms.IMcms {
	return p.mcms
}

func NewMCMS(address string, client suiclient.ClientImpl) (MCMS, error) {
	mcmsContract, err := module_mcms.NewMcms(address, client)
	if err != nil {
		return nil, err
	}

	packageId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	return MCMSPackage{
		address: *packageId,
		mcms:    mcmsContract,
	}, nil
}

func PublishMCMS(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (MCMS, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.MCMS, map[string]string{
		"mcms":       "0x0",
		"mcms_owner": "0x2",
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

	contract, err := NewMCMS(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
