package adapters

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	module_router "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_router"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	rmnops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
)

var (
	_ fastcurse.CurseAdapter        = (*CurseAdapter)(nil)
	_ fastcurse.CurseSubjectAdapter = (*CurseAdapter)(nil)
)

// CurseAdapter implements fastcurse.CurseAdapter and fastcurse.CurseSubjectAdapter for Sui.
type CurseAdapter struct {
	CCIPAddress          string
	CCIPObjectRef        string
	CCIPOwnerCapObjectID string
	RouterAddress        string
	RouterStateObjectID  string
}

// NewCurseAdapter returns a new, uninitialized CurseAdapter.
func NewCurseAdapter() *CurseAdapter {
	return &CurseAdapter{}
}

// Initialize populates the adapter's state fields from the on-chain state for the given selector.
func (c *CurseAdapter) Initialize(e cldf.Environment, selector uint64) error {
	stateMap, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("failed to load Sui on-chain state: %w", err)
	}
	state, ok := stateMap[selector]
	if !ok {
		return fmt.Errorf("Sui chain %d not found in state", selector)
	}
	c.CCIPAddress = state.CCIPAddress
	c.CCIPObjectRef = state.CCIPObjectRef
	c.CCIPOwnerCapObjectID = state.CCIPOwnerCapObjectId
	c.RouterAddress = state.CCIPRouterAddress
	c.RouterStateObjectID = state.CCIPRouterStateObjectID
	return nil
}

// IsSubjectCursedOnChain returns true if the given subject is cursed on the Sui chain.
func (c *CurseAdapter) IsSubjectCursedOnChain(e cldf.Environment, selector uint64, subject fastcurse.Subject) (bool, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return false, fmt.Errorf("Sui chain %d not found in environment", selector)
	}
	contract, err := module_rmn_remote.NewRmnRemote(c.CCIPAddress, chain.Client)
	if err != nil {
		return false, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}
	return contract.DevInspect().IsCursed(
		context.Background(),
		&bind.CallOpts{Signer: chain.Signer},
		bind.Object{Id: c.CCIPObjectRef},
		subject[:],
	)
}

// IsChainConnectedToTargetChain returns true if targetSelector is a configured destination
// on the Sui router for the chain identified by selector.
func (c *CurseAdapter) IsChainConnectedToTargetChain(e cldf.Environment, selector uint64, targetSelector uint64) (bool, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return false, fmt.Errorf("Sui chain %d not found in environment", selector)
	}
	routerContract, err := module_router.NewRouter(c.RouterAddress, chain.Client)
	if err != nil {
		return false, fmt.Errorf("failed to create router contract: %w", err)
	}
	connected, err := routerContract.DevInspect().IsChainSupported(
		context.Background(),
		&bind.CallOpts{Signer: chain.Signer},
		bind.Object{Id: c.RouterStateObjectID},
		targetSelector,
	)
	if err != nil {
		return false, fmt.Errorf("failed to check if chain %d is connected to chain %d: %w", selector, targetSelector, err)
	}
	return connected, nil
}

// IsCurseEnabledForChain always returns true for Sui — cursing is always available.
func (c *CurseAdapter) IsCurseEnabledForChain(cldf.Environment, uint64) (bool, error) {
	return true, nil
}

// SubjectToSelector converts a fastcurse.Subject to a chain selector using BigEndian encoding.
func (c *CurseAdapter) SubjectToSelector(subject fastcurse.Subject) (uint64, error) {
	return fastcurse.GenericSubjectToSelector(subject)
}

// SelectorToSubject converts a chain selector to a fastcurse.Subject using BigEndian encoding.
// Sui uses the same encoding as the generic (EVM-default) case.
func (c *CurseAdapter) SelectorToSubject(selector uint64) fastcurse.Subject {
	return fastcurse.GenericSelectorToSubject(selector)
}

// DeriveCurseAdapterVersion returns the RMN adapter version for this Sui deployment.
func (c *CurseAdapter) DeriveCurseAdapterVersion(cldf.Environment, uint64) (*semver.Version, error) {
	return semver.MustParse("1.6.0"), nil
}

// Curse returns a sequence that curses the given subjects on the specified Sui chain.
func (c *CurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		rmnops.CurseSequence.ID(),
		semver.MustParse("1.0.0"),
		rmnops.CurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (sequences.OnChainOutput, error) {
			seqInput := rmnops.CurseUncurseSeqInput{
				CCIPAddress:          c.CCIPAddress,
				CCIPObjectRef:        c.CCIPObjectRef,
				CCIPOwnerCapObjectID: c.CCIPOwnerCapObjectID,
				ChainSelector:        in.ChainSelector,
				Subjects:             in.Subjects,
			}
			seqReport, err := cldf_ops.ExecuteSequence(b, rmnops.CurseSequence, chains, seqInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute curse sequence on Sui chain %d: %w", in.ChainSelector, err)
			}
			return seqReport.Output, nil
		},
	)
}

// Uncurse returns a sequence that lifts the curse on given subjects on the specified Sui chain.
func (c *CurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		rmnops.UncurseSequence.ID(),
		semver.MustParse("1.0.0"),
		rmnops.UncurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (sequences.OnChainOutput, error) {
			seqInput := rmnops.CurseUncurseSeqInput{
				CCIPAddress:          c.CCIPAddress,
				CCIPObjectRef:        c.CCIPObjectRef,
				CCIPOwnerCapObjectID: c.CCIPOwnerCapObjectID,
				ChainSelector:        in.ChainSelector,
				Subjects:             in.Subjects,
			}
			seqReport, err := cldf_ops.ExecuteSequence(b, rmnops.UncurseSequence, chains, seqInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute uncurse sequence on Sui chain %d: %w", in.ChainSelector, err)
			}
			return seqReport.Output, nil
		},
	)
}

// ListConnectedChains returns all destination chain selectors configured in the Sui router.
func (c *CurseAdapter) ListConnectedChains(e cldf.Environment, selector uint64) ([]uint64, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return nil, fmt.Errorf("Sui chain %d not found in environment", selector)
	}
	routerContract, err := module_router.NewRouter(c.RouterAddress, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create router contract: %w", err)
	}
	connectedChains, err := routerContract.DevInspect().GetDestChains(
		context.Background(),
		&bind.CallOpts{Signer: chain.Signer},
		bind.Object{Id: c.RouterStateObjectID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected chains for chain %d: %w", selector, err)
	}
	return connectedChains, nil
}
