package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type InitTARObjects struct {
	TARStateObjectID string
}

type InitTARInput struct {
	CCIPPackageID      string
	StateObjectID      string
	OwnerCapObjectID   string
	LocalChainSelector uint64
}

var initTarHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitTARInput) (output sui_ops.OpTxResult[InitTARObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitTARObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Initialize(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.OwnerCapObjectID},
	)
	if err != nil {
		return sui_ops.OpTxResult[InitTARObjects]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "token_admin_registry", "TokenAdminRegistryState")
	if err1 != nil {
		return sui_ops.OpTxResult[InitTARObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[InitTARObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageID,
		Objects: InitTARObjects{
			TARStateObjectID: obj1,
		},
	}, err
}

var TokenAdminRegistryInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Token Admin Registry contract",
	initTarHandler,
)
