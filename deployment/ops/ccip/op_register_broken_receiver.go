package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type RegisterBrokenReceiverInput struct {
	OwnerCapObjectId        string
	CCIPObjectRefObjectId   string
	BrokenReceiverPackageId string
}

type RegisterBrokenReceiverObjects struct {
	// No specific objects are returned from registration
}

// registerBrokenReceiverHandler registers the broken receiver with the OffRamp
// receiver registry. There is no generated Go binding for the broken receiver
// (it is a test-only contract), so we build the register_receiver call directly
// against a bound contract.
var registerBrokenReceiverHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RegisterBrokenReceiverInput) (output sui_ops.OpTxResult[RegisterBrokenReceiverObjects], err error) {
	contract, err := bind.NewBoundContract(input.BrokenReceiverPackageId, "ccip_broken_receiver", "broken_receiver", deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RegisterBrokenReceiverObjects]{}, fmt.Errorf("failed to create broken receiver bound contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	b.Logger.Debugw("Registering broken receiver", "input", input)

	// broken_receiver::register_receiver(owner_cap: &OwnerCap, ref: &mut CCIPObjectRef)
	encoded, err := contract.EncodeCallArgsWithGenerics(
		"register_receiver",
		[]string{}, // typeArgs
		[]string{}, // typeParams
		[]string{"&OwnerCap", "&mut CCIPObjectRef"},
		[]any{
			bind.Object{Id: input.OwnerCapObjectId},
			bind.Object{Id: input.CCIPObjectRefObjectId},
		},
		nil,
	)
	if err != nil {
		return sui_ops.OpTxResult[RegisterBrokenReceiverObjects]{}, fmt.Errorf("failed to encode broken receiver register_receiver call: %w", err)
	}

	tx, err := contract.ExecuteTransaction(b.GetContext(), opts, encoded)
	if err != nil {
		return sui_ops.OpTxResult[RegisterBrokenReceiverObjects]{}, fmt.Errorf("failed to execute broken receiver registration: %w", err)
	}

	b.Logger.Infow("Broken receiver registered",
		"brokenReceiverPackageId", input.BrokenReceiverPackageId,
	)

	return sui_ops.OpTxResult[RegisterBrokenReceiverObjects]{
		Digest:    tx.Digest,
		PackageId: input.BrokenReceiverPackageId,
		Objects:   RegisterBrokenReceiverObjects{},
	}, nil
}

var RegisterBrokenReceiverOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-broken-receiver", "broken_receiver", "register_receiver"),
	semver.MustParse("0.1.0"),
	"Registers the CCIP broken receiver with the receiver registry (testing only)",
	registerBrokenReceiverHandler,
)
