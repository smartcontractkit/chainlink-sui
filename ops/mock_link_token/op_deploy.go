package mocklinktokenops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mocklinktoken "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/mock_link_token/mock_link_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployMockLinkTokenObjects struct {
	CoinMetadataObjectID string
	TreasuryCapObjectID  string
	UpgradeCapObjectID   string
}

var handler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployMockLinkTokenObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	artifact, err := bind.CompilePackage(contracts.MockLinkToken, map[string]string{
		"mock_link_token": "0x0",
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployMockLinkTokenObjects]{}, err
	}

	packageID, tx, err := bind.PublishPackage(b.GetContext(), opts, deps.Client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployMockLinkTokenObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "coin", "CoinMetadata")
	if err1 != nil {
		return sui_ops.OpTxResult[DeployMockLinkTokenObjects]{}, fmt.Errorf("failed to find CoinMetadata object ID in publish tx: %w", err1)
	}

	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "coin", "TreasuryCap")
	if err2 != nil {
		return sui_ops.OpTxResult[DeployMockLinkTokenObjects]{}, fmt.Errorf("failed to find TreasuryCap object ID in publish tx: %w", err2)
	}

	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err3 != nil {
		return sui_ops.OpTxResult[DeployMockLinkTokenObjects]{}, fmt.Errorf("failed to find UpgradeCap object ID in publish tx: %w", err3)
	}

	return sui_ops.OpTxResult[DeployMockLinkTokenObjects]{
		Digest:    tx.Digest,
		PackageID: packageID,
		Objects: DeployMockLinkTokenObjects{
			CoinMetadataObjectID: obj1,
			TreasuryCapObjectID:  obj2,
			UpgradeCapObjectID:   obj3,
		},
	}, err
}

type MintMockLinkTokenInput struct {
	MockLinkTokenPackageID string
	TreasuryCapID          string
	Amount                 uint64
}

type MintMockLinkTokenOutput struct {
	MintedMockLinkTokenObjectID string
}

var handlerMint = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MintMockLinkTokenInput) (output sui_ops.OpTxResult[MintMockLinkTokenOutput], err error) {
	mockLinkToken, err := mocklinktoken.NewMockLinkToken(input.MockLinkTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[MintMockLinkTokenOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Get the signer address to transfer the minted coin to
	signerAddress, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[MintMockLinkTokenOutput]{}, fmt.Errorf("failed to get signer address: %w", err)
	}

	// Use MintAndTransfer instead of Mint to ensure the coin is transferred and visible
	tx, err := mockLinkToken.MintAndTransfer(b.GetContext(), opts, bind.Object{Id: input.TreasuryCapID}, input.Amount, signerAddress)
	if err != nil {
		return sui_ops.OpTxResult[MintMockLinkTokenOutput]{}, fmt.Errorf("failed to execute MintAndTransfer on MockLinkToken: %w", err)
	}

	// Use the correct function for finding coin objects and provide the coin type
	coinType := input.MockLinkTokenPackageID + "::mock_link_token::MOCK_LINK_TOKEN"
	obj1, err1 := bind.FindCoinObjectIdFromTx(*tx, coinType)
	if err1 != nil {
		return sui_ops.OpTxResult[MintMockLinkTokenOutput]{}, fmt.Errorf("failed to find minted coin object: %w", err1)
	}

	return sui_ops.OpTxResult[MintMockLinkTokenOutput]{
		Digest:    tx.Digest,
		PackageID: input.MockLinkTokenPackageID,
		Objects: MintMockLinkTokenOutput{
			MintedMockLinkTokenObjectID: obj1,
		},
	}, err
}

var DeployMockLinkTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mock_link_token", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the Mock LINK Token contract",
	handler,
)

var MintMockLinkTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mock_link_token", "package", "mint"),
	semver.MustParse("0.1.0"),
	"Mint the deployed MockLinkToken",
	handlerMint,
)
