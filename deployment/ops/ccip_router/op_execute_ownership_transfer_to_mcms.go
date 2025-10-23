package routerops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsRouterInput struct {
	RouterPackageId     string
	OwnerCapObjectId    string
	RouterStateObjectId string
	RegistryObjectId    string
	To                  string
}

type ExecuteOwnershipTransferToMcmsRouterObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsRouterHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsRouterInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects], err error) {
	contract, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{}, fmt.Errorf("failed to create Router contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.RouterStateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.RouterStateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on Router as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{
			Digest:    "",
			PackageId: input.RouterPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsRouterObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on Router: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, bind.Object{Id: input.RouterStateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{}, fmt.Errorf("failed to get new owner for Router: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for Router: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for Router", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsRouterObjects]{
		Digest:    tx.Digest,
		PackageId: input.RouterPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsRouterObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsRouterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "router", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP Router",
	executeOwnershipTransferToMcmsRouterHandler,
)
