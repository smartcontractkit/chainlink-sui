package lanes

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
)

// packageIDToBytes encodes a Sui package ID as a 32-byte, left-padded address for
// cross-chain CCIP references (same convention as EVM offramp source onramp config).
func packageIDToBytes(address string) ([]byte, error) {
	out, err := suideploy.StrTo32(address)
	if err != nil {
		return nil, fmt.Errorf("invalid Sui package id %q: %w", address, err)
	}
	return out, nil
}

func loadChainState(env cldf.Environment, chainSelector uint64) (suideploy.CCIPChainState, error) {
	stateMap, err := suideploy.LoadOnchainStatesui(env)
	if err != nil {
		return suideploy.CCIPChainState{}, fmt.Errorf("load sui onchain state: %w", err)
	}
	state, ok := stateMap[chainSelector]
	if !ok {
		return suideploy.CCIPChainState{}, fmt.Errorf("sui chain %d not found in address book state", chainSelector)
	}
	return state, nil
}

func resolveOnRampPackageID(env cldf.Environment, chainSelector uint64) ([]byte, error) {
	state, err := loadChainState(env, chainSelector)
	if err != nil {
		return nil, err
	}
	if state.OnRampAddress == "" {
		return nil, fmt.Errorf("no SuiOnRamp package for chain %d in address book", chainSelector)
	}
	return packageIDToBytes(state.OnRampAddress)
}

func resolveOffRampPackageID(env cldf.Environment, chainSelector uint64) ([]byte, error) {
	state, err := loadChainState(env, chainSelector)
	if err != nil {
		return nil, err
	}
	if state.OffRampAddress == "" {
		return nil, fmt.Errorf("no SuiOffRamp package for chain %d in address book", chainSelector)
	}
	return packageIDToBytes(state.OffRampAddress)
}

func resolveCCIPPackageID(env cldf.Environment, chainSelector uint64) ([]byte, error) {
	state, err := loadChainState(env, chainSelector)
	if err != nil {
		return nil, err
	}
	if state.CCIPAddress == "" {
		return nil, fmt.Errorf("no SuiCCIP package for chain %d in address book", chainSelector)
	}
	return packageIDToBytes(state.CCIPAddress)
}

func resolveRouterPackageID(env cldf.Environment, chainSelector uint64) ([]byte, error) {
	state, err := loadChainState(env, chainSelector)
	if err != nil {
		return nil, err
	}
	if state.CCIPRouterAddress == "" {
		return nil, fmt.Errorf("no SuiRouter package for chain %d in address book", chainSelector)
	}
	return packageIDToBytes(state.CCIPRouterAddress)
}
