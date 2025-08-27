# Move & Sui Challenges for CCIP Off-Chain Relayer

> **Note**: Please start by reviewing the Object basics introduced in the On-Chain documentation.

### Challenge 1: Object Addresses
Given that Sui's network uses an object-centric model, all calls made to contracts require knowledge of the location (address) of the objects that a given method interacts with. This introduces a challenge that can be handled in several different ways.

First, let's consider a basic example of a counter contract. In a language such as Solidity, we only need to be aware of the method (e.g., `getCount()`) and we can simply call it without knowing anything about the underlying object it fetches or mutates.

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title A simple counter contract
contract Counter {
    // The current count
    uint256 private count;

    /// @param initialCount The starting value for the counter
    constructor(uint256 initialCount) {
        count = initialCount;
    }

    /// @notice Returns the current count
    /// @return The current value of the counter
    function getCount() public view returns (uint256) {
        return count;
    }
}
```


In Move, we must pass a reference to the object (its address - MoveVM finds the actual reference) in order to read or mutate an object.

```move
module 0x2::counter {
    use sui::object::{UID, Self, ID};
    use sui::object;
    use sui::signer;
    use sui::tx_context::{TxContext};

    /// A simple counter resource
    struct Counter has key, store {
        id: UID,
        value: u64,
    }

    /// Create a new Counter owned by `owner` with an initial value
    public(entry) fun create(
        owner: &signer,
        initial: u64,
        ctx: &mut TxContext
    ) {
        let id = UID::new(owner);
        let counter = Counter { id, value: initial };
        object::move_to(owner, counter);
        emit<Event<CounterEvent>>(CounterEvent { new_value: initial }, ctx);
    }

    /// Read the current counter value
    public(view) fun get_count(
        counter_ref: &Counter
    ): u64 {
        counter_ref.value
    }
}
```


With this, we are able to tell the contract which object we are referring to. Object addresses are returned during the publish step and must be kept track of immediately. Note that a contract cannot see objects it does not own or objects that have not been shared with the contract.

> **Note**: The Sui RPC endpoints include methods for fetching an object's details given its ID as well as finding all the objects that a certain address (EOA or contract) owns or has access to. The restriction of needing to own an object to read it does not exist off-chain, only within contracts.

#### Solution: Object Pointers

Since there are many objects and states that we may need to keep track of, we use what we call "object pointers" - single objects (single address) that contain addresses to all other objects we are interested in. This way, we can make a single call (usually per contract) to get the references we need.

This is important because when the off-chain relayer starts, we are only aware of the off-ramp package address and nothing else. We use the object pointer in the off-ramp contract to find the locations of other packages as well as object references (e.g., state object).

In ChainReader, we store a reference to the field within a pointer object that we may need to fetch before making calls to contracts.

> **TODO**: Add link - See code example for details.

### Challenge 2: Lack of Dynamic Dispatch
There are two main ways to orchestrate multi-step logic in smart contracts: dynamic dispatch (as in EVM) and Programmable Transaction Blocks (as in Sui).

The main differences are outlined in the table below:


| Aspect | Programmable Transaction Block | Dynamic Dispatch |
|--------|-------------------------------|------------------|
| **Core idea** | Pack an ordered list of calls into one atomic transaction. The list is explicit at submission time. | Call a function by name or type at run-time; the system resolves which concrete implementation to run. |
| **Where you see it** | Sui, Aptos, Solana & other object-centric chains. | Traditional OO languages (Java/C++ virtual methods), EVM delegatecall, plugin architectures, many runtime-polymorphic APIs. |
| **When it links** | At build/client time: the PTB lists every module/function ID up front. It must also have the values of every object that will be referenced in arguments to those method calls. | At execution time: a dispatcher, v-table, or address lookup chooses the target. |
| **Atomicity** | All steps succeed or the whole PTB reverts—perfect for complex state changes. | Each dispatched call is just another function call; atomicity depends on surrounding transaction logic. |
| **Security surface** | Smaller: every callee is known and type-checked before the tx hits the chain. | Larger: the resolver must check permissions at runtime; wrong target can be exploitable. |
| **Flexibility** | Caller can chain any sequence of modules—even ones that don't know each other—without those modules supporting dispatch. | Callee evolution: you can add new implementations without changing callers. Excellent for plugin systems. |

Due to the lack of dynamic dispatch, the off-chain relayer plays an important role in orchestrating multi-step calls to contracts in Sui. The relayer is responsible for constructing the PTB, which comprises multiple commands (i.e., sub-transactions) that are conceptually similar to database transactions.

One of the challenges with this approach is that all object IDs (references) must be known during PTB construction. In practice, we don't keep track of these object IDs in a centralized database and often need to make calls to the contracts' object pointers (mentioned above) to get all the references needed.

#### Solution: PTB Constructor

In the implementation of the Sui relayer, we created an abstraction called `PTBConstructor` which is responsible for reading a specific configuration from the ChainWriter configuration and building a PTB from it.

Below is an example of what a ChainWriter configuration would look like for a 3-command PTB. In practice, it gets much more complex than this.

```go
config := chainwriter.ChainWriterConfig{
    Modules: map[string]*chainwriter.ChainWriterModule{
        "counter": {
            Name:     "counter",
            ModuleID: packageId,
            Functions: map[string]*chainwriter.ChainWriterFunction{
                "ptb_example_operation": {
                    Name:      "ptb_example_operation",
                    PublicKey: publicKeyBytes,
                    PTBCommands: []chainwriter.ChainWriterPTBCommand{
                        {
                            Type:      codec.SuiPTBCommandMoveCall,
                            PackageId: &packageId,
                            ModuleId:  stringPointer("counter"),
                            Function:  stringPointer("increment"),
                            Params: []codec.SuiFunctionParam{
                                {
                                    Name:     "counter_id",
                                    Type:     "object_id",
                                    Required: true,
                                },
                            },
                        },
                        {
                            Type:      codec.SuiPTBCommandMoveCall,
                            PackageId: &packageId,
                            ModuleId:  stringPointer("counter"),
                            Function:  stringPointer("increment_by"),
                            Params: []codec.SuiFunctionParam{
                                {
                                    Name:     "counter_id",
                                    Type:     "object_id",
                                    Required: true,
                                },
                                {
                                    Name:         "increment_by",
                                    Type:         "u64",
                                    Required:     true,
                                    DefaultValue: uint64(10),
                                },
                            },
                        },
                        {
                            Type:      codec.SuiPTBCommandMoveCall,
                            PackageId: &packageId,
                            ModuleId:  stringPointer("counter"),
                            Function:  stringPointer("get_count"),
                            Params: []codec.SuiFunctionParam{
                                {
                                    Name:     "counter_id",
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
```


Notice that each PTB has a name which can be referenced during the construction step along with the arguments that are necessary for each command. Constructing the PTB can then be done as follows:

```go
// Construct the PTB
ptb, cError := constructor.BuildPTBCommands(ctx, "counter", "ptb_example_operation", args, nil)
require.NoError(t, cError)

// Get the results
ptbResult, err := ptbClient.FinishPTBAndSend(ctx, &txnSigner, ptb, client.WaitForLocalExecution)
```


> **TODO**: Add code references and diagram

### Challenge 3: BCS Encoding & Decoding 

When you call into Sui Move contracts from Go, every argument and return value must be serialized and deserialized with exactly the same Binary Canonical Serialization (BCS) rules the VM uses. A few of the sharp edges you'll encounter are highlighted below.

> **Note**: When decoding values in BCS, the exact expected type must be known. For example, if we are deserializing the BCS bytes of Counter from our example above, we must know exactly what the struct looks like when deserializing. Alternatively, we can deserialize values one by one manually without any mismatched types in the order of deserialization.

#### Common BCS Challenges

**Minimal LEB128 lengths**

BCS uses ULEB128 to encode all integer and length fields, and rejects any encoding that isn't the shortest possible form. If your Go code uses `binary.PutUvarint` (or hand-rolled loops) to emit a varint, it's easy to accidentally emit extra continuation bytes. The result? "Deserialization error: ULEB128 encoding was not minimal in size" when the node tries to read your payload.

**Pure vs. object arguments**

For `vector<u8>`, `u64`, `bool`, or `struct<T>`, you must pass them as "pure" BCS blobs. For any on-chain object (a Move resource, shared or owned), you hand over a three-field object reference. Mixing them up – or forgetting to mark a `vector<address>` as pure – will either get your tx rejected or crash your Go SDK with a reflect-based panic.

**Generics & type-tags**

Every generic Move function needs concrete TypeTags (e.g., `0x2::sui::SUI`) in the same order the module declares them. Your Go code must build and BCS-encode those tags before the pure arguments. A mismatch in order or missing a nested `vector<T>` tag will trigger a "type argument mismatch" or corrupt the whole payload.

**Duplication and caching**

Many Go SDK builders dedupe identical pure values via an internal index map, panicking if you call `ptb.Pure(blob)` twice with exactly the same bytes. To work around this, you either deep-copy every slice (so pointers differ) or cache and reuse the same Argument handle for repeated literals.

**Big-integer sizes**

Move's `u128` or higher-precision types must be marshaled as little-endian 16-byte fields, often via `math/big`. Off-by-one in byte order or forgetting to zero-pad quickly leads to "invalid BCS bytes" errors that can be tough to debug.



While these issues are all resolved in the TypeScript and Rust SDKs built by the Mysten Labs team, there is no officially supported Go SDK. We chose to go with the BlockVision community SDK and add the encoding/decoding features to it.

#### Solution: Codec Package

All the issues above are now automatically handled in the relayer via the PTBClient and PTBConstructor abstractions but would otherwise need to be done for a Go integration that does not import the relayer.

The approach taken in the relayer is to create a codec package that handles type manipulation, parsing and encoding data to/from BCS. We also leveraged the BCS deserializer that was built for Aptos due to its user-friendly interface.

Another important note is that the Sui RPC offers endpoints to get the normalized types within modules. This is tremendously helpful because responses from contract methods return BCS encoded data and must be deserialized. The normalized module value includes the equivalent of an ABI in EVM chains which can be used to deserialize any response coming back from the contracts as well as figure out complex generic types that contract methods expect to be called with.

> **TODO**: Add code snippets from the decoder

### Challenge 4: Dynamic Number Token Pools

> **TODO**: Add the docs for the PTB expander here.





### Challenge 5: Interfacing over LOOP

> **TODO**: Add detailed contents





### Challenge 6: Events and Transactions Indexing

In this section, two separate but related challenges are discussed, both originating from the topic of Sui RPC node reliance and possible limitations.