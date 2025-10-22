package ownershipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	usdctokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_usdc_token_pool"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

// ContractType represents the different contract types that support ownership transfer to MCMS
type ContractType string

const (
	ContractTypeStateObject          ContractType = "state_object"
	ContractTypeOnRamp               ContractType = "onramp"
	ContractTypeOffRamp              ContractType = "offramp"
	ContractTypeRouter               ContractType = "router"
	ContractTypeManagedToken         ContractType = "managed_token"
	ContractTypeBurnMintTokenPool    ContractType = "burn_mint_token_pool"
	ContractTypeLockReleaseTokenPool ContractType = "lock_release_token_pool"
	ContractTypeUsdcTokenPool        ContractType = "usdc_token_pool"
	ContractTypeManagedTokenPool     ContractType = "managed_token_pool"
	ContractTypeMCMS                 ContractType = "mcms"
)

type ExecuteOwnershipTransferToMcmsSeqInput struct {
	// Specific input for each contract type - only include the ones you want to transfer
	MCMS                 *mcmsops.MCMSExecuteTransferOwnershipInput                                       `json:"mcms,omitempty"`
	StateObject          *ccipops.ExecuteOwnershipTransferToMcmsStateObjectInput                          `json:"state_object,omitempty"`
	OnRamp               *onrampops.ExecuteOwnershipTransferToMcmsOnRampInput                             `json:"onramp,omitempty"`
	OffRamp              *offrampops.ExecuteOwnershipTransferToMcmsOffRampInput                           `json:"offramp,omitempty"`
	Router               *routerops.ExecuteOwnershipTransferToMcmsRouterInput                             `json:"router,omitempty"`
	ManagedToken         *managedtokenops.ExecuteOwnershipTransferToMcmsManagedTokenInput                 `json:"managed_token,omitempty"`
	BurnMintTokenPool    *burnminttokenpoolops.ExecuteOwnershipTransferToMcmsBurnMintTokenPoolInput       `json:"burn_mint_token_pool,omitempty"`
	LockReleaseTokenPool *lockreleasetokenpoolops.ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolInput `json:"lock_release_token_pool,omitempty"`
	UsdcTokenPool        *usdctokenpoolops.ExecuteOwnershipTransferToMcmsUsdcTokenPoolInput               `json:"usdc_token_pool,omitempty"`
	ManagedTokenPool     *managedtokenpoolops.ExecuteOwnershipTransferToMcmsManagedTokenPoolInput         `json:"managed_token_pool,omitempty"`
}

type ExecuteOwnershipTransferToMcmsSeqOutput struct {
	// Map of contract type to execution results
	Results map[ContractType]string // Transaction digests
}

var ExecuteOwnershipTransferToMcmsSequence = cld_ops.NewSequence(
	"sui-execute-ownership-transfer-to-mcms-seq",
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for specified CCIP contracts",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsSeqInput) (ExecuteOwnershipTransferToMcmsSeqOutput, error) {
		results := make(map[ContractType]string)

		// Execute MCMS ownership transfer if provided
		if input.MCMS != nil {
			report, err := cld_ops.ExecuteOperation(env, mcmsops.MCMSExecuteTransferOwnershipOp, deps, *input.MCMS)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for MCMS: %w", err)
			}
			results[ContractTypeMCMS] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for MCMS",
				"packageId", input.MCMS.McmsPackageID,
				"to", input.MCMS.McmsPackageID,
				"digest", report.Output.Digest)
		}
		// Execute StateObject ownership transfer if provided
		if input.StateObject != nil {
			report, err := cld_ops.ExecuteOperation(env, ccipops.ExecuteOwnershipTransferToMcmsStateObjectOp, deps, *input.StateObject)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for StateObject: %w", err)
			}
			results[ContractTypeStateObject] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for StateObject",
				"packageId", input.StateObject.CCIPPackageId,
				"to", input.StateObject.To,
				"digest", report.Output.Digest)
		}

		// Execute OnRamp ownership transfer if provided
		if input.OnRamp != nil {
			report, err := cld_ops.ExecuteOperation(env, onrampops.ExecuteOwnershipTransferToMcmsOnRampOp, deps, *input.OnRamp)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for OnRamp: %w", err)
			}
			results[ContractTypeOnRamp] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for OnRamp",
				"packageId", input.OnRamp.OnRampPackageId,
				"to", input.OnRamp.To,
				"digest", report.Output.Digest)
		}

		// Execute OffRamp ownership transfer if provided
		if input.OffRamp != nil {
			report, err := cld_ops.ExecuteOperation(env, offrampops.ExecuteOwnershipTransferToMcmsOffRampOp, deps, *input.OffRamp)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for OffRamp: %w", err)
			}
			results[ContractTypeOffRamp] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for OffRamp",
				"packageId", input.OffRamp.OffRampPackageId,
				"to", input.OffRamp.To,
				"digest", report.Output.Digest)
		}

		// Execute Router ownership transfer if provided
		if input.Router != nil {
			report, err := cld_ops.ExecuteOperation(env, routerops.ExecuteOwnershipTransferToMcmsRouterOp, deps, *input.Router)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for Router: %w", err)
			}
			results[ContractTypeRouter] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for Router",
				"packageId", input.Router.RouterPackageId,
				"to", input.Router.To,
				"digest", report.Output.Digest)
		}

		// Execute ManagedToken ownership transfer if provided
		if input.ManagedToken != nil {
			report, err := cld_ops.ExecuteOperation(env, managedtokenops.ExecuteOwnershipTransferToMcmsManagedTokenOp, deps, *input.ManagedToken)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for ManagedToken: %w", err)
			}
			results[ContractTypeManagedToken] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for ManagedToken",
				"packageId", input.ManagedToken.ManagedTokenPackageId,
				"to", input.ManagedToken.To,
				"digest", report.Output.Digest)
		}

		// Execute BurnMintTokenPool ownership transfer if provided
		if input.BurnMintTokenPool != nil {
			report, err := cld_ops.ExecuteOperation(env, burnminttokenpoolops.ExecuteOwnershipTransferToMcmsBurnMintTokenPoolOp, deps, *input.BurnMintTokenPool)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for BurnMintTokenPool: %w", err)
			}
			results[ContractTypeBurnMintTokenPool] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for BurnMintTokenPool",
				"packageId", input.BurnMintTokenPool.BurnMintTokenPoolPackageId,
				"to", input.BurnMintTokenPool.To,
				"digest", report.Output.Digest)
		}

		// Execute LockReleaseTokenPool ownership transfer if provided
		if input.LockReleaseTokenPool != nil {
			report, err := cld_ops.ExecuteOperation(env, lockreleasetokenpoolops.ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolOp, deps, *input.LockReleaseTokenPool)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for LockReleaseTokenPool: %w", err)
			}
			results[ContractTypeLockReleaseTokenPool] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for LockReleaseTokenPool",
				"packageId", input.LockReleaseTokenPool.LockReleaseTokenPoolPackageId,
				"to", input.LockReleaseTokenPool.To,
				"digest", report.Output.Digest)
		}

		// Execute UsdcTokenPool ownership transfer if provided
		if input.UsdcTokenPool != nil {
			report, err := cld_ops.ExecuteOperation(env, usdctokenpoolops.ExecuteOwnershipTransferToMcmsUsdcTokenPoolOp, deps, *input.UsdcTokenPool)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for UsdcTokenPool: %w", err)
			}
			results[ContractTypeUsdcTokenPool] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for UsdcTokenPool",
				"packageId", input.UsdcTokenPool.UsdcTokenPoolPackageId,
				"to", input.UsdcTokenPool.To,
				"digest", report.Output.Digest)
		}

		// Execute ManagedTokenPool ownership transfer if provided
		if input.ManagedTokenPool != nil {
			report, err := cld_ops.ExecuteOperation(env, managedtokenpoolops.ExecuteOwnershipTransferToMcmsManagedTokenPoolOp, deps, *input.ManagedTokenPool)
			if err != nil {
				return ExecuteOwnershipTransferToMcmsSeqOutput{}, fmt.Errorf("failed to execute ownership transfer for ManagedTokenPool: %w", err)
			}
			results[ContractTypeManagedTokenPool] = report.Output.Digest
			env.Logger.Infow("Successfully executed ownership transfer to MCMS for ManagedTokenPool",
				"packageId", input.ManagedTokenPool.ManagedTokenPoolPackageId,
				"to", input.ManagedTokenPool.To,
				"digest", report.Output.Digest)
		}

		return ExecuteOwnershipTransferToMcmsSeqOutput{
			Results: results,
		}, nil
	},
)
