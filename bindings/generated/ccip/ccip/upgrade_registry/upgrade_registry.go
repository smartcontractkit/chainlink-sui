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

const FunctionInfo = `[{"package":"ccip","module":"upgrade_registry","name":"block_function","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"module_name","type":"0x1::string::String"},{"name":"function_name","type":"0x1::string::String"},{"name":"version","type":"u8"}]},{"package":"ccip","module":"upgrade_registry","name":"block_version","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"module_name","type":"0x1::string::String"},{"name":"version","type":"u8"}]},{"package":"ccip","module":"upgrade_registry","name":"get_module_restrictions","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"module_name","type":"0x1::string::String"}]},{"package":"ccip","module":"upgrade_registry","name":"initialize","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"}]},{"package":"ccip","module":"upgrade_registry","name":"is_function_allowed","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"module_name","type":"0x1::string::String"},{"name":"function_name","type":"0x1::string::String"},{"name":"version","type":"u8"}]},{"package":"ccip","module":"upgrade_registry","name":"unblock_function","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"module_name","type":"0x1::string::String"},{"name":"function_name","type":"0x1::string::String"},{"name":"version","type":"u8"}]},{"package":"ccip","module":"upgrade_registry","name":"unblock_version","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"module_name","type":"0x1::string::String"},{"name":"version","type":"u8"}]},{"package":"ccip","module":"upgrade_registry","name":"verify_function_allowed","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"module_name","type":"0x1::string::String"},{"name":"function_name","type":"0x1::string::String"},{"name":"version","type":"u8"}]}]`

type IUpgradeRegistry interface {
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	BlockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*models.SuiTransactionBlockResponse, error)
	UnblockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*models.SuiTransactionBlockResponse, error)
	BlockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error)
	UnblockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error)
	GetModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string) (*models.SuiTransactionBlockResponse, error)
	IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error)
	VerifyFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error)
	McmsBlockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsUnblockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsBlockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsUnblockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IUpgradeRegistryDevInspect
	Encoder() UpgradeRegistryEncoder
	Bound() bind.IBoundContract
}

type IUpgradeRegistryDevInspect interface {
	GetModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string) ([][]byte, error)
	IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, version byte) (bool, error)
}

type UpgradeRegistryEncoder interface {
	Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	BlockVersion(ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*bind.EncodedCall, error)
	BlockVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	UnblockVersion(ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*bind.EncodedCall, error)
	UnblockVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	BlockFunction(ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error)
	BlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	UnblockFunction(ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error)
	UnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	GetModuleRestrictions(ref bind.Object, moduleName string) (*bind.EncodedCall, error)
	GetModuleRestrictionsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsFunctionAllowed(ref bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error)
	IsFunctionAllowedWithArgs(args ...any) (*bind.EncodedCall, error)
	VerifyFunctionAllowed(ref bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error)
	VerifyFunctionAllowedWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsBlockVersion(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsBlockVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsUnblockVersion(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsUnblockVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsBlockFunction(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsBlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsUnblockFunction(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsUnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewUpgradeRegistry(packageID string, client sui.ISuiAPI) (IUpgradeRegistry, error) {
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

func (c *UpgradeRegistryContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *UpgradeRegistryContract) Encoder() UpgradeRegistryEncoder {
	return c.upgradeRegistryEncoder
}

func (c *UpgradeRegistryContract) DevInspect() IUpgradeRegistryDevInspect {
	return c.devInspect
}

type VersionBlocked struct {
	ModuleName string `move:"0x1::string::String"`
	Version    byte   `move:"u8"`
}

type VersionUnblocked struct {
	ModuleName string `move:"0x1::string::String"`
	Version    byte   `move:"u8"`
}

type FunctionBlocked struct {
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Version      byte   `move:"u8"`
}

type FunctionUnblocked struct {
	ModuleName   string `move:"0x1::string::String"`
	FunctionName string `move:"0x1::string::String"`
	Version      byte   `move:"u8"`
}

type UpgradeRegistry struct {
	Id                   string      `move:"sui::object::UID"`
	FunctionRestrictions bind.Object `move:"Table<String, vector<vector<u8>>>"`
}

type McmsCallback struct {
}

func init() {
	bind.RegisterStructDecoder("ccip::upgrade_registry::VersionBlocked", func(data []byte) (interface{}, error) {
		var result VersionBlocked
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for VersionBlocked
	bind.RegisterStructDecoder("vector<ccip::upgrade_registry::VersionBlocked>", func(data []byte) (interface{}, error) {
		var results []VersionBlocked
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::VersionUnblocked", func(data []byte) (interface{}, error) {
		var result VersionUnblocked
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for VersionUnblocked
	bind.RegisterStructDecoder("vector<ccip::upgrade_registry::VersionUnblocked>", func(data []byte) (interface{}, error) {
		var results []VersionUnblocked
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::FunctionBlocked", func(data []byte) (interface{}, error) {
		var result FunctionBlocked
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for FunctionBlocked
	bind.RegisterStructDecoder("vector<ccip::upgrade_registry::FunctionBlocked>", func(data []byte) (interface{}, error) {
		var results []FunctionBlocked
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::FunctionUnblocked", func(data []byte) (interface{}, error) {
		var result FunctionUnblocked
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for FunctionUnblocked
	bind.RegisterStructDecoder("vector<ccip::upgrade_registry::FunctionUnblocked>", func(data []byte) (interface{}, error) {
		var results []FunctionUnblocked
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::UpgradeRegistry", func(data []byte) (interface{}, error) {
		var result UpgradeRegistry
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for UpgradeRegistry
	bind.RegisterStructDecoder("vector<ccip::upgrade_registry::UpgradeRegistry>", func(data []byte) (interface{}, error) {
		var results []UpgradeRegistry
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::upgrade_registry::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsCallback
	bind.RegisterStructDecoder("vector<ccip::upgrade_registry::McmsCallback>", func(data []byte) (interface{}, error) {
		var results []McmsCallback
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
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

// BlockVersion executes the block_version Move function.
func (c *UpgradeRegistryContract) BlockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.BlockVersion(ref, ownerCap, moduleName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UnblockVersion executes the unblock_version Move function.
func (c *UpgradeRegistryContract) UnblockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.UnblockVersion(ref, ownerCap, moduleName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// BlockFunction executes the block_function Move function.
func (c *UpgradeRegistryContract) BlockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.BlockFunction(ref, ownerCap, moduleName, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UnblockFunction executes the unblock_function Move function.
func (c *UpgradeRegistryContract) UnblockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.UnblockFunction(ref, ownerCap, moduleName, functionName, version)
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

// IsFunctionAllowed executes the is_function_allowed Move function.
func (c *UpgradeRegistryContract) IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.IsFunctionAllowed(ref, moduleName, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// VerifyFunctionAllowed executes the verify_function_allowed Move function.
func (c *UpgradeRegistryContract) VerifyFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, version byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.VerifyFunctionAllowed(ref, moduleName, functionName, version)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsBlockVersion executes the mcms_block_version Move function.
func (c *UpgradeRegistryContract) McmsBlockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.McmsBlockVersion(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsUnblockVersion executes the mcms_unblock_version Move function.
func (c *UpgradeRegistryContract) McmsUnblockVersion(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.McmsUnblockVersion(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsBlockFunction executes the mcms_block_function Move function.
func (c *UpgradeRegistryContract) McmsBlockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.McmsBlockFunction(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsUnblockFunction executes the mcms_unblock_function Move function.
func (c *UpgradeRegistryContract) McmsUnblockFunction(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.upgradeRegistryEncoder.McmsUnblockFunction(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetModuleRestrictions executes the get_module_restrictions Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *UpgradeRegistryDevInspect) GetModuleRestrictions(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string) ([][]byte, error) {
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
	result, ok := results[0].([][]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected [][]byte, got %T", results[0])
	}
	return result, nil
}

// IsFunctionAllowed executes the is_function_allowed Move function using DevInspect to get return values.
//
// Returns: bool
func (d *UpgradeRegistryDevInspect) IsFunctionAllowed(ctx context.Context, opts *bind.CallOpts, ref bind.Object, moduleName string, functionName string, version byte) (bool, error) {
	encoded, err := d.contract.upgradeRegistryEncoder.IsFunctionAllowed(ref, moduleName, functionName, version)
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

// BlockVersion encodes a call to the block_version Move function.
func (c upgradeRegistryEncoder) BlockVersion(ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("block_version", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"u8",
	}, []any{
		ref,
		ownerCap,
		moduleName,
		version,
	}, nil)
}

// BlockVersionWithArgs encodes a call to the block_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) BlockVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("block_version", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// UnblockVersion encodes a call to the unblock_version Move function.
func (c upgradeRegistryEncoder) UnblockVersion(ref bind.Object, ownerCap bind.Object, moduleName string, version byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unblock_version", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"u8",
	}, []any{
		ref,
		ownerCap,
		moduleName,
		version,
	}, nil)
}

// UnblockVersionWithArgs encodes a call to the unblock_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) UnblockVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unblock_version", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// BlockFunction encodes a call to the block_function Move function.
func (c upgradeRegistryEncoder) BlockFunction(ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("block_function", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"0x1::string::String",
		"u8",
	}, []any{
		ref,
		ownerCap,
		moduleName,
		functionName,
		version,
	}, nil)
}

// BlockFunctionWithArgs encodes a call to the block_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) BlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"0x1::string::String",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("block_function", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// UnblockFunction encodes a call to the unblock_function Move function.
func (c upgradeRegistryEncoder) UnblockFunction(ref bind.Object, ownerCap bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unblock_function", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"0x1::string::String",
		"u8",
	}, []any{
		ref,
		ownerCap,
		moduleName,
		functionName,
		version,
	}, nil)
}

// UnblockFunctionWithArgs encodes a call to the unblock_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) UnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"0x1::string::String",
		"0x1::string::String",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unblock_function", typeArgsList, typeParamsList, expectedParams, args, nil)
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
		"vector<vector<u8>>",
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
		"vector<vector<u8>>",
	})
}

// IsFunctionAllowed encodes a call to the is_function_allowed Move function.
func (c upgradeRegistryEncoder) IsFunctionAllowed(ref bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_function_allowed", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"0x1::string::String",
		"u8",
	}, []any{
		ref,
		moduleName,
		functionName,
		version,
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
		"u8",
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

// VerifyFunctionAllowed encodes a call to the verify_function_allowed Move function.
func (c upgradeRegistryEncoder) VerifyFunctionAllowed(ref bind.Object, moduleName string, functionName string, version byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify_function_allowed", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"0x1::string::String",
		"u8",
	}, []any{
		ref,
		moduleName,
		functionName,
		version,
	}, nil)
}

// VerifyFunctionAllowedWithArgs encodes a call to the verify_function_allowed Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) VerifyFunctionAllowedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"0x1::string::String",
		"0x1::string::String",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("verify_function_allowed", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsBlockVersion encodes a call to the mcms_block_version Move function.
func (c upgradeRegistryEncoder) McmsBlockVersion(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_block_version", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsBlockVersionWithArgs encodes a call to the mcms_block_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) McmsBlockVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_block_version", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsUnblockVersion encodes a call to the mcms_unblock_version Move function.
func (c upgradeRegistryEncoder) McmsUnblockVersion(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_unblock_version", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsUnblockVersionWithArgs encodes a call to the mcms_unblock_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) McmsUnblockVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_unblock_version", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsBlockFunction encodes a call to the mcms_block_function Move function.
func (c upgradeRegistryEncoder) McmsBlockFunction(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_block_function", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsBlockFunctionWithArgs encodes a call to the mcms_block_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) McmsBlockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_block_function", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsUnblockFunction encodes a call to the mcms_unblock_function Move function.
func (c upgradeRegistryEncoder) McmsUnblockFunction(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_unblock_function", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsUnblockFunctionWithArgs encodes a call to the mcms_unblock_function Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c upgradeRegistryEncoder) McmsUnblockFunctionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_unblock_function", typeArgsList, typeParamsList, expectedParams, args, nil)
}
