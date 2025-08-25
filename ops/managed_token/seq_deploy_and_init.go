package managedtokenops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployManagedTokenObjects struct {
	OwnerCapObjectID string
	StateObjectID    string
}

type DeployManagedTokenOutput struct {
	ManagedTokenPackageID string
	Objects               DeployManagedTokenObjects
}

type DeployAndInitManagedTokenInput struct {
	ManagedTokenDeployInput
	// init
	CoinObjectTypeArg   string
	TreasuryCapObjectID string
	DenyCapObjectID     string // Optional - can be empty
	// configure_new_minter
	MinterAddress string
	Allowance     uint64
	IsUnlimited   bool
}

var DeployAndInitManagedTokenSequence = cld_ops.NewSequence(
	"sui-deploy-managed-token-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial managed token configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitManagedTokenInput) (DeployManagedTokenOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPManagedTokenOp, deps, input.ManagedTokenDeployInput)
		if err != nil {
			return DeployManagedTokenOutput{}, err
		}

		initReport, err := cld_ops.ExecuteOperation(
			env,
			ManagedTokenInitializeOp,
			deps,
			ManagedTokenInitializeInput{
				ManagedTokenPackageID: deployReport.Output.PackageId,
				CoinObjectTypeArg:     input.CoinObjectTypeArg,
				TreasuryCapObjectID:   input.TreasuryCapObjectID,
				DenyCapObjectID:       input.DenyCapObjectID,
			},
		)
		if err != nil {
			return DeployManagedTokenOutput{}, err
		}

		// Configure a new minter if specified
		if input.MinterAddress != "" {
			_, err = cld_ops.ExecuteOperation(
				env,
				ManagedTokenConfigureNewMinterOp,
				deps,
				ManagedTokenConfigureNewMinterInput{
					ManagedTokenPackageID: deployReport.Output.PackageId,
					CoinObjectTypeArg:     input.CoinObjectTypeArg,
					StateObjectID:         initReport.Output.Objects.StateObjectID,
					OwnerCapObjectID:      initReport.Output.Objects.OwnerCapObjectID,
					MinterAddress:         input.MinterAddress,
					Allowance:             input.Allowance,
					IsUnlimited:           input.IsUnlimited,
				},
			)
			if err != nil {
				return DeployManagedTokenOutput{}, err
			}
		}

		return DeployManagedTokenOutput{
			ManagedTokenPackageID: deployReport.Output.PackageId,
			Objects: DeployManagedTokenObjects{
				OwnerCapObjectID: initReport.Output.Objects.OwnerCapObjectID,
				StateObjectID:    initReport.Output.Objects.StateObjectID,
			},
		}, nil
	},
)
