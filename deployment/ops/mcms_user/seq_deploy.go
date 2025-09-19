package mcmsuserops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployMCMSUserSeqInput struct {
	McmsPackageID        string `json:"mcmsPackageID" yaml:"mcmsPackageID"`
	McmsRegistryObjectID string `json:"mcmsRegistryObjectID" yaml:"mcmsRegistryObjectID"`
}

var DeployMCMSUserSequence = cld_ops.NewSequence(
	"sui-deploy-mcms-user-seq",
	semver.MustParse("0.1.0"),
	"Deploys and registers mcms entrypoint for the MCMS user contract",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployMCMSUserSeqInput) (sui_ops.OpTxResult[DeployMCMSUserObjects], error) {

		signerAddress, err := deps.Signer.GetAddress()
		if err != nil {
			return sui_ops.OpTxResult[DeployMCMSUserObjects]{}, fmt.Errorf("failed to get signer address: %w", err)
		}
		// Deploy MCMS User contract first
		deployInput := DeployMCMSUserInput{
			McmsPackageID:     input.McmsPackageID,
			McmsOwnerObjectID: signerAddress,
		}

		deployReport, err := cld_ops.ExecuteOperation(env, DeployMCMSUserOp, deps, deployInput)
		if err != nil {
			return sui_ops.OpTxResult[DeployMCMSUserObjects]{}, fmt.Errorf("failed to deploy MCMS user contract: %w", err)
		}

		// Register MCMS entrypoint
		registerInput := MCMSUserRegisterEntrypointInput{
			McmsUserPackageID:        deployReport.Output.PackageId,
			McmsUserOwnerCapObjectID: deployReport.Output.Objects.McmsUserOwnerCapObjectID,
			McmsRegistryObjectID:     input.McmsRegistryObjectID,
			McmsUserDataObjectID:     deployReport.Output.Objects.McmsUserDataObjectID,
		}

		registerReport, err := cld_ops.ExecuteOperation(env, RegisterMCMSEntrypointOp, deps, registerInput)
		if err != nil {
			return sui_ops.OpTxResult[DeployMCMSUserObjects]{}, fmt.Errorf("failed to register mcms entrypoint: %w", err)
		}

		env.Logger.Infow("Successfully registered MCMS entrypoint", "digest", registerReport.Output.Digest)

		return sui_ops.OpTxResult[DeployMCMSUserObjects]{
			Digest:    registerReport.Output.Digest, // Use the register transaction digest as the final result
			PackageId: deployReport.Output.PackageId,
			Objects:   deployReport.Output.Objects,
		}, nil
	},
)
