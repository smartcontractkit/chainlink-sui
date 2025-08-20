package managedtokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type SeqDeployAndInitManagedTokenPoolInput struct {
	// deploy
	CCIPPackageId          string
	CCIPTokenPoolPackageId string
	ManagedTokenPackageId  string
	MCMSAddress            string
	MCMSOwnerAddress       string
	// initialize
	CoinObjectTypeArg         string
	CCIPObjectRefObjectId     string
	ManagedTokenStateObjectId string
	ManagedTokenOwnerCapId    string
	CoinMetadataObjectId      string
	MintCapObjectId           string
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
	OwnerCapObjectId string
	StateObjectId    string
}

type DeployManagedTokenPoolOutput struct {
	ManagedTPPackageId string
	Objects            DeployManagedTokenPoolObjects
}

var DeployAndInitManagedTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-managed-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial managed token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input SeqDeployAndInitManagedTokenPoolInput) (DeployManagedTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPManagedTokenPoolOp, deps, ManagedTokenPoolDeployInput{
			CCIPPackageId:          input.CCIPPackageId,
			CCIPTokenPoolPackageId: input.CCIPTokenPoolPackageId,
			ManagedTokenPackageId:  input.ManagedTokenPackageId,
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
				ManagedTokenPoolPackageId: deployReport.Output.PackageId,
				CoinObjectTypeArg:         input.CoinObjectTypeArg,
				CCIPObjectRefObjectId:     input.CCIPObjectRefObjectId,
				ManagedTokenStateObjectId: input.ManagedTokenStateObjectId,
				ManagedTokenOwnerCapId:    input.ManagedTokenOwnerCapId,
				CoinMetadataObjectId:      input.CoinMetadataObjectId,
				MintCapObjectId:           input.MintCapObjectId,
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
				ManagedTokenPoolPackageId:    deployReport.Output.PackageId,
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
			return DeployManagedTokenPoolOutput{}, err
		}

		_, err = cld_ops.ExecuteOperation(
			env,
			ManagedTokenPoolSetChainRateLimiterOp,
			deps,
			ManagedTokenPoolSetChainRateLimiterInput{
				ManagedTokenPoolPackageId: deployReport.Output.PackageId,
				CoinObjectTypeArg:         input.CoinObjectTypeArg,
				StateObjectId:             initReport.Output.Objects.StateObjectId,
				OwnerCap:                  initReport.Output.Objects.OwnerCapObjectId,
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
			ManagedTPPackageId: deployReport.Output.PackageId,
			Objects: DeployManagedTokenPoolObjects{
				OwnerCapObjectId: initReport.Output.Objects.OwnerCapObjectId,
				StateObjectId:    initReport.Output.Objects.StateObjectId,
			},
		}, nil
	},
)
