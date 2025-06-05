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
	mcmsPackage, tx, err := link.PublishLinkToken(
		b.GetContext(),
		deps.GetTxOpts(),
		deps.Signer,
		deps.Client,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployLinkObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "coin", "CoinMetadata")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "coin", "TreasuryCap")
	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")

	if err1 != nil || err2 != nil || err3 != nil {
		return sui_ops.OpTxResult[DeployLinkObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployLinkObjects]{
		Digest:    tx.Digest.String(),
		PackageId: mcmsPackage.Address().String(),
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
