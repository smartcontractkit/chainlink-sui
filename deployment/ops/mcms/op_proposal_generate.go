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

// mcmsModuleName is the on-chain module that owns the timelock admin surface.
const mcmsModuleName = "mcms"

// bypassForbiddenMcmsFunctions is the set of mcms-module functions that MUST
// NOT be reachable through the BYPASSER-signed execute path. They mutate the
// timelock's own scheduling / blocklist / signer state and are semantically
// "timelock administering itself"; the BYPASSER role is meant to fast-execute
// external application calls, not to reconfigure the timelock. The Move
// wrappers currently hardcode TIMELOCK_ROLE and cannot distinguish the
// originating role from the ExecutingCallbackParams they consume, so the
// deploy-side belt is to refuse to build such a proposal in the first place.
// See finding F30 in findings-deduped.md.
var bypassForbiddenMcmsFunctions = map[string]struct{}{
	"mcms_timelock_update_min_delay": {},
	"mcms_timelock_block_function":   {},
	"mcms_timelock_unblock_function": {},
	"mcms_timelock_cancel":           {},
	"mcms_set_config":                {},
}

// assertBypassAllowsCall rejects calls that must not flow through the
// BYPASSER path. It is a no-op for any other timelock action.
func assertBypassAllowsCall(action mcmstypes.TimelockAction, call sui_ops.TransactionCall) error {
	if action != mcmstypes.TimelockActionBypass {
		return nil
	}
	if call.Module != mcmsModuleName {
		return nil
	}
	if _, forbidden := bypassForbiddenMcmsFunctions[call.Function]; !forbidden {
		return nil
	}
	return fmt.Errorf(
		"F30 defense: refusing to build bypass proposal that calls mcms::%s; timelock-admin functions must be routed through the schedule path",
		call.Function,
	)
}

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

		if err := assertBypassAllowsCall(input.TimelockConfig.MCMSAction, call); err != nil {
			return mcms.TimelockProposal{}, fmt.Errorf("operation %s: %w", def.ID, err)
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
		// If the op signalled an upgraded package, record it so the timelock converter
		// propagates it as InternalLatestPackageIDs in the outer batch operation.
		if call.LatestPackageID != "" {
			if setErr := suisdk.SetLatestPackageID(&tx, call.LatestPackageID); setErr != nil {
				return mcms.TimelockProposal{}, fmt.Errorf("failed to set latest package ID for operation %s: %w", def.ID, setErr)
			}
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
