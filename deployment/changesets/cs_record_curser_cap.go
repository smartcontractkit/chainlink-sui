package changesets

import (
	"errors"
	"fmt"
	"strings"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type RecordCurserCapConfig struct {
	SuiChainSelector  uint64 `yaml:"suiChainSelector"`
	TxDigest          string `yaml:"txDigest,omitempty"`
	CurserCapObjectId string `yaml:"curserCapObjectId,omitempty"`
	// ReplaceExisting allows this changeset to take a datastore key already held by a
	// *different* cap object. Without it that is an error, raised before anything is written;
	// re-recording the cap already under the key is a no-op and needs no flag.
	ReplaceExisting bool `yaml:"replaceExisting"`
}

var _ cldf.ChangeSetV2[RecordCurserCapConfig] = RecordCurserCap{}

type RecordCurserCap struct{}

func (c RecordCurserCap) VerifyPreconditions(e cldf.Environment, cfg RecordCurserCapConfig) error {
	if strings.TrimSpace(cfg.TxDigest) == "" && strings.TrimSpace(cfg.CurserCapObjectId) == "" {
		return errors.New("either txDigest or curserCapObjectId is required")
	}
	if strings.TrimSpace(cfg.TxDigest) != "" {
		if _, ok := e.BlockChains.SuiChains()[cfg.SuiChainSelector]; !ok {
			return fmt.Errorf("no Sui chain client for selector %d (required to resolve txDigest)", cfg.SuiChainSelector)
		}
	}
	return deployment.ValidateNoDatastoreConflicts(e, cfg.SuiChainSelector, cfg.ReplaceExisting,
		func() ([]deployment.PlannedRef, error) {
			// The cap ID is known upfront (import path), so the plan carries it: re-recording
			// the same cap is then a no-op rather than a conflict.
			capID, err := resolveCurserCapIDForRecord(e, cfg)
			if err != nil {
				return nil, err
			}
			return []deployment.PlannedRef{
				{Type: deployment.SuiCurserCapObjectIDType, Qualifier: deployment.ChainSingletonQualifier, Address: capID},
			}, nil
		})
}

func (c RecordCurserCap) Apply(e cldf.Environment, cfg RecordCurserCapConfig) (cldf.ChangesetOutput, error) {
	capID, err := resolveCurserCapIDForRecord(e, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("load onchain state: %w", err)
	}
	chainState, ok := state[cfg.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for selector %d", cfg.SuiChainSelector)
	}
	if registered := chainState.CurserCapObjectId; registered != "" && registered != capID {
		return cldf.ChangesetOutput{}, fmt.Errorf(
			"curserCapObjectId %q conflicts with registered CurserCap %q",
			capID, registered,
		)
	}

	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	tv := cldf.NewTypeAndVersion(deployment.SuiCurserCapObjectIDType, deployment.Version1_0_0)
	// chain singleton; empty qualifier
	if err := deployment.SaveSuiAddress(ab, ds.Addresses(), cfg.SuiChainSelector, capID, tv, deployment.ChainSingletonQualifier); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("save CurserCap to address book: %w", err)
	}

	return cldf.ChangesetOutput{AddressBook: ab, DataStore: ds}, nil
}

func resolveCurserCapIDForRecord(e cldf.Environment, cfg RecordCurserCapConfig) (string, error) {
	if id := strings.TrimSpace(cfg.CurserCapObjectId); id != "" {
		return id, nil
	}
	suiChain, ok := e.BlockChains.SuiChains()[cfg.SuiChainSelector]
	if !ok {
		return "", fmt.Errorf("no Sui chain client for selector %d", cfg.SuiChainSelector)
	}
	return utils.FindCurserCapObjectIDFromTx(e.GetContext(), suiChain.Client, strings.TrimSpace(cfg.TxDigest))
}
