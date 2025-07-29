package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_dummy_receiver "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_dummy_receiver/ccip_dummy_receiver"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type RegisterDummyReceiverInput struct {
	CCIPObjectRefObjectId      string
	DummyReceiverPackageId     string
	DummyReceiverStateObjectId string
	ReceiverStateParams        []string
}

type RegisterDummyReceiverObjects struct {
	// No specific objects are returned from registration
}

var registerDummyReceiverHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RegisterDummyReceiverInput) (output sui_ops.OpTxResult[RegisterDummyReceiverObjects], err error) {
	// Create a CCIP dummy receiver contract instance using the generated binding
	contract, err := module_dummy_receiver.NewDummyReceiver(input.DummyReceiverPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RegisterDummyReceiverObjects]{}, fmt.Errorf("failed to create dummy receiver contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	b.Logger.Debugw("Registering dummy receiver", "input", input)

	// Call the register_receiver function using the generated binding
	tx, err := contract.RegisterReceiver(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
		input.DummyReceiverStateObjectId,
		input.ReceiverStateParams,
	)
	if err != nil {
		return sui_ops.OpTxResult[RegisterDummyReceiverObjects]{}, fmt.Errorf("failed to execute dummy receiver registration: %w", err)
	}

	b.Logger.Infow("Dummy receiver registered",
		"dummyReceiverPackageId", input.DummyReceiverPackageId,
		"receiverStateId", input.DummyReceiverStateObjectId,
	)

	return sui_ops.OpTxResult[RegisterDummyReceiverObjects]{
		Digest:    tx.Digest,
		PackageId: input.DummyReceiverPackageId,
		Objects:   RegisterDummyReceiverObjects{},
	}, nil
}

var RegisterDummyReceiverOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-dummy-receiver", "dummy_receiver", "register_receiver"),
	semver.MustParse("0.1.0"),
	"Registers the CCIP dummy receiver with the receiver registry using the dummy receiver's register_receiver function",
	registerDummyReceiverHandler,
)
