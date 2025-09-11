# Chainlink CCIP (Cross-Chain Interoperability Protocol) for Sui

This package provides the core implementation of Chainlink's Cross-Chain Interoperability Protocol (CCIP) for the Sui blockchain. CCIP enables secure cross-chain communication and token transfers between different blockchain networks.

## Table of Contents

- [Overview](#overview)
- [Core Components](#core-components)
  - [Token Admin Registry](#token-admin-registry-token_admin_registrymove)
  - [Receiver Registry](#receiver-registry-receiver_registrymove)
  - [OnRamp State Helper](#onramp-state-helper-onramp_state_helpermove)
  - [Offramp State Helper](#offramp-state-helper-offramp_state_helpermove)
  - [RMN Remote](#rmn-remote-rmn_remotemove)
  - [Client](#client-clientmove)
  - [Upgrade Registry](#upgrade-registry-upgrade_registrymove)
- [Additional Utilities](#additional-utilities)
- [Using CCIP as a Dependency](#using-ccip-as-a-dependency)
- [Architecture Overview](#architecture-overview)
- [Security Considerations](#security-considerations)
- [Package Information](#package-information)
- [Recent Updates](#recent-updates)
- [Getting Help](#getting-help)
- [Contributing](#contributing)

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
- `register_pool()`: Register a new token pool
- `get_pools()`: Retrieve token pool addresses for given tokens
- `get_token_configs()`: Get detailed information about registered token pools
- `transfer_admin_role()`: Transfer administrative control of a token

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

### OnRamp State Helper (`onramp_state_helper.move`)

The OnRamp State Helper handles token transfer parameters for outgoing cross-chain messages. It provides a secure way to construct and validate token transfer requests for onramp operations.

**Key Features:**
- Create token transfer parameters for cross-chain operations
- Type-safe token pool validation through proof system
- Support for token transfers in outgoing messages
- Destination chain selector validation

**Key Functions:**
- `create_token_transfer_params()`: Create new token transfer parameters
- `add_token_transfer_param()`: Add a token transfer to the parameters
- `deconstruct_token_params()`: Extract token transfer data (permissioned)
- `has_token_transfer()`: Check if token transfer exists in parameters

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
- `get_dest_token_transfer_data()`: Retrieve token transfer information

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

### Upgrade Registry (`upgrade_registry.move`)

The Upgrade Registry provides controlled contract upgrade management and function restriction capabilities for CCIP modules. This is a critical security feature that allows administrators to block specific functions or entire versions of modules to prevent unauthorized or problematic upgrades.

**Key Features:**
- **Version Blocking**: Block entire versions of modules to prevent problematic upgrades
- **Function Blocking**: Block specific functions within modules for granular control
- **Access Control**: Owner-controlled upgrade restrictions with proper authorization
- **Event Emission**: Track upgrade restrictions through emitted events
- **Query Interface**: Check if functions or versions are allowed before execution

**Key Functions:**
- `initialize()`: Initialize the upgrade registry
- `block_version()`: Block an entire version of a module
- `block_function()`: Block a specific function within a module
- `get_module_restrictions()`: Retrieve all restrictions for a module
- `is_function_allowed()`: Check if a function is allowed to be called
- `verify_function_allowed()`: Assert that a function is allowed (with abort on failure)

**Usage Example:**
```move
// Block version 2 of the fee_quoter module
upgrade_registry::block_version(
    ccip_ref,
    owner_cap,
    b"fee_quoter".to_string(),
    2,
    ctx
);

// Block a specific function in version 1
upgrade_registry::block_function(
    ccip_ref,
    owner_cap,
    b"offramp".to_string(),
    b"execute_message".to_string(),
    1,
    ctx
);

// Check if a function is allowed before calling it
assert!(
    upgrade_registry::is_function_allowed(
        ccip_ref,
        b"fee_quoter".to_string(),
        b"get_fee".to_string(),
        1
    ),
    EFunctionNotAllowed
);
```

**Security Benefits:**
- Prevents execution of known vulnerable functions
- Allows gradual rollout of new versions
- Provides emergency stop mechanism for problematic upgrades
- Maintains audit trail through event emissions

## Additional Utilities

### Other Important Modules

- **Fee Quoter** (`fee_quoter.move`): Calculates fees for cross-chain operations
- **Nonce Manager** (`nonce_manager.move`): Manages message nonces for ordering
- **State Object** (`state_object.move`): Provides shared state management
- **Allowlist** (`allowlist.move`): Manages authorized addresses
- **Merkle Proof** (`merkle_proof.move`): Provides merkle tree verification utilities
- **ETH ABI** (`eth_abi.move`): Ethereum ABI encoding utilities
- **Upgrade Registry** (`upgrade_registry.move`): Manages contract upgrade restrictions and function blocking

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
use ccip::onramp_state_helper;
use ccip::offramp_state_helper;
use ccip::upgrade_registry;
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
let token_params = onramp_state_helper::create_token_transfer_params(token_receiver);

// Add token transfers
let token_params = onramp_state_helper::add_token_transfer_param<TOKEN_POOL_PROOF>(
    ccip_ref,
    token_params,
    destination_chain_selector,
    amount,
    source_token_address,
    dest_token_address,
    extra_data,
    TOKEN_POOL_PROOF {}
);

// Send through onramp (implementation depends on your onramp integration)
```

### 6. Using the Upgrade Registry

The Upgrade Registry allows you to control which functions and versions can be executed in your CCIP modules:

```move
module your_package::upgrade_management {
    use ccip::upgrade_registry;
    use ccip::state_object::CCIPObjectRef;
    use ccip::ownable::OwnerCap;

    // Block a problematic version
    public fun block_problematic_version(
        ccip_ref: &mut CCIPObjectRef,
        owner_cap: &OwnerCap,
        ctx: &mut TxContext
    ) {
        upgrade_registry::block_version(
            ccip_ref,
            owner_cap,
            b"fee_quoter".to_string(),
            2, // Block version 2
            ctx
        );
    }

    // Block a specific function
    public fun block_vulnerable_function(
        ccip_ref: &mut CCIPObjectRef,
        owner_cap: &OwnerCap,
        ctx: &mut TxContext
    ) {
        upgrade_registry::block_function(
            ccip_ref,
            owner_cap,
            b"offramp".to_string(),
            b"execute_message".to_string(),
            1, // Block version 1 of this function
            ctx
        );
    }

    // Check if a function is allowed before calling it
    public fun safe_function_call(
        ccip_ref: &CCIPObjectRef,
        module_name: vector<u8>,
        function_name: vector<u8>,
        version: u8
    ) {
        assert!(
            upgrade_registry::is_function_allowed(
                ccip_ref,
                module_name,
                function_name,
                version
            ),
            EFunctionNotAllowed
        );
        
        // Proceed with function call
    }
}
```

## Architecture Overview

The CCIP package follows a modular architecture where:

1. **Token Admin Registry** manages token pool configurations
2. **Receiver Registry** handles message receiver registration
3. **OnRamp State Helper** manages outgoing token transfers
4. **Offramp State Helper** processes incoming transfers and messages
5. **RMN Remote** provides security and risk management
6. **Client** defines core message structures
7. **Upgrade Registry** provides controlled upgrade management and function restrictions

This design ensures type safety, proper access control, and secure cross-chain operations while maintaining flexibility for different use cases. The Upgrade Registry adds an additional layer of security by allowing administrators to control which functions and versions can be executed, providing protection against vulnerable or unauthorized upgrades.

## Security Considerations

- All operations use type-safe proof systems to prevent unauthorized access
- Token transfers are validated against registered token pools
- Message receivers must be properly registered before receiving messages
- RMN signatures are verified before processing critical operations
- Administrative functions are protected by ownership controls
- **Upgrade Registry** provides additional security by allowing function and version blocking
- Controlled upgrade management prevents execution of vulnerable or unauthorized code
- Event emission provides audit trail for all upgrade restrictions

## Package Information

- **Package Name**: ChainlinkCCIP
- **Version**: 1.6.0
- **Edition**: 2024.beta
- **Dependencies**: ChainlinkManyChainMultisig (MCMS)

## Recent Updates

### Upgrade Registry Feature
The latest version includes the **Upgrade Registry** module, which provides:
- Controlled contract upgrade management
- Function and version blocking capabilities
- Enhanced security for CCIP operations
- Audit trail through event emissions

This feature is particularly important for production deployments where controlled upgrades and security are critical.

## Getting Help

For more information about specific functions and their usage, refer to:
- Individual module documentation in the source files
- The broader Chainlink CCIP protocol documentation
- [Chainlink Documentation](https://docs.chain.link/ccip)
- [Sui Documentation](https://docs.sui.io/)

## Contributing

Contributions to the CCIP package are welcome! Please ensure that:
- All new functions include proper documentation
- Security considerations are addressed
- Tests are included for new functionality
- The Upgrade Registry is considered for any new modules that may need upgrade control 