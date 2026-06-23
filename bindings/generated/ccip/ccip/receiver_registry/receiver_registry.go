// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_receiver_registry

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

var (
	_ = big.NewInt
)

const FunctionInfo = `[{"package":"ccip","module":"receiver_registry","name":"get_receiver_config","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"receiver_package_id","type":"address"}]},{"package":"ccip","module":"receiver_registry","name":"get_receiver_config_fields","parameters":[{"name":"rc","type":"ReceiverConfig"}]},{"package":"ccip","module":"receiver_registry","name":"get_receiver_info","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"receiver_package_id","type":"address"}]},{"package":"ccip","module":"receiver_registry","name":"initialize","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"}]},{"package":"ccip","module":"receiver_registry","name":"is_registered_receiver","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"receiver_package_id","type":"address"}]},{"package":"ccip","module":"receiver_registry","name":"register_receiver","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"publisher_wrapper","type":"PublisherWrapper<ProofType>"},{"name":"_proof","type":"ProofType"}]},{"package":"ccip","module":"receiver_registry","name":"type_and_version","parameters":null},{"package":"ccip","module":"receiver_registry","name":"unregister_receiver","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"},{"name":"receiver_package_id","type":"address"}]}]`

type IReceiverRegistry interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	RegisterReceiver(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, publisherWrapper bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	UnregisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error)
	IsRegisteredReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error)
	GetReceiverConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error)
	GetReceiverConfigFields(ctx context.Context, opts *bind.CallOpts, rc ReceiverConfig) (*models.SuiTransactionBlockResponse, error)
	GetReceiverInfo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error)
	McmsUnregisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IReceiverRegistryDevInspect
	Encoder() ReceiverRegistryEncoder
	Bound() bind.IBoundContract
}

type IReceiverRegistryDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	IsRegisteredReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (bool, error)
	GetReceiverConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (ReceiverConfig, error)
	GetReceiverConfigFields(ctx context.Context, opts *bind.CallOpts, rc ReceiverConfig) ([]any, error)
	GetReceiverInfo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) ([]any, error)
}

type ReceiverRegistryEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterReceiver(typeArgs []string, ref bind.Object, publisherWrapper bind.Object, proof bind.Object) (*bind.EncodedCall, error)
	RegisterReceiverWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	UnregisterReceiver(ref bind.Object, ownerCap bind.Object, receiverPackageId string) (*bind.EncodedCall, error)
	UnregisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	IsRegisteredReceiver(ref bind.Object, receiverPackageId string) (*bind.EncodedCall, error)
	IsRegisteredReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	GetReceiverConfig(ref bind.Object, receiverPackageId string) (*bind.EncodedCall, error)
	GetReceiverConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetReceiverConfigFields(rc ReceiverConfig) (*bind.EncodedCall, error)
	GetReceiverConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetReceiverInfo(ref bind.Object, receiverPackageId string) (*bind.EncodedCall, error)
	GetReceiverInfoWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsUnregisterReceiver(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsUnregisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
}

type ReceiverRegistryContract struct {
	*bind.BoundContract
	receiverRegistryEncoder
	devInspect *ReceiverRegistryDevInspect
}

type ReceiverRegistryDevInspect struct {
	contract *ReceiverRegistryContract
}

var _ IReceiverRegistry = (*ReceiverRegistryContract)(nil)
var _ IReceiverRegistryDevInspect = (*ReceiverRegistryDevInspect)(nil)

func NewReceiverRegistry(packageID string, chainClient client.BindingsClient) (IReceiverRegistry, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "receiver_registry", chainClient)
	if err != nil {
		return nil, err
	}

	c := &ReceiverRegistryContract{
		BoundContract:           contract,
		receiverRegistryEncoder: receiverRegistryEncoder{BoundContract: contract},
	}
	c.devInspect = &ReceiverRegistryDevInspect{contract: c}
	return c, nil
}

func (c *ReceiverRegistryContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *ReceiverRegistryContract) Encoder() ReceiverRegistryEncoder {
	return c.receiverRegistryEncoder
}

func (c *ReceiverRegistryContract) DevInspect() IReceiverRegistryDevInspect {
	return c.devInspect
}

type ReceiverConfig struct {
	ModuleName    string `move:"0x1::string::String"`
	ProofTypename string `move:"ascii::String"`
}

type ReceiverRegistry struct {
	Id              string      `move:"sui::object::UID"`
	ReceiverConfigs bind.Object `move:"LinkedTable<address, ReceiverConfig>"`
}

type ReceiverRegistered struct {
	ReceiverPackageId  string `move:"address"`
	ReceiverModuleName string `move:"0x1::string::String"`
	ProofTypename      string `move:"ascii::String"`
}

type ReceiverUnregistered struct {
	ReceiverPackageId string `move:"address"`
}

// TypeAndVersion executes the type_and_version Move function.
func (c *ReceiverRegistryContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *ReceiverRegistryContract) Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.Initialize(ref, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterReceiver executes the register_receiver Move function.
func (c *ReceiverRegistryContract) RegisterReceiver(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, publisherWrapper bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.RegisterReceiver(typeArgs, ref, publisherWrapper, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UnregisterReceiver executes the unregister_receiver Move function.
func (c *ReceiverRegistryContract) UnregisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.UnregisterReceiver(ref, ownerCap, receiverPackageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsRegisteredReceiver executes the is_registered_receiver Move function.
func (c *ReceiverRegistryContract) IsRegisteredReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.IsRegisteredReceiver(ref, receiverPackageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetReceiverConfig executes the get_receiver_config Move function.
func (c *ReceiverRegistryContract) GetReceiverConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.GetReceiverConfig(ref, receiverPackageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetReceiverConfigFields executes the get_receiver_config_fields Move function.
func (c *ReceiverRegistryContract) GetReceiverConfigFields(ctx context.Context, opts *bind.CallOpts, rc ReceiverConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.GetReceiverConfigFields(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetReceiverInfo executes the get_receiver_info Move function.
func (c *ReceiverRegistryContract) GetReceiverInfo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.GetReceiverInfo(ref, receiverPackageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsUnregisterReceiver executes the mcms_unregister_receiver Move function.
func (c *ReceiverRegistryContract) McmsUnregisterReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.receiverRegistryEncoder.McmsUnregisterReceiver(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *ReceiverRegistryDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.receiverRegistryEncoder.TypeAndVersion()
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
	var result string
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return "", fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// IsRegisteredReceiver executes the is_registered_receiver Move function using DevInspect to get return values.
//
// Returns: bool
func (d *ReceiverRegistryDevInspect) IsRegisteredReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (bool, error) {
	encoded, err := d.contract.receiverRegistryEncoder.IsRegisteredReceiver(ref, receiverPackageId)
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
	var result bool
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return false, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// GetReceiverConfig executes the get_receiver_config Move function using DevInspect to get return values.
//
// Returns: ReceiverConfig
func (d *ReceiverRegistryDevInspect) GetReceiverConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) (ReceiverConfig, error) {
	encoded, err := d.contract.receiverRegistryEncoder.GetReceiverConfig(ref, receiverPackageId)
	if err != nil {
		return ReceiverConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return ReceiverConfig{}, err
	}
	if len(results) == 0 {
		return ReceiverConfig{}, fmt.Errorf("no return value")
	}
	var result ReceiverConfig
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return ReceiverConfig{}, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// GetReceiverConfigFields executes the get_receiver_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: 0x1::string::String
//	[1]: ascii::String
func (d *ReceiverRegistryDevInspect) GetReceiverConfigFields(ctx context.Context, opts *bind.CallOpts, rc ReceiverConfig) ([]any, error) {
	encoded, err := d.contract.receiverRegistryEncoder.GetReceiverConfigFields(rc)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) != 2 {
		return nil, fmt.Errorf("expected 2 return values, got %d", len(results))
	}
	decoded := make([]any, 2)
	var ret0 string
	if err := bind.DecodeJSONReturn(results[0], &ret0); err != nil {
		return nil, fmt.Errorf("failed to decode return value 0: %w", err)
	}
	decoded[0] = ret0
	var ret1 string
	if err := bind.DecodeJSONReturn(results[1], &ret1); err != nil {
		return nil, fmt.Errorf("failed to decode return value 1: %w", err)
	}
	decoded[1] = ret1
	return decoded, nil
}

// GetReceiverInfo executes the get_receiver_info Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: 0x1::string::String
//	[1]: ascii::String
func (d *ReceiverRegistryDevInspect) GetReceiverInfo(ctx context.Context, opts *bind.CallOpts, ref bind.Object, receiverPackageId string) ([]any, error) {
	encoded, err := d.contract.receiverRegistryEncoder.GetReceiverInfo(ref, receiverPackageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return nil, err
	}
	if len(results) != 2 {
		return nil, fmt.Errorf("expected 2 return values, got %d", len(results))
	}
	decoded := make([]any, 2)
	var ret0 string
	if err := bind.DecodeJSONReturn(results[0], &ret0); err != nil {
		return nil, fmt.Errorf("failed to decode return value 0: %w", err)
	}
	decoded[0] = ret0
	var ret1 string
	if err := bind.DecodeJSONReturn(results[1], &ret1); err != nil {
		return nil, fmt.Errorf("failed to decode return value 1: %w", err)
	}
	decoded[1] = ret1
	return decoded, nil
}

type receiverRegistryEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c receiverRegistryEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
	})
}

// Initialize encodes a call to the initialize Move function.
func (c receiverRegistryEncoder) Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
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
func (c receiverRegistryEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// RegisterReceiver encodes a call to the register_receiver Move function.
func (c receiverRegistryEncoder) RegisterReceiver(typeArgs []string, ref bind.Object, publisherWrapper bind.Object, proof bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"ProofType",
	}
	return c.EncodeCallArgsWithGenerics("register_receiver", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"PublisherWrapper<ProofType>",
		"ProofType",
	}, []any{
		ref,
		publisherWrapper,
		proof,
	}, nil)
}

// RegisterReceiverWithArgs encodes a call to the register_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) RegisterReceiverWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"PublisherWrapper<ProofType>",
		"ProofType",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"ProofType",
	}
	return c.EncodeCallArgsWithGenerics("register_receiver", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// UnregisterReceiver encodes a call to the unregister_receiver Move function.
func (c receiverRegistryEncoder) UnregisterReceiver(ref bind.Object, ownerCap bind.Object, receiverPackageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unregister_receiver", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"address",
	}, []any{
		ref,
		ownerCap,
		receiverPackageId,
	}, nil)
}

// UnregisterReceiverWithArgs encodes a call to the unregister_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) UnregisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("unregister_receiver", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IsRegisteredReceiver encodes a call to the is_registered_receiver Move function.
func (c receiverRegistryEncoder) IsRegisteredReceiver(ref bind.Object, receiverPackageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_registered_receiver", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		receiverPackageId,
	}, []string{
		"bool",
	})
}

// IsRegisteredReceiverWithArgs encodes a call to the is_registered_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) IsRegisteredReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_registered_receiver", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetReceiverConfig encodes a call to the get_receiver_config Move function.
func (c receiverRegistryEncoder) GetReceiverConfig(ref bind.Object, receiverPackageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_receiver_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		receiverPackageId,
	}, []string{
		"ccip::receiver_registry::ReceiverConfig",
	})
}

// GetReceiverConfigWithArgs encodes a call to the get_receiver_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) GetReceiverConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_receiver_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::receiver_registry::ReceiverConfig",
	})
}

// GetReceiverConfigFields encodes a call to the get_receiver_config_fields Move function.
func (c receiverRegistryEncoder) GetReceiverConfigFields(rc ReceiverConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_receiver_config_fields", typeArgsList, typeParamsList, []string{
		"ccip::receiver_registry::ReceiverConfig",
	}, []any{
		rc,
	}, []string{
		"0x1::string::String",
		"ascii::String",
	})
}

// GetReceiverConfigFieldsWithArgs encodes a call to the get_receiver_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) GetReceiverConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip::receiver_registry::ReceiverConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_receiver_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
		"ascii::String",
	})
}

// GetReceiverInfo encodes a call to the get_receiver_info Move function.
func (c receiverRegistryEncoder) GetReceiverInfo(ref bind.Object, receiverPackageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_receiver_info", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		receiverPackageId,
	}, []string{
		"0x1::string::String",
		"ascii::String",
	})
}

// GetReceiverInfoWithArgs encodes a call to the get_receiver_info Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) GetReceiverInfoWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_receiver_info", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::string::String",
		"ascii::String",
	})
}

// McmsUnregisterReceiver encodes a call to the mcms_unregister_receiver Move function.
func (c receiverRegistryEncoder) McmsUnregisterReceiver(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_unregister_receiver", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsUnregisterReceiverWithArgs encodes a call to the mcms_unregister_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c receiverRegistryEncoder) McmsUnregisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_unregister_receiver", typeArgsList, typeParamsList, expectedParams, args, nil)
}
