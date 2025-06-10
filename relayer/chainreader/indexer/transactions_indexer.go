package indexer

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
)

type TransactionsIndexer struct {
	db              *database.DBStore
	pollingInterval time.Duration
	syncTimeout     time.Duration
}

type TransactionsIndexerApi interface {
	Start(ctx context.Context) error
}

func NewTransactionsIndexer(
	db *database.DBStore,
	pollingInterval time.Duration,
	syncTimeout time.Duration,
) TransactionsIndexerApi {
	return &TransactionsIndexer{
		db:              db,
		pollingInterval: pollingInterval,
		syncTimeout:     syncTimeout,
	}
}

func (*TransactionsIndexer) Start(ctx context.Context) error {
	// TODO: IMPLEMENT
	return nil
}
