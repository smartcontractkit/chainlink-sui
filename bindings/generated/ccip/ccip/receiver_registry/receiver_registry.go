// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_receiver_registry

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

type IReceiverRegistry interface {
	TypeAndVersion() bind.IMethod
	Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap) bind.IMethod
	GetReceiverConfig(ref module_common.CCIPObjectRef, receiverPackageId string) bind.IMethod
	GetReceiverConfigFields(rc ReceiverConfig) bind.IMethod
	IsRegisteredReceiver(ref module_common.CCIPObjectRef, receiverPackageId string) bind.IMethod
	GetReceiverModuleAndState(ref module_common.CCIPObjectRef, receiverPackageId string) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type ReceiverRegistryContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IReceiverRegistry = (*ReceiverRegistryContract)(nil)

func NewReceiverRegistry(packageID string, client suiclient.ClientImpl) (*ReceiverRegistryContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &ReceiverRegistryContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *ReceiverRegistryContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type ReceiverConfig struct {
	ModuleName      string `move:"0x1::string::String"`
	FunctionName    string `move:"0x1::string::String"`
	ReceiverStateId string `move:"address"`
}

type ReceiverRegistry struct {
	Id string `move:"sui::object::UID"`
}

type ReceiverRegistered struct {
	ReceiverPackageId  string `move:"address"`
	ReceiverStateId    string `move:"address"`
	ReceiverModuleName string `move:"0x1::string::String"`
}

// Functions

func (c *ReceiverRegistryContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "receiver_registry", "type_and_version", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "receiver_registry", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ReceiverRegistryContract) Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "receiver_registry", "initialize", false, "", "", ref, param)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "receiver_registry", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ReceiverRegistryContract) GetReceiverConfig(ref module_common.CCIPObjectRef, receiverPackageId string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "receiver_registry", "get_receiver_config", false, "", "", ref, receiverPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "receiver_registry", "get_receiver_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ReceiverRegistryContract) GetReceiverConfigFields(rc ReceiverConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "receiver_registry", "get_receiver_config_fields", false, "", "", rc)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "receiver_registry", "get_receiver_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ReceiverRegistryContract) IsRegisteredReceiver(ref module_common.CCIPObjectRef, receiverPackageId string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "receiver_registry", "is_registered_receiver", false, "", "", ref, receiverPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "receiver_registry", "is_registered_receiver", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *ReceiverRegistryContract) GetReceiverModuleAndState(ref module_common.CCIPObjectRef, receiverPackageId string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "receiver_registry", "get_receiver_module_and_state", false, "", "", ref, receiverPackageId)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "receiver_registry", "get_receiver_module_and_state", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
