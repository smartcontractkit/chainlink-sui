package lockreleasetokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// LRTP -- INITIALIZE
type LockReleaseTokenPoolInitializeObjects struct {
	StateObjectId         string
	RebalancerCapObjectId string
}

type LockReleaseTokenPoolInitializeInput struct {
	LockReleasePackageId   string
	OwnerCapObjectId       string
	CoinObjectTypeArg      string
	StateObjectId          string
	CoinMetadataObjectId   string
	TreasuryCapObjectId    string
	TokenPoolAdministrator string
	Rebalancer             string
}

var initLRTPHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolInitializeInput) (output sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleasePackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Initialize(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.CoinMetadataObjectId},
		bind.Object{Id: input.TreasuryCapObjectId},
		input.TokenPoolAdministrator,
		input.Rebalancer,
	)
	if err != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute lock release token pool initialization: %w", err)
	}

	stateObj, err1 := bind.FindObjectIdFromPublishTx(*tx, "lock_release_token_pool", "LockReleaseTokenPoolState")
	rebalancerObj, err2 := bind.FindObjectIdFromPublishTx(*tx, "lock_release_token_pool", "RebalancerCap")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[LockReleaseTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleasePackageId,
		Objects: LockReleaseTokenPoolInitializeObjects{
			StateObjectId:         stateObj,
			RebalancerCapObjectId: rebalancerObj,
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
	LockReleasePackageId         string
	CoinObjectTypeArg            string
	StateObjectId                string
	OwnerCap                     string
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][]string
	RemoteTokenAddressesToAdd    []string
}

var applyChainUpdates = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolApplyChainUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleasePackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	// Convert [][]string to [][][]byte for RemotePoolAddressesToAdd
	remotePoolAddressesBytes := make([][][]byte, len(input.RemotePoolAddressesToAdd))
	for i, addresses := range input.RemotePoolAddressesToAdd {
		remotePoolAddressesBytes[i] = make([][]byte, len(addresses))
		for j, address := range addresses {
			b, err := deployment.StrToBytes(address)
			if err != nil {
				return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("bad remote pool address [%d][%d]: %w", i, j, err)
			}
			remotePoolAddressesBytes[i][j] = b
		}
	}

	// Convert []string to [][]byte for RemoteTokenAddressesToAdd
	remoteTokenAddressesBytes := make([][]byte, len(input.RemoteTokenAddressesToAdd))
	for i, address := range input.RemoteTokenAddressesToAdd {
		b32, err := deployment.StrTo32(address)
		if err != nil {
			return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("bad remote token address [%d]: %w", i, err)
		}
		remoteTokenAddressesBytes[i] = b32
	}

	encodedCall, err := contract.Encoder().ApplyChainUpdates(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelectorsToRemove,
		input.RemoteChainSelectorsToAdd,
		remotePoolAddressesBytes,
		remoteTokenAddressesBytes,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode ApplyChainUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.CoinObjectTypeArg})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyChainUpdates on LockReleaseTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.LockReleasePackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool apply chain updates: %w", err)
	}

	b.Logger.Infow("ApplyChainUpdates on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleasePackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleasePackageId,
		Objects:   NoObjects{},
		Call:      call,
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
	LockReleasePackageId string
	CoinObjectTypeArg    string
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
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleasePackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	encodedCall, err := contract.Encoder().SetChainRateLimiterConfigs(
		[]string{input.CoinObjectTypeArg},
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
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode SetChainRateLimiterConfigs call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.CoinObjectTypeArg})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetChainRateLimiterConfigs on LockReleaseTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.LockReleasePackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool set configs rate limiter: %w", err)
	}

	b.Logger.Infow("SetChainRateLimiter on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleasePackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleasePackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var LockReleaseTokenPoolSetChainRateLimiterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "set_chain_rate_limiter_configs"),
	semver.MustParse("0.1.0"),
	"Sets chain rate limiter configs in the CCIP Lock Release Token Pool contract",
	setChainRateLimiterHandler,
)

// LRTP -- provide_liquidity
type LockReleaseTokenPoolProvideLiquidityInput struct {
	LockReleaseTokenPoolPackageId string
	CoinObjectTypeArg             string
	StateObjectId                 string
	RebalancerCapObjectId         string
	Coin                          string
}

var provideLiquidityHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolProvideLiquidityInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ProvideLiquidity(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.RebalancerCapObjectId},
		bind.Object{Id: input.Coin},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to provide liquidity to lock release token pool: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   NoObjects{},
	}, err
}

var LockReleaseTokenPoolProvideLiquidityOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "provide_liquidity"),
	semver.MustParse("0.1.0"),
	"Provide liquidity CCIP Lock Release Token Pool contract",
	provideLiquidityHandler,
)

// LRTP -- add_remote_pool
type LockReleaseTokenPoolAddRemotePoolInput struct {
	LockReleaseTokenPoolPackageId string
	CoinObjectTypeArg             string
	StateObjectId                 string
	OwnerCap                      string
	RemoteChainSelector           uint64
	RemotePoolAddress             string
}

var addRemotePoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolAddRemotePoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release token pool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().AddRemotePool(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelector,
		[]byte(input.RemotePoolAddress),
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode AddRemotePool call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.CoinObjectTypeArg})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AddRemotePool on LockReleaseTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.LockReleaseTokenPoolPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool add remote pool: %w", err)
	}

	b.Logger.Infow("AddRemotePool on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleaseTokenPoolPackageId, "Chain:", input.RemoteChainSelector)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var LockReleaseTokenPoolAddRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "add_remote_pool"),
	semver.MustParse("0.1.0"),
	"Adds a remote pool in the CCIP LockRelease Token Pool contract",
	addRemotePoolHandler,
)

// LRTP -- set_allowlist_enabled
type LockReleaseTokenPoolSetAllowlistEnabledInput struct {
	LockReleaseTokenPoolPackageId string
	CoinObjectTypeArg             string
	StateObjectId                 string
	OwnerCap                      string
	Enabled                       bool
}

var setAllowlistEnabledHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolSetAllowlistEnabledInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release token pool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().SetAllowlistEnabled(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.Enabled,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode SetAllowlistEnabled call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.CoinObjectTypeArg})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetAllowlistEnabled on LockReleaseTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.LockReleaseTokenPoolPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool set allowlist enabled: %w", err)
	}

	b.Logger.Infow("SetAllowlistEnabled on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleaseTokenPoolPackageId, "Enabled:", input.Enabled)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var LockReleaseTokenPoolSetAllowlistEnabledOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "set_allowlist_enabled"),
	semver.MustParse("0.1.0"),
	"Sets allowlist enabled in the CCIP LockRelease Token Pool contract",
	setAllowlistEnabledHandler,
)

// LRTP -- apply_allowlist_updates
type LockReleaseTokenPoolApplyAllowlistUpdatesInput struct {
	LockReleaseTokenPoolPackageId string
	CoinObjectTypeArg             string
	StateObjectId                 string
	OwnerCap                      string
	Removes                       []string
	Adds                          []string
}

var applyAllowlistUpdatesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolApplyAllowlistUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release token pool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ApplyAllowlistUpdates(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.Removes,
		input.Adds,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode ApplyAllowlistUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.CoinObjectTypeArg})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyAllowlistUpdates on LockReleaseTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.LockReleaseTokenPoolPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool apply allowlist updates: %w", err)
	}

	b.Logger.Infow("ApplyAllowlistUpdates on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleaseTokenPoolPackageId, "Removes:", len(input.Removes), "Adds:", len(input.Adds))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var LockReleaseTokenPoolApplyAllowlistUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "apply_allowlist_updates"),
	semver.MustParse("0.1.0"),
	"Applies allowlist updates in the CCIP LockRelease Token Pool contract",
	applyAllowlistUpdatesHandler,
)

// LRTP -- remove_remote_pool
type LockReleaseTokenPoolRemoveRemotePoolInput struct {
	LockReleaseTokenPoolPackageId string
	CoinObjectTypeArg             string
	StateObjectId                 string
	OwnerCap                      string
	RemoteChainSelector           uint64
	RemotePoolAddress             string
}

var removeRemotePoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input LockReleaseTokenPoolRemoveRemotePoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(input.LockReleaseTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create lock release token pool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().RemoveRemotePool(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelector,
		[]byte(input.RemotePoolAddress),
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode RemoveRemotePool call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.CoinObjectTypeArg})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of RemoveRemotePool on LockReleaseTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.LockReleaseTokenPoolPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute lock release token pool remove remote pool: %w", err)
	}

	b.Logger.Infow("RemoveRemotePool on LockReleaseTokenPool", "LockReleaseTokenPool PackageId:", input.LockReleaseTokenPoolPackageId, "Chain:", input.RemoteChainSelector)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.LockReleaseTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var LockReleaseTokenPoolRemoveRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "lock_release_token_pool", "remove_remote_pool"),
	semver.MustParse("0.1.0"),
	"Removes a remote pool in the CCIP LockRelease Token Pool contract",
	removeRemotePoolHandler,
)
