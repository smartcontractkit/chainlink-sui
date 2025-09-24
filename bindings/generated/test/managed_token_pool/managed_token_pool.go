// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_managed_token_pool

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

type IManagedTokenPool interface {
	EmitTokenLockedOrBurnedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitChainRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitRemotePoolRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitTokenReleasedOrMintedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error)
	CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error)
	CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IManagedTokenPoolDevInspect
	Encoder() ManagedTokenPoolEncoder
	Bound() bind.IBoundContract
}

type IManagedTokenPoolDevInspect interface {
	CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
	CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
	CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
}

type ManagedTokenPoolEncoder interface {
	EmitTokenLockedOrBurnedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitTokenLockedOrBurnedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitChainAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitChainAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitChainRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitChainRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitRemotePoolAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitRemotePoolAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitRemotePoolRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitRemotePoolRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitTokenReleasedOrMintedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error)
	EmitTokenReleasedOrMintedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	CreatePendingTransferRequest(from string, to string) (*bind.EncodedCall, error)
	CreatePendingTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateAcceptedTransferRequest(from string, to string) (*bind.EncodedCall, error)
	CreateAcceptedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateRejectedTransferRequest(from string, to string) (*bind.EncodedCall, error)
	CreateRejectedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error)
}

type ManagedTokenPoolContract struct {
	*bind.BoundContract
	managedTokenPoolEncoder
	devInspect *ManagedTokenPoolDevInspect
}

type ManagedTokenPoolDevInspect struct {
	contract *ManagedTokenPoolContract
}

var _ IManagedTokenPool = (*ManagedTokenPoolContract)(nil)
var _ IManagedTokenPoolDevInspect = (*ManagedTokenPoolDevInspect)(nil)

func NewManagedTokenPool(packageID string, client sui.ISuiAPI) (IManagedTokenPool, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "managed_token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &ManagedTokenPoolContract{
		BoundContract:           contract,
		managedTokenPoolEncoder: managedTokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &ManagedTokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *ManagedTokenPoolContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *ManagedTokenPoolContract) Encoder() ManagedTokenPoolEncoder {
	return c.managedTokenPoolEncoder
}

func (c *ManagedTokenPoolContract) DevInspect() IManagedTokenPoolDevInspect {
	return c.devInspect
}

type TokenLockedOrBurned struct {
	Amount              uint64 `move:"u64"`
	RemoteChainSelector uint64 `move:"u64"`
	Token               string `move:"address"`
}

type TokenReleasedOrMinted struct {
	Receiver            string `move:"address"`
	Amount              uint64 `move:"u64"`
	RemoteChainSelector uint64 `move:"u64"`
}

type ChainAdded struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemoteTokenAddress  []byte `move:"vector<u8>"`
}

type ChainRemoved struct {
	RemoteChainSelector uint64 `move:"u64"`
}

type RemotePoolAdded struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemotePoolAddress   []byte `move:"vector<u8>"`
}

type RemotePoolRemoved struct {
	RemoteChainSelector uint64 `move:"u64"`
	RemotePoolAddress   []byte `move:"vector<u8>"`
}

type LiquidityAdded struct {
	LocalToken string `move:"address"`
	Provider   string `move:"address"`
	Amount     uint64 `move:"u64"`
}

type TokenPoolState struct {
	Id                   string      `move:"sui::object::UID"`
	Token                string      `move:"address"`
	LocalDecimals        byte        `move:"u8"`
	RemoteChainSelectors []uint64    `move:"vector<u64>"`
	RemotePools          bind.Object `move:"Table<u64, vector<vector<u8>>>"`
	RemoteTokens         bind.Object `move:"Table<u64, vector<u8>>"`
	AllowlistEnabled     bool        `move:"bool"`
	Allowlist            []string    `move:"vector<address>"`
	RateLimiters         bind.Object `move:"Table<u64, RateLimiter>"`
}

type MintCap struct {
	Id string `move:"sui::object::UID"`
}

type OwnableState struct {
	Id              string           `move:"sui::object::UID"`
	Owner           string           `move:"address"`
	PendingOwner    *string          `move:"0x1::option::Option<address>"`
	PendingTransfer *TransferRequest `move:"0x1::option::Option<TransferRequest>"`
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

type bcsTokenLockedOrBurned struct {
	Amount              uint64
	RemoteChainSelector uint64
	Token               [32]byte
}

func convertTokenLockedOrBurnedFromBCS(bcs bcsTokenLockedOrBurned) (TokenLockedOrBurned, error) {

	return TokenLockedOrBurned{
		Amount:              bcs.Amount,
		RemoteChainSelector: bcs.RemoteChainSelector,
		Token:               fmt.Sprintf("0x%x", bcs.Token),
	}, nil
}

type bcsTokenReleasedOrMinted struct {
	Receiver            [32]byte
	Amount              uint64
	RemoteChainSelector uint64
}

func convertTokenReleasedOrMintedFromBCS(bcs bcsTokenReleasedOrMinted) (TokenReleasedOrMinted, error) {

	return TokenReleasedOrMinted{
		Receiver:            fmt.Sprintf("0x%x", bcs.Receiver),
		Amount:              bcs.Amount,
		RemoteChainSelector: bcs.RemoteChainSelector,
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

type bcsTokenPoolState struct {
	Id                   string
	Token                [32]byte
	LocalDecimals        byte
	RemoteChainSelectors []uint64
	RemotePools          bind.Object
	RemoteTokens         bind.Object
	AllowlistEnabled     bool
	Allowlist            [][32]byte
	RateLimiters         bind.Object
}

func convertTokenPoolStateFromBCS(bcs bcsTokenPoolState) (TokenPoolState, error) {

	return TokenPoolState{
		Id:                   bcs.Id,
		Token:                fmt.Sprintf("0x%x", bcs.Token),
		LocalDecimals:        bcs.LocalDecimals,
		RemoteChainSelectors: bcs.RemoteChainSelectors,
		RemotePools:          bcs.RemotePools,
		RemoteTokens:         bcs.RemoteTokens,
		AllowlistEnabled:     bcs.AllowlistEnabled,
		Allowlist: func() []string {
			addrs := make([]string, len(bcs.Allowlist))
			for i, addr := range bcs.Allowlist {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		RateLimiters: bcs.RateLimiters,
	}, nil
}

type bcsOwnableState struct {
	Id              string
	Owner           [32]byte
	PendingOwner    *string
	PendingTransfer *TransferRequest
}

func convertOwnableStateFromBCS(bcs bcsOwnableState) (OwnableState, error) {

	return OwnableState{
		Id:              bcs.Id,
		Owner:           fmt.Sprintf("0x%x", bcs.Owner),
		PendingOwner:    bcs.PendingOwner,
		PendingTransfer: bcs.PendingTransfer,
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
	bind.RegisterStructDecoder("test::managed_token_pool::TokenLockedOrBurned", func(data []byte) (interface{}, error) {
		var temp bcsTokenLockedOrBurned
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenLockedOrBurnedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::TokenReleasedOrMinted", func(data []byte) (interface{}, error) {
		var temp bcsTokenReleasedOrMinted
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenReleasedOrMintedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::ChainAdded", func(data []byte) (interface{}, error) {
		var result ChainAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::ChainRemoved", func(data []byte) (interface{}, error) {
		var result ChainRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::RemotePoolAdded", func(data []byte) (interface{}, error) {
		var result RemotePoolAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::RemotePoolRemoved", func(data []byte) (interface{}, error) {
		var result RemotePoolRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::LiquidityAdded", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::managed_token_pool::TokenPoolState", func(data []byte) (interface{}, error) {
		var temp bcsTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenPoolStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::MintCap", func(data []byte) (interface{}, error) {
		var result MintCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::managed_token_pool::OwnableState", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::managed_token_pool::TransferRequest", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::managed_token_pool::RateLimiter", func(data []byte) (interface{}, error) {
		var result RateLimiter
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// EmitTokenLockedOrBurnedEvent executes the emit_token_locked_or_burned_event Move function.
func (c *ManagedTokenPoolContract) EmitTokenLockedOrBurnedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.EmitTokenLockedOrBurnedEvent(amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitChainAddedEvent executes the emit_chain_added_event Move function.
func (c *ManagedTokenPoolContract) EmitChainAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.EmitChainAddedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitChainRemovedEvent executes the emit_chain_removed_event Move function.
func (c *ManagedTokenPoolContract) EmitChainRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.EmitChainRemovedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRemotePoolAddedEvent executes the emit_remote_pool_added_event Move function.
func (c *ManagedTokenPoolContract) EmitRemotePoolAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.EmitRemotePoolAddedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRemotePoolRemovedEvent executes the emit_remote_pool_removed_event Move function.
func (c *ManagedTokenPoolContract) EmitRemotePoolRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.EmitRemotePoolRemovedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitTokenReleasedOrMintedEvent executes the emit_token_released_or_minted_event Move function.
func (c *ManagedTokenPoolContract) EmitTokenReleasedOrMintedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.EmitTokenReleasedOrMintedEvent(amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePendingTransferRequest executes the create_pending_transfer_request Move function.
func (c *ManagedTokenPoolContract) CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.CreatePendingTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateAcceptedTransferRequest executes the create_accepted_transfer_request Move function.
func (c *ManagedTokenPoolContract) CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.CreateAcceptedTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateRejectedTransferRequest executes the create_rejected_transfer_request Move function.
func (c *ManagedTokenPoolContract) CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.managedTokenPoolEncoder.CreateRejectedTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePendingTransferRequest executes the create_pending_transfer_request Move function using DevInspect to get return values.
//
// Returns: TransferRequest
func (d *ManagedTokenPoolDevInspect) CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.CreatePendingTransferRequest(from, to)
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
func (d *ManagedTokenPoolDevInspect) CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.CreateAcceptedTransferRequest(from, to)
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
func (d *ManagedTokenPoolDevInspect) CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.managedTokenPoolEncoder.CreateRejectedTransferRequest(from, to)
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

type managedTokenPoolEncoder struct {
	*bind.BoundContract
}

// EmitTokenLockedOrBurnedEvent encodes a call to the emit_token_locked_or_burned_event Move function.
func (c managedTokenPoolEncoder) EmitTokenLockedOrBurnedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_locked_or_burned_event", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
	}, []any{
		amount,
		remoteChainSelector,
	}, nil)
}

// EmitTokenLockedOrBurnedEventWithArgs encodes a call to the emit_token_locked_or_burned_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) EmitTokenLockedOrBurnedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_locked_or_burned_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitChainAddedEvent encodes a call to the emit_chain_added_event Move function.
func (c managedTokenPoolEncoder) EmitChainAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c managedTokenPoolEncoder) EmitChainAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c managedTokenPoolEncoder) EmitChainRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c managedTokenPoolEncoder) EmitChainRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// EmitRemotePoolAddedEvent encodes a call to the emit_remote_pool_added_event Move function.
func (c managedTokenPoolEncoder) EmitRemotePoolAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c managedTokenPoolEncoder) EmitRemotePoolAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c managedTokenPoolEncoder) EmitRemotePoolRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c managedTokenPoolEncoder) EmitRemotePoolRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// EmitTokenReleasedOrMintedEvent encodes a call to the emit_token_released_or_minted_event Move function.
func (c managedTokenPoolEncoder) EmitTokenReleasedOrMintedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_released_or_minted_event", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
	}, []any{
		amount,
		remoteChainSelector,
	}, nil)
}

// EmitTokenReleasedOrMintedEventWithArgs encodes a call to the emit_token_released_or_minted_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) EmitTokenReleasedOrMintedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_released_or_minted_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CreatePendingTransferRequest encodes a call to the create_pending_transfer_request Move function.
func (c managedTokenPoolEncoder) CreatePendingTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_pending_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::managed_token_pool::TransferRequest",
	})
}

// CreatePendingTransferRequestWithArgs encodes a call to the create_pending_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) CreatePendingTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
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
		"test::managed_token_pool::TransferRequest",
	})
}

// CreateAcceptedTransferRequest encodes a call to the create_accepted_transfer_request Move function.
func (c managedTokenPoolEncoder) CreateAcceptedTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_accepted_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::managed_token_pool::TransferRequest",
	})
}

// CreateAcceptedTransferRequestWithArgs encodes a call to the create_accepted_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) CreateAcceptedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
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
		"test::managed_token_pool::TransferRequest",
	})
}

// CreateRejectedTransferRequest encodes a call to the create_rejected_transfer_request Move function.
func (c managedTokenPoolEncoder) CreateRejectedTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_rejected_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::managed_token_pool::TransferRequest",
	})
}

// CreateRejectedTransferRequestWithArgs encodes a call to the create_rejected_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c managedTokenPoolEncoder) CreateRejectedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
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
		"test::managed_token_pool::TransferRequest",
	})
}
