package link

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_link "github.com/smartcontractkit/chainlink-sui/bindings/generated/link/link"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type Link interface {
	Address() string
	Link() module_link.ILink
}

var _ Link = LinkPackage{}

type LinkPackage struct {
	address string

	linkToken module_link.ILink
}

func (p LinkPackage) Address() string {
	return p.address
}

func (p LinkPackage) Link() module_link.ILink {
	return p.linkToken
}

func NewLink(address string, client sui.ISuiAPI) (Link, error) {
	pkgObjectId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	linkTokenContract, err := module_link.NewLink(address, client)
	if err != nil {
		return nil, err
	}

	return LinkPackage{
		address:   pkgObjectId,
		linkToken: linkTokenContract,
	}, nil
}

func PublishLink(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI) (Link, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.LINK, map[string]string{
		"link": "0x0",
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

	contract, err := NewLink(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
