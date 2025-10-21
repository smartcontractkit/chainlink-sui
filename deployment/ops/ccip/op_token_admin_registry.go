package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
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

// ================================================================
// |                    Unregister Pool                          |
// ================================================================

type UnregisterPoolInput struct {
	CCIPPackageId       string
	StateObjectId       string
	CoinMetadataAddress string
}

var unregisterPoolHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input UnregisterPoolInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.UnregisterPool(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.CoinMetadataAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute unregister pool: %w", err)
	}

	b.Logger.Infow("UnregisterPool on TokenAdminRegistry", "PackageId:", input.CCIPPackageId, "CoinMetadataAddress:", input.CoinMetadataAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
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
	StateObjectId       string
	CoinMetadataAddress string
	NewAdmin            string
}

var transferAdminRoleHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferAdminRoleInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.TransferAdminRole(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.CoinMetadataAddress,
		input.NewAdmin,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute transfer admin role: %w", err)
	}

	b.Logger.Infow("TransferAdminRole on TokenAdminRegistry", "PackageId:", input.CCIPPackageId, "CoinMetadataAddress:", input.CoinMetadataAddress, "NewAdmin:", input.NewAdmin)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
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
	StateObjectId       string
	CoinMetadataAddress string
}

var acceptAdminRoleHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptAdminRoleInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create token admin registry contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.AcceptAdminRole(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		input.CoinMetadataAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute accept admin role: %w", err)
	}

	b.Logger.Infow("AcceptAdminRole on TokenAdminRegistry", "PackageId:", input.CCIPPackageId, "CoinMetadataAddress:", input.CoinMetadataAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
	}, nil
}

var TokenAdminRegistryAcceptAdminRoleOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "token_admin_registry", "accept_admin_role"),
	semver.MustParse("0.1.0"),
	"Accepts admin role for a token in the CCIP Token Admin Registry",
	acceptAdminRoleHandler,
)
