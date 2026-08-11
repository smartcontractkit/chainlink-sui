package burnminttokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployBurnMintTokenPoolObjects struct {
	OwnerCapObjectId   string
	StateObjectId      string
	UpgradeCapObjectId string
}

type DeployBurnMintTokenPoolOutput struct {
	TokenSymbol         string
	BurnMintTPPackageID string
	Objects             DeployBurnMintTokenPoolObjects
}

type DeployAndInitBurnMintTokenPoolInput struct {
	BurnMintTokenPoolDeployInput
	// init
	CoinObjectTypeArg      string
	CCIPObjectRefObjectId  string
	CoinMetadataObjectId   string
	TreasuryCapObjectId    string
	TokenPoolAdministrator string
	// apply chain updates
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][]string
	RemoteTokenAddressesToAdd    []string
	// set chain rate limiter configs
	RemoteChainSelectors []uint64
	OutboundIsEnableds   []bool
	OutboundCapacities   []uint64
	OutboundRates        []uint64
	InboundIsEnableds    []bool
	InboundCapacities    []uint64
	InboundRates         []uint64
}

var DeployAndInitBurnMintTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-burn-mint-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial burn mint token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitBurnMintTokenPoolInput) (DeployBurnMintTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPBurnMintTokenPoolOp, deps, input.BurnMintTokenPoolDeployInput)
		if err != nil {
			return DeployBurnMintTokenPoolOutput{}, err
		}

		initReport, err := cld_ops.ExecuteOperation(
			env,
			BurnMintTokenPoolInitializeOp,
			deps,
			BurnMintTokenPoolInitializeInput{
				BurnMintPackageId:      deployReport.Output.PackageId,
				OwnerCapObjectId:       deployReport.Output.Objects.OwnerCapObjectId,
				CoinObjectTypeArg:      input.CoinObjectTypeArg,
				StateObjectId:          input.CCIPObjectRefObjectId,
				CoinMetadataObjectId:   input.CoinMetadataObjectId,
				TreasuryCapObjectId:    input.TreasuryCapObjectId,
				TokenPoolAdministrator: input.TokenPoolAdministrator,
			},
		)
		if err != nil {
			return DeployBurnMintTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			BurnMintTokenPoolApplyChainUpdatesOp,
			deps,
			BurnMintTokenPoolApplyChainUpdatesInput{
				BurnMintPackageId:            deployReport.Output.PackageId,
				CoinObjectTypeArg:            input.CoinObjectTypeArg,
				StateObjectId:                initReport.Output.Objects.StateObjectId,
				OwnerCap:                     deployReport.Output.Objects.OwnerCapObjectId,
				RemoteChainSelectorsToRemove: input.RemoteChainSelectorsToRemove,
				RemoteChainSelectorsToAdd:    input.RemoteChainSelectorsToAdd,
				RemotePoolAddressesToAdd:     input.RemotePoolAddressesToAdd,
				RemoteTokenAddressesToAdd:    input.RemoteTokenAddressesToAdd,
			},
		)
		if err != nil {
			return DeployBurnMintTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			BurnMintTokenPoolSetChainRateLimiterOp,
			deps,
			BurnMintTokenPoolSetChainRateLimiterInput{
				BurnMintPackageId:    deployReport.Output.PackageId,
				CoinObjectTypeArg:    input.CoinObjectTypeArg,
				StateObjectId:        initReport.Output.Objects.StateObjectId,
				OwnerCap:             deployReport.Output.Objects.OwnerCapObjectId,
				RemoteChainSelectors: input.RemoteChainSelectors,
				OutboundIsEnableds:   input.OutboundIsEnableds,
				OutboundCapacities:   input.OutboundCapacities,
				OutboundRates:        input.OutboundRates,
				InboundIsEnableds:    input.InboundIsEnableds,
				InboundCapacities:    input.InboundCapacities,
				InboundRates:         input.InboundRates,
			},
		)
		if err != nil {
			return DeployBurnMintTokenPoolOutput{}, err
		}

		// transfer ownership to MCMS
		_, err = cld_ops.ExecuteOperation(
			env,
			TransferOwnershipBurnMintTokenPoolOp,
			deps,
			TransferOwnershipBurnMintTokenPoolInput{
				BurnMintTokenPoolPackageId: deployReport.Output.PackageId,
				TypeArgs:                   []string{input.CoinObjectTypeArg},
				StateObjectId:              initReport.Output.Objects.StateObjectId,
				OwnerCapObjectId:           deployReport.Output.Objects.OwnerCapObjectId,
				To:                         input.BurnMintTokenPoolDeployInput.MCMSAddress,
			},
		)

		return DeployBurnMintTokenPoolOutput{
			BurnMintTPPackageID: deployReport.Output.PackageId,
			Objects: DeployBurnMintTokenPoolObjects{
				OwnerCapObjectId:   deployReport.Output.Objects.OwnerCapObjectId,
				StateObjectId:      initReport.Output.Objects.StateObjectId,
				UpgradeCapObjectId: deployReport.Output.Objects.UpgradeCapObjectId,
			},
		}, nil
	},
)
