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

type GetPackageIdsOffRampInput struct {
	OffRampPackageId string
	StateObjectId    string
}

type GetPackageIdsOffRampOutput struct {
	PackageIds []string
}

var getPackageIdsOffRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetPackageIdsOffRampInput) (output sui_ops.OpTxResult[GetPackageIdsOffRampOutput], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetPackageIdsOffRampOutput]{}, err
	}

	opts := deps.GetCallOpts()
	packageIds, err := offRampPackage.DevInspect().GetPackageIds(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetPackageIdsOffRampOutput]{}, fmt.Errorf("failed to get package IDs from OffRamp: %w", err)
	}

	b.Logger.Infow("Package IDs retrieved from OffRamp", "packageIds", packageIds)

	return sui_ops.OpTxResult[GetPackageIdsOffRampOutput]{
		Digest:    "",
		PackageId: input.OffRampPackageId,
		Objects: GetPackageIdsOffRampOutput{
			PackageIds: packageIds,
		},
	}, nil
}

var GetPackageIdsOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-get-package-ids", "package", "query"),
	semver.MustParse("0.1.0"),
	"Gets all package IDs from the OffRamp state",
	getPackageIdsOffRampHandler,
)

type GetInitialPackageIdOffRampInput struct {
	OffRampPackageId string
	StateObjectId    string
}

type GetInitialPackageIdOffRampOutput struct {
	InitialPackageId string
}

var getInitialPackageIdOffRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetInitialPackageIdOffRampInput) (output sui_ops.OpTxResult[GetInitialPackageIdOffRampOutput], err error) {
	offRampPackage, err := module_offramp.NewOfframp(input.OffRampPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetInitialPackageIdOffRampOutput]{}, err
	}

	opts := deps.GetCallOpts()
	initialPackageId, err := offRampPackage.DevInspect().GetInitialPackageId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetInitialPackageIdOffRampOutput]{}, fmt.Errorf("failed to get initial package ID from OffRamp: %w", err)
	}

	b.Logger.Infow("Initial package ID retrieved from OffRamp", "initialPackageId", initialPackageId)

	return sui_ops.OpTxResult[GetInitialPackageIdOffRampOutput]{
		Digest:    "",
		PackageId: input.OffRampPackageId,
		Objects: GetInitialPackageIdOffRampOutput{
			InitialPackageId: initialPackageId,
		},
	}, nil
}

var GetInitialPackageIdOffRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-offramp-get-initial-package-id", "package", "query"),
	semver.MustParse("0.1.0"),
	"Gets the initial package ID from the OffRamp state",
	getInitialPackageIdOffRampHandler,
)
