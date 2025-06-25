// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_lock_release_token_pool

import (
	"context"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_common "github.com/smartcontractkit/chainlink-sui/bindings/common"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
	_ = uint256.NewInt
)

type ILockReleaseTokenPool interface {
	TypeAndVersion() bind.IMethod
	Initialize(typeArgs string, ref module_common.CCIPObjectRef, coinMetadata bind.Object, treasuryCap bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) bind.IMethod
	InitializeByCcipAdmin(typeArgs string, ref module_common.CCIPObjectRef, coinMetadata bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) bind.IMethod
	GetToken(state bind.Object) bind.IMethod
	GetTokenDecimals(state bind.Object) bind.IMethod
	GetRemotePools(state bind.Object, remoteChainSelector uint64) bind.IMethod
	IsRemotePool(state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	GetRemoteToken(state bind.Object, remoteChainSelector uint64) bind.IMethod
	AddRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	RemoveRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	IsSupportedChain(state bind.Object, remoteChainSelector uint64) bind.IMethod
	GetSupportedChains(state bind.Object) bind.IMethod
	ApplyChainUpdates(state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) bind.IMethod
	GetAllowlistEnabled(state bind.Object) bind.IMethod
	GetAllowlist(state bind.Object) bind.IMethod
	SetAllowlistEnabled(state bind.Object, ownerCap bind.Object, enabled bool) bind.IMethod
	ApplyAllowlistUpdates(state bind.Object, ownerCap bind.Object, removes []string, adds []string) bind.IMethod
	LockOrBurn(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, state bind.Object, c_ bind.Object, tokenParams module_common.TokenParams) bind.IMethod
	ReleaseOrMint(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, pool bind.Object, receiverParams module_common.ReceiverParams, index uint64) bind.IMethod
	SetChainRateLimiterConfigs(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) bind.IMethod
	SetChainRateLimiterConfig(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) bind.IMethod
	ProvideLiquidity(typeArgs string, state bind.Object, c_ bind.Object) bind.IMethod
	WithdrawLiquidity(typeArgs string, state bind.Object, amount uint64) bind.IMethod
	SetRebalancer(ownerCap bind.Object, state bind.Object, rebalancer string) bind.IMethod
	GetRebalancer(state bind.Object) bind.IMethod
	GetBalance(typeArgs string, state bind.Object) bind.IMethod
	DestroyTokenPool(typeArgs string, state bind.Object, ownerCap bind.Object) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type LockReleaseTokenPoolContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ ILockReleaseTokenPool = (*LockReleaseTokenPoolContract)(nil)

func NewLockReleaseTokenPool(packageID string, client suiclient.ClientImpl) (*LockReleaseTokenPoolContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &LockReleaseTokenPoolContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *LockReleaseTokenPoolContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OwnerCap struct {
	Id      string      `move:"sui::object::UID"`
	StateId bind.Object `move:"ID"`
}

type LockReleaseTokenPoolState struct {
	Id         string `move:"sui::object::UID"`
	Rebalancer string `move:"address"`
}

type TypeProof struct {
}

// Functions

func (c *LockReleaseTokenPoolContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "type_and_version", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) Initialize(typeArgs string, ref module_common.CCIPObjectRef, coinMetadata bind.Object, treasuryCap bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "initialize", false, "", typeArgs, ref, coinMetadata, treasuryCap, tokenPoolPackageId, tokenPoolAdministrator, rebalancer)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) InitializeByCcipAdmin(typeArgs string, ref module_common.CCIPObjectRef, coinMetadata bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string, rebalancer string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "initialize_by_ccip_admin", false, "", typeArgs, ref, coinMetadata, tokenPoolPackageId, tokenPoolAdministrator, rebalancer)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "initialize_by_ccip_admin", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetToken(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_token", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_token", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetTokenDecimals(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_token_decimals", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_token_decimals", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetRemotePools(state bind.Object, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_remote_pools", false, "", "", state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_remote_pools", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) IsRemotePool(state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "is_remote_pool", false, "", "", state, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "is_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetRemoteToken(state bind.Object, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_remote_token", false, "", "", state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_remote_token", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) AddRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "add_remote_pool", false, "", "", state, ownerCap, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "add_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) RemoveRemotePool(state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "remove_remote_pool", false, "", "", state, ownerCap, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "remove_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) IsSupportedChain(state bind.Object, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "is_supported_chain", false, "", "", state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "is_supported_chain", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetSupportedChains(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_supported_chains", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_supported_chains", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) ApplyChainUpdates(state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "apply_chain_updates", false, "", "", state, ownerCap, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "apply_chain_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetAllowlistEnabled(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_allowlist_enabled", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_allowlist_enabled", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetAllowlist(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_allowlist", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_allowlist", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) SetAllowlistEnabled(state bind.Object, ownerCap bind.Object, enabled bool) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "set_allowlist_enabled", false, "", "", state, ownerCap, enabled)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "set_allowlist_enabled", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) ApplyAllowlistUpdates(state bind.Object, ownerCap bind.Object, removes []string, adds []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "apply_allowlist_updates", false, "", "", state, ownerCap, removes, adds)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "apply_allowlist_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) LockOrBurn(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, state bind.Object, c_ bind.Object, tokenParams module_common.TokenParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "lock_or_burn", false, "", typeArgs, ref, clock, state, c_, tokenParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "lock_or_burn", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) ReleaseOrMint(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, pool bind.Object, receiverParams module_common.ReceiverParams, index uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "release_or_mint", false, "", typeArgs, ref, clock, pool, receiverParams, index)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "release_or_mint", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) SetChainRateLimiterConfigs(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "set_chain_rate_limiter_configs", false, "", "", state, ownerCap, clock, remoteChainSelectors, outboundIsEnableds, outboundCapacities, outboundRates, inboundIsEnableds, inboundCapacities, inboundRates)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "set_chain_rate_limiter_configs", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) SetChainRateLimiterConfig(state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "set_chain_rate_limiter_config", false, "", "", state, ownerCap, clock, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "set_chain_rate_limiter_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) ProvideLiquidity(typeArgs string, state bind.Object, c_ bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "provide_liquidity", false, "", typeArgs, state, c_)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "provide_liquidity", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) WithdrawLiquidity(typeArgs string, state bind.Object, amount uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "withdraw_liquidity", false, "", typeArgs, state, amount)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "withdraw_liquidity", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) SetRebalancer(ownerCap bind.Object, state bind.Object, rebalancer string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "set_rebalancer", false, "", "", ownerCap, state, rebalancer)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "set_rebalancer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetRebalancer(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_rebalancer", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_rebalancer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) GetBalance(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "get_balance", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "get_balance", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *LockReleaseTokenPoolContract) DestroyTokenPool(typeArgs string, state bind.Object, ownerCap bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "lock_release_token_pool", "destroy_token_pool", false, "", typeArgs, state, ownerCap)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "lock_release_token_pool", "destroy_token_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
