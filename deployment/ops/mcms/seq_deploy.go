package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
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

type DeployMCMSSeqOutput struct {
	AcceptOwnershipProposal mcms.TimelockProposal `json:"acceptOwnershipProposal"`
	PackageId               string                `json:"packageId"`
	Objects                 DeployMCMSObjects     `json:"objects"`
}

var DeployMCMSSequence = cld_ops.NewSequence(
	"sui-deploy-mcms-seq",
	semver.MustParse("0.1.0"),
	"Deploys the MCMS package, sets the initial configuration, init the ownership transfer to self and generates the proposal to accept the ownership",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployMCMSSeqInput) (DeployMCMSSeqOutput, error) {
		// Deploy MCMS first
		deployReport, err := cld_ops.ExecuteOperation(env, DeployMCMSOp, deps, cld_ops.EmptyInput{})
		if err != nil {
			return DeployMCMSSeqOutput{}, fmt.Errorf("failed to deploy MCMS: %w", err)
		}

		// Configure each timelock role if config is provided
		cfgMCMSInput := ConfigureMCMSSeqInput{
			ChainSelector:               input.ChainSelector,
			PackageId:                   deployReport.Output.PackageId,
			McmsAccountOwnerCapObjectId: deployReport.Output.Objects.McmsAccountOwnerCapObjectId,
			McmsAccountStateObjectId:    deployReport.Output.Objects.McmsAccountStateObjectId,
			McmsMultisigStateObjectId:   deployReport.Output.Objects.McmsMultisigStateObjectId,
			Bypasser:                    input.Bypasser,
			Proposer:                    input.Proposer,
			Canceller:                   input.Canceller,
		}
		_, err = cld_ops.ExecuteSequence(env, ConfigureMCMSSequence, deps, cfgMCMSInput)
		if err != nil {
			return DeployMCMSSeqOutput{}, fmt.Errorf("failed to configure MCMS: %w", err)
		}

		// Init the ownership transfer to self
		transferOwnershipInput := MCMSTransferOwnershipInput{
			McmsPackageID:   deployReport.Output.PackageId,
			OwnerCap:        deployReport.Output.Objects.McmsAccountOwnerCapObjectId,
			AccountObjectID: deployReport.Output.Objects.McmsAccountStateObjectId,
		}
		_, err = cld_ops.ExecuteOperation(env, MCMSTransferOwnershipOp, deps, transferOwnershipInput)
		if err != nil {
			return DeployMCMSSeqOutput{}, fmt.Errorf("failed to transfer ownership to MCMS: %w", err)
		}

		// Generate the proposal to accept the ownership
		proposalInput := ProposalGenerateInput{
			Defs: []cld_ops.Definition{
				MCMSAcceptOwnershipOp.Def(),
			},
			Inputs: []any{
				MCMSAcceptOwnershipInput{
					McmsPackageID:   deployReport.Output.PackageId,
					AccountObjectID: deployReport.Output.Objects.McmsAccountStateObjectId,
				},
			},
			// MCMS related
			MmcsPackageID:      deployReport.Output.PackageId,
			McmsStateObjID:     deployReport.Output.Objects.McmsMultisigStateObjectId,
			TimelockObjID:      deployReport.Output.Objects.TimelockObjectId,
			AccountObjID:       deployReport.Output.Objects.McmsAccountStateObjectId,
			RegistryObjID:      deployReport.Output.Objects.McmsRegistryObjectId,
			DeployerStateObjID: deployReport.Output.Objects.McmsDeployerStateObjectId,
			ChainSelector:      uint64(input.ChainSelector),
			// Proposal
			TimelockConfig: utils.TimelockConfig{
				MCMSAction:   types.TimelockActionSchedule,
				MinDelay:     0,
				OverrideRoot: false,
			},
		}

		acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(env, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		if err != nil {
			return DeployMCMSSeqOutput{}, fmt.Errorf("failed to generate accept ownership proposal: %w", err)
		}

		output := DeployMCMSSeqOutput{
			AcceptOwnershipProposal: acceptOwnershipProposalReport.Output,
			PackageId:               deployReport.Output.PackageId,
			Objects:                 deployReport.Output.Objects,
		}

		return output, nil
	},
)

// ConfigureMCMSSeqInput defines the input for configuring MCMS
type ConfigureMCMSSeqInput struct {
	ChainSelector               uint64 `json:"chainSelector" yaml:"chainSelector"`
	PackageId                   string `json:"packageId" yaml:"packageId"`
	McmsAccountOwnerCapObjectId string `json:"mcmsAccountOwnerCapObjectId" yaml:"mcmsAccountOwnerCapObjectId"`
	McmsAccountStateObjectId    string `json:"mcmsAccountStateObjectId" yaml:"mcmsAccountStateObjectId"`
	McmsMultisigStateObjectId   string `json:"mcmsMultisigStateObjectId" yaml:"mcmsMultisigStateObjectId"`

	// Optional configs for each timelock role
	// If nil, the role will not be configured
	Bypasser  *types.Config `json:"bypasser,omitempty" yaml:"bypasser,omitempty"`
	Proposer  *types.Config `json:"proposer,omitempty" yaml:"proposer,omitempty"`
	Canceller *types.Config `json:"canceller,omitempty" yaml:"canceller,omitempty"`
}

type ConfigureMCMSSeqOutput struct {
	Reports []cld_ops.Report[any, any]
}

var ConfigureMCMSSequence = cld_ops.NewSequence(
	"sui-configure-mcms-seq",
	semver.MustParse("0.1.0"),
	"Configures the MCMS package with the provided timelock roles configuration",
	configureMCMS,
)

func configureMCMS(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input ConfigureMCMSSeqInput) (ConfigureMCMSSeqOutput, error) {
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

	opReports := make([]cld_ops.Report[any, any], 0)
	for _, roleConfig := range roleConfigs {
		if roleConfig.config == nil {
			continue
		}

		setConfigInput := MCMSSetConfigInput{
			ChainSelector: input.ChainSelector,
			McmsPackageID: input.PackageId,
			OwnerCap:      input.McmsAccountOwnerCapObjectId,
			McmsObjectID:  input.McmsMultisigStateObjectId,
			Role:          roleConfig.role,
			Config:        *roleConfig.config,
		}

		report, err := cld_ops.ExecuteOperation(env, SetConfigMCMSOp, deps, setConfigInput)
		if err != nil {
			return ConfigureMCMSSeqOutput{}, fmt.Errorf("failed to set config for role %s: %w", roleConfig.name, err)
		}
		opReports = append(opReports, report.ToGenericReport())
		env.Logger.Infow("Set MCMS config", "role", roleConfig.name, "chainSelector", input.ChainSelector)
	}

	return ConfigureMCMSSeqOutput{Reports: opReports}, nil
}
