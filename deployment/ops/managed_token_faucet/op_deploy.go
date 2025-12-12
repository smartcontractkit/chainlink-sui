package managedtokenfaucetops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployManagedTokenFaucetInput struct {
	ManagedTokenPackageId string
	MCMSAddress           string
	MCMSOwnerAddress      string
}

type DeployManagedTokenFaucetObjects struct {
	UpgradeCapObjectId string
}

var deployManagedTokenFaucetOp = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployManagedTokenFaucetInput) (output sui_ops.OpTxResult[DeployManagedTokenFaucetObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[DeployManagedTokenFaucetObjects]{}, err
	}

	artifact, err := bind.CompilePackage(contracts.ManagedTokenFaucet, map[string]string{
		"managed_token_faucet": "0x0",
		"managed_token":        input.ManagedTokenPackageId,
		"mcms":                 input.MCMSAddress,
		"mcms_owner":           input.MCMSOwnerAddress,
		"signer":               signerAddr,
	}, false, deps.SuiRPC)
	if err != nil {
		return sui_ops.OpTxResult[DeployManagedTokenFaucetObjects]{}, fmt.Errorf("failed to compile managed token faucet package: %w", err)
	}

	packageID, tx, err := bind.PublishPackage(b.GetContext(), opts, deps.Client, bind.PublishRequest{
		CompiledModules: artifact.Modules,
		Dependencies:    artifact.Dependencies,
	})
	if err != nil {
		return sui_ops.OpTxResult[DeployManagedTokenFaucetObjects]{}, fmt.Errorf("failed to publish managed token faucet package: %w", err)
	}

	upgradeCapID, err := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	if err != nil {
		return sui_ops.OpTxResult[DeployManagedTokenFaucetObjects]{}, fmt.Errorf("failed to find UpgradeCap object ID: %w", err)
	}

	return sui_ops.OpTxResult[DeployManagedTokenFaucetObjects]{
		Digest:    tx.Digest,
		PackageId: packageID,
		Objects: DeployManagedTokenFaucetObjects{
			UpgradeCapObjectId: upgradeCapID,
		},
	}, nil
}

var DeployManagedTokenFaucetOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-faucet", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP managed token faucet package",
	deployManagedTokenFaucetOp,
)
