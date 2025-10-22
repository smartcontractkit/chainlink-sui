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
	OwnerCapObjectId string
	StateObjectId    string
}

type BurnMintTokenPoolInitializeInput struct {
	BurnMintPackageId      string
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
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.CoinMetadataObjectId},
		bind.Object{Id: input.TreasuryCapObjectId},
		input.TokenPoolAdministrator,
	)
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute burn mint token pool initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "burn_mint_token_pool", "BurnMintTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects: BurnMintTokenPoolInitializeObjects{
			OwnerCapObjectId: obj1,
			StateObjectId:    obj2,
		},
	}, err
}

var BurnMintTokenPoolInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Burn Mint Token Pool contract",
	initBMTPHandler,
)

// BMTP -- INITIALIZE BY CCIP ADMIN
type BurnMintTokenPoolInitializeByCcipAdminInput struct {
	BurnMintPackageId      string
	CoinObjectTypeArg      string
	StateObjectId          string
	CoinMetadataObjectId   string
	OwnerCapObjectId       string
	TreasuryCapObjectId    string
	TokenPoolAdministrator string
}

var initByCcipAdminBMTPHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolInitializeByCcipAdminInput) (output sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create burn mint contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.InitializeByCcipAdmin(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.CoinMetadataObjectId},
		bind.Object{Id: input.TreasuryCapObjectId},
		input.TokenPoolAdministrator,
	)
	if err != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute burn mint token pool initialization by ccip admin: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "burn_mint_token_pool", "BurnMintTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[BurnMintTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects: BurnMintTokenPoolInitializeObjects{
			OwnerCapObjectId: obj1,
			StateObjectId:    obj2,
		},
	}, err
}

var BurnMintTokenPoolInitializeByCcipAdminOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "initialize_by_ccip_admin"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Burn Mint Token Pool contract by CCIP admin",
	initByCcipAdminBMTPHandler,
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
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyChainUpdates on BurnMintTokenPool as per no Signer provided",
			"chainsToRemove", len(input.RemoteChainSelectorsToRemove), "chainsToAdd", len(input.RemoteChainSelectorsToAdd))
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ApplyChainUpdates on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Chain updates applied on BurnMintTokenPool", "PackageId", input.BurnMintPackageId,
		"ChainsRemoved", len(input.RemoteChainSelectorsToRemove), "ChainsAdded", len(input.RemoteChainSelectorsToAdd))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
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
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetChainRateLimiterConfigs on BurnMintTokenPool as per no Signer provided",
			"chains", len(input.RemoteChainSelectors))
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute SetChainRateLimiterConfigs on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Chain rate limiter configs set on BurnMintTokenPool", "PackageId", input.BurnMintPackageId,
		"Chains", len(input.RemoteChainSelectors))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
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
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AddRemotePool on BurnMintTokenPool as per no Signer provided",
			"chain", input.RemoteChainSelector)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute AddRemotePool on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Remote pool added to BurnMintTokenPool", "PackageId", input.BurnMintTokenPoolPackageId,
		"Chain", input.RemoteChainSelector)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var BurnMintTokenPoolAddRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "add_remote_pool"),
	semver.MustParse("0.1.0"),
	"Adds a remote pool in the CCIP BurnMint Token Pool contract",
	addRemotePoolHandler,
)

// BMTP -- set_pool
type BurnMintTokenPoolSetPoolInput struct {
	BurnMintTokenPoolPackageId string
	CoinObjectTypeArg          string
	RefObjectId                string
	StateObjectId              string
	OwnerCap                   string
	CoinMetadataAddress        string
}

var setPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolSetPoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint token pool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().SetPool(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.RefObjectId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.CoinMetadataAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode SetPool call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetPool on BurnMintTokenPool as per no Signer provided",
			"coinMetadata", input.CoinMetadataAddress)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute SetPool on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Pool set on BurnMintTokenPool", "PackageId", input.BurnMintTokenPoolPackageId,
		"CoinMetadata", input.CoinMetadataAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var BurnMintTokenPoolSetPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "set_pool"),
	semver.MustParse("0.1.0"),
	"Sets the pool in the token admin registry for the CCIP Burn Mint Token Pool",
	setPoolHandler,
)

// BMTP -- set_allowlist_enabled
type BurnMintTokenPoolSetAllowlistEnabledInput struct {
	BurnMintTokenPoolPackageId string
	CoinObjectTypeArg          string
	StateObjectId              string
	OwnerCap                   string
	Enabled                    bool
}

func setAllowlistEnabledHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolSetAllowlistEnabledInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint token pool contract: %w", err)
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
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetAllowlistEnabled on BurnMintTokenPool as per no Signer provided", "enabled", input.Enabled)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute SetAllowlistEnabled on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Allowlist enabled state set on BurnMintTokenPool", "PackageId", input.BurnMintTokenPoolPackageId, "Enabled", input.Enabled)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var BurnMintTokenPoolSetAllowlistEnabledOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "set_allowlist_enabled"),
	semver.MustParse("0.1.0"),
	"Sets the allowlist enabled state for the CCIP Burn Mint Token Pool",
	setAllowlistEnabledHandler,
)

// BMTP -- apply_allowlist_updates
type BurnMintTokenPoolApplyAllowlistUpdatesInput struct {
	BurnMintTokenPoolPackageId string
	CoinObjectTypeArg          string
	StateObjectId              string
	OwnerCap                   string
	Removes                    []string
	Adds                       []string
}

func applyAllowlistUpdatesHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolApplyAllowlistUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint token pool contract: %w", err)
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
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyAllowlistUpdates on BurnMintTokenPool as per no Signer provided",
			"removes", len(input.Removes), "adds", len(input.Adds))
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ApplyAllowlistUpdates on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Allowlist updates applied on BurnMintTokenPool", "PackageId", input.BurnMintTokenPoolPackageId,
		"Removes", len(input.Removes), "Adds", len(input.Adds))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var BurnMintTokenPoolApplyAllowlistUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "apply_allowlist_updates"),
	semver.MustParse("0.1.0"),
	"Applies allowlist updates for the CCIP Burn Mint Token Pool",
	applyAllowlistUpdatesHandler,
)

// BMTP -- remove_remote_pool
type BurnMintTokenPoolRemoveRemotePoolInput struct {
	BurnMintTokenPoolPackageId string
	CoinObjectTypeArg          string
	StateObjectId              string
	OwnerCap                   string
	RemoteChainSelector        uint64
	RemotePoolAddress          string
}

func removeRemotePoolHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolRemoveRemotePoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint token pool contract: %w", err)
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
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of RemoveRemotePool on BurnMintTokenPool as per no Signer provided",
			"chain", input.RemoteChainSelector)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute RemoveRemotePool on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Remote pool removed from BurnMintTokenPool", "PackageId", input.BurnMintTokenPoolPackageId,
		"Chain", input.RemoteChainSelector)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var BurnMintTokenPoolRemoveRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "remove_remote_pool"),
	semver.MustParse("0.1.0"),
	"Removes a remote pool from the CCIP Burn Mint Token Pool",
	removeRemotePoolHandler,
)

// BMTP -- transfer_ownership
type BurnMintTokenPoolTransferOwnershipInput struct {
	BurnMintTokenPoolPackageId string
	CoinObjectTypeArg          string
	StateObjectId              string
	OwnerCap                   string
	NewOwner                   string
}

func transferOwnershipHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolTransferOwnershipInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint token pool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().TransferOwnership(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.NewOwner,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode TransferOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of TransferOwnership on BurnMintTokenPool as per no Signer provided",
			"newOwner", input.NewOwner)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute TransferOwnership on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated on BurnMintTokenPool", "PackageId", input.BurnMintTokenPoolPackageId,
		"NewOwner", input.NewOwner)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var BurnMintTokenPoolTransferOwnershipOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "transfer_ownership"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the CCIP Burn Mint Token Pool",
	transferOwnershipHandler,
)

// BMTP -- execute_ownership_transfer
type BurnMintTokenPoolExecuteOwnershipTransferInput struct {
	BurnMintTokenPoolPackageId string
	CoinObjectTypeArg          string
	OwnerCap                   string
	StateObjectId              string
	To                         string
}

var executeOwnershipTransferHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BurnMintTokenPoolExecuteOwnershipTransferInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(input.BurnMintTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create burn mint token pool contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransfer(
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.OwnerCap},
		bind.Object{Id: input.StateObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransfer call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransfer on BurnMintTokenPool as per no Signer provided",
			"to", input.To)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransfer on BurnMintTokenPool: %w", err)
	}

	b.Logger.Infow("Ownership transfer executed on BurnMintTokenPool", "PackageId", input.BurnMintTokenPoolPackageId,
		"To", input.To)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintTokenPoolPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var BurnMintTokenPoolExecuteOwnershipTransferOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "execute_ownership_transfer"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer for the CCIP Burn Mint Token Pool",
	executeOwnershipTransferHandler,
)
