package mcmsops

import (
	"github.com/Masterminds/semver/v3"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/ops"
)

var DeployMCMSSequence = cld_ops.NewSequence(
	"sui-deploy-mcms-seq",
	semver.MustParse("0.1.0"),
	"Deploys and sets initial MCMS configuration",
	func(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input cld_ops.EmptyInput) (sui_ops.OpTxResult[DeployMCMSObjects], error) {
		deployReport, err := cld_ops.ExecuteOperation(env, DeployMCMSOp, deps, input)
		if err != nil {
			return sui_ops.OpTxResult[DeployMCMSObjects]{}, err
		}

		// TODO: Add more operations to the sequence as needed
		return deployReport.Output, nil
	},
)
