package mcmsops

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

var DefaultTimelockExpirationInHours = 72

type ProposalGenerateInput struct {
	// Ops Related
	// Order matters, each definition should correspond to the input at the same index
	Defs   []cld_ops.Definition
	Inputs []any // Each element should be the corresponding input type for its operation

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

var generateProposalHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ProposalGenerateInput) (output mcms.TimelockProposal, err error) {
	if len(input.Defs) != len(input.Inputs) {
		return mcms.TimelockProposal{}, fmt.Errorf("number of definitions (%d) does not match number of inputs (%d)", len(input.Defs), len(input.Inputs))
	}
	mcmsTxs := make([]mcmstypes.Transaction, len(input.Defs))

	for i, def := range input.Defs {
		op, err := b.OperationRegistry.Retrieve(def)
		if err != nil {
			return mcms.TimelockProposal{}, err
		}
		// Remove the signer to make the operations read-only, and prevent accidental tx sends during execution
		deps.Signer = nil
		res, err := cld_ops.ExecuteOperation(b, op, any(deps), input.Inputs[i])
		if err != nil {
			return mcms.TimelockProposal{}, fmt.Errorf("failed to execute operation %s: %w", def.ID, err)
		}
		// Use reflection to extract the Call field from the OpTxResult regardless of its specific type
		outputValue := reflect.ValueOf(res.Output)
		if outputValue.Kind() != reflect.Struct {
			return mcms.TimelockProposal{}, fmt.Errorf("operation %s did not return a struct output", def.ID)
		}

		callField := outputValue.FieldByName("Call")
		if !callField.IsValid() {
			return mcms.TimelockProposal{}, fmt.Errorf("operation %s output does not have a Call field", def.ID)
		}

		call, ok := callField.Interface().(sui_ops.TransactionCall)
		if !ok {
			return mcms.TimelockProposal{}, fmt.Errorf("operation %s Call field is not a TransactionCall", def.ID)
		}

		tx, err := suisdk.NewTransactionWithStateObj(
			call.Module,
			call.Function,
			call.PackageID,
			call.Data,
			call.Module,
			[]string{},
			call.StateObjID,
		)
		mcmsTxs[i] = tx
	}

	op := types.BatchOperation{
		ChainSelector: types.ChainSelector(input.ChainSelector),
		Transactions:  mcmsTxs,
	}

	validUntilMs := uint32(time.Now().Add(time.Duration(DefaultTimelockExpirationInHours) * time.Hour).Unix())

	// Get OP Count from inspector
	inspector, err := suisdk.NewInspector(deps.Client, deps.Signer, input.MmcsPackageID, input.Role)
	if err != nil {
		return mcms.TimelockProposal{}, fmt.Errorf("failed to create inspector: %w", err)
	}
	// TODO: somehow callopts with signer is required for get functions.
	opCount, err := inspector.GetOpCount(b.GetContext(), input.MmcsPackageID)
	if err != nil {
		return mcms.TimelockProposal{}, fmt.Errorf("failed to get op count: %w", err)
	}

	metadata, err := suisdk.NewChainMetadata(opCount, input.Role, input.MmcsPackageID, input.McmsStateObjID, input.AccountObjID, input.RegistryObjID, input.TimelockObjID)
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
	case suisdk.TimelockRoleCanceller:
		action = types.TimelockActionCancel
	default:
		// NewChainMetadata will always error on invalid role, but this is a safeguard
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

var MCMSDynamicProposalGenerateSeq = cld_ops.NewSequence(
	sui_ops.NewSuiOperationName("mcms", "proposal", "generate"),
	semver.MustParse("0.1.0"),
	"Generates an MCMS timelock proposal that batches multiple operations based on the provided definitions and inputs",
	generateProposalHandler,
)
