// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_offramp

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

type IOfframp interface {
	EmitStaticConfigSetEvent(ctx context.Context, opts *bind.CallOpts, chainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitDynamicConfigSetEvent(ctx context.Context, opts *bind.CallOpts, dynamicConfig DynamicConfig) (*models.SuiTransactionBlockResponse, error)
	EmitSourceChainConfigSetEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, sourceChainConfig SourceChainConfig) (*models.SuiTransactionBlockResponse, error)
	EmitSkippedAlreadyExecutedEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error)
	EmitExecutionStateChangedEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId []byte, messageHash []byte, state byte) (*models.SuiTransactionBlockResponse, error)
	EmitCommitReportAcceptedEvent(ctx context.Context, opts *bind.CallOpts, blessedMerkleRoots []MerkleRoot, unblessedMerkleRoots []MerkleRoot, priceUpdates PriceUpdates) (*models.SuiTransactionBlockResponse, error)
	EmitSkippedReportExecutionEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	EmitOcrConfigEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	AddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error)
	RemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error)
	GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*models.SuiTransactionBlockResponse, error)
	FinishExecute(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOfframpDevInspect
	Encoder() OfframpEncoder
	Bound() bind.IBoundContract
}

type IOfframpDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts) ([]any, error)
}

type OfframpEncoder interface {
	EmitStaticConfigSetEvent(chainSelector uint64) (*bind.EncodedCall, error)
	EmitStaticConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitDynamicConfigSetEvent(dynamicConfig DynamicConfig) (*bind.EncodedCall, error)
	EmitDynamicConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitSourceChainConfigSetEvent(sourceChainSelector uint64, sourceChainConfig SourceChainConfig) (*bind.EncodedCall, error)
	EmitSourceChainConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitSkippedAlreadyExecutedEvent(sourceChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error)
	EmitSkippedAlreadyExecutedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitExecutionStateChangedEvent(sourceChainSelector uint64, sequenceNumber uint64, messageId []byte, messageHash []byte, state byte) (*bind.EncodedCall, error)
	EmitExecutionStateChangedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitCommitReportAcceptedEvent(blessedMerkleRoots []MerkleRoot, unblessedMerkleRoots []MerkleRoot, priceUpdates PriceUpdates) (*bind.EncodedCall, error)
	EmitCommitReportAcceptedEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitSkippedReportExecutionEvent(sourceChainSelector uint64) (*bind.EncodedCall, error)
	EmitSkippedReportExecutionEventWithArgs(args ...any) (*bind.EncodedCall, error)
	EmitOcrConfigEvent() (*bind.EncodedCall, error)
	EmitOcrConfigEventWithArgs(args ...any) (*bind.EncodedCall, error)
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	AddPackageId(state bind.Object, packageId string) (*bind.EncodedCall, error)
	AddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	RemovePackageId(state bind.Object, packageId string) (*bind.EncodedCall, error)
	RemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllSourceChainConfigs() (*bind.EncodedCall, error)
	GetAllSourceChainConfigsWithArgs(args ...any) (*bind.EncodedCall, error)
	InitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*bind.EncodedCall, error)
	InitExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	FinishExecute() (*bind.EncodedCall, error)
	FinishExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
}

type OfframpContract struct {
	*bind.BoundContract
	offrampEncoder
	devInspect *OfframpDevInspect
}

type OfframpDevInspect struct {
	contract *OfframpContract
}

var _ IOfframp = (*OfframpContract)(nil)
var _ IOfframpDevInspect = (*OfframpDevInspect)(nil)

func NewOfframp(packageID string, client sui.ISuiAPI) (IOfframp, error) {
	contract, err := bind.NewBoundContract(packageID, "test", "offramp", client)
	if err != nil {
		return nil, err
	}

	c := &OfframpContract{
		BoundContract:  contract,
		offrampEncoder: offrampEncoder{BoundContract: contract},
	}
	c.devInspect = &OfframpDevInspect{contract: c}
	return c, nil
}

func (c *OfframpContract) Bound() bind.IBoundContract {
	return c.BoundContract
}

func (c *OfframpContract) Encoder() OfframpEncoder {
	return c.offrampEncoder
}

func (c *OfframpContract) DevInspect() IOfframpDevInspect {
	return c.devInspect
}

type CCIPObjectRef struct {
	Id string `move:"sui::object::UID"`
}

type OffRampState struct {
	Id         string   `move:"sui::object::UID"`
	PackageIds []string `move:"vector<address>"`
}

type OffRampStatePointer struct {
	Id             string `move:"sui::object::UID"`
	OffRampStateId string `move:"address"`
	OwnerCapId     string `move:"address"`
}

type SourceChainConfig struct {
	Router                    string `move:"address"`
	IsEnabled                 bool   `move:"bool"`
	MinSeqNr                  uint64 `move:"u64"`
	IsRmnVerificationDisabled bool   `move:"bool"`
	OnRamp                    []byte `move:"vector<u8>"`
}

type RampMessageHeader struct {
	MessageId           []byte `move:"vector<u8>"`
	SourceChainSelector uint64 `move:"u64"`
	DestChainSelector   uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
	Nonce               uint64 `move:"u64"`
}

type Any2SuiRampMessage struct {
	Header       RampMessageHeader      `move:"RampMessageHeader"`
	Sender       []byte                 `move:"vector<u8>"`
	Data         []byte                 `move:"vector<u8>"`
	Receiver     string                 `move:"address"`
	GasLimit     *big.Int               `move:"u256"`
	TokenAmounts []Any2SuiTokenTransfer `move:"vector<Any2SuiTokenTransfer>"`
}

type Any2SuiTokenTransfer struct {
	SourcePoolAddress []byte   `move:"vector<u8>"`
	DestTokenAddress  string   `move:"address"`
	DestGasAmount     uint32   `move:"u32"`
	ExtraData         []byte   `move:"vector<u8>"`
	Amount            *big.Int `move:"u256"`
}

type ExecutionReport struct {
	SourceChainSelector uint64             `move:"u64"`
	Message             Any2SuiRampMessage `move:"Any2SuiRampMessage"`
	OffchainTokenData   [][]byte           `move:"vector<vector<u8>>"`
	Proofs              [][]byte           `move:"vector<vector<u8>>"`
}

type CommitReport struct {
	PriceUpdates         PriceUpdates `move:"PriceUpdates"`
	BlessedMerkleRoots   []MerkleRoot `move:"vector<MerkleRoot>"`
	UnblessedMerkleRoots []MerkleRoot `move:"vector<MerkleRoot>"`
	RmnSignatures        [][]byte     `move:"vector<vector<u8>>"`
}

type PriceUpdates struct {
	TokenPriceUpdates []TokenPriceUpdate `move:"vector<TokenPriceUpdate>"`
	GasPriceUpdates   []GasPriceUpdate   `move:"vector<GasPriceUpdate>"`
}

type TokenPriceUpdate struct {
	SourceToken string   `move:"address"`
	UsdPerToken *big.Int `move:"u256"`
}

type GasPriceUpdate struct {
	DestChainSelector uint64   `move:"u64"`
	UsdPerUnitGas     *big.Int `move:"u256"`
}

type MerkleRoot struct {
	SourceChainSelector uint64 `move:"u64"`
	OnRampAddress       []byte `move:"vector<u8>"`
	MinSeqNr            uint64 `move:"u64"`
	MaxSeqNr            uint64 `move:"u64"`
	MerkleRoot          []byte `move:"vector<u8>"`
}

type StaticConfig struct {
	ChainSelector      uint64 `move:"u64"`
	RmnRemote          string `move:"address"`
	TokenAdminRegistry string `move:"address"`
	NonceManager       string `move:"address"`
}

type DynamicConfig struct {
	FeeQuoter                               string `move:"address"`
	PermissionlessExecutionThresholdSeconds uint32 `move:"u32"`
}

type StaticConfigSet struct {
	ChainSelector uint64 `move:"u64"`
}

type DynamicConfigSet struct {
	DynamicConfig DynamicConfig `move:"DynamicConfig"`
}

type SourceChainConfigSet struct {
	SourceChainSelector uint64            `move:"u64"`
	SourceChainConfig   SourceChainConfig `move:"SourceChainConfig"`
}

type SkippedAlreadyExecuted struct {
	SourceChainSelector uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
}

type ExecutionStateChanged struct {
	SourceChainSelector uint64 `move:"u64"`
	SequenceNumber      uint64 `move:"u64"`
	MessageId           []byte `move:"vector<u8>"`
	MessageHash         []byte `move:"vector<u8>"`
	State               byte   `move:"u8"`
}

type CommitReportAccepted struct {
	BlessedMerkleRoots   []MerkleRoot `move:"vector<MerkleRoot>"`
	UnblessedMerkleRoots []MerkleRoot `move:"vector<MerkleRoot>"`
	PriceUpdates         PriceUpdates `move:"PriceUpdates"`
}

type SkippedReportExecution struct {
	SourceChainSelector uint64 `move:"u64"`
}

type OFFRAMP struct {
}

type bcsOffRampState struct {
	Id         string
	PackageIds [][32]byte
}

func convertOffRampStateFromBCS(bcs bcsOffRampState) (OffRampState, error) {

	return OffRampState{
		Id: bcs.Id,
		PackageIds: func() []string {
			addrs := make([]string, len(bcs.PackageIds))
			for i, addr := range bcs.PackageIds {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
	}, nil
}

type bcsOffRampStatePointer struct {
	Id             string
	OffRampStateId [32]byte
	OwnerCapId     [32]byte
}

func convertOffRampStatePointerFromBCS(bcs bcsOffRampStatePointer) (OffRampStatePointer, error) {

	return OffRampStatePointer{
		Id:             bcs.Id,
		OffRampStateId: fmt.Sprintf("0x%x", bcs.OffRampStateId),
		OwnerCapId:     fmt.Sprintf("0x%x", bcs.OwnerCapId),
	}, nil
}

type bcsSourceChainConfig struct {
	Router                    [32]byte
	IsEnabled                 bool
	MinSeqNr                  uint64
	IsRmnVerificationDisabled bool
	OnRamp                    []byte
}

func convertSourceChainConfigFromBCS(bcs bcsSourceChainConfig) (SourceChainConfig, error) {

	return SourceChainConfig{
		Router:                    fmt.Sprintf("0x%x", bcs.Router),
		IsEnabled:                 bcs.IsEnabled,
		MinSeqNr:                  bcs.MinSeqNr,
		IsRmnVerificationDisabled: bcs.IsRmnVerificationDisabled,
		OnRamp:                    bcs.OnRamp,
	}, nil
}

type bcsAny2SuiRampMessage struct {
	Header       RampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     [32]byte
	GasLimit     [32]byte
	TokenAmounts []Any2SuiTokenTransfer
}

func convertAny2SuiRampMessageFromBCS(bcs bcsAny2SuiRampMessage) (Any2SuiRampMessage, error) {
	GasLimitField, err := bind.DecodeU256Value(bcs.GasLimit)
	if err != nil {
		return Any2SuiRampMessage{}, fmt.Errorf("failed to decode u256 field GasLimit: %w", err)
	}

	return Any2SuiRampMessage{
		Header:       bcs.Header,
		Sender:       bcs.Sender,
		Data:         bcs.Data,
		Receiver:     fmt.Sprintf("0x%x", bcs.Receiver),
		GasLimit:     GasLimitField,
		TokenAmounts: bcs.TokenAmounts,
	}, nil
}

type bcsAny2SuiTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  [32]byte
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            [32]byte
}

func convertAny2SuiTokenTransferFromBCS(bcs bcsAny2SuiTokenTransfer) (Any2SuiTokenTransfer, error) {
	AmountField, err := bind.DecodeU256Value(bcs.Amount)
	if err != nil {
		return Any2SuiTokenTransfer{}, fmt.Errorf("failed to decode u256 field Amount: %w", err)
	}

	return Any2SuiTokenTransfer{
		SourcePoolAddress: bcs.SourcePoolAddress,
		DestTokenAddress:  fmt.Sprintf("0x%x", bcs.DestTokenAddress),
		DestGasAmount:     bcs.DestGasAmount,
		ExtraData:         bcs.ExtraData,
		Amount:            AmountField,
	}, nil
}

type bcsExecutionReport struct {
	SourceChainSelector uint64
	Message             bcsAny2SuiRampMessage
	OffchainTokenData   [][]byte
	Proofs              [][]byte
}

func convertExecutionReportFromBCS(bcs bcsExecutionReport) (ExecutionReport, error) {
	MessageField, err := convertAny2SuiRampMessageFromBCS(bcs.Message)
	if err != nil {
		return ExecutionReport{}, fmt.Errorf("failed to convert nested struct Message: %w", err)
	}

	return ExecutionReport{
		SourceChainSelector: bcs.SourceChainSelector,
		Message:             MessageField,
		OffchainTokenData:   bcs.OffchainTokenData,
		Proofs:              bcs.Proofs,
	}, nil
}

type bcsTokenPriceUpdate struct {
	SourceToken [32]byte
	UsdPerToken [32]byte
}

func convertTokenPriceUpdateFromBCS(bcs bcsTokenPriceUpdate) (TokenPriceUpdate, error) {
	UsdPerTokenField, err := bind.DecodeU256Value(bcs.UsdPerToken)
	if err != nil {
		return TokenPriceUpdate{}, fmt.Errorf("failed to decode u256 field UsdPerToken: %w", err)
	}

	return TokenPriceUpdate{
		SourceToken: fmt.Sprintf("0x%x", bcs.SourceToken),
		UsdPerToken: UsdPerTokenField,
	}, nil
}

type bcsGasPriceUpdate struct {
	DestChainSelector uint64
	UsdPerUnitGas     [32]byte
}

func convertGasPriceUpdateFromBCS(bcs bcsGasPriceUpdate) (GasPriceUpdate, error) {
	UsdPerUnitGasField, err := bind.DecodeU256Value(bcs.UsdPerUnitGas)
	if err != nil {
		return GasPriceUpdate{}, fmt.Errorf("failed to decode u256 field UsdPerUnitGas: %w", err)
	}

	return GasPriceUpdate{
		DestChainSelector: bcs.DestChainSelector,
		UsdPerUnitGas:     UsdPerUnitGasField,
	}, nil
}

type bcsStaticConfig struct {
	ChainSelector      uint64
	RmnRemote          [32]byte
	TokenAdminRegistry [32]byte
	NonceManager       [32]byte
}

func convertStaticConfigFromBCS(bcs bcsStaticConfig) (StaticConfig, error) {

	return StaticConfig{
		ChainSelector:      bcs.ChainSelector,
		RmnRemote:          fmt.Sprintf("0x%x", bcs.RmnRemote),
		TokenAdminRegistry: fmt.Sprintf("0x%x", bcs.TokenAdminRegistry),
		NonceManager:       fmt.Sprintf("0x%x", bcs.NonceManager),
	}, nil
}

type bcsDynamicConfig struct {
	FeeQuoter                               [32]byte
	PermissionlessExecutionThresholdSeconds uint32
}

func convertDynamicConfigFromBCS(bcs bcsDynamicConfig) (DynamicConfig, error) {

	return DynamicConfig{
		FeeQuoter:                               fmt.Sprintf("0x%x", bcs.FeeQuoter),
		PermissionlessExecutionThresholdSeconds: bcs.PermissionlessExecutionThresholdSeconds,
	}, nil
}

type bcsDynamicConfigSet struct {
	DynamicConfig bcsDynamicConfig
}

func convertDynamicConfigSetFromBCS(bcs bcsDynamicConfigSet) (DynamicConfigSet, error) {
	DynamicConfigField, err := convertDynamicConfigFromBCS(bcs.DynamicConfig)
	if err != nil {
		return DynamicConfigSet{}, fmt.Errorf("failed to convert nested struct DynamicConfig: %w", err)
	}

	return DynamicConfigSet{
		DynamicConfig: DynamicConfigField,
	}, nil
}

type bcsSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceChainConfig   bcsSourceChainConfig
}

func convertSourceChainConfigSetFromBCS(bcs bcsSourceChainConfigSet) (SourceChainConfigSet, error) {
	SourceChainConfigField, err := convertSourceChainConfigFromBCS(bcs.SourceChainConfig)
	if err != nil {
		return SourceChainConfigSet{}, fmt.Errorf("failed to convert nested struct SourceChainConfig: %w", err)
	}

	return SourceChainConfigSet{
		SourceChainSelector: bcs.SourceChainSelector,
		SourceChainConfig:   SourceChainConfigField,
	}, nil
}

func init() {
	bind.RegisterStructDecoder("test::offramp::CCIPObjectRef", func(data []byte) (interface{}, error) {
		var result CCIPObjectRef
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::OffRampState", func(data []byte) (interface{}, error) {
		var temp bcsOffRampState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOffRampStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::OffRampStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsOffRampStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOffRampStatePointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::SourceChainConfig", func(data []byte) (interface{}, error) {
		var temp bcsSourceChainConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSourceChainConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::RampMessageHeader", func(data []byte) (interface{}, error) {
		var result RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::Any2SuiRampMessage", func(data []byte) (interface{}, error) {
		var temp bcsAny2SuiRampMessage
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAny2SuiRampMessageFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::Any2SuiTokenTransfer", func(data []byte) (interface{}, error) {
		var temp bcsAny2SuiTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertAny2SuiTokenTransferFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::ExecutionReport", func(data []byte) (interface{}, error) {
		var temp bcsExecutionReport
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertExecutionReportFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::CommitReport", func(data []byte) (interface{}, error) {
		var result CommitReport
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::PriceUpdates", func(data []byte) (interface{}, error) {
		var result PriceUpdates
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::TokenPriceUpdate", func(data []byte) (interface{}, error) {
		var temp bcsTokenPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTokenPriceUpdateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::GasPriceUpdate", func(data []byte) (interface{}, error) {
		var temp bcsGasPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertGasPriceUpdateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::MerkleRoot", func(data []byte) (interface{}, error) {
		var result MerkleRoot
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::StaticConfig", func(data []byte) (interface{}, error) {
		var temp bcsStaticConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertStaticConfigFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::DynamicConfig", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("test::offramp::StaticConfigSet", func(data []byte) (interface{}, error) {
		var result StaticConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::DynamicConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsDynamicConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertDynamicConfigSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::SourceChainConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsSourceChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertSourceChainConfigSetFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::SkippedAlreadyExecuted", func(data []byte) (interface{}, error) {
		var result SkippedAlreadyExecuted
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::ExecutionStateChanged", func(data []byte) (interface{}, error) {
		var result ExecutionStateChanged
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::CommitReportAccepted", func(data []byte) (interface{}, error) {
		var result CommitReportAccepted
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::SkippedReportExecution", func(data []byte) (interface{}, error) {
		var result SkippedReportExecution
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("test::offramp::OFFRAMP", func(data []byte) (interface{}, error) {
		var result OFFRAMP
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// EmitStaticConfigSetEvent executes the emit_static_config_set_event Move function.
func (c *OfframpContract) EmitStaticConfigSetEvent(ctx context.Context, opts *bind.CallOpts, chainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitStaticConfigSetEvent(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitDynamicConfigSetEvent executes the emit_dynamic_config_set_event Move function.
func (c *OfframpContract) EmitDynamicConfigSetEvent(ctx context.Context, opts *bind.CallOpts, dynamicConfig DynamicConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitDynamicConfigSetEvent(dynamicConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitSourceChainConfigSetEvent executes the emit_source_chain_config_set_event Move function.
func (c *OfframpContract) EmitSourceChainConfigSetEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, sourceChainConfig SourceChainConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitSourceChainConfigSetEvent(sourceChainSelector, sourceChainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitSkippedAlreadyExecutedEvent executes the emit_skipped_already_executed_event Move function.
func (c *OfframpContract) EmitSkippedAlreadyExecutedEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitSkippedAlreadyExecutedEvent(sourceChainSelector, sequenceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitExecutionStateChangedEvent executes the emit_execution_state_changed_event Move function.
func (c *OfframpContract) EmitExecutionStateChangedEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, sequenceNumber uint64, messageId []byte, messageHash []byte, state byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitExecutionStateChangedEvent(sourceChainSelector, sequenceNumber, messageId, messageHash, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitCommitReportAcceptedEvent executes the emit_commit_report_accepted_event Move function.
func (c *OfframpContract) EmitCommitReportAcceptedEvent(ctx context.Context, opts *bind.CallOpts, blessedMerkleRoots []MerkleRoot, unblessedMerkleRoots []MerkleRoot, priceUpdates PriceUpdates) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitCommitReportAcceptedEvent(blessedMerkleRoots, unblessedMerkleRoots, priceUpdates)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitSkippedReportExecutionEvent executes the emit_skipped_report_execution_event Move function.
func (c *OfframpContract) EmitSkippedReportExecutionEvent(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitSkippedReportExecutionEvent(sourceChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// EmitOcrConfigEvent executes the emit_ocr_config_event Move function.
func (c *OfframpContract) EmitOcrConfigEvent(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.EmitOcrConfigEvent()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function.
func (c *OfframpContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AddPackageId executes the add_package_id Move function.
func (c *OfframpContract) AddPackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.AddPackageId(state, packageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// RemovePackageId executes the remove_package_id Move function.
func (c *OfframpContract) RemovePackageId(ctx context.Context, opts *bind.CallOpts, state bind.Object, packageId string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.RemovePackageId(state, packageId)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllSourceChainConfigs executes the get_all_source_chain_configs Move function.
func (c *OfframpContract) GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetAllSourceChainConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// InitExecute executes the init_execute Move function.
func (c *OfframpContract) InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.InitExecute(ref, state, clock, reportContext, report)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// FinishExecute executes the finish_execute Move function.
func (c *OfframpContract) FinishExecute(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.FinishExecute()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *OfframpDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.offrampEncoder.TypeAndVersion()
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

// GetAllSourceChainConfigs executes the get_all_source_chain_configs Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<u64>
//	[1]: vector<SourceChainConfig>
func (d *OfframpDevInspect) GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetAllSourceChainConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

type offrampEncoder struct {
	*bind.BoundContract
}

// EmitStaticConfigSetEvent encodes a call to the emit_static_config_set_event Move function.
func (c offrampEncoder) EmitStaticConfigSetEvent(chainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_static_config_set_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		chainSelector,
	}, nil)
}

// EmitStaticConfigSetEventWithArgs encodes a call to the emit_static_config_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitStaticConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_static_config_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitDynamicConfigSetEvent encodes a call to the emit_dynamic_config_set_event Move function.
func (c offrampEncoder) EmitDynamicConfigSetEvent(dynamicConfig DynamicConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dynamic_config_set_event", typeArgsList, typeParamsList, []string{
		"test::offramp::DynamicConfig",
	}, []any{
		dynamicConfig,
	}, nil)
}

// EmitDynamicConfigSetEventWithArgs encodes a call to the emit_dynamic_config_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitDynamicConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"test::offramp::DynamicConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_dynamic_config_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitSourceChainConfigSetEvent encodes a call to the emit_source_chain_config_set_event Move function.
func (c offrampEncoder) EmitSourceChainConfigSetEvent(sourceChainSelector uint64, sourceChainConfig SourceChainConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_source_chain_config_set_event", typeArgsList, typeParamsList, []string{
		"u64",
		"test::offramp::SourceChainConfig",
	}, []any{
		sourceChainSelector,
		sourceChainConfig,
	}, nil)
}

// EmitSourceChainConfigSetEventWithArgs encodes a call to the emit_source_chain_config_set_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitSourceChainConfigSetEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"test::offramp::SourceChainConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_source_chain_config_set_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitSkippedAlreadyExecutedEvent encodes a call to the emit_skipped_already_executed_event Move function.
func (c offrampEncoder) EmitSkippedAlreadyExecutedEvent(sourceChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_skipped_already_executed_event", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
	}, []any{
		sourceChainSelector,
		sequenceNumber,
	}, nil)
}

// EmitSkippedAlreadyExecutedEventWithArgs encodes a call to the emit_skipped_already_executed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitSkippedAlreadyExecutedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_skipped_already_executed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitExecutionStateChangedEvent encodes a call to the emit_execution_state_changed_event Move function.
func (c offrampEncoder) EmitExecutionStateChangedEvent(sourceChainSelector uint64, sequenceNumber uint64, messageId []byte, messageHash []byte, state byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_execution_state_changed_event", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"u8",
	}, []any{
		sourceChainSelector,
		sequenceNumber,
		messageId,
		messageHash,
		state,
	}, nil)
}

// EmitExecutionStateChangedEventWithArgs encodes a call to the emit_execution_state_changed_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitExecutionStateChangedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"u8",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_execution_state_changed_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitCommitReportAcceptedEvent encodes a call to the emit_commit_report_accepted_event Move function.
func (c offrampEncoder) EmitCommitReportAcceptedEvent(blessedMerkleRoots []MerkleRoot, unblessedMerkleRoots []MerkleRoot, priceUpdates PriceUpdates) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_commit_report_accepted_event", typeArgsList, typeParamsList, []string{
		"vector<test::offramp::MerkleRoot>",
		"vector<test::offramp::MerkleRoot>",
		"test::offramp::PriceUpdates",
	}, []any{
		blessedMerkleRoots,
		unblessedMerkleRoots,
		priceUpdates,
	}, nil)
}

// EmitCommitReportAcceptedEventWithArgs encodes a call to the emit_commit_report_accepted_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitCommitReportAcceptedEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<test::offramp::MerkleRoot>",
		"vector<test::offramp::MerkleRoot>",
		"test::offramp::PriceUpdates",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_commit_report_accepted_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitSkippedReportExecutionEvent encodes a call to the emit_skipped_report_execution_event Move function.
func (c offrampEncoder) EmitSkippedReportExecutionEvent(sourceChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_skipped_report_execution_event", typeArgsList, typeParamsList, []string{
		"u64",
	}, []any{
		sourceChainSelector,
	}, nil)
}

// EmitSkippedReportExecutionEventWithArgs encodes a call to the emit_skipped_report_execution_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitSkippedReportExecutionEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_skipped_report_execution_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// EmitOcrConfigEvent encodes a call to the emit_ocr_config_event Move function.
func (c offrampEncoder) EmitOcrConfigEvent() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_ocr_config_event", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// EmitOcrConfigEventWithArgs encodes a call to the emit_ocr_config_event Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) EmitOcrConfigEventWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("emit_ocr_config_event", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c offrampEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// AddPackageId encodes a call to the add_package_id Move function.
func (c offrampEncoder) AddPackageId(state bind.Object, packageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_package_id", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"address",
	}, []any{
		state,
		packageId,
	}, nil)
}

// AddPackageIdWithArgs encodes a call to the add_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) AddPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("add_package_id", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// RemovePackageId encodes a call to the remove_package_id Move function.
func (c offrampEncoder) RemovePackageId(state bind.Object, packageId string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_package_id", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"address",
	}, []any{
		state,
		packageId,
	}, nil)
}

// RemovePackageIdWithArgs encodes a call to the remove_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) RemovePackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("remove_package_id", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetAllSourceChainConfigs encodes a call to the get_all_source_chain_configs Move function.
func (c offrampEncoder) GetAllSourceChainConfigs() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_source_chain_configs", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"vector<u64>",
		"vector<test::offramp::SourceChainConfig>",
	})
}

// GetAllSourceChainConfigsWithArgs encodes a call to the get_all_source_chain_configs Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetAllSourceChainConfigsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_source_chain_configs", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
		"vector<test::offramp::SourceChainConfig>",
	})
}

// InitExecute encodes a call to the init_execute Move function.
func (c offrampEncoder) InitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("init_execute", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
	}, []any{
		ref,
		state,
		clock,
		reportContext,
		report,
	}, nil)
}

// InitExecuteWithArgs encodes a call to the init_execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) InitExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("init_execute", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// FinishExecute encodes a call to the finish_execute Move function.
func (c offrampEncoder) FinishExecute() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("finish_execute", typeArgsList, typeParamsList, []string{}, []any{}, nil)
}

// FinishExecuteWithArgs encodes a call to the finish_execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) FinishExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("finish_execute", typeArgsList, typeParamsList, expectedParams, args, nil)
}
