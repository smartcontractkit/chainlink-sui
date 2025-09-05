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
	RegisterWithMCMSInput
	RegisterUpgradeCapInput
}

type DeployCCIPOnRampSeqObjects struct {
	StateObjectId string
}

type DeployCCIPOnRampSeqOutput struct {
	CCIPOnRampPackageId string
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
		input.OnRampInitializeInput.OnRampPackageId = deployReport.Output.PackageId
		input.OnRampInitializeInput.OnRampStateId = deployReport.Output.Objects.CCIPOnrampStateObjectId
		input.OnRampInitializeInput.OwnerCapObjectId = deployReport.Output.Objects.OwnerCapObjectId

		_, err = cld_ops.ExecuteOperation(env, OnRampInitializeOP, deps, input.OnRampInitializeInput)
		if err != nil {
			return DeployCCIPOnRampSeqOutput{}, err
		}

		applyDestChainConfigUpdateInput := ApplyDestChainConfigureOnRampInput{
			OnRampPackageId:           deployReport.Output.PackageId,
			OwnerCapObjectId:          deployReport.Output.Objects.OwnerCapObjectId,
			StateObjectId:             deployReport.Output.Objects.CCIPOnrampStateObjectId,
			DestChainSelector:         input.ApplyDestChainConfigureOnRampInput.DestChainSelector,
			DestChainEnabled:          input.ApplyDestChainConfigureOnRampInput.DestChainEnabled,
			DestChainAllowListEnabled: input.ApplyDestChainConfigureOnRampInput.DestChainAllowListEnabled,
		}

		_, err = cld_ops.ExecuteOperation(env, ApplyDestChainConfigUpdateOp, deps, applyDestChainConfigUpdateInput)
		if err != nil {
			return DeployCCIPOnRampSeqOutput{}, err
		}

		// Register OnRamp with MCMS registry first
		registerWithMCMSInput := RegisterWithMCMSInput{
			OnRampPackageId:  deployReport.Output.PackageId,
			OwnerCapObjectId: deployReport.Output.Objects.OwnerCapObjectId,
			RegistryId:       input.RegisterWithMCMSInput.RegistryId,
		}

		_, err = cld_ops.ExecuteOperation(env, RegisterWithMCMSOp, deps, registerWithMCMSInput)
		if err != nil {
			return DeployCCIPOnRampSeqOutput{}, err
		}

		// Register UpgradeCap with MCMS deployer
		registerUpgradeCapInput := RegisterUpgradeCapInput{
			OnRampPackageId:    deployReport.Output.PackageId,
			UpgradeCapObjectId: deployReport.Output.Objects.UpgradeCapObjectId,
			RegistryId:         input.RegisterUpgradeCapInput.RegistryId,
			DeployerStateId:    input.RegisterUpgradeCapInput.DeployerStateId,
		}

		_, err = cld_ops.ExecuteOperation(env, RegisterUpgradeCapOp, deps, registerUpgradeCapInput)
		if err != nil {
			return DeployCCIPOnRampSeqOutput{}, err
		}

		return DeployCCIPOnRampSeqOutput{
			CCIPOnRampPackageId: deployReport.Output.PackageId,
			Objects: DeployCCIPOnRampSeqObjects{
				StateObjectId: deployReport.Output.Objects.CCIPOnrampStateObjectId,
			},
		}, nil
	},
)
