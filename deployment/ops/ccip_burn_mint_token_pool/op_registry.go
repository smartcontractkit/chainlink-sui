package burnminttokenpoolops

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsBurnMintTP = []any{
	// Deployment Operations
	*DeployCCIPBurnMintTokenPoolOp,
	*TransferOwnershipBurnMintTokenPoolOp,
	*AcceptOwnershipBurnMintTokenPoolOp,
	// Token Pool Operations
	*BurnMintTokenPoolInitializeOp,
	*BurnMintTokenPoolApplyChainUpdatesOp,
	*BurnMintTokenPoolSetChainRateLimiterOp,
	*BurnMintTokenPoolAddRemotePoolOp,
	*BurnMintTokenPoolSetAllowlistEnabledOp,
	*BurnMintTokenPoolApplyAllowlistUpdatesOp,
	*BurnMintTokenPoolRemoveRemotePoolOp,
	// MCMS Operations
	*ExecuteOwnershipTransferToMcmsBurnMintTokenPoolOp,
}
