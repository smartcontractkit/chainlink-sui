package onrampops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/onramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"

	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
)

type DeployCCIPOnRampObjects struct {
	// State Object
	OwnerCapObjectId        string
	CCIPOnrampStateObjectId string
}

type DeployCCIPOnRampInput struct {
	CCIPPackageId      string
	MCMSPackageId      string
	MCMSOwnerPackageId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployCCIPOnRampInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	onRampPackage, tx, err := onramp.PublishOnramp(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.MCMSPackageId,
		input.MCMSOwnerPackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "onramp", "OnRampState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: onRampPackage.Address(),
		Objects: DeployCCIPOnRampObjects{
			OwnerCapObjectId:        obj1,
			CCIPOnrampStateObjectId: obj2,
		},
	}, err
}

type OnRampInitializeInput struct {
	OnRampPackageId           string
	OnRampStateId             string
	OwnerCapObjectId          string
	NonceManagerCapId         string
	SourceTransferCapId       string
	ChainSelector             uint64
	FeeAggregator             string
	AllowListAdmin            string
	DestChainSelectors        []uint64
	DestChainEnabled          []bool
	DestChainAllowListEnabled []bool
}

var InitializeHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input OnRampInitializeInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.Initialize(
		b.GetContext(),
		opts,
		bind.Object{Id: input.OnRampStateId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.NonceManagerCapId},
		bind.Object{Id: input.SourceTransferCapId},
		input.ChainSelector,
		input.FeeAggregator,
		input.AllowListAdmin,
		input.DestChainSelectors,
		input.DestChainEnabled,
		input.DestChainAllowListEnabled,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute onRamp initialization: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
	}, err
}

type ApplyDestChainConfigureOnRampInput struct {
	OnRampPackageId           string
	OwnerCapObjectId          string
	StateObjectId             string
	DestChainSelector         []uint64
	DestChainEnabled          []bool
	DestChainAllowListEnabled []bool
}

var ApplyDestChainUpdateHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ApplyDestChainConfigureOnRampInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.ApplyDestChainConfigUpdates(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.DestChainSelector,
		input.DestChainEnabled,
		input.DestChainAllowListEnabled,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute ApplyDestChainUpdate on onRamp: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
	}, err
}

type ApplyAllowListUpdatesInput struct {
	OnRampPackageId               string
	OwnerCapObjectId              string
	StateObjectId                 string
	DestChainSelector             []uint64
	DestChainAllowListEnabled     []bool
	DestChainAddAllowedSenders    [][]string
	DestChainRemoveAllowedSenders [][]string
}

var ApplyAllowListUpdatesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ApplyAllowListUpdatesInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.ApplyAllowlistUpdates(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.DestChainSelector,
		input.DestChainAllowListEnabled,
		input.DestChainAddAllowedSenders,
		input.DestChainRemoveAllowedSenders,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute  ApplyAllowListUpdates on onRamp: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
	}, err
}

type IsChainSupportedInput struct {
	OnRampPackageId   string
	StateObjectId     string
	DestChainSelector uint64
}

type IsChainSupportedOutput struct {
	IsChainSupported bool
}

var IsChainSupportedHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IsChainSupportedInput) (output sui_ops.OpTxResult[IsChainSupportedOutput], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	isSupported, err := onRampPackage.DevInspect().IsChainSupported(b.GetContext(), opts, bind.Object{Id: input.StateObjectId}, input.DestChainSelector)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects: IsChainSupportedOutput{
			IsChainSupported: isSupported,
		},
	}, nil
}

// Note: Shares the same input as IsChainSupported
// TODO: maybe rename the input to make it more generic
var GetDestChainConfigHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IsChainSupportedInput) (output sui_ops.OpTxResult[IsChainSupportedOutput], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	config, err := onRampPackage.DevInspect().GetDestChainConfig(b.GetContext(), opts, bind.Object{Id: input.StateObjectId}, input.DestChainSelector)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	// The first return value is isEnabled (bool)
	isEnabled, ok := config[0].(bool)
	if !ok {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to parse isEnabled from config")
	}

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects: IsChainSupportedOutput{
			IsChainSupported: isEnabled,
		},
	}, nil
}

type GetFeeInput struct {
	OnRampPackageId   string
	TypeArgs          string
	CCIPObjectRef     string
	DestChainSelector uint64
	Receiver          []byte
	Data              []byte
	TokenAddress      []string
	TokenAmounts      []uint64
	FeeToken          string
	ExtraArgs         []byte
}

var GetFee = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetFeeInput) (output sui_ops.OpTxResult[IsChainSupportedOutput], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	inspectResp, err := onRampPackage.GetFee(b.GetContext(), opts, []string{input.TypeArgs}, bind.Object{Id: input.CCIPObjectRef}, bind.Object{Id: "0x6"}, input.DestChainSelector, input.Receiver, input.Data, input.TokenAddress, input.TokenAmounts, bind.Object{Id: input.FeeToken}, input.ExtraArgs)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	b.Logger.Infow("getFee returned fee", "fee", inspectResp.Results[0])

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects:   IsChainSupportedOutput{},
	}, err
}

var ApplyAllowListUpdateOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-apply-allow-list-updates", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Runs ApplyAllowListUpdates on OnRamp",
	ApplyAllowListUpdatesHandler,
)

var DeployCCIPOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-on-ramp", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP onRamp package",
	deployHandler,
)

var OnRampInitializeOP = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-on-ramp", "package", "initialize"),
	semver.MustParse("0.1.0"),
	"Initialize the CCIP onRamp package",
	InitializeHandler,
)

var ApplyDestChainConfigUpdateOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-apply-dest-chain-update", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Runs ApplyDestChainConfig update on OnRamp",
	ApplyDestChainUpdateHandler,
)

var IsChainSupportedOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-is-chain-supported", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Runs IsChainSupported OnRamp",
	IsChainSupportedHandler,
)

var GetDestChainConfigOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-get-dest-chain-config", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Runs GetDestChainConfig OnRamp",
	GetDestChainConfigHandler,
)
