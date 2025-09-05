// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_upgrade_registry

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

type IUpgradeRegistry interface {
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	UpdateFunctionRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, moduleName string, functionName string, blockedVersions []uint64) (*models.SuiTransactionBlockResponse, error)
	GetFunctionRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string) (*models.SuiTransactionBlockResponse, error)
	IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, contractVersion uint64) (*models.SuiTransactionBlockResponse, error)
	UpdateModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, moduleName string, blockedVersions []uint64) (*models.SuiTransactionBlockResponse, error)
	GetModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string) (*models.SuiTransactionBlockResponse, error)
	IsModuleAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, contractVersion uint64) (*models.SuiTransactionBlockResponse, error)
	GetPackageHistory(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageName string) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IUpgradeRegistryDevInspect
	Encoder() UpgradeRegistryEncoder
}

type IUpgradeRegistryDevInspect interface {
	GetFunctionRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string) ([]uint64, error)
	IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, contractVersion uint64) (bool, error)
	GetModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string) ([]uint64, error)
	IsModuleAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, contractVersion uint64) (bool, error)
	GetPackageHistory(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageName string) ([]any, error)
}

type UpgradeRegistryEncoder interface {
	Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	UpdateFunctionRestrictions(ref bind.Object, param bind.Object, moduleName string, functionName string, blockedVersions []uint64) (*bind.EncodedCall, error)
	UpdateFunctionRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetFunctionRestrictions(ref bind.Object, moduleName string, functionName string) (*bind.EncodedCall, error)
	GetFunctionRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsFunctionAllowed(ref bind.Object, moduleName string, functionName string, contractVersion uint64) (*bind.EncodedCall, error)
	IsFunctionAllowedWithArgs(args ...any) (*bind.EncodedCall, error)
	UpdateModuleRestrictions(ref bind.Object, param bind.Object, moduleName string, blockedVersions []uint64) (*bind.EncodedCall, error)
	UpdateModuleRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetModuleRestrictions(ref bind.Object, moduleName string) (*bind.EncodedCall, error)
	GetModuleRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsModuleAllowed(ref bind.Object, moduleName string, contractVersion uint64) (*bind.EncodedCall, error)
	IsModuleAllowedWithArgs(args ...any) (*bind.EncodedCall, error)
	GetPackageHistory(ref bind.Object, packageName string) (*bind.EncodedCall, error)
	GetPackageHistoryWithArgs(args ...any) (*bind.EncodedCall, error)
}

type UpgradeRegistryContract struct {
	*bind.BoundContract
	upgradeRegistryEncoder
	devInspect *UpgradeRegistryDevInspect
}

type UpgradeRegistryDevInspect struct {
	contract *UpgradeRegistryContract
}

var _ IUpgradeRegistry = (*UpgradeRegistryContract)(nil)
var _ IUpgradeRegistryDevInspect = (*UpgradeRegistryDevInspect)(nil)

func NewUpgradeRegistry(packageID string, client sui.ISuiAPI) (*UpgradeRegistryContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "upgrade_registry", client)
	if err != nil {
		return nil, err
	}

	c := &UpgradeRegistryContract{
		BoundContract:          contract,
		upgradeRegistryEncoder: upgradeRegistryEncoder{BoundContract: contract},
	}
	c.devInspect = &UpgradeRegistryDevInspect{contract: c}
	return c, nil
}

func (c *UpgradeRegistryContract) Encoder() UpgradeRegistryEncoder {
	return c.upgradeRegistryEncoder
}

func (c *UpgradeRegistryContract) DevInspect() IUpgradeRegistryDevInspect {
	return c.devInspect
}

type FunctionRestrictionsUpdated struct {
	ModuleName      string   `move:"0x1::string::String"`
	FunctionName    string   `move:"0x1::string::String"`
	BlockedVersions []uint64 `move:"vector<u64>"`
}

type ModuleRestrictionsUpdated struct {
	ModuleName      string   `move:"0x1::string::String"`
	BlockedVersions []uint64 `move:"vector<u64>"`
}

type PackageHistory struct {
	PackageId string `move:"address"`
	Version   uint64 `move:"u64"`
	Timestamp uint64 `move:"u64"`
}

type FunctionKey struct {
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
}

type UpgradeRegistry struct {
	Id                   string      `move:"sui::object::UID"`
	FunctionRestrictions bind.Object `move:"Table<FunctionKey, vector<u64>>"`
	ModuleRestrictions   bind.Object `move:"Table<String, vector<u64>>"`
	PackageHistory       bind.Object `move:"Table<String, vector<PackageHistory>>"`
}

type bcsPackageHistory struct {
	PackageId [32]byte
	Version   uint64
	Timestamp uint64
}

func convertPackageHistoryFromBCS(bcs bcsPackageHistory) (PackageHistory, error) {

	return PackageHistory{
		PackageId: fmt.Sprintf("0x%x", bcs.PackageId),
		Version:   bcs.Version,
		Timestamp: bcs.Timestamp,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip::upgrade_registry::FunctionRestrictionsUpdated", func(data []byte) (interface{}, error) {
		var result FunctionRestrictionsUpdated
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::ModuleRestrictionsUpdated", func(data []byte) (interface{}, error) {
		var result ModuleRestrictionsUpdated
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::PackageHistory", func(data []byte) (interface{}, error) {
		var temp bcsPackageHistory
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertPackageHistoryFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::FunctionKey", func(data []byte) (interface{}, error) {
		var result FunctionKey
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::UpgradeRegistry", func(data []byte) (interface{}, error) {
		var result UpgradeRegistry
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// Initialize executes the initialize Move function.
func (c *UpgradeRegistryContract) Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.Initialize(ref, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UpdateFunctionRestrictions executes the update_function_restrictions Move function.
func (c *UpgradeRegistryContract) UpdateFunctionRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, moduleName string, functionName string, blockedVersions []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.UpdateFunctionRestrictions(ref, param, moduleName, functionName, blockedVersions)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFunctionRestrictions executes the get_function_restrictions Move function.
func (c *UpgradeRegistryContract) GetFunctionRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.GetFunctionRestrictions(ref, moduleName, functionName)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsFunctionAllowed executes the is_function_allowed Move function.
func (c *UpgradeRegistryContract) IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, contractVersion uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.IsFunctionAllowed(ref, moduleName, functionName, contractVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UpdateModuleRestrictions executes the update_module_restrictions Move function.
func (c *UpgradeRegistryContract) UpdateModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, moduleName string, blockedVersions []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.UpdateModuleRestrictions(ref, param, moduleName, blockedVersions)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetModuleRestrictions executes the get_module_restrictions Move function.
func (c *UpgradeRegistryContract) GetModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.GetModuleRestrictions(ref, moduleName)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsModuleAllowed executes the is_module_allowed Move function.
func (c *UpgradeRegistryContract) IsModuleAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, contractVersion uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.IsModuleAllowed(ref, moduleName, contractVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetPackageHistory executes the get_package_history Move function.
func (c *UpgradeRegistryContract) GetPackageHistory(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageName string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.GetPackageHistory(ref, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFunctionRestrictions executes the get_function_restrictions Move function using DevInspect to get return values.
//
// Returns: vector<u64>
func (d *UpgradeRegistryDevInspect) GetFunctionRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string) ([]uint64, error) {
	encoded, err := d.contract.upgradeRegistryEncoder.GetFunctionRestrictions(ref, moduleName, functionName)
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
	result, ok := results[0].([]uint64)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []uint64, got %T", results[0])
	}
	return result, nil
}

// IsFunctionAllowed executes the is_function_allowed Move function using DevInspect to get return values.
//
// Returns: bool
func (d *UpgradeRegistryDevInspect) IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, contractVersion uint64) (bool, error) {
	encoded, err := d.contract.upgradeRegistryEncoder.IsFunctionAllowed(ref, moduleName, functionName, contractVersion)
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

// GetModuleRestrictions executes the get_module_restrictions Move function using DevInspect to get return values.
//
// Returns: vector<u64>
func (d *UpgradeRegistryDevInspect) GetModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string) ([]uint64, error) {
	encoded, err := d.contract.upgradeRegistryEncoder.GetModuleRestrictions(ref, moduleName)
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
	result, ok := results[0].([]uint64)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []uint64, got %T", results[0])
	}
	return result, nil
}

// IsModuleAllowed executes the is_module_allowed Move function using DevInspect to get return values.
//
// Returns: bool
func (d *UpgradeRegistryDevInspect) IsModuleAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, contractVersion uint64) (bool, error) {
	encoded, err := d.contract.upgradeRegistryEncoder.IsModuleAllowed(ref, moduleName, contractVersion)
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

// GetPackageHistory executes the get_package_history Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<address>
//	[1]: vector<u64>
//	[2]: vector<u64>
func (d *UpgradeRegistryDevInspect) GetPackageHistory(ctx context.Context, opts *bind.CallOpts, ref bind.Object, packageName string) ([]any, error) {
	encoded, err := d.contract.upgradeRegistryEncoder.GetPackageHistory(ref, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

type upgradeRegistryEncoder struct {
	*bind.BoundContract
}

// Initialize encodes a call to the initialize Move function.
func (c upgradeRegistryEncoder) Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}, []any{
		ref,
		ownerCap,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// UpdateFunctionRestrictions encodes a call to the update_function_restrictions Move function.
func (c upgradeRegistryEncoder) UpdateFunctionRestrictions(ref bind.Object, param bind.Object, moduleName string, functionName string, blockedVersions []uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("update_function_restrictions", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u64>",
	}, []any{
		ref,
		param,
		moduleName,
		functionName,
		blockedVersions,
	}, nil)
}

// UpdateFunctionRestrictionsWithArgs encodes a call to the update_function_restrictions Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) UpdateFunctionRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u64>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("update_function_restrictions", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetFunctionRestrictions encodes a call to the get_function_restrictions Move function.
func (c upgradeRegistryEncoder) GetFunctionRestrictions(ref bind.Object, moduleName string, functionName string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_function_restrictions", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"0x1::string::String",
	}, []any{
		ref,
		moduleName,
		functionName,
	}, []string{
		"vector<u64>",
	})
}

// GetFunctionRestrictionsWithArgs encodes a call to the get_function_restrictions Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) GetFunctionRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"0x1::string::String",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_function_restrictions", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
	})
}

// IsFunctionAllowed encodes a call to the is_function_allowed Move function.
func (c upgradeRegistryEncoder) IsFunctionAllowed(ref bind.Object, moduleName string, functionName string, contractVersion uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_function_allowed", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"0x1::string::String",
		"u64",
	}, []any{
		ref,
		moduleName,
		functionName,
		contractVersion,
	}, []string{
		"bool",
	})
}

// IsFunctionAllowedWithArgs encodes a call to the is_function_allowed Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) IsFunctionAllowedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"0x1::string::String",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_function_allowed", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// UpdateModuleRestrictions encodes a call to the update_module_restrictions Move function.
func (c upgradeRegistryEncoder) UpdateModuleRestrictions(ref bind.Object, param bind.Object, moduleName string, blockedVersions []uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("update_module_restrictions", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"vector<u64>",
	}, []any{
		ref,
		param,
		moduleName,
		blockedVersions,
	}, nil)
}

// UpdateModuleRestrictionsWithArgs encodes a call to the update_module_restrictions Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) UpdateModuleRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"vector<u64>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("update_module_restrictions", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetModuleRestrictions encodes a call to the get_module_restrictions Move function.
func (c upgradeRegistryEncoder) GetModuleRestrictions(ref bind.Object, moduleName string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_module_restrictions", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"0x1::string::String",
	}, []any{
		ref,
		moduleName,
	}, []string{
		"vector<u64>",
	})
}

// GetModuleRestrictionsWithArgs encodes a call to the get_module_restrictions Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) GetModuleRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"0x1::string::String",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_module_restrictions", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
	})
}

// IsModuleAllowed encodes a call to the is_module_allowed Move function.
func (c upgradeRegistryEncoder) IsModuleAllowed(ref bind.Object, moduleName string, contractVersion uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_module_allowed", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"u64",
	}, []any{
		ref,
		moduleName,
		contractVersion,
	}, []string{
		"bool",
	})
}

// IsModuleAllowedWithArgs encodes a call to the is_module_allowed Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) IsModuleAllowedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_module_allowed", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetPackageHistory encodes a call to the get_package_history Move function.
func (c upgradeRegistryEncoder) GetPackageHistory(ref bind.Object, packageName string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_package_history", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"0x1::string::String",
	}, []any{
		ref,
		packageName,
	}, []string{
		"vector<address>",
		"vector<u64>",
		"vector<u64>",
	})
}

// GetPackageHistoryWithArgs encodes a call to the get_package_history Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) GetPackageHistoryWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"0x1::string::String",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_package_history", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
		"vector<u64>",
		"vector<u64>",
	})
}
