package usdctokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_usdc_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/usdc_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsUsdcTokenPoolInput struct {
	UsdcTokenPoolPackageId string
	TypeArgs               []string
	OwnerCapObjectId       string
	StateObjectId          string
	RegistryObjectId       string
	To                     string
}

type ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsUsdcTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsUsdcTokenPoolInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects], err error) {
	contract, err := module_usdc_token_pool.NewUsdcTokenPool(input.UsdcTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{}, fmt.Errorf("failed to create UsdcTokenPool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		input.TypeArgs,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on UsdcTokenPool as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{
			Digest:    "",
			PackageId: input.UsdcTokenPoolPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on UsdcTokenPool: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{}, fmt.Errorf("failed to get new owner for UsdcTokenPool: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for UsdcTokenPool: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for UsdcTokenPool", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.UsdcTokenPoolPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsUsdcTokenPoolObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsUsdcTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "usdc_token_pool", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP UsdcTokenPool",
	executeOwnershipTransferToMcmsUsdcTokenPoolHandler,
)
