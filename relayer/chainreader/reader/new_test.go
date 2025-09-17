package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/stretchr/testify/require"
)

func TestGetLatestValue(t *testing.T) {
	log := logger.Test(t)
	// initialize CR

	t.Helper()
	ctx := context.Background()

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, _ := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	relayerClient, clientErr := client.NewPTBClient(log, "", nil, 10*time.Second, keystoreInstance, 5, "WaitForLocalExecution")
	require.NoError(t, clientErr)

	ccipObjectRefStatePointer := "_::state_object::CCIPObjectRefPointer::object_ref_id"
	offRampStatePointer := "_::offramp::OffRampStatePointer::off_ramp_state_id"

	chainReaderConfig := config.ChainReaderConfig{
		IsLoopPlugin: true,
		EventsIndexer: config.EventsIndexerConfig{
			PollingInterval: 10 * time.Second,
			SyncTimeout:     10 * time.Second,
		},
		TransactionsIndexer: config.TransactionsIndexerConfig{
			PollingInterval: 10 * time.Second,
			SyncTimeout:     10 * time.Second,
		},
		Modules: map[string]*config.ChainReaderModule{
			"OffRamp": {
				Name: "offramp",
				Functions: map[string]*config.ChainReaderFunction{
					consts.MethodNameGetSourceChainConfig: {
						Name:          "get_source_chain_config",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "object_ref_id",
								Type:       "object_id",
								PointerTag: &ccipObjectRefStatePointer,
								Required:   true,
							},
							{
								Name:       "off_ramp_state_id",
								PointerTag: &offRampStatePointer,
								Type:       "object_id",
								Required:   true,
							},
							{
								Name:     "sourceChainSelector",
								Type:     "u64",
								Required: true,
							},
						},
					},
				},
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "OffRamp",
		Address: "0xab21eb88ffdd8ba2eabed19dbfdf0b2f94da5edd34441e6a9da6c0850c3be284", // Package ID of the deployed counter contract
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

	chainReader, err := NewChainReader(ctx, log, relayerClient, chainReaderConfig, db, indexerInstance)
	require.NoError(t, err)

	err = chainReader.Bind(context.Background(), []types.BoundContract{counterBinding})
	require.NoError(t, err)

	log.Debugw("ChainReader setup complete")

	go func() {
		err = chainReader.Start(ctx)
		require.NoError(t, err)
		log.Debugw("ChainReader started")
	}()
	go func() {
		err = indexerInstance.Start(ctx)
		require.NoError(t, err)
		log.Debugw("Indexers started")
	}()

	params := map[string]any{
		"sourceChainSelector": "16015286601757825753",
	}

	paramBytes, _ := json.Marshal(params)

	err = chainReader.GetLatestValue(ctx, "0xab21eb88ffdd8ba2eabed19dbfdf0b2f94da5edd34441e6a9da6c0850c3be284-OffRamp-GetSourceChainConfig", "", &paramBytes, "")

	fmt.Println("ERROR: ", err)
}
