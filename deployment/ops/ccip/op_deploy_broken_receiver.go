package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// DeployBrokenReceiverObjects holds the objects produced by deploying the broken
// receiver. Unlike the dummy receiver, the broken receiver's init only creates an
// OwnerCap (it has no shared receiver-state object).
type DeployBrokenReceiverObjects struct {
	OwnerCapObjectId string
}

type DeployBrokenReceiverInput struct {
	CCIPPackageId string
	McmsPackageId string
	McmsOwner     string
}

var deployBrokenReceiverHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployBrokenReceiverInput) (output sui_ops.OpTxResult[DeployBrokenReceiverObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[DeployBrokenReceiverObjects]{}, err
	}

	// Compile the broken receiver package. Its ccip_receive<T> signature produces a
	// normalized ABI containing {"Vector": {"TypeParameter": 0}}, the poison shape
	// from report 71024.
	artifact, err := bind.CompilePackage(contracts.CCIPBrokenReceiver, map[string]string{
		"ccip":                 input.CCIPPackageId,
		"ccip_broken_receiver": "0x0",
		"mcms":                 input.McmsPackageId,
		"mcms_owner":           input.McmsOwner,

		"signer": signerAddr,
	}, false, deps.SuiRPC)
	if err != nil {
		return sui_ops.OpTxResult[DeployBrokenReceiverObjects]{}, fmt.Errorf("failed to compile broken receiver package: %w", err)
	}

	// Publish the package
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
		return sui_ops.OpTxResult[DeployBrokenReceiverObjects]{}, fmt.Errorf("failed to publish broken receiver package: %w", err)
	}

	// The init function creates and transfers an OwnerCap to the sender.
	ownerCapObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "broken_receiver", "OwnerCap")
	if err != nil {
		return sui_ops.OpTxResult[DeployBrokenReceiverObjects]{}, fmt.Errorf("failed to find OwnerCap object ID in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployBrokenReceiverObjects]{
		Digest:    tx.Digest,
		PackageId: packageId,
		Objects: DeployBrokenReceiverObjects{
			OwnerCapObjectId: ownerCapObjectId,
		},
	}, nil
}

var DeployCCIPBrokenReceiverOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-broken-receiver", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP broken receiver package (testing only; generic ccip_receive)",
	deployBrokenReceiverHandler,
)
