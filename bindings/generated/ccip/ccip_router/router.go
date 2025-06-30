// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_router

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

type IRouter interface {
	TypeAndVersion() bind.IMethod
	IsChainSupported(router bind.Object, destChainSelector uint64) bind.IMethod
	GetOnRampInfo(router bind.Object, destChainSelector uint64) bind.IMethod
	GetOnRampInfos(router bind.Object, destChainSelectors []uint64) bind.IMethod
	GetOnRampVersion(info OnRampInfo) bind.IMethod
	GetOnRampAddress(info OnRampInfo) bind.IMethod
	SetOnRampInfos(ownerCap module_common.OwnerCap, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) bind.IMethod
	Owner(state bind.Object) bind.IMethod
	HasPendingTransfer(state bind.Object) bind.IMethod
	PendingTransferFrom(state bind.Object) bind.IMethod
	PendingTransferTo(state bind.Object) bind.IMethod
	PendingTransferAccepted(state bind.Object) bind.IMethod
	TransferOwnership(state bind.Object, ownerCap module_common.OwnerCap, newOwner string) bind.IMethod
	AcceptOwnership(state bind.Object) bind.IMethod
	AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type RouterContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IRouter = (*RouterContract)(nil)

func NewRouter(packageID string, client suiclient.ClientImpl) (*RouterContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &RouterContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *RouterContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type ROUTER struct {
}

type OnRampSet struct {
	DestChainSelector uint64     `move:"u64"`
	OnRampInfo        OnRampInfo `move:"OnRampInfo"`
}

type OnRampInfo struct {
	OnrampAddress string `move:"address"`
	OnrampVersion []byte `move:"vector<u8>"`
}

type RouterState struct {
	Id string `move:"sui::object::UID"`
}

type McmsCallback struct {
}

// Functions

func (c *RouterContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "type_and_version", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) IsChainSupported(router bind.Object, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "is_chain_supported", false, "", "", router, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "is_chain_supported", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) GetOnRampInfo(router bind.Object, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_info", false, "", "", router, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "get_on_ramp_info", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) GetOnRampInfos(router bind.Object, destChainSelectors []uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_infos", false, "", "", router, destChainSelectors)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "get_on_ramp_infos", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) GetOnRampVersion(info OnRampInfo) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_version", false, "", "", info)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "get_on_ramp_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) GetOnRampAddress(info OnRampInfo) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_address", false, "", "", info)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "get_on_ramp_address", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) SetOnRampInfos(ownerCap module_common.OwnerCap, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "set_on_ramp_infos", false, "", "", ownerCap, router, destChainSelectors, onRampAddresses, onRampVersions)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "set_on_ramp_infos", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) Owner(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "owner", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "owner", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) HasPendingTransfer(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "has_pending_transfer", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "has_pending_transfer", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) PendingTransferFrom(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "pending_transfer_from", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "pending_transfer_from", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) PendingTransferTo(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "pending_transfer_to", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "pending_transfer_to", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) PendingTransferAccepted(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "pending_transfer_accepted", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "pending_transfer_accepted", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) TransferOwnership(state bind.Object, ownerCap module_common.OwnerCap, newOwner string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "transfer_ownership", false, "", "", state, ownerCap, newOwner)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "transfer_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) AcceptOwnership(state bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "accept_ownership", false, "", "", state)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "accept_ownership", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) AcceptOwnershipFromObject(state bind.Object, from string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "accept_ownership_from_object", false, "", "", state, from)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "accept_ownership_from_object", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
