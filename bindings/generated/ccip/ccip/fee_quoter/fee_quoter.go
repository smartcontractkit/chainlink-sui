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
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error)
	Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, maxFeeJuelsPerMsg *big.Int, linkToken string, tokenPriceStalenessThreshold uint64, feeTokens []string) (*models.SuiTransactionBlockResponse, error)
	IssueFeeQuoterCap(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error)
	NewFeeQuoterCap(ctx context.Context, opts *bind.CallOpts, param bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetTokenPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (*models.SuiTransactionBlockResponse, error)
	GetTimestampedPriceFields(ctx context.Context, opts *bind.CallOpts, tp TimestampedPrice) (*models.SuiTransactionBlockResponse, error)
	GetTokenPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokens []string) (*models.SuiTransactionBlockResponse, error)
	GetDestChainGasPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetTokenAndGasPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, token string, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	ConvertTokenAmount(ctx context.Context, opts *bind.CallOpts, ref bind.Object, fromToken string, fromTokenAmount uint64, toToken string) (*models.SuiTransactionBlockResponse, error)
	GetFeeTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	ApplyFeeTokenUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, feeTokensToRemove []string, feeTokensToAdd []string) (*models.SuiTransactionBlockResponse, error)
	GetTokenTransferFeeConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, token string) (*models.SuiTransactionBlockResponse, error)
	ApplyTokenTransferFeeConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, destChainSelector uint64, addTokens []string, addMinFeeUsdCents []uint32, addMaxFeeUsdCents []uint32, addDeciBps []uint16, addDestGasOverhead []uint32, addDestBytesOverhead []uint32, addIsEnabled []bool, removeTokens []string) (*models.SuiTransactionBlockResponse, error)
	UpdatePrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, clock bind.Object, sourceTokens []string, sourceUsdPerToken []*big.Int, gasDestChainSelectors []uint64, gasUsdPerUnitGas []*big.Int) (*models.SuiTransactionBlockResponse, error)
	GetValidatedFee(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) (*models.SuiTransactionBlockResponse, error)
	ApplyPremiumMultiplierWeiPerEthUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, tokens []string, premiumMultiplierWeiPerEth []uint64) (*models.SuiTransactionBlockResponse, error)
	GetPremiumMultiplierWeiPerEth(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (*models.SuiTransactionBlockResponse, error)
	GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, extraArgs []byte, messageReceiver []byte) (*models.SuiTransactionBlockResponse, error)
	ProcessMessageArgs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) (*models.SuiTransactionBlockResponse, error)
	GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error)
	GetDestChainConfigFields(ctx context.Context, opts *bind.CallOpts, destChainConfig DestChainConfig) (*models.SuiTransactionBlockResponse, error)
	ApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, destChainSelector uint64, isEnabled bool, maxNumberOfTokensPerMsg uint16, maxDataBytes uint32, maxPerMsgGasLimit uint32, destGasOverhead uint32, destGasPerPayloadByteBase byte, destGasPerPayloadByteHigh byte, destGasPerPayloadByteThreshold uint16, destDataAvailabilityOverheadGas uint32, destGasPerDataAvailabilityByte uint16, destDataAvailabilityMultiplierBps uint16, chainFamilySelector []byte, enforceOutOfOrder bool, defaultTokenFeeUsdCents uint16, defaultTokenDestGasOverhead uint32, defaultTxGasLimit uint32, gasMultiplierWeiPerEth uint64, gasPriceStalenessThreshold uint32, networkFeeUsdCents uint32) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error)
	GetTokenTransferFeeConfigFields(ctx context.Context, opts *bind.CallOpts, cfg TokenTransferFeeConfig) (*models.SuiTransactionBlockResponse, error)
	McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error)
	DevInspect() IFeeQuoterDevInspect
	Encoder() FeeQuoterEncoder
}

type IFeeQuoterDevInspect interface {
	TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error)
	NewFeeQuoterCap(ctx context.Context, opts *bind.CallOpts, param bind.Object) (bind.Object, error)
	GetTokenPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (TimestampedPrice, error)
	GetTimestampedPriceFields(ctx context.Context, opts *bind.CallOpts, tp TimestampedPrice) ([]any, error)
	GetTokenPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokens []string) ([]TimestampedPrice, error)
	GetDestChainGasPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (TimestampedPrice, error)
	GetTokenAndGasPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, token string, destChainSelector uint64) ([]any, error)
	ConvertTokenAmount(ctx context.Context, opts *bind.CallOpts, ref bind.Object, fromToken string, fromTokenAmount uint64, toToken string) (uint64, error)
	GetFeeTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([]string, error)
	GetTokenTransferFeeConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, token string) (TokenTransferFeeConfig, error)
	GetValidatedFee(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) (uint64, error)
	GetPremiumMultiplierWeiPerEth(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (uint64, error)
	GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, extraArgs []byte, messageReceiver []byte) ([]byte, error)
	ProcessMessageArgs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) ([]any, error)
	GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (DestChainConfig, error)
	GetDestChainConfigFields(ctx context.Context, opts *bind.CallOpts, destChainConfig DestChainConfig) ([]any, error)
	GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (StaticConfig, error)
	GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) ([]any, error)
	GetTokenTransferFeeConfigFields(ctx context.Context, opts *bind.CallOpts, cfg TokenTransferFeeConfig) ([]any, error)
}

type FeeQuoterEncoder interface {
	TypeAndVersion() (*bind.EncodedCall, error)
	TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error)
	Initialize(ref bind.Object, ownerCap bind.Object, maxFeeJuelsPerMsg *big.Int, linkToken string, tokenPriceStalenessThreshold uint64, feeTokens []string) (*bind.EncodedCall, error)
	InitializeWithArgs(args ...any) (*bind.EncodedCall, error)
	IssueFeeQuoterCap(ownerCap bind.Object) (*bind.EncodedCall, error)
	IssueFeeQuoterCapWithArgs(args ...any) (*bind.EncodedCall, error)
	NewFeeQuoterCap(param bind.Object) (*bind.EncodedCall, error)
	NewFeeQuoterCapWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenPrice(ref bind.Object, token string) (*bind.EncodedCall, error)
	GetTokenPriceWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTimestampedPriceFields(tp TimestampedPrice) (*bind.EncodedCall, error)
	GetTimestampedPriceFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenPrices(ref bind.Object, tokens []string) (*bind.EncodedCall, error)
	GetTokenPricesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestChainGasPrice(ref bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetDestChainGasPriceWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenAndGasPrices(ref bind.Object, clock bind.Object, token string, destChainSelector uint64) (*bind.EncodedCall, error)
	GetTokenAndGasPricesWithArgs(args ...any) (*bind.EncodedCall, error)
	ConvertTokenAmount(ref bind.Object, fromToken string, fromTokenAmount uint64, toToken string) (*bind.EncodedCall, error)
	ConvertTokenAmountWithArgs(args ...any) (*bind.EncodedCall, error)
	GetFeeTokens(ref bind.Object) (*bind.EncodedCall, error)
	GetFeeTokensWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyFeeTokenUpdates(ref bind.Object, ownerCap bind.Object, feeTokensToRemove []string, feeTokensToAdd []string) (*bind.EncodedCall, error)
	ApplyFeeTokenUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenTransferFeeConfig(ref bind.Object, destChainSelector uint64, token string) (*bind.EncodedCall, error)
	GetTokenTransferFeeConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyTokenTransferFeeConfigUpdates(ref bind.Object, ownerCap bind.Object, destChainSelector uint64, addTokens []string, addMinFeeUsdCents []uint32, addMaxFeeUsdCents []uint32, addDeciBps []uint16, addDestGasOverhead []uint32, addDestBytesOverhead []uint32, addIsEnabled []bool, removeTokens []string) (*bind.EncodedCall, error)
	ApplyTokenTransferFeeConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	UpdatePrices(ref bind.Object, param bind.Object, clock bind.Object, sourceTokens []string, sourceUsdPerToken []*big.Int, gasDestChainSelectors []uint64, gasUsdPerUnitGas []*big.Int) (*bind.EncodedCall, error)
	UpdatePricesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetValidatedFee(ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) (*bind.EncodedCall, error)
	GetValidatedFeeWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyPremiumMultiplierWeiPerEthUpdates(ref bind.Object, ownerCap bind.Object, tokens []string, premiumMultiplierWeiPerEth []uint64) (*bind.EncodedCall, error)
	ApplyPremiumMultiplierWeiPerEthUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetPremiumMultiplierWeiPerEth(ref bind.Object, token string) (*bind.EncodedCall, error)
	GetPremiumMultiplierWeiPerEthWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenReceiver(ref bind.Object, destChainSelector uint64, extraArgs []byte, messageReceiver []byte) (*bind.EncodedCall, error)
	GetTokenReceiverWithArgs(args ...any) (*bind.EncodedCall, error)
	ProcessMessageArgs(ref bind.Object, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) (*bind.EncodedCall, error)
	ProcessMessageArgsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestChainConfig(ref bind.Object, destChainSelector uint64) (*bind.EncodedCall, error)
	GetDestChainConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetDestChainConfigFields(destChainConfig DestChainConfig) (*bind.EncodedCall, error)
	GetDestChainConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	ApplyDestChainConfigUpdates(ref bind.Object, ownerCap bind.Object, destChainSelector uint64, isEnabled bool, maxNumberOfTokensPerMsg uint16, maxDataBytes uint32, maxPerMsgGasLimit uint32, destGasOverhead uint32, destGasPerPayloadByteBase byte, destGasPerPayloadByteHigh byte, destGasPerPayloadByteThreshold uint16, destDataAvailabilityOverheadGas uint32, destGasPerDataAvailabilityByte uint16, destDataAvailabilityMultiplierBps uint16, chainFamilySelector []byte, enforceOutOfOrder bool, defaultTokenFeeUsdCents uint16, defaultTokenDestGasOverhead uint32, defaultTxGasLimit uint32, gasMultiplierWeiPerEth uint64, gasPriceStalenessThreshold uint32, networkFeeUsdCents uint32) (*bind.EncodedCall, error)
	ApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfig(ref bind.Object) (*bind.EncodedCall, error)
	GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error)
	GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error)
	GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	GetTokenTransferFeeConfigFields(cfg TokenTransferFeeConfig) (*bind.EncodedCall, error)
	GetTokenTransferFeeConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error)
	McmsEntrypoint(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error)
	McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error)
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

func NewFeeQuoter(packageID string, client sui.ISuiAPI) (*FeeQuoterContract, error) {
	contract, err := bind.NewBoundContract(packageID, "ccip", "fee_quoter", client)
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

func (c *FeeQuoterContract) Encoder() FeeQuoterEncoder {
	return c.feeQuoterEncoder
}

func (c *FeeQuoterContract) DevInspect() IFeeQuoterDevInspect {
	return c.devInspect
}

type FeeQuoterState struct {
	Id                           string      `move:"sui::object::UID"`
	MaxFeeJuelsPerMsg            *big.Int    `move:"u256"`
	LinkToken                    string      `move:"address"`
	TokenPriceStalenessThreshold uint64      `move:"u64"`
	FeeTokens                    []string    `move:"vector<address>"`
	UsdPerUnitGasByDestChain     bind.Object `move:"table::Table<u64, TimestampedPrice>"`
	UsdPerToken                  bind.Object `move:"table::Table<address, TimestampedPrice>"`
	DestChainConfigs             bind.Object `move:"table::Table<u64, DestChainConfig>"`
	TokenTransferFeeConfigs      bind.Object `move:"table::Table<u64, table::Table<address, TokenTransferFeeConfig>>"`
	PremiumMultiplierWeiPerEth   bind.Object `move:"table::Table<address, u64>"`
}

type FeeQuoterCap struct {
	Id string `move:"sui::object::UID"`
}

type StaticConfig struct {
	MaxFeeJuelsPerMsg            *big.Int `move:"u256"`
	LinkToken                    string   `move:"address"`
	TokenPriceStalenessThreshold uint64   `move:"u64"`
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

type TokenTransferFeeConfig struct {
	MinFeeUsdCents    uint32 `move:"u32"`
	MaxFeeUsdCents    uint32 `move:"u32"`
	DeciBps           uint16 `move:"u16"`
	DestGasOverhead   uint32 `move:"u32"`
	DestBytesOverhead uint32 `move:"u32"`
	IsEnabled         bool   `move:"bool"`
}

type TimestampedPrice struct {
	Value     *big.Int `move:"u256"`
	Timestamp uint64   `move:"u64"`
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

type CCIPAdminProof struct {
}

type McmsCallback struct {
}

type bcsFeeQuoterState struct {
	Id                           string
	MaxFeeJuelsPerMsg            [32]byte
	LinkToken                    [32]byte
	TokenPriceStalenessThreshold uint64
	FeeTokens                    [][32]byte
	UsdPerUnitGasByDestChain     bind.Object
	UsdPerToken                  bind.Object
	DestChainConfigs             bind.Object
	TokenTransferFeeConfigs      bind.Object
	PremiumMultiplierWeiPerEth   bind.Object
}

func convertFeeQuoterStateFromBCS(bcs bcsFeeQuoterState) (FeeQuoterState, error) {
	MaxFeeJuelsPerMsgField, err := bind.DecodeU256Value(bcs.MaxFeeJuelsPerMsg)
	if err != nil {
		return FeeQuoterState{}, fmt.Errorf("failed to decode u256 field MaxFeeJuelsPerMsg: %w", err)
	}

	return FeeQuoterState{
		Id:                           bcs.Id,
		MaxFeeJuelsPerMsg:            MaxFeeJuelsPerMsgField,
		LinkToken:                    fmt.Sprintf("0x%x", bcs.LinkToken),
		TokenPriceStalenessThreshold: bcs.TokenPriceStalenessThreshold,
		FeeTokens: func() []string {
			addrs := make([]string, len(bcs.FeeTokens))
			for i, addr := range bcs.FeeTokens {
				addrs[i] = fmt.Sprintf("0x%x", addr)
			}
			return addrs
		}(),
		UsdPerUnitGasByDestChain:   bcs.UsdPerUnitGasByDestChain,
		UsdPerToken:                bcs.UsdPerToken,
		DestChainConfigs:           bcs.DestChainConfigs,
		TokenTransferFeeConfigs:    bcs.TokenTransferFeeConfigs,
		PremiumMultiplierWeiPerEth: bcs.PremiumMultiplierWeiPerEth,
	}, nil
}

type bcsStaticConfig struct {
	MaxFeeJuelsPerMsg            [32]byte
	LinkToken                    [32]byte
	TokenPriceStalenessThreshold uint64
}

func convertStaticConfigFromBCS(bcs bcsStaticConfig) (StaticConfig, error) {
	MaxFeeJuelsPerMsgField, err := bind.DecodeU256Value(bcs.MaxFeeJuelsPerMsg)
	if err != nil {
		return StaticConfig{}, fmt.Errorf("failed to decode u256 field MaxFeeJuelsPerMsg: %w", err)
	}

	return StaticConfig{
		MaxFeeJuelsPerMsg:            MaxFeeJuelsPerMsgField,
		LinkToken:                    fmt.Sprintf("0x%x", bcs.LinkToken),
		TokenPriceStalenessThreshold: bcs.TokenPriceStalenessThreshold,
	}, nil
}

type bcsTimestampedPrice struct {
	Value     [32]byte
	Timestamp uint64
}

func convertTimestampedPriceFromBCS(bcs bcsTimestampedPrice) (TimestampedPrice, error) {
	ValueField, err := bind.DecodeU256Value(bcs.Value)
	if err != nil {
		return TimestampedPrice{}, fmt.Errorf("failed to decode u256 field Value: %w", err)
	}

	return TimestampedPrice{
		Value:     ValueField,
		Timestamp: bcs.Timestamp,
	}, nil
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
	bind.RegisterStructDecoder("ccip::fee_quoter::FeeQuoterState", func(data []byte) (interface{}, error) {
		var temp bcsFeeQuoterState
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertFeeQuoterStateFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::FeeQuoterCap", func(data []byte) (interface{}, error) {
		var result FeeQuoterCap
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::StaticConfig", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::DestChainConfig", func(data []byte) (interface{}, error) {
		var result DestChainConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::TokenTransferFeeConfig", func(data []byte) (interface{}, error) {
		var result TokenTransferFeeConfig
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::TimestampedPrice", func(data []byte) (interface{}, error) {
		var temp bcsTimestampedPrice
		_, err := mystenbcs.Unmarshal(data, &temp)
		if err != nil {
			return nil, err
		}

		result, err := convertTimestampedPriceFromBCS(temp)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::FeeTokenAdded", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::FeeTokenRemoved", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::TokenTransferFeeConfigAdded", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::TokenTransferFeeConfigRemoved", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::UsdPerTokenUpdated", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::UsdPerUnitGasUpdated", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::DestChainAdded", func(data []byte) (interface{}, error) {
		var result DestChainAdded
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::DestChainConfigUpdated", func(data []byte) (interface{}, error) {
		var result DestChainConfigUpdated
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::PremiumMultiplierWeiPerEthUpdated", func(data []byte) (interface{}, error) {
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
	bind.RegisterStructDecoder("ccip::fee_quoter::CCIPAdminProof", func(data []byte) (interface{}, error) {
		var result CCIPAdminProof
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
	bind.RegisterStructDecoder("ccip::fee_quoter::McmsCallback", func(data []byte) (interface{}, error) {
		var result McmsCallback
		_, err := mystenbcs.Unmarshal(data, &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	})
}

// TypeAndVersion executes the type_and_version Move function.
func (c *FeeQuoterContract) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.TypeAndVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// Initialize executes the initialize Move function.
func (c *FeeQuoterContract) Initialize(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, maxFeeJuelsPerMsg *big.Int, linkToken string, tokenPriceStalenessThreshold uint64, feeTokens []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.Initialize(ref, ownerCap, maxFeeJuelsPerMsg, linkToken, tokenPriceStalenessThreshold, feeTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// IssueFeeQuoterCap executes the issue_fee_quoter_cap Move function.
func (c *FeeQuoterContract) IssueFeeQuoterCap(ctx context.Context, opts *bind.CallOpts, ownerCap bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.IssueFeeQuoterCap(ownerCap)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// NewFeeQuoterCap executes the new_fee_quoter_cap Move function.
func (c *FeeQuoterContract) NewFeeQuoterCap(ctx context.Context, opts *bind.CallOpts, param bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.NewFeeQuoterCap(param)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenPrice executes the get_token_price Move function.
func (c *FeeQuoterContract) GetTokenPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetTokenPrice(ref, token)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTimestampedPriceFields executes the get_timestamped_price_fields Move function.
func (c *FeeQuoterContract) GetTimestampedPriceFields(ctx context.Context, opts *bind.CallOpts, tp TimestampedPrice) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetTimestampedPriceFields(tp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenPrices executes the get_token_prices Move function.
func (c *FeeQuoterContract) GetTokenPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokens []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetTokenPrices(ref, tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestChainGasPrice executes the get_dest_chain_gas_price Move function.
func (c *FeeQuoterContract) GetDestChainGasPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetDestChainGasPrice(ref, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenAndGasPrices executes the get_token_and_gas_prices Move function.
func (c *FeeQuoterContract) GetTokenAndGasPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, token string, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetTokenAndGasPrices(ref, clock, token, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ConvertTokenAmount executes the convert_token_amount Move function.
func (c *FeeQuoterContract) ConvertTokenAmount(ctx context.Context, opts *bind.CallOpts, ref bind.Object, fromToken string, fromTokenAmount uint64, toToken string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.ConvertTokenAmount(ref, fromToken, fromTokenAmount, toToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetFeeTokens executes the get_fee_tokens Move function.
func (c *FeeQuoterContract) GetFeeTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetFeeTokens(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyFeeTokenUpdates executes the apply_fee_token_updates Move function.
func (c *FeeQuoterContract) ApplyFeeTokenUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, feeTokensToRemove []string, feeTokensToAdd []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.ApplyFeeTokenUpdates(ref, ownerCap, feeTokensToRemove, feeTokensToAdd)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenTransferFeeConfig executes the get_token_transfer_fee_config Move function.
func (c *FeeQuoterContract) GetTokenTransferFeeConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, token string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetTokenTransferFeeConfig(ref, destChainSelector, token)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyTokenTransferFeeConfigUpdates executes the apply_token_transfer_fee_config_updates Move function.
func (c *FeeQuoterContract) ApplyTokenTransferFeeConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, destChainSelector uint64, addTokens []string, addMinFeeUsdCents []uint32, addMaxFeeUsdCents []uint32, addDeciBps []uint16, addDestGasOverhead []uint32, addDestBytesOverhead []uint32, addIsEnabled []bool, removeTokens []string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.ApplyTokenTransferFeeConfigUpdates(ref, ownerCap, destChainSelector, addTokens, addMinFeeUsdCents, addMaxFeeUsdCents, addDeciBps, addDestGasOverhead, addDestBytesOverhead, addIsEnabled, removeTokens)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// UpdatePrices executes the update_prices Move function.
func (c *FeeQuoterContract) UpdatePrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, param bind.Object, clock bind.Object, sourceTokens []string, sourceUsdPerToken []*big.Int, gasDestChainSelectors []uint64, gasUsdPerUnitGas []*big.Int) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.UpdatePrices(ref, param, clock, sourceTokens, sourceUsdPerToken, gasDestChainSelectors, gasUsdPerUnitGas)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetValidatedFee executes the get_validated_fee Move function.
func (c *FeeQuoterContract) GetValidatedFee(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetValidatedFee(ref, clock, destChainSelector, receiver, data, localTokenAddresses, localTokenAmounts, feeToken, extraArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyPremiumMultiplierWeiPerEthUpdates executes the apply_premium_multiplier_wei_per_eth_updates Move function.
func (c *FeeQuoterContract) ApplyPremiumMultiplierWeiPerEthUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, tokens []string, premiumMultiplierWeiPerEth []uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.ApplyPremiumMultiplierWeiPerEthUpdates(ref, ownerCap, tokens, premiumMultiplierWeiPerEth)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetPremiumMultiplierWeiPerEth executes the get_premium_multiplier_wei_per_eth Move function.
func (c *FeeQuoterContract) GetPremiumMultiplierWeiPerEth(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetPremiumMultiplierWeiPerEth(ref, token)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenReceiver executes the get_token_receiver Move function.
func (c *FeeQuoterContract) GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, extraArgs []byte, messageReceiver []byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetTokenReceiver(ref, destChainSelector, extraArgs, messageReceiver)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ProcessMessageArgs executes the process_message_args Move function.
func (c *FeeQuoterContract) ProcessMessageArgs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.ProcessMessageArgs(ref, destChainSelector, feeToken, feeTokenAmount, extraArgs, localTokenAddresses, destTokenAddresses, destPoolDatas)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestChainConfig executes the get_dest_chain_config Move function.
func (c *FeeQuoterContract) GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetDestChainConfig(ref, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetDestChainConfigFields executes the get_dest_chain_config_fields Move function.
func (c *FeeQuoterContract) GetDestChainConfigFields(ctx context.Context, opts *bind.CallOpts, destChainConfig DestChainConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetDestChainConfigFields(destChainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// ApplyDestChainConfigUpdates executes the apply_dest_chain_config_updates Move function.
func (c *FeeQuoterContract) ApplyDestChainConfigUpdates(ctx context.Context, opts *bind.CallOpts, ref bind.Object, ownerCap bind.Object, destChainSelector uint64, isEnabled bool, maxNumberOfTokensPerMsg uint16, maxDataBytes uint32, maxPerMsgGasLimit uint32, destGasOverhead uint32, destGasPerPayloadByteBase byte, destGasPerPayloadByteHigh byte, destGasPerPayloadByteThreshold uint16, destDataAvailabilityOverheadGas uint32, destGasPerDataAvailabilityByte uint16, destDataAvailabilityMultiplierBps uint16, chainFamilySelector []byte, enforceOutOfOrder bool, defaultTokenFeeUsdCents uint16, defaultTokenDestGasOverhead uint32, defaultTxGasLimit uint32, gasMultiplierWeiPerEth uint64, gasPriceStalenessThreshold uint32, networkFeeUsdCents uint32) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.ApplyDestChainConfigUpdates(ref, ownerCap, destChainSelector, isEnabled, maxNumberOfTokensPerMsg, maxDataBytes, maxPerMsgGasLimit, destGasOverhead, destGasPerPayloadByteBase, destGasPerPayloadByteHigh, destGasPerPayloadByteThreshold, destDataAvailabilityOverheadGas, destGasPerDataAvailabilityByte, destDataAvailabilityMultiplierBps, chainFamilySelector, enforceOutOfOrder, defaultTokenFeeUsdCents, defaultTokenDestGasOverhead, defaultTxGasLimit, gasMultiplierWeiPerEth, gasPriceStalenessThreshold, networkFeeUsdCents)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function.
func (c *FeeQuoterContract) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetStaticConfig(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetStaticConfigFields executes the get_static_config_fields Move function.
func (c *FeeQuoterContract) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetStaticConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// GetTokenTransferFeeConfigFields executes the get_token_transfer_fee_config_fields Move function.
func (c *FeeQuoterContract) GetTokenTransferFeeConfigFields(ctx context.Context, opts *bind.CallOpts, cfg TokenTransferFeeConfig) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.GetTokenTransferFeeConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// McmsEntrypoint executes the mcms_entrypoint Move function.
func (c *FeeQuoterContract) McmsEntrypoint(ctx context.Context, opts *bind.CallOpts, ref bind.Object, registry bind.Object, params bind.Object) (*models.SuiTransactionBlockResponse, error) {
	encoded, err := c.feeQuoterEncoder.McmsEntrypoint(ref, registry, params)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}

	return c.ExecuteTransaction(ctx, opts, encoded)
}

// TypeAndVersion executes the type_and_version Move function using DevInspect to get return values.
//
// Returns: 0x1::string::String
func (d *FeeQuoterDevInspect) TypeAndVersion(ctx context.Context, opts *bind.CallOpts) (string, error) {
	encoded, err := d.contract.feeQuoterEncoder.TypeAndVersion()
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

// NewFeeQuoterCap executes the new_fee_quoter_cap Move function using DevInspect to get return values.
//
// Returns: FeeQuoterCap
func (d *FeeQuoterDevInspect) NewFeeQuoterCap(ctx context.Context, opts *bind.CallOpts, param bind.Object) (bind.Object, error) {
	encoded, err := d.contract.feeQuoterEncoder.NewFeeQuoterCap(param)
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

// GetTokenPrice executes the get_token_price Move function using DevInspect to get return values.
//
// Returns: TimestampedPrice
func (d *FeeQuoterDevInspect) GetTokenPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (TimestampedPrice, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetTokenPrice(ref, token)
	if err != nil {
		return TimestampedPrice{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TimestampedPrice{}, err
	}
	if len(results) == 0 {
		return TimestampedPrice{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TimestampedPrice)
	if !ok {
		return TimestampedPrice{}, fmt.Errorf("unexpected return type: expected TimestampedPrice, got %T", results[0])
	}
	return result, nil
}

// GetTimestampedPriceFields executes the get_timestamped_price_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u256
//	[1]: u64
func (d *FeeQuoterDevInspect) GetTimestampedPriceFields(ctx context.Context, opts *bind.CallOpts, tp TimestampedPrice) ([]any, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetTimestampedPriceFields(tp)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetTokenPrices executes the get_token_prices Move function using DevInspect to get return values.
//
// Returns: vector<TimestampedPrice>
func (d *FeeQuoterDevInspect) GetTokenPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, tokens []string) ([]TimestampedPrice, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetTokenPrices(ref, tokens)
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
	result, ok := results[0].([]TimestampedPrice)
	if !ok {
		return nil, fmt.Errorf("unexpected return type: expected []TimestampedPrice, got %T", results[0])
	}
	return result, nil
}

// GetDestChainGasPrice executes the get_dest_chain_gas_price Move function using DevInspect to get return values.
//
// Returns: TimestampedPrice
func (d *FeeQuoterDevInspect) GetDestChainGasPrice(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (TimestampedPrice, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetDestChainGasPrice(ref, destChainSelector)
	if err != nil {
		return TimestampedPrice{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TimestampedPrice{}, err
	}
	if len(results) == 0 {
		return TimestampedPrice{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TimestampedPrice)
	if !ok {
		return TimestampedPrice{}, fmt.Errorf("unexpected return type: expected TimestampedPrice, got %T", results[0])
	}
	return result, nil
}

// GetTokenAndGasPrices executes the get_token_and_gas_prices Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u256
//	[1]: u256
func (d *FeeQuoterDevInspect) GetTokenAndGasPrices(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, token string, destChainSelector uint64) ([]any, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetTokenAndGasPrices(ref, clock, token, destChainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// ConvertTokenAmount executes the convert_token_amount Move function using DevInspect to get return values.
//
// Returns: u64
func (d *FeeQuoterDevInspect) ConvertTokenAmount(ctx context.Context, opts *bind.CallOpts, ref bind.Object, fromToken string, fromTokenAmount uint64, toToken string) (uint64, error) {
	encoded, err := d.contract.feeQuoterEncoder.ConvertTokenAmount(ref, fromToken, fromTokenAmount, toToken)
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

// GetFeeTokens executes the get_fee_tokens Move function using DevInspect to get return values.
//
// Returns: vector<address>
func (d *FeeQuoterDevInspect) GetFeeTokens(ctx context.Context, opts *bind.CallOpts, ref bind.Object) ([]string, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetFeeTokens(ref)
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

// GetTokenTransferFeeConfig executes the get_token_transfer_fee_config Move function using DevInspect to get return values.
//
// Returns: TokenTransferFeeConfig
func (d *FeeQuoterDevInspect) GetTokenTransferFeeConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, token string) (TokenTransferFeeConfig, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetTokenTransferFeeConfig(ref, destChainSelector, token)
	if err != nil {
		return TokenTransferFeeConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return TokenTransferFeeConfig{}, err
	}
	if len(results) == 0 {
		return TokenTransferFeeConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(TokenTransferFeeConfig)
	if !ok {
		return TokenTransferFeeConfig{}, fmt.Errorf("unexpected return type: expected TokenTransferFeeConfig, got %T", results[0])
	}
	return result, nil
}

// GetValidatedFee executes the get_validated_fee Move function using DevInspect to get return values.
//
// Returns: u64
func (d *FeeQuoterDevInspect) GetValidatedFee(ctx context.Context, opts *bind.CallOpts, ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) (uint64, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetValidatedFee(ref, clock, destChainSelector, receiver, data, localTokenAddresses, localTokenAmounts, feeToken, extraArgs)
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

// GetPremiumMultiplierWeiPerEth executes the get_premium_multiplier_wei_per_eth Move function using DevInspect to get return values.
//
// Returns: u64
func (d *FeeQuoterDevInspect) GetPremiumMultiplierWeiPerEth(ctx context.Context, opts *bind.CallOpts, ref bind.Object, token string) (uint64, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetPremiumMultiplierWeiPerEth(ref, token)
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

// GetTokenReceiver executes the get_token_receiver Move function using DevInspect to get return values.
//
// Returns: vector<u8>
func (d *FeeQuoterDevInspect) GetTokenReceiver(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, extraArgs []byte, messageReceiver []byte) ([]byte, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetTokenReceiver(ref, destChainSelector, extraArgs, messageReceiver)
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

// ProcessMessageArgs executes the process_message_args Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u256
//	[1]: bool
//	[2]: vector<u8>
//	[3]: vector<vector<u8>>
func (d *FeeQuoterDevInspect) ProcessMessageArgs(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) ([]any, error) {
	encoded, err := d.contract.feeQuoterEncoder.ProcessMessageArgs(ref, destChainSelector, feeToken, feeTokenAmount, extraArgs, localTokenAddresses, destTokenAddresses, destPoolDatas)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetDestChainConfig executes the get_dest_chain_config Move function using DevInspect to get return values.
//
// Returns: DestChainConfig
func (d *FeeQuoterDevInspect) GetDestChainConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object, destChainSelector uint64) (DestChainConfig, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetDestChainConfig(ref, destChainSelector)
	if err != nil {
		return DestChainConfig{}, fmt.Errorf("failed to encode function call: %w", err)
	}
	results, err := d.contract.Call(ctx, opts, encoded)
	if err != nil {
		return DestChainConfig{}, err
	}
	if len(results) == 0 {
		return DestChainConfig{}, fmt.Errorf("no return value")
	}
	result, ok := results[0].(DestChainConfig)
	if !ok {
		return DestChainConfig{}, fmt.Errorf("unexpected return type: expected DestChainConfig, got %T", results[0])
	}
	return result, nil
}

// GetDestChainConfigFields executes the get_dest_chain_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: bool
//	[1]: u16
//	[2]: u32
//	[3]: u32
//	[4]: u32
//	[5]: u8
//	[6]: u8
//	[7]: u16
//	[8]: u32
//	[9]: u16
//	[10]: u16
//	[11]: vector<u8>
//	[12]: bool
//	[13]: u16
//	[14]: u32
//	[15]: u32
//	[16]: u64
//	[17]: u32
//	[18]: u32
func (d *FeeQuoterDevInspect) GetDestChainConfigFields(ctx context.Context, opts *bind.CallOpts, destChainConfig DestChainConfig) ([]any, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetDestChainConfigFields(destChainConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetStaticConfig executes the get_static_config Move function using DevInspect to get return values.
//
// Returns: StaticConfig
func (d *FeeQuoterDevInspect) GetStaticConfig(ctx context.Context, opts *bind.CallOpts, ref bind.Object) (StaticConfig, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetStaticConfig(ref)
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
//	[0]: u256
//	[1]: address
//	[2]: u64
func (d *FeeQuoterDevInspect) GetStaticConfigFields(ctx context.Context, opts *bind.CallOpts, cfg StaticConfig) ([]any, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetStaticConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

// GetTokenTransferFeeConfigFields executes the get_token_transfer_fee_config_fields Move function using DevInspect to get return values.
//
// Returns:
//
//	[0]: u32
//	[1]: u32
//	[2]: u16
//	[3]: u32
//	[4]: u32
//	[5]: bool
func (d *FeeQuoterDevInspect) GetTokenTransferFeeConfigFields(ctx context.Context, opts *bind.CallOpts, cfg TokenTransferFeeConfig) ([]any, error) {
	encoded, err := d.contract.feeQuoterEncoder.GetTokenTransferFeeConfigFields(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to encode function call: %w", err)
	}
	return d.contract.Call(ctx, opts, encoded)
}

type feeQuoterEncoder struct {
	*bind.BoundContract
}

// TypeAndVersion encodes a call to the type_and_version Move function.
func (c feeQuoterEncoder) TypeAndVersion() (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("type_and_version", typeArgsList, typeParamsList, []string{}, []any{}, []string{
		"0x1::string::String",
	})
}

// TypeAndVersionWithArgs encodes a call to the type_and_version Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) TypeAndVersionWithArgs(args ...any) (*bind.EncodedCall, error) {
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
func (c feeQuoterEncoder) Initialize(ref bind.Object, ownerCap bind.Object, maxFeeJuelsPerMsg *big.Int, linkToken string, tokenPriceStalenessThreshold uint64, feeTokens []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u256",
		"address",
		"u64",
		"vector<address>",
	}, []any{
		ref,
		ownerCap,
		maxFeeJuelsPerMsg,
		linkToken,
		tokenPriceStalenessThreshold,
		feeTokens,
	}, nil)
}

// InitializeWithArgs encodes a call to the initialize Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) InitializeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u256",
		"address",
		"u64",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("initialize", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// IssueFeeQuoterCap encodes a call to the issue_fee_quoter_cap Move function.
func (c feeQuoterEncoder) IssueFeeQuoterCap(ownerCap bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("issue_fee_quoter_cap", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
	}, []any{
		ownerCap,
	}, nil)
}

// IssueFeeQuoterCapWithArgs encodes a call to the issue_fee_quoter_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) IssueFeeQuoterCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("issue_fee_quoter_cap", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// NewFeeQuoterCap encodes a call to the new_fee_quoter_cap Move function.
func (c feeQuoterEncoder) NewFeeQuoterCap(param bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_fee_quoter_cap", typeArgsList, typeParamsList, []string{
		"&OwnerCap",
	}, []any{
		param,
	}, []string{
		"ccip::fee_quoter::FeeQuoterCap",
	})
}

// NewFeeQuoterCapWithArgs encodes a call to the new_fee_quoter_cap Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) NewFeeQuoterCapWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&OwnerCap",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("new_fee_quoter_cap", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::fee_quoter::FeeQuoterCap",
	})
}

// GetTokenPrice encodes a call to the get_token_price Move function.
func (c feeQuoterEncoder) GetTokenPrice(ref bind.Object, token string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_price", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		token,
	}, []string{
		"ccip::fee_quoter::TimestampedPrice",
	})
}

// GetTokenPriceWithArgs encodes a call to the get_token_price Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetTokenPriceWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_price", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::fee_quoter::TimestampedPrice",
	})
}

// GetTimestampedPriceFields encodes a call to the get_timestamped_price_fields Move function.
func (c feeQuoterEncoder) GetTimestampedPriceFields(tp TimestampedPrice) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_timestamped_price_fields", typeArgsList, typeParamsList, []string{
		"ccip::fee_quoter::TimestampedPrice",
	}, []any{
		tp,
	}, []string{
		"u256",
		"u64",
	})
}

// GetTimestampedPriceFieldsWithArgs encodes a call to the get_timestamped_price_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetTimestampedPriceFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip::fee_quoter::TimestampedPrice",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_timestamped_price_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u256",
		"u64",
	})
}

// GetTokenPrices encodes a call to the get_token_prices Move function.
func (c feeQuoterEncoder) GetTokenPrices(ref bind.Object, tokens []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_prices", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"vector<address>",
	}, []any{
		ref,
		tokens,
	}, []string{
		"vector<ccip::fee_quoter::TimestampedPrice>",
	})
}

// GetTokenPricesWithArgs encodes a call to the get_token_prices Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetTokenPricesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_prices", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<ccip::fee_quoter::TimestampedPrice>",
	})
}

// GetDestChainGasPrice encodes a call to the get_dest_chain_gas_price Move function.
func (c feeQuoterEncoder) GetDestChainGasPrice(ref bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_gas_price", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
	}, []any{
		ref,
		destChainSelector,
	}, []string{
		"ccip::fee_quoter::TimestampedPrice",
	})
}

// GetDestChainGasPriceWithArgs encodes a call to the get_dest_chain_gas_price Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetDestChainGasPriceWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_gas_price", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::fee_quoter::TimestampedPrice",
	})
}

// GetTokenAndGasPrices encodes a call to the get_token_and_gas_prices Move function.
func (c feeQuoterEncoder) GetTokenAndGasPrices(ref bind.Object, clock bind.Object, token string, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_and_gas_prices", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&clock::Clock",
		"address",
		"u64",
	}, []any{
		ref,
		clock,
		token,
		destChainSelector,
	}, []string{
		"u256",
		"u256",
	})
}

// GetTokenAndGasPricesWithArgs encodes a call to the get_token_and_gas_prices Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetTokenAndGasPricesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&clock::Clock",
		"address",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_and_gas_prices", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u256",
		"u256",
	})
}

// ConvertTokenAmount encodes a call to the convert_token_amount Move function.
func (c feeQuoterEncoder) ConvertTokenAmount(ref bind.Object, fromToken string, fromTokenAmount uint64, toToken string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("convert_token_amount", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
		"u64",
		"address",
	}, []any{
		ref,
		fromToken,
		fromTokenAmount,
		toToken,
	}, []string{
		"u64",
	})
}

// ConvertTokenAmountWithArgs encodes a call to the convert_token_amount Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) ConvertTokenAmountWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
		"u64",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("convert_token_amount", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetFeeTokens encodes a call to the get_fee_tokens Move function.
func (c feeQuoterEncoder) GetFeeTokens(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_fee_tokens", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"vector<address>",
	})
}

// GetFeeTokensWithArgs encodes a call to the get_fee_tokens Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetFeeTokensWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_fee_tokens", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<address>",
	})
}

// ApplyFeeTokenUpdates encodes a call to the apply_fee_token_updates Move function.
func (c feeQuoterEncoder) ApplyFeeTokenUpdates(ref bind.Object, ownerCap bind.Object, feeTokensToRemove []string, feeTokensToAdd []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_fee_token_updates", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
		"vector<address>",
	}, []any{
		ref,
		ownerCap,
		feeTokensToRemove,
		feeTokensToAdd,
	}, nil)
}

// ApplyFeeTokenUpdatesWithArgs encodes a call to the apply_fee_token_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) ApplyFeeTokenUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_fee_token_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetTokenTransferFeeConfig encodes a call to the get_token_transfer_fee_config Move function.
func (c feeQuoterEncoder) GetTokenTransferFeeConfig(ref bind.Object, destChainSelector uint64, token string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_transfer_fee_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"address",
	}, []any{
		ref,
		destChainSelector,
		token,
	}, []string{
		"ccip::fee_quoter::TokenTransferFeeConfig",
	})
}

// GetTokenTransferFeeConfigWithArgs encodes a call to the get_token_transfer_fee_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetTokenTransferFeeConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
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
	return c.EncodeCallArgsWithGenerics("get_token_transfer_fee_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::fee_quoter::TokenTransferFeeConfig",
	})
}

// ApplyTokenTransferFeeConfigUpdates encodes a call to the apply_token_transfer_fee_config_updates Move function.
func (c feeQuoterEncoder) ApplyTokenTransferFeeConfigUpdates(ref bind.Object, ownerCap bind.Object, destChainSelector uint64, addTokens []string, addMinFeeUsdCents []uint32, addMaxFeeUsdCents []uint32, addDeciBps []uint16, addDestGasOverhead []uint32, addDestBytesOverhead []uint32, addIsEnabled []bool, removeTokens []string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_token_transfer_fee_config_updates", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u64",
		"vector<address>",
		"vector<u32>",
		"vector<u32>",
		"vector<u16>",
		"vector<u32>",
		"vector<u32>",
		"vector<bool>",
		"vector<address>",
	}, []any{
		ref,
		ownerCap,
		destChainSelector,
		addTokens,
		addMinFeeUsdCents,
		addMaxFeeUsdCents,
		addDeciBps,
		addDestGasOverhead,
		addDestBytesOverhead,
		addIsEnabled,
		removeTokens,
	}, nil)
}

// ApplyTokenTransferFeeConfigUpdatesWithArgs encodes a call to the apply_token_transfer_fee_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) ApplyTokenTransferFeeConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u64",
		"vector<address>",
		"vector<u32>",
		"vector<u32>",
		"vector<u16>",
		"vector<u32>",
		"vector<u32>",
		"vector<bool>",
		"vector<address>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_token_transfer_fee_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// UpdatePrices encodes a call to the update_prices Move function.
func (c feeQuoterEncoder) UpdatePrices(ref bind.Object, param bind.Object, clock bind.Object, sourceTokens []string, sourceUsdPerToken []*big.Int, gasDestChainSelectors []uint64, gasUsdPerUnitGas []*big.Int) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("update_prices", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&FeeQuoterCap",
		"&clock::Clock",
		"vector<address>",
		"vector<u256>",
		"vector<u64>",
		"vector<u256>",
	}, []any{
		ref,
		param,
		clock,
		sourceTokens,
		sourceUsdPerToken,
		gasDestChainSelectors,
		gasUsdPerUnitGas,
	}, nil)
}

// UpdatePricesWithArgs encodes a call to the update_prices Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) UpdatePricesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&FeeQuoterCap",
		"&clock::Clock",
		"vector<address>",
		"vector<u256>",
		"vector<u64>",
		"vector<u256>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("update_prices", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetValidatedFee encodes a call to the get_validated_fee Move function.
func (c feeQuoterEncoder) GetValidatedFee(ref bind.Object, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_validated_fee", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"&clock::Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"vector<address>",
		"vector<u64>",
		"address",
		"vector<u8>",
	}, []any{
		ref,
		clock,
		destChainSelector,
		receiver,
		data,
		localTokenAddresses,
		localTokenAmounts,
		feeToken,
		extraArgs,
	}, []string{
		"u64",
	})
}

// GetValidatedFeeWithArgs encodes a call to the get_validated_fee Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetValidatedFeeWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"&clock::Clock",
		"u64",
		"vector<u8>",
		"vector<u8>",
		"vector<address>",
		"vector<u64>",
		"address",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_validated_fee", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// ApplyPremiumMultiplierWeiPerEthUpdates encodes a call to the apply_premium_multiplier_wei_per_eth_updates Move function.
func (c feeQuoterEncoder) ApplyPremiumMultiplierWeiPerEthUpdates(ref bind.Object, ownerCap bind.Object, tokens []string, premiumMultiplierWeiPerEth []uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_premium_multiplier_wei_per_eth_updates", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
		"vector<u64>",
	}, []any{
		ref,
		ownerCap,
		tokens,
		premiumMultiplierWeiPerEth,
	}, nil)
}

// ApplyPremiumMultiplierWeiPerEthUpdatesWithArgs encodes a call to the apply_premium_multiplier_wei_per_eth_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) ApplyPremiumMultiplierWeiPerEthUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"vector<address>",
		"vector<u64>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_premium_multiplier_wei_per_eth_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetPremiumMultiplierWeiPerEth encodes a call to the get_premium_multiplier_wei_per_eth Move function.
func (c feeQuoterEncoder) GetPremiumMultiplierWeiPerEth(ref bind.Object, token string) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_premium_multiplier_wei_per_eth", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"address",
	}, []any{
		ref,
		token,
	}, []string{
		"u64",
	})
}

// GetPremiumMultiplierWeiPerEthWithArgs encodes a call to the get_premium_multiplier_wei_per_eth Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetPremiumMultiplierWeiPerEthWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"address",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_premium_multiplier_wei_per_eth", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u64",
	})
}

// GetTokenReceiver encodes a call to the get_token_receiver Move function.
func (c feeQuoterEncoder) GetTokenReceiver(ref bind.Object, destChainSelector uint64, extraArgs []byte, messageReceiver []byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_receiver", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"vector<u8>",
		"vector<u8>",
	}, []any{
		ref,
		destChainSelector,
		extraArgs,
		messageReceiver,
	}, []string{
		"vector<u8>",
	})
}

// GetTokenReceiverWithArgs encodes a call to the get_token_receiver Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetTokenReceiverWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
		"vector<u8>",
		"vector<u8>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_receiver", typeArgsList, typeParamsList, expectedParams, args, []string{
		"vector<u8>",
	})
}

// ProcessMessageArgs encodes a call to the process_message_args Move function.
func (c feeQuoterEncoder) ProcessMessageArgs(ref bind.Object, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("process_message_args", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
		"address",
		"u64",
		"vector<u8>",
		"vector<address>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
	}, []any{
		ref,
		destChainSelector,
		feeToken,
		feeTokenAmount,
		extraArgs,
		localTokenAddresses,
		destTokenAddresses,
		destPoolDatas,
	}, []string{
		"u256",
		"bool",
		"vector<u8>",
		"vector<vector<u8>>",
	})
}

// ProcessMessageArgsWithArgs encodes a call to the process_message_args Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) ProcessMessageArgsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
		"address",
		"u64",
		"vector<u8>",
		"vector<address>",
		"vector<vector<u8>>",
		"vector<vector<u8>>",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("process_message_args", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u256",
		"bool",
		"vector<u8>",
		"vector<vector<u8>>",
	})
}

// GetDestChainConfig encodes a call to the get_dest_chain_config Move function.
func (c feeQuoterEncoder) GetDestChainConfig(ref bind.Object, destChainSelector uint64) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
		"u64",
	}, []any{
		ref,
		destChainSelector,
	}, []string{
		"ccip::fee_quoter::DestChainConfig",
	})
}

// GetDestChainConfigWithArgs encodes a call to the get_dest_chain_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetDestChainConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
		"u64",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::fee_quoter::DestChainConfig",
	})
}

// GetDestChainConfigFields encodes a call to the get_dest_chain_config_fields Move function.
func (c feeQuoterEncoder) GetDestChainConfigFields(destChainConfig DestChainConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config_fields", typeArgsList, typeParamsList, []string{
		"ccip::fee_quoter::DestChainConfig",
	}, []any{
		destChainConfig,
	}, []string{
		"bool",
		"u16",
		"u32",
		"u32",
		"u32",
		"u8",
		"u8",
		"u16",
		"u32",
		"u16",
		"u16",
		"vector<u8>",
		"bool",
		"u16",
		"u32",
		"u32",
		"u64",
		"u32",
		"u32",
	})
}

// GetDestChainConfigFieldsWithArgs encodes a call to the get_dest_chain_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetDestChainConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip::fee_quoter::DestChainConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_dest_chain_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"bool",
		"u16",
		"u32",
		"u32",
		"u32",
		"u8",
		"u8",
		"u16",
		"u32",
		"u16",
		"u16",
		"vector<u8>",
		"bool",
		"u16",
		"u32",
		"u32",
		"u64",
		"u32",
		"u32",
	})
}

// ApplyDestChainConfigUpdates encodes a call to the apply_dest_chain_config_updates Move function.
func (c feeQuoterEncoder) ApplyDestChainConfigUpdates(ref bind.Object, ownerCap bind.Object, destChainSelector uint64, isEnabled bool, maxNumberOfTokensPerMsg uint16, maxDataBytes uint32, maxPerMsgGasLimit uint32, destGasOverhead uint32, destGasPerPayloadByteBase byte, destGasPerPayloadByteHigh byte, destGasPerPayloadByteThreshold uint16, destDataAvailabilityOverheadGas uint32, destGasPerDataAvailabilityByte uint16, destDataAvailabilityMultiplierBps uint16, chainFamilySelector []byte, enforceOutOfOrder bool, defaultTokenFeeUsdCents uint16, defaultTokenDestGasOverhead uint32, defaultTxGasLimit uint32, gasMultiplierWeiPerEth uint64, gasPriceStalenessThreshold uint32, networkFeeUsdCents uint32) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_dest_chain_config_updates", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u64",
		"bool",
		"u16",
		"u32",
		"u32",
		"u32",
		"u8",
		"u8",
		"u16",
		"u32",
		"u16",
		"u16",
		"vector<u8>",
		"bool",
		"u16",
		"u32",
		"u32",
		"u64",
		"u32",
		"u32",
	}, []any{
		ref,
		ownerCap,
		destChainSelector,
		isEnabled,
		maxNumberOfTokensPerMsg,
		maxDataBytes,
		maxPerMsgGasLimit,
		destGasOverhead,
		destGasPerPayloadByteBase,
		destGasPerPayloadByteHigh,
		destGasPerPayloadByteThreshold,
		destDataAvailabilityOverheadGas,
		destGasPerDataAvailabilityByte,
		destDataAvailabilityMultiplierBps,
		chainFamilySelector,
		enforceOutOfOrder,
		defaultTokenFeeUsdCents,
		defaultTokenDestGasOverhead,
		defaultTxGasLimit,
		gasMultiplierWeiPerEth,
		gasPriceStalenessThreshold,
		networkFeeUsdCents,
	}, nil)
}

// ApplyDestChainConfigUpdatesWithArgs encodes a call to the apply_dest_chain_config_updates Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) ApplyDestChainConfigUpdatesWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
		"&OwnerCap",
		"u64",
		"bool",
		"u16",
		"u32",
		"u32",
		"u32",
		"u8",
		"u8",
		"u16",
		"u32",
		"u16",
		"u16",
		"vector<u8>",
		"bool",
		"u16",
		"u32",
		"u32",
		"u64",
		"u32",
		"u32",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("apply_dest_chain_config_updates", typeArgsList, typeParamsList, expectedParams, args, nil)
}

// GetStaticConfig encodes a call to the get_static_config Move function.
func (c feeQuoterEncoder) GetStaticConfig(ref bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, []string{
		"&CCIPObjectRef",
	}, []any{
		ref,
	}, []string{
		"ccip::fee_quoter::StaticConfig",
	})
}

// GetStaticConfigWithArgs encodes a call to the get_static_config Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetStaticConfigWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&CCIPObjectRef",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config", typeArgsList, typeParamsList, expectedParams, args, []string{
		"ccip::fee_quoter::StaticConfig",
	})
}

// GetStaticConfigFields encodes a call to the get_static_config_fields Move function.
func (c feeQuoterEncoder) GetStaticConfigFields(cfg StaticConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, []string{
		"ccip::fee_quoter::StaticConfig",
	}, []any{
		cfg,
	}, []string{
		"u256",
		"address",
		"u64",
	})
}

// GetStaticConfigFieldsWithArgs encodes a call to the get_static_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetStaticConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip::fee_quoter::StaticConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_static_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u256",
		"address",
		"u64",
	})
}

// GetTokenTransferFeeConfigFields encodes a call to the get_token_transfer_fee_config_fields Move function.
func (c feeQuoterEncoder) GetTokenTransferFeeConfigFields(cfg TokenTransferFeeConfig) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_transfer_fee_config_fields", typeArgsList, typeParamsList, []string{
		"ccip::fee_quoter::TokenTransferFeeConfig",
	}, []any{
		cfg,
	}, []string{
		"u32",
		"u32",
		"u16",
		"u32",
		"u32",
		"bool",
	})
}

// GetTokenTransferFeeConfigFieldsWithArgs encodes a call to the get_token_transfer_fee_config_fields Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) GetTokenTransferFeeConfigFieldsWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"ccip::fee_quoter::TokenTransferFeeConfig",
	}

	if len(args) != len(expectedParams) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(expectedParams), len(args))
	}
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("get_token_transfer_fee_config_fields", typeArgsList, typeParamsList, expectedParams, args, []string{
		"u32",
		"u32",
		"u16",
		"u32",
		"u32",
		"bool",
	})
}

// McmsEntrypoint encodes a call to the mcms_entrypoint Move function.
func (c feeQuoterEncoder) McmsEntrypoint(ref bind.Object, registry bind.Object, params bind.Object) (*bind.EncodedCall, error) {
	typeArgsList := []string{}
	typeParamsList := []string{}
	return c.EncodeCallArgsWithGenerics("mcms_entrypoint", typeArgsList, typeParamsList, []string{
		"&mut CCIPObjectRef",
		"&mut Registry",
		"ExecutingCallbackParams",
	}, []any{
		ref,
		registry,
		params,
	}, nil)
}

// McmsEntrypointWithArgs encodes a call to the mcms_entrypoint Move function using arbitrary arguments.
// This method allows passing both regular values and transaction.Argument values for PTB chaining.
func (c feeQuoterEncoder) McmsEntrypointWithArgs(args ...any) (*bind.EncodedCall, error) {
	expectedParams := []string{
		"&mut CCIPObjectRef",
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
