package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
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

// =================== Get Package IDs Operations =================== //

type GetPackageIdsStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetPackageIdsStateObjectOutput struct {
	PackageIds []string
}

var getPackageIdsStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetPackageIdsStateObjectInput) (output sui_ops.OpTxResult[GetPackageIdsStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetPackageIdsStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	packageIds, err := contract.DevInspect().GetPackageIds(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetPackageIdsStateObjectOutput]{}, fmt.Errorf("failed to get package IDs from StateObject: %w", err)
	}

	b.Logger.Infow("Package IDs retrieved from CCIP StateObject",
		"packageIds", packageIds,
	)

	return sui_ops.OpTxResult[GetPackageIdsStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetPackageIdsStateObjectOutput{
			PackageIds: packageIds,
		},
	}, nil
}

var GetPackageIdsStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_package_ids"),
	semver.MustParse("0.1.0"),
	"Gets all package IDs from the CCIP StateObject",
	getPackageIdsStateObjectHandler,
)

// =================== Get Initial Package ID Operations =================== //

type GetInitialPackageIdStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetInitialPackageIdStateObjectOutput struct {
	InitialPackageId string
}

var getInitialPackageIdStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetInitialPackageIdStateObjectInput) (output sui_ops.OpTxResult[GetInitialPackageIdStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetInitialPackageIdStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	initialPackageId, err := contract.DevInspect().GetInitialPackageId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetInitialPackageIdStateObjectOutput]{}, fmt.Errorf("failed to get initial package ID from StateObject: %w", err)
	}

	b.Logger.Infow("Initial package ID retrieved from CCIP StateObject",
		"initialPackageId", initialPackageId,
	)

	return sui_ops.OpTxResult[GetInitialPackageIdStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetInitialPackageIdStateObjectOutput{
			InitialPackageId: initialPackageId,
		},
	}, nil
}

var GetInitialPackageIdStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_initial_package_id"),
	semver.MustParse("0.1.0"),
	"Gets the initial package ID from the CCIP StateObject",
	getInitialPackageIdStateObjectHandler,
)
