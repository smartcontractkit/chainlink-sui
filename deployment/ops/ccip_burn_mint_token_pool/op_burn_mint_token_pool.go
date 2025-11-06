package burnminttokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// BMTP -- INITIALIZE
type BurnMintTokenPoolInitializeObjects struct {
	StateObjectId string
}

type BurnMintTokenPoolInitializeInput struct {
	BurnMintPackageId      string
	OwnerCapObjectId       string
	CoinObjectTypeArg      string
	StateObjectId          string
	CoinMetadataObjectId   string
	TreasuryCapObjectId    string
	TokenPoolAdministrator string
}

var initBMTPHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolInitializeInput) (output sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create burn mint contract: %w", err)
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
	)
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute burn mint token pool initialization: %w", err)
	}

	stateObj, err := bind.FindObjectIdFromPublishTx(*tx, "burn_mint_token_pool", "BurnMintTokenPoolState")
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects: BurnMintTokenPoolInitializeObjects{
			StateObjectId: stateObj,
		},
	}, err
}

var BurnMintTokenPoolInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Burn Mint Token Pool contract",
	initBMTPHandler,
)

// BMTP -- apply_chain_updates
type NoObjects struct {
}

type BurnMintTokenPoolApplyChainUpdatesInput struct {
	BurnMintPackageId            string
	CoinObjectTypeArg            string
	StateObjectId                string
	OwnerCap                     string
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][]string
	RemoteTokenAddressesToAdd    []string
}

var applyChainUpdates = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolApplyChainUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint contract: %w", err)
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
		b.Logger.Infow("Skipping execution of ApplyChainUpdates on BurnMintTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.BurnMintPackageId,
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool apply chain updates: %w", err)
	}

	b.Logger.Infow("ApplyChainUpdates on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var BurnMintTokenPoolApplyChainUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "apply_chain_updates"),
	semver.MustParse("0.1.0"),
	"Applies chain updates in the CCIP Burn Mint Token Pool contract",
	applyChainUpdates,
)

// BMTP -- set_chain_rate_limiter_configs
type BurnMintTokenPoolSetChainRateLimiterInput struct {
	BurnMintPackageId    string
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

var setChainRateLimiterHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolSetChainRateLimiterInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint contract: %w", err)
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
		b.Logger.Infow("Skipping execution of SetChainRateLimiterConfigs on BurnMintTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.BurnMintPackageId,
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool set configs rate limiter: %w", err)
	}

	b.Logger.Infow("SetChainRateLimiter on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var BurnMintTokenPoolSetChainRateLimiterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "set_chain_rate_limiter_configs"),
	semver.MustParse("0.1.0"),
	"Sets chain rate limiter configs in the CCIP Burn Mint Token Pool contract",
	setChainRateLimiterHandler,
)

// BMTP -- add_remote_pool
type BurnMintTokenPoolAddRemotePoolInput struct {
	BurnMintTokenPoolPackageId string
	CoinObjectTypeArg          string
	StateObjectId              string
	OwnerCap                   string
	RemoteChainSelector        uint64
	RemotePoolAddress          string
}

var addRemotePoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolAddRemotePoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint token pool contract: %w", err)
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
		b.Logger.Infow("Skipping execution of AddRemotePool on BurnMintTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.BurnMintTokenPoolPackageId,
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool add remote pool: %w", err)
	}

	b.Logger.Infow("AddRemotePool on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintTokenPoolPackageId, "Chain:", input.RemoteChainSelector)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var BurnMintTokenPoolAddRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "add_remote_pool"),
	semver.MustParse("0.1.0"),
	"Adds a remote pool in the CCIP BurnMint Token Pool contract",
	addRemotePoolHandler,
)

// BMTP -- set_allowlist_enabled
type BurnMintTokenPoolSetAllowlistEnabledInput struct {
	BurnMintPackageId string `json:"burn_mint_package_id"`
	StateObjectId     string `json:"state_object_id"`
	OwnerCap          string `json:"owner_cap"`
	CoinObjectTypeArg string `json:"coin_object_type_arg"`
	Enabled           bool   `json:"enabled"`
}

var setAllowlistEnabledHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolSetAllowlistEnabledInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint contract: %w", err)
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
		b.Logger.Infow("Skipping execution of SetAllowlistEnabled on BurnMintTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.BurnMintPackageId,
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool set allowlist enabled: %w", err)
	}

	b.Logger.Infow("SetAllowlistEnabled on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var BurnMintTokenPoolSetAllowlistEnabledOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "set_allowlist_enabled"),
	semver.MustParse("0.1.0"),
	"Sets allowlist enabled in the CCIP Burn Mint Token Pool contract",
	setAllowlistEnabledHandler,
)

// BMTP -- apply_allowlist_updates
type BurnMintTokenPoolApplyAllowlistUpdatesInput struct {
	BurnMintPackageId string   `json:"burn_mint_package_id"`
	StateObjectId     string   `json:"state_object_id"`
	OwnerCap          string   `json:"owner_cap"`
	CoinObjectTypeArg string   `json:"coin_object_type_arg"`
	Removes           []string `json:"removes"`
	Adds              []string `json:"adds"`
}

var applyAllowlistUpdatesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolApplyAllowlistUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint contract: %w", err)
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
		b.Logger.Infow("Skipping execution of ApplyAllowlistUpdates on BurnMintTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.BurnMintPackageId,
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool apply allowlist updates: %w", err)
	}

	b.Logger.Infow("ApplyAllowlistUpdates on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var BurnMintTokenPoolApplyAllowlistUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "apply_allowlist_updates"),
	semver.MustParse("0.1.0"),
	"Applies allowlist updates in the CCIP Burn Mint Token Pool contract",
	applyAllowlistUpdatesHandler,
)

// BMTP -- remove_remote_pool
type BurnMintTokenPoolRemoveRemotePoolInput struct {
	BurnMintPackageId   string `json:"burn_mint_package_id"`
	StateObjectId       string `json:"state_object_id"`
	OwnerCap            string `json:"owner_cap"`
	CoinObjectTypeArg   string `json:"coin_object_type_arg"`
	RemoteChainSelector uint64 `json:"remote_chain_selector"`
	RemotePoolAddress   string `json:"remote_pool_address"`
}

var removeRemotePoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolRemoveRemotePoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint contract: %w", err)
	}

	// Convert string to bytes for RemotePoolAddress
	remotePoolAddressBytes, err := deployment.StrToBytes(input.RemotePoolAddress)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert remote pool address to bytes: %w", err)
	}

	encodedCall, err := contract.Encoder().RemoveRemotePool(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelector,
		remotePoolAddressBytes,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode RemoveRemotePool call: %w", err)
	}
	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.CoinObjectTypeArg})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of RemoveRemotePool on BurnMintTokenPool as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.BurnMintPackageId,
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool remove remote pool: %w", err)
	}

	b.Logger.Infow("RemoveRemotePool on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, err
}

var BurnMintTokenPoolRemoveRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "remove_remote_pool"),
	semver.MustParse("0.1.0"),
	"Removes remote pool in the CCIP Burn Mint Token Pool contract",
	removeRemotePoolHandler,
)
