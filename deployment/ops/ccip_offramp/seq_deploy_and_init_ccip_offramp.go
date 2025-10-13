package offrampops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployAndInitCCIPOffRampSeqInput struct {
	DeployCCIPOffRampInput
	InitializeOffRampInput
	CCIPObjectRefId     string
	CommitOCR3Config    SetOCR3ConfigInput
	ExecutionOCR3Config SetOCR3ConfigInput
}

type DeployCCIPOffRampSeqObjects struct {
	StateObjectId      string
	OwnerCapId         string
	UpgradeCapObjectId string
}

type DeployCCIPOffRampSeqOutput struct {
	CCIPOffRampPackageId string
	Objects              DeployCCIPOffRampSeqObjects
}

var DeployAndInitCCIPOffRampSequence = cld_ops.NewSequence(
	"sui-deploy-ccip-offramp-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial CCIP offRamp configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitCCIPOffRampSeqInput) (DeployCCIPOffRampSeqOutput, error) {

		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPOffRampOp, deps, input.DeployCCIPOffRampInput)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, err
		}

		input.InitializeOffRampInput.OffRampPackageId = deployReport.Output.PackageId
		input.InitializeOffRampInput.OwnerCapObjectId = deployReport.Output.Objects.OwnerCapObjectId
		input.InitializeOffRampInput.OffRampStateId = deployReport.Output.Objects.CCIPOffRampStateObjectId

		_, err = cld_ops.ExecuteOperation(env, InitializeOffRampOp, deps, input.InitializeOffRampInput)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, err
		}

		// init transfer to mcms
		executeOwnershipTransferToMcmsInput := TransferOwnershipOffRampInput{
			OffRampPackageId:     deployReport.Output.PackageId,
			CCIPObjectRefId:      input.CCIPObjectRefId,
			OffRampStateObjectId: deployReport.Output.Objects.CCIPOffRampStateObjectId,
			OwnerCapObjectId:     deployReport.Output.Objects.OwnerCapObjectId,
			To:                   input.DeployCCIPOffRampInput.MCMSPackageId,
		}
		_, err = cld_ops.ExecuteOperation(env, TransferOwnershipOffRampOp, deps, executeOwnershipTransferToMcmsInput)
		if err != nil {
			return DeployCCIPOffRampSeqOutput{}, err
		}

		return DeployCCIPOffRampSeqOutput{
			CCIPOffRampPackageId: deployReport.Output.PackageId,
			Objects: DeployCCIPOffRampSeqObjects{
				StateObjectId:      deployReport.Output.Objects.CCIPOffRampStateObjectId,
				OwnerCapId:         deployReport.Output.Objects.OwnerCapObjectId,
				UpgradeCapObjectId: deployReport.Output.Objects.UpgradeCapObjectId,
			},
		}, nil
	},
)
