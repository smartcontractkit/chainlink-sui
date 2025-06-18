// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_managed_token_pool

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

type IManagedTokenPool interface {
	TypeAndVersion() bind.IMethod
	Initialize(typeArgs string, ref module_common.CCIPObjectRef, treasuryCap bind.Object, coinMetadata bind.Object, mintCap bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string) bind.IMethod
	InitializeByCcipAdmin(typeArgs string, ref module_common.CCIPObjectRef, coinMetadata bind.Object, mintCap bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string) bind.IMethod
	AddRemotePool(typeArgs string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	RemoveRemotePool(typeArgs string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	IsSupportedChain(typeArgs string, state bind.Object, remoteChainSelector uint64) bind.IMethod
	GetSupportedChains(typeArgs string, state bind.Object) bind.IMethod
	ApplyChainUpdates(typeArgs string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) bind.IMethod
	GetAllowlistEnabled(typeArgs string, state bind.Object) bind.IMethod
	GetAllowlist(typeArgs string, state bind.Object) bind.IMethod
	ApplyAllowlistUpdates(typeArgs string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) bind.IMethod
	GetToken(typeArgs string, state bind.Object) bind.IMethod
	GetRouter() bind.IMethod
	GetTokenDecimals(typeArgs string, state bind.Object) bind.IMethod
	GetRemotePools(typeArgs string, state bind.Object, remoteChainSelector uint64) bind.IMethod
	IsRemotePool(typeArgs string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod
	GetRemoteToken(typeArgs string, state bind.Object, remoteChainSelector uint64) bind.IMethod
	LockOrBurn(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, state bind.Object, denyList bind.Object, tokenState bind.Object, c_ bind.Object, tokenParams module_common.TokenParams) bind.IMethod
	ReleaseOrMint(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, pool bind.Object, tokenState bind.Object, denyList bind.Object, receiverParams module_common.ReceiverParams, index uint64) bind.IMethod
	SetChainRateLimiterConfigs(typeArgs string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) bind.IMethod
	SetChainRateLimiterConfig(typeArgs string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type ManagedTokenPoolContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IManagedTokenPool = (*ManagedTokenPoolContract)(nil)

func NewManagedTokenPool(packageID string, client suiclient.ClientImpl) (*ManagedTokenPoolContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &ManagedTokenPoolContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *ManagedTokenPoolContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OwnerCap struct {
	Id      string      `move:"sui::object::UID"`
	StateId bind.Object `move:"ID"`
}

type ManagedTokenPoolState struct {
	Id      string      `move:"sui::object::UID"`
	MintCap bind.Object `move:"MintCap<T>"`
}

type TypeProof struct {
}

// Functions

func (c *ManagedTokenPoolContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "type_and_version", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) Initialize(typeArgs string, ref module_common.CCIPObjectRef, treasuryCap bind.Object, coinMetadata bind.Object, mintCap bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "initialize", false, "", typeArgs, ref, treasuryCap, coinMetadata, mintCap, tokenPoolPackageId, tokenPoolAdministrator)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) InitializeByCcipAdmin(typeArgs string, ref module_common.CCIPObjectRef, coinMetadata bind.Object, mintCap bind.Object, tokenPoolPackageId string, tokenPoolAdministrator string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "initialize_by_ccip_admin", false, "", typeArgs, ref, coinMetadata, mintCap, tokenPoolPackageId, tokenPoolAdministrator)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "initialize_by_ccip_admin", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) AddRemotePool(typeArgs string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "add_remote_pool", false, "", typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "add_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) RemoveRemotePool(typeArgs string, state bind.Object, ownerCap bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "remove_remote_pool", false, "", typeArgs, state, ownerCap, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "remove_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) IsSupportedChain(typeArgs string, state bind.Object, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "is_supported_chain", false, "", typeArgs, state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "is_supported_chain", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetSupportedChains(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_supported_chains", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_supported_chains", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) ApplyChainUpdates(typeArgs string, state bind.Object, ownerCap bind.Object, remoteChainSelectorsToRemove []uint64, remoteChainSelectorsToAdd []uint64, remotePoolAddressesToAdd [][][]byte, remoteTokenAddressesToAdd [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "apply_chain_updates", false, "", typeArgs, state, ownerCap, remoteChainSelectorsToRemove, remoteChainSelectorsToAdd, remotePoolAddressesToAdd, remoteTokenAddressesToAdd)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "apply_chain_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetAllowlistEnabled(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_allowlist_enabled", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_allowlist_enabled", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetAllowlist(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_allowlist", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_allowlist", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) ApplyAllowlistUpdates(typeArgs string, state bind.Object, ownerCap bind.Object, removes []string, adds []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "apply_allowlist_updates", false, "", typeArgs, state, ownerCap, removes, adds)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "apply_allowlist_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetToken(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_token", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_token", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetRouter() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_router", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_router", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetTokenDecimals(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_token_decimals", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_token_decimals", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetRemotePools(typeArgs string, state bind.Object, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_remote_pools", false, "", typeArgs, state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_remote_pools", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) IsRemotePool(typeArgs string, state bind.Object, remoteChainSelector uint64, remotePoolAddress []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "is_remote_pool", false, "", typeArgs, state, remoteChainSelector, remotePoolAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "is_remote_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) GetRemoteToken(typeArgs string, state bind.Object, remoteChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "get_remote_token", false, "", typeArgs, state, remoteChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "get_remote_token", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) LockOrBurn(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, state bind.Object, denyList bind.Object, tokenState bind.Object, c_ bind.Object, tokenParams module_common.TokenParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "lock_or_burn", false, "", typeArgs, ref, clock, state, denyList, tokenState, c_, tokenParams)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "lock_or_burn", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) ReleaseOrMint(typeArgs string, ref module_common.CCIPObjectRef, clock bind.Object, pool bind.Object, tokenState bind.Object, denyList bind.Object, receiverParams module_common.ReceiverParams, index uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "release_or_mint", false, "", typeArgs, ref, clock, pool, tokenState, denyList, receiverParams, index)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "release_or_mint", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) SetChainRateLimiterConfigs(typeArgs string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelectors []uint64, outboundIsEnableds []bool, outboundCapacities []uint64, outboundRates []uint64, inboundIsEnableds []bool, inboundCapacities []uint64, inboundRates []uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "set_chain_rate_limiter_configs", false, "", typeArgs, state, ownerCap, clock, remoteChainSelectors, outboundIsEnableds, outboundCapacities, outboundRates, inboundIsEnableds, inboundCapacities, inboundRates)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "set_chain_rate_limiter_configs", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenPoolContract) SetChainRateLimiterConfig(typeArgs string, state bind.Object, ownerCap bind.Object, clock bind.Object, remoteChainSelector uint64, outboundIsEnabled bool, outboundCapacity uint64, outboundRate uint64, inboundIsEnabled bool, inboundCapacity uint64, inboundRate uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token_pool", "set_chain_rate_limiter_config", false, "", typeArgs, state, ownerCap, clock, remoteChainSelector, outboundIsEnabled, outboundCapacity, outboundRate, inboundIsEnabled, inboundCapacity, inboundRate)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token_pool", "set_chain_rate_limiter_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
