// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_lock_release_token_pool

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

type ILockReleaseTokenPool interface {
	EmitLockedOrBurnedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitReleasedOrMintedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitLiquidityAddedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64) (*models.SuiTransactionBlockResponse, error)
	EmitLiquidityRemovedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRebalancerSetEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error)
	CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error)
	CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error)
	DevInspect() ILockReleaseTokenPoolDevInspect
	Encoder() LockReleaseTokenPoolEncoder
	Bound() bind.IBoundContract
}

type ILockReleaseTokenPoolDevInspect interface {
	CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
	CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
	CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
}

type LockReleaseTokenPoolEncoder interface {
	EmitLockedOrBurnedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitLockedOrBurnedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitReleasedOrMintedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitReleasedOrMintedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitRemotePoolAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitRemotePoolAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitRemotePoolRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitRemotePoolRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitChainAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitChainAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitChainRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitChainRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitLiquidityAddedEvent(amount uint64) (*bind.EncodedCall, error)
	EmitLiquidityAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitLiquidityRemovedEvent(amount uint64) (*bind.EncodedCall, error)
	EmitLiquidityRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitRebalancerSetEvent() (*bind.EncodedCall, error)
	EmitRebalancerSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
	CreatePendingTransferRequest(from string, to string) (*bind.EncodedCall, error)
	CreatePendingTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateAcceptedTransferRequest(from string, to string) (*bind.EncodedCall, error)
	CreateAcceptedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateRejectedTransferRequest(from string, to string) (*bind.EncodedCall, error)
	CreateRejectedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error)
}

type LockReleaseTokenPoolContract struct {
	*bind.BoundContract
	lockReleaseTokenPoolEncoder
	devInspect *LockReleaseTokenPoolDevInspect
}

type LockReleaseTokenPoolDevInspect struct {
	contract *LockReleaseTokenPoolContract
}

var _ ILockReleaseTokenPool = (*LockReleaseTokenPoolContract)(nil)
var _ ILockReleaseTokenPoolDevInspect = (*LockReleaseTokenPoolDevInspect)(nil)

func NewLockReleaseTokenPool(packageID string, client sui.ISuiAPI) (ILockReleaseTokenPool, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "lock_release_token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &LockReleaseTokenPoolContract{
		BoundContract:               contract,
		lockReleaseTokenPoolEncoder: lockReleaseTokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &LockReleaseTokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *LockReleaseTokenPoolContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *LockReleaseTokenPoolContract) Encoder() LockReleaseTokenPoolEncoder {
	return c.lockReleaseTokenPoolEncoder
}

func (c *LockReleaseTokenPoolContract) DevInspect() ILockReleaseTokenPoolDevInspect {
	return c.devInspect
}

type ReleasedOrMinted struct {
	RemoteChainSelector uint64 `move:"u64"`
	LocalToken          string `move:"address"`
	Recipient           string `move:"address"`
	Amount              uint64 `move:"u64"`
}

type RemotePoolAdded struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemotePoolAddress   []byte `move:"vector<u8>"`
}

type RemotePoolRemoved struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemotePoolAddress   []byte `move:"vector<u8>"`
}

type ChainAdded struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemoteTokenAddress  []byte `move:"vector<u8>"`
}

type ChainRemoved struct {
	RemoteChainSelector uint64 `move:"u64"`
}

type LiquidityAdded struct {
	LocalToken string `move:"address"`
	Provider   string `move:"address"`
	Amount     uint64 `move:"u64"`
}

type LiquidityRemoved struct {
	LocalToken string `move:"address"`
	Provider   string `move:"address"`
	Amount     uint64 `move:"u64"`
}

type RebalancerSet struct {
	LocalToken         string `move:"address"`
	PreviousRebalancer string `move:"address"`
	Rebalancer         string `move:"address"`
}

type TransferRequest struct {
	From     string `move:"address"`
	To       string `move:"address"`
	Accepted *bool  `move:"0x1::option::Option<bool>"`
}

type RateLimiter struct {
	OutboundIsEnabled bool   `move:"bool"`
	OutboundCapacity  uint64 `move:"u64"`
	OutboundRate      uint64 `move:"u64"`
	OutboundCurrent   uint64 `move:"u64"`
	OutboundLastReset uint64 `move:"u64"`
	InboundIsEnabled  bool   `move:"bool"`
	InboundCapacity   uint64 `move:"u64"`
	InboundRate       uint64 `move:"u64"`
	InboundCurrent    uint64 `move:"u64"`
	InboundLastReset  uint64 `move:"u64"`
}

type bcsReleasedOrMinted struct {
	RemoteChainSelector uint64
	LocalToken          [32]byte
	Recipient           [32]byte
	Amount              uint64
}

func convertReleasedOrMintedFromBCS(bcs bcsReleasedOrMinted) (ReleasedOrMinted, error) {

	return ReleasedOrMinted{
		RemoteChainSelector: bcs.RemoteChainSelector,
		LocalToken:          fmt.Sprintf("0x%x", bcs.LocalToken),
		Recipient:           fmt.Sprintf("0x%x", bcs.Recipient),
		Amount:              bcs.Amount,
	}, nil
}

type bcsLiquidityAdded struct {
	LocalToken [32]byte
	Provider   [32]byte
	Amount     uint64
}

func convertLiquidityAddedFromBCS(bcs bcsLiquidityAdded) (LiquidityAdded, error) {

	return LiquidityAdded{
		LocalToken: fmt.Sprintf("0x%x", bcs.LocalToken),
		Provider:   fmt.Sprintf("0x%x", bcs.Provider),
		Amount:     bcs.Amount,
	}, nil
}

type bcsLiquidityRemoved struct {
	LocalToken [32]byte
	Provider   [32]byte
	Amount     uint64
}

func convertLiquidityRemovedFromBCS(bcs bcsLiquidityRemoved) (LiquidityRemoved, error) {

	return LiquidityRemoved{
		LocalToken: fmt.Sprintf("0x%x", bcs.LocalToken),
		Provider:   fmt.Sprintf("0x%x", bcs.Provider),
		Amount:     bcs.Amount,
	}, nil
}

type bcsRebalancerSet struct {
	LocalToken         [32]byte
	PreviousRebalancer [32]byte
	Rebalancer         [32]byte
}

func convertRebalancerSetFromBCS(bcs bcsRebalancerSet) (RebalancerSet, error) {

	return RebalancerSet{
		LocalToken:         fmt.Sprintf("0x%x", bcs.LocalToken),
		PreviousRebalancer: fmt.Sprintf("0x%x", bcs.PreviousRebalancer),
		Rebalancer:         fmt.Sprintf("0x%x", bcs.Rebalancer),
	}, nil
}

type bcsTransferRequest struct {
	From     [32]byte
	To       [32]byte
	Accepted *bool
}

func convertTransferRequestFromBCS(bcs bcsTransferRequest) (TransferRequest, error) {

	return TransferRequest{
		From:     fmt.Sprintf("0x%x", bcs.From),
		To:       fmt.Sprintf("0x%x", bcs.To),
		Accepted: bcs.Accepted,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("test::lock_release_token_pool::ReleasedOrMinted", func(data []byte) (interface{}, error) {
		var temp bcsReleasedOrMinted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertReleasedOrMintedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::RemotePoolAdded", func(data []byte) (interface{}, error) {
		var result RemotePoolAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::RemotePoolRemoved", func(data []byte) (interface{}, error) {
		var result RemotePoolRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::ChainAdded", func(data []byte) (interface{}, error) {
		var result ChainAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::ChainRemoved", func(data []byte) (interface{}, error) {
		var result ChainRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::LiquidityAdded", func(data []byte) (interface{}, error) {
		var temp bcsLiquidityAdded
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertLiquidityAddedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::LiquidityRemoved", func(data []byte) (interface{}, error) {
		var temp bcsLiquidityRemoved
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertLiquidityRemovedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::RebalancerSet", func(data []byte) (interface{}, error) {
		var temp bcsRebalancerSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertRebalancerSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::TransferRequest", func(data []byte) (interface{}, error) {
		var temp bcsTransferRequest
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTransferRequestFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::lock_release_token_pool::RateLimiter", func(data []byte) (interface{}, error) {
		var result RateLimiter
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// EmitLockedOrBurnedEvent executes the emit_locked_or_burned_event Move function.
func (c *LockReleaseTokenPoolContract) EmitLockedOrBurnedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitLockedOrBurnedEvent(amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitReleasedOrMintedEvent executes the emit_released_or_minted_event Move function.
func (c *LockReleaseTokenPoolContract) EmitReleasedOrMintedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitReleasedOrMintedEvent(amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRemotePoolAddedEvent executes the emit_remote_pool_added_event Move function.
func (c *LockReleaseTokenPoolContract) EmitRemotePoolAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitRemotePoolAddedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRemotePoolRemovedEvent executes the emit_remote_pool_removed_event Move function.
func (c *LockReleaseTokenPoolContract) EmitRemotePoolRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitRemotePoolRemovedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitChainAddedEvent executes the emit_chain_added_event Move function.
func (c *LockReleaseTokenPoolContract) EmitChainAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitChainAddedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitChainRemovedEvent executes the emit_chain_removed_event Move function.
func (c *LockReleaseTokenPoolContract) EmitChainRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitChainRemovedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitLiquidityAddedEvent executes the emit_liquidity_added_event Move function.
func (c *LockReleaseTokenPoolContract) EmitLiquidityAddedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitLiquidityAddedEvent(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitLiquidityRemovedEvent executes the emit_liquidity_removed_event Move function.
func (c *LockReleaseTokenPoolContract) EmitLiquidityRemovedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitLiquidityRemovedEvent(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRebalancerSetEvent executes the emit_rebalancer_set_event Move function.
func (c *LockReleaseTokenPoolContract) EmitRebalancerSetEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.EmitRebalancerSetEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePendingTransferRequest executes the create_pending_transfer_request Move function.
func (c *LockReleaseTokenPoolContract) CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.CreatePendingTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateAcceptedTransferRequest executes the create_accepted_transfer_request Move function.
func (c *LockReleaseTokenPoolContract) CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.CreateAcceptedTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateRejectedTransferRequest executes the create_rejected_transfer_request Move function.
func (c *LockReleaseTokenPoolContract) CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.lockReleaseTokenPoolEncoder.CreateRejectedTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePendingTransferRequest executes the create_pending_transfer_request Move function using DevInspect to get return values.
//
// Returns: TransferRequest
func (d *LockReleaseTokenPoolDevInspect) CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.CreatePendingTransferRequest(from, to)
	if err != nil {
		return TransferRequest{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TransferRequest{}, err
	}
	if len(results) == 0 {
		return TransferRequest{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TransferRequest)
	if !ok {
		return TransferRequest{}, fmt.Errorf("unexpected return type: expected TransferRequest, got %T", results[0])
	}
	return result, nil
}

// CreateAcceptedTransferRequest executes the create_accepted_transfer_request Move function using DevInspect to get return values.
//
// Returns: TransferRequest
func (d *LockReleaseTokenPoolDevInspect) CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.CreateAcceptedTransferRequest(from, to)
	if err != nil {
		return TransferRequest{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TransferRequest{}, err
	}
	if len(results) == 0 {
		return TransferRequest{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TransferRequest)
	if !ok {
		return TransferRequest{}, fmt.Errorf("unexpected return type: expected TransferRequest, got %T", results[0])
	}
	return result, nil
}

// CreateRejectedTransferRequest executes the create_rejected_transfer_request Move function using DevInspect to get return values.
//
// Returns: TransferRequest
func (d *LockReleaseTokenPoolDevInspect) CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.lockReleaseTokenPoolEncoder.CreateRejectedTransferRequest(from, to)
	if err != nil {
		return TransferRequest{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TransferRequest{}, err
	}
	if len(results) == 0 {
		return TransferRequest{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TransferRequest)
	if !ok {
		return TransferRequest{}, fmt.Errorf("unexpected return type: expected TransferRequest, got %T", results[0])
	}
	return result, nil
}

type lockReleaseTokenPoolEncoder struct {
	*bind.BoundContract
}

// EmitLockedOrBurnedEvent encodes a call to the emit_locked_or_burned_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitLockedOrBurnedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_locked_or_burned_event", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
	}, []any{
		amount,
		remoteChainSelector,
	}, nil)
}

// EmitLockedOrBurnedEventWithArgs encodes a call to the emit_locked_or_burned_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitLockedOrBurnedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_locked_or_burned_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitReleasedOrMintedEvent encodes a call to the emit_released_or_minted_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitReleasedOrMintedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_released_or_minted_event", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
	}, []any{
		amount,
		remoteChainSelector,
	}, nil)
}

// EmitReleasedOrMintedEventWithArgs encodes a call to the emit_released_or_minted_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitReleasedOrMintedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_released_or_minted_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitRemotePoolAddedEvent encodes a call to the emit_remote_pool_added_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitRemotePoolAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_remote_pool_added_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		remoteChainSelector,
	}, nil)
}

// EmitRemotePoolAddedEventWithArgs encodes a call to the emit_remote_pool_added_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitRemotePoolAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_remote_pool_added_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitRemotePoolRemovedEvent encodes a call to the emit_remote_pool_removed_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitRemotePoolRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_remote_pool_removed_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		remoteChainSelector,
	}, nil)
}

// EmitRemotePoolRemovedEventWithArgs encodes a call to the emit_remote_pool_removed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitRemotePoolRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_remote_pool_removed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitChainAddedEvent encodes a call to the emit_chain_added_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitChainAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_chain_added_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		remoteChainSelector,
	}, nil)
}

// EmitChainAddedEventWithArgs encodes a call to the emit_chain_added_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitChainAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_chain_added_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitChainRemovedEvent encodes a call to the emit_chain_removed_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitChainRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_chain_removed_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		remoteChainSelector,
	}, nil)
}

// EmitChainRemovedEventWithArgs encodes a call to the emit_chain_removed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitChainRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_chain_removed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitLiquidityAddedEvent encodes a call to the emit_liquidity_added_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitLiquidityAddedEvent(amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_added_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		amount,
	}, nil)
}

// EmitLiquidityAddedEventWithArgs encodes a call to the emit_liquidity_added_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitLiquidityAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_added_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitLiquidityRemovedEvent encodes a call to the emit_liquidity_removed_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitLiquidityRemovedEvent(amount uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_removed_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		amount,
	}, nil)
}

// EmitLiquidityRemovedEventWithArgs encodes a call to the emit_liquidity_removed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitLiquidityRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_liquidity_removed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitRebalancerSetEvent encodes a call to the emit_rebalancer_set_event Move function.
func (c lockReleaseTokenPoolEncoder) EmitRebalancerSetEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_rebalancer_set_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitRebalancerSetEventWithArgs encodes a call to the emit_rebalancer_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) EmitRebalancerSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_rebalancer_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CreatePendingTransferRequest encodes a call to the create_pending_transfer_request Move function.
func (c lockReleaseTokenPoolEncoder) CreatePendingTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_pending_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::lock_release_token_pool::TransferRequest",
	})
}

// CreatePendingTransferRequestWithArgs encodes a call to the create_pending_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) CreatePendingTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_pending_transfer_request", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::lock_release_token_pool::TransferRequest",
	})
}

// CreateAcceptedTransferRequest encodes a call to the create_accepted_transfer_request Move function.
func (c lockReleaseTokenPoolEncoder) CreateAcceptedTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_accepted_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::lock_release_token_pool::TransferRequest",
	})
}

// CreateAcceptedTransferRequestWithArgs encodes a call to the create_accepted_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) CreateAcceptedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_accepted_transfer_request", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::lock_release_token_pool::TransferRequest",
	})
}

// CreateRejectedTransferRequest encodes a call to the create_rejected_transfer_request Move function.
func (c lockReleaseTokenPoolEncoder) CreateRejectedTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_rejected_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::lock_release_token_pool::TransferRequest",
	})
}

// CreateRejectedTransferRequestWithArgs encodes a call to the create_rejected_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c lockReleaseTokenPoolEncoder) CreateRejectedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_rejected_transfer_request", typeArgsList, typeParamsList, expectedParams, args, []string{
		"test::lock_release_token_pool::TransferRequest",
	})
}
