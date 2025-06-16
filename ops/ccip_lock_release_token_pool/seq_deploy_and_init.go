package lockreleasetokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployLockReleaseTokenPoolObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type DeployLockReleaseTokenPoolOutput struct {
	CCIPPackageId string
	Objects       DeployLockReleaseTokenPoolObjects
}

type DeployAndInitLockReleaseTokenPoolInput struct {
	LockReleaseTokenPoolDeployInput
	// init
	CCIPObjectRefObjectId  string
	CoinMetadataObjectId   string
	TreasuryCapObjectId    string
	TokenPoolPackageId     string
	TokenPoolAdministrator string
	Rebalancer             string
	// apply chain updates
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][][]byte
	RemoteTokenAddressesToAdd    [][]byte
	// set chain rate limiter configs
	RemoteChainSelectors []uint64
	OutboundIsEnableds   []bool
	OutboundCapacities   []uint64
	OutboundRates        []uint64
	InboundIsEnableds    []bool
	InboundCapacities    []uint64
	InboundRates         []uint64
}

var DeployAndInitLockReleaseTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-lock-release-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial lock release token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitLockReleaseTokenPoolInput) (DeployLockReleaseTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPTokenPoolOp, deps, input.LockReleaseTokenPoolDeployInput)
		if err != nil {
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		initReport, err := cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolInitializeOp,
			deps,
			LockReleaseTokenPoolInitializeInput{
				CCIPPackageId:          deployReport.Output.PackageId,
				StateObjectId:          input.CCIPObjectRefObjectId,
				CoinMetadataObjectId:   input.CoinMetadataObjectId,
				TreasuryCapObjectId:    input.TreasuryCapObjectId,
				TokenPoolPackageId:     input.TokenPoolPackageId,
				TokenPoolAdministrator: input.TokenPoolAdministrator,
				Rebalancer:             input.Rebalancer,
			},
		)
		if err != nil {
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolApplyChainUpdatesOp,
			deps,
			LockReleaseTokenPoolApplyChainUpdatesInput{
				CCIPPackageId:                deployReport.Output.PackageId,
				StateObjectId:                initReport.Output.Objects.StateObjectId,
				OwnerCap:                     initReport.Output.Objects.OwnerCapObjectId,
				RemoteChainSelectorsToRemove: input.RemoteChainSelectorsToRemove,
				RemoteChainSelectorsToAdd:    input.RemoteChainSelectorsToAdd,
				RemotePoolAddressesToAdd:     input.RemotePoolAddressesToAdd,
				RemoteTokenAddressesToAdd:    input.RemoteTokenAddressesToAdd,
			},
		)
		if err != nil {
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolSetChainRateLimiterOp,
			deps,
			LockReleaseTokenPoolSetChainRateLimiterInput{
				CCIPPackageId:        deployReport.Output.PackageId,
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
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		return DeployLockReleaseTokenPoolOutput{
			CCIPPackageId: deployReport.Output.PackageId,
			Objects: DeployLockReleaseTokenPoolObjects{
				OwnerCapObjectId: initReport.Output.Objects.OwnerCapObjectId,
				StateObjectId:    initReport.Output.Objects.StateObjectId,
			},
		}, nil
	},
)
