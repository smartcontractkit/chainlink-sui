package burnminttokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployBurnMintTokenPoolObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type DeployBurnMintTokenPoolOutput struct {
	CCIPPackageId string
	Objects       DeployBurnMintTokenPoolObjects
}

type DeployAndInitBurnMintTokenPoolInput struct {
	BurnMintTokenPoolDeployInput
	// init
	CoinObjectTypeArg      string
	CCIPObjectRefObjectId  string
	CoinMetadataObjectId   string
	TreasuryCapObjectId    string
	TokenPoolPackageId     string
	TokenPoolAdministrator string
	LockOrBurnParams       []string
	ReleaseOrMintParams    []string
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
				BurnMintPackageId:      deployReport.Output.PackageId,
				StateObjectId:          input.CCIPObjectRefObjectId,
				CoinMetadataObjectId:   input.CoinMetadataObjectId,
				TreasuryCapObjectId:    input.TreasuryCapObjectId,
				TokenPoolPackageId:     input.TokenPoolPackageId,
				TokenPoolAdministrator: input.TokenPoolAdministrator,
				LockOrBurnParams:       input.LockOrBurnParams,
				ReleaseOrMintParams:    input.ReleaseOrMintParams,
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
				OwnerCap:                     initReport.Output.Objects.OwnerCapObjectId,
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
				OwnerCap:             initReport.Output.Objects.OwnerCapObjectId,
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
			CCIPPackageId: deployReport.Output.PackageId,
			Objects: DeployBurnMintTokenPoolObjects{
				OwnerCapObjectId: initReport.Output.Objects.OwnerCapObjectId,
				StateObjectId:    initReport.Output.Objects.StateObjectId,
			},
		}, nil
	},
)
