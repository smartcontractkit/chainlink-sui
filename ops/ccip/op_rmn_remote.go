package ccipops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type InitRMNRemoteObjects struct {
	RMNRemoteStateObjectId string
}

type InitRMNRemoteInput struct {
	CCIPPackageId      string
	StateObjectId      string
	OwnerCapObjectId   string
	LocalChainSelector uint64
}

var initRMNRemoteHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitRMNRemoteInput) (output sui_ops.OpTxResult[InitRMNRemoteObjects], err error) {
	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitRMNRemoteObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	method := contract.Initialize(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.LocalChainSelector,
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitRMNRemoteObjects]{}, fmt.Errorf("failed to execute RMN Remote initialization: %w", err)
	}

	obj1, err1 := bind.FindObjectIdFromPublishTx(*tx, "rmn_remote", "RMNRemoteState")
	if err1 != nil {
		return sui_ops.OpTxResult[InitRMNRemoteObjects]{}, fmt.Errorf("failed to find object IDs in tx: %w", err)
	}

	return sui_ops.OpTxResult[InitRMNRemoteObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects: InitRMNRemoteObjects{
			RMNRemoteStateObjectId: obj1,
		},
	}, err
}

var RMNRemoteInitializeOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP RMN Remote contract",
	initRMNRemoteHandler,
)

type RMNRemoteSetConfigInput struct {
	CCIPPackageId               string
	StateObjectId               string
	OwnerCapObjectId            string
	RmnHomeContractConfigDigest []byte
	SignerOnchainPublicKeys     [][]byte
	NodeIndexes                 []uint64
	FSign                       uint64
}

var handlerSetconfig = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input RMNRemoteSetConfigInput) (output sui_ops.OpTxResult[NoObjects], err error) {
	contract, err := module_rmn_remote.NewRmnRemote(input.CCIPPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}

	method := contract.SetConfig(
		bind.Object{Id: input.StateObjectId},
		bind.Object{Id: input.OwnerCapObjectId},
		input.RmnHomeContractConfigDigest,
		input.SignerOnchainPublicKeys,
		input.NodeIndexes,
		input.FSign,
	)
	tx, err := method.Execute(b.GetContext(), deps.GetTxOpts(), deps.Signer, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[NoObjects]{}, fmt.Errorf("failed to execute set config in RMN Remote: %w", err)
	}

	return sui_ops.OpTxResult[NoObjects]{
		Digest:    tx.Digest.String(),
		PackageId: input.CCIPPackageId,
		Objects:   NoObjects{},
	}, err
}

var RMNRemoteSetConfigOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip", "rmn_remote", "set_config"),
	semver.MustParse("0.1.0"),
	"Sets config the CCIP RMN Remote contract",
	handlerSetconfig,
)
