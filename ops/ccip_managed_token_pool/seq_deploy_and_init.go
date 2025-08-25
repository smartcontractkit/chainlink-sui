package managedtokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type SeqDeployAndInitManagedTokenPoolInput struct {
	// deploy
	CCIPPackageID          string
	CCIPTokenPoolPackageID string
	ManagedTokenPackageID  string
	MCMSAddress            string
	MCMSOwnerAddress       string
	// initialize
	CoinObjectTypeArg         string
	CCIPObjectRefObjectID     string
	ManagedTokenStateObjectID string
	ManagedTokenOwnerCapID    string
	CoinMetadataObjectID      string
	MintCapObjectID           string
	TokenPoolAdministrator    string
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

type DeployManagedTokenPoolObjects struct {
	OwnerCapObjectID string
	StateObjectID    string
}

type DeployManagedTokenPoolOutput struct {
	ManagedTPPackageID string
	Objects            DeployManagedTokenPoolObjects
}

var DeployAndInitManagedTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-managed-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial managed token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input SeqDeployAndInitManagedTokenPoolInput) (DeployManagedTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPManagedTokenPoolOp, deps, ManagedTokenPoolDeployInput{
			CCIPPackageID:          input.CCIPPackageID,
			CCIPTokenPoolPackageID: input.CCIPTokenPoolPackageID,
			ManagedTokenPackageID:  input.ManagedTokenPackageID,
			MCMSAddress:            input.MCMSAddress,
			MCMSOwnerAddress:       input.MCMSOwnerAddress,
		})
		if err != nil {
			return DeployManagedTokenPoolOutput{}, err
		}

		initReport, err := cld_ops.ExecuteOperation(
			env,
			ManagedTokenPoolInitializeOp,
			deps,
			ManagedTokenPoolInitializeInput{
				ManagedTokenPoolPackageID: deployReport.Output.PackageId,
				CoinObjectTypeArg:         input.CoinObjectTypeArg,
				CCIPObjectRefObjectID:     input.CCIPObjectRefObjectID,
				ManagedTokenStateObjectID: input.ManagedTokenStateObjectID,
				ManagedTokenOwnerCapID:    input.ManagedTokenOwnerCapID,
				CoinMetadataObjectID:      input.CoinMetadataObjectID,
				MintCapObjectID:           input.MintCapObjectID,
				TokenPoolAdministrator:    input.TokenPoolAdministrator,
			},
		)
		if err != nil {
			return DeployManagedTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			ManagedTokenPoolApplyChainUpdatesOp,
			deps,
			ManagedTokenPoolApplyChainUpdatesInput{
				ManagedTokenPoolPackageID:    deployReport.Output.PackageId,
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
			return DeployManagedTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			ManagedTokenPoolSetChainRateLimiterOp,
			deps,
			ManagedTokenPoolSetChainRateLimiterInput{
				ManagedTokenPoolPackageID: deployReport.Output.PackageId,
				CoinObjectTypeArg:         input.CoinObjectTypeArg,
				StateObjectID:             initReport.Output.Objects.StateObjectID,
				OwnerCap:                  initReport.Output.Objects.OwnerCapObjectID,
				RemoteChainSelectors:      input.RemoteChainSelectors,
				OutboundIsEnableds:        input.OutboundIsEnableds,
				OutboundCapacities:        input.OutboundCapacities,
				OutboundRates:             input.OutboundRates,
				InboundIsEnableds:         input.InboundIsEnableds,
				InboundCapacities:         input.InboundCapacities,
				InboundRates:              input.InboundRates,
			},
		)
		if err != nil {
			return DeployManagedTokenPoolOutput{}, err
		}

		return DeployManagedTokenPoolOutput{
			ManagedTPPackageID: deployReport.Output.PackageId,
			Objects: DeployManagedTokenPoolObjects{
				OwnerCapObjectID: initReport.Output.Objects.OwnerCapObjectID,
				StateObjectID:    initReport.Output.Objects.StateObjectID,
			},
		}, nil
	},
)
