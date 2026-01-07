package managedtokenpoolops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployAndInitManagedTokenPoolInput struct {
	// deploy
	CCIPPackageId         string
	ManagedTokenPackageId string // ManagedToken
	MCMSAddress           string
	MCMSOwnerAddress      string
	// initialize
	CoinObjectTypeArg         string // CCIPBnM Token TypeArgs
	CCIPObjectRefObjectId     string // CCIP ObjectRef
	ManagedTokenStateObjectId string // ManagedToken
	ManagedTokenOwnerCapId    string // ManagedToken
	CoinMetadataObjectId      string // CCIPBnM Token
	MintCapObjectId           string // ManagedToken
	TokenPoolAdministrator    string // put yourself
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
	OwnerCapObjectId   string
	StateObjectId      string
	UpgradeCapObjectId string
}

type DeployManagedTokenPoolOutput struct {
	TokenSymbol        string
	ManagedTPPackageId string
	Objects            DeployManagedTokenPoolObjects
}

var DeployAndInitManagedTokenPoolSequence = cld_ops.NewSequence(
	"sui-deploy-managed-token-pool-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial managed token pool configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitManagedTokenPoolInput) (DeployManagedTokenPoolOutput, error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployCCIPManagedTokenPoolOp, deps, ManagedTokenPoolDeployInput{
			CCIPPackageId:         input.CCIPPackageId,
			ManagedTokenPackageId: input.ManagedTokenPackageId,
			MCMSAddress:           input.MCMSAddress,
			MCMSOwnerAddress:      input.MCMSOwnerAddress,
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
				OwnerCapObjectId:          deployReport.Output.Objects.OwnerCapObjectId,
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
				OwnerCap:                     deployReport.Output.Objects.OwnerCapObjectId,
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
				OwnerCap:                  deployReport.Output.Objects.OwnerCapObjectId,
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

		// init ownership transfer to MCMS
		_, err = cld_ops.ExecuteOperation(
			env,
			TransferOwnershipManagedTokenPoolOp,
			deps,
			TransferOwnershipManagedTokenPoolInput{
				ManagedTokenPoolPackageId: deployReport.Output.PackageId,
				TypeArgs:                  []string{input.CoinObjectTypeArg},
				StateObjectId:             initReport.Output.Objects.StateObjectId,
				OwnerCapObjectId:          deployReport.Output.Objects.OwnerCapObjectId,
				To:                        input.MCMSAddress,
			},
		)
		if err != nil {
			return DeployManagedTokenPoolOutput{}, err
		}

		return DeployManagedTokenPoolOutput{
			ManagedTPPackageId: deployReport.Output.PackageId,
			Objects: DeployManagedTokenPoolObjects{
				OwnerCapObjectId:   deployReport.Output.Objects.OwnerCapObjectId,
				StateObjectId:      initReport.Output.Objects.StateObjectId,
				UpgradeCapObjectId: deployReport.Output.Objects.UpgradeCapObjectId,
			},
		}, nil
	},
)
