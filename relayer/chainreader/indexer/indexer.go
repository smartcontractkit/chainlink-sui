package indexer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	"github.com/smartcontractkit/chainlink-common/pkg/types/sui"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// checkpointCursorID namespaces the ChainPoller's row in sui.checkpoint_cursors.
const checkpointCursorID = "chain_poller"

// Indexer orchestrates the ChainPoller and the two consumer indexers (EventsIndexer and TransactionsIndexer).
// The ChainPoller fetches checkpoint data from the chain and fans it out to the consumers via channels.
type Indexer struct {
	log     logger.Logger
	starter services.StateMachine

	// dbStore, when set (production wiring via NewIndexer), has its schema ensured before the
	// ChainPoller starts, since the poller reads the checkpoint cursor at startup.
	dbStore *database.DBStore

	chainPoller ChainPollerAPI

	eventsIndexer       EventsIndexerApi
	eventsIndexerCancel *context.CancelFunc
	eventsIndexerErr    atomic.Value // stores error from events indexer goroutine

	transactionIndexer       TransactionsIndexerApi
	transactionIndexerCancel *context.CancelFunc
	transactionIndexerErr    atomic.Value // stores error from transaction indexer goroutine

	pollerCancel *context.CancelFunc
	pollerErr    atomic.Value // stores error from poller goroutine

	wg sync.WaitGroup // wait for poller + both indexer goroutines to exit
}

// IndexerApi defines the interface for the combined indexer orchestration.
type IndexerApi interface {
	Name() string
	Start(ctx context.Context) error
	Ready() error
	HealthReport() map[string]error
	Close() error
	GetEventIndexer() EventsIndexerApi
	GetTransactionIndexer() TransactionsIndexerApi
	// RescanFromCheckpoint rewinds the ChainPoller so a bounded window of checkpoints starting at
	// fromSeq is re-scanned. The relayer calls this when the core node requests a Replay, so events
	// and transactions the poller already processed (or missed) are picked up on a second pass.
	RescanFromCheckpoint(ctx context.Context, fromSeq uint64) error
}

// Params holds the dependencies needed to construct a fully-wired Indexer via NewIndexer.
type Params struct {
	Logger logger.Logger
	DB     sqlutil.DataSource
	Client client.SuiPTBClient
	// Poller cursor ID (defaults to "chain_poller")
	PollerCursorId string
	PollerConfig   sui.ChainPollerConfig
	// PollerWorkers is the number of concurrent checkpoint catch-up workers; non-positive
	// values use the ChainPoller default.
	PollerWorkers int
	// PollerChunkSize is the number of checkpoints per catch-up chunk; non-positive values use
	// the ChainPoller default.
	PollerChunkSize int
	// PollerReplayCheckpointCount is how many checkpoints a Replay request re-scans, starting from
	// the requested checkpoint; zero re-scans all the way to the latest checkpoint.
	PollerReplayCheckpointCount uint64
	// EventSelectors optionally seeds the EventsIndexer. Selectors are normally registered later
	// when the ChainReader binds contracts, so this is usually left nil/empty.
	EventSelectors []*sui.EventFilterByMoveEventModule
	// TransactionConfigs optionally seeds the TransactionsIndexer. Normally left nil/empty and
	// populated later via the ChainReader.
	TransactionConfigs map[string]*sui.ChainReaderEvent
}

// NewIndexer constructs a fully-wired Indexer: it creates the TransactionsIndexer, EventsIndexer,
// and ChainPoller, and connects the poller's SelectorProvider to the EventsIndexer so that event
// selectors registered at bind time are honored during polling. This is the production entry point;
// tests that need to inject mocks can use NewIndexerFromComponents instead.
func NewIndexer(p Params) *Indexer {
	eventSelectors := p.EventSelectors
	if eventSelectors == nil {
		eventSelectors = []*sui.EventFilterByMoveEventModule{}
	}
	txnConfigs := p.TransactionConfigs
	if txnConfigs == nil {
		txnConfigs = map[string]*sui.ChainReaderEvent{}
	}

	txnIndexer := NewTransactionsIndexer(p.DB, p.Logger, txnConfigs)
	eventsIndexer := NewEventIndexer(p.DB, p.Logger, eventSelectors)

	// The poller pulls the live selector set from the events indexer on each checkpoint, so
	// selectors added later (e.g. during Bind) are picked up without re-wiring.
	dbStore := database.NewDBStore(p.DB, p.Logger)

	// Allow overwriting the name of the default cursor ID (stored in DB) to enable running
	// multiple instances of the indexer using the same database.
	checkpointCursor := checkpointCursorID
	if p.PollerCursorId != "" {
		checkpointCursor = p.PollerCursorId
	}

	chainPoller := NewChainPoller(
		p.Client,
		p.Logger,
		p.PollerConfig,
		eventsIndexer.GetEventSelectors,
		WithWorkerPool(p.PollerWorkers, p.PollerChunkSize),
		WithCursorStore(dbStore, checkpointCursor),
		WithRescanCheckpointCount(p.PollerReplayCheckpointCount),
	)

	idx := NewIndexerFromComponents(p.Logger, chainPoller, eventsIndexer, txnIndexer)
	idx.dbStore = dbStore

	return idx
}

// NewIndexerFromComponents creates an Indexer from already-constructed components. Prefer NewIndexer
// for production wiring; this constructor exists for tests that inject mocks or need direct handles
// to the poller / consumer indexers.
func NewIndexerFromComponents(
	l logger.Logger,
	chainPoller ChainPollerAPI,
	eventsIndexer EventsIndexerApi,
	transactionIndexer TransactionsIndexerApi,
) *Indexer {
	return &Indexer{
		log:                logger.Named(l, "SuiIndexers"),
		chainPoller:        chainPoller,
		eventsIndexer:      eventsIndexer,
		transactionIndexer: transactionIndexer,
	}
}

func (i *Indexer) Name() string {
	return i.log.Name()
}

// Start begins the ChainPoller and then starts the consumer indexers.
// The ChainPoller fetches checkpoints and sends data over channels to the consumers.
// When the poller stops, it closes the channels, which signals the consumers to exit.
func (i *Indexer) Start(_ context.Context) error {
	return i.starter.StartOnce(i.Name(), func() error {
		// Ensure the database schema exists before the poller starts: the poller reads the
		// persisted checkpoint cursor at startup. (EventsIndexer.Start also ensures the schema,
		// idempotently, for tests that start it standalone.)
		if i.dbStore != nil {
			schemaCtx, schemaCancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer schemaCancel()
			if err := i.dbStore.EnsureSchema(schemaCtx); err != nil {
				return fmt.Errorf("failed to ensure schema: %w", err)
			}
		}

		// Chain poller - runs first and creates the channels
		//nolint:gosec // G118: pollerCancel is invoked from Stop()
		pollerCtx, pollerCancel := context.WithCancel(context.Background())
		i.pollerCancel = &pollerCancel

		if err := i.chainPoller.Start(pollerCtx); err != nil {
			return err
		}

		// Events indexer - consumes from the chainPoller's events channel
		eventsIndexerCtx, eventsIndexerCancel := context.WithCancel(context.Background())
		i.eventsIndexerCancel = &eventsIndexerCancel

		if err := i.eventsIndexer.Start(eventsIndexerCtx, i.chainPoller.EventsChannel()); err != nil {
			return err
		}

		// Transaction indexer Start blocks until the OffRamp package is bound; run in the
		// background so Indexer.Start can return while the poller and events indexer run.
		txnIndexerCtx, txnIndexerCancel := context.WithCancel(context.Background())
		i.transactionIndexerCancel = &txnIndexerCancel

		i.wg.Go(func() {
			if err := i.transactionIndexer.Start(txnIndexerCtx, i.chainPoller.TransactionsChannel()); err != nil {
				i.log.Errorw("Transaction indexer failed", "error", err)
				i.transactionIndexerErr.Store(err)
			}
		})

		return nil
	})
}

func (i *Indexer) Ready() error {
	if err := i.starter.Ready(); err != nil {
		return err
	}

	// Check if poller has failed
	if err := i.pollerErr.Load(); err != nil {
		return err.(error)
	}

	// Check if either indexer has failed
	if err := i.eventsIndexerErr.Load(); err != nil {
		return err.(error)
	}
	if err := i.transactionIndexerErr.Load(); err != nil {
		return err.(error)
	}

	return nil
}

func (i *Indexer) HealthReport() map[string]error {
	report := map[string]error{
		i.Name(): i.starter.Healthy(),
	}

	if err := i.pollerErr.Load(); err != nil {
		report["ChainPoller"] = err.(error)
	}
	if err := i.eventsIndexerErr.Load(); err != nil {
		report["EventsIndexer"] = err.(error)
	}
	if err := i.transactionIndexerErr.Load(); err != nil {
		report["TransactionIndexer"] = err.(error)
	}

	return report
}

// Close stops all components. The stop order is:
// 1. Stop the ChainPoller (cancels its context, which closes channels)
// 2. Wait for consumers (they will exit when channels are closed)
// 3. All goroutines complete via WaitGroup
func (i *Indexer) Close() error {
	return i.starter.StopOnce(i.Name(), func() error {
		if i.pollerCancel != nil {
			(*i.pollerCancel)()
		}
		if i.eventsIndexerCancel != nil {
			(*i.eventsIndexerCancel)()
		}
		if i.transactionIndexerCancel != nil {
			(*i.transactionIndexerCancel)()
		}

		i.log.Info("Waiting for ChainPoller and indexers to stop...")

		var closeErr error
		if err := i.chainPoller.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
		if err := i.eventsIndexer.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}
		if err := i.transactionIndexer.Close(); err != nil {
			closeErr = errors.Join(closeErr, err)
		}

		// Background TransactionsIndexer.Start goroutine
		i.wg.Wait()

		i.log.Info("All indexers stopped")

		return closeErr
	})
}

func (i *Indexer) GetEventIndexer() EventsIndexerApi {
	if i.eventsIndexer == nil {
		return nil
	}
	return i.eventsIndexer
}

func (i *Indexer) GetTransactionIndexer() TransactionsIndexerApi {
	if i.transactionIndexer == nil {
		return nil
	}
	return i.transactionIndexer
}

// RescanFromCheckpoint rewinds the ChainPoller so a bounded window of checkpoints starting at
// fromSeq is re-scanned. See the IndexerApi docs: the relayer calls this when the core node
// requests a Replay.
func (i *Indexer) RescanFromCheckpoint(ctx context.Context, fromSeq uint64) error {
	if i.chainPoller == nil {
		return errors.New("chain poller not configured")
	}
	return i.chainPoller.RescanFrom(ctx, fromSeq)
}
