package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
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
	deployMCMS,
)

func deployMCMS(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployMCMSSeqInput) (DeployMCMSSeqOutput, error) {
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

	// Generate accept ownership proposal
	acceptOwnershipInput := AcceptMCMSOwnershipSeqInput{
		ChainSelector:             input.ChainSelector,
		PackageId:                 deployReport.Output.PackageId,
		McmsAccountStateObjectId:  deployReport.Output.Objects.McmsAccountStateObjectId,
		McmsDeployerStateObjectId: deployReport.Output.Objects.McmsDeployerStateObjectId,
		McmsMultisigStateObjectId: deployReport.Output.Objects.McmsMultisigStateObjectId,
		McmsRegistryObjectId:      deployReport.Output.Objects.McmsRegistryObjectId,
		TimelockObjectId:          deployReport.Output.Objects.TimelockObjectId,
	}

	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(env, AcceptMCMSOwnershipSequence, deps, acceptOwnershipInput)
	if err != nil {
		return DeployMCMSSeqOutput{}, fmt.Errorf("failed to generate accept ownership proposal: %w", err)
	}

	output := DeployMCMSSeqOutput{
		AcceptOwnershipProposal: acceptOwnershipProposalReport.Output,
		PackageId:               deployReport.Output.PackageId,
		Objects:                 deployReport.Output.Objects,
	}

	return output, nil
}
