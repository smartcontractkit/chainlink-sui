package indexer

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-sui/relayer/monitor"
)

type Indexer struct {
	log     logger.Logger
	starter services.StateMachine

	eventsIndexer       EventsIndexerApi
	eventsIndexerCancel *context.CancelFunc
	eventsIndexerErr    atomic.Value // stores error from events indexer goroutine

	transactionIndexer       TransactionsIndexerApi
	transactionIndexerCancel *context.CancelFunc
	transactionIndexerErr    atomic.Value // stores error from transaction indexer goroutine

	wg sync.WaitGroup // wait for both indexer goroutines to exit

	// Health metrics for monitoring (optional)
	healthMetrics *monitor.HealthMetrics
}

type IndexerApi interface {
	Name() string
	Start(ctx context.Context) error
	Ready() error
	HealthReport() map[string]error
	Close() error
	GetEventIndexer() EventsIndexerApi
	GetTransactionIndexer() TransactionsIndexerApi
}

func NewIndexer(
	l logger.Logger,
	eventsIndexer EventsIndexerApi,
	transactionIndexer TransactionsIndexerApi,
) *Indexer {
	return &Indexer{
		log:                      logger.Named(l, "SuiIndexers"),
		eventsIndexer:            eventsIndexer,
		eventsIndexerCancel:      nil,
		transactionIndexer:       transactionIndexer,
		transactionIndexerCancel: nil,
	}
}

func (i *Indexer) Name() string {
	return i.log.Name()
}

func (i *Indexer) Start(_ context.Context) error {
	return i.starter.StartOnce(i.Name(), func() error {
		// Set up health metrics callbacks if health metrics are configured
		if i.healthMetrics != nil {
			i.eventsIndexer.SetOnSyncSuccess(func(ctx context.Context) {
				i.RecordEventsIndexerSuccess(ctx)
			})
			i.transactionIndexer.SetOnSyncSuccess(func(ctx context.Context) {
				i.RecordTransactionsIndexerSuccess(ctx)
			})
		}

		// Events indexer
		eventsIndexerCtx, eventsIndexerCancel := context.WithCancel(context.Background())
		i.eventsIndexerCancel = &eventsIndexerCancel

		i.wg.Add(1)
		go func() {
			defer i.wg.Done()
			defer eventsIndexerCancel()

			if err := i.eventsIndexer.Start(eventsIndexerCtx); err != nil {
				i.log.Errorw("Events indexer failed", "error", err)
				i.eventsIndexerErr.Store(err)
				return
			}
			i.log.Info("Events indexer exited cleanly")
		}()

		// Transaction indexer
		// context.Background() so the TxIndexer's wait loop isn't killed by the parent context
		txnIndexerCtx, txnIndexerCancel := context.WithCancel(context.Background())
		i.transactionIndexerCancel = &txnIndexerCancel

		i.wg.Add(1)
		go func() {
			defer i.wg.Done()
			defer txnIndexerCancel()

			if err := i.transactionIndexer.Start(txnIndexerCtx); err != nil {
				i.log.Errorw("Transaction indexer failed", "error", err)
				i.transactionIndexerErr.Store(err)
				return
			}
			i.log.Info("Transaction indexer exited cleanly")
		}()

		return nil
	})
}

func (i *Indexer) Ready() error {
	if err := i.starter.Ready(); err != nil {
		return err
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

	if err := i.eventsIndexerErr.Load(); err != nil {
		report["EventsIndexer"] = err.(error)
	}
	if err := i.transactionIndexerErr.Load(); err != nil {
		report["TransactionIndexer"] = err.(error)
	}

	return report
}

func (i *Indexer) Close() error {
	return i.starter.StopOnce(i.Name(), func() error {
		// Signal both indexers to stop
		if i.eventsIndexerCancel != nil {
			(*i.eventsIndexerCancel)()
		}
		if i.transactionIndexerCancel != nil {
			(*i.transactionIndexerCancel)()
		}

		i.log.Info("Waiting for indexers to stop...")

		// Wait for both goroutines to exit
		i.wg.Wait()

		i.log.Info("All indexers stopped")

		return nil
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

// SetHealthMetrics sets the health metrics instance for the indexer.
// This should be called after creating the indexer to enable health metrics reporting.
func (i *Indexer) SetHealthMetrics(hm *monitor.HealthMetrics) {
	i.healthMetrics = hm
}

// RecordEventsIndexerSuccess records a successful events indexer sync operation.
func (i *Indexer) RecordEventsIndexerSuccess(ctx context.Context) {
	if i.healthMetrics != nil {
		i.healthMetrics.RecordLastSuccess(ctx, monitor.ComponentEventsIndexer)
	}
}

// RecordTransactionsIndexerSuccess records a successful transactions indexer sync operation.
func (i *Indexer) RecordTransactionsIndexerSuccess(ctx context.Context) {
	if i.healthMetrics != nil {
		i.healthMetrics.RecordLastSuccess(ctx, monitor.ComponentTransactionsIndexer)
	}
}
