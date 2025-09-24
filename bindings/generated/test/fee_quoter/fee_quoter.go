// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_fee_quoter

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

type IFeeQuoter interface {
	EmitFeeTokenAddedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	EmitFeeTokenRemovedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	EmitTokenTransferFeeConfigAddedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitTokenTransferFeeConfigRemovedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitUsdPerTokenUpdatedEvent(ctx context.Context, opts *bind.CallOpts, clock bind.Object) (*models.SuiTransactionBlockResponse, error)
	EmitUsdPerUnitGasUpdatedEvent(ctx context.Context, opts *bind.CallOpts, clock bind.Object) (*models.SuiTransactionBlockResponse, error)
	EmitDestChainAddedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitDestChainConfigUpdatedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitPremiumMultiplierWeiPerEthUpdatedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IFeeQuoterDevInspect
	Encoder() FeeQuoterEncoder
	Bound() bind.IBoundContract
}

type IFeeQuoterDevInspect interface {
}

type FeeQuoterEncoder interface {
	EmitFeeTokenAddedEvent() (*bind.EncodedCall, error)
	EmitFeeTokenAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitFeeTokenRemovedEvent() (*bind.EncodedCall, error)
	EmitFeeTokenRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitTokenTransferFeeConfigAddedEvent(destChainSelector uint64) (*bind.EncodedCall, error)
	EmitTokenTransferFeeConfigAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitTokenTransferFeeConfigRemovedEvent(destChainSelector uint64) (*bind.EncodedCall, error)
	EmitTokenTransferFeeConfigRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitUsdPerTokenUpdatedEvent(clock bind.Object) (*bind.EncodedCall, error)
	EmitUsdPerTokenUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitUsdPerUnitGasUpdatedEvent(clock bind.Object) (*bind.EncodedCall, error)
	EmitUsdPerUnitGasUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitDestChainAddedEvent(destChainSelector uint64) (*bind.EncodedCall, error)
	EmitDestChainAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitDestChainConfigUpdatedEvent(destChainSelector uint64) (*bind.EncodedCall, error)
	EmitDestChainConfigUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitPremiumMultiplierWeiPerEthUpdatedEvent() (*bind.EncodedCall, error)
	EmitPremiumMultiplierWeiPerEthUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error)
}

type FeeQuoterContract struct {
	*bind.BoundContract
	feeQuoterEncoder
	devInspect *FeeQuoterDevInspect
}

type FeeQuoterDevInspect struct {
	contract *FeeQuoterContract
}

var _ IFeeQuoter = (*FeeQuoterContract)(nil)
var _ IFeeQuoterDevInspect = (*FeeQuoterDevInspect)(nil)

func NewFeeQuoter(packageID string, client sui.ISuiAPI) (IFeeQuoter, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "fee_quoter", client)
	if err != nil {
		return nil, err
	}

	c := &FeeQuoterContract{
		BoundContract:    contract,
		feeQuoterEncoder: feeQuoterEncoder{BoundContract: contract},
	}
	c.devInspect = &FeeQuoterDevInspect{contract: c}
	return c, nil
}

func (c *FeeQuoterContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *FeeQuoterContract) Encoder() FeeQuoterEncoder {
	return c.feeQuoterEncoder
}

func (c *FeeQuoterContract) DevInspect() IFeeQuoterDevInspect {
	return c.devInspect
}

type FeeTokenAdded struct {
	FeeToken string `move:"address"`
}

type FeeTokenRemoved struct {
	FeeToken string `move:"address"`
}

type TokenTransferFeeConfigAdded struct {
	DestChainSelector      uint64                 `move:"u64"`
	Token                  string                 `move:"address"`
	TokenTransferFeeConfig TokenTransferFeeConfig `move:"TokenTransferFeeConfig"`
}

type TokenTransferFeeConfigRemoved struct {
	DestChainSelector uint64 `move:"u64"`
	Token             string `move:"address"`
}

type UsdPerTokenUpdated struct {
	Token       string   `move:"address"`
	UsdPerToken *big.Int `move:"u256"`
	Timestamp   uint64   `move:"u64"`
}

type UsdPerUnitGasUpdated struct {
	DestChainSelector uint64   `move:"u64"`
	UsdPerUnitGas     *big.Int `move:"u256"`
	Timestamp         uint64   `move:"u64"`
}

type DestChainAdded struct {
	DestChainSelector uint64          `move:"u64"`
	DestChainConfig   DestChainConfig `move:"DestChainConfig"`
}

type DestChainConfigUpdated struct {
	DestChainSelector uint64          `move:"u64"`
	DestChainConfig   DestChainConfig `move:"DestChainConfig"`
}

type PremiumMultiplierWeiPerEthUpdated struct {
	Token                      string `move:"address"`
	PremiumMultiplierWeiPerEth uint64 `move:"u64"`
}

type TokenTransferFeeConfig struct {
	MinFeeUsdCents    uint32 `move:"u32"`
	MaxFeeUsdCents    uint32 `move:"u32"`
	DeciBps           uint16 `move:"u16"`
	DestGasOverhead   uint32 `move:"u32"`
	DestBytesOverhead uint32 `move:"u32"`
	IsEnabled         bool   `move:"bool"`
}

type DestChainConfig struct {
	IsEnabled                         bool   `move:"bool"`
	MaxNumberOfTokensPerMsg           uint16 `move:"u16"`
	MaxDataBytes                      uint32 `move:"u32"`
	MaxPerMsgGasLimit                 uint32 `move:"u32"`
	DestGasOverhead                   uint32 `move:"u32"`
	DestGasPerPayloadByteBase         byte   `move:"u8"`
	DestGasPerPayloadByteHigh         byte   `move:"u8"`
	DestGasPerPayloadByteThreshold    uint16 `move:"u16"`
	DestDataAvailabilityOverheadGas   uint32 `move:"u32"`
	DestGasPerDataAvailabilityByte    uint16 `move:"u16"`
	DestDataAvailabilityMultiplierBps uint16 `move:"u16"`
	ChainFamilySelector               []byte `move:"vector<u8>"`
	EnforceOutOfOrder                 bool   `move:"bool"`
	DefaultTokenFeeUsdCents           uint16 `move:"u16"`
	DefaultTokenDestGasOverhead       uint32 `move:"u32"`
	DefaultTxGasLimit                 uint32 `move:"u32"`
	GasMultiplierWeiPerEth            uint64 `move:"u64"`
	GasPriceStalenessThreshold        uint32 `move:"u32"`
	NetworkFeeUsdCents                uint32 `move:"u32"`
}

type bcsFeeTokenAdded struct {
	FeeToken [32]byte
}

func convertFeeTokenAddedFromBCS(bcs bcsFeeTokenAdded) (FeeTokenAdded, error) {

	return FeeTokenAdded{
		FeeToken: fmt.Sprintf("0x%x", bcs.FeeToken),
	}, nil
}

type bcsFeeTokenRemoved struct {
	FeeToken [32]byte
}

func convertFeeTokenRemovedFromBCS(bcs bcsFeeTokenRemoved) (FeeTokenRemoved, error) {

	return FeeTokenRemoved{
		FeeToken: fmt.Sprintf("0x%x", bcs.FeeToken),
	}, nil
}

type bcsTokenTransferFeeConfigAdded struct {
	DestChainSelector      uint64
	Token                  [32]byte
	TokenTransferFeeConfig TokenTransferFeeConfig
}

func convertTokenTransferFeeConfigAddedFromBCS(bcs bcsTokenTransferFeeConfigAdded) (TokenTransferFeeConfigAdded, error) {

	return TokenTransferFeeConfigAdded{
		DestChainSelector:      bcs.DestChainSelector,
		Token:                  fmt.Sprintf("0x%x", bcs.Token),
		TokenTransferFeeConfig: bcs.TokenTransferFeeConfig,
	}, nil
}

type bcsTokenTransferFeeConfigRemoved struct {
	DestChainSelector uint64
	Token             [32]byte
}

func convertTokenTransferFeeConfigRemovedFromBCS(bcs bcsTokenTransferFeeConfigRemoved) (TokenTransferFeeConfigRemoved, error) {

	return TokenTransferFeeConfigRemoved{
		DestChainSelector: bcs.DestChainSelector,
		Token:             fmt.Sprintf("0x%x", bcs.Token),
	}, nil
}

type bcsUsdPerTokenUpdated struct {
	Token       [32]byte
	UsdPerToken [32]byte
	Timestamp   uint64
}

func convertUsdPerTokenUpdatedFromBCS(bcs bcsUsdPerTokenUpdated) (UsdPerTokenUpdated, error) {
	UsdPerTokenField, err := bind.DecodeU256Value(bcs.UsdPerToken)
	if err != nil {
		return UsdPerTokenUpdated{}, fmt.Errorf("failed to decode u256 field UsdPerToken: %w", err)
	}

	return UsdPerTokenUpdated{
		Token:       fmt.Sprintf("0x%x", bcs.Token),
		UsdPerToken: UsdPerTokenField,
		Timestamp:   bcs.Timestamp,
	}, nil
}

type bcsUsdPerUnitGasUpdated struct {
	DestChainSelector uint64
	UsdPerUnitGas     [32]byte
	Timestamp         uint64
}

func convertUsdPerUnitGasUpdatedFromBCS(bcs bcsUsdPerUnitGasUpdated) (UsdPerUnitGasUpdated, error) {
	UsdPerUnitGasField, err := bind.DecodeU256Value(bcs.UsdPerUnitGas)
	if err != nil {
		return UsdPerUnitGasUpdated{}, fmt.Errorf("failed to decode u256 field UsdPerUnitGas: %w", err)
	}

	return UsdPerUnitGasUpdated{
		DestChainSelector: bcs.DestChainSelector,
		UsdPerUnitGas:     UsdPerUnitGasField,
		Timestamp:         bcs.Timestamp,
	}, nil
}

type bcsPremiumMultiplierWeiPerEthUpdated struct {
	Token                      [32]byte
	PremiumMultiplierWeiPerEth uint64
}

func convertPremiumMultiplierWeiPerEthUpdatedFromBCS(bcs bcsPremiumMultiplierWeiPerEthUpdated) (PremiumMultiplierWeiPerEthUpdated, error) {

	return PremiumMultiplierWeiPerEthUpdated{
		Token:                      fmt.Sprintf("0x%x", bcs.Token),
		PremiumMultiplierWeiPerEth: bcs.PremiumMultiplierWeiPerEth,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("test::fee_quoter::FeeTokenAdded", func(data []byte) (interface{}, error) {
		var temp bcsFeeTokenAdded
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFeeTokenAddedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::FeeTokenRemoved", func(data []byte) (interface{}, error) {
		var temp bcsFeeTokenRemoved
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFeeTokenRemovedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::TokenTransferFeeConfigAdded", func(data []byte) (interface{}, error) {
		var temp bcsTokenTransferFeeConfigAdded
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenTransferFeeConfigAddedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::TokenTransferFeeConfigRemoved", func(data []byte) (interface{}, error) {
		var temp bcsTokenTransferFeeConfigRemoved
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenTransferFeeConfigRemovedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::UsdPerTokenUpdated", func(data []byte) (interface{}, error) {
		var temp bcsUsdPerTokenUpdated
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertUsdPerTokenUpdatedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::UsdPerUnitGasUpdated", func(data []byte) (interface{}, error) {
		var temp bcsUsdPerUnitGasUpdated
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertUsdPerUnitGasUpdatedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::DestChainAdded", func(data []byte) (interface{}, error) {
		var result DestChainAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::DestChainConfigUpdated", func(data []byte) (interface{}, error) {
		var result DestChainConfigUpdated
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::PremiumMultiplierWeiPerEthUpdated", func(data []byte) (interface{}, error) {
		var temp bcsPremiumMultiplierWeiPerEthUpdated
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertPremiumMultiplierWeiPerEthUpdatedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::TokenTransferFeeConfig", func(data []byte) (interface{}, error) {
		var result TokenTransferFeeConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::fee_quoter::DestChainConfig", func(data []byte) (interface{}, error) {
		var result DestChainConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// EmitFeeTokenAddedEvent executes the emit_fee_token_added_event Move function.
func (c *FeeQuoterContract) EmitFeeTokenAddedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitFeeTokenAddedEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitFeeTokenRemovedEvent executes the emit_fee_token_removed_event Move function.
func (c *FeeQuoterContract) EmitFeeTokenRemovedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitFeeTokenRemovedEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitTokenTransferFeeConfigAddedEvent executes the emit_token_transfer_fee_config_added_event Move function.
func (c *FeeQuoterContract) EmitTokenTransferFeeConfigAddedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitTokenTransferFeeConfigAddedEvent(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitTokenTransferFeeConfigRemovedEvent executes the emit_token_transfer_fee_config_removed_event Move function.
func (c *FeeQuoterContract) EmitTokenTransferFeeConfigRemovedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitTokenTransferFeeConfigRemovedEvent(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitUsdPerTokenUpdatedEvent executes the emit_usd_per_token_updated_event Move function.
func (c *FeeQuoterContract) EmitUsdPerTokenUpdatedEvent(ctx context.Context, opts *bind.CallOpts, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitUsdPerTokenUpdatedEvent(clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitUsdPerUnitGasUpdatedEvent executes the emit_usd_per_unit_gas_updated_event Move function.
func (c *FeeQuoterContract) EmitUsdPerUnitGasUpdatedEvent(ctx context.Context, opts *bind.CallOpts, clock bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitUsdPerUnitGasUpdatedEvent(clock)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitDestChainAddedEvent executes the emit_dest_chain_added_event Move function.
func (c *FeeQuoterContract) EmitDestChainAddedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitDestChainAddedEvent(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitDestChainConfigUpdatedEvent executes the emit_dest_chain_config_updated_event Move function.
func (c *FeeQuoterContract) EmitDestChainConfigUpdatedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitDestChainConfigUpdatedEvent(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitPremiumMultiplierWeiPerEthUpdatedEvent executes the emit_premium_multiplier_wei_per_eth_updated_event Move function.
func (c *FeeQuoterContract) EmitPremiumMultiplierWeiPerEthUpdatedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.EmitPremiumMultiplierWeiPerEthUpdatedEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type feeQuoterEncoder struct {
	*bind.BoundContract
}

// EmitFeeTokenAddedEvent encodes a call to the emit_fee_token_added_event Move function.
func (c feeQuoterEncoder) EmitFeeTokenAddedEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_fee_token_added_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitFeeTokenAddedEventWithArgs encodes a call to the emit_fee_token_added_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitFeeTokenAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_fee_token_added_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitFeeTokenRemovedEvent encodes a call to the emit_fee_token_removed_event Move function.
func (c feeQuoterEncoder) EmitFeeTokenRemovedEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_fee_token_removed_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitFeeTokenRemovedEventWithArgs encodes a call to the emit_fee_token_removed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitFeeTokenRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_fee_token_removed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitTokenTransferFeeConfigAddedEvent encodes a call to the emit_token_transfer_fee_config_added_event Move function.
func (c feeQuoterEncoder) EmitTokenTransferFeeConfigAddedEvent(destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_transfer_fee_config_added_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		destChainSelector,
	}, nil)
}

// EmitTokenTransferFeeConfigAddedEventWithArgs encodes a call to the emit_token_transfer_fee_config_added_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitTokenTransferFeeConfigAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_transfer_fee_config_added_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitTokenTransferFeeConfigRemovedEvent encodes a call to the emit_token_transfer_fee_config_removed_event Move function.
func (c feeQuoterEncoder) EmitTokenTransferFeeConfigRemovedEvent(destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_transfer_fee_config_removed_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		destChainSelector,
	}, nil)
}

// EmitTokenTransferFeeConfigRemovedEventWithArgs encodes a call to the emit_token_transfer_fee_config_removed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitTokenTransferFeeConfigRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_token_transfer_fee_config_removed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitUsdPerTokenUpdatedEvent encodes a call to the emit_usd_per_token_updated_event Move function.
func (c feeQuoterEncoder) EmitUsdPerTokenUpdatedEvent(clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_usd_per_token_updated_event", typeArgsList, typeParamsList, []string{
		"&clock::Clock",
	}, []any{
		clock,
	}, nil)
}

// EmitUsdPerTokenUpdatedEventWithArgs encodes a call to the emit_usd_per_token_updated_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitUsdPerTokenUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&clock::Clock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_usd_per_token_updated_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitUsdPerUnitGasUpdatedEvent encodes a call to the emit_usd_per_unit_gas_updated_event Move function.
func (c feeQuoterEncoder) EmitUsdPerUnitGasUpdatedEvent(clock bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_usd_per_unit_gas_updated_event", typeArgsList, typeParamsList, []string{
		"&clock::Clock",
	}, []any{
		clock,
	}, nil)
}

// EmitUsdPerUnitGasUpdatedEventWithArgs encodes a call to the emit_usd_per_unit_gas_updated_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitUsdPerUnitGasUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&clock::Clock",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_usd_per_unit_gas_updated_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitDestChainAddedEvent encodes a call to the emit_dest_chain_added_event Move function.
func (c feeQuoterEncoder) EmitDestChainAddedEvent(destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dest_chain_added_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		destChainSelector,
	}, nil)
}

// EmitDestChainAddedEventWithArgs encodes a call to the emit_dest_chain_added_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitDestChainAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dest_chain_added_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitDestChainConfigUpdatedEvent encodes a call to the emit_dest_chain_config_updated_event Move function.
func (c feeQuoterEncoder) EmitDestChainConfigUpdatedEvent(destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dest_chain_config_updated_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		destChainSelector,
	}, nil)
}

// EmitDestChainConfigUpdatedEventWithArgs encodes a call to the emit_dest_chain_config_updated_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitDestChainConfigUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dest_chain_config_updated_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitPremiumMultiplierWeiPerEthUpdatedEvent encodes a call to the emit_premium_multiplier_wei_per_eth_updated_event Move function.
func (c feeQuoterEncoder) EmitPremiumMultiplierWeiPerEthUpdatedEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_premium_multiplier_wei_per_eth_updated_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitPremiumMultiplierWeiPerEthUpdatedEventWithArgs encodes a call to the emit_premium_multiplier_wei_per_eth_updated_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) EmitPremiumMultiplierWeiPerEthUpdatedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_premium_multiplier_wei_per_eth_updated_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}
