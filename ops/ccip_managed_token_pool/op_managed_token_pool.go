package managedtokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

// MTP -- INITIALIZE_WITH_MANAGED_TOKEN
type ManagedTokenPoolInitializeObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type ManagedTokenPoolInitializeInput struct {
	ManagedTokenPoolPackageId string
	CoinObjectTypeArg         string
	CCIPObjectRefObjectId     string
	ManagedTokenStateObjectId string
	ManagedTokenOwnerCapId    string
	CoinMetadataObjectId      string
	MintCapObjectId           string
	TokenPoolAdministrator    string
}

var initMTPHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolInitializeInput) (output sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.InitializeWithManagedToken(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.CCIPObjectRefObjectId},
		bind.Object{Id: input.ManagedTokenStateObjectId},
		bind.Object{Id: input.ManagedTokenOwnerCapId},
		bind.Object{Id: input.CoinMetadataObjectId},
		bind.Object{Id: input.MintCapObjectId},
		input.TokenPoolAdministrator,
	)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute managed token pool initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "managed_token_pool", "ManagedTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects: ManagedTokenPoolInitializeObjects{
			OwnerCapObjectId: obj1,
			StateObjectId:    obj2,
		},
	}, err
}

var ManagedTokenPoolInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Managed Token Pool contract",
	initMTPHandler,
)

// MTP -- INITIALIZE BY CCIP ADMIN
type ManagedTokenPoolInitializeByCcipAdminInput struct {
	ManagedTokenPoolPackageId string
	CoinObjectTypeArg         string
	CCIPObjectRefObjectId     string
	OwnerCapObjectId          string
	CoinMetadataObjectId      string
	MintCapObjectId           string
	ManagedTokenStateObjectId string
	TokenPoolAdministrator    string
}

var initByCcipAdminManagedTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolInitializeByCcipAdminInput) (output sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.InitializeByCcipAdmin(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.CCIPObjectRefObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.CoinMetadataObjectId},
		bind.Object{Id: input.MintCapObjectId},
		input.ManagedTokenStateObjectId,
		input.TokenPoolAdministrator,
	)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute managed token pool initialization by ccip admin: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "managed_token_pool", "ManagedTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[ManagedTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects: ManagedTokenPoolInitializeObjects{
			OwnerCapObjectId: obj1,
			StateObjectId:    obj2,
		},
	}, err
}

var ManagedTokenPoolInitializeByCcipAdminOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "initialize_by_ccip_admin"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Managed Token Pool contract by CCIP admin",
	initByCcipAdminManagedTokenPoolHandler,
)

// MTP -- apply_chain_updates
type NoObjects struct {
}

type ManagedTokenPoolApplyChainUpdatesInput struct {
	ManagedTokenPoolPackageId    string
	CoinObjectTypeArg            string
	StateObjectId                string
	OwnerCap                     string
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][]string
	RemoteTokenAddressesToAdd    []string
}

var applyChainUpdates = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolApplyChainUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pool apply chain updates: %w", err)
	}

	b.Logger.Infow("ApplyChainUpdates on ManagedTokenPool", "ManagedTokenPool PackageId:", input.ManagedTokenPoolPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenPoolApplyChainUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "apply_chain_updates"),
	semver.MustParse("0.1.0"),
	"Applies chain updates in the CCIP Managed Token Pool contract",
	applyChainUpdates,
)

// MTP -- add_remote_pool
type ManagedTokenPoolAddRemotePoolInput struct {
	ManagedTokenPoolPackageId string
	CoinObjectTypeArg         string
	StateObjectId             string
	OwnerCap                  string
	RemoteChainSelector       uint64
	RemotePoolAddress         string
}

var addRemotePoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolAddRemotePoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.AddRemotePool(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelector,
		[]byte(input.RemotePoolAddress),
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pool add remote pool: %w", err)
	}

	b.Logger.Infow("AddRemotePool on ManagedTokenPool", "ManagedTokenPool PackageId:", input.ManagedTokenPoolPackageId, "Chain:", input.RemoteChainSelector)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenPoolAddRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "add_remote_pool"),
	semver.MustParse("0.1.0"),
	"Adds a remote pool in the CCIP Managed Token Pool contract",
	addRemotePoolHandler,
)

// MTP -- remove_remote_pool
type ManagedTokenPoolRemoveRemotePoolInput struct {
	ManagedTokenPoolPackageId string
	CoinObjectTypeArg         string
	StateObjectId             string
	OwnerCap                  string
	RemoteChainSelector       uint64
	RemotePoolAddress         string
}

var removeRemotePoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolRemoveRemotePoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.RemoveRemotePool(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelector,
		[]byte(input.RemotePoolAddress),
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pool remove remote pool: %w", err)
	}

	b.Logger.Infow("RemoveRemotePool on ManagedTokenPool", "ManagedTokenPool PackageId:", input.ManagedTokenPoolPackageId, "Chain:", input.RemoteChainSelector)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenPoolRemoveRemotePoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "remove_remote_pool"),
	semver.MustParse("0.1.0"),
	"Removes a remote pool in the CCIP Managed Token Pool contract",
	removeRemotePoolHandler,
)

// MTP -- set_chain_rate_limiter_configs
type ManagedTokenPoolSetChainRateLimiterInput struct {
	ManagedTokenPoolPackageId string
	CoinObjectTypeArg         string
	StateObjectId             string
	OwnerCap                  string
	RemoteChainSelectors      []uint64
	OutboundIsEnableds        []bool
	OutboundCapacities        []uint64
	OutboundRates             []uint64
	InboundIsEnableds         []bool
	InboundCapacities         []uint64
	InboundRates              []uint64
}

var setChainRateLimiterHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolSetChainRateLimiterInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pool set configs rate limiter: %w", err)
	}

	b.Logger.Infow("SetChainRateLimiter on ManagedTokenPool", "ManagedTokenPool PackageId:", input.ManagedTokenPoolPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenPoolSetChainRateLimiterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "set_chain_rate_limiter_configs"),
	semver.MustParse("0.1.0"),
	"Sets chain rate limiter configs in the CCIP Managed Token Pool contract",
	setChainRateLimiterHandler,
)

// MTP -- set_allowlist_enabled
type ManagedTokenPoolSetAllowlistEnabledInput struct {
	ManagedTokenPoolPackageId string
	CoinObjectTypeArg         string
	StateObjectId             string
	OwnerCap                  string
	Enabled                   bool
}

var setAllowlistEnabledHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolSetAllowlistEnabledInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.SetAllowlistEnabled(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.Enabled,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pool set allowlist enabled: %w", err)
	}

	b.Logger.Infow("SetAllowlistEnabled on ManagedTokenPool", "ManagedTokenPool PackageId:", input.ManagedTokenPoolPackageId, "Enabled:", input.Enabled)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenPoolSetAllowlistEnabledOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "set_allowlist_enabled"),
	semver.MustParse("0.1.0"),
	"Sets allowlist enabled in the CCIP Managed Token Pool contract",
	setAllowlistEnabledHandler,
)

// MTP -- apply_allowlist_updates
type ManagedTokenPoolApplyAllowlistUpdatesInput struct {
	ManagedTokenPoolPackageId string
	CoinObjectTypeArg         string
	StateObjectId             string
	OwnerCap                  string
	Removes                   []string
	Adds                      []string
}

var applyAllowlistUpdatesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenPoolApplyAllowlistUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_managed_token_pool.NewManagedTokenPool(input.ManagedTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create managed token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ApplyAllowlistUpdates(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.Removes,
		input.Adds,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute managed token pool apply allowlist updates: %w", err)
	}

	b.Logger.Infow("ApplyAllowlistUpdates on ManagedTokenPool", "ManagedTokenPool PackageId:", input.ManagedTokenPoolPackageId, "Removes:", len(input.Removes), "Adds:", len(input.Adds))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenPoolPackageId,
		Objects:   NoObjects{},
	}, err
}

var ManagedTokenPoolApplyAllowlistUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "managed_token_pool", "apply_allowlist_updates"),
	semver.MustParse("0.1.0"),
	"Applies allowlist updates in the CCIP Managed Token Pool contract",
	applyAllowlistUpdatesHandler,
)
