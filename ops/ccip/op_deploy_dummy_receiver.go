package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployDummyReceiverObjects struct {
	OwnerCapObjectID          string
	CCIPReceiverStateObjectID string
}

type DeployDummyReceiverInput struct {
	CCIPPackageID string
	McmsPackageID string
	McmsOwner     string
}

var deployDummyReceiverHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployDummyReceiverInput) (output sui_ops.OpTxResult[DeployDummyReceiverObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Compile the dummy receiver package
	artifact, err := bind.CompilePackage(contracts.CCIPDummyReceiver, map[string]string{
		"ccip":                input.CCIPPackageID,
		"ccip_dummy_receiver": "0x0",
		"mcms":                input.McmsPackageID,
		"mcms_owner":          input.McmsOwner,
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployDummyReceiverObjects]{}, fmt.Errorf("failed to compile dummy receiver package: %w", err)
	}

	// Publish the package
	packageID, tx, err := bind.PublishPackage(
		b.GetContext(),
		opts,
		deps.Client,
		bind.PublishRequest{
			CompiledModules: artifact.Modules,
			Dependencies:    artifact.Dependencies,
		},
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployDummyReceiverObjects]{}, fmt.Errorf("failed to publish dummy receiver package: %w", err)
	}

	// Extract object IDs from the publish transaction
	// The init function creates:
	// 1. CCIPReceiverState (shared object)
	// 2. OwnerCap (transferred to sender)
	ownerCapObjectID, err1 := bind.FindObjectIdFromPublishTx(*tx, "dummy_receiver", "OwnerCap")
	if err1 != nil {
		return sui_ops.OpTxResult[DeployDummyReceiverObjects]{}, fmt.Errorf("failed to find OwnerCap object ID in publish tx: %w", err1)
	}

	receiverStateObjectID, err2 := bind.FindObjectIdFromPublishTx(*tx, "dummy_receiver", "CCIPReceiverState")
	if err2 != nil {
		return sui_ops.OpTxResult[DeployDummyReceiverObjects]{}, fmt.Errorf("failed to find CCIPReceiverState object ID in publish tx: %w", err2)
	}

	return sui_ops.OpTxResult[DeployDummyReceiverObjects]{
		Digest:    tx.Digest,
		PackageID: packageID,
		Objects: DeployDummyReceiverObjects{
			OwnerCapObjectID:          ownerCapObjectID,
			CCIPReceiverStateObjectID: receiverStateObjectID,
		},
	}, nil
}

var DeployCCIPDummyReceiverOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-dummy-receiver", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP dummy receiver package",
	deployDummyReceiverHandler,
)
