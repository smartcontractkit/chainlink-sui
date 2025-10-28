package burnminttokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsBurnMintTokenPoolInput struct {
	BurnMintTokenPoolPackageId string
	TypeArgs                   []string
	OwnerCapObjectId           string
	StateObjectId              string
	RegistryObjectId           string
	To                         string
}

type ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsBurnMintTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsBurnMintTokenPoolInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to create BurnMintTokenPool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		input.TypeArgs,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, input.TypeArgs)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on BurnMintTokenPool as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{
			Digest:    "",
			PackageId: input.BurnMintTokenPoolPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on BurnMintTokenPool: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, input.TypeArgs, bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{}, fmt.Errorf("failed to get new owner for BurnMintTokenPool: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for BurnMintTokenPool: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for BurnMintTokenPool", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsBurnMintTokenPoolObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsBurnMintTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP BurnMintTokenPool",
	executeOwnershipTransferToMcmsBurnMintTokenPoolHandler,
)
