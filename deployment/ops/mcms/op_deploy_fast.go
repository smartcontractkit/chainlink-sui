package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/mcms"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployFastMCMSObjects = DeployMCMSObjects

var deployFastHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployFastMCMSObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcmsPackage, tx, err := mcms.PublishFastMCMS(
		b.GetContext(),
		opts,
		deps.Client,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployFastMCMSObjects]{}, err
	}

	mcmsObject, err1 := bind.FindObjectIdFromPublishTx(*tx, "mcms", "MultisigState")
	timelockObj, err2 := bind.FindObjectIdFromPublishTx(*tx, "mcms", "Timelock")
	depState, err3 := bind.FindObjectIdFromPublishTx(*tx, "mcms_deployer", "DeployerState")
	reg, err4 := bind.FindObjectIdFromPublishTx(*tx, "mcms_registry", "Registry")
	acc, err5 := bind.FindObjectIdFromPublishTx(*tx, "mcms_account", "AccountState")
	ownCap, err6 := bind.FindObjectIdFromPublishTx(*tx, "mcms_account", "OwnerCap")

	if err1 != nil {
		return sui_ops.OpTxResult[DeployFastMCMSObjects]{}, fmt.Errorf("failed to find MultisigState object ID in publish tx: %w", err1)
	}
	if err2 != nil {
		return sui_ops.OpTxResult[DeployFastMCMSObjects]{}, fmt.Errorf("failed to find Timelock object ID in publish tx: %w", err2)
	}
	if err3 != nil {
		return sui_ops.OpTxResult[DeployFastMCMSObjects]{}, fmt.Errorf("failed to find DeployerState object ID in publish tx: %w", err3)
	}
	if err4 != nil {
		return sui_ops.OpTxResult[DeployFastMCMSObjects]{}, fmt.Errorf("failed to find Registry object ID in publish tx: %w", err4)
	}
	if err5 != nil {
		return sui_ops.OpTxResult[DeployFastMCMSObjects]{}, fmt.Errorf("failed to find AccountState object ID in publish tx: %w", err5)
	}
	if err6 != nil {
		return sui_ops.OpTxResult[DeployFastMCMSObjects]{}, fmt.Errorf("failed to find OwnerCap object ID in publish tx: %w", err6)
	}

	return sui_ops.OpTxResult[DeployFastMCMSObjects]{
		Digest:    tx.Digest,
		PackageId: mcmsPackage.Address(),
		Objects: DeployFastMCMSObjects{
			McmsMultisigStateObjectId:   mcmsObject,
			TimelockObjectId:            timelockObj,
			McmsDeployerStateObjectId:   depState,
			McmsRegistryObjectId:        reg,
			McmsAccountStateObjectId:    acc,
			McmsAccountOwnerCapObjectId: ownCap,
		},
	}, err
}

var DeployFastMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "fast_package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the fast MCMS contract package",
	deployFastHandler,
)
