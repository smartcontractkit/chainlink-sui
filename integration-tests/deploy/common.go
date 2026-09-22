//go:build integration

package deploy

import (
	"context"
	"fmt"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldfsui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/stretchr/testify/suite"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	bindutils "github.com/smartcontractkit/chainlink-sui/bindings/utils"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

type DeployTestSuite struct {
	suite.Suite
	lggr   logger.Logger
	signer bindutils.SuiSigner
	client client.SuiPTBClient
	env    cldf.Environment
	// ds is the mutable backing store of env.DataStore: the state loader is
	// datastore-only, so each phase's changeset output must be merged into it.
	ds *fdatastore.MemoryDataStore

	// Cached deployment addresses
	linkTokenPackageID     string
	linkTokenMetadataID    string
	linkTokenTreasuryCapID string
	ccipPackageID          string
	ccipObjectRef          string
	mcmsPackageID          string
	fastMcmsPackageID      string
	deployerAddr           string
}

func (s *DeployTestSuite) SetupSuite() {
	s.signer, s.client = testenv.SetupEnvironment(s.T())
	s.lggr = logger.Test(s.T())

	// Setup operation registry
	registry := cld_ops.NewOperationRegistry(opregistry.AllOperations...)

	bundle := cld_ops.NewBundle(
		func() context.Context { return s.T().Context() },
		s.lggr,
		cld_ops.NewMemoryReporter(),
		cld_ops.WithOperationRegistry(registry),
	)

	s.ds = fdatastore.NewMemoryDataStore()

	s.env = cldf.Environment{
		Name:              "test",
		Logger:            s.lggr,
		ExistingAddresses: cldf.NewMemoryAddressBook(),
		DataStore:         s.ds.Seal(),
		BlockChains: chain.NewBlockChains(
			map[uint64]chain.BlockChain{
				cselectors.SUI_LOCALNET.Selector: cldfsui.Chain{
					ChainMetadata: cldfsui.ChainMetadata{
						Selector: cselectors.SUI_LOCALNET.Selector,
					},
					Client: s.client,
					Signer: s.signer,
				},
			}),
		OperationsBundle: bundle,
	}
}

// mergeChangesetOutput merges a changeset output into the suite environment: the address
// book and the datastore. The state loader is datastore-only, so the datastore merge is
// what makes each phase's deployments visible to the next phase's state loads.
// (The sealed env.DataStore wraps the same underlying stores, so merging into s.ds is
// observed through the environment without resealing.)
func (s *DeployTestSuite) mergeChangesetOutput(out cldf.ChangesetOutput) {
	s.T().Helper()
	if out.AddressBook != nil {
		s.Require().NoError(s.env.ExistingAddresses.Merge(out.AddressBook), "failed to merge address book")
	}
	if out.DataStore != nil {
		s.Require().NoError(s.ds.Merge(out.DataStore.Seal()), "failed to merge datastore")
	}
}

// suiDatastoreRefs returns the chain's datastore address refs — the datastore-only
// replacement for address-book lookups.
func (s *DeployTestSuite) suiDatastoreRefs(selector uint64) []fdatastore.AddressRef {
	return s.ds.Addresses().Filter(fdatastore.AddressRefByChainSelector(selector))
}

// findUnusedManagedTokenMinterCapID finds the mint cap ID that wasn't consumed by the faucet.
// The faucet consumes its mint cap during initialization, so we check which mint caps still
// exist on-chain. The one that exists is the unused one (from ConfigureDeployerAsMinter).
func (s *DeployTestSuite) findUnusedManagedTokenMinterCapID() (string, error) {
	ctx := s.T().Context()
	var unusedMintCapID string

	// Find all mint caps and check which ones still exist on-chain
	for _, ref := range s.suiDatastoreRefs(SuiChainSelector) {
		if ref.Type == fdatastore.ContractType(deployment.SuiManagedTokenMinterCapID) &&
			ref.Labels.Contains(changesets.CCIPBnMSymbol) {
			// Check if this object still exists on-chain (not consumed/deleted by faucet)
			resp, err := bind.ReadObject(ctx, ref.Address, s.client)
			// If the object exists (no error and has data), it's the unused one
			if err == nil && resp != nil && resp.Data != nil {
				unusedMintCapID = ref.Address
				s.T().Logf("Found unused managed token minter cap ID: %s", ref.Address)
				break
			}
		}
	}

	if unusedMintCapID == "" {
		return "", fmt.Errorf("no unused managed token minter cap ID found (all may have been consumed)")
	}

	return unusedMintCapID, nil
}
