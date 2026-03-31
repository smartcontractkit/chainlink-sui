package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	module_mcms_deployer "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_deployer"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type RegisterCCIPUpgradeCapMcmsInput struct {
	CCIPPackageId         string
	UpgradeCapObjectId    string
	RegistryObjectId      string
	DeployerStateObjectId string
	McmsPackageId         string
}

type RegisterCCIPUpgradeCapMcmsObjects struct{}

var registerCCIPUpgradeCapMcmsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RegisterCCIPUpgradeCapMcmsInput) (output sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	encodedCall, err := contract.Encoder().McmsRegisterUpgradeCap(
		bind.Object{Id: input.UpgradeCapObjectId},
		bind.Object{Id: input.RegistryObjectId},
		bind.Object{Id: input.DeployerStateObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{}, fmt.Errorf("failed to encode McmsRegisterUpgradeCap: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.DeployerStateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping McmsRegisterUpgradeCap execution (no signer)",
			"ccipPackageId", input.CCIPPackageId,
			"upgradeCapObjectId", input.UpgradeCapObjectId,
		)
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   RegisterCCIPUpgradeCapMcmsObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{}, fmt.Errorf("failed to execute McmsRegisterUpgradeCap: %w", err)
	}

	deployer, err := module_mcms_deployer.NewMcmsDeployer(input.McmsPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{}, fmt.Errorf("failed to create mcms deployer binding: %w", err)
	}

	registered, err := deployer.DevInspect().HasUpgradeCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.DeployerStateObjectId},
		input.CCIPPackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{}, fmt.Errorf("failed to verify upgrade cap registration: %w", err)
	}
	if !registered {
		return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{}, fmt.Errorf(
			"upgrade cap not found in MCMS deployer state for package %s after registration tx %s",
			input.CCIPPackageId, tx.Digest,
		)
	}

	b.Logger.Infow("Registered CCIP UpgradeCap with MCMS deployer",
		"ccipPackageId", input.CCIPPackageId,
		"digest", tx.Digest,
	)

	return sui_ops.OpTxResult[RegisterCCIPUpgradeCapMcmsObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   RegisterCCIPUpgradeCapMcmsObjects{},
		Call:      call,
	}, nil
}

var RegisterCCIPUpgradeCapMcmsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "register_upgrade_cap_mcms"),
	semver.MustParse("0.1.0"),
	"Registers the CCIP package UpgradeCap with mcms::mcms_deployer (required before MCMS-authorized upgrades)",
	registerCCIPUpgradeCapMcmsHandler,
)
