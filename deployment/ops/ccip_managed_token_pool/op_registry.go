package managedtokenpoolops

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsManagedTP = []cld_ops.Operation[any, any, any]{
	*AcceptOwnershipManagedTokenPoolOp.AsUntyped(),
	*ExecuteOwnershipTransferToMcmsManagedTokenPoolOp.AsUntyped(),
	*DeployCCIPManagedTokenPoolOp.AsUntyped(),
	*ManagedTokenPoolInitializeOp.AsUntyped(),
	*ManagedTokenPoolAddRemotePoolOp.AsUntyped(),
	*ManagedTokenPoolApplyChainUpdatesOp.AsUntyped(),
	*ManagedTokenPoolAddRemotePoolOp.AsUntyped(),
	*ManagedTokenPoolRemoveRemotePoolOp.AsUntyped(),
	*ManagedTokenPoolSetChainRateLimiterOp.AsUntyped(),
	*ManagedTokenPoolSetAllowlistEnabledOp.AsUntyped(),
	*ManagedTokenPoolApplyAllowlistUpdatesOp.AsUntyped(),
}
