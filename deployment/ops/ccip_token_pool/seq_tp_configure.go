package tokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
)

// ConfigureAllTokenPoolsInput contains configuration for all supported token pool types.
// This replicates the logic from the DeployTPAndConfigure changeset, allowing deployment
// of multiple token pool types in a single sequence operation.
type ConfigureAllTokenPoolsInput struct {
	// Common configuration
	SuiChainSelector uint64
	TokenPoolTypes   []string // supported: "bnm", "lnr", "managed"

	// Token pool specific inputs
	ManagedTPInput     managedtokenpoolops.ConfigureManagedTokenPoolInput
	LockReleaseTPInput lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolInput
	BurnMintTpInput    burnminttokenpoolops.ConfigureBurnMintTokenPoolInput
}

// ConfigureAllTokenPoolsOutput contains results for all deployed token pools
type ConfigureAllTokenPoolsOutput struct {
	burnminttokenpoolops.DeployBurnMintTokenPoolOutput
	lockreleasetokenpoolops.DeployLockReleaseTokenPoolOutput
	managedtokenpoolops.DeployManagedTokenPoolOutput

	Reports []cld_ops.Report[any, any]
}

// ConfigureAllTokenPoolsSequence provides a unified sequence for deploying and
// initializing multiple token pool types. This sequence replicates the core logic
// of the DeployTPAndConfigure changeset by conditionally deploying burn mint,
// lock release, and managed token pools based on the provided configuration.
var ConfigureAllTokenPoolsSequence = cld_ops.NewSequence(
	"sui-deploy-and-init-all-token-pools-seq",
	semver.MustParse("0.1.0"),
	"Deploys and initializes multiple token pool types based on configuration",
	ConfigureAllTokenPoolsSeq,
)

func ConfigureAllTokenPoolsSeq(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input ConfigureAllTokenPoolsInput) (ConfigureAllTokenPoolsOutput, error) {
	output := ConfigureAllTokenPoolsOutput{}
	for _, tokenPoolType := range input.TokenPoolTypes {
		switch tokenPoolType {
		case "bnm":
			report, err := cld_ops.ExecuteSequence(env, burnminttokenpoolops.ConfigureBurnMintTokenPoolSequence, deps, input.BurnMintTpInput)
			if err != nil {
				return ConfigureAllTokenPoolsOutput{}, fmt.Errorf("failed to deploy burn mint token pool: %w", err)
			}

			symbol, err := getSymbol(env, deps, input.BurnMintTpInput.CoinObjectTypeArg)
			if err != nil {
				return ConfigureAllTokenPoolsOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
			}
			output.DeployBurnMintTokenPoolOutput.TokenSymbol = symbol
			output.Reports = append(output.Reports, report.Output.Reports...)
		case "lnr":
			// todo
		case "managed":
			report, err := cld_ops.ExecuteSequence(env, managedtokenpoolops.ConfigureManagedTokenPoolSequence, deps, input.ManagedTPInput)
			if err != nil {
				return ConfigureAllTokenPoolsOutput{}, fmt.Errorf("failed to deploy burn mint token pool: %w", err)
			}

			symbol, err := getSymbol(env, deps, input.ManagedTPInput.CoinObjectTypeArg)
			if err != nil {
				return ConfigureAllTokenPoolsOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
			}
			output.DeployManagedTokenPoolOutput.TokenSymbol = symbol
			output.Reports = append(output.Reports, report.Output.Reports...)
		default:
			return ConfigureAllTokenPoolsOutput{}, fmt.Errorf("unsupported token pool type: %s", tokenPoolType)
		}
	}

	return output, nil
}
