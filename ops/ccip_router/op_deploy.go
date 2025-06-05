package routerops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/packages/router"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	routerPackage, tx, err := router.PublishCCIPRouter(
		b.GetContext(),
		deps.GetTxOpts(),
		deps.Signer,
		deps.Client,
	)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    tx.Digest.String(),
		PackageId: routerPackage.Address().String(),
	}, err
}

var DeployCCIPRouterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP router package",
	deployHandler,
)
