package codec

type PTBCommandDependency struct {
	CommandIndex uint16
	ResultIndex  *uint16
}

// SuiFunctionParam defines a parameter for a Sui function call
type SuiFunctionParam struct {
	// Name of the parameter
	Name string
	// Type of the parameter (e.g., "u64", "String", "vector<u8>", "ptb_dependency")
	Type string
	// IsMutable specifies if the object is mutable or not (optional - defaults to true)
	IsMutable *bool
	// Whether the parameter is required

	// IsGeneric specifies if the parameter is a generic argument
	IsGeneric bool

	Required bool
	// Default value to use if not provided
	DefaultValue any
	// Result from a previous PTB Command (optional). It is used for expressive construction of PTB commands
	PTBDependency *PTBCommandDependency
}

type SuiPTBCommandType string

const (
	SuiPTBCommandMoveCall SuiPTBCommandType = "move_call"
	SuiPTBCommandPublish  SuiPTBCommandType = "publish"
	SuiPTBCommandTransfer SuiPTBCommandType = "transfer"
)

// A generic argument can come from 3 places:
//  1. A constant TypeTag string   (e.g. "0x2::sui::SUI")
//  2. A user-supplied parameter   (so you can decide the type at run time)
//  3. A previous PTB command      (mostly when the generic is an object type)
type GenericArg struct {
	// Constant value – the most common case
	TypeTag *string `json:"type_tag,omitempty"`

	// Map to the *name* of a Param in the same Function/Command
	ParamName *string `json:"param_name,omitempty"`

	// Pull the TypeTag from the Nth result of a previous command
	PTBDependency *PTBCommandDependency `json:"ptb_dependency,omitempty"`
}
