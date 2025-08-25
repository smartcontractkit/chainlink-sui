package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_receiver_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/receiver_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type InitRecRegObjects struct {
	ReceiverRegistryStateObjectID string
}

type InitRecRegInput struct {
	CCIPPackageID    string
	StateObjectID    string
	OwnerCapObjectID string
}

var initRecRegHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitRecRegInput) (output sui_ops.OpTxResult[InitRecRegObjects], err error) {
	contract, err := module_receiver_registry.NewReceiverRegistry(input.CCIPPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitRecRegObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
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
		return sui_ops.OpTxResult[InitRecRegObjects]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "receiver_registry", "ReceiverRegistry")
	if err1 != nil {
		return sui_ops.OpTxResult[InitRecRegObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[InitRecRegObjects]{
		Digest:    tx.Digest,
		PackageID: input.CCIPPackageID,
		Objects: InitRecRegObjects{
			ReceiverRegistryStateObjectID: obj1,
		},
	}, err
}

var ReceiverRegistryInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "receiver_registry", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Receiver Registry contract",
	initRecRegHandler,
)
