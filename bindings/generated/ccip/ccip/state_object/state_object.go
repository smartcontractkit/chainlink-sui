// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_state_object

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

type IStateObject interface {
	OwnerCapId(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	Add(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object, obj bind.Object) (*models.SuiTransactionBlockResponse, error)
	Contains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	Remove(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	Borrow(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	BorrowMut(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsProofEntrypoint(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IStateObjectDevInspect
	Encoder() StateObjectEncoder
}

type IStateObjectDevInspect interface {
	OwnerCapId(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (bind.Object, error)
	Contains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (bool, error)
	Remove(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object) (any, error)
	Borrow(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (bind.Object, error)
	BorrowMut(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (bind.Object, error)
	Owner(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*bool, error)
	McmsProofEntrypoint(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (CCIPAdminProof, error)
}

type StateObjectEncoder interface {
	OwnerCapId(ref bind.Object) (*bind.EncodedCall, error)
	OwnerCapIdWithArgs(args ...any) (*bind.EncodedCall, error)
	Add(typeArgs []string, ref bind.Object, ownerCap bind.Object, obj bind.Object) (*bind.EncodedCall, error)
	AddWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Contains(typeArgs []string, ref bind.Object) (*bind.EncodedCall, error)
	ContainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Remove(typeArgs []string, ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	RemoveWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	Borrow(typeArgs []string, ref bind.Object) (*bind.EncodedCall, error)
	BorrowWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	BorrowMut(typeArgs []string, ref bind.Object) (*bind.EncodedCall, error)
	BorrowMutWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	TransferOwnership(ref bind.Object, ownerCap bind.Object, to string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(ref bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ref bind.Object, ownerCap bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(ref bind.Object, ownerCap bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	Owner(ref bind.Object) (*bind.EncodedCall, error)
	OwnerWithArgs(args ...any) (*bind.EncodedCall, error)
	HasPendingTransfer(ref bind.Object) (*bind.EncodedCall, error)
	HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(ref bind.Object) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(ref bind.Object) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(ref bind.Object) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsProofEntrypoint(registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsProofEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
}

type StateObjectContract struct {
	*bind.BoundContract
	stateObjectEncoder
	devInspect *StateObjectDevInspect
}

type StateObjectDevInspect struct {
	contract *StateObjectContract
}

var _ IStateObject = (*StateObjectContract)(nil)
var _ IStateObjectDevInspect = (*StateObjectDevInspect)(nil)

func NewStateObject(packageID string, client sui.ISuiAPI) (*StateObjectContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "state_object", client, nil)
	if err != nil {
		return nil, err
	}

	c := &StateObjectContract{
		BoundContract:      contract,
		stateObjectEncoder: stateObjectEncoder{BoundContract: contract},
	}
	c.devInspect = &StateObjectDevInspect{contract: c}
	return c, nil
}

func (c *StateObjectContract) Encoder() StateObjectEncoder {
	return c.stateObjectEncoder
}

func (c *StateObjectContract) DevInspect() IStateObjectDevInspect {
	return c.devInspect
}

type CCIPObjectRef struct {
	Id           string      `move:"sui::object::UID"`
	OwnableState bind.Object `move:"OwnableState"`
}

type CCIPObjectRefPointer struct {
	Id          string `move:"sui::object::UID"`
	ObjectRefId string `move:"address"`
	OwnerCapId  string `move:"address"`
}

type STATE_OBJECT struct {
}

type CCIPAdminProof struct {
}

type McmsCallback struct {
}

type bcsCCIPObjectRefPointer struct {
	Id          string
	ObjectRefId [32]byte
	OwnerCapId  [32]byte
}

func convertCCIPObjectRefPointerFromBCS(bcs bcsCCIPObjectRefPointer) (CCIPObjectRefPointer, error) {

	return CCIPObjectRefPointer{
		Id:          bcs.Id,
		ObjectRefId: fmt.Sprintf("0x%x", bcs.ObjectRefId),
		OwnerCapId:  fmt.Sprintf("0x%x", bcs.OwnerCapId),
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip::state_object::CCIPObjectRef", func(data []byte) (interface{}, error) {
		var result CCIPObjectRef
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::state_object::CCIPObjectRefPointer", func(data []byte) (interface{}, error) {
		var temp bcsCCIPObjectRefPointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCCIPObjectRefPointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::state_object::STATE_OBJECT", func(data []byte) (interface{}, error) {
		var result STATE_OBJECT
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::state_object::CCIPAdminProof", func(data []byte) (interface{}, error) {
		var result CCIPAdminProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::state_object::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// OwnerCapId executes the owner_cap_id Move function.
func (c *StateObjectContract) OwnerCapId(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.OwnerCapId(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Add executes the add Move function.
func (c *StateObjectContract) Add(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object, obj bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.Add(typeArgs, ref, ownerCap, obj)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Contains executes the contains Move function.
func (c *StateObjectContract) Contains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.Contains(typeArgs, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Remove executes the remove Move function.
func (c *StateObjectContract) Remove(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.Remove(typeArgs, ref, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Borrow executes the borrow Move function.
func (c *StateObjectContract) Borrow(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.Borrow(typeArgs, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// BorrowMut executes the borrow_mut Move function.
func (c *StateObjectContract) BorrowMut(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.BorrowMut(typeArgs, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *StateObjectContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.TransferOwnership(ref, ownerCap, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *StateObjectContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.AcceptOwnership(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *StateObjectContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.ExecuteOwnershipTransfer(ref, ownerCap, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *StateObjectContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.ExecuteOwnershipTransferToMcms(ref, ownerCap, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *StateObjectContract) Owner(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.Owner(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *StateObjectContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.HasPendingTransfer(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *StateObjectContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.PendingTransferFrom(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *StateObjectContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.PendingTransferTo(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *StateObjectContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.PendingTransferAccepted(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *StateObjectContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.McmsEntrypoint(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsProofEntrypoint executes the mcms_proof_entrypoint Move function.
func (c *StateObjectContract) McmsProofEntrypoint(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.stateObjectEncoder.McmsProofEntrypoint(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// OwnerCapId executes the owner_cap_id Move function using DevInspect to get return values.
//
// Returns: ID
func (d *StateObjectDevInspect) OwnerCapId(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (bind.Object, error) {
	encoded, err := d.contract.stateObjectEncoder.OwnerCapId(ref)
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

// Contains executes the contains Move function using DevInspect to get return values.
//
// Returns: bool
func (d *StateObjectDevInspect) Contains(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (bool, error) {
	encoded, err := d.contract.stateObjectEncoder.Contains(typeArgs, ref)
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

// Remove executes the remove Move function using DevInspect to get return values.
//
// Returns: T
func (d *StateObjectDevInspect) Remove(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, ownerCap bind.Object) (any, error) {
	encoded, err := d.contract.stateObjectEncoder.Remove(typeArgs, ref, ownerCap)
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

// Borrow executes the borrow Move function using DevInspect to get return values.
//
// Returns: &T
func (d *StateObjectDevInspect) Borrow(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (bind.Object, error) {
	encoded, err := d.contract.stateObjectEncoder.Borrow(typeArgs, ref)
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

// BorrowMut executes the borrow_mut Move function using DevInspect to get return values.
//
// Returns: &mut T
func (d *StateObjectDevInspect) BorrowMut(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object) (bind.Object, error) {
	encoded, err := d.contract.stateObjectEncoder.BorrowMut(typeArgs, ref)
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

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *StateObjectDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (string, error) {
	encoded, err := d.contract.stateObjectEncoder.Owner(ref)
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

// HasPendingTransfer executes the has_pending_transfer Move function using DevInspect to get return values.
//
// Returns: bool
func (d *StateObjectDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (bool, error) {
	encoded, err := d.contract.stateObjectEncoder.HasPendingTransfer(ref)
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

// PendingTransferFrom executes the pending_transfer_from Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *StateObjectDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*string, error) {
	encoded, err := d.contract.stateObjectEncoder.PendingTransferFrom(ref)
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
func (d *StateObjectDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*string, error) {
	encoded, err := d.contract.stateObjectEncoder.PendingTransferTo(ref)
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
func (d *StateObjectDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*bool, error) {
	encoded, err := d.contract.stateObjectEncoder.PendingTransferAccepted(ref)
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

// McmsProofEntrypoint executes the mcms_proof_entrypoint Move function using DevInspect to get return values.
//
// Returns: CCIPAdminProof
func (d *StateObjectDevInspect) McmsProofEntrypoint(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (CCIPAdminProof, error) {
	encoded, err := d.contract.stateObjectEncoder.McmsProofEntrypoint(registry, params)
	if err != nil {
		return CCIPAdminProof{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return CCIPAdminProof{}, err
	}
	if len(results) == 0 {
		return CCIPAdminProof{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(CCIPAdminProof)
	if !ok {
		return CCIPAdminProof{}, fmt.Errorf("unexpected return type: expected CCIPAdminProof, got %T", results[0])
	}
	return result, nil
}

type stateObjectEncoder struct {
	*bind.BoundContract
}

// OwnerCapId encodes a call to the owner_cap_id Move function.
func (c stateObjectEncoder) OwnerCapId(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner_cap_id", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"ID",
	})
}

// OwnerCapIdWithArgs encodes a call to the owner_cap_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) OwnerCapIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner_cap_id", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ID",
	})
}

// Add encodes a call to the add Move function.
func (c stateObjectEncoder) Add(typeArgs []string, ref bind.Object, ownerCap bind.Object, obj bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"T",
	}, []any{
		ref,
		ownerCap,
		obj,
	}, nil)
}

// AddWithArgs encodes a call to the add Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) AddWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"T",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("add", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Contains encodes a call to the contains Move function.
func (c stateObjectEncoder) Contains(typeArgs []string, ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("contains", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"bool",
	})
}

// ContainsWithArgs encodes a call to the contains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) ContainsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("contains", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// Remove encodes a call to the remove Move function.
func (c stateObjectEncoder) Remove(typeArgs []string, ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}, []any{
		ref,
		ownerCap,
	}, []string{
		"T",
	})
}

// RemoveWithArgs encodes a call to the remove Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) RemoveWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("remove", typeArgsList, typeParamsList, expectedParams, args, []string{
		"T",
	})
}

// Borrow encodes a call to the borrow Move function.
func (c stateObjectEncoder) Borrow(typeArgs []string, ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("borrow", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"&T",
	})
}

// BorrowWithArgs encodes a call to the borrow Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) BorrowWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("borrow", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&T",
	})
}

// BorrowMut encodes a call to the borrow_mut Move function.
func (c stateObjectEncoder) BorrowMut(typeArgs []string, ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("borrow_mut", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"&mut T",
	})
}

// BorrowMutWithArgs encodes a call to the borrow_mut Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) BorrowMutWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("borrow_mut", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&mut T",
	})
}

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c stateObjectEncoder) TransferOwnership(ref bind.Object, ownerCap bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"address",
	}, []any{
		ref,
		ownerCap,
		to,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnership encodes a call to the accept_ownership Move function.
func (c stateObjectEncoder) AcceptOwnership(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
	}, []any{
		ref,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c stateObjectEncoder) ExecuteOwnershipTransfer(ref bind.Object, ownerCap bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"OwnerCap",
		"address",
	}, []any{
		ref,
		ownerCap,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransferToMcms encodes a call to the execute_ownership_transfer_to_mcms Move function.
func (c stateObjectEncoder) ExecuteOwnershipTransferToMcms(ref bind.Object, ownerCap bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"OwnerCap",
		"&mut Registry",
		"address",
	}, []any{
		ref,
		ownerCap,
		registry,
		to,
	}, nil)
}

// ExecuteOwnershipTransferToMcmsWithArgs encodes a call to the execute_ownership_transfer_to_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"OwnerCap",
		"&mut Registry",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Owner encodes a call to the owner Move function.
func (c stateObjectEncoder) Owner(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// HasPendingTransfer encodes a call to the has_pending_transfer Move function.
func (c stateObjectEncoder) HasPendingTransfer(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// PendingTransferFrom encodes a call to the pending_transfer_from Move function.
func (c stateObjectEncoder) PendingTransferFrom(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
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
func (c stateObjectEncoder) PendingTransferTo(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
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
func (c stateObjectEncoder) PendingTransferAccepted(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
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

// McmsEntrypoint encodes a call to the mcms_entrypoint Move function.
func (c stateObjectEncoder) McmsEntrypoint(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsEntrypointWithArgs encodes a call to the mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsProofEntrypoint encodes a call to the mcms_proof_entrypoint Move function.
func (c stateObjectEncoder) McmsProofEntrypoint(registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_proof_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		params,
	}, []string{
		"ccip::state_object::CCIPAdminProof",
	})
}

// McmsProofEntrypointWithArgs encodes a call to the mcms_proof_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c stateObjectEncoder) McmsProofEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_proof_entrypoint", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::state_object::CCIPAdminProof",
	})
}
