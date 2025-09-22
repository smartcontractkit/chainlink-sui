package mcmsuserops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	mcmsuser "github.com/smartcontractkit/chainlink-sui/bindings/packages/mcms/mcms_user"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type MCMSUserRegisterEntrypointInput struct {
	McmsUserPackageID        string `json:"mcmsUserPackageID"`
	McmsUserOwnerCapObjectID string `json:"mcmsUserOwnerCapObjectID"`
	McmsRegistryObjectID     string `json:"mcmsRegistryObjectID"`
	McmsUserDataObjectID     string `json:"mcmsUserDataObjectID"`
}

var registerEntrypointHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSUserRegisterEntrypointInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	contract, err := mcmsuser.NewMCMSUser(input.McmsUserPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	suiTx, err := contract.MCMSUser().RegisterMcmsEntrypoint(
		b.GetContext(),
		opts,
		bind.Object{Id: input.McmsUserOwnerCapObjectID},
		bind.Object{Id: input.McmsRegistryObjectID},
		bind.Object{Id: input.McmsUserDataObjectID},
	)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to call set_config on mcms: %w", err)
	}

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    suiTx.Digest,
		PackageId: input.McmsUserPackageID,
	}, err
}

var RegisterMCMSEntrypointOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms_user", "mcms_user", "register_entrypoint"),
	semver.MustParse("0.1.0"),
	"Register entrypoint in the MCMS User contract",
	registerEntrypointHandler,
)
