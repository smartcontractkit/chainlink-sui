package chainaccessor

import (
	"fmt"
	"sync"
)

// CCIP contract names. These mirror the contract names used by the CCIP plugin
// (consts.ContractNameOnRamp etc.) and are the keys passed to Sync /
// GetContractAddress.
const (
	ContractNameOnRamp       = "OnRamp"
	ContractNameOffRamp      = "OffRamp"
	ContractNameFeeQuoter    = "FeeQuoter"
	ContractNameNonceManager = "NonceManager"
	ContractNameRouter       = "Router"
)

// CCIP Move event type names (the struct name emitted on-chain).
const (
	EventNameCCIPMessageSent       = "CCIPMessageSent"
	EventNameCommitReportAccepted  = "CommitReportAccepted"
	EventNameExecutionStateChanged = "ExecutionStateChanged"
)

// Sui Move module names for each contract. FeeQuoter and NonceManager are
// modules within the shared CCIP package and read their state from the shared
// CCIPObjectRef object rather than a dedicated state object.
var contractModule = map[string]string{
	ContractNameOnRamp:       "onramp",
	ContractNameOffRamp:      "offramp",
	ContractNameFeeQuoter:    "fee_quoter",
	ContractNameNonceManager: "nonce_manager",
	ContractNameRouter:       "router",
}

// modulesUsingCCIPObjectRef are modules whose state lives in the shared
// CCIPObjectRef object (the CCIP package) rather than a dedicated state object.
var modulesUsingCCIPObjectRef = map[string]bool{
	"fee_quoter":    true,
	"nonce_manager": true,
}

// bindingsCache stores the resolved on-chain locations for the CCIP contracts
// bound via Sync. It is the Sui analog of the Solana accessor's pdaCache:
// expensive object-ID resolution happens once at Sync time and read methods
// only consult the cache afterwards.
type bindingsCache struct {
	mu sync.RWMutex
	// packageAddress maps contract name -> bound package address (the address
	// passed to Sync, i.e. the original package ID used as the event-store key).
	packageAddress map[string]string
	// stateObjectID maps Sui module name -> resolved state object ID
	// (e.g. "offramp" -> OffRampState object ID).
	stateObjectID map[string]string
	// ccipObjectRefID is the shared CCIPObjectRef object ID used by fee_quoter
	// and nonce_manager reads. Empty until a CCIP-package contract is synced.
	ccipObjectRefID string
}

func newBindingsCache() *bindingsCache {
	return &bindingsCache{
		packageAddress: make(map[string]string),
		stateObjectID:  make(map[string]string),
	}
}

func (b *bindingsCache) setPackageAddress(contractName, address string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.packageAddress[contractName] = address
}

func (b *bindingsCache) getPackageAddress(contractName string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	addr, ok := b.packageAddress[contractName]
	if !ok {
		return "", fmt.Errorf("%w: contract %q is not bound (call Sync first)", ErrNotBound, contractName)
	}
	return addr, nil
}

func (b *bindingsCache) setStateObjectID(module, objectID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.stateObjectID[module] = objectID
}

func (b *bindingsCache) getStateObjectID(module string) (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	id, ok := b.stateObjectID[module]
	if !ok {
		return "", fmt.Errorf("%w: state object for module %q is not resolved (call Sync first)", ErrNotBound, module)
	}
	return id, nil
}

func (b *bindingsCache) setCCIPObjectRefID(objectID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ccipObjectRefID = objectID
}

func (b *bindingsCache) getCCIPObjectRefID() (string, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.ccipObjectRefID == "" {
		return "", fmt.Errorf("%w: CCIPObjectRef is not resolved (Sync a CCIP-package contract first)", ErrNotBound)
	}
	return b.ccipObjectRefID, nil
}
