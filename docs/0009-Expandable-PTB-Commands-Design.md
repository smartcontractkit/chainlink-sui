## ADR: Enhancing ChainWriter for Expandable PTB Commands

- [Design Document: Enhancing ChainWriter for Expandable PTB Commands](#design-document-enhancing-chainwriter-for-expandable-ptb-commands)
    - [1. Introduction](#1-introduction)
    - [2. Goals](#2-goals)
    - [3. Accepted Design: Specific hardcoded logic for each flow](#3-accepted-design-specific-hardcoded-logic-for-each-flow)
        - [3.1 Core Idea](#31-core-idea)
        - [3.2 Structure Modifications](#32-structure-modifications)
        - [3.3 ChainWriter Instance Enhancements](#33-chainwriter-instance-enhancements)
        - [3.4 PTB Generation Logic](#34-ptb-generation-logic)
        - [3.5 Implementing/Handling Expansion Logic](#35-implementinghandling-expansion-logic)
        - [3.6 Pros](#36-pros)
        - [3.7 Cons](#37-cons)
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

This document outlines design options for achieving this functionality.

**2. Goals**

*   Enable specific PTB commands within a `ChainWriterFunction` configuration to be marked as expandable.
*   Provide a clear and flexible mechanism to define the logic for expanding these commands.
*   Ensure that PTB argument dependencies from subsequent commands to the outputs of expanded commands are correctly resolved.
*   Allow dynamic provision of arguments for each instance of an expanded command, with clear sourcing.
*   Maintain backward compatibility for existing non-expandable command configurations.

**3. Accepted Design: Specific hardcoded logic for each flow**

This design makes the following assumptions:
- The only expandable flow is the offRamp flow.
- The expander function only needs the OCR report to be able to expand the PTB.
- The signature of the `ccip_receive` function receives a hot potato and returns another hot potato.
- The PTB will only add the receiver call command if the receiver is registered.
- Sequential execution of token pool commands (even though they are independent and parallelizable in nature)
  - Each token pool command mutates and passes the hot potato to the next command
  - Same thing happens with the receiver call command

**3.1 Core Idea**

This option introduces specific hardcoded logic for each flow. It is inspired by a similar approach that was followed by the Solana team and it can be found [here](https://github.com/smartcontractkit/chainlink-solana/blob/3c6bdae1a2f72144becfd63a7332b5c982ed687b/pkg/solana/chainwriter/transform_registry.go#L30). This makes the implementation effort simpler by sacrificing the flexibility of the design. We will only use this approach for the offRamp flow since it is the only flow that is expandable.

This option introduces an expander function for the OffRamp flow that will transform the input arguments to the expected format for the expander. Since we know the flow that the PTB is modeling, we also know which command(s) can be expandable and what data structure to use for the expansion logic.

The offRamp flow can be defined as follows;

1. Call OffRamp `init` function to create a hot potato;
2. If the Report contains token pool messages do:
   2.1. Call `lock_or_burn` function for each token pool message. Ensure each token pool uses a refernece to the hot potato from the previous command;
3. If the report contains a receiver, and if that the receiver is registered, do:
   3.1. Call function to extract the Any2SuiMessage from the hot potato;
   3.2  Call the receiver `ccip_receive` implementation with the Any2SuiMessage;
4. Call OffRamp function, lets call it `offramp_finalize`, to finalise the flow and pass it the hot potato;

**3.2 Structure Modifications**

This solution does not require any new data structures but the following functions:

1. `FindPTExpander` that will return the PTB expander for a given flow.
2. Implement a `PTExpander` function that will have the following signature:

```go
func PTBExpander(args any, config chainwriter.ChainWriterConfig) (ptbCommand chainwriter.ChainWriterPTBCommand, updatedArgs any, err error)
```

**3.3 ChainWriter Instance Enhancements**

No changes are needed to the ChainWriter struct. PTBs that model known flows will need to have a specific name that will be used to find the expander function.

**3.4 PTB Generation Logic**

The OffRamp ChainWriter config will only have two known PTB commands: the first one (let's call it `offramp_init` function) and the last one (let's call it `offramp_finalize` function). All the other commands will be derived from the OCR report. The `PTBExpander` function will receive the OCR report and the ChainWriter config and will generate the final PTB commands to that specific report execution.

**3.5 Implementing/Handling Expansion Logic**

```go

type OffRampCCIPExecuteArgs struct {
    ocrReportInfo ExecuteReportInfo
    receiverAddress string
    data []byte
    destinationChainSelector uint64
}

func OffRampPTBExpander(args OffRampCCIPExecuteArgs, config chainwriter.ChainWriterConfig) (ptbCommands []chainwriter.ChainWriterPTBCommand, updatedArgs any, err error) {
    for idx, tokenAmount := range args.ocrReportInfo.tokenAmounts {
        // dynamically generate the ptb command to call the respective token pool
    }

    // update the args with the new values
    // update the hot potato reference index of the last token pool command.
}
```

**3.6 Pros**

*   Simpler to implement when compared to the alternative design
*   No need to create any new data structures or interfaces

**3.7 Cons**

*   Less flexible than the alternative design

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

Accepted Design

*   **Function needs to be marked as repeatable:** Achieved by naming convention - functions that need expansion are identified by their name pattern and mapped to specific `PTBExpander` functions.
*   **Check downstream PTB argument dependencies:** Dependencies are handled within the `PTBExpander` function itself, which has full control over the expanded commands and their argument mapping.
*   **Dynamically provide values:** The `PTBExpander` function receives the original arguments and can transform them as needed, returning both the expanded commands and updated arguments.
*   **ChainWriter needs expansion logic:** Achieved through the `FindPTExpander` function that maps function names to their corresponding expander implementations.
*   **Struct for repeatability and expansion logic:** No additional structures needed - expansion logic is encapsulated in the `PTBExpander` function signature and implementation.

Alternative Design

*   **Function needs to be marked as repeatable:** Achieved by `ChainWriterPTBCommand` having `ExpansionSetup` with a valid `HandlerKey` pointing to a registered `CommandExpander`.
*   **Check downstream PTB argument dependencies:** The two-stage system with `commandMetadata` resolves dependencies to absolute PTB indices, targeting the last instance of an expanded block.
*   **Dynamically provide values:** The `CommandExpander.ExpandCommands` method returns `ExpandedCommandArguments` containing values and type tags for parameters of the generated commands. These are sourced from global `chainWriterInputArgs`, `reportData`, or handler-internal logic.
*   **ChainWriter needs expansion logic:** Achieved with the `CommandExpander` interface and a registry in `ChainWriter`.
*   **Struct for repeatability and expansion logic:** `ChainWriterPTBCommand.ExpansionSetup` links to the expander.

