package managedtokenops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	managedtoken "github.com/smartcontractkit/chainlink-sui/bindings/packages/managed_token"
	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

type ManagedTokenDeployInput struct {
	MCMSAddress      string
	MCMSOwnerAddress string
}

type ManagedTokenDeployOutput struct {
}

var deployHandler = func(b cld_ops.Bundle, deps sui_ops.OpTxDeps, input ManagedTokenDeployInput) (output sui_ops.OpTxResult[ManagedTokenDeployOutput], err error) {
	opts := deps.GetCallOpts()
	opts.Signer = deps.Signer
	managedTokenPackage, tx, err := managedtoken.PublishCCIPManagedToken(
		b.GetContext(),
		opts,
		deps.Client,
		input.MCMSAddress,
		input.MCMSOwnerAddress,
	)
	if err != nil {
		return sui_ops.OpTxResult[ManagedTokenDeployOutput]{}, err
	}

	return sui_ops.OpTxResult[ManagedTokenDeployOutput]{
		Digest:    tx.Digest,
		PackageId: managedTokenPackage.Address(),
	}, err
}

var DeployCCIPManagedTokenOp = cld_ops.NewOperation(
	sui_ops.NewSuiOperationName("ccip-managed-token", "package", "deploy"),
	semver.MustParse("0.1.0"),
	"Deploys the CCIP managed token package",
	deployHandler,
)
