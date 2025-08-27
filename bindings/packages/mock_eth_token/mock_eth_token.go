package mockethtoken

import (
	"context"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mock_eth_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/mock_eth_token/mock_eth_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
)

type MockEthToken interface {
	Address() string
	MockEthToken() module_mock_eth_token.IMockEthToken
}

var _ MockEthToken = MockEthTokenPackage{}

type MockEthTokenPackage struct {
	address string

	mockEthToken module_mock_eth_token.IMockEthToken
}

func (p MockEthTokenPackage) Address() string {
	return p.address
}

func (p MockEthTokenPackage) MockEthToken() module_mock_eth_token.IMockEthToken {
	return p.mockEthToken
}

func NewMockEthToken(address string, client sui.ISuiAPI) (MockEthToken, error) {
	pkgObjectId, err := bind.ToSuiAddress(address)
	if err != nil {
		return nil, err
	}

	mockEthTokenContract, err := module_mock_eth_token.NewMockEthToken(address, client)
	if err != nil {
		return nil, err
	}

	return MockEthTokenPackage{
		address:      pkgObjectId,
		mockEthToken: mockEthTokenContract,
	}, nil
}

func PublishMockEthToken(ctx context.Context, opts *bind.CallOpts, client sui.ISuiAPI) (MockEthToken, *models.SuiTransactionBlockResponse, error) {
	artifact, err := bind.CompilePackage(contracts.MockEthToken, map[string]string{
		"mock_eth_token": "0x0",
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

	contract, err := NewMockEthToken(packageId, client)
	if err != nil {
		return nil, nil, err
	}

	return contract, tx, nil
}
