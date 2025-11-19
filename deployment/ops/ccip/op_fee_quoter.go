package ccipops

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
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

	feeQuoterCapObjectId, err1 := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterCap")
	if err1 != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to find fee quoter cap object ID in tx: %w", err1)
	}
	feeQuoterStateObjectId, err2 := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterState")
	if err2 != nil {
		return sui_ops.OpTxResult[InitFeeQuoterObjects]{}, fmt.Errorf("failed to find fee quoter state object ID in tx: %w", err2)
	}

	return sui_ops.OpTxResult[InitFeeQuoterObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects: InitFeeQuoterObjects{
			FeeQuoterCapObjectId:   feeQuoterCapObjectId,
			FeeQuoterStateObjectId: feeQuoterStateObjectId,
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

	encodedCall, err := contract.Encoder().ApplyFeeTokenUpdates(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.FeeTokensToRemove,
		input.FeeTokensToAdd,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode ApplyFeeTokenUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyFeeTokenUpdates on FeeQuoter as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ApplyFeeTokenUpdates on FeeQuoter: %w", err)
	}

	b.Logger.Infow("Fee token updates applied to CCIP FeeQuoter")

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
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

	encodedCall, err := contract.Encoder().ApplyTokenTransferFeeConfigUpdates(
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode ApplyTokenTransferFeeConfigUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyTokenTransferFeeConfigUpdates on FeeQuoter as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ApplyTokenTransferFeeConfigUpdates on FeeQuoter: %w", err)
	}

	b.Logger.Infow("Token transfer fee config updates applied to CCIP FeeQuoter")

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
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

	encodedCall, err := contract.Encoder().ApplyDestChainConfigUpdates(
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
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyDestChainConfigUpdates on FeeQuoter as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ApplyDestChainConfigUpdates on FeeQuoter: %w", err)
	}

	b.Logger.Infow("Destination chain config updates applied to CCIP FeeQuoter")

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
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

	encodedCall, err := contract.Encoder().ApplyPremiumMultiplierWeiPerEthUpdates(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.Tokens,
		input.PremiumMultiplierWeiPerEth,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode ApplyPremiumMultiplierWeiPerEthUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyPremiumMultiplierWeiPerEthUpdates on FeeQuoter as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ApplyPremiumMultiplierWeiPerEthUpdates on FeeQuoter: %w", err)
	}

	b.Logger.Infow("Premium multiplier wei per eth updates applied to CCIP FeeQuoter")

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
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
	SourceTokens          []string
	SourceUsdPerToken     []*big.Int
	GasDestChainSelectors []uint64
	GasUsdPerUnitGas      []*big.Int

	FeeQuoterCapId string // optional: only provide if you're running direct UpdateTokenPrice
	OwnerCapId     string // optional: only provide if you're running UpdateTokenPricesWithOwnerCap
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
	}, err
}

var FeeQuoterUpdateTokenPricesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "update_prices"),
	semver.MustParse("0.1.0"),
	"Apply update prices in CCIP Fee Quoter contract",
	updateTokenPrices,
)

var updateTokenPricesWithOwnerCap = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input FeeQuoterUpdateTokenPricesInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	fmt.Println("INPUTSS: ", input)
	tx, err := contract.UpdatePricesWithOwnerCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		bind.Object{Id: input.OwnerCapId},
		bind.Object{Id: "0x6"}, // Clock object
		input.SourceTokens,
		input.SourceUsdPerToken,
		input.GasDestChainSelectors,
		input.GasUsdPerUnitGas,
	)

	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute updateTokenPrices  on SUI: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, err
}

var FeeQuoterUpdateTokenPricesWithOwnerCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "update_prices_with_owner_cap"),
	semver.MustParse("0.1.0"),
	"Apply update prices with ownerCap in CCIP Fee Quoter contract",
	updateTokenPricesWithOwnerCap,
)

// FEE QUOTER -- new_fee_quoter_cap
type NewFeeQuoterCapObjects struct {
	FeeQuoterCapObjectId string
}

type NewFeeQuoterCapInput struct {
	CCIPPackageId    string
	CCIPObjectRef    string
	OwnerCapObjectId string
}

var newFeeQuoterCapHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input NewFeeQuoterCapInput) (output sui_ops.OpTxResult[NewFeeQuoterCapObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NewFeeQuoterCapObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.NewFeeQuoterCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		bind.Object{Id: input.OwnerCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[NewFeeQuoterCapObjects]{}, fmt.Errorf("failed to execute new_fee_quoter_cap: %w", err)
	}

	feeQuoterCapObjectId, err1 := bind.FindObjectIdFromPublishTx(*tx, "fee_quoter", "FeeQuoterCap")
	if err1 != nil {
		return sui_ops.OpTxResult[NewFeeQuoterCapObjects]{}, fmt.Errorf("failed to find fee quoter cap object ID in tx: %w", err1)
	}

	return sui_ops.OpTxResult[NewFeeQuoterCapObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects: NewFeeQuoterCapObjects{
			FeeQuoterCapObjectId: feeQuoterCapObjectId,
		},
	}, err
}

var FeeQuoterNewFeeQuoterCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "new_fee_quoter_cap"),
	semver.MustParse("0.1.0"),
	"Create a new fee quoter cap in the CCIP Fee Quoter contract",
	newFeeQuoterCapHandler,
)

// FEE QUOTER -- destroy_fee_quoter_cap
type DestroyFeeQuoterCapInput struct {
	CCIPPackageId        string
	CCIPObjectRef        string
	OwnerCapObjectId     string
	FeeQuoterCapObjectId string
}

var destroyFeeQuoterCapHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DestroyFeeQuoterCapInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.DestroyFeeQuoterCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.FeeQuoterCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute destroy_fee_quoter_cap: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, err
}

var FeeQuoterDestroyFeeQuoterCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "destroy_fee_quoter_cap"),
	semver.MustParse("0.1.0"),
	"Destroy a fee quoter cap in the CCIP Fee Quoter contract",
	destroyFeeQuoterCapHandler,
)

// FEE QUOTER -- update_prices_with_owner_cap
type FeeQuoterUpdatePricesWithOwnerCapInput struct {
	CCIPPackageId         string
	CCIPObjectRef         string
	OwnerCapObjectId      string
	SourceTokens          []string
	SourceUsdPerToken     []*big.Int
	GasDestChainSelectors []uint64
	GasUsdPerUnitGas      []*big.Int
}

var updatePricesWithOwnerCapHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input FeeQuoterUpdatePricesWithOwnerCapInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_fee_quoter.NewFeeQuoter(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UpdatePricesWithOwnerCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: "0x6"}, // Clock object
		input.SourceTokens,
		input.SourceUsdPerToken,
		input.GasDestChainSelectors,
		input.GasUsdPerUnitGas,
	)

	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute updatePricesWithOwnerCap on SUI: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
	}, err
}

var FeeQuoterUpdatePricesWithOwnerCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "fee_quoter", "update_prices_with_owner_cap"),
	semver.MustParse("0.1.0"),
	"Update prices using owner cap in CCIP Fee Quoter contract",
	updatePricesWithOwnerCapHandler,
)
