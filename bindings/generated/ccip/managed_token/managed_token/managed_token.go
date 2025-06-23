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
	GetAllMintCaps(typeArgs string, state bind.Object) bind.IMethod
	MintAndTransfer(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) bind.IMethod
	Mint(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, amount uint64, recipient string) bind.IMethod
	Burn(typeArgs string, state bind.Object, mintCap bind.Object, denyList bind.Object, coin bind.Object, from string) bind.IMethod
	Owner(typeArgs string, state bind.Object) bind.IMethod
	HasPendingTransfer(typeArgs string, state bind.Object) bind.IMethod
	PendingTransferFrom(typeArgs string, state bind.Object) bind.IMethod
	PendingTransferTo(typeArgs string, state bind.Object) bind.IMethod
	PendingTransferAccepted(typeArgs string, state bind.Object) bind.IMethod
	AcceptOwnership(typeArgs string, state bind.Object) bind.IMethod
	AcceptOwnershipFromObject(typeArgs string, state bind.Object, from string) bind.IMethod
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

type McmsCallback struct {
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

func (c *ManagedTokenContract) Owner(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "owner", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "owner", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) HasPendingTransfer(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "has_pending_transfer", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "has_pending_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) PendingTransferFrom(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "pending_transfer_from", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "pending_transfer_from", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) PendingTransferTo(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "pending_transfer_to", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "pending_transfer_to", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) PendingTransferAccepted(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "pending_transfer_accepted", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "pending_transfer_accepted", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) AcceptOwnership(typeArgs string, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "accept_ownership", false, "", typeArgs, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "accept_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ManagedTokenContract) AcceptOwnershipFromObject(typeArgs string, state bind.Object, from string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "managed_token", "accept_ownership_from_object", false, "", typeArgs, state, from)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "managed_token", "accept_ownership_from_object", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
