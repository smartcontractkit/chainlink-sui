package ptb

import (
	"context"
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/transaction"

	cwConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainwriter/ptb/offramp"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

// PTBConstructor handles building programmable transactions based on configuration.
// It provides methods to construct PTBs by mapping arguments to their respective commands
// and handling dependencies between commands.
type PTBConstructor struct {
	config cwConfig.ChainWriterConfig // Configuration for building PTBs
	client client.SuiPTBClient        // Client for interacting with Sui PTB functionality
	log    logger.Logger              // Logger for debugging and error reporting
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
func NewPTBConstructor(config cwConfig.ChainWriterConfig, ptbClient client.SuiPTBClient, log logger.Logger) *PTBConstructor {
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
func (p *PTBConstructor) BuildPTBCommands(ctx context.Context, moduleName string, function string, arguments cwConfig.Arguments, toAddress string, txnConfig *cwConfig.ChainWriterFunction) (*transaction.Transaction, error) {
	p.log.Debugw("Building PTB commands", "module", moduleName, "function", function)

	// Create a new transaction builder
	sdkClient := p.client.GetClient()
	ptb := transaction.NewTransaction()
	ptb.SetSuiClient(sdkClient.(*sui.Client))

	signerAddress, err := client.GetAddressFromPublicKey(txnConfig.PublicKey)
	if err != nil {
		return nil, err
	}

	// If the function is Execute, then we need to build the PTB using bespoke code rather than using the configs to programmatically build the PTB commands.
	if function == cwConfig.CCIPExecute {
		addressMappings, err := offramp.GetOfframpAddressMappings(ctx, p.log, p.client, toAddress, txnConfig.PublicKey)
		if err != nil {
			p.log.Errorw("Error setting up address mappings", "error", err)
			return nil, err
		}

		// Construct the entire PTB transaction for offramp execute without CW configs
		err = offramp.BuildOffRampExecutePTB(ctx, p.log, p.client, ptb, arguments, signerAddress, addressMappings)
		if err != nil {
			p.log.Errorw("Error building OffRamp execute PTB", "error", err)
			return nil, err
		}

		return ptb, nil
	} else if function == cwConfig.CCIPCommit {
		// If it's just a commit, then we just need to get the address mappings and use the regular
		// PTB builder to build the PTB.
		addressMappings, err := offramp.GetOfframpAddressMappings(ctx, p.log, p.client, toAddress, txnConfig.PublicKey)
		if err != nil {
			p.log.Errorw("Error setting up address mappings", "error", err)
			return nil, err
		}
		// Add values from address mappings to the arguments received from core to enable the regualar
		// PTB building flow for commit.
		arguments.Args["ccip_object_ref"] = addressMappings.CcipObjectRef
		arguments.Args["state"] = addressMappings.OffRampState
		arguments.Args["clock"] = addressMappings.ClockObject
	}

	// Create a map for caching objects
	cachedArgs := make(map[string]transaction.Argument)

	// Process each command in order
	for _, cmd := range txnConfig.PTBCommands {
		cmd.PackageId = &toAddress
		// Process the command based on its type
		switch cmd.Type {
		case codec.SuiPTBCommandMoveCall:
			_, err := p.ProcessMoveCall(ctx, ptb, cmd, &arguments, &cachedArgs)
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

	return ptb, nil
}

// ProcessMoveCall handles constructing move call commands and adds it to the PTB `builder` instance.
func (p *PTBConstructor) ProcessMoveCall(
	ctx context.Context,
	builder *transaction.Transaction,
	cmd cwConfig.ChainWriterPTBCommand,
	arguments *cwConfig.Arguments,
	cachedArgs *map[string]transaction.Argument,
) (*transaction.Argument, error) {
	p.log.Debugw("Processing move call", "Command", cmd, "Args", arguments)

	// All three fields below are required for a successful move call
	if cmd.PackageId == nil {
		return nil, fmt.Errorf("missing required parameter 'PackageId' for move call PTB command")
	}
	if cmd.ModuleId == nil {
		return nil, fmt.Errorf("missing required parameter 'ModuleId' for move call PTB command")
	}
	if cmd.Function == nil {
		return nil, fmt.Errorf("missing required parameter 'Function' for move call PTB command")
	}

	// Convert package ID to Address
	packageId := models.SuiAddress(*cmd.PackageId)

	// Process arguments
	processedArgs, err := p.ProcessArgsForCommand(ctx, builder, cmd.Params, arguments, cachedArgs)
	if err != nil {
		return nil, err
	}

	processedArgTypes, err := p.ResolveGenericTypeTags(cmd.Params)
	if err != nil {
		return nil, err
	}

	p.log.Debugw("Processed Type Tags", "Type Tags", processedArgTypes)
	p.log.Debugw("Processed args", "Args", processedArgs)
	// Add the move call to the builder
	ptbArgument := builder.MoveCall(packageId, *cmd.ModuleId, *cmd.Function, processedArgTypes, processedArgs)

	return &ptbArgument, nil
}

// ProcessArgsForCommand converts parametedsr specifications into concrete arguments
func (p *PTBConstructor) ProcessArgsForCommand(
	ctx context.Context,
	builder *transaction.Transaction,
	params []codec.SuiFunctionParam,
	arguments *cwConfig.Arguments,
	cachedArgs *map[string]transaction.Argument,
) ([]transaction.Argument, error) {
	processedArgs := make([]transaction.Argument, 0, len(params))

	// if someone passed the wrapper here
	if wrap, ok := arguments.Args["Args"].(map[string]any); ok {
		// try to pull ArgTypes from the same wrapper map (or however it was passed)
		if rawAT, ok := arguments.Args["ArgTypes"]; ok {
			switch m := rawAT.(type) {
			case map[string]string:
				arguments.ArgTypes = m
			case map[string]any:
				at := make(map[string]string, len(m))
				for k, v := range m {
					if s, ok := v.(string); ok {
						at[k] = s
					}
				}
				arguments.ArgTypes = at
			}
		}
		arguments.Args = wrap
	}

	for _, param := range params {
		p.log.Debugw("Processing PTB parameter", "Param", param)

		// specify if the value is Mutable, this is used specifically for object PTB args
		isMutable := true
		if param.IsMutable != nil {
			isMutable = *param.IsMutable
		}

		// check if this is a PTB result dependency
		if param.PTBDependency != nil {
			// if the config does not specify a ResultIndex, then the dependency is
			// on the entire result of the dependee command
			if param.PTBDependency.ResultIndex == nil {
				processedArgs = append(processedArgs, transaction.Argument{
					Result: &param.PTBDependency.CommandIndex,
				})

				continue
			}

			// otherwise, we need a specific result from the dependee command
			processedArgs = append(processedArgs, transaction.Argument{
				NestedResult: &transaction.NestedResult{
					Index:       param.PTBDependency.CommandIndex,
					ResultIndex: *param.PTBDependency.ResultIndex,
				},
			})

			continue
		}

		// otherwise, check if the parameter is in the provided args
		if argRawValue, exists := arguments.Args[param.Name]; exists {
			// check if the param has already been converted and cached
			if cachedArg, exists := (*cachedArgs)[param.Name]; exists {
				processedArgs = append(processedArgs, cachedArg)
				continue
			}

			if param.Type == "object_id" {
				id, ok := argRawValue.(string)
				if !ok {
					return nil, fmt.Errorf("expected string for object id for param %s, got %T", param.Name, argRawValue)
				}
				argRawValue = id
			}

			// append to the array of args
			processedArgValue, err := p.client.(*client.PTBClient).TransformTransactionArg(ctx, builder, argRawValue, param.Type, isMutable)
			if err != nil {
				return nil, fmt.Errorf("failed to build argument for %s: %w, %s", param.Name, err, argRawValue)
			}
			processedArgs = append(processedArgs, *processedArgValue)
			// add the processed arg to the cache
			(*cachedArgs)[param.Name] = *processedArgValue

			continue
		}

		// fallback to the default value. Some assumptions:
		// - arguments with default values do NOT have type arguments
		// - objects do not have default values
		if param.DefaultValue != nil {
			value := param.DefaultValue
			if param.Type == "object_id" {
				id, ok := param.DefaultValue.(string)
				if !ok {
					return nil, fmt.Errorf("expected string for object id for param %s, got %T", param.Name, param.DefaultValue)
				}
				value = id
			}
			ptbArg, err := p.client.(*client.PTBClient).TransformTransactionArg(ctx, builder, value, param.Type, isMutable)
			if err != nil {
				return nil, fmt.Errorf("failed to build default value for %s: %w", param.Name, err)
			}
			// append to the array of args
			processedArgs = append(processedArgs, *ptbArg)

			continue
		}

		// Value not found for required param
		if param.Required {
			return nil, fmt.Errorf("required parameter %s has no value", param.Name)
		}

		// append an empty argument since it is not required and no value found
		processedArgs = append(processedArgs, transaction.Argument{})
	}

	return processedArgs, nil
}

// FetchPrereqObjects fetches each pre-requisite object and its details, then populates the args map with its values
func (p *PTBConstructor) FetchPrereqObjects(ctx context.Context, prereqObjects []cwConfig.PrerequisiteObject, args *map[string]any, ownerFallback *string) error {
	for _, prereq := range prereqObjects {
		// set the owner fallback if the ownerId is not provided
		if prereq.OwnerId == nil {
			if ownerFallback == nil {
				return fmt.Errorf("ownerId or ownerFallback required for pre-requisite object %s", prereq.Name)
			}

			prereq.OwnerId = ownerFallback
		}
		// fetch owned objects
		ownedObjects, err := p.client.ReadOwnedObjects(ctx, *prereq.OwnerId, nil)
		if err != nil {
			return err
		}

		// check each returned object
		for _, ownedObject := range ownedObjects {
			// object tag matches
			if ownedObject.Data != nil && ownedObject.Data.Type != "" && strings.Contains(ownedObject.Data.Type, prereq.Tag) {
				p.log.Debugw("Found pre-requisite object", "Object", ownedObject.Data, "Prereq", prereq)
				// object must be parsed and its keys added to the args map
				if prereq.SetKeys {
					// parse the object into a map
					if ownedObject.Data.Content != nil && ownedObject.Data.Content.Fields != nil {
						// add each key and value to the args map
						for key, value := range ownedObject.Data.Content.Fields {
							(*args)[key] = value
						}
					}
				} else {
					// add the object id to the args map
					(*args)[prereq.Name] = ownedObject.Data.ObjectId
				}
			}
		}
	}

	return nil
}
