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
	SetOnRampInfos(param bind.Object, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) bind.IMethod
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

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
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

// Functions

func (c *RouterContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "type_and_version", false, "")
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "is_chain_supported", false, "", router, destChainSelector)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_info", false, "", router, destChainSelector)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_infos", false, "", router, destChainSelectors)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_version", false, "", info)
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
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "get_on_ramp_address", false, "", info)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "get_on_ramp_address", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *RouterContract) SetOnRampInfos(param bind.Object, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "router", "set_on_ramp_infos", false, "", param, router, destChainSelectors, onRampAddresses, onRampVersions)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "router", "set_on_ramp_infos", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
