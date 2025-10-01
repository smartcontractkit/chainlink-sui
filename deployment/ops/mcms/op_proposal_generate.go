package mcmsops

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

type Input struct {
	// Ops Related
	Defs   []cld_ops.Definition
	Inputs []sui_ops.OpTxInput[any]

	// MCMS related
	MmcsPackageID  string `json:"mcmsPackageID"`
	McmsStateObjID string `json:"mcmsStateObjID"`
	TimelockObjID  string `json:"timelockObjID"`
	AccountObjID   string `json:"accountObjID"`
	RegistryObjID  string `json:"registryObjID"`

	// Proposal related
	Role  suisdk.TimelockRole `json:"role"`
	Delay time.Duration       `json:"delay"`

	// Chain related
	ChainSelector uint64 `json:"chainSelector"`
}

var GenerateProposalHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input Input) (output mcms.TimelockProposal, err error) {
	mcmsTxs := make([]mcmstypes.Transaction, len(input.Defs))

	for i, def := range input.Defs {
		op, err := b.OperationRegistry.Retrieve(def)
		if err != nil {
			return mcms.TimelockProposal{}, err
		}
		res, err := cld_ops.ExecuteOperation(b, op, any(deps), any(sui_ops.OpTxInput[any]{
			Input:     input.Inputs[i].Input,
			NoExecute: true, // we don't want to execute, just encode
		}))
		if err != nil {
			return mcms.TimelockProposal{}, fmt.Errorf("failed to execute operation %s: %w", def.ID, err)
		}
		output, ok := res.Output.(sui_ops.OpTxResult[any])
		if !ok {
			return mcms.TimelockProposal{}, fmt.Errorf("operation %s did not return OpTxResult output", def.ID)
		}
		tx, err := suisdk.NewTransactionWithStateObj(
			output.Call.Module,
			output.Call.Function,
			output.Call.PackageID,
			output.Call.Data,
			output.Call.Module,
			[]string{},
			output.Call.StateObjID,
		)
		mcmsTxs[i] = tx
	}

	op := types.BatchOperation{
		ChainSelector: types.ChainSelector(input.ChainSelector),
		Transactions:  mcmsTxs,
	}

	validUntilMs := uint32(time.Now().Add(time.Hour * 24).Unix())
	metadata, err := suisdk.NewChainMetadata(0, input.Role, input.MmcsPackageID, input.McmsStateObjID, input.AccountObjID, input.RegistryObjID, input.TimelockObjID)
	if err != nil {
		return mcms.TimelockProposal{}, fmt.Errorf("failed to create chain metadata: %w", err)
	}

	var action types.TimelockAction
	var delay *types.Duration
	switch input.Role {
	case suisdk.TimelockRoleProposer:
		action = types.TimelockActionSchedule
		delayDuration := types.NewDuration(input.Delay)
		delay = &delayDuration
	case suisdk.TimelockRoleBypasser:
		action = types.TimelockActionBypass
	default:
		return mcms.TimelockProposal{}, fmt.Errorf("unsupported role: %v", input.Role)
	}

	var description string = "Invokes the following set of operations: "
	for i, def := range input.Defs {
		if i > 0 {
			description += ", "
		}
		description += def.ID
	}

	builder := mcms.NewTimelockProposalBuilder().
		SetVersion("v1").
		SetValidUntil(validUntilMs).
		SetDescription(description).
		AddTimelockAddress(types.ChainSelector(input.ChainSelector), input.TimelockObjID).
		AddChainMetadata(types.ChainSelector(input.ChainSelector), metadata).
		AddOperation(op).
		SetAction(action)

	if delay != nil {
		builder.SetDelay(*delay)
	}

	timelockProposal, err := builder.Build()
	if err != nil {
		return mcms.TimelockProposal{}, fmt.Errorf("failed to build proposal: %w", err)
	}

	return *timelockProposal, nil
}

var MCMSDynamicProposalGenerateSeq = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "proposal", "generate"),
	semver.MustParse("0.1.0"),
	"Generates an MCMS timelock proposal that batches multiple operations based on the provided definitions and inputs",
	GenerateProposalHandler,
)
