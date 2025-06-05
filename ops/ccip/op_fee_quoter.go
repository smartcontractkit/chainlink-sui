package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/holiman/uint256"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

// FEE QUOTER -- INITIALIZE
type InitFeeQuoterObjects struct {
	FeeQuoterCapObjectId   string
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

	maxFeeJuels256, err := uint256.FromDecimal(input.MaxFeeJuelsPerMsg)
	if err != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to convert big.Int to uint256: %s", input.MaxFeeJuelsPerMsg)
	}

	method := contract.Initialize(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		*maxFeeJuels256,
		input.LinkTokenCoinMetadataObjectId,
		input.TokenPriceStalenessThreshold,
		input.FeeTokens,
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[InitFeeQuoterObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects: InitFeeQuoterObjects{
			FeeQuoterCapObjectId:   obj1,
			FeeQuoterStateObjectId: obj2,
		},
	}, err
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

	method := contract.ApplyFeeTokenUpdates(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.FeeTokensToRemove,
		input.FeeTokensToAdd,
	)

	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_fee_token_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
	}, err
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

	method := contract.ApplyTokenTransferFeeConfigUpdates(
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

	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_token_transfer_fee_config_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
	}, err
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

	method := contract.ApplyDestChainConfigUpdates(
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

	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_dest_chain_config_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
	}, err
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

	method := contract.ApplyPremiumMultiplierWeiPerEthUpdates(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.Tokens,
		input.PremiumMultiplierWeiPerEth,
	)

	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute fee quoter apply_premium_multiplier_wei_per_eth_updates: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
	}, err
}

var FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "apply_premium_multiplier_wei_per_eth_updates"),
	semver.MustParse("0.1.0"),
	"Apply premium multiplier wei per eth updates in the CCIP Fee Quoter contract",
	applyPremiumMultiplierHandler,
)
