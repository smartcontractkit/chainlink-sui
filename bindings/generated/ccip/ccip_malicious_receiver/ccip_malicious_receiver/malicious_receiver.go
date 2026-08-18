// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_malicious_receiver

import (
	"context"
	"fmt"
	"math/big"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

var (
	_ = big.NewInt
)

const FunctionInfo = `[{"package":"ccip_malicious_receiver","module":"malicious_receiver","name":"ccip_receive","parameters":[{"name":"expected_message_id","type":"vector<u8>"},{"name":"ref","type":"CCIPObjectRef"},{"name":"message","type":"client::Any2SuiMessage"},{"name":"clock","type":"Clock"},{"name":"state","type":"CCIPReceiverState"},{"name":"drain_coin","type":"Coin<sui::sui::SUI>"}]},{"package":"ccip_malicious_receiver","module":"malicious_receiver","name":"get_counter","parameters":[{"name":"state","type":"CCIPReceiverState"}]},{"package":"ccip_malicious_receiver","module":"malicious_receiver","name":"register_receiver","parameters":[{"name":"owner_cap","type":"OwnerCap"},{"name":"ref","type":"CCIPObjectRef"}]},{"package":"ccip_malicious_receiver","module":"malicious_receiver","name":"type_and_version","parameters":null}]`

type IMaliciousReceiver interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	CcipReceive(ctx context.Context, opts *bind.CallOpts, expectedMessageId []byte, ref bind.Object, message bind.Object, clock bind.Object, state bind.Object, drainCoin bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IMaliciousReceiverDevInspect
	Encoder() MaliciousReceiverEncoder
	Bound() bind.IBoundContract
}

type IMaliciousReceiverDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (uint64, error)
}

type MaliciousReceiverEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	RegisterReceiver(ownerCap bind.Object, ref bind.Object) (*bind.EncodedCall, error)
	RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCounter(state bind.Object) (*bind.EncodedCall, error)
	GetCounterWithArgs(args ...any) (*bind.EncodedCall, error)
	CcipReceive(expectedMessageId []byte, ref bind.Object, message bind.Object, clock bind.Object, state bind.Object, drainCoin bind.Object) (*bind.EncodedCall, error)
	CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error)
}

type MaliciousReceiverContract struct {
	*bind.BoundContract
	maliciousReceiverEncoder
	devInspect *MaliciousReceiverDevInspect
}

type MaliciousReceiverDevInspect struct {
	contract *MaliciousReceiverContract
}

var _ IMaliciousReceiver = (*MaliciousReceiverContract)(nil)
var _ IMaliciousReceiverDevInspect = (*MaliciousReceiverDevInspect)(nil)

func NewMaliciousReceiver(packageID string, chainClient client.BindingsClient) (IMaliciousReceiver, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_malicious_receiver", "malicious_receiver", chainClient)
	if err != nil {
		return nil, err
	}

	c := &MaliciousReceiverContract{
		BoundContract:            contract,
		maliciousReceiverEncoder: maliciousReceiverEncoder{BoundContract: contract},
	}
	c.devInspect = &MaliciousReceiverDevInspect{contract: c}
	return c, nil
}

func (c *MaliciousReceiverContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *MaliciousReceiverContract) Encoder() MaliciousReceiverEncoder {
	return c.maliciousReceiverEncoder
}

func (c *MaliciousReceiverContract) DevInspect() IMaliciousReceiverDevInspect {
	return c.devInspect
}

type MALICIOUS_RECEIVER struct {
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

type MaliciousReceiverProof struct {
}

type PublisherKey struct {
}

type TokenAmount struct {
	Token  string   `move:"address"`
	Amount *big.Int `move:"u256"`
}

// TypeAndVersion executes the type_and_version Move function.
func (c *MaliciousReceiverContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.maliciousReceiverEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RegisterReceiver executes the register_receiver Move function.
func (c *MaliciousReceiverContract) RegisterReceiver(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.maliciousReceiverEncoder.RegisterReceiver(ownerCap, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCounter executes the get_counter Move function.
func (c *MaliciousReceiverContract) GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.maliciousReceiverEncoder.GetCounter(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CcipReceive executes the ccip_receive Move function.
func (c *MaliciousReceiverContract) CcipReceive(ctx context.Context, opts *bind.CallOpts, expectedMessageId []byte, ref bind.Object, message bind.Object, clock bind.Object, state bind.Object, drainCoin bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.maliciousReceiverEncoder.CcipReceive(expectedMessageId, ref, message, clock, state, drainCoin)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *MaliciousReceiverDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.maliciousReceiverEncoder.TypeAndVersion()
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
	var result string
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return "", fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

// GetCounter executes the get_counter Move function using DevInspect to get return values.
//
// Returns: u64
func (d *MaliciousReceiverDevInspect) GetCounter(ctx context.Context, opts *bind.CallOpts, state bind.Object) (uint64, error) {
	encoded, err := d.contract.maliciousReceiverEncoder.GetCounter(state)
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
	var result uint64
	if err := bind.DecodeJSONReturn(results[0], &result); err != nil {
		return 0, fmt.Errorf("failed to decode return value: %w", err)
	}
	return result, nil
}

type maliciousReceiverEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c maliciousReceiverEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c maliciousReceiverEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c maliciousReceiverEncoder) RegisterReceiver(ownerCap bind.Object, ref bind.Object) (*bind.EncodedCall, error) {
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
func (c maliciousReceiverEncoder) RegisterReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c maliciousReceiverEncoder) GetCounter(state bind.Object) (*bind.EncodedCall, error) {
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
func (c maliciousReceiverEncoder) GetCounterWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// CcipReceive encodes a call to the ccip_receive Move function.
func (c maliciousReceiverEncoder) CcipReceive(expectedMessageId []byte, ref bind.Object, message bind.Object, clock bind.Object, state bind.Object, drainCoin bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"&CCIPObjectRef",
		"client::Any2SuiMessage",
		"&Clock",
		"&mut CCIPReceiverState",
		"&mut Coin<sui::sui::SUI>",
	}, []any{
		expectedMessageId,
		ref,
		message,
		clock,
		state,
		drainCoin,
	}, nil)
}

// CcipReceiveWithArgs encodes a call to the ccip_receive Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c maliciousReceiverEncoder) CcipReceiveWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"&CCIPObjectRef",
		"client::Any2SuiMessage",
		"&Clock",
		"&mut CCIPReceiverState",
		"&mut Coin<sui::sui::SUI>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("ccip_receive", typeArgsList, typeParamsList, expectedParams, args, nil)
}
