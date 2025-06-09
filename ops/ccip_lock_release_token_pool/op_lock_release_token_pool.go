package lockreleasetokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

// LRTP -- INITIALIZE
type LockReleaseTokenPoolInitializeObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type LockReleaseTokenPoolInitializeInput struct {
	CCIPPackageId        string
	StateObjectId        string
	CoinMetadataObjectId string
	TreasuryCapObjectId  string
	TokenPoolPackageId   string
	Rebalancer           string
}

var initLRTPHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolInitializeInput) (output sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	method := contract.Initialize(
		"",
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.CoinMetadataObjectId},
		bind.Object{Id: input.TreasuryCapObjectId},
		input.TokenPoolPackageId,
		input.Rebalancer,
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute lock release token pool initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "lock_release_token_pool", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "lock_release_token_pool", "LockReleaseTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects: LockReleaseTokenPoolInitializeObjects{
			OwnerCapObjectId: obj1,
			StateObjectId:    obj2,
		},
	}, err
}

var LockReleaseTokenPoolInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Lock Release Token Pool contract",
	initLRTPHandler,
)

// LRTP -- apply_chain_updates
type NoObjects struct {
}

type LockReleaseTokenPoolApplyChainUpdatesInput struct {
	CCIPPackageId                string
	StateObjectId                string
	OwnerCap                     string
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	// TODO: Provide a more human readable type for these
	RemotePoolAddressesToAdd  [][][]byte
	RemoteTokenAddressesToAdd [][]byte
}

var applyChainUpdates = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolApplyChainUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	method := contract.ApplyChainUpdates(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelectorsToRemove,
		input.RemoteChainSelectorsToAdd,
		input.RemotePoolAddressesToAdd,
		input.RemoteTokenAddressesToAdd,
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool apply chain updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
	}, err
}

var LockReleaseTokenPoolApplyChainUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "apply_chain_updates"),
	semver.MustParse("0.1.0"),
	"Applies chain updates in the CCIP Lock Release Token Pool contract",
	applyChainUpdates,
)

// LRTP -- set_chain_rate_limiter_configs
type LockReleaseTokenPoolSetChainRateLimiterInput struct {
	CCIPPackageId        string
	StateObjectId        string
	OwnerCap             string
	RemoteChainSelectors []uint64
	OutboundIsEnableds   []bool
	OutboundCapacities   []uint64
	OutboundRates        []uint64
	InboundIsEnableds    []bool
	InboundCapacities    []uint64
	InboundRates         []uint64
}

var setChainRateLimiterHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolSetChainRateLimiterInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	method := contract.SetChainRateLimiterConfigs(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		bind.Object{Id: "0x6"}, // Clock object
		input.RemoteChainSelectors,
		input.OutboundIsEnableds,
		input.OutboundCapacities,
		input.OutboundRates,
		input.InboundIsEnableds,
		input.InboundCapacities,
		input.InboundRates,
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool set configs rate limiter: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
	}, err
}

var LockReleaseTokenPoolSetChainRateLimiterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "set_chain_rate_limiter_configs"),
	semver.MustParse("0.1.0"),
	"Sets chain rate limiter configs in the CCIP Lock Release Token Pool contract",
	setChainRateLimiterHandler,
)

// LRTP -- provide_liquidity
type LockReleaseTokenPoolProviderLiquidityInput struct {
	CCIPPackageId string
	StateObjectId string
	Coin          string
}

var providerLiquidityHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolProviderLiquidityInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	method := contract.ProvideLiquidity(
		"",
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.Coin},
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool set configs rate limiter: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
	}, err
}

var LockReleaseTokenPoolProviderLiquidityOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "provide_liquidity"),
	semver.MustParse("0.1.0"),
	"Provide liquidity CCIP Lock Release Token Pool contract",
	providerLiquidityHandler,
)
