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
	EmitOnRampSetEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IRouterDevInspect
	Encoder() RouterEncoder
	Bound() bind.IBoundContract
}

type IRouterDevInspect interface {
}

type RouterEncoder interface {
	EmitOnRampSetEvent(destChainSelector uint64) (*bind.EncodedCall, error)
	EmitOnRampSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
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
	contract, err := bind.NewBoundContract(packageID, "test", "router", client)
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

type OnRampSet struct {
	DestChainSelector uint64     `move:"u64"`
	OnRampInfo        OnRampInfo `move:"OnRampInfo"`
}

type RouterState struct {
	Id          string       `move:"sui::object::UID"`
	OnRampInfos []OnRampInfo `move:"vector<OnRampInfo>"`
}

type OnRampInfo struct {
	OnrampAddress string `move:"address"`
	OnrampVersion []byte `move:"vector<u8>"`
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
	bind.RegisterStructDecoder("test::router::OnRampSet", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::router::RouterState", func(data []byte) (interface{}, error) {
		var result RouterState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::router::OnRampInfo", func(data []byte) (interface{}, error) {
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
}

// EmitOnRampSetEvent executes the emit_on_ramp_set_event Move function.
func (c *RouterContract) EmitOnRampSetEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.routerEncoder.EmitOnRampSetEvent(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type routerEncoder struct {
	*bind.BoundContract
}

// EmitOnRampSetEvent encodes a call to the emit_on_ramp_set_event Move function.
func (c routerEncoder) EmitOnRampSetEvent(destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_on_ramp_set_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		destChainSelector,
	}, nil)
}

// EmitOnRampSetEventWithArgs encodes a call to the emit_on_ramp_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c routerEncoder) EmitOnRampSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_on_ramp_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}
