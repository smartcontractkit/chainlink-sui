package reader

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/loop"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/stretchr/testify/require"
)

var (
	rmnDigestHeader       cciptypes.RMNDigestHeader
	commitLatestOCRConfig cciptypes.OCRConfigResponse
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

	// ccipObjectRefStatePointer := "_::state_object::CCIPObjectRefPointer::object_ref_id"
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
					// consts.MethodNameGetSourceChainConfig: {
					// 	Name:          "get_source_chain_config",
					// 	SignerAddress: accountAddress,
					// 	Params: []codec.SuiFunctionParam{
					// 		{
					// 			Name:       "object_ref_id",
					// 			Type:       "object_id",
					// 			PointerTag: &ccipObjectRefStatePointer,
					// 			Required:   true,
					// 		},
					// 		{
					// 			Name:       "off_ramp_state_id",
					// 			PointerTag: &offRampStatePointer,
					// 			Type:       "object_id",
					// 			Required:   true,
					// 		},
					// 		{
					// 			Name:     "sourceChainSelector",
					// 			Type:     "u64",
					// 			Required: true,
					// 		},
					// 	},
					// },

					consts.MethodNameOffRampLatestConfigDetails: {
						Name:          "latest_config_details",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "off_ramp_state_id",
								PointerTag: &offRampStatePointer,
								Type:       "object_id",
								Required:   true,
							},
							{
								Name:     "ocrPluginType",
								Type:     "u8",
								Required: true,
							},
						},
						ResultTupleToStruct: []string{"ocr_config"},
					},
				},
			},
		},
	}

	counterBinding := types.BoundContract{
		Name:    "OffRamp",
		Address: "0xcea836b828b80a8dc5eee57a7d61feb8e1ce96a6a88d3cf08e393e193b9aabb3", // Package ID of the deployed counter contract
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

	lrReader := loop.NewLoopChainReader(log, chainReader)

	err = lrReader.Bind(context.Background(), []types.BoundContract{counterBinding})
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
		"ocrPluginType": consts.PluginTypeCommit,
	}

	err = lrReader.GetLatestValue(ctx, "0xcea836b828b80a8dc5eee57a7d61feb8e1ce96a6a88d3cf08e393e193b9aabb3-OffRamp-OffRampLatestConfigDetails", "", params, &commitLatestOCRConfig)

	fmt.Println("ERROR: ", err)
}
