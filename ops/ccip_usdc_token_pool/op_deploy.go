package usdctokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	usdctokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/usdc_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type USDCTokenPoolDeployInput struct {
	CCIPPackageId                     string
	CCIPTokenPoolPackageId            string
	USDCCoinMetadataObjectId          string
	TokenMessengerMinterPackageId     string
	TokenMessengerMinterStateObjectId string
	MessageTransmitterPackageId       string
	MessageTransmitterStateObjectId   string
	TreasuryObjectId                  string
	MCMSAddress                       string
	MCMSOwnerAddress                  string
}

type USDCTokenPoolDeployOutput struct {
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input USDCTokenPoolDeployInput) (output sui_ops.OpTxResult[USDCTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := usdctokenpool.PublishCCIPUSDCTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.CCIPTokenPoolPackageId,
		input.USDCCoinMetadataObjectId,
		input.TokenMessengerMinterPackageId,
		input.TokenMessengerMinterStateObjectId,
		input.MessageTransmitterPackageId,
		input.MessageTransmitterStateObjectId,
		input.TreasuryObjectId,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[USDCTokenPoolDeployOutput]{}, err
	}

	return sui_ops.OpTxResult[USDCTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageId: tokenPoolPackage.Address(),
	}, err
}

var DeployCCIPUSDCTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-usdc-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP USDC token pool package",
	deployHandler,
)
