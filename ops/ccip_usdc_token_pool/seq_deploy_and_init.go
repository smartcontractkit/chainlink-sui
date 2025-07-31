package usdctokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployUSDCTokenPoolObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type DeployUSDCTokenPoolOutput struct {
	CCIPPackageId string
	Objects       DeployUSDCTokenPoolObjects
}

type DeployAndInitUSDCTokenPoolInput struct {
	USDCTokenPoolDeployInput
	// init
	CoinObjectTypeArg      string
	CCIPObjectRefObjectId  string
	OwnerCapObjectId       string
	CoinMetadataObjectId   string
	LocalDomainIdentifier  uint32
	TokenPoolPackageId     string
	TokenPoolAdministrator string
	// set domains
	RemoteChainSelectors    []uint64
	RemoteDomainIdentifiers []uint32
	AllowedRemoteCallers    [][]byte
	DomainsEnableds         []bool
	// apply chain updates
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][]string
	RemoteTokenAddressesToAdd    []string
	// set chain rate limiter configs
	RateLimiterRemoteChainSelectors []uint64
	OutboundIsEnableds              []bool
	OutboundCapacities              []uint64
	OutboundRates                   []uint64
	InboundIsEnableds               []bool
	InboundCapacities               []uint64
	InboundRates                    []uint64
	ClockObjectId                   string
}

var DeployAndInitUSDCTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-usdc-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial USDC token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitUSDCTokenPoolInput) (DeployUSDCTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPUSDCTokenPoolOp, deps, input.USDCTokenPoolDeployInput)
		if err != nil {
			return DeployUSDCTokenPoolOutput{}, err
		}

		initReport, err := cld_ops.ExecuteOperation(
			env,
			USDCTokenPoolInitializeOp,
			deps,
			USDCTokenPoolInitializeInput{
				USDCTokenPoolPackageId: deployReport.Output.PackageId,
				CoinObjectTypeArg:      input.CoinObjectTypeArg,
				StateObjectId:          input.CCIPObjectRefObjectId,
				OwnerCapObjectId:       input.OwnerCapObjectId,
				CoinMetadataObjectId:   input.CoinMetadataObjectId,
				LocalDomainIdentifier:  input.LocalDomainIdentifier,
				TokenPoolPackageId:     input.TokenPoolPackageId,
				TokenPoolAdministrator: input.TokenPoolAdministrator,
			},
		)
		if err != nil {
			return DeployUSDCTokenPoolOutput{}, err
		}

		// Set domains configuration
		_, err = cld_ops.ExecuteOperation(
			env,
			USDCTokenPoolSetDomainsOp,
			deps,
			USDCTokenPoolSetDomainsInput{
				USDCTokenPoolPackageId:  deployReport.Output.PackageId,
				CoinObjectTypeArg:       input.CoinObjectTypeArg,
				StateObjectId:           initReport.Output.Objects.StateObjectId,
				OwnerCap:                initReport.Output.Objects.OwnerCapObjectId,
				RemoteChainSelectors:    input.RemoteChainSelectors,
				RemoteDomainIdentifiers: input.RemoteDomainIdentifiers,
				AllowedRemoteCallers:    input.AllowedRemoteCallers,
				Enableds:                input.DomainsEnableds,
			},
		)
		if err != nil {
			return DeployUSDCTokenPoolOutput{}, err
		}

		// Apply chain updates
		_, err = cld_ops.ExecuteOperation(
			env,
			USDCTokenPoolApplyChainUpdatesOp,
			deps,
			USDCTokenPoolApplyChainUpdatesInput{
				USDCTokenPoolPackageId:       deployReport.Output.PackageId,
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
			return DeployUSDCTokenPoolOutput{}, err
		}

		// Set chain rate limiter configurations
		_, err = cld_ops.ExecuteOperation(
			env,
			USDCTokenPoolSetChainRateLimiterOp,
			deps,
			USDCTokenPoolSetChainRateLimiterInput{
				USDCTokenPoolPackageId: deployReport.Output.PackageId,
				CoinObjectTypeArg:      input.CoinObjectTypeArg,
				StateObjectId:          initReport.Output.Objects.StateObjectId,
				OwnerCap:               initReport.Output.Objects.OwnerCapObjectId,
				ClockObjectId:          input.ClockObjectId,
				RemoteChainSelectors:   input.RateLimiterRemoteChainSelectors,
				OutboundIsEnableds:     input.OutboundIsEnableds,
				OutboundCapacities:     input.OutboundCapacities,
				OutboundRates:          input.OutboundRates,
				InboundIsEnableds:      input.InboundIsEnableds,
				InboundCapacities:      input.InboundCapacities,
				InboundRates:           input.InboundRates,
			},
		)
		if err != nil {
			return DeployUSDCTokenPoolOutput{}, err
		}

		return DeployUSDCTokenPoolOutput{
			CCIPPackageId: deployReport.Output.PackageId,
			Objects: DeployUSDCTokenPoolObjects{
				OwnerCapObjectId: initReport.Output.Objects.OwnerCapObjectId,
				StateObjectId:    initReport.Output.Objects.StateObjectId,
			},
		}, nil
	},
)
