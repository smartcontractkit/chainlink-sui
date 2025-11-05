// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_dummy_receiver

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

const FunctionInfo = `[{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"ccip_receive","parameters":[{"name":"expected_message_id","type":"vector<u8>"},{"name":"ref","type":"CCIPObjectRef"},{"name":"message","type":"client::Any2SuiMessage"},{"name":"_","type":"Clock"},{"name":"state","type":"CCIPReceiverState"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"get_counter","parameters":[{"name":"state","type":"CCIPReceiverState"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"get_dest_token_amounts","parameters":[{"name":"state","type":"CCIPReceiverState"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"get_token_amount_amount","parameters":[{"name":"token_amount","type":"TokenAmount"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"get_token_amount_token","parameters":[{"name":"token_amount","type":"TokenAmount"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"get_token_receiver","parameters":[{"name":"state","type":"CCIPReceiverState"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"receive_and_send_coin","parameters":[{"name":"state","type":"CCIPReceiverState"},{"name":"_","type":"OwnerCap"},{"name":"coin_receiving","type":"Receiving<Coin<T>>"},{"name":"recipient","type":"address"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"receive_and_send_coin_no_owner_cap","parameters":[{"name":"state","type":"CCIPReceiverState"},{"name":"coin_receiving","type":"Receiving<Coin<T>>"},{"name":"recipient","type":"address"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"receive_coin","parameters":[{"name":"state","type":"CCIPReceiverState"},{"name":"_","type":"OwnerCap"},{"name":"coin_receiving","type":"Receiving<Coin<T>>"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"receive_coin_no_owner_cap","parameters":[{"name":"state","type":"CCIPReceiverState"},{"name":"coin_receiving","type":"Receiving<Coin<T>>"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"register_receiver","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"ref","type":"CCIPObjectRef"}]},{"package":"ccip_dummy_receiver","module":"dummy_receiver","name":"type_and_version","parameters":null}]`

type IDummyReceiver interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error)
	GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error)
	ReceiveAndSendCoin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error)
	ReceiveCoin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object) (*models.SuiTransactionBlockResponse, error)
	ReceiveAndSendCoinNoOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, coinReceiving bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error)
	ReceiveCoinNoOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, coinReceiving bind.Object) (*models.SuiTransactionBlockResponse, error)
	CcipReceive(ctx context.Context, opts *bind.CallOpts, expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IDummyReceiverDevInspect
	Encoder() DummyReceiverEncoder
	Bound() bind.IBoundContract
}

type IDummyReceiverDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (uint64, error)
	GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]TokenAmount, error)
	GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (string, error)
	GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*big.Int, error)
	ReceiveCoin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object) (any, error)
	ReceiveCoinNoOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, coinReceiving bind.Object) (any, error)
}

type DummyReceiverEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterReceiver(ownerCap bind.Object, ref bind.Object) (*bind.EncodedCall, error)
	RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCounter(state bind.Object) (*bind.EncodedCall, error)
	GetCounterWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestTokenAmounts(state bind.Object) (*bind.EncodedCall, error)
	GetDestTokenAmountsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenReceiver(state bind.Object) (*bind.EncodedCall, error)
	GetTokenReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenAmountToken(tokenAmount TokenAmount) (*bind.EncodedCall, error)
	GetTokenAmountTokenWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenAmountAmount(tokenAmount TokenAmount) (*bind.EncodedCall, error)
	GetTokenAmountAmountWithArgs(args ...any) (*bind.EncodedCall, error)
	ReceiveAndSendCoin(typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object, recipient string) (*bind.EncodedCall, error)
	ReceiveAndSendCoinWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReceiveCoin(typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object) (*bind.EncodedCall, error)
	ReceiveCoinWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReceiveAndSendCoinNoOwnerCap(typeArgs []string, state bind.Object, coinReceiving bind.Object, recipient string) (*bind.EncodedCall, error)
	ReceiveAndSendCoinNoOwnerCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ReceiveCoinNoOwnerCap(typeArgs []string, state bind.Object, coinReceiving bind.Object) (*bind.EncodedCall, error)
	ReceiveCoinNoOwnerCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	CcipReceive(expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*bind.EncodedCall, error)
	CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error)
}

type DummyReceiverContract struct {
	*bind.BoundContract
	dummyReceiverEncoder
	devInspect *DummyReceiverDevInspect
}

type DummyReceiverDevInspect struct {
	contract *DummyReceiverContract
}

var _ IDummyReceiver = (*DummyReceiverContract)(nil)
var _ IDummyReceiverDevInspect = (*DummyReceiverDevInspect)(nil)

func NewDummyReceiver(packageID string, client sui.ISuiAPI) (IDummyReceiver, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_dummy_receiver", "dummy_receiver", client)
	if err != nil {
		return nil, err
	}

	c := &DummyReceiverContract{
		BoundContract:        contract,
		dummyReceiverEncoder: dummyReceiverEncoder{BoundContract: contract},
	}
	c.devInspect = &DummyReceiverDevInspect{contract: c}
	return c, nil
}

func (c *DummyReceiverContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *DummyReceiverContract) Encoder() DummyReceiverEncoder {
	return c.dummyReceiverEncoder
}

func (c *DummyReceiverContract) DevInspect() IDummyReceiverDevInspect {
	return c.devInspect
}

type DUMMY_RECEIVER struct {
}

type OwnerCap struct {
	Id              string `move:"sui::object::UID"`
	ReceiverAddress string `move:"address"`
}

type ReceivedMessage struct {
	MessageId               []byte        `move:"vector<u8>"`
	SourceChainSelector     uint64        `move:"u64"`
	Sender                  []byte        `move:"vector<u8>"`
	Data                    []byte        `move:"vector<u8>"`
	DestTokenTransferLength uint64        `move:"u64"`
	DestTokenAmounts        []TokenAmount `move:"vector<TokenAmount>"`
}

type CCIPReceiverState struct {
	Id                      string        `move:"sui::object::UID"`
	Counter                 uint64        `move:"u64"`
	MessageId               []byte        `move:"vector<u8>"`
	SourceChainSelector     uint64        `move:"u64"`
	Sender                  []byte        `move:"vector<u8>"`
	Data                    []byte        `move:"vector<u8>"`
	MessageReceiver         string        `move:"address"`
	TokenReceiver           string        `move:"address"`
	DestTokenTransferLength uint64        `move:"u64"`
	DestTokenAmounts        []TokenAmount `move:"vector<TokenAmount>"`
}

type DummyReceiverProof struct {
}

type PublisherKey struct {
}

type TokenAmount struct {
	Token  string   `move:"address"`
	Amount *big.Int `move:"u256"`
}

type bcsOwnerCap struct {
	Id              string
	ReceiverAddress [32]byte
}

func convertOwnerCapFromBCS(bcs bcsOwnerCap) (OwnerCap, error) {

	return OwnerCap{
		Id:              bcs.Id,
		ReceiverAddress: fmt.Sprintf("0x%x", bcs.ReceiverAddress),
	}, nil
}

type bcsCCIPReceiverState struct {
	Id                      string
	Counter                 uint64
	MessageId               []byte
	SourceChainSelector     uint64
	Sender                  []byte
	Data                    []byte
	MessageReceiver         [32]byte
	TokenReceiver           [32]byte
	DestTokenTransferLength uint64
	DestTokenAmounts        []TokenAmount
}

func convertCCIPReceiverStateFromBCS(bcs bcsCCIPReceiverState) (CCIPReceiverState, error) {

	return CCIPReceiverState{
		Id:                      bcs.Id,
		Counter:                 bcs.Counter,
		MessageId:               bcs.MessageId,
		SourceChainSelector:     bcs.SourceChainSelector,
		Sender:                  bcs.Sender,
		Data:                    bcs.Data,
		MessageReceiver:         fmt.Sprintf("0x%x", bcs.MessageReceiver),
		TokenReceiver:           fmt.Sprintf("0x%x", bcs.TokenReceiver),
		DestTokenTransferLength: bcs.DestTokenTransferLength,
		DestTokenAmounts:        bcs.DestTokenAmounts,
	}, nil
}

type bcsTokenAmount struct {
	Token  [32]byte
	Amount [32]byte
}

func convertTokenAmountFromBCS(bcs bcsTokenAmount) (TokenAmount, error) {
	AmountField, err := bind.DecodeU256Value(bcs.Amount)
	if err != nil {
		return TokenAmount{}, fmt.Errorf("failed to decode u256 field Amount: %w", err)
	}

	return TokenAmount{
		Token:  fmt.Sprintf("0x%x", bcs.Token),
		Amount: AmountField,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::DUMMY_RECEIVER", func(data []byte) (interface{}, error) {
		var result DUMMY_RECEIVER
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DUMMY_RECEIVER
	bind.RegisterStructDecoder("vector<ccip_dummy_receiver::dummy_receiver::DUMMY_RECEIVER>", func(data []byte) (interface{}, error) {
		var results []DUMMY_RECEIVER
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::OwnerCap", func(data []byte) (interface{}, error) {
		var temp bcsOwnerCap
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOwnerCapFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OwnerCap
	bind.RegisterStructDecoder("vector<ccip_dummy_receiver::dummy_receiver::OwnerCap>", func(data []byte) (interface{}, error) {
		var temps []bcsOwnerCap
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]OwnerCap, len(temps))
		for i, temp := range temps {
			result, err := convertOwnerCapFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::ReceivedMessage", func(data []byte) (interface{}, error) {
		var result ReceivedMessage
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ReceivedMessage
	bind.RegisterStructDecoder("vector<ccip_dummy_receiver::dummy_receiver::ReceivedMessage>", func(data []byte) (interface{}, error) {
		var results []ReceivedMessage
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::CCIPReceiverState", func(data []byte) (interface{}, error) {
		var temp bcsCCIPReceiverState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCCIPReceiverStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for CCIPReceiverState
	bind.RegisterStructDecoder("vector<ccip_dummy_receiver::dummy_receiver::CCIPReceiverState>", func(data []byte) (interface{}, error) {
		var temps []bcsCCIPReceiverState
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]CCIPReceiverState, len(temps))
		for i, temp := range temps {
			result, err := convertCCIPReceiverStateFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::DummyReceiverProof", func(data []byte) (interface{}, error) {
		var result DummyReceiverProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DummyReceiverProof
	bind.RegisterStructDecoder("vector<ccip_dummy_receiver::dummy_receiver::DummyReceiverProof>", func(data []byte) (interface{}, error) {
		var results []DummyReceiverProof
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::PublisherKey", func(data []byte) (interface{}, error) {
		var result PublisherKey
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for PublisherKey
	bind.RegisterStructDecoder("vector<ccip_dummy_receiver::dummy_receiver::PublisherKey>", func(data []byte) (interface{}, error) {
		var results []PublisherKey
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip_dummy_receiver::dummy_receiver::TokenAmount", func(data []byte) (interface{}, error) {
		var temp bcsTokenAmount
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenAmountFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for TokenAmount
	bind.RegisterStructDecoder("vector<ccip_dummy_receiver::dummy_receiver::TokenAmount>", func(data []byte) (interface{}, error) {
		var temps []bcsTokenAmount
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]TokenAmount, len(temps))
		for i, temp := range temps {
			result, err := convertTokenAmountFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *DummyReceiverContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterReceiver executes the register_receiver Move function.
func (c *DummyReceiverContract) RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.RegisterReceiver(ownerCap, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCounter executes the get_counter Move function.
func (c *DummyReceiverContract) GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetCounter(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestTokenAmounts executes the get_dest_token_amounts Move function.
func (c *DummyReceiverContract) GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetDestTokenAmounts(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenReceiver executes the get_token_receiver Move function.
func (c *DummyReceiverContract) GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetTokenReceiver(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenAmountToken executes the get_token_amount_token Move function.
func (c *DummyReceiverContract) GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetTokenAmountToken(tokenAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenAmountAmount executes the get_token_amount_amount Move function.
func (c *DummyReceiverContract) GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.GetTokenAmountAmount(tokenAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReceiveAndSendCoin executes the receive_and_send_coin Move function.
func (c *DummyReceiverContract) ReceiveAndSendCoin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.ReceiveAndSendCoin(typeArgs, state, param, coinReceiving, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReceiveCoin executes the receive_coin Move function.
func (c *DummyReceiverContract) ReceiveCoin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.ReceiveCoin(typeArgs, state, param, coinReceiving)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReceiveAndSendCoinNoOwnerCap executes the receive_and_send_coin_no_owner_cap Move function.
func (c *DummyReceiverContract) ReceiveAndSendCoinNoOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, coinReceiving bind.Object, recipient string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.ReceiveAndSendCoinNoOwnerCap(typeArgs, state, coinReceiving, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ReceiveCoinNoOwnerCap executes the receive_coin_no_owner_cap Move function.
func (c *DummyReceiverContract) ReceiveCoinNoOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, coinReceiving bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.ReceiveCoinNoOwnerCap(typeArgs, state, coinReceiving)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CcipReceive executes the ccip_receive Move function.
func (c *DummyReceiverContract) CcipReceive(ctx context.Context, opts *bind.CallOpts, expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.dummyReceiverEncoder.CcipReceive(expectedMessageId, ref, message, param, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *DummyReceiverDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.dummyReceiverEncoder.TypeAndVersion()
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

// GetCounter executes the get_counter Move function using DevInspect to get return values.
//
// Returns: u64
func (d *DummyReceiverDevInspect) GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (uint64, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetCounter(state)
	if err != nil {
		return 0, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, fmt.Errorf("no return value")
	}
	result, ok := results[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected uint64, got %T", results[0])
	}
	return result, nil
}

// GetDestTokenAmounts executes the get_dest_token_amounts Move function using DevInspect to get return values.
//
// Returns: vector<TokenAmount>
func (d *DummyReceiverDevInspect) GetDestTokenAmounts(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]TokenAmount, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetDestTokenAmounts(state)
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
	result, ok := results[0].([]TokenAmount)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []TokenAmount, got %T", results[0])
	}
	return result, nil
}

// GetTokenReceiver executes the get_token_receiver Move function using DevInspect to get return values.
//
// Returns: address
func (d *DummyReceiverDevInspect) GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetTokenReceiver(state)
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

// GetTokenAmountToken executes the get_token_amount_token Move function using DevInspect to get return values.
//
// Returns: address
func (d *DummyReceiverDevInspect) GetTokenAmountToken(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (string, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetTokenAmountToken(tokenAmount)
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

// GetTokenAmountAmount executes the get_token_amount_amount Move function using DevInspect to get return values.
//
// Returns: u256
func (d *DummyReceiverDevInspect) GetTokenAmountAmount(ctx context.Context, opts *bind.CallOpts, tokenAmount TokenAmount) (*big.Int, error) {
	encoded, err := d.contract.dummyReceiverEncoder.GetTokenAmountAmount(tokenAmount)
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
	result, ok := results[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *big.Int, got %T", results[0])
	}
	return result, nil
}

// ReceiveCoin executes the receive_coin Move function using DevInspect to get return values.
//
// Returns: Coin<T>
func (d *DummyReceiverDevInspect) ReceiveCoin(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object) (any, error) {
	encoded, err := d.contract.dummyReceiverEncoder.ReceiveCoin(typeArgs, state, param, coinReceiving)
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

// ReceiveCoinNoOwnerCap executes the receive_coin_no_owner_cap Move function using DevInspect to get return values.
//
// Returns: Coin<T>
func (d *DummyReceiverDevInspect) ReceiveCoinNoOwnerCap(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, coinReceiving bind.Object) (any, error) {
	encoded, err := d.contract.dummyReceiverEncoder.ReceiveCoinNoOwnerCap(typeArgs, state, coinReceiving)
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

type dummyReceiverEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c dummyReceiverEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// RegisterReceiver encodes a call to the register_receiver Move function.
func (c dummyReceiverEncoder) RegisterReceiver(ownerCap bind.Object, ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_receiver", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
		"&mut CCIPObjectRef",
	}, []any{
		ownerCap,
		ref,
	}, nil)
}

// RegisterReceiverWithArgs encodes a call to the register_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
		"&mut CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("register_receiver", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCounter encodes a call to the get_counter Move function.
func (c dummyReceiverEncoder) GetCounter(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_counter", typeArgsList, typeParamsList, []string{
		"&CCIPReceiverState",
	}, []any{
		state,
	}, []string{
		"u64",
	})
}

// GetCounterWithArgs encodes a call to the get_counter Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetCounterWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_counter", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetDestTokenAmounts encodes a call to the get_dest_token_amounts Move function.
func (c dummyReceiverEncoder) GetDestTokenAmounts(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_token_amounts", typeArgsList, typeParamsList, []string{
		"&CCIPReceiverState",
	}, []any{
		state,
	}, []string{
		"vector<ccip_dummy_receiver::dummy_receiver::TokenAmount>",
	})
}

// GetDestTokenAmountsWithArgs encodes a call to the get_dest_token_amounts Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetDestTokenAmountsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_token_amounts", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ccip_dummy_receiver::dummy_receiver::TokenAmount>",
	})
}

// GetTokenReceiver encodes a call to the get_token_receiver Move function.
func (c dummyReceiverEncoder) GetTokenReceiver(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_receiver", typeArgsList, typeParamsList, []string{
		"&CCIPReceiverState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// GetTokenReceiverWithArgs encodes a call to the get_token_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetTokenReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_receiver", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetTokenAmountToken encodes a call to the get_token_amount_token Move function.
func (c dummyReceiverEncoder) GetTokenAmountToken(tokenAmount TokenAmount) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_token", typeArgsList, typeParamsList, []string{
		"&TokenAmount",
	}, []any{
		tokenAmount,
	}, []string{
		"address",
	})
}

// GetTokenAmountTokenWithArgs encodes a call to the get_token_amount_token Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetTokenAmountTokenWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenAmount",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_token", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// GetTokenAmountAmount encodes a call to the get_token_amount_amount Move function.
func (c dummyReceiverEncoder) GetTokenAmountAmount(tokenAmount TokenAmount) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_amount", typeArgsList, typeParamsList, []string{
		"&TokenAmount",
	}, []any{
		tokenAmount,
	}, []string{
		"u256",
	})
}

// GetTokenAmountAmountWithArgs encodes a call to the get_token_amount_amount Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) GetTokenAmountAmountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&TokenAmount",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_amount_amount", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u256",
	})
}

// ReceiveAndSendCoin encodes a call to the receive_and_send_coin Move function.
func (c dummyReceiverEncoder) ReceiveAndSendCoin(typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_and_send_coin", typeArgsList, typeParamsList, []string{
		"&mut CCIPReceiverState",
		"&OwnerCap",
		"Receiving<Coin<T>>",
		"address",
	}, []any{
		state,
		param,
		coinReceiving,
		recipient,
	}, nil)
}

// ReceiveAndSendCoinWithArgs encodes a call to the receive_and_send_coin Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) ReceiveAndSendCoinWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPReceiverState",
		"&OwnerCap",
		"Receiving<Coin<T>>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_and_send_coin", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ReceiveCoin encodes a call to the receive_coin Move function.
func (c dummyReceiverEncoder) ReceiveCoin(typeArgs []string, state bind.Object, param bind.Object, coinReceiving bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_coin", typeArgsList, typeParamsList, []string{
		"&mut CCIPReceiverState",
		"&OwnerCap",
		"Receiving<Coin<T>>",
	}, []any{
		state,
		param,
		coinReceiving,
	}, []string{
		"Coin<T>",
	})
}

// ReceiveCoinWithArgs encodes a call to the receive_coin Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) ReceiveCoinWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPReceiverState",
		"&OwnerCap",
		"Receiving<Coin<T>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_coin", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Coin<T>",
	})
}

// ReceiveAndSendCoinNoOwnerCap encodes a call to the receive_and_send_coin_no_owner_cap Move function.
func (c dummyReceiverEncoder) ReceiveAndSendCoinNoOwnerCap(typeArgs []string, state bind.Object, coinReceiving bind.Object, recipient string) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_and_send_coin_no_owner_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPReceiverState",
		"Receiving<Coin<T>>",
		"address",
	}, []any{
		state,
		coinReceiving,
		recipient,
	}, nil)
}

// ReceiveAndSendCoinNoOwnerCapWithArgs encodes a call to the receive_and_send_coin_no_owner_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) ReceiveAndSendCoinNoOwnerCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPReceiverState",
		"Receiving<Coin<T>>",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_and_send_coin_no_owner_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ReceiveCoinNoOwnerCap encodes a call to the receive_coin_no_owner_cap Move function.
func (c dummyReceiverEncoder) ReceiveCoinNoOwnerCap(typeArgs []string, state bind.Object, coinReceiving bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_coin_no_owner_cap", typeArgsList, typeParamsList, []string{
		"&mut CCIPReceiverState",
		"Receiving<Coin<T>>",
	}, []any{
		state,
		coinReceiving,
	}, []string{
		"Coin<T>",
	})
}

// ReceiveCoinNoOwnerCapWithArgs encodes a call to the receive_coin_no_owner_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) ReceiveCoinNoOwnerCapWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPReceiverState",
		"Receiving<Coin<T>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("receive_coin_no_owner_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Coin<T>",
	})
}

// CcipReceive encodes a call to the ccip_receive Move function.
func (c dummyReceiverEncoder) CcipReceive(expectedMessageId []byte, ref bind.Object, message bind.Object, param bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"&CCIPObjectRef",
		"client::Any2SuiMessage",
		"&Clock",
		"&mut CCIPReceiverState",
	}, []any{
		expectedMessageId,
		ref,
		message,
		param,
		state,
	}, nil)
}

// CcipReceiveWithArgs encodes a call to the ccip_receive Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c dummyReceiverEncoder) CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"&CCIPObjectRef",
		"client::Any2SuiMessage",
		"&Clock",
		"&mut CCIPReceiverState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, expectedParams, args, nil)
}
