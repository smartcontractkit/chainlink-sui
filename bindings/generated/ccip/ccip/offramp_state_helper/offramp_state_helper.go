// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_offramp_state_helper

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

const FunctionInfo = `[{"package":"ccip","module":"offramp_state_helper","name":"add_dest_token_transfer","parameters":[{"name":"_","type":"DestTransferCap"},{"name":"receiver_params","type":"ReceiverParams"},{"name":"token_receiver","type":"address"},{"name":"remote_chain_selector","type":"u64"},{"name":"source_amount","type":"u256"},{"name":"dest_token_address","type":"address"},{"name":"dest_token_pool_package_id","type":"address"},{"name":"source_pool_address","type":"vector<u8>"},{"name":"source_pool_data","type":"vector<u8>"},{"name":"offchain_data","type":"vector<u8>"}]},{"package":"ccip","module":"offramp_state_helper","name":"complete_token_transfer","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"receiver_params","type":"ReceiverParams"},{"name":"_","type":"TypeProof"}]},{"package":"ccip","module":"offramp_state_helper","name":"consume_any2sui_message","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"message","type":"Any2SuiMessage"},{"name":"_","type":"TypeProof"}]},{"package":"ccip","module":"offramp_state_helper","name":"create_receiver_params","parameters":[{"name":"_","type":"DestTransferCap"},{"name":"source_chain_selector","type":"u64"}]},{"package":"ccip","module":"offramp_state_helper","name":"deconstruct_receiver_params","parameters":[{"name":"_","type":"DestTransferCap"},{"name":"receiver_params","type":"ReceiverParams"}]},{"package":"ccip","module":"offramp_state_helper","name":"extract_any2sui_message","parameters":[{"name":"receiver_params","type":"ReceiverParams"}]},{"package":"ccip","module":"offramp_state_helper","name":"get_dest_token_transfer_data","parameters":[{"name":"receiver_params","type":"ReceiverParams"}]},{"package":"ccip","module":"offramp_state_helper","name":"get_source_chain_selector","parameters":[{"name":"receiver_params","type":"ReceiverParams"}]},{"package":"ccip","module":"offramp_state_helper","name":"get_token_param_data","parameters":[{"name":"receiver_params","type":"ReceiverParams"}]},{"package":"ccip","module":"offramp_state_helper","name":"new_any2sui_message","parameters":[{"name":"_","type":"DestTransferCap"},{"name":"message_id","type":"vector<u8>"},{"name":"source_chain_selector","type":"u64"},{"name":"sender","type":"vector<u8>"},{"name":"data","type":"vector<u8>"},{"name":"message_receiver","type":"address"},{"name":"token_receiver","type":"address"},{"name":"token_addresses","type":"vector<address>"},{"name":"token_amounts","type":"vector<u256>"}]},{"package":"ccip","module":"offramp_state_helper","name":"new_dest_transfer_cap","parameters":[{"name":"ref","type":"CCIPObjectRef"},{"name":"owner_cap","type":"OwnerCap"}]},{"package":"ccip","module":"offramp_state_helper","name":"populate_message","parameters":[{"name":"_","type":"DestTransferCap"},{"name":"receiver_params","type":"ReceiverParams"},{"name":"any2sui_message","type":"Any2SuiMessage"}]}]`

type IOfframpStateHelper interface {
	NewDestTransferCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	CreateReceiverParams(ctx context.Context, opts *bind.CallOpts, param bind.Object, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetSourceChainSelector(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error)
	AddDestTokenTransfer(ctx context.Context, opts *bind.CallOpts, param bind.Object, receiverParams ReceiverParams, tokenReceiver string, remoteChainSelector uint64, sourceAmount *big.Int, destTokenAddress string, destTokenPoolPackageId string, sourcePoolAddress []byte, sourcePoolData []byte, offchainData []byte) (*models.SuiTransactionBlockResponse, error)
	PopulateMessage(ctx context.Context, opts *bind.CallOpts, param bind.Object, receiverParams ReceiverParams, any2suiMessage bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDestTokenTransferData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error)
	GetTokenParamData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error)
	CompleteTokenTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams ReceiverParams, param bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExtractAny2suiMessage(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error)
	NewAny2suiMessage(ctx context.Context, opts *bind.CallOpts, param bind.Object, messageId []byte, sourceChainSelector uint64, sender []byte, data []byte, messageReceiver string, tokenReceiver string, tokenAddresses []string, tokenAmounts []*big.Int) (*models.SuiTransactionBlockResponse, error)
	ConsumeAny2suiMessage(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, message bind.Object, param bind.Object) (*models.SuiTransactionBlockResponse, error)
	DeconstructReceiverParams(ctx context.Context, opts *bind.CallOpts, param bind.Object, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOfframpStateHelperDevInspect
	Encoder() OfframpStateHelperEncoder
	Bound() bind.IBoundContract
}

type IOfframpStateHelperDevInspect interface {
	NewDestTransferCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (bind.Object, error)
	CreateReceiverParams(ctx context.Context, opts *bind.CallOpts, param bind.Object, sourceChainSelector uint64) (ReceiverParams, error)
	GetSourceChainSelector(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (uint64, error)
	GetDestTokenTransferData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) ([]any, error)
	GetTokenParamData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) ([]any, error)
	ExtractAny2suiMessage(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (bind.Object, error)
	NewAny2suiMessage(ctx context.Context, opts *bind.CallOpts, param bind.Object, messageId []byte, sourceChainSelector uint64, sender []byte, data []byte, messageReceiver string, tokenReceiver string, tokenAddresses []string, tokenAmounts []*big.Int) (bind.Object, error)
	ConsumeAny2suiMessage(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, message bind.Object, param bind.Object) ([]any, error)
}

type OfframpStateHelperEncoder interface {
	NewDestTransferCap(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error)
	NewDestTransferCapWithArgs(args ...any) (*bind.EncodedCall, error)
	CreateReceiverParams(param bind.Object, sourceChainSelector uint64) (*bind.EncodedCall, error)
	CreateReceiverParamsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetSourceChainSelector(receiverParams ReceiverParams) (*bind.EncodedCall, error)
	GetSourceChainSelectorWithArgs(args ...any) (*bind.EncodedCall, error)
	AddDestTokenTransfer(param bind.Object, receiverParams ReceiverParams, tokenReceiver string, remoteChainSelector uint64, sourceAmount *big.Int, destTokenAddress string, destTokenPoolPackageId string, sourcePoolAddress []byte, sourcePoolData []byte, offchainData []byte) (*bind.EncodedCall, error)
	AddDestTokenTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PopulateMessage(param bind.Object, receiverParams ReceiverParams, any2suiMessage bind.Object) (*bind.EncodedCall, error)
	PopulateMessageWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestTokenTransferData(receiverParams ReceiverParams) (*bind.EncodedCall, error)
	GetDestTokenTransferDataWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenParamData(receiverParams ReceiverParams) (*bind.EncodedCall, error)
	GetTokenParamDataWithArgs(args ...any) (*bind.EncodedCall, error)
	CompleteTokenTransfer(typeArgs []string, ref bind.Object, receiverParams ReceiverParams, param bind.Object) (*bind.EncodedCall, error)
	CompleteTokenTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	ExtractAny2suiMessage(receiverParams ReceiverParams) (*bind.EncodedCall, error)
	ExtractAny2suiMessageWithArgs(args ...any) (*bind.EncodedCall, error)
	NewAny2suiMessage(param bind.Object, messageId []byte, sourceChainSelector uint64, sender []byte, data []byte, messageReceiver string, tokenReceiver string, tokenAddresses []string, tokenAmounts []*big.Int) (*bind.EncodedCall, error)
	NewAny2suiMessageWithArgs(args ...any) (*bind.EncodedCall, error)
	ConsumeAny2suiMessage(typeArgs []string, ref bind.Object, message bind.Object, param bind.Object) (*bind.EncodedCall, error)
	ConsumeAny2suiMessageWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	DeconstructReceiverParams(param bind.Object, receiverParams ReceiverParams) (*bind.EncodedCall, error)
	DeconstructReceiverParamsWithArgs(args ...any) (*bind.EncodedCall, error)
}

type OfframpStateHelperContract struct {
	*bind.BoundContract
	offrampStateHelperEncoder
	devInspect *OfframpStateHelperDevInspect
}

type OfframpStateHelperDevInspect struct {
	contract *OfframpStateHelperContract
}

var _ IOfframpStateHelper = (*OfframpStateHelperContract)(nil)
var _ IOfframpStateHelperDevInspect = (*OfframpStateHelperDevInspect)(nil)

func NewOfframpStateHelper(packageID string, client sui.ISuiAPI) (IOfframpStateHelper, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "offramp_state_helper", client)
	if err != nil {
		return nil, err
	}

	c := &OfframpStateHelperContract{
		BoundContract:             contract,
		offrampStateHelperEncoder: offrampStateHelperEncoder{BoundContract: contract},
	}
	c.devInspect = &OfframpStateHelperDevInspect{contract: c}
	return c, nil
}

func (c *OfframpStateHelperContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *OfframpStateHelperContract) Encoder() OfframpStateHelperEncoder {
	return c.offrampStateHelperEncoder
}

func (c *OfframpStateHelperContract) DevInspect() IOfframpStateHelperDevInspect {
	return c.devInspect
}

type OFFRAMP_STATE_HELPER struct {
}

type ReceiverParams struct {
	TokenTransfer       *DestTokenTransfer          `move:"0x1::option::Option<DestTokenTransfer>"`
	Message             *bind.Object                `move:"0x1::option::Option<Any2SuiMessage>"`
	SourceChainSelector uint64                      `move:"u64"`
	Receipt             *CompletedDestTokenTransfer `move:"0x1::option::Option<CompletedDestTokenTransfer>"`
}

type DestTransferCap struct {
	Id string `move:"sui::object::UID"`
}

type CompletedDestTokenTransfer struct {
	TokenReceiver    string `move:"address"`
	DestTokenAddress string `move:"address"`
}

type DestTokenTransfer struct {
	TokenReceiver          string   `move:"address"`
	RemoteChainSelector    uint64   `move:"u64"`
	SourceAmount           *big.Int `move:"u256"`
	DestTokenAddress       string   `move:"address"`
	DestTokenPoolPackageId string   `move:"address"`
	SourcePoolAddress      []byte   `move:"vector<u8>"`
	SourcePoolData         []byte   `move:"vector<u8>"`
	OffchainTokenData      []byte   `move:"vector<u8>"`
}

type bcsCompletedDestTokenTransfer struct {
	TokenReceiver    [32]byte
	DestTokenAddress [32]byte
}

func convertCompletedDestTokenTransferFromBCS(bcs bcsCompletedDestTokenTransfer) (CompletedDestTokenTransfer, error) {

	return CompletedDestTokenTransfer{
		TokenReceiver:    fmt.Sprintf("0x%x", bcs.TokenReceiver),
		DestTokenAddress: fmt.Sprintf("0x%x", bcs.DestTokenAddress),
	}, nil
}

type bcsDestTokenTransfer struct {
	TokenReceiver          [32]byte
	RemoteChainSelector    uint64
	SourceAmount           [32]byte
	DestTokenAddress       [32]byte
	DestTokenPoolPackageId [32]byte
	SourcePoolAddress      []byte
	SourcePoolData         []byte
	OffchainTokenData      []byte
}

func convertDestTokenTransferFromBCS(bcs bcsDestTokenTransfer) (DestTokenTransfer, error) {
	SourceAmountField, err := bind.DecodeU256Value(bcs.SourceAmount)
	if err != nil {
		return DestTokenTransfer{}, fmt.Errorf("failed to decode u256 field SourceAmount: %w", err)
	}

	return DestTokenTransfer{
		TokenReceiver:          fmt.Sprintf("0x%x", bcs.TokenReceiver),
		RemoteChainSelector:    bcs.RemoteChainSelector,
		SourceAmount:           SourceAmountField,
		DestTokenAddress:       fmt.Sprintf("0x%x", bcs.DestTokenAddress),
		DestTokenPoolPackageId: fmt.Sprintf("0x%x", bcs.DestTokenPoolPackageId),
		SourcePoolAddress:      bcs.SourcePoolAddress,
		SourcePoolData:         bcs.SourcePoolData,
		OffchainTokenData:      bcs.OffchainTokenData,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("ccip::offramp_state_helper::OFFRAMP_STATE_HELPER", func(data []byte) (interface{}, error) {
		var result OFFRAMP_STATE_HELPER
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for OFFRAMP_STATE_HELPER
	bind.RegisterStructDecoder("vector<ccip::offramp_state_helper::OFFRAMP_STATE_HELPER>", func(data []byte) (interface{}, error) {
		var results []OFFRAMP_STATE_HELPER
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::offramp_state_helper::ReceiverParams", func(data []byte) (interface{}, error) {
		var result ReceiverParams
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for ReceiverParams
	bind.RegisterStructDecoder("vector<ccip::offramp_state_helper::ReceiverParams>", func(data []byte) (interface{}, error) {
		var results []ReceiverParams
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::offramp_state_helper::DestTransferCap", func(data []byte) (interface{}, error) {
		var result DestTransferCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DestTransferCap
	bind.RegisterStructDecoder("vector<ccip::offramp_state_helper::DestTransferCap>", func(data []byte) (interface{}, error) {
		var results []DestTransferCap
		_, err := mystenbcs.Unmarshal(data, &results)
		if err != nil {
			return nil, err
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::offramp_state_helper::CompletedDestTokenTransfer", func(data []byte) (interface{}, error) {
		var temp bcsCompletedDestTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertCompletedDestTokenTransferFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for CompletedDestTokenTransfer
	bind.RegisterStructDecoder("vector<ccip::offramp_state_helper::CompletedDestTokenTransfer>", func(data []byte) (interface{}, error) {
		var temps []bcsCompletedDestTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]CompletedDestTokenTransfer, len(temps))
		for i, temp := range temps {
			result, err := convertCompletedDestTokenTransferFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
	bind.RegisterStructDecoder("ccip::offramp_state_helper::DestTokenTransfer", func(data []byte) (interface{}, error) {
		var temp bcsDestTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDestTokenTransferFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	// Register vector decoder for DestTokenTransfer
	bind.RegisterStructDecoder("vector<ccip::offramp_state_helper::DestTokenTransfer>", func(data []byte) (interface{}, error) {
		var temps []bcsDestTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temps)
		if err != nil {
			return nil, err
		}

		results := make([]DestTokenTransfer, len(temps))
		for i, temp := range temps {
			result, err := convertDestTokenTransferFromBCS(temp)
			if err != nil {
				return nil, fmt.Errorf("failed to convert element %d: %w", i, err)
			}
			results[i] = result
		}
		return results, nil
	})
}

// NewDestTransferCap executes the new_dest_transfer_cap Move function.
func (c *OfframpStateHelperContract) NewDestTransferCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.NewDestTransferCap(ref, ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CreateReceiverParams executes the create_receiver_params Move function.
func (c *OfframpStateHelperContract) CreateReceiverParams(ctx context.Context, opts *bind.CallOpts, param bind.Object, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.CreateReceiverParams(param, sourceChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSourceChainSelector executes the get_source_chain_selector Move function.
func (c *OfframpStateHelperContract) GetSourceChainSelector(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.GetSourceChainSelector(receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddDestTokenTransfer executes the add_dest_token_transfer Move function.
func (c *OfframpStateHelperContract) AddDestTokenTransfer(ctx context.Context, opts *bind.CallOpts, param bind.Object, receiverParams ReceiverParams, tokenReceiver string, remoteChainSelector uint64, sourceAmount *big.Int, destTokenAddress string, destTokenPoolPackageId string, sourcePoolAddress []byte, sourcePoolData []byte, offchainData []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.AddDestTokenTransfer(param, receiverParams, tokenReceiver, remoteChainSelector, sourceAmount, destTokenAddress, destTokenPoolPackageId, sourcePoolAddress, sourcePoolData, offchainData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PopulateMessage executes the populate_message Move function.
func (c *OfframpStateHelperContract) PopulateMessage(ctx context.Context, opts *bind.CallOpts, param bind.Object, receiverParams ReceiverParams, any2suiMessage bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.PopulateMessage(param, receiverParams, any2suiMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestTokenTransferData executes the get_dest_token_transfer_data Move function.
func (c *OfframpStateHelperContract) GetDestTokenTransferData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.GetDestTokenTransferData(receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenParamData executes the get_token_param_data Move function.
func (c *OfframpStateHelperContract) GetTokenParamData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.GetTokenParamData(receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CompleteTokenTransfer executes the complete_token_transfer Move function.
func (c *OfframpStateHelperContract) CompleteTokenTransfer(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, receiverParams ReceiverParams, param bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.CompleteTokenTransfer(typeArgs, ref, receiverParams, param)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExtractAny2suiMessage executes the extract_any2sui_message Move function.
func (c *OfframpStateHelperContract) ExtractAny2suiMessage(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.ExtractAny2suiMessage(receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// NewAny2suiMessage executes the new_any2sui_message Move function.
func (c *OfframpStateHelperContract) NewAny2suiMessage(ctx context.Context, opts *bind.CallOpts, param bind.Object, messageId []byte, sourceChainSelector uint64, sender []byte, data []byte, messageReceiver string, tokenReceiver string, tokenAddresses []string, tokenAmounts []*big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.NewAny2suiMessage(param, messageId, sourceChainSelector, sender, data, messageReceiver, tokenReceiver, tokenAddresses, tokenAmounts)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConsumeAny2suiMessage executes the consume_any2sui_message Move function.
func (c *OfframpStateHelperContract) ConsumeAny2suiMessage(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, message bind.Object, param bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.ConsumeAny2suiMessage(typeArgs, ref, message, param)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// DeconstructReceiverParams executes the deconstruct_receiver_params Move function.
func (c *OfframpStateHelperContract) DeconstructReceiverParams(ctx context.Context, opts *bind.CallOpts, param bind.Object, receiverParams ReceiverParams) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampStateHelperEncoder.DeconstructReceiverParams(param, receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// NewDestTransferCap executes the new_dest_transfer_cap Move function using DevInspect to get return values.
//
// Returns: DestTransferCap
func (d *OfframpStateHelperDevInspect) NewDestTransferCap(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object) (bind.Object, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.NewDestTransferCap(ref, ownerCap)
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

// CreateReceiverParams executes the create_receiver_params Move function using DevInspect to get return values.
//
// Returns: ReceiverParams
func (d *OfframpStateHelperDevInspect) CreateReceiverParams(ctx context.Context, opts *bind.CallOpts, param bind.Object, sourceChainSelector uint64) (ReceiverParams, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.CreateReceiverParams(param, sourceChainSelector)
	if err != nil {
		return ReceiverParams{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return ReceiverParams{}, err
	}
	if len(results) == 0 {
		return ReceiverParams{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(ReceiverParams)
	if !ok {
		return ReceiverParams{}, fmt.Errorf("unexpected return type: expected ReceiverParams, got %T", results[0])
	}
	return result, nil
}

// GetSourceChainSelector executes the get_source_chain_selector Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OfframpStateHelperDevInspect) GetSourceChainSelector(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (uint64, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.GetSourceChainSelector(receiverParams)
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

// GetDestTokenTransferData executes the get_dest_token_transfer_data Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: u64
//	[2]: u256
//	[3]: address
//	[4]: address
//	[5]: vector<u8>
//	[6]: vector<u8>
//	[7]: vector<u8>
func (d *OfframpStateHelperDevInspect) GetDestTokenTransferData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) ([]any, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.GetDestTokenTransferData(receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetTokenParamData executes the get_token_param_data Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: u256
//	[2]: address
//	[3]: vector<u8>
//	[4]: vector<u8>
//	[5]: vector<u8>
func (d *OfframpStateHelperDevInspect) GetTokenParamData(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) ([]any, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.GetTokenParamData(receiverParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// ExtractAny2suiMessage executes the extract_any2sui_message Move function using DevInspect to get return values.
//
// Returns: Any2SuiMessage
func (d *OfframpStateHelperDevInspect) ExtractAny2suiMessage(ctx context.Context, opts *bind.CallOpts, receiverParams ReceiverParams) (bind.Object, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.ExtractAny2suiMessage(receiverParams)
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

// NewAny2suiMessage executes the new_any2sui_message Move function using DevInspect to get return values.
//
// Returns: Any2SuiMessage
func (d *OfframpStateHelperDevInspect) NewAny2suiMessage(ctx context.Context, opts *bind.CallOpts, param bind.Object, messageId []byte, sourceChainSelector uint64, sender []byte, data []byte, messageReceiver string, tokenReceiver string, tokenAddresses []string, tokenAmounts []*big.Int) (bind.Object, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.NewAny2suiMessage(param, messageId, sourceChainSelector, sender, data, messageReceiver, tokenReceiver, tokenAddresses, tokenAmounts)
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

// ConsumeAny2suiMessage executes the consume_any2sui_message Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<u8>
//	[1]: u64
//	[2]: vector<u8>
//	[3]: vector<u8>
//	[4]: address
//	[5]: address
//	[6]: vector<Any2SuiTokenAmount>
func (d *OfframpStateHelperDevInspect) ConsumeAny2suiMessage(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, message bind.Object, param bind.Object) ([]any, error) {
	encoded, err := d.contract.offrampStateHelperEncoder.ConsumeAny2suiMessage(typeArgs, ref, message, param)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

type offrampStateHelperEncoder struct {
	*bind.BoundContract
}

// NewDestTransferCap encodes a call to the new_dest_transfer_cap Move function.
func (c offrampStateHelperEncoder) NewDestTransferCap(ref bind.Object, ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_dest_transfer_cap", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&OwnerCap",
	}, []any{
		ref,
		ownerCap,
	}, []string{
		"ccip::offramp_state_helper::DestTransferCap",
	})
}

// NewDestTransferCapWithArgs encodes a call to the new_dest_transfer_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) NewDestTransferCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_dest_transfer_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::offramp_state_helper::DestTransferCap",
	})
}

// CreateReceiverParams encodes a call to the create_receiver_params Move function.
func (c offrampStateHelperEncoder) CreateReceiverParams(param bind.Object, sourceChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_receiver_params", typeArgsList, typeParamsList, []string{
		"&DestTransferCap",
		"u64",
	}, []any{
		param,
		sourceChainSelector,
	}, []string{
		"ccip::offramp_state_helper::ReceiverParams",
	})
}

// CreateReceiverParamsWithArgs encodes a call to the create_receiver_params Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) CreateReceiverParamsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&DestTransferCap",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("create_receiver_params", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::offramp_state_helper::ReceiverParams",
	})
}

// GetSourceChainSelector encodes a call to the get_source_chain_selector Move function.
func (c offrampStateHelperEncoder) GetSourceChainSelector(receiverParams ReceiverParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_selector", typeArgsList, typeParamsList, []string{
		"&ReceiverParams",
	}, []any{
		receiverParams,
	}, []string{
		"u64",
	})
}

// GetSourceChainSelectorWithArgs encodes a call to the get_source_chain_selector Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) GetSourceChainSelectorWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ReceiverParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_selector", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// AddDestTokenTransfer encodes a call to the add_dest_token_transfer Move function.
func (c offrampStateHelperEncoder) AddDestTokenTransfer(param bind.Object, receiverParams ReceiverParams, tokenReceiver string, remoteChainSelector uint64, sourceAmount *big.Int, destTokenAddress string, destTokenPoolPackageId string, sourcePoolAddress []byte, sourcePoolData []byte, offchainData []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_dest_token_transfer", typeArgsList, typeParamsList, []string{
		"&DestTransferCap",
		"&mut ReceiverParams",
		"address",
		"u64",
		"u256",
		"address",
		"address",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
	}, []any{
		param,
		receiverParams,
		tokenReceiver,
		remoteChainSelector,
		sourceAmount,
		destTokenAddress,
		destTokenPoolPackageId,
		sourcePoolAddress,
		sourcePoolData,
		offchainData,
	}, nil)
}

// AddDestTokenTransferWithArgs encodes a call to the add_dest_token_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) AddDestTokenTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&DestTransferCap",
		"&mut ReceiverParams",
		"address",
		"u64",
		"u256",
		"address",
		"address",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_dest_token_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// PopulateMessage encodes a call to the populate_message Move function.
func (c offrampStateHelperEncoder) PopulateMessage(param bind.Object, receiverParams ReceiverParams, any2suiMessage bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("populate_message", typeArgsList, typeParamsList, []string{
		"&DestTransferCap",
		"&mut ReceiverParams",
		"Any2SuiMessage",
	}, []any{
		param,
		receiverParams,
		any2suiMessage,
	}, nil)
}

// PopulateMessageWithArgs encodes a call to the populate_message Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) PopulateMessageWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&DestTransferCap",
		"&mut ReceiverParams",
		"Any2SuiMessage",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("populate_message", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetDestTokenTransferData encodes a call to the get_dest_token_transfer_data Move function.
func (c offrampStateHelperEncoder) GetDestTokenTransferData(receiverParams ReceiverParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_token_transfer_data", typeArgsList, typeParamsList, []string{
		"&ReceiverParams",
	}, []any{
		receiverParams,
	}, []string{
		"address",
		"u64",
		"u256",
		"address",
		"address",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
	})
}

// GetDestTokenTransferDataWithArgs encodes a call to the get_dest_token_transfer_data Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) GetDestTokenTransferDataWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ReceiverParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_token_transfer_data", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"u64",
		"u256",
		"address",
		"address",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
	})
}

// GetTokenParamData encodes a call to the get_token_param_data Move function.
func (c offrampStateHelperEncoder) GetTokenParamData(receiverParams ReceiverParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_param_data", typeArgsList, typeParamsList, []string{
		"&ReceiverParams",
	}, []any{
		receiverParams,
	}, []string{
		"address",
		"u256",
		"address",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
	})
}

// GetTokenParamDataWithArgs encodes a call to the get_token_param_data Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) GetTokenParamDataWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&ReceiverParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_param_data", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"u256",
		"address",
		"vector<u8>",
		"vector<u8>",
		"vector<u8>",
	})
}

// CompleteTokenTransfer encodes a call to the complete_token_transfer Move function.
func (c offrampStateHelperEncoder) CompleteTokenTransfer(typeArgs []string, ref bind.Object, receiverParams ReceiverParams, param bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("complete_token_transfer", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut ReceiverParams",
		"TypeProof",
	}, []any{
		ref,
		receiverParams,
		param,
	}, nil)
}

// CompleteTokenTransferWithArgs encodes a call to the complete_token_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) CompleteTokenTransferWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut ReceiverParams",
		"TypeProof",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("complete_token_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExtractAny2suiMessage encodes a call to the extract_any2sui_message Move function.
func (c offrampStateHelperEncoder) ExtractAny2suiMessage(receiverParams ReceiverParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("extract_any2sui_message", typeArgsList, typeParamsList, []string{
		"&mut ReceiverParams",
	}, []any{
		receiverParams,
	}, []string{
		"Any2SuiMessage",
	})
}

// ExtractAny2suiMessageWithArgs encodes a call to the extract_any2sui_message Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) ExtractAny2suiMessageWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut ReceiverParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("extract_any2sui_message", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Any2SuiMessage",
	})
}

// NewAny2suiMessage encodes a call to the new_any2sui_message Move function.
func (c offrampStateHelperEncoder) NewAny2suiMessage(param bind.Object, messageId []byte, sourceChainSelector uint64, sender []byte, data []byte, messageReceiver string, tokenReceiver string, tokenAddresses []string, tokenAmounts []*big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_any2sui_message", typeArgsList, typeParamsList, []string{
		"&DestTransferCap",
		"vector<u8>",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"address",
		"address",
		"vector<address>",
		"vector<u256>",
	}, []any{
		param,
		messageId,
		sourceChainSelector,
		sender,
		data,
		messageReceiver,
		tokenReceiver,
		tokenAddresses,
		tokenAmounts,
	}, []string{
		"Any2SuiMessage",
	})
}

// NewAny2suiMessageWithArgs encodes a call to the new_any2sui_message Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) NewAny2suiMessageWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&DestTransferCap",
		"vector<u8>",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"address",
		"address",
		"vector<address>",
		"vector<u256>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_any2sui_message", typeArgsList, typeParamsList, expectedParams, args, []string{
		"Any2SuiMessage",
	})
}

// ConsumeAny2suiMessage encodes a call to the consume_any2sui_message Move function.
func (c offrampStateHelperEncoder) ConsumeAny2suiMessage(typeArgs []string, ref bind.Object, message bind.Object, param bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("consume_any2sui_message", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"Any2SuiMessage",
		"TypeProof",
	}, []any{
		ref,
		message,
		param,
	}, []string{
		"vector<u8>",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"address",
		"address",
		"vector<Any2SuiTokenAmount>",
	})
}

// ConsumeAny2suiMessageWithArgs encodes a call to the consume_any2sui_message Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) ConsumeAny2suiMessageWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"Any2SuiMessage",
		"TypeProof",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"TypeProof",
	}
	return c.EncodeCallArgsWithGenerics("consume_any2sui_message", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"address",
		"address",
		"vector<Any2SuiTokenAmount>",
	})
}

// DeconstructReceiverParams encodes a call to the deconstruct_receiver_params Move function.
func (c offrampStateHelperEncoder) DeconstructReceiverParams(param bind.Object, receiverParams ReceiverParams) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("deconstruct_receiver_params", typeArgsList, typeParamsList, []string{
		"&DestTransferCap",
		"ccip::offramp_state_helper::ReceiverParams",
	}, []any{
		param,
		receiverParams,
	}, nil)
}

// DeconstructReceiverParamsWithArgs encodes a call to the deconstruct_receiver_params Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampStateHelperEncoder) DeconstructReceiverParamsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&DestTransferCap",
		"ccip::offramp_state_helper::ReceiverParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("deconstruct_receiver_params", typeArgsList, typeParamsList, expectedParams, args, nil)
}
