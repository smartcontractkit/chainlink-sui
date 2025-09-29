// Package routerops provides operations for deploying and managing CCIP Router contracts.
//
// Available operations:
//   - DeployCCIPRouterOp: Deploys the CCIP router package
//   - SetOnRampsOp: Sets on-ramp addresses for destination chains
//
// Example usage:
//
//	// Deploy router
//	reportRouter, err := cld_ops.ExecuteOperation(bundle, DeployCCIPRouterOp, deps, DeployCCIPRouterInput{
//	    McmsPackageId: mcmsPackageId,
//	    McmsOwner:     ownerAddress,
//	})
//
//	// Set on-ramps
//	_, err = cld_ops.ExecuteOperation(bundle, SetOnRampsOp, deps, SetOnRampsInput{
//	    RouterPackageId:     reportRouter.Output.PackageId,
//	    RouterStateObjectId: reportRouter.Output.Objects.RouterStateObjectId,
//	    OwnerCapObjectId:    reportRouter.Output.Objects.OwnerCapObjectId,
//	    DestChainSelectors:  []uint64{5009297550715157269}, // ETH chain selector
//	    OnRampAddresses:     []string{"0x1111111111111111111111111111111111111111"},
//	})
package routerops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/router"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type DeployCCIPRouterInput struct {
	McmsPackageId string
	McmsOwner     string
}
type DeployCCIPRouterObjects struct {
	OwnerCapObjectId           string
	RouterStateObjectId        string
	RouterStatePointerObjectId string
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployCCIPRouterInput) (output sui_ops.OpTxResult[DeployCCIPRouterObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	routerPackage, tx, err := router.PublishCCIPRouter(
		b.GetContext(),
		opts,
		deps.Client,
		input.McmsPackageId,
		input.McmsOwner,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, err
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "router", "RouterState")
	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "router", "RouterStatePointer")
	if err1 != nil || err2 != nil || err3 != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPRouterObjects]{
		Digest:    tx.Digest,
		PackageId: routerPackage.Address(),
		Objects: DeployCCIPRouterObjects{
			OwnerCapObjectId:           obj1,
			RouterStateObjectId:        obj2,
			RouterStatePointerObjectId: obj3,
		},
	}, err
}

type SetOnRampsInput struct {
	RouterPackageId     string
	RouterStateObjectId string
	OwnerCapObjectId    string
	DestChainSelectors  []uint64
	OnRampAddresses     []string
}

type SetOnRampsObjects struct {
	// No specific objects are returned from set_on_ramps
}

var setOnRampsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input SetOnRampsInput) (output sui_ops.OpTxResult[SetOnRampsObjects], err error) {
	routerPackage, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[SetOnRampsObjects]{}, err
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := routerPackage.SetOnRamps(
		b.GetContext(),
		opts,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.RouterStateObjectId},
		input.DestChainSelectors,
		input.OnRampAddresses,
	)
	if err != nil {
		return sui_ops.OpTxResult[SetOnRampsObjects]{}, fmt.Errorf("failed to execute set_on_ramps: %w", err)
	}

	b.Logger.Infow("On-ramps set successfully",
		"destChainSelectors", input.DestChainSelectors,
		"onRampAddresses", input.OnRampAddresses)

	return sui_ops.OpTxResult[SetOnRampsObjects]{
		Digest:    tx.Digest,
		PackageId: input.RouterPackageId,
		Objects:   SetOnRampsObjects{},
	}, nil
}

var DeployCCIPRouterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP router package",
	deployHandler,
)

var SetOnRampsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "set-on-ramps"),
	semver.MustParse("0.1.0"),
	"Sets on-ramp addresses for destination chains in the CCIP router",
	setOnRampsHandler,
)
