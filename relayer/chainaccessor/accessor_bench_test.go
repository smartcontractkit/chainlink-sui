//go:build integration

// Package-level benchmark comparing the native SuiAccessor read path against the
// equivalent ChainReader read, against a running Sui node.
//
// Required environment:
//   - GRPC_TARGET / GRPC_TOKEN: Sui node gRPC endpoint
//   - TEST_DB_URL: Postgres URL for the event store / indexer
//   - ONRAMP_PACKAGE_ID: deployed CCIP OnRamp package address
//   - DEST_CHAIN_SELECTOR: destination chain selector to query (uint64)
//
// Run with, e.g.:
//
//	GRPC_TARGET=... GRPC_TOKEN=... TEST_DB_URL=... ONRAMP_PACKAGE_ID=0x... \
//	DEST_CHAIN_SELECTOR=1234 go test -tags integration -bench=. -benchmem \
//	./relayer/chainaccessor/...
package chainaccessor

import (
	"context"
	"crypto/ed25519"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	crConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	chainreader "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"
)

// benchHarness holds everything both read paths need.
type benchHarness struct {
	ctx             context.Context
	accessor        *SuiAccessor
	reader          types.ContractReader
	onrampPackageID string
	destChainSel    ccipocr3.ChainSelector
	readIdentifier  string
}

func setupBenchHarness(tb testing.TB) *benchHarness {
	log := logger.Test(tb)
	ctx := context.Background()

	grpcTarget := os.Getenv("GRPC_TARGET")
	grpcToken := os.Getenv("GRPC_TOKEN")
	datastoreURL := os.Getenv("TEST_DB_URL")
	onrampPackageID := os.Getenv("ONRAMP_PACKAGE_ID")
	destRaw := os.Getenv("DEST_CHAIN_SELECTOR")
	if grpcTarget == "" || datastoreURL == "" || onrampPackageID == "" || destRaw == "" {
		tb.Skip("set GRPC_TARGET, TEST_DB_URL, ONRAMP_PACKAGE_ID and DEST_CHAIN_SELECTOR to run accessor benchmarks")
	}
	destSel, err := strconv.ParseUint(destRaw, 10, 64)
	require.NoError(tb, err)

	// devInspect reads use an internal random signer, so a minimal (empty)
	// keystore and placeholder signer address are sufficient for read benchmarks.
	keystoreInstance := &testutils.TestKeystore{Keys: map[string]ed25519.PrivateKey{}}
	accountAddress := "0x0000000000000000000000000000000000000000000000000000000000000000"

	relayerClient, err := client.NewPTBClient(log, client.PTBClientConfig{
		GrpcTarget:            grpcTarget,
		GrpcToken:             grpcToken,
		TransactionTimeout:    10 * time.Second,
		MaxConcurrentRequests: 10,
		KeystoreService:       keystoreInstance,
	})
	require.NoError(tb, err)

	db := sqltest.NewDB(tb, datastoreURL)
	dbStore := database.NewDBStore(db, log)
	require.NoError(tb, dbStore.EnsureSchema(ctx))

	evIndexer := indexer.NewEventIndexer(db, log, []*client.EventSelector{})
	txnIndexer := indexer.NewTransactionsIndexer(db, log, map[string]*crConfig.ChainReaderEvent{})
	chainPoller := indexer.NewChainPoller(relayerClient, log, crConfig.ChainPollerConfig{
		PollingInterval: 1 * time.Second,
		SyncTimeout:     60 * time.Second,
	}, evIndexer.GetEventSelectors)
	indexerInstance := indexer.NewIndexerFromComponents(log, chainPoller, evIndexer, txnIndexer)

	// --- accessor under test ---
	accessor, err := NewSuiAccessor(log, ccipocr3.ChainSelector(destSel), relayerClient, dbStore, indexerInstance)
	require.NoError(tb, err)
	pkgBytes, err := suiAddressToBytes(onrampPackageID)
	require.NoError(tb, err)
	require.NoError(tb, accessor.Sync(ctx, ContractNameOnRamp, pkgBytes))

	// --- equivalent ChainReader read path ---
	onRampStatePointer := &codec.PointerTag{
		Module:        "onramp",
		PointerName:   "OnRampStatePointer",
		DerivationKey: "OnRampState",
		FieldName:     "on_ramp_object_id",
	}
	readerConfig := crConfig.ChainReaderConfig{
		Modules: map[string]*crConfig.ChainReaderModule{
			ContractNameOnRamp: {
				Name: "onramp",
				Functions: map[string]*crConfig.ChainReaderFunction{
					"GetExpectedNextSequenceNumber": {
						Name:          "get_expected_next_sequence_number",
						SignerAddress: accountAddress,
						Params: []codec.SuiFunctionParam{
							{Name: "state", Type: "object_id", PointerTag: onRampStatePointer, Required: true},
							{Name: "dest_chain_selector", Type: "u64", Required: true},
						},
					},
				},
			},
		},
	}
	reader, err := chainreader.NewChainReader(ctx, log, relayerClient, readerConfig, db, indexerInstance)
	require.NoError(tb, err)
	require.NoError(tb, reader.Bind(ctx, []types.BoundContract{{
		Name:    ContractNameOnRamp,
		Address: onrampPackageID,
	}}))

	return &benchHarness{
		ctx:             ctx,
		accessor:        accessor,
		reader:          reader,
		onrampPackageID: onrampPackageID,
		destChainSel:    ccipocr3.ChainSelector(destSel),
		readIdentifier: strings.Join(
			[]string{onrampPackageID, ContractNameOnRamp, "GetExpectedNextSequenceNumber"}, "-"),
	}
}

// BenchmarkGetExpectedNextSequenceNumber_Accessor measures the native accessor
// read path (devInspect, no config-driven indirection).
func BenchmarkGetExpectedNextSequenceNumber_Accessor(b *testing.B) {
	h := setupBenchHarness(b)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := h.accessor.GetExpectedNextSequenceNumber(h.ctx, h.destChainSel); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetExpectedNextSequenceNumber_ChainReader measures the equivalent
// read through the ChainReader (readIdentifier parse → bind validate → pointer
// resolution → devInspect).
func BenchmarkGetExpectedNextSequenceNumber_ChainReader(b *testing.B) {
	h := setupBenchHarness(b)
	params := map[string]any{"dest_chain_selector": uint64(h.destChainSel)}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result uint64
		if err := h.reader.GetLatestValue(h.ctx, h.readIdentifier, primitives.Unconfirmed, params, &result); err != nil {
			b.Fatal(err)
		}
	}
}
