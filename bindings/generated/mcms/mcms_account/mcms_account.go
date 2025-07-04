// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_account

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

type IMcmsAccount interface {
	TransferOwnership(param bind.Object, state bind.Object, to string) bind.IMethod
	TransferOwnershipToSelf(ownerCap bind.Object, state bind.Object) bind.IMethod
	AcceptOwnership(state bind.Object) bind.IMethod
	AcceptOwnershipAsTimelock(state bind.Object) bind.IMethod
	AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod
	PendingTransferFrom(state bind.Object) bind.IMethod
	PendingTransferTo(state bind.Object) bind.IMethod
	PendingTransferAccepted(state bind.Object) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type McmsAccountContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IMcmsAccount = (*McmsAccountContract)(nil)

func NewMcmsAccount(packageID string, client suiclient.ClientImpl) (*McmsAccountContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &McmsAccountContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *McmsAccountContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type AccountState struct {
	Id              string           `move:"sui::object::UID"`
	Owner           string           `move:"address"`
	PendingTransfer *PendingTransfer `move:"0x1::option::Option<PendingTransfer>"`
}

type PendingTransfer struct {
	From     string `move:"address"`
	To       string `move:"address"`
	Accepted bool   `move:"bool"`
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

type MCMS_ACCOUNT struct {
}

// Functions

func (c *McmsAccountContract) TransferOwnership(param bind.Object, state bind.Object, to string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "transfer_ownership", false, "", "", param, state, to)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "transfer_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsAccountContract) TransferOwnershipToSelf(ownerCap bind.Object, state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "transfer_ownership_to_self", false, "", "", ownerCap, state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "transfer_ownership_to_self", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsAccountContract) AcceptOwnership(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "accept_ownership", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "accept_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsAccountContract) AcceptOwnershipAsTimelock(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "accept_ownership_as_timelock", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "accept_ownership_as_timelock", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsAccountContract) AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "accept_ownership_from_object", false, "", "", state, from)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "accept_ownership_from_object", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsAccountContract) PendingTransferFrom(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "pending_transfer_from", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "pending_transfer_from", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsAccountContract) PendingTransferTo(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "pending_transfer_to", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "pending_transfer_to", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsAccountContract) PendingTransferAccepted(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_account", "pending_transfer_accepted", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_account", "pending_transfer_accepted", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
