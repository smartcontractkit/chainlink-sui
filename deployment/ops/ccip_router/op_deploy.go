// Package routerops provides operations for deploying and managing CCIP Router contracts.
//
// Available operations:
//   - DeployCCIPRouterOp: Deploys the CCIP router package
//   - SetOnRampsOp: Sets on-ramp addresses for destination chains
//   - GetOnRampOp: Gets the on-ramp address for a specific destination chain
//   - IsChainSupportedOp: Checks if a destination chain is supported
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
	OwnerCapObjectId    string
	RouterStateObjectId string
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
	if err1 != nil || err2 != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	return sui_ops.OpTxResult[DeployCCIPRouterObjects]{
		Digest:    tx.Digest,
		PackageId: routerPackage.Address(),
		Objects: DeployCCIPRouterObjects{
			OwnerCapObjectId:    obj1,
			RouterStateObjectId: obj2,
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

type GetOnRampInput struct {
	RouterPackageId     string
	RouterStateObjectId string
	DestChainSelector   uint64
}

type GetOnRampOutput struct {
	OnRampAddress string
}

var getOnRampHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetOnRampInput) (output sui_ops.OpTxResult[GetOnRampOutput], err error) {
	routerPackage, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetOnRampOutput]{}, err
	}

	opts := deps.GetCallOpts()
	onRampAddress, err := routerPackage.DevInspect().GetOnRamp(
		b.GetContext(),
		opts,
		bind.Object{Id: input.RouterStateObjectId},
		input.DestChainSelector,
	)
	if err != nil {
		return sui_ops.OpTxResult[GetOnRampOutput]{}, fmt.Errorf("failed to get on-ramp address: %w", err)
	}

	b.Logger.Infow("Retrieved on-ramp address",
		"destChainSelector", input.DestChainSelector,
		"onRampAddress", onRampAddress)

	return sui_ops.OpTxResult[GetOnRampOutput]{
		Digest:    "",
		PackageId: input.RouterPackageId,
		Objects: GetOnRampOutput{
			OnRampAddress: onRampAddress,
		},
	}, nil
}

type IsChainSupportedInput struct {
	RouterPackageId     string
	RouterStateObjectId string
	DestChainSelector   uint64
}

type IsChainSupportedOutput struct {
	IsSupported bool
}

var isChainSupportedHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input IsChainSupportedInput) (output sui_ops.OpTxResult[IsChainSupportedOutput], err error) {
	routerPackage, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, err
	}

	opts := deps.GetCallOpts()
	isSupported, err := routerPackage.DevInspect().IsChainSupported(
		b.GetContext(),
		opts,
		bind.Object{Id: input.RouterStateObjectId},
		input.DestChainSelector,
	)
	if err != nil {
		return sui_ops.OpTxResult[IsChainSupportedOutput]{}, fmt.Errorf("failed to check if chain is supported: %w", err)
	}

	b.Logger.Infow("Checked chain support",
		"destChainSelector", input.DestChainSelector,
		"isSupported", isSupported)

	return sui_ops.OpTxResult[IsChainSupportedOutput]{
		Digest:    "",
		PackageId: input.RouterPackageId,
		Objects: IsChainSupportedOutput{
			IsSupported: isSupported,
		},
	}, nil
}

var SetOnRampsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "set-on-ramps"),
	semver.MustParse("0.1.0"),
	"Sets on-ramp addresses for destination chains in the CCIP router",
	setOnRampsHandler,
)

var GetOnRampOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "get-on-ramp"),
	semver.MustParse("0.1.0"),
	"Gets the on-ramp address for a destination chain from the CCIP router",
	getOnRampHandler,
)

var IsChainSupportedOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "is-chain-supported"),
	semver.MustParse("0.1.0"),
	"Checks if a destination chain is supported by the CCIP router",
	isChainSupportedHandler,
)
