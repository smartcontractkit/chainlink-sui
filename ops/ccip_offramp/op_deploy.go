package offrampops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/offramp"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployCCIPOffRampObjects struct {
	// State Object
	OwnerCapObjectId         string
	CCIPOffRampStateObjectId string
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
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "offramp", "OffRampState")

	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: err1=%w, err2=%w", err1, err2)
	}

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: offRampPackage.Address(),
		Objects: DeployCCIPOffRampObjects{
			OwnerCapObjectId:         obj1,
			CCIPOffRampStateObjectId: obj2,
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

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := offRampPackage.SetOcr3Config(
		b.GetContext(),
		opts,
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
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to execute set ocr3 config in offramp: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest,
		PackageId: input.OffRampPackageId,
		Objects:   DeployCCIPOffRampObjects{},
	}, err
}

var DeployCCIPOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-off-ramp", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP offramp package",
	deployHandler,
)

var InitializeOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-off-ramp", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Initialize the CCIP offramp package",
	initializeHandler,
)

var SetOCR3ConfigOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-off-ramp", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Initialize the CCIP setOCR3Config package",
	setOCR3ConfigHandler,
)

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

var AddPackageIdOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-add-package-id", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Adds a new package ID to the OffRamp state for upgrade tracking",
	addPackageIdOffRampHandler,
)

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

var RemovePackageIdOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-remove-package-id", "package", "configure"),
	semver.MustParse("0.1.0"),
	"Removes a package ID from the OffRamp state for upgrade tracking",
	removePackageIdOffRampHandler,
)
