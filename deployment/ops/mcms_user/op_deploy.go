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
		deps.SuiRPC,
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

type MCMSUserFunctionOneInput struct {
	McmsUserPackageID        string `json:"mcmsUserPackageID"`
	McmsUserOwnerCapObjectID string `json:"mcmsUserOwnerCapObjectID"`
	McmsRegistryObjectID     string `json:"mcmsRegistryObjectID"`
	McmsUserDataObjectID     string `json:"mcmsUserDataObjectID"`
	Arg1                     string `json:"arg1"`
	Arg2                     []byte `json:"arg2"`
}

var functionOneHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSUserFunctionOneInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	contract, err := mcmsuser.NewMCMSUser(input.McmsUserPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	encodedCall, err := contract.MCMSUser().Encoder().FunctionOne(bind.Object{Id: input.McmsUserDataObjectID}, bind.Object{Id: input.McmsUserOwnerCapObjectID}, input.Arg1, input.Arg2)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to encode RemovePackageId call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.McmsUserDataObjectID)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of MCMS User Function One as per no Signer provided", "packageId", input.McmsUserPackageID)
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{
			Digest:    "",
			PackageId: input.McmsUserPackageID,
			Objects:   cld_ops.EmptyInput{},
			Call:      call,
		}, nil
	}

	suiTx, err := contract.MCMSUser().Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to call set_config on mcms: %w", err)
	}

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    suiTx.Digest,
		PackageId: input.McmsUserPackageID,
	}, err
}

var FunctionOneOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms_user", "mcms_user", "function_one"),
	semver.MustParse("0.1.0"),
	"Function one in the MCMS User contract",
	functionOneHandler,
)
