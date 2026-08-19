package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_malicious_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_malicious_token_pool/malicious_token_pool"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// DeployMaliciousTokenPoolObjects are the objects created by publishing and
// initializing the TEST-only malicious token pool.
type DeployMaliciousTokenPoolObjects struct {
	OwnerCapObjectId   string
	StateObjectId      string
	UpgradeCapObjectId string
}

// DeployMaliciousTokenPoolInput publishes the TEST-only malicious token pool
// package and initializes it. The pool is a burn-mint pool whose release_or_mint
// declares an extra &mut Coin<SUI> drain tail; initialize forwards
// ReleaseOrMintParams into register_pool as the pool's release_or_mint_params,
// which is the attacker-controlled field the offchain ReleaseOrMintParams
// ownership guard defends.
type DeployMaliciousTokenPoolInput struct {
	// publish
	CCIPPackageId    string
	MCMSAddress      string
	FastMcmsAddress  string
	MCMSOwnerAddress string
	// initialize
	CoinObjectTypeArg      string
	CCIPObjectRefObjectId  string
	CoinMetadataObjectId   string
	TreasuryCapObjectId    string
	TokenPoolAdministrator string
	// 89201: attacker-controlled release_or_mint_params, e.g. executor-transmitter-owned SUI coin ids.
	ReleaseOrMintParams []string
}

var deployMaliciousTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployMaliciousTokenPoolInput) (output sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	signerAddr, err := opts.Signer.GetAddress()
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, err
	}

	// Compile the TEST-only malicious token pool package.
	artifact, err := bind.CompilePackage(contracts.CCIPMaliciousTokenPool, map[string]string{
		"ccip":                      input.CCIPPackageId,
		"ccip_malicious_token_pool": "0x0",
		"mcms":                      input.MCMSAddress,
		"fast_mcms":                 input.FastMcmsAddress,
		"mcms_owner":                input.MCMSOwnerAddress,
		"signer":                    signerAddr,
	}, false, deps.SuiRPC)
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, fmt.Errorf("failed to compile malicious token pool package: %w", err)
	}

	// Publish the package. init creates the OwnerCap and the package UpgradeCap.
	packageId, publishTx, err := bind.PublishPackage(
		b.GetContext(),
		opts,
		deps.Client,
		bind.PublishRequest{
			CompiledModules: artifact.Modules,
			Dependencies:    artifact.Dependencies,
		},
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, fmt.Errorf("failed to publish malicious token pool package: %w", err)
	}

	ownerCapObj, err := bind.FindObjectIdFromPublishTx(*publishTx, "ownable", "OwnerCap")
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, fmt.Errorf("failed to find OwnerCap object ID: %w", err)
	}

	upgradeCapObj, err := bind.FindObjectIdFromPublishTx(*publishTx, "package", "UpgradeCap")
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, fmt.Errorf("failed to find UpgradeCap object ID: %w", err)
	}

	// Initialize the pool, forwarding the attacker ReleaseOrMintParams into
	// register_pool. The pool self-registers here, same as all Sui pools.
	contract, err := module_malicious_token_pool.NewMaliciousTokenPool(packageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, fmt.Errorf("failed to create malicious token pool contract: %w", err)
	}

	initTx, err := contract.Initialize(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: ownerCapObj},
		bind.Object{Id: input.CCIPObjectRefObjectId},
		bind.Object{Id: input.CoinMetadataObjectId},
		bind.Object{Id: input.TreasuryCapObjectId},
		input.TokenPoolAdministrator,
		input.ReleaseOrMintParams,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, fmt.Errorf("failed to execute malicious token pool initialization: %w", err)
	}

	stateObj, err := bind.FindObjectIdFromPublishTx(*initTx, "malicious_token_pool", "BurnMintTokenPoolState")
	if err != nil {
		return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{}, fmt.Errorf("failed to find BurnMintTokenPoolState object ID: %w", err)
	}

	return sui_ops.OpTxResult[DeployMaliciousTokenPoolObjects]{
		Digest:    initTx.Digest,
		PackageId: packageId,
		Objects: DeployMaliciousTokenPoolObjects{
			OwnerCapObjectId:   ownerCapObj,
			StateObjectId:      stateObj,
			UpgradeCapObjectId: upgradeCapObj,
		},
	}, nil
}

var DeployCCIPMaliciousTokenPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-malicious-token-pool", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys and initializes the TEST-only CCIP malicious token pool package (ReleaseOrMintParams guard E2E smoke test)",
	deployMaliciousTokenPoolHandler,
)
