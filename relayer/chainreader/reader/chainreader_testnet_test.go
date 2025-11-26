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
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

func TestChainReaderTestnet(t *testing.T) {
	log := logger.Test(t)
	rpcUrl := testutils.TestnetUrl

	offrampContractName := "OffRamp"
	offrampPackageId := "0x50ff1c5a49f012f9360de2fa5065efe7185c22bbeee6254e68ef695b0b0d0f40"

	tokenAdminRegistryContractName := "TokenAdminRegistry"
	tokenAdminRegistryPackageId := "0x9de8a33d158e26f0b51f199da8be1a22e9755510b705cfb88230b257187da733"

	burnMintTokenPoolContractName := "BurnMintTokenPool"
	burnMintTokenPoolPackageId := "0xfeff675b624e55da49f80fda3b676fe1ef5a957a8334cb675ca35de8918f612d"
	burnMintTokenPoolIdentifier := strings.Join([]string{burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_token"}, "-")
	burnMintTokenPoolGetSupportedChainsIdentifier := strings.Join([]string{burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_supported_chains"}, "-")
	burnMintTokenPoolGetRemotePoolsIdentifier := strings.Join([]string{burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_remote_pools"}, "-")
	burnMintTokenPoolGetCurrentInboundRateLimiterStateIdentifier := strings.Join([]string{
		burnMintTokenPoolPackageId, burnMintTokenPoolContractName, "get_current_inbound_rate_limiter_state",
	}, "-")

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
			offrampContractName: {
				Name:      "offramp",
				Functions: map[string]*config.ChainReaderFunction{},
			},
			burnMintTokenPoolContractName: {
				Name: "token_pool",
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
		},
	}

	db := sqltest.NewDB(t, os.Getenv("TEST_DB_URL"))

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
		Name:    burnMintTokenPoolContractName,
		Address: burnMintTokenPoolPackageId,
	}, {
		Name:    tokenAdminRegistryContractName,
		Address: tokenAdminRegistryPackageId,
	}})
	require.NoError(t, err)

	t.Run("get_token_pool_state_type generic dependency for BurnMintTokenPool", func(t *testing.T) {
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
}
