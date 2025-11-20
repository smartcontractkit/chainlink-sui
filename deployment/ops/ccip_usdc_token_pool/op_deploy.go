package usdctokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_usdctokenpool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/usdc_token_pool"
	usdctokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/usdc_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type USDCTokenPoolDeployInput struct {
	CCIPPackageId                     string
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
	OwnerCapObjectId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input USDCTokenPoolDeployInput) (output sui_ops.OpTxResult[USDCTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := usdctokenpool.PublishCCIPUSDCTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.USDCCoinMetadataObjectId,
		input.TokenMessengerMinterPackageId,
		input.TokenMessengerMinterStateObjectId,
		input.MessageTransmitterPackageId,
		input.MessageTransmitterStateObjectId,
		input.TreasuryObjectId,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[USDCTokenPoolDeployOutput]{}, err
	}

	ownerCapObj, err := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	if err != nil {
		return sui_ops.OpTxResult[USDCTokenPoolDeployOutput]{}, fmt.Errorf("failed to find OwnerCap object ID: %w", err)
	}

	return sui_ops.OpTxResult[USDCTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageId: tokenPoolPackage.Address(),
		Objects: USDCTokenPoolDeployOutput{
			OwnerCapObjectId: ownerCapObj,
		},
	}, err
}

var DeployCCIPUSDCTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-usdc-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP USDC token pool package",
	deployHandler,
)

type TransferOwnershipUsdcTokenPoolInput struct {
	UsdcTokenPoolPackageId string
	TypeArgs               []string
	StateObjectId          string
	OwnerCapObjectId       string
	To                     string
}

type TransferOwnershipUsdcTokenPoolObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipUsdcTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipUsdcTokenPoolInput) (output sui_ops.OpTxResult[TransferOwnershipUsdcTokenPoolObjects], err error) {
	usdcTokenPoolPackage, err := module_usdctokenpool.NewUsdcTokenPool(input.UsdcTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipUsdcTokenPoolObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := usdcTokenPoolPackage.TransferOwnership(
		b.GetContext(),
		opts,
		input.TypeArgs,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipUsdcTokenPoolObjects]{}, fmt.Errorf("failed to execute TransferOwnership on UsdcTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for UsdcTokenPool", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipUsdcTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.UsdcTokenPoolPackageId,
		Objects:   TransferOwnershipUsdcTokenPoolObjects{},
	}, nil
}

var TransferOwnershipUsdcTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-usdc-token-pool-transfer-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the UsdcTokenPool",
	transferOwnershipUsdcTokenPoolHandler,
)

type AcceptOwnershipUsdcTokenPoolInput struct {
	UsdcTokenPoolPackageId string
	TypeArgs               []string
	StateObjectId          string
}

type AcceptOwnershipUsdcTokenPoolObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipUsdcTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipUsdcTokenPoolInput) (output sui_ops.OpTxResult[AcceptOwnershipUsdcTokenPoolObjects], err error) {
	usdcTokenPoolPackage, err := module_usdctokenpool.NewUsdcTokenPool(input.UsdcTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipUsdcTokenPoolObjects]{}, err
	}

	encodedCall, err := usdcTokenPoolPackage.Encoder().AcceptOwnership(input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipUsdcTokenPoolObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipUsdcTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on UsdcTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[AcceptOwnershipUsdcTokenPoolObjects]{
			Digest:    "",
			PackageId: input.UsdcTokenPoolPackageId,
			Objects:   AcceptOwnershipUsdcTokenPoolObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := usdcTokenPoolPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipUsdcTokenPoolObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on UsdcTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership accepted for UsdcTokenPool")

	return sui_ops.OpTxResult[AcceptOwnershipUsdcTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.UsdcTokenPoolPackageId,
		Objects:   AcceptOwnershipUsdcTokenPoolObjects{},
		Call:      call,
	}, nil
}

var AcceptOwnershipUsdcTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-usdc-token-pool-accept-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the UsdcTokenPool",
	acceptOwnershipUsdcTokenPoolHandler,
)
