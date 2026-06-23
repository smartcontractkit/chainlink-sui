package ccipops

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_state_object "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/state_object"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// =================== Add Package ID Operations =================== //

type AddPackageIdStateObjectInput struct {
	PackageId             string // original package ID (MCMS registry identity; used as binary when LatestPackageId is "")
	LatestPackageId       string // optional: upgraded package ID (PTB execution target when set)
	CCIPObjectRefObjectId string
	OwnerCapObjectId      string
	NewPackageId          string // the package ID to register in the CCIP StateObject
}

type AddPackageIdStateObjectObjects struct {
	// No specific objects are returned from add_package_id
}

var addPackageIdStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AddPackageIdStateObjectInput) (output sui_ops.OpTxResult[AddPackageIdStateObjectObjects], err error) {
	// When the package has been upgraded, PTB must target the latest bytecode for execution.
	binaryPkgId := input.PackageId
	if input.LatestPackageId != "" {
		binaryPkgId = input.LatestPackageId
	}
	contract, err := module_state_object.NewStateObject(binaryPkgId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	encodedCall, err := contract.Encoder().AddPackageId(bind.Object{Id: input.CCIPObjectRefObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.NewPackageId)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{}, fmt.Errorf("failed to encode AddPackageId call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRefObjectId)
	if err != nil {
		return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	// When the package has been upgraded, the on-chain MCMS registry still holds the original package's
	// proof type, so tx.To must be the original package ID. The PTB MoveCall must target the latest package
	// so upgraded bytecode runs. Use LatestPackageId so the proposal generator can separate the two.
	if input.LatestPackageId != "" {
		call.LatestPackageID = call.PackageID // current PackageID is the latest (from binaryPkgId)
		call.PackageID = input.PackageId      // replace with original for on-chain identity
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AddPackageId on StateObject as per no Signer provided", "newPackageId", input.NewPackageId)
		return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{
			Digest:    "",
			PackageId: input.PackageId,
			Objects:   AddPackageIdStateObjectObjects{},
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
		return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{}, fmt.Errorf("failed to execute AddPackageId on StateObject: %w", err)
	}

	b.Logger.Infow("Package ID added to CCIP StateObject", "newPackageId", input.NewPackageId)

	return sui_ops.OpTxResult[AddPackageIdStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.PackageId,
		Objects:   AddPackageIdStateObjectObjects{},
		Call:      call,
	}, nil
}

var AddPackageIdStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "add_package_id"),
	semver.MustParse("0.1.0"),
	"Adds a new package ID to the CCIP StateObject for upgrade tracking",
	addPackageIdStateObjectHandler,
)

// =================== Remove Package ID Operations =================== //

type RemovePackageIdStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
	OwnerCapObjectId      string
	PackageId             string
}

type RemovePackageIdStateObjectObjects struct {
	// No specific objects are returned from remove_package_id
}

var removePackageIdStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RemovePackageIdStateObjectInput) (output sui_ops.OpTxResult[RemovePackageIdStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	encodedCall, err := contract.Encoder().RemovePackageId(bind.Object{Id: input.CCIPObjectRefObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.PackageId)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{}, fmt.Errorf("failed to encode RemovePackageId call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRefObjectId)
	if err != nil {
		return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of RemovePackageId on StateObject as per no Signer provided", "packageId", input.PackageId)
		return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   RemovePackageIdStateObjectObjects{},
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
		return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{}, fmt.Errorf("failed to execute RemovePackageId on StateObject: %w", err)
	}

	b.Logger.Infow("Package ID removed from CCIP StateObject", "packageId", input.PackageId)

	return sui_ops.OpTxResult[RemovePackageIdStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   RemovePackageIdStateObjectObjects{},
		Call:      call,
	}, nil
}

var RemovePackageIdStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "remove_package_id"),
	semver.MustParse("0.1.0"),
	"Removes a package ID from the CCIP StateObject",
	removePackageIdStateObjectHandler,
)

// =================== Get Owner Cap ID Operations =================== //

type GetOwnerCapIdStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetOwnerCapIdStateObjectOutput struct {
	OwnerCapId string
}

var getOwnerCapIdStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetOwnerCapIdStateObjectInput) (output sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	ownerCapId, err := contract.DevInspect().OwnerCapId(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput]{}, fmt.Errorf("failed to get owner cap ID from StateObject: %w", err)
	}

	b.Logger.Infow("Owner cap ID retrieved from CCIP StateObject",
		"ownerCapId", ownerCapId.Id,
	)

	return sui_ops.OpTxResult[GetOwnerCapIdStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetOwnerCapIdStateObjectOutput{
			OwnerCapId: ownerCapId.Id,
		},
	}, nil
}

var GetOwnerCapIdStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_owner_cap_id"),
	semver.MustParse("0.1.0"),
	"Gets the owner cap ID from the CCIP StateObject",
	getOwnerCapIdStateObjectHandler,
)

// =================== Get Owner Operations =================== //

type GetOwnerStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetOwnerStateObjectOutput struct {
	Owner string
}

var getOwnerStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetOwnerStateObjectInput) (output sui_ops.OpTxResult[GetOwnerStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	owner, err := contract.DevInspect().Owner(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetOwnerStateObjectOutput]{}, fmt.Errorf("failed to get owner from StateObject: %w", err)
	}

	b.Logger.Infow("Owner retrieved from CCIP StateObject",
		"owner", owner,
	)

	return sui_ops.OpTxResult[GetOwnerStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetOwnerStateObjectOutput{
			Owner: owner,
		},
	}, nil
}

var GetOwnerStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_owner"),
	semver.MustParse("0.1.0"),
	"Gets the owner from the CCIP StateObject",
	getOwnerStateObjectHandler,
)

// =================== Get Pending Transfer Operations =================== //

type GetPendingTransferStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type GetPendingTransferStateObjectOutput struct {
	HasPendingTransfer      bool
	PendingTransferFrom     *string
	PendingTransferTo       *string
	PendingTransferAccepted *bool
}

var getPendingTransferStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input GetPendingTransferStateObjectInput) (output sui_ops.OpTxResult[GetPendingTransferStateObjectOutput], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer

	hasPending, err := contract.DevInspect().HasPendingTransfer(
		b.GetContext(),
		opts,
		bind.Object{Id: input.CCIPObjectRefObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to check pending transfer: %w", err)
	}

	var pendingFrom *string
	var pendingTo *string
	var pendingAccepted *bool

	if hasPending {
		pendingFrom, err = contract.DevInspect().PendingTransferFrom(
			b.GetContext(),
			opts,
			bind.Object{Id: input.CCIPObjectRefObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to get pending transfer from: %w", err)
		}

		pendingTo, err = contract.DevInspect().PendingTransferTo(
			b.GetContext(),
			opts,
			bind.Object{Id: input.CCIPObjectRefObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to get pending transfer to: %w", err)
		}

		pendingAccepted, err = contract.DevInspect().PendingTransferAccepted(
			b.GetContext(),
			opts,
			bind.Object{Id: input.CCIPObjectRefObjectId},
		)
		if err != nil {
			return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{}, fmt.Errorf("failed to get pending transfer accepted: %w", err)
		}
	}

	b.Logger.Infow("Pending transfer info retrieved from CCIP StateObject",
		"hasPending", hasPending,
		"pendingFrom", pendingFrom,
		"pendingTo", pendingTo,
		"pendingAccepted", pendingAccepted,
	)

	return sui_ops.OpTxResult[GetPendingTransferStateObjectOutput]{
		Digest:    "",
		PackageId: input.CCIPPackageId,
		Objects: GetPendingTransferStateObjectOutput{
			HasPendingTransfer:      hasPending,
			PendingTransferFrom:     pendingFrom,
			PendingTransferTo:       pendingTo,
			PendingTransferAccepted: pendingAccepted,
		},
	}, nil
}

var GetPendingTransferStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "get_pending_transfer"),
	semver.MustParse("0.1.0"),
	"Gets pending transfer information from the CCIP StateObject",
	getPendingTransferStateObjectHandler,
)

// =================== Transfer Ownership Operations =================== //

type TransferOwnershipStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
	OwnerCapObjectId      string
	To                    string
}

type TransferOwnershipStateObjectObjects struct {
	// No specific objects are returned from transfer_ownership
}

var transferOwnershipStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input TransferOwnershipStateObjectInput) (output sui_ops.OpTxResult[TransferOwnershipStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	encodedCall, err := contract.Encoder().TransferOwnership(bind.Object{Id: input.CCIPObjectRefObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.To)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{}, fmt.Errorf("failed to encode TransferOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRefObjectId)
	if err != nil {
		return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of TransferOwnership on StateObject as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   TransferOwnershipStateObjectObjects{},
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
		return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{}, fmt.Errorf("failed to execute TransferOwnership on StateObject: %w", err)
	}

	b.Logger.Infow("Ownership transfer initiated for CCIP StateObject", "to", input.To)

	return sui_ops.OpTxResult[TransferOwnershipStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   TransferOwnershipStateObjectObjects{},
		Call:      call,
	}, nil
}

var TransferOwnershipStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "transfer_ownership"),
	semver.MustParse("0.1.0"),
	"Transfers ownership of the CCIP StateObject",
	transferOwnershipStateObjectHandler,
)

// =================== Accept Ownership Operations =================== //

type AcceptOwnershipStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
}

type AcceptOwnershipStateObjectObjects struct {
	// No specific objects are returned from accept_ownership
}

var acceptOwnershipStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AcceptOwnershipStateObjectInput) (output sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	encodedCall, err := contract.Encoder().AcceptOwnership(bind.Object{Id: input.CCIPObjectRefObjectId})
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRefObjectId)
	if err != nil {
		return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on StateObject as per no Signer provided")
		return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   AcceptOwnershipStateObjectObjects{},
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
		return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{}, fmt.Errorf("failed to execute AcceptOwnership on StateObject: %w", err)
	}

	b.Logger.Infow("Ownership accepted for CCIP StateObject")

	return sui_ops.OpTxResult[AcceptOwnershipStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   AcceptOwnershipStateObjectObjects{},
		Call:      call,
	}, nil
}

var AcceptOwnershipStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "accept_ownership"),
	semver.MustParse("0.1.0"),
	"Accepts ownership of the CCIP StateObject",
	acceptOwnershipStateObjectHandler,
)

// =================== Execute Ownership Transfer To MCMS Operations =================== //

type ExecuteOwnershipTransferToMcmsStateObjectInput struct {
	CCIPPackageId         string
	CCIPObjectRefObjectId string
	OwnerCapObjectId      string
	RegistryObjectId      string
	To                    string
}

type ExecuteOwnershipTransferToMcmsStateObjectObjects struct {
	// No specific objects are returned from execute_ownership_transfer_to_mcms
}

var executeOwnershipTransferToMcmsStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ExecuteOwnershipTransferToMcmsStateObjectInput) (output sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects], err error) {
	contract, err := module_state_object.NewStateObject(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{}, fmt.Errorf("failed to create StateObject contract: %w", err)
	}

	encodedCall, err := contract.Encoder().ExecuteOwnershipTransferToMcms(
		bind.Object{Id: input.CCIPObjectRefObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		bind.Object{Id: input.RegistryObjectId},
		input.To,
	)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{}, fmt.Errorf("failed to encode ExecuteOwnershipTransferToMcms call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.CCIPObjectRefObjectId)
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of ExecuteOwnershipTransferToMcms on StateObject as per no Signer provided", "to", input.To)
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{
			Digest:    "",
			PackageId: input.CCIPPackageId,
			Objects:   ExecuteOwnershipTransferToMcmsStateObjectObjects{},
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
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{}, fmt.Errorf("failed to execute ExecuteOwnershipTransferToMcms on StateObject: %w", err)
	}

	newOwner, err := contract.DevInspect().Owner(b.GetContext(), opts, bind.Object{Id: input.CCIPObjectRefObjectId})
	if err != nil {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{}, fmt.Errorf("failed to get new owner for StateObject: %w", err)
	}

	if newOwner != input.To {
		return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{}, fmt.Errorf("ownership transfer to MCMS failed for StateObject: expected new owner %s, got %s", input.To, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully for CCIP StateObject", "to", input.To)

	return sui_ops.OpTxResult[ExecuteOwnershipTransferToMcmsStateObjectObjects]{
		Digest:    tx.Digest,
		PackageId: input.CCIPPackageId,
		Objects:   ExecuteOwnershipTransferToMcmsStateObjectObjects{},
		Call:      call,
	}, nil
}

var ExecuteOwnershipTransferToMcmsStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "execute_ownership_transfer_to_mcms"),
	semver.MustParse("0.1.0"),
	"Executes ownership transfer to MCMS for the CCIP StateObject",
	executeOwnershipTransferToMcmsStateObjectHandler,
)

// AddAllowedModulesStateObjectInput configures mcms_registry::add_allowed_modules<T>
type AddAllowedModulesStateObjectInput struct {
	PackageId         string // original CCIP package ID (MCMS registry identity; used as binary when LatestPackageId is "")
	LatestPackageId   string // optional: upgraded CCIP package ID (PTB execution target when set)
	MCMSRegistryObjId string // mcms_registry Registry shared object ID
	AllowedModules    []string
}

type AddAllowedModulesStateObjectObjects struct{}

// bcsAppendULEB128 appends n encoded as an unsigned LEB128 integer to data (BCS length encoding).
func bcsAppendULEB128(data []byte, n uint64) []byte {
	for {
		b := byte(n & 0x7F)
		n >>= 7
		if n == 0 {
			return append(data, b)
		}
		data = append(data, b|0x80)
	}
}

var addAllowedModulesStateObjectHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AddAllowedModulesStateObjectInput) (output sui_ops.OpTxResult[AddAllowedModulesStateObjectObjects], err error) {
	if strings.TrimSpace(input.PackageId) == "" {
		return sui_ops.OpTxResult[AddAllowedModulesStateObjectObjects]{}, fmt.Errorf("PackageId is required")
	}
	if len(input.AllowedModules) == 0 {
		return sui_ops.OpTxResult[AddAllowedModulesStateObjectObjects]{}, fmt.Errorf("AllowedModules must not be empty")
	}

	// When the package has been upgraded, PTB must target the latest bytecode for execution.
	binaryPkgId := input.PackageId
	if input.LatestPackageId != "" {
		binaryPkgId = input.LatestPackageId
	}

	// state_object::mcms_add_allowed_modules routes to mcms_registry::add_allowed_modules internally.
	// There is no direct (non-MCMS) add_allowed_modules function; this op is for proposal generation only.
	// The function BCS-decodes data as:
	//   1. validate_obj_addr(registry) → 32 raw address bytes
	//   2. deserialize_vector<vector<u8>>(modules) → ULEB128(N) + each: ULEB128(len)+bytes
	registryHex := strings.TrimPrefix(strings.TrimSpace(input.MCMSRegistryObjId), "0x")
	if len(registryHex)%2 != 0 {
		registryHex = "0" + registryHex
	}
	registryBytes, err := hex.DecodeString(registryHex)
	if err != nil {
		return sui_ops.OpTxResult[AddAllowedModulesStateObjectObjects]{}, fmt.Errorf("invalid MCMSRegistryObjId: %w", err)
	}
	if len(registryBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(registryBytes):], registryBytes)
		registryBytes = padded
	}
	if len(registryBytes) != 32 {
		return sui_ops.OpTxResult[AddAllowedModulesStateObjectObjects]{}, fmt.Errorf("MCMSRegistryObjId must be a 32-byte Sui address, got %d bytes", len(registryBytes))
	}

	data := make([]byte, 0, 32+1+len(input.AllowedModules)*16)
	data = append(data, registryBytes...)
	data = bcsAppendULEB128(data, uint64(len(input.AllowedModules)))
	for _, m := range input.AllowedModules {
		modBytes := []byte(m)
		data = bcsAppendULEB128(data, uint64(len(modBytes)))
		data = append(data, modBytes...)
	}

	call := sui_ops.TransactionCall{
		PackageID:  binaryPkgId,
		Module:     "state_object",
		Function:   "add_allowed_modules",
		Data:       data,
		StateObjID: input.MCMSRegistryObjId,
		TypeArgs:   []string{},
	}
	// When the package has been upgraded, the on-chain MCMS registry still holds the original package's
	// proof type, so tx.To must be the original package ID. The PTB MoveCall must target the latest package
	// so upgraded bytecode runs. Use LatestPackageID so the proposal generator can separate the two.
	if input.LatestPackageId != "" {
		call.LatestPackageID = call.PackageID // current PackageID is the latest (from binaryPkgId)
		call.PackageID = input.PackageId      // replace with original for on-chain MCMS identity
	}

	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of add_allowed_modules on StateObject as per no Signer provided",
			"registry", input.MCMSRegistryObjId, "modules", input.AllowedModules)
		return sui_ops.OpTxResult[AddAllowedModulesStateObjectObjects]{
			Digest:    "",
			PackageId: input.PackageId,
			Objects:   AddAllowedModulesStateObjectObjects{},
			Call:      call,
		}, nil
	}

	return sui_ops.OpTxResult[AddAllowedModulesStateObjectObjects]{}, fmt.Errorf(
		"direct execution of state_object::add_allowed_modules is not supported; the function can only be invoked through the MCMS executor",
	)
}

var AddAllowedModulesStateObjectOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "state_object", "add_allowed_modules"),
	semver.MustParse("0.1.0"),
	"Adds module names to the MCMS registry allowlist via the CCIP state_object MCMS entrypoint (state_object::add_allowed_modules)",
	addAllowedModulesStateObjectHandler,
)
