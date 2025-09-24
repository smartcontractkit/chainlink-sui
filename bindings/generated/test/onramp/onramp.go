// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_onramp

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

type IOnramp interface {
	EmitConfigSetEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	EmitDestChainConfigSetEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	EmitCcipMessageSentEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error)
	EmitAllowlistSendersAddedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitAllowlistSendersRemovedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitFeeTokenWithdrawnEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOnrampDevInspect
	Encoder() OnrampEncoder
	Bound() bind.IBoundContract
}

type IOnrampDevInspect interface {
}

type OnrampEncoder interface {
	EmitConfigSetEvent() (*bind.EncodedCall, error)
	EmitConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitDestChainConfigSetEvent() (*bind.EncodedCall, error)
	EmitDestChainConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitCcipMessageSentEvent(destChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error)
	EmitCcipMessageSentEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitAllowlistSendersAddedEvent(destChainSelector uint64) (*bind.EncodedCall, error)
	EmitAllowlistSendersAddedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitAllowlistSendersRemovedEvent(destChainSelector uint64) (*bind.EncodedCall, error)
	EmitAllowlistSendersRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitFeeTokenWithdrawnEvent() (*bind.EncodedCall, error)
	EmitFeeTokenWithdrawnEventWithArgs(args ...any) (*bind.EncodedCall, error)
}

type OnrampContract struct {
	*bind.BoundContract
	onrampEncoder
	devInspect *OnrampDevInspect
}

type OnrampDevInspect struct {
	contract *OnrampContract
}

var _ IOnramp = (*OnrampContract)(nil)
var _ IOnrampDevInspect = (*OnrampDevInspect)(nil)

func NewOnramp(packageID string, client sui.ISuiAPI) (IOnramp, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "onramp", client)
	if err != nil {
		return nil, err
	}

	c := &OnrampContract{
		BoundContract: contract,
		onrampEncoder: onrampEncoder{BoundContract: contract},
	}
	c.devInspect = &OnrampDevInspect{contract: c}
	return c, nil
}

func (c *OnrampContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *OnrampContract) Encoder() OnrampEncoder {
	return c.onrampEncoder
}

func (c *OnrampContract) DevInspect() IOnrampDevInspect {
	return c.devInspect
}

type ConfigSet struct {
	StaticConfig  StaticConfig  `move:"StaticConfig"`
	DynamicConfig DynamicConfig `move:"DynamicConfig"`
}

type DestChainConfigSet struct {
	DestChainSelector uint64 `move:"u64"`
	IsEnabled         bool   `move:"bool"`
	SequenceNumber    uint64 `move:"u64"`
	AllowlistEnabled  bool   `move:"bool"`
}

type CCIPMessageSent struct {
	DestChainSelector uint64             `move:"u64"`
	SequenceNumber    uint64             `move:"u64"`
	Message           Sui2AnyRampMessage `move:"Sui2AnyRampMessage"`
}

type AllowlistSendersAdded struct {
	DestChainSelector uint64   `move:"u64"`
	Senders           []string `move:"vector<address>"`
}

type AllowlistSendersRemoved struct {
	DestChainSelector uint64   `move:"u64"`
	Senders           []string `move:"vector<address>"`
}

type FeeTokenWithdrawn struct {
	FeeAggregator string `move:"address"`
	FeeToken      string `move:"address"`
	Amount        uint64 `move:"u64"`
}

type OnRampState struct {
	Id             string `move:"sui::object::UID"`
	ChainSelector  uint64 `move:"u64"`
	FeeAggregator  string `move:"address"`
	AllowlistAdmin string `move:"address"`
}

type StaticConfig struct {
	ChainSelector uint64 `move:"u64"`
}

type DynamicConfig struct {
	FeeAggregator  string `move:"address"`
	AllowlistAdmin string `move:"address"`
}

type DestChainConfig struct {
	IsEnabled        bool     `move:"bool"`
	SequenceNumber   uint64   `move:"u64"`
	AllowlistEnabled bool     `move:"bool"`
	AllowedSenders   []string `move:"vector<address>"`
}

type RampMessageHeader struct {
	MessageId           []byte `move:"vector<u8>"`
	SourceChainSelector uint64 `move:"u64"`
	DestChainSelector   uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
	Nonce               uint64 `move:"u64"`
}

type Sui2AnyRampMessage struct {
	Header         RampMessageHeader      `move:"RampMessageHeader"`
	Sender         string                 `move:"address"`
	Data           []byte                 `move:"vector<u8>"`
	Receiver       []byte                 `move:"vector<u8>"`
	ExtraArgs      []byte                 `move:"vector<u8>"`
	FeeToken       string                 `move:"address"`
	FeeTokenAmount uint64                 `move:"u64"`
	FeeValueJuels  *big.Int               `move:"u256"`
	TokenAmounts   []Sui2AnyTokenTransfer `move:"vector<Sui2AnyTokenTransfer>"`
}

type Sui2AnyTokenTransfer struct {
	SourcePoolAddress string `move:"address"`
	DestTokenAddress  []byte `move:"vector<u8>"`
	ExtraData         []byte `move:"vector<u8>"`
	Amount            uint64 `move:"u64"`
	DestExecData      []byte `move:"vector<u8>"`
}

type bcsConfigSet struct {
	StaticConfig  StaticConfig
	DynamicConfig bcsDynamicConfig
}

func convertConfigSetFromBCS(bcs bcsConfigSet) (ConfigSet, error) {
	DynamicConfigField, err := convertDynamicConfigFromBCS(bcs.DynamicConfig)
	if err != nil {
		return ConfigSet{}, fmt.Errorf("failed to convert nested struct DynamicConfig: %w", err)
	}

	return ConfigSet{
		StaticConfig:  bcs.StaticConfig,
		DynamicConfig: DynamicConfigField,
	}, nil
}

type bcsCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           bcsSui2AnyRampMessage
}

func convertCCIPMessageSentFromBCS(bcs bcsCCIPMessageSent) (CCIPMessageSent, error) {
	MessageField, err := convertSui2AnyRampMessageFromBCS(bcs.Message)
	if err != nil {
		return CCIPMessageSent{}, fmt.Errorf("failed to convert nested struct Message: %w", err)
	}

	return CCIPMessageSent{
		DestChainSelector: bcs.DestChainSelector,
		SequenceNumber:    bcs.SequenceNumber,
		Message:           MessageField,
	}, nil
}

type bcsAllowlistSendersAdded struct {
	DestChainSelector uint64
	Senders           [][32]byte
}

func convertAllowlistSendersAddedFromBCS(bcs bcsAllowlistSendersAdded) (AllowlistSendersAdded, error) {

	return AllowlistSendersAdded{
		DestChainSelector: bcs.DestChainSelector,
		Senders: func() []string {
			addrs := make([]string, len(bcs.Senders))
			for i, addr := range bcs.Senders {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}, nil
}

type bcsAllowlistSendersRemoved struct {
	DestChainSelector uint64
	Senders           [][32]byte
}

func convertAllowlistSendersRemovedFromBCS(bcs bcsAllowlistSendersRemoved) (AllowlistSendersRemoved, error) {

	return AllowlistSendersRemoved{
		DestChainSelector: bcs.DestChainSelector,
		Senders: func() []string {
			addrs := make([]string, len(bcs.Senders))
			for i, addr := range bcs.Senders {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}, nil
}

type bcsFeeTokenWithdrawn struct {
	FeeAggregator [32]byte
	FeeToken      [32]byte
	Amount        uint64
}

func convertFeeTokenWithdrawnFromBCS(bcs bcsFeeTokenWithdrawn) (FeeTokenWithdrawn, error) {

	return FeeTokenWithdrawn{
		FeeAggregator: fmt.Sprintf("0x%x", bcs.FeeAggregator),
		FeeToken:      fmt.Sprintf("0x%x", bcs.FeeToken),
		Amount:        bcs.Amount,
	}, nil
}

type bcsOnRampState struct {
	Id             string
	ChainSelector  uint64
	FeeAggregator  [32]byte
	AllowlistAdmin [32]byte
}

func convertOnRampStateFromBCS(bcs bcsOnRampState) (OnRampState, error) {

	return OnRampState{
		Id:             bcs.Id,
		ChainSelector:  bcs.ChainSelector,
		FeeAggregator:  fmt.Sprintf("0x%x", bcs.FeeAggregator),
		AllowlistAdmin: fmt.Sprintf("0x%x", bcs.AllowlistAdmin),
	}, nil
}

type bcsDynamicConfig struct {
	FeeAggregator  [32]byte
	AllowlistAdmin [32]byte
}

func convertDynamicConfigFromBCS(bcs bcsDynamicConfig) (DynamicConfig, error) {

	return DynamicConfig{
		FeeAggregator:  fmt.Sprintf("0x%x", bcs.FeeAggregator),
		AllowlistAdmin: fmt.Sprintf("0x%x", bcs.AllowlistAdmin),
	}, nil
}

type bcsDestChainConfig struct {
	IsEnabled        bool
	SequenceNumber   uint64
	AllowlistEnabled bool
	AllowedSenders   [][32]byte
}

func convertDestChainConfigFromBCS(bcs bcsDestChainConfig) (DestChainConfig, error) {

	return DestChainConfig{
		IsEnabled:        bcs.IsEnabled,
		SequenceNumber:   bcs.SequenceNumber,
		AllowlistEnabled: bcs.AllowlistEnabled,
		AllowedSenders: func() []string {
			addrs := make([]string, len(bcs.AllowedSenders))
			for i, addr := range bcs.AllowedSenders {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}, nil
}

type bcsSui2AnyRampMessage struct {
	Header         RampMessageHeader
	Sender         [32]byte
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       [32]byte
	FeeTokenAmount uint64
	FeeValueJuels  [32]byte
	TokenAmounts   []Sui2AnyTokenTransfer
}

func convertSui2AnyRampMessageFromBCS(bcs bcsSui2AnyRampMessage) (Sui2AnyRampMessage, error) {
	FeeValueJuelsField, err := bind.DecodeU256Value(bcs.FeeValueJuels)
	if err != nil {
		return Sui2AnyRampMessage{}, fmt.Errorf("failed to decode u256 field FeeValueJuels: %w", err)
	}

	return Sui2AnyRampMessage{
		Header:         bcs.Header,
		Sender:         fmt.Sprintf("0x%x", bcs.Sender),
		Data:           bcs.Data,
		Receiver:       bcs.Receiver,
		ExtraArgs:      bcs.ExtraArgs,
		FeeToken:       fmt.Sprintf("0x%x", bcs.FeeToken),
		FeeTokenAmount: bcs.FeeTokenAmount,
		FeeValueJuels:  FeeValueJuelsField,
		TokenAmounts:   bcs.TokenAmounts,
	}, nil
}

type bcsSui2AnyTokenTransfer struct {
	SourcePoolAddress [32]byte
	DestTokenAddress  []byte
	ExtraData         []byte
	Amount            uint64
	DestExecData      []byte
}

func convertSui2AnyTokenTransferFromBCS(bcs bcsSui2AnyTokenTransfer) (Sui2AnyTokenTransfer, error) {

	return Sui2AnyTokenTransfer{
		SourcePoolAddress: fmt.Sprintf("0x%x", bcs.SourcePoolAddress),
		DestTokenAddress:  bcs.DestTokenAddress,
		ExtraData:         bcs.ExtraData,
		Amount:            bcs.Amount,
		DestExecData:      bcs.DestExecData,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("test::onramp::ConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertConfigSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::DestChainConfigSet", func(data []byte) (interface{}, error) {
		var result DestChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::CCIPMessageSent", func(data []byte) (interface{}, error) {
		var temp bcsCCIPMessageSent
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCCIPMessageSentFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::AllowlistSendersAdded", func(data []byte) (interface{}, error) {
		var temp bcsAllowlistSendersAdded
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAllowlistSendersAddedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::AllowlistSendersRemoved", func(data []byte) (interface{}, error) {
		var temp bcsAllowlistSendersRemoved
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAllowlistSendersRemovedFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::FeeTokenWithdrawn", func(data []byte) (interface{}, error) {
		var temp bcsFeeTokenWithdrawn
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFeeTokenWithdrawnFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::OnRampState", func(data []byte) (interface{}, error) {
		var temp bcsOnRampState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOnRampStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::StaticConfig", func(data []byte) (interface{}, error) {
		var result StaticConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::DynamicConfig", func(data []byte) (interface{}, error) {
		var temp bcsDynamicConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDynamicConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::DestChainConfig", func(data []byte) (interface{}, error) {
		var temp bcsDestChainConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDestChainConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::RampMessageHeader", func(data []byte) (interface{}, error) {
		var result RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::Sui2AnyRampMessage", func(data []byte) (interface{}, error) {
		var temp bcsSui2AnyRampMessage
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSui2AnyRampMessageFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::onramp::Sui2AnyTokenTransfer", func(data []byte) (interface{}, error) {
		var temp bcsSui2AnyTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSui2AnyTokenTransferFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// EmitConfigSetEvent executes the emit_config_set_event Move function.
func (c *OnrampContract) EmitConfigSetEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.EmitConfigSetEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitDestChainConfigSetEvent executes the emit_dest_chain_config_set_event Move function.
func (c *OnrampContract) EmitDestChainConfigSetEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.EmitDestChainConfigSetEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitCcipMessageSentEvent executes the emit_ccip_message_sent_event Move function.
func (c *OnrampContract) EmitCcipMessageSentEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.EmitCcipMessageSentEvent(destChainSelector, sequenceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitAllowlistSendersAddedEvent executes the emit_allowlist_senders_added_event Move function.
func (c *OnrampContract) EmitAllowlistSendersAddedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.EmitAllowlistSendersAddedEvent(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitAllowlistSendersRemovedEvent executes the emit_allowlist_senders_removed_event Move function.
func (c *OnrampContract) EmitAllowlistSendersRemovedEvent(ctx context.Context, opts *bind.CallOpts, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.EmitAllowlistSendersRemovedEvent(destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitFeeTokenWithdrawnEvent executes the emit_fee_token_withdrawn_event Move function.
func (c *OnrampContract) EmitFeeTokenWithdrawnEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.EmitFeeTokenWithdrawnEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type onrampEncoder struct {
	*bind.BoundContract
}

// EmitConfigSetEvent encodes a call to the emit_config_set_event Move function.
func (c onrampEncoder) EmitConfigSetEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_config_set_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitConfigSetEventWithArgs encodes a call to the emit_config_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) EmitConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_config_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitDestChainConfigSetEvent encodes a call to the emit_dest_chain_config_set_event Move function.
func (c onrampEncoder) EmitDestChainConfigSetEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dest_chain_config_set_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitDestChainConfigSetEventWithArgs encodes a call to the emit_dest_chain_config_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) EmitDestChainConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dest_chain_config_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitCcipMessageSentEvent encodes a call to the emit_ccip_message_sent_event Move function.
func (c onrampEncoder) EmitCcipMessageSentEvent(destChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_ccip_message_sent_event", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
	}, []any{
		destChainSelector,
		sequenceNumber,
	}, nil)
}

// EmitCcipMessageSentEventWithArgs encodes a call to the emit_ccip_message_sent_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) EmitCcipMessageSentEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_ccip_message_sent_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitAllowlistSendersAddedEvent encodes a call to the emit_allowlist_senders_added_event Move function.
func (c onrampEncoder) EmitAllowlistSendersAddedEvent(destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_allowlist_senders_added_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		destChainSelector,
	}, nil)
}

// EmitAllowlistSendersAddedEventWithArgs encodes a call to the emit_allowlist_senders_added_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) EmitAllowlistSendersAddedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_allowlist_senders_added_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitAllowlistSendersRemovedEvent encodes a call to the emit_allowlist_senders_removed_event Move function.
func (c onrampEncoder) EmitAllowlistSendersRemovedEvent(destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_allowlist_senders_removed_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		destChainSelector,
	}, nil)
}

// EmitAllowlistSendersRemovedEventWithArgs encodes a call to the emit_allowlist_senders_removed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) EmitAllowlistSendersRemovedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_allowlist_senders_removed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitFeeTokenWithdrawnEvent encodes a call to the emit_fee_token_withdrawn_event Move function.
func (c onrampEncoder) EmitFeeTokenWithdrawnEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_fee_token_withdrawn_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitFeeTokenWithdrawnEventWithArgs encodes a call to the emit_fee_token_withdrawn_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) EmitFeeTokenWithdrawnEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_fee_token_withdrawn_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}
