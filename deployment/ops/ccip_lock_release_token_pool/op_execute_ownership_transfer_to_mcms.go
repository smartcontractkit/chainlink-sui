package lockreleasetokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolInput struct {
	LockReleaseTokenPoolPackageId string
	TypeArgs                      []string
	OwnerCapObjectId              string
	StateObjectId                 string
	RegistryObjectId              string
	To                            string
}

type ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsLockReleaseTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to create LockReleaseTokenPool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		input.TypeArgs,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on LockReleaseTokenPool as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{
			Digest:    "",
			PackageId: input.LockReleaseTokenPoolPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on LockReleaseTokenPool: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{}, fmt.Errorf("failed to get new owner for LockReleaseTokenPool: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for LockReleaseTokenPool: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for LockReleaseTokenPool", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP LockReleaseTokenPool",
	executeOwnershipTransferToMcmsLockReleaseTokenPoolHandler,
)
