package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// UpdateMinDelayMCMSInput builds an mcms::mcms_timelock_update_min_delay call
// that sets a timelock's global min_delay to NewMinDelaySec seconds. The op
// only prepares the call for inclusion in an MCMS timelock proposal; it does
// not execute it, because updating min_delay must itself go through the
// timelock it targets.
type UpdateMinDelayMCMSInput struct {
	// McmsPackageID is the MCMS Move package exposing mcms_timelock_update_min_delay.
	McmsPackageID string `yaml:"mcmsPackageID" json:"mcmsPackageID"`
	// TimelockObjID is the timelock object whose min_delay is updated. It is also
	// the state object the proposal dispatches the call against.
	TimelockObjID string `yaml:"timelockObjID" json:"timelockObjID"`
	// NewMinDelaySec is the new min_delay in seconds. assertUpdateMinDelayWithinCap
	// in the proposal generator rejects values above MaxTimelockScheduleDelay.
	NewMinDelaySec uint64 `yaml:"newMinDelaySec" json:"newMinDelaySec"`
}

var UpdateMinDelayMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "timelock", "update_min_delay"),
	semver.MustParse("0.1.0"),
	"Update a timelock's global min_delay via mcms::mcms_timelock_update_min_delay",
	updateMinDelayHandler,
)

var updateMinDelayHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UpdateMinDelayMCMSInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	if input.McmsPackageID == "" {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("mcmsPackageID is required")
	}
	if input.TimelockObjID == "" {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("timelockObjID is required")
	}

	configurer := suisdk.NewTimelockConfigurer(input.McmsPackageID)
	result, err := configurer.UpdateDelay(b.GetContext(), input.TimelockObjID, input.NewMinDelaySec)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("encoding mcms_timelock_update_min_delay: %w", err)
	}

	tx, ok := result.RawData.(types.Transaction)
	if !ok {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("unexpected transaction type %T from TimelockConfigurer.UpdateDelay", result.RawData)
	}

	// The configurer serializes the inner callback data: a 32-byte timelock
	// address followed by a u64 new_min_delay. The call carries the inner
	// callback name because the timelock stores it verbatim and the
	// mcms::mcms_timelock_update_min_delay entry asserts it on execute; the
	// executor maps the inner name back to that entry when building the PTB.
	// The named constant also lets assertUpdateMinDelayWithinCap cap-check it.
	call := sui_ops.TransactionCall{
		PackageID:  input.McmsPackageID,
		Module:     mcmsModuleName,
		Function:   updateMinDelayFunctionName,
		Data:       tx.Data,
		StateObjID: input.TimelockObjID,
		TypeArgs:   []string{},
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of mcms_timelock_update_min_delay; returning call for proposal",
			"timelockObjID", input.TimelockObjID, "newMinDelaySec", input.NewMinDelaySec)
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{
			PackageId: input.McmsPackageID,
			Call:      call,
		}, nil
	}

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf(
		"direct execution of mcms_timelock_update_min_delay is not supported by this op; route it through a timelock proposal",
	)
}
