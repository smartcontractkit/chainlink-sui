package config

import (
	"time"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"

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
	EventFieldRenames map[string]RenamedField

	// Renames provided filters to match the event field names (optional). When not provided, the filters are used as-is.
	EventFilterRenames map[string]string
}

type RenamedField struct {
	// The new field name (optional). This does not need to be provided if this field does not need
	// to be renamed.
	NewName string

	// Rename sub-fields. This assumes that the event field value is a struct or a map with string keys.
	SubFieldRenames map[string]RenamedField
}

type EventsIndexerConfig struct {
	PollingInterval time.Duration
	SyncTimeout     time.Duration
}

type TransactionsIndexerConfig struct {
	PollingInterval time.Duration
	SyncTimeout     time.Duration
}
