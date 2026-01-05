package config

import (
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	aptosCRConfig "github.com/smartcontractkit/chainlink-aptos/relayer/chainreader/config"

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
	// Defines a way to transform a tuple result into a JSON object
	ResultTupleToStruct []string
	// Defines a mapping for renaming response fields
	ResultFieldRenames map[string]aptosCRConfig.RenamedField
	// Static response
	StaticResponse []any
	// Response from inputs
	ResponseFromInputs []string
}

type ChainReaderEvent struct {
	// The event name (optional). When not provided, the key in the map under which this event
	// is stored is used.
	Name      string
	EventType string
	// EventSelector specifies how the event is tagged within a package, and it includes
	// the 3 fields of the tag `packageId::moduleId::eventId`
	client.EventSelector

	// Renames of event field names (optional). When not provided, the field names are used as-is.
	EventFieldRenames map[string]aptosCRConfig.RenamedField

	// Renames provided filters to match the event field names (optional). When not provided, the filters are used as-is.
	EventFilterRenames map[string]string

	// The expected event type (optional). When not provided, the event type is used as-is.
	ExpectedEventType any

	// A fallback for events selectors with no offset recorded in the DB and a starting point
	// earlier than the pruning cutoff of the RPC
	EventSelectorDefaultOffset *client.EventId
}

type EventsIndexerConfig struct {
	PollingInterval time.Duration
	SyncTimeout     time.Duration
}

type TransactionsIndexerConfig struct {
	PollingInterval time.Duration
	SyncTimeout     time.Duration
}
