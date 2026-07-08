package mcmsops

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/mcms/sdk"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	modulemcms "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type MCMSSetConfigInput struct {
	ChainSelector uint64 `yaml:"chainSelector"`
	// MCMS related
	McmsPackageID string `yaml:"mcmsPackageID"`
	OwnerCap      string `yaml:"ownerCap"`
	McmsObjectID  string `yaml:"mcmsObjectID"`
	// TimelockObjectID enables the defensive on-chain min_delay check (F5/F8
	// Interpretation B): before submitting a new signer set we read the
	// current min_delay and refuse to touch a timelock whose state has already
	// been driven out of the safe range. Optional — leave empty to skip the
	// read (kept as an escape hatch for callers that have not yet deployed
	// their MCMS state alongside a timelock).
	TimelockObjectID string `yaml:"timelockObjectID,omitempty"`
	// Timelock related
	Role suisdk.TimelockRole `yaml:"role"`
	// Config related
	Config    types.Config `yaml:"config"`
	ClearRoot bool         `yaml:"clearRoot"`
}

var SetConfigMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "mcms", "set_config"),
	semver.MustParse("0.1.0"),
	"Set config in the MCMS contract",
	setConfigMcmsHandler,
)

var setConfigMcmsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSSetConfigInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcms, err := modulemcms.NewMcms(input.McmsPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	// Interpretation B: read the current on-chain min_delay before rotating
	// signers so we do not hand a fresh signer set to a timelock whose
	// scheduling contract is already unusable. Complements
	// assertUpdateMinDelayWithinCap (which blocks outbound bricking writes) by
	// catching pre-existing bad state that a compromised BYPASSER could have
	// left behind. Skipped when TimelockObjectID is empty.
	if input.TimelockObjectID != "" {
		// DevInspect only needs an address, not a real signature. In
		// proposal-only mode deps.Signer is nil, so fall back to a
		// DevInspectSigner so the check runs regardless of execution mode.
		inspectOpts := *opts
		if inspectOpts.Signer == nil {
			inspectOpts.Signer = signer.NewDevInspectSigner("0x0")
		}
		currentMinDelay, delayErr := mcms.DevInspect().TimelockMinDelay(
			b.GetContext(),
			&inspectOpts,
			bind.Object{Id: input.TimelockObjectID},
		)
		if delayErr != nil {
			return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf(
				"failed to read current timelock min_delay for defensive check (F5/F8): %w", delayErr,
			)
		}
		if err := utils.AssertMinDelayWithinCap(currentMinDelay); err != nil {
			return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf(
				"refusing set_config on timelock %s: %w",
				input.TimelockObjectID, err,
			)
		}
		if currentMinDelay == 0 {
			b.Logger.Warnw(
				"F5: set_config invoked while timelock min_delay is 0; the bootstrap zero-delay window is still open, so any batch this signer set later schedules can be executed with delay=0 until an mcms_timelock_update_min_delay op lands",
				"role", input.Role,
				"chainSelector", input.ChainSelector,
				"timelockObjectID", input.TimelockObjectID,
			)
		}
	}

	chainID, err := cselectors.SuiChainIdFromSelector(input.ChainSelector)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}
	chainIDBig := new(big.Int).SetUint64(chainID)
	groupQuorum, groupParents, signerAddresses, signerGroups, err := sdk.ExtractSetConfigInputs(&input.Config)

	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("unable to extract set config inputs: %w", err)
	}
	signers := make([][]byte, len(signerAddresses))
	for i, addr := range signerAddresses {
		signers[i] = addr.Bytes()
	}

	if !input.ClearRoot {
		b.Logger.Warnw(
			"F7: set_config invoked with ClearRoot=false; if this changes the signer set, the previous root and its remaining [pre_op_count, post_op_count) window survive and the old (or newly-installed) signers can still consume it",
			"role", input.Role,
			"chainSelector", input.ChainSelector,
			"newSignerCount", len(signers),
		)
	}

	encodedCall, err := mcms.Encoder().SetConfig(
		bind.Object{Id: input.OwnerCap},
		bind.Object{Id: input.McmsObjectID},
		input.Role.Byte(),
		chainIDBig,
		signers,
		signerGroups,
		groupQuorum[:],
		groupParents[:],
		input.ClearRoot,
	)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to encode SetConfig call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.McmsObjectID)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on StateObject as per no Signer provided")
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{
			Digest:    "",
			PackageId: input.McmsPackageID,
			Call:      call,
		}, nil
	}

	tx, err := mcms.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Call:      call,
		Digest:    tx.Digest,
		PackageId: input.McmsPackageID,
	}, err
}
