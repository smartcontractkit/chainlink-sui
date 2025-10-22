package onrampops

import cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

var AllOperationsOnramp = []cld_ops.Operation[any, any, any]{
	*TransferOwnershipOnRampOp.AsUntyped(),
	*AcceptOwnershipOnRampOp.AsUntyped(),
	*ExecuteOwnershipTransferToMcmsOnRampOp.AsUntyped(),
}
