package opregistry

import (
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
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

	// TP Operations
	operations = append(operations, lockreleasetokenpoolops.AllOperationsLockReleaseTP...)
	operations = append(operations, burnminttokenpoolops.AllOperationsBurnMintTP...)
	// Add more operation slices here as needed:
	// operations = append(operations, anotherops.AllOperations...)

	return operations
}()
