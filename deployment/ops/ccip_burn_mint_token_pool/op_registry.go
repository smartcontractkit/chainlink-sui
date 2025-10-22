package burnminttokenpoolops

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsBurnMintTP = []cld_ops.Operation[any, any, any]{
	// Deployment Operations
	*DeployCCIPBurnMintTokenPoolOp.AsUntyped(),
	*TransferOwnershipBurnMintTokenPoolOp.AsUntyped(),
	*AcceptOwnershipBurnMintTokenPoolOp.AsUntyped(),
	// Token Pool Operations
	*BurnMintTokenPoolInitializeOp.AsUntyped(),
	*BurnMintTokenPoolInitializeByCcipAdminOp.AsUntyped(),
	*BurnMintTokenPoolApplyChainUpdatesOp.AsUntyped(),
	*BurnMintTokenPoolSetChainRateLimiterOp.AsUntyped(),
	*BurnMintTokenPoolAddRemotePoolOp.AsUntyped(),
	*BurnMintTokenPoolSetPoolOp.AsUntyped(),
	*BurnMintTokenPoolSetAllowlistEnabledOp.AsUntyped(),
	*BurnMintTokenPoolApplyAllowlistUpdatesOp.AsUntyped(),
	*BurnMintTokenPoolRemoveRemotePoolOp.AsUntyped(),
	// MCMS Operations
	*ExecuteOwnershipTransferToMcmsBurnMintTokenPoolOp.AsUntyped(),
}
