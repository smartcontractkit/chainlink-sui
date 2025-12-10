package view

import (
	"context"
	"fmt"
	"reflect"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	module_managed_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool"
	module_usdc_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/usdc_token_pool"
)

type TokenBucketWrapper interface {
	module_burn_mint_token_pool.TokenBucketWrapper |
		module_lock_release_token_pool.TokenBucketWrapper |
		module_managed_token_pool.TokenBucketWrapper |
		module_usdc_token_pool.TokenBucketWrapper
}

// ITokenPoolDevInspect defines the common DevInspect methods across pool types
type ITokenPoolDevInspect[T TokenBucketWrapper] interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error)
	GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (T, error)
	GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (T, error)
}

type TokenPoolView struct {
	ContractMetaData

	Token              string                       `json:"token"`
	RemoteChainConfigs map[uint64]RemoteChainConfig `json:"remoteChainConfigs"`
	AllowList          []string                     `json:"allowList"`
	AllowListEnabled   bool                         `json:"allowListEnabled"`
	RebalancerCapIds   []string                     `json:"rebalancerCapIds,omitempty"` // only applicable for LR TP
}

type RemoteChainConfig struct {
	RemoteTokenAddress        string            `json:"remoteTokenAddress"`
	RemotePoolAddresses       []string          `json:"remotePoolAddresses"`
	InboundRateLimiterConfig  RateLimiterConfig `json:"inboundRateLimiterConfig"`
	OutboundRateLimiterConfig RateLimiterConfig `json:"outboundRateLimiterConfig"`
}

type RateLimiterConfig struct {
	IsEnabled bool   `json:"isEnabled"`
	Capacity  uint64 `json:"capacity"`
	Rate      uint64 `json:"rate"`
}

func GenerateTokenPoolView[T TokenBucketWrapper](
	ctx context.Context,
	chain sui.Chain,
	poolPackageID string,
	poolStateObjectID string,
	tokenConfigs map[string]TokenConfigView,
	poolDevInspect ITokenPoolDevInspect[T],
	lggr logger.Logger,
) (TokenPoolView, error) {
	callOpts := &bind.CallOpts{Signer: chain.Signer}

	// Get token type args from TokenAdminRegistry
	var typeArgs []string
	for _, c := range tokenConfigs {
		if c.TokenPoolPackageId == poolPackageID {
			typeArgs = append(typeArgs, c.TokenType)
			break
		}
	}
	if len(typeArgs) == 0 {
		return TokenPoolView{}, fmt.Errorf("no token config found for token pool on TokenAdminRegistry, package ID: %s", poolPackageID)
	}

	// Get TypeAndVersion
	typeAndVersion, err := poolDevInspect.TypeAndVersion(ctx, callOpts)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get type and version: %w", err)
	}

	// Get Owner
	poolStateObj := bind.Object{Id: poolStateObjectID}
	owner, err := poolDevInspect.Owner(ctx, callOpts, typeArgs, poolStateObj)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get owner: %w", err)
	}

	// Get Token
	token, err := poolDevInspect.GetToken(ctx, callOpts, typeArgs, poolStateObj)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get token: %w", err)
	}

	// Get AllowList Enabled
	allowlistEnabled, err := poolDevInspect.GetAllowlistEnabled(ctx, callOpts, typeArgs, poolStateObj)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get allowlist enabled: %w", err)
	}

	// Get AllowList
	allowlist, err := poolDevInspect.GetAllowlist(ctx, callOpts, typeArgs, poolStateObj)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get allowlist: %w", err)
	}

	// Get Supported Chains
	supportedChains, err := poolDevInspect.GetSupportedChains(ctx, callOpts, typeArgs, poolStateObj)
	if err != nil {
		return TokenPoolView{}, fmt.Errorf("failed to get supported chains: %w", err)
	}

	// Get Remote Chain Configs for each supported chain
	remoteChainConfigs := make(map[uint64]RemoteChainConfig)
	clockObj := bind.Object{Id: "0x6"}

	for _, chainSelector := range supportedChains {
		// Get Remote Pools
		remotePools, err := poolDevInspect.GetRemotePools(ctx, callOpts, typeArgs, poolStateObj, chainSelector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get remote pools for chain %d: %w", chainSelector, err)
		}

		// Get Remote Token
		remoteToken, err := poolDevInspect.GetRemoteToken(ctx, callOpts, typeArgs, poolStateObj, chainSelector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get remote token for chain %d: %w", chainSelector, err)
		}

		// Get Inbound Rate Limiter State
		inboundRateLimiter, err := poolDevInspect.GetCurrentInboundRateLimiterState(ctx, callOpts, typeArgs, clockObj, poolStateObj, chainSelector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get inbound rate limiter for chain %d: %w", chainSelector, err)
		}

		// Get Outbound Rate Limiter State
		outboundRateLimiter, err := poolDevInspect.GetCurrentOutboundRateLimiterState(ctx, callOpts, typeArgs, clockObj, poolStateObj, chainSelector)
		if err != nil {
			return TokenPoolView{}, fmt.Errorf("failed to get outbound rate limiter for chain %d: %w", chainSelector, err)
		}

		// Convert remote pools to hex strings
		remotePoolAddresses := make([]string, len(remotePools))
		for i, pool := range remotePools {
			remotePoolAddresses[i] = fmt.Sprintf("0x%x", pool)
		}

		// Parse rate limiter states
		inboundConfig := parseRateLimiterConfig(inboundRateLimiter)
		outboundConfig := parseRateLimiterConfig(outboundRateLimiter)

		remoteChainConfigs[chainSelector] = RemoteChainConfig{
			RemoteTokenAddress:        fmt.Sprintf("0x%x", remoteToken),
			RemotePoolAddresses:       remotePoolAddresses,
			InboundRateLimiterConfig:  inboundConfig,
			OutboundRateLimiterConfig: outboundConfig,
		}
	}

	return TokenPoolView{
		ContractMetaData: ContractMetaData{
			TypeAndVersion: typeAndVersion,
			Owner:          owner,
			Address:        poolPackageID,
			StateObjectID:  poolStateObjectID,
		},
		Token:              token,
		RemoteChainConfigs: remoteChainConfigs,
		AllowList:          allowlist,
		AllowListEnabled:   allowlistEnabled,
	}, nil
}

func parseRateLimiterConfig[T TokenBucketWrapper](rateLimiterObj T) RateLimiterConfig {
	v := reflect.ValueOf(rateLimiterObj)
	return RateLimiterConfig{
		IsEnabled: v.FieldByName("IsEnabled").Bool(),
		Capacity:  v.FieldByName("Capacity").Uint(),
		Rate:      v.FieldByName("Rate").Uint(),
	}
}
