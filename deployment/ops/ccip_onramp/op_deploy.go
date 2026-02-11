package onrampops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/onramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"

	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_ownable "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/ownable"
)

type DeployCCIPOnRampObjects struct {
	// State Object
	OwnerCapObjectId        string
	CCIPOnrampStateObjectId string
	UpgradeCapObjectId      string
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
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "onramp", "OnRampState")
	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")

	if err1 != nil || err2 != nil || err3 != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: onRampPackage.Address(),
		Objects: DeployCCIPOnRampObjects{
			OwnerCapObjectId:        obj1,
			CCIPOnrampStateObjectId: obj2,
			UpgradeCapObjectId:      obj3,
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
	DestChainAllowListEnabled []bool
	DestChainRouters          []string
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
		input.DestChainAllowListEnabled,
		input.DestChainRouters,
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
	CCIPObjectRefId           string
	OwnerCapObjectId          string
	StateObjectId             string
	DestChainSelector         []uint64
	DestChainAllowListEnabled []bool
	DestChainRouters          []string
}

var ApplyDestChainUpdateHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ApplyDestChainConfigureOnRampInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	encodedCall, err := onRampPackage.Encoder().ApplyDestChainConfigUpdates(
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.DestChainSelector,
		input.DestChainAllowListEnabled,
		input.DestChainRouters,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to encode ApplyDestChainConfigUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyDestChainConfigUpdates on OnRamp as per no Signer provided")
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   DeployCCIPOnRampObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute ApplyDestChainConfigUpdates on OnRamp: %w", err)
	}

	b.Logger.Infow("Destination chain config updates applied on OnRamp")

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
		Call:      call,
	}, nil
}

type ApplyAllowListUpdatesInput struct {
	OnRampPackageId               string
	CCIPObjectRefId               string
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

	encodedCall, err := onRampPackage.Encoder().ApplyAllowlistUpdates(
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.DestChainSelector,
		input.DestChainAllowListEnabled,
		input.DestChainAddAllowedSenders,
		input.DestChainRemoveAllowedSenders,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to encode ApplyAllowlistUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplyAllowlistUpdates on OnRamp as per no Signer provided")
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   DeployCCIPOnRampObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute ApplyAllowlistUpdates on OnRamp: %w", err)
	}

	b.Logger.Infow("Allowlist updates applied on OnRamp")

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
		Call:      call,
	}, nil
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

var GetDestChainConfigHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IsChainSupportedInput) (output sui_ops.OpTxResult[IsChainSupportedOutput], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	config, err := onRampPackage.DevInspect().GetDestChainConfig(b.GetContext(), opts, bind.Object{Id: input.StateObjectId}, input.DestChainSelector)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to get dest chain config: %w", err)
	}

	// GetDestChainConfig returns (sequence_number: u64, allowlist_enabled: bool, router: address)
	// The router address being non-zero indicates the destination chain is enabled
	router, ok := config[2].(string)
	if !ok {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to parse router address from config")
	}

	// Chain is supported if router is not zero address
	isSupported := router != "0x0" && router != "0x0000000000000000000000000000000000000000000000000000000000000000"

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects: IsChainSupportedOutput{
			IsChainSupported: isSupported,
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
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to get fee: %w", err)
	}

	b.Logger.Infow("getFee returned fee", "fee", inspectResp.Results[0])

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects:   IsChainSupportedOutput{},
	}, err
}

type SetDynamicConfigInput struct {
	OnRampPackageId  string
	CCIPObjectRefId  string
	StateObjectId    string
	OwnerCapObjectId string
	FeeAggregator    string
	AllowListAdmin   string
}

var SetDynamicConfigHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input SetDynamicConfigInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	encodedCall, err := onRampPackage.Encoder().SetDynamicConfig(
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.FeeAggregator,
		input.AllowListAdmin,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to encode SetDynamicConfig call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetDynamicConfig on OnRamp as per no Signer provided")
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   DeployCCIPOnRampObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute SetDynamicConfig on OnRamp: %w", err)
	}

	b.Logger.Infow("Dynamic config set on OnRamp")

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
		Call:      call,
	}, nil
}

var ApplyAllowListUpdateOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-apply-allow-list-updates", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Runs ApplyAllowListUpdates on OnRamp",
	ApplyAllowListUpdatesHandler,
)

type WithdrawFeeTokensInput struct {
	OnRampPackageId    string
	CCIPObjectRefId    string
	StateObjectId      string
	OwnerCapObjectId   string
	FeeTokenMetadataId string
	TypeArg            string
}

var WithdrawFeeTokensHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input WithdrawFeeTokensInput) (output sui_ops.OpTxResult[DeployCCIPOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, err
	}

	encodedCall, err := onRampPackage.Encoder().WithdrawFeeTokens(
		[]string{input.TypeArg},
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.FeeTokenMetadataId},
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to encode WithdrawFeeTokens call: %w", err)
	}

	call, err := sui_ops.ToTransactionCallWithTypeArgs(encodedCall, input.StateObjectId, []string{input.TypeArg})
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of WithdrawFeeTokens on OnRamp as per no Signer provided")
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   DeployCCIPOnRampObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{}, fmt.Errorf("failed to execute WithdrawFeeTokens on OnRamp: %w", err)
	}

	b.Logger.Infow("Fee tokens withdrawn on OnRamp")

	return sui_ops.OpTxResult[DeployCCIPOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   DeployCCIPOnRampObjects{},
		Call:      call,
	}, nil
}

var WithdrawFeeTokensOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-withdraw-fee-tokens", "package", "withdraw"),
	semver.MustParse("0.1.0"),
	"Withdraws fee tokens from the OnRamp",
	WithdrawFeeTokensHandler,
)

var SetDynamicConfigOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-set-dynamic-config", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Runs set_dynamic_config on OnRamp",
	SetDynamicConfigHandler,
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

type AddPackageIdInput struct {
	OnRampPackageId  string
	StateObjectId    string
	OwnerCapObjectId string
	PackageId        string
}

type AddPackageIdObjects struct {
	// No specific objects are returned from add_package_id
}

var addPackageIdHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AddPackageIdInput) (output sui_ops.OpTxResult[AddPackageIdObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdObjects]{}, err
	}

	encodedCall, err := onRampPackage.Encoder().AddPackageId(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.PackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdObjects]{}, fmt.Errorf("failed to encode AddPackageId call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AddPackageId on OnRamp as per no Signer provided", "packageId", input.PackageId)
		return sui_ops.OpTxResult[AddPackageIdObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   AddPackageIdObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdObjects]{}, fmt.Errorf("failed to execute AddPackageId on OnRamp: %w", err)
	}

	b.Logger.Infow("Package ID added to OnRamp", "packageId", input.PackageId)

	return sui_ops.OpTxResult[AddPackageIdObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   AddPackageIdObjects{},
		Call:      call,
	}, nil
}

var AddPackageIdOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-add-package-id", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Adds a new package ID to the OnRamp state for upgrade tracking",
	addPackageIdHandler,
)

type RemovePackageIdOnRampInput struct {
	OnRampPackageId  string
	StateObjectId    string
	OwnerCapObjectId string
	PackageId        string
}

type RemovePackageIdOnRampObjects struct {
	// No specific objects are returned from remove_package_id
}

var removePackageIdOnRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RemovePackageIdOnRampInput) (output sui_ops.OpTxResult[RemovePackageIdOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdOnRampObjects]{}, err
	}

	encodedCall, err := onRampPackage.Encoder().RemovePackageId(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.PackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdOnRampObjects]{}, fmt.Errorf("failed to encode RemovePackageId call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdOnRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of RemovePackageId on OnRamp as per no Signer provided", "packageId", input.PackageId)
		return sui_ops.OpTxResult[RemovePackageIdOnRampObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   RemovePackageIdOnRampObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdOnRampObjects]{}, fmt.Errorf("failed to execute RemovePackageId on OnRamp: %w", err)
	}

	b.Logger.Infow("Package ID removed from OnRamp", "packageId", input.PackageId)

	return sui_ops.OpTxResult[RemovePackageIdOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   RemovePackageIdOnRampObjects{},
		Call:      call,
	}, nil
}

var RemovePackageIdOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-remove-package-id", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Removes a package ID from the OnRamp state",
	removePackageIdOnRampHandler,
)

type GetOwnerOnRampInput struct {
	OnRampPackageId string
	StateObjectId   string
}

type GetOwnerOnRampOutput struct {
	Owner string
}

var getOwnerOnRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetOwnerOnRampInput) (output sui_ops.OpTxResult[GetOwnerOnRampOutput], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerOnRampOutput]{}, err
	}

	opts := deps.GetCallOpts()
	owner, err := onRampPackage.DevInspect().Owner(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerOnRampOutput]{}, fmt.Errorf("failed to get owner from OnRamp: %w", err)
	}

	b.Logger.Infow("Owner retrieved from OnRamp", "owner", owner)

	return sui_ops.OpTxResult[GetOwnerOnRampOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects: GetOwnerOnRampOutput{
			Owner: owner,
		},
	}, nil
}

var GetOwnerOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-get-owner", "package", "query"),
	semver.MustParse("0.1.0"),
	"Gets the owner from the OnRamp state",
	getOwnerOnRampHandler,
)

type GetPendingTransferOnRampInput struct {
	OnRampPackageId string
	StateObjectId   string
}

type GetPendingTransferOnRampOutput struct {
	HasPendingTransfer      bool
	PendingTransferFrom     *string
	PendingTransferTo       *string
	PendingTransferAccepted *bool
}

var getPendingTransferOnRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetPendingTransferOnRampInput) (output sui_ops.OpTxResult[GetPendingTransferOnRampOutput], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetPendingTransferOnRampOutput]{}, err
	}

	opts := deps.GetCallOpts()

	hasPending, err := onRampPackage.DevInspect().HasPendingTransfer(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetPendingTransferOnRampOutput]{}, fmt.Errorf("failed to check pending transfer: %w", err)
	}

	var pendingFrom *string
	var pendingTo *string
	var pendingAccepted *bool

	if hasPending {
		pendingFrom, err = onRampPackage.DevInspect().PendingTransferFrom(
			b.GetContext(),
			opts,
			bind.Object{Id: input.StateObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferOnRampOutput]{}, fmt.Errorf("failed to get pending transfer from: %w", err)
		}

		pendingTo, err = onRampPackage.DevInspect().PendingTransferTo(
			b.GetContext(),
			opts,
			bind.Object{Id: input.StateObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferOnRampOutput]{}, fmt.Errorf("failed to get pending transfer to: %w", err)
		}

		pendingAccepted, err = onRampPackage.DevInspect().PendingTransferAccepted(
			b.GetContext(),
			opts,
			bind.Object{Id: input.StateObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferOnRampOutput]{}, fmt.Errorf("failed to get pending transfer accepted: %w", err)
		}
	}

	b.Logger.Infow("Pending transfer info retrieved from OnRamp",
		"hasPending", hasPending,
		"pendingFrom", pendingFrom,
		"pendingTo", pendingTo,
		"pendingAccepted", pendingAccepted,
	)

	return sui_ops.OpTxResult[GetPendingTransferOnRampOutput]{
		Digest:    "",
		PackageId: input.OnRampPackageId,
		Objects: GetPendingTransferOnRampOutput{
			HasPendingTransfer:      hasPending,
			PendingTransferFrom:     pendingFrom,
			PendingTransferTo:       pendingTo,
			PendingTransferAccepted: pendingAccepted,
		},
	}, nil
}

var GetPendingTransferOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-get-pending-transfer", "package", "query"),
	semver.MustParse("0.1.0"),
	"Gets pending transfer information from the OnRamp state",
	getPendingTransferOnRampHandler,
)

type TransferOwnershipOnRampInput struct {
	OnRampPackageId  string
	CCIPObjectRefId  string
	StateObjectId    string
	OwnerCapObjectId string
	To               string
}

type TransferOwnershipOnRampObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipOnRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipOnRampInput) (output sui_ops.OpTxResult[TransferOwnershipOnRampObjects], err error) {
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipOnRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.TransferOwnership(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipOnRampObjects]{}, fmt.Errorf("failed to execute TransferOwnership on OnRamp: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for OnRamp", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipOnRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
		Objects:   TransferOwnershipOnRampObjects{},
	}, nil
}

var TransferOwnershipOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-transfer-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the OnRamp",
	transferOwnershipOnRampHandler,
)

type AcceptOwnershipOnRampInput struct {
	OnRampPackageId string
	CCIPObjectRefId string
	StateObjectId   string
}

type NoObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipOnRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipOnRampInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	// mcms_accept_ownership uses the Ownable module
	ownablePackage, err := module_ownable.NewOwnable(input.OnRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, err
	}

	// Use AcceptOwnershipWithArgs since OwnableState no longer has key ability
	// We pass the full OnRampState object which contains OwnableState as a field
	encodedCall, err := ownablePackage.Encoder().AcceptOwnershipWithArgs(bind.Object{Id: input.StateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	// We need ownable encoding, but we call the onramp
	encodedCall.Module.ModuleName = "onramp"
	// we set the ref object as the state of the call, so we can access it in the entrypoint encoder
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRefId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on OnRamp as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.OnRampPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	// If we have a signer, this is accepting the ownership from EOA, we use the onramp directly
	onRampPackage, err := module_onramp.NewOnramp(input.OnRampPackageId, deps.Client)
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := onRampPackage.AcceptOwnership(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.StateObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on StateObject: %w", err)
	}

	b.Logger.Infow("Ownership accepted for CCIP StateObject")

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.OnRampPackageId,
	}, nil
}

var AcceptOwnershipOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-onramp-accept-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the OnRamp",
	acceptOwnershipOnRampHandler,
)
