package mcmsops

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/smartcontractkit/mcms/types"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

// ConfigureMCMSSeqInput defines the input for configuring MCMS
type ConfigureMCMSSeqInput struct {
	ChainSelector               uint64 `yaml:"chainSelector"`
	PackageId                   string `yaml:"packageId"`
	McmsAccountOwnerCapObjectId string `yaml:"mcmsAccountOwnerCapObjectId"`
	McmsAccountStateObjectId    string `yaml:"mcmsAccountStateObjectId"`
	McmsMultisigStateObjectId   string `yaml:"mcmsMultisigStateObjectId"`
	// McmsTimelockObjectId opts each SetConfigMCMSOp invocation into the
	// defensive on-chain min_delay check (see MCMSSetConfigInput.TimelockObjectID
	// for details on F5/F8 Interpretation B). Optional — leave empty to skip.
	McmsTimelockObjectId string `yaml:"mcmsTimelockObjectId,omitempty"`
	// ClearRoot forwards MCMSSetConfigInput.ClearRoot to each SetConfigMCMSOp
	// invocation. Set true when rotating the signer set under a suspected
	// compromise to revoke the previously installed root and close its remaining
	// op-count window on-chain (see the warning in op_set_config.go). Defaults to
	// false, preserving the prior behavior where the old root stays executable
	// until it expires or exhausts its op-count range.
	ClearRoot bool `yaml:"clearRoot,omitempty"`

	// Optional configs for each timelock role
	// If nil, the role will not be configured
	Bypasser  *types.Config `yaml:"bypasser,omitempty"`
	Proposer  *types.Config `yaml:"proposer,omitempty"`
	Canceller *types.Config `yaml:"canceller,omitempty"`
}

type ConfigureMCMSSeqOutput struct {
	Reports []cld_ops.Report[any, any]
}

var ConfigureMCMSSequence = cld_ops.NewSequence(
	"sui-configure-mcms-seq",
	semver.MustParse("0.1.0"),
	"Configures the MCMS package with the provided timelock roles configuration",
	configureMCMS,
)

func configureMCMS(env cld_ops.Bundle, deps sui_ops.OpTxDeps, input ConfigureMCMSSeqInput) (ConfigureMCMSSeqOutput, error) {
	// Configure each timelock role if config is provided
	roleConfigs := []struct {
		config *types.Config
		role   suisdk.TimelockRole
		name   string
	}{
		{input.Bypasser, suisdk.TimelockRoleBypasser, "Bypasser"},
		{input.Canceller, suisdk.TimelockRoleCanceller, "Canceller"},
		{input.Proposer, suisdk.TimelockRoleProposer, "Proposer"},
	}

	opReports := make([]cld_ops.Report[any, any], 0)
	for _, roleConfig := range roleConfigs {
		if roleConfig.config == nil {
			continue
		}

		setConfigInput := MCMSSetConfigInput{
			ChainSelector:    input.ChainSelector,
			McmsPackageID:    input.PackageId,
			OwnerCap:         input.McmsAccountOwnerCapObjectId,
			McmsObjectID:     input.McmsMultisigStateObjectId,
			TimelockObjectID: input.McmsTimelockObjectId,
			Role:             roleConfig.role,
			Config:           *roleConfig.config,
			ClearRoot:        input.ClearRoot,
		}

		report, err := cld_ops.ExecuteOperation(env, SetConfigMCMSOp, deps, setConfigInput)
		if err != nil {
			return ConfigureMCMSSeqOutput{}, fmt.Errorf("failed to set config for role %s: %w", roleConfig.name, err)
		}
		opReports = append(opReports, report.ToGenericReport())
		env.Logger.Infow("Set MCMS config", "role", roleConfig.name, "chainSelector", input.ChainSelector)
	}

	return ConfigureMCMSSeqOutput{Reports: opReports}, nil
}
