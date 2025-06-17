// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_managed_token

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

type IManagedToken interface {
	TypeAndVersion() bind.IMethod
	Initialize(typeArgs string, treasuryCap bind.Object) bind.IMethod
	InitializeWithDenyCap(typeArgs string, treasuryCap bind.Object, denyCap bind.Object) bind.IMethod
	MintAllowance(typeArgs string, state bind.Object, mintCap bind.Object) bind.IMethod
	TotalSupply(typeArgs string, state bind.Object) bind.IMethod
	IsAuthorizedMintCap(typeArgs string, state bind.Object, id bind.Object) bind.IMethod
	ConfigureNewMinter(typeArgs string, state bind.Object, param bind.Object, minter string, allowance uint64, isUnlimited bool) bind.IMethod
	IncrementMintAllowance(typeArgs string, state bind.Object, param bind.Object, mintCapId bind.Object, denyList bind.Object, allowanceIncrement uint64) bind.IMethod
	SetUnlimitedMintAllowances(typeArgs string, state bind.Object, param bind.Object, mintCapId bind.Object, denyList bind.Object) bind.IMethod
	GetAllMintCaps(typeArgs string, state bind.Object) bind.IMethod
	MintAndTransfer(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) bind.IMethod
	Mint(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) bind.IMethod
	Burn(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, coin bind.Object, from string) bind.IMethod
	Blocklist(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) bind.IMethod
	Unblocklist(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) bind.IMethod
	Pause(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object) bind.IMethod
	Unpause(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object) bind.IMethod
	DestroyManagedToken(typeArgs string, ownerCap bind.Object, state bind.Object) bind.IMethod
	BorrowTreasuryCap(typeArgs string, ownerCap bind.Object, state bind.Object) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type ManagedTokenContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IManagedToken = (*ManagedTokenContract)(nil)

func NewManagedToken(packageID string, client suiclient.ClientImpl) (*ManagedTokenContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &ManagedTokenContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *ManagedTokenContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type TokenState struct {
	Id                string      `move:"sui::object::UID"`
	TreasuryCap       bind.Object `move:"TreasuryCap<T>"`
	DenyCap           bind.Object `move:"Option<DenyCapV2<T>>"`
	MintAllowancesMap bind.Object `move:"VecMap<ID, MintAllowance<T>>"`
}

type OwnerCap struct {
	Id      string      `move:"sui::object::UID"`
	StateId bind.Object `move:"ID"`
}

type MintCap struct {
	Id string `move:"sui::object::UID"`
}

type MintCapCreated struct {
	MintCap bind.Object `move:"ID"`
}

type MinterConfigured struct {
	MintCapOwner string      `move:"address"`
	MintCap      bind.Object `move:"ID"`
	Allowance    uint64      `move:"u64"`
	IsUnlimited  bool        `move:"bool"`
}

type Minted struct {
	MintCap bind.Object `move:"ID"`
	Minter  string      `move:"address"`
	To      string      `move:"address"`
	Amount  uint64      `move:"u64"`
}

type Burnt struct {
	MintCap bind.Object `move:"ID"`
	Burner  string      `move:"address"`
	From    string      `move:"address"`
	Amount  uint64      `move:"u64"`
}

type Blocklisted struct {
	Address string `move:"address"`
}

type Unblocklisted struct {
	Address string `move:"address"`
}

type Paused struct {
}

type Unpaused struct {
}

type MinterAllowanceIncremented struct {
	MintCap            bind.Object `move:"ID"`
	AllowanceIncrement uint64      `move:"u64"`
	NewAllowance       uint64      `move:"u64"`
}

type MinterUnlimitedAllowanceSet struct {
	MintCap bind.Object `move:"ID"`
}

// Functions

func (c *ManagedTokenContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "type_and_version", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) Initialize(typeArgs string, treasuryCap bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "initialize", false, "", typeArgs, treasuryCap)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) InitializeWithDenyCap(typeArgs string, treasuryCap bind.Object, denyCap bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "initialize_with_deny_cap", false, "", typeArgs, treasuryCap, denyCap)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "initialize_with_deny_cap", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) MintAllowance(typeArgs string, state bind.Object, mintCap bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "mint_allowance", false, "", typeArgs, state, mintCap)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "mint_allowance", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) TotalSupply(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "total_supply", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "total_supply", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) IsAuthorizedMintCap(typeArgs string, state bind.Object, id bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "is_authorized_mint_cap", false, "", typeArgs, state, id)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "is_authorized_mint_cap", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) ConfigureNewMinter(typeArgs string, state bind.Object, param bind.Object, minter string, allowance uint64, isUnlimited bool) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "configure_new_minter", false, "", typeArgs, state, param, minter, allowance, isUnlimited)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "configure_new_minter", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) IncrementMintAllowance(typeArgs string, state bind.Object, param bind.Object, mintCapId bind.Object, denyList bind.Object, allowanceIncrement uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "increment_mint_allowance", false, "", typeArgs, state, param, mintCapId, denyList, allowanceIncrement)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "increment_mint_allowance", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) SetUnlimitedMintAllowances(typeArgs string, state bind.Object, param bind.Object, mintCapId bind.Object, denyList bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "set_unlimited_mint_allowances", false, "", typeArgs, state, param, mintCapId, denyList)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "set_unlimited_mint_allowances", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) GetAllMintCaps(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "get_all_mint_caps", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "get_all_mint_caps", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) MintAndTransfer(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "mint_and_transfer", false, "", typeArgs, state, mintCap, denyList, amount, recipient)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "mint_and_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) Mint(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "mint", false, "", typeArgs, state, mintCap, denyList, amount, recipient)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "mint", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) Burn(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, coin bind.Object, from string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "burn", false, "", typeArgs, state, mintCap, denyList, coin, from)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "burn", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) Blocklist(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "blocklist", false, "", typeArgs, state, ownerCap, denyList, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "blocklist", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) Unblocklist(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object, addr string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "unblocklist", false, "", typeArgs, state, ownerCap, denyList, addr)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "unblocklist", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) Pause(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "pause", false, "", typeArgs, state, ownerCap, denyList)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "pause", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) Unpause(typeArgs string, state bind.Object, ownerCap bind.Object, denyList bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "unpause", false, "", typeArgs, state, ownerCap, denyList)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "unpause", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) DestroyManagedToken(typeArgs string, ownerCap bind.Object, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "destroy_managed_token", false, "", typeArgs, ownerCap, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "destroy_managed_token", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) BorrowTreasuryCap(typeArgs string, ownerCap bind.Object, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "borrow_treasury_cap", false, "", typeArgs, ownerCap, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "borrow_treasury_cap", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
