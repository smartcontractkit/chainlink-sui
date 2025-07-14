## ADR: Enhancing ChainWriter for Expandable PTB Commands

- [Design Document: Enhancing ChainWriter for Expandable PTB Commands](#design-document-enhancing-chainwriter-for-expandable-ptb-commands)
    - [1. Introduction](#1-introduction)
    - [2. Goals](#2-goals)
    - [3. Implemented Design: Generic PTBExpander Interface](#3-implemented-design-generic-ptbexpander-interface)
        - [3.1 Core Idea](#31-core-idea)
        - [3.2 Generic Interface Structure](#32-generic-interface-structure)
        - [3.3 Concrete OffRamp Implementation](#33-concrete-offramp-implementation)
        - [3.4 Address Mapping System](#34-address-mapping-system)
        - [3.5 PTB Expansion Logic](#35-ptb-expansion-logic)
        - [3.6 Token Pool and Receiver Handling](#36-token-pool-and-receiver-handling)
        - [3.7 Pros](#37-pros)
        - [3.8 Cons](#38-cons)
    - [4. Alternative Design: CommandExpander Interface Approach](#4-alternative-design-commandexpander-interface-approach)
        - [4.1 Core Idea](#41-core-idea)
        - [4.2 Structure Modifications](#42-structure-modifications)
        - [4.3 ChainWriter Instance Enhancements](#43-chainwriter-instance-enhancements)
        - [4.4 PTB Generation Logic Update](#44-ptb-generation-logic-update)
        - [4.5 Implementing Command Expanders](#45-implementing-command-expanders)
    - [5. Key Requirements Addressed](#5-key-requirements-addressed)

**1. Introduction**

The ChainWriter component is responsible for constructing Sui Programmable Transaction Blocks (PTBs) based on predefined configurations. Currently, each command in the configuration maps to a single command in the PTB. This document proposes enhancements to ChainWriter to support "expandable" commands. An expandable command is a command defined once in the configuration that can be dynamically repeated or instantiated multiple times during PTB generation. This expansion will be driven by runtime logic, potentially using external data like a "report," allowing for more flexible and dynamic PTB construction.

A key use case is processing multiple instances of a similar operation, such as handling multiple token pool calls on the offRamp flow.

This document outlines the implemented design for achieving this functionality.

**2. Goals**

*   Enable specific PTB commands within a `ChainWriterFunction` configuration to be marked as expandable.
*   Provide a clear and flexible mechanism to define the logic for expanding these commands.
*   Ensure that PTB argument dependencies from subsequent commands to the outputs of expanded commands are correctly resolved.
*   Allow dynamic provision of arguments for each instance of an expanded command, with clear sourcing.
*   Maintain backward compatibility for existing non-expandable command configurations.

**3. Implemented Design: Generic PTBExpander Interface**

This design was implemented to provide a balance between flexibility and simplicity. It uses Go generics to provide type-safe PTB expansion operations while maintaining specific implementations for known flows.

**3.1 Core Idea**

The implemented approach introduces a generic `PTBExpander` interface that can be specialized for different types of operations. This design provides:

1. **Type Safety**: Uses Go generics to ensure compile-time type safety for different expansion operations
2. **Flexibility**: Allows for different argument and result types for different expansion scenarios
3. **Simplicity**: Maintains a single `Expand` method interface while supporting various operation types
4. **Extensibility**: New expansion types can be added by implementing the generic interface

The core interface is defined as:

```go
type PTBExpander[T any, R any] interface {
    Expand(
        ctx context.Context,
        lggr logger.Logger,
        args T,
        signerPublicKey []byte,
    ) (R, error)
}
```

**3.2 Generic Interface Structure**

The generic interface allows for type-safe implementations:

```go
// PTBExpander defines a generic interface for expanding PTB commands
type PTBExpander[T any, R any] interface {
    Expand(
        ctx context.Context,
        lggr logger.Logger,
        args T,
        signerPublicKey []byte,
    ) (R, error)
}
```

This interface can be specialized for different operations:
- `PTBExpander[OffRampPTBArgs, OffRampPTBResult]` for OffRamp operations
- `PTBExpander[OnRampPTBArgs, OnRampPTBResult]` for potential OnRamp operations
- Other specialized implementations as needed

**3.3 Concrete OffRamp Implementation**

The `SuiPTBExpander` implements the generic interface specifically for OffRamp operations:

```go
type SuiPTBExpander struct {
    lggr            logger.Logger
    ptbClient       client.SuiPTBClient
    AddressMappings map[string]string
}

// Implements PTBExpander[OffRampPTBArgs, OffRampPTBResult]
func (s *SuiPTBExpander) Expand(
    ctx context.Context,
    lggr logger.Logger,
    args OffRampPTBArgs,
    signerPublicKey []byte,
) (OffRampPTBResult, error)
```

**Input and Output Types:**

```go
// OffRampPTBArgs represents arguments for OffRamp PTB expansion
type OffRampPTBArgs struct {
    ExecArgs   SuiOffRampExecCallArgs
    PTBConfigs *config.ChainWriterFunction
}

// OffRampPTBResult represents the result of OffRamp PTB expansion
type OffRampPTBResult struct {
    PTBCommands []config.ChainWriterPTBCommand
    UpdatedArgs map[string]any
    TypeArgs    map[string]string
}
```

**3.4 Address Mapping System**

A critical component of the implementation is the address mapping system that discovers and maintains references to on-chain objects:

```go
func SetupAddressMappings(
    ctx context.Context,
    lggr logger.Logger,
    ptbClient client.SuiPTBClient,
    offRampPackageId string,
    publicKey []byte,
) (map[string]string, error)
```

This function performs a 4-step discovery process:
1. Uses the OffRamp package ID to discover the CCIP package ID
2. Reads owned objects to locate the OffRamp state pointer and extract the state address
3. Reads CCIP package objects to find the CCIP object reference and owner capability addresses
4. Assembles a complete address mapping required for PTB operations

The resulting mappings include:
- `ccipPackageId`: Main CCIP package identifier
- `ccipObjectRef`: Reference to the main CCIP state object
- `ccipOwnerCap`: Owner capability object for privileged operations
- `clockObject`: Sui system clock object (fixed at 0x6)
- `offRampPackageId`: The OffRamp package identifier
- `offRampState`: The OffRamp state object address

**3.5 PTB Expansion Logic**

The PTB expansion process follows a structured approach:

1. **Base Commands**: Start with predefined commands from configuration:
   - `init_execute`: Initializes the OffRamp execution
   - `finish_execute`: Finalizes the execution

2. **Dynamic Token Pool Commands**: Generate commands for each token transfer:
   - Extract token amounts from the execution report
   - Look up token pool information for each token
   - Generate PTB commands for token pool operations

3. **Receiver Call Commands**: Generate commands for message receivers:
   - Check if receivers are registered
   - Generate receiver call commands for valid receivers
   - Handle message data passing

4. **Command Sequencing**: Ensure proper ordering and dependencies:
   - Token pool commands are inserted between init and finish
   - Receiver commands follow token pool commands
   - PTB dependencies are updated to maintain proper references

**3.6 Token Pool and Receiver Handling**

**Token Pool Operations:**
- Query token pool information using `get_pool_infos` function
- Generate `release_or_mint` commands for each token pool
- Handle token type arguments and state addresses
- Maintain proper PTB dependency chains

**Receiver Operations:**
- Parse receiver addresses in `packageID::moduleID::functionName` format
- Check receiver registration using `is_registered_receiver`
- Generate receiver call commands with proper message data
- Handle receiver parameters and dependencies

**3.7 Pros**

*   **Type Safety**: Go generics provide compile-time type safety
*   **Flexibility**: Can support multiple expansion types through generic interface
*   **Simplicity**: Single method interface is easy to understand and implement
*   **Extensibility**: New expansion types can be added without changing existing code
*   **Maintainability**: Clear separation between generic interface and concrete implementations
*   **Performance**: Efficient address discovery and caching system

**3.8 Cons**

*   **Generic Complexity**: Requires understanding of Go generics
*   **Single Implementation**: Currently only OffRamp expansion is implemented
*   **Address Discovery Overhead**: Initial setup requires multiple on-chain reads

**4. Alternative Design: CommandExpander Interface Approach**

**4.1 Core Idea**

The core idea is to define a `CommandExpander` interface that explicitly separates the planning and generation phases of command expansion. `ChainWriterPTBCommand` will have an `ExpansionSetup` to link to a registered `CommandExpander` implementation. The PTB generation logic will be updated to use this interface. We need a planning phase to determine the number of commands that will be generated and the PTB indices of the commands. This is needed in case commands receive arguments from outputs of previous commands. This approach allows for complex expansion logic.

**4.2 Structure Modifications**

We'll define `CommandExpander` interface, `ExpandedCommandArguments` struct, and `CommandExpansionSetup` struct.

```go
package codec // Or your relevant package for these type definitions

// ExpandedCommandArguments holds the arguments and their type tags for commands generated by an expansion handler.
type ExpandedCommandArguments struct {
    Args     map[string]interface{} // Map of argument name to value
    ArgTypes map[string]string      // Map of argument name to Sui type tag (for generic types)
}

// CommandExpander defines the interface for logic that can expand a PTB command.
type CommandExpander interface {
    // PlanLayout determines how many actual PTB commands will be generated from the base command.
    // It's called during Stage 1 of PTB generation.
    PlanLayout(
        baseCommand ChainWriterPTBCommand,
        expansionSetup *CommandExpansionSetup,
        chainWriterInputArgs map[string]interface{}, // Global input args
        reportData interface{},                      // External data
        originalConfigIndex int,
        currentPtbOffset int, // Current offset in PTB being built during planning
    ) (numExpanded int, err error)

    // ExpandCommands generates the actual PTB command instances and their arguments.
    // It's called during Stage 2 of PTB generation.
    ExpandCommands(
        baseCommand ChainWriterPTBCommand,
        expansionSetup *CommandExpansionSetup,
        chainWriterInputArgs map[string]interface{}, // Global input args
        reportData interface{},                      // External data
        originalConfigIndex int,
        actualStartIndexInPtb int, // The calculated start index for this block from Stage 1
    ) (
        expandedCommands []ChainWriterPTBCommand,
        instanceArguments *ExpandedCommandArguments, // Arguments for the generated expandedCommands
        err error,
    )
}

// CommandExpansionSetup holds the configuration for how a PTB command can be expanded.
type CommandExpansionSetup struct {
    // HandlerKey is a string key to look up a registered CommandExpander
    // implementation in the ChainWriter instance. Mandatory for expandable commands.
    HandlerKey string `json:"handlerKey"`
    // Optional: Context or static parameters for the handler, can be decoded by the handler.
    // HandlerContext map[string]interface{} `json:"handlerContext,omitempty"`
}

type ChainWriterPTBCommand struct {
    Type        SuiPTBCommandType // e.g., SuiPTBCommandMoveCall
    PackageId   *string
    ModuleId    *string
    Function    *string
    Params      []SuiFunctionParam
    // ... other existing fields ...

    // ExpansionSetup defines if and how this command should be expanded.
    // If nil, the command is not expandable.
    ExpansionSetup *CommandExpansionSetup `json:"expansionSetup,omitempty"`
}

type SuiFunctionParam struct {
    Name          string
    Type          string // e.g., "object_id", "u64", "ptb_dependency", "vector<u8>"
    Required      bool
    IsMutable     *bool `json:"isMutable,omitempty"`
    IsGeneric     bool  `json:"isGeneric,omitempty"`
    PTBDependency *PTBCommandDependency `json:"ptbDependency,omitempty"`
}

type PTBCommandDependency struct {
    CommandIndex int  `json:"commandIndex"` // Refers to index in original config; resolved to actual PTB index.
    ResultIndex  *int `json:"resultIndex,omitempty"` // For tuple results, defaults to 0.
}
```

**4.3 ChainWriter Instance Enhancements**

The `ChainWriter` will need a registry for `CommandExpander` implementations:

```go
type ChainWriter struct {
    // ... other fields like configurations, SUI client ...
    commandExpanders map[string]CommandExpander // Registry for CommandExpander implementations
}

func NewChainWriter(/*... deps ...*/) *ChainWriter {
    return &ChainWriter{
        commandExpanders: make(map[string]CommandExpander),
        // ... initialize other fields ...
    }
}

func (cw *ChainWriter) RegisterCommandExpander(key string, expander CommandExpander) {
    cw.commandExpanders[key] = expander
}
``` 

**4.4 PTB Generation Logic Update**

The two-stage process is maintained, now utilizing the `CommandExpander` interface.

**Stage 1: Planning the Layout (Determine Expansion Counts & PTB Indices)**

*   **Goal:** Determine how many PTB commands each configured command will produce and their final PTB indices.
*   **How:**
    1.  Iterate through each `configuredCommand` in the original configuration.
    2.  If `configuredCommand.ExpansionSetup` is not nil:
        *   Retrieve the `CommandExpander` using `ExpansionSetup.HandlerKey`.
        *   Call `expander.PlanLayout(...)`. The returned `numExpanded` is used.
    3.  If not expandable, `numExpanded` is 1.
    4.  Maintain `commandMetadata` (mapping original config index to `ActualStartIndex` in PTB and `NumInstances` which is `numExpanded`).
    5.  Calculate `offset` based on `NumInstances` to correctly position subsequent commands.

**Stage 2: Building the Commands (Generate Commands & Resolve Dependencies)**

*   **Goal:** Create actual PTB command instances with resolved parameters.
*   **How:**
    1.  Iterate through the original configuration again, using `commandMetadata`.
    2.  If a `configuredCommand` has `ExpansionSetup`:
        *   Retrieve the `CommandExpander`.
        *   Get `actualStartIndexInPtb` from `commandMetadata`.
        *   Call `expander.ExpandCommands(..., actualStartIndexInPtb)`.
        *   The returned `expandedCommands` are added to the final PTB list. The associated `instanceArguments` are stored with these commands for parameter value lookup.
        *   Ensure each command in `expandedCommands` has `ExpansionSetup = nil`.
    3.  If not expandable, prepare the single command instance and add it to the final list.
    4.  **For every command instance generated (expanded or not):**
        *   **Resolve `PTBDependency` parameters:** Update `PTBDependency.CommandIndex` to the absolute PTB index using `commandMetadata` (pointing to the last instance if the depended-upon command was expanded).
        *   **Resolve direct value parameters:**
            *   If from an expansion (has `instanceArguments`): Source value from `instanceArguments.Args[param.Name]` and type tag from `instanceArguments.ArgTypes[param.Name]` if generic.
            *   If not expanded: Source value from global `chainWriterInputArgs.Args[param.Name]` and type tag from `chainWriterInputArgs.ArgTypes[param.Name]` if generic.
    5.  Collect all fully formed commands into the final PTB list.

**4.5 Implementing Command Expanders**

Developers provide implementations of the `CommandExpander` interface.

```go
// Hypothetical struct for input data to the expander
type TargetTokenInfo struct {
    RefObjectID        string `json:"refObjectId"`
    StateObjectID      string `json:"stateObjectId"`
    TokenObjectID      string `json:"tokenObjectId"`
    TokenTypeTag       string `json:"tokenTypeTag"`
    // Optional: PackageId, ModuleId, FunctionName if they vary per instance
}

// Example: LockOrBurnExpander
type LockOrBurnExpander struct {
    // Can have its own configuration or dependencies if needed, injected on creation
}

// (Assume deepCopy and deepCopySlice utility functions are available)

func (e *LockOrBurnExpander) PlanLayout(
    baseCommand codec.ChainWriterPTBCommand,
    expansionSetup *codec.CommandExpansionSetup,
    chainWriterInputArgs map[string]interface{},
    reportData interface{},
    originalConfigIndex int,
    currentPtbOffset int,
) (numExpanded int, err error) {
    tokensToLockVal, ok := chainWriterInputArgs["tokens_for_lock_or_burn"]
    if !ok {
        return 0, fmt.Errorf("[PlanLayout] 'tokens_for_lock_or_burn' missing in chainWriterInputArgs")
    }
    tokensToLock, ok := tokensToLockVal.([]TargetTokenInfo)
    if !ok {
        // Or attempt a more flexible decode, e.g., from []map[string]interface{}
        return 0, fmt.Errorf("[PlanLayout] 'tokens_for_lock_or_burn' not []TargetTokenInfo")
    }
    return len(tokensToLock), nil
}

func (e *LockOrBurnExpander) ExpandCommands(
    baseCommand codec.ChainWriterPTBCommand,
    expansionSetup *codec.CommandExpansionSetup,
    chainWriterInputArgs map[string]interface{},
    reportData interface{},
    originalConfigIndex int,
    actualStartIndexInPtb int,
) (
    expandedCmds []codec.ChainWriterPTBCommand,
    instanceArgs *codec.ExpandedCommandArguments,
    err error,
) {
    tokensToLockVal, ok := chainWriterInputArgs["tokens_for_lock_or_burn"]
    if !ok { // Should have been caught in PlanLayout, but good to check
        return nil, nil, fmt.Errorf("[ExpandCommands] 'tokens_for_lock_or_burn' missing")
    }
    tokensToLock, ok := tokensToLockVal.([]TargetTokenInfo)
    if !ok {
        return nil, nil, fmt.Errorf("[ExpandCommands] 'tokens_for_lock_or_burn' not []TargetTokenInfo")
    }

    N := len(tokensToLock)
    if N == 0 {
        return []codec.ChainWriterPTBCommand{}, nil, nil
    }

    actualInstanceArgs := &codec.ExpandedCommandArguments{
        Args:     make(map[string]interface{}),
        ArgTypes: make(map[string]string),
    }

    for i := 0; i < N; i++ {
        tokenInfo := tokensToLock[i]
        instanceCmd := deepCopy(baseCommand) // Deep copy base command
        instanceCmd.ExpansionSetup = nil       // Expanded instance is not further expandable

        originalBaseParams := deepCopySlice(baseCommand.Params)

        for j, baseParamConfig := range originalBaseParams {
            instanceParam := &instanceCmd.Params[j]
            uniqueInstanceParamName := fmt.Sprintf("%s_%d", baseParamConfig.Name, i)
            instanceParam.Name = uniqueInstanceParamName

            if instanceParam.PTBDependency != nil {
                continue
            }

            switch baseParamConfig.Name {
            case "ref":
                actualInstanceArgs.Args[uniqueInstanceParamName] = tokenInfo.RefObjectID
            case "state":
                actualInstanceArgs.Args[uniqueInstanceParamName] = tokenInfo.StateObjectID
            case "c": // Coin
                actualInstanceArgs.Args[uniqueInstanceParamName] = tokenInfo.TokenObjectID
                if instanceParam.IsGeneric {
                    actualInstanceArgs.ArgTypes[uniqueInstanceParamName] = tokenInfo.TokenTypeTag
                }
            case "remote_chain_selector":
                if sharedVal, e := chainWriterInputArgs["shared_remote_chain_selector"]; e {
                    actualInstanceArgs.Args[uniqueInstanceParamName] = sharedVal
                } // else error or default
            case "clock":
                if clockVal, e := chainWriterInputArgs["global_clock_object"]; e {
                    actualInstanceArgs.Args[uniqueInstanceParamName] = clockVal
                } // else error or default
            }
        }
        expandedCmds = append(expandedCmds, instanceCmd)
    }
    return expandedCmds, actualInstanceArgs, nil
}

// Registration:
// lockExpander := &LockOrBurnExpander{}
// chainWriter.RegisterCommandExpander("LockOrBurnExpanderKey", lockExpander)

```
The `chainWriterInputArgs` for `SubmitTransaction` would provide data like `tokens_for_lock_or_burn: []TargetTokenInfo{...}`.

**5. Key Requirements Addressed**

**Implemented Design (Generic PTBExpander Interface)**

*   **Function needs to be marked as repeatable:** Achieved through the generic `PTBExpander` interface. Functions requiring expansion are handled by specific implementations that determine expansion logic based on runtime data (e.g., token amounts in OffRamp reports).

*   **Check downstream PTB argument dependencies:** Dependencies are resolved during PTB generation in the `getOffRampPTB` method. The implementation:
    - Tracks command indices as PTB commands are generated
    - Updates `PTBDependency.CommandIndex` values to maintain proper references
    - Ensures the final `finish_execute` command references the last generated token pool or receiver command

*   **Dynamically provide values:** The `Expand` method processes runtime data to generate:
    - Token pool arguments from `getTokenPoolByTokenAddress` lookups
    - Receiver arguments from message processing
    - Type arguments for generic Move function calls
    - Address mappings for on-chain object references

*   **ChainWriter needs expansion logic:** Achieved through:
    - The generic `PTBExpander[T, R]` interface that can be specialized for different operations
    - Concrete implementations like `SuiPTBExpander` for OffRamp operations
    - Address mapping discovery system in `SetupAddressMappings`
    - Integration with existing ChainWriter PTB generation flow

*   **Struct for repeatability and expansion logic:** Provided through:
    - `OffRampPTBArgs` for input parameters including execution arguments and PTB configurations
    - `OffRampPTBResult` for output including generated PTB commands, updated arguments, and type arguments
    - `TokenPool` struct for token pool information and expansion logic
    - Address mapping system for maintaining on-chain object references

**Additional Benefits of Implemented Design:**

*   **Address Discovery:** Automatic discovery and caching of critical on-chain addresses reduces manual configuration
*   **Receiver Filtering:** Automatic filtering of registered receivers prevents failed transactions
*   **Type Safety:** Go generics ensure compile-time type checking for different expansion scenarios
*   **Extensibility:** New expansion types can be added by implementing the generic interface with different type parameters
*   **Maintainability:** Clear separation of concerns between generic interface and specific implementations

**Integration with ChainWriter:**

The implemented design integrates with the existing ChainWriter by:
1. Using the existing `ChainWriterFunction` configuration as input
2. Generating compatible `ChainWriterPTBCommand` structures
3. Providing argument and type argument mappings that work with existing PTB generation logic
4. Maintaining backward compatibility with non-expandable commands

This approach successfully addresses all original requirements while providing a clean, extensible architecture for future expansion types.

