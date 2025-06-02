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
      - [Handling Generic Types in PTBs](#handling-generic-types-in-ptbs)
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

    args := chainwriter.Arguments{
        Args: map[string]any{
            "counter": counterObjectId,
        },
    }
    
    // Submit the transaction
    err := writer.SubmitTransaction(
        context.Background(),
        "counter", // Module name
        "increment", // Function name
        args, // Arguments
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

                        },
                    },
                },
            },
        },
    },
}
```

### Calling a PTB in a ChainWriter instance


You can combine these different argument types in a single PTB. The recommended way to pass arguments to PTB commands is using the `chainwriter.Arguments` struct. This approach simplifies argument handling by automatically mapping the arguments to the appropriate parameters in your PTB commands based on their names.

The `Arguments` struct has two fields:
- `Args`: A map of argument names to their values
- `ArgTypes`: A map of argument names to their generic types (for handling generic type parameters). By default this is empty.

Here's how to use the `chainwriter.Arguments` approach:

```go
// Create an Arguments struct with your argument values
args := chainwriter.Arguments{
    Args: map[string]any{
        "counter":      "0x123...", // Object ID for a counter object
        "increment_by": uint64(10), // Value to increment by
        "description":  "metadata", // String metadata
        "signers":      [][]byte{signerA, signerB, signerC}, // Vector/array argument
    },
}

// Submit the transaction with the Arguments struct
err := writer.SubmitTransaction(
    context.Background(),
    chainwriter.PTBChainWriterModuleName, // Module name - for PTBs always use this constant
    "complex_operation",                  // Function name
    args,                                 // Arguments struct
    txID,                                 // Transaction ID
    "",                                   // To address (not used in Sui)
    &commonTypes.TxMeta{GasLimit: 10000000}, // Transaction metadata
    nil,                                  // Value (not used in Sui)
)
```

This approach has several advantages:
1. **Simplicity**: No need to manually map arguments to command indices
2. **Automatic mapping**: Arguments are automatically mapped to parameters based on name
3. **Type safety**: The arguments are validated against the expected parameter types
4. **Dependency handling**: Dependencies between commands are handled internally

#### Handling Generic Types in PTBs

When working with generic types in Sui Move functions, you need to specify both the value arguments and the type arguments. The `chainwriter.Arguments` struct supports this through the `ArgTypes` field.

Here's a real-world example from a CCIP (Cross-Chain Interoperability Protocol) implementation that uses generics to handle different token types:

```go
// ChainWriter configuration with multiple commands and generic types
func configureChainWriterForCCIP(addresses ContractAddresses, publicKeyBytes []byte) chainwriter.ChainWriterConfig {
    // Define generic type variables
    coinTypeTag := "0x2::coin::Coin"
    coinParamName := "c"
    feeTokenParamName := "fee_token"
    
    return chainwriter.ChainWriterConfig{
        Modules: map[string]*chainwriter.ChainWriterModule{
            chainwriter.PTBChainWriterModuleName: {
                Name:     chainwriter.PTBChainWriterModuleName,
                ModuleID: "0x123",
                Functions: map[string]*chainwriter.ChainWriterFunction{
                    "ccip_send": {
                        Name:      "ccip_send",
                        PublicKey: publicKeyBytes,
						PTBCommands: []chainwriter.ChainWriterPTBCommand{
							// First command: create token params
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: strPtr(addresses.CCIPPackageID),
								ModuleId:  strPtr("dynamic_dispatcher"),
								Function:  strPtr("create_token_params"),
								Params:    []codec.SuiFunctionParam{},
							},
							// Second command: lock tokens in the token pool
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: strPtr(addresses.LinkLockReleaseTokenPool),
								ModuleId:  strPtr("lock_release_token_pool"),
								Function:  strPtr("lock_or_burn"),
								Params: []codec.SuiFunctionParam{
									{
										Name:     "ref",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:      "clock",
										Type:      "object_id",
										Required:  true,
										IsMutable: testutils.BoolPointer(false),
									},
									{
										Name:     "state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "c",
										Type:     "object_id",
										Required: true,
										IsGeneric: true,
									},
									{
										Name:     "remote_chain_selector",
										Type:     "u64",
										Required: true,
									},
									{
										Name:     "token_params",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 0,
										},
									},
								},
							},
							{
								Type:      codec.SuiPTBCommandMoveCall,
								PackageId: strPtr(addresses.CCIPOnrampPackageID),
								ModuleId:  strPtr("onramp"),
								Function:  strPtr("ccip_send"),
								GenericTypeArgs: []codec.GenericArg{
									{
										TypeTag:   &coinTypeTag,
										ParamName: &feeTokenParamName,
									},
								},
								Params: []codec.SuiFunctionParam{
									{
										Name:     "ref",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:     "onramp_state",
										Type:     "object_id",
										Required: true,
									},
									{
										Name:      "clock",
										Type:      "object_id",
										Required:  true,
										IsMutable: testutils.BoolPointer(false),
									},
									{
										Name:     "dest_chain_selector",
										Type:     "u64",
										Required: true,
									},
									{
										Name:     "receiver",
										Type:     "vector<u8>",
										Required: true,
									},
									{
										Name:     "data",
										Type:     "vector<u8>",
										Required: true,
									},
									{
										Name:     "token_params",
										Type:     "ptb_dependency",
										Required: true,
										PTBDependency: &codec.PTBCommandDependency{
											CommandIndex: 1,
										},
									},
									{
										Name:      "fee_token_metadata",
										Type:      "object_id",
										Required:  true,
										IsMutable: testutils.BoolPointer(false),
									},
									{
										Name:     "fee_token",
										Type:     "object_id",
										Required: true,
										IsGeneric: true,
									},
									{
										Name:     "extra_args",
										Type:     "vector<u8>",
										Required: true,
									},
								},
							},
						},
                },
            },
        },
    }
}

// Create Arguments with generic type information for the CCIP send operation
func createCCIPSendArguments(addresses ContractAddresses) chainwriter.Arguments {
	// Define a destination chain selector (e.g., Ethereum Sepolia)
	destChainSelector := uint64(2)
	linkTokenTypeTag := "0xe3c005c4195ec60a3468ce01238df650e4fedbd36e517bf75b9d2ee90cce8a8b::link_token::LINK_TOKEN"

	return chainwriter.Arguments{
		Args: map[string]any{
			"ref":                   addresses.CCIPStateRef,
			"clock":                 addresses.ClockObject,
			"remote_chain_selector": destChainSelector,
			"dest_chain_selector":   destChainSelector,
			"state":                 addresses.LinkLockReleaseTokenPoolState,
			"c":                     addresses.LinkCoinObjects[0],
			"onramp_state":          addresses.CCIPOnrampState,
			"receiver":              []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"data":                  []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			"fee_token_metadata":    addresses.LinkTokenCoinMetadata,
			"fee_token":             addresses.LinkCoinObjects[1],
			"extra_args":            []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		ArgTypes: map[string]string{
			"c":         linkTokenTypeTag,
			"fee_token": linkTokenTypeTag,
		},
	}
}
```

When the PTB is executed:

1. The first command creates token parameters with no generic types
2. The `c` argument is passed as a generic type to the second command. That command will receive the `c` argument as a `Coin<T>` object.
3. In the third command, the `fee_token` argument is passed as a generic type. That command will receive the `fee_token` argument as a `Coin<T>` object.

This allows the same PTB configuration to work with different token types by simply changing the values in the `ArgTypes` map, making your code more flexible and reusable. Because the generic types are the same, only one argtype will be passed to the PTB execution engine.

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