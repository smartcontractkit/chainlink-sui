package routerops

import cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

var AllOperationsRouter = []cld_ops.Operation[any, any, any]{
	*TransferOwnershipOp.AsUntyped(),
	*AcceptOwnershipOp.AsUntyped(),
	*ExecuteOwnershipTransferOp.AsUntyped(),
	*SetOnRampsOp.AsUntyped(),
}
