package burnminttokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burnminttokenpool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	burnminttokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/burn_mint_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type BurnMintTokenPoolDeployInput struct {
	CCIPPackageId    string
	MCMSAddress      string
	MCMSOwnerAddress string
}

type BurnMintTokenPoolDeployOutput struct {
	OwnerCapObjectId   string
	UpgradeCapObjectId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolDeployInput) (output sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := burnminttokenpool.PublishCCIPBurnMintTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput]{}, err
	}

	ownerCapObj, err := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput]{}, fmt.Errorf("failed to find OwnerCap object ID: %w", err)
	}

	upgradeCapObj, err := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput]{}, fmt.Errorf("failed to find UpgradeCap object ID: %w", err)
	}

	return sui_ops.OpTxResult[BurnMintTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageId: tokenPoolPackage.Address(),
		Objects: BurnMintTokenPoolDeployOutput{
			OwnerCapObjectId:   ownerCapObj,
			UpgradeCapObjectId: upgradeCapObj,
		},
	}, err
}

var DeployCCIPBurnMintTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-burn-mint-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP burn mint token pool package",
	deployHandler,
)

type TransferOwnershipBurnMintTokenPoolInput struct {
	BurnMintTokenPoolPackageId string
	TypeArgs                   []string
	StateObjectId              string
	OwnerCapObjectId           string
	To                         string
}

type TransferOwnershipBurnMintTokenPoolObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipBurnMintTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipBurnMintTokenPoolInput) (output sui_ops.OpTxResult[TransferOwnershipBurnMintTokenPoolObjects], err error) {
	burnMintTokenPoolPackage, err := module_burnminttokenpool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipBurnMintTokenPoolObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := burnMintTokenPoolPackage.TransferOwnership(
		b.GetContext(),
		opts,
		input.TypeArgs,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to execute TransferOwnership on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for BurnMintTokenPool", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipBurnMintTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   TransferOwnershipBurnMintTokenPoolObjects{},
	}, nil
}

var TransferOwnershipBurnMintTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-burn-mint-token-pool-transfer-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the BurnMintTokenPool",
	transferOwnershipBurnMintTokenPoolHandler,
)

type AcceptOwnershipBurnMintTokenPoolInput struct {
	BurnMintTokenPoolPackageId string
	TypeArgs                   []string
	StateObjectId              string
}

type AcceptOwnershipBurnMintTokenPoolObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipBurnMintTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipBurnMintTokenPoolInput) (output sui_ops.OpTxResult[AcceptOwnershipBurnMintTokenPoolObjects], err error) {
	burnMintTokenPoolPackage, err := module_burnminttokenpool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipBurnMintTokenPoolObjects]{}, err
	}

	encodedCall, err := burnMintTokenPoolPackage.Encoder().AcceptOwnership(input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on BurnMintTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[AcceptOwnershipBurnMintTokenPoolObjects]{
			Digest:    "",
			PackageId: input.BurnMintTokenPoolPackageId,
			Objects:   AcceptOwnershipBurnMintTokenPoolObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := burnMintTokenPoolPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership accepted for BurnMintTokenPool")

	return sui_ops.OpTxResult[AcceptOwnershipBurnMintTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   AcceptOwnershipBurnMintTokenPoolObjects{},
		Call:      call,
	}, nil
}

var AcceptOwnershipBurnMintTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-burn-mint-token-pool-accept-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the BurnMintTokenPool",
	acceptOwnershipBurnMintTokenPoolHandler,
)
