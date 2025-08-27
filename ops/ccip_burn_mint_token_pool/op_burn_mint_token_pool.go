package burnminttokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
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
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelectorsToRemove,
		input.RemoteChainSelectorsToAdd,
		remotePoolAddressesBytes,
		remoteTokenAddressesBytes,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool apply chain updates: %w", err)
	}

	b.Logger.Infow("ApplyChainUpdates on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
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

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.SetChainRateLimiterConfigs(
		b.GetContext(),
		opts,
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute burn mint token pool set configs rate limiter: %w", err)
	}

	b.Logger.Infow("SetChainRateLimiter on BurnMintTokenPool", "BurnMintTokenPool PackageId:", input.BurnMintPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.BurnMintPackageId,
		Objects:   NoObjects{},
	}, err
}

var BurnMintTokenPoolSetChainRateLimiterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "burn_mint_token_pool", "set_chain_rate_limiter_configs"),
	semver.MustParse("0.1.0"),
	"Sets chain rate limiter configs in the CCIP Burn Mint Token Pool contract",
	setChainRateLimiterHandler,
)
