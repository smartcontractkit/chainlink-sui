package managedtokenops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
)

type DeployManagedTokenObjects struct {
	OwnerCapObjectId   string
	StateObjectId      string
	MinterCapObjectId  string
	PublisherObjectId  string
	UpgradeCapObjectId string
}

type DeployManagedTokenOutput struct {
	ManagedTokenPackageId string
	TokenSymbol           string
	Objects               DeployManagedTokenObjects
}

type DeployAndInitManagedTokenInput struct {
	ManagedTokenDeployInput
	// init
	CoinObjectTypeArg   string
	TreasuryCapObjectId string
	DenyCapObjectId     string // Optional - can be empty
	PublisherObjectId   string
	// configure_new_minter
	MinterAddress string
	Allowance     uint64
	IsUnlimited   bool
	Source        string
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
				ManagedTokenPackageId: deployReport.Output.PackageId,
				CoinObjectTypeArg:     input.CoinObjectTypeArg,
				TreasuryCapObjectId:   input.TreasuryCapObjectId,
				DenyCapObjectId:       input.DenyCapObjectId,
				PublisherObjectId:     deployReport.Output.Objects.PublisherObjectId,
			},
		)
		if err != nil {
			return DeployManagedTokenOutput{}, err
		}

		minterObjectId := ""
		// Configure a new minter if specified
		if input.MinterAddress != "" {
			minterReport, err := cld_ops.ExecuteOperation(
				env,
				ManagedTokenConfigureNewMinterOp,
				deps,
				ManagedTokenConfigureNewMinterInput{
					ManagedTokenPackageId: deployReport.Output.PackageId,
					CoinObjectTypeArg:     input.CoinObjectTypeArg,
					StateObjectId:         initReport.Output.Objects.StateObjectId,
					OwnerCapObjectId:      initReport.Output.Objects.OwnerCapObjectId,
					MinterAddress:         input.MinterAddress,
					Allowance:             input.Allowance,
					IsUnlimited:           input.IsUnlimited,
					Source:                input.Source,
				},
			)
			if err != nil {
				return DeployManagedTokenOutput{}, err
			}

			minterObjectId = minterReport.Output.Objects.MinterCapObjectId

		}

		symbol, err := getSymbol(env, deps, input.CoinObjectTypeArg)
		if err != nil {
			return DeployManagedTokenOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
		}

		return DeployManagedTokenOutput{
			ManagedTokenPackageId: deployReport.Output.PackageId,
			TokenSymbol:           symbol,
			Objects: DeployManagedTokenObjects{
				OwnerCapObjectId:   initReport.Output.Objects.OwnerCapObjectId,
				StateObjectId:      initReport.Output.Objects.StateObjectId,
				PublisherObjectId:  deployReport.Output.Objects.PublisherObjectId,
				MinterCapObjectId:  minterObjectId,
				UpgradeCapObjectId: deployReport.Output.Objects.UpgradeCapObjectId,
			},
		}, nil
	},
)

func getSymbol(env cld_ops.Bundle, deps sui_ops.OpTxDeps, coinObjectTypeArg string) (string, error) {
	symbolReport, err := cld_ops.ExecuteOperation(env, coin_ops.GetCoinSymbolOp, deps, coinObjectTypeArg)
	if err != nil {
		return "", err
	}
	return symbolReport.Output.Symbol, nil
}
