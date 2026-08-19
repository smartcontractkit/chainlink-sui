package ownershipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type AcceptTokenPoolOwnershipInput struct {
	// MCMS related
	MCMSPackageId          string
	MCMSStateObjId         string
	MCMSTimelockObjId      string
	MCMSAccountObjId       string
	MCMSRegistryObjId      string
	MCMSDeployerStateObjId string

	// Proposal
	TimelockConfig utils.TimelockConfig
	ChainSelector  uint64

	// Token pools to accept ownership for. Nil entries are skipped, so callers
	// accept only the pools they select.
	ManagedTokenPool     *managedtokenpoolops.AcceptOwnershipManagedTokenPoolInput
	BurnMintTokenPool    *burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolInput
	LockReleaseTokenPool *lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolInput
}

var AcceptTokenPoolOwnershipSeq = cld_ops.NewSequence(
	"sui-accept-ownership-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Creates accept ownership proposal from MCMS for the selected token pools",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptTokenPoolOwnershipInput) (mcms.TimelockProposal, error) {
		defs := make([]cld_ops.Definition, 0, 3)
		inputs := make([]any, 0, 3)

		if input.ManagedTokenPool != nil {
			defs = append(defs, managedtokenpoolops.AcceptOwnershipManagedTokenPoolOp.Def())
			inputs = append(inputs, *input.ManagedTokenPool)
		}
		if input.BurnMintTokenPool != nil {
			defs = append(defs, burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolOp.Def())
			inputs = append(inputs, *input.BurnMintTokenPool)
		}
		if input.LockReleaseTokenPool != nil {
			defs = append(defs, lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolOp.Def())
			inputs = append(inputs, *input.LockReleaseTokenPool)
		}

		if len(defs) == 0 {
			return mcms.TimelockProposal{}, fmt.Errorf("no token pools selected for accept ownership")
		}

		proposalInput := mcmsops.ProposalGenerateInput{
			Defs:   defs,
			Inputs: inputs,

			// MCMS related
			MmcsPackageID:      input.MCMSPackageId,
			McmsStateObjID:     input.MCMSStateObjId,
			TimelockObjID:      input.MCMSTimelockObjId,
			AccountObjID:       input.MCMSAccountObjId,
			RegistryObjID:      input.MCMSRegistryObjId,
			DeployerStateObjID: input.MCMSDeployerStateObjId,

			// Proposal
			ChainSelector:  input.ChainSelector,
			TimelockConfig: input.TimelockConfig,
		}

		report, err := cld_ops.ExecuteSequence(env, mcmsops.MCMSDynamicProposalGenerateSeq, deps, proposalInput)
		if err != nil {
			return mcms.TimelockProposal{}, err
		}

		return report.Output, nil
	},
)
