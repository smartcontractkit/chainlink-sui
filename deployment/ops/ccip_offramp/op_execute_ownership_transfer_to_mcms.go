package offrampops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsOffRampInput struct {
	OffRampPackageId     string
	OffRampRefObjectId   string
	OwnerCapObjectId     string
	OffRampStateObjectId string
	RegistryObjectId     string
	To                   string
}

type ExecuteOwnershipTransferToMcmsOffRampObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsOffRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsOffRampInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects], err error) {
	contract, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{}, fmt.Errorf("failed to create OffRamp contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		bind.Object{Id: input.OffRampRefObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.OffRampStateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.OffRampStateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on OffRamp as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{
			Digest:    "",
			PackageId: input.OffRampPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsOffRampObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on OffRamp: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, bind.Object{Id: input.OffRampStateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{}, fmt.Errorf("failed to get new owner for OffRamp: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for OffRamp: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for OffRamp", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsOffRampObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "offramp", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP OffRamp",
	executeOwnershipTransferToMcmsOffRampHandler,
)
