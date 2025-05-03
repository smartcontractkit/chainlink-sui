package chainreader

import (
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

type ChainReaderConfig struct {
	Modules map[string]*ChainReaderModule
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
}
