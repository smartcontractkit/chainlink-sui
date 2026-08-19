//go:build integration

// Package chainreader_test contains a controlled reproduction of the EVM->Sui commit blocker
// observed in CI (chainlink-sui PR #446 / chainlink PR #23036).
//
// Background / root cause this test targets
// -----------------------------------------
// In the failing EVM->Sui runs the message is sent but never committed on Sui. The Sui commit
// plugin rejects every report in ShouldAccept/ShouldTransmit with "config digest mismatch":
//
//	commit/report.go:222 "my config digest doesn't match offramp's config digest, not accepting report"
//	    myConfigDigest:      000ab1886...            (the DON's running config)
//	    offRampConfigDigest: 000000000...0000        (what it read back from the OffRamp)
//
// The OffRamp OCR config IS set on-chain (the OCR ContractConfigTracker reads it fine from the
// ConfigSet event, and individual latest_config_details reads always decode the correct non-zero
// digest). The problem is that the CCIP config poller's batch snapshot of the Sui OffRamp comes
// back all-zeros ~60% of the time. Those zero snapshots are NOT stale cached values (every
// SUCCESSFUL read returns the correct digest) — they are the result of the underlying Sui read
// being cancelled ("rpc error: code = Canceled") when the config poller's batch budget is
// exceeded. On a per-read failure the CCIP config processor leaves the OCR config at its zero
// value, so the commit plugin sees a zero digest and refuses to commit.
//
// In short: transient read cancellation -> zero-value OCR config -> persistent config digest
// mismatch -> EVM->Sui never commits.
//
// What this test does
// -------------------
// It uses the counter test contract's get_ocr_config view, which returns the exact OCRConfig shape
// (config_info{config_digest, big_f, n, is_signature_verification_enabled}, signers, transmitters)
// that the OffRamp's latest_config_details returns, so it exercises the same read + decode path the
// commit plugin depends on. It then:
//
//  1. (Baseline) proves a healthy read returns the correct non-zero config digest.
//  2. (Reproduce) proves that with the read cache DISABLED (production default today), a transient
//     context cancellation makes the read fail — which upstream becomes the zero-value OCR config
//     that blocks EVM->Sui commits.
//  3. (Fix) proves that with the read-result cache ENABLED, the same transient cancellation is
//     absorbed: the last-known-good config digest is still returned, so the commit plugin would
//     see a matching digest.
//
// It is intentionally driven by context cancellation (deterministic, fast) rather than by trying to
// reproduce real RPC slowness, so we can iterate on the resilience fix quickly.
package chainreader_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	aptosCRConfig "github.com/smartcontractkit/chainlink-common/pkg/types/aptos"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// ocrConfigInfo mirrors the Move offramp OCRConfig.config_info struct. Note the JSON tag `f`: the
// Sui node renders the field as `big_f`, and the chainreader's ResultFieldRenames (below) renames it
// to `f` before decoding. config_digest is the field that decodes to zero on a failed read.
type ocrConfigInfo struct {
	ConfigDigest                   []byte `json:"config_digest"`
	F                              uint64 `json:"f"`
	N                              uint64 `json:"n"`
	IsSignatureVerificationEnabled bool   `json:"is_signature_verification_enabled"`
}

type ocrConfig struct {
	ConfigInfo   ocrConfigInfo `json:"config_info"`
	Signers      [][]byte      `json:"signers"`
	Transmitters [][]byte      `json:"transmitters"`
}

type ocrConfigWrapped struct {
	OCRConfig ocrConfig `json:"OCRConfig"`
}

func isZeroDigest(b []byte) bool {
	for _, x := range b {
		if x != 0 {
			return false
		}
	}
	return true
}

//nolint:paralleltest // spins up a dedicated local Sui node; must not run in parallel with others.
func TestSuiConfigReadResilience(t *testing.T) {
	log := logger.Test(t)
	ctx := context.Background()

	datastoreURL := os.Getenv("TEST_DB_URL")
	if datastoreURL == "" {
		t.Skip("Skipping: TEST_DB_URL is not set")
	}

	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	testutils.CleanupTestContracts()
	t.Cleanup(func() {
		testutils.CleanupTestContracts()
		if cmd.Process != nil {
			if perr := cmd.Process.Kill(); perr != nil {
				t.Logf("failed to kill Sui node process: %v", perr)
			}
		}
	})

	keystoreInstance := testutils.NewTestKeystore(t)
	accountAddress, _ := testutils.GetAccountAndKeyFromSui(keystoreInstance)

	relayerClient, err := client.NewPTBClient(log, client.PTBClientConfig{
		GrpcTarget:            testutils.LocalGrpcURL,
		GrpcToken:             "test",
		TransactionTimeout:    10 * time.Second,
		MaxConcurrentRequests: 5,
		KeystoreService:       keystoreInstance,
		DefaultRequestType:    client.TransactionRequestType("WaitForLocalExecution"),
	})
	require.NoError(t, err)

	require.NoError(t, testutils.FundWithFaucet(log, testutils.SuiLocalnet, accountAddress))

	chainID, err := testutils.GetChainIdentifier(testutils.LocalGrpcURL)
	require.NoError(t, err)
	testutils.PatchEnvironmentTOML("contracts/test", "local", chainID)

	contractPath := testutils.BuildSetup(t, "contracts/test")
	gasBudget := int(2_000_000_000)
	packageID, _, err := testutils.PublishContract(t, "counter", contractPath, accountAddress, &gasBudget)
	require.NoError(t, err)
	require.NotEmpty(t, packageID)

	// The OffRamp latest_config_details read has empty explicit params (its object args are resolved
	// via pointer tags), and so does get_ocr_config — which keeps this reproduction focused on the
	// devInspect read itself rather than on param resolution.
	crConfig := config.ChainReaderConfig{
		IsLoopPlugin: false,
		Modules: map[string]*config.ChainReaderModule{
			"counter": {
				Name: "counter",
				Functions: map[string]*config.ChainReaderFunction{
					"get_ocr_config": {
						Name:                "get_ocr_config",
						SignerAddress:       accountAddress,
						Params:              []codec.SuiFunctionParam{},
						ResultTupleToStruct: []string{"OCRConfig"},
						// big_f -> f so it lands on ocrConfigInfo.F; mirrors the OffRamp config in
						// chainlink core's Sui contract_reader.go.
						ResultFieldRenames: map[string]aptosCRConfig.RenamedField{
							"OCRConfig": {
								SubFieldRenames: map[string]aptosCRConfig.RenamedField{
									"config_info": {
										SubFieldRenames: map[string]aptosCRConfig.RenamedField{
											"big_f": {NewName: "f"},
										},
									},
								},
							},
						},
					},
				},
				Events: map[string]*config.ChainReaderEvent{},
			},
		},
		EventsIndexer: config.EventsIndexerConfig{
			PollingInterval: 2 * time.Second,
			SyncTimeout:     60 * time.Second,
		},
		TransactionsIndexer: config.TransactionsIndexerConfig{
			PollingInterval: 2 * time.Second,
			SyncTimeout:     60 * time.Second,
		},
	}

	counterBinding := types.BoundContract{Name: "counter", Address: packageID}
	readID := strings.Join([]string{packageID, counterBinding.Name, "get_ocr_config"}, "-")

	db := sqltest.NewDB(t, datastoreURL)

	// buildReader constructs a bound ContractReader with the provided read cache (nil == production
	// default, where the read-result cache is disabled).
	buildReader := func(rc *reader.Cache) types.ContractReader {
		evIndexer := indexer.NewEventIndexer(db, log, []*client.EventSelector{})
		txnIndexer := indexer.NewTransactionsIndexer(db, log, map[string]*config.ChainReaderEvent{})
		chainPoller := indexer.NewChainPoller(relayerClient, log, config.ChainPollerConfig{
			PollingInterval:         1 * time.Second,
			SyncTimeout:             60 * time.Second,
			BackfillCheckpointCount: testutils.Uint64Pointer(uint64(1)),
		}, evIndexer.GetEventSelectors)
		indexerInstance := indexer.NewIndexerFromComponents(log, chainPoller, evIndexer, txnIndexer)

		cr, cErr := reader.NewChainReader(ctx, log, relayerClient, crConfig, db, indexerInstance, rc)
		require.NoError(t, cErr)
		require.NoError(t, cr.Bind(ctx, []types.BoundContract{counterBinding}))
		require.NoError(t, indexerInstance.Start(ctx))
		t.Cleanup(func() { _ = indexerInstance.Close() })
		return cr
	}

	readConfig := func(cr types.ContractReader, c context.Context) (ocrConfigWrapped, error) {
		var out ocrConfigWrapped
		e := cr.GetLatestValue(c, readID, primitives.Finalized, map[string]any{}, &out)
		return out, e
	}

	cancelledCtx := func() context.Context {
		c, cancel := context.WithCancel(ctx)
		cancel()
		return c
	}

	// (1) Baseline: a healthy read returns the correct non-zero config digest. Uses the production
	// default (no read cache) to show the happy path works.
	noCacheReader := buildReader(nil)
	t.Run("Baseline_HealthyReadReturnsNonZeroDigest", func(t *testing.T) {
		out, rErr := readConfig(noCacheReader, ctx)
		require.NoError(t, rErr)
		require.Falsef(t, isZeroDigest(out.OCRConfig.ConfigInfo.ConfigDigest),
			"healthy read must return a non-zero config digest, got %x", out.OCRConfig.ConfigInfo.ConfigDigest)
		require.Equal(t, []byte{2, 3, 4, 5}, out.OCRConfig.ConfigInfo.ConfigDigest[:4])
	})

	// (2) Reproduce the CI failure: with the read cache disabled, a transient cancellation makes the
	// read fail. Upstream (CCIP config processor) this failure is turned into a zero-value OCR config
	// in the snapshot — the exact offRampConfigDigest=0x0 that makes the commit plugin reject reports.
	t.Run("Reproduce_CancelledReadFailsWithoutCache", func(t *testing.T) {
		_, rErr := readConfig(noCacheReader, cancelledCtx())
		require.Error(t, rErr,
			"without a read cache a cancelled read fails, and CCIP records a zero-value OCR config for it")
	})

	// (3) Fix direction: with the read-result cache enabled and warmed by a prior healthy read, the
	// same transient cancellation is absorbed — the last-known-good config digest is still returned,
	// so the commit plugin would see a matching digest instead of 0x0.
	t.Run("Fix_CachedConfigSurvivesTransientCancellation", func(t *testing.T) {
		cachedReader := buildReader(reader.NewCache(log, reader.CacheConfig{
			ObjectCacheEnabled: true,
			ReadCacheEnabled:   true,
			ReadTTL:            30 * time.Second,
		}))

		// Warm the cache with a healthy read.
		warm, rErr := readConfig(cachedReader, ctx)
		require.NoError(t, rErr)
		require.False(t, isZeroDigest(warm.OCRConfig.ConfigInfo.ConfigDigest))

		// A cancelled read must now be served from cache (no zero digest, no error).
		out, rErr := readConfig(cachedReader, cancelledCtx())
		require.NoErrorf(t, rErr, "cached config read must survive a transient cancellation")
		require.Falsef(t, isZeroDigest(out.OCRConfig.ConfigInfo.ConfigDigest),
			"cached config read must not return a zero digest under cancellation, got %x",
			out.OCRConfig.ConfigInfo.ConfigDigest)
		require.Equal(t, warm.OCRConfig.ConfigInfo.ConfigDigest, out.OCRConfig.ConfigInfo.ConfigDigest)
	})

	// (4) Serve-stale after the fresh TTL expires: once the fresh window lapses, a transient
	// cancellation must still return the last-known-good config (not a zero/failed read). This is the
	// window that a plain read cache (fresh-only) would not cover.
	t.Run("Fix_ServeStaleAfterTTLExpiry", func(t *testing.T) {
		cachedReader := buildReader(reader.NewCache(log, reader.CacheConfig{
			ObjectCacheEnabled: true,
			ReadCacheEnabled:   true,
			ReadTTL:            200 * time.Millisecond, // short fresh window
			StaleReadTTL:       5 * time.Minute,
		}))

		warm, rErr := readConfig(cachedReader, ctx)
		require.NoError(t, rErr)
		require.False(t, isZeroDigest(warm.OCRConfig.ConfigInfo.ConfigDigest))

		time.Sleep(400 * time.Millisecond) // let the fresh entry expire so a re-fetch is forced

		out, rErr := readConfig(cachedReader, cancelledCtx())
		require.NoErrorf(t, rErr, "after TTL expiry a cancelled read must be served from the stale cache")
		require.Falsef(t, isZeroDigest(out.OCRConfig.ConfigInfo.ConfigDigest),
			"serve-stale must not return a zero digest, got %x", out.OCRConfig.ConfigInfo.ConfigDigest)
		require.Equal(t, warm.OCRConfig.ConfigInfo.ConfigDigest, out.OCRConfig.ConfigInfo.ConfigDigest)
	})

	// (5) Simulate the CI 60/40 mix: a burst of reads where a fraction are cancelled. With a warmed
	// read cache, EVERY read must yield the correct non-zero digest (the commit plugin never sees 0x0).
	t.Run("Fix_MixedHealthyAndCancelledReadsAllResolve", func(t *testing.T) {
		cachedReader := buildReader(reader.NewCache(log, reader.CacheConfig{
			ObjectCacheEnabled: true,
			ReadCacheEnabled:   true,
			ReadTTL:            30 * time.Second,
		}))
		warm, rErr := readConfig(cachedReader, ctx)
		require.NoError(t, rErr)
		want := warm.OCRConfig.ConfigInfo.ConfigDigest

		const iterations = 20
		for i := 0; i < iterations; i++ {
			c := ctx
			if i%2 == 0 {
				c = cancelledCtx() // simulate the frequent transient cancellations seen in CI
			}
			out, e := readConfig(cachedReader, c)
			require.NoErrorf(t, e, "iteration %d: read must not fail", i)
			require.Falsef(t, isZeroDigest(out.OCRConfig.ConfigInfo.ConfigDigest),
				"iteration %d: digest must never be zero", i)
			require.Equalf(t, want, out.OCRConfig.ConfigInfo.ConfigDigest,
				"iteration %d: digest must be stable", i)
		}
	})
}
