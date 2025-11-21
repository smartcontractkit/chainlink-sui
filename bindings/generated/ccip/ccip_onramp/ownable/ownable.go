// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_ownable

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

const FunctionInfo = `[{"package":"ccip_onramp","module":"ownable","name":"accept_ownership","parameters":[{"name":"state","type":"OwnableState"}]},{"package":"ccip_onramp","module":"ownable","name":"accept_ownership_from_object","parameters":[{"name":"state","type":"OwnableState"},{"name":"from","type":"sui::object::UID"}]},{"package":"ccip_onramp","module":"ownable","name":"attach_publisher","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"publisher","type":"Publisher"}]},{"package":"ccip_onramp","module":"ownable","name":"borrow_publisher","parameters":[{"name":"owner_cap","type":"OwnerCap"}]},{"package":"ccip_onramp","module":"ownable","name":"default_key","parameters":null},{"package":"ccip_onramp","module":"ownable","name":"execute_ownership_transfer","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"OwnableState"},{"name":"to","type":"address"}]},{"package":"ccip_onramp","module":"ownable","name":"execute_ownership_transfer_to_mcms","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"OwnableState"},{"name":"registry","type":"Registry"},{"name":"to","type":"address"},{"name":"publisher_wrapper","type":"PublisherWrapper<T>"},{"name":"proof","type":"T"},{"name":"allowed_modules","type":"vector<vector<u8>>"}]},{"package":"ccip_onramp","module":"ownable","name":"has_pending_transfer","parameters":[{"name":"state","type":"OwnableState"}]},{"package":"ccip_onramp","module":"ownable","name":"new","parameters":[{"name":"uid","type":"sui::object::UID"}]},{"package":"ccip_onramp","module":"ownable","name":"new_with_key","parameters":[{"name":"uid","type":"sui::object::UID"},{"name":"key","type":"K"}]},{"package":"ccip_onramp","module":"ownable","name":"owner","parameters":[{"name":"state","type":"OwnableState"}]},{"package":"ccip_onramp","module":"ownable","name":"owner_cap_id","parameters":[{"name":"state","type":"OwnableState"}]},{"package":"ccip_onramp","module":"ownable","name":"pending_transfer_accepted","parameters":[{"name":"state","type":"OwnableState"}]},{"package":"ccip_onramp","module":"ownable","name":"pending_transfer_from","parameters":[{"name":"state","type":"OwnableState"}]},{"package":"ccip_onramp","module":"ownable","name":"pending_transfer_to","parameters":[{"name":"state","type":"OwnableState"}]},{"package":"ccip_onramp","module":"ownable","name":"transfer_ownership","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"OwnableState"},{"name":"to","type":"address"}]}]`

type IOwnable interface {
	DefaultKey(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	New(ctx context.Context, opts *bind.CallOpts, uid string) (*models.SuiTransactionBlockResponse, error)
	NewWithKey(ctx context.Context, opts *bind.CallOpts, typeArgs []string, uid string, key bind.Object) (*models.SuiTransactionBlockResponse, error)
	OwnerCapId(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error)
	AttachPublisher(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, publisher bind.Object) (*models.SuiTransactionBlockResponse, error)
	BorrowPublisher(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state OwnableState, to string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state OwnableState, from string) (*models.SuiTransactionBlockResponse, error)
	McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, state OwnableState, mcms string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state OwnableState, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state OwnableState, registry bind.Object, to string, publisherWrapper bind.Object, proof bind.Object, allowedModules [][]byte) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOwnableDevInspect
	Encoder() OwnableEncoder
	Bound() bind.IBoundContract
}

type IOwnableDevInspect interface {
	DefaultKey(ctx context.Context, opts *bind.CallOpts) ([]byte, error)
	New(ctx context.Context, opts *bind.CallOpts, uid string) ([]any, error)
	NewWithKey(ctx context.Context, opts *bind.CallOpts, typeArgs []string, uid string, key bind.Object) ([]any, error)
	OwnerCapId(ctx context.Context, opts *bind.CallOpts, state OwnableState) (bind.Object, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state OwnableState) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state OwnableState) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*bool, error)
	BorrowPublisher(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object) (bind.Object, error)
}

type OwnableEncoder interface {
	DefaultKey() (*bind.EncodedCall, error)
	DefaultKeyWithArgs(args ...any) (*bind.EncodedCall, error)
	New(uid string) (*bind.EncodedCall, error)
	NewWithArgs(args ...any) (*bind.EncodedCall, error)
	NewWithKey(typeArgs []string, uid string, key bind.Object) (*bind.EncodedCall, error)
	NewWithKeyWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	OwnerCapId(state OwnableState) (*bind.EncodedCall, error)
	OwnerCapIdWithArgs(args ...any) (*bind.EncodedCall, error)
	Owner(state OwnableState) (*bind.EncodedCall, error)
	OwnerWithArgs(args ...any) (*bind.EncodedCall, error)
	HasPendingTransfer(state OwnableState) (*bind.EncodedCall, error)
	HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(state OwnableState) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(state OwnableState) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(state OwnableState) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error)
	AttachPublisher(ownerCap bind.Object, publisher bind.Object) (*bind.EncodedCall, error)
	AttachPublisherWithArgs(args ...any) (*bind.EncodedCall, error)
	BorrowPublisher(ownerCap bind.Object) (*bind.EncodedCall, error)
	BorrowPublisherWithArgs(args ...any) (*bind.EncodedCall, error)
	TransferOwnership(ownerCap bind.Object, state OwnableState, to string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(state OwnableState) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipFromObject(state OwnableState, from string) (*bind.EncodedCall, error)
	AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsAcceptOwnership(state OwnableState, mcms string) (*bind.EncodedCall, error)
	McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ownerCap bind.Object, state OwnableState, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(typeArgs []string, ownerCap bind.Object, state OwnableState, registry bind.Object, to string, publisherWrapper bind.Object, proof bind.Object, allowedModules [][]byte) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
}

type OwnableContract struct {
	*bind.BoundContract
	ownableEncoder
	devInspect *OwnableDevInspect
}

type OwnableDevInspect struct {
	contract *OwnableContract
}

var _ IOwnable = (*OwnableContract)(nil)
var _ IOwnableDevInspect = (*OwnableDevInspect)(nil)

func NewOwnable(packageID string, client sui.ISuiAPI) (IOwnable, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_onramp", "ownable", client)
	if err != nil {
		return nil, err
	}

	c := &OwnableContract{
		BoundContract:  contract,
		ownableEncoder: ownableEncoder{BoundContract: contract},
	}
	c.devInspect = &OwnableDevInspect{contract: c}
	return c, nil
}

func (c *OwnableContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *OwnableContract) Encoder() OwnableEncoder {
	return c.ownableEncoder
}

func (c *OwnableContract) DevInspect() IOwnableDevInspect {
	return c.devInspect
}

type OwnerCap struct {
	Id string `move:"sui::object::UID"`
}

type OwnableState struct {
	Owner           string           `move:"address"`
	PendingTransfer *PendingTransfer `move:"0x1::option::Option<PendingTransfer>"`
	OwnerCapId      bind.Object      `move:"ID"`
}

type PendingTransfer struct {
	From     string `move:"address"`
	To       string `move:"address"`
	Accepted bool   `move:"bool"`
}

type PublisherKey struct {
}

type NewOwnableStateEvent struct {
	OwnerCapId bind.Object `move:"ID"`
	Owner      string      `move:"address"`
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

type bcsOwnableState struct {
	Owner           [32]byte
	PendingTransfer *PendingTransfer
	OwnerCapId      bind.Object
}

func convertOwnableStateFromBCS(bcs bcsOwnableState) (OwnableState, error) {

	return OwnableState{
		Owner:           fmt.Sprintf("0x%x", bcs.Owner),
		PendingTransfer: bcs.PendingTransfer,
		OwnerCapId:      bcs.OwnerCapId,
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

type bcsNewOwnableStateEvent struct {
	OwnerCapId bind.Object
	Owner      [32]byte
}

func convertNewOwnableStateEventFromBCS(bcs bcsNewOwnableStateEvent) (NewOwnableStateEvent, error) {

	return NewOwnableStateEvent{
		OwnerCapId: bcs.OwnerCapId,
		Owner:      fmt.Sprintf("0x%x", bcs.Owner),
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
	bind.RegisterStructDecoder("ccip_onramp::ownable::OwnerCap", func(data []byte) (interface{}, error) {
		var result OwnerCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OwnerCap
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::OwnerCap>", func(data []byte) (interface{}, error) {
		var results []OwnerCap
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::ownable::OwnableState", func(data []byte) (interface{}, error) {
		var temp bcsOwnableState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOwnableStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OwnableState
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::OwnableState>", func(data []byte) (interface{}, error) {
		var temps []bcsOwnableState
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OwnableState, len(temps))
		for i, temp := range temps {
			result, err := convertOwnableStateFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::ownable::PendingTransfer", func(data []byte) (interface{}, error) {
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
	// Register vector decoder for PendingTransfer
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::PendingTransfer>", func(data []byte) (interface{}, error) {
		var temps []bcsPendingTransfer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]PendingTransfer, len(temps))
		for i, temp := range temps {
			result, err := convertPendingTransferFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::ownable::PublisherKey", func(data []byte) (interface{}, error) {
		var result PublisherKey
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for PublisherKey
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::PublisherKey>", func(data []byte) (interface{}, error) {
		var results []PublisherKey
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::ownable::NewOwnableStateEvent", func(data []byte) (interface{}, error) {
		var temp bcsNewOwnableStateEvent
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertNewOwnableStateEventFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for NewOwnableStateEvent
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::NewOwnableStateEvent>", func(data []byte) (interface{}, error) {
		var temps []bcsNewOwnableStateEvent
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]NewOwnableStateEvent, len(temps))
		for i, temp := range temps {
			result, err := convertNewOwnableStateEventFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::ownable::OwnershipTransferRequested", func(data []byte) (interface{}, error) {
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
	// Register vector decoder for OwnershipTransferRequested
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::OwnershipTransferRequested>", func(data []byte) (interface{}, error) {
		var temps []bcsOwnershipTransferRequested
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OwnershipTransferRequested, len(temps))
		for i, temp := range temps {
			result, err := convertOwnershipTransferRequestedFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::ownable::OwnershipTransferAccepted", func(data []byte) (interface{}, error) {
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
	// Register vector decoder for OwnershipTransferAccepted
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::OwnershipTransferAccepted>", func(data []byte) (interface{}, error) {
		var temps []bcsOwnershipTransferAccepted
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OwnershipTransferAccepted, len(temps))
		for i, temp := range temps {
			result, err := convertOwnershipTransferAcceptedFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::ownable::OwnershipTransferred", func(data []byte) (interface{}, error) {
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
	// Register vector decoder for OwnershipTransferred
	bind.RegisterStructDecoder("vector<ccip_onramp::ownable::OwnershipTransferred>", func(data []byte) (interface{}, error) {
		var temps []bcsOwnershipTransferred
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OwnershipTransferred, len(temps))
		for i, temp := range temps {
			result, err := convertOwnershipTransferredFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
}

// DefaultKey executes the default_key Move function.
func (c *OwnableContract) DefaultKey(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.DefaultKey()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// New executes the new Move function.
func (c *OwnableContract) New(ctx context.Context, opts *bind.CallOpts, uid string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.New(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// NewWithKey executes the new_with_key Move function.
func (c *OwnableContract) NewWithKey(ctx context.Context, opts *bind.CallOpts, typeArgs []string, uid string, key bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.NewWithKey(typeArgs, uid, key)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// OwnerCapId executes the owner_cap_id Move function.
func (c *OwnableContract) OwnerCapId(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.OwnerCapId(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *OwnableContract) Owner(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.Owner(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *OwnableContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.HasPendingTransfer(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *OwnableContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *OwnableContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *OwnableContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AttachPublisher executes the attach_publisher Move function.
func (c *OwnableContract) AttachPublisher(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, publisher bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.AttachPublisher(ownerCap, publisher)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// BorrowPublisher executes the borrow_publisher Move function.
func (c *OwnableContract) BorrowPublisher(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.BorrowPublisher(ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *OwnableContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state OwnableState, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.TransferOwnership(ownerCap, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *OwnableContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.AcceptOwnership(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *OwnableContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state OwnableState, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.AcceptOwnershipFromObject(state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAcceptOwnership executes the mcms_accept_ownership Move function.
func (c *OwnableContract) McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, state OwnableState, mcms string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.McmsAcceptOwnership(state, mcms)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *OwnableContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state OwnableState, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.ExecuteOwnershipTransfer(ownerCap, state, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *OwnableContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ownerCap bind.Object, state OwnableState, registry bind.Object, to string, publisherWrapper bind.Object, proof bind.Object, allowedModules [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.ownableEncoder.ExecuteOwnershipTransferToMcms(typeArgs, ownerCap, state, registry, to, publisherWrapper, proof, allowedModules)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DefaultKey executes the default_key Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OwnableDevInspect) DefaultKey(ctx context.Context, opts *bind.CallOpts) ([]byte, error) {
	encoded, err := d.contract.ownableEncoder.DefaultKey()
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

// New executes the new Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: OwnableState
//	[1]: OwnerCap
func (d *OwnableDevInspect) New(ctx context.Context, opts *bind.CallOpts, uid string) ([]any, error) {
	encoded, err := d.contract.ownableEncoder.New(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// NewWithKey executes the new_with_key Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: OwnableState
//	[1]: OwnerCap
func (d *OwnableDevInspect) NewWithKey(ctx context.Context, opts *bind.CallOpts, typeArgs []string, uid string, key bind.Object) ([]any, error) {
	encoded, err := d.contract.ownableEncoder.NewWithKey(typeArgs, uid, key)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// OwnerCapId executes the owner_cap_id Move function using DevInspect to get return values.
//
// Returns: ID
func (d *OwnableDevInspect) OwnerCapId(ctx context.Context, opts *bind.CallOpts, state OwnableState) (bind.Object, error) {
	encoded, err := d.contract.ownableEncoder.OwnerCapId(state)
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
func (d *OwnableDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, state OwnableState) (string, error) {
	encoded, err := d.contract.ownableEncoder.Owner(state)
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
func (d *OwnableDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state OwnableState) (bool, error) {
	encoded, err := d.contract.ownableEncoder.HasPendingTransfer(state)
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
func (d *OwnableDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*string, error) {
	encoded, err := d.contract.ownableEncoder.PendingTransferFrom(state)
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
func (d *OwnableDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*string, error) {
	encoded, err := d.contract.ownableEncoder.PendingTransferTo(state)
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
func (d *OwnableDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state OwnableState) (*bool, error) {
	encoded, err := d.contract.ownableEncoder.PendingTransferAccepted(state)
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

// BorrowPublisher executes the borrow_publisher Move function using DevInspect to get return values.
//
// Returns: &Publisher
func (d *OwnableDevInspect) BorrowPublisher(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object) (bind.Object, error) {
	encoded, err := d.contract.ownableEncoder.BorrowPublisher(ownerCap)
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

type ownableEncoder struct {
	*bind.BoundContract
}

// DefaultKey encodes a call to the default_key Move function.
func (c ownableEncoder) DefaultKey() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("default_key", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"vector<u8>",
	})
}

// DefaultKeyWithArgs encodes a call to the default_key Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) DefaultKeyWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("default_key", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// New encodes a call to the new Move function.
func (c ownableEncoder) New(uid string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new", typeArgsList, typeParamsList, []string{
		"&mut UID",
	}, []any{
		uid,
	}, []string{
		"ccip_onramp::ownable::OwnableState",
		"ccip_onramp::ownable::OwnerCap",
	})
}

// NewWithArgs encodes a call to the new Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) NewWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_onramp::ownable::OwnableState",
		"ccip_onramp::ownable::OwnerCap",
	})
}

// NewWithKey encodes a call to the new_with_key Move function.
func (c ownableEncoder) NewWithKey(typeArgs []string, uid string, key bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"K",
	}
	return c.EncodeCallArgsWithGenerics("new_with_key", typeArgsList, typeParamsList, []string{
		"&mut UID",
		"K",
	}, []any{
		uid,
		key,
	}, []string{
		"ccip_onramp::ownable::OwnableState",
		"ccip_onramp::ownable::OwnerCap",
	})
}

// NewWithKeyWithArgs encodes a call to the new_with_key Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) NewWithKeyWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut UID",
		"K",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"K",
	}
	return c.EncodeCallArgsWithGenerics("new_with_key", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_onramp::ownable::OwnableState",
		"ccip_onramp::ownable::OwnerCap",
	})
}

// OwnerCapId encodes a call to the owner_cap_id Move function.
func (c ownableEncoder) OwnerCapId(state OwnableState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner_cap_id", typeArgsList, typeParamsList, []string{
		"&OwnableState",
	}, []any{
		state,
	}, []string{
		"ID",
	})
}

// OwnerCapIdWithArgs encodes a call to the owner_cap_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) OwnerCapIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnableState",
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

// Owner encodes a call to the owner Move function.
func (c ownableEncoder) Owner(state OwnableState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&OwnableState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnableState",
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
func (c ownableEncoder) HasPendingTransfer(state OwnableState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&OwnableState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnableState",
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
func (c ownableEncoder) PendingTransferFrom(state OwnableState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&OwnableState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnableState",
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
func (c ownableEncoder) PendingTransferTo(state OwnableState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&OwnableState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnableState",
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
func (c ownableEncoder) PendingTransferAccepted(state OwnableState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&OwnableState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnableState",
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

// AttachPublisher encodes a call to the attach_publisher Move function.
func (c ownableEncoder) AttachPublisher(ownerCap bind.Object, publisher bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("attach_publisher", typeArgsList, typeParamsList, []string{
		"&mut OwnerCap",
		"Publisher",
	}, []any{
		ownerCap,
		publisher,
	}, nil)
}

// AttachPublisherWithArgs encodes a call to the attach_publisher Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) AttachPublisherWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OwnerCap",
		"Publisher",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("attach_publisher", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// BorrowPublisher encodes a call to the borrow_publisher Move function.
func (c ownableEncoder) BorrowPublisher(ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("borrow_publisher", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
	}, []any{
		ownerCap,
	}, []string{
		"&Publisher",
	})
}

// BorrowPublisherWithArgs encodes a call to the borrow_publisher Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) BorrowPublisherWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("borrow_publisher", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&Publisher",
	})
}

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c ownableEncoder) TransferOwnership(ownerCap bind.Object, state OwnableState, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut OwnableState",
		"address",
	}, []any{
		ownerCap,
		state,
		to,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut OwnableState",
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
func (c ownableEncoder) AcceptOwnership(state OwnableState) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut OwnableState",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OwnableState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c ownableEncoder) AcceptOwnershipFromObject(state OwnableState, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut OwnableState",
		"&UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OwnableState",
		"&UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAcceptOwnership encodes a call to the mcms_accept_ownership Move function.
func (c ownableEncoder) McmsAcceptOwnership(state OwnableState, mcms string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut OwnableState",
		"address",
	}, []any{
		state,
		mcms,
	}, nil)
}

// McmsAcceptOwnershipWithArgs encodes a call to the mcms_accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OwnableState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c ownableEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, state OwnableState, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"ccip_onramp::ownable::OwnerCap",
		"&mut OwnableState",
		"address",
	}, []any{
		ownerCap,
		state,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_onramp::ownable::OwnerCap",
		"&mut OwnableState",
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
func (c ownableEncoder) ExecuteOwnershipTransferToMcms(typeArgs []string, ownerCap bind.Object, state OwnableState, registry bind.Object, to string, publisherWrapper bind.Object, proof bind.Object, allowedModules [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"ccip_onramp::ownable::OwnerCap",
		"&mut OwnableState",
		"&mut Registry",
		"address",
		"PublisherWrapper<T>",
		"T",
		"vector<vector<u8>>",
	}, []any{
		ownerCap,
		state,
		registry,
		to,
		publisherWrapper,
		proof,
		allowedModules,
	}, nil)
}

// ExecuteOwnershipTransferToMcmsWithArgs encodes a call to the execute_ownership_transfer_to_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c ownableEncoder) ExecuteOwnershipTransferToMcmsWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_onramp::ownable::OwnerCap",
		"&mut OwnableState",
		"&mut Registry",
		"address",
		"PublisherWrapper<T>",
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
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}
