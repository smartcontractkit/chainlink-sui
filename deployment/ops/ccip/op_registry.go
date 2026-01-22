package ccipops

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsCCIP = []any{
	// Fee Quoter Operations
	*FeeQuoterInitializeOp,
	*FeeQuoterApplyFeeTokenUpdatesOp,
	*FeeQuoterApplyTokenTransferFeeConfigUpdatesOp,
	*FeeQuoterApplyDestChainConfigUpdatesOp,
	*FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp,
	*FeeQuoterUpdateTokenPricesOp,
	*FeeQuoterNewFeeQuoterCapOp,
	*FeeQuoterDestroyFeeQuoterCapOp,
	*FeeQuoterUpdatePricesWithOwnerCapOp,
	// State Object Operations
	*AddPackageIdStateObjectOp,
	*RemovePackageIdStateObjectOp,
	*GetOwnerCapIdStateObjectOp,
	*GetOwnerStateObjectOp,
	*GetPendingTransferStateObjectOp,
	*TransferOwnershipStateObjectOp,
	*AcceptOwnershipStateObjectOp,
	*ExecuteOwnershipTransferToMcmsStateObjectOp,
	// Token Admin Registry Operations
	*TokenAdminRegistryInitializeOp,
	*TokenAdminRegistryUnregisterPoolOp,
	*TokenAdminRegistryTransferAdminRoleOp,
	*TokenAdminRegistryAcceptAdminRoleOp,
	// Upgrade Registry Operations
	*UpgradeRegistryInitializeOp,
	*BlockVersionOp,
	*BlockFunctionOp,
	*GetModuleRestrictionsOp,
	*IsFunctionAllowedOp,
	*VerifyFunctionAllowedOp,
	// RMN Remote Operations
	*RMNRemoteInitializeOp,
	*RMNRemoteSetConfigOp,
	*RMNRemoteCurseOp,
	*RMNRemoteCurseMultipleOp,
	*RMNRemoteUncurseOp,
	*RMNRemoteUncurseMultipleOp,
}
