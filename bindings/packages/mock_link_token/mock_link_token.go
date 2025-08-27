package mocklinktoken

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mock_link_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/mock_link_token/mock_link_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type MockLinkToken interface {
	Address() string
	MockLinkToken() module_mock_link_token.IMockLinkToken
}

var _ MockLinkToken = MockLinkTokenPackage{}

type MockLinkTokenPackage struct {
	address string

	mockLinkToken module_mock_link_token.IMockLinkToken
}

func (p MockLinkTokenPackage) Address() string {
	return p.address
}

func (p MockLinkTokenPackage) MockLinkToken() module_mock_link_token.IMockLinkToken {
	return p.mockLinkToken
}

func NewMockLinkToken(address string, client sui.ISuiAPI) (MockLinkToken, error) {
	pkgObjectId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	mockLinkTokenContract, err := module_mock_link_token.NewMockLinkToken(address, client)
	if err != nil {
		return nil, err
	}

	return MockLinkTokenPackage{
		address:       pkgObjectId,
		mockLinkToken: mockLinkTokenContract,
	}, nil
}

func PublishMockLinkToken(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI) (MockLinkToken, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.MockLinkToken, map[string]string{
		"mock_link_token": "0x0",
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

	contract, err := NewMockLinkToken(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
