package mcmsuserops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mcmsuser "github.com/smartcontractkit/chainlink-sui/bindings/packages/mcms/mcms_user"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployMCMSUserInput struct {
	McmsPackageID     string
	McmsOwnerObjectID string
}

type DeployMCMSUserObjects struct {
	McmsUserDataObjectID     string
	McmsUserOwnerCapObjectID string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployMCMSUserInput) (output sui_ops.OpTxResult[DeployMCMSUserObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcmsPackage, tx, err := mcmsuser.PublishMCMSUser(
		b.GetContext(),
		opts,
		deps.Client,
		input.McmsPackageID,
		input.McmsOwnerObjectID,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployMCMSUserObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	mcmsUserDataObject, err1 := bind.FindObjectIdFromPublishTx(*tx, "mcms_user", "UserData")
	mcmsUserOwnerCapObject, err2 := bind.FindObjectIdFromPublishTx(*tx, "mcms_user", "OwnerCap")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[DeployMCMSUserObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployMCMSUserObjects]{
		Digest:    tx.Digest,
		PackageId: mcmsPackage.Address(),
		Objects: DeployMCMSUserObjects{
			McmsUserDataObjectID:     mcmsUserDataObject,
			McmsUserOwnerCapObjectID: mcmsUserOwnerCapObject,
		},
	}, err
}

var DeployMCMSUserOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms_user", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the MCMS User contract",
	deployHandler,
)
