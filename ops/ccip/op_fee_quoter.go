package ccipops

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

// FEE QUOTER -- INITIALIZE
type InitFeeQuoterObjects struct {
	FeeQuoterCapObjectId string
	FeeQuoterStateObjectId string
}

type InitFeeQuoterInput struct {
	CCIPPackageId                 string
	StateObjectId                 string
	OwnerCapObjectId              string
	MaxFeeJuelsPerMsg             string
	LinkTokenCoinMetadataObjectId string
	TokenPriceStalenessThreshold  uint64
	FeeTokens                     []string
}

var initFQHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitFeeQuoterInput) (output sui_ops.OpTxResult[InitFeeQuoterObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	const decimalBase = 10
	maxFeeJuels, ok := new(big.Int).SetString(input.MaxFeeJuelsPerMsg, decimalBase)
	if !ok {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to parse MaxFeeJuelsPerMsg: %s", input.MaxFeeJuelsPerMsg)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Initialize(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		maxFeeJuels,
		input.LinkTokenCoinMetadataObjectId,
		input.TokenPriceStalenessThreshold,
		input.FeeTokens,
	)
	if err != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	capObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterCap")
	if err != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to find FeeQuoterCap object ID in tx: %w", err)
	}

	stateObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterState")
	if err != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to find FeeQuoterState object ID in tx: %w", err)
	}

	return sui_ops.OpTxResult[InitFeeQuoterObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects: InitFeeQuoterObjects{
			FeeQuoterCapObjectId: capObjectId,
			FeeQuoterStateObjectId: stateObjectId,
		},
	}, nil
}

var FeeQuoterInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Fee Quoter contract",
	initFQHandler,
)

// FEE QUOTER -- apply_fee_token_updates
type NoObjects struct{}

type FeeQuoterApplyFeeTokenUpdatesInput struct {
	CCIPPackageId     string
	StateObjectId     string
	OwnerCapObjectId  string
	FeeTokensToRemove []string
	FeeTokensToAdd    []string
}

var applyUpdatesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input FeeQuoterApplyFeeTokenUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ApplyFeeTokenUpdates(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.FeeTokensToRemove,
		input.FeeTokensToAdd,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_fee_token_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, nil
}

var FeeQuoterApplyFeeTokenUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "apply_fee_token_updates"),
	semver.MustParse("0.1.0"),
	"Apply fee token updates in the CCIP Fee Quoter contract",
	applyUpdatesHandler,
)

// FEE QUOTER -- apply_token_transfer_fee_config_updates
type FeeQuoterApplyTokenTransferFeeConfigUpdatesInput struct {
	CCIPPackageId        string
	StateObjectId        string
	OwnerCapObjectId     string
	DestChainSelector    uint64
	AddTokens            []string
	AddMinFeeUsdCents    []uint32
	AddMaxFeeUsdCents    []uint32
	AddDeciBps           []uint16
	AddDestGasOverhead   []uint32
	AddDestBytesOverhead []uint32
	AddIsEnabled         []bool
	RemoveTokens         []string
}

var applyTokenTransferFeeHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input FeeQuoterApplyTokenTransferFeeConfigUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ApplyTokenTransferFeeConfigUpdates(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.DestChainSelector,
		input.AddTokens,
		input.AddMinFeeUsdCents,
		input.AddMaxFeeUsdCents,
		input.AddDeciBps,
		input.AddDestGasOverhead,
		input.AddDestBytesOverhead,
		input.AddIsEnabled,
		input.RemoveTokens,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_token_transfer_fee_config_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, nil
}

var FeeQuoterApplyTokenTransferFeeConfigUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "apply_token_transfer_fee_config_updates"),
	semver.MustParse("0.1.0"),
	"Apply transfer fee config updates in the CCIP Fee Quoter contract",
	applyTokenTransferFeeHandler,
)

// FEE QUOTER -- apply_dest_chain_config_updates
type FeeQuoterApplyDestChainConfigUpdatesInput struct {
	CCIPPackageId                     string
	StateObjectId                     string
	OwnerCapObjectId                  string
	DestChainSelector                 uint64
	IsEnabled                         bool
	MaxNumberOfTokensPerMsg           uint16
	MaxDataBytes                      uint32
	MaxPerMsgGasLimit                 uint32
	DestGasOverhead                   uint32
	DestGasPerPayloadByteBase         byte
	DestGasPerPayloadByteHigh         byte
	DestGasPerPayloadByteThreshold    uint16
	DestDataAvailabilityOverheadGas   uint32
	DestGasPerDataAvailabilityByte    uint16
	DestDataAvailabilityMultiplierBps uint16
	ChainFamilySelector               []byte
	EnforceOutOfOrder                 bool
	DefaultTokenFeeUsdCents           uint16
	DefaultTokenDestGasOverhead       uint32
	DefaultTxGasLimit                 uint32
	GasMultiplierWeiPerEth            uint64
	GasPriceStalenessThreshold        uint32
	NetworkFeeUsdCents                uint32
}

var applyDestChainConfigHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input FeeQuoterApplyDestChainConfigUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ApplyDestChainConfigUpdates(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.DestChainSelector,
		input.IsEnabled,
		input.MaxNumberOfTokensPerMsg,
		input.MaxDataBytes,
		input.MaxPerMsgGasLimit,
		input.DestGasOverhead,
		input.DestGasPerPayloadByteBase,
		input.DestGasPerPayloadByteHigh,
		input.DestGasPerPayloadByteThreshold,
		input.DestDataAvailabilityOverheadGas,
		input.DestGasPerDataAvailabilityByte,
		input.DestDataAvailabilityMultiplierBps,
		input.ChainFamilySelector,
		input.EnforceOutOfOrder,
		input.DefaultTokenFeeUsdCents,
		input.DefaultTokenDestGasOverhead,
		input.DefaultTxGasLimit,
		input.GasMultiplierWeiPerEth,
		input.GasPriceStalenessThreshold,
		input.NetworkFeeUsdCents,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_dest_chain_config_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, nil
}

var FeeQuoterApplyDestChainConfigUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "apply_dest_chain_config_updates"),
	semver.MustParse("0.1.0"),
	"Apply destination chain config updates in the CCIP Fee Quoter contract",
	applyDestChainConfigHandler,
)

// FEE QUOTER -- apply_premium_multiplier_wei_per_eth_updates
type FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput struct {
	CCIPPackageId              string
	StateObjectId              string
	OwnerCapObjectId           string
	Tokens                     []string
	PremiumMultiplierWeiPerEth []uint64
}

var applyPremiumMultiplierHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.ApplyPremiumMultiplierWeiPerEthUpdates(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.Tokens,
		input.PremiumMultiplierWeiPerEth,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_premium_multiplier_wei_per_eth_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, nil
}

var FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "apply_premium_multiplier_wei_per_eth_updates"),
	semver.MustParse("0.1.0"),
	"Apply premium multiplier wei per eth updates in the CCIP Fee Quoter contract",
	applyPremiumMultiplierHandler,
)

type FeeQuoterUpdateTokenPricesInput struct {
	CCIPPackageId         string
	CCIPObjectRef         string
	FeeQuoterCapId        string
	SourceTokens          []string
	SourceUsdPerToken     []*big.Int
	GasDestChainSelectors []uint64
	GasUsdPerUnitGas      []*big.Int
}

var updateTokenPrices = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input FeeQuoterUpdateTokenPricesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UpdatePrices(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		bind.Object{Id: input.FeeQuoterCapId},
		bind.Object{Id: "0x6"}, // Clock object
		input.SourceTokens,
		input.SourceUsdPerToken,
		input.GasDestChainSelectors,
		input.GasUsdPerUnitGas,
	)

	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute updateTokenPrices on SUI: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, nil
}

var FeeQuoterUpdateTokenPricesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "update_prices"),
	semver.MustParse("0.1.0"),
	"Apply update prices in CCIP Fee Quoter contract",
	updateTokenPrices,
)

// FEE QUOTER -- issue_fee_quoter_cap
type IssueFeeQuoterCapObjects struct {
	FeeQuoterCapObjectId string
}

type IssueFeeQuoterCapInput struct {
	CCIPPackageId    string
	OwnerCapObjectId string
}

var issueFeeQuoterCapHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IssueFeeQuoterCapInput) (output sui_ops.OpTxResult[IssueFeeQuoterCapObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IssueFeeQuoterCapObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.IssueFeeQuoterCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.OwnerCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[IssueFeeQuoterCapObjects]{}, fmt.Errorf("failed to execute fee quoter issue_fee_quoter_cap: %w", err)
	}

	capObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterCap")
	if err != nil {
		return sui_ops.OpTxResult[IssueFeeQuoterCapObjects]{}, fmt.Errorf("failed to find FeeQuoterCap object ID in tx: %w", err)
	}

	return sui_ops.OpTxResult[IssueFeeQuoterCapObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects: IssueFeeQuoterCapObjects{
			FeeQuoterCapObjectId: capObjectId,
		},
	}, nil
}

var FeeQuoterIssueFeeQuoterCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "issue_fee_quoter_cap"),
	semver.MustParse("0.1.0"),
	"Issue a new FeeQuoterCap in the CCIP Fee Quoter contract",
	issueFeeQuoterCapHandler,
)
