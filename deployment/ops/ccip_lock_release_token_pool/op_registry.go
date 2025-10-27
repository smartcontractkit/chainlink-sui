package lockreleasetokenpoolops

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsLockReleaseTP = []cld_ops.Operation[any, any, any]{
	// Deployment Operations
	*DeployCCIPLockReleaseTokenPoolOp.AsUntyped(),
	*TransferOwnershipLockReleaseTokenPoolOp.AsUntyped(),
	*AcceptOwnershipLockReleaseTokenPoolOp.AsUntyped(),
	// Token Pool Operations
	*LockReleaseTokenPoolInitializeOp.AsUntyped(),
	*LockReleaseTokenPoolApplyChainUpdatesOp.AsUntyped(),
	*LockReleaseTokenPoolSetChainRateLimiterOp.AsUntyped(),
	*LockReleaseTokenPoolProviderLiquidityOp.AsUntyped(),
	*LockReleaseTokenPoolAddRemotePoolOp.AsUntyped(),
	*LockReleaseTokenPoolSetPoolOp.AsUntyped(),
	*LockReleaseTokenPoolSetRebalancerOp.AsUntyped(),
	*LockReleaseTokenPoolSetAllowlistEnabledOp.AsUntyped(),
	*LockReleaseTokenPoolApplyAllowlistUpdatesOp.AsUntyped(),
	*LockReleaseTokenPoolRemoveRemotePoolOp.AsUntyped(),
	// MCMS Operations
	*ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolOp.AsUntyped(),
}
