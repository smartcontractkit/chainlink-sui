package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	cslclient "github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

var DefaultTimelockExpirationInHours = 72

// MaxTimelockScheduleDelay caps every timelock delay value produced by
// tooling in this package. It applies in two places:
//
//  1. The per-op scheduling delay carried in a Schedule proposal
//     (TimelockConfig.MinDelay), enforced by TimelockConfig.Validate.
//  2. The new_min_delay payload of any mcms::mcms_timelock_update_min_delay
//     batch entry, enforced by assertUpdateMinDelayWithinCap in
//     op_proposal_generate.
//
// The Move source has no upper bound on either value. This cap is the
// deploy-side belt against the compound F8 finding: values large enough to
// trip the u64 overflow in `get_timestamp_seconds(clock) + delay` (in
// timelock_schedule), or to make min_delay > any u64 the schedule cap would
// allow, permanently brick the timelock.
const MaxTimelockScheduleDelay = 30 * 24 * time.Hour

// TimelockConfig is based on chainlink/deployment proposal utils
type TimelockConfig struct {
	MCMSAction   types.TimelockAction `json:"mcmsAction"`
	MinDelay     time.Duration        `json:"minDelay"`     // delay for timelock worker to execute the transfers.
	OverrideRoot bool                 `json:"overrideRoot"` // if true, override the previous root with the new one.
}

// Validate rejects configurations that would produce a proposal the Sui MCMS
// timelock cannot safely execute. Currently enforces the F8 scheduling-delay
// ceiling; a negative delay is also rejected since time.Duration is signed and
// would silently underflow when converted to u64 seconds on chain.
func (c TimelockConfig) Validate() error {
	if c.MCMSAction != types.TimelockActionSchedule {
		return nil
	}
	if c.MinDelay < 0 {
		return fmt.Errorf("timelock MinDelay must be non-negative, got %s", c.MinDelay)
	}
	if c.MinDelay > MaxTimelockScheduleDelay {
		return fmt.Errorf(
			"timelock MinDelay %s exceeds MaxTimelockScheduleDelay %s (F8 defense: larger values risk overflowing get_timestamp_seconds(clock) + delay in timelock_schedule)",
			c.MinDelay, MaxTimelockScheduleDelay,
		)
	}
	return nil
}

// AssertMinDelayWithinCap rejects any candidate on-chain min_delay value above
// MaxTimelockScheduleDelay. It is the single source of truth for the numeric
// bound so both the preventative payload check (in
// op_proposal_generate.assertUpdateMinDelayWithinCap) and the defensive
// state-read check (in op_set_config, Interpretation B) share the same
// definition of "too large". See finding F8 in findings-deduped.md.
func AssertMinDelayWithinCap(newMinDelaySeconds uint64) error {
	maxSeconds := uint64(MaxTimelockScheduleDelay / time.Second)
	if newMinDelaySeconds > maxSeconds {
		return fmt.Errorf(
			"F8 defense: min_delay %ds (~%s) exceeds MaxTimelockScheduleDelay %ds (%s); larger values risk bricking the timelock",
			newMinDelaySeconds, time.Duration(newMinDelaySeconds)*time.Second,
			maxSeconds, MaxTimelockScheduleDelay,
		)
	}
	return nil
}

type GenerateProposalInput struct {
	ChainSelector      uint64
	Client             cslclient.BindingsClient
	MCMSPackageID      string
	MCMSStateObjID     string
	AccountObjID       string
	RegistryObjID      string
	TimelockObjID      string
	DeployerStateObjID string
	Description        string
	BatchOp            types.BatchOperation
	TimelockConfig     TimelockConfig
}

func GenerateProposal(ctx context.Context, input GenerateProposalInput) (*mcms.TimelockProposal, error) {
	if err := input.TimelockConfig.Validate(); err != nil {
		return nil, fmt.Errorf("invalid timelock config: %w", err)
	}

	// Get action and delay from role
	var delay *types.Duration
	if input.TimelockConfig.MCMSAction == types.TimelockActionSchedule {
		delayDuration := types.NewDuration(input.TimelockConfig.MinDelay)
		delay = &delayDuration
	}
	role, err := getRoleFromAction(input.TimelockConfig.MCMSAction)
	if err != nil {
		return nil, fmt.Errorf("failed to get action from role: %w", err)
	}

	// Get OP Count from inspector
	devInspectSigner := signer.NewDevInspectSigner("0x0")
	inspector, err := suisdk.NewInspector(input.Client, devInspectSigner, input.MCMSPackageID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create inspector: %w", err)
	}
	opCount, err := inspector.GetOpCount(ctx, input.MCMSStateObjID)
	if err != nil {
		return nil, fmt.Errorf("failed to get op count: %w", err)
	}

	// Build metadata
	metadata, err := suisdk.NewChainMetadata(opCount, role, input.MCMSPackageID, input.MCMSStateObjID, input.AccountObjID, input.RegistryObjID, input.TimelockObjID, input.DeployerStateObjID)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain metadata: %w", err)
	}

	// Build proposal
	validUntilMs := uint32(time.Now().Add(time.Duration(DefaultTimelockExpirationInHours) * time.Hour).Unix())
	builder := mcms.NewTimelockProposalBuilder().
		SetVersion("v1").
		SetValidUntil(validUntilMs).
		SetDescription(input.Description).
		AddTimelockAddress(types.ChainSelector(input.ChainSelector), input.TimelockObjID).
		AddChainMetadata(types.ChainSelector(input.ChainSelector), metadata).
		AddOperation(input.BatchOp).
		SetAction(input.TimelockConfig.MCMSAction).
		SetOverridePreviousRoot(input.TimelockConfig.OverrideRoot)

	if delay != nil {
		builder.SetDelay(*delay)
	}

	return builder.Build()
}

// TransactionCallToMCMSTransaction converts an encoded Sui op call into an MCMS transaction.
// When call.LatestPackageID is set, it is propagated for upgraded-package PTB routing.
func TransactionCallToMCMSTransaction(call sui_ops.TransactionCall) (types.Transaction, error) {
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
		return types.Transaction{}, fmt.Errorf("create MCMS transaction: %w", err)
	}
	if call.LatestPackageID != "" {
		if setErr := suisdk.SetLatestPackageID(&tx, call.LatestPackageID); setErr != nil {
			return types.Transaction{}, fmt.Errorf("set latest package ID: %w", setErr)
		}
	}
	return tx, nil
}

// TransactionLatestPackageID returns the optional upgraded package ID stored on an MCMS transaction.
func TransactionLatestPackageID(tx types.Transaction) (string, error) {
	var fields suisdk.AdditionalFields
	if err := json.Unmarshal(tx.AdditionalFields, &fields); err != nil {
		return "", fmt.Errorf("unmarshal transaction additional fields: %w", err)
	}
	return fields.LatestPackageID, nil
}

func ExtractTransactionCall(output interface{}, operationID string) (sui_ops.TransactionCall, error) {
	jsonBytes, err := json.Marshal(output)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to marshal operation %s output: %w", operationID, err)
	}

	var outputMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &outputMap); err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to unmarshal operation %s output: %w", operationID, err)
	}

	callInterface, exists := outputMap["Call"]
	if !exists {
		return sui_ops.TransactionCall{}, fmt.Errorf("operation %s output does not have a Call field", operationID)
	}

	callBytes, err := json.Marshal(callInterface)
	if err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to marshal Call field for operation %s: %w", operationID, err)
	}

	var call sui_ops.TransactionCall
	if err := json.Unmarshal(callBytes, &call); err != nil {
		return sui_ops.TransactionCall{}, fmt.Errorf("failed to unmarshal Call field for operation %s: %w", operationID, err)
	}

	return call, nil
}

func getRoleFromAction(action types.TimelockAction) (suisdk.TimelockRole, error) {
	switch action {
	case types.TimelockActionSchedule:
		return suisdk.TimelockRoleProposer, nil
	case types.TimelockActionBypass:
		return suisdk.TimelockRoleBypasser, nil
	case types.TimelockActionCancel:
		return suisdk.TimelockRoleCanceller, nil
	default:
		// NewChainMetadata will always error on invalid action, but this is a safeguard
		return 0, fmt.Errorf("unsupported action: %v", action)
	}
}
