package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_malicious_receiver "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_malicious_receiver/ccip_malicious_receiver"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type RegisterMaliciousReceiverInput struct {
	OwnerCapObjectId           string
	CCIPObjectRefObjectId      string
	MaliciousReceiverPackageId string
}

type RegisterMaliciousReceiverObjects struct {
	// No specific objects are returned from registration.
}

var registerMaliciousReceiverHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RegisterMaliciousReceiverInput) (output sui_ops.OpTxResult[RegisterMaliciousReceiverObjects], err error) {
	// Create a malicious receiver contract instance using the generated binding.
	contract, err := module_malicious_receiver.NewMaliciousReceiver(input.MaliciousReceiverPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RegisterMaliciousReceiverObjects]{}, fmt.Errorf("failed to create malicious receiver contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	b.Logger.Debugw("Registering malicious receiver", "input", input)

	// Call register_receiver using the generated binding.
	tx, err := contract.RegisterReceiver(
		b.GetContext(),
		opts,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[RegisterMaliciousReceiverObjects]{}, fmt.Errorf("failed to execute malicious receiver registration: %w", err)
	}

	b.Logger.Infow("Malicious receiver registered",
		"maliciousReceiverPackageId", input.MaliciousReceiverPackageId,
	)

	return sui_ops.OpTxResult[RegisterMaliciousReceiverObjects]{
		Digest:    tx.Digest,
		PackageId: input.MaliciousReceiverPackageId,
		Objects:   RegisterMaliciousReceiverObjects{},
	}, nil
}

var RegisterMaliciousReceiverOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-malicious-receiver", "malicious_receiver", "register_receiver"),
	semver.MustParse("0.1.0"),
	"Registers the TEST-only CCIP malicious receiver in the OffRamp receiver registry",
	registerMaliciousReceiverHandler,
)
