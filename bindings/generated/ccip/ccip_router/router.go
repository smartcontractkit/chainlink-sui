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

const FunctionInfo = `[{"package":"ccip_router","module":"router","name":"accept_ownership","parameters":[{"name":"state","type":"RouterState"}]},{"package":"ccip_router","module":"router","name":"accept_ownership_from_object","parameters":[{"name":"state","type":"RouterState"},{"name":"from","type":"sui::object::UID"}]},{"package":"ccip_router","module":"router","name":"execute_ownership_transfer","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"RouterState"},{"name":"to","type":"address"}]},{"package":"ccip_router","module":"router","name":"execute_ownership_transfer_to_mcms","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"state","type":"RouterState"},{"name":"registry","type":"Registry"},{"name":"to","type":"address"}]},{"package":"ccip_router","module":"router","name":"get_dest_chains","parameters":[{"name":"router","type":"RouterState"}]},{"package":"ccip_router","module":"router","name":"get_on_ramp","parameters":[{"name":"router","type":"RouterState"},{"name":"dest_chain_selector","type":"u64"}]},{"package":"ccip_router","module":"router","name":"get_uid","parameters":[{"name":"router_object","type":"RouterObject"}]},{"package":"ccip_router","module":"router","name":"has_pending_transfer","parameters":[{"name":"state","type":"RouterState"}]},{"package":"ccip_router","module":"router","name":"is_chain_supported","parameters":[{"name":"router","type":"RouterState"},{"name":"dest_chain_selector","type":"u64"}]},{"package":"ccip_router","module":"router","name":"owner","parameters":[{"name":"state","type":"RouterState"}]},{"package":"ccip_router","module":"router","name":"pending_transfer_accepted","parameters":[{"name":"state","type":"RouterState"}]},{"package":"ccip_router","module":"router","name":"pending_transfer_from","parameters":[{"name":"state","type":"RouterState"}]},{"package":"ccip_router","module":"router","name":"pending_transfer_to","parameters":[{"name":"state","type":"RouterState"}]},{"package":"ccip_router","module":"router","name":"set_on_ramps","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"router","type":"RouterState"},{"name":"dest_chain_selectors","type":"vector<u64>"},{"name":"on_ramp_package_ids","type":"vector<address>"}]},{"package":"ccip_router","module":"router","name":"transfer_ownership","parameters":[{"name":"state","type":"RouterState"},{"name":"owner_cap","type":"OwnerCap"},{"name":"new_owner","type":"address"}]},{"package":"ccip_router","module":"router","name":"type_and_version","parameters":null}]`

type IRouter interface {
	GetUid(ctx context.Context, opts *bind.CallOpts, routerObject bind.Object) (*models.SuiTransactionBlockResponse, error)
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetOnRamp(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetDestChains(ctx context.Context, opts *bind.CallOpts, router bind.Object) (*models.SuiTransactionBlockResponse, error)
	SetOnRamps(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampPackageIds []string) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsSetOnRamps(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IRouterDevInspect
	Encoder() RouterEncoder
	Bound() bind.IBoundContract
}

type IRouterDevInspect interface {
	GetUid(ctx context.Context, opts *bind.CallOpts, routerObject bind.Object) (bind.Object, error)
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (bool, error)
	GetOnRamp(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (string, error)
	GetDestChains(ctx context.Context, opts *bind.CallOpts, router bind.Object) ([]uint64, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type RouterEncoder interface {
	GetUid(routerObject bind.Object) (*bind.EncodedCall, error)
	GetUidWithArgs(args ...any) (*bind.EncodedCall, error)
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	IsChainSupported(router bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	IsChainSupportedWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOnRamp(router bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetOnRampWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestChains(router bind.Object) (*bind.EncodedCall, error)
	GetDestChainsWithArgs(args ...any) (*bind.EncodedCall, error)
	SetOnRamps(ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampPackageIds []string) (*bind.EncodedCall, error)
	SetOnRampsWithArgs(args ...any) (*bind.EncodedCall, error)
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
	McmsAcceptOwnership(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsSetOnRamps(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsSetOnRampsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsTransferOwnership(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransfer(state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsAddAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAddAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRemoveAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsRemoveAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewRouter(packageID string, client sui.ISuiAPI) (IRouter, error) {
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

func (c *RouterContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *RouterContract) Encoder() RouterEncoder {
	return c.routerEncoder
}

func (c *RouterContract) DevInspect() IRouterDevInspect {
	return c.devInspect
}

type ROUTER struct {
}

type RouterObject struct {
	Id string `move:"sui::object::UID"`
}

type OnRampSet struct {
	DestChainSelector uint64 `move:"u64"`
	OnRampPackageId   string `move:"address"`
}

type RouterState struct {
	Id               string      `move:"sui::object::UID"`
	OwnableState     bind.Object `move:"OwnableState"`
	OnRampPackageIds bind.Object `move:"VecMap<u64, address>"`
}

type RouterStatePointer struct {
	Id             string `move:"sui::object::UID"`
	RouterObjectId string `move:"address"`
}

type McmsCallback struct {
}

type McmsAcceptOwnershipProof struct {
}

type bcsOnRampSet struct {
	DestChainSelector uint64
	OnRampPackageId   [32]byte
}

func convertOnRampSetFromBCS(bcs bcsOnRampSet) (OnRampSet, error) {

	return OnRampSet{
		DestChainSelector: bcs.DestChainSelector,
		OnRampPackageId:   fmt.Sprintf("0x%x", bcs.OnRampPackageId),
	}, nil
}

type bcsRouterStatePointer struct {
	Id             string
	RouterObjectId [32]byte
}

func convertRouterStatePointerFromBCS(bcs bcsRouterStatePointer) (RouterStatePointer, error) {

	return RouterStatePointer{
		Id:             bcs.Id,
		RouterObjectId: fmt.Sprintf("0x%x", bcs.RouterObjectId),
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
	// Register vector decoder for ROUTER
	bind.RegisterStructDecoder("vector<ccip_router::router::ROUTER>", func(data []byte) (interface{}, error) {
		var results []ROUTER
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::RouterObject", func(data []byte) (interface{}, error) {
		var result RouterObject
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for RouterObject
	bind.RegisterStructDecoder("vector<ccip_router::router::RouterObject>", func(data []byte) (interface{}, error) {
		var results []RouterObject
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
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
	// Register vector decoder for OnRampSet
	bind.RegisterStructDecoder("vector<ccip_router::router::OnRampSet>", func(data []byte) (interface{}, error) {
		var temps []bcsOnRampSet
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OnRampSet, len(temps))
		for i, temp := range temps {
			result, err := convertOnRampSetFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::RouterState", func(data []byte) (interface{}, error) {
		var result RouterState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for RouterState
	bind.RegisterStructDecoder("vector<ccip_router::router::RouterState>", func(data []byte) (interface{}, error) {
		var results []RouterState
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::RouterStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsRouterStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertRouterStatePointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for RouterStatePointer
	bind.RegisterStructDecoder("vector<ccip_router::router::RouterStatePointer>", func(data []byte) (interface{}, error) {
		var temps []bcsRouterStatePointer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]RouterStatePointer, len(temps))
		for i, temp := range temps {
			result, err := convertRouterStatePointerFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsCallback
	bind.RegisterStructDecoder("vector<ccip_router::router::McmsCallback>", func(data []byte) (interface{}, error) {
		var results []McmsCallback
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_router::router::McmsAcceptOwnershipProof", func(data []byte) (interface{}, error) {
		var result McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for McmsAcceptOwnershipProof
	bind.RegisterStructDecoder("vector<ccip_router::router::McmsAcceptOwnershipProof>", func(data []byte) (interface{}, error) {
		var results []McmsAcceptOwnershipProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
}

// GetUid executes the get_uid Move function.
func (c *RouterContract) GetUid(ctx context.Context, opts *bind.CallOpts, routerObject bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.GetUid(routerObject)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
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

// GetOnRamp executes the get_on_ramp Move function.
func (c *RouterContract) GetOnRamp(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.GetOnRamp(router, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestChains executes the get_dest_chains Move function.
func (c *RouterContract) GetDestChains(ctx context.Context, opts *bind.CallOpts, router bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.GetDestChains(router)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetOnRamps executes the set_on_ramps Move function.
func (c *RouterContract) SetOnRamps(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampPackageIds []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.SetOnRamps(ownerCap, router, destChainSelectors, onRampPackageIds)
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

// McmsAcceptOwnership executes the mcms_accept_ownership Move function.
func (c *RouterContract) McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsAcceptOwnership(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *RouterContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.ExecuteOwnershipTransfer(ownerCap, state, to)
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

// McmsSetOnRamps executes the mcms_set_on_ramps Move function.
func (c *RouterContract) McmsSetOnRamps(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsSetOnRamps(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsTransferOwnership executes the mcms_transfer_ownership Move function.
func (c *RouterContract) McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsTransferOwnership(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsExecuteOwnershipTransfer executes the mcms_execute_ownership_transfer Move function.
func (c *RouterContract) McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsExecuteOwnershipTransfer(state, registry, deployerState, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAddAllowedModules executes the mcms_add_allowed_modules Move function.
func (c *RouterContract) McmsAddAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsAddAllowedModules(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRemoveAllowedModules executes the mcms_remove_allowed_modules Move function.
func (c *RouterContract) McmsRemoveAllowedModules(ctx context.Context, opts *bind.CallOpts, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.McmsRemoveAllowedModules(registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetUid executes the get_uid Move function using DevInspect to get return values.
//
// Returns: &mut UID
func (d *RouterDevInspect) GetUid(ctx context.Context, opts *bind.CallOpts, routerObject bind.Object) (bind.Object, error) {
	encoded, err := d.contract.routerEncoder.GetUid(routerObject)
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

// GetOnRamp executes the get_on_ramp Move function using DevInspect to get return values.
//
// Returns: address
func (d *RouterDevInspect) GetOnRamp(ctx context.Context, opts *bind.CallOpts, router bind.Object, destChainSelector uint64) (string, error) {
	encoded, err := d.contract.routerEncoder.GetOnRamp(router, destChainSelector)
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

// GetDestChains executes the get_dest_chains Move function using DevInspect to get return values.
//
// Returns: vector<u64>
func (d *RouterDevInspect) GetDestChains(ctx context.Context, opts *bind.CallOpts, router bind.Object) ([]uint64, error) {
	encoded, err := d.contract.routerEncoder.GetDestChains(router)
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

// GetUid encodes a call to the get_uid Move function.
func (c routerEncoder) GetUid(routerObject bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_uid", typeArgsList, typeParamsList, []string{
		"&mut RouterObject",
	}, []any{
		routerObject,
	}, []string{
		"&mut UID",
	})
}

// GetUidWithArgs encodes a call to the get_uid Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) GetUidWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut RouterObject",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_uid", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&mut UID",
	})
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

// GetOnRamp encodes a call to the get_on_ramp Move function.
func (c routerEncoder) GetOnRamp(router bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp", typeArgsList, typeParamsList, []string{
		"&RouterState",
		"u64",
	}, []any{
		router,
		destChainSelector,
	}, []string{
		"address",
	})
}

// GetOnRampWithArgs encodes a call to the get_on_ramp Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) GetOnRampWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_on_ramp", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetDestChains encodes a call to the get_dest_chains Move function.
func (c routerEncoder) GetDestChains(router bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chains", typeArgsList, typeParamsList, []string{
		"&RouterState",
	}, []any{
		router,
	}, []string{
		"vector<u64>",
	})
}

// GetDestChainsWithArgs encodes a call to the get_dest_chains Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) GetDestChainsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&RouterState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chains", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
	})
}

// SetOnRamps encodes a call to the set_on_ramps Move function.
func (c routerEncoder) SetOnRamps(ownerCap bind.Object, router bind.Object, destChainSelectors []uint64, onRampPackageIds []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_on_ramps", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut RouterState",
		"vector<u64>",
		"vector<address>",
	}, []any{
		ownerCap,
		router,
		destChainSelectors,
		onRampPackageIds,
	}, nil)
}

// SetOnRampsWithArgs encodes a call to the set_on_ramps Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) SetOnRampsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut RouterState",
		"vector<u64>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_on_ramps", typeArgsList, typeParamsList, expectedParams, args, nil)
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

// McmsAcceptOwnership encodes a call to the mcms_accept_ownership Move function.
func (c routerEncoder) McmsAcceptOwnership(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsAcceptOwnershipWithArgs encodes a call to the mcms_accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c routerEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, state bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut RouterState",
		"address",
	}, []any{
		ownerCap,
		state,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut RouterState",
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

// McmsSetOnRamps encodes a call to the mcms_set_on_ramps Move function.
func (c routerEncoder) McmsSetOnRamps(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_on_ramps", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsSetOnRampsWithArgs encodes a call to the mcms_set_on_ramps Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsSetOnRampsWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_set_on_ramps", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsTransferOwnership encodes a call to the mcms_transfer_ownership Move function.
func (c routerEncoder) McmsTransferOwnership(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsTransferOwnershipWithArgs encodes a call to the mcms_transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsExecuteOwnershipTransfer encodes a call to the mcms_execute_ownership_transfer Move function.
func (c routerEncoder) McmsExecuteOwnershipTransfer(state bind.Object, registry bind.Object, deployerState bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&mut RouterState",
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		deployerState,
		params,
	}, nil)
}

// McmsExecuteOwnershipTransferWithArgs encodes a call to the mcms_execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut RouterState",
		"&mut Registry",
		"&mut DeployerState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAddAllowedModules encodes a call to the mcms_add_allowed_modules Move function.
func (c routerEncoder) McmsAddAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_add_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		params,
	}, nil)
}

// McmsAddAllowedModulesWithArgs encodes a call to the mcms_add_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsAddAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_add_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRemoveAllowedModules encodes a call to the mcms_remove_allowed_modules Move function.
func (c routerEncoder) McmsRemoveAllowedModules(registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_remove_allowed_modules", typeArgsList, typeParamsList, []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		registry,
		params,
	}, nil)
}

// McmsRemoveAllowedModulesWithArgs encodes a call to the mcms_remove_allowed_modules Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) McmsRemoveAllowedModulesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_remove_allowed_modules", typeArgsList, typeParamsList, expectedParams, args, nil)
}
