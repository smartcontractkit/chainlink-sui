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

// PrerequisiteObject represents a structure defining requirements or dependencies needed before constructing the PTB.
// These requirements refer to object details that need to be fetched with `SuiX_GetOwnedObjects` and then populated
// into the arguments map provided for PTB construction.
//
// The usage flow is that a request is made to get all the owned objects by "OwnerId" and then picking the one
// that matches the Tag
type PrerequisiteObject struct {
	OwnerId *string
	Name    string // the key under which the value is inserted in the args, must match one of the arg names used in the PTB commands
	Tag     string
	SetKeys bool // optionally set the keys of the object in the arg map instead of name
}

type ChainWriterFunction struct {
	// The function name (optional). When not provided, the key in the map under which this function
	// is stored is used.
	Name string
	// The public key of the account that will sign and submit the transaction.
	PublicKey []byte
	// The values that need to be loaded into the args by making SuiX_GetOwnedObjects calls
	PrerequisiteObjects []PrerequisiteObject
	Params              []codec.SuiFunctionParam
	// The set of PTB commands to run as part of this function call.
	// This field is used in replacement of `Params` above.
	PTBCommands []ChainWriterPTBCommand
}

// ConfigOverrides contains fields with dynamic values to override the default configs
type ConfigOverrides struct {
	// ToAddress specifies an override for the owner address in PrerequisiteObject such that if it is
	// empty in the config, the value can be passed in from the chainwriter.SendTransaction method
	ToAddress string
}

type ChainWriterSignal struct {
}
