package managedtokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsManagedTokenPoolInput struct {
	ManagedTokenPoolPackageId string
	TypeArgs                  []string
	OwnerCapObjectId          string
	StateObjectId             string
	RegistryObjectId          string
	To                        string
}

type ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsManagedTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsManagedTokenPoolInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{}, fmt.Errorf("failed to create ManagedTokenPool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		input.TypeArgs,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on ManagedTokenPool as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{
			Digest:    "",
			PackageId: input.ManagedTokenPoolPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on ManagedTokenPool: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{}, fmt.Errorf("failed to get new owner for ManagedTokenPool: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for ManagedTokenPool: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for ManagedTokenPool", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsManagedTokenPoolObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsManagedTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP ManagedTokenPool",
	executeOwnershipTransferToMcmsManagedTokenPoolHandler,
)
