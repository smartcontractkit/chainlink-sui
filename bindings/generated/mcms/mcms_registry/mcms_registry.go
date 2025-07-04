// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_registry

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

type IMcmsRegistry interface {
	BorrowOwnerCap(typeArgs string, registry bind.Object) bind.IMethod
	GetCallbackParamsForMcms(params ExecutingCallbackParams) bind.IMethod
	CreateExecutingCallbackParams(target string, moduleName string, functionName string, data []byte) bind.IMethod
	IsPackageRegistered(registry bind.Object, packageAddress string) bind.IMethod
	IsModuleRegistered(registry bind.Object, moduleName string) bind.IMethod
	Target(params ExecutingCallbackParams) bind.IMethod
	ModuleName(params ExecutingCallbackParams) bind.IMethod
	FunctionName(params ExecutingCallbackParams) bind.IMethod
	Data(params ExecutingCallbackParams) bind.IMethod
	CreateMcmsProof() bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type McmsRegistryContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IMcmsRegistry = (*McmsRegistryContract)(nil)

func NewMcmsRegistry(packageID string, client suiclient.ClientImpl) (*McmsRegistryContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &McmsRegistryContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *McmsRegistryContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type Registry struct {
	Id string `move:"sui::object::UID"`
}

type ExecutingCallbackParams struct {
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Data         []byte `move:"vector<u8>"`
}

type RegisteredModule struct {
}

type EntrypointRegistered struct {
	RegistryId     bind.Object `move:"ID"`
	AccountAddress string      `move:"address"`
	ModuleName     string      `move:"0x1::string::String"`
}

type MCMS_REGISTRY struct {
}

type McmsProof struct {
}

// Functions

func (c *McmsRegistryContract) BorrowOwnerCap(typeArgs string, registry bind.Object) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "borrow_owner_cap", false, "", typeArgs, registry)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "borrow_owner_cap", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) GetCallbackParamsForMcms(params ExecutingCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "get_callback_params_for_mcms", false, "", "", params)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "get_callback_params_for_mcms", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) CreateExecutingCallbackParams(target string, moduleName string, functionName string, data []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "create_executing_callback_params", false, "", "", target, moduleName, functionName, data)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "create_executing_callback_params", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) IsPackageRegistered(registry bind.Object, packageAddress string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "is_package_registered", false, "", "", registry, packageAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "is_package_registered", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) IsModuleRegistered(registry bind.Object, moduleName string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "is_module_registered", false, "", "", registry, moduleName)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "is_module_registered", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) Target(params ExecutingCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "target", false, "", "", params)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "target", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) ModuleName(params ExecutingCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "module_name", false, "", "", params)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "module_name", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) FunctionName(params ExecutingCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "function_name", false, "", "", params)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "function_name", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) Data(params ExecutingCallbackParams) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "data", false, "", "", params)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "data", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *McmsRegistryContract) CreateMcmsProof() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "mcms_registry", "create_mcms_proof", false, "", "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "mcms_registry", "create_mcms_proof", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
