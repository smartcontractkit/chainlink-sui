// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_registry

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/block-vision/sui-go-sdk/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

var (
	_ = big.NewInt
)

type IMcmsRegistry interface {
	RegisterEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, packageCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCallbackParams(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, witness bind.Object) (*models.SuiTransactionBlockResponse, error)
	BorrowOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCallbackParamsForMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, params ExecutingCallbackParams, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte) (*models.SuiTransactionBlockResponse, error)
	IsPackageRegistered(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	Target(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	ModuleName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	FunctionName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	Data(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	GetMultisigAddress(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	CreateMcmsProof(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMcmsRegistryDevInspect
	Encoder() McmsRegistryEncoder
}

type IMcmsRegistryDevInspect interface {
	GetCallbackParams(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) ([]any, error)
	ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, witness bind.Object) (any, error)
	BorrowOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object) (bind.Object, error)
	GetCallbackParamsForMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, params ExecutingCallbackParams, proof bind.Object) ([]any, error)
	GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) ([]any, error)
	CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte) (ExecutingCallbackParams, error)
	IsPackageRegistered(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (bool, error)
	Target(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error)
	ModuleName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error)
	FunctionName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error)
	Data(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) ([]byte, error)
	GetMultisigAddress(ctx context.Context, opts *bind.CallOpts) (string, error)
	CreateMcmsProof(ctx context.Context, opts *bind.CallOpts) (McmsProof, error)
}

type McmsRegistryEncoder interface {
	RegisterEntrypoint(typeArgs []string, registry bind.Object, proof bind.Object, packageCap bind.Object) (*bind.EncodedCall, error)
	RegisterEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetCallbackParams(typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*bind.EncodedCall, error)
	GetCallbackParamsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReleaseCap(typeArgs []string, registry bind.Object, witness bind.Object) (*bind.EncodedCall, error)
	ReleaseCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	BorrowOwnerCap(typeArgs []string, registry bind.Object) (*bind.EncodedCall, error)
	BorrowOwnerCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetCallbackParamsForMcms(typeArgs []string, params ExecutingCallbackParams, proof bind.Object) (*bind.EncodedCall, error)
	GetCallbackParamsForMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetCallbackParamsFromMcms(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	GetCallbackParamsFromMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateExecutingCallbackParams(target string, moduleName string, functionName string, data []byte) (*bind.EncodedCall, error)
	CreateExecutingCallbackParamsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsPackageRegistered(registry bind.Object, packageAddress string) (*bind.EncodedCall, error)
	IsPackageRegisteredWithArgs(args ...any) (*bind.EncodedCall, error)
	Target(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	TargetWithArgs(args ...any) (*bind.EncodedCall, error)
	ModuleName(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	ModuleNameWithArgs(args ...any) (*bind.EncodedCall, error)
	FunctionName(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	FunctionNameWithArgs(args ...any) (*bind.EncodedCall, error)
	Data(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	DataWithArgs(args ...any) (*bind.EncodedCall, error)
	GetMultisigAddress() (*bind.EncodedCall, error)
	GetMultisigAddressWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateMcmsProof() (*bind.EncodedCall, error)
	CreateMcmsProofWithArgs(args ...any) (*bind.EncodedCall, error)
}

type McmsRegistryContract struct {
	*bind.BoundContract
	mcmsRegistryEncoder
	devInspect *McmsRegistryDevInspect
}

type McmsRegistryDevInspect struct {
	contract *McmsRegistryContract
}

var _ IMcmsRegistry = (*McmsRegistryContract)(nil)
var _ IMcmsRegistryDevInspect = (*McmsRegistryDevInspect)(nil)

func NewMcmsRegistry(packageID string, client sui.ISuiAPI) (*McmsRegistryContract, error) {
	contract, err := bind.NewBoundContract(packageID, "mcms", "mcms_registry", client)
	if err != nil {
		return nil, err
	}

	c := &McmsRegistryContract{
		BoundContract:       contract,
		mcmsRegistryEncoder: mcmsRegistryEncoder{BoundContract: contract},
	}
	c.devInspect = &McmsRegistryDevInspect{contract: c}
	return c, nil
}

func (c *McmsRegistryContract) Encoder() McmsRegistryEncoder {
	return c.mcmsRegistryEncoder
}

func (c *McmsRegistryContract) DevInspect() IMcmsRegistryDevInspect {
	return c.devInspect
}

type Registry struct {
	Id          string      `move:"sui::object::UID"`
	PackageCaps bind.Object `move:"Bag"`
}

type ExecutingCallbackParams struct {
	Target       string `move:"address"`
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Data         []byte `move:"vector<u8>"`
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

type bcsExecutingCallbackParams struct {
	Target       [32]byte
	ModuleName   string
	FunctionName string
	Data         []byte
}

func convertExecutingCallbackParamsFromBCS(bcs bcsExecutingCallbackParams) (ExecutingCallbackParams, error) {

	return ExecutingCallbackParams{
		Target:       fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:   bcs.ModuleName,
		FunctionName: bcs.FunctionName,
		Data:         bcs.Data,
	}, nil
}

type bcsEntrypointRegistered struct {
	RegistryId     bind.Object
	AccountAddress [32]byte
	ModuleName     string
}

func convertEntrypointRegisteredFromBCS(bcs bcsEntrypointRegistered) (EntrypointRegistered, error) {

	return EntrypointRegistered{
		RegistryId:     bcs.RegistryId,
		AccountAddress: fmt.Sprintf("0x%x", bcs.AccountAddress),
		ModuleName:     bcs.ModuleName,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("mcms::mcms_registry::Registry", func(data []byte) (interface{}, error) {
		var result Registry
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::ExecutingCallbackParams", func(data []byte) (interface{}, error) {
		var temp bcsExecutingCallbackParams
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertExecutingCallbackParamsFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::EntrypointRegistered", func(data []byte) (interface{}, error) {
		var temp bcsEntrypointRegistered
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertEntrypointRegisteredFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::MCMS_REGISTRY", func(data []byte) (interface{}, error) {
		var result MCMS_REGISTRY
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::McmsProof", func(data []byte) (interface{}, error) {
		var result McmsProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// RegisterEntrypoint executes the register_entrypoint Move function.
func (c *McmsRegistryContract) RegisterEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, packageCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.RegisterEntrypoint(typeArgs, registry, proof, packageCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCallbackParams executes the get_callback_params Move function.
func (c *McmsRegistryContract) GetCallbackParams(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetCallbackParams(typeArgs, registry, proof, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseCap executes the release_cap Move function.
func (c *McmsRegistryContract) ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, witness bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.ReleaseCap(typeArgs, registry, witness)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// BorrowOwnerCap executes the borrow_owner_cap Move function.
func (c *McmsRegistryContract) BorrowOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.BorrowOwnerCap(typeArgs, registry)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCallbackParamsForMcms executes the get_callback_params_for_mcms Move function.
func (c *McmsRegistryContract) GetCallbackParamsForMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, params ExecutingCallbackParams, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetCallbackParamsForMcms(typeArgs, params, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCallbackParamsFromMcms executes the get_callback_params_from_mcms Move function.
func (c *McmsRegistryContract) GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetCallbackParamsFromMcms(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateExecutingCallbackParams executes the create_executing_callback_params Move function.
func (c *McmsRegistryContract) CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.CreateExecutingCallbackParams(target, moduleName, functionName, data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsPackageRegistered executes the is_package_registered Move function.
func (c *McmsRegistryContract) IsPackageRegistered(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.IsPackageRegistered(registry, packageAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Target executes the target Move function.
func (c *McmsRegistryContract) Target(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.Target(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ModuleName executes the module_name Move function.
func (c *McmsRegistryContract) ModuleName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.ModuleName(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// FunctionName executes the function_name Move function.
func (c *McmsRegistryContract) FunctionName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.FunctionName(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Data executes the data Move function.
func (c *McmsRegistryContract) Data(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.Data(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetMultisigAddress executes the get_multisig_address Move function.
func (c *McmsRegistryContract) GetMultisigAddress(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetMultisigAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateMcmsProof executes the create_mcms_proof Move function.
func (c *McmsRegistryContract) CreateMcmsProof(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.CreateMcmsProof()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCallbackParams executes the get_callback_params Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: &C
//	[1]: 0x1::string::String
//	[2]: vector<u8>
func (d *McmsRegistryDevInspect) GetCallbackParams(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) ([]any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetCallbackParams(typeArgs, registry, proof, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// ReleaseCap executes the release_cap Move function using DevInspect to get return values.
//
// Returns: C
func (d *McmsRegistryDevInspect) ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, witness bind.Object) (any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.ReleaseCap(typeArgs, registry, witness)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	return results[0], nil
}

// BorrowOwnerCap executes the borrow_owner_cap Move function using DevInspect to get return values.
//
// Returns: &C
func (d *McmsRegistryDevInspect) BorrowOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object) (bind.Object, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.BorrowOwnerCap(typeArgs, registry)
	if err != nil {
		return bind.Object{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return bind.Object{}, err
	}
	if len(results) == 0 {
		return bind.Object{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bind.Object)
	if !ok {
		return bind.Object{}, fmt.Errorf("unexpected return type: expected bind.Object, got %T", results[0])
	}
	return result, nil
}

// GetCallbackParamsForMcms executes the get_callback_params_for_mcms Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: 0x1::string::String
//	[2]: 0x1::string::String
//	[3]: vector<u8>
func (d *McmsRegistryDevInspect) GetCallbackParamsForMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, params ExecutingCallbackParams, proof bind.Object) ([]any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetCallbackParamsForMcms(typeArgs, params, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetCallbackParamsFromMcms executes the get_callback_params_from_mcms Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: 0x1::string::String
//	[2]: 0x1::string::String
//	[3]: vector<u8>
func (d *McmsRegistryDevInspect) GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) ([]any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetCallbackParamsFromMcms(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// CreateExecutingCallbackParams executes the create_executing_callback_params Move function using DevInspect to get return values.
//
// Returns: ExecutingCallbackParams
func (d *McmsRegistryDevInspect) CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte) (ExecutingCallbackParams, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.CreateExecutingCallbackParams(target, moduleName, functionName, data)
	if err != nil {
		return ExecutingCallbackParams{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return ExecutingCallbackParams{}, err
	}
	if len(results) == 0 {
		return ExecutingCallbackParams{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(ExecutingCallbackParams)
	if !ok {
		return ExecutingCallbackParams{}, fmt.Errorf("unexpected return type: expected ExecutingCallbackParams, got %T", results[0])
	}
	return result, nil
}

// IsPackageRegistered executes the is_package_registered Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsRegistryDevInspect) IsPackageRegistered(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (bool, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.IsPackageRegistered(registry, packageAddress)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// Target executes the target Move function using DevInspect to get return values.
//
// Returns: address
func (d *McmsRegistryDevInspect) Target(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.Target(params)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// ModuleName executes the module_name Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *McmsRegistryDevInspect) ModuleName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.ModuleName(params)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// FunctionName executes the function_name Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *McmsRegistryDevInspect) FunctionName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.FunctionName(params)
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// Data executes the data Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsRegistryDevInspect) Data(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) ([]byte, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.Data(params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no return value")
	}
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// GetMultisigAddress executes the get_multisig_address Move function using DevInspect to get return values.
//
// Returns: address
func (d *McmsRegistryDevInspect) GetMultisigAddress(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetMultisigAddress()
	if err != nil {
		return "", fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no return value")
	}
	result, ok := results[0].(string)
	if !ok {
		return "", fmt.Errorf("unexpected return type: expected string, got %T", results[0])
	}
	return result, nil
}

// CreateMcmsProof executes the create_mcms_proof Move function using DevInspect to get return values.
//
// Returns: McmsProof
func (d *McmsRegistryDevInspect) CreateMcmsProof(ctx context.Context, opts *bind.CallOpts) (McmsProof, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.CreateMcmsProof()
	if err != nil {
		return McmsProof{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return McmsProof{}, err
	}
	if len(results) == 0 {
		return McmsProof{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(McmsProof)
	if !ok {
		return McmsProof{}, fmt.Errorf("unexpected return type: expected McmsProof, got %T", results[0])
	}
	return result, nil
}

type mcmsRegistryEncoder struct {
	*bind.BoundContract
}

// RegisterEntrypoint encodes a call to the register_entrypoint Move function.
func (c mcmsRegistryEncoder) RegisterEntrypoint(typeArgs []string, registry bind.Object, proof bind.Object, packageCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("register_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"T",
		"C",
	}, []any{
		registry,
		proof,
		packageCap,
	}, nil)
}

// RegisterEntrypointWithArgs encodes a call to the register_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) RegisterEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"T",
		"C",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("register_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCallbackParams encodes a call to the get_callback_params Move function.
func (c mcmsRegistryEncoder) GetCallbackParams(typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("get_callback_params", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"T",
		"mcms::mcms_registry::ExecutingCallbackParams",
	}, []any{
		registry,
		proof,
		params,
	}, []string{
		"&C",
		"0x1::string::String",
		"vector<u8>",
	})
}

// GetCallbackParamsWithArgs encodes a call to the get_callback_params Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetCallbackParamsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"T",
		"mcms::mcms_registry::ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("get_callback_params", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&C",
		"0x1::string::String",
		"vector<u8>",
	})
}

// ReleaseCap encodes a call to the release_cap Move function.
func (c mcmsRegistryEncoder) ReleaseCap(typeArgs []string, registry bind.Object, witness bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("release_cap", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"T",
	}, []any{
		registry,
		witness,
	}, []string{
		"C",
	})
}

// ReleaseCapWithArgs encodes a call to the release_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) ReleaseCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"T",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("release_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"C",
	})
}

// BorrowOwnerCap encodes a call to the borrow_owner_cap Move function.
func (c mcmsRegistryEncoder) BorrowOwnerCap(typeArgs []string, registry bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"C",
	}
	return c.EncodeCallArgsWithGenerics("borrow_owner_cap", typeArgsList, typeParamsList, []string{
		"&Registry",
	}, []any{
		registry,
	}, []string{
		"&C",
	})
}

// BorrowOwnerCapWithArgs encodes a call to the borrow_owner_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) BorrowOwnerCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Registry",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"C",
	}
	return c.EncodeCallArgsWithGenerics("borrow_owner_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&C",
	})
}

// GetCallbackParamsForMcms encodes a call to the get_callback_params_for_mcms Move function.
func (c mcmsRegistryEncoder) GetCallbackParamsForMcms(typeArgs []string, params ExecutingCallbackParams, proof bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_callback_params_for_mcms", typeArgsList, typeParamsList, []string{
		"mcms::mcms_registry::ExecutingCallbackParams",
		"T",
	}, []any{
		params,
		proof,
	}, []string{
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
	})
}

// GetCallbackParamsForMcmsWithArgs encodes a call to the get_callback_params_for_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetCallbackParamsForMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms_registry::ExecutingCallbackParams",
		"T",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_callback_params_for_mcms", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
	})
}

// GetCallbackParamsFromMcms encodes a call to the get_callback_params_from_mcms Move function.
func (c mcmsRegistryEncoder) GetCallbackParamsFromMcms(params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_callback_params_from_mcms", typeArgsList, typeParamsList, []string{
		"mcms::mcms_registry::ExecutingCallbackParams",
	}, []any{
		params,
	}, []string{
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
	})
}

// GetCallbackParamsFromMcmsWithArgs encodes a call to the get_callback_params_from_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetCallbackParamsFromMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms_registry::ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_callback_params_from_mcms", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
	})
}

// CreateExecutingCallbackParams encodes a call to the create_executing_callback_params Move function.
func (c mcmsRegistryEncoder) CreateExecutingCallbackParams(target string, moduleName string, functionName string, data []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_executing_callback_params", typeArgsList, typeParamsList, []string{
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
	}, []any{
		target,
		moduleName,
		functionName,
		data,
	}, []string{
		"mcms::mcms_registry::ExecutingCallbackParams",
	})
}

// CreateExecutingCallbackParamsWithArgs encodes a call to the create_executing_callback_params Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) CreateExecutingCallbackParamsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_executing_callback_params", typeArgsList, typeParamsList, expectedParams, args, []string{
		"mcms::mcms_registry::ExecutingCallbackParams",
	})
}

// IsPackageRegistered encodes a call to the is_package_registered Move function.
func (c mcmsRegistryEncoder) IsPackageRegistered(registry bind.Object, packageAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_package_registered", typeArgsList, typeParamsList, []string{
		"&Registry",
		"address",
	}, []any{
		registry,
		packageAddress,
	}, []string{
		"bool",
	})
}

// IsPackageRegisteredWithArgs encodes a call to the is_package_registered Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) IsPackageRegisteredWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Registry",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_package_registered", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// Target encodes a call to the target Move function.
func (c mcmsRegistryEncoder) Target(params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("target", typeArgsList, typeParamsList, []string{
		"&ExecutingCallbackParams",
	}, []any{
		params,
	}, []string{
		"address",
	})
}

// TargetWithArgs encodes a call to the target Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) TargetWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("target", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// ModuleName encodes a call to the module_name Move function.
func (c mcmsRegistryEncoder) ModuleName(params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("module_name", typeArgsList, typeParamsList, []string{
		"&ExecutingCallbackParams",
	}, []any{
		params,
	}, []string{
		"0x1::string::String",
	})
}

// ModuleNameWithArgs encodes a call to the module_name Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) ModuleNameWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("module_name", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// FunctionName encodes a call to the function_name Move function.
func (c mcmsRegistryEncoder) FunctionName(params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_name", typeArgsList, typeParamsList, []string{
		"&ExecutingCallbackParams",
	}, []any{
		params,
	}, []string{
		"0x1::string::String",
	})
}

// FunctionNameWithArgs encodes a call to the function_name Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) FunctionNameWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("function_name", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// Data encodes a call to the data Move function.
func (c mcmsRegistryEncoder) Data(params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("data", typeArgsList, typeParamsList, []string{
		"&ExecutingCallbackParams",
	}, []any{
		params,
	}, []string{
		"vector<u8>",
	})
}

// DataWithArgs encodes a call to the data Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) DataWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("data", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// GetMultisigAddress encodes a call to the get_multisig_address Move function.
func (c mcmsRegistryEncoder) GetMultisigAddress() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_multisig_address", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"address",
	})
}

// GetMultisigAddressWithArgs encodes a call to the get_multisig_address Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetMultisigAddressWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_multisig_address", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// CreateMcmsProof encodes a call to the create_mcms_proof Move function.
func (c mcmsRegistryEncoder) CreateMcmsProof() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_mcms_proof", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"mcms::mcms_registry::McmsProof",
	})
}

// CreateMcmsProofWithArgs encodes a call to the create_mcms_proof Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) CreateMcmsProofWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_mcms_proof", typeArgsList, typeParamsList, expectedParams, args, []string{
		"mcms::mcms_registry::McmsProof",
	})
}
