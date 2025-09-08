// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_token_pool

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

type ITokenPool interface {
	Initialize(ctx context.Context, opts *bind.CallOpts, coinMetadataAddress string, localDecimals byte, allowlist []string) (*models.SuiTransactionBlockResponse, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	AddRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error)
	ValidateLockOrBurn(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, state TokenPoolState, sender string, remoteChainSelector uint64, localAmount uint64) (*models.SuiTransactionBlockResponse, error)
	ValidateReleaseOrMint(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, state TokenPoolState, remoteChainSelector uint64, destTokenAddress string, sourcePoolAddress []byte, localAmount uint64) (*models.SuiTransactionBlockResponse, error)
	EmitReleasedOrMinted(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, recipient string, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitLockedOrBurned(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitLiquidityAdded(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, provider string, amount uint64) (*models.SuiTransactionBlockResponse, error)
	EmitLiquidityRemoved(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, provider string, amount uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRebalancerSet(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, previousRebalancer string, rebalancer string) (*models.SuiTransactionBlockResponse, error)
	GetLocalDecimals(ctx context.Context, opts *bind.CallOpts, pool TokenPoolState) (*models.SuiTransactionBlockResponse, error)
	EncodeLocalDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) (*models.SuiTransactionBlockResponse, error)
	ParseRemoteDecimals(ctx context.Context, opts *bind.CallOpts, sourcePoolData []byte, localDecimals byte) (*models.SuiTransactionBlockResponse, error)
	CalculateLocalAmount(ctx context.Context, opts *bind.CallOpts, remoteAmount *big.Int, remoteDecimals byte, localDecimals byte) (*models.SuiTransactionBlockResponse, error)
	CalculateReleaseOrMintAmount(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) (*models.SuiTransactionBlockResponse, error)
	SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state TokenPoolState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error)
	SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, enabled bool) (*models.SuiTransactionBlockResponse, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error)
	DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ITokenPoolDevInspect
	Encoder() TokenPoolEncoder
	Bound() bind.IBoundContract
}

type ITokenPoolDevInspect interface {
	Initialize(ctx context.Context, opts *bind.CallOpts, coinMetadataAddress string, localDecimals byte, allowlist []string) (TokenPoolState, error)
	GetToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (string, error)
	GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) (byte, error)
	GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) ([]uint64, error)
	IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (bool, error)
	GetRemotePools(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) ([][]byte, error)
	IsRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error)
	GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) ([]byte, error)
	ValidateLockOrBurn(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, state TokenPoolState, sender string, remoteChainSelector uint64, localAmount uint64) ([]byte, error)
	GetLocalDecimals(ctx context.Context, opts *bind.CallOpts, pool TokenPoolState) (byte, error)
	EncodeLocalDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) ([]byte, error)
	ParseRemoteDecimals(ctx context.Context, opts *bind.CallOpts, sourcePoolData []byte, localDecimals byte) (byte, error)
	CalculateLocalAmount(ctx context.Context, opts *bind.CallOpts, remoteAmount *big.Int, remoteDecimals byte, localDecimals byte) (uint64, error)
	CalculateReleaseOrMintAmount(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) (uint64, error)
	GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (bool, error)
	GetAllowlist(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) ([]string, error)
}

type TokenPoolEncoder interface {
	Initialize(coinMetadataAddress string, localDecimals byte, allowlist []string) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	GetToken(state TokenPoolState) (*bind.EncodedCall, error)
	GetTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenDecimals(typeArgs []string, coinMetadata bind.Object) (*bind.EncodedCall, error)
	GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetSupportedChains(state TokenPoolState) (*bind.EncodedCall, error)
	GetSupportedChainsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsSupportedChain(state TokenPoolState, remoteChainSelector uint64) (*bind.EncodedCall, error)
	IsSupportedChainWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyChainUpdates(state TokenPoolState, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error)
	ApplyChainUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRemotePools(state TokenPoolState, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemotePoolsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	IsRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRemoteToken(state TokenPoolState, remoteChainSelector uint64) (*bind.EncodedCall, error)
	GetRemoteTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	AddRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	AddRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error)
	RemoveRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error)
	RemoveRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error)
	ValidateLockOrBurn(ref bind.Object, clock bind.Object, state TokenPoolState, sender string, remoteChainSelector uint64, localAmount uint64) (*bind.EncodedCall, error)
	ValidateLockOrBurnWithArgs(args ...any) (*bind.EncodedCall, error)
	ValidateReleaseOrMint(ref bind.Object, clock bind.Object, state TokenPoolState, remoteChainSelector uint64, destTokenAddress string, sourcePoolAddress []byte, localAmount uint64) (*bind.EncodedCall, error)
	ValidateReleaseOrMintWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitReleasedOrMinted(state TokenPoolState, recipient string, amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitReleasedOrMintedWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitLockedOrBurned(state TokenPoolState, amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitLockedOrBurnedWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitLiquidityAdded(state TokenPoolState, provider string, amount uint64) (*bind.EncodedCall, error)
	EmitLiquidityAddedWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitLiquidityRemoved(state TokenPoolState, provider string, amount uint64) (*bind.EncodedCall, error)
	EmitLiquidityRemovedWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitRebalancerSet(state TokenPoolState, previousRebalancer string, rebalancer string) (*bind.EncodedCall, error)
	EmitRebalancerSetWithArgs(args ...any) (*bind.EncodedCall, error)
	GetLocalDecimals(pool TokenPoolState) (*bind.EncodedCall, error)
	GetLocalDecimalsWithArgs(args ...any) (*bind.EncodedCall, error)
	EncodeLocalDecimals(typeArgs []string, coinMetadata bind.Object) (*bind.EncodedCall, error)
	EncodeLocalDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ParseRemoteDecimals(sourcePoolData []byte, localDecimals byte) (*bind.EncodedCall, error)
	ParseRemoteDecimalsWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateLocalAmount(remoteAmount *big.Int, remoteDecimals byte, localDecimals byte) (*bind.EncodedCall, error)
	CalculateLocalAmountWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateReleaseOrMintAmount(state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) (*bind.EncodedCall, error)
	CalculateReleaseOrMintAmountWithArgs(args ...any) (*bind.EncodedCall, error)
	SetChainRateLimiterConfig(clock bind.Object, state TokenPoolState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error)
	SetChainRateLimiterConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowlistEnabled(state TokenPoolState) (*bind.EncodedCall, error)
	GetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error)
	SetAllowlistEnabled(state TokenPoolState, enabled bool) (*bind.EncodedCall, error)
	SetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowlist(state TokenPoolState) (*bind.EncodedCall, error)
	GetAllowlistWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyAllowlistUpdates(state TokenPoolState, removes []string, adds []string) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	DestroyTokenPool(state TokenPoolState) (*bind.EncodedCall, error)
	DestroyTokenPoolWithArgs(args ...any) (*bind.EncodedCall, error)
}

type TokenPoolContract struct {
	*bind.BoundContract
	tokenPoolEncoder
	devInspect *TokenPoolDevInspect
}

type TokenPoolDevInspect struct {
	contract *TokenPoolContract
}

var _ ITokenPool = (*TokenPoolContract)(nil)
var _ ITokenPoolDevInspect = (*TokenPoolDevInspect)(nil)

func NewTokenPool(packageID string, client sui.ISuiAPI) (ITokenPool, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_token_pool", "token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &TokenPoolContract{
		BoundContract:    contract,
		tokenPoolEncoder: tokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &TokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *TokenPoolContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *TokenPoolContract) Encoder() TokenPoolEncoder {
	return c.tokenPoolEncoder
}

func (c *TokenPoolContract) DevInspect() ITokenPoolDevInspect {
	return c.devInspect
}

type TokenPoolState struct {
	AllowlistState     bind.Object `move:"allowlist::AllowlistState"`
	CoinMetadata       string      `move:"address"`
	LocalDecimals      byte        `move:"u8"`
	RemoteChainConfigs bind.Object `move:"VecMap<u64, RemoteChainConfig>"`
	RateLimiterConfig  bind.Object `move:"token_pool_rate_limiter::RateLimitState"`
}

type RemoteChainConfig struct {
	RemoteTokenAddress []byte   `move:"vector<u8>"`
	RemotePools        [][]byte `move:"vector<vector<u8>>"`
}

type LockedOrBurned struct {
	RemoteChainSelector uint64 `move:"u64"`
	LocalToken          string `move:"address"`
	Amount              uint64 `move:"u64"`
}

type ReleasedOrMinted struct {
	RemoteChainSelector uint64 `move:"u64"`
	LocalToken          string `move:"address"`
	Recipient           string `move:"address"`
	Amount              uint64 `move:"u64"`
}

type RemotePoolAdded struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemotePoolAddress   []byte `move:"vector<u8>"`
}

type RemotePoolRemoved struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemotePoolAddress   []byte `move:"vector<u8>"`
}

type ChainAdded struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemoteTokenAddress  []byte `move:"vector<u8>"`
}

type ChainRemoved struct {
	RemoteChainSelector uint64 `move:"u64"`
}

type LiquidityAdded struct {
	LocalToken string `move:"address"`
	Provider   string `move:"address"`
	Amount     uint64 `move:"u64"`
}

type LiquidityRemoved struct {
	LocalToken string `move:"address"`
	Provider   string `move:"address"`
	Amount     uint64 `move:"u64"`
}

type RebalancerSet struct {
	LocalToken         string `move:"address"`
	PreviousRebalancer string `move:"address"`
	Rebalancer         string `move:"address"`
}

type bcsTokenPoolState struct {
	AllowlistState     bind.Object
	CoinMetadata       [32]byte
	LocalDecimals      byte
	RemoteChainConfigs bind.Object
	RateLimiterConfig  bind.Object
}

func convertTokenPoolStateFromBCS(bcs bcsTokenPoolState) (TokenPoolState, error) {

	return TokenPoolState{
		AllowlistState:     bcs.AllowlistState,
		CoinMetadata:       fmt.Sprintf("0x%x", bcs.CoinMetadata),
		LocalDecimals:      bcs.LocalDecimals,
		RemoteChainConfigs: bcs.RemoteChainConfigs,
		RateLimiterConfig:  bcs.RateLimiterConfig,
	}, nil
}

type bcsLockedOrBurned struct {
	RemoteChainSelector uint64
	LocalToken          [32]byte
	Amount              uint64
}

func convertLockedOrBurnedFromBCS(bcs bcsLockedOrBurned) (LockedOrBurned, error) {

	return LockedOrBurned{
		RemoteChainSelector: bcs.RemoteChainSelector,
		LocalToken:          fmt.Sprintf("0x%x", bcs.LocalToken),
		Amount:              bcs.Amount,
	}, nil
}

type bcsReleasedOrMinted struct {
	RemoteChainSelector uint64
	LocalToken          [32]byte
	Recipient           [32]byte
	Amount              uint64
}

func convertReleasedOrMintedFromBCS(bcs bcsReleasedOrMinted) (ReleasedOrMinted, error) {

	return ReleasedOrMinted{
		RemoteChainSelector: bcs.RemoteChainSelector,
		LocalToken:          fmt.Sprintf("0x%x", bcs.LocalToken),
		Recipient:           fmt.Sprintf("0x%x", bcs.Recipient),
		Amount:              bcs.Amount,
	}, nil
}

type bcsLiquidityAdded struct {
	LocalToken [32]byte
	Provider   [32]byte
	Amount     uint64
}

func convertLiquidityAddedFromBCS(bcs bcsLiquidityAdded) (LiquidityAdded, error) {

	return LiquidityAdded{
		LocalToken: fmt.Sprintf("0x%x", bcs.LocalToken),
		Provider:   fmt.Sprintf("0x%x", bcs.Provider),
		Amount:     bcs.Amount,
	}, nil
}

type bcsLiquidityRemoved struct {
	LocalToken [32]byte
	Provider   [32]byte
	Amount     uint64
}

func convertLiquidityRemovedFromBCS(bcs bcsLiquidityRemoved) (LiquidityRemoved, error) {

	return LiquidityRemoved{
		LocalToken: fmt.Sprintf("0x%x", bcs.LocalToken),
		Provider:   fmt.Sprintf("0x%x", bcs.Provider),
		Amount:     bcs.Amount,
	}, nil
}

type bcsRebalancerSet struct {
	LocalToken         [32]byte
	PreviousRebalancer [32]byte
	Rebalancer         [32]byte
}

func convertRebalancerSetFromBCS(bcs bcsRebalancerSet) (RebalancerSet, error) {

	return RebalancerSet{
		LocalToken:         fmt.Sprintf("0x%x", bcs.LocalToken),
		PreviousRebalancer: fmt.Sprintf("0x%x", bcs.PreviousRebalancer),
		Rebalancer:         fmt.Sprintf("0x%x", bcs.Rebalancer),
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::TokenPoolState", func(data []byte) (interface{}, error) {
		var temp bcsTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenPoolStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::RemoteChainConfig", func(data []byte) (interface{}, error) {
		var result RemoteChainConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::LockedOrBurned", func(data []byte) (interface{}, error) {
		var temp bcsLockedOrBurned
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertLockedOrBurnedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::ReleasedOrMinted", func(data []byte) (interface{}, error) {
		var temp bcsReleasedOrMinted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertReleasedOrMintedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::RemotePoolAdded", func(data []byte) (interface{}, error) {
		var result RemotePoolAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::RemotePoolRemoved", func(data []byte) (interface{}, error) {
		var result RemotePoolRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::ChainAdded", func(data []byte) (interface{}, error) {
		var result ChainAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::ChainRemoved", func(data []byte) (interface{}, error) {
		var result ChainRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::LiquidityAdded", func(data []byte) (interface{}, error) {
		var temp bcsLiquidityAdded
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertLiquidityAddedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::LiquidityRemoved", func(data []byte) (interface{}, error) {
		var temp bcsLiquidityRemoved
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertLiquidityRemovedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_token_pool::token_pool::RebalancerSet", func(data []byte) (interface{}, error) {
		var temp bcsRebalancerSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertRebalancerSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// Initialize executes the initialize Move function.
func (c *TokenPoolContract) Initialize(ctx context.Context, opts *bind.CallOpts, coinMetadataAddress string, localDecimals byte, allowlist []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.Initialize(coinMetadataAddress, localDecimals, allowlist)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetToken executes the get_token Move function.
func (c *TokenPoolContract) GetToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetToken(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenDecimals executes the get_token_decimals Move function.
func (c *TokenPoolContract) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetTokenDecimals(typeArgs, coinMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSupportedChains executes the get_supported_chains Move function.
func (c *TokenPoolContract) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetSupportedChains(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsSupportedChain executes the is_supported_chain Move function.
func (c *TokenPoolContract) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.IsSupportedChain(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyChainUpdates executes the apply_chain_updates Move function.
func (c *TokenPoolContract) ApplyChainUpdates(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.ApplyChainUpdates(state, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemotePools executes the get_remote_pools Move function.
func (c *TokenPoolContract) GetRemotePools(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetRemotePools(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsRemotePool executes the is_remote_pool Move function.
func (c *TokenPoolContract) IsRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.IsRemotePool(state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetRemoteToken executes the get_remote_token Move function.
func (c *TokenPoolContract) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetRemoteToken(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddRemotePool executes the add_remote_pool Move function.
func (c *TokenPoolContract) AddRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.AddRemotePool(state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemoveRemotePool executes the remove_remote_pool Move function.
func (c *TokenPoolContract) RemoveRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.RemoveRemotePool(state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ValidateLockOrBurn executes the validate_lock_or_burn Move function.
func (c *TokenPoolContract) ValidateLockOrBurn(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, state TokenPoolState, sender string, remoteChainSelector uint64, localAmount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.ValidateLockOrBurn(ref, clock, state, sender, remoteChainSelector, localAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ValidateReleaseOrMint executes the validate_release_or_mint Move function.
func (c *TokenPoolContract) ValidateReleaseOrMint(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, state TokenPoolState, remoteChainSelector uint64, destTokenAddress string, sourcePoolAddress []byte, localAmount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.ValidateReleaseOrMint(ref, clock, state, remoteChainSelector, destTokenAddress, sourcePoolAddress, localAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitReleasedOrMinted executes the emit_released_or_minted Move function.
func (c *TokenPoolContract) EmitReleasedOrMinted(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, recipient string, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.EmitReleasedOrMinted(state, recipient, amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitLockedOrBurned executes the emit_locked_or_burned Move function.
func (c *TokenPoolContract) EmitLockedOrBurned(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.EmitLockedOrBurned(state, amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitLiquidityAdded executes the emit_liquidity_added Move function.
func (c *TokenPoolContract) EmitLiquidityAdded(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, provider string, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.EmitLiquidityAdded(state, provider, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitLiquidityRemoved executes the emit_liquidity_removed Move function.
func (c *TokenPoolContract) EmitLiquidityRemoved(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, provider string, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.EmitLiquidityRemoved(state, provider, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRebalancerSet executes the emit_rebalancer_set Move function.
func (c *TokenPoolContract) EmitRebalancerSet(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, previousRebalancer string, rebalancer string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.EmitRebalancerSet(state, previousRebalancer, rebalancer)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetLocalDecimals executes the get_local_decimals Move function.
func (c *TokenPoolContract) GetLocalDecimals(ctx context.Context, opts *bind.CallOpts, pool TokenPoolState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetLocalDecimals(pool)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EncodeLocalDecimals executes the encode_local_decimals Move function.
func (c *TokenPoolContract) EncodeLocalDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.EncodeLocalDecimals(typeArgs, coinMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ParseRemoteDecimals executes the parse_remote_decimals Move function.
func (c *TokenPoolContract) ParseRemoteDecimals(ctx context.Context, opts *bind.CallOpts, sourcePoolData []byte, localDecimals byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.ParseRemoteDecimals(sourcePoolData, localDecimals)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateLocalAmount executes the calculate_local_amount Move function.
func (c *TokenPoolContract) CalculateLocalAmount(ctx context.Context, opts *bind.CallOpts, remoteAmount *big.Int, remoteDecimals byte, localDecimals byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.CalculateLocalAmount(remoteAmount, remoteDecimals, localDecimals)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateReleaseOrMintAmount executes the calculate_release_or_mint_amount Move function.
func (c *TokenPoolContract) CalculateReleaseOrMintAmount(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.CalculateReleaseOrMintAmount(state, sourcePoolData, sourceAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetChainRateLimiterConfig executes the set_chain_rate_limiter_config Move function.
func (c *TokenPoolContract) SetChainRateLimiterConfig(ctx context.Context, opts *bind.CallOpts, clock bind.Object, state TokenPoolState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.SetChainRateLimiterConfig(clock, state, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function.
func (c *TokenPoolContract) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetAllowlistEnabled(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetAllowlistEnabled executes the set_allowlist_enabled Move function.
func (c *TokenPoolContract) SetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, enabled bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.SetAllowlistEnabled(state, enabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowlist executes the get_allowlist Move function.
func (c *TokenPoolContract) GetAllowlist(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.GetAllowlist(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *TokenPoolContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, removes []string, adds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.ApplyAllowlistUpdates(state, removes, adds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DestroyTokenPool executes the destroy_token_pool Move function.
func (c *TokenPoolContract) DestroyTokenPool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.tokenPoolEncoder.DestroyTokenPool(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function using DevInspect to get return values.
//
// Returns: TokenPoolState
func (d *TokenPoolDevInspect) Initialize(ctx context.Context, opts *bind.CallOpts, coinMetadataAddress string, localDecimals byte, allowlist []string) (TokenPoolState, error) {
	encoded, err := d.contract.tokenPoolEncoder.Initialize(coinMetadataAddress, localDecimals, allowlist)
	if err != nil {
		return TokenPoolState{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TokenPoolState{}, err
	}
	if len(results) == 0 {
		return TokenPoolState{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TokenPoolState)
	if !ok {
		return TokenPoolState{}, fmt.Errorf("unexpected return type: expected TokenPoolState, got %T", results[0])
	}
	return result, nil
}

// GetToken executes the get_token Move function using DevInspect to get return values.
//
// Returns: address
func (d *TokenPoolDevInspect) GetToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (string, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetToken(state)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// GetTokenDecimals executes the get_token_decimals Move function using DevInspect to get return values.
//
// Returns: u8
func (d *TokenPoolDevInspect) GetTokenDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) (byte, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetTokenDecimals(typeArgs, coinMetadata)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// GetSupportedChains executes the get_supported_chains Move function using DevInspect to get return values.
//
// Returns: vector<u64>
func (d *TokenPoolDevInspect) GetSupportedChains(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) ([]uint64, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetSupportedChains(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]uint64)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []uint64, got %T", results[0])
	}
	return result, nil
}

// IsSupportedChain executes the is_supported_chain Move function using DevInspect to get return values.
//
// Returns: bool
func (d *TokenPoolDevInspect) IsSupportedChain(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) (bool, error) {
	encoded, err := d.contract.tokenPoolEncoder.IsSupportedChain(state, remoteChainSelector)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// GetRemotePools executes the get_remote_pools Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *TokenPoolDevInspect) GetRemotePools(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) ([][]byte, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetRemotePools(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([][]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected [][]byte, got %T", results[0])
	}
	return result, nil
}

// IsRemotePool executes the is_remote_pool Move function using DevInspect to get return values.
//
// Returns: bool
func (d *TokenPoolDevInspect) IsRemotePool(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (bool, error) {
	encoded, err := d.contract.tokenPoolEncoder.IsRemotePool(state, remoteChainSelector, remotePoolAddress)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// GetRemoteToken executes the get_remote_token Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *TokenPoolDevInspect) GetRemoteToken(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, remoteChainSelector uint64) ([]byte, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetRemoteToken(state, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// ValidateLockOrBurn executes the validate_lock_or_burn Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *TokenPoolDevInspect) ValidateLockOrBurn(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, state TokenPoolState, sender string, remoteChainSelector uint64, localAmount uint64) ([]byte, error) {
	encoded, err := d.contract.tokenPoolEncoder.ValidateLockOrBurn(ref, clock, state, sender, remoteChainSelector, localAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// GetLocalDecimals executes the get_local_decimals Move function using DevInspect to get return values.
//
// Returns: u8
func (d *TokenPoolDevInspect) GetLocalDecimals(ctx context.Context, opts *bind.CallOpts, pool TokenPoolState) (byte, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetLocalDecimals(pool)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// EncodeLocalDecimals executes the encode_local_decimals Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *TokenPoolDevInspect) EncodeLocalDecimals(ctx context.Context, opts *bind.CallOpts, typeArgs []string, coinMetadata bind.Object) ([]byte, error) {
	encoded, err := d.contract.tokenPoolEncoder.EncodeLocalDecimals(typeArgs, coinMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// ParseRemoteDecimals executes the parse_remote_decimals Move function using DevInspect to get return values.
//
// Returns: u8
func (d *TokenPoolDevInspect) ParseRemoteDecimals(ctx context.Context, opts *bind.CallOpts, sourcePoolData []byte, localDecimals byte) (byte, error) {
	encoded, err := d.contract.tokenPoolEncoder.ParseRemoteDecimals(sourcePoolData, localDecimals)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// CalculateLocalAmount executes the calculate_local_amount Move function using DevInspect to get return values.
//
// Returns: u64
func (d *TokenPoolDevInspect) CalculateLocalAmount(ctx context.Context, opts *bind.CallOpts, remoteAmount *big.Int, remoteDecimals byte, localDecimals byte) (uint64, error) {
	encoded, err := d.contract.tokenPoolEncoder.CalculateLocalAmount(remoteAmount, remoteDecimals, localDecimals)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// CalculateReleaseOrMintAmount executes the calculate_release_or_mint_amount Move function using DevInspect to get return values.
//
// Returns: u64
func (d *TokenPoolDevInspect) CalculateReleaseOrMintAmount(ctx context.Context, opts *bind.CallOpts, state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) (uint64, error) {
	encoded, err := d.contract.tokenPoolEncoder.CalculateReleaseOrMintAmount(state, sourcePoolData, sourceAmount)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// GetAllowlistEnabled executes the get_allowlist_enabled Move function using DevInspect to get return values.
//
// Returns: bool
func (d *TokenPoolDevInspect) GetAllowlistEnabled(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) (bool, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetAllowlistEnabled(state)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// GetAllowlist executes the get_allowlist Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *TokenPoolDevInspect) GetAllowlist(ctx context.Context, opts *bind.CallOpts, state TokenPoolState) ([]string, error) {
	encoded, err := d.contract.tokenPoolEncoder.GetAllowlist(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []string, got %T", results[0])
	}
	return result, nil
}

type tokenPoolEncoder struct {
	*bind.BoundContract
}

// Initialize encodes a call to the initialize Move function.
func (c tokenPoolEncoder) Initialize(coinMetadataAddress string, localDecimals byte, allowlist []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"address",
		"u8",
		"vector<address>",
	}, []any{
		coinMetadataAddress,
		localDecimals,
		allowlist,
	}, []string{
		"ccip_token_pool::token_pool::TokenPoolState",
	})
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"address",
		"u8",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_token_pool::token_pool::TokenPoolState",
	})
}

// GetToken encodes a call to the get_token Move function.
func (c tokenPoolEncoder) GetToken(state TokenPoolState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetTokenWithArgs encodes a call to the get_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetTokenDecimals encodes a call to the get_token_decimals Move function.
func (c tokenPoolEncoder) GetTokenDecimals(typeArgs []string, coinMetadata bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, []string{
		"&CoinMetadata<T>",
	}, []any{
		coinMetadata,
	}, []string{
		"u8",
	})
}

// GetTokenDecimalsWithArgs encodes a call to the get_token_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetTokenDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CoinMetadata<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_token_decimals", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// GetSupportedChains encodes a call to the get_supported_chains Move function.
func (c tokenPoolEncoder) GetSupportedChains(state TokenPoolState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
	}, []any{
		state,
	}, []string{
		"vector<u64>",
	})
}

// GetSupportedChainsWithArgs encodes a call to the get_supported_chains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetSupportedChainsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_supported_chains", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
	})
}

// IsSupportedChain encodes a call to the is_supported_chain Move function.
func (c tokenPoolEncoder) IsSupportedChain(state TokenPoolState, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
		"u64",
	}, []any{
		state,
		remoteChainSelector,
	}, []string{
		"bool",
	})
}

// IsSupportedChainWithArgs encodes a call to the is_supported_chain Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) IsSupportedChainWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_supported_chain", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// ApplyChainUpdates encodes a call to the apply_chain_updates Move function.
func (c tokenPoolEncoder) ApplyChainUpdates(state TokenPoolState, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"vector<u64>",
		"vector<u64>",
		"vector<vector<vector<u8>>>",
		"vector<vector<u8>>",
	}, []any{
		state,
		remoteChainSelectorsToRemove,
		remoteChainSelectorsToAdd,
		remotePoolAddressesToAdd,
		remoteTokenAddressesToAdd,
	}, nil)
}

// ApplyChainUpdatesWithArgs encodes a call to the apply_chain_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) ApplyChainUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"vector<u64>",
		"vector<u64>",
		"vector<vector<vector<u8>>>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_chain_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetRemotePools encodes a call to the get_remote_pools Move function.
func (c tokenPoolEncoder) GetRemotePools(state TokenPoolState, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
		"u64",
	}, []any{
		state,
		remoteChainSelector,
	}, []string{
		"vector<vector<u8>>",
	})
}

// GetRemotePoolsWithArgs encodes a call to the get_remote_pools Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetRemotePoolsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_pools", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
	})
}

// IsRemotePool encodes a call to the is_remote_pool Move function.
func (c tokenPoolEncoder) IsRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
		"u64",
		"vector<u8>",
	}, []any{
		state,
		remoteChainSelector,
		remotePoolAddress,
	}, []string{
		"bool",
	})
}

// IsRemotePoolWithArgs encodes a call to the is_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) IsRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_remote_pool", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetRemoteToken encodes a call to the get_remote_token Move function.
func (c tokenPoolEncoder) GetRemoteToken(state TokenPoolState, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
		"u64",
	}, []any{
		state,
		remoteChainSelector,
	}, []string{
		"vector<u8>",
	})
}

// GetRemoteTokenWithArgs encodes a call to the get_remote_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetRemoteTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_remote_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// AddRemotePool encodes a call to the add_remote_pool Move function.
func (c tokenPoolEncoder) AddRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"u64",
		"vector<u8>",
	}, []any{
		state,
		remoteChainSelector,
		remotePoolAddress,
	}, nil)
}

// AddRemotePoolWithArgs encodes a call to the add_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) AddRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RemoveRemotePool encodes a call to the remove_remote_pool Move function.
func (c tokenPoolEncoder) RemoveRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"u64",
		"vector<u8>",
	}, []any{
		state,
		remoteChainSelector,
		remotePoolAddress,
	}, nil)
}

// RemoveRemotePoolWithArgs encodes a call to the remove_remote_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) RemoveRemotePoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_remote_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ValidateLockOrBurn encodes a call to the validate_lock_or_burn Move function.
func (c tokenPoolEncoder) ValidateLockOrBurn(ref bind.Object, clock bind.Object, state TokenPoolState, sender string, remoteChainSelector uint64, localAmount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("validate_lock_or_burn", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&Clock",
		"&mut TokenPoolState",
		"address",
		"u64",
		"u64",
	}, []any{
		ref,
		clock,
		state,
		sender,
		remoteChainSelector,
		localAmount,
	}, []string{
		"vector<u8>",
	})
}

// ValidateLockOrBurnWithArgs encodes a call to the validate_lock_or_burn Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) ValidateLockOrBurnWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&Clock",
		"&mut TokenPoolState",
		"address",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("validate_lock_or_burn", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// ValidateReleaseOrMint encodes a call to the validate_release_or_mint Move function.
func (c tokenPoolEncoder) ValidateReleaseOrMint(ref bind.Object, clock bind.Object, state TokenPoolState, remoteChainSelector uint64, destTokenAddress string, sourcePoolAddress []byte, localAmount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("validate_release_or_mint", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&Clock",
		"&mut TokenPoolState",
		"u64",
		"address",
		"vector<u8>",
		"u64",
	}, []any{
		ref,
		clock,
		state,
		remoteChainSelector,
		destTokenAddress,
		sourcePoolAddress,
		localAmount,
	}, nil)
}

// ValidateReleaseOrMintWithArgs encodes a call to the validate_release_or_mint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) ValidateReleaseOrMintWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&Clock",
		"&mut TokenPoolState",
		"u64",
		"address",
		"vector<u8>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("validate_release_or_mint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitReleasedOrMinted encodes a call to the emit_released_or_minted Move function.
func (c tokenPoolEncoder) EmitReleasedOrMinted(state TokenPoolState, recipient string, amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_released_or_minted", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"address",
		"u64",
		"u64",
	}, []any{
		state,
		recipient,
		amount,
		remoteChainSelector,
	}, nil)
}

// EmitReleasedOrMintedWithArgs encodes a call to the emit_released_or_minted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) EmitReleasedOrMintedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"address",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_released_or_minted", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitLockedOrBurned encodes a call to the emit_locked_or_burned Move function.
func (c tokenPoolEncoder) EmitLockedOrBurned(state TokenPoolState, amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_locked_or_burned", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"u64",
		"u64",
	}, []any{
		state,
		amount,
		remoteChainSelector,
	}, nil)
}

// EmitLockedOrBurnedWithArgs encodes a call to the emit_locked_or_burned Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) EmitLockedOrBurnedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_locked_or_burned", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitLiquidityAdded encodes a call to the emit_liquidity_added Move function.
func (c tokenPoolEncoder) EmitLiquidityAdded(state TokenPoolState, provider string, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_added", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"address",
		"u64",
	}, []any{
		state,
		provider,
		amount,
	}, nil)
}

// EmitLiquidityAddedWithArgs encodes a call to the emit_liquidity_added Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) EmitLiquidityAddedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"address",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_added", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitLiquidityRemoved encodes a call to the emit_liquidity_removed Move function.
func (c tokenPoolEncoder) EmitLiquidityRemoved(state TokenPoolState, provider string, amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_removed", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"address",
		"u64",
	}, []any{
		state,
		provider,
		amount,
	}, nil)
}

// EmitLiquidityRemovedWithArgs encodes a call to the emit_liquidity_removed Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) EmitLiquidityRemovedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"address",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_removed", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitRebalancerSet encodes a call to the emit_rebalancer_set Move function.
func (c tokenPoolEncoder) EmitRebalancerSet(state TokenPoolState, previousRebalancer string, rebalancer string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_rebalancer_set", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"address",
		"address",
	}, []any{
		state,
		previousRebalancer,
		rebalancer,
	}, nil)
}

// EmitRebalancerSetWithArgs encodes a call to the emit_rebalancer_set Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) EmitRebalancerSetWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_rebalancer_set", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetLocalDecimals encodes a call to the get_local_decimals Move function.
func (c tokenPoolEncoder) GetLocalDecimals(pool TokenPoolState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_local_decimals", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
	}, []any{
		pool,
	}, []string{
		"u8",
	})
}

// GetLocalDecimalsWithArgs encodes a call to the get_local_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetLocalDecimalsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_local_decimals", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// EncodeLocalDecimals encodes a call to the encode_local_decimals Move function.
func (c tokenPoolEncoder) EncodeLocalDecimals(typeArgs []string, coinMetadata bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("encode_local_decimals", typeArgsList, typeParamsList, []string{
		"&CoinMetadata<T>",
	}, []any{
		coinMetadata,
	}, []string{
		"vector<u8>",
	})
}

// EncodeLocalDecimalsWithArgs encodes a call to the encode_local_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) EncodeLocalDecimalsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CoinMetadata<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("encode_local_decimals", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// ParseRemoteDecimals encodes a call to the parse_remote_decimals Move function.
func (c tokenPoolEncoder) ParseRemoteDecimals(sourcePoolData []byte, localDecimals byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("parse_remote_decimals", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"u8",
	}, []any{
		sourcePoolData,
		localDecimals,
	}, []string{
		"u8",
	})
}

// ParseRemoteDecimalsWithArgs encodes a call to the parse_remote_decimals Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) ParseRemoteDecimalsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("parse_remote_decimals", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// CalculateLocalAmount encodes a call to the calculate_local_amount Move function.
func (c tokenPoolEncoder) CalculateLocalAmount(remoteAmount *big.Int, remoteDecimals byte, localDecimals byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_local_amount", typeArgsList, typeParamsList, []string{
		"u256",
		"u8",
		"u8",
	}, []any{
		remoteAmount,
		remoteDecimals,
		localDecimals,
	}, []string{
		"u64",
	})
}

// CalculateLocalAmountWithArgs encodes a call to the calculate_local_amount Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) CalculateLocalAmountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u256",
		"u8",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_local_amount", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// CalculateReleaseOrMintAmount encodes a call to the calculate_release_or_mint_amount Move function.
func (c tokenPoolEncoder) CalculateReleaseOrMintAmount(state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_release_or_mint_amount", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
		"vector<u8>",
		"u64",
	}, []any{
		state,
		sourcePoolData,
		sourceAmount,
	}, []string{
		"u64",
	})
}

// CalculateReleaseOrMintAmountWithArgs encodes a call to the calculate_release_or_mint_amount Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) CalculateReleaseOrMintAmountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
		"vector<u8>",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_release_or_mint_amount", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// SetChainRateLimiterConfig encodes a call to the set_chain_rate_limiter_config Move function.
func (c tokenPoolEncoder) SetChainRateLimiterConfig(clock bind.Object, state TokenPoolState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, []string{
		"&Clock",
		"&mut TokenPoolState",
		"u64",
		"bool",
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	}, []any{
		clock,
		state,
		remoteChainSelector,
		outboundIsEnabled,
		outboundCapacity,
		outboundRate,
		inboundIsEnabled,
		inboundCapacity,
		inboundRate,
	}, nil)
}

// SetChainRateLimiterConfigWithArgs encodes a call to the set_chain_rate_limiter_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) SetChainRateLimiterConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Clock",
		"&mut TokenPoolState",
		"u64",
		"bool",
		"u64",
		"u64",
		"bool",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_chain_rate_limiter_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetAllowlistEnabled encodes a call to the get_allowlist_enabled Move function.
func (c tokenPoolEncoder) GetAllowlistEnabled(state TokenPoolState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// GetAllowlistEnabledWithArgs encodes a call to the get_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist_enabled", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// SetAllowlistEnabled encodes a call to the set_allowlist_enabled Move function.
func (c tokenPoolEncoder) SetAllowlistEnabled(state TokenPoolState, enabled bool) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"bool",
	}, []any{
		state,
		enabled,
	}, nil)
}

// SetAllowlistEnabledWithArgs encodes a call to the set_allowlist_enabled Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) SetAllowlistEnabledWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"bool",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_allowlist_enabled", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetAllowlist encodes a call to the get_allowlist Move function.
func (c tokenPoolEncoder) GetAllowlist(state TokenPoolState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, []string{
		"&TokenPoolState",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// GetAllowlistWithArgs encodes a call to the get_allowlist Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) GetAllowlistWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowlist", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// ApplyAllowlistUpdates encodes a call to the apply_allowlist_updates Move function.
func (c tokenPoolEncoder) ApplyAllowlistUpdates(state TokenPoolState, removes []string, adds []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut TokenPoolState",
		"vector<address>",
		"vector<address>",
	}, []any{
		state,
		removes,
		adds,
	}, nil)
}

// ApplyAllowlistUpdatesWithArgs encodes a call to the apply_allowlist_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut TokenPoolState",
		"vector<address>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// DestroyTokenPool encodes a call to the destroy_token_pool Move function.
func (c tokenPoolEncoder) DestroyTokenPool(state TokenPoolState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, []string{
		"ccip_token_pool::token_pool::TokenPoolState",
	}, []any{
		state,
	}, nil)
}

// DestroyTokenPoolWithArgs encodes a call to the destroy_token_pool Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c tokenPoolEncoder) DestroyTokenPoolWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_token_pool::token_pool::TokenPoolState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("destroy_token_pool", typeArgsList, typeParamsList, expectedParams, args, nil)
}
