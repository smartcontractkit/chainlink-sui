# Chainlink Sui Relayer

- [Chainlink Sui Relayer](#chainlink-sui-relayer)
  - [Relayer Configuration](#relayer-configuration)
  - [Initializing the Relayer](#initializing-the-relayer)
  - [Creating a ChainReader](#creating-a-chainreader)
  - [Creating a ChainWriter](#creating-a-chainwriter)
  - [Programmable Transaction Blocks (PTBs)](#programmable-transaction-blocks-ptbs)
      - [Configuring a ChainWriter to use PTBs](#configuring-a-chainwriter-to-use-ptbs)
      - [Using PTB Dependencies](#using-ptb-dependencies)
    - [Calling a PTB in a ChainWriter instance](#calling-a-ptb-in-a-chainwriter-instance)
  - [Closing Resources](#closing-resources)


The Chainlink Sui integration provides a relayer plugin that enables communication with the Sui blockchain. The relayer offers two main components:

1. **ChainReader**: Reads data from the Sui blockchain (querying objects, calling view functions, and listening for events)
2. **ChainWriter**: Writes data to the Sui blockchain (submitting transactions to call smart contract functions)

## Relayer Configuration

To launch the relayer, you need a TOML configuration file that specifies the Sui chain and node details:

```toml
[[Chains]]
ChainID = '0x1'  # The Sui network ID
Enabled = true

# Transaction settings
BroadcastChanSize = 4096     # Size of the broadcast channel buffer
ConfirmPollPeriod = '500ms'  # Time between transaction confirmation checks
MaxConcurrentRequests = 5    # Maximum number of concurrent RPC requests
TransactionTimeout = '10s'   # Timeout for transaction requests
NumberRetries = 5            # Number of retries for failed requests
GasLimit = 10000000          # Maximum gas limit for transactions
RequestType = 'WaitForLocalExecution'  # Transaction execution mode (WaitForLocalExecution or WaitForEffectsCert)

# Transaction Manager settings
[TransactionManager]
BroadcastChanSize = 100        # Size of the broadcast channel buffer
ConfirmPollSecs = 2           # Time between transaction confirmation checks
DefaultMaxGasAmount = 200000  # Default maximum gas amount for transactions
MaxSimulateAttempts = 5       # Maximum number of simulation attempts
MaxSubmitRetryAttempts = 10   # Maximum number of submission retry attempts
MaxTxRetryAttempts = 5        # Maximum number of transaction retry attempts
PruneIntervalSecs = 14400     # Interval for pruning old transactions (4 hours)
PruneTxExpirationSecs = 7200  # Age threshold for pruning transactions (2 hours)
SubmitDelayDuration = 3       # Delay between retries in seconds
TxExpirationSecs = 10         # Transaction expiration time in seconds

# Node configurations
[[Chains.Nodes]]
Name = 'sui-node-1'
URL = 'http://localhost:9000'  # For local development
SolidityURL = 'http://localhost:9000'

# Optional backup nodes
[[Chains.Nodes]]
Name = 'sui-node-2'
URL = 'https://sui-rpc-testnet.example.com'
SolidityURL = 'https://sui-rpc-testnet.example.com'
```

## Initializing the Relayer

Here's how to initialize the Sui relayer with your configuration:

```go
import (
    "github.com/smartcontractkit/chainlink-common/pkg/logger"
    "github.com/smartcontractkit/chainlink-common/pkg/types/core"
    "github.com/smartcontractkit/chainlink-sui/relayer/plugin"
    "github.com/smartcontractkit/chainlink-sui/relayer/keystore"
)

func initRelayer() (*plugin.SuiRelayer, error) {
    log := logger.Sugar() // Initialize your logger
    
    // Load configuration from TOML
    cfg := &plugin.TOMLConfig{
        // Load from file or set programmatically
    }
    
    // Initialize keystore
    ks, err := keystore.NewSuiKeystore(log, "")
    if err != nil {
        return nil, err
    }
    
    // Create the relayer
    relayer, err := plugin.NewRelayer(cfg, log, ks)
    if err != nil {
        return nil, err
    }
    
    // Start the relayer
    if err := relayer.Start(context.Background()); err != nil {
        return nil, err
    }
    
    return relayer, nil
}
```

## Creating a ChainReader

The ChainReader allows you to interact with the Sui blockchain to read data. Here's how to configure and use it:

```go
import (
    "encoding/json"
    "github.com/smartcontractkit/chainlink-common/pkg/types"
    "github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
    "github.com/smartcontractkit/chainlink-sui/relayer/chainreader"
    "github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func createChainReader(relayer *plugin.SuiRelayer) (types.ContractReader, error) {
    // Define the ChainReader configuration
    config := chainreader.ChainReaderConfig{
        Modules: map[string]*chainreader.ChainReaderModule{
            "counter": {
                Name: "counter",
                Functions: map[string]*chainreader.ChainReaderFunction{
                    "get_count": {
                        Name:          "get_count",
                        // This is required for reader/view functions that leverage DevInspect. 
                        // No gas will be paid for these transactions but the spec requires a signer address.
                        SignerAddress: "0x123...",
                        Params: []codec.SuiFunctionParam{
                            {
                                Type:         "address",
                                Name:         "counter_id",
                                DefaultValue: "0x456...", // Counter object ID
                                Required:     true,
                            },
                        },
                    },
                },
                Events: map[string]*chainreader.ChainReaderEvent{
                    "counter_incremented": {
                        Name:      "counter_incremented",
                        EventType: "CounterIncremented",
                    },
                },
            },
        },
    }
    
    // Serialize the configuration
    configBytes, err := json.Marshal(config)
    if err != nil {
        return nil, err
    }
    
    // Create the ChainReader
    reader, err := relayer.NewContractReader(context.Background(), configBytes)
    if err != nil {
        return nil, err
    }
    
    return reader, nil
}

func readCounterValue(reader types.ContractReader, packageId, counterObjectId string) (uint64, error) {
    // Create the read identifier (packageId-moduleName-objectId)
    readIdentifier := packageId + "-counter-" + counterObjectId
    
    var counterValue uint64
    err := reader.GetLatestValue(
        context.Background(),
        readIdentifier,
        primitives.Finalized,
        struct{}{}, // No parameters needed for object reads
        &counterValue,
    )
    if err != nil {
        return 0, err
    }
    
    return counterValue, nil
}

func callGetCountFunction(reader types.ContractReader, packageId string, counterObjectId string) (uint64, error) {
    // Create the read identifier (packageId-moduleName-functionName)
    readIdentifier := packageId + "-counter-get_count"
    
    var counterValue uint64
    err := reader.GetLatestValue(
        context.Background(),
        readIdentifier,
        primitives.Finalized,
        map[string]any{
            "counter_id": counterObjectId,
        },
        &counterValue,
    )
    if err != nil {
        return 0, err
    }
    
    return counterValue, nil
}
```

## Creating a ChainWriter

The ChainWriter allows you to submit transactions to the Sui blockchain:

```go
import (
    "encoding/json"
    "github.com/smartcontractkit/chainlink-common/pkg/types"
    "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
    "github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func createChainWriter(relayer *plugin.SuiRelayer) (types.ContractWriter, error) {
    // Define the ChainWriter configuration
    config := chainwriter.ChainWriterConfig{
        Modules: map[string]*chainwriter.ChainWriterModule{
            "counter": {
                Name:     "counter",
                ModuleID: "0x123...", // Package ID
                Functions: map[string]*chainwriter.ChainWriterFunction{
                    "increment": {
                        Name:      "increment",
                        PublicKey: []byte{/* Public key bytes */},
                        Params: []codec.SuiFunctionParam{
                            {
                                Name:     "counter",
                                Type:     "address",
                                Required: true,
                            },
                        },
                    },
                },
            },
        },
    }
    
    // Serialize the configuration
    configBytes, err := json.Marshal(config)
    if err != nil {
        return nil, err
    }
    
    // Create the ChainWriter
    writer, err := relayer.NewContractWriter(context.Background(), configBytes)
    if err != nil {
        return nil, err
    }
    
    return writer, nil
}

func incrementCounter(writer types.ContractWriter, counterObjectId string) error {
    // Generate a unique transaction ID
    txID := uuid.New().String()
    
    // Submit the transaction
    err := writer.SubmitTransaction(
        context.Background(),
        "counter", // Module name
        "increment", // Function name
        map[string]any{
            "counter": counterObjectId,
        },
        txID, // Transaction ID
        "", // To address (not used in Sui)
        &commonTypes.TxMeta{GasLimit: 10000000}, // Transaction metadata
        nil, // Value (not used in Sui)
    )
    if err != nil {
        return err
    }
    
    return nil
}
```

## Programmable Transaction Blocks (PTBs)

PTBs allow you to pipeline multiple Sui instructions (move calls, coin transfers, etc.) atomically. In this section we'll look at how to configure and use them.

#### Configuring a ChainWriter to use PTBs

To configure a ChainWriter to use PTBs, you need to specify the `PTBChainWriterModuleName` as the module name. ChainWriter will treat all the functions in this module as PTBs. The example below shows how to configure a ChainWriter instance that has a `simple_operation` function that receives an argument called `counter`.

```go
import (
    "encoding/json"
    "github.com/smartcontractkit/chainlink-common/pkg/types"
    "github.com/smartcontractkit/chainlink-sui/relayer/chainwriter"
    "github.com/smartcontractkit/chainlink-sui/relayer/codec"
)

func createChainWriterWithPTB(relayer *plugin.SuiRelayer) (types.ContractWriter, error) {
    // Define the ChainWriter configuration with PTB support
    config := chainwriter.ChainWriterConfig{
        Modules: map[string]*chainwriter.ChainWriterModule{
            chainwriter.PTBChainWriterModuleName: {
                Name:     chainwriter.PTBChainWriterModuleName,
                ModuleID: "0x123", // Package ID, not used for PTBs
                Functions: map[string]*chainwriter.ChainWriterFunction{
                    "simple_operation": {
                        Name:      "simple_operation",
                        PublicKey: []byte{/* Public key bytes */},
                        Params:    []codec.SuiFunctionParam{},
                        PTBCommands: []chainwriter.ChainWriterPTBCommand{
                            {
                                Type:      codec.SuiPTBCommandMoveCall,
                                PackageId: strPtr("0x123..."),
                                ModuleId:  strPtr("counter"),
                                Function:  strPtr("increment"),
                                Params: []codec.SuiFunctionParam{
                                    {
                                        Name:     "counter",
                                        Type:     "object_id",
                                        Required: true,
                                    },
                                },
                                Order: 1,
                            },
                        },
                    },
                },
            },
        },
    }
    
    // Serialize the configuration
    configBytes, err := json.Marshal(config)
    if err != nil {
        return nil, err
    }
    
    // Create the ChainWriter
    writer, err := relayer.NewContractWriter(context.Background(), configBytes)
    if err != nil {
        return nil, err
    }
    
    return writer, nil
}
```

#### Using PTB Dependencies

PTBs allow you to chain multiple commands together and to use the outputs of one command as arguments for another. The example below shows how to configure a ChainWriter instance that has a `complex_operation` function that receives an argument called `counter` and a `previous_result` argument that is the result of the previous command. the previous result is specified using the `PTBDependency` field and needs to have two fields: `CommandIndex` (the index of the command in the PTB) and `ResultIndex` (the index of the result in the command).

```go
config := chainwriter.ChainWriterConfig{
    Modules: map[string]*chainwriter.ChainWriterModule{
        "counter": {
            Name:     "counter",
            ModuleID: "0x123...",
            Functions: map[string]*chainwriter.ChainWriterFunction{
                "complex_operation": {
                    Name:      "complex_operation",
                    PublicKey: []byte{/* Public key bytes */},
                    PTBCommands: []chainwriter.ChainWriterPTBCommand{
                        {
                            Type:      codec.SuiPTBCommandMoveCall,
                            PackageId: "0x123...",
                            ModuleId:  "0xABC",
                            Function:  "function_1",
                            Params: []codec.SuiFunctionParam{
                                {
                                    Name:     "counter",
                                    Type:     "object_id",
                                    Required: true,
                                },
                            },
                            Order: 1,
                        },
                        {
                            Type:      codec.SuiPTBCommandMoveCall,
                            PackageId: "0x456...",
                            ModuleId:  "0xDEF",
                            Function:  "function_2",
                            Params: []codec.SuiFunctionParam{
                                {
                                    Name:     "counter",
                                    Type:     "object_id",
                                    Required: true,
                                },
                                {
                                    Name:     "previous_result",
                                    Type:     "ptb_dependency",
                                    Required: true,
                                    PTBDependency: &codec.PTBCommandDependency{
                                        CommandIndex: 0,
                                        ResultIndex:  0,
                                    },
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

### Calling a PTB in a ChainWriter instance

When working with PTBs, you can pass different types of arguments to your commands. The Sui relayer supports the following argument types:

1. **Object Arguments** (`ObjectArgType`) - Used for Sui object references:
   ```go
   {
       Type: chainwriter.ObjectArgType,
       Content: chainwriter.PTBArgContent{
           ID: "0x123...", // Object ID
           MapTo: []chainwriter.PTBArgLocation{
               {
                   CommandIndex: 0, // Index of the command in the PTB
                   Param: "counter", // Parameter name in the command
                   CommandName: "0x123.counter.increment", // Full command name (optional)
               },
           },
       },
   }
   ```

2. **Scalar Arguments** (`ScalarArgType`) - Used for basic value types like integers, strings, booleans, etc.:
   ```go
   {
       Type: chainwriter.ScalarArgType,
       Content: chainwriter.PTBArgContent{
           Value: uint64(10), // Any value type
           MapTo: []chainwriter.PTBArgLocation{
               {
                   CommandIndex: 1,
                   Param: "increment_by",
                   CommandName: "0x123.counter.increment_by",
               },
           },
       },
   }
   ```

3. **Vector Arguments** (`VectorArgType`) - Used for arrays/vectors of values:
   ```go
   {
       Type: chainwriter.VectorArgType,
       Content: chainwriter.PTBArgContent{
           Value: []string{"value1", "value2"}, // Array of values
           MapTo: []chainwriter.PTBArgLocation{
               {
                   CommandIndex: 0,
                   Param: "string_list",
                   CommandName: "0x123.list.process",
               },
           },
       },
   }
   ```

You can combine these different argument types in a single PTB. Each argument can be mapped to multiple commands within the same PTB using the `MapTo` field. Commands that depend on the result of previous commands are handled internally through the PTB dependency mechanism.

Here's a complete example that combines multiple argument types:

```go
ptbArgs := chainwriter.PTBArgMapping{
    Args: []chainwriter.PTBArg{
        {
            Type: chainwriter.ObjectArgType,
            Content: chainwriter.PTBArgContent{
                ID: "0x123...",
                MapTo: []chainwriter.PTBArgLocation{
                    {
                        CommandIndex: 0,
                        Param:        "counter",
                        CommandName:  "0x123.counter.increment",
                    },
                },
            },
        },
        {
            Type: chainwriter.ScalarArgType,
            Content: chainwriter.PTBArgContent{
                Value: uint64(10),
                MapTo: []chainwriter.PTBArgLocation{
                    {
                        CommandIndex: 1,
                        Param:        "increment_by",
                        CommandName:  "0x123.counter.increment_by",
                    },
                },
            },
        },
        {
            Type: chainwriter.ScalarArgType,
            Content: chainwriter.PTBArgContent{
                Value: "metadata",
                MapTo: []chainwriter.PTBArgLocation{
                    {
                        CommandIndex: 1,
                        Param:        "description",
                        CommandName:  "0x123.counter.increment_by",
                    },
                },
            },
        },
    },
}
```

When submitting a PTB transaction, use the `SubmitTransaction` method with your arguments:

```go
err := writer.SubmitTransaction(
    context.Background(),
    "counter",           // Module name
    "complex_operation", // Function name
    ptbArgs,            // PTB arguments
    txID,               // Transaction ID
    "",                 // To address (not used in Sui)
    &commonTypes.TxMeta{GasLimit: 10000000}, // Transaction metadata
    nil,                // Value (not used in Sui)
)
```

## Closing Resources

Don't forget to properly close resources when you're done:

```go
func cleanup(relayer *plugin.SuiRelayer, reader types.ContractReader, writer types.ContractWriter) {
    if writer != nil {
        writer.Close()
    }
    
    if reader != nil {
        reader.Close()
    }
    
    if relayer != nil {
        relayer.Close()
    }
}
```

The Sui relayer plugin provides a flexible and powerful way to interact with the Sui blockchain, allowing you to both read from and write to Sui smart contracts in your Chainlink integrations. 