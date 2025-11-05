package mcmsops

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	modulemcmsaccount "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_account"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type MCMSTransferOwnershipInput struct {
	McmsPackageID   string `json:"mcmsPackageID"`
	OwnerCap        string `json:"ownerCap"`
	AccountObjectID string `json:"accountObjectID"`
}

var transferOwnershipHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSTransferOwnershipInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcmsAccount, err := modulemcmsaccount.NewMcmsAccount(input.McmsPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	tx, err := mcmsAccount.TransferOwnershipToSelf(b.GetContext(), opts, bind.Object{Id: input.OwnerCap}, bind.Object{Id: input.AccountObjectID})
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	b.Logger.Infow("Transfer ownership to MCMS executed successfully", "to", input.McmsPackageID)

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    tx.Digest,
		PackageId: input.McmsPackageID,
	}, err
}

var MCMSTransferOwnershipOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "mcms_account", "transfer_ownership"),
	semver.MustParse("0.1.0"),
	"Init the transfer ownership of the MCMS contract to itself",
	transferOwnershipHandler,
)

type MCMSAcceptOwnershipInput struct {
	McmsPackageID   string `json:"mcmsPackageID"`
	AccountObjectID string `json:"accountObjectID"`
}

var acceptOwnershipHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSAcceptOwnershipInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcmsAccount, err := modulemcmsaccount.NewMcmsAccount(input.McmsPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	encodedCall, err := mcmsAccount.Encoder().AcceptOwnershipAsTimelock(bind.Object{Id: input.AccountObjectID})
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to encode AcceptOwnership call: %w", err)
	}
	call, err := sui_ops.ToTransactionCall(encodedCall, input.AccountObjectID)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to convert encoded call to TransactionCall: %w", err)
	}

	if deps.Signer != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, errors.New("Accept Ownership as timelock cannot be called directly, only through MCMS proposal")
	}

	b.Logger.Infow("Skipping execution of AcceptOwnership on MCMS Account as per no Signer provided")
	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    "",
		PackageId: input.McmsPackageID,
		Objects:   cld_ops.EmptyInput{},
		Call:      call,
	}, nil
}

var MCMSAcceptOwnershipOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "mcms_account", "accept_ownership"),
	semver.MustParse("0.1.0"),
	"Accept ownership of the MCMS contract as MCMS",
	acceptOwnershipHandler,
)

type MCMSExecuteTransferOwnershipInput struct {
	// MCMS related
	McmsPackageID         string `json:"mcmsPackageID"`
	OwnerCap              string `json:"ownerCap"`
	AccountObjectID       string `json:"accountObjectID"`
	RegistryObjectID      string `json:"registryObjectID"`
	DeployerStateObjectID string `json:"deployerStateObjectID"`
}

var executeTransferOwnershipHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input MCMSExecuteTransferOwnershipInput) (output sui_ops.OpTxResult[cld_ops.EmptyInput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	mcmsAccount, err := modulemcmsaccount.NewMcmsAccount(input.McmsPackageID, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	tx, err := mcmsAccount.ExecuteOwnershipTransfer(b.GetContext(), opts, bind.Object{Id: input.OwnerCap}, bind.Object{Id: input.AccountObjectID}, bind.Object{Id: input.RegistryObjectID}, input.McmsPackageID)
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, err
	}

	newOwner, err := mcmsAccount.DevInspect().Owner(b.GetContext(), opts, bind.Object{Id: input.AccountObjectID})
	if err != nil {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("failed to get new owner for MCMS: %w", err)
	}

	if newOwner != input.McmsPackageID {
		return sui_ops.OpTxResult[cld_ops.EmptyInput]{}, fmt.Errorf("ownership transfer to MCMS failed: expected new owner %s, got %s", input.McmsPackageID, newOwner)
	}

	b.Logger.Infow("Ownership transfer to MCMS executed successfully", "to", input.McmsPackageID)

	return sui_ops.OpTxResult[cld_ops.EmptyInput]{
		Digest:    tx.Digest,
		PackageId: input.McmsPackageID,
	}, err
}

var MCMSExecuteTransferOwnershipOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("mcms", "mcms_account", "execute_transfer_ownership"),
	semver.MustParse("0.1.0"),
	"Execute transfer ownership of the MCMS contract to itself",
	executeTransferOwnershipHandler,
)
