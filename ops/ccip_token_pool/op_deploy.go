package tokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	tokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type TokenPoolDeployInput struct {
	CCIPPackageId    string
	MCMSAddress      string
	MCMSOwnerAddress string
}

type TokenPoolDeployOutput struct {
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TokenPoolDeployInput) (output sui_ops.OpTxResult[TokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := tokenpool.PublishCCIPTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[TokenPoolDeployOutput]{}, err
	}

	return sui_ops.OpTxResult[TokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageId: tokenPoolPackage.Address(),
	}, err
}

var DeployCCIPTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP tokenPool package",
	deployHandler,
)
