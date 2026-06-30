package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
)

type InitTARObjects struct {
	TARStateObjectId string
}

type InitTARInput struct {
	CCIPPackageId      string
	StateObjectId      string
	OwnerCapObjectId   string
	LocalChainSelector uint64
}

var initTarHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitTARInput) (output sui_ops.OpTxResult[InitTARObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitTARObjects]{}, fmt.Errorf("failed to create fee quoter contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Initialize(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[InitTARObjects]{}, fmt.Errorf("failed to execute fee quoter initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "token_admin_registry", "TokenAdminRegistryState")
	if err1 != nil {
		return sui_ops.OpTxResult[InitTARObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[InitTARObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects: InitTARObjects{
			TARStateObjectId: obj1,
		},
	}, err
}

var TokenAdminRegistryInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP Token Admin Registry contract",
	initTarHandler,
)

type InitLocalDecimalsInput struct {
	CCIPPackageId    string // original CCIP package = MCMS on-chain identity (proposal target)
	LatestPackageId  string // upgraded CCIP package = PTB dispatch / direct-exec binary (empty when not upgraded)
	StateObjectId    string
	OwnerCapObjectId string
}

var initLocalDecimalsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitLocalDecimalsInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	// Encode/dispatch against the upgraded package; the on-chain MCMS identity stays the original.
	binaryPkgId := input.CCIPPackageId
	if input.LatestPackageId != "" {
		binaryPkgId = input.LatestPackageId
	}

	// MCMS callback validate_obj_addrs expects ref then owner_cap (see mcms_initialize_local_decimals).
	data, err := SerializeMcmsObjectAddrs(input.StateObjectId, input.OwnerCapObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to serialize initialize_local_decimals MCMS data: %w", err)
	}
	call := sui_ops.TransactionCall{
		PackageID:  input.CCIPPackageId, // original = MCMS on-chain identity (validated against with_original_ids)
		Module:     "token_admin_registry",
		Function:   "initialize_local_decimals",
		Data:       data,
		StateObjID: input.StateObjectId,
		TypeArgs:   []string{},
	}
	if input.LatestPackageId != "" {
		call.LatestPackageID = input.LatestPackageId // upgraded = PTB MoveCall dispatch target
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of InitializeLocalDecimals on TokenAdminRegistry as per no Signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	contract, err := module_token_admin_registry.NewTokenAdminRegistry(binaryPkgId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	encodedCall, err := contract.Encoder().InitializeLocalDecimals(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode InitializeLocalDecimals call: %w", err)
	}
	call, err = sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if input.LatestPackageId != "" {
		call.LatestPackageID = call.PackageID // latest (binaryPkgId) → PTB dispatch
		call.PackageID = input.CCIPPackageId  // original → on-chain MCMS identity
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.InitializeLocalDecimals(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute local decimals initialization: %w", err)
	}

	b.Logger.Infow("InitializeLocalDecimals on TokenAdminRegistry", "packageId", input.CCIPPackageId)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var TokenAdminRegistryInitializeLocalDecimalsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "initialize_local_decimals"),
	semver.MustParse("0.1.0"),
	"Initializes local token decimals state on the CCIP Token Admin Registry",
	initLocalDecimalsHandler,
)

type BackfillLocalDecimalsInput struct {
	CCIPPackageId       string // original CCIP package = MCMS on-chain identity (proposal target)
	LatestPackageId     string // upgraded CCIP package = PTB dispatch / direct-exec binary (empty when not upgraded)
	StateObjectId       string
	OwnerCapObjectId    string
	CoinMetadataAddress string
	TokenType           string
	LocalDecimals       *byte
}

var backfillLocalDecimalsHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input BackfillLocalDecimalsInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	binaryPkgId := input.CCIPPackageId
	if input.LatestPackageId != "" {
		binaryPkgId = input.LatestPackageId
	}
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(binaryPkgId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	localDecimals, err := ResolveLocalDecimals(b.GetContext(), deps.Client, input.TokenType, input.LocalDecimals)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, err
	}

	encodedCall, err := contract.Encoder().BackfillLocalDecimals(
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.StateObjectId},
		input.CoinMetadataAddress,
		localDecimals,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode BackfillLocalDecimals call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if input.LatestPackageId != "" {
		call.LatestPackageID = call.PackageID // latest (binaryPkgId) → PTB dispatch
		call.PackageID = input.CCIPPackageId  // original → on-chain MCMS identity
	}

	if deps.Signer == nil {
		b.Logger.Infow(
			"Skipping execution of BackfillLocalDecimals on TokenAdminRegistry as per no Signer provided",
			"coinMetadataAddress", input.CoinMetadataAddress,
			"localDecimals", localDecimals,
		)
		return sui_ops.OpTxResult[NoObjects]{
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(
		b.GetContext(),
		opts,
		encodedCall,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute backfill local decimals: %w", err)
	}

	b.Logger.Infow(
		"BackfillLocalDecimals on TokenAdminRegistry",
		"packageId", input.CCIPPackageId,
		"coinMetadataAddress", input.CoinMetadataAddress,
		"localDecimals", localDecimals,
	)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var TokenAdminRegistryBackfillLocalDecimalsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "backfill_local_decimals"),
	semver.MustParse("0.1.0"),
	"Backfills local token decimals for an already-registered pool",
	backfillLocalDecimalsHandler,
)

// ================================================================
// |                    Unregister Pool                          |
// ================================================================

type UnregisterPoolInput struct {
	CCIPPackageId       string
	CCIPObjectRef       string
	OwnerCapObjectId    string
	CoinMetadataAddress string
}

var unregisterPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UnregisterPoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	data, err := rmn.SerializeMcmsObjectAddrs(
		input.OwnerCapObjectId,
		input.CCIPObjectRef,
		input.CoinMetadataAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode unregister_pool MCMS callback data: %w", err)
	}
	call := sui_ops.TransactionCall{
		PackageID:  input.CCIPPackageId,
		Module:     "token_admin_registry",
		Function:   "unregister_pool",
		Data:       data,
		StateObjID: input.CCIPObjectRef,
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of UnregisterPool on TokenAdminRegistry as per no Signer provided", "CoinMetadataAddress", input.CoinMetadataAddress)
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UnregisterPool(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		input.CoinMetadataAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute UnregisterPool on TokenAdminRegistry: %w", err)
	}

	b.Logger.Infow("UnregisterPool on TokenAdminRegistry", "PackageId:", input.CCIPPackageId, "CoinMetadataAddress:", input.CoinMetadataAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var TokenAdminRegistryUnregisterPoolOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "unregister_pool"),
	semver.MustParse("0.1.0"),
	"Unregisters a token pool from the CCIP Token Admin Registry",
	unregisterPoolHandler,
)

// ================================================================
// |                  Transfer Admin Role                        |
// ================================================================

type TransferAdminRoleInput struct {
	CCIPPackageId       string
	CCIPObjectRef       string
	CoinMetadataAddress string
	NewAdmin            string
}

var transferAdminRoleHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferAdminRoleInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	encodedCall, err := contract.Encoder().TransferAdminRole(bind.Object{Id: input.CCIPObjectRef}, input.CoinMetadataAddress, input.NewAdmin)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode TransferAdminRole call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRef)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of TransferAdminRole on TokenAdminRegistry as per no Signer provided", "CoinMetadataAddress", input.CoinMetadataAddress, "NewAdmin", input.NewAdmin)
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.TransferAdminRole(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		input.CoinMetadataAddress,
		input.NewAdmin,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute TransferAdminRole on TokenAdminRegistry: %w", err)
	}

	b.Logger.Infow("TransferAdminRole on TokenAdminRegistry", "PackageId:", input.CCIPPackageId, "CoinMetadataAddress:", input.CoinMetadataAddress, "NewAdmin:", input.NewAdmin)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var TokenAdminRegistryTransferAdminRoleOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "transfer_admin_role"),
	semver.MustParse("0.1.0"),
	"Transfers admin role for a token in the CCIP Token Admin Registry",
	transferAdminRoleHandler,
)

// ================================================================
// |                   Accept Admin Role                         |
// ================================================================

type AcceptAdminRoleInput struct {
	CCIPPackageId       string
	CCIPObjectRef       string
	CoinMetadataAddress string
}

var acceptAdminRoleHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptAdminRoleInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	encodedCall, err := contract.Encoder().AcceptAdminRole(bind.Object{Id: input.CCIPObjectRef}, input.CoinMetadataAddress)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode AcceptAdminRole call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRef)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptAdminRole on TokenAdminRegistry as per no Signer provided", "CoinMetadataAddress", input.CoinMetadataAddress)
		return sui_ops.OpTxResult[NoObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   NoObjects{},
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.AcceptAdminRole(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRef},
		input.CoinMetadataAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute AcceptAdminRole on TokenAdminRegistry: %w", err)
	}

	b.Logger.Infow("AcceptAdminRole on TokenAdminRegistry", "PackageId:", input.CCIPPackageId, "CoinMetadataAddress:", input.CoinMetadataAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
		Call:      call,
	}, nil
}

var TokenAdminRegistryAcceptAdminRoleOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "accept_admin_role"),
	semver.MustParse("0.1.0"),
	"Accepts admin role for a token in the CCIP Token Admin Registry",
	acceptAdminRoleHandler,
)
