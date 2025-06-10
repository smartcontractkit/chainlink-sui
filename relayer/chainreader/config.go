package chainreader

import (
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

type ChainReaderConfig struct {
	IsLoopPlugin        bool
	EventsIndexer       EventsIndexerConfig
	TransactionsIndexer TransactionsIndexerConfig
	Modules             map[string]*ChainReaderModule
}

type ChainReaderModule struct {
	// The module name (optional). When not provided, the key in the map under which this module
	// is stored is used.
	Name      string
	Functions map[string]*ChainReaderFunction
	Events    map[string]*ChainReaderEvent
}

type ChainReaderFunction struct {
	// The function name (optional). When not provided, the key in the map under which this function
	// is stored is used.
	Name          string
	SignerAddress string
	Params        []codec.SuiFunctionParam
}

type ChainReaderEvent struct {
	// The event name (optional). When not provided, the key in the map under which this event
	// is stored is used.
	Name      string
	EventType string
	// EventSelector specifies how the event is tagged within a package, and it includes
	// the 3 fields of the tag `packageId::moduleId::eventId`
	client.EventSelector
}

type SequenceWithMetadata struct {
	Sequence  pkgtypes.Sequence
	TxVersion uint64
	TxHash    string
}

type EventsIndexerConfig struct {
	PollingInterval time.Duration
	SyncTimeout     time.Duration
}

type TransactionsIndexerConfig struct {
	PollingInterval time.Duration
	SyncTimeout     time.Duration
}
