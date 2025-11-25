package ccipops

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsCCIP = []cld_ops.Operation[any, any, any]{
	// Fee Quoter Operations
	*FeeQuoterInitializeOp.AsUntyped(),
	*FeeQuoterApplyFeeTokenUpdatesOp.AsUntyped(),
	*FeeQuoterApplyTokenTransferFeeConfigUpdatesOp.AsUntyped(),
	*FeeQuoterApplyDestChainConfigUpdatesOp.AsUntyped(),
	*FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp.AsUntyped(),
	*FeeQuoterUpdateTokenPricesOp.AsUntyped(),
	*FeeQuoterNewFeeQuoterCapOp.AsUntyped(),
	*FeeQuoterDestroyFeeQuoterCapOp.AsUntyped(),
	*FeeQuoterUpdatePricesWithOwnerCapOp.AsUntyped(),
	// State Object Operations
	*AddPackageIdStateObjectOp.AsUntyped(),
	*RemovePackageIdStateObjectOp.AsUntyped(),
	*GetOwnerCapIdStateObjectOp.AsUntyped(),
	*GetOwnerStateObjectOp.AsUntyped(),
	*GetPendingTransferStateObjectOp.AsUntyped(),
	*TransferOwnershipStateObjectOp.AsUntyped(),
	*AcceptOwnershipStateObjectOp.AsUntyped(),
	*ExecuteOwnershipTransferToMcmsStateObjectOp.AsUntyped(),
	// Token Admin Registry Operations
	*TokenAdminRegistryInitializeOp.AsUntyped(),
	*TokenAdminRegistryUnregisterPoolOp.AsUntyped(),
	*TokenAdminRegistryTransferAdminRoleOp.AsUntyped(),
	*TokenAdminRegistryAcceptAdminRoleOp.AsUntyped(),
	// Upgrade Registry Operations
	*UpgradeRegistryInitializeOp.AsUntyped(),
	*BlockVersionOp.AsUntyped(),
	*BlockFunctionOp.AsUntyped(),
	*GetModuleRestrictionsOp.AsUntyped(),
	*IsFunctionAllowedOp.AsUntyped(),
	*VerifyFunctionAllowedOp.AsUntyped(),
	// RMN Remote Operations
	*RMNRemoteInitializeOp.AsUntyped(),
	*RMNRemoteSetConfigOp.AsUntyped(),
	*RMNRemoteCurseOp.AsUntyped(),
	*RMNRemoteCurseMultipleOp.AsUntyped(),
	*RMNRemoteUncurseOp.AsUntyped(),
	*RMNRemoteUncurseMultipleOp.AsUntyped(),
}
