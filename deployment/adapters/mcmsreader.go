package adapters

import (
	"fmt"

	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	mcms_utils "github.com/smartcontractkit/chainlink-ccip/deployment/utils/mcms"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
)

type MCMSReader struct{}

func (r *MCMSReader) GetChainMetadata(e cldf.Environment, chainSelector uint64, input mcms_utils.Input) (mcmstypes.ChainMetadata, error) {
	chain, ok := e.BlockChains.SuiChains()[chainSelector]
	if !ok {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("sui chain with selector %d not found", chainSelector)
	}

	stateMap, err := suideploy.LoadOnchainStatesui(e)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to load sui onchain state: %w", err)
	}
	state, ok := stateMap[chainSelector]
	if !ok {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("sui chain %d not found in state", chainSelector)
	}

	role, err := timelockRoleFromAction(input.TimelockAction)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to get role from action: %w", err)
	}

	inspector, err := suisdk.NewInspector(chain.Client, chain.Signer, state.MCMSPackageID, role)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to create sui mcms inspector for chain %d: %w", chainSelector, err)
	}

	opCount, err := inspector.GetOpCount(e.GetContext(), state.MCMSStateObjectID)
	if err != nil {
		return mcmstypes.ChainMetadata{}, fmt.Errorf("failed to get opCount for MCMS at %s on chain %d: %w", state.MCMSStateObjectID, chainSelector, err)
	}

	return suisdk.NewChainMetadata(
		opCount,
		role,
		state.MCMSPackageID,
		state.MCMSStateObjectID,
		state.MCMSAccountStateObjectID,
		state.MCMSRegistryObjectID,
		state.MCMSTimelockObjectID,
		state.MCMSDeployerStateObjectID,
	)
}

func (r *MCMSReader) GetTimelockRef(e cldf.Environment, chainSelector uint64, input mcms_utils.Input) (datastore.AddressRef, error) {
	stateMap, err := suideploy.LoadOnchainStatesui(e)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("failed to load sui onchain state: %w", err)
	}
	state, ok := stateMap[chainSelector]
	if !ok {
		return datastore.AddressRef{}, fmt.Errorf("sui chain %d not found in state", chainSelector)
	}
	return datastore.AddressRef{
		Address:       state.MCMSTimelockObjectID,
		ChainSelector: chainSelector,
	}, nil
}

func (r *MCMSReader) GetMCMSRef(e cldf.Environment, chainSelector uint64, input mcms_utils.Input) (datastore.AddressRef, error) {
	stateMap, err := suideploy.LoadOnchainStatesui(e)
	if err != nil {
		return datastore.AddressRef{}, fmt.Errorf("failed to load sui onchain state: %w", err)
	}
	state, ok := stateMap[chainSelector]
	if !ok {
		return datastore.AddressRef{}, fmt.Errorf("sui chain %d not found in state", chainSelector)
	}
	return datastore.AddressRef{
		Address:       state.MCMSStateObjectID,
		ChainSelector: chainSelector,
	}, nil
}

// timelockRoleFromAction converts a TimelockAction to the corresponding Sui TimelockRole.
func timelockRoleFromAction(action mcmstypes.TimelockAction) (suisdk.TimelockRole, error) {
	switch action {
	case mcmstypes.TimelockActionSchedule:
		return suisdk.TimelockRoleProposer, nil
	case mcmstypes.TimelockActionBypass:
		return suisdk.TimelockRoleBypasser, nil
	case mcmstypes.TimelockActionCancel:
		return suisdk.TimelockRoleCanceller, nil
	case "":
		// Default case for empty action to avoid breaking changes
		return suisdk.TimelockRoleProposer, nil
	default:
		return suisdk.TimelockRoleProposer, fmt.Errorf("invalid timelock action: %s", action)
	}
}
