package mcmsops

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsMCMS = []cld_ops.Operation[any, any, any]{
	*MCMSAcceptOwnershipOp.AsUntyped(),
	*MCMSTransferOwnershipOp.AsUntyped(),
	*MCMSExecuteTransferOwnershipOp.AsUntyped(),
	*SetConfigMCMSOp.AsUntyped(),
}
