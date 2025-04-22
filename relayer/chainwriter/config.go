package chainwriter

import "github.com/smartcontractkit/chainlink-sui/relayer/codec"

type ChainWriterConfig struct {
	Modules map[string]*ChainWriterModule
}

type ChainWriterModule struct {
	// The module name (optional). When not provided, the key in the map under which this module
	// is stored is used.
	Name      string
	ModuleID  string
	Functions map[string]*ChainWriterFunction
}

type ChainWriterFunction struct {
	// The function name (optional). When not provided, the key in the map under which this function
	// is stored is used.
	Name string
	// The account address (optional). When not provided, the address is calculated
	// from the public key.
	FromAddress string
	Params      []codec.SuiFunctionParam
}
