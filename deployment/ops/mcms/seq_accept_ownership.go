package mcmsops

import (
    "fmt"

    "github.com/Masterminds/semver/v3"
    "github.com/smartcontractkit/mcms"
    "github.com/smartcontractkit/mcms/types"

    cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

    sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
    "github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type AcceptMCMSOwnershipSeqInput struct {
    ChainSelector             uint64 `json:"chainSelector" yaml:"chainSelector"`
    PackageId                 string `json:"packageId" yaml:"packageId"`
    McmsMultisigStateObjectId string `json:"mcmsMultisigStateObjectId" yaml:"mcmsMultisigStateObjectId"`
    TimelockObjectId          string `json:"timelockObjectId" yaml:"timelockObjectId"`
    McmsAccountStateObjectId  string `json:"mcmsAccountStateObjectId" yaml:"mcmsAccountStateObjectId"`
    McmsRegistryObjectId      string `json:"mcmsRegistryObjectId" yaml:"mcmsRegistryObjectId"`
    McmsDeployerStateObjectId string `json:"mcmsDeployerStateObjectId" yaml:"mcmsDeployerStateObjectId"`
}

var AcceptMCMSOwnershipSequence = cld_ops.NewSequence(
    "sui-accept-mcms-ownership-seq",
    semver.MustParse("0.1.0"),
    "Generates the MCMS proposal to accept MCMS ownership via the timelock",
    acceptMCMSOwnership,
)

func acceptMCMSOwnership(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptMCMSOwnershipSeqInput) (mcms.TimelockProposal, error) {
    proposalInput := ProposalGenerateInput{
        Defs: []cld_ops.Definition{
            MCMSAcceptOwnershipOp.Def(),
        },
        Inputs: []any{
            MCMSAcceptOwnershipInput{
                McmsPackageID:   input.PackageId,
                AccountObjectID: input.McmsAccountStateObjectId,
            },
        },
        MmcsPackageID:      input.PackageId,
        McmsStateObjID:     input.McmsMultisigStateObjectId,
        TimelockObjID:      input.TimelockObjectId,
        AccountObjID:       input.McmsAccountStateObjectId,
        RegistryObjID:      input.McmsRegistryObjectId,
        DeployerStateObjID: input.McmsDeployerStateObjectId,
        ChainSelector:      input.ChainSelector,
        TimelockConfig: utils.TimelockConfig{
            MCMSAction:   types.TimelockActionSchedule,
            MinDelay:     0,
            OverrideRoot: false,
        },
    }

    report, err := cld_ops.ExecuteSequence(env, MCMSDynamicProposalGenerateSeq, deps, proposalInput)
    if err != nil {
        return mcms.TimelockProposal{}, fmt.Errorf("failed to generate accept ownership proposal: %w", err)
    }

    return report.Output, nil
}
