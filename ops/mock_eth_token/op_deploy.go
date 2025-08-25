package mockethtokenops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mockethtoken "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/mock_eth_token/mock_eth_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployMockEthTokenObjects struct {
	CoinMetadataObjectID string
	TreasuryCapObjectID  string
	UpgradeCapObjectID   string
}

var handler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployMockEthTokenObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	artifact, err := bind.CompilePackage(contracts.MockEthToken, map[string]string{
		"mock_eth_token": "0x0",
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployMockEthTokenObjects]{}, err
	}

	packageID, tx, err := bind.PublishPackage(b.GetContext(), opts, deps.Client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployMockEthTokenObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "coin", "CoinMetadata")
	if err1 != nil {
		return sui_ops.OpTxResult[DeployMockEthTokenObjects]{}, fmt.Errorf("failed to find CoinMetadata object ID in publish tx: %w", err1)
	}

	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "coin", "TreasuryCap")
	if err2 != nil {
		return sui_ops.OpTxResult[DeployMockEthTokenObjects]{}, fmt.Errorf("failed to find TreasuryCap object ID in publish tx: %w", err2)
	}

	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err3 != nil {
		return sui_ops.OpTxResult[DeployMockEthTokenObjects]{}, fmt.Errorf("failed to find UpgradeCap object ID in publish tx: %w", err3)
	}

	return sui_ops.OpTxResult[DeployMockEthTokenObjects]{
		Digest:    tx.Digest,
		PackageID: packageID,
		Objects: DeployMockEthTokenObjects{
			CoinMetadataObjectID: obj1,
			TreasuryCapObjectID:  obj2,
			UpgradeCapObjectID:   obj3,
		},
	}, err
}

type MintMockEthTokenInput struct {
	MockEthTokenPackageID string
	TreasuryCapID         string
	Amount                uint64
}

type MintMockEthTokenOutput struct {
	MintedMockEthTokenObjectID string
}

var handlerMint = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MintMockEthTokenInput) (output sui_ops.OpTxResult[MintMockEthTokenOutput], err error) {
	mockEthToken, err := mockethtoken.NewMockEthToken(input.MockEthTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[MintMockEthTokenOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Get the signer address to transfer the minted coin to
	signerAddress, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[MintMockEthTokenOutput]{}, fmt.Errorf("failed to get signer address: %w", err)
	}

	// Use MintAndTransfer instead of Mint to ensure the coin is transferred and visible
	tx, err := mockEthToken.MintAndTransfer(b.GetContext(), opts, bind.Object{Id: input.TreasuryCapID}, input.Amount, signerAddress)
	if err != nil {
		return sui_ops.OpTxResult[MintMockEthTokenOutput]{}, fmt.Errorf("failed to execute MintAndTransfer on MockEthToken: %w", err)
	}

	// Use the correct function for finding coin objects and provide the coin type
	coinType := input.MockEthTokenPackageID + "::mock_eth_token::MOCK_ETH_TOKEN"
	obj1, err1 := bind.FindCoinObjectIdFromTx(*tx, coinType)
	if err1 != nil {
		return sui_ops.OpTxResult[MintMockEthTokenOutput]{}, fmt.Errorf("failed to find minted coin object: %w", err1)
	}

	return sui_ops.OpTxResult[MintMockEthTokenOutput]{
		Digest:    tx.Digest,
		PackageID: input.MockEthTokenPackageID,
		Objects: MintMockEthTokenOutput{
			MintedMockEthTokenObjectID: obj1,
		},
	}, err
}

var DeployMockEthTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mock_eth_token", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the Mock ETH Token contract",
	handler,
)

var MintMockEthTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mock_eth_token", "package", "mint"),
	semver.MustParse("0.1.0"),
	"Mint the deployed MockEthToken",
	handlerMint,
)
