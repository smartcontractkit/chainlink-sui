package indexer

import suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"

// CheckpointMeta contains metadata for a checkpoint
type CheckpointMeta struct {
	SequenceNumber uint64
	Digest         string
	TimestampMs    uint64 // from checkpoint summary
}

// CheckpointEventItem represents a single event within a checkpoint's transaction
type CheckpointEventItem struct {
	Event      *suirpcv2.Event
	TxDigest   string
	TxIndex    uint32 // position of the transaction within the checkpoint, as returned by the node
	EventIndex uint32 // position within tx events list (for offset stability)
}

// CheckpointEventsBatch is sent over the events channel from ChainPoller to EventsIndexer
type CheckpointEventsBatch struct {
	Checkpoint CheckpointMeta
	Events     []CheckpointEventItem
}

// CheckpointTransactionsBatch is sent over the transactions channel from ChainPoller to TransactionsIndexer
type CheckpointTransactionsBatch struct {
	Checkpoint   CheckpointMeta
	Transactions []*suirpcv2.ExecutedTransaction
}
