package mcmsops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"
)

// DeployMCMSSeqInput defines the input for deploying MCMS with timelock roles configuration
type DeployMCMSSeqInput struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`

	// Optional configs for each timelock role
	// If nil, the role will not be configured
	Bypasser  *types.Config `json:"bypasser,omitempty" yaml:"bypasser,omitempty"`
	Proposer  *types.Config `json:"proposer,omitempty" yaml:"proposer,omitempty"`
	Canceller *types.Config `json:"canceller,omitempty" yaml:"canceller,omitempty"`
}

var DeployMCMSSequence = cld_ops.NewSequence(
	"sui-deploy-mcms-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial MCMS configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployMCMSSeqInput) (sui_ops.OpTxResult[DeployMCMSObjects], error) {
		// Deploy MCMS first
		deployReport, err := cld_ops.ExecuteOperation(env, DeployMCMSOp, deps, cld_ops.EmptyInput{})
		if err != nil {
			return sui_ops.OpTxResult[DeployMCMSObjects]{}, err
		}

		// Configure each timelock role if config is provided
		roleConfigs := []struct {
			config *types.Config
			role   suisdk.TimelockRole
			name   string
		}{
			{input.Bypasser, suisdk.TimelockRoleBypasser, "Bypasser"},
			{input.Canceller, suisdk.TimelockRoleCanceller, "Canceller"},
			{input.Proposer, suisdk.TimelockRoleProposer, "Proposer"},
		}

		for _, roleConfig := range roleConfigs {
			if roleConfig.config != nil {
				setConfigInput := MCMSSetConfigInput{
					ChainSelector: input.ChainSelector,
					McmsPackageID: deployReport.Output.PackageId,
					OwnerCap:      deployReport.Output.Objects.McmsAccountOwnerCapObjectId,
					McmsObjectID:  deployReport.Output.Objects.McmsMultisigStateObjectId,
					Role:          roleConfig.role,
					Config:        *roleConfig.config,
				}

				_, err = cld_ops.ExecuteOperation(env, SetConfigMCMSOp, deps, setConfigInput)
				if err != nil {
					return sui_ops.OpTxResult[DeployMCMSObjects]{}, err
				}

				env.Logger.Infow("Set MCMS config", "role", roleConfig.name, "chainSelector", input.ChainSelector)
			}
		}

		return deployReport.Output, nil
	},
)
