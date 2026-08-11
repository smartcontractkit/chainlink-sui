package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mcms_deployer "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_deployer"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// RegisterUpgradeCapInput describes the parameters for registering an UpgradeCap with MCMS
type RegisterUpgradeCapInput struct {
	MCMSPackageID         string `json:"mcmsPackageID" validate:"required"`
	UpgradeCapObjectID    string `json:"upgradeCapObjectID" validate:"required"`
	RegistryObjectID      string `json:"registryObjectID" validate:"required"`
	DeployerStateObjectID string `json:"deployerStateObjectID" validate:"required"`
	PackageName           string `json:"packageName" validate:"required"`
}

var registerUpgradeCapHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RegisterUpgradeCapInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	// Create the MCMS deployer contract instance
	mcmsDeployer, err := module_mcms_deployer.NewMcmsDeployer(input.MCMSPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to create MCMS deployer contract: %w", err)
	}

	// Call RegisterUpgradeCap on the MCMS deployer
	// Parameter order: state, registry, upgradeCap
	tx, err := mcmsDeployer.RegisterUpgradeCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.DeployerStateObjectID},
		bind.Object{Id: input.RegistryObjectID},
		bind.Object{Id: input.UpgradeCapObjectID},
	)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to register upgrade cap for package %s: %w", input.PackageName, err)
	}

	b.Logger.Infow("Registered upgrade cap with MCMS",
		"packageName", input.PackageName,
		"upgradeCapID", input.UpgradeCapObjectID,
		"digest", tx.Digest,
	)

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    tx.Digest,
		PackageId: input.MCMSPackageID,
	}, nil
}

var RegisterUpgradeCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "deployer", "register_upgrade_cap"),
	semver.MustParse("0.1.0"),
	"Register a package's upgrade capability with the MCMS deployer for future upgrades",
	registerUpgradeCapHandler,
)
