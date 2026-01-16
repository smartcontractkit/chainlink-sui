package mcmsops

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/transaction"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type AddModulesMCMSInput struct {
	MCMSPackageId     string
	MCMSRegistryObjId string
	AllowedModules    []string
}

var addModulesHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input AddModulesMCMSInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {

	registryBytes, err := hex.DecodeString(strings.TrimPrefix(input.MCMSRegistryObjId, "0x"))
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to convert registry object ID to bytes: %w", err)
	}
	nameBytes, err := encodeAddModulesCall(input.AllowedModules)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to encode add modules call: %w", err)
	}
	encodedCall := bind.EncodedCall{
		Module: bind.ModuleInformation{
			PackageID:   input.MCMSPackageId,
			ModuleName:  "mcms_registry",
			PackageName: "mcms",
		},
		Function: "add_allowed_modules",
		CallArgs: []*bind.EncodedCallArgument{
			{
				CallArg: &transaction.CallArg{
					Pure: &transaction.Pure{
						Bytes: registryBytes,
					},
				},
			},
			{
				CallArg: &transaction.CallArg{
					Pure: &transaction.Pure{
						Bytes: nameBytes,
					},
				},
			},
		},
	}
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(&encodedCall, input.MCMSRegistryObjId)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}
	if deps.Signer == nil {
		b.Logger.Infow("Skipping execution of AcceptOwnership on StateObject as per no Signer provided")
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{
			Digest:    "",
			PackageId: input.MCMSPackageId,
			Call:      call,
		}, nil
	}

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("cannot execute this call directly: %w", err)
}

var AddModulesMCMSOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "package", "add_modules"),
	semver.MustParse("0.1.0"),
	"Add modules to the MCMS registry",
	addModulesHandler,
)

func encodeAddModulesCall(moduleNames []string) ([]byte, error) {
	// Create a BCS serializer
	serializer := &bcs.Serializer{}

	// Serialize the length of the vector using ULEB128 (BCS standard for vector lengths)
	serializer.Uleb128(uint32(len(moduleNames)))

	// Serialize each module name as a vector<u8> (since strings in Move are vector<u8>)
	for _, moduleName := range moduleNames {
		// Each string is serialized as:
		// 1. Length of the string (ULEB128)
		// 2. The bytes of the string
		moduleBytes := []byte(moduleName)
		serializer.Uleb128(uint32(len(moduleBytes)))
		serializer.FixedBytes(moduleBytes)
	}

	return serializer.ToBytes(), nil
}

// Exports every operation available so they can be registered to be used in dynamic changesets
var AllOperationsMCMS = []any{
	*MCMSAcceptOwnershipOp,
	*MCMSTransferOwnershipOp,
	*MCMSExecuteTransferOwnershipOp,
	*SetConfigMCMSOp,
	*AddModulesMCMSOp,
}
