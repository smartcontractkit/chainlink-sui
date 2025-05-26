// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_token_admin_registry

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_common "github.com/smartcontractkit/chainlink-sui/bindings/common"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
)

type ITokenAdminRegistry interface {
	TypeAndVersion() bind.IMethod
	Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap) bind.IMethod
	GetPools(ref module_common.CCIPObjectRef, coinMetadataAddresses []string) bind.IMethod
	GetPool(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod
	GetTokenConfig(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod
	GetAllConfiguredTokens(ref module_common.CCIPObjectRef, startKey string, maxCount uint64) bind.IMethod
	UnregisterPool(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod
	TransferAdminRole(ref module_common.CCIPObjectRef, coinMetadataAddress string, newAdmin string) bind.IMethod
	AcceptAdminRole(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod
	IsAdministrator(ref module_common.CCIPObjectRef, coinMetadataAddress string, administrator string) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type TokenAdminRegistryContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ ITokenAdminRegistry = (*TokenAdminRegistryContract)(nil)

func NewTokenAdminRegistry(packageID string, client suiclient.ClientImpl) (*TokenAdminRegistryContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &TokenAdminRegistryContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *TokenAdminRegistryContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type TokenAdminRegistryState struct {
	Id string `move:"sui::object::UID"`
}

type TokenConfig struct {
	TokenPoolAddress     string `move:"address"`
	Administrator        string `move:"address"`
	PendingAdministrator string `move:"address"`
}

type PoolSet struct {
	CoinMetadataAddress string `move:"address"`
	PreviousPoolAddress string `move:"address"`
	NewPoolAddress      string `move:"address"`
}

type TokenUnregistered struct {
	LocalToken          string `move:"address"`
	PreviousPoolAddress string `move:"address"`
}

type AdministratorTransferRequested struct {
	CoinMetadataAddress string `move:"address"`
	CurrentAdmin        string `move:"address"`
	NewAdmin            string `move:"address"`
}

type AdministratorTransferred struct {
	CoinMetadataAddress string `move:"address"`
	NewAdmin            string `move:"address"`
}

// Functions

func (c *TokenAdminRegistryContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "type_and_version", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "initialize", false, "", ref, param)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) GetPools(ref module_common.CCIPObjectRef, coinMetadataAddresses []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "get_pools", false, "", ref, coinMetadataAddresses)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "get_pools", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) GetPool(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "get_pool", false, "", ref, coinMetadataAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "get_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) GetTokenConfig(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "get_token_config", false, "", ref, coinMetadataAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "get_token_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) GetAllConfiguredTokens(ref module_common.CCIPObjectRef, startKey string, maxCount uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "get_all_configured_tokens", false, "", ref, startKey, maxCount)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "get_all_configured_tokens", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) UnregisterPool(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "unregister_pool", false, "", ref, coinMetadataAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "unregister_pool", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) TransferAdminRole(ref module_common.CCIPObjectRef, coinMetadataAddress string, newAdmin string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "transfer_admin_role", false, "", ref, coinMetadataAddress, newAdmin)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "transfer_admin_role", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) AcceptAdminRole(ref module_common.CCIPObjectRef, coinMetadataAddress string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "accept_admin_role", false, "", ref, coinMetadataAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "accept_admin_role", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *TokenAdminRegistryContract) IsAdministrator(ref module_common.CCIPObjectRef, coinMetadataAddress string, administrator string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "token_admin_registry", "is_administrator", false, "", ref, coinMetadataAddress, administrator)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "token_admin_registry", "is_administrator", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
