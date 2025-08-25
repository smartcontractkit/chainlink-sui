package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/ccip"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type DeployCCIPObjects struct {
	// State Object
	OwnerCapObjectID             string
	CCIPObjectRefPointerObjectID string
	CCIPObjectRefObjectID        string
	// onramp_state_helper
	SourceTransferCapObjectID string
	// offramp_state_helper
	DestTransferCapObjectID string
}

type DeployCCIPInput struct {
	McmsPackageID string
	McmsOwner     string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployCCIPInput) (output sui_ops.OpTxResult[DeployCCIPObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	ccipPackage, tx, err := ccip.PublishCCIP(
		b.GetContext(),
		opts,
		deps.Client,
		input.McmsPackageID,
		input.McmsOwner,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPObjects]{}, err
	}

	// TODO: We should move the object ID finding logic into the binding package
	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "state_object", "CCIPObjectRefPointer")
	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "state_object", "CCIPObjectRef")
	obj4, err4 := bind.FindObjectIdFromPublishTx(*tx, "onramp_state_helper", "SourceTransferCap")
	obj5, err5 := bind.FindObjectIdFromPublishTx(*tx, "offramp_state_helper", "DestTransferCap")

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return sui_ops.OpTxResult[DeployCCIPObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	b.Logger.Infow("CCIP package deployed", "packageId", ccipPackage.Address(), "CCIP Object Ref", obj3)

	return sui_ops.OpTxResult[DeployCCIPObjects]{
		Digest:    tx.Digest,
		PackageId: ccipPackage.Address(),
		Objects: DeployCCIPObjects{
			OwnerCapObjectID:             obj1,
			CCIPObjectRefPointerObjectID: obj2,
			CCIPObjectRefObjectID:        obj3,
			SourceTransferCapObjectID:    obj4,
			DestTransferCapObjectID:      obj5,
		},
	}, err
}

var DeployCCIPOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP package",
	deployHandler,
)
