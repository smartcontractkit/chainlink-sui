// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_OwnershipTransferRequested

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

type IOwnershipTransferRequested interface {
	Contains(ref bind.Object) bind.IMethod
	Remove(ref bind.Object) bind.IMethod
	Borrow(ref bind.Object) bind.IMethod
	BorrowMut(ref bind.Object) bind.IMethod
	TransferOwnership(ref bind.Object, to string) bind.IMethod
	AcceptOwnership(ref bind.Object) bind.IMethod
	ExecuteOwnershipTransfer(ref bind.Object, to string) bind.IMethod
	GetCurrentOwner(ref bind.Object) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type OwnershipTransferRequestedContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IOwnershipTransferRequested = (*OwnershipTransferRequestedContract)(nil)

func NewOwnershipTransferRequested(packageID string, client suiclient.ClientImpl) (*OwnershipTransferRequestedContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &OwnershipTransferRequestedContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *OwnershipTransferRequestedContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OwnershipTransferAccepted struct {
	From string `move:"address"`
	To   string `move:"address"`
}

type OwnershipTransferred struct {
	From string `move:"address"`
	To   string `move:"address"`
}

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type CCIPObjectRef struct {
	Id              string           `move:"sui::object::UID"`
	CurrentOwner    string           `move:"address"`
	PendingTransfer *PendingTransfer `move:"0x1::option::Option<PendingTransfer>"`
}

type CCIPObjectRefPointer struct {
	Id          string `move:"sui::object::UID"`
	ObjectRefId string `move:"address"`
	OwnerCapId  string `move:"address"`
}

type PendingTransfer struct {
	From     string `move:"address"`
	To       string `move:"address"`
	Accepted bool   `move:"bool"`
}

type STATE_OBJECT struct {
}

// Functions

func (c *OwnershipTransferRequestedContract) Contains(ref bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "contains", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "contains", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnershipTransferRequestedContract) Remove(ref bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "remove", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "remove", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnershipTransferRequestedContract) Borrow(ref bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "borrow", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "borrow", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnershipTransferRequestedContract) BorrowMut(ref bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "borrow_mut", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "borrow_mut", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnershipTransferRequestedContract) TransferOwnership(ref bind.Object, to string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "transfer_ownership", false, "", ref, to)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "transfer_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnershipTransferRequestedContract) AcceptOwnership(ref bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "accept_ownership", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "accept_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnershipTransferRequestedContract) ExecuteOwnershipTransfer(ref bind.Object, to string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "execute_ownership_transfer", false, "", ref, to)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "execute_ownership_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *OwnershipTransferRequestedContract) GetCurrentOwner(ref bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "OwnershipTransferRequested", "get_current_owner", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "OwnershipTransferRequested", "get_current_owner", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
