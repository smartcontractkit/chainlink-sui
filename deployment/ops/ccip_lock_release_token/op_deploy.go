package lnrops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	lnrtoken "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_lock_release_token/ccip_lock_release_token"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployLnRObjects struct {
	CoinMetadataObjectId string
	TreasuryCapObjectId  string
	UpgradeCapObjectId   string
}

var DeployLnROp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("lnr", "token", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the LnR Token contract",
	deployLnROp,
)

var deployLnROp = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployLnRObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[DeployLnRObjects]{}, err
	}

	artifact, err := bind.CompilePackage(contracts.CCIPLnR, map[string]string{
		"ccip_lock_release_token": "0x0",
		"signer":                  signerAddr,
	}, false, deps.SuiRPC)
	if err != nil {
		return sui_ops.OpTxResult[DeployLnRObjects]{}, err
	}

	packageId, tx, err := bind.PublishPackage(b.GetContext(), opts, deps.Client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployLnRObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "coin", "CoinMetadata")
	if err1 != nil {
		return sui_ops.OpTxResult[DeployLnRObjects]{}, fmt.Errorf("failed to find CoinMetadata object ID in publish tx: %w", err1)
	}

	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "coin", "TreasuryCap")
	if err2 != nil {
		return sui_ops.OpTxResult[DeployLnRObjects]{}, fmt.Errorf("failed to find TreasuryCap object ID in publish tx: %w", err2)
	}

	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err3 != nil {
		return sui_ops.OpTxResult[DeployLnRObjects]{}, fmt.Errorf("failed to find UpgradeCap object ID in publish tx: %w", err3)
	}

	return sui_ops.OpTxResult[DeployLnRObjects]{
		Digest:    tx.Digest,
		PackageId: packageId,
		Objects: DeployLnRObjects{
			CoinMetadataObjectId: obj1,
			TreasuryCapObjectId:  obj2,
			UpgradeCapObjectId:   obj3,
		},
	}, err
}

type MintLnRTokenInput struct {
	LnRTokenPackageId string
	TreasuryCapId     string
	Amount            uint64
	ToAddress         string
}

type MintLnRTokenOutput struct {
	MintedLnRTokenObjectId string
}

var MintLnROp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("lnr", "token", "mint"),
	semver.MustParse("0.1.0"),
	"Mint the deployed LnR Token",
	mintLnROp,
)

var mintLnROp = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MintLnRTokenInput) (output sui_ops.OpTxResult[MintLnRTokenOutput], err error) {
	lnrToken, err := lnrtoken.NewCcipLockReleaseToken(input.LnRTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[MintLnRTokenOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Use MintAndTransfer instead of Mint to ensure the coin is transferred and visible
	tx, err := lnrToken.MintAndTransfer(b.GetContext(), opts, bind.Object{Id: input.TreasuryCapId}, input.Amount, input.ToAddress)
	if err != nil {
		return sui_ops.OpTxResult[MintLnRTokenOutput]{}, fmt.Errorf("failed to execute MintAndTransfer on LnRToken: %w", err)
	}

	// Use the correct function for finding coin objects and provide the coin type
	coinType := fmt.Sprintf("%s::ccip_lock_release_token::CCIP_LOCK_RELEASE_TOKEN", input.LnRTokenPackageId)
	obj1, err1 := bind.FindCoinObjectIdFromTx(*tx, coinType)
	if err1 != nil {
		return sui_ops.OpTxResult[MintLnRTokenOutput]{}, fmt.Errorf("failed to find minted coin object: %w", err1)
	}

	return sui_ops.OpTxResult[MintLnRTokenOutput]{
		Digest:    tx.Digest,
		PackageId: input.LnRTokenPackageId,
		Objects: MintLnRTokenOutput{
			MintedLnRTokenObjectId: obj1,
		},
	}, err
}
