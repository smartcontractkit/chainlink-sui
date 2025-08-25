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
	OwnerCapObjectID string
	StateObjectID    string
}

type LockReleaseTokenPoolInitializeInput struct {
	LockReleasePackageID   string
	CoinObjectTypeArg      string
	StateObjectID          string
	CoinMetadataObjectID   string
	TreasuryCapObjectID    string
	TokenPoolAdministrator string
	Rebalancer             string
}

var initLRTPHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolInitializeInput) (output sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleasePackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Initialize(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.CoinMetadataObjectID},
		bind.Object{Id: input.TreasuryCapObjectID},
		input.TokenPoolAdministrator,
		input.Rebalancer,
	)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute lock release token pool initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "lock_release_token_pool", "LockReleaseTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleasePackageID,
		Objects: LockReleaseTokenPoolInitializeObjects{
			OwnerCapObjectID: obj1,
			StateObjectID:    obj2,
		},
	}, err
}

var LockReleaseTokenPoolInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Lock Release Token Pool contract",
	initLRTPHandler,
)

// LRTP -- INITIALIZE BY CCIP ADMIN
type LockReleaseTokenPoolInitializeByCcipAdminInput struct {
	LockReleasePackageID   string
	CoinObjectTypeArg      string
	StateObjectID          string
	CoinMetadataObjectID   string
	OwnerCapObjectID       string
	TokenPoolAdministrator string
	Rebalancer             string
}

var initByCcipAdminLRTPHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolInitializeByCcipAdminInput) (output sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleasePackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.InitializeByCcipAdmin(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.OwnerCapObjectID},
		bind.Object{Id: input.CoinMetadataObjectID},
		input.TokenPoolAdministrator,
		input.Rebalancer,
	)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute lock release token pool initialization by ccip admin: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "lock_release_token_pool", "LockReleaseTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleasePackageID,
		Objects: LockReleaseTokenPoolInitializeObjects{
			OwnerCapObjectID: obj1,
			StateObjectID:    obj2,
		},
	}, err
}

var LockReleaseTokenPoolInitializeByCcipAdminOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "initialize_by_ccip_admin"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Lock Release Token Pool contract by CCIP admin",
	initByCcipAdminLRTPHandler,
)

// LRTP -- apply_chain_updates
type NoObjects struct {
}

type LockReleaseTokenPoolApplyChainUpdatesInput struct {
	LockReleasePackageID         string
	CoinObjectTypeArg            string
	StateObjectID                string
	OwnerCap                     string
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][]string
	RemoteTokenAddressesToAdd    []string
}

var applyChainUpdates = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolApplyChainUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleasePackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	// Convert [][]string to [][][]byte for RemotePoolAddressesToAdd
	remotePoolAddressesBytes := make([][][]byte, len(input.RemotePoolAddressesToAdd))
	for i, addresses := range input.RemotePoolAddressesToAdd {
		remotePoolAddressesBytes[i] = make([][]byte, len(addresses))
		for j, address := range addresses {
			remotePoolAddressesBytes[i][j] = []byte(address)
		}
	}

	// Convert []string to [][]byte for RemoteTokenAddressesToAdd
	remoteTokenAddressesBytes := make([][]byte, len(input.RemoteTokenAddressesToAdd))
	for i, address := range input.RemoteTokenAddressesToAdd {
		remoteTokenAddressesBytes[i] = []byte(address)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ApplyChainUpdates(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelectorsToRemove,
		input.RemoteChainSelectorsToAdd,
		remotePoolAddressesBytes,
		remoteTokenAddressesBytes,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool apply chain updates: %w", err)
	}

	b.Logger.Infow("ApplyChainUpdates on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleasePackageID)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleasePackageID,
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
	LockReleasePackageID string
	CoinObjectTypeArg    string
	StateObjectID        string
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
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleasePackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.SetChainRateLimiterConfigs(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
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
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool set configs rate limiter: %w", err)
	}

	b.Logger.Infow("SetChainRateLimiter on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleasePackageID)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleasePackageID,
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
	LockReleaseTokenPoolPackageID string
	CoinObjectTypeArg             string
	StateObjectID                 string
	Coin                          string
}

var providerLiquidityHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolProviderLiquidityInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ProvideLiquidity(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectID},
		bind.Object{Id: input.Coin},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to provide liquidity to lock release token pool: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageID,
		Objects:   NoObjects{},
	}, err
}

var LockReleaseTokenPoolProviderLiquidityOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "provide_liquidity"),
	semver.MustParse("0.1.0"),
	"Provide liquidity CCIP Lock Release Token Pool contract",
	providerLiquidityHandler,
)
