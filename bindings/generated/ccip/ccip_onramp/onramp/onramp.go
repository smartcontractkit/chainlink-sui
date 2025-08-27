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
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*models.SuiTransactionBlockResponse, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	WithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, feeTokenMetadata bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error)
	SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, feeAggregator string, allowlistAdmin string) (*models.SuiTransactionBlockResponse, error)
	ApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*models.SuiTransactionBlockResponse, error)
	GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error)
	ApplyAllowlistUpdatesByAdmin(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error)
	GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*models.SuiTransactionBlockResponse, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*models.SuiTransactionBlockResponse, error)
	CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error)
	GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error)
	McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error)
	McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsSetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsInitialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	McmsWithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IOnrampDevInspect
	Encoder() OnrampEncoder
}

type IOnrampDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (bool, error)
	GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (uint64, error)
	GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (uint64, error)
	GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error)
	GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error)
	GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (StaticConfig, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (uint64, error)
	GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (DynamicConfig, error)
	GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) ([]any, error)
	CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) ([]byte, error)
	CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) ([]byte, error)
	CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) ([]byte, error)
	GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error)
	Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error)
	HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error)
	PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error)
	PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error)
}

type OnrampEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(state bind.Object, param bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	IsChainSupported(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	IsChainSupportedWithArgs(args ...any) (*bind.EncodedCall, error)
	GetExpectedNextSequenceNumber(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetExpectedNextSequenceNumberWithArgs(args ...any) (*bind.EncodedCall, error)
	WithdrawFeeTokens(typeArgs []string, state bind.Object, param bind.Object, feeTokenMetadata bind.Object) (*bind.EncodedCall, error)
	WithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	GetFee(typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error)
	GetFeeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
	SetDynamicConfig(state bind.Object, param bind.Object, feeAggregator string, allowlistAdmin string) (*bind.EncodedCall, error)
	SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyDestChainConfigUpdates(state bind.Object, param bind.Object, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*bind.EncodedCall, error)
	ApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestChainConfig(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetDestChainConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetAllowedSendersList(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetAllowedSendersListWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyAllowlistUpdates(state bind.Object, param bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesByAdmin(state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error)
	ApplyAllowlistUpdatesByAdminWithArgs(args ...any) (*bind.EncodedCall, error)
	GetOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error)
	GetOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfig(state bind.Object) (*bind.EncodedCall, error)
	GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error)
	GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfig(state bind.Object) (*bind.EncodedCall, error)
	GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDynamicConfigFields(cfg DynamicConfig) (*bind.EncodedCall, error)
	GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMessageHash(onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*bind.EncodedCall, error)
	CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error)
	CalculateMetadataHash(sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*bind.EncodedCall, error)
	CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error)
	CcipSend(typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error)
	CcipSendWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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
	McmsAcceptOwnership(state bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error)
	ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error)
	McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsSetDynamicConfig(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsSetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsApplyDestChainConfigUpdates(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsApplyAllowlistUpdates(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsTransferOwnership(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransfer(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsInitialize(state bind.Object, registry bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsInitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsWithdrawFeeTokens(typeArgs []string, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsWithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error)
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

func NewOnramp(packageID string, client sui.ISuiAPI) (*OnrampContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip_onramp", "onramp", client)
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

func (c *OnrampContract) Encoder() OnrampEncoder {
	return c.onrampEncoder
}

func (c *OnrampContract) DevInspect() IOnrampDevInspect {
	return c.devInspect
}

type OnRampState struct {
	Id                string       `move:"sui::object::UID"`
	ChainSelector     uint64       `move:"u64"`
	FeeAggregator     string       `move:"address"`
	AllowlistAdmin    string       `move:"address"`
	DestChainConfigs  bind.Object  `move:"Table<u64, DestChainConfig>"`
	FeeTokens         bind.Object  `move:"Bag"`
	NonceManagerCap   *bind.Object `move:"0x1::option::Option<NonceManagerCap>"`
	SourceTransferCap *bind.Object `move:"0x1::option::Option<osh::SourceTransferCap>"`
	OwnableState      bind.Object  `move:"OwnableState"`
}

type OnRampStatePointer struct {
	Id            string `move:"sui::object::UID"`
	OnRampStateId string `move:"address"`
	OwnerCapId    string `move:"address"`
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

type StaticConfig struct {
	ChainSelector uint64 `move:"u64"`
}

type DynamicConfig struct {
	FeeAggregator  string `move:"address"`
	AllowlistAdmin string `move:"address"`
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

type ONRAMP struct {
}

type McmsCallback struct {
}

type bcsOnRampState struct {
	Id                string
	ChainSelector     uint64
	FeeAggregator     [32]byte
	AllowlistAdmin    [32]byte
	DestChainConfigs  bind.Object
	FeeTokens         bind.Object
	NonceManagerCap   *bind.Object
	SourceTransferCap *bind.Object
	OwnableState      bind.Object
}

func convertOnRampStateFromBCS(bcs bcsOnRampState) (OnRampState, error) {

	return OnRampState{
		Id:                bcs.Id,
		ChainSelector:     bcs.ChainSelector,
		FeeAggregator:     fmt.Sprintf("0x%x", bcs.FeeAggregator),
		AllowlistAdmin:    fmt.Sprintf("0x%x", bcs.AllowlistAdmin),
		DestChainConfigs:  bcs.DestChainConfigs,
		FeeTokens:         bcs.FeeTokens,
		NonceManagerCap:   bcs.NonceManagerCap,
		SourceTransferCap: bcs.SourceTransferCap,
		OwnableState:      bcs.OwnableState,
	}, nil
}

type bcsOnRampStatePointer struct {
	Id            string
	OnRampStateId [32]byte
	OwnerCapId    [32]byte
}

func convertOnRampStatePointerFromBCS(bcs bcsOnRampStatePointer) (OnRampStatePointer, error) {

	return OnRampStatePointer{
		Id:            bcs.Id,
		OnRampStateId: fmt.Sprintf("0x%x", bcs.OnRampStateId),
		OwnerCapId:    fmt.Sprintf("0x%x", bcs.OwnerCapId),
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

func init() {
	bind.RegisterStructDecoder("ccip_onramp::onramp::OnRampState", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::OnRampStatePointer", func(data []byte) (interface{}, error) {
		var temp bcsOnRampStatePointer
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertOnRampStatePointerFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::DestChainConfig", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::RampMessageHeader", func(data []byte) (interface{}, error) {
		var result RampMessageHeader
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::Sui2AnyRampMessage", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::Sui2AnyTokenTransfer", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::StaticConfig", func(data []byte) (interface{}, error) {
		var result StaticConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::DynamicConfig", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::ConfigSet", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::DestChainConfigSet", func(data []byte) (interface{}, error) {
		var result DestChainConfigSet
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::CCIPMessageSent", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::AllowlistSendersAdded", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::AllowlistSendersRemoved", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::FeeTokenWithdrawn", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip_onramp::onramp::ONRAMP", func(data []byte) (interface{}, error) {
		var result ONRAMP
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip_onramp::onramp::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *OnrampContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *OnrampContract) Initialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.Initialize(state, param, nonceManagerCap, sourceTransferCap, chainSelector, feeAggregator, allowlistAdmin, destChainSelectors, destChainEnabled, destChainAllowlistEnabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IsChainSupported executes the is_chain_supported Move function.
func (c *OnrampContract) IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.IsChainSupported(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetExpectedNextSequenceNumber executes the get_expected_next_sequence_number Move function.
func (c *OnrampContract) GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetExpectedNextSequenceNumber(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// WithdrawFeeTokens executes the withdraw_fee_tokens Move function.
func (c *OnrampContract) WithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, param bind.Object, feeTokenMetadata bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.WithdrawFeeTokens(typeArgs, state, param, feeTokenMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFee executes the get_fee Move function.
func (c *OnrampContract) GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetFee(typeArgs, ref, clock, destChainSelector, receiver, data, tokenAddresses, tokenAmounts, feeToken, extraArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// SetDynamicConfig executes the set_dynamic_config Move function.
func (c *OnrampContract) SetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, feeAggregator string, allowlistAdmin string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.SetDynamicConfig(state, param, feeAggregator, allowlistAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyDestChainConfigUpdates executes the apply_dest_chain_config_updates Move function.
func (c *OnrampContract) ApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ApplyDestChainConfigUpdates(state, param, destChainSelectors, destChainEnabled, destChainAllowlistEnabled)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestChainConfig executes the get_dest_chain_config Move function.
func (c *OnrampContract) GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetDestChainConfig(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetAllowedSendersList executes the get_allowed_senders_list Move function.
func (c *OnrampContract) GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetAllowedSendersList(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdates executes the apply_allowlist_updates Move function.
func (c *OnrampContract) ApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, param bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ApplyAllowlistUpdates(state, param, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyAllowlistUpdatesByAdmin executes the apply_allowlist_updates_by_admin Move function.
func (c *OnrampContract) ApplyAllowlistUpdatesByAdmin(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ApplyAllowlistUpdatesByAdmin(state, destChainSelectors, destChainAllowlistEnabled, destChainAddAllowedSenders, destChainRemoveAllowedSenders)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetOutboundNonce executes the get_outbound_nonce Move function.
func (c *OnrampContract) GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetOutboundNonce(ref, destChainSelector, sender)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function.
func (c *OnrampContract) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetStaticConfig(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfigFields executes the get_static_config_fields Move function.
func (c *OnrampContract) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetStaticConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfig executes the get_dynamic_config Move function.
func (c *OnrampContract) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetDynamicConfig(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDynamicConfigFields executes the get_dynamic_config_fields Move function.
func (c *OnrampContract) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetDynamicConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMessageHash executes the calculate_message_hash Move function.
func (c *OnrampContract) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.CalculateMessageHash(onRampAddress, messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, data, feeToken, feeTokenAmount, sourcePoolAddresses, destTokenAddresses, extraDatas, amounts, destExecDatas, extraArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CalculateMetadataHash executes the calculate_metadata_hash Move function.
func (c *OnrampContract) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.CalculateMetadataHash(sourceChainSelector, destChainSelector, onRampAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// CcipSend executes the ccip_send Move function.
func (c *OnrampContract) CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.CcipSend(typeArgs, ref, state, clock, destChainSelector, receiver, data, tokenParams, feeTokenMetadata, feeToken, extraArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetCcipPackageId executes the get_ccip_package_id Move function.
func (c *OnrampContract) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.GetCcipPackageId()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Owner executes the owner Move function.
func (c *OnrampContract) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.Owner(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// HasPendingTransfer executes the has_pending_transfer Move function.
func (c *OnrampContract) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.HasPendingTransfer(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferFrom executes the pending_transfer_from Move function.
func (c *OnrampContract) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.PendingTransferFrom(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferTo executes the pending_transfer_to Move function.
func (c *OnrampContract) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.PendingTransferTo(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// PendingTransferAccepted executes the pending_transfer_accepted Move function.
func (c *OnrampContract) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.PendingTransferAccepted(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TransferOwnership executes the transfer_ownership Move function.
func (c *OnrampContract) TransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, ownerCap bind.Object, newOwner string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.TransferOwnership(state, ownerCap, newOwner)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnership executes the accept_ownership Move function.
func (c *OnrampContract) AcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.AcceptOwnership(state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// AcceptOwnershipFromObject executes the accept_ownership_from_object Move function.
func (c *OnrampContract) AcceptOwnershipFromObject(ctx context.Context, opts *bind.CallOpts, state bind.Object, from string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.AcceptOwnershipFromObject(state, from)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsAcceptOwnership executes the mcms_accept_ownership Move function.
func (c *OnrampContract) McmsAcceptOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsAcceptOwnership(state, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransfer executes the execute_ownership_transfer Move function.
func (c *OnrampContract) ExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, ownableState bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ExecuteOwnershipTransfer(ownerCap, ownableState, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ExecuteOwnershipTransferToMcms executes the execute_ownership_transfer_to_mcms Move function.
func (c *OnrampContract) ExecuteOwnershipTransferToMcms(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.ExecuteOwnershipTransferToMcms(ownerCap, state, registry, to)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsRegisterUpgradeCap executes the mcms_register_upgrade_cap Move function.
func (c *OnrampContract) McmsRegisterUpgradeCap(ctx context.Context, opts *bind.CallOpts, upgradeCap bind.Object, registry bind.Object, state bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsRegisterUpgradeCap(upgradeCap, registry, state)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsSetDynamicConfig executes the mcms_set_dynamic_config Move function.
func (c *OnrampContract) McmsSetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsSetDynamicConfig(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyDestChainConfigUpdates executes the mcms_apply_dest_chain_config_updates Move function.
func (c *OnrampContract) McmsApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsApplyDestChainConfigUpdates(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsApplyAllowlistUpdates executes the mcms_apply_allowlist_updates Move function.
func (c *OnrampContract) McmsApplyAllowlistUpdates(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsApplyAllowlistUpdates(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsTransferOwnership executes the mcms_transfer_ownership Move function.
func (c *OnrampContract) McmsTransferOwnership(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsTransferOwnership(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsExecuteOwnershipTransfer executes the mcms_execute_ownership_transfer Move function.
func (c *OnrampContract) McmsExecuteOwnershipTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsExecuteOwnershipTransfer(state, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsInitialize executes the mcms_initialize Move function.
func (c *OnrampContract) McmsInitialize(ctx context.Context, opts *bind.CallOpts, state bind.Object, registry bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsInitialize(state, registry, nonceManagerCap, sourceTransferCap, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsWithdrawFeeTokens executes the mcms_withdraw_fee_tokens Move function.
func (c *OnrampContract) McmsWithdrawFeeTokens(ctx context.Context, opts *bind.CallOpts, typeArgs []string, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.onrampEncoder.McmsWithdrawFeeTokens(typeArgs, state, registry, feeTokenMetadata, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *OnrampDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.onrampEncoder.TypeAndVersion()
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

// IsChainSupported executes the is_chain_supported Move function using DevInspect to get return values.
//
// Returns: bool
func (d *OnrampDevInspect) IsChainSupported(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (bool, error) {
	encoded, err := d.contract.onrampEncoder.IsChainSupported(state, destChainSelector)
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

// GetExpectedNextSequenceNumber executes the get_expected_next_sequence_number Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OnrampDevInspect) GetExpectedNextSequenceNumber(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetExpectedNextSequenceNumber(state, destChainSelector)
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

// GetFee executes the get_fee Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OnrampDevInspect) GetFee(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetFee(typeArgs, ref, clock, destChainSelector, receiver, data, tokenAddresses, tokenAmounts, feeToken, extraArgs)
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

// GetDestChainConfig executes the get_dest_chain_config Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: bool
//	[1]: u64
//	[2]: bool
//	[3]: vector<address>
func (d *OnrampDevInspect) GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error) {
	encoded, err := d.contract.onrampEncoder.GetDestChainConfig(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetAllowedSendersList executes the get_allowed_senders_list Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: bool
//	[1]: vector<address>
func (d *OnrampDevInspect) GetAllowedSendersList(ctx context.Context, opts *bind.CallOpts, state bind.Object, destChainSelector uint64) ([]any, error) {
	encoded, err := d.contract.onrampEncoder.GetAllowedSendersList(state, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetOutboundNonce executes the get_outbound_nonce Move function using DevInspect to get return values.
//
// Returns: u64
func (d *OnrampDevInspect) GetOutboundNonce(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, sender string) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetOutboundNonce(ref, destChainSelector, sender)
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

// GetStaticConfig executes the get_static_config Move function using DevInspect to get return values.
//
// Returns: StaticConfig
func (d *OnrampDevInspect) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (StaticConfig, error) {
	encoded, err := d.contract.onrampEncoder.GetStaticConfig(state)
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
// Returns: u64
func (d *OnrampDevInspect) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (uint64, error) {
	encoded, err := d.contract.onrampEncoder.GetStaticConfigFields(cfg)
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

// GetDynamicConfig executes the get_dynamic_config Move function using DevInspect to get return values.
//
// Returns: DynamicConfig
func (d *OnrampDevInspect) GetDynamicConfig(ctx context.Context, opts *bind.CallOpts, state bind.Object) (DynamicConfig, error) {
	encoded, err := d.contract.onrampEncoder.GetDynamicConfig(state)
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
//	[1]: address
func (d *OnrampDevInspect) GetDynamicConfigFields(ctx context.Context, opts *bind.CallOpts, cfg DynamicConfig) ([]any, error) {
	encoded, err := d.contract.onrampEncoder.GetDynamicConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// CalculateMessageHash executes the calculate_message_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OnrampDevInspect) CalculateMessageHash(ctx context.Context, opts *bind.CallOpts, onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) ([]byte, error) {
	encoded, err := d.contract.onrampEncoder.CalculateMessageHash(onRampAddress, messageId, sourceChainSelector, destChainSelector, sequenceNumber, nonce, sender, receiver, data, feeToken, feeTokenAmount, sourcePoolAddresses, destTokenAddresses, extraDatas, amounts, destExecDatas, extraArgs)
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

// CalculateMetadataHash executes the calculate_metadata_hash Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OnrampDevInspect) CalculateMetadataHash(ctx context.Context, opts *bind.CallOpts, sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) ([]byte, error) {
	encoded, err := d.contract.onrampEncoder.CalculateMetadataHash(sourceChainSelector, destChainSelector, onRampAddress)
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

// CcipSend executes the ccip_send Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *OnrampDevInspect) CcipSend(ctx context.Context, opts *bind.CallOpts, typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) ([]byte, error) {
	encoded, err := d.contract.onrampEncoder.CcipSend(typeArgs, ref, state, clock, destChainSelector, receiver, data, tokenParams, feeTokenMetadata, feeToken, extraArgs)
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

// GetCcipPackageId executes the get_ccip_package_id Move function using DevInspect to get return values.
//
// Returns: address
func (d *OnrampDevInspect) GetCcipPackageId(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.onrampEncoder.GetCcipPackageId()
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
func (d *OnrampDevInspect) Owner(ctx context.Context, opts *bind.CallOpts, state bind.Object) (string, error) {
	encoded, err := d.contract.onrampEncoder.Owner(state)
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
func (d *OnrampDevInspect) HasPendingTransfer(ctx context.Context, opts *bind.CallOpts, state bind.Object) (bool, error) {
	encoded, err := d.contract.onrampEncoder.HasPendingTransfer(state)
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
func (d *OnrampDevInspect) PendingTransferFrom(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.onrampEncoder.PendingTransferFrom(state)
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
func (d *OnrampDevInspect) PendingTransferTo(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*string, error) {
	encoded, err := d.contract.onrampEncoder.PendingTransferTo(state)
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
func (d *OnrampDevInspect) PendingTransferAccepted(ctx context.Context, opts *bind.CallOpts, state bind.Object) (*bool, error) {
	encoded, err := d.contract.onrampEncoder.PendingTransferAccepted(state)
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

type onrampEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c onrampEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) Initialize(state bind.Object, param bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, chainSelector uint64, feeAggregator string, allowlistAdmin string, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&OwnerCap",
		"NonceManagerCap",
		"osh::SourceTransferCap",
		"u64",
		"address",
		"address",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
	}, []any{
		state,
		param,
		nonceManagerCap,
		sourceTransferCap,
		chainSelector,
		feeAggregator,
		allowlistAdmin,
		destChainSelectors,
		destChainEnabled,
		destChainAllowlistEnabled,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&OwnerCap",
		"NonceManagerCap",
		"osh::SourceTransferCap",
		"u64",
		"address",
		"address",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IsChainSupported encodes a call to the is_chain_supported Move function.
func (c onrampEncoder) IsChainSupported(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_chain_supported", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"bool",
	})
}

// IsChainSupportedWithArgs encodes a call to the is_chain_supported Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) IsChainSupportedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("is_chain_supported", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
	})
}

// GetExpectedNextSequenceNumber encodes a call to the get_expected_next_sequence_number Move function.
func (c onrampEncoder) GetExpectedNextSequenceNumber(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_expected_next_sequence_number", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"u64",
	})
}

// GetExpectedNextSequenceNumberWithArgs encodes a call to the get_expected_next_sequence_number Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetExpectedNextSequenceNumberWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_expected_next_sequence_number", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// WithdrawFeeTokens encodes a call to the withdraw_fee_tokens Move function.
func (c onrampEncoder) WithdrawFeeTokens(typeArgs []string, state bind.Object, param bind.Object, feeTokenMetadata bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("withdraw_fee_tokens", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&OwnerCap",
		"&CoinMetadata<T>",
	}, []any{
		state,
		param,
		feeTokenMetadata,
	}, nil)
}

// WithdrawFeeTokensWithArgs encodes a call to the withdraw_fee_tokens Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) WithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&OwnerCap",
		"&CoinMetadata<T>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("withdraw_fee_tokens", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetFee encodes a call to the get_fee Move function.
func (c onrampEncoder) GetFee(typeArgs []string, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenAddresses []string, tokenAmounts []uint64, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_fee", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"vector<address>",
		"vector<u64>",
		"&CoinMetadata<T>",
		"vector<u8>",
	}, []any{
		ref,
		clock,
		destChainSelector,
		receiver,
		data,
		tokenAddresses,
		tokenAmounts,
		feeToken,
		extraArgs,
	}, []string{
		"u64",
	})
}

// GetFeeWithArgs encodes a call to the get_fee Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetFeeWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"vector<address>",
		"vector<u64>",
		"&CoinMetadata<T>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("get_fee", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// SetDynamicConfig encodes a call to the set_dynamic_config Move function.
func (c onrampEncoder) SetDynamicConfig(state bind.Object, param bind.Object, feeAggregator string, allowlistAdmin string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&OwnerCap",
		"address",
		"address",
	}, []any{
		state,
		param,
		feeAggregator,
		allowlistAdmin,
	}, nil)
}

// SetDynamicConfigWithArgs encodes a call to the set_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) SetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&OwnerCap",
		"address",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("set_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplyDestChainConfigUpdates encodes a call to the apply_dest_chain_config_updates Move function.
func (c onrampEncoder) ApplyDestChainConfigUpdates(state bind.Object, param bind.Object, destChainSelectors []uint64, destChainEnabled []bool, destChainAllowlistEnabled []bool) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_dest_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
	}, []any{
		state,
		param,
		destChainSelectors,
		destChainEnabled,
		destChainAllowlistEnabled,
	}, nil)
}

// ApplyDestChainConfigUpdatesWithArgs encodes a call to the apply_dest_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) ApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<bool>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_dest_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetDestChainConfig encodes a call to the get_dest_chain_config Move function.
func (c onrampEncoder) GetDestChainConfig(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"bool",
		"u64",
		"bool",
		"vector<address>",
	})
}

// GetDestChainConfigWithArgs encodes a call to the get_dest_chain_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetDestChainConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
		"u64",
		"bool",
		"vector<address>",
	})
}

// GetAllowedSendersList encodes a call to the get_allowed_senders_list Move function.
func (c onrampEncoder) GetAllowedSendersList(state bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_senders_list", typeArgsList, typeParamsList, []string{
		"&OnRampState",
		"u64",
	}, []any{
		state,
		destChainSelector,
	}, []string{
		"bool",
		"vector<address>",
	})
}

// GetAllowedSendersListWithArgs encodes a call to the get_allowed_senders_list Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetAllowedSendersListWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_allowed_senders_list", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
		"vector<address>",
	})
}

// ApplyAllowlistUpdates encodes a call to the apply_allowlist_updates Move function.
func (c onrampEncoder) ApplyAllowlistUpdates(state bind.Object, param bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}, []any{
		state,
		param,
		destChainSelectors,
		destChainAllowlistEnabled,
		destChainAddAllowedSenders,
		destChainRemoveAllowedSenders,
	}, nil)
}

// ApplyAllowlistUpdatesWithArgs encodes a call to the apply_allowlist_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) ApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&OwnerCap",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ApplyAllowlistUpdatesByAdmin encodes a call to the apply_allowlist_updates_by_admin Move function.
func (c onrampEncoder) ApplyAllowlistUpdatesByAdmin(state bind.Object, destChainSelectors []uint64, destChainAllowlistEnabled []bool, destChainAddAllowedSenders [][]string, destChainRemoveAllowedSenders [][]string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates_by_admin", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}, []any{
		state,
		destChainSelectors,
		destChainAllowlistEnabled,
		destChainAddAllowedSenders,
		destChainRemoveAllowedSenders,
	}, nil)
}

// ApplyAllowlistUpdatesByAdminWithArgs encodes a call to the apply_allowlist_updates_by_admin Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) ApplyAllowlistUpdatesByAdminWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"vector<u64>",
		"vector<bool>",
		"vector<vector<address>>",
		"vector<vector<address>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_allowlist_updates_by_admin", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetOutboundNonce encodes a call to the get_outbound_nonce Move function.
func (c onrampEncoder) GetOutboundNonce(ref bind.Object, destChainSelector uint64, sender string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_outbound_nonce", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"address",
	}, []any{
		ref,
		destChainSelector,
		sender,
	}, []string{
		"u64",
	})
}

// GetOutboundNonceWithArgs encodes a call to the get_outbound_nonce Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetOutboundNonceWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_outbound_nonce", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetStaticConfig encodes a call to the get_static_config Move function.
func (c onrampEncoder) GetStaticConfig(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"ccip_onramp::onramp::StaticConfig",
	})
}

// GetStaticConfigWithArgs encodes a call to the get_static_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_onramp::onramp::StaticConfig",
	})
}

// GetStaticConfigFields encodes a call to the get_static_config_fields Move function.
func (c onrampEncoder) GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, []string{
		"ccip_onramp::onramp::StaticConfig",
	}, []any{
		cfg,
	}, []string{
		"u64",
	})
}

// GetStaticConfigFieldsWithArgs encodes a call to the get_static_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_onramp::onramp::StaticConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetDynamicConfig encodes a call to the get_dynamic_config Move function.
func (c onrampEncoder) GetDynamicConfig(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"ccip_onramp::onramp::DynamicConfig",
	})
}

// GetDynamicConfigWithArgs encodes a call to the get_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip_onramp::onramp::DynamicConfig",
	})
}

// GetDynamicConfigFields encodes a call to the get_dynamic_config_fields Move function.
func (c onrampEncoder) GetDynamicConfigFields(cfg DynamicConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, []string{
		"ccip_onramp::onramp::DynamicConfig",
	}, []any{
		cfg,
	}, []string{
		"address",
		"address",
	})
}

// GetDynamicConfigFieldsWithArgs encodes a call to the get_dynamic_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetDynamicConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip_onramp::onramp::DynamicConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dynamic_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"address",
		"address",
	})
}

// CalculateMessageHash encodes a call to the calculate_message_hash Move function.
func (c onrampEncoder) CalculateMessageHash(onRampAddress string, messageId []byte, sourceChainSelector uint64, destChainSelector uint64, sequenceNumber uint64, nonce uint64, sender string, receiver []byte, data []byte, feeToken string, feeTokenAmount uint64, sourcePoolAddresses []string, destTokenAddresses [][]byte, extraDatas [][]byte, amounts []uint64, destExecDatas [][]byte, extraArgs []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_message_hash", typeArgsList, typeParamsList, []string{
		"address",
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"address",
		"vector<u8>",
		"vector<u8>",
		"address",
		"u64",
		"vector<address>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<u8>",
	}, []any{
		onRampAddress,
		messageId,
		sourceChainSelector,
		destChainSelector,
		sequenceNumber,
		nonce,
		sender,
		receiver,
		data,
		feeToken,
		feeTokenAmount,
		sourcePoolAddresses,
		destTokenAddresses,
		extraDatas,
		amounts,
		destExecDatas,
		extraArgs,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMessageHashWithArgs encodes a call to the calculate_message_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) CalculateMessageHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"address",
		"vector<u8>",
		"u64",
		"u64",
		"u64",
		"u64",
		"address",
		"vector<u8>",
		"vector<u8>",
		"address",
		"u64",
		"vector<address>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
		"vector<u64>",
		"vector<vector<u8>>",
		"vector<u8>",
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

// CalculateMetadataHash encodes a call to the calculate_metadata_hash Move function.
func (c onrampEncoder) CalculateMetadataHash(sourceChainSelector uint64, destChainSelector uint64, onRampAddress string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("calculate_metadata_hash", typeArgsList, typeParamsList, []string{
		"u64",
		"u64",
		"address",
	}, []any{
		sourceChainSelector,
		destChainSelector,
		onRampAddress,
	}, []string{
		"vector<u8>",
	})
}

// CalculateMetadataHashWithArgs encodes a call to the calculate_metadata_hash Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) CalculateMetadataHashWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"u64",
		"u64",
		"address",
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

// CcipSend encodes a call to the ccip_send Move function.
func (c onrampEncoder) CcipSend(typeArgs []string, ref bind.Object, state bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, tokenParams bind.Object, feeTokenMetadata bind.Object, feeToken bind.Object, extraArgs []byte) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("ccip_send", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut OnRampState",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"TokenTransferParams",
		"&CoinMetadata<T>",
		"&mut Coin<T>",
		"vector<u8>",
	}, []any{
		ref,
		state,
		clock,
		destChainSelector,
		receiver,
		data,
		tokenParams,
		feeTokenMetadata,
		feeToken,
		extraArgs,
	}, []string{
		"vector<u8>",
	})
}

// CcipSendWithArgs encodes a call to the ccip_send Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) CcipSendWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&mut OnRampState",
		"&Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"TokenTransferParams",
		"&CoinMetadata<T>",
		"&mut Coin<T>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("ccip_send", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// GetCcipPackageId encodes a call to the get_ccip_package_id Move function.
func (c onrampEncoder) GetCcipPackageId() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_ccip_package_id", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"address",
	})
}

// GetCcipPackageIdWithArgs encodes a call to the get_ccip_package_id Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) GetCcipPackageIdWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) Owner(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("owner", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"address",
	})
}

// OwnerWithArgs encodes a call to the owner Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) OwnerWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) HasPendingTransfer(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("has_pending_transfer", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"bool",
	})
}

// HasPendingTransferWithArgs encodes a call to the has_pending_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) HasPendingTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) PendingTransferFrom(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_from", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferFromWithArgs encodes a call to the pending_transfer_from Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) PendingTransferFromWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) PendingTransferTo(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_to", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<address>",
	})
}

// PendingTransferToWithArgs encodes a call to the pending_transfer_to Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) PendingTransferToWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) PendingTransferAccepted(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("pending_transfer_accepted", typeArgsList, typeParamsList, []string{
		"&OnRampState",
	}, []any{
		state,
	}, []string{
		"0x1::option::Option<bool>",
	})
}

// PendingTransferAcceptedWithArgs encodes a call to the pending_transfer_accepted Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) PendingTransferAcceptedWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OnRampState",
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
func (c onrampEncoder) TransferOwnership(state bind.Object, ownerCap bind.Object, newOwner string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
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
func (c onrampEncoder) TransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
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
func (c onrampEncoder) AcceptOwnership(state bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
	}, []any{
		state,
	}, nil)
}

// AcceptOwnershipWithArgs encodes a call to the accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) AcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// AcceptOwnershipFromObject encodes a call to the accept_ownership_from_object Move function.
func (c onrampEncoder) AcceptOwnershipFromObject(state bind.Object, from string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut UID",
	}, []any{
		state,
		from,
	}, nil)
}

// AcceptOwnershipFromObjectWithArgs encodes a call to the accept_ownership_from_object Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) AcceptOwnershipFromObjectWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut UID",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("accept_ownership_from_object", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsAcceptOwnership encodes a call to the mcms_accept_ownership Move function.
func (c onrampEncoder) McmsAcceptOwnership(state bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"ExecutingCallbackParams",
	}, []any{
		state,
		params,
	}, nil)
}

// McmsAcceptOwnershipWithArgs encodes a call to the mcms_accept_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsAcceptOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_accept_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// ExecuteOwnershipTransfer encodes a call to the execute_ownership_transfer Move function.
func (c onrampEncoder) ExecuteOwnershipTransfer(ownerCap bind.Object, ownableState bind.Object, to string) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) ExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) ExecuteOwnershipTransferToMcms(ownerCap bind.Object, state bind.Object, registry bind.Object, to string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("execute_ownership_transfer_to_mcms", typeArgsList, typeParamsList, []string{
		"OwnerCap",
		"&mut OnRampState",
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
func (c onrampEncoder) ExecuteOwnershipTransferToMcmsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"OwnerCap",
		"&mut OnRampState",
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
func (c onrampEncoder) McmsRegisterUpgradeCap(upgradeCap bind.Object, registry bind.Object, state bind.Object) (*bind.EncodedCall, error) {
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
func (c onrampEncoder) McmsRegisterUpgradeCapWithArgs(args ...any) (*bind.EncodedCall, error) {
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

// McmsSetDynamicConfig encodes a call to the mcms_set_dynamic_config Move function.
func (c onrampEncoder) McmsSetDynamicConfig(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_dynamic_config", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsSetDynamicConfigWithArgs encodes a call to the mcms_set_dynamic_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsSetDynamicConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_set_dynamic_config", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsApplyDestChainConfigUpdates encodes a call to the mcms_apply_dest_chain_config_updates Move function.
func (c onrampEncoder) McmsApplyDestChainConfigUpdates(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_dest_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsApplyDestChainConfigUpdatesWithArgs encodes a call to the mcms_apply_dest_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_dest_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsApplyAllowlistUpdates encodes a call to the mcms_apply_allowlist_updates Move function.
func (c onrampEncoder) McmsApplyAllowlistUpdates(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_allowlist_updates", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsApplyAllowlistUpdatesWithArgs encodes a call to the mcms_apply_allowlist_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsApplyAllowlistUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_apply_allowlist_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsTransferOwnership encodes a call to the mcms_transfer_ownership Move function.
func (c onrampEncoder) McmsTransferOwnership(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsTransferOwnershipWithArgs encodes a call to the mcms_transfer_ownership Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsTransferOwnershipWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_transfer_ownership", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsExecuteOwnershipTransfer encodes a call to the mcms_execute_ownership_transfer Move function.
func (c onrampEncoder) McmsExecuteOwnershipTransfer(state bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		params,
	}, nil)
}

// McmsExecuteOwnershipTransferWithArgs encodes a call to the mcms_execute_ownership_transfer Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsExecuteOwnershipTransferWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut Registry",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_execute_ownership_transfer", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsInitialize encodes a call to the mcms_initialize Move function.
func (c onrampEncoder) McmsInitialize(state bind.Object, registry bind.Object, nonceManagerCap bind.Object, sourceTransferCap bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_initialize", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut Registry",
		"NonceManagerCap",
		"osh::SourceTransferCap",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		nonceManagerCap,
		sourceTransferCap,
		params,
	}, nil)
}

// McmsInitializeWithArgs encodes a call to the mcms_initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsInitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut Registry",
		"NonceManagerCap",
		"osh::SourceTransferCap",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// McmsWithdrawFeeTokens encodes a call to the mcms_withdraw_fee_tokens Move function.
func (c onrampEncoder) McmsWithdrawFeeTokens(typeArgs []string, state bind.Object, registry bind.Object, feeTokenMetadata bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_withdraw_fee_tokens", typeArgsList, typeParamsList, []string{
		"&mut OnRampState",
		"&mut Registry",
		"&CoinMetadata<T>",
		"ExecutingCallbackParams",
	}, []any{
		state,
		registry,
		feeTokenMetadata,
		params,
	}, nil)
}

// McmsWithdrawFeeTokensWithArgs encodes a call to the mcms_withdraw_fee_tokens Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c onrampEncoder) McmsWithdrawFeeTokensWithArgs(typeArgs []string, args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut OnRampState",
		"&mut Registry",
		"&CoinMetadata<T>",
		"ExecutingCallbackParams",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := typeArgs
	typeParamsList := []string{
		"T",
	}
	return c.EncodeCallArgsWithGenerics("mcms_withdraw_fee_tokens", typeArgsList, typeParamsList, expectedParams, args, nil)
}
