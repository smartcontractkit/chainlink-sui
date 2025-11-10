package managedtokenops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managed_token "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/managed_token/managed_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsManagedTokenInput struct {
	ManagedTokenPackageId string
	TypeArgs              []string
	OwnerCapObjectId      string
	StateObjectId         string
	RegistryObjectId      string
	To                    string
}

type ExecuteOwnershipTransferToMcmsManagedTokenObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsManagedTokenHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsManagedTokenInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects], err error) {
	contract, err := module_managed_token.NewManagedToken(input.ManagedTokenPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{}, fmt.Errorf("failed to create ManagedToken contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		input.TypeArgs,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on ManagedToken as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{
			Digest:    "",
			PackageId: input.ManagedTokenPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsManagedTokenObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on ManagedToken: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{}, fmt.Errorf("failed to get new owner for ManagedToken: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for ManagedToken: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for ManagedToken", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsManagedTokenObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsManagedTokenObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsManagedTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP ManagedToken",
	executeOwnershipTransferToMcmsManagedTokenHandler,
)
