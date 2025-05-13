package chainwriter

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

var PTBChainWriterModuleName = "cll://component=cw/type=ptb_builder"

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

type ChainWriterPTBCommand struct {
	Type codec.SuiPTBCommandType
	// The package ID to call (optional). This may not be needed in the case
	// that the type of PTB command does not require it (e.g. Publish).
	PackageId *string                  `json:"package_id,omitempty"`
	ModuleId  *string                  `json:"module_id,omitempty"`
	Function  *string                  `json:"function,omitempty"`
	Params    []codec.SuiFunctionParam `json:"params,omitempty"`
	// TODO: is this needed? is order of array items maintained?
	Order int `json:"order"`
}

// GetParamKey returns the key for a parameter in the PTB command in a map of arguments.
// The key is a string that uniquely identifies the parameter within the map of arguments.
// The key is formatted as follows:
// "packageId::moduleId::functionName::parameterName"
// This format allows associating specific argument values with their target
// Move function call and parameter name within a potentially complex PTB.
func (c ChainWriterPTBCommand) GetParamKey(paramName string) string {
	return fmt.Sprintf("%s.%s.%s.%s", *c.PackageId, *c.ModuleId, *c.Function, paramName)
}

type ChainWriterFunction struct {
	// The function name (optional). When not provided, the key in the map under which this function
	// is stored is used.
	Name string
	// The public key of the account that will sign and submit the transaction.
	PublicKey []byte
	Params    []codec.SuiFunctionParam
	// The set of PTB commands to run as part of this function call.
	// This field is used in replacement of `Params` above.
	PTBCommands []ChainWriterPTBCommand
}

type ChainWriterSignal struct {
}
