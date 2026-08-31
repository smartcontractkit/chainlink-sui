package lockreleasetokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type ConfigureLockReleaseTokenPoolOutput struct {
	TokenSymbol string
	Objects     DeployLockReleaseTokenPoolObjects
	Reports     []cld_ops.Report[any, any]
}

type ConfigureLockReleaseTokenPoolInput struct {
	// pool object ids for an already-deployed lock release token pool
	TokenPoolPkgID    string
	StateObjectId     string
	OwnerCap          string
	CoinObjectTypeArg string

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

// ConfigureLockReleaseTokenPoolSequence configures an already-deployed lock
// release token pool by applying chain updates and setting rate limiter configs.
// It reuses the direct OwnerCap-gated ops, which execute EOA-direct when a
// Signer is provided and emit an MCMS-encoded call when Signer is nil.
var ConfigureLockReleaseTokenPoolSequence = cld_ops.NewSequence(
	"sui-configure-lock-release-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Configures an already-deployed lock release token pool",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input ConfigureLockReleaseTokenPoolInput) (ConfigureLockReleaseTokenPoolOutput, error) {
		seqReports := make([]cld_ops.Report[any, any], 0)

		report, err := cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolApplyChainUpdatesOp,
			deps,
			LockReleaseTokenPoolApplyChainUpdatesInput{
				LockReleasePackageId:         input.TokenPoolPkgID,
				CoinObjectTypeArg:            input.CoinObjectTypeArg,
				StateObjectId:                input.StateObjectId,
				OwnerCap:                     input.OwnerCap,
				RemoteChainSelectorsToRemove: input.RemoteChainSelectorsToRemove,
				RemoteChainSelectorsToAdd:    input.RemoteChainSelectorsToAdd,
				RemotePoolAddressesToAdd:     input.RemotePoolAddressesToAdd,
				RemoteTokenAddressesToAdd:    input.RemoteTokenAddressesToAdd,
			},
		)
		if err != nil {
			return ConfigureLockReleaseTokenPoolOutput{}, err
		}
		seqReports = append(seqReports, report.ToGenericReport())

		report2, err := cld_ops.ExecuteOperation(
			env,
			LockReleaseTokenPoolSetChainRateLimiterOp,
			deps,
			LockReleaseTokenPoolSetChainRateLimiterInput{
				LockReleasePackageId: input.TokenPoolPkgID,
				CoinObjectTypeArg:    input.CoinObjectTypeArg,
				StateObjectId:        input.StateObjectId,
				OwnerCap:             input.OwnerCap,
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
			return ConfigureLockReleaseTokenPoolOutput{}, err
		}
		seqReports = append(seqReports, report2.ToGenericReport())

		return ConfigureLockReleaseTokenPoolOutput{
			Objects: DeployLockReleaseTokenPoolObjects{
				OwnerCapObjectId: input.OwnerCap,
				StateObjectId:    input.StateObjectId,
			},
			Reports: seqReports,
		}, nil
	},
)
