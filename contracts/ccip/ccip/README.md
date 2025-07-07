# Chainlink CCIP (Cross-Chain Interoperability Protocol) for Sui

This package provides the core implementation of Chainlink's Cross-Chain Interoperability Protocol (CCIP) for the Sui blockchain. CCIP enables secure cross-chain communication and token transfers between different blockchain networks.

## Overview

The CCIP package serves as the foundational layer for cross-chain messaging and token transfers on Sui. It provides a comprehensive set of modules that handle message routing, token administration, receiver registration, and security features through Risk Management Network (RMN) integration.

## Core Components

### Token Admin Registry (`token_admin_registry.move`)

The Token Admin Registry manages token pool configurations and administrative control for cross-chain token transfers. It maintains a registry of token pools and their associated metadata.

**Key Features:**
- Register and manage token pools for cross-chain transfers
- Configure token administrators and handle ownership transfers
- Retrieve token pool information and validate token configurations
- Support for pending administrator transfers with two-step ownership model

**Key Functions:**
- `initialize()`: Initialize the token admin registry
- `register_token_pool()`: Register a new token pool
- `get_pools()`: Retrieve token pool addresses for given tokens
- `get_pool_infos()`: Get detailed information about registered token pools
- `transfer_administrator()`: Transfer administrative control of a token

### Receiver Registry (`receiver_registry.move`)

The Receiver Registry manages the registration of CCIP message receivers, enabling applications to receive and process cross-chain messages.

**Key Features:**
- Register receiver contracts to handle incoming CCIP messages
- Type-safe receiver validation through proof system
- Support for both stateless and stateful receivers
- Manage receiver configurations and module information

**Key Functions:**
- `initialize()`: Initialize the receiver registry
- `register_receiver()`: Register a new CCIP message receiver
- `unregister_receiver()`: Remove a receiver from the registry
- `get_receiver_config()`: Retrieve receiver configuration
- `get_receiver_module_and_state()`: Get receiver module name and state

### Dynamic Dispatcher (`dynamic_dispatcher.move`)

The Dynamic Dispatcher handles token transfer parameters for outgoing cross-chain messages. It provides a secure way to construct and validate token transfer requests.

**Key Features:**
- Create token transfer parameters for cross-chain operations
- Type-safe token pool validation through proof system
- Support for multiple token transfers in a single message
- Destination chain selector validation

**Key Functions:**
- `create_token_params()`: Create new token transfer parameters
- `add_source_token_transfer()`: Add a token transfer to the parameters
- `deconstruct_token_params()`: Extract token transfer data (permissioned)
- `get_source_token_transfer_data()`: Retrieve token transfer information

### Offramp State Helper (`offramp_state_helper.move`)

The Offramp State Helper manages incoming token transfers and message delivery on the destination chain. It handles the completion of cross-chain operations.

**Key Features:**
- Manage destination token transfers for received CCIP messages
- Handle message extraction and delivery to receivers
- Track completion status of token transfers
- Type-safe token pool validation

**Key Functions:**
- `create_receiver_params()`: Create parameters for incoming messages
- `add_dest_token_transfer()`: Add destination token transfer
- `complete_token_transfer()`: Mark token transfer as completed
- `extract_any2sui_message()`: Extract message for receiver processing
- `populate_message()`: Add message data to receiver parameters

### RMN Remote (`rmn_remote.move`)

The RMN (Risk Management Network) Remote provides signature verification and risk management capabilities for CCIP operations.

**Key Features:**
- Verify multi-signature reports from RMN nodes
- Manage RMN configuration and signer sets
- Handle cursing and uncursing of subjects
- Merkle root verification for cross-chain state validation

**Key Functions:**
- `initialize()`: Initialize RMN remote state
- `verify()`: Verify RMN signatures on reports
- `set_config()`: Update RMN configuration
- `curse()`: Curse subjects to halt operations
- `uncurse()`: Remove curse from subjects

### Client (`client.move`)

The Client module provides core data structures and utilities for CCIP messaging, including message formats and encoding functions.

**Key Features:**
- Define standard message structures for cross-chain communication
- Provide encoding utilities for extra arguments
- Support for both generic and SVM-specific message formats
- Token amount structures for cross-chain transfers

**Key Structures:**
- `Any2SuiMessage`: Standard cross-chain message format
- `Any2SuiTokenAmount`: Token amount specification
- Extra args encoding functions for different chain types

## Additional Utilities

### Other Important Modules

- **Fee Quoter** (`fee_quoter.move`): Calculates fees for cross-chain operations
- **Nonce Manager** (`nonce_manager.move`): Manages message nonces for ordering
- **State Object** (`state_object.move`): Provides shared state management
- **Allowlist** (`allowlist.move`): Manages authorized addresses
- **Merkle Proof** (`merkle_proof.move`): Provides merkle tree verification utilities
- **ETH ABI** (`eth_abi.move`): Ethereum ABI encoding utilities

## Using CCIP as a Dependency

To use the CCIP package as a dependency in your Move project, follow these steps:

### 1. Add Dependency to Move.toml

```toml
[dependencies]
ChainlinkCCIP = { local = "../path/to/ccip" }
# or if using a specific version/branch
# ChainlinkCCIP = { git = "https://github.com/smartcontractkit/chainlink-sui.git", subdir = "contracts/ccip/ccip", branch = "main" }
```

### 2. Import Required Modules

```move
use ccip::client;
use ccip::receiver_registry;
use ccip::token_admin_registry;
use ccip::dynamic_dispatcher;
use ccip::offramp_state_helper;
```

### 3. Implement a CCIP Receiver

```move
module your_package::ccip_receiver {
    use ccip::client::{Any2SuiMessage};
    use ccip::offramp_state_helper::{ReceiverParams};
    use ccip::receiver_registry;
    use ccip::state_object::CCIPObjectRef;

    public struct YourReceiver has key {
        id: UID,
        // your receiver state
    }

    public struct CCIP_RECEIVER_PROOF has drop {}

    public fun ccip_receive(
        _ref: &CCIPObjectRef,
        _receiver_state: &mut YourReceiver,
        _receiver_package_id: address,
        receiver_params: ReceiverParams
    ): ReceiverParams {
        // Process incoming CCIP message
        // Extract message if needed
        // Handle token transfers
        receiver_params
    }
}
```

### 4. Register Your Receiver

```move
// In your initialization function
receiver_registry::register_receiver<CCIP_RECEIVER_PROOF>(
    ccip_ref,
    receiver_state_id,
    CCIP_RECEIVER_PROOF {}
);
```

### 5. Send Cross-Chain Messages

```move
// Create token transfer parameters
let token_params = dynamic_dispatcher::create_token_params(destination_chain_selector);

// Add token transfers
let token_params = dynamic_dispatcher::add_source_token_transfer<TOKEN_POOL_PROOF>(
    ccip_ref,
    token_params,
    amount,
    source_token_address,
    dest_token_address,
    extra_data,
    TOKEN_POOL_PROOF {}
);

// Send through onramp (implementation depends on your onramp integration)
```

## Architecture Overview

The CCIP package follows a modular architecture where:

1. **Token Admin Registry** manages token pool configurations
2. **Receiver Registry** handles message receiver registration
3. **Dynamic Dispatcher** manages outgoing token transfers
4. **Offramp State Helper** processes incoming transfers and messages
5. **RMN Remote** provides security and risk management
6. **Client** defines core message structures

This design ensures type safety, proper access control, and secure cross-chain operations while maintaining flexibility for different use cases.

## Security Considerations

- All operations use type-safe proof systems to prevent unauthorized access
- Token transfers are validated against registered token pools
- Message receivers must be properly registered before receiving messages
- RMN signatures are verified before processing critical operations
- Administrative functions are protected by ownership controls

## Package Information

- **Package Name**: ChainlinkCCIP
- **Version**: 1.6.0
- **Edition**: 2024.beta
- **Dependencies**: ChainlinkManyChainMultisig (MCMS)

For more information about specific functions and their usage, refer to the individual module documentation and the broader Chainlink CCIP protocol documentation. 