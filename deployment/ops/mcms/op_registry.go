package mcmsops

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mcms_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_registry"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// AddModulesMCMSInput encodes mcms_registry::add_allowed_modules<T> with
// T = <CCIPPackageId>::state_object::McmsCallback (original CCIP package id).
type AddModulesMCMSInput struct {
	MCMSPackageId     string
	MCMSRegistryObjId string
	CCIPPackageId     string
	AllowedModules    []string
}

func ccipStateObjectMcmsCallbackTypeArg(ccipPackageID string) string {
	id := strings.TrimPrefix(strings.TrimSpace(ccipPackageID), "0x")
	return "0x" + id + "::state_object::McmsCallback"
}

var addModulesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AddModulesMCMSInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	if strings.TrimSpace(input.CCIPPackageId) == "" {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("CCIPPackageId is required (original CCIP package for state_object::McmsCallback type arg)")
	}
	if len(input.AllowedModules) == 0 {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("AllowedModules must not be empty")
	}

	contract, err := module_mcms_registry.NewMcmsRegistry(input.MCMSPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to create McmsRegistry contract: %w", err)
	}

	typeArg := ccipStateObjectMcmsCallbackTypeArg(input.CCIPPackageId)
	newMods := make([][]byte, len(input.AllowedModules))
	for i, m := range input.AllowedModules {
		newMods[i] = []byte(m)
	}

	// Generated AddAllowedModules types witness T as bind.Object; McmsCallback is an empty drop struct
	// and must be Pure (empty BCS). WithArgs + EmptyMoveStructWitness matches the Move ABI.
	encodedCall, err := contract.Encoder().AddAllowedModulesWithArgs(
		[]string{typeArg},
		bind.Object{Id: input.MCMSRegistryObjId},
		bind.EmptyMoveStructWitness{},
		newMods,
	)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to encode add_allowed_modules: %w", err)
	}

	// Build the TransactionCall targeting the CCIP package's state_object::mcms_add_allowed_modules
	// wrapper instead of the MCMS package's mcms_registry::add_allowed_modules. This routes execution
	// through processDestinationContractCall → EncodeEntryPointArg (same path as add_package_id),
	// avoiding the processMCMSMainModuleCall path that requires mcms_dispatch_to_registry support.
	// The BCS data is compatible: (registry_addr [32 bytes] || vector<vector<u8>>).
	call, err := sui_ops.ToTransactionCall(encodedCall, input.MCMSRegistryObjId)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	call.PackageID = input.CCIPPackageId
	call.Module = "state_object"
	call.Function = "add_allowed_modules"

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of add_allowed_modules (no signer); returning call for proposal",
			"registry", input.MCMSRegistryObjId, "modules", input.AllowedModules)
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(b.GetContext(), opts, encodedCall)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to execute add_allowed_modules: %w", err)
	}

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}

var AddModulesMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "package", "add_modules"),
	semver.MustParse("0.1.0"),
	"Add module names to the MCMS registry allowlist for the CCIP package (mcms_registry::add_allowed_modules<state_object::McmsCallback>)",
	addModulesHandler,
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsMCMS = []any{
	*MCMSAcceptOwnershipOp,
	*MCMSTransferOwnershipOp,
	*MCMSExecuteTransferOwnershipOp,
	*SetConfigMCMSOp,
	*AddModulesMCMSOp,
	*UpdateMinDelayMCMSOp,
}
