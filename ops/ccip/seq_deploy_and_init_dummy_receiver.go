package ccipops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployAndInitDummyReceiverInput struct {
	// For deployment
	DeployDummyReceiverInput
	// For registration
	CCIPObjectRefObjectID string
}

type DeployDummyReceiverSeqObjects struct {
	OwnerCapObjectID          string
	CCIPReceiverStateObjectID string
}

type DeployDummyReceiverSeqOutput struct {
	DummyReceiverPackageID string
	Objects                DeployDummyReceiverSeqObjects
}

var DeployAndInitDummyReceiverSequence = cld_ops.NewSequence(
	"sui-deploy-dummy-receiver-seq",
	semver.MustParse("0.1.0"),
	"Deploys CCIP dummy receiver and registers it with the receiver registry",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitDummyReceiverInput) (DeployDummyReceiverSeqOutput, error) {

		lggr := env.Logger
		lggr.Debugw("Deploying dummy receiver", "input", input)

		// Step 1: Deploy the dummy receiver
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPDummyReceiverOp, deps, input.DeployDummyReceiverInput)
		if err != nil {
			return DeployDummyReceiverSeqOutput{}, err
		}

		env.Logger.Infow("Dummy receiver deployed successfully",
			"packageId", deployReport.Output.PackageId,
			"receiverStateId", deployReport.Output.Objects.CCIPReceiverStateObjectID,
		)

		// Step 2: Register the dummy receiver with the receiver registry
		_, err = cld_ops.ExecuteOperation(
			env,
			RegisterDummyReceiverOp,
			deps,
			RegisterDummyReceiverInput{
				CCIPObjectRefObjectID:  input.CCIPObjectRefObjectID,
				DummyReceiverPackageID: deployReport.Output.PackageId,
			},
		)
		if err != nil {
			return DeployDummyReceiverSeqOutput{}, err
		}

		env.Logger.Infow("Dummy receiver registered successfully with receiver registry")

		return DeployDummyReceiverSeqOutput{
			DummyReceiverPackageID: deployReport.Output.PackageId,
			Objects: DeployDummyReceiverSeqObjects{
				OwnerCapObjectID:          deployReport.Output.Objects.OwnerCapObjectID,
				CCIPReceiverStateObjectID: deployReport.Output.Objects.CCIPReceiverStateObjectID,
			},
		}, nil
	},
)
