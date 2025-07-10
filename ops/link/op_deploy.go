package linkops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	link "github.com/smartcontractkit/chainlink-sui/bindings/packages/link_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployLinkObjects struct {
	CoinMetadataObjectId string
	TreasuryCapObjectId  string
	UpgradeCapObjectId   string
}

var handler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployLinkObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcmsPackage, tx, err := link.PublishLinkToken(
		b.GetContext(),
		opts,
		deps.Client,
	)
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
		PackageId: mcmsPackage.Address(),
		Objects: DeployLinkObjects{
			CoinMetadataObjectId: obj1,
			TreasuryCapObjectId:  obj2,
			UpgradeCapObjectId:   obj3,
		},
	}, err
}

var DeployLINKOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("link", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the LINK Token contract",
	handler,
)
