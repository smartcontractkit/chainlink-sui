package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// DeployMaliciousReceiverObjects are the objects created by publishing the malicious receiver.
type DeployMaliciousReceiverObjects struct {
	OwnerCapObjectId          string
	CCIPReceiverStateObjectId string
}

type DeployMaliciousReceiverInput struct {
	CCIPPackageId     string
	McmsPackageId     string
	FastMcmsPackageId string
	McmsOwner         string
}

var deployMaliciousReceiverHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployMaliciousReceiverInput) (output sui_ops.OpTxResult[DeployMaliciousReceiverObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousReceiverObjects]{}, err
	}

	// Compile the malicious receiver package (TEST-only; used by the guard E2E smoke test).
	artifact, err := bind.CompilePackage(contracts.CCIPMaliciousReceiver, map[string]string{
		"ccip":                    input.CCIPPackageId,
		"ccip_malicious_receiver": "0x0",
		"mcms":                    input.McmsPackageId,
		"fast_mcms":               input.FastMcmsPackageId,
		"mcms_owner":              input.McmsOwner,
		"signer":                  signerAddr,
	}, false, deps.SuiRPC)
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousReceiverObjects]{}, fmt.Errorf("failed to compile malicious receiver package: %w", err)
	}

	// Publish the package.
	packageId, tx, err := bind.PublishPackage(
		b.GetContext(),
		opts,
		deps.Client,
		bind.PublishRequest{
			CompiledModules: artifact.Modules,
			Dependencies:    artifact.Dependencies,
		},
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousReceiverObjects]{}, fmt.Errorf("failed to publish malicious receiver package: %w", err)
	}

	// The init function creates CCIPReceiverState (shared) and OwnerCap (transferred to sender).
	ownerCapObjectId, err1 := bind.FindObjectIdFromPublishTx(*tx, "malicious_receiver", "OwnerCap")
	if err1 != nil {
		return sui_ops.OpTxResult[DeployMaliciousReceiverObjects]{}, fmt.Errorf("failed to find OwnerCap object ID in publish tx: %w", err1)
	}

	receiverStateObjectId, err2 := bind.FindObjectIdFromPublishTx(*tx, "malicious_receiver", "CCIPReceiverState")
	if err2 != nil {
		return sui_ops.OpTxResult[DeployMaliciousReceiverObjects]{}, fmt.Errorf("failed to find CCIPReceiverState object ID in publish tx: %w", err2)
	}

	return sui_ops.OpTxResult[DeployMaliciousReceiverObjects]{
		Digest:    tx.Digest,
		PackageId: packageId,
		Objects: DeployMaliciousReceiverObjects{
			OwnerCapObjectId:          ownerCapObjectId,
			CCIPReceiverStateObjectId: receiverStateObjectId,
		},
	}, nil
}

var DeployCCIPMaliciousReceiverOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-malicious-receiver", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the TEST-only CCIP malicious receiver package (guard E2E smoke test)",
	deployMaliciousReceiverHandler,
)
