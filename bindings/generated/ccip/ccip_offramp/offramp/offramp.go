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
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*models.SuiTransactionBlockResponse, error)
	FinishExecute(ctx context.Context, opts *bind.CallOpts, state bind.Object, receiverParams bind.Object, completedTransfers []bind.Object) (*models.SuiTransactionBlockResponse, error)
	ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*models.SuiTransactionBlockResponse, error)
	GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*models.SuiTransactionBlockResponse, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*models.SuiTransactionBlockResponse, error)
	SetOcr3Config(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*models.SuiTransactionBlockResponse, error)
	ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	Commit(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (*models.SuiTransactionBlockResponse, error)
	GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, sourceChainConfig SourceChainConfig) (*models.SuiTransactionBlockResponse, error)
	GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error)
	SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, permissionlessExecutionThresholdSeconds uint32) (*models.SuiTransactionBlockResponse, error)
	ApplySourceChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipAsMcms(ctx context.Context, opts *bind.CallOpts, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOfframpDevInspect
	Encoder() OfframpEncoder
}

type IOfframpDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bind.Object, error)
	InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (bind.Object, error)
	ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (bind.Object, error)
	GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (byte, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) ([]byte, error)
	ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([][]byte, error)
	ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]string, error)
	GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (uint64, error)
	GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64) (SourceChainConfig, error)
	GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, sourceChainConfig SourceChainConfig) ([]any, error)
	GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]any, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (StaticConfig, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) ([]any, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (DynamicConfig, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) ([]any, error)
	GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type OfframpEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(state bind.Object, param bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOcr3Base(state bind.Object) (*bind.EncodedCall, error)
	GetOcr3BaseWithArgs(args ...any) (*bind.EncodedCall, error)
	InitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (*bind.EncodedCall, error)
	InitExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	FinishExecute(state bind.Object, receiverParams bind.Object, completedTransfers []bind.Object) (*bind.EncodedCall, error)
	FinishExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	ManuallyInitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*bind.EncodedCall, error)
	ManuallyInitExecuteWithArgs(args ...any) (*bind.EncodedCall, error)
	GetExecutionState(state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error)
	GetExecutionStateWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMetadataHash(sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*bind.EncodedCall, error)
	CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMessageHash(messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*bind.EncodedCall, error)
	CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error)
	SetOcr3Config(state bind.Object, param bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*bind.EncodedCall, error)
	SetOcr3ConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	ConfigSigners(state bind.Object) (*bind.EncodedCall, error)
	ConfigSignersWithArgs(args ...any) (*bind.EncodedCall, error)
	ConfigTransmitters(state bind.Object) (*bind.EncodedCall, error)
	ConfigTransmittersWithArgs(args ...any) (*bind.EncodedCall, error)
	Commit(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*bind.EncodedCall, error)
	CommitWithArgs(args ...any) (*bind.EncodedCall, error)
	GetMerkleRoot(state bind.Object, root []byte) (*bind.EncodedCall, error)
	GetMerkleRootWithArgs(args ...any) (*bind.EncodedCall, error)
	GetSourceChainConfig(state bind.Object, sourceChainSelector uint64) (*bind.EncodedCall, error)
	GetSourceChainConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetSourceChainConfigFields(sourceChainConfig SourceChainConfig) (*bind.EncodedCall, error)
	GetSourceChainConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllSourceChainConfigs(state bind.Object) (*bind.EncodedCall, error)
	GetAllSourceChainConfigsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfig(state bind.Object) (*bind.EncodedCall, error)
	GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error)
	GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfig(state bind.Object) (*bind.EncodedCall, error)
	GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfigFields(cfg DynamicConfig) (*bind.EncodedCall, error)
	GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	SetDynamicConfig(state bind.Object, param bind.Object, permissionlessExecutionThresholdSeconds uint32) (*bind.EncodedCall, error)
	SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplySourceChainConfigUpdates(state bind.Object, param bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error)
	ApplySourceChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetCcipPackageId() (*bind.EncodedCall, error)
	GetCcipPackageIdWithArgs(args ...any) (*bind.EncodedCall, error)
	Owner(state bind.Object) (*bind.EncodedCall, error)
	OwnerWithArgs(args ...any) (*bind.EncodedCall, error)
	HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error)
	HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferTo(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error)
	PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error)
	PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error)
	TransferOwnership(state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error)
	TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnership(state bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error)
	AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error)
	AcceptOwnershipAsMcms(state bind.Object, params bind.Object) (*bind.EncodedCall, error)
	AcceptOwnershipAsMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewOfframp(packageID string, client sui.ISuiAPI) (*OfframpContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_offramp", "offramp", client)
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

func (c *OfframpContract) Encoder() OfframpEncoder {
	return c.offrampEncoder
}

func (c *OfframpContract) DevInspect() IOfframpDevInspect {
	return c.devInspect
}

type OffRampState struct {
	Id                                      string       `move:"sui::object::UID"`
	Ocr3BaseState                           bind.Object  `move:"OCR3BaseState"`
	ChainSelector                           uint64       `move:"u64"`
	PermissionlessExecutionThresholdSeconds uint32       `move:"u32"`
	SourceChainConfigs                      bind.Object  `move:"VecMap<u64, SourceChainConfig>"`
	ExecutionStates                         bind.Object  `move:"Table<u64, Table<u64, u8>>"`
	Roots                                   bind.Object  `move:"Table<vector<u8>, u64>"`
	LatestPriceSequenceNumber               uint64       `move:"u64"`
	FeeQuoterCap                            *bind.Object `move:"0x1::option::Option<FeeQuoterCap>"`
	DestTransferCap                         *bind.Object `move:"0x1::option::Option<osh::DestTransferCap>"`
	OwnableState                            bind.Object  `move:"OwnableState"`
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

type McmsCallback struct {
}

type bcsOffRampStatePointer struct {
	Id             string
	OffRampStateId [32]byte
	OwnerCapId     [32]byte
}

func convertOffRampStatePointerFromBCS(bcs bcsOffRampStatePointer) OffRampStatePointer {
	return OffRampStatePointer{
		Id:             bcs.Id,
		OffRampStateId: fmt.Sprintf("0x%x", bcs.OffRampStateId),
		OwnerCapId:     fmt.Sprintf("0x%x", bcs.OwnerCapId),
	}
}

type bcsSourceChainConfig struct {
	Router                    [32]byte
	IsEnabled                 bool
	MinSeqNr                  uint64
	IsRmnVerificationDisabled bool
	OnRamp                    []byte
}

func convertSourceChainConfigFromBCS(bcs bcsSourceChainConfig) SourceChainConfig {
	return SourceChainConfig{
		Router:                    fmt.Sprintf("0x%x", bcs.Router),
		IsEnabled:                 bcs.IsEnabled,
		MinSeqNr:                  bcs.MinSeqNr,
		IsRmnVerificationDisabled: bcs.IsRmnVerificationDisabled,
		OnRamp:                    bcs.OnRamp,
	}
}

type bcsAny2SuiRampMessage struct {
	Header       RampMessageHeader
	Sender       []byte
	Data         []byte
	Receiver     [32]byte
	GasLimit     *big.Int
	TokenAmounts []Any2SuiTokenTransfer
}

func convertAny2SuiRampMessageFromBCS(bcs bcsAny2SuiRampMessage) Any2SuiRampMessage {
	return Any2SuiRampMessage{
		Header:       bcs.Header,
		Sender:       bcs.Sender,
		Data:         bcs.Data,
		Receiver:     fmt.Sprintf("0x%x", bcs.Receiver),
		GasLimit:     bcs.GasLimit,
		TokenAmounts: bcs.TokenAmounts,
	}
}

type bcsAny2SuiTokenTransfer struct {
	SourcePoolAddress []byte
	DestTokenAddress  [32]byte
	DestGasAmount     uint32
	ExtraData         []byte
	Amount            *big.Int
}

func convertAny2SuiTokenTransferFromBCS(bcs bcsAny2SuiTokenTransfer) Any2SuiTokenTransfer {
	return Any2SuiTokenTransfer{
		SourcePoolAddress: bcs.SourcePoolAddress,
		DestTokenAddress:  fmt.Sprintf("0x%x", bcs.DestTokenAddress),
		DestGasAmount:     bcs.DestGasAmount,
		ExtraData:         bcs.ExtraData,
		Amount:            bcs.Amount,
	}
}

type bcsExecutionReport struct {
	SourceChainSelector uint64
	Message             bcsAny2SuiRampMessage
	OffchainTokenData   [][]byte
	Proofs              [][]byte
}

func convertExecutionReportFromBCS(bcs bcsExecutionReport) ExecutionReport {
	return ExecutionReport{
		SourceChainSelector: bcs.SourceChainSelector,
		Message:             convertAny2SuiRampMessageFromBCS(bcs.Message),
		OffchainTokenData:   bcs.OffchainTokenData,
		Proofs:              bcs.Proofs,
	}
}

type bcsTokenPriceUpdate struct {
	SourceToken [32]byte
	UsdPerToken *big.Int
}

func convertTokenPriceUpdateFromBCS(bcs bcsTokenPriceUpdate) TokenPriceUpdate {
	return TokenPriceUpdate{
		SourceToken: fmt.Sprintf("0x%x", bcs.SourceToken),
		UsdPerToken: bcs.UsdPerToken,
	}
}

type bcsStaticConfig struct {
	ChainSelector      uint64
	RmnRemote          [32]byte
	TokenAdminRegistry [32]byte
	NonceManager       [32]byte
}

func convertStaticConfigFromBCS(bcs bcsStaticConfig) StaticConfig {
	return StaticConfig{
		ChainSelector:      bcs.ChainSelector,
		RmnRemote:          fmt.Sprintf("0x%x", bcs.RmnRemote),
		TokenAdminRegistry: fmt.Sprintf("0x%x", bcs.TokenAdminRegistry),
		NonceManager:       fmt.Sprintf("0x%x", bcs.NonceManager),
	}
}

type bcsDynamicConfig struct {
	FeeQuoter                               [32]byte
	PermissionlessExecutionThresholdSeconds uint32
}

func convertDynamicConfigFromBCS(bcs bcsDynamicConfig) DynamicConfig {
	return DynamicConfig{
		FeeQuoter:                               fmt.Sprintf("0x%x", bcs.FeeQuoter),
		PermissionlessExecutionThresholdSeconds: bcs.PermissionlessExecutionThresholdSeconds,
	}
}

type bcsDynamicConfigSet struct {
	DynamicConfig bcsDynamicConfig
}

func convertDynamicConfigSetFromBCS(bcs bcsDynamicConfigSet) DynamicConfigSet {
	return DynamicConfigSet{
		DynamicConfig: convertDynamicConfigFromBCS(bcs.DynamicConfig),
	}
}

type bcsSourceChainConfigSet struct {
	SourceChainSelector uint64
	SourceChainConfig   bcsSourceChainConfig
}

func convertSourceChainConfigSetFromBCS(bcs bcsSourceChainConfigSet) SourceChainConfigSet {
	return SourceChainConfigSet{
		SourceChainSelector: bcs.SourceChainSelector,
		SourceChainConfig:   convertSourceChainConfigFromBCS(bcs.SourceChainConfig),
	}
}

func init() {
	bind.RegisterStructDecoder("ccip_offramp::offramp::OffRampState", func(data []byte) (interface{}, error) {
		var result OffRampState
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::OffRampStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsOffRampStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertOffRampStatePointerFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SourceChainConfig", func(data []byte) (interface{}, error) {
		var temp bcsSourceChainConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertSourceChainConfigFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::RampMessageHeader", func(data []byte) (interface{}, error) {
		var result RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::Any2SuiRampMessage", func(data []byte) (interface{}, error) {
		var temp bcsAny2SuiRampMessage
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertAny2SuiRampMessageFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::Any2SuiTokenTransfer", func(data []byte) (interface{}, error) {
		var temp bcsAny2SuiTokenTransfer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertAny2SuiTokenTransferFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::ExecutionReport", func(data []byte) (interface{}, error) {
		var temp bcsExecutionReport
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertExecutionReportFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::CommitReport", func(data []byte) (interface{}, error) {
		var result CommitReport
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::PriceUpdates", func(data []byte) (interface{}, error) {
		var result PriceUpdates
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::TokenPriceUpdate", func(data []byte) (interface{}, error) {
		var temp bcsTokenPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertTokenPriceUpdateFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::GasPriceUpdate", func(data []byte) (interface{}, error) {
		var result GasPriceUpdate
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::MerkleRoot", func(data []byte) (interface{}, error) {
		var result MerkleRoot
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::StaticConfig", func(data []byte) (interface{}, error) {
		var temp bcsStaticConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertStaticConfigFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::DynamicConfig", func(data []byte) (interface{}, error) {
		var temp bcsDynamicConfig
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertDynamicConfigFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::StaticConfigSet", func(data []byte) (interface{}, error) {
		var result StaticConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::DynamicConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsDynamicConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertDynamicConfigSetFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SourceChainConfigSet", func(data []byte) (interface{}, error) {
		var temp bcsSourceChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result := convertSourceChainConfigSetFromBCS(temp)
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SkippedAlreadyExecuted", func(data []byte) (interface{}, error) {
		var result SkippedAlreadyExecuted
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::ExecutionStateChanged", func(data []byte) (interface{}, error) {
		var result ExecutionStateChanged
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::CommitReportAccepted", func(data []byte) (interface{}, error) {
		var result CommitReportAccepted
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::SkippedReportExecution", func(data []byte) (interface{}, error) {
		var result SkippedReportExecution
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::OFFRAMP", func(data []byte) (interface{}, error) {
		var result OFFRAMP
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_offramp::offramp::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *OfframpContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *OfframpContract) Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.Initialize(state, param, feeQuoterCap, destTransferCap, chainSelector, permissionlessExecutionThresholdSeconds, sourceChainsSelectors, sourceChainsIsEnabled, sourceChainsIsRmnVerificationDisabled, sourceChainsOnRamp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOcr3Base executes the get_ocr3_base Move function.
func (c *OfframpContract) GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetOcr3Base(state)
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
func (c *OfframpContract) FinishExecute(ctx context.Context, opts *bind.CallOpts, state bind.Object, receiverParams bind.Object, completedTransfers []bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.FinishExecute(state, receiverParams, completedTransfers)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ManuallyInitExecute executes the manually_init_execute Move function.
func (c *OfframpContract) ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ManuallyInitExecute(ref, state, clock, reportBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetExecutionState executes the get_execution_state Move function.
func (c *OfframpContract) GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetExecutionState(state, sourceChainSelector, sequenceNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMetadataHash executes the calculate_metadata_hash Move function.
func (c *OfframpContract) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.CalculateMetadataHash(sourceChainSelector, destChainSelector, onRamp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMessageHash executes the calculate_message_hash Move function.
func (c *OfframpContract) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.CalculateMessageHash(messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, onRamp, data, gasLimit, sourcePoolAddresses, destTokenAddresses, destGasAmounts, extraDatas, amounts)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetOcr3Config executes the set_ocr3_config Move function.
func (c *OfframpContract) SetOcr3Config(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.SetOcr3Config(state, param, configDigest, ocrPluginType, bigF, isSignatureVerificationEnabled, signers, transmitters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigSigners executes the config_signers Move function.
func (c *OfframpContract) ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ConfigSigners(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConfigTransmitters executes the config_transmitters Move function.
func (c *OfframpContract) ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ConfigTransmitters(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Commit executes the commit Move function.
func (c *OfframpContract) Commit(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.Commit(ref, state, clock, reportContext, report, signatures)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetMerkleRoot executes the get_merkle_root Move function.
func (c *OfframpContract) GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetMerkleRoot(state, root)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSourceChainConfig executes the get_source_chain_config Move function.
func (c *OfframpContract) GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetSourceChainConfig(state, sourceChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetSourceChainConfigFields executes the get_source_chain_config_fields Move function.
func (c *OfframpContract) GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, sourceChainConfig SourceChainConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetSourceChainConfigFields(sourceChainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllSourceChainConfigs executes the get_all_source_chain_configs Move function.
func (c *OfframpContract) GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetAllSourceChainConfigs(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function.
func (c *OfframpContract) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetStaticConfig(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfigFields executes the get_static_config_fields Move function.
func (c *OfframpContract) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetStaticConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfig executes the get_dynamic_config Move function.
func (c *OfframpContract) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetDynamicConfig(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfigFields executes the get_dynamic_config_fields Move function.
func (c *OfframpContract) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetDynamicConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetDynamicConfig executes the set_dynamic_config Move function.
func (c *OfframpContract) SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, permissionlessExecutionThresholdSeconds uint32) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.SetDynamicConfig(state, param, permissionlessExecutionThresholdSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplySourceChainConfigUpdates executes the apply_source_chain_config_updates Move function.
func (c *OfframpContract) ApplySourceChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ApplySourceChainConfigUpdates(state, param, sourceChainsSelector, sourceChainsIsEnabled, sourceChainsIsRmnVerificationDisabled, sourceChainsOnRamp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCcipPackageId executes the get_ccip_package_id Move function.
func (c *OfframpContract) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.GetCcipPackageId()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *OfframpContract) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.Owner(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *OfframpContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.HasPendingTransfer(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *OfframpContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *OfframpContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *OfframpContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *OfframpContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.TransferOwnership(state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *OfframpContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.AcceptOwnership(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *OfframpContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.AcceptOwnershipFromObject(state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipAsMcms executes the accept_ownership_as_mcms Move function.
func (c *OfframpContract) AcceptOwnershipAsMcms(ctx context.Context, opts *bind.CallOpts, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.AcceptOwnershipAsMcms(state, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *OfframpContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ExecuteOwnershipTransfer(ownerCap, ownableState, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *OfframpContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.ExecuteOwnershipTransferToMcms(ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *OfframpContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *OfframpContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.offrampEncoder.McmsEntrypoint(state, registry, params)
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

// GetOcr3Base executes the get_ocr3_base Move function using DevInspect to get return values.
//
// Returns: &OCR3BaseState
func (d *OfframpDevInspect) GetOcr3Base(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bind.Object, error) {
	encoded, err := d.contract.offrampEncoder.GetOcr3Base(state)
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

// InitExecute executes the init_execute Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *OfframpDevInspect) InitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte) (bind.Object, error) {
	encoded, err := d.contract.offrampEncoder.InitExecute(ref, state, clock, reportContext, report)
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

// ManuallyInitExecute executes the manually_init_execute Move function using DevInspect to get return values.
//
// Returns: osh::ReceiverParams
func (d *OfframpDevInspect) ManuallyInitExecute(ctx context.Context, opts *bind.CallOpts, ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (bind.Object, error) {
	encoded, err := d.contract.offrampEncoder.ManuallyInitExecute(ref, state, clock, reportBytes)
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

// GetExecutionState executes the get_execution_state Move function using DevInspect to get return values.
//
// Returns: u8
func (d *OfframpDevInspect) GetExecutionState(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (byte, error) {
	encoded, err := d.contract.offrampEncoder.GetExecutionState(state, sourceChainSelector, sequenceNumber)
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
	result, ok := results[0].(byte)
	if !ok {
		return 0, fmt.Errorf("unexpected return type: expected byte, got %T", results[0])
	}
	return result, nil
}

// CalculateMetadataHash executes the calculate_metadata_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OfframpDevInspect) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) ([]byte, error) {
	encoded, err := d.contract.offrampEncoder.CalculateMetadataHash(sourceChainSelector, destChainSelector, onRamp)
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
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// CalculateMessageHash executes the calculate_message_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OfframpDevInspect) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) ([]byte, error) {
	encoded, err := d.contract.offrampEncoder.CalculateMessageHash(messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, onRamp, data, gasLimit, sourcePoolAddresses, destTokenAddresses, destGasAmounts, extraDatas, amounts)
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
	result, ok := results[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []byte, got %T", results[0])
	}
	return result, nil
}

// ConfigSigners executes the config_signers Move function using DevInspect to get return values.
//
// Returns: vector<vector<u8>>
func (d *OfframpDevInspect) ConfigSigners(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([][]byte, error) {
	encoded, err := d.contract.offrampEncoder.ConfigSigners(state)
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
	result, ok := results[0].([][]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected [][]byte, got %T", results[0])
	}
	return result, nil
}

// ConfigTransmitters executes the config_transmitters Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *OfframpDevInspect) ConfigTransmitters(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]string, error) {
	encoded, err := d.contract.offrampEncoder.ConfigTransmitters(state)
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
	result, ok := results[0].([]string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []string, got %T", results[0])
	}
	return result, nil
}

// GetMerkleRoot executes the get_merkle_root Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OfframpDevInspect) GetMerkleRoot(ctx context.Context, opts *bind.CallOpts, state bind.Object, root []byte) (uint64, error) {
	encoded, err := d.contract.offrampEncoder.GetMerkleRoot(state, root)
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

// GetSourceChainConfig executes the get_source_chain_config Move function using DevInspect to get return values.
//
// Returns: SourceChainConfig
func (d *OfframpDevInspect) GetSourceChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, sourceChainSelector uint64) (SourceChainConfig, error) {
	encoded, err := d.contract.offrampEncoder.GetSourceChainConfig(state, sourceChainSelector)
	if err != nil {
		return SourceChainConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return SourceChainConfig{}, err
	}
	if len(results) == 0 {
		return SourceChainConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(SourceChainConfig)
	if !ok {
		return SourceChainConfig{}, fmt.Errorf("unexpected return type: expected SourceChainConfig, got %T", results[0])
	}
	return result, nil
}

// GetSourceChainConfigFields executes the get_source_chain_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: bool
//	[2]: u64
//	[3]: bool
//	[4]: vector<u8>
func (d *OfframpDevInspect) GetSourceChainConfigFields(ctx context.Context, opts *bind.CallOpts, sourceChainConfig SourceChainConfig) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetSourceChainConfigFields(sourceChainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetAllSourceChainConfigs executes the get_all_source_chain_configs Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: vector<u64>
//	[1]: vector<SourceChainConfig>
func (d *OfframpDevInspect) GetAllSourceChainConfigs(ctx context.Context, opts *bind.CallOpts, state bind.Object) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetAllSourceChainConfigs(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function using DevInspect to get return values.
//
// Returns: StaticConfig
func (d *OfframpDevInspect) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (StaticConfig, error) {
	encoded, err := d.contract.offrampEncoder.GetStaticConfig(state)
	if err != nil {
		return StaticConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return StaticConfig{}, err
	}
	if len(results) == 0 {
		return StaticConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(StaticConfig)
	if !ok {
		return StaticConfig{}, fmt.Errorf("unexpected return type: expected StaticConfig, got %T", results[0])
	}
	return result, nil
}

// GetStaticConfigFields executes the get_static_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u64
//	[1]: address
//	[2]: address
//	[3]: address
func (d *OfframpDevInspect) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetStaticConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetDynamicConfig executes the get_dynamic_config Move function using DevInspect to get return values.
//
// Returns: DynamicConfig
func (d *OfframpDevInspect) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (DynamicConfig, error) {
	encoded, err := d.contract.offrampEncoder.GetDynamicConfig(state)
	if err != nil {
		return DynamicConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return DynamicConfig{}, err
	}
	if len(results) == 0 {
		return DynamicConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(DynamicConfig)
	if !ok {
		return DynamicConfig{}, fmt.Errorf("unexpected return type: expected DynamicConfig, got %T", results[0])
	}
	return result, nil
}

// GetDynamicConfigFields executes the get_dynamic_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: address
//	[1]: u32
func (d *OfframpDevInspect) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) ([]any, error) {
	encoded, err := d.contract.offrampEncoder.GetDynamicConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetCcipPackageId executes the get_ccip_package_id Move function using DevInspect to get return values.
//
// Returns: address
func (d *OfframpDevInspect) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.offrampEncoder.GetCcipPackageId()
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

// Owner executes the owner Move function using DevInspect to get return values.
//
// Returns: address
func (d *OfframpDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.offrampEncoder.Owner(state)
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

// HasPendingTransfer executes the has_pending_transfer Move function using DevInspect to get return values.
//
// Returns: bool
func (d *OfframpDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error) {
	encoded, err := d.contract.offrampEncoder.HasPendingTransfer(state)
	if err != nil {
		return false, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return false, err
	}
	if len(results) == 0 {
		return false, fmt.Errorf("no return value")
	}
	result, ok := results[0].(bool)
	if !ok {
		return false, fmt.Errorf("unexpected return type: expected bool, got %T", results[0])
	}
	return result, nil
}

// PendingTransferFrom executes the pending_transfer_from Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *OfframpDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.offrampEncoder.PendingTransferFrom(state)
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
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferTo executes the pending_transfer_to Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<address>
func (d *OfframpDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.offrampEncoder.PendingTransferTo(state)
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
	result, ok := results[0].(*string)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *string, got %T", results[0])
	}
	return result, nil
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function using DevInspect to get return values.
//
// Returns: 0x1::option::Option<bool>
func (d *OfframpDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error) {
	encoded, err := d.contract.offrampEncoder.PendingTransferAccepted(state)
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
	result, ok := results[0].(*bool)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected *bool, got %T", results[0])
	}
	return result, nil
}

type offrampEncoder struct {
	*bind.BoundContract
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

// Initialize encodes a call to the initialize Move function.
func (c offrampEncoder) Initialize(state bind.Object, param bind.Object, feeQuoterCap bind.Object, destTransferCap bind.Object, chainSelector uint64, permissionlessExecutionThresholdSeconds uint32, sourceChainsSelectors []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"FeeQuoterCap",
		"osh::DestTransferCap",
		"u64",
		"u32",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}, []any{
		state,
		param,
		feeQuoterCap,
		destTransferCap,
		chainSelector,
		permissionlessExecutionThresholdSeconds,
		sourceChainsSelectors,
		sourceChainsIsEnabled,
		sourceChainsIsRmnVerificationDisabled,
		sourceChainsOnRamp,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"FeeQuoterCap",
		"osh::DestTransferCap",
		"u64",
		"u32",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetOcr3Base encodes a call to the get_ocr3_base Move function.
func (c offrampEncoder) GetOcr3Base(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ocr3_base", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"&OCR3BaseState",
	})
}

// GetOcr3BaseWithArgs encodes a call to the get_ocr3_base Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetOcr3BaseWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ocr3_base", typeArgsList, typeParamsList, expectedParams, args, []string{
		"&OCR3BaseState",
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
	}, []string{
		"osh::ReceiverParams",
	})
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
	return c.EncodeCallArgsWithGenerics("init_execute", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}

// FinishExecute encodes a call to the finish_execute Move function.
func (c offrampEncoder) FinishExecute(state bind.Object, receiverParams bind.Object, completedTransfers []bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("finish_execute", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"osh::ReceiverParams",
		"vector<osh::CompletedDestTokenTransfer>",
	}, []any{
		state,
		receiverParams,
		completedTransfers,
	}, nil)
}

// FinishExecuteWithArgs encodes a call to the finish_execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) FinishExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"osh::ReceiverParams",
		"vector<osh::CompletedDestTokenTransfer>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("finish_execute", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ManuallyInitExecute encodes a call to the manually_init_execute Move function.
func (c offrampEncoder) ManuallyInitExecute(ref bind.Object, state bind.Object, clock bind.Object, reportBytes []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("manually_init_execute", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<u8>",
	}, []any{
		ref,
		state,
		clock,
		reportBytes,
	}, []string{
		"osh::ReceiverParams",
	})
}

// ManuallyInitExecuteWithArgs encodes a call to the manually_init_execute Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ManuallyInitExecuteWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("manually_init_execute", typeArgsList, typeParamsList, expectedParams, args, []string{
		"osh::ReceiverParams",
	})
}

// GetExecutionState encodes a call to the get_execution_state Move function.
func (c offrampEncoder) GetExecutionState(state bind.Object, sourceChainSelector uint64, sequenceNumber uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_execution_state", typeArgsList, typeParamsList, []string{
		"&OffRampState",
		"u64",
		"u64",
	}, []any{
		state,
		sourceChainSelector,
		sequenceNumber,
	}, []string{
		"u8",
	})
}

// GetExecutionStateWithArgs encodes a call to the get_execution_state Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetExecutionStateWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
		"u64",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_execution_state", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u8",
	})
}

// CalculateMetadataHash encodes a call to the calculate_metadata_hash Move function.
func (c offrampEncoder) CalculateMetadataHash(sourceChainSelector uint64, destChainSelector uint64, onRamp []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_metadata_hash", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
		"vector<u8>",
	}, []any{
		sourceChainSelector,
		destChainSelector,
		onRamp,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMetadataHashWithArgs encodes a call to the calculate_metadata_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_metadata_hash", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// CalculateMessageHash encodes a call to the calculate_message_hash Move function.
func (c offrampEncoder) CalculateMessageHash(messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender []byte, receiver string, onRamp []byte, data []byte, gasLimit *big.Int, sourcePoolAddresses [][]byte, destTokenAddresses []string, destGasAmounts []uint32, extraDatas [][]byte, amounts []*big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_message_hash", typeArgsList, typeParamsList, []string{
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"vector<u8>",
		"address",
		"vector<u8>",
		"vector<u8>",
		"u256",
		"vector<vector<u8>>",
		"vector<address>",
		"vector<u32>",
		"vector<vector<u8>>",
		"vector<u256>",
	}, []any{
		messageId,
		sourceChainSelector,
		destChainSelector,
		sequenceNumber,
		nonce,
		sender,
		receiver,
		onRamp,
		data,
		gasLimit,
		sourcePoolAddresses,
		destTokenAddresses,
		destGasAmounts,
		extraDatas,
		amounts,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMessageHashWithArgs encodes a call to the calculate_message_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"vector<u8>",
		"address",
		"vector<u8>",
		"vector<u8>",
		"u256",
		"vector<vector<u8>>",
		"vector<address>",
		"vector<u32>",
		"vector<vector<u8>>",
		"vector<u256>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_message_hash", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// SetOcr3Config encodes a call to the set_ocr3_config Move function.
func (c offrampEncoder) SetOcr3Config(state bind.Object, param bind.Object, configDigest []byte, ocrPluginType byte, bigF byte, isSignatureVerificationEnabled bool, signers [][]byte, transmitters []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_ocr3_config", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u8>",
		"u8",
		"u8",
		"bool",
		"vector<vector<u8>>",
		"vector<address>",
	}, []any{
		state,
		param,
		configDigest,
		ocrPluginType,
		bigF,
		isSignatureVerificationEnabled,
		signers,
		transmitters,
	}, nil)
}

// SetOcr3ConfigWithArgs encodes a call to the set_ocr3_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) SetOcr3ConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u8>",
		"u8",
		"u8",
		"bool",
		"vector<vector<u8>>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_ocr3_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ConfigSigners encodes a call to the config_signers Move function.
func (c offrampEncoder) ConfigSigners(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_signers", typeArgsList, typeParamsList, []string{
		"&OCRConfig",
	}, []any{
		state,
	}, []string{
		"vector<vector<u8>>",
	})
}

// ConfigSignersWithArgs encodes a call to the config_signers Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ConfigSignersWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OCRConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_signers", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<vector<u8>>",
	})
}

// ConfigTransmitters encodes a call to the config_transmitters Move function.
func (c offrampEncoder) ConfigTransmitters(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_transmitters", typeArgsList, typeParamsList, []string{
		"&OCRConfig",
	}, []any{
		state,
	}, []string{
		"vector<address>",
	})
}

// ConfigTransmittersWithArgs encodes a call to the config_transmitters Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ConfigTransmittersWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OCRConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("config_transmitters", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// Commit encodes a call to the commit Move function.
func (c offrampEncoder) Commit(ref bind.Object, state bind.Object, clock bind.Object, reportContext [][]byte, report []byte, signatures [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("commit", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<vector<u8>>",
	}, []any{
		ref,
		state,
		clock,
		reportContext,
		report,
		signatures,
	}, nil)
}

// CommitWithArgs encodes a call to the commit Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) CommitWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut OffRampState",
		"&clock::Clock",
		"vector<vector<u8>>",
		"vector<u8>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("commit", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetMerkleRoot encodes a call to the get_merkle_root Move function.
func (c offrampEncoder) GetMerkleRoot(state bind.Object, root []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_merkle_root", typeArgsList, typeParamsList, []string{
		"&OffRampState",
		"vector<u8>",
	}, []any{
		state,
		root,
	}, []string{
		"u64",
	})
}

// GetMerkleRootWithArgs encodes a call to the get_merkle_root Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetMerkleRootWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_merkle_root", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetSourceChainConfig encodes a call to the get_source_chain_config Move function.
func (c offrampEncoder) GetSourceChainConfig(state bind.Object, sourceChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config", typeArgsList, typeParamsList, []string{
		"&OffRampState",
		"u64",
	}, []any{
		state,
		sourceChainSelector,
	}, []string{
		"ccip_offramp::offramp::SourceChainConfig",
	})
}

// GetSourceChainConfigWithArgs encodes a call to the get_source_chain_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetSourceChainConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_offramp::offramp::SourceChainConfig",
	})
}

// GetSourceChainConfigFields encodes a call to the get_source_chain_config_fields Move function.
func (c offrampEncoder) GetSourceChainConfigFields(sourceChainConfig SourceChainConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config_fields", typeArgsList, typeParamsList, []string{
		"ccip_offramp::offramp::SourceChainConfig",
	}, []any{
		sourceChainConfig,
	}, []string{
		"address",
		"bool",
		"u64",
		"bool",
		"vector<u8>",
	})
}

// GetSourceChainConfigFieldsWithArgs encodes a call to the get_source_chain_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetSourceChainConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_offramp::offramp::SourceChainConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_source_chain_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"bool",
		"u64",
		"bool",
		"vector<u8>",
	})
}

// GetAllSourceChainConfigs encodes a call to the get_all_source_chain_configs Move function.
func (c offrampEncoder) GetAllSourceChainConfigs(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_source_chain_configs", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"vector<u64>",
		"vector<ccip_offramp::offramp::SourceChainConfig>",
	})
}

// GetAllSourceChainConfigsWithArgs encodes a call to the get_all_source_chain_configs Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetAllSourceChainConfigsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_all_source_chain_configs", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u64>",
		"vector<ccip_offramp::offramp::SourceChainConfig>",
	})
}

// GetStaticConfig encodes a call to the get_static_config Move function.
func (c offrampEncoder) GetStaticConfig(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"ccip_offramp::offramp::StaticConfig",
	})
}

// GetStaticConfigWithArgs encodes a call to the get_static_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_offramp::offramp::StaticConfig",
	})
}

// GetStaticConfigFields encodes a call to the get_static_config_fields Move function.
func (c offrampEncoder) GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, []string{
		"ccip_offramp::offramp::StaticConfig",
	}, []any{
		cfg,
	}, []string{
		"u64",
		"address",
		"address",
		"address",
	})
}

// GetStaticConfigFieldsWithArgs encodes a call to the get_static_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_offramp::offramp::StaticConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
		"address",
		"address",
		"address",
	})
}

// GetDynamicConfig encodes a call to the get_dynamic_config Move function.
func (c offrampEncoder) GetDynamicConfig(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"ccip_offramp::offramp::DynamicConfig",
	})
}

// GetDynamicConfigWithArgs encodes a call to the get_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_offramp::offramp::DynamicConfig",
	})
}

// GetDynamicConfigFields encodes a call to the get_dynamic_config_fields Move function.
func (c offrampEncoder) GetDynamicConfigFields(cfg DynamicConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, []string{
		"ccip_offramp::offramp::DynamicConfig",
	}, []any{
		cfg,
	}, []string{
		"address",
		"u32",
	})
}

// GetDynamicConfigFieldsWithArgs encodes a call to the get_dynamic_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_offramp::offramp::DynamicConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"u32",
	})
}

// SetDynamicConfig encodes a call to the set_dynamic_config Move function.
func (c offrampEncoder) SetDynamicConfig(state bind.Object, param bind.Object, permissionlessExecutionThresholdSeconds uint32) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"u32",
	}, []any{
		state,
		param,
		permissionlessExecutionThresholdSeconds,
	}, nil)
}

// SetDynamicConfigWithArgs encodes a call to the set_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"u32",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplySourceChainConfigUpdates encodes a call to the apply_source_chain_config_updates Move function.
func (c offrampEncoder) ApplySourceChainConfigUpdates(state bind.Object, param bind.Object, sourceChainsSelector []uint64, sourceChainsIsEnabled []bool, sourceChainsIsRmnVerificationDisabled []bool, sourceChainsOnRamp [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_source_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}, []any{
		state,
		param,
		sourceChainsSelector,
		sourceChainsIsEnabled,
		sourceChainsIsRmnVerificationDisabled,
		sourceChainsOnRamp,
	}, nil)
}

// ApplySourceChainConfigUpdatesWithArgs encodes a call to the apply_source_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ApplySourceChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_source_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetCcipPackageId encodes a call to the get_ccip_package_id Move function.
func (c offrampEncoder) GetCcipPackageId() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ccip_package_id", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"address",
	})
}

// GetCcipPackageIdWithArgs encodes a call to the get_ccip_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) GetCcipPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ccip_package_id", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// Owner encodes a call to the owner Move function.
func (c offrampEncoder) Owner(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
	})
}

// HasPendingTransfer encodes a call to the has_pending_transfer Move function.
func (c offrampEncoder) HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// PendingTransferFrom encodes a call to the pending_transfer_from Move function.
func (c offrampEncoder) PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferTo encodes a call to the pending_transfer_to Move function.
func (c offrampEncoder) PendingTransferTo(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferAccepted encodes a call to the pending_transfer_accepted Move function.
func (c offrampEncoder) PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&OffRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, expectedParams, args, []string{
		"0x1::option::Option<bool>",
	})
}

// TransferOwnership encodes a call to the transfer_ownership Move function.
func (c offrampEncoder) TransferOwnership(state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}, []any{
		state,
		ownerCap,
		newOwner,
	}, nil)
}

// TransferOwnershipWithArgs encodes a call to the transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&OwnerCap",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnership encodes a call to the accept_ownership Move function.
func (c offrampEncoder) AcceptOwnership(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c offrampEncoder) AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipAsMcms encodes a call to the accept_ownership_as_mcms Move function.
func (c offrampEncoder) AcceptOwnershipAsMcms(state bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_mcms", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"ExecutingCallbackParams",
	}, []any{
		state,
		params,
	}, nil)
}

// AcceptOwnershipAsMcmsWithArgs encodes a call to the accept_ownership_as_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) AcceptOwnershipAsMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_as_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c offrampEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut OwnableState",
		"address",
	}, []any{
		ownerCap,
		ownableState,
		to,
	}, nil)
}

// ExecuteOwnershipTransferWithArgs encodes a call to the execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut OwnableState",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransferToMcms encodes a call to the execute_ownership_transfer_to_mcms Move function.
func (c offrampEncoder) ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut OffRampState",
		"&mut Registry",
		"address",
	}, []any{
		ownerCap,
		state,
		registry,
		to,
	}, nil)
}

// ExecuteOwnershipTransferToMcmsWithArgs encodes a call to the execute_ownership_transfer_to_mcms Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut OffRampState",
		"&mut Registry",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsRegisterUpgradeCap encodes a call to the mcms_register_upgrade_cap Move function.
func (c offrampEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, []string{
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}, []any{
		upgradeCap,
		registry,
		state,
	}, nil)
}

// McmsRegisterUpgradeCapWithArgs encodes a call to the mcms_register_upgrade_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"UpgradeCap",
		"&mut Registry",
		"&mut DeployerState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_register_upgrade_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsEntrypoint encodes a call to the mcms_entrypoint Move function.
func (c offrampEncoder) McmsEntrypoint(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsEntrypointWithArgs encodes a call to the mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c offrampEncoder) McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OffRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, expectedParams, args, nil)
}
