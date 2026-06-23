package managedtokenpoolops

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsManagedTP = []any{
	*AcceptOwnershipManagedTokenPoolOp,
	*ExecuteOwnershipTransferToMcmsManagedTokenPoolOp,
	*DeployCCIPManagedTokenPoolOp,
	*ManagedTokenPoolInitializeOp,
	*ManagedTokenPoolAddRemotePoolOp,
	*ManagedTokenPoolApplyChainUpdatesOp,
	*ManagedTokenPoolAddRemotePoolOp,
	*ManagedTokenPoolRemoveRemotePoolOp,
	*ManagedTokenPoolSetChainRateLimiterOp,
	*ManagedTokenPoolSetAllowlistEnabledOp,
	*ManagedTokenPoolApplyAllowlistUpdatesOp,
}
