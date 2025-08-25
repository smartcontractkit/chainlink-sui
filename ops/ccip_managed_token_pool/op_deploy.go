package managedtokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	managedtokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/managed_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type ManagedTokenPoolDeployInput struct {
	CCIPPackageID          string
	CCIPTokenPoolPackageID string
	ManagedTokenPackageID  string
	MCMSAddress            string
	MCMSOwnerAddress       string
}

type ManagedTokenPoolDeployOutput struct {
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolDeployInput) (output sui_ops.OpTxResult[ManagedTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := managedtokenpool.PublishCCIPManagedTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageID,
		input.CCIPTokenPoolPackageID,
		input.ManagedTokenPackageID,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolDeployOutput]{}, err
	}

	return sui_ops.OpTxResult[ManagedTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageID: tokenPoolPackage.Address(),
	}, err
}

var DeployCCIPManagedTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP managed token pool package",
	deployHandler,
)
