package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
	mcmstypes "github.com/smartcontractkit/mcms/types"
)

// ConfigureFastCurseMCMSConfig reconfigures one or more timelock roles on an
// already-deployed fastcurse MCMS instance (qualifier RMNMCMS). It is a thin
// preset of ConfigureMCMS that forces IsFastCurse=true and IsInitialConfig=false
// (the rotation path, so the F7 ClearRoot=false warning still fires).
//
// Nil role configs are left untouched on-chain: only the roles whose *Config is
// non-nil receive a mcms::set_config call. At least one role must be set.
//
// Execution mode is chosen per call via TimelockConfig:
//   - nil           => execute directly; the env signer must hold the MCMS OwnerCap.
//     Once ownership has been transferred off the deployer EOA this
//     is not available and the call will fail to find the cap.
//   - non-nil       => generate a proposer-scheduled timelock proposal (sign with the
//     fastcurse proposer set, schedule, wait min_delay, execute).
//
// NOTE: set_config is on the bypass-forbidden (F30) list, so a rotation proposal
// MUST be scheduled via the proposer role + min_delay — it can never be bypassed.
// This is intentional: bypass cannot be used to change who can bypass.
type ConfigureFastCurseMCMSConfig struct {
	ChainSelector  uint64                `json:"chainSelector" yaml:"chainSelector"`
	Bypasser       *mcmstypes.Config     `json:"bypasser,omitempty" yaml:"bypasser,omitempty"`
	Proposer       *mcmstypes.Config     `json:"proposer,omitempty" yaml:"proposer,omitempty"`
	Canceller      *mcmstypes.Config     `json:"canceller,omitempty" yaml:"canceller,omitempty"`
	ClearRoot      bool                  `json:"clearRoot,omitempty" yaml:"clearRoot,omitempty"`
	TimelockConfig *utils.TimelockConfig `json:"timelockConfig,omitempty" yaml:"timelockConfig,omitempty"`
}

var _ cldf.ChangeSetV2[ConfigureFastCurseMCMSConfig] = ConfigureFastCurseMCMS{}

type ConfigureFastCurseMCMS struct{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (ConfigureFastCurseMCMS) VerifyPreconditions(e cldf.Environment, config ConfigureFastCurseMCMSConfig) error {
	if config.Bypasser == nil && config.Proposer == nil && config.Canceller == nil {
		return fmt.Errorf("configure fast curse mcms: at least one of bypasser/proposer/canceller must be set")
	}
	return nil
}

// Apply implements deployment.ChangeSetV2. It delegates to ConfigureMCMS with
// IsFastCurse=true; object IDs auto-fill from on-chain state because PackageId
// is left empty (ConfigureMCMS.Apply calls MCMSState(true)).
func (ConfigureFastCurseMCMS) Apply(e cldf.Environment, config ConfigureFastCurseMCMSConfig) (cldf.ChangesetOutput, error) {
	return ConfigureMCMS{}.Apply(e, ConfigureMCMSConfig{
		IsFastCurse: true,
		ConfigureMCMSSeqInput: mcmsops.ConfigureMCMSSeqInput{
			ChainSelector:   config.ChainSelector,
			Bypasser:        config.Bypasser,
			Proposer:        config.Proposer,
			Canceller:       config.Canceller,
			ClearRoot:       config.ClearRoot,
			IsInitialConfig: false,
		},
		TimelockConfig: config.TimelockConfig,
	})
}
