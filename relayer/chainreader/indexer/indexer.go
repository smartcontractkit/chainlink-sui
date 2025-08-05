package indexer

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
)

type Indexer struct {
	log     logger.Logger
	starter services.StateMachine

	eventsIndexer       EventsIndexerApi
	eventsIndexerCancel *context.CancelFunc

	transactionIndexer       TransactionsIndexerApi
	transactionIndexerCancel *context.CancelFunc
}

type IndexerApi interface {
	Name() string
	Start(ctx context.Context) error
	Ready() error
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
		log:                      logger.Named(l, "Indexers"),
		eventsIndexer:            eventsIndexer,
		eventsIndexerCancel:      nil,
		transactionIndexer:       transactionIndexer,
		transactionIndexerCancel: nil,
	}
}

func (i *Indexer) Name() string {
	return i.log.Name()
}

func (i *Indexer) Start(ctx context.Context) error {
	return i.starter.StartOnce(i.Name(), func() error {
		txnIndexerCtx, txnIndexerCancel := context.WithCancel(ctx)

		go func() {
			if err := i.transactionIndexer.Start(txnIndexerCtx); err != nil {
				i.log.Error("Transaction indexer failed to start", "error", err)
				txnIndexerCancel()
				return
			}

			i.log.Info("Events indexer started")
			// set the cancel function
			i.transactionIndexerCancel = &txnIndexerCancel
		}()

		eventsIndexerCtx, eventsIndexerCancel := context.WithCancel(ctx)

		go func() {
			if err := i.eventsIndexer.Start(eventsIndexerCtx); err != nil {
				i.log.Error("Events indexer failed to start", "error", err)
				eventsIndexerCancel()
				return
			}

			i.log.Info("Events indexer started")
			// set the cancel function
			i.eventsIndexerCancel = &eventsIndexerCancel
		}()

		// If either of the indexers failed to start, we return an error
		if i.transactionIndexerCancel == nil || i.eventsIndexerCancel == nil {
			return fmt.Errorf("Indexers failed to start, cancel functions are nil")
		}

		return nil
	})
}

func (i *Indexer) Ready() error {
	return i.starter.Ready()
}

func (i *Indexer) Close() error {
	return i.starter.StopOnce(i.Name(), func() error {
		// stop events indexer
		if i.eventsIndexerCancel != nil {
			(*i.eventsIndexerCancel)()
		}
		i.log.Info("Events indexer stopped")

		// stop transactions indexer
		if i.transactionIndexerCancel != nil {
			(*i.transactionIndexerCancel)()
		}
		i.log.Info("Transactions indexer stopped")

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
