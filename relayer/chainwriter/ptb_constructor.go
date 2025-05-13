package chainwriter

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// PTBArgType represents the type of argument in a Programmable Transaction Block
type PTBArgType string

const (
	// ObjectArgType represents a reference to a Sui object
	ObjectArgType PTBArgType = "object"
	// ScalarArgType represents a basic value type (number, string, etc.)
	ScalarArgType PTBArgType = "scalar"
	// VectorArgType represents a vector of arguments
	VectorArgType PTBArgType = "vector"
)

// PTBArg represents a single argument in a Programmable Transaction Block
type PTBArg struct {
	// Type indicates the kind of argument (object, scalar, or vector)
	Type PTBArgType `json:"type"`
	// Content contains the actual argument data
	Content PTBArgContent `json:"content"`
}

// PTBArgContent contains the actual data for a PTB argument
type PTBArgContent struct {
	// ID is the object ID for object-type arguments
	ID string `json:"id,omitempty"`
	// Value contains the actual value for scalar-type arguments
	Value any `json:"value,omitempty"`
	// MapTo specifies how this argument maps to parameters in the PTB
	MapTo []PTBArgLocation `json:"map_to,omitempty"`
}

// PTBArgLocation specifies how an argument maps to a parameter in a PTB command
type PTBArgLocation struct {
	// CommandIndex is the index of the command in the PTB
	CommandIndex int `json:"command_index"`
	// Param is the name of the parameter in the command
	Param string `json:"param"`
	// CommandName is the full name of the command (e.g., "0x2::coin::transfer")
	CommandName string `json:"command_name,omitempty"`
}

// PTBArgMapping represents a collection of arguments for a Programmable Transaction Block
type PTBArgMapping struct {
	// Args is the list of arguments in the PTB
	Args []PTBArg `json:"args"`
}

// PTBConstructor handles building programmable transactions based on configuration.
// It provides methods to construct PTBs by mapping arguments to their respective commands
// and handling dependencies between commands.
type PTBConstructor struct {
	config ChainWriterConfig   // Configuration for building PTBs
	client client.SuiPTBClient // Client for interacting with Sui PTB functionality
	log    logger.Logger       // Logger for debugging and error reporting
}

// NewPTBConstructor creates a new PTB constructor with the given configuration.
// It initializes a PTBConstructor with the provided config, client, and logger.
//
// Parameters:
//   - config: The ChainWriterConfig containing module and function definitions
//   - ptbClient: The SuiPTBClient for interacting with Sui PTB functionality
//   - log: Logger for debugging and error reporting
//
// Returns:
//   - *PTBConstructor: A new instance of PTBConstructor
func NewPTBConstructor(config ChainWriterConfig, ptbClient client.SuiPTBClient, log logger.Logger) *PTBConstructor {
	return &PTBConstructor{
		config: config,
		client: ptbClient,
		log:    log,
	}
}

/*
BuildPTBCommands builds a set of PTB commands based on a signal specified in the ChainWriter configuration.
The function first builds all PTB arguments (both object and scalar) before constructing the commands.
This ensures that each argument is only built once, even if it's used in multiple commands.

The process follows these steps:
1. Builds all object arguments first, storing them in a map keyed by their IDs
2. Builds all scalar arguments, storing them in a map with generated keys
3. Processes each command in order, mapping the pre-built arguments to their respective parameters
4. Handles PTB dependencies between commands
5. Constructs the final PTB with all commands and their arguments

An example of a ChainWriter config is shown below to illustrate how expressive configuration can define a set
of PTB commands that will be execution from a single signal. In the example below, the ChainWriter will receive
a function indicating that the function "use_ptb" within the "ptb_builder" module should be called.

The BuildPTBCommands method should then be called by the ChainWriter to convert that signal into N commands (in this case 2 commands)
internally and return a single transaction that can be executed on the Sui node.

Each `PTBCommand` (chainwriter.ChainWriterPTBCommand) within the configuration defines a possible PTB action (e.g. MoveCall or Publish)
along with the necessary parameters (arguments) to run it (codec.SuiFunctionParam).

Each parameter can have an optional `PTBDependency` field (codec.PTBCommandDependency) which defines a dependency on the results
of previous commands within the same PTB.

```

	chainWriterConfig := chainwriter.ChainWriterConfig{
		Modules: map[string]*chainwriter.ChainWriterModule{
			"ptb_builder": {
				Name:     "ptb_builder",
				ModuleID: "...",
				Functions: map[string]*chainwriter.ChainWriterFunction{
					"use_ptb": {
						Name:        "use_ptb",
						FromAddress: testState.AccountAddress,
						Params:      []codec.SuiFunctionParam{},
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							{
								Type:       codec.SuiPTBCommandMoveCall,
								PackageId: "...",
								ModuleId:   "...",
								Function:   "get_counter",
								Params: []codec.SuiFunctionParam{
									{
										Name:         "counter_object_id",
										Type:         "object_id",
										Required:     true,
										DefaultValue: nil,
									},
								},
								Order: 1,
							},
							{
								Type:       codec.SuiPTBCommandMoveCall,
								PackageId: "...",
								ModuleId:   "...",
								Function:   "increment_by",
								Params: []codec.SuiFunctionParam{
									{
										Name:         "counter_object_id",
										Type:         "object_id",
										Required:     true,
										DefaultValue: nil,
									},
									{
										Name: "new_counter_value",
										Type: "ptb_dependency",
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 1,
											ResultIndex:  1,
										},
										Required:     true,
										DefaultValue: nil,
									},
								},
								Order: 2,
							},
						},
					},
				},
			},
		},
	}

```

Parameters:
  - ctx: context
  - moduleName: the name of the module key in the config
  - function: the name of the signal (virtual function) which does not actually map to a single contract call
  - argMapping: a structured representation of the arguments for various commands within PTB, containing both object and scalar arguments
*/
func (p *PTBConstructor) BuildPTBCommands(ctx context.Context, moduleName string, function string, argMapping PTBArgMapping) (*suiptb.ProgrammableTransactionBuilder, error) {
	p.log.Debugw("Building PTB commands", "module", moduleName, "function", function, "args", argMapping)

	// Look up the module
	module, ok := p.config.Modules[moduleName]
	if !ok {
		return nil, fmt.Errorf("missing module %s not found in configuration", moduleName)
	}

	// Look up the transaction
	txnConfig, ok := module.Functions[function]
	if !ok {
		return nil, fmt.Errorf("missing function config (%s) not found in module (%s)", function, moduleName)
	}

	// Create a new PTB builder
	builder := suiptb.NewTransactionDataTransactionBuilder()

	// Build all PTB arguments first
	ptbArgs := make(map[string]suiptb.Argument)

	// Process object arguments
	for _, obj := range argMapping.Args {
		if obj.Type == ObjectArgType {
			p.log.Debugw("Processing object argument", "ID", obj.Content.ID)
			ptbArg, err := p.client.ToPTBArg(ctx, builder, obj.Content.ID)
			if err != nil {
				p.log.Errorw("Error processing object argument", "ID", obj.Content.ID, "Error", err)
				return nil, err
			}
			ptbArgs[obj.Content.ID] = ptbArg
		}
	}

	// Process scalar arguments
	for _, scalar := range argMapping.Args {
		if scalar.Type == ScalarArgType {
			scalarKey := fmt.Sprintf("scalar_%d", len(ptbArgs))
			ptbArg, err := p.client.ToPTBArg(ctx, builder, scalar.Content.Value)
			if err != nil {
				p.log.Errorw("Error processing scalar argument", "Value", scalar.Content.Value, "Error", err)
				return nil, err
			}
			ptbArgs[scalarKey] = ptbArg
		}
	}

	// Process each command in order
	for cmdIdx, cmd := range txnConfig.PTBCommands {
		args := make(map[string]suiptb.Argument)

		// Map object arguments to command parameters
		for _, obj := range argMapping.Args {
			if obj.Type == ObjectArgType {
				for _, loc := range obj.Content.MapTo {
					if loc.CommandIndex == cmdIdx {
						args[loc.Param] = ptbArgs[obj.Content.ID]
					}
				}
			}
		}

		// Map scalar arguments to command parameters
		for _, scalar := range argMapping.Args {
			if scalar.Type == ScalarArgType {
				for _, loc := range scalar.Content.MapTo {
					if loc.CommandIndex == cmdIdx {
						args[loc.Param] = ptbArgs[fmt.Sprintf("scalar_%d", len(ptbArgs)-1)]
					}
				}
			}
		}

		// Process PTB dependencies
		for _, param := range cmd.Params {
			if param.PTBDependency != nil {
				args[param.Name] = suiptb.Argument{
					NestedResult: &suiptb.NestedResult{
						Cmd:    param.PTBDependency.CommandIndex,
						Result: param.PTBDependency.ResultIndex,
					},
				}
			}
		}

		switch cmd.Type {
		case codec.SuiPTBCommandMoveCall:
			_, err := p.ProcessMoveCall(ctx, builder, cmd, args)
			if err != nil {
				p.log.Errorw("Error processing move call", "Error", err)
				return nil, err
			}
		case codec.SuiPTBCommandPublish:
			return nil, fmt.Errorf("publishing is not supported yet")
		case codec.SuiPTBCommandTransfer:
			return nil, fmt.Errorf("transfers are not supported yet")
		default:
			return nil, fmt.Errorf("unsupported command type: %v", cmd.Type)
		}
	}

	return builder, nil
}

// ProcessMoveCall handles constructing move call commands and adds it to the PTB `builder` instance.
func (p *PTBConstructor) ProcessMoveCall(
	ctx context.Context,
	builder *suiptb.ProgrammableTransactionBuilder,
	cmd ChainWriterPTBCommand,
	args map[string]suiptb.Argument,
) (*suiptb.Argument, error) {
	p.log.Debugw("Processing move call", "Command", cmd, "Args", args)

	// All three fields below are required for a successful move call
	if cmd.Params == nil {
		return nil, fmt.Errorf("missing required parameter 'PackageId' for move call PTB command")
	}
	if cmd.ModuleId == nil {
		return nil, fmt.Errorf("missing required parameter 'ModuleId' for move call PTB command")
	}
	if cmd.Function == nil {
		return nil, fmt.Errorf("missing required parameter 'Function' for move call PTB command")
	}

	// Convert package ID to Address
	packageId, err := sui.AddressFromHex(*cmd.PackageId)
	if err != nil {
		return nil, fmt.Errorf("failed to build package address from hex (%s): %w", *cmd.PackageId, err)
	}

	// Process arguments
	processedArgTypes := []sui.TypeTag{}
	processedArgs, err := p.ProcessParams(ctx, builder, cmd.Params, args)
	if err != nil {
		return nil, err
	}

	p.log.Debugw("Processed args", "Args", processedArgs)
	// Add the move call to the builder
	ptbArgument := builder.ProgrammableMoveCall(packageId, *cmd.ModuleId, *cmd.Function, processedArgTypes, processedArgs)

	return &ptbArgument, nil
}

// ProcessParams converts parameter specifications into concrete arguments
func (p *PTBConstructor) ProcessParams(
	ctx context.Context,
	builder *suiptb.ProgrammableTransactionBuilder,
	params []codec.SuiFunctionParam,
	args map[string]suiptb.Argument,
) ([]suiptb.Argument, error) {
	processedArgs := make([]suiptb.Argument, 0, len(params))

	for _, param := range params {
		// check if this is a PTB result dependency
		if param.PTBDependency != nil {
			processedArgs = append(processedArgs, suiptb.Argument{
				NestedResult: &suiptb.NestedResult{
					Cmd:    param.PTBDependency.CommandIndex,
					Result: param.PTBDependency.ResultIndex,
				},
			})

			continue
		}

		// otherwise, check if the parameter is in the provided args
		if providedArg, exists := args[param.Name]; exists {
			// append to the array of args
			processedArgs = append(processedArgs, providedArg)
			continue
		}

		// fallback to the default value
		if param.DefaultValue != nil {
			ptbArg, err := p.client.ToPTBArg(ctx, builder, param.DefaultValue)
			if err != nil {
				return nil, err
			}
			// append to the array of args
			processedArgs = append(processedArgs, ptbArg)

			continue
		}

		// Value not found for required param
		if param.Required {
			return nil, fmt.Errorf("required parameter %s has no value", param.Name)
		}

		// append an empty argument since it is not required and no value found
		processedArgs = append(processedArgs, suiptb.Argument{})
	}

	return processedArgs, nil
}

// ParseArgsToPTBArgMapping converts the input arguments to a PTBArgMapping structure.
// It handles both map[string]any and PTBArgMapping input types.
func (p *PTBConstructor) ParseArgsToPTBArgMapping(args any) (PTBArgMapping, error) {
	p.log.Debugw("Parsing args to PTBArgMapping", "args", args)

	var mapping PTBArgMapping

	err := mapstructure.Decode(args, &mapping)
	if err != nil {
		p.log.Errorw("Error decoding args", "error", err)
		return PTBArgMapping{}, err
	}

	return mapping, nil
}
