package managedtokenfaucetops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_faucet "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/managed_token_faucet/managed_token_faucet"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

type InitializeManagedTokenFaucetInput struct {
	ManagedTokenFaucetPackageId string
	CoinObjectTypeArg           string
	MintCapObjectId             string
}

type InitializeManagedTokenFaucetObjects struct {
	FaucetStateObjectId string
}

var initializeManagedTokenFaucetHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input InitializeManagedTokenFaucetInput) (output sui_ops.OpTxResult[InitializeManagedTokenFaucetObjects], err error) {
	contract, err := module_faucet.NewFaucet(input.ManagedTokenFaucetPackageId, deps.Client)
	if err != nil {
		return sui_ops.OpTxResult[InitializeManagedTokenFaucetObjects]{}, fmt.Errorf("failed to create managed token faucet contract: %w", err)
	}

	ops := deps.GetCallOpts()
	ops.Signer = deps.Signer

	tx, err := contract.Initialize(
		b.GetContext(),
		ops,
		[]string{input.CoinObjectTypeArg},
		bind.Object{Id: input.MintCapObjectId},
	)
	if err != nil {
		return sui_ops.OpTxResult[InitializeManagedTokenFaucetObjects]{}, fmt.Errorf("failed to execute managed token faucet initialize: %w", err)
	}

	faucetStateId, err := bind.FindObjectIdFromPublishTx(*tx, "faucet", "FaucetState")
	if err != nil {
		return sui_ops.OpTxResult[InitializeManagedTokenFaucetObjects]{}, fmt.Errorf("failed to find FaucetState object ID: %w", err)
	}

	return sui_ops.OpTxResult[InitializeManagedTokenFaucetObjects]{
		Digest:    tx.Digest,
		PackageId: input.ManagedTokenFaucetPackageId,
		Objects: InitializeManagedTokenFaucetObjects{
			FaucetStateObjectId: faucetStateId,
		},
	}, nil
}

var InitializeManagedTokenFaucetOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token-faucet", "package", "initialize"),
	semver.MustParse("0.1.0"),
	"Initializes the CCIP managed token faucet",
	initializeManagedTokenFaucetHandler,
)
