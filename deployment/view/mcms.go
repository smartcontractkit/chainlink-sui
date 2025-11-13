package view

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	suimcms "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mcms "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms"
	module_mcms_account "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_account"
)

const mcmsTypeAndVersion = "MCMS 1.6.0" // TODO: define correctly type and version for MCMS contracts

type MCMSWithTimelockView struct {
	ContractMetaData

	Bypasser  mcmstypes.Config `json:"bypasser"`
	Proposer  mcmstypes.Config `json:"proposer"`
	Canceller mcmstypes.Config `json:"canceller"`

	TimelockMinDelay         uint64                    `json:"timelockMinDelay"`
	TimelockBlockedFunctions []TimelockBlockedFunction `json:"timelockBlockedFunctions"`
}

type TimelockBlockedFunction struct {
	Target       string `json:"target"`
	ModuleName   string `json:"moduleName"`
	FunctionName string `json:"functionName"`
}

// GenerateMCMSWithTimelockView generates an MCMS with timelock view by querying the on-chain state
func GenerateMCMSWithTimelockView(
	ctx context.Context,
	chain sui.Chain,
	mcmsPackageID string,
	mcmsStateObjectID string,
	timelockObjectID string,
	accountStateObjectID string,
) (MCMSWithTimelockView, error) {
	if mcmsPackageID == "" || mcmsStateObjectID == "" {
		return MCMSWithTimelockView{}, fmt.Errorf("mcmsPackageID and mcmsStateObjectID cannot be empty")
	}

	mcmsStateObj := bind.Object{Id: mcmsStateObjectID}
	accountStateObj := bind.Object{Id: accountStateObjectID}
	callOpts := &bind.CallOpts{Signer: chain.Signer}

	// Create MCMS contract binding
	mcmsContract, err := module_mcms.NewMcms(mcmsPackageID, chain.Client)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to create mcms contract binding: %w", err)
	}

	// Create MCMS account contract binding to get owner
	mcmsAccountContract, err := module_mcms_account.NewMcmsAccount(mcmsPackageID, chain.Client)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to create mcms account contract binding: %w", err)
	}

	// Get owner
	owner, err := mcmsAccountContract.DevInspect().Owner(ctx, callOpts, accountStateObj)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get owner: %w", err)
	}

	// Get role constants
	bypasserRole, err := mcmsContract.DevInspect().BypasserRole(ctx, callOpts)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get bypasser role: %w", err)
	}

	proposerRole, err := mcmsContract.DevInspect().ProposerRole(ctx, callOpts)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get proposer role: %w", err)
	}

	cancellerRole, err := mcmsContract.DevInspect().CancellerRole(ctx, callOpts)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get canceller role: %w", err)
	}

	configTransformer := suimcms.NewConfigTransformer()

	// Get config for each role
	bypasserCfg, err := mcmsContract.DevInspect().GetConfig(ctx, callOpts, mcmsStateObj, bypasserRole)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get config for role %d: %w", bypasserRole, err)
	}
	tBypasserCfg, err := configTransformer.ToConfig(bypasserCfg)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to transform config for role %d: %w", bypasserRole, err)
	}

	proposerConfig, err := mcmsContract.DevInspect().GetConfig(ctx, callOpts, mcmsStateObj, proposerRole)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get proposer config: %w", err)
	}
	tProposerCfg, err := configTransformer.ToConfig(proposerConfig)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to transform config for role %d: %w", proposerRole, err)
	}

	cancellerConfig, err := mcmsContract.DevInspect().GetConfig(ctx, callOpts, mcmsStateObj, cancellerRole)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get canceller config: %w", err)
	}
	tCancellerCfg, err := configTransformer.ToConfig(cancellerConfig)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to transform config for role %d: %w", cancellerRole, err)
	}

	// Get timelock data if available
	var timelockMinDelay uint64
	var timelockBlockedFunctions []TimelockBlockedFunction

	if timelockObjectID != "" {
		timelockObj := bind.Object{Id: timelockObjectID}

		// Get minimum delay
		timelockMinDelay, err = mcmsContract.DevInspect().TimelockMinDelay(ctx, callOpts, timelockObj)
		if err != nil {
			return MCMSWithTimelockView{}, fmt.Errorf("failed to get timelock min delay: %w", err)
		}

		// Get blocked functions
		blockedFunctionsRaw, err := mcmsContract.DevInspect().TimelockGetBlockedFunctions(ctx, callOpts, timelockObj)
		if err != nil {
			return MCMSWithTimelockView{}, fmt.Errorf("failed to get timelock blocked functions: %w", err)
		}

		timelockBlockedFunctions = make([]TimelockBlockedFunction, 0, len(blockedFunctionsRaw))
		for _, fn := range blockedFunctionsRaw {
			timelockBlockedFunctions = append(timelockBlockedFunctions, TimelockBlockedFunction{
				Target:       fn.Target,
				ModuleName:   fn.ModuleName,
				FunctionName: fn.FunctionName,
			})
		}
	}

	return MCMSWithTimelockView{
		ContractMetaData: ContractMetaData{
			Address:        mcmsPackageID,
			Owner:          owner,
			TypeAndVersion: mcmsTypeAndVersion,
			StateObjectID:  mcmsStateObjectID,
		},
		Bypasser:                 *tBypasserCfg,
		Proposer:                 *tProposerCfg,
		Canceller:                *tCancellerCfg,
		TimelockMinDelay:         timelockMinDelay,
		TimelockBlockedFunctions: timelockBlockedFunctions,
	}, nil
}
