package burnminttokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployBurnMintTokenPoolObjects struct {
	OwnerCapObjectID string
	StateObjectID    string
}

type DeployBurnMintTokenPoolOutput struct {
	BurnMintTPPackageID string
	Objects             DeployBurnMintTokenPoolObjects
}

type DeployAndInitBurnMintTokenPoolInput struct {
	BurnMintTokenPoolDeployInput
	// init
	CoinObjectTypeArg      string
	CCIPObjectRefObjectID  string
	CoinMetadataObjectID   string
	TreasuryCapObjectID    string
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
				CoinObjectTypeArg:      input.CoinObjectTypeArg,
				BurnMintPackageID:      deployReport.Output.PackageID,
				StateObjectID:          input.CCIPObjectRefObjectID,
				CoinMetadataObjectID:   input.CoinMetadataObjectID,
				TreasuryCapObjectID:    input.TreasuryCapObjectID,
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
				BurnMintPackageID:            deployReport.Output.PackageID,
				CoinObjectTypeArg:            input.CoinObjectTypeArg,
				StateObjectID:                initReport.Output.Objects.StateObjectID,
				OwnerCap:                     initReport.Output.Objects.OwnerCapObjectID,
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
				BurnMintPackageID:    deployReport.Output.PackageID,
				CoinObjectTypeArg:    input.CoinObjectTypeArg,
				StateObjectID:        initReport.Output.Objects.StateObjectID,
				OwnerCap:             initReport.Output.Objects.OwnerCapObjectID,
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

		return DeployBurnMintTokenPoolOutput{
			BurnMintTPPackageID: deployReport.Output.PackageID,
			Objects: DeployBurnMintTokenPoolObjects{
				OwnerCapObjectID: initReport.Output.Objects.OwnerCapObjectID,
				StateObjectID:    initReport.Output.Objects.StateObjectID,
			},
		}, nil
	},
)
