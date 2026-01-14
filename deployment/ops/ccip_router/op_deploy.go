// Package routerops provides operations for deploying and managing CCIP Router contracts.
//
// Available operations:
//   - DeployCCIPRouterOp: Deploys the CCIP router package
//   - SetOnRampsOp: Sets on-ramp addresses for destination chains
//   - AcceptOwnershipOp: Accepts ownership transfer for the router
//   - TransferOwnershipOp: Transfers ownership of the router to a new owner
//   - ExecuteOwnershipTransferOp: Executes ownership transfer for the router
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
//
//	// Transfer ownership
//	_, err = cld_ops.ExecuteOperation(bundle, TransferOwnershipOp, deps, TransferOwnershipInput{
//	    RouterPackageId:     reportRouter.Output.PackageId,
//	    RouterStateObjectId: reportRouter.Output.Objects.RouterStateObjectId,
//	    OwnerCapObjectId:    reportRouter.Output.Objects.OwnerCapObjectId,
//	    NewOwner:            newOwnerAddress,
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
	RouterObjectId             string
	OwnerCapObjectId           string
	RouterStateObjectId        string
	RouterStatePointerObjectId string
	UpgradeCapObjectId         string
}

var DeployCCIPRouterOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP router package",
	deployHandler,
)

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input DeployCCIPRouterInput) (output sui_ops.OpTxResult[DeployCCIPRouterObjects], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	routerPackage, tx, err := router.PublishCCIPRouter(
		b.GetContext(),
		opts,
		deps.Client,
		input.McmsPackageId,
		input.McmsOwner,
		deps.SuiRPC,
	)
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, err
	}

	routerObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "router", "RouterObject")
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to find RouterObject ID in publish tx: %w", err)
	}

	ownerCapId, err := bind.DeriveObjectIDWithVectorU8Key(routerObjectId, []byte("CCIP_OWNABLE"))
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to derive OwnerCap ID: %w", err)
	}

	routerStateId, err := bind.DeriveObjectIDWithVectorU8Key(routerObjectId, []byte("RouterState"))
	if err != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to derive RouterState ID: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	obj2, err2 := bind.FindObjectIdFromPublishTx(*tx, "router", "RouterState")
	obj3, err3 := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")

	if err1 != nil || err2 != nil || err3 != nil {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("failed to find object IDs in publish tx: %w", err)
	}

	// Validate derived IDs match the created IDs
	if ownerCapId != obj1 {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("derived OwnerCap ID mismatch: %s != %s", ownerCapId, obj1)
	}
	if routerStateId != obj2 {
		return sui_ops.OpTxResult[DeployCCIPRouterObjects]{}, fmt.Errorf("derived RouterState ID mismatch: %s != %s", routerStateId, obj2)
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
			UpgradeCapObjectId:         obj3,
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

var SetOnRampsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "set-on-ramps"),
	semver.MustParse("0.1.0"),
	"Sets on-ramp addresses for destination chains in the CCIP router",
	setOnRampsHandler,
)

var setOnRampsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input SetOnRampsInput) (output sui_ops.OpTxResult[SetOnRampsObjects], err error) {
	routerPackage, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[SetOnRampsObjects]{}, err
	}

	encodedCall, err := routerPackage.Encoder().SetOnRamps(
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.RouterStateObjectId},
		input.DestChainSelectors,
		input.OnRampAddresses,
	)
	if err != nil {
		return sui_ops.OpTxResult[SetOnRampsObjects]{}, fmt.Errorf("failed to encode SetOnRamps call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.RouterStateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[SetOnRampsObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of SetOnRamps on Router as per no Signer provided",
			"destChainSelectors", input.DestChainSelectors,
			"onRampAddresses", input.OnRampAddresses)
		return sui_ops.OpTxResult[SetOnRampsObjects]{
			Digest:    "",
			PackageId: input.RouterPackageId,
			Objects:   SetOnRampsObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := routerPackage.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[SetOnRampsObjects]{}, fmt.Errorf("failed to execute SetOnRamps on Router: %w", err)
	}

	b.Logger.Infow("On-ramps set successfully",
		"destChainSelectors", input.DestChainSelectors,
		"onRampAddresses", input.OnRampAddresses)

	return sui_ops.OpTxResult[SetOnRampsObjects]{
		Digest:    tx.Digest,
		PackageId: input.RouterPackageId,
		Objects:   SetOnRampsObjects{},
		Call:      call,
	}, nil
}

// NoObjects is used for operations that don't return any specific objects
type NoObjects struct{}

// ================================================================
// |                   Accept Ownership                          |
// ================================================================

type AcceptOwnershipInput struct {
	RouterPackageId     string
	RouterStateObjectId string
}

var AcceptOwnershipOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "accept-ownership"),
	semver.MustParse("0.1.0"),
	"Accepts ownership transfer for the CCIP router",
	acceptOwnershipHandler,
)

var acceptOwnershipHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	routerContract, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create router contract: %w", err)
	}

	encodedCall, err := routerContract.Encoder().AcceptOwnership(bind.Object{Id: input.RouterStateObjectId})
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.RouterStateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on Router as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.RouterPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := routerContract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on StateObject: %w", err)
	}

	b.Logger.Infow("Ownership accepted for CCIP StateObject")

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.RouterPackageId,
		Call:      call,
	}, nil
}

// ================================================================
// |                  Transfer Ownership                         |
// ================================================================

type TransferOwnershipInput struct {
	RouterPackageId     string
	RouterStateObjectId string
	OwnerCapObjectId    string
	NewOwner            string
}

var TransferOwnershipOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "transfer-ownership"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the CCIP router to a new owner",
	transferOwnershipHandler,
)

var transferOwnershipHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	routerContract, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create router contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := routerContract.TransferOwnership(
		b.GetContext(),
		opts,
		bind.Object{Id: input.RouterStateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.NewOwner,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute transfer ownership: %w", err)
	}

	b.Logger.Infow("TransferOwnership on Router", "PackageId:", input.RouterPackageId, "NewOwner:", input.NewOwner)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.RouterPackageId,
		Objects:   NoObjects{},
	}, nil
}

// ================================================================
// |               Execute Ownership Transfer                     |
// ================================================================

type ExecuteOwnershipTransferInput struct {
	RouterPackageId     string
	RouterStateObjectId string
	OwnerCapObjectId    string
	To                  string
}

var ExecuteOwnershipTransferOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-router", "package", "execute-ownership-transfer"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer for the CCIP router",
	executeOwnershipTransferHandler,
)

var executeOwnershipTransferHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	routerContract, err := module_router.NewRouter(input.RouterPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create router contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := routerContract.ExecuteOwnershipTransfer(
		b.GetContext(),
		opts,
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.RouterStateObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute ownership transfer: %w", err)
	}

	b.Logger.Infow("ExecuteOwnershipTransfer on Router", "PackageId:", input.RouterPackageId, "To:", input.To)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.RouterPackageId,
		Objects:   NoObjects{},
	}, nil
}
