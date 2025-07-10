package link

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_link_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/link_token/link_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type LinkToken interface {
	Address() string
	LinkToken() module_link_token.ILinkToken
}

var _ LinkToken = LinkTokenPackage{}

type LinkTokenPackage struct {
	address string

	linkToken module_link_token.ILinkToken
}

func (p LinkTokenPackage) Address() string {
	return p.address
}

func (p LinkTokenPackage) LinkToken() module_link_token.ILinkToken {
	return p.linkToken
}

func NewLinkToken(address string, client sui.ISuiAPI) (LinkToken, error) {
	pkgObjectId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	linkTokenContract, err := module_link_token.NewLinkToken(address, client)
	if err != nil {
		return nil, err
	}

	return LinkTokenPackage{
		address:   pkgObjectId,
		linkToken: linkTokenContract,
	}, nil
}

func PublishLinkToken(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI) (LinkToken, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.LINKToken, map[string]string{
		"link_token": "0x0",
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

	contract, err := NewLinkToken(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
