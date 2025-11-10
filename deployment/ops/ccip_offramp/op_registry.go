package offrampops

import cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

var AllOperationsOfframp = []cld_ops.Operation[any, any, any]{
	*TransferOwnershipOffRampOp.AsUntyped(),
	*AcceptOwnershipOffRampOp.AsUntyped(),
	*ExecuteOwnershipTransferToMcmsOffRampOp.AsUntyped(),
	*ApplySourceChainConfigUpdatesOp.AsUntyped(),
	*SetOCR3ConfigOp.AsUntyped(),
}
