package ownershipops

import (
	"encoding/json"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
)

type AcceptCCIPOwnershipInput struct {
	// MCMS related
	MCMSPackageId     string
	MCMSStateObjId    string
	MCMSTimelockObjId string
	MCMSAccountObjId  string
	MCMSRegistryObjId string

	// Proposal
	Role          suisdk.TimelockRole
	ChainSelector uint64

	// CCIP related
	CCIPPackageId string
	CCIPObjectRef string

	// Router related
	RouterPackageId     string
	RouterStateObjectId string

	// OnRamp related
	OnRampPackageId     string
	OnRampStateObjectId string

	// OffRamp related
	OffRampPackageId     string
	OffRampStateObjectId string
}

var AcceptCCIPOwnershipSeq = cld_ops.NewSequence(
	"sui-accept-ownership-ccip-seq",
	semver.MustParse("0.1.0"),
	"Creates accept ownership proposal from MCMS for every CCIP contract",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptCCIPOwnershipInput) (mcms.TimelockProposal, error) {
		// Generate the proposal to accept the ownership of the deployed contracts
		proposalInput := mcmsops.ProposalGenerateInput{
			Defs: []operations.Definition{
				ccipops.AcceptOwnershipStateObjectOp.Def(),
				routerops.AcceptOwnershipOp.Def(),
				onrampops.AcceptOwnershipOnRampOp.Def(),
				offrampops.AcceptOwnershipOffRampOp.Def(),
			},
			Inputs: []any{
				ccipops.AcceptOwnershipStateObjectInput{
					CCIPPackageId:         input.CCIPPackageId,
					CCIPObjectRefObjectId: input.CCIPObjectRef,
				},
				routerops.AcceptOwnershipInput{
					RouterPackageId:     input.RouterPackageId,
					RouterStateObjectId: input.RouterStateObjectId,
				},
				onrampops.AcceptOwnershipOnRampInput{
					OnRampPackageId: input.OnRampPackageId,
					CCIPObjectRefId: input.CCIPObjectRef,
					StateObjectId:   input.OnRampStateObjectId,
				},
				offrampops.AcceptOwnershipOffRampInput{
					OffRampPackageId:     input.OffRampPackageId,
					OffRampRefObjectId:   input.CCIPObjectRef,
					OffRampStateObjectId: input.OffRampStateObjectId,
				},
			},
			// MCMS related
			MmcsPackageID:  input.MCMSPackageId,
			McmsStateObjID: input.MCMSStateObjId,
			TimelockObjID:  input.MCMSTimelockObjId,
			AccountObjID:   input.MCMSAccountObjId,
			RegistryObjID:  input.MCMSRegistryObjId,

			// Proposal
			Role: input.Role,

			ChainSelector: input.ChainSelector,
		}

		jsonbytes, _ := json.MarshalIndent(proposalInput, "", " ")
		fmt.Println("Accept CCIP Ownership Proposal Input: ", string(jsonbytes))

		acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(env, mcmsops.MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		if err != nil {
			return mcms.TimelockProposal{}, err
		}

		return acceptOwnershipProposalReport.Output, nil
	},
)
