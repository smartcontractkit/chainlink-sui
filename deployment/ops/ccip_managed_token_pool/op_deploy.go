package managedtokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managedtokenpool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	managedtokenpool "github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip_token_pools/managed_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type ManagedTokenPoolDeployInput struct {
	CCIPPackageId         string
	ManagedTokenPackageId string
	MCMSAddress           string
	MCMSOwnerAddress      string
}

type ManagedTokenPoolDeployOutput struct {
	OwnerCapObjectId   string
	UpgradeCapObjectId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolDeployInput) (output sui_ops.OpTxResult[ManagedTokenPoolDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tokenPoolPackage, tx, err := managedtokenpool.PublishCCIPManagedTokenPool(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.ManagedTokenPackageId,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolDeployOutput]{}, err
	}

	ownerCapObj, err := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolDeployOutput]{}, fmt.Errorf("failed to find OwnerCap object ID: %w", err)
	}

	upgradeCap, err := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolDeployOutput]{}, fmt.Errorf("failed to find UpgradeCap object ID: %w", err)
	}

	return sui_ops.OpTxResult[ManagedTokenPoolDeployOutput]{
		Digest:    tx.Digest,
		PackageId: tokenPoolPackage.Address(),
		Objects: ManagedTokenPoolDeployOutput{
			OwnerCapObjectId:   ownerCapObj,
			UpgradeCapObjectId: upgradeCap,
		},
	}, err
}

var DeployCCIPManagedTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP managed token pool package",
	deployHandler,
)

type TransferOwnershipManagedTokenPoolInput struct {
	ManagedTokenPoolPackageId string
	TypeArgs                  []string
	StateObjectId             string
	OwnerCapObjectId          string
	To                        string
}

type TransferOwnershipManagedTokenPoolObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipManagedTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipManagedTokenPoolInput) (output sui_ops.OpTxResult[TransferOwnershipManagedTokenPoolObjects], err error) {
	managedTokenPoolPackage, err := module_managedtokenpool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipManagedTokenPoolObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := managedTokenPoolPackage.TransferOwnership(
		b.GetContext(),
		opts,
		input.TypeArgs,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipManagedTokenPoolObjects]{}, fmt.Errorf("failed to execute TransferOwnership on ManagedTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for ManagedTokenPool", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipManagedTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   TransferOwnershipManagedTokenPoolObjects{},
	}, nil
}

var TransferOwnershipManagedTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-pool-transfer-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the ManagedTokenPool",
	transferOwnershipManagedTokenPoolHandler,
)

type AcceptOwnershipManagedTokenPoolInput struct {
	ManagedTokenPoolPackageId string
	TypeArgs                  []string
	StateObjectId             string
}

type AcceptOwnershipManagedTokenPoolObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipManagedTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipManagedTokenPoolInput) (output sui_ops.OpTxResult[AcceptOwnershipManagedTokenPoolObjects], err error) {
	managedTokenPoolPackage, err := module_managedtokenpool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenPoolObjects]{}, err
	}

	encodedCall, err := managedTokenPoolPackage.Encoder().AcceptOwnership(input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenPoolObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on ManagedTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenPoolObjects]{
			Digest:    "",
			PackageId: input.ManagedTokenPoolPackageId,
			Objects:   AcceptOwnershipManagedTokenPoolObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := managedTokenPoolPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenPoolObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on ManagedTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership accepted for ManagedTokenPool")

	return sui_ops.OpTxResult[AcceptOwnershipManagedTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   AcceptOwnershipManagedTokenPoolObjects{},
		Call:      call,
	}, nil
}

var AcceptOwnershipManagedTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-pool-accept-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the ManagedTokenPool",
	acceptOwnershipManagedTokenPoolHandler,
)
