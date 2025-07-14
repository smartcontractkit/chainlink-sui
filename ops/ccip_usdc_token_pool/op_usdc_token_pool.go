package usdctokenpoolops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_usdc_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/usdc_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

// USDC Token Pool -- INITIALIZE
type USDCTokenPoolInitializeObjects struct {
	OwnerCapObjectId string
	StateObjectId    string
}

type USDCTokenPoolInitializeInput struct {
	USDCTokenPoolPackageId string
	CoinObjectTypeArg      string
	StateObjectId          string
	OwnerCapObjectId       string
	CoinMetadataObjectId   string
	LocalDomainIdentifier  uint32
	TokenPoolPackageId     string
	TokenPoolAdministrator string
	LockOrBurnParams       []string
	ReleaseOrMintParams    []string
}

var initUSDCTokenPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input USDCTokenPoolInitializeInput) (output sui_ops.OpTxResult[USDCTokenPoolInitializeObjects], err error) {
	contract, err := module_usdc_token_pool.NewUsdcTokenPool(input.USDCTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[USDCTokenPoolInitializeObjects]{}, fmt.Errorf("failed to create USDC token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Initialize(
		b.GetContext(),
		opts,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.CoinMetadataObjectId},
		input.LocalDomainIdentifier,
		input.TokenPoolPackageId,
		input.TokenPoolAdministrator,
		input.LockOrBurnParams,
		input.ReleaseOrMintParams,
	)
	if err != nil {
		return sui_ops.OpTxResult[USDCTokenPoolInitializeObjects]{}, fmt.Errorf("failed to execute USDC token pool initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "usdc_token_pool", "USDCTokenPoolState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[USDCTokenPoolInitializeObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[USDCTokenPoolInitializeObjects]{
		Digest:    tx.Digest,
		PackageId: input.USDCTokenPoolPackageId,
		Objects: USDCTokenPoolInitializeObjects{
			OwnerCapObjectId: obj1,
			StateObjectId:    obj2,
		},
	}, err
}

var USDCTokenPoolInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "usdc_token_pool", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP USDC Token Pool contract",
	initUSDCTokenPoolHandler,
)

// USDC Token Pool -- SET_DOMAINS
type NoObjects struct {
}

type USDCTokenPoolSetDomainsInput struct {
	USDCTokenPoolPackageId  string
	CoinObjectTypeArg       string
	StateObjectId           string
	OwnerCap                string
	RemoteChainSelectors    []uint64
	RemoteDomainIdentifiers []uint32
	AllowedRemoteCallers    [][]byte
	Enableds                []bool
}

var setDomainsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input USDCTokenPoolSetDomainsInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_usdc_token_pool.NewUsdcTokenPool(input.USDCTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create USDC token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.SetDomains(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelectors,
		input.RemoteDomainIdentifiers,
		input.AllowedRemoteCallers,
		input.Enableds,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute set domains: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.USDCTokenPoolPackageId,
	}, err
}

var USDCTokenPoolSetDomainsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "usdc_token_pool", "set_domains"),
	semver.MustParse("0.1.0"),
	"Sets domain configurations for the USDC Token Pool",
	setDomainsHandler,
)

// USDC Token Pool -- APPLY_CHAIN_UPDATES
type USDCTokenPoolApplyChainUpdatesInput struct {
	USDCTokenPoolPackageId       string
	CoinObjectTypeArg            string
	StateObjectId                string
	OwnerCap                     string
	RemoteChainSelectorsToRemove []uint64
	RemoteChainSelectorsToAdd    []uint64
	RemotePoolAddressesToAdd     [][]string
	RemoteTokenAddressesToAdd    []string
}

var applyChainUpdatesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input USDCTokenPoolApplyChainUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_usdc_token_pool.NewUsdcTokenPool(input.USDCTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create USDC token pool contract: %w", err)
	}

	// Convert string arrays to byte arrays for remote pool addresses
	remotePoolAddresses := make([][][]byte, len(input.RemotePoolAddressesToAdd))
	for i, addresses := range input.RemotePoolAddressesToAdd {
		remotePoolAddresses[i] = make([][]byte, len(addresses))
		for j, addr := range addresses {
			remotePoolAddresses[i][j] = []byte(addr)
		}
	}

	// Convert string addresses to byte arrays for remote token addresses
	remoteTokenAddresses := make([][]byte, len(input.RemoteTokenAddressesToAdd))
	for i, addr := range input.RemoteTokenAddressesToAdd {
		remoteTokenAddresses[i] = []byte(addr)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ApplyChainUpdates(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		input.RemoteChainSelectorsToRemove,
		input.RemoteChainSelectorsToAdd,
		remotePoolAddresses,
		remoteTokenAddresses,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute apply chain updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.USDCTokenPoolPackageId,
	}, err
}

var USDCTokenPoolApplyChainUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "usdc_token_pool", "apply_chain_updates"),
	semver.MustParse("0.1.0"),
	"Applies chain updates to the USDC Token Pool",
	applyChainUpdatesHandler,
)

// USDC Token Pool -- SET_CHAIN_RATE_LIMITER_CONFIGS
type USDCTokenPoolSetChainRateLimiterInput struct {
	USDCTokenPoolPackageId string
	CoinObjectTypeArg      string
	StateObjectId          string
	OwnerCap               string
	ClockObjectId          string
	RemoteChainSelectors   []uint64
	OutboundIsEnableds     []bool
	OutboundCapacities     []uint64
	OutboundRates          []uint64
	InboundIsEnableds      []bool
	InboundCapacities      []uint64
	InboundRates           []uint64
}

var setChainRateLimiterHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input USDCTokenPoolSetChainRateLimiterInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_usdc_token_pool.NewUsdcTokenPool(input.USDCTokenPoolPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create USDC token pool contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.SetChainRateLimiterConfigs(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCap},
		bind.Object{Id: input.ClockObjectId},
		input.RemoteChainSelectors,
		input.OutboundIsEnableds,
		input.OutboundCapacities,
		input.OutboundRates,
		input.InboundIsEnableds,
		input.InboundCapacities,
		input.InboundRates,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute set chain rate limiter configs: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.USDCTokenPoolPackageId,
	}, err
}

var USDCTokenPoolSetChainRateLimiterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "usdc_token_pool", "set_chain_rate_limiter"),
	semver.MustParse("0.1.0"),
	"Sets chain rate limiter configurations for the USDC Token Pool",
	setChainRateLimiterHandler,
)
