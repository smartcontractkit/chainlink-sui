package lockreleasetokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lockreleasetokenpool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	lockreleasetokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/lock_release_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type LockReleaseTokenPoolDeployInput struct {
	CCIPPackageId    string
	MCMSAddress      string
	MCMSOwnerAddress string
}

type LockReleaseTokenPoolDeployOutput struct {
	OwnerCapObjectId   string
	UpgradeCapObjectId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolDeployInput) (output sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := lockreleasetokenpool.PublishCCIPLockReleaseTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput]{}, err
	}

	ownerCapObj, err := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput]{}, fmt.Errorf("failed to find OwnerCap object ID: %w", err)
	}

	upgradeCapObj, err := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput]{}, fmt.Errorf("failed to find UpgradeCap object ID: %w", err)
	}

	return sui_ops.OpTxResult[LockReleaseTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageId: tokenPoolPackage.Address(),
		Objects: LockReleaseTokenPoolDeployOutput{
			OwnerCapObjectId:   ownerCapObj,
			UpgradeCapObjectId: upgradeCapObj,
		},
	}, err
}

var DeployCCIPLockReleaseTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-lock-release-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP lock release token pool package",
	deployHandler,
)

type TransferOwnershipLockReleaseTokenPoolInput struct {
	LockReleaseTokenPoolPackageId string
	TypeArgs                      []string
	StateObjectId                 string
	OwnerCapObjectId              string
	To                            string
}

type TransferOwnershipLockReleaseTokenPoolObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipLockReleaseTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipLockReleaseTokenPoolInput) (output sui_ops.OpTxResult[TransferOwnershipLockReleaseTokenPoolObjects], err error) {
	lockReleaseTokenPoolPackage, err := module_lockreleasetokenpool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipLockReleaseTokenPoolObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := lockReleaseTokenPoolPackage.TransferOwnership(
		b.GetContext(),
		opts,
		input.TypeArgs,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to execute TransferOwnership on LockReleaseTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for LockReleaseTokenPool", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipLockReleaseTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   TransferOwnershipLockReleaseTokenPoolObjects{},
	}, nil
}

var TransferOwnershipLockReleaseTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-lock-release-token-pool-transfer-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the LockReleaseTokenPool",
	transferOwnershipLockReleaseTokenPoolHandler,
)

type AcceptOwnershipLockReleaseTokenPoolInput struct {
	LockReleaseTokenPoolPackageId string
	TypeArgs                      []string
	StateObjectId                 string
}

type AcceptOwnershipLockReleaseTokenPoolObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipLockReleaseTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipLockReleaseTokenPoolInput) (output sui_ops.OpTxResult[AcceptOwnershipLockReleaseTokenPoolObjects], err error) {
	lockReleaseTokenPoolPackage, err := module_lockreleasetokenpool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipLockReleaseTokenPoolObjects]{}, err
	}

	encodedCall, err := lockReleaseTokenPoolPackage.Encoder().AcceptOwnership(input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on LockReleaseTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[AcceptOwnershipLockReleaseTokenPoolObjects]{
			Digest:    "",
			PackageId: input.LockReleaseTokenPoolPackageId,
			Objects:   AcceptOwnershipLockReleaseTokenPoolObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := lockReleaseTokenPoolPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on LockReleaseTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership accepted for LockReleaseTokenPool")

	return sui_ops.OpTxResult[AcceptOwnershipLockReleaseTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   AcceptOwnershipLockReleaseTokenPoolObjects{},
		Call:      call,
	}, nil
}

var AcceptOwnershipLockReleaseTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-lock-release-token-pool-accept-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the LockReleaseTokenPool",
	acceptOwnershipLockReleaseTokenPoolHandler,
)
