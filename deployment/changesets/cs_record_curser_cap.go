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
	return nil
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
	if err := deployment.SaveSuiAddress(ab, ds.Addresses(), cfg.SuiChainSelector, capID, tv); err != nil {
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
