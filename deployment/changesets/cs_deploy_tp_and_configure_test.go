package changesets

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

type stubSigner struct{}

func (stubSigner) Sign([]byte) ([]string, error) { return nil, nil }
func (stubSigner) GetAddress() (string, error)   { return "0x1", nil }

// stubCoinMetadataPTBClient supplies the metadata read used by preconditions.
type stubCoinMetadataPTBClient struct {
	testutils.FakeSuiPTBClient
}

func (stubCoinMetadataPTBClient) GetCoinMetadata(context.Context, string) (models.CoinMetadataResponse, error) {
	return models.CoinMetadataResponse{Symbol: "CCIP-BnM"}, nil
}

const testBurnMintCoinType = "0xcoin::ccip_burn_mint::CCIP_BNM"

func deployTPEnv(t *testing.T, ds fdatastore.DataStore, selector uint64) cldf.Environment {
	t.Helper()
	return cldf.Environment{
		DataStore:        ds,
		OperationsBundle: mcmstest.Bundle(t),
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{Signer: stubSigner{}, Client: &stubCoinMetadataPTBClient{}},
		}),
	}
}

func storeSlowMCMS(t *testing.T, ds *fdatastore.MemoryDataStore, selector uint64) {
	t.Helper()
	require.NoError(t, deployment.StoreMCMSInAddressBook(cldf.NewMemoryAddressBook(), ds.Addresses(), selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xslow_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xslow_state",
			McmsRegistryObjectId:        "0xslow_registry",
			McmsAccountStateObjectId:    "0xslow_account",
			McmsAccountOwnerCapObjectId: "0xslow_owner_cap",
			TimelockObjectId:            "0xslow_timelock",
			McmsDeployerStateObjectId:   "0xslow_deployer",
		},
	}, deployment.MCMSInstanceSlow))
}

func storeFastMCMS(t *testing.T, ds *fdatastore.MemoryDataStore, selector uint64) {
	t.Helper()
	require.NoError(t, deployment.StoreMCMSInAddressBook(cldf.NewMemoryAddressBook(), ds.Addresses(), selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xfast_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xfast_state",
			McmsRegistryObjectId:        "0xfast_registry",
			McmsAccountStateObjectId:    "0xfast_account",
			McmsAccountOwnerCapObjectId: "0xfast_owner_cap",
			TimelockObjectId:            "0xfast_timelock",
			McmsDeployerStateObjectId:   "0xfast_deployer",
		},
	}, deployment.MCMSInstanceFastCurse))
}

func TestDeployTPAndConfigure_VerifyPreconditions_RequiresFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ds := fdatastore.NewMemoryDataStore()
	storeSlowMCMS(t, ds, selector)

	cs := DeployTPAndConfigure{}
	err := cs.VerifyPreconditions(deployTPEnv(t, ds.Seal(), selector), DeployTPAndConfigureConfig{
		SuiChainSelector: selector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeManaged},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fast MCMS package not deployed")
}

func TestDeployTPAndConfigure_VerifyPreconditions_SucceedsWithFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ds := fdatastore.NewMemoryDataStore()
	storeSlowMCMS(t, ds, selector)
	storeFastMCMS(t, ds, selector)

	cs := DeployTPAndConfigure{}
	err := cs.VerifyPreconditions(deployTPEnv(t, ds.Seal(), selector), DeployTPAndConfigureConfig{
		SuiChainSelector: selector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeBurnMint},
		BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
			CoinObjectTypeArg: testBurnMintCoinType,
		},
	})
	require.NoError(t, err)
}

func TestDeployTPAndConfigure_Apply_RequiresFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ds := fdatastore.NewMemoryDataStore()
	storeSlowMCMS(t, ds, selector)

	_, err := DeployTPAndConfigure{}.Apply(deployTPEnv(t, ds.Seal(), selector), DeployTPAndConfigureConfig{
		SuiChainSelector: selector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeLockRelease},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fast MCMS package not deployed")
}
