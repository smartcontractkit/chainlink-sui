package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_upgrade_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/upgrade_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

// =================== Initialize Operations =================== //

type InitUpgradeRegistryObjects struct {
	UpgradeRegistryObjectId string
}

type InitUpgradeRegistryInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
}

var initUpgradeRegistryHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitUpgradeRegistryInput) (output sui_ops.OpTxResult[InitUpgradeRegistryObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitUpgradeRegistryObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Initialize(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[InitUpgradeRegistryObjects]{}, fmt.Errorf("failed to execute UpgradeRegistry initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "upgrade_registry", "UpgradeRegistry")
	if err1 != nil {
		return sui_ops.OpTxResult[InitUpgradeRegistryObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	b.Logger.Infow("UpgradeRegistry initialized", "upgradeRegistryObjectId", obj1)

	return sui_ops.OpTxResult[InitUpgradeRegistryObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects: InitUpgradeRegistryObjects{
			UpgradeRegistryObjectId: obj1,
		},
	}, err
}

var UpgradeRegistryInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP UpgradeRegistry contract",
	initUpgradeRegistryHandler,
)

// =================== Function Restrictions Operations =================== //

type UpdateFunctionRestrictionsInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	FunctionName     string
	BlockedVersions  []uint64
}

type UpdateFunctionRestrictionsObjects struct {
	// No specific objects are returned from update operations
}

var updateFunctionRestrictionsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UpdateFunctionRestrictionsInput) (output sui_ops.OpTxResult[UpdateFunctionRestrictionsObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[UpdateFunctionRestrictionsObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UpdateFunctionRestrictions(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ModuleName,
		input.FunctionName,
		input.BlockedVersions,
	)
	if err != nil {
		return sui_ops.OpTxResult[UpdateFunctionRestrictionsObjects]{}, fmt.Errorf("failed to execute UpdateFunctionRestrictions: %w", err)
	}

	b.Logger.Infow("Function restrictions updated",
		"moduleName", input.ModuleName,
		"functionName", input.FunctionName,
		"blockedVersions", input.BlockedVersions,
	)

	return sui_ops.OpTxResult[UpdateFunctionRestrictionsObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   UpdateFunctionRestrictionsObjects{},
	}, nil
}

var UpdateFunctionRestrictionsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "update_function_restrictions"),
	semver.MustParse("0.1.0"),
	"Updates function restrictions in the UpgradeRegistry",
	updateFunctionRestrictionsHandler,
)

type GetFunctionRestrictionsInput struct {
	CCIPPackageId string
	StateObjectId string
	ModuleName    string
	FunctionName  string
}

type GetFunctionRestrictionsOutput struct {
	BlockedVersions []uint64
}

var getFunctionRestrictionsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetFunctionRestrictionsInput) (output sui_ops.OpTxResult[GetFunctionRestrictionsOutput], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetFunctionRestrictionsOutput]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	blockedVersions, err := contract.DevInspect().GetFunctionRestrictions(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.ModuleName,
		input.FunctionName,
	)
	if err != nil {
		return sui_ops.OpTxResult[GetFunctionRestrictionsOutput]{}, fmt.Errorf("failed to get function restrictions: %w", err)
	}

	b.Logger.Infow("Function restrictions retrieved",
		"moduleName", input.ModuleName,
		"functionName", input.FunctionName,
		"blockedVersions", blockedVersions,
	)

	return sui_ops.OpTxResult[GetFunctionRestrictionsOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetFunctionRestrictionsOutput{
			BlockedVersions: blockedVersions,
		},
	}, nil
}

var GetFunctionRestrictionsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "get_function_restrictions"),
	semver.MustParse("0.1.0"),
	"Gets function restrictions from the UpgradeRegistry",
	getFunctionRestrictionsHandler,
)

type IsFunctionAllowedInput struct {
	CCIPPackageId   string
	StateObjectId   string
	ModuleName      string
	FunctionName    string
	ContractVersion uint64
}

type IsFunctionAllowedOutput struct {
	IsAllowed bool
}

var isFunctionAllowedHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IsFunctionAllowedInput) (output sui_ops.OpTxResult[IsFunctionAllowedOutput], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsFunctionAllowedOutput]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	isAllowed, err := contract.DevInspect().IsFunctionAllowed(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.ModuleName,
		input.FunctionName,
		input.ContractVersion,
	)
	if err != nil {
		return sui_ops.OpTxResult[IsFunctionAllowedOutput]{}, fmt.Errorf("failed to check if function is allowed: %w", err)
	}

	b.Logger.Infow("Function allowed check completed",
		"moduleName", input.ModuleName,
		"functionName", input.FunctionName,
		"contractVersion", input.ContractVersion,
		"isAllowed", isAllowed,
	)

	return sui_ops.OpTxResult[IsFunctionAllowedOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: IsFunctionAllowedOutput{
			IsAllowed: isAllowed,
		},
	}, nil
}

var IsFunctionAllowedOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "is_function_allowed"),
	semver.MustParse("0.1.0"),
	"Checks if a function is allowed in the UpgradeRegistry",
	isFunctionAllowedHandler,
)

// =================== Module Restrictions Operations =================== //

type UpdateModuleRestrictionsInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	BlockedVersions  []uint64
}

type UpdateModuleRestrictionsObjects struct {
	// No specific objects are returned from update operations
}

var updateModuleRestrictionsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UpdateModuleRestrictionsInput) (output sui_ops.OpTxResult[UpdateModuleRestrictionsObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[UpdateModuleRestrictionsObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UpdateModuleRestrictions(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ModuleName,
		input.BlockedVersions,
	)
	if err != nil {
		return sui_ops.OpTxResult[UpdateModuleRestrictionsObjects]{}, fmt.Errorf("failed to execute UpdateModuleRestrictions: %w", err)
	}

	b.Logger.Infow("Module restrictions updated",
		"moduleName", input.ModuleName,
		"blockedVersions", input.BlockedVersions,
	)

	return sui_ops.OpTxResult[UpdateModuleRestrictionsObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   UpdateModuleRestrictionsObjects{},
	}, nil
}

var UpdateModuleRestrictionsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "update_module_restrictions"),
	semver.MustParse("0.1.0"),
	"Updates module restrictions in the UpgradeRegistry",
	updateModuleRestrictionsHandler,
)

type GetModuleRestrictionsInput struct {
	CCIPPackageId string
	StateObjectId string
	ModuleName    string
}

type GetModuleRestrictionsOutput struct {
	BlockedVersions []uint64
}

var getModuleRestrictionsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetModuleRestrictionsInput) (output sui_ops.OpTxResult[GetModuleRestrictionsOutput], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetModuleRestrictionsOutput]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	blockedVersions, err := contract.DevInspect().GetModuleRestrictions(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.ModuleName,
	)
	if err != nil {
		return sui_ops.OpTxResult[GetModuleRestrictionsOutput]{}, fmt.Errorf("failed to get module restrictions: %w", err)
	}

	b.Logger.Infow("Module restrictions retrieved",
		"moduleName", input.ModuleName,
		"blockedVersions", blockedVersions,
	)

	return sui_ops.OpTxResult[GetModuleRestrictionsOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetModuleRestrictionsOutput{
			BlockedVersions: blockedVersions,
		},
	}, nil
}

var GetModuleRestrictionsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "get_module_restrictions"),
	semver.MustParse("0.1.0"),
	"Gets module restrictions from the UpgradeRegistry",
	getModuleRestrictionsHandler,
)

type IsModuleAllowedInput struct {
	CCIPPackageId   string
	StateObjectId   string
	ModuleName      string
	ContractVersion uint64
}

type IsModuleAllowedOutput struct {
	IsAllowed bool
}

var isModuleAllowedHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IsModuleAllowedInput) (output sui_ops.OpTxResult[IsModuleAllowedOutput], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsModuleAllowedOutput]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	isAllowed, err := contract.DevInspect().IsModuleAllowed(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.ModuleName,
		input.ContractVersion,
	)
	if err != nil {
		return sui_ops.OpTxResult[IsModuleAllowedOutput]{}, fmt.Errorf("failed to check if module is allowed: %w", err)
	}

	b.Logger.Infow("Module allowed check completed",
		"moduleName", input.ModuleName,
		"contractVersion", input.ContractVersion,
		"isAllowed", isAllowed,
	)

	return sui_ops.OpTxResult[IsModuleAllowedOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: IsModuleAllowedOutput{
			IsAllowed: isAllowed,
		},
	}, nil
}

var IsModuleAllowedOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "is_module_allowed"),
	semver.MustParse("0.1.0"),
	"Checks if a module is allowed in the UpgradeRegistry",
	isModuleAllowedHandler,
)

// =================== Package History Operations =================== //

type GetPackageHistoryInput struct {
	CCIPPackageId string
	StateObjectId string
	PackageName   string
}

type GetPackageHistoryOutput struct {
	PackageIds []string
	Versions   []uint64
	Timestamps []uint64
}

var getPackageHistoryHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetPackageHistoryInput) (output sui_ops.OpTxResult[GetPackageHistoryOutput], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetPackageHistoryOutput]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	results, err := contract.DevInspect().GetPackageHistory(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.PackageName,
	)
	if err != nil {
		return sui_ops.OpTxResult[GetPackageHistoryOutput]{}, fmt.Errorf("failed to get package history: %w", err)
	}

	// The results should contain 3 vectors: package_ids, versions, timestamps
	if len(results) != 3 {
		return sui_ops.OpTxResult[GetPackageHistoryOutput]{}, fmt.Errorf("expected 3 return values, got %d", len(results))
	}

	packageIds, ok1 := results[0].([]string)
	versions, ok2 := results[1].([]uint64)
	timestamps, ok3 := results[2].([]uint64)

	if !ok1 || !ok2 || !ok3 {
		return sui_ops.OpTxResult[GetPackageHistoryOutput]{}, fmt.Errorf("failed to parse package history results")
	}

	b.Logger.Infow("Package history retrieved",
		"packageName", input.PackageName,
		"packageIds", packageIds,
		"versions", versions,
		"timestamps", timestamps,
	)

	return sui_ops.OpTxResult[GetPackageHistoryOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetPackageHistoryOutput{
			PackageIds: packageIds,
			Versions:   versions,
			Timestamps: timestamps,
		},
	}, nil
}

var GetPackageHistoryOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "get_package_history"),
	semver.MustParse("0.1.0"),
	"Gets package history from the UpgradeRegistry",
	getPackageHistoryHandler,
)
