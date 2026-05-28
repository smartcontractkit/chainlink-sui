package indexer

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

// Indexer orchestrates the ChainPoller and the two consumer indexers (EventsIndexer and TransactionsIndexer).
// The ChainPoller fetches checkpoint data from the chain and fans it out to the consumers via channels.
type Indexer struct {
	log     logger.Logger
	starter services.StateMachine

	chainPoller ChainPollerApi

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
}

// NewIndexer creates a new Indexer instance that orchestrates the ChainPoller and consumer indexers.
// The ChainPoller is responsible for fetching checkpoint data and fanning it out to the consumer indexers.
// The channels are created inside the ChainPoller and exposed via EventsChannel() and TransactionsChannel().
func NewIndexer(
	l logger.Logger,
	chainPoller ChainPollerApi,
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
		// Chain poller - runs first and creates the channels
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

		i.wg.Add(1)
		go func() {
			defer i.wg.Done()

			if err := i.transactionIndexer.Start(txnIndexerCtx, i.chainPoller.TransactionsChannel()); err != nil {
				i.log.Errorw("Transaction indexer failed", "error", err)
				i.transactionIndexerErr.Store(err)
			}
		}()

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
