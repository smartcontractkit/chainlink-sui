package offrampops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_ownable "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/ownable"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/offramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployCCIPOffRampObjects struct {
	// State Object
	OwnerCapObjectId         string
	CCIPOffRampStateObjectId string
	UpgradeCapObjectId       string
}

type DeployCCIPOffRampInput struct {
	CCIPPackageId string
	MCMSPackageId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployCCIPOffRampInput) (output sui_ops.OpTxResult[DeployCCIPOffRampObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	offRampPackage, tx, err := offramp.PublishOfframp(
		b.GetContext(),
		opts,
		deps.Client,
		input.CCIPPackageId,
		input.MCMSPackageId,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "offramp", "OffRampState")
	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")

	if err1 != nil || err2 != nil || err3 != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: err1=%w, err2=%w", err1, err2)
	}

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: offRampPackage.Address(),
		Objects: DeployCCIPOffRampObjects{
			OwnerCapObjectId:         obj1,
			CCIPOffRampStateObjectId: obj2,
			UpgradeCapObjectId:       obj3,
		},
	}, err
}

type InitializeOffRampInput struct {
	OffRampPackageId                      string
	OffRampStateId                        string
	OwnerCapObjectId                      string
	FeeQuoterCapId                        string
	DestTransferCapId                     string
	ChainSelector                         uint64
	PremissionExecThresholdSeconds        uint32
	SourceChainSelectors                  []uint64
	SourceChainsIsEnabled                 []bool
	SourceChainsIsRMNVerificationDisabled []bool
	SourceChainsOnRamp                    [][]byte
}

var initializeHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitializeOffRampInput) (output sui_ops.OpTxResult[DeployCCIPOffRampObjects], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.Initialize(
		b.GetContext(),
		opts,
		bind.Object{Id: input.OffRampStateId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.FeeQuoterCapId},
		bind.Object{Id: input.DestTransferCapId},
		input.ChainSelector,
		input.PremissionExecThresholdSeconds,
		input.SourceChainSelectors,
		input.SourceChainsIsEnabled,
		input.SourceChainsIsRMNVerificationDisabled,
		input.SourceChainsOnRamp,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to execute Offramp initialization: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   DeployCCIPOffRampObjects{},
	}, err
}

type SetOCR3ConfigInput struct {
	OffRampPackageId               string
	OffRampStateId                 string
	CCIPObjectRefId                string
	OwnerCapObjectId               string
	ConfigDigest                   []byte
	OCRPluginType                  byte
	BigF                           byte
	IsSignatureVerificationEnabled bool
	Signers                        [][]byte
	Transmitters                   []string
}

var setOCR3ConfigHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input SetOCR3ConfigInput) (output sui_ops.OpTxResult[DeployCCIPOffRampObjects], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, err
	}

	encodedCall, err := offRampPackage.Encoder().SetOcr3Config(
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.OffRampStateId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ConfigDigest,
		input.OCRPluginType,
		input.BigF,
		input.IsSignatureVerificationEnabled,
		input.Signers,
		input.Transmitters,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to encode SetOcr3Config call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.OffRampStateId)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetOcr3Config on OffRamp as per no Signer provided")
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
			Digest:    "",
			PackageId: input.OffRampPackageId,
			Objects:   DeployCCIPOffRampObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to execute SetOcr3Config on OffRamp: %w", err)
	}

	b.Logger.Infow("OCR3 config set on OffRamp")

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   DeployCCIPOffRampObjects{},
		Call:      call,
	}, nil
}

type ApplySourceChainConfigUpdateInput struct {
	CCIPObjectRef                         string
	OffRampPackageId                      string
	OffRampStateId                        string
	OwnerCapObjectId                      string
	SourceChainsSelectors                 []uint64
	SourceChainsIsEnabled                 []bool
	SourceChainsIsRMNVerificationDisabled []bool
	SourceChainsOnRamp                    [][]byte
}

var applySourceChainConfigUpdateHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ApplySourceChainConfigUpdateInput) (output sui_ops.OpTxResult[DeployCCIPOffRampObjects], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, err
	}

	encodedCall, err := offRampPackage.Encoder().ApplySourceChainConfigUpdates(
		bind.Object{Id: input.CCIPObjectRef},
		bind.Object{Id: input.OffRampStateId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.SourceChainsSelectors,
		input.SourceChainsIsEnabled,
		input.SourceChainsIsRMNVerificationDisabled,
		input.SourceChainsOnRamp,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to encode ApplySourceChainConfigUpdates call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.OffRampStateId)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ApplySourceChainConfigUpdates on OffRamp as per no Signer provided")
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
			Digest:    "",
			PackageId: input.OffRampPackageId,
			Objects:   DeployCCIPOffRampObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to execute ApplySourceChainConfigUpdates on OffRamp: %w", err)
	}

	b.Logger.Infow("Source chain config updates applied on OffRamp")

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   DeployCCIPOffRampObjects{},
		Call:      call,
	}, nil
}

type AddPackageIdOffRampInput struct {
	OffRampPackageId string
	StateObjectId    string
	OwnerCapObjectId string
	PackageId        string
}

type AddPackageIdOffRampObjects struct {
	// No specific objects are returned from add_package_id
}

var addPackageIdOffRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AddPackageIdOffRampInput) (output sui_ops.OpTxResult[AddPackageIdOffRampObjects], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdOffRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.AddPackageId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.PackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdOffRampObjects]{}, fmt.Errorf("failed to execute AddPackageId on offRamp: %w", err)
	}

	b.Logger.Infow("Package ID added to OffRamp", "packageId", input.PackageId)

	return sui_ops.OpTxResult[AddPackageIdOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   AddPackageIdOffRampObjects{},
	}, nil
}

type RemovePackageIdOffRampInput struct {
	OffRampPackageId string
	StateObjectId    string
	OwnerCapObjectId string
	PackageId        string
}

type RemovePackageIdOffRampObjects struct {
	// No specific objects are returned from remove_package_id
}

var removePackageIdOffRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RemovePackageIdOffRampInput) (output sui_ops.OpTxResult[RemovePackageIdOffRampObjects], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdOffRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.RemovePackageId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.PackageId,
	)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdOffRampObjects]{}, fmt.Errorf("failed to execute RemovePackageId on offRamp: %w", err)
	}

	b.Logger.Infow("Package ID removed from OffRamp", "packageId", input.PackageId)

	return sui_ops.OpTxResult[RemovePackageIdOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   RemovePackageIdOffRampObjects{},
	}, nil
}

var DeployCCIPOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-off-ramp", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP offramp package",
	deployHandler,
)

var InitializeOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-off-ramp", "package", "initialize"),
	semver.MustParse("0.1.0"),
	"Initialize the CCIP offramp package",
	initializeHandler,
)

var SetOCR3ConfigOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-off-ramp", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Running CCIP setOCR3Config package",
	setOCR3ConfigHandler,
)

var ApplySourceChainConfigUpdatesOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-off-ramp", "package", "applysourcechainconfigupdate"),
	semver.MustParse("0.1.0"),
	"Running Offramp ApplySourceChainConfigUpdate operation",
	applySourceChainConfigUpdateHandler,
)

var AddPackageIdOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-add-package-id", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Adds a new package ID to the OffRamp state for upgrade tracking",
	addPackageIdOffRampHandler,
)

var RemovePackageIdOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-remove-package-id", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Removes a package ID from the OffRamp state for upgrade tracking",
	removePackageIdOffRampHandler,
)

type TransferOwnershipOffRampInput struct {
	OffRampPackageId     string
	CCIPObjectRefId      string
	OffRampStateObjectId string
	OwnerCapObjectId     string
	To                   string
}

type TransferOwnershipOffRampObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipOffRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipOffRampInput) (output sui_ops.OpTxResult[TransferOwnershipOffRampObjects], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipOffRampObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.TransferOwnership(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefId},
		bind.Object{Id: input.OffRampStateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipOffRampObjects]{}, fmt.Errorf("failed to execute TransferOwnership on OffRamp: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for OffRamp", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   TransferOwnershipOffRampObjects{},
	}, nil
}

var TransferOwnershipOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-transfer-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the OffRamp",
	transferOwnershipOffRampHandler,
)

type AcceptOwnershipOffRampInput struct {
	OffRampPackageId     string
	OffRampRefObjectId   string
	OffRampStateObjectId string
}

type AcceptOwnershipOffRampObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipOffRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipOffRampInput) (output sui_ops.OpTxResult[AcceptOwnershipOffRampObjects], err error) {
	// mcms_accept_ownership uses the Ownable module
	ownablePackage, err := module_ownable.NewOwnable(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipOffRampObjects]{}, err
	}

	// Use AcceptOwnershipWithArgs since OwnableState no longer has key ability
	// We pass the full OffRampState object which contains OwnableState as a field
	encodedCall, err := ownablePackage.Encoder().AcceptOwnershipWithArgs(bind.Object{Id: input.OffRampStateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipOffRampObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	// We need ownable encoding, but we call the offramp
	encodedCall.Module.ModuleName = "offramp"
	// we set the ref object as the state of the call, so we can access it in the entrypoint encoder
	call, err := sui_ops.ToTransactionCall(encodedCall, input.OffRampRefObjectId)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipOffRampObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on OffRamp as per no Signer provided")
		return sui_ops.OpTxResult[AcceptOwnershipOffRampObjects]{
			Digest:    "",
			PackageId: input.OffRampPackageId,
			Objects:   AcceptOwnershipOffRampObjects{},
			Call:      call,
		}, nil
	}

	// If we have a signer, this is accepting the ownership from EOA, we use the offramp directly
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipOffRampObjects]{}, err
	}
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.AcceptOwnership(
		b.GetContext(),
		opts,
		bind.Object{Id: input.OffRampRefObjectId},
		bind.Object{Id: input.OffRampStateObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipOffRampObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on OffRamp: %w", err)
	}

	b.Logger.Infow("Ownership accepted for OffRamp")

	return sui_ops.OpTxResult[AcceptOwnershipOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   AcceptOwnershipOffRampObjects{},
	}, nil
}

var AcceptOwnershipOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-accept-ownership", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the OffRamp",
	acceptOwnershipOffRampHandler,
)
