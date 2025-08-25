package linkops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	linktoken "github.com/smartcontractkit/chainlink-sui/bindings/generated/link/link"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployLinkObjects struct {
	CoinMetadataObjectID string
	TreasuryCapObjectID  string
	UpgradeCapObjectID   string
}

var handler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployLinkObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	artifact, err := bind.CompilePackage(contracts.LINK, map[string]string{
		"link": "0x0",
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployLinkObjects]{}, err
	}

	packageID, tx, err := bind.PublishPackage(b.GetContext(), opts, deps.Client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployLinkObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "coin", "CoinMetadata")
	if err1 != nil {
		return sui_ops.OpTxResult[DeployLinkObjects]{}, fmt.Errorf("failed to find CoinMetadata object ID in publish tx: %w", err1)
	}

	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "coin", "TreasuryCap")
	if err2 != nil {
		return sui_ops.OpTxResult[DeployLinkObjects]{}, fmt.Errorf("failed to find TreasuryCap object ID in publish tx: %w", err2)
	}

	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err3 != nil {
		return sui_ops.OpTxResult[DeployLinkObjects]{}, fmt.Errorf("failed to find UpgradeCap object ID in publish tx: %w", err3)
	}

	return sui_ops.OpTxResult[DeployLinkObjects]{
		Digest:    tx.Digest,
		PackageID: packageID,
		Objects: DeployLinkObjects{
			CoinMetadataObjectID: obj1,
			TreasuryCapObjectID:  obj2,
			UpgradeCapObjectID:   obj3,
		},
	}, err
}

type MintLinkTokenInput struct {
	LinkTokenPackageID string
	TreasuryCapID      string
	Amount             uint64
}

type MintLinkTokenOutput struct {
	MintedLinkTokenObjectID string
}

var handlerMint = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MintLinkTokenInput) (output sui_ops.OpTxResult[MintLinkTokenOutput], err error) {
	linkToken, err := linktoken.NewLink(input.LinkTokenPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[MintLinkTokenOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Get the signer address to transfer the minted coin to
	signerAddress, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to get signer address: %w", err)
	}

	// Use MintAndTransfer instead of Mint to ensure the coin is transferred and visible
	tx, err := linkToken.MintAndTransfer(b.GetContext(), opts, bind.Object{Id: input.TreasuryCapID}, input.Amount, signerAddress)
	if err != nil {
		return sui_ops.OpTxResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to execute MintAndTransfer on LinkToken: %w", err)
	}

	// Use the correct function for finding coin objects and provide the coin type
	coinType := input.LinkTokenPackageID + "::link::LINK"
	obj1, err1 := bind.FindCoinObjectIdFromTx(*tx, coinType)
	if err1 != nil {
		return sui_ops.OpTxResult[MintLinkTokenOutput]{}, fmt.Errorf("failed to find minted coin object: %w", err1)
	}

	return sui_ops.OpTxResult[MintLinkTokenOutput]{
		Digest:    tx.Digest,
		PackageID: input.LinkTokenPackageID,
		Objects: MintLinkTokenOutput{
			MintedLinkTokenObjectID: obj1,
		},
	}, err
}

var DeployLINKOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("link", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the LINK Token contract",
	handler,
)

var MintLinkOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("link", "package", "mint"),
	semver.MustParse("0.1.0"),
	"Mint the deployed LinkToken",
	handlerMint,
)
