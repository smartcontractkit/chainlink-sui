// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package module_fee_quoter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_common "github.com/smartcontractkit/chainlink-sui/bindings/common"
)

// Unused vars used for unused imports
var (
	_ = big.NewInt
	_ = uint256.NewInt
)

type IFeeQuoter interface {
	TypeAndVersion() bind.IMethod
	Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap, maxFeeJuelsPerMsg uint256.Int, linkToken string, tokenPriceStalenessThreshold uint64, feeTokens []string) bind.IMethod
	ApplyFeeTokenUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, feeTokensToRemove []string, feeTokensToAdd []string) bind.IMethod
	ApplyTokenTransferFeeConfigUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, destChainSelector uint64, addTokens []string, addMinFeeUsdCents []uint32, addMaxFeeUsdCents []uint32, addDeciBps []uint16, addDestGasOverhead []uint32, addDestBytesOverhead []uint32, addIsEnabled []bool, removeTokens []string) bind.IMethod
	ApplyDestChainConfigUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, destChainSelector uint64, isEnabled bool, maxNumberOfTokensPerMsg uint16, maxDataBytes uint32, maxPerMsgGasLimit uint32, destGasOverhead uint32, destGasPerPayloadByteBase byte, destGasPerPayloadByteHigh byte, destGasPerPayloadByteThreshold uint16, destDataAvailabilityOverheadGas uint32, destGasPerDataAvailabilityByte uint16, destDataAvailabilityMultiplierBps uint16, chainFamilySelector []byte, enforceOutOfOrder bool, defaultTokenFeeUsdCents uint16, defaultTokenDestGasOverhead uint32, defaultTxGasLimit uint32, gasMultiplierWeiPerEth uint64, gasPriceStalenessThreshold uint32, networkFeeUsdCents uint32) bind.IMethod
	ApplyPremiumMultiplierWeiPerEthUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, tokens []string, premiumMultiplierWeiPerEth []uint64) bind.IMethod
	GetStaticConfig(ref module_common.CCIPObjectRef) bind.IMethod
	GetStaticConfigFields(cfg StaticConfig) bind.IMethod
	GetTokenTransferFeeConfig(ref module_common.CCIPObjectRef, destChainSelector uint64, token string) bind.IMethod
	GetTokenTransferFeeConfigFields(cfg TokenTransferFeeConfig) bind.IMethod
	GetTokenPrice(ref module_common.CCIPObjectRef, token string) bind.IMethod
	GetTimestampedPriceFields(tp TimestampedPrice) bind.IMethod
	GetTokenPrices(ref module_common.CCIPObjectRef, tokens []string) bind.IMethod
	GetDestChainGasPrice(ref module_common.CCIPObjectRef, destChainSelector uint64) bind.IMethod
	GetTokenAndGasPrices(ref module_common.CCIPObjectRef, clock bind.Object, token string, destChainSelector uint64) bind.IMethod
	GetDestChainConfig(ref module_common.CCIPObjectRef, destChainSelector uint64) bind.IMethod
	GetDestChainConfigFields(destChainConfig DestChainConfig) bind.IMethod
	UpdatePrices(ref module_common.CCIPObjectRef, param bind.Object, clock bind.Object, sourceTokens []string, sourceUsdPerToken []uint256.Int, gasDestChainSelectors []uint64, gasUsdPerUnitGas []uint256.Int) bind.IMethod
	GetValidatedFee(ref module_common.CCIPObjectRef, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) bind.IMethod
	GetPremiumMultiplierWeiPerEth(ref module_common.CCIPObjectRef, token string) bind.IMethod
	ProcessMessageArgs(ref module_common.CCIPObjectRef, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) bind.IMethod
	GetFeeTokens(ref module_common.CCIPObjectRef) bind.IMethod
	// Connect adds/changes the client used in the contract
	Connect(client suiclient.ClientImpl)
}

type FeeQuoterContract struct {
	packageID *sui.Address
	client    suiclient.ClientImpl
}

var _ IFeeQuoter = (*FeeQuoterContract)(nil)

func NewFeeQuoter(packageID string, client suiclient.ClientImpl) (*FeeQuoterContract, error) {
	pkgObjectId, err := bind.ToSuiAddress(packageID)
	if err != nil {
		return nil, fmt.Errorf("package ID is not a Sui address: %w", err)
	}

	return &FeeQuoterContract{
		packageID: pkgObjectId,
		client:    client,
	}, nil
}

func (c *FeeQuoterContract) Connect(client suiclient.ClientImpl) {
	c.client = client
}

// Structs

type FeeQuoterState struct {
	Id                           string      `move:"sui::object::UID"`
	MaxFeeJuelsPerMsg            uint256.Int `move:"u256"`
	LinkToken                    string      `move:"address"`
	TokenPriceStalenessThreshold uint64      `move:"u64"`
	FeeTokens                    []string    `move:"vector<address>"`
}

type FeeQuoterCap struct {
	Id string `move:"sui::object::UID"`
}

type StaticConfig struct {
	MaxFeeJuelsPerMsg            uint256.Int `move:"u256"`
	LinkToken                    string      `move:"address"`
	TokenPriceStalenessThreshold uint64      `move:"u64"`
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
	Value     uint256.Int `move:"u256"`
	Timestamp uint64      `move:"u64"`
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
	Token       string      `move:"address"`
	UsdPerToken uint256.Int `move:"u256"`
	Timestamp   uint64      `move:"u64"`
}

type UsdPerUnitGasUpdated struct {
	DestChainSelector uint64      `move:"u64"`
	UsdPerUnitGas     uint256.Int `move:"u256"`
	Timestamp         uint64      `move:"u64"`
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

// Functions

func (c *FeeQuoterContract) TypeAndVersion() bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "type_and_version", false, "")
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "type_and_version", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) Initialize(ref module_common.CCIPObjectRef, param module_common.OwnerCap, maxFeeJuelsPerMsg uint256.Int, linkToken string, tokenPriceStalenessThreshold uint64, feeTokens []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "initialize", false, "", ref, param, maxFeeJuelsPerMsg, linkToken, tokenPriceStalenessThreshold, feeTokens)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "initialize", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) ApplyFeeTokenUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, feeTokensToRemove []string, feeTokensToAdd []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "apply_fee_token_updates", false, "", ref, param, feeTokensToRemove, feeTokensToAdd)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "apply_fee_token_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) ApplyTokenTransferFeeConfigUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, destChainSelector uint64, addTokens []string, addMinFeeUsdCents []uint32, addMaxFeeUsdCents []uint32, addDeciBps []uint16, addDestGasOverhead []uint32, addDestBytesOverhead []uint32, addIsEnabled []bool, removeTokens []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "apply_token_transfer_fee_config_updates", false, "", ref, param, destChainSelector, addTokens, addMinFeeUsdCents, addMaxFeeUsdCents, addDeciBps, addDestGasOverhead, addDestBytesOverhead, addIsEnabled, removeTokens)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "apply_token_transfer_fee_config_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) ApplyDestChainConfigUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, destChainSelector uint64, isEnabled bool, maxNumberOfTokensPerMsg uint16, maxDataBytes uint32, maxPerMsgGasLimit uint32, destGasOverhead uint32, destGasPerPayloadByteBase byte, destGasPerPayloadByteHigh byte, destGasPerPayloadByteThreshold uint16, destDataAvailabilityOverheadGas uint32, destGasPerDataAvailabilityByte uint16, destDataAvailabilityMultiplierBps uint16, chainFamilySelector []byte, enforceOutOfOrder bool, defaultTokenFeeUsdCents uint16, defaultTokenDestGasOverhead uint32, defaultTxGasLimit uint32, gasMultiplierWeiPerEth uint64, gasPriceStalenessThreshold uint32, networkFeeUsdCents uint32) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "apply_dest_chain_config_updates", false, "", ref, param, destChainSelector, isEnabled, maxNumberOfTokensPerMsg, maxDataBytes, maxPerMsgGasLimit, destGasOverhead, destGasPerPayloadByteBase, destGasPerPayloadByteHigh, destGasPerPayloadByteThreshold, destDataAvailabilityOverheadGas, destGasPerDataAvailabilityByte, destDataAvailabilityMultiplierBps, chainFamilySelector, enforceOutOfOrder, defaultTokenFeeUsdCents, defaultTokenDestGasOverhead, defaultTxGasLimit, gasMultiplierWeiPerEth, gasPriceStalenessThreshold, networkFeeUsdCents)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "apply_dest_chain_config_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) ApplyPremiumMultiplierWeiPerEthUpdates(ref module_common.CCIPObjectRef, param module_common.OwnerCap, tokens []string, premiumMultiplierWeiPerEth []uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "apply_premium_multiplier_wei_per_eth_updates", false, "", ref, param, tokens, premiumMultiplierWeiPerEth)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "apply_premium_multiplier_wei_per_eth_updates", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetStaticConfig(ref module_common.CCIPObjectRef) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_static_config", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_static_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetStaticConfigFields(cfg StaticConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_static_config_fields", false, "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_static_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetTokenTransferFeeConfig(ref module_common.CCIPObjectRef, destChainSelector uint64, token string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_token_transfer_fee_config", false, "", ref, destChainSelector, token)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_token_transfer_fee_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetTokenTransferFeeConfigFields(cfg TokenTransferFeeConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_token_transfer_fee_config_fields", false, "", cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_token_transfer_fee_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetTokenPrice(ref module_common.CCIPObjectRef, token string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_token_price", false, "", ref, token)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_token_price", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetTimestampedPriceFields(tp TimestampedPrice) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_timestamped_price_fields", false, "", tp)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_timestamped_price_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetTokenPrices(ref module_common.CCIPObjectRef, tokens []string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_token_prices", false, "", ref, tokens)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_token_prices", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetDestChainGasPrice(ref module_common.CCIPObjectRef, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_dest_chain_gas_price", false, "", ref, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_dest_chain_gas_price", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetTokenAndGasPrices(ref module_common.CCIPObjectRef, clock bind.Object, token string, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_token_and_gas_prices", false, "", ref, clock, token, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_token_and_gas_prices", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetDestChainConfig(ref module_common.CCIPObjectRef, destChainSelector uint64) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_dest_chain_config", false, "", ref, destChainSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_dest_chain_config", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetDestChainConfigFields(destChainConfig DestChainConfig) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_dest_chain_config_fields", false, "", destChainConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_dest_chain_config_fields", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) UpdatePrices(ref module_common.CCIPObjectRef, param bind.Object, clock bind.Object, sourceTokens []string, sourceUsdPerToken []uint256.Int, gasDestChainSelectors []uint64, gasUsdPerUnitGas []uint256.Int) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "update_prices", false, "", ref, param, clock, sourceTokens, sourceUsdPerToken, gasDestChainSelectors, gasUsdPerUnitGas)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "update_prices", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetValidatedFee(ref module_common.CCIPObjectRef, clock bind.Object, destChainSelector uint64, receiver []byte, data []byte, localTokenAddresses []string, localTokenAmounts []uint64, feeToken string, extraArgs []byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_validated_fee", false, "", ref, clock, destChainSelector, receiver, data, localTokenAddresses, localTokenAmounts, feeToken, extraArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_validated_fee", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetPremiumMultiplierWeiPerEth(ref module_common.CCIPObjectRef, token string) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_premium_multiplier_wei_per_eth", false, "", ref, token)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_premium_multiplier_wei_per_eth", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) ProcessMessageArgs(ref module_common.CCIPObjectRef, destChainSelector uint64, feeToken string, feeTokenAmount uint64, extraArgs []byte, localTokenAddresses []string, destTokenAddresses [][]byte, destPoolDatas [][]byte) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "process_message_args", false, "", ref, destChainSelector, feeToken, feeTokenAmount, extraArgs, localTokenAddresses, destTokenAddresses, destPoolDatas)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "process_message_args", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}

func (c *FeeQuoterContract) GetFeeTokens(ref module_common.CCIPObjectRef) bind.IMethod {
	build := func(ctx context.Context) (*suiptb.ProgrammableTransactionBuilder, error) {
		// TODO: Object creation is always set to false. Contract analyzer should check if the function uses ::transfer
		ptb, err := bind.BuildPTBFromArgs(ctx, c.client, c.packageID, "fee_quoter", "get_fee_tokens", false, "", ref)
		if err != nil {
			return nil, fmt.Errorf("failed to build PTB for moudule %v in function %v: %w", "fee_quoter", "get_fee_tokens", err)
		}

		return ptb, nil
	}

	return bind.NewMethod(build, bind.MakeExecute(build), bind.MakeInspect(build))
}
