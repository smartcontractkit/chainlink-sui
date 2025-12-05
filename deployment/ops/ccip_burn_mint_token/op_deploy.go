package bnmops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	bnmtoken "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_burn_mint_token/ccip_burn_mint_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployBnMObjects struct {
	CoinMetadataObjectId string
	TreasuryCapObjectId  string
	UpgradeCapObjectId   string
}

var DeployBnMOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("bnm", "token", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the BnM Token contract",
	deployBnMOp,
)

var deployBnMOp = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployBnMObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[DeployBnMObjects]{}, err
	}

	artifact, err := bind.CompilePackage(contracts.CCIPBnM, map[string]string{
		"ccip_burn_mint_token": "0x0",
		"signer":               signerAddr,
	}, false, deps.SuiRPC)
	if err != nil {
		return sui_ops.OpTxResult[DeployBnMObjects]{}, err
	}

	packageId, tx, err := bind.PublishPackage(b.GetContext(), opts, deps.Client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployBnMObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "coin", "CoinMetadata")
	if err1 != nil {
		return sui_ops.OpTxResult[DeployBnMObjects]{}, fmt.Errorf("failed to find CoinMetadata object ID in publish tx: %w", err1)
	}

	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "coin", "TreasuryCap")
	if err2 != nil {
		return sui_ops.OpTxResult[DeployBnMObjects]{}, fmt.Errorf("failed to find TreasuryCap object ID in publish tx: %w", err2)
	}

	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err3 != nil {
		return sui_ops.OpTxResult[DeployBnMObjects]{}, fmt.Errorf("failed to find UpgradeCap object ID in publish tx: %w", err3)
	}

	return sui_ops.OpTxResult[DeployBnMObjects]{
		Digest:    tx.Digest,
		PackageId: packageId,
		Objects: DeployBnMObjects{
			CoinMetadataObjectId: obj1,
			TreasuryCapObjectId:  obj2,
			UpgradeCapObjectId:   obj3,
		},
	}, err
}

type MintBnMTokenInput struct {
	BnMTokenPackageId string
	TreasuryCapId     string
	Amount            uint64
	ToAddress         string
}

type MintBnMTokenOutput struct {
	MintedBnMTokenObjectId string
}

var MintBnMOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("bnm", "token", "mint"),
	semver.MustParse("0.1.0"),
	"Mint the deployed BnM Token",
	mintBnMOp,
)

var mintBnMOp = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MintBnMTokenInput) (output sui_ops.OpTxResult[MintBnMTokenOutput], err error) {
	bnmToken, err := bnmtoken.NewCcipBurnMintToken(input.BnMTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[MintBnMTokenOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Use MintAndTransfer instead of Mint to ensure the coin is transferred and visible
	tx, err := bnmToken.MintAndTransfer(b.GetContext(), opts, bind.Object{Id: input.TreasuryCapId}, input.Amount, input.ToAddress)
	if err != nil {
		return sui_ops.OpTxResult[MintBnMTokenOutput]{}, fmt.Errorf("failed to execute MintAndTransfer on BnMToken: %w", err)
	}

	// Use the correct function for finding coin objects and provide the coin type
	coinType := fmt.Sprintf("%s::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", input.BnMTokenPackageId)
	obj1, err1 := bind.FindCoinObjectIdFromTx(*tx, coinType)
	if err1 != nil {
		return sui_ops.OpTxResult[MintBnMTokenOutput]{}, fmt.Errorf("failed to find minted coin object: %w", err1)
	}

	return sui_ops.OpTxResult[MintBnMTokenOutput]{
		Digest:    tx.Digest,
		PackageId: input.BnMTokenPackageId,
		Objects: MintBnMTokenOutput{
			MintedBnMTokenObjectId: obj1,
		},
	}, err
}
