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
	CCIPObjectRefObjectId string
	ReceiverStateParams   []string
}

type DeployDummyReceiverSeqObjects struct {
	OwnerCapObjectId          string
	CCIPReceiverStateObjectId string
}

type DeployDummyReceiverSeqOutput struct {
	DummyReceiverPackageId string
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
			"receiverStateId", deployReport.Output.Objects.CCIPReceiverStateObjectId,
		)

		// Step 2: Register the dummy receiver with the receiver registry
		_, err = cld_ops.ExecuteOperation(
			env,
			RegisterDummyReceiverOp,
			deps,
			RegisterDummyReceiverInput{
				CCIPObjectRefObjectId:  input.CCIPObjectRefObjectId,
				DummyReceiverPackageId: deployReport.Output.PackageId,
				ReceiverStateParams:    input.ReceiverStateParams,
			},
		)
		if err != nil {
			return DeployDummyReceiverSeqOutput{}, err
		}

		env.Logger.Infow("Dummy receiver registered successfully with receiver registry")

		return DeployDummyReceiverSeqOutput{
			DummyReceiverPackageId: deployReport.Output.PackageId,
			Objects: DeployDummyReceiverSeqObjects{
				OwnerCapObjectId:          deployReport.Output.Objects.OwnerCapObjectId,
				CCIPReceiverStateObjectId: deployReport.Output.Objects.CCIPReceiverStateObjectId,
			},
		}, nil
	},
)
