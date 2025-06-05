package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_nonce_manager "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/nonce_manager"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type InitNMObjects struct {
	NonceManagerStateObjectId string
	NonceManagerCapObjectId   string
}

type InitNMInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
}

var initNMHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitNMInput) (output sui_ops.OpTxResult[InitNMObjects], err error) {
	contract, err := module_nonce_manager.NewNonceManager(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitNMObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	method := contract.Initialize(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitNMObjects]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "nonce_manager", "NonceManagerState")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "nonce_manager", "NonceManagerCap")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[InitNMObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[InitNMObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects: InitNMObjects{
			NonceManagerStateObjectId: obj1,
			NonceManagerCapObjectId:   obj2,
		},
	}, err
}

var NonceManagerInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "nonce_manager", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Nonce Manager contract",
	initNMHandler,
)
