package onrampops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployAndInitCCIPOnRampSeqInput struct {
	DeployCCIPOnRampInput
	OnRampInitializeInput
	ApplyDestChainConfigureOnRampInput
	ApplyAllowListUpdatesInput
}

type DeployCCIPOnRampSeqObjects struct {
	StateObjectID string
}

type DeployCCIPOnRampSeqOutput struct {
	CCIPOnRampPackageID string
	Objects             DeployCCIPOnRampSeqObjects
}

var DeployAndInitCCIPOnRampSequence = cld_ops.NewSequence(
	"sui-deploy-ccip-onramp-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial CCIP onRamp configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitCCIPOnRampSeqInput) (DeployCCIPOnRampSeqOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPOnRampOp, deps, input.DeployCCIPOnRampInput)
		if err != nil {
			return DeployCCIPOnRampSeqOutput{}, err
		}

		// Prepare updated input for initialization
		input.OnRampInitializeInput.OnRampPackageID = deployReport.Output.PackageId
		input.OnRampInitializeInput.OnRampStateID = deployReport.Output.Objects.CCIPOnrampStateObjectID
		input.OnRampInitializeInput.OwnerCapObjectID = deployReport.Output.Objects.OwnerCapObjectID

		_, err = cld_ops.ExecuteOperation(env, OnRampInitializeOP, deps, input.OnRampInitializeInput)
		if err != nil {
			return DeployCCIPOnRampSeqOutput{}, err
		}

		applyDestChainConfigUpdateInput := ApplyDestChainConfigureOnRampInput{
			OnRampPackageID:           deployReport.Output.PackageId,
			OwnerCapObjectID:          deployReport.Output.Objects.OwnerCapObjectID,
			StateObjectID:             deployReport.Output.Objects.CCIPOnrampStateObjectID,
			DestChainSelector:         input.ApplyDestChainConfigureOnRampInput.DestChainSelector,
			DestChainEnabled:          input.ApplyDestChainConfigureOnRampInput.DestChainEnabled,
			DestChainAllowListEnabled: input.ApplyDestChainConfigureOnRampInput.DestChainAllowListEnabled,
		}

		_, err = cld_ops.ExecuteOperation(env, ApplyDestChainConfigUpdateOp, deps, applyDestChainConfigUpdateInput)
		if err != nil {
			return DeployCCIPOnRampSeqOutput{}, err
		}

		return DeployCCIPOnRampSeqOutput{
			CCIPOnRampPackageID: deployReport.Output.PackageId,
			Objects: DeployCCIPOnRampSeqObjects{
				StateObjectID: deployReport.Output.Objects.CCIPOnrampStateObjectID,
			},
		}, nil
	},
)
