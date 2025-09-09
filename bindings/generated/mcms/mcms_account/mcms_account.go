// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_mcms_account

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

type IMcmsAccount interface {
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	TransferOwnershipToSelf(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipAsTimelock(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMcmsAccountDevInspect
	Encoder() McmsAccountEncoder
	Bound() bind.IBoundContract
}

type IMcmsAccountDevInspect interface {
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type McmsAccountEncoder interface {
	TransferOwnership(param bind.Object, state bind.Object, to string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	TransferOwnershipToSelf(ownerCap bind.Object, state bind.Object) (*bind.EncodedCall, error)
	TransferOwnershipToSelfWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(state bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipAsTimelock(state bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipAsTimelockWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error)
	AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error)
}

type McmsAccountContract struct {
	*bind.BoundContract
	mcmsAccountEncoder
	devInspect *McmsAccountDevInspect
}

type McmsAccountDevInspect struct {
	contract *McmsAccountContract
}

var _ IMcmsAccount = (*McmsAccountContract)(nil)
var _ IMcmsAccountDevInspect = (*McmsAccountDevInspect)(nil)

func NewMcmsAccount(packageID string, client sui.ISuiAPI) (IMcmsAccount, error) {
	contract, err := bind.NewBoundContract(packageID, "mcms", "mcms_account", client)
	if err != nil {
		return nil, err
	}

	c := &McmsAccountContract{
		BoundContract:      contract,
		mcmsAccountEncoder: mcmsAccountEncoder{BoundContract: contract},
	}
	c.devInspect = &McmsAccountDevInspect{contract: c}
	return c, nil
}

func (c *McmsAccountContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *McmsAccountContract) Encoder() McmsAccountEncoder {
	return c.mcmsAccountEncoder
}

func (c *McmsAccountContract) DevInspect() IMcmsAccountDevInspect {
	return c.devInspect
}

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type AccountState struct {
	Id              string           `move:"sui::object::UID"`
	Owner           string           `move:"address"`
	PendingTransfer *PendingTransfer `move:"0x1::option::Option<PendingTransfer>"`
}

type PendingTransfer struct {
	From     string `move:"address"`
	To       string `move:"address"`
	Accepted bool   `move:"bool"`
}

type OwnershipTransferRequested struct {
	From string `move:"address"`
	To   string `move:"address"`
}

type OwnershipTransferAccepted struct {
	From string `move:"address"`
	To   string `move:"address"`
}

type OwnershipTransferred struct {
	From string `move:"address"`
	To   string `move:"address"`
}

type MCMS_ACCOUNT struct {
}

type bcsAccountState struct {
	Id              string
	Owner           [32]byte
	PendingTransfer *PendingTransfer
}

func convertAccountStateFromBCS(bcs bcsAccountState) (AccountState, error) {

	return AccountState{
		Id:              bcs.Id,
		Owner:           fmt.Sprintf("0x%x", bcs.Owner),
		PendingTransfer: bcs.PendingTransfer,
	}, nil
}

type bcsPendingTransfer struct {
	From     [32]byte
	To       [32]byte
	Accepted bool
}

func convertPendingTransferFromBCS(bcs bcsPendingTransfer) (PendingTransfer, error) {

	return PendingTransfer{
		From:     fmt.Sprintf("0x%x", bcs.From),
		To:       fmt.Sprintf("0x%x", bcs.To),
		Accepted: bcs.Accepted,
	}, nil
}

type bcsOwnershipTransferRequested struct {
	From [32]byte
	To   [32]byte
}

func convertOwnershipTransferRequestedFromBCS(bcs bcsOwnershipTransferRequested) (OwnershipTransferRequested, error) {

	return OwnershipTransferRequested{
		From: fmt.Sprintf("0x%x", bcs.From),
		To:   fmt.Sprintf("0x%x", bcs.To),
	}, nil
}

type bcsOwnershipTransferAccepted struct {
	From [32]byte
	To   [32]byte
}

func convertOwnershipTransferAcceptedFromBCS(bcs bcsOwnershipTransferAccepted) (OwnershipTransferAccepted, error) {

	return OwnershipTransferAccepted{
		From: fmt.Sprintf("0x%x", bcs.From),
		To:   fmt.Sprintf("0x%x", bcs.To),
	}, nil
}

type bcsOwnershipTransferred struct {
	From [32]byte
	To   [32]byte
}

func convertOwnershipTransferredFromBCS(bcs bcsOwnershipTransferred) (OwnershipTransferred, error) {

	return OwnershipTransferred{
		From: fmt.Sprintf("0x%x", bcs.From),
		To:   fmt.Sprintf("0x%x", bcs.To),
	}, nil
}

func init() {
	bind.RegisterStructDecoder("mcms::mcms_account::OwnerCap", func(data []byte) (interface{}, error) {
		var result OwnerCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_account::AccountState", func(data []byte) (interface{}, error) {
		var temp bcsAccountState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAccountStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_account::PendingTransfer", func(data []byte) (interface{}, error) {
		var temp bcsPendingTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertPendingTransferFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_account::OwnershipTransferRequested", func(data []byte) (interface{}, error) {
		var temp bcsOwnershipTransferRequested
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOwnershipTransferRequestedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_account::OwnershipTransferAccepted", func(data []byte) (interface{}, error) {
		var temp bcsOwnershipTransferAccepted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOwnershipTransferAcceptedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_account::OwnershipTransferred", func(data []byte) (interface{}, error) {
		var temp bcsOwnershipTransferred
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOwnershipTransferredFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("mcms::mcms_account::MCMS_ACCOUNT", func(data []byte) (interface{}, error) {
		var result MCMS_ACCOUNT
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *McmsAccountContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, param bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.TransferOwnership(param, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnershipToSelf executes the transfer_ownership_to_self Move function.
func (c *McmsAccountContract) TransferOwnershipToSelf(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.TransferOwnershipToSelf(ownerCap, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *McmsAccountContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.AcceptOwnership(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipAsTimelock executes the accept_ownership_as_timelock Move function.
func (c *McmsAccountContract) AcceptOwnershipAsTimelock(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.AcceptOwnershipAsTimelock(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *McmsAccountContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.AcceptOwnershipFromObject(state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *McmsAccountContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.ExecuteOwnershipTransfer(ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *McmsAccountContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *McmsAccountContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *McmsAccountContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.mcmsAccountEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *McmsAccountDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.mcmsAccountEncoder.PendingTransferFrom(state)
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
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferTo executes the pending_transfer_to Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *McmsAccountDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.mcmsAccountEncoder.PendingTransferTo(state)
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
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<bool>
func (d *McmsAccountDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error) {
	encoded, err := d.contract.mcmsAccountEncoder.PendingTransferAccepted(state)
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
	result, ok := results[0].(*bool)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *bool, got %T", results[0])
	}
	return result, nil
}

type mcmsAccountEncoder struct {
	*bind.BoundContract
}

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c mcmsAccountEncoder) TransferOwnership(param bind.Object, state bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut AccountState",
		"address",
	}, []any{
		param,
		state,
		to,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut AccountState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// TransferOwnershipToSelf encodes a call to the transfer_ownership_to_self Move function.
func (c mcmsAccountEncoder) TransferOwnershipToSelf(ownerCap bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership_to_self", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut AccountState",
	}, []any{
		ownerCap,
		state,
	}, nil)
}

// TransferOwnershipToSelfWithArgs encodes a call to the transfer_ownership_to_self Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) TransferOwnershipToSelfWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut AccountState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership_to_self", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnership encodes a call to the accept_ownership Move function.
func (c mcmsAccountEncoder) AcceptOwnership(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut AccountState",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut AccountState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipAsTimelock encodes a call to the accept_ownership_as_timelock Move function.
func (c mcmsAccountEncoder) AcceptOwnershipAsTimelock(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_timelock", typeArgsList, typeParamsList, []string{
		"&mut AccountState",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipAsTimelockWithArgs encodes a call to the accept_ownership_as_timelock Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) AcceptOwnershipAsTimelockWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut AccountState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_timelock", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c mcmsAccountEncoder) AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut AccountState",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut AccountState",
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c mcmsAccountEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"mcms::mcms_account::OwnerCap",
		"&mut AccountState",
		"&mut Registry",
		"address",
	}, []any{
		ownerCap,
		state,
		registry,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"mcms::mcms_account::OwnerCap",
		"&mut AccountState",
		"&mut Registry",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// PendingTransferFrom encodes a call to the pending_transfer_from Move function.
func (c mcmsAccountEncoder) PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&AccountState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&AccountState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferTo encodes a call to the pending_transfer_to Move function.
func (c mcmsAccountEncoder) PendingTransferTo(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&AccountState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&AccountState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferAccepted encodes a call to the pending_transfer_accepted Move function.
func (c mcmsAccountEncoder) PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&AccountState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c mcmsAccountEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&AccountState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<bool>",
	})
}
