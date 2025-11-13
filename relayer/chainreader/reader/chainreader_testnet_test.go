//go:build testnet

package reader

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestChainReaderTestnet(t *testing.T) {
	log := logger.Test(t)
	rpcUrl := testutils.TestnetUrl

	tokenAdminRegistryContractName := "TokenAdminRegistry"
	tokenAdminRegistryPackageId := "0x9de8a33d158e26f0b51f199da8be1a22e9755510b705cfb88230b257187da733"

	burnMintTokenPoolContractName := "BurnMintTokenPool"
	burnMintTokenPoolPackageId := "0xfeff675b624e55da49f80fda3b676fe1ef5a957a8334cb675ca35de8918f612d"
	burnMintTokenPoolIdentifier := strings.Join([]string{burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_token"}, "-")

	t.Helper()
	ctx := context.Background()

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, _ := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	relayerClient, clientErr := client.NewPTBClient(log, rpcUrl, nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, clientErr)

	chainReaderConfig := config.ChainReaderConfig{
		IsLoopPlugin: false,
		EventsIndexer: config.EventsIndexerConfig{
			PollingInterval: 10 * time.Second,
			SyncTimeout:     10 * time.Second,
		},
		TransactionsIndexer: config.TransactionsIndexerConfig{
			PollingInterval: 10 * time.Second,
			SyncTimeout:     10 * time.Second,
		},
		Modules: map[string]*config.ChainReaderModule{
			tokenAdminRegistryContractName: {
				Name:      "token_admin_registry",
				Functions: map[string]*config.ChainReaderFunction{},
			},
			burnMintTokenPoolContractName: {
				Name: "burn_mint_token_pool",
				Functions: map[string]*config.ChainReaderFunction{
					"get_token": {
						Name:          "get_token",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:              "object_id",
								Name:              "state_pointer",
								GenericDependency: testutils.StringPointer("get_token_pool_state_type"),
								PointerTag: &codec.PointerTag{
									Module:        "burn_mint_token_pool",
									PointerName:   "BurnMintTokenPoolStatePointer",
									DerivationKey: "BurnMintTokenPoolState",
									FieldName:     "burn_mint_token_pool_object_id",
								},
								Required:  true,
								IsMutable: testutils.BoolPointer(true),
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{},
			},
		},
	}

	db := sqltest.NewNoOpDataSource()

	// Create the indexers
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		log,
		relayerClient,
		chainReaderConfig.TransactionsIndexer.PollingInterval,
		chainReaderConfig.TransactionsIndexer.SyncTimeout,
		// start without any configs, they will be set when ChainReader is initialized and gets a reference
		// to the transaction indexer to avoid having to reading ChainReader configs here as well
		map[string]*config.ChainReaderEvent{},
	)
	evIndexer := indexer.NewEventIndexer(
		db,
		log,
		relayerClient,
		// start without any selectors, they will be added during .Bind() calls on ChainReader
		[]*client.EventSelector{},
		chainReaderConfig.EventsIndexer.PollingInterval,
		chainReaderConfig.EventsIndexer.SyncTimeout,
	)
	indexerInstance := indexer.NewIndexer(
		log,
		evIndexer,
		txnIndexer,
	)

	// ChainReader in non-loop mode
	chainReader, err := NewChainReader(ctx, log, relayerClient, chainReaderConfig, db, indexerInstance)
	require.NoError(t, err)

	err = chainReader.Bind(context.Background(), []types.BoundContract{{
		Name:    tokenAdminRegistryContractName,
		Address: tokenAdminRegistryPackageId,
	}, {
		Name:    burnMintTokenPoolContractName,
		Address: burnMintTokenPoolPackageId,
	}})
	require.NoError(t, err)

	var retAddress string
	err = chainReader.GetLatestValue(ctx, burnMintTokenPoolIdentifier, primitives.Finalized, nil, &retAddress)
	require.NoError(t, err)
	require.Equal(t, len(retAddress), 66)

	var retAddress2 string
	err = chainReader.GetLatestValue(ctx, burnMintTokenPoolIdentifier, primitives.Finalized, nil, &retAddress2)
	require.NoError(t, err)
	require.Equal(t, len(retAddress2), 66)

	var retAddress3 string
	nilParams := make(map[string]any)
	err = chainReader.GetLatestValue(ctx, burnMintTokenPoolIdentifier, primitives.Finalized, &nilParams, &retAddress3)
	require.NoError(t, err)
	require.Equal(t, len(retAddress3), 66)
}
