// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_ownable

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

type IOwnable interface {
	New() bind.IMethod
	Owner(state bind.Object) bind.IMethod
	HasPendingTransfer(state bind.Object) bind.IMethod
	PendingTransferFrom(state bind.Object) bind.IMethod
	PendingTransferTo(state bind.Object) bind.IMethod
	PendingTransferAccepted(state bind.Object) bind.IMethod
	SetOwner(ownerCap bind.Object, state bind.Object, owner string) bind.IMethod
	TransferOwnership(ownerCap bind.Object, state bind.Object, to string) bind.IMethod
	AcceptOwnership(state bind.Object) bind.IMethod
	AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod
	AcceptOwnershipAsMcms(state bind.Object, mcms string) bind.IMethod
	ExecuteOwnershipTransfer(ownerCap bind.Object, state bind.Object, to string) bind.IMethod
	DestroyOwnableState(state bind.Object) bind.IMethod
	DestroyOwnerCap(ownerCap bind.Object) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type OwnableContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IOwnable = (*OwnableContract)(nil)

func NewOwnable(packageID string, client suiclient.ClientImpl) (*OwnableContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &OwnableContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *OwnableContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type OwnableState struct {
	Id              string           `move:"sui::object::UID"`
	Owner           string           `move:"address"`
	PendingTransfer *PendingTransfer `move:"0x1::option::Option<PendingTransfer>"`
	OwnerCapId      bind.Object      `move:"ID"`
}

type PendingTransfer struct {
	From     string `move:"address"`
	To       string `move:"address"`
	Accepted bool   `move:"bool"`
}

type McmsCallback struct {
}

type NewOwnableStateEvent struct {
	OwnableStateId bind.Object `move:"ID"`
	OwnerCapId     bind.Object `move:"ID"`
	Owner          string      `move:"address"`
}

type OwnershipTransferRequested struct {
	From string `move:"address"`
	To   string `move:"address"`
}

type OwnershipTransferAccepted struct {
	From string `move:"address"`
	To   string `move:"address"`
}

type OwnershipTransferred struct {
	From string `move:"address"`
	To   string `move:"address"`
}

// Functions

func (c *OwnableContract) New() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "new", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "new", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) Owner(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "owner", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "owner", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) HasPendingTransfer(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "has_pending_transfer", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "has_pending_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) PendingTransferFrom(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "pending_transfer_from", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "pending_transfer_from", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) PendingTransferTo(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "pending_transfer_to", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "pending_transfer_to", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) PendingTransferAccepted(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "pending_transfer_accepted", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "pending_transfer_accepted", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) SetOwner(ownerCap bind.Object, state bind.Object, owner string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "set_owner", false, "", "", ownerCap, state, owner)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "set_owner", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) TransferOwnership(ownerCap bind.Object, state bind.Object, to string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "transfer_ownership", false, "", "", ownerCap, state, to)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "transfer_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) AcceptOwnership(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "accept_ownership", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "accept_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "accept_ownership_from_object", false, "", "", state, from)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "accept_ownership_from_object", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) AcceptOwnershipAsMcms(state bind.Object, mcms string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "accept_ownership_as_mcms", false, "", "", state, mcms)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "accept_ownership_as_mcms", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) ExecuteOwnershipTransfer(ownerCap bind.Object, state bind.Object, to string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "execute_ownership_transfer", false, "", "", ownerCap, state, to)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "execute_ownership_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) DestroyOwnableState(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "destroy_ownable_state", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "destroy_ownable_state", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnableContract) DestroyOwnerCap(ownerCap bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "ownable", "destroy_owner_cap", false, "", "", ownerCap)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "ownable", "destroy_owner_cap", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
