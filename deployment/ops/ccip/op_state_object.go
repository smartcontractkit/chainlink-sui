package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Add Package ID Operations =================== //

type AddPackageIdStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
	OwnerCapObjectId      string
	PackageId             string
}

type AddPackageIdStateObjectObjects struct {
	// No specific objects are returned from add_package_id
}

var addPackageIdStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AddPackageIdStateObjectInput) (output sui_ops.OpTxResult[AddPackageIdStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.AddPackageId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.PackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{}, fmt.Errorf("failed to execute AddPackageId on StateObject: %w", err)
	}

	b.Logger.Infow("Package ID added to CCIP StateObject", "packageId", input.PackageId)

	return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   AddPackageIdStateObjectObjects{},
	}, nil
}

var AddPackageIdStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "add_package_id"),
	semver.MustParse("0.1.0"),
	"Adds a new package ID to the CCIP StateObject for upgrade tracking",
	addPackageIdStateObjectHandler,
)

// =================== Remove Package ID Operations =================== //

type RemovePackageIdStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
	OwnerCapObjectId      string
	PackageId             string
}

type RemovePackageIdStateObjectObjects struct {
	// No specific objects are returned from remove_package_id
}

var removePackageIdStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RemovePackageIdStateObjectInput) (output sui_ops.OpTxResult[RemovePackageIdStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.RemovePackageId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.PackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{}, fmt.Errorf("failed to execute RemovePackageId on StateObject: %w", err)
	}

	b.Logger.Infow("Package ID removed from CCIP StateObject", "packageId", input.PackageId)

	return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   RemovePackageIdStateObjectObjects{},
	}, nil
}

var RemovePackageIdStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "remove_package_id"),
	semver.MustParse("0.1.0"),
	"Removes a package ID from the CCIP StateObject",
	removePackageIdStateObjectHandler,
)

// =================== Get Owner Cap ID Operations =================== //

type GetOwnerCapIdStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetOwnerCapIdStateObjectOutput struct {
	OwnerCapId string
}

var getOwnerCapIdStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetOwnerCapIdStateObjectInput) (output sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	ownerCapId, err := contract.DevInspect().OwnerCapId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput]{}, fmt.Errorf("failed to get owner cap ID from StateObject: %w", err)
	}

	b.Logger.Infow("Owner cap ID retrieved from CCIP StateObject",
		"ownerCapId", ownerCapId.Id,
	)

	return sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetOwnerCapIdStateObjectOutput{
			OwnerCapId: ownerCapId.Id,
		},
	}, nil
}

var GetOwnerCapIdStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_owner_cap_id"),
	semver.MustParse("0.1.0"),
	"Gets the owner cap ID from the CCIP StateObject",
	getOwnerCapIdStateObjectHandler,
)

// =================== Get Owner Operations =================== //

type GetOwnerStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetOwnerStateObjectOutput struct {
	Owner string
}

var getOwnerStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetOwnerStateObjectInput) (output sui_ops.OpTxResult[GetOwnerStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	owner, err := contract.DevInspect().Owner(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerStateObjectOutput]{}, fmt.Errorf("failed to get owner from StateObject: %w", err)
	}

	b.Logger.Infow("Owner retrieved from CCIP StateObject",
		"owner", owner,
	)

	return sui_ops.OpTxResult[GetOwnerStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetOwnerStateObjectOutput{
			Owner: owner,
		},
	}, nil
}

var GetOwnerStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_owner"),
	semver.MustParse("0.1.0"),
	"Gets the owner from the CCIP StateObject",
	getOwnerStateObjectHandler,
)

// =================== Get Pending Transfer Operations =================== //

type GetPendingTransferStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetPendingTransferStateObjectOutput struct {
	HasPendingTransfer      bool
	PendingTransferFrom     *string
	PendingTransferTo       *string
	PendingTransferAccepted *bool
}

var getPendingTransferStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetPendingTransferStateObjectInput) (output sui_ops.OpTxResult[GetPendingTransferStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	hasPending, err := contract.DevInspect().HasPendingTransfer(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to check pending transfer: %w", err)
	}

	var pendingFrom *string
	var pendingTo *string
	var pendingAccepted *bool

	if hasPending {
		pendingFrom, err = contract.DevInspect().PendingTransferFrom(
			b.GetContext(),
			opts,
			bind.Object{Id: input.CCIPObjectRefObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to get pending transfer from: %w", err)
		}

		pendingTo, err = contract.DevInspect().PendingTransferTo(
			b.GetContext(),
			opts,
			bind.Object{Id: input.CCIPObjectRefObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to get pending transfer to: %w", err)
		}

		pendingAccepted, err = contract.DevInspect().PendingTransferAccepted(
			b.GetContext(),
			opts,
			bind.Object{Id: input.CCIPObjectRefObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to get pending transfer accepted: %w", err)
		}
	}

	b.Logger.Infow("Pending transfer info retrieved from CCIP StateObject",
		"hasPending", hasPending,
		"pendingFrom", pendingFrom,
		"pendingTo", pendingTo,
		"pendingAccepted", pendingAccepted,
	)

	return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetPendingTransferStateObjectOutput{
			HasPendingTransfer:      hasPending,
			PendingTransferFrom:     pendingFrom,
			PendingTransferTo:       pendingTo,
			PendingTransferAccepted: pendingAccepted,
		},
	}, nil
}

var GetPendingTransferStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_pending_transfer"),
	semver.MustParse("0.1.0"),
	"Gets pending transfer information from the CCIP StateObject",
	getPendingTransferStateObjectHandler,
)

// =================== Transfer Ownership Operations =================== //

type TransferOwnershipStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
	OwnerCapObjectId      string
	To                    string
}

type TransferOwnershipStateObjectObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipStateObjectInput) (output sui_ops.OpTxResult[TransferOwnershipStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.TransferOwnership(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{}, fmt.Errorf("failed to execute TransferOwnership on StateObject: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for CCIP StateObject", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   TransferOwnershipStateObjectObjects{},
	}, nil
}

var TransferOwnershipStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "transfer_ownership"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the CCIP StateObject",
	transferOwnershipStateObjectHandler,
)

// =================== Accept Ownership Operations =================== //

type AcceptOwnershipStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type AcceptOwnershipStateObjectObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipStateObjectInput) (output sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.AcceptOwnership(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on StateObject: %w", err)
	}

	b.Logger.Infow("Ownership accepted for CCIP StateObject")

	return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   AcceptOwnershipStateObjectObjects{},
	}, nil
}

var AcceptOwnershipStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "accept_ownership"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the CCIP StateObject",
	acceptOwnershipStateObjectHandler,
)
