package lanes

import (
	"fmt"

	"github.com/Masterminds/semver/v3"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccip_offramp_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
)

// ConfigureLaneLegAsDest wires a remote source chain into the Sui OffRamp.
// Sui has no router offramp step for inbound lanes (unlike Solana).
var ConfigureLaneLegAsDest = cldf_ops.NewSequence(
	"ConfigureLaneLegAsDest",
	semver.MustParse("1.6.0"),
	"Configures lane leg as destination on CCIP 1.6.0 for Sui",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, input laneapi.UpdateLanesInput) (sequences.OnChainOutput, error) {
		if input.Source == nil || input.Dest == nil {
			return sequences.OnChainOutput{}, fmt.Errorf("ConfigureLaneLegAsDest requires Source and Dest chain definitions")
		}
		if len(input.Source.OnRamp) == 0 {
			return sequences.OnChainOutput{}, fmt.Errorf("source OnRamp address required to configure Sui as destination for chain %d", input.Source.Selector)
		}

		env, err := connectChainsEnvironment()
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		state, err := loadOffRampState(env, input.Dest.Selector)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		deps, err := opTxDepsForChain(chains, input.Dest.Selector)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		latestIDs := resolveLatestPackageIDs(input.Dest.Selector)

		opInput := ccip_offramp_ops.ApplySourceChainConfigUpdateInput{
			CCIPObjectRef:                         state.CCIPObjectRef,
			OffRampPackageId:                      state.OffRampAddress,
			LatestPackageId:                       latestIDs.OffRamp,
			OffRampStateId:                        state.OffRampStateObjectId,
			OwnerCapObjectId:                      state.OffRampOwnerCapId,
			SourceChainsSelectors:                 []uint64{input.Source.Selector},
			SourceChainsIsEnabled:                 []bool{!input.IsDisabled},
			SourceChainsIsRMNVerificationDisabled: []bool{!input.Source.RMNVerificationEnabled},
			SourceChainsOnRamp:                    [][]byte{append([]byte(nil), input.Source.OnRamp...)},
		}

		report, err := cldf_ops.ExecuteOperation(b, ccip_offramp_ops.ApplySourceChainConfigUpdatesOp, deps, opInput)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf(
				"failed to apply source chain config on Sui OffRamp for source %d: %w",
				input.Source.Selector, err,
			)
		}

		var out sequences.OnChainOutput
		if err := appendMCMSBatchOpFromCall(&out, input.Dest.Selector, report.Output.Call, deps); err != nil {
			return sequences.OnChainOutput{}, err
		}
		return out, nil
	},
)

func (a *SuiAdapter) ConfigureLaneLegAsDest() *cldf_ops.Sequence[laneapi.UpdateLanesInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return ConfigureLaneLegAsDest
}
