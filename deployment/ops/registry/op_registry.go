package opregistry

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperations = func() []cld_ops.Operation[any, any, any] {
	var operations []cld_ops.Operation[any, any, any]

	// Add CCIP operations
	operations = append(operations, ccipops.AllOperationsCCIP...)
	operations = append(operations, offrampops.AllOperationsOfframp...)
	operations = append(operations, onrampops.AllOperationsOnramp...)
	operations = append(operations, routerops.AllOperationsRouter...)

	// MCMS Operations
	operations = append(operations, mcmsops.AllOperationsMCMS...)

	operations = append(operations, *mcmsops.AddModulesMCMSOp.AsUntyped())
	// Add more operation slices here as needed:
	// operations = append(operations, anotherops.AllOperations...)

	return operations
}()
