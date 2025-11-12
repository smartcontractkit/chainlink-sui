//go:build testnet

package reader

import (
	"context"
	"fmt"
	"os"
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
			burnMintTokenPoolContractName: {
				Name: "burn_mint_token_pool",
				Functions: map[string]*config.ChainReaderFunction{
					"get_token": {
						Name:          "get_token",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type: "object_id",
								Name: "state_pointer",
								PointerTag: &codec.PointerTag{
									Module:        "burn_mint_token_pool",
									PointerName:   "BurnMintTokenPoolStatePointer",
									DerivationKey: "BurnMintTokenPoolState",
									FieldName:     "burn_mint_token_pool_object_id",
								},
								Required:    true,
								IsMutable:   testutils.BoolPointer(true),
								GenericType: testutils.StringPointer("0x406964d837bc10f14e26a40b4462c58cce0b9e57fc5265a6272dd2aba57e15a4::link::LINK"),
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{},
			},
		},
	}

	datastoreUrl := os.Getenv("TEST_DB_URL")
	if datastoreUrl == "" {
		t.Skip("Skipping persistent tests as TEST_DB_URL is not set in CI")
	}
	db := sqltest.NewDB(t, datastoreUrl)

	// attempt to connect
	_, err := db.Connx(ctx)
	require.NoError(t, err)

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
		Name:    burnMintTokenPoolContractName,
		Address: burnMintTokenPoolPackageId,
	}})
	require.NoError(t, err)

	var retAddress string
	err = chainReader.GetLatestValue(ctx, burnMintTokenPoolIdentifier, primitives.Finalized, map[string]any{}, &retAddress)
	require.NoError(t, err)

	fmt.Println(retAddress)
}
