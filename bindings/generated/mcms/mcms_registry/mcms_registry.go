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

const FunctionInfo = `[{"package":"mcms","module":"mcms_registry","name":"add_allowed_modules","parameters":[{"name":"registry","type":"Registry"},{"name":"_proof","type":"T"},{"name":"new_allowed_modules","type":"vector<vector<u8>>"}]},{"package":"mcms","module":"mcms_registry","name":"borrow_owner_cap","parameters":[{"name":"registry","type":"Registry"}]},{"package":"mcms","module":"mcms_registry","name":"create_executing_callback_params","parameters":[{"name":"target","type":"address"},{"name":"module_name","type":"0x1::string::String"},{"name":"function_name","type":"0x1::string::String"},{"name":"data","type":"vector<u8>"},{"name":"batch_id","type":"vector<u8>"},{"name":"sequence_number","type":"u64"},{"name":"total_in_batch","type":"u64"}]},{"package":"mcms","module":"mcms_registry","name":"create_publisher_wrapper","parameters":[{"name":"publisher","type":"Publisher"},{"name":"_proof","type":"T"}]},{"package":"mcms","module":"mcms_registry","name":"data","parameters":[{"name":"params","type":"ExecutingCallbackParams"}]},{"package":"mcms","module":"mcms_registry","name":"function_name","parameters":[{"name":"params","type":"ExecutingCallbackParams"}]},{"package":"mcms","module":"mcms_registry","name":"get_accept_ownership_data","parameters":[{"name":"registry","type":"Registry"},{"name":"params","type":"ExecutingCallbackParams"},{"name":"_proof","type":"T"}]},{"package":"mcms","module":"mcms_registry","name":"get_allowed_modules","parameters":[{"name":"registry","type":"Registry"},{"name":"package_address","type":"ascii::String"}]},{"package":"mcms","module":"mcms_registry","name":"get_callback_params_from_mcms","parameters":[{"name":"registry","type":"Registry"},{"name":"params","type":"ExecutingCallbackParams"}]},{"package":"mcms","module":"mcms_registry","name":"get_callback_params_with_caps","parameters":[{"name":"registry","type":"Registry"},{"name":"_proof","type":"T"},{"name":"params","type":"ExecutingCallbackParams"}]},{"package":"mcms","module":"mcms_registry","name":"get_multisig_address","parameters":null},{"package":"mcms","module":"mcms_registry","name":"get_multisig_address_ascii","parameters":null},{"package":"mcms","module":"mcms_registry","name":"get_next_expected_sequence","parameters":[{"name":"registry","type":"Registry"},{"name":"batch_id","type":"vector<u8>"}]},{"package":"mcms","module":"mcms_registry","name":"get_registered_proof_type","parameters":[{"name":"registry","type":"Registry"},{"name":"package_address","type":"ascii::String"}]},{"package":"mcms","module":"mcms_registry","name":"is_batch_completed","parameters":[{"name":"registry","type":"Registry"},{"name":"batch_id","type":"vector<u8>"}]},{"package":"mcms","module":"mcms_registry","name":"is_package_registered","parameters":[{"name":"registry","type":"Registry"},{"name":"package_address","type":"ascii::String"}]},{"package":"mcms","module":"mcms_registry","name":"module_name","parameters":[{"name":"params","type":"ExecutingCallbackParams"}]},{"package":"mcms","module":"mcms_registry","name":"register_entrypoint","parameters":[{"name":"registry","type":"Registry"},{"name":"publisher_wrapper","type":"PublisherWrapper<T>"},{"name":"_proof","type":"T"},{"name":"package_cap","type":"C"},{"name":"allowed_modules","type":"vector<vector<u8>>"}]},{"package":"mcms","module":"mcms_registry","name":"release_cap","parameters":[{"name":"registry","type":"Registry"},{"name":"_proof","type":"T"}]},{"package":"mcms","module":"mcms_registry","name":"remove_allowed_modules","parameters":[{"name":"registry","type":"Registry"},{"name":"_proof","type":"T"},{"name":"modules_to_remove","type":"vector<vector<u8>>"}]},{"package":"mcms","module":"mcms_registry","name":"target","parameters":[{"name":"params","type":"ExecutingCallbackParams"}]}]`

type IMcmsRegistry interface {
	CreatePublisherWrapper(ctx context.Context, opts *bind.CallOpts, typeArgs []string, publisher bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	RegisterEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, publisherWrapper PublisherWrapper, proof bind.Object, packageCap bind.Object, allowedModules [][]byte) (*models.SuiTransactionBlockResponse, error)
	AddAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, newAllowedModules [][]byte) (*models.SuiTransactionBlockResponse, error)
	RemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, modulesToRemove [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetCallbackParamsWithCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	BorrowOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetAcceptOwnershipData(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params ExecutingCallbackParams, proof bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte, batchId []byte, sequenceNumber uint64, totalInBatch uint64) (*models.SuiTransactionBlockResponse, error)
	IsPackageRegistered(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	GetRegisteredProofType(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	GetAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error)
	Target(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	ModuleName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	FunctionName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	Data(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error)
	IsBatchCompleted(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (*models.SuiTransactionBlockResponse, error)
	GetNextExpectedSequence(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (*models.SuiTransactionBlockResponse, error)
	GetMultisigAddress(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	GetMultisigAddressAscii(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMcmsRegistryDevInspect
	Encoder() McmsRegistryEncoder
	Bound() bind.IBoundContract
}

type IMcmsRegistryDevInspect interface {
	CreatePublisherWrapper(ctx context.Context, opts *bind.CallOpts, typeArgs []string, publisher bind.Object, proof bind.Object) (any, error)
	GetCallbackParamsWithCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) ([]any, error)
	ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object) (any, error)
	BorrowOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object) (bind.Object, error)
	GetAcceptOwnershipData(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params ExecutingCallbackParams, proof bind.Object) ([]byte, error)
	GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params ExecutingCallbackParams) ([]any, error)
	CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte, batchId []byte, sequenceNumber uint64, totalInBatch uint64) (ExecutingCallbackParams, error)
	IsPackageRegistered(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (bool, error)
	GetRegisteredProofType(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (bind.Object, error)
	GetAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) ([][]byte, error)
	Target(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error)
	ModuleName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error)
	FunctionName(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) (string, error)
	Data(ctx context.Context, opts *bind.CallOpts, params ExecutingCallbackParams) ([]byte, error)
	IsBatchCompleted(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (bool, error)
	GetNextExpectedSequence(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (uint64, error)
	GetMultisigAddress(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetMultisigAddressAscii(ctx context.Context, opts *bind.CallOpts) (string, error)
}

type McmsRegistryEncoder interface {
	CreatePublisherWrapper(typeArgs []string, publisher bind.Object, proof bind.Object) (*bind.EncodedCall, error)
	CreatePublisherWrapperWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	RegisterEntrypoint(typeArgs []string, registry bind.Object, publisherWrapper PublisherWrapper, proof bind.Object, packageCap bind.Object, allowedModules [][]byte) (*bind.EncodedCall, error)
	RegisterEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	AddAllowedModules(typeArgs []string, registry bind.Object, proof bind.Object, newAllowedModules [][]byte) (*bind.EncodedCall, error)
	AddAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	RemoveAllowedModules(typeArgs []string, registry bind.Object, proof bind.Object, modulesToRemove [][]byte) (*bind.EncodedCall, error)
	RemoveAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetCallbackParamsWithCaps(typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*bind.EncodedCall, error)
	GetCallbackParamsWithCapsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReleaseCap(typeArgs []string, registry bind.Object, proof bind.Object) (*bind.EncodedCall, error)
	ReleaseCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	BorrowOwnerCap(typeArgs []string, registry bind.Object) (*bind.EncodedCall, error)
	BorrowOwnerCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetAcceptOwnershipData(typeArgs []string, registry bind.Object, params ExecutingCallbackParams, proof bind.Object) (*bind.EncodedCall, error)
	GetAcceptOwnershipDataWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetCallbackParamsFromMcms(registry bind.Object, params ExecutingCallbackParams) (*bind.EncodedCall, error)
	GetCallbackParamsFromMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateExecutingCallbackParams(target string, moduleName string, functionName string, data []byte, batchId []byte, sequenceNumber uint64, totalInBatch uint64) (*bind.EncodedCall, error)
	CreateExecutingCallbackParamsWithArgs(args ...any) (*bind.EncodedCall, error)
	IsPackageRegistered(registry bind.Object, packageAddress string) (*bind.EncodedCall, error)
	IsPackageRegisteredWithArgs(args ...any) (*bind.EncodedCall, error)
	GetRegisteredProofType(registry bind.Object, packageAddress string) (*bind.EncodedCall, error)
	GetRegisteredProofTypeWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowedModules(registry bind.Object, packageAddress string) (*bind.EncodedCall, error)
	GetAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error)
	Target(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	TargetWithArgs(args ...any) (*bind.EncodedCall, error)
	ModuleName(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	ModuleNameWithArgs(args ...any) (*bind.EncodedCall, error)
	FunctionName(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	FunctionNameWithArgs(args ...any) (*bind.EncodedCall, error)
	Data(params ExecutingCallbackParams) (*bind.EncodedCall, error)
	DataWithArgs(args ...any) (*bind.EncodedCall, error)
	IsBatchCompleted(registry bind.Object, batchId []byte) (*bind.EncodedCall, error)
	IsBatchCompletedWithArgs(args ...any) (*bind.EncodedCall, error)
	GetNextExpectedSequence(registry bind.Object, batchId []byte) (*bind.EncodedCall, error)
	GetNextExpectedSequenceWithArgs(args ...any) (*bind.EncodedCall, error)
	GetMultisigAddress() (*bind.EncodedCall, error)
	GetMultisigAddressWithArgs(args ...any) (*bind.EncodedCall, error)
	GetMultisigAddressAscii() (*bind.EncodedCall, error)
	GetMultisigAddressAsciiWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewMcmsRegistry(packageID string, client sui.ISuiAPI) (IMcmsRegistry, error) {
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

func (c *McmsRegistryContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *McmsRegistryContract) Encoder() McmsRegistryEncoder {
	return c.mcmsRegistryEncoder
}

func (c *McmsRegistryContract) DevInspect() IMcmsRegistryDevInspect {
	return c.devInspect
}

type Registry struct {
	Id                   string      `move:"sui::object::UID"`
	PackageCaps          bind.Object `move:"Bag"`
	RegisteredProofTypes bind.Object `move:"Table<ascii::String, TypeName>"`
	ProofTypeToPackage   bind.Object `move:"Table<TypeName, ascii::String>"`
	AllowedModules       bind.Object `move:"Table<ascii::String, vector<vector<u8>>>"`
	BatchExecution       bind.Object `move:"Table<vector<u8>, BatchExecutionState>"`
	CompletedBatches     bind.Object `move:"Table<vector<u8>, bool>"`
}

type PublisherWrapper struct {
	PackageAddress string `move:"ascii::String"`
}

type BatchExecutionState struct {
	TotalCallbacks       uint64 `move:"u64"`
	NextExpectedSequence uint64 `move:"u64"`
}

type ExecutingCallbackParams struct {
	Target         string `move:"address"`
	ModuleName     string `move:"0x1::string::String"`
	FunctionName   string `move:"0x1::string::String"`
	Data           []byte `move:"vector<u8>"`
	BatchId        []byte `move:"vector<u8>"`
	SequenceNumber uint64 `move:"u64"`
	TotalInBatch   uint64 `move:"u64"`
}

type EntrypointRegistered struct {
	RegistryId     bind.Object `move:"ID"`
	AccountAddress string      `move:"ascii::String"`
	AllowedModules [][]byte    `move:"vector<vector<u8>>"`
	ProofType      bind.Object `move:"TypeName"`
}

type ModulesAdded struct {
	RegistryId     bind.Object `move:"ID"`
	PackageAddress string      `move:"ascii::String"`
	ModuleNames    [][]byte    `move:"vector<vector<u8>>"`
}

type ModulesRemoved struct {
	RegistryId     bind.Object `move:"ID"`
	PackageAddress string      `move:"ascii::String"`
	ModuleNames    [][]byte    `move:"vector<vector<u8>>"`
}

type MCMS_REGISTRY struct {
}

type bcsExecutingCallbackParams struct {
	Target         [32]byte
	ModuleName     string
	FunctionName   string
	Data           []byte
	BatchId        []byte
	SequenceNumber uint64
	TotalInBatch   uint64
}

func convertExecutingCallbackParamsFromBCS(bcs bcsExecutingCallbackParams) (ExecutingCallbackParams, error) {

	return ExecutingCallbackParams{
		Target:         fmt.Sprintf("0x%x", bcs.Target),
		ModuleName:     bcs.ModuleName,
		FunctionName:   bcs.FunctionName,
		Data:           bcs.Data,
		BatchId:        bcs.BatchId,
		SequenceNumber: bcs.SequenceNumber,
		TotalInBatch:   bcs.TotalInBatch,
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
	// Register vector decoder for Registry
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::Registry>", func(data []byte) (interface{}, error) {
		var results []Registry
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::PublisherWrapper", func(data []byte) (interface{}, error) {
		var result PublisherWrapper
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for PublisherWrapper
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::PublisherWrapper>", func(data []byte) (interface{}, error) {
		var results []PublisherWrapper
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::BatchExecutionState", func(data []byte) (interface{}, error) {
		var result BatchExecutionState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for BatchExecutionState
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::BatchExecutionState>", func(data []byte) (interface{}, error) {
		var results []BatchExecutionState
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
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
	// Register vector decoder for ExecutingCallbackParams
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::ExecutingCallbackParams>", func(data []byte) (interface{}, error) {
		var temps []bcsExecutingCallbackParams
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]ExecutingCallbackParams, len(temps))
		for i, temp := range temps {
			result, err := convertExecutingCallbackParamsFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::EntrypointRegistered", func(data []byte) (interface{}, error) {
		var result EntrypointRegistered
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for EntrypointRegistered
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::EntrypointRegistered>", func(data []byte) (interface{}, error) {
		var results []EntrypointRegistered
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::ModulesAdded", func(data []byte) (interface{}, error) {
		var result ModulesAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ModulesAdded
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::ModulesAdded>", func(data []byte) (interface{}, error) {
		var results []ModulesAdded
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::ModulesRemoved", func(data []byte) (interface{}, error) {
		var result ModulesRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ModulesRemoved
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::ModulesRemoved>", func(data []byte) (interface{}, error) {
		var results []ModulesRemoved
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_registry::MCMS_REGISTRY", func(data []byte) (interface{}, error) {
		var result MCMS_REGISTRY
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for MCMS_REGISTRY
	bind.RegisterStructDecoder("vector<mcms::mcms_registry::MCMS_REGISTRY>", func(data []byte) (interface{}, error) {
		var results []MCMS_REGISTRY
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// CreatePublisherWrapper executes the create_publisher_wrapper Move function.
func (c *McmsRegistryContract) CreatePublisherWrapper(ctx context.Context, opts *bind.CallOpts, typeArgs []string, publisher bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.CreatePublisherWrapper(typeArgs, publisher, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterEntrypoint executes the register_entrypoint Move function.
func (c *McmsRegistryContract) RegisterEntrypoint(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, publisherWrapper PublisherWrapper, proof bind.Object, packageCap bind.Object, allowedModules [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.RegisterEntrypoint(typeArgs, registry, publisherWrapper, proof, packageCap, allowedModules)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddAllowedModules executes the add_allowed_modules Move function.
func (c *McmsRegistryContract) AddAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, newAllowedModules [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.AddAllowedModules(typeArgs, registry, proof, newAllowedModules)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemoveAllowedModules executes the remove_allowed_modules Move function.
func (c *McmsRegistryContract) RemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, modulesToRemove [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.RemoveAllowedModules(typeArgs, registry, proof, modulesToRemove)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCallbackParamsWithCaps executes the get_callback_params_with_caps Move function.
func (c *McmsRegistryContract) GetCallbackParamsWithCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetCallbackParamsWithCaps(typeArgs, registry, proof, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReleaseCap executes the release_cap Move function.
func (c *McmsRegistryContract) ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.ReleaseCap(typeArgs, registry, proof)
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

// GetAcceptOwnershipData executes the get_accept_ownership_data Move function.
func (c *McmsRegistryContract) GetAcceptOwnershipData(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params ExecutingCallbackParams, proof bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetAcceptOwnershipData(typeArgs, registry, params, proof)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCallbackParamsFromMcms executes the get_callback_params_from_mcms Move function.
func (c *McmsRegistryContract) GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params ExecutingCallbackParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetCallbackParamsFromMcms(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateExecutingCallbackParams executes the create_executing_callback_params Move function.
func (c *McmsRegistryContract) CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte, batchId []byte, sequenceNumber uint64, totalInBatch uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.CreateExecutingCallbackParams(target, moduleName, functionName, data, batchId, sequenceNumber, totalInBatch)
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

// GetRegisteredProofType executes the get_registered_proof_type Move function.
func (c *McmsRegistryContract) GetRegisteredProofType(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetRegisteredProofType(registry, packageAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowedModules executes the get_allowed_modules Move function.
func (c *McmsRegistryContract) GetAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetAllowedModules(registry, packageAddress)
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

// IsBatchCompleted executes the is_batch_completed Move function.
func (c *McmsRegistryContract) IsBatchCompleted(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.IsBatchCompleted(registry, batchId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetNextExpectedSequence executes the get_next_expected_sequence Move function.
func (c *McmsRegistryContract) GetNextExpectedSequence(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetNextExpectedSequence(registry, batchId)
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

// GetMultisigAddressAscii executes the get_multisig_address_ascii Move function.
func (c *McmsRegistryContract) GetMultisigAddressAscii(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsRegistryEncoder.GetMultisigAddressAscii()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePublisherWrapper executes the create_publisher_wrapper Move function using DevInspect to get return values.
//
// Returns: PublisherWrapper<T>
func (d *McmsRegistryDevInspect) CreatePublisherWrapper(ctx context.Context, opts *bind.CallOpts, typeArgs []string, publisher bind.Object, proof bind.Object) (any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.CreatePublisherWrapper(typeArgs, publisher, proof)
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

// GetCallbackParamsWithCaps executes the get_callback_params_with_caps Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: &mut C
//	[1]: 0x1::string::String
//	[2]: vector<u8>
func (d *McmsRegistryDevInspect) GetCallbackParamsWithCaps(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) ([]any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetCallbackParamsWithCaps(typeArgs, registry, proof, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// ReleaseCap executes the release_cap Move function using DevInspect to get return values.
//
// Returns: C
func (d *McmsRegistryDevInspect) ReleaseCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, proof bind.Object) (any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.ReleaseCap(typeArgs, registry, proof)
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

// GetAcceptOwnershipData executes the get_accept_ownership_data Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *McmsRegistryDevInspect) GetAcceptOwnershipData(ctx context.Context, opts *bind.CallOpts, typeArgs []string, registry bind.Object, params ExecutingCallbackParams, proof bind.Object) ([]byte, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetAcceptOwnershipData(typeArgs, registry, params, proof)
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

// GetCallbackParamsFromMcms executes the get_callback_params_from_mcms Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: 0x1::string::String
//	[2]: 0x1::string::String
//	[3]: vector<u8>
func (d *McmsRegistryDevInspect) GetCallbackParamsFromMcms(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params ExecutingCallbackParams) ([]any, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetCallbackParamsFromMcms(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// CreateExecutingCallbackParams executes the create_executing_callback_params Move function using DevInspect to get return values.
//
// Returns: ExecutingCallbackParams
func (d *McmsRegistryDevInspect) CreateExecutingCallbackParams(ctx context.Context, opts *bind.CallOpts, target string, moduleName string, functionName string, data []byte, batchId []byte, sequenceNumber uint64, totalInBatch uint64) (ExecutingCallbackParams, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.CreateExecutingCallbackParams(target, moduleName, functionName, data, batchId, sequenceNumber, totalInBatch)
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

// GetRegisteredProofType executes the get_registered_proof_type Move function using DevInspect to get return values.
//
// Returns: TypeName
func (d *McmsRegistryDevInspect) GetRegisteredProofType(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) (bind.Object, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetRegisteredProofType(registry, packageAddress)
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

// GetAllowedModules executes the get_allowed_modules Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *McmsRegistryDevInspect) GetAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, packageAddress string) ([][]byte, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetAllowedModules(registry, packageAddress)
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

// IsBatchCompleted executes the is_batch_completed Move function using DevInspect to get return values.
//
// Returns: bool
func (d *McmsRegistryDevInspect) IsBatchCompleted(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (bool, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.IsBatchCompleted(registry, batchId)
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

// GetNextExpectedSequence executes the get_next_expected_sequence Move function using DevInspect to get return values.
//
// Returns: u64
func (d *McmsRegistryDevInspect) GetNextExpectedSequence(ctx context.Context, opts *bind.CallOpts, registry bind.Object, batchId []byte) (uint64, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetNextExpectedSequence(registry, batchId)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
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

// GetMultisigAddressAscii executes the get_multisig_address_ascii Move function using DevInspect to get return values.
//
// Returns: ascii::String
func (d *McmsRegistryDevInspect) GetMultisigAddressAscii(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.mcmsRegistryEncoder.GetMultisigAddressAscii()
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

type mcmsRegistryEncoder struct {
	*bind.BoundContract
}

// CreatePublisherWrapper encodes a call to the create_publisher_wrapper Move function.
func (c mcmsRegistryEncoder) CreatePublisherWrapper(typeArgs []string, publisher bind.Object, proof bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_publisher_wrapper", typeArgsList, typeParamsList, []string{
		"&Publisher",
		"T",
	}, []any{
		publisher,
		proof,
	}, []string{
		"PublisherWrapper<T>",
	})
}

// CreatePublisherWrapperWithArgs encodes a call to the create_publisher_wrapper Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) CreatePublisherWrapperWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Publisher",
		"T",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("create_publisher_wrapper", typeArgsList, typeParamsList, expectedParams, args, []string{
		"PublisherWrapper<T>",
	})
}

// RegisterEntrypoint encodes a call to the register_entrypoint Move function.
func (c mcmsRegistryEncoder) RegisterEntrypoint(typeArgs []string, registry bind.Object, publisherWrapper PublisherWrapper, proof bind.Object, packageCap bind.Object, allowedModules [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("register_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"PublisherWrapper<T>",
		"T",
		"C",
		"vector<vector<u8>>",
	}, []any{
		registry,
		publisherWrapper,
		proof,
		packageCap,
		allowedModules,
	}, nil)
}

// RegisterEntrypointWithArgs encodes a call to the register_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) RegisterEntrypointWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"PublisherWrapper<T>",
		"T",
		"C",
		"vector<vector<u8>>",
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

// AddAllowedModules encodes a call to the add_allowed_modules Move function.
func (c mcmsRegistryEncoder) AddAllowedModules(typeArgs []string, registry bind.Object, proof bind.Object, newAllowedModules [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"T",
		"vector<vector<u8>>",
	}, []any{
		registry,
		proof,
		newAllowedModules,
	}, nil)
}

// AddAllowedModulesWithArgs encodes a call to the add_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) AddAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"T",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RemoveAllowedModules encodes a call to the remove_allowed_modules Move function.
func (c mcmsRegistryEncoder) RemoveAllowedModules(typeArgs []string, registry bind.Object, proof bind.Object, modulesToRemove [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"T",
		"vector<vector<u8>>",
	}, []any{
		registry,
		proof,
		modulesToRemove,
	}, nil)
}

// RemoveAllowedModulesWithArgs encodes a call to the remove_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) RemoveAllowedModulesWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"T",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCallbackParamsWithCaps encodes a call to the get_callback_params_with_caps Move function.
func (c mcmsRegistryEncoder) GetCallbackParamsWithCaps(typeArgs []string, registry bind.Object, proof bind.Object, params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
		"C",
	}
	return c.EncodeCallArgsWithGenerics("get_callback_params_with_caps", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"T",
		"mcms::mcms_registry::ExecutingCallbackParams",
	}, []any{
		registry,
		proof,
		params,
	}, []string{
		"&mut C",
		"0x1::string::String",
		"vector<u8>",
	})
}

// GetCallbackParamsWithCapsWithArgs encodes a call to the get_callback_params_with_caps Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetCallbackParamsWithCapsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("get_callback_params_with_caps", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&mut C",
		"0x1::string::String",
		"vector<u8>",
	})
}

// ReleaseCap encodes a call to the release_cap Move function.
func (c mcmsRegistryEncoder) ReleaseCap(typeArgs []string, registry bind.Object, proof bind.Object) (*bind.EncodedCall, error) {
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
		proof,
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

// GetAcceptOwnershipData encodes a call to the get_accept_ownership_data Move function.
func (c mcmsRegistryEncoder) GetAcceptOwnershipData(typeArgs []string, registry bind.Object, params ExecutingCallbackParams, proof bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_accept_ownership_data", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"mcms::mcms_registry::ExecutingCallbackParams",
		"T",
	}, []any{
		registry,
		params,
		proof,
	}, []string{
		"vector<u8>",
	})
}

// GetAcceptOwnershipDataWithArgs encodes a call to the get_accept_ownership_data Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetAcceptOwnershipDataWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
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
	return c.EncodeCallArgsWithGenerics("get_accept_ownership_data", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// GetCallbackParamsFromMcms encodes a call to the get_callback_params_from_mcms Move function.
func (c mcmsRegistryEncoder) GetCallbackParamsFromMcms(registry bind.Object, params ExecutingCallbackParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_callback_params_from_mcms", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"mcms::mcms_registry::ExecutingCallbackParams",
	}, []any{
		registry,
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
		"&mut Registry",
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
func (c mcmsRegistryEncoder) CreateExecutingCallbackParams(target string, moduleName string, functionName string, data []byte, batchId []byte, sequenceNumber uint64, totalInBatch uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_executing_callback_params", typeArgsList, typeParamsList, []string{
		"address",
		"0x1::string::String",
		"0x1::string::String",
		"vector<u8>",
		"vector<u8>",
		"u64",
		"u64",
	}, []any{
		target,
		moduleName,
		functionName,
		data,
		batchId,
		sequenceNumber,
		totalInBatch,
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
		"vector<u8>",
		"u64",
		"u64",
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
		"ascii::String",
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
		"ascii::String",
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

// GetRegisteredProofType encodes a call to the get_registered_proof_type Move function.
func (c mcmsRegistryEncoder) GetRegisteredProofType(registry bind.Object, packageAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_registered_proof_type", typeArgsList, typeParamsList, []string{
		"&Registry",
		"ascii::String",
	}, []any{
		registry,
		packageAddress,
	}, []string{
		"TypeName",
	})
}

// GetRegisteredProofTypeWithArgs encodes a call to the get_registered_proof_type Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetRegisteredProofTypeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Registry",
		"ascii::String",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_registered_proof_type", typeArgsList, typeParamsList, expectedParams, args, []string{
		"TypeName",
	})
}

// GetAllowedModules encodes a call to the get_allowed_modules Move function.
func (c mcmsRegistryEncoder) GetAllowedModules(registry bind.Object, packageAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_modules", typeArgsList, typeParamsList, []string{
		"&Registry",
		"ascii::String",
	}, []any{
		registry,
		packageAddress,
	}, []string{
		"vector<vector<u8>>",
	})
}

// GetAllowedModulesWithArgs encodes a call to the get_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Registry",
		"ascii::String",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
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

// IsBatchCompleted encodes a call to the is_batch_completed Move function.
func (c mcmsRegistryEncoder) IsBatchCompleted(registry bind.Object, batchId []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_batch_completed", typeArgsList, typeParamsList, []string{
		"&Registry",
		"vector<u8>",
	}, []any{
		registry,
		batchId,
	}, []string{
		"bool",
	})
}

// IsBatchCompletedWithArgs encodes a call to the is_batch_completed Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) IsBatchCompletedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Registry",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_batch_completed", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetNextExpectedSequence encodes a call to the get_next_expected_sequence Move function.
func (c mcmsRegistryEncoder) GetNextExpectedSequence(registry bind.Object, batchId []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_next_expected_sequence", typeArgsList, typeParamsList, []string{
		"&Registry",
		"vector<u8>",
	}, []any{
		registry,
		batchId,
	}, []string{
		"u64",
	})
}

// GetNextExpectedSequenceWithArgs encodes a call to the get_next_expected_sequence Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetNextExpectedSequenceWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&Registry",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_next_expected_sequence", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
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

// GetMultisigAddressAscii encodes a call to the get_multisig_address_ascii Move function.
func (c mcmsRegistryEncoder) GetMultisigAddressAscii() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_multisig_address_ascii", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"ascii::String",
	})
}

// GetMultisigAddressAsciiWithArgs encodes a call to the get_multisig_address_ascii Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsRegistryEncoder) GetMultisigAddressAsciiWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_multisig_address_ascii", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ascii::String",
	})
}
