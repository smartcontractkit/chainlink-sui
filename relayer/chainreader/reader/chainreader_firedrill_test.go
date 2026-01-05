//go:build integration

package reader

import (
	"context"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestChainReaderFiredrill(t *testing.T) {
	t.Skip("skipping ChainReaderFiredrill test, this is used as a sanity check only and not to be included in CI")

	log := logger.Test(t)
	rpcUrl := "https://sui-testnet-rpc.publicnode.com" // testutils.TestnetUrl

	offrampContractName := "OffRamp"
	offrampPackageId := "0xe2d83f15195acd57b798610d167dc241fcb30b5cc3808af497c33d97512b7970"

	onrampContractName := "OnRamp"
	onrampPackageId := "0x30e087460af8a8aacccbc218aa358cdcde8d43faf61ec0638d71108e276e2f1d"

	rmnRemoteContractName := "RMNRemote"
	ccipPackageAddress := "0x324c505732fadfa5ac2877cdca28a6be28910009e100de8e6e16eb33ed1218dc"

	t.Helper()
	ctx := context.Background()

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, _ := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	clientMaxConcurrentRequests := int64(10)
	if envClientMaxConcurrentRequests := os.Getenv("MAX_CONCURRENT_REQUESTS"); envClientMaxConcurrentRequests != "" {
		if parsed, err := strconv.Atoi(envClientMaxConcurrentRequests); err == nil {
			clientMaxConcurrentRequests = int64(parsed)
		}
	}
	relayerClient, clientErr := client.NewPTBClient(log, rpcUrl, nil, 120*time.Second, keystoreInstance, clientMaxConcurrentRequests, "WaitForLocalExecution")
	require.NoError(t, clientErr)

	type SourceChainConfig struct {
		Router                    string
		IsEnabled                 bool
		MinSeqNr                  uint64
		IsRMNVerificationDisabled bool
		OnRamp                    string
	}

	type SourceChainConfigSetEvent struct {
		SourceChainSelector uint64
		SourceChainConfig   SourceChainConfig
	}

	ccipObjectRefStatePointer := &codec.PointerTag{
		Module:        "state_object",
		PointerName:   "CCIPObjectRefPointer",
		FieldName:     "ccip_object_id",
		DerivationKey: "CCIPObjectRef",
	}

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
			offrampContractName: {
				Name:      "offramp",
				Functions: map[string]*config.ChainReaderFunction{},
				Events: map[string]*config.ChainReaderEvent{
					"SourceChainConfigSet": {
						Name:      "SourceChainConfigSet",
						EventType: "SourceChainConfigSet",
						EventSelector: client.EventSelector{
							Module: "offramp",
							Event:  "SourceChainConfigSet",
						},
						ExpectedEventType: &SourceChainConfigSetEvent{},
					},
				},
			},
			rmnRemoteContractName: {
				Name: "rmn_remote",
				Functions: map[string]*config.ChainReaderFunction{
					"GetVersionedConfig": {
						Name:          "get_versioned_config",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Name:       "object_ref_id",
								Type:       "object_id",
								PointerTag: ccipObjectRefStatePointer,
								Required:   true,
							},
						},
						ResultTupleToStruct: []string{"version", "config"},
					},
				},
				Events: map[string]*config.ChainReaderEvent{},
			},
			onrampContractName: {
				Name:      "onramp",
				Functions: map[string]*config.ChainReaderFunction{},
				Events: map[string]*config.ChainReaderEvent{
					"CCIPMessageSent": {
						Name:      "CCIPMessageSent",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventSelector{
							Module: "onramp",
							Event:  "CCIPMessageSent",
						},
						EventSelectorDefaultOffset: &client.EventId{
							TxDigest: "CpFQ8JsaHwTEuNLCfeJQopu3eM3ipViowkWmg23k4fNk",
							EventSeq: "0",
						},
					},
				},
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
		Name:    offrampContractName,
		Address: offrampPackageId,
	}, {
		Name:    rmnRemoteContractName,
		Address: ccipPackageAddress,
	}, {
		Name:    onrampContractName,
		Address: onrampPackageId,
	}})
	require.NoError(t, err)

	err = chainReader.Start(ctx)
	require.NoError(t, err)
	defer chainReader.Close()

	err = indexerInstance.Start(ctx)
	require.NoError(t, err)
	defer indexerInstance.Close()

	t.Run("sanity check for SourceChainConfigSet event", func(t *testing.T) {
		t.Skip("skipping SourceChainConfigSet event test")
		var seqType any
		events, err := chainReader.QueryKey(ctx, types.BoundContract{
			Name:    offrampContractName,
			Address: offrampPackageId,
		}, query.KeyFilter{Key: "SourceChainConfigSet"}, query.LimitAndSort{}, &seqType)

		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, events, "events")
	})

	t.Run("sanity check for GetVersionedConfig function", func(t *testing.T) {
		t.Skip("skipping GetVersionedConfig function test")
		var expectedVersionedConfig map[string]any
		err := chainReader.GetLatestValue(
			context.Background(),
			strings.Join([]string{ccipPackageAddress, rmnRemoteContractName, "GetVersionedConfig"}, "-"),
			primitives.Finalized,
			map[string]any{},
			&expectedVersionedConfig,
		)
		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, expectedVersionedConfig, "expectedVersionedConfig")
	})

	t.Run("sanity check for CCIPMessageSent event with offset override from configs", func(t *testing.T) {
		var ccipMessageSent any
		sequences, err := chainReader.QueryKey(ctx, types.BoundContract{
			Name:    onrampContractName,
			Address: onrampPackageId,
		}, query.KeyFilter{Key: "CCIPMessageSent"}, query.LimitAndSort{}, &ccipMessageSent)
		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, sequences, "sequences")
	})
}
