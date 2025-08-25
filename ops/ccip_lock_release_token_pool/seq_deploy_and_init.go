package lockreleasetokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployLockReleaseTokenPoolObjects struct {
	OwnerCapObjectID string
	StateObjectID    string
}

type DeployLockReleaseTokenPoolOutput struct {
	LockReleaseTPPackageID string
	Objects                DeployLockReleaseTokenPoolObjects
}

type DeployAndInitLockReleaseTokenPoolInput struct {
	LockReleaseTokenPoolDeployInput
	// init
	CoinObjectTypeArg      string
	CCIPObjectRefObjectID  string
	CoinMetadataObjectID   string
	TreasuryCapObjectID    string
	TokenPoolAdministrator string
	Rebalancer             string
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

var DeployAndInitLockReleaseTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-lock-release-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial lock release token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitLockReleaseTokenPoolInput) (DeployLockReleaseTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPLockReleaseTokenPoolOp, deps, input.LockReleaseTokenPoolDeployInput)
		if err != nil {
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		initReport, err := cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolInitializeOp,
			deps,
			LockReleaseTokenPoolInitializeInput{
				CoinObjectTypeArg:      input.CoinObjectTypeArg,
				LockReleasePackageID:   deployReport.Output.PackageID,
				StateObjectID:          input.CCIPObjectRefObjectID,
				CoinMetadataObjectID:   input.CoinMetadataObjectID,
				TreasuryCapObjectID:    input.TreasuryCapObjectID,
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
				LockReleasePackageID:         deployReport.Output.PackageID,
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
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolSetChainRateLimiterOp,
			deps,
			LockReleaseTokenPoolSetChainRateLimiterInput{
				LockReleasePackageID: deployReport.Output.PackageID,
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
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		return DeployLockReleaseTokenPoolOutput{
			LockReleaseTPPackageID: deployReport.Output.PackageID,
			Objects: DeployLockReleaseTokenPoolObjects{
				OwnerCapObjectID: initReport.Output.Objects.OwnerCapObjectID,
				StateObjectID:    initReport.Output.Objects.StateObjectID,
			},
		}, nil
	},
)

// DEPLOY AND INIT BY CCIP ADMIN SEQUENCE
type DeployAndInitLockReleaseTokenPoolByCcipAdminInput struct {
	LockReleaseTokenPoolDeployInput
	// init by ccip admin
	CoinObjectTypeArg      string
	CCIPObjectRefObjectID  string
	CoinMetadataObjectID   string
	OwnerCapObjectID       string
	TokenPoolAdministrator string
	Rebalancer             string
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

var DeployAndInitLockReleaseTokenPoolByCcipAdminSequence = cld_ops.NewSequence(
	"sui-deploy-lock-release-token-pool-by-ccip-admin-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial lock release token pool configuration using CCIP admin",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitLockReleaseTokenPoolByCcipAdminInput) (DeployLockReleaseTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPLockReleaseTokenPoolOp, deps, input.LockReleaseTokenPoolDeployInput)
		if err != nil {
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		initReport, err := cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolInitializeByCcipAdminOp,
			deps,
			LockReleaseTokenPoolInitializeByCcipAdminInput{
				CoinObjectTypeArg:      input.CoinObjectTypeArg,
				LockReleasePackageID:   deployReport.Output.PackageID,
				StateObjectID:          input.CCIPObjectRefObjectID,
				CoinMetadataObjectID:   input.CoinMetadataObjectID,
				OwnerCapObjectID:       input.OwnerCapObjectID,
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
				LockReleasePackageID:         deployReport.Output.PackageID,
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
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolSetChainRateLimiterOp,
			deps,
			LockReleaseTokenPoolSetChainRateLimiterInput{
				LockReleasePackageID: deployReport.Output.PackageID,
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
			return DeployLockReleaseTokenPoolOutput{}, err
		}

		return DeployLockReleaseTokenPoolOutput{
			LockReleaseTPPackageID: deployReport.Output.PackageID,
			Objects: DeployLockReleaseTokenPoolObjects{
				OwnerCapObjectID: initReport.Output.Objects.OwnerCapObjectID,
				StateObjectID:    initReport.Output.Objects.StateObjectID,
			},
		}, nil
	},
)
