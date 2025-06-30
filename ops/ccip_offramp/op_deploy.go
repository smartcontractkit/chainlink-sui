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
	offRampPackage, tx, err := offramp.PublishOfframp(
		b.GetContext(),
		deps.GetTxOpts(),
		deps.Signer,
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
		Digest:    tx.Digest.String(),
		PackageId: offRampPackage.Address().String(),
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

	call := offRampPackage.Initialize(
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

	tx, err := call.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to execute Offramp initialization: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.OffRampPackageId,
		Objects:   DeployCCIPOffRampObjects{},
	}, err
}

type SetOCR3ConfigInput struct {
	OffRampPackageId               string
	OffRampStateId                 string
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

	call := offRampPackage.SetOcr3Config(
		bind.Object{Id: input.OffRampStateId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.ConfigDigest,
		input.OCRPluginType,
		input.BigF,
		input.IsSignatureVerificationEnabled,
		input.Signers,
		input.Transmitters,
	)

	tx, err := call.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{}, fmt.Errorf("failed to execute set ocr3 config in offramp: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPOffRampObjects]{
		Digest:    tx.Digest.String(),
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
