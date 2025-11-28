package burnminttokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type ConfigureBurnMintTokenPoolObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type ConfigureBurnMintTokenPoolOutput struct {
	TokenSymbol string
	Objects     DeployBurnMintTokenPoolObjects
	Reports     []cld_ops.Report[any, any]
}

type ConfigureBurnMintTokenPoolInput struct {
	BurnMintTokenPoolDeployInput
	// init
	TokenPoolPkgID         string
	TokenPoolStateObjectID string
	TokenOwnerCapID        string
	CoinObjectTypeArg      string

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

var ConfigureBurnMintTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-burn-mint-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial burn mint token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input ConfigureBurnMintTokenPoolInput) (ConfigureBurnMintTokenPoolOutput, error) {

		seqReports := make([]cld_ops.Report[any, any], 0)
		report, err := cld_ops.ExecuteOperation(
			env,
			BurnMintTokenPoolApplyChainUpdatesOp,
			deps,
			BurnMintTokenPoolApplyChainUpdatesInput{
				BurnMintPackageId:            input.TokenPoolPkgID,
				CoinObjectTypeArg:            input.CoinObjectTypeArg,
				StateObjectId:                input.TokenPoolStateObjectID,
				OwnerCap:                     input.TokenOwnerCapID,
				RemoteChainSelectorsToRemove: input.RemoteChainSelectorsToRemove,
				RemoteChainSelectorsToAdd:    input.RemoteChainSelectorsToAdd,
				RemotePoolAddressesToAdd:     input.RemotePoolAddressesToAdd,
				RemoteTokenAddressesToAdd:    input.RemoteTokenAddressesToAdd,
			},
		)
		if err != nil {
			return ConfigureBurnMintTokenPoolOutput{}, err
		}
		seqReports = append(seqReports, report.ToGenericReport())

		report2, err := cld_ops.ExecuteOperation(
			env,
			BurnMintTokenPoolSetChainRateLimiterOp,
			deps,
			BurnMintTokenPoolSetChainRateLimiterInput{
				BurnMintPackageId:    input.TokenPoolPkgID,
				CoinObjectTypeArg:    input.CoinObjectTypeArg,
				StateObjectId:        input.TokenPoolStateObjectID,
				OwnerCap:             input.TokenOwnerCapID,
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
			return ConfigureBurnMintTokenPoolOutput{}, err
		}
		seqReports = append(seqReports, report2.ToGenericReport())

		for i, chainSelector := range input.RemoteChainSelectors {
			report, err := cld_ops.ExecuteOperation(
				env,
				BurnMintTokenPoolAddRemotePoolOp,
				deps,
				BurnMintTokenPoolAddRemotePoolInput{
					BurnMintTokenPoolPackageId: input.TokenPoolPkgID,
					CoinObjectTypeArg:          input.CoinObjectTypeArg,
					StateObjectId:              input.TokenPoolStateObjectID,
					OwnerCap:                   input.TokenOwnerCapID,
					RemoteChainSelector:        chainSelector,
					RemotePoolAddress:          input.RemotePoolAddressesToAdd[i][0], // one address at a time
				},
			)
			if err != nil {
				return ConfigureBurnMintTokenPoolOutput{}, err
			}
			seqReports = append(seqReports, report.ToGenericReport())
		}

		return ConfigureBurnMintTokenPoolOutput{
			Objects: DeployBurnMintTokenPoolObjects{
				OwnerCapObjectId: input.TokenOwnerCapID,
				StateObjectId:    input.TokenPoolStateObjectID,
			},
			Reports: seqReports,
		}, nil
	},
)
