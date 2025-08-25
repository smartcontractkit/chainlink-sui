package offrampops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployAndInitCCIPOffRampSeqInput struct {
	DeployCCIPOffRampInput
	InitializeOffRampInput
	CommitOCR3Config    SetOCR3ConfigInput
	ExecutionOCR3Config SetOCR3ConfigInput
}

type DeployCCIPOffRampSeqObjects struct {
	StateObjectID string
	OwnerCapID    string
}

type DeployCCIPOffRampSeqOutput struct {
	CCIPOffRampPackageID string
	Objects              DeployCCIPOffRampSeqObjects
}

var DeployAndInitCCIPOffRampSequence = cld_ops.NewSequence(
	"sui-deploy-ccip-offramp-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial CCIP offRamp configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitCCIPOffRampSeqInput) (DeployCCIPOffRampSeqOutput, error) {

		lggr := env.Logger

		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPOffRampOp, deps, input.DeployCCIPOffRampInput)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, err
		}

		input.InitializeOffRampInput.OffRampPackageID = deployReport.Output.PackageID
		input.InitializeOffRampInput.OwnerCapObjectID = deployReport.Output.Objects.OwnerCapObjectID
		input.InitializeOffRampInput.OffRampStateID = deployReport.Output.Objects.CCIPOffRampStateObjectID

		_, err = cld_ops.ExecuteOperation(env, InitializeOffRampOp, deps, input.InitializeOffRampInput)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, err
		}

		lggr.Infow("SetOCR3Config for COMMIT")
		input.CommitOCR3Config.OffRampPackageID = deployReport.Output.PackageID
		input.CommitOCR3Config.OwnerCapObjectID = deployReport.Output.Objects.OwnerCapObjectID
		input.CommitOCR3Config.OffRampStateID = deployReport.Output.Objects.CCIPOffRampStateObjectID
		_, err = cld_ops.ExecuteOperation(env, SetOCR3ConfigOp, deps, input.CommitOCR3Config)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, fmt.Errorf("failed to set COMMIT OCR3 config: %w", err)
		}

		lggr.Infow("SetOCR3Config for EXECUTION")
		input.ExecutionOCR3Config.OffRampPackageID = deployReport.Output.PackageID
		input.ExecutionOCR3Config.OwnerCapObjectID = deployReport.Output.Objects.OwnerCapObjectID
		input.ExecutionOCR3Config.OffRampStateID = deployReport.Output.Objects.CCIPOffRampStateObjectID
		_, err = cld_ops.ExecuteOperation(env, SetOCR3ConfigOp, deps, input.ExecutionOCR3Config)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, fmt.Errorf("failed to set EXECUTION OCR3 config: %w", err)
		}

		return DeployCCIPOffRampSeqOutput{
			CCIPOffRampPackageID: deployReport.Output.PackageID,
			Objects: DeployCCIPOffRampSeqObjects{
				StateObjectID: deployReport.Output.Objects.CCIPOffRampStateObjectID,
				OwnerCapID:    deployReport.Output.Objects.OwnerCapObjectID,
			},
		}, nil
	},
)
