package ownershipops

import (
	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

// AcceptCCIPStateObjectOnlyInput carries MCMS + CCIP state object fields for a single accept_ownership tx in a proposal.
type AcceptCCIPStateObjectOnlyInput struct {
	MCMSPackageId          string
	MCMSStateObjId         string
	MCMSTimelockObjId      string
	MCMSAccountObjId       string
	MCMSRegistryObjId      string
	MCMSDeployerStateObjId string

	TimelockConfig utils.TimelockConfig
	ChainSelector  uint64

	CCIPPackageId string
	CCIPObjectRef string
}

// AcceptCCIPStateObjectOwnershipOnlySeq builds an MCMS timelock proposal with only
// ccip::state_object::accept_ownership (no router / onramp / offramp).
var AcceptCCIPStateObjectOwnershipOnlySeq = cld_ops.NewSequence(
	"sui-accept-ownership-ccip-state-object-only-seq",
	semver.MustParse("0.1.0"),
	"MCMS proposal: accept CCIP package state-object ownership only",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptCCIPStateObjectOnlyInput) (mcms.TimelockProposal, error) {
		proposalInput := mcmsops.ProposalGenerateInput{
			Defs: []cld_ops.Definition{
				ccipops.AcceptOwnershipStateObjectOp.Def(),
			},
			Inputs: []any{
				ccipops.AcceptOwnershipStateObjectInput{
					CCIPPackageId:         input.CCIPPackageId,
					CCIPObjectRefObjectId: input.CCIPObjectRef,
				},
			},
			MmcsPackageID:      input.MCMSPackageId,
			McmsStateObjID:     input.MCMSStateObjId,
			TimelockObjID:      input.MCMSTimelockObjId,
			AccountObjID:       input.MCMSAccountObjId,
			RegistryObjID:      input.MCMSRegistryObjId,
			DeployerStateObjID: input.MCMSDeployerStateObjId,
			ChainSelector:      input.ChainSelector,
			TimelockConfig:     input.TimelockConfig,
		}

		report, err := cld_ops.ExecuteSequence(env, mcmsops.MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		if err != nil {
			return mcms.TimelockProposal{}, err
		}

		return report.Output, nil
	},
)
