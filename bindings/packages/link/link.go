package link

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

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

func NewLink(address string, chainClient client.BindingsClient) (Link, error) {
	pkgObjectId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	linkTokenContract, err := module_link.NewLink(address, chainClient)
	if err != nil {
		return nil, err
	}

	return LinkPackage{
		address:   pkgObjectId,
		linkToken: linkTokenContract,
	}, nil
}

func PublishLink(ctx context.Context, opts *bind.CallOpts, chainClient client.BindingsClient, suiRPC string) (Link, *models.SuiTransactionBlockResponse, error) {
	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return nil, nil, err
	}

	artifact, err := bind.CompilePackage(contracts.LINK, map[string]string{
		"link":   "0x0",
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

	contract, err := NewLink(packageId, chainClient)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
