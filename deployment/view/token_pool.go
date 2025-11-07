package view

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// ITokenPool defines the common interface for all token pool types
type ITokenPool interface {
	DevInspect() ITokenPoolDevInspect
}

// ITokenPoolDevInspect defines the common DevInspect methods across pool types
type ITokenPoolDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	Owner(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (string, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) (bool, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]string, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object) ([]uint64, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([][]byte, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, remoteChainSelector uint64) ([]byte, error)
	GetCurrentInboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (bind.Object, error)
	GetCurrentOutboundRateLimiterState(ctx context.Context, opts *bind.CallOpts, typeArgs []string, clock bind.Object, state bind.Object, remoteChainSelector uint64) (bind.Object, error)
}

type TokenPoolView struct {
	ContractMetaData

	Token              string                       `json:"token"`
	RemoteChainConfigs map[uint64]RemoteChainConfig `json:"remoteChainConfigs"`
	AllowList          []string                     `json:"allowList"`
	AllowListEnabled   bool                         `json:"allowListEnabled"`
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

func GenerateTokenPoolView(
	ctx context.Context,
	chain sui.Chain,
	poolPackageID string,
	poolStateObjectID string,
	tokenConfigs map[string]TokenConfigView,
	poolDevInspect ITokenPoolDevInspect,
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
			lggr.Warnw("Failed to get remote pools", "chainSelector", chainSelector, "error", err)
			continue
		}

		// Get Remote Token
		remoteToken, err := poolDevInspect.GetRemoteToken(ctx, callOpts, typeArgs, poolStateObj, chainSelector)
		if err != nil {
			lggr.Warnw("Failed to get remote token", "chainSelector", chainSelector, "error", err)
			continue
		}

		// Get Inbound Rate Limiter State
		inboundRateLimiter, err := poolDevInspect.GetCurrentInboundRateLimiterState(ctx, callOpts, typeArgs, clockObj, poolStateObj, chainSelector)
		if err != nil {
			lggr.Warnw("Failed to get inbound rate limiter", "chainSelector", chainSelector, "error", err)
			// TODO: should we continue?
		}

		// Get Outbound Rate Limiter State
		outboundRateLimiter, err := poolDevInspect.GetCurrentOutboundRateLimiterState(ctx, callOpts, typeArgs, clockObj, poolStateObj, chainSelector)
		if err != nil {
			lggr.Warnw("Failed to get outbound rate limiter", "chainSelector", chainSelector, "error", err)
			// TODO: should we continue?
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
		},
		Token:              token,
		RemoteChainConfigs: remoteChainConfigs,
		AllowList:          allowlist,
		AllowListEnabled:   allowlistEnabled,
	}, nil
}

// TODO: rate limiter is coming as bind.Object instead of proper struct...
// parseRateLimiterConfig extracts rate limiter config from the bind.Object
func parseRateLimiterConfig(rateLimiterObj bind.Object) RateLimiterConfig {
	config := RateLimiterConfig{}

	return config
}
