// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_rmn_remote

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

type IRmnRemote interface {
	EmitConfigSetEvent(ctx context.Context, opts *bind.CallOpts, version uint32) (*models.SuiTransactionBlockResponse, error)
	EmitCursedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	EmitCursedMultipleEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	EmitUncursedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	EmitUncursedMultipleEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IRmnRemoteDevInspect
	Encoder() RmnRemoteEncoder
	Bound() bind.IBoundContract
}

type IRmnRemoteDevInspect interface {
}

type RmnRemoteEncoder interface {
	EmitConfigSetEvent(version uint32) (*bind.EncodedCall, error)
	EmitConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitCursedEvent() (*bind.EncodedCall, error)
	EmitCursedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitCursedMultipleEvent() (*bind.EncodedCall, error)
	EmitCursedMultipleEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitUncursedEvent() (*bind.EncodedCall, error)
	EmitUncursedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitUncursedMultipleEvent() (*bind.EncodedCall, error)
	EmitUncursedMultipleEventWithArgs(args ...any) (*bind.EncodedCall, error)
}

type RmnRemoteContract struct {
	*bind.BoundContract
	rmnRemoteEncoder
	devInspect *RmnRemoteDevInspect
}

type RmnRemoteDevInspect struct {
	contract *RmnRemoteContract
}

var _ IRmnRemote = (*RmnRemoteContract)(nil)
var _ IRmnRemoteDevInspect = (*RmnRemoteDevInspect)(nil)

func NewRmnRemote(packageID string, client sui.ISuiAPI) (IRmnRemote, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "rmn_remote", client)
	if err != nil {
		return nil, err
	}

	c := &RmnRemoteContract{
		BoundContract:    contract,
		rmnRemoteEncoder: rmnRemoteEncoder{BoundContract: contract},
	}
	c.devInspect = &RmnRemoteDevInspect{contract: c}
	return c, nil
}

func (c *RmnRemoteContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *RmnRemoteContract) Encoder() RmnRemoteEncoder {
	return c.rmnRemoteEncoder
}

func (c *RmnRemoteContract) DevInspect() IRmnRemoteDevInspect {
	return c.devInspect
}

type RMNRemoteState struct {
	Id                 string      `move:"sui::object::UID"`
	LocalChainSelector uint64      `move:"u64"`
	Config             Config      `move:"Config"`
	ConfigCount        uint32      `move:"u32"`
	Signers            bind.Object `move:"VecMap<vector<u8>, bool>"`
	CursedSubjects     bind.Object `move:"VecMap<vector<u8>, bool>"`
}

type Config struct {
	RmnHomeContractConfigDigest []byte   `move:"vector<u8>"`
	Signers                     []Signer `move:"vector<Signer>"`
	FSign                       uint64   `move:"u64"`
}

type Signer struct {
	OnchainPublicKey []byte `move:"vector<u8>"`
	NodeIndex        uint64 `move:"u64"`
}

type Report struct {
	DestChainSelector           uint64       `move:"u64"`
	RmnRemoteContractAddress    string       `move:"address"`
	OffRampAddress              string       `move:"address"`
	RmnHomeContractConfigDigest []byte       `move:"vector<u8>"`
	MerkleRoots                 []MerkleRoot `move:"vector<MerkleRoot>"`
}

type MerkleRoot struct {
	SourceChainSelector uint64 `move:"u64"`
	OnRampAddress       []byte `move:"vector<u8>"`
	MinSeqNr            uint64 `move:"u64"`
	MaxSeqNr            uint64 `move:"u64"`
	MerkleRoot          []byte `move:"vector<u8>"`
}

type ConfigSet struct {
	Version uint32 `move:"u32"`
	Config  Config `move:"Config"`
}

type Cursed struct {
	Subjects [][]byte `move:"vector<vector<u8>>"`
}

type Uncursed struct {
	Subjects [][]byte `move:"vector<vector<u8>>"`
}

type bcsReport struct {
	DestChainSelector           uint64
	RmnRemoteContractAddress    [32]byte
	OffRampAddress              [32]byte
	RmnHomeContractConfigDigest []byte
	MerkleRoots                 []MerkleRoot
}

func convertReportFromBCS(bcs bcsReport) (Report, error) {

	return Report{
		DestChainSelector:           bcs.DestChainSelector,
		RmnRemoteContractAddress:    fmt.Sprintf("0x%x", bcs.RmnRemoteContractAddress),
		OffRampAddress:              fmt.Sprintf("0x%x", bcs.OffRampAddress),
		RmnHomeContractConfigDigest: bcs.RmnHomeContractConfigDigest,
		MerkleRoots:                 bcs.MerkleRoots,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("test::rmn_remote::RMNRemoteState", func(data []byte) (interface{}, error) {
		var result RMNRemoteState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::rmn_remote::Config", func(data []byte) (interface{}, error) {
		var result Config
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::rmn_remote::Signer", func(data []byte) (interface{}, error) {
		var result Signer
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::rmn_remote::Report", func(data []byte) (interface{}, error) {
		var temp bcsReport
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertReportFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::rmn_remote::MerkleRoot", func(data []byte) (interface{}, error) {
		var result MerkleRoot
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::rmn_remote::ConfigSet", func(data []byte) (interface{}, error) {
		var result ConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::rmn_remote::Cursed", func(data []byte) (interface{}, error) {
		var result Cursed
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::rmn_remote::Uncursed", func(data []byte) (interface{}, error) {
		var result Uncursed
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// EmitConfigSetEvent executes the emit_config_set_event Move function.
func (c *RmnRemoteContract) EmitConfigSetEvent(ctx context.Context, opts *bind.CallOpts, version uint32) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.EmitConfigSetEvent(version)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitCursedEvent executes the emit_cursed_event Move function.
func (c *RmnRemoteContract) EmitCursedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.EmitCursedEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitCursedMultipleEvent executes the emit_cursed_multiple_event Move function.
func (c *RmnRemoteContract) EmitCursedMultipleEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.EmitCursedMultipleEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitUncursedEvent executes the emit_uncursed_event Move function.
func (c *RmnRemoteContract) EmitUncursedEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.EmitUncursedEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitUncursedMultipleEvent executes the emit_uncursed_multiple_event Move function.
func (c *RmnRemoteContract) EmitUncursedMultipleEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.rmnRemoteEncoder.EmitUncursedMultipleEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

type rmnRemoteEncoder struct {
	*bind.BoundContract
}

// EmitConfigSetEvent encodes a call to the emit_config_set_event Move function.
func (c rmnRemoteEncoder) EmitConfigSetEvent(version uint32) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_config_set_event", typeArgsList, typeParamsList, []string{
		"u32",
	}, []any{
		version,
	}, nil)
}

// EmitConfigSetEventWithArgs encodes a call to the emit_config_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) EmitConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u32",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_config_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitCursedEvent encodes a call to the emit_cursed_event Move function.
func (c rmnRemoteEncoder) EmitCursedEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_cursed_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitCursedEventWithArgs encodes a call to the emit_cursed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) EmitCursedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_cursed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitCursedMultipleEvent encodes a call to the emit_cursed_multiple_event Move function.
func (c rmnRemoteEncoder) EmitCursedMultipleEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_cursed_multiple_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitCursedMultipleEventWithArgs encodes a call to the emit_cursed_multiple_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) EmitCursedMultipleEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_cursed_multiple_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitUncursedEvent encodes a call to the emit_uncursed_event Move function.
func (c rmnRemoteEncoder) EmitUncursedEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_uncursed_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitUncursedEventWithArgs encodes a call to the emit_uncursed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) EmitUncursedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_uncursed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitUncursedMultipleEvent encodes a call to the emit_uncursed_multiple_event Move function.
func (c rmnRemoteEncoder) EmitUncursedMultipleEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_uncursed_multiple_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitUncursedMultipleEventWithArgs encodes a call to the emit_uncursed_multiple_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c rmnRemoteEncoder) EmitUncursedMultipleEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_uncursed_multiple_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}
