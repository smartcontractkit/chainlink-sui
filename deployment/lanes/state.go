package lanes

import (
	"fmt"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	suideploy "github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

func opTxDepsForChain(chains cldf_chain.BlockChains, chainSelector uint64) (sui_ops.OpTxDeps, error) {
	chain, ok := chains.SuiChains()[chainSelector]
	if !ok {
		return sui_ops.OpTxDeps{}, fmt.Errorf("sui chain with selector %d not found in environment", chainSelector)
	}

	return sui_ops.OpTxDeps{
		Client: chain.Client,
		Signer: chain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			gasBudget := uint64(400_000_000)
			return &bind.CallOpts{WaitForExecution: true, GasBudget: &gasBudget}
		},
		SuiRPC: chain.URL,
	}, nil
}

func loadOffRampState(env cldf.Environment, chainSelector uint64) (suideploy.CCIPChainState, error) {
	state, err := loadChainState(env, chainSelector)
	if err != nil {
		return suideploy.CCIPChainState{}, err
	}
	if state.CCIPObjectRef == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing CCIPObjectRef for sui chain %d in address book", chainSelector)
	}
	if state.OffRampAddress == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiOffRamp package for sui chain %d in address book", chainSelector)
	}
	if state.OffRampStateObjectId == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiOffRampStateObjectID for sui chain %d in address book", chainSelector)
	}
	if state.OffRampOwnerCapId == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiOffRampOwnerCapObjectID for sui chain %d in address book", chainSelector)
	}
	return state, nil
}

func loadSourceChainState(env cldf.Environment, chainSelector uint64) (suideploy.CCIPChainState, error) {
	state, err := loadChainState(env, chainSelector)
	if err != nil {
		return suideploy.CCIPChainState{}, err
	}
	if state.CCIPObjectRef == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing CCIPObjectRef for sui chain %d in address book", chainSelector)
	}
	if state.CCIPAddress == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiCCIP package for sui chain %d in address book", chainSelector)
	}
	if state.CCIPOwnerCapObjectId == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiCCIPOwnerCapObjectID for sui chain %d in address book", chainSelector)
	}
	if state.LinkTokenCoinMetadataId == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiLinkTokenObjectMetadataID for sui chain %d in address book", chainSelector)
	}
	if state.OnRampAddress == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiOnRamp package for sui chain %d in address book", chainSelector)
	}
	if state.OnRampStateObjectId == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiOnRampStateObjectID for sui chain %d in address book", chainSelector)
	}
	if state.OnRampOwnerCapObjectId == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiOnRampOwnerCapObjectID for sui chain %d in address book", chainSelector)
	}
	if state.CCIPRouterAddress == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiRouter package for sui chain %d in address book", chainSelector)
	}
	if state.CCIPRouterStateObjectID == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiRouterStateObjectID for sui chain %d in address book", chainSelector)
	}
	if state.CCIPRouterOwnerCapObjectId == "" {
		return suideploy.CCIPChainState{}, fmt.Errorf("missing SuiRouterOwnerCapObjectID for sui chain %d in address book", chainSelector)
	}
	return state, nil
}
