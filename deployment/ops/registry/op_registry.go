package opregistry

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperations = func() []cld_ops.Operation[any, any, any] {
	var operations []cld_ops.Operation[any, any, any]

	// Add CCIP operations
	operations = append(operations, ccipops.AllOperationsFeeQuoter...)
	operations = append(operations, ccipops.AllOperationsStateObject...)
	operations = append(operations, ccipops.AllOperationsTokenAdminRegistry...)
	operations = append(operations, ccipops.AllOperationsUpgradeRegistry...)

	// Add more operation slices here as needed:
	// operations = append(operations, anotherops.AllOperations...)

	return operations
}()
