package routerops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/router"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployCCIPRouterInput struct {
	McmsPackageId string
	McmsOwner     string
}
type DeployCCIPRouterObjects struct {
	OwnerCapObjectId    string
	RouterStateObjectId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployCCIPRouterInput) (output sui_ops.OpTxResult[DeployCCIPRouterObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	routerPackage, tx, err := router.PublishCCIPRouter(
		b.GetContext(),
		opts,
		deps.Client,
		input.McmsPackageId,
		input.McmsOwner,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "router", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "router", "RouterState")
	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPRouterObjects]{
		Digest:    tx.Digest,
		PackageId: routerPackage.Address(),
		Objects: DeployCCIPRouterObjects{
			OwnerCapObjectId:    obj1,
			RouterStateObjectId: obj2,
		},
	}, err
}

var DeployCCIPRouterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP router package",
	deployHandler,
)
