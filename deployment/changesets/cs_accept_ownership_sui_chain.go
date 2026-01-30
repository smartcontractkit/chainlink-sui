package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type AcceptOwnershipCCIPConfig struct {
	SuiChainSelector uint64
}

var _ cldf.ChangeSetV2[AcceptOwnershipCCIPConfig] = AcceptOwnershipCCIP{}

// AcceptOwnershipCCIP deploys Sui chain packages and modules
type AcceptOwnershipCCIP struct{}

// Apply implements deployment.ChangeSetV2.
func (d AcceptOwnershipCCIP) Apply(e cldf.Environment, config AcceptOwnershipCCIPConfig) (cldf.ChangesetOutput, error) {
	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]
	signer := suiChain.Signer

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(1_000_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	// in case the registry is not loaded with all operations. Needed to build accept ownership proposals
	for i := range opregistry.AllOperations {
		cld_ops.RegisterOperation(e.OperationsBundle.OperationRegistry, opregistry.AllOperations[i])
	}

	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.SuiChainSelector]

	// Generate the proposal to accept the ownership of the CCIP contracts. Get the addresses from AB
	proposalInput := ownershipops.AcceptCCIPOwnershipInput{
		ChainSelector: config.SuiChainSelector,

		// MCMS related
		MCMSPackageId:          state.MCMSPackageID,
		MCMSStateObjId:         state.MCMSStateObjectID,
		MCMSTimelockObjId:      state.MCMSTimelockObjectID,
		MCMSAccountObjId:       state.MCMSAccountStateObjectID,
		MCMSRegistryObjId:      state.MCMSRegistryObjectID,
		MCMSDeployerStateObjId: state.MCMSDeployerStateObjectID,

		CCIPPackageId: state.CCIPAddress,
		CCIPObjectRef: state.CCIPObjectRef,

		RouterPackageId:     state.CCIPRouterAddress,
		RouterStateObjectId: state.CCIPRouterStateObjectID,

		OnRampPackageId:     state.OnRampAddress,
		OnRampStateObjectId: state.OnRampStateObjectId,

		OffRampPackageId:     state.OffRampAddress,
		OffRampStateObjectId: state.OffRampStateObjectId,

		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionSchedule,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}

	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, ownershipops.AcceptCCIPOwnershipSeq, deps, proposalInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{acceptOwnershipProposalReport.Output},
	}, nil
}

// TODO
// VerifyPreconditions imsplements deployment.ChangeSetV2.
func (d AcceptOwnershipCCIP) VerifyPreconditions(e cldf.Environment, config AcceptOwnershipCCIPConfig) error {
	return nil
}
