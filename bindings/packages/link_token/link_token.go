package link

import (
	"context"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_link_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/link_token/link_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type LinkToken interface {
	Address() sui.Address
	LinkToken() module_link_token.ILinkToken
}

var _ LinkToken = LinkTokenPackage{}

type LinkTokenPackage struct {
	address sui.Address

	linkToken module_link_token.ILinkToken
}

func (p LinkTokenPackage) Address() sui.Address {
	return p.address
}

func (p LinkTokenPackage) LinkToken() module_link_token.ILinkToken {
	return p.linkToken
}

func NewLinkToken(address string, client suiclient.ClientImpl) (LinkToken, error) {
	pkgObjectId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	linkTokenContract, err := module_link_token.NewLinkToken(address, client)
	if err != nil {
		return nil, err
	}

	return LinkTokenPackage{
		address:   *pkgObjectId,
		linkToken: linkTokenContract,
	}, nil
}

func PublishLinkToken(ctx context.Context, opts bind.TxOpts, signer rel.SuiSigner, client suiclient.ClientImpl) (LinkToken, *suiclient.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.LINKToken, map[string]string{
		"link_token": "0x0",
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

	contract, err := NewLinkToken(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
