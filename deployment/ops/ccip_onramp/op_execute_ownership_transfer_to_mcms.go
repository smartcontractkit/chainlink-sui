package onrampops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsOnRampInput struct {
	OnRampPackageId     string
	OnRampRefObjectId   string
	OwnerCapObjectId    string
	OnRampStateObjectId string
	RegistryObjectId    string
	To                  string
}

type ExecuteOwnershipTransferToMcmsOnRampObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsOnRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsOnRampInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects], err error) {
	contract, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{}, fmt.Errorf("failed to create OnRamp contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		bind.Object{Id: input.OnRampRefObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.OnRampStateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.OnRampStateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on OnRamp as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsOnRampObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on OnRamp: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, bind.Object{Id: input.OnRampStateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{}, fmt.Errorf("failed to get new owner for OnRamp: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for OnRamp: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for OnRamp", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsOnRampObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "onramp", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP OnRamp",
	executeOwnershipTransferToMcmsOnRampHandler,
)
