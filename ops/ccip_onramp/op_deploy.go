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
	onRampPackage, tx, err := onramp.PublishOnramp(
		b.GetContext(),
		deps.GetTxOpts(),
		deps.Signer,
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
		Digest:    tx.Digest.String(),
		PackageId: onRampPackage.Address().String(),
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

	call := onRampPackage.Initialize(
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

	tx, err := call.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute onRamp initialization: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest.String(),
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

	call := onRampPackage.ApplyDestChainConfigUpdates(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.DestChainSelector,
		input.DestChainEnabled,
		input.DestChainAllowListEnabled,
	)

	tx, err := call.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest.String(),
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

	call := onRampPackage.IsChainSupported(bind.Object{Id: input.StateObjectId}, input.DestChainSelector)

	inspectResp, err := call.Inspect(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	IsChainSupported := inspectResp.Results[0].ReturnValues[0].([]any)[0].([]any)[0].(float64)

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects: IsChainSupportedOutput{
			IsChainSupported: IsChainSupported != 0,
		},
	}, err
}

// Note: Shares the same input as IsChainSupported
// TODO: maybe rename the input to make it more generic
var GetDestChainConfigHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IsChainSupportedInput) (output sui_ops.OpTxResult[IsChainSupportedOutput], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, err
	}

	call := onRampPackage.GetDestChainConfig(bind.Object{Id: input.StateObjectId}, input.DestChainSelector)

	inspectResp, err := call.Inspect(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	isEnabled := inspectResp.Results[0].ReturnValues[0].([]any)[0].([]any)[0].(float64)

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects: IsChainSupportedOutput{
			IsChainSupported: isEnabled != 0,
		},
	}, err
}

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
