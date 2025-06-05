package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/mcms"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployMCMSObjects struct {
	// MCMS
	McmsMultisigStateObjectId string
	TimelockObjectId          string
	// MCMS Deployer
	McmsDeployerObjectId string
	// MCMS Registry
	McmsRegistryObjectId string
	// MCMS Account
	McmsAccountStateObjectId    string
	McmsAccountOwnerCapObjectId string
}

var handler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (output sui_ops.OpTxResult[DeployMCMSObjects], err error) {
	mcmsPackage, tx, err := mcms.PublishMCMS(
		b.GetContext(),
		deps.GetTxOpts(),
		deps.Signer,
		deps.Client,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployMCMSObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	mcmsObject, err1 := bind.FindObjectIdFromPublishTx(*tx, "mcms", "MultisigState")
	timelockObj, err2 := bind.FindObjectIdFromPublishTx(*tx, "mcms", "Timelock")
	depState, err3 := bind.FindObjectIdFromPublishTx(*tx, "mcms_deployer", "DeployerState")
	reg, err4 := bind.FindObjectIdFromPublishTx(*tx, "mcms_registry", "Registry")
	acc, err5 := bind.FindObjectIdFromPublishTx(*tx, "mcms_account", "AccountState")
	ownCap, err6 := bind.FindObjectIdFromPublishTx(*tx, "mcms_account", "OwnerCap")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil {
		return sui_ops.OpTxResult[DeployMCMSObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployMCMSObjects]{
		Digest:    tx.Digest.String(),
		PackageId: mcmsPackage.Address().String(),
		Objects: DeployMCMSObjects{
			McmsMultisigStateObjectId:   mcmsObject,
			TimelockObjectId:            timelockObj,
			McmsDeployerObjectId:        depState,
			McmsRegistryObjectId:        reg,
			McmsAccountStateObjectId:    acc,
			McmsAccountOwnerCapObjectId: ownCap,
		},
	}, err
}

var DeployMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the MCMS contract",
	handler,
)
