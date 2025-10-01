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
	"github.com/block-vision/sui-go-sdk/mystenbcs"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	"github.com/smartcontractkit/chainlink-sui/bindings/packages/router"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

type DeployCCIPRouterInput struct {
	McmsPackageId string
	McmsOwner     string
}
type DeployCCIPRouterObjects struct {
	RouterObjectId             string
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

	routerObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "router", "RouterObject")
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to find RouterObject ID in publish tx: %w", err)
	}

	ownerCapId, err := deriveRouterOwnerCapId(routerObjectId, routerPackage.Address())
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to derive OwnerCap ID: %w", err)
	}

	routerStateId, err := deriveRouterStateId(routerObjectId, routerPackage.Address())
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to derive RouterState ID: %w", err)
	}

	b.Logger.Infow("Router objects calculated deterministically",
		"routerObjectId", routerObjectId,
		"ownerCapId", ownerCapId,
		"routerStateId", routerStateId,
	)

	routerStatePointerId, err := bind.FindObjectIdFromPublishTx(*tx, "router", "RouterStatePointer")
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to find RouterStatePointer ID in publish tx: %w", err)
	}

	routerStatePointerResp, err := bind.ReadObject(b.GetContext(), routerStatePointerId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to read RouterStatePointer object: %w", err)
	}

	// Decode the RouterStatePointer struct from the response
	var routerStatePointer module_router.RouterStatePointer
	if routerStatePointerResp.Data == nil || routerStatePointerResp.Data.Content == nil ||
		routerStatePointerResp.Data.Content.SuiMoveObject.Fields == nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("RouterStatePointer object has no content")
	}

	fields := routerStatePointerResp.Data.Content.SuiMoveObject.Fields
	if routerObjectIdField, ok := fields["router_object_id"].(string); ok {
		routerStatePointer.RouterObjectId = routerObjectIdField
	} else {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to decode router_object_id from RouterStatePointer")
	}

	// Validate that the RouterObjectId in RouterStatePointer matches what we found in the tx
	if routerStatePointer.RouterObjectId != routerObjectId {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf(
			"RouterObjectId mismatch: found %s in tx, but RouterStatePointer contains %s",
			routerObjectId,
			routerStatePointer.RouterObjectId,
		)
	}

	b.Logger.Infow("RouterStatePointer validated",
		"routerStatePointerId", routerStatePointerId,
		"storedRouterObjectId", routerStatePointer.RouterObjectId,
		"derivedOwnerCapId", ownerCapId,
		"derivedRouterStateId", routerStateId,
	)

	return sui_ops.OpTxResult[DeployCCIPRouterObjects]{
		Digest:    tx.Digest,
		PackageId: routerPackage.Address(),
		Objects: DeployCCIPRouterObjects{
			RouterObjectId:             routerObjectId,
			OwnerCapObjectId:           ownerCapId,
			RouterStateObjectId:        routerStateId,
			RouterStatePointerObjectId: routerStatePointerId,
		},
	}, nil
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

type RouterOwnableKey bool
type RouterKey bool

func deriveRouterOwnerCapId(routerObjectId string, routerPackageId string) (string, error) {
	key := RouterOwnableKey(false)
	keyBytes, err := mystenbcs.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("failed to BCS serialize key: %w", err)
	}

	return codec.DeriveDerivedObjectID(routerObjectId, routerPackageId, "router", "RouterOwnableKey", keyBytes)
}

func deriveRouterStateId(routerObjectId string, routerPackageId string) (string, error) {
	key := RouterKey(false)
	keyBytes, err := mystenbcs.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("failed to BCS serialize key: %w", err)
	}

	return codec.DeriveDerivedObjectID(routerObjectId, routerPackageId, "router", "RouterKey", keyBytes)
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
