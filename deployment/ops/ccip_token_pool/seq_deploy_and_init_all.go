package tokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
)

// DeployAndInitAllTokenPoolsInput contains configuration for all supported token pool types.
// This replicates the logic from the DeployTPAndConfigure changeset, allowing deployment
// of multiple token pool types in a single sequence operation.
type DeployAndInitAllTokenPoolsInput struct {
	// Common configuration
	SuiChainSelector uint64
	TokenPoolTypes   []deployment.TokenPoolType

	// Token pool specific inputs
	ManagedTPInput     managedtokenpoolops.DeployAndInitManagedTokenPoolInput
	LockReleaseTPInput lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolInput
	BurnMintTpInput    burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput
}

// DeployAndInitAllTokenPoolsOutput contains results for all deployed token pools
type DeployAndInitAllTokenPoolsOutput struct {
	burnminttokenpoolops.DeployBurnMintTokenPoolOutput
	lockreleasetokenpoolops.DeployLockReleaseTokenPoolOutput
	managedtokenpoolops.DeployManagedTokenPoolOutput
}

// DeployAndInitAllTokenPoolsSequence provides a unified sequence for deploying and
// initializing multiple token pool types. This sequence replicates the core logic
// of the DeployTPAndConfigure changeset by conditionally deploying burn mint,
// lock release, and managed token pools based on the provided configuration.
var DeployAndInitAllTokenPoolsSequence = cld_ops.NewSequence(
	"sui-deploy-and-init-all-token-pools-seq",
	semver.MustParse("0.1.0"),
	"Deploys and initializes multiple token pool types based on configuration",
	deployAndInitAllTokenPoolsSeq,
)

func deployAndInitAllTokenPoolsSeq(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployAndInitAllTokenPoolsInput) (DeployAndInitAllTokenPoolsOutput, error) {
	output := DeployAndInitAllTokenPoolsOutput{}
	for _, tokenPoolType := range input.TokenPoolTypes {
		switch tokenPoolType {
		case deployment.TokenPoolTypeBurnMint:
			burnMintReport, err := cld_ops.ExecuteSequence(env, burnminttokenpoolops.DeployAndInitBurnMintTokenPoolSequence, deps, input.BurnMintTpInput)
			if err != nil {
				return DeployAndInitAllTokenPoolsOutput{}, fmt.Errorf("failed to deploy burn mint token pool: %w", err)
			}
			output.DeployBurnMintTokenPoolOutput = burnMintReport.Output
			symbol, err := getSymbol(env, deps, input.BurnMintTpInput.CoinObjectTypeArg)
			if err != nil {
				return DeployAndInitAllTokenPoolsOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
			}
			output.DeployBurnMintTokenPoolOutput.TokenSymbol = symbol
		case deployment.TokenPoolTypeLockRelease:
			lockReleaseReport, err := cld_ops.ExecuteSequence(env, lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolSequence, deps, input.LockReleaseTPInput)
			if err != nil {
				return DeployAndInitAllTokenPoolsOutput{}, fmt.Errorf("failed to deploy lock release token pool: %w", err)
			}
			output.DeployLockReleaseTokenPoolOutput = lockReleaseReport.Output
			symbol, err := getSymbol(env, deps, input.LockReleaseTPInput.CoinObjectTypeArg)
			if err != nil {
				return DeployAndInitAllTokenPoolsOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
			}
			output.DeployLockReleaseTokenPoolOutput.TokenSymbol = symbol
		case deployment.TokenPoolTypeManaged:
			managedReport, err := cld_ops.ExecuteSequence(env, managedtokenpoolops.DeployAndInitManagedTokenPoolSequence, deps, input.ManagedTPInput)
			if err != nil {
				return DeployAndInitAllTokenPoolsOutput{}, fmt.Errorf("failed to deploy managed token pool: %w", err)
			}
			output.DeployManagedTokenPoolOutput = managedReport.Output
			symbol, err := getSymbol(env, deps, input.ManagedTPInput.CoinObjectTypeArg)
			if err != nil {
				return DeployAndInitAllTokenPoolsOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
			}
			output.DeployManagedTokenPoolOutput.TokenSymbol = symbol
		default:
			return DeployAndInitAllTokenPoolsOutput{}, fmt.Errorf("unsupported token pool type: %s", tokenPoolType)
		}
	}

	return output, nil
}

func getSymbol(env cld_ops.Bundle, deps sui_ops.OpTxDeps, coinObjectTypeArg string) (string, error) {
	symbolReport, err := cld_ops.ExecuteOperation(env, coin_ops.GetCoinSymbolOp, deps, coinObjectTypeArg)
	if err != nil {
		return "", err
	}
	return symbolReport.Output.Symbol, nil
}
