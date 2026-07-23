//go:build testnet

package reader

import (
	"context"
	"fmt"
	"math/rand"
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
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestChainReaderTestnet(t *testing.T) {
	log := logger.Test(t)

	burnMintTokenPoolContractName := "BurnMintTokenPool"
	burnMintTokenPoolPackageId := "0xfeff675b624e55da49f80fda3b676fe1ef5a957a8334cb675ca35de8918f612d"
	burnMintTokenPoolIdentifier := strings.Join([]string{burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_token"}, "-")
	burnMintTokenPoolGetSupportedChainsIdentifier := strings.Join([]string{burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_supported_chains"}, "-")
	burnMintTokenPoolGetRemotePoolsIdentifier := strings.Join([]string{burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_remote_pools"}, "-")
	burnMintTokenPoolGetCurrentInboundRateLimiterStateIdentifier := strings.Join([]string{
		burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_current_inbound_rate_limiter_state",
	}, "-")
	onRampContractName := "OnRamp"
	onRampPackageId := "0x30e087460af8a8aacccbc218aa358cdcde8d43faf61ec0638d71108e276e2f1d"

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

	grpcTarget := os.Getenv("GRPC_TARGET")
	if grpcTarget == "" {
		t.Fatal("evn value for GRPC_TARGET is not set")
	}
	grpcToken := os.Getenv("GRPC_TOKEN")
	if grpcToken == "" {
		t.Fatal("evn value for GRPC_TOKEN is not set")
	}

	testCfg := client.PTBClientConfig{
		GrpcTarget:            grpcTarget,
		GrpcToken:             grpcToken,
		TransactionTimeout:    10 * time.Second,
		MaxConcurrentRequests: clientMaxConcurrentRequests,
		KeystoreService:       keystoreInstance,
	}

	relayerClient, clientErr := client.NewPTBClient(log, testCfg)
	require.NoError(t, clientErr)

	chainReaderConfig := config.ChainReaderConfig{
		IsLoopPlugin: false,
		EventsIndexer: config.EventsIndexerConfig{
			PollingInterval: 1 * time.Second,
			SyncTimeout:     60 * time.Second,
		},
		TransactionsIndexer: config.TransactionsIndexerConfig{
			PollingInterval: 1 * time.Second,
			SyncTimeout:     60 * time.Second,
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
					"type_and_version": {
						Name:          "type_and_version",
						SignerAddress: accountAddress,
						Params:        []codec.SuiFunctionParam{},
					},
					"get_supported_chains": {
						Name:          "get_supported_chains",
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
					"get_remote_pools": {
						Name:          "get_remote_pools",
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
							{
								Type:     "u64",
								Name:     "remote_chain_selector",
								Required: true,
							},
						},
					},
					"get_current_inbound_rate_limiter_state": {
						Name:          "get_current_inbound_rate_limiter_state",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{
								Type:         "object_id",
								Name:         "clock",
								Required:     false,
								DefaultValue: "0x06",
							},
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
							{
								Type:     "u64",
								Name:     "remote_chain_selector",
								Required: true,
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{
					"released_or_minted": {
						Name:      "released_or_minted",
						EventType: "ReleasedOrMinted",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "token_pool",
							Event:  "ReleasedOrMinted",
						},
					},
				},
			},
			onRampContractName: {
				Events: map[string]*config.ChainReaderEvent{
					"ccip_message_sent": {
						Name:      "onramp",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventFilterByMoveEventModule{
							Module: "onramp",
							Event:  "CCIPMessageSent",
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
	dbStore := database.NewDBStore(db, log)
	require.NoError(t, dbStore.EnsureSchema(ctx))

	// Create the indexers
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		log,
		// start without any configs, they will be set when ChainReader is initialized and gets a reference
		// to the transaction indexer to avoid having to reading ChainReader configs here as well
		map[string]*config.ChainReaderEvent{},
	)

	evIndexer := indexer.NewEventIndexer(
		db,
		log,
		// start without any selectors, they will be added during .Bind() calls on ChainReader
		[]*client.EventSelector{},
	)

	startingCheckpointSequence := uint64(362867769)
	chainPoller := indexer.NewChainPoller(
		relayerClient,
		log,
		config.ChainPollerConfig{
			PollingInterval:         1 * time.Second,
			SyncTimeout:             60 * time.Second,
			BackfillCheckpointCount: nil,
			StartCheckpointSequence: &startingCheckpointSequence,
			ChannelBufferSize:       16,
		},
		evIndexer.GetEventSelectors,
	)

	indexerInstance := indexer.NewIndexerFromComponents(
		log,
		chainPoller,
		evIndexer,
		txnIndexer,
	)

	indexerInstance.Start(ctx)

	// ChainReader in non-loop mode
	chainReader, err := NewChainReader(ctx, log, relayerClient, chainReaderConfig, db, indexerInstance, nil)
	require.NoError(t, err)

	err = chainReader.Bind(context.Background(), []types.BoundContract{{
		Name:    burnMintTokenPoolContractName,
		Address: burnMintTokenPoolPackageId,
	}})
	require.NoError(t, err)

	err = chainReader.Bind(context.Background(), []types.BoundContract{{
		Name:    onRampContractName,
		Address: onRampPackageId,
	}})
	require.NoError(t, err)

	t.Run("get_token_pool_state_type generic dependency for BurnMintTokenPool", func(t *testing.T) {
		t.Skip("...")
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

		var retAddress4 string
		var params map[string]any
		err = chainReader.GetLatestValue(ctx, burnMintTokenPoolIdentifier, primitives.Finalized, &params, &retAddress4)
		require.NoError(t, err)
		require.Equal(t, len(retAddress4), 66)

		var retSupportedChains []uint64
		err = chainReader.GetLatestValue(ctx, burnMintTokenPoolGetSupportedChainsIdentifier, primitives.Finalized, nil, &retSupportedChains)
		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, retSupportedChains, "retSupportedChains")

		var retRemotePools any
		params = map[string]any{
			"remote_chain_selector": uint64(14767482510784806043),
		}
		err = chainReader.GetLatestValue(ctx, burnMintTokenPoolGetRemotePoolsIdentifier, primitives.Finalized, &params, &retRemotePools)
		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, retRemotePools, "retRemotePools")

		var retCurrentInboundRateLimiterState any
		params = map[string]any{
			"remote_chain_selector": uint64(14767482510784806043),
		}
		err = chainReader.GetLatestValue(ctx, burnMintTokenPoolGetCurrentInboundRateLimiterStateIdentifier, primitives.Finalized, &params, &retCurrentInboundRateLimiterState)
		require.NoError(t, err)
		testutils.PrettyPrintDebug(log, retCurrentInboundRateLimiterState, "retCurrentInboundRateLimiterState")
	})

	t.Run("client load test GetObjectId", func(t *testing.T) {
		t.Skip("...")
		numRequests := 10
		if envNumRequests := os.Getenv("NUM_REQUESTS"); envNumRequests != "" {
			if parsed, err := strconv.Atoi(envNumRequests); err == nil {
				numRequests = parsed
			}
		}
		errChan := make(chan error, numRequests)

		for i := range numRequests {
			go func() {
				// Random sleep to simulate real-world load
				sleepDuration := time.Duration(100+rand.Intn(1500)) * time.Millisecond
				time.Sleep(sleepDuration)

				start := time.Now()
				response, err := relayerClient.ReadObjectId(ctx, burnMintTokenPoolPackageId)
				if err != nil {
					elapsed := time.Since(start)
					log.Infow("Request completed", "request", i, "elapsed", elapsed)

					errChan <- fmt.Errorf("failed to get value at request %d: %w", i, err)
					return
				}

				elapsed := time.Since(start)
				log.Infow("Request completed", "request", i, "elapsed", elapsed, "response", response.ObjectId)

				errChan <- nil
			}()
		}

		// Collect all results
		errorCount := 0
		processedCount := 0
		for range numRequests {
			err := <-errChan
			if err != nil {
				errorCount++
				log.Errorw("ReadObjectId Test Error", "error", err)
			}
			processedCount++
		}

		log.Infof("Completed %d requests, %d errors", processedCount, errorCount)
	})

	t.Run("chainreader load test GetLatestValue", func(t *testing.T) {
		t.Skip("...")
		numRequests := 10
		if envNumRequests := os.Getenv("NUM_REQUESTS"); envNumRequests != "" {
			if parsed, err := strconv.Atoi(envNumRequests); err == nil {
				numRequests = parsed
			}
		}
		errChan := make(chan error, numRequests)

		for i := range numRequests {
			go func() {
				var retAddress string

				// Random sleep to simulate real-world load
				sleepDuration := time.Duration(100+rand.Intn(1500)) * time.Millisecond
				time.Sleep(sleepDuration)

				start := time.Now()
				err := chainReader.GetLatestValue(ctx, burnMintTokenPoolIdentifier, primitives.Finalized, nil, &retAddress)
				if err != nil {
					elapsed := time.Since(start)
					log.Infow("Request completed", "request", i, "elapsed", elapsed)

					errChan <- fmt.Errorf("failed to get value at request %d: %w", i, err)
					return
				}

				elapsed := time.Since(start)
				log.Infow("Request completed", "request", i, "elapsed", elapsed)

				errChan <- nil
			}()
		}

		// Collect all results
		errorCount := 0
		processedCount := 0
		for range numRequests {
			err := <-errChan
			if err != nil {
				errorCount++
				log.Errorw("Error", "error", err)
			}
			processedCount++
		}

		log.Infof("Completed %d requests, %d errors", processedCount, errorCount)
	})

	t.Run("token pool events", func(t *testing.T) {
		t.Skip("...")
		var retReleasedOrMinted map[string]any

		sequences, err := chainReader.QueryKey(ctx, types.BoundContract{
			Name:    burnMintTokenPoolContractName,
			Address: burnMintTokenPoolPackageId,
		}, query.KeyFilter{
			Key: "released_or_minted",
		}, query.LimitAndSort{
			Limit: query.Limit{
				Count: 100,
			},
		}, &retReleasedOrMinted)

		testutils.PrettyPrintDebug(log, sequences, "sequences")
		require.NoError(t, err)
	})

	t.Run("onramp events", func(t *testing.T) {
		require.Eventually(t, func() bool {
			events, err := chainReader.QueryKey(ctx, types.BoundContract{
				Name:    onRampContractName,
				Address: onRampPackageId,
			}, query.KeyFilter{
				Key: "ccip_message_sent",
			}, query.LimitAndSort{
				Limit: query.Limit{
					Count: 10,
				},
			}, nil)

			isFound := err == nil && len(events) > 0

			return isFound && false
		}, 300*time.Second, 2*time.Second, "Expected to find at least one ccip_message_sent event")
	})
}
