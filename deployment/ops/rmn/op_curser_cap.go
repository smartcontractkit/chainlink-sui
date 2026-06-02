package rmn

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// CreateCurserCapInput mints a CurserCap using CCIP OwnerCap.
//
// Direct path: EOA holds OwnerCap and signs the PTB (deps.Signer set).
// MCMS path: slow MCMS holds OwnerCap in its Registry; proposal leaf uses inner
// function create_curser_cap, routed to mcms_create_curser_cap by bindings/mcms_encoder.go.
type CreateCurserCapInput struct {
	CCIPPackageId    string
	StateObjectId    string
	OwnerCapObjectId string
}

type CreateCurserCapObjects struct {
	CurserCapObjectId string
}

var CreateCurserCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "create_curser_cap"),
	semver.MustParse("0.1.0"),
	"Mint a CurserCap via CCIP OwnerCap",
	createCurserCapHandler,
)

func createCurserCapHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input CreateCurserCapInput) (sui_ops.OpTxResult[CreateCurserCapObjects], error) {
	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[CreateCurserCapObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	ref := bind.Object{Id: input.StateObjectId}
	ownerCap := bind.Object{Id: input.OwnerCapObjectId}

	encodedCall, err := contract.Encoder().CreateCurserCap(ref, ownerCap)
	if err != nil {
		return sui_ops.OpTxResult[CreateCurserCapObjects]{}, fmt.Errorf("failed to encode create_curser_cap: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[CreateCurserCapObjects]{}, fmt.Errorf("failed to build transaction call: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of create_curser_cap as no signer provided")
		return sui_ops.OpTxResult[CreateCurserCapObjects]{
			PackageId: input.CCIPPackageId,
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(b.GetContext(), opts, encodedCall)
	if err != nil {
		return sui_ops.OpTxResult[CreateCurserCapObjects]{}, fmt.Errorf("failed to execute create_curser_cap: %w", err)
	}

	curserCapObjectID, err := bind.FindObjectIdFromPublishTx(*tx, "rmn_remote", "CurserCap")
	if err != nil {
		return sui_ops.OpTxResult[CreateCurserCapObjects]{}, fmt.Errorf("failed to find CurserCap object ID in tx: %w", err)
	}

	b.Logger.Infow("CurserCap minted", "digest", tx.Digest, "curserCapObjectId", curserCapObjectID)

	return sui_ops.OpTxResult[CreateCurserCapObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Call:      call,
		Objects: CreateCurserCapObjects{
			CurserCapObjectId: curserCapObjectID,
		},
	}, nil
}

// McmsMintAndRegisterCurserCapInput builds a slow-MCMS proposal leaf that atomically
// mints a CurserCap and registers it in the fast MCMS Registry.
//
// MCMS-only: there is no direct Move entrypoint for this flow. Use CreateCurserCapOp
// when an EOA holds OwnerCap and only needs to mint the cap object.
type McmsMintAndRegisterCurserCapInput struct {
	CCIPPackageId        string
	StateObjectId        string
	SlowOwnerCapObjectId string
	FastRegistryObjectId string
}

var McmsMintAndRegisterCurserCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "mcms_mint_and_register_curser_cap"),
	semver.MustParse("0.1.0"),
	"Register a CurserCap in the fast MCMS Registry via slow MCMS",
	mcmsMintAndRegisterCurserCapHandler,
)

func mcmsMintAndRegisterCurserCapHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input McmsMintAndRegisterCurserCapInput) (sui_ops.OpTxResult[NoObjects], error) {
	if deps.Signer != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("mint_and_register_curser_cap must run through slow MCMS proposal execution, not direct signer PTB")
	}

	data, err := SerializeMcmsObjectAddrs(
		input.StateObjectId,
		input.SlowOwnerCapObjectId,
		input.FastRegistryObjectId,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to serialize mint_and_register_curser_cap data: %w", err)
	}

	call := sui_ops.TransactionCall{
		PackageID:  input.CCIPPackageId,
		Module:     "rmn_remote",
		Function:     "mint_and_register_curser_cap",
		Data:         data,
		StateObjID:   input.StateObjectId,
		TypeArgs:     []string{},
	}

	b.Logger.Infow("Encoded mint_and_register_curser_cap MCMS leaf",
		"fastRegistry", input.FastRegistryObjectId,
	)

	return sui_ops.OpTxResult[NoObjects]{
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}

// CurseWithCurserCapInput curses subjects using CurserCap.
//
// Direct path: EOA (or any signer) holds CurserCap and signs the PTB.
// MCMS path: fast MCMS holds CurserCap in its Registry; proposal leaf uses inner
// function curse_multiple_with_curser_cap, routed by bindings/mcms_encoder.go.
type CurseWithCurserCapInput struct {
	CCIPPackageId     string
	StateObjectId     string
	CurserCapObjectId string
	Subjects          [][]byte
}

var CurseWithCurserCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "curse_with_curser_cap"),
	semver.MustParse("0.1.0"),
	"Curse subjects via CurserCap",
	curseWithCurserCapHandler,
)

func curseWithCurserCapHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input CurseWithCurserCapInput) (sui_ops.OpTxResult[NoObjects], error) {
	if len(input.Subjects) == 0 {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("at least one subject is required to curse")
	}

	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	encodedCall, err := contract.Encoder().CurseMultipleWithCurserCap(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.CurserCapObjectId},
		input.Subjects,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode curse_with_curser_cap: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to build transaction call: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of curse_with_curser_cap as no signer provided")
		return sui_ops.OpTxResult[NoObjects]{
			PackageId: input.CCIPPackageId,
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.Bound().ExecuteTransaction(b.GetContext(), opts, encodedCall)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute curse_with_curser_cap: %w", err)
	}

	b.Logger.Infow("Subjects cursed via CurserCap", "digest", tx.Digest, "count", len(input.Subjects))

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}
