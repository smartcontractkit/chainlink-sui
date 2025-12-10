package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type ProposalGenerateInput struct {
	// Ops Related
	// Order matters, each definition should correspond to the input at the same index
	Defs   []cld_ops.Definition
	Inputs []any // Each element should be the corresponding input type for its operation

	// MCMS related
	MmcsPackageID      string `json:"mcmsPackageID"`
	McmsStateObjID     string `json:"mcmsStateObjID"`
	TimelockObjID      string `json:"timelockObjID"`
	AccountObjID       string `json:"accountObjID"`
	RegistryObjID      string `json:"registryObjID"`
	DeployerStateObjID string `json:"deployerStateObjID"`

	// Chain related
	ChainSelector uint64 `json:"chainSelector"`

	// Timelock related
	TimelockConfig utils.TimelockConfig `json:"timelockConfig"`
}

var MCMSDynamicProposalGenerateSeq = cld_ops.NewSequence(
	sui_ops.NewSuiOperationName("mcms", "proposal", "generate"),
	semver.MustParse("0.1.0"),
	"Generates an MCMS timelock proposal that batches multiple operations based on the provided definitions and inputs",
	generateProposalHandler,
)

var generateProposalHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ProposalGenerateInput) (output mcms.TimelockProposal, err error) {
	if len(input.Defs) != len(input.Inputs) {
		return mcms.TimelockProposal{}, fmt.Errorf("number of definitions (%d) does not match number of inputs (%d)", len(input.Defs), len(input.Inputs))
	}

	mcmsTxs := make([]mcmstypes.Transaction, len(input.Defs))

	for i, def := range input.Defs {
		op, err := b.OperationRegistry.Retrieve(def)
		if err != nil {
			return mcms.TimelockProposal{}, fmt.Errorf("failed to retrieve operation %s: %w", def.ID, err)
		}
		// Remove the signer to make the operations read-only, and prevent accidental tx sends during execution
		deps.Signer = nil
		res, err := cld_ops.ExecuteOperation(b, op, any(deps), input.Inputs[i])
		if err != nil {
			return mcms.TimelockProposal{}, fmt.Errorf("failed to execute operation %s: %w", def.ID, err)
		}
		// Extract the Call field
		call, err := utils.ExtractTransactionCall(res.Output, def.ID)
		if err != nil {
			return mcms.TimelockProposal{}, err
		}

		tx, err := suisdk.NewTransactionWithStateObj(
			call.Module,
			call.Function,
			call.PackageID,
			call.Data,
			call.Module,
			[]string{},
			call.StateObjID,
			call.TypeArgs,
		)
		if err != nil {
			return mcms.TimelockProposal{}, fmt.Errorf("failed to create transaction for operation %s: %w", def.ID, err)
		}
		mcmsTxs[i] = tx
	}

	op := mcmstypes.BatchOperation{
		ChainSelector: mcmstypes.ChainSelector(input.ChainSelector),
		Transactions:  mcmsTxs,
	}

	var description string = "Invokes the following set of operations: "
	for i, def := range input.Defs {
		if i > 0 {
			description += ", "
		}
		description += def.ID
	}

	proposalInput := utils.GenerateProposalInput{
		Client:             deps.Client,
		MCMSPackageID:      input.MmcsPackageID,
		MCMSStateObjID:     input.McmsStateObjID,
		TimelockObjID:      input.TimelockObjID,
		AccountObjID:       input.AccountObjID,
		RegistryObjID:      input.RegistryObjID,
		DeployerStateObjID: input.DeployerStateObjID,
		ChainSelector:      input.ChainSelector,
		TimelockConfig:     input.TimelockConfig,
		Description:        description,
		BatchOp:            op,
	}
	timelockProposal, err := utils.GenerateProposal(b.GetContext(), proposalInput)
	if err != nil {
		return mcms.TimelockProposal{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return *timelockProposal, nil
}
