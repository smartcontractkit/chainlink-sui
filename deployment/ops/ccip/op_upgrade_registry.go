package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_upgrade_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/upgrade_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
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

// =================== Version Blocking Operations =================== //

type BlockVersionInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	Version          uint8
}

type BlockVersionObjects struct {
	// No specific objects are returned from block operations
}

var blockVersionHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BlockVersionInput) (output sui_ops.OpTxResult[BlockVersionObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[BlockVersionObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.BlockVersion(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ModuleName,
		input.Version,
	)
	if err != nil {
		return sui_ops.OpTxResult[BlockVersionObjects]{}, fmt.Errorf("failed to execute BlockVersion: %w", err)
	}

	b.Logger.Infow("Version blocked",
		"moduleName", input.ModuleName,
		"version", input.Version,
	)

	return sui_ops.OpTxResult[BlockVersionObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   BlockVersionObjects{},
	}, nil
}

var BlockVersionOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "block_version"),
	semver.MustParse("0.1.0"),
	"Blocks an entire version of a module in the UpgradeRegistry",
	blockVersionHandler,
)

// =================== Unblock Version Operations =================== //

type UnblockVersionInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	Version          uint8
}

type UnblockVersionObjects struct {
	// No specific objects are returned from unblock operations
}

var unblockVersionHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UnblockVersionInput) (output sui_ops.OpTxResult[UnblockVersionObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[UnblockVersionObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UnblockVersion(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ModuleName,
		input.Version,
	)
	if err != nil {
		return sui_ops.OpTxResult[UnblockVersionObjects]{}, fmt.Errorf("failed to execute UnblockVersion: %w", err)
	}

	b.Logger.Infow("Version unblocked",
		"moduleName", input.ModuleName,
		"version", input.Version,
	)

	return sui_ops.OpTxResult[UnblockVersionObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   UnblockVersionObjects{},
	}, nil
}

var UnblockVersionOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "unblock_version"),
	semver.MustParse("0.1.0"),
	"Unblocks an entire version of a module in the UpgradeRegistry",
	unblockVersionHandler,
)

// =================== Function Blocking Operations =================== //

type BlockFunctionInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	FunctionName     string
	Version          uint8
}

type BlockFunctionObjects struct {
	// No specific objects are returned from block operations
}

var blockFunctionHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BlockFunctionInput) (output sui_ops.OpTxResult[BlockFunctionObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[BlockFunctionObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.BlockFunction(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ModuleName,
		input.FunctionName,
		input.Version,
	)
	if err != nil {
		return sui_ops.OpTxResult[BlockFunctionObjects]{}, fmt.Errorf("failed to execute BlockFunction: %w", err)
	}

	b.Logger.Infow("Function blocked",
		"moduleName", input.ModuleName,
		"functionName", input.FunctionName,
		"version", input.Version,
	)

	return sui_ops.OpTxResult[BlockFunctionObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   BlockFunctionObjects{},
	}, nil
}

var BlockFunctionOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "block_function"),
	semver.MustParse("0.1.0"),
	"Blocks a specific function in a specific version in the UpgradeRegistry",
	blockFunctionHandler,
)

// =================== Unblock Function Operations =================== //

type UnblockFunctionInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
	ModuleName       string
	FunctionName     string
	Version          uint8
}

type UnblockFunctionObjects struct {
	// No specific objects are returned from unblock operations
}

var unblockFunctionHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UnblockFunctionInput) (output sui_ops.OpTxResult[UnblockFunctionObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[UnblockFunctionObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UnblockFunction(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ModuleName,
		input.FunctionName,
		input.Version,
	)
	if err != nil {
		return sui_ops.OpTxResult[UnblockFunctionObjects]{}, fmt.Errorf("failed to execute UnblockFunction: %w", err)
	}

	b.Logger.Infow("Function unblocked",
		"moduleName", input.ModuleName,
		"functionName", input.FunctionName,
		"version", input.Version,
	)

	return sui_ops.OpTxResult[UnblockFunctionObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   UnblockFunctionObjects{},
	}, nil
}

var UnblockFunctionOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "unblock_function"),
	semver.MustParse("0.1.0"),
	"Unblocks a specific function in a specific version in the UpgradeRegistry",
	unblockFunctionHandler,
)

// =================== Module Restrictions Operations =================== //

type GetModuleRestrictionsInput struct {
	CCIPPackageId string
	StateObjectId string
	ModuleName    string
}

type GetModuleRestrictionsOutput struct {
	Restrictions [][]byte
}

var getModuleRestrictionsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetModuleRestrictionsInput) (output sui_ops.OpTxResult[GetModuleRestrictionsOutput], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetModuleRestrictionsOutput]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	restrictions, err := contract.DevInspect().GetModuleRestrictions(
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
		"restrictions", restrictions,
	)

	return sui_ops.OpTxResult[GetModuleRestrictionsOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetModuleRestrictionsOutput{
			Restrictions: restrictions,
		},
	}, nil
}

var GetModuleRestrictionsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "get_module_restrictions"),
	semver.MustParse("0.1.0"),
	"Gets module restrictions from the UpgradeRegistry",
	getModuleRestrictionsHandler,
)

// =================== Function Permission Operations =================== //

type IsFunctionAllowedInput struct {
	CCIPPackageId   string
	StateObjectId   string
	ModuleName      string
	FunctionName    string
	ContractVersion uint8
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
	opts.Signer = deps.Signer
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

// =================== Function Verification Operations =================== //

type VerifyFunctionAllowedInput struct {
	CCIPPackageId   string
	StateObjectId   string
	ModuleName      string
	FunctionName    string
	ContractVersion uint8
}

type VerifyFunctionAllowedObjects struct {
	// No specific objects are returned from verification operations
}

var verifyFunctionAllowedHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input VerifyFunctionAllowedInput) (output sui_ops.OpTxResult[VerifyFunctionAllowedObjects], err error) {
	contract, err := module_upgrade_registry.NewUpgradeRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[VerifyFunctionAllowedObjects]{}, fmt.Errorf("failed to create UpgradeRegistry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.VerifyFunctionAllowed(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.ModuleName,
		input.FunctionName,
		input.ContractVersion,
	)
	if err != nil {
		return sui_ops.OpTxResult[VerifyFunctionAllowedObjects]{}, fmt.Errorf("failed to verify function allowed: %w", err)
	}

	b.Logger.Infow("Function verification completed",
		"moduleName", input.ModuleName,
		"functionName", input.FunctionName,
		"contractVersion", input.ContractVersion,
	)

	return sui_ops.OpTxResult[VerifyFunctionAllowedObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   VerifyFunctionAllowedObjects{},
	}, nil
}

var VerifyFunctionAllowedOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "upgrade_registry", "verify_function_allowed"),
	semver.MustParse("0.1.0"),
	"Verifies that a function is allowed in the UpgradeRegistry (throws error if not allowed)",
	verifyFunctionAllowedHandler,
)
