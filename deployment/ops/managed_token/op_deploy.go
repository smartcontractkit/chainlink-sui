package managedtokenops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managedtoken "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/managed_token/managed_token"
	managedtoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/managed_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type ManagedTokenDeployInput struct {
	MCMSAddress      string
	MCMSOwnerAddress string
}

type ManagedTokenDeployOutput struct {
	PublisherObjectId  string
	UpgradeCapObjectId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenDeployInput) (output sui_ops.OpTxResult[ManagedTokenDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	managedTokenPackage, tx, err := managedtoken.PublishCCIPManagedToken(
		b.GetContext(),
		opts,
		deps.Client,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenDeployOutput]{}, err
	}

	publisherObj, err := bind.FindObjectIdFromPublishTx(*tx, "package", "Publisher")
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenDeployOutput]{}, fmt.Errorf("failed to find Publisher object ID: %w", err)
	}

	upgradeCapObj, err := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenDeployOutput]{}, fmt.Errorf("failed to find UpgradeCap object ID: %w", err)
	}

	return sui_ops.OpTxResult[ManagedTokenDeployOutput]{
		Digest:    tx.Digest,
		PackageId: managedTokenPackage.Address(),
		Objects: ManagedTokenDeployOutput{
			PublisherObjectId:  publisherObj,
			UpgradeCapObjectId: upgradeCapObj,
		},
	}, err
}

var DeployCCIPManagedTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP managed token package",
	deployHandler,
)

type TransferOwnershipManagedTokenInput struct {
	ManagedTokenPackageId string
	TypeArgs              []string
	StateObjectId         string
	OwnerCapObjectId      string
	To                    string
}

type TransferOwnershipManagedTokenObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipManagedTokenHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipManagedTokenInput) (output sui_ops.OpTxResult[TransferOwnershipManagedTokenObjects], err error) {
	managedTokenPackage, err := module_managedtoken.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipManagedTokenObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := managedTokenPackage.TransferOwnership(
		b.GetContext(),
		opts,
		input.TypeArgs,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipManagedTokenObjects]{}, fmt.Errorf("failed to execute TransferOwnership on ManagedToken: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for ManagedToken", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipManagedTokenObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPackageId,
		Objects:   TransferOwnershipManagedTokenObjects{},
	}, nil
}

var TransferOwnershipManagedTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-transfer-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the ManagedToken",
	transferOwnershipManagedTokenHandler,
)

type AcceptOwnershipManagedTokenInput struct {
	ManagedTokenPackageId string
	TypeArgs              []string
	StateObjectId         string
}

type AcceptOwnershipManagedTokenObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipManagedTokenHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipManagedTokenInput) (output sui_ops.OpTxResult[AcceptOwnershipManagedTokenObjects], err error) {
	managedTokenPackage, err := module_managedtoken.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenObjects]{}, err
	}

	encodedCall, err := managedTokenPackage.Encoder().AcceptOwnership(input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on ManagedToken as per no Signer provided")
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenObjects]{
			Digest:    "",
			PackageId: input.ManagedTokenPackageId,
			Objects:   AcceptOwnershipManagedTokenObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := managedTokenPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipManagedTokenObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on ManagedToken: %w", err)
	}

	b.Logger.Infow("Ownership accepted for ManagedToken")

	return sui_ops.OpTxResult[AcceptOwnershipManagedTokenObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPackageId,
		Objects:   AcceptOwnershipManagedTokenObjects{},
		Call:      call,
	}, nil
}

var AcceptOwnershipManagedTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-accept-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the ManagedToken",
	acceptOwnershipManagedTokenHandler,
)
