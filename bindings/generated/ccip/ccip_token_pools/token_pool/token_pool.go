// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_token_pool

import (
	"context"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
	_ = uint256.NewInt
)

type ITokenPool interface {
	Initialize(coinMetadataAddress string, localDecimals byte, allowlist []string) bind.IMethod
	GetRouter() bind.IMethod
	GetToken(state TokenPoolState) bind.IMethod
	GetTokenDecimals(typeArgs string, coinMetadata bind.Object) bind.IMethod
	GetSupportedChains(state TokenPoolState) bind.IMethod
	IsSupportedChain(state TokenPoolState, remoteChainSelector uint64) bind.IMethod
	ApplyChainUpdates(state TokenPoolState, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) bind.IMethod
	GetRemotePools(state TokenPoolState, remoteChainSelector uint64) bind.IMethod
	IsRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	GetRemoteToken(state TokenPoolState, remoteChainSelector uint64) bind.IMethod
	AddRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	RemoveRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	EmitReleasedOrMinted(state TokenPoolState, recipient string, amount uint64, remoteChainSelector uint64) bind.IMethod
	EmitLockedOrBurned(state TokenPoolState, amount uint64, remoteChainSelector uint64) bind.IMethod
	EmitLiquidityAdded(state TokenPoolState, provider string, amount uint64) bind.IMethod
	EmitLiquidityRemoved(state TokenPoolState, provider string, amount uint64) bind.IMethod
	EmitRebalancerSet(state TokenPoolState, previousRebalancer string, rebalancer string) bind.IMethod
	GetLocalDecimals(pool TokenPoolState) bind.IMethod
	EncodeLocalDecimals(typeArgs string, coinMetadata bind.Object) bind.IMethod
	ParseRemoteDecimals(sourcePoolData []byte, localDecimals byte) bind.IMethod
	CalculateLocalAmount(remoteAmount uint256.Int, remoteDecimals byte, localDecimals byte) bind.IMethod
	CalculateReleaseOrMintAmount(state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) bind.IMethod
	SetChainRateLimiterConfig(clock bind.Object, state TokenPoolState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) bind.IMethod
	GetAllowlistEnabled(state TokenPoolState) bind.IMethod
	SetAllowlistEnabled(state TokenPoolState, enabled bool) bind.IMethod
	GetAllowlist(state TokenPoolState) bind.IMethod
	ApplyAllowlistUpdates(state TokenPoolState, removes []string, adds []string) bind.IMethod
	DestroyTokenPool(state TokenPoolState) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type TokenPoolContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ ITokenPool = (*TokenPoolContract)(nil)

func NewTokenPool(packageID string, client suiclient.ClientImpl) (*TokenPoolContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &TokenPoolContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *TokenPoolContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type TokenPoolState struct {
	CoinMetadata  string `move:"address"`
	LocalDecimals byte   `move:"u8"`
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

// Functions

func (c *TokenPoolContract) Initialize(coinMetadataAddress string, localDecimals byte, allowlist []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "initialize", false, "", "", coinMetadataAddress, localDecimals, allowlist)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetRouter() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_router", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_router", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetToken(state TokenPoolState) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_token", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_token", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetTokenDecimals(typeArgs string, coinMetadata bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_token_decimals", false, "", typeArgs, coinMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_token_decimals", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetSupportedChains(state TokenPoolState) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_supported_chains", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_supported_chains", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) IsSupportedChain(state TokenPoolState, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "is_supported_chain", false, "", "", state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "is_supported_chain", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) ApplyChainUpdates(state TokenPoolState, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "apply_chain_updates", false, "", "", state, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "apply_chain_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetRemotePools(state TokenPoolState, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_remote_pools", false, "", "", state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_remote_pools", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) IsRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "is_remote_pool", false, "", "", state, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "is_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetRemoteToken(state TokenPoolState, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_remote_token", false, "", "", state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_remote_token", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) AddRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "add_remote_pool", false, "", "", state, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "add_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) RemoveRemotePool(state TokenPoolState, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "remove_remote_pool", false, "", "", state, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "remove_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) EmitReleasedOrMinted(state TokenPoolState, recipient string, amount uint64, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "emit_released_or_minted", false, "", "", state, recipient, amount, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "emit_released_or_minted", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) EmitLockedOrBurned(state TokenPoolState, amount uint64, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "emit_locked_or_burned", false, "", "", state, amount, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "emit_locked_or_burned", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) EmitLiquidityAdded(state TokenPoolState, provider string, amount uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "emit_liquidity_added", false, "", "", state, provider, amount)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "emit_liquidity_added", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) EmitLiquidityRemoved(state TokenPoolState, provider string, amount uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "emit_liquidity_removed", false, "", "", state, provider, amount)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "emit_liquidity_removed", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) EmitRebalancerSet(state TokenPoolState, previousRebalancer string, rebalancer string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "emit_rebalancer_set", false, "", "", state, previousRebalancer, rebalancer)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "emit_rebalancer_set", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetLocalDecimals(pool TokenPoolState) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_local_decimals", false, "", "", pool)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_local_decimals", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) EncodeLocalDecimals(typeArgs string, coinMetadata bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "encode_local_decimals", false, "", typeArgs, coinMetadata)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "encode_local_decimals", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) ParseRemoteDecimals(sourcePoolData []byte, localDecimals byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "parse_remote_decimals", false, "", "", sourcePoolData, localDecimals)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "parse_remote_decimals", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) CalculateLocalAmount(remoteAmount uint256.Int, remoteDecimals byte, localDecimals byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "calculate_local_amount", false, "", "", remoteAmount, remoteDecimals, localDecimals)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "calculate_local_amount", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) CalculateReleaseOrMintAmount(state TokenPoolState, sourcePoolData []byte, sourceAmount uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "calculate_release_or_mint_amount", false, "", "", state, sourcePoolData, sourceAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "calculate_release_or_mint_amount", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) SetChainRateLimiterConfig(clock bind.Object, state TokenPoolState, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "set_chain_rate_limiter_config", false, "", "", clock, state, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "set_chain_rate_limiter_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetAllowlistEnabled(state TokenPoolState) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_allowlist_enabled", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_allowlist_enabled", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) SetAllowlistEnabled(state TokenPoolState, enabled bool) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "set_allowlist_enabled", false, "", "", state, enabled)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "set_allowlist_enabled", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) GetAllowlist(state TokenPoolState) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "get_allowlist", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "get_allowlist", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) ApplyAllowlistUpdates(state TokenPoolState, removes []string, adds []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "apply_allowlist_updates", false, "", "", state, removes, adds)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "apply_allowlist_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenPoolContract) DestroyTokenPool(state TokenPoolState) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_pool", "destroy_token_pool", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_pool", "destroy_token_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
