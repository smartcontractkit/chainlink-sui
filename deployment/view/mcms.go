package view

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_mcms "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms"
	module_mcms_account "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_account"
)

const typeAndVersion = "MCMS 1.6.0" // TODO: define correctly type and version for MCMS contracts

type MCMSWithTimelockView struct {
	ContractMetaData

	Bypasser  MCMSConfig `json:"bypasser"`
	Proposer  MCMSConfig `json:"proposer"`
	Canceller MCMSConfig `json:"canceller"`

	TimelockMinDelay         uint64                    `json:"timelockMinDelay"`
	TimelockBlockedFunctions []TimelockBlockedFunction `json:"timelockBlockedFunctions"`
}

type MCMSConfig struct {
	Signers     []MCMSSigner `json:"signers"`
	GroupQuorum []uint8      `json:"group_quorum"`
	GroupParent []uint8      `json:"group_parent"`
}

type MCMSSigner struct {
	Signer     string `json:"signer"`
	EvmSigner  string `json:"evm_signer"`
	Index      uint8  `json:"index"`
	GroupIndex uint8  `json:"group_index"`
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

	// Get config for each role
	bypasserConfig, err := getMCMSConfig(ctx, mcmsContract, callOpts, mcmsStateObj, bypasserRole)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get bypasser config: %w", err)
	}

	proposerConfig, err := getMCMSConfig(ctx, mcmsContract, callOpts, mcmsStateObj, proposerRole)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get proposer config: %w", err)
	}

	cancellerConfig, err := getMCMSConfig(ctx, mcmsContract, callOpts, mcmsStateObj, cancellerRole)
	if err != nil {
		return MCMSWithTimelockView{}, fmt.Errorf("failed to get canceller config: %w", err)
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
			TypeAndVersion: typeAndVersion,
		},
		Bypasser:                 bypasserConfig,
		Proposer:                 proposerConfig,
		Canceller:                cancellerConfig,
		TimelockMinDelay:         timelockMinDelay,
		TimelockBlockedFunctions: timelockBlockedFunctions,
	}, nil
}

func getMCMSConfig(
	ctx context.Context,
	mcmsContract module_mcms.IMcms,
	callOpts *bind.CallOpts,
	mcmsStateObj bind.Object,
	role byte,
) (MCMSConfig, error) {
	// Get config for the role
	config, err := mcmsContract.DevInspect().GetConfig(ctx, callOpts, mcmsStateObj, role)
	if err != nil {
		return MCMSConfig{}, fmt.Errorf("failed to get config for role %d: %w", role, err)
	}

	// Parse signers
	signers := make([]MCMSSigner, 0, len(config.Signers))
	for _, signer := range config.Signers {
		signers = append(signers, MCMSSigner{
			Signer:     "0x" + hex.EncodeToString(signer.Addr),
			EvmSigner:  "0x" + hex.EncodeToString(signer.Addr), // In Sui MCMS, both are the same address bytes
			Index:      signer.Index,
			GroupIndex: signer.Group,
		})
	}

	return MCMSConfig{
		Signers:     signers,
		GroupQuorum: config.GroupQuorums,
		GroupParent: config.GroupParents,
	}, nil
}
