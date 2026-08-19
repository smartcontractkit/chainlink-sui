package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

// DeployMaliciousTokenPoolConfig configures the TEST-only malicious token pool
// deploy used by the ReleaseOrMintParams ownership guard E2E smoke test. The pool
// is registered with ReleaseOrMintParams as its release_or_mint_params; that is the
// attacker-controlled field the offchain guard defends.
type DeployMaliciousTokenPoolConfig struct {
	SuiChainSelector       uint64
	McmsOwner              string
	CoinObjectTypeArg      string
	CCIPObjectRefObjectId  string
	CoinMetadataObjectId   string
	TreasuryCapObjectId    string
	TokenPoolAdministrator string
	ReleaseOrMintParams    []string
}

var _ cldf.ChangeSetV2[DeployMaliciousTokenPoolConfig] = DeployMaliciousTokenPool{}

// DeployMaliciousTokenPool deploys and initializes the TEST-only malicious token
// pool used by the ReleaseOrMintParams ownership guard E2E smoke test.
type DeployMaliciousTokenPool struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployMaliciousTokenPool) Apply(e cldf.Environment, config DeployMaliciousTokenPoolConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	seqReports := make([]operations.Report[any, any], 0)

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	deployOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.DeployCCIPMaliciousTokenPoolOp, deps, ccipops.DeployMaliciousTokenPoolInput{
		CCIPPackageId:          state[config.SuiChainSelector].CCIPAddress,
		MCMSAddress:            state[config.SuiChainSelector].MCMSPackageID,
		FastMcmsAddress:        state[config.SuiChainSelector].FastCurseMCMSPackageID,
		MCMSOwnerAddress:       config.McmsOwner,
		CoinObjectTypeArg:      config.CoinObjectTypeArg,
		CCIPObjectRefObjectId:  config.CCIPObjectRefObjectId,
		CoinMetadataObjectId:   config.CoinMetadataObjectId,
		TreasuryCapObjectId:    config.TreasuryCapObjectId,
		TokenPoolAdministrator: config.TokenPoolAdministrator,
		ReleaseOrMintParams:    config.ReleaseOrMintParams,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy malicious token pool for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, deployOp.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployMaliciousTokenPool) VerifyPreconditions(e cldf.Environment, config DeployMaliciousTokenPoolConfig) error {
	return nil
}
