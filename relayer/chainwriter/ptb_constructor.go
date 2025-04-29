package chainwriter

import (
	"context"
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/sui/suiptb"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// PTBConstructor handles building programmable transactions based on configuration
type PTBConstructor struct {
	config ChainWriterConfig
	client client.SuiPTBClient
	log    logger.Logger
}

// NewPTBConstructor creates a new PTB constructor with the given configuration
func NewPTBConstructor(config ChainWriterConfig, ptbClient client.SuiPTBClient, log logger.Logger) *PTBConstructor {
	return &PTBConstructor{
		config: config,
		client: ptbClient,
		log:    log,
	}
}

/*
BuildPTBCommands builds a set of PTB commands based on a signal specified in the ChainWriter configuration.
The goal is to enable ChainWriter to receive a single signal and have the full expressiveness be described
in the config.

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
  - args: a map of named arguments that can be used to fill in the param values of various arguments within PTB calls
*/
func (p *PTBConstructor) BuildPTBCommands(ctx context.Context, moduleName string, function string, args map[string]any) (*suiptb.ProgrammableTransactionBuilder, error) {
	p.log.Debugw("Building PTB commands", "module", moduleName, "function", function, "args", args)

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

	processedArgs := make(map[string]suiptb.Argument)
	// Iterate through the args map values and process each as a PTBArgument
	for key, value := range args {
		processedArg, err := p.client.ToPTBArg(ctx, builder, value)
		if err != nil {
			return nil, err
		}
		processedArgs[key] = processedArg
	}

	// Process each command in order
	for _, cmd := range txnConfig.PTBCommands {
		switch cmd.Type {
		case codec.SuiPTBCommandMoveCall:
			_, err := p.ProcessMoveCall(ctx, builder, cmd, processedArgs)
			if err != nil {
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
