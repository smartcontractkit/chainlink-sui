package adapters

import (
	"bytes"
	"context"
	"fmt"
	"sync"

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

// SuiChainState holds the per-chain on-chain state a CurseAdapter needs to curse/uncurse and
// verify curses for a single Sui chain selector. It is the selector-scoped counterpart of the
// flat fields the adapter previously held; exporting it lets tests and integration suites
// populate state directly via SetChainState without a real Environment.
type SuiChainState struct {
	CCIPAddress          string
	LatestCCIPPackageID  string
	CCIPObjectRef        string
	CCIPOwnerCapObjectID string
	CurserCapObjectID    string
	RouterAddress        string
	RouterStateObjectID  string
}

// CurseAdapter implements fastcurse.CurseAdapter and fastcurse.CurseSubjectAdapter for Sui.
//
// Per-chain state is cached in a selector-keyed map so that initializing a second Sui selector
// on the shared adapter instance (the fastcurse registry stores one adapter per family+version)
// does not clobber the first — mirroring the EVM adapter's rmnAddressCache/routerAddressCache
// maps. Initialize is additive and idempotent; all stateful reads resolve state by selector.
//
// Curse() routes through FastCurseCurseSequence and requires a registered CurserCap.
// Uncurse uses CCIP OwnerCap via UncurseSequence for slow MCMS proposals.
type CurseAdapter struct {
	mu     sync.RWMutex
	states map[uint64]SuiChainState
}

// NewCurseAdapter returns a new, uninitialized CurseAdapter.
func NewCurseAdapter() *CurseAdapter {
	return &CurseAdapter{states: make(map[uint64]SuiChainState)}
}

// suiChainStateFromCCIP converts a deployment.CCIPChainState (the on-chain state shape returned by
// LoadOnchainStatesui) into a SuiChainState, preserving the field-name casing the adapter exposes.
func suiChainStateFromCCIP(state deployment.CCIPChainState) SuiChainState {
	return SuiChainState{
		CCIPAddress:          state.CCIPAddress,
		LatestCCIPPackageID:  state.LatestCCIPPackageID,
		CCIPObjectRef:        state.CCIPObjectRef,
		CCIPOwnerCapObjectID: state.CCIPOwnerCapObjectId,
		CurserCapObjectID:    state.CurserCapObjectId,
		RouterAddress:        state.CCIPRouterAddress,
		RouterStateObjectID:  state.CCIPRouterStateObjectID,
	}
}

// Initialize populates the adapter's per-selector state from the on-chain state for the given
// selector. Chain metadata is resolved via LoadOnchainStatesui, which prefers datastore address
// refs (address_refs.json) and falls back to the legacy address book (addresses.json).
//
// Initialize is additive and idempotent: a selector already cached is left untouched (a second
// init for a different selector does not clobber the first). This mirrors the EVM adapter and
// makes the shared adapter instance safe for any selector cardinality.
func (c *CurseAdapter) Initialize(e cldf.Environment, selector uint64) error {
	stateMap, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("failed to load Sui on-chain state: %w", err)
	}
	state, ok := stateMap[selector]
	if !ok {
		return fmt.Errorf("Sui chain %d not found in state", selector)
	}
	st := suiChainStateFromCCIP(state)
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.states[selector]; exists {
		return nil
	}
	c.states[selector] = st
	return nil
}

// SetChainState stores per-selector state directly, bypassing Initialize. It is intended for
// tests and integration suites that construct state without a real Environment.
func (c *CurseAdapter) SetChainState(selector uint64, state SuiChainState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states[selector] = state
}

// ChainState returns the cached state for selector and whether it was found.
func (c *CurseAdapter) ChainState(selector uint64) (SuiChainState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	st, ok := c.states[selector]
	return st, ok
}

// stateFor returns the cached state for selector, erroring if the selector was never initialized.
func (c *CurseAdapter) stateFor(selector uint64) (SuiChainState, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	st, ok := c.states[selector]
	if !ok {
		return SuiChainState{}, fmt.Errorf("sui curse adapter not initialized for selector %d", selector)
	}
	return st, nil
}

// IsSubjectCursedOnChain returns true when subject is explicitly present in the on-chain
// cursed-subjects set. A global curse does not cause lane subjects to appear cursed here;
// call with GlobalCurseSubject() to test global curse state.
func (c *CurseAdapter) IsSubjectCursedOnChain(e cldf.Environment, selector uint64, subject fastcurse.Subject) (bool, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return false, fmt.Errorf("Sui chain %d not found in environment", selector)
	}
	st, err := c.stateFor(selector)
	if err != nil {
		return false, err
	}
	contract, err := module_rmn_remote.NewRmnRemote(st.CCIPAddress, chain.Client)
	if err != nil {
		return false, fmt.Errorf("failed to create RMN Remote contract: %w", err)
	}
	cursedSubjects, err := contract.DevInspect().GetCursedSubjects(
		context.Background(),
		&bind.CallOpts{Signer: chain.Signer},
		bind.Object{Id: st.CCIPObjectRef},
	)
	if err != nil {
		return false, fmt.Errorf("failed to get cursed subjects on Sui chain %d: %w", selector, err)
	}
	return subjectInCursedSubjects(cursedSubjects, subject), nil
}

func subjectInCursedSubjects(cursedSubjects [][]byte, subject fastcurse.Subject) bool {
	for _, cursedSubject := range cursedSubjects {
		if len(cursedSubject) == len(subject) && bytes.Equal(cursedSubject, subject[:]) {
			return true
		}
	}
	return false
}

// IsChainConnectedToTargetChain returns true if targetSelector is a configured destination
// on the Sui router for the chain identified by selector.
func (c *CurseAdapter) IsChainConnectedToTargetChain(e cldf.Environment, selector uint64, targetSelector uint64) (bool, error) {
	chain, ok := e.BlockChains.SuiChains()[selector]
	if !ok {
		return false, fmt.Errorf("Sui chain %d not found in environment", selector)
	}
	st, err := c.stateFor(selector)
	if err != nil {
		return false, err
	}
	routerContract, err := module_router.NewRouter(st.RouterAddress, chain.Client)
	if err != nil {
		return false, fmt.Errorf("failed to create router contract: %w", err)
	}
	connected, err := routerContract.DevInspect().IsChainSupported(
		context.Background(),
		&bind.CallOpts{Signer: chain.Signer},
		bind.Object{Id: st.RouterStateObjectID},
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

// Curse returns a sequence that curses the given subjects via CurserCap on the fast MCMS path.
//
// The adapter is only invoked by the generic fastcurse framework (fastcurse.CurseChangeset /
// fastcurse.GloballyCurseChainChangeset), which always builds an MCMS proposal via
// OutputBuilder.Build(cfg.MCMS). We therefore force ProposalOnly: true so the underlying
// operation never attempts direct execution — the CurserCap lives inside the fast MCMS
// Registry (registered via mint_and_register_curser_cap) and has no top-level owner, so
// direct PTB assembly against its object ID would fail with "Object not found".
//
// State is resolved by in.ChainSelector at execution time (not captured at sequence-build time),
// so the shared adapter instance remains safe when multiple Sui selectors are cursed in one run.
func (c *CurseAdapter) Curse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		rmnops.FastCurseCurseSequence.ID(),
		semver.MustParse("1.0.0"),
		rmnops.FastCurseCurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (sequences.OnChainOutput, error) {
			st, err := c.stateFor(in.ChainSelector)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}
			if st.CurserCapObjectID == "" {
				return sequences.OnChainOutput{}, fmt.Errorf(
					"registered CurserCap is required to curse on Sui chain %d; run sui_register_curser_cap and sui_record_curser_cap first",
					in.ChainSelector,
				)
			}
			seqInput := rmnops.FastCurseSeqInput{
				CCIPAddress:         st.CCIPAddress,
				LatestCCIPPackageID: st.LatestCCIPPackageID,
				CCIPObjectRef:       st.CCIPObjectRef,
				CurserCapObjectID:   st.CurserCapObjectID,
				ChainSelector:       in.ChainSelector,
				Subjects:            in.Subjects,
				ProposalOnly:        true,
			}
			seqReport, err := cldf_ops.ExecuteSequence(b, rmnops.FastCurseCurseSequence, chains, seqInput)
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("failed to execute fast curse sequence on Sui chain %d: %w", in.ChainSelector, err)
			}
			return seqReport.Output, nil
		},
	)
}

// Uncurse returns a sequence that lifts the curse on given subjects on the specified Sui chain.
// Uncurse always uses OwnerCap via slow MCMS; fast MCMS cannot uncurse.
//
// Like Curse, this adapter is only invoked by the generic fastcurse framework in MCMS-proposal
// mode, so we set ProposalOnly: true. In production the OwnerCap is owned by the slow MCMS
// timelock — the loaded Sui deployer signer cannot authorize direct execution anyway, so the
// only useful output here is an MCMS proposal leaf.
//
// State is resolved by in.ChainSelector at execution time (not captured at sequence-build time),
// so the shared adapter instance remains safe when multiple Sui selectors are cursed in one run.
func (c *CurseAdapter) Uncurse() *cldf_ops.Sequence[fastcurse.CurseInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return cldf_ops.NewSequence(
		rmnops.UncurseSequence.ID(),
		semver.MustParse("1.0.0"),
		rmnops.UncurseSequence.Description(),
		func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, in fastcurse.CurseInput) (sequences.OnChainOutput, error) {
			st, err := c.stateFor(in.ChainSelector)
			if err != nil {
				return sequences.OnChainOutput{}, err
			}
			seqInput := rmnops.CurseUncurseSeqInput{
				CCIPAddress:          st.CCIPAddress,
				LatestCCIPPackageID:  st.LatestCCIPPackageID,
				CCIPObjectRef:        st.CCIPObjectRef,
				CCIPOwnerCapObjectID: st.CCIPOwnerCapObjectID,
				ChainSelector:        in.ChainSelector,
				Subjects:             in.Subjects,
				ProposalOnly:         true,
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
	st, err := c.stateFor(selector)
	if err != nil {
		return nil, err
	}
	routerContract, err := module_router.NewRouter(st.RouterAddress, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create router contract: %w", err)
	}
	connectedChains, err := routerContract.DevInspect().GetDestChains(
		context.Background(),
		&bind.CallOpts{Signer: chain.Signer},
		bind.Object{Id: st.RouterStateObjectID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected chains for chain %d: %w", selector, err)
	}
	return connectedChains, nil
}
