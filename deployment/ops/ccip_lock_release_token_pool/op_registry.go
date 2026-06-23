package lockreleasetokenpoolops

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsLockReleaseTP = []any{
	// Deployment Operations
	*DeployCCIPLockReleaseTokenPoolOp,
	*TransferOwnershipLockReleaseTokenPoolOp,
	*AcceptOwnershipLockReleaseTokenPoolOp,
	// Token Pool Operations
	*LockReleaseTokenPoolInitializeOp,
	*LockReleaseTokenPoolApplyChainUpdatesOp,
	*LockReleaseTokenPoolSetChainRateLimiterOp,
	*LockReleaseTokenPoolProvideLiquidityOp,
	*LockReleaseTokenPoolAddRemotePoolOp,
	*LockReleaseTokenPoolSetAllowlistEnabledOp,
	*LockReleaseTokenPoolApplyAllowlistUpdatesOp,
	*LockReleaseTokenPoolRemoveRemotePoolOp,
	// MCMS Operations
	*ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolOp,
}
