package lockreleasetokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	lockreleasetokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/lock_release_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type LockReleaseTokenPoolDeployInput struct {
	CCIPPackageId          string
	CCIPTokenPoolPackageId string
	MCMSAddress            string
	MCMSOwnerAddress       string
}

type LockReleaseTokenPoolDeployOutput struct {
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolDeployInput) (output sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := lockreleasetokenpool.PublishCCIPLockReleaseTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.CCIPTokenPoolPackageId,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput]{}, err
	}

	return sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageId: tokenPoolPackage.Address(),
	}, err
}

var DeployCCIPLockReleaseTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-lock-release-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP lock release token pool package",
	deployHandler,
)
