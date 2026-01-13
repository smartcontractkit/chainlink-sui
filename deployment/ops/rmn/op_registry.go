package rmn

import cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

var AllOperationsRMN = []cld_ops.Operation[any, any, any]{
	*CurseChainOp.AsUntyped(),
	*UncurseChainOp.AsUntyped(),
}
