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
	tx, err := contract.CreateCurserCap(b.GetContext(), opts, ref, ownerCap)
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
		Function:   "mint_and_register_curser_cap",
		Data:       data,
		StateObjID: input.StateObjectId,
		TypeArgs:   []string{},
	}

	b.Logger.Infow("Encoded mint_and_register_curser_cap MCMS leaf",
		"fastRegistry", input.FastRegistryObjectId,
	)

	return sui_ops.OpTxResult[NoObjects]{
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}

// McmsCreateCurserCapAndTransferInput mints a CurserCap and transfers it to RecipientAddress.
//
// Direct path: EOA holds OwnerCap and signs the PTB (deps.Signer set).
// MCMS path: proposal leaf uses inner function create_curser_cap_and_transfer, routed by bindings/mcms_encoder.go.
type McmsCreateCurserCapAndTransferInput struct {
	CCIPPackageId        string
	StateObjectId        string
	SlowOwnerCapObjectId string
	RecipientAddress     string
}

var McmsCreateCurserCapAndTransferOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "mcms_create_curser_cap_and_transfer"),
	semver.MustParse("0.1.0"),
	"Mint a CurserCap and transfer it to a recipient via slow MCMS",
	mcmsCreateCurserCapAndTransferHandler,
)

func mcmsCreateCurserCapAndTransferHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input McmsCreateCurserCapAndTransferInput) (sui_ops.OpTxResult[NoObjects], error) {
	if input.RecipientAddress == "" {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("recipient address is required for create_curser_cap_and_transfer")
	}

	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	ref := bind.Object{Id: input.StateObjectId}
	ownerCap := bind.Object{Id: input.SlowOwnerCapObjectId}

	encodedCall, err := contract.Encoder().CreateCurserCapAndTransfer(ref, ownerCap, input.RecipientAddress)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode create_curser_cap_and_transfer: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to build transaction call: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of create_curser_cap_and_transfer as no signer provided",
			"recipient", input.RecipientAddress,
		)
		return sui_ops.OpTxResult[NoObjects]{
			PackageId: input.CCIPPackageId,
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.CreateCurserCapAndTransfer(b.GetContext(), opts, ref, ownerCap, input.RecipientAddress)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute create_curser_cap_and_transfer: %w", err)
	}

	b.Logger.Infow("CurserCap minted and transferred", "digest", tx.Digest, "recipient", input.RecipientAddress)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}

// McmsRegisterCurserCapInput builds a slow-MCMS proposal leaf that registers an
// existing on-chain CurserCap in the fast MCMS Registry. The cap object ID is
// pinned in callback data and validated on-chain at execution.
//
// MCMS-only: use McmsMintAndRegisterCurserCapOp when the cap does not exist yet.
type McmsRegisterCurserCapInput struct {
	CCIPPackageId        string
	StateObjectId        string
	SlowOwnerCapObjectId string
	FastRegistryObjectId string
	CurserCapObjectId    string
}

var McmsRegisterCurserCapOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "mcms_register_curser_cap"),
	semver.MustParse("0.1.0"),
	"Register an existing CurserCap in the fast MCMS Registry via slow MCMS",
	mcmsRegisterCurserCapHandler,
)

func mcmsRegisterCurserCapHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input McmsRegisterCurserCapInput) (sui_ops.OpTxResult[NoObjects], error) {
	if deps.Signer != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("register_curser_cap must run through slow MCMS proposal execution, not direct signer PTB")
	}
	if input.CurserCapObjectId == "" {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("curser cap object id is required for register_curser_cap")
	}

	data, err := SerializeMcmsObjectAddrs(
		input.StateObjectId,
		input.SlowOwnerCapObjectId,
		input.FastRegistryObjectId,
		input.CurserCapObjectId,
	)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to serialize register_curser_cap data: %w", err)
	}

	call := sui_ops.TransactionCall{
		PackageID:  input.CCIPPackageId,
		Module:     "rmn_remote",
		Function:   "register_curser_cap",
		Data:       data,
		StateObjID: input.StateObjectId,
		TypeArgs:   []string{},
	}

	b.Logger.Infow("Encoded register_curser_cap MCMS leaf",
		"fastRegistry", input.FastRegistryObjectId,
		"curserCap", input.CurserCapObjectId,
	)

	return sui_ops.OpTxResult[NoObjects]{
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}

// McmsInitializeAllowedCurserCapsInput initializes the CurserCap allowlist with an optional initial cap ID set.
//
// Direct path: EOA holds OwnerCap and signs the PTB (deps.Signer set).
// MCMS path: proposal leaf uses inner function initialize_allowed_curser_caps, routed by bindings/mcms_encoder.go.
type McmsInitializeAllowedCurserCapsInput struct {
	CCIPPackageId        string
	StateObjectId        string
	SlowOwnerCapObjectId string
	InitialCurserCapIds  []string
}

var McmsInitializeAllowedCurserCapsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "mcms_initialize_allowed_curser_caps"),
	semver.MustParse("0.1.0"),
	"Initialize the CurserCap allowlist via slow MCMS",
	mcmsInitializeAllowedCurserCapsHandler,
)

func mcmsInitializeAllowedCurserCapsHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input McmsInitializeAllowedCurserCapsInput) (sui_ops.OpTxResult[NoObjects], error) {
	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	ref := bind.Object{Id: input.StateObjectId}
	ownerCap := bind.Object{Id: input.SlowOwnerCapObjectId}

	encodedCall, err := contract.Encoder().InitializeAllowedCurserCaps(ref, ownerCap, input.InitialCurserCapIds)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode initialize_allowed_curser_caps: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to build transaction call: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of initialize_allowed_curser_caps as no signer provided",
			"initialCurserCapIds", input.InitialCurserCapIds,
		)
		return sui_ops.OpTxResult[NoObjects]{
			PackageId: input.CCIPPackageId,
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.InitializeAllowedCurserCaps(b.GetContext(), opts, ref, ownerCap, input.InitialCurserCapIds)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute initialize_allowed_curser_caps: %w", err)
	}

	b.Logger.Infow("CurserCap allowlist initialized", "digest", tx.Digest, "initialCurserCapIds", input.InitialCurserCapIds)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}

// McmsRegisterCurserCapIdsInput adds cap object IDs to the on-chain allowlist.
//
// Direct path: EOA holds OwnerCap and signs the PTB (deps.Signer set).
// MCMS path: proposal leaf uses inner function register_curser_cap_ids, routed by bindings/mcms_encoder.go.
type McmsRegisterCurserCapIdsInput struct {
	CCIPPackageId        string
	StateObjectId        string
	SlowOwnerCapObjectId string
	CurserCapObjectIds   []string
}

var McmsRegisterCurserCapIdsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "mcms_register_curser_cap_ids"),
	semver.MustParse("0.1.0"),
	"Register CurserCap object IDs on the allowlist via slow MCMS",
	mcmsRegisterCurserCapIdsHandler,
)

func mcmsRegisterCurserCapIdsHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input McmsRegisterCurserCapIdsInput) (sui_ops.OpTxResult[NoObjects], error) {
	if len(input.CurserCapObjectIds) == 0 {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("at least one curser cap object id is required")
	}

	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	ref := bind.Object{Id: input.StateObjectId}
	ownerCap := bind.Object{Id: input.SlowOwnerCapObjectId}

	encodedCall, err := contract.Encoder().RegisterCurserCapIds(ref, ownerCap, input.CurserCapObjectIds)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode register_curser_cap_ids: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to build transaction call: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of register_curser_cap_ids as no signer provided",
			"curserCapIds", input.CurserCapObjectIds,
		)
		return sui_ops.OpTxResult[NoObjects]{
			PackageId: input.CCIPPackageId,
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.RegisterCurserCapIds(b.GetContext(), opts, ref, ownerCap, input.CurserCapObjectIds)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute register_curser_cap_ids: %w", err)
	}

	b.Logger.Infow("CurserCap IDs registered on allowlist", "digest", tx.Digest, "curserCapIds", input.CurserCapObjectIds)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Call:      call,
	}, nil
}

// McmsDeregisterCurserCapIdsInput removes CurserCap object IDs from the on-chain allowlist.
//
// Direct path: EOA holds OwnerCap and signs the PTB (deps.Signer set).
// MCMS path: proposal leaf uses inner function deregister_curser_cap_ids, routed by bindings/mcms_encoder.go.
type McmsDeregisterCurserCapIdsInput struct {
	CCIPPackageId        string
	StateObjectId        string
	SlowOwnerCapObjectId string
	CurserCapObjectIds   []string
}

var McmsDeregisterCurserCapIdsOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "mcms_deregister_curser_cap_ids"),
	semver.MustParse("0.1.0"),
	"Revoke CurserCap curse authority via slow MCMS allowlist deregistration",
	mcmsDeregisterCurserCapIdsHandler,
)

func mcmsDeregisterCurserCapIdsHandler(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input McmsDeregisterCurserCapIdsInput) (sui_ops.OpTxResult[NoObjects], error) {
	if len(input.CurserCapObjectIds) == 0 {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("at least one curser cap object id is required")
	}

	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	ref := bind.Object{Id: input.StateObjectId}
	ownerCap := bind.Object{Id: input.SlowOwnerCapObjectId}

	encodedCall, err := contract.Encoder().DeregisterCurserCapIds(ref, ownerCap, input.CurserCapObjectIds)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to encode deregister_curser_cap_ids: %w", err)
	}

	call, err := sui_ops.ToTransactionCall(encodedCall, input.StateObjectId)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to build transaction call: %w", err)
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of deregister_curser_cap_ids as no signer provided",
			"curserCapIds", input.CurserCapObjectIds,
		)
		return sui_ops.OpTxResult[NoObjects]{
			PackageId: input.CCIPPackageId,
			Call:      call,
		}, nil
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	tx, err := contract.DeregisterCurserCapIds(b.GetContext(), opts, ref, ownerCap, input.CurserCapObjectIds)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute deregister_curser_cap_ids: %w", err)
	}

	b.Logger.Infow("CurserCap IDs deregistered from allowlist", "digest", tx.Digest, "curserCapIds", input.CurserCapObjectIds)

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest,
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
	tx, err := contract.CurseMultipleWithCurserCap(
		b.GetContext(),
		opts,
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.CurserCapObjectId},
		input.Subjects,
	)
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
