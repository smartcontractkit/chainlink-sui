// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_router

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

type IRouter interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetOnRampInfo(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetOnRampInfos(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelectors []uint64) (*models.SuiTransactionBlockResponse, error)
	GetOnRampVersion(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) (*models.SuiTransactionBlockResponse, error)
	GetOnRampAddress(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) (*models.SuiTransactionBlockResponse, error)
	SetOnRampInfos(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipAsMcms(ctx context.Context, opts *bind.CallOpts, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IRouterDevInspect
	Encoder() RouterEncoder
}

type IRouterDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (bool, error)
	GetOnRampInfo(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) ([]any, error)
	GetOnRampInfos(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelectors []uint64) ([]OnRampInfo, error)
	GetOnRampVersion(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) ([]byte, error)
	GetOnRampAddress(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) (string, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type RouterEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	IsChainSupported(router bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	IsChainSupportedWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOnRampInfo(router bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetOnRampInfoWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOnRampInfos(router bind.Object, destChainSelectors []uint64) (*bind.EncodedCall, error)
	GetOnRampInfosWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOnRampVersion(info OnRampInfo) (*bind.EncodedCall, error)
	GetOnRampVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOnRampAddress(info OnRampInfo) (*bind.EncodedCall, error)
	GetOnRampAddressWithArgs(args ...any) (*bind.EncodedCall, error)
	SetOnRampInfos(ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) (*bind.EncodedCall, error)
	SetOnRampInfosWithArgs(args ...any) (*bind.EncodedCall, error)
	Owner(state bind.Object) (*bind.EncodedCall, error)
	OwnerWithArgs(args ...any) (*bind.EncodedCall, error)
	HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error)
	HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error)
	TransferOwnership(state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(state bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error)
	AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipAsMcms(state bind.Object, params bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipAsMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
}

type RouterContract struct {
	*bind.BoundContract
	routerEncoder
	devInspect *RouterDevInspect
}

type RouterDevInspect struct {
	contract *RouterContract
}

var _ IRouter = (*RouterContract)(nil)
var _ IRouterDevInspect = (*RouterDevInspect)(nil)

func NewRouter(packageID string, client sui.ISuiAPI) (*RouterContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_router", "router", client)
	if err != nil {
		return nil, err
	}

	c := &RouterContract{
		BoundContract: contract,
		routerEncoder: routerEncoder{BoundContract: contract},
	}
	c.devInspect = &RouterDevInspect{contract: c}
	return c, nil
}

func (c *RouterContract) Encoder() RouterEncoder {
	return c.routerEncoder
}

func (c *RouterContract) DevInspect() IRouterDevInspect {
	return c.devInspect
}

type ROUTER struct {
}

type OnRampSet struct {
	DestChainSelector uint64     `move:"u64"`
	OnRampInfo        OnRampInfo `move:"OnRampInfo"`
}

type OnRampInfo struct {
	OnrampAddress string `move:"address"`
	OnrampVersion []byte `move:"vector<u8>"`
}

type RouterState struct {
	Id           string      `move:"sui::object::UID"`
	OwnableState bind.Object `move:"OwnableState"`
	OnRampInfos  bind.Object `move:"Table<u64, OnRampInfo>"`
}

type McmsCallback struct {
}

type bcsOnRampSet struct {
	DestChainSelector uint64
	OnRampInfo        bcsOnRampInfo
}

func convertOnRampSetFromBCS(bcs bcsOnRampSet) (OnRampSet, error) {
	OnRampInfoField, err := convertOnRampInfoFromBCS(bcs.OnRampInfo)
	if err != nil {
		return OnRampSet{}, fmt.Errorf("failed to convert nested struct OnRampInfo: %w", err)
	}

	return OnRampSet{
		DestChainSelector: bcs.DestChainSelector,
		OnRampInfo:        OnRampInfoField,
	}, nil
}

type bcsOnRampInfo struct {
	OnrampAddress [32]byte
	OnrampVersion []byte
}

func convertOnRampInfoFromBCS(bcs bcsOnRampInfo) (OnRampInfo, error) {

	return OnRampInfo{
		OnrampAddress: fmt.Sprintf("0x%x", bcs.OnrampAddress),
		OnrampVersion: bcs.OnrampVersion,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip_router::router::ROUTER", func(data []byte) (interface{}, error) {
		var result ROUTER
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::OnRampSet", func(data []byte) (interface{}, error) {
		var temp bcsOnRampSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOnRampSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::OnRampInfo", func(data []byte) (interface{}, error) {
		var temp bcsOnRampInfo
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOnRampInfoFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::RouterState", func(data []byte) (interface{}, error) {
		var result RouterState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *RouterContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsChainSupported executes the is_chain_supported Move function.
func (c *RouterContract) IsChainSupported(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.IsChainSupported(router, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOnRampInfo executes the get_on_ramp_info Move function.
func (c *RouterContract) GetOnRampInfo(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.GetOnRampInfo(router, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOnRampInfos executes the get_on_ramp_infos Move function.
func (c *RouterContract) GetOnRampInfos(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelectors []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.GetOnRampInfos(router, destChainSelectors)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOnRampVersion executes the get_on_ramp_version Move function.
func (c *RouterContract) GetOnRampVersion(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.GetOnRampVersion(info)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOnRampAddress executes the get_on_ramp_address Move function.
func (c *RouterContract) GetOnRampAddress(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.GetOnRampAddress(info)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetOnRampInfos executes the set_on_ramp_infos Move function.
func (c *RouterContract) SetOnRampInfos(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.SetOnRampInfos(ownerCap, router, destChainSelectors, onRampAddresses, onRampVersions)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *RouterContract) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.Owner(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *RouterContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.HasPendingTransfer(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *RouterContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *RouterContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *RouterContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *RouterContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.TransferOwnership(state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *RouterContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.AcceptOwnership(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *RouterContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.AcceptOwnershipFromObject(state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipAsMcms executes the accept_ownership_as_mcms Move function.
func (c *RouterContract) AcceptOwnershipAsMcms(ctx context.Context, opts *bind.CallOpts, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.AcceptOwnershipAsMcms(state, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *RouterContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.ExecuteOwnershipTransfer(ownerCap, ownableState, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *RouterContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.ExecuteOwnershipTransferToMcms(ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *RouterContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *RouterContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsEntrypoint(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *RouterDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.routerEncoder.TypeAndVersion()
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

// IsChainSupported executes the is_chain_supported Move function using DevInspect to get return values.
//
// Returns: bool
func (d *RouterDevInspect) IsChainSupported(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (bool, error) {
	encoded, err := d.contract.routerEncoder.IsChainSupported(router, destChainSelector)
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

// GetOnRampInfo executes the get_on_ramp_info Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: vector<u8>
func (d *RouterDevInspect) GetOnRampInfo(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) ([]any, error) {
	encoded, err := d.contract.routerEncoder.GetOnRampInfo(router, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetOnRampInfos executes the get_on_ramp_infos Move function using DevInspect to get return values.
//
// Returns: vector<OnRampInfo>
func (d *RouterDevInspect) GetOnRampInfos(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelectors []uint64) ([]OnRampInfo, error) {
	encoded, err := d.contract.routerEncoder.GetOnRampInfos(router, destChainSelectors)
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
	result, ok := results[0].([]OnRampInfo)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []OnRampInfo, got %T", results[0])
	}
	return result, nil
}

// GetOnRampVersion executes the get_on_ramp_version Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *RouterDevInspect) GetOnRampVersion(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) ([]byte, error) {
	encoded, err := d.contract.routerEncoder.GetOnRampVersion(info)
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

// GetOnRampAddress executes the get_on_ramp_address Move function using DevInspect to get return values.
//
// Returns: address
func (d *RouterDevInspect) GetOnRampAddress(ctx context.Context, opts *bind.CallOpts, info OnRampInfo) (string, error) {
	encoded, err := d.contract.routerEncoder.GetOnRampAddress(info)
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

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *RouterDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.routerEncoder.Owner(state)
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
func (d *RouterDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error) {
	encoded, err := d.contract.routerEncoder.HasPendingTransfer(state)
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
func (d *RouterDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.routerEncoder.PendingTransferFrom(state)
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
func (d *RouterDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.routerEncoder.PendingTransferTo(state)
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
func (d *RouterDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error) {
	encoded, err := d.contract.routerEncoder.PendingTransferAccepted(state)
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

type routerEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c routerEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// IsChainSupported encodes a call to the is_chain_supported Move function.
func (c routerEncoder) IsChainSupported(router bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_chain_supported", typeArgsList, typeParamsList, []string{
		"&RouterState",
		"u64",
	}, []any{
		router,
		destChainSelector,
	}, []string{
		"bool",
	})
}

// IsChainSupportedWithArgs encodes a call to the is_chain_supported Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) IsChainSupportedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_chain_supported", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetOnRampInfo encodes a call to the get_on_ramp_info Move function.
func (c routerEncoder) GetOnRampInfo(router bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_info", typeArgsList, typeParamsList, []string{
		"&RouterState",
		"u64",
	}, []any{
		router,
		destChainSelector,
	}, []string{
		"address",
		"vector<u8>",
	})
}

// GetOnRampInfoWithArgs encodes a call to the get_on_ramp_info Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) GetOnRampInfoWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_info", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"vector<u8>",
	})
}

// GetOnRampInfos encodes a call to the get_on_ramp_infos Move function.
func (c routerEncoder) GetOnRampInfos(router bind.Object, destChainSelectors []uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_infos", typeArgsList, typeParamsList, []string{
		"&RouterState",
		"vector<u64>",
	}, []any{
		router,
		destChainSelectors,
	}, []string{
		"vector<ccip_router::router::OnRampInfo>",
	})
}

// GetOnRampInfosWithArgs encodes a call to the get_on_ramp_infos Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) GetOnRampInfosWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
		"vector<u64>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_infos", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ccip_router::router::OnRampInfo>",
	})
}

// GetOnRampVersion encodes a call to the get_on_ramp_version Move function.
func (c routerEncoder) GetOnRampVersion(info OnRampInfo) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_version", typeArgsList, typeParamsList, []string{
		"ccip_router::router::OnRampInfo",
	}, []any{
		info,
	}, []string{
		"vector<u8>",
	})
}

// GetOnRampVersionWithArgs encodes a call to the get_on_ramp_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) GetOnRampVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_router::router::OnRampInfo",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_version", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// GetOnRampAddress encodes a call to the get_on_ramp_address Move function.
func (c routerEncoder) GetOnRampAddress(info OnRampInfo) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_address", typeArgsList, typeParamsList, []string{
		"ccip_router::router::OnRampInfo",
	}, []any{
		info,
	}, []string{
		"address",
	})
}

// GetOnRampAddressWithArgs encodes a call to the get_on_ramp_address Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) GetOnRampAddressWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_router::router::OnRampInfo",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp_address", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// SetOnRampInfos encodes a call to the set_on_ramp_infos Move function.
func (c routerEncoder) SetOnRampInfos(ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampAddresses []string, onRampVersions [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_on_ramp_infos", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut RouterState",
		"vector<u64>",
		"vector<address>",
		"vector<vector<u8>>",
	}, []any{
		ownerCap,
		router,
		destChainSelectors,
		onRampAddresses,
		onRampVersions,
	}, nil)
}

// SetOnRampInfosWithArgs encodes a call to the set_on_ramp_infos Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) SetOnRampInfosWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut RouterState",
		"vector<u64>",
		"vector<address>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_on_ramp_infos", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// Owner encodes a call to the owner Move function.
func (c routerEncoder) Owner(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&RouterState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
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
func (c routerEncoder) HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&RouterState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
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
func (c routerEncoder) PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&RouterState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
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
func (c routerEncoder) PendingTransferTo(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&RouterState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
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
func (c routerEncoder) PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&RouterState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
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

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c routerEncoder) TransferOwnership(state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"&OwnerCap",
		"address",
	}, []any{
		state,
		ownerCap,
		newOwner,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut RouterState",
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
func (c routerEncoder) AcceptOwnership(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut RouterState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c routerEncoder) AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut RouterState",
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipAsMcms encodes a call to the accept_ownership_as_mcms Move function.
func (c routerEncoder) AcceptOwnershipAsMcms(state bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_mcms", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"ExecutingCallbackParams",
	}, []any{
		state,
		params,
	}, nil)
}

// AcceptOwnershipAsMcmsWithArgs encodes a call to the accept_ownership_as_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) AcceptOwnershipAsMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut RouterState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c routerEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut OwnableState",
		"address",
	}, []any{
		ownerCap,
		ownableState,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
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
func (c routerEncoder) ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut RouterState",
		"&mut Registry",
		"address",
	}, []any{
		ownerCap,
		state,
		registry,
		to,
	}, nil)
}

// ExecuteOwnershipTransferToMcmsWithArgs encodes a call to the execute_ownership_transfer_to_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut RouterState",
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

// McmsRegisterUpgradeCap encodes a call to the mcms_register_upgrade_cap Move function.
func (c routerEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, []string{
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}, []any{
		upgradeCap,
		registry,
		state,
	}, nil)
}

// McmsRegisterUpgradeCapWithArgs encodes a call to the mcms_register_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsEntrypoint encodes a call to the mcms_entrypoint Move function.
func (c routerEncoder) McmsEntrypoint(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsEntrypointWithArgs encodes a call to the mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut RouterState",
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
