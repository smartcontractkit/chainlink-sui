// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_nonce_manager

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

type INonceManager interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error)
	GetIncrementedOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error)
	DevInspect() INonceManagerDevInspect
	Encoder() NonceManagerEncoder
}

type INonceManagerDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error)
	GetIncrementedOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error)
}

type NonceManagerEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error)
	GetOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error)
	GetIncrementedOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error)
	GetIncrementedOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error)
}

type NonceManagerContract struct {
	*bind.BoundContract
	nonceManagerEncoder
	devInspect *NonceManagerDevInspect
}

type NonceManagerDevInspect struct {
	contract *NonceManagerContract
}

var _ INonceManager = (*NonceManagerContract)(nil)
var _ INonceManagerDevInspect = (*NonceManagerDevInspect)(nil)

func NewNonceManager(packageID string, client sui.ISuiAPI) (*NonceManagerContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "nonce_manager", client)
	if err != nil {
		return nil, err
	}

	c := &NonceManagerContract{
		BoundContract:       contract,
		nonceManagerEncoder: nonceManagerEncoder{BoundContract: contract},
	}
	c.devInspect = &NonceManagerDevInspect{contract: c}
	return c, nil
}

func (c *NonceManagerContract) Encoder() NonceManagerEncoder {
	return c.nonceManagerEncoder
}

func (c *NonceManagerContract) DevInspect() INonceManagerDevInspect {
	return c.devInspect
}

type NonceManagerState struct {
	Id             string      `move:"sui::object::UID"`
	OutboundNonces bind.Object `move:"Table<u64, Table<address, u64>>"`
}

func init() {
	bind.RegisterStructDecoder("ccip::nonce_manager::NonceManagerState", func(data []byte) (interface{}, error) {
		var result NonceManagerState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *NonceManagerContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.nonceManagerEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *NonceManagerContract) Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.nonceManagerEncoder.Initialize(ref, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOutboundNonce executes the get_outbound_nonce Move function.
func (c *NonceManagerContract) GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.nonceManagerEncoder.GetOutboundNonce(ref, destChainSelector, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetIncrementedOutboundNonce executes the get_incremented_outbound_nonce Move function.
func (c *NonceManagerContract) GetIncrementedOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.nonceManagerEncoder.GetIncrementedOutboundNonce(ref, destChainSelector, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *NonceManagerDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.nonceManagerEncoder.TypeAndVersion()
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

// GetOutboundNonce executes the get_outbound_nonce Move function using DevInspect to get return values.
//
// Returns: u64
func (d *NonceManagerDevInspect) GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error) {
	encoded, err := d.contract.nonceManagerEncoder.GetOutboundNonce(ref, destChainSelector, sender)
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

// GetIncrementedOutboundNonce executes the get_incremented_outbound_nonce Move function using DevInspect to get return values.
//
// Returns: u64
func (d *NonceManagerDevInspect) GetIncrementedOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error) {
	encoded, err := d.contract.nonceManagerEncoder.GetIncrementedOutboundNonce(ref, destChainSelector, sender)
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

type nonceManagerEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c nonceManagerEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c nonceManagerEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c nonceManagerEncoder) Initialize(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
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
func (c nonceManagerEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// GetOutboundNonce encodes a call to the get_outbound_nonce Move function.
func (c nonceManagerEncoder) GetOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_outbound_nonce", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"address",
	}, []any{
		ref,
		destChainSelector,
		sender,
	}, []string{
		"u64",
	})
}

// GetOutboundNonceWithArgs encodes a call to the get_outbound_nonce Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c nonceManagerEncoder) GetOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_outbound_nonce", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetIncrementedOutboundNonce encodes a call to the get_incremented_outbound_nonce Move function.
func (c nonceManagerEncoder) GetIncrementedOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_incremented_outbound_nonce", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"u64",
		"address",
	}, []any{
		ref,
		destChainSelector,
		sender,
	}, []string{
		"u64",
	})
}

// GetIncrementedOutboundNonceWithArgs encodes a call to the get_incremented_outbound_nonce Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c nonceManagerEncoder) GetIncrementedOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_incremented_outbound_nonce", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}
