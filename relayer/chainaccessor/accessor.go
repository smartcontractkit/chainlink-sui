// Package chainaccessor provides a native Sui implementation of the
// ccipocr3.ChainAccessor interface.
//
// Unlike the DefaultAccessor (which wraps an existing ContractReader /
// ContractWriter), this accessor talks directly to the Sui chain: on-chain
// state is read via devInspect Move view calls and events are read from the
// indexer-populated event store. It mirrors the structure of chainlink-solana's
// native SolanaAccessor.
//
// This package currently implements the high-value CCIP read methods. Methods
// that are not yet implemented natively return ErrNotImplemented so that
// *SuiAccessor still satisfies the full ccipocr3.ChainAccessor interface.
package chainaccessor

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/common"
)

// SuiAccessor is a native ccipocr3.ChainAccessor implementation for Sui.
type SuiAccessor struct {
	lggr          logger.Logger
	chainSelector ccipocr3.ChainSelector

	// client performs devInspect Move view calls and resolves package / object
	// IDs (GetLatestPackageId, GetParentObjectID, GetCCIPPackageID).
	client *client.PTBClient
	// dbStore holds events populated by the indexer; event reads query it.
	dbStore *database.DBStore
	// indexer is used at Sync time to register the event selectors that should
	// be polled and persisted to dbStore.
	indexer indexer.IndexerApi

	bindings *bindingsCache
}

// compile-time assertion that *SuiAccessor satisfies the full interface.
var _ ccipocr3.ChainAccessor = (*SuiAccessor)(nil)

// NewSuiAccessor constructs a native Sui ChainAccessor. Contracts must be bound
// via Sync before reads against them will succeed.
func NewSuiAccessor(
	lggr logger.Logger,
	chainSelector ccipocr3.ChainSelector,
	ptbClient *client.PTBClient,
	dbStore *database.DBStore,
	idx indexer.IndexerApi,
) (*SuiAccessor, error) {
	if ptbClient == nil {
		return nil, fmt.Errorf("ptbClient must not be nil")
	}
	if dbStore == nil {
		return nil, fmt.Errorf("dbStore must not be nil")
	}
	if idx == nil {
		return nil, fmt.Errorf("indexer must not be nil")
	}

	return &SuiAccessor{
		lggr:          logger.Named(lggr, "SuiAccessor"),
		chainSelector: chainSelector,
		client:        ptbClient,
		dbStore:       dbStore,
		indexer:       idx,
		bindings:      newBindingsCache(),
	}, nil
}

// GetContractAddress returns the bound package address for the provided contract
// name as raw bytes.
//
// Contract: N/A
func (a *SuiAccessor) GetContractAddress(contractName string) ([]byte, error) {
	addr, err := a.bindings.getPackageAddress(contractName)
	if err != nil {
		return nil, err
	}
	return suiAddressToBytes(addr)
}

// Sync binds a contract by name and address, resolves and caches the on-chain
// state object IDs needed by subsequent reads, and registers the contract's
// event selectors with the indexer.
//
// Contract: Many
func (a *SuiAccessor) Sync(ctx context.Context, contractName string, contractAddress ccipocr3.UnknownAddress) error {
	packageAddr := suiAddressFromBytes(contractAddress)
	if !bind.IsSuiAddress(packageAddr) {
		return fmt.Errorf("invalid sui address for contract %q: %q", contractName, packageAddr)
	}

	module, ok := contractModule[contractName]
	if !ok {
		return fmt.Errorf("unsupported contract name %q", contractName)
	}

	a.lggr.Debugw("Sync: binding contract", "contract", contractName, "module", module, "address", packageAddr)
	a.bindings.setPackageAddress(contractName, packageAddr)

	// Resolve the state object IDs this contract's reads will need.
	if modulesUsingCCIPObjectRef[module] {
		// fee_quoter / nonce_manager read the shared CCIPObjectRef in the CCIP
		// package. The bound address is the CCIP package address.
		refID, err := a.resolveStateObjectID(ctx, packageAddr, "state_object")
		if err != nil {
			return fmt.Errorf("failed to resolve CCIPObjectRef for %q: %w", contractName, err)
		}
		a.bindings.setCCIPObjectRefID(refID)
	} else {
		stateID, err := a.resolveStateObjectID(ctx, packageAddr, module)
		if err != nil {
			return fmt.Errorf("failed to resolve state object for %q: %w", contractName, err)
		}
		a.bindings.setStateObjectID(module, stateID)

		// OffRamp also needs the shared CCIPObjectRef (e.g. get_source_chain_config).
		if module == "offramp" {
			ccipPkgID, err := a.client.GetCCIPPackageID(ctx, packageAddr)
			if err != nil {
				return fmt.Errorf("failed to resolve CCIP package id from offramp: %w", err)
			}
			refID, err := a.resolveStateObjectID(ctx, ccipPkgID, "state_object")
			if err != nil {
				return fmt.Errorf("failed to resolve CCIPObjectRef from offramp: %w", err)
			}
			a.bindings.setCCIPObjectRefID(refID)
		}
	}

	// Register event selectors so the indexer polls and persists these events.
	if err := a.registerEventSelectors(ctx, contractName, module, packageAddr); err != nil {
		return fmt.Errorf("failed to register event selectors for %q: %w", contractName, err)
	}

	return nil
}

// resolveStateObjectID reproduces the ChainReader's pointer-tag resolution
// natively: it fetches the pointer object's parent ID and derives the state
// object's ID from it. The pointer recipe is taken from the shared
// common.PointerConfigs registry (the single source of truth).
func (a *SuiAccessor) resolveStateObjectID(ctx context.Context, packageAddr, module string) (string, error) {
	configs := common.GetPointerConfigsByContract(module)
	if len(configs) == 0 {
		return "", fmt.Errorf("no pointer config for module %q", module)
	}
	pointer := configs[0]

	parentObjectID, err := a.client.GetParentObjectID(ctx, packageAddr, pointer.Module, pointer.Pointer)
	if err != nil {
		return "", fmt.Errorf("failed to get parent object id (module=%s pointer=%s): %w", pointer.Module, pointer.Pointer, err)
	}

	derivationKey := common.GetStateObjectNameByModule(pointer.Module)
	if derivationKey == "" {
		return "", fmt.Errorf("no state object name for module %q", pointer.Module)
	}

	stateObjectID, err := client.DeriveObjectIDWithVectorU8Key(parentObjectID, []byte(derivationKey))
	if err != nil {
		return "", fmt.Errorf("failed to derive state object id (key=%s): %w", derivationKey, err)
	}

	a.lggr.Debugw("resolved state object", "module", pointer.Module, "pointer", pointer.Pointer, "stateObjectID", stateObjectID)
	return stateObjectID, nil
}

// registerEventSelectors registers the CCIP event selectors emitted by the given
// contract with the events indexer so they are polled and persisted.
func (a *SuiAccessor) registerEventSelectors(ctx context.Context, contractName, module, packageAddr string) error {
	var events []string
	switch contractName {
	case ContractNameOnRamp:
		events = []string{EventNameCCIPMessageSent}
	case ContractNameOffRamp:
		events = []string{EventNameCommitReportAccepted, EventNameExecutionStateChanged}
	default:
		return nil
	}

	eventIndexer := a.indexer.GetEventIndexer()
	for _, ev := range events {
		selector := &client.EventSelector{
			Package: packageAddr,
			Module:  module,
			Event:   ev,
		}
		if err := eventIndexer.AddEventSelector(ctx, selector); err != nil {
			return fmt.Errorf("failed to add event selector %s::%s::%s: %w", packageAddr, module, ev, err)
		}
	}
	return nil
}

// GetAllConfigsLegacy is not yet implemented natively; it would compose the
// individual view calls (offramp/onramp/fee_quoter/rmn) into a snapshot.
//
// Contract: Many
func (a *SuiAccessor) GetAllConfigsLegacy(
	ctx context.Context,
	destChainSelector ccipocr3.ChainSelector,
	sourceChainSelectors []ccipocr3.ChainSelector,
) (ccipocr3.ChainConfigSnapshot, map[ccipocr3.ChainSelector]ccipocr3.SourceChainConfig, error) {
	return ccipocr3.ChainConfigSnapshot{}, nil, ErrNotImplemented
}

// GetChainFeeComponents is not yet implemented natively (would be sourced from
// the txm / fee estimator).
//
// Contract: N/A
func (a *SuiAccessor) GetChainFeeComponents(ctx context.Context) (ccipocr3.ChainFeeComponents, error) {
	return ccipocr3.ChainFeeComponents{}, ErrNotImplemented
}
