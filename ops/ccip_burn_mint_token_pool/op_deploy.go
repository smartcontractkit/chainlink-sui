package burnminttokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	burnminttokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/burn_mint_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type BurnMintTokenPoolDeployInput struct {
	CCIPPackageID          string
	CCIPTokenPoolPackageID string
	MCMSAddress            string
	MCMSOwnerAddress       string
}

type BurnMintTokenPoolDeployOutput struct {
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolDeployInput) (output sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := burnminttokenpool.PublishCCIPBurnMintTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageID,
		input.CCIPTokenPoolPackageID,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput]{}, err
	}

	return sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageID: tokenPoolPackage.Address(),
	}, err
}

var DeployCCIPBurnMintTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-burn-mint-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP burn mint token pool package",
	deployHandler,
)
