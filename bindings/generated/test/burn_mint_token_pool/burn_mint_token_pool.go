// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_burn_mint_token_pool

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

type IBurnMintTokenPool interface {
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
	DevInspect() IBurnMintTokenPoolDevInspect
	Encoder() BurnMintTokenPoolEncoder
	Bound() bind.IBoundContract
}

type IBurnMintTokenPoolDevInspect interface {
	CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
	CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
	CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error)
}

type BurnMintTokenPoolEncoder interface {
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

type BurnMintTokenPoolContract struct {
	*bind.BoundContract
	burnMintTokenPoolEncoder
	devInspect *BurnMintTokenPoolDevInspect
}

type BurnMintTokenPoolDevInspect struct {
	contract *BurnMintTokenPoolContract
}

var _ IBurnMintTokenPool = (*BurnMintTokenPoolContract)(nil)
var _ IBurnMintTokenPoolDevInspect = (*BurnMintTokenPoolDevInspect)(nil)

func NewBurnMintTokenPool(packageID string, client sui.ISuiAPI) (IBurnMintTokenPool, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "burn_mint_token_pool", client)
	if err != nil {
		return nil, err
	}

	c := &BurnMintTokenPoolContract{
		BoundContract:            contract,
		burnMintTokenPoolEncoder: burnMintTokenPoolEncoder{BoundContract: contract},
	}
	c.devInspect = &BurnMintTokenPoolDevInspect{contract: c}
	return c, nil
}

func (c *BurnMintTokenPoolContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *BurnMintTokenPoolContract) Encoder() BurnMintTokenPoolEncoder {
	return c.burnMintTokenPoolEncoder
}

func (c *BurnMintTokenPoolContract) DevInspect() IBurnMintTokenPoolDevInspect {
	return c.devInspect
}

type LockedOrBurned struct {
	RemoteChainSelector uint64 `move:"u64"`
	LocalToken          string `move:"address"`
	Amount              uint64 `move:"u64"`
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

type BurnMintTokenPoolState struct {
	Id             string      `move:"sui::object::UID"`
	TokenPoolState bind.Object `move:"TokenPoolState"`
	TreasuryCap    bind.Object `move:"TreasuryCap<T>"`
	OwnableState   bind.Object `move:"OwnableState"`
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

type TreasuryCap struct {
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

type bcsLockedOrBurned struct {
	RemoteChainSelector uint64
	LocalToken          [32]byte
	Amount              uint64
}

func convertLockedOrBurnedFromBCS(bcs bcsLockedOrBurned) (LockedOrBurned, error) {

	return LockedOrBurned{
		RemoteChainSelector: bcs.RemoteChainSelector,
		LocalToken:          fmt.Sprintf("0x%x", bcs.LocalToken),
		Amount:              bcs.Amount,
	}, nil
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

type bcsBurnMintTokenPoolState struct {
	Id             string
	TokenPoolState bcsTokenPoolState
	TreasuryCap    bind.Object
	OwnableState   bcsOwnableState
}

func convertBurnMintTokenPoolStateFromBCS(bcs bcsBurnMintTokenPoolState) (BurnMintTokenPoolState, error) {
	TokenPoolStateField, err := convertTokenPoolStateFromBCS(bcs.TokenPoolState)
	if err != nil {
		return BurnMintTokenPoolState{}, fmt.Errorf("failed to convert nested struct TokenPoolState: %w", err)
	}
	OwnableStateField, err := convertOwnableStateFromBCS(bcs.OwnableState)
	if err != nil {
		return BurnMintTokenPoolState{}, fmt.Errorf("failed to convert nested struct OwnableState: %w", err)
	}

	return BurnMintTokenPoolState{
		Id:             bcs.Id,
		TokenPoolState: TokenPoolStateField,
		TreasuryCap:    bcs.TreasuryCap,
		OwnableState:   OwnableStateField,
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::LockedOrBurned", func(data []byte) (interface{}, error) {
		var temp bcsLockedOrBurned
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertLockedOrBurnedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::burn_mint_token_pool::ReleasedOrMinted", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::RemotePoolAdded", func(data []byte) (interface{}, error) {
		var result RemotePoolAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::burn_mint_token_pool::RemotePoolRemoved", func(data []byte) (interface{}, error) {
		var result RemotePoolRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::burn_mint_token_pool::ChainAdded", func(data []byte) (interface{}, error) {
		var result ChainAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::burn_mint_token_pool::ChainRemoved", func(data []byte) (interface{}, error) {
		var result ChainRemoved
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::burn_mint_token_pool::LiquidityAdded", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::LiquidityRemoved", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::RebalancerSet", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::BurnMintTokenPoolState", func(data []byte) (interface{}, error) {
		var temp bcsBurnMintTokenPoolState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertBurnMintTokenPoolStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::burn_mint_token_pool::TokenPoolState", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::TreasuryCap", func(data []byte) (interface{}, error) {
		var result TreasuryCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::burn_mint_token_pool::OwnableState", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::TransferRequest", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::burn_mint_token_pool::RateLimiter", func(data []byte) (interface{}, error) {
		var result RateLimiter
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// EmitLockedOrBurnedEvent executes the emit_locked_or_burned_event Move function.
func (c *BurnMintTokenPoolContract) EmitLockedOrBurnedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitLockedOrBurnedEvent(amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitReleasedOrMintedEvent executes the emit_released_or_minted_event Move function.
func (c *BurnMintTokenPoolContract) EmitReleasedOrMintedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitReleasedOrMintedEvent(amount, remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRemotePoolAddedEvent executes the emit_remote_pool_added_event Move function.
func (c *BurnMintTokenPoolContract) EmitRemotePoolAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitRemotePoolAddedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRemotePoolRemovedEvent executes the emit_remote_pool_removed_event Move function.
func (c *BurnMintTokenPoolContract) EmitRemotePoolRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitRemotePoolRemovedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitChainAddedEvent executes the emit_chain_added_event Move function.
func (c *BurnMintTokenPoolContract) EmitChainAddedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitChainAddedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitChainRemovedEvent executes the emit_chain_removed_event Move function.
func (c *BurnMintTokenPoolContract) EmitChainRemovedEvent(ctx context.Context, opts *bind.CallOpts, remoteChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitChainRemovedEvent(remoteChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitLiquidityAddedEvent executes the emit_liquidity_added_event Move function.
func (c *BurnMintTokenPoolContract) EmitLiquidityAddedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitLiquidityAddedEvent(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitLiquidityRemovedEvent executes the emit_liquidity_removed_event Move function.
func (c *BurnMintTokenPoolContract) EmitLiquidityRemovedEvent(ctx context.Context, opts *bind.CallOpts, amount uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitLiquidityRemovedEvent(amount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitRebalancerSetEvent executes the emit_rebalancer_set_event Move function.
func (c *BurnMintTokenPoolContract) EmitRebalancerSetEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.EmitRebalancerSetEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePendingTransferRequest executes the create_pending_transfer_request Move function.
func (c *BurnMintTokenPoolContract) CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.CreatePendingTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateAcceptedTransferRequest executes the create_accepted_transfer_request Move function.
func (c *BurnMintTokenPoolContract) CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.CreateAcceptedTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateRejectedTransferRequest executes the create_rejected_transfer_request Move function.
func (c *BurnMintTokenPoolContract) CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.burnMintTokenPoolEncoder.CreateRejectedTransferRequest(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreatePendingTransferRequest executes the create_pending_transfer_request Move function using DevInspect to get return values.
//
// Returns: TransferRequest
func (d *BurnMintTokenPoolDevInspect) CreatePendingTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.CreatePendingTransferRequest(from, to)
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
func (d *BurnMintTokenPoolDevInspect) CreateAcceptedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.CreateAcceptedTransferRequest(from, to)
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
func (d *BurnMintTokenPoolDevInspect) CreateRejectedTransferRequest(ctx context.Context, opts *bind.CallOpts, from string, to string) (TransferRequest, error) {
	encoded, err := d.contract.burnMintTokenPoolEncoder.CreateRejectedTransferRequest(from, to)
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

type burnMintTokenPoolEncoder struct {
	*bind.BoundContract
}

// EmitLockedOrBurnedEvent encodes a call to the emit_locked_or_burned_event Move function.
func (c burnMintTokenPoolEncoder) EmitLockedOrBurnedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitLockedOrBurnedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitReleasedOrMintedEvent(amount uint64, remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitReleasedOrMintedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitRemotePoolAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitRemotePoolAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitRemotePoolRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitRemotePoolRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitChainAddedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitChainAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitChainRemovedEvent(remoteChainSelector uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitChainRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitLiquidityAddedEvent(amount uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitLiquidityAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitLiquidityRemovedEvent(amount uint64) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitLiquidityRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c burnMintTokenPoolEncoder) EmitRebalancerSetEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_rebalancer_set_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitRebalancerSetEventWithArgs encodes a call to the emit_rebalancer_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) EmitRebalancerSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_rebalancer_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// CreatePendingTransferRequest encodes a call to the create_pending_transfer_request Move function.
func (c burnMintTokenPoolEncoder) CreatePendingTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_pending_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::burn_mint_token_pool::TransferRequest",
	})
}

// CreatePendingTransferRequestWithArgs encodes a call to the create_pending_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) CreatePendingTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
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
		"test::burn_mint_token_pool::TransferRequest",
	})
}

// CreateAcceptedTransferRequest encodes a call to the create_accepted_transfer_request Move function.
func (c burnMintTokenPoolEncoder) CreateAcceptedTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_accepted_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::burn_mint_token_pool::TransferRequest",
	})
}

// CreateAcceptedTransferRequestWithArgs encodes a call to the create_accepted_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) CreateAcceptedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
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
		"test::burn_mint_token_pool::TransferRequest",
	})
}

// CreateRejectedTransferRequest encodes a call to the create_rejected_transfer_request Move function.
func (c burnMintTokenPoolEncoder) CreateRejectedTransferRequest(from string, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_rejected_transfer_request", typeArgsList, typeParamsList, []string{
		"address",
		"address",
	}, []any{
		from,
		to,
	}, []string{
		"test::burn_mint_token_pool::TransferRequest",
	})
}

// CreateRejectedTransferRequestWithArgs encodes a call to the create_rejected_transfer_request Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c burnMintTokenPoolEncoder) CreateRejectedTransferRequestWithArgs(args ...any) (*bind.EncodedCall, error) {
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
		"test::burn_mint_token_pool::TransferRequest",
	})
}
