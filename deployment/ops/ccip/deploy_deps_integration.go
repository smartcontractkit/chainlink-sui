//go:build integration

package ccipops

import (
	"fmt"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

// DeployCCIPDependencyPackages publishes slow and fast MCMS packages required by CCIP.
func DeployCCIPDependencyPackages(
	bundle cld_ops.Bundle,
	deps sui_ops.OpTxDeps,
) (DeployCCIPInput, error) {
	signerAddress, err := deps.Signer.GetAddress()
	if err != nil {
		return DeployCCIPInput{}, fmt.Errorf("get signer address: %w", err)
	}

	mcmsReport, err := cld_ops.ExecuteOperation(bundle, mcmsops.DeployMCMSOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return DeployCCIPInput{}, fmt.Errorf("deploy MCMS package: %w", err)
	}

	fastMcmsReport, err := cld_ops.ExecuteOperation(bundle, mcmsops.DeployFastMCMSOp, deps, cld_ops.EmptyInput{})
	if err != nil {
		return DeployCCIPInput{}, fmt.Errorf("deploy fast MCMS package: %w", err)
	}

	return DeployCCIPInput{
		McmsPackageId:     mcmsReport.Output.PackageId,
		FastMcmsPackageId: fastMcmsReport.Output.PackageId,
		McmsOwner:         signerAddress,
	}, nil
}
