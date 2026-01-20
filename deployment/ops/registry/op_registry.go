package opregistry

import (
	"fmt"
	"reflect"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	rmnops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
)

var AllOperationsTyped = func() []any {
	operations := []any{}

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
	operations = append(operations, managedtokenpoolops.AllOperationsManagedTP...)

	// RMN Operations
	operations = append(operations, rmnops.AllOperationsRMN...)

	return operations
}()

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperations = func() []*cld_ops.Operation[any, any, any] {
	typedOps := AllOperationsTyped

	operations := make([]*cld_ops.Operation[any, any, any], len(typedOps))
	for i, op := range typedOps {
		// Use reflection to call AsUntyped method
		// The operations are stored as values, so we need to get their address to call pointer methods
		opVal := reflect.ValueOf(op)

		// If AsUntyped is a pointer receiver, we need to get the address
		if opVal.CanAddr() {
			opVal = opVal.Addr()
		} else {
			// If we can't get the address directly (e.g., because it's stored as an interface value),
			// we need to create a new addressable value
			opPtr := reflect.New(opVal.Type())
			opPtr.Elem().Set(opVal)
			opVal = opPtr
		}

		asUntypedMethod := opVal.MethodByName("AsUntypedRelaxed")
		if !asUntypedMethod.IsValid() {
			panic(fmt.Sprintf("operation %v does not have AsUntypedRelaxed method", opVal.Type()))
		}
		results := asUntypedMethod.Call(nil)
		if len(results) != 1 {
			panic("AsUntypedRelaxed should return exactly one value")
		}
		// The result is a pointer to Operation[any, any, any]
		untypedOp := results[0].Interface().(*cld_ops.Operation[any, any, any])
		operations[i] = untypedOp
	}

	// add safeguard to ensure no nil/valid operations
	for i, op := range operations {
		if op == nil {
			panic(fmt.Sprintf("operation at index %d is nil", i))
		}
		// try to access Def to ensure it's valid
		_ = op.Def()
	}

	return operations
}()
