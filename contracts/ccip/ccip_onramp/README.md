# Chainlink CCIP OnRamp

The CCIP OnRamp enables secure cross-chain message and token transfers from Sui to other supported blockchain networks. It serves as the entry point for all outbound cross-chain operations on Sui.

## Overview

The OnRamp handles the complete lifecycle of cross-chain operations, including message construction, fee calculation, token transfers, and destination chain routing. It integrates with Sui's Programmable Transaction Blocks (PTBs) to provide atomic, multi-step cross-chain operations.

**Important**: Many OnRamp functions, particularly `ccip_send`, are designed to be used exclusively within Programmable Transaction Blocks (PTBs) to ensure atomic execution and proper state management.

## Key Features

- **Cross-Chain Messaging**: Send messages and tokens to other blockchains
- **PTB Integration**: Designed for use within Programmable Transaction Blocks
- **Fee Management**: Automatic fee calculation and collection with multiple token support
- **Token Transfers**: Support for multiple token types in single messages
- **Access Control**: Configurable allowlists and owner controls
- **Security**: RMN integration and message integrity validation
- **Message Ordering**: Sequence numbers and nonce management for reliable delivery
- **Multi-Chain Support**: Configurable support for multiple destination chains
- **Upgrade Registry**: Integration with CCIP upgrade management system

## Core Functions

### Send Messages
**⚠️ PTB Required**: This function must be used within a Programmable Transaction Block.

```move
// Send a cross-chain message with optional token transfers
let message_id = onramp::ccip_send<T>(
    &mut ccip_ref,
    &mut onramp_state,
    &clock,
    dest_chain_selector,
    receiver_address,
    message_data,
    token_params,
    &fee_token_metadata,
    &mut fee_token_coin,
    extra_args,
    &mut ctx
);
```

**Process Flow:**
1. Validates destination chain and sender authorization
2. Calculates required fees based on message content and token transfers
3. Collects fee payment from sender
4. Increments sequence number and manages nonce
5. Constructs and emits the cross-chain message
6. Returns unique message ID for tracking

**Parameters:**
- `ccip_ref`: Reference to CCIP state object
- `onramp_state`: OnRamp state object
- `clock`: Sui clock for timestamp validation
- `dest_chain_selector`: Destination chain selector
- `receiver_address`: Destination address on target chain
- `message_data`: Message payload (arbitrary bytes)
- `token_params`: Token transfer parameters (from onramp_state_helper)
- `fee_token_metadata`: Metadata for fee payment token
- `fee_token_coin`: Fee payment token (will be consumed)
- `extra_args`: Additional execution parameters
- `ctx`: Transaction context

### Calculate Fees
**Safe for standalone calls**: This function can be called outside of PTBs for fee estimation.

```move
// Get required fee before sending
let required_fee = onramp::get_fee<T>(
    &ccip_ref,
    &clock,
    dest_chain_selector,
    receiver_address,
    message_data,
    token_addresses,
    token_amounts,
    &fee_token_metadata,
    extra_args
);
```

**Returns**: Required fee amount in the specified token denomination.

**Use Cases:**
- Pre-transaction fee estimation
- UI fee display
- Validation before sending messages

### Configuration Management
**Owner-only functions**: These require owner capabilities and should be used in PTBs for atomic updates.

```move
// Update destination chain configuration (owner only)
onramp::apply_dest_chain_config_updates(
    &mut onramp_state,
    &owner_cap,
    dest_chain_selectors,
    dest_chain_enabled,
    dest_chain_allowlist_enabled
);

// Update dynamic configuration (owner only)
onramp::set_dynamic_config(
    &mut onramp_state,
    &owner_cap,
    fee_aggregator,
    allowlist_admin
);

// Withdraw collected fees (owner only)
onramp::withdraw_fee_tokens<T>(
    &mut onramp_state,
    &owner_cap,
    &mut fee_tokens,
    amount
);
```

### Allowlist Management
**Allowlist admin functions**: Manage sender allowlists for specific destination chains.

```move
// Update allowlists (allowlist admin only)
onramp::apply_allowlist_updates(
    &mut onramp_state,
    &allowlist_admin_cap,
    dest_chain_selectors,
    dest_chain_allowlist_enabled,
    dest_chain_add_allowed_senders,
    dest_chain_remove_allowed_senders
);
```

### Query Functions
**Safe for standalone calls**: These functions can be called outside of PTBs for information retrieval.

```move
// Check if a chain is supported
let is_supported = onramp::is_chain_supported(&onramp_state, dest_chain_selector);

// Get expected next sequence number
let next_seq = onramp::get_expected_next_sequence_number(&onramp_state, dest_chain_selector);

// Get destination chain configuration
let config = onramp::get_dest_chain_config(&onramp_state, dest_chain_selector);

// Get allowed senders list
let allowed_senders = onramp::get_allowed_senders_list(&onramp_state, dest_chain_selector);

// Get outbound nonce for a sender
let nonce = onramp::get_outbound_nonce(&ccip_ref, dest_chain_selector, sender_address);
```

## State Structure

### OnRampState
The main state object that manages all onramp operations and configurations.

- `id`: Unique object identifier
- `package_ids`: Vector of package IDs for upgrade tracking
- `chain_selector`: Sui chain identifier (unique per network)
- `fee_aggregator`: Address where collected fees are sent
- `allowlist_admin`: Administrator for managing sender allowlists
- `dest_chain_configs`: Table mapping chain selectors to configurations
- `fee_tokens`: Bag storage for collected fee tokens
- `nonce_manager_cap`: Capability for managing message nonces
- `source_transfer_cap`: Capability for token transfer operations
- `ownable_state`: Ownership management state

### DestChainConfig
Configuration for each supported destination chain.

- `is_enabled`: Whether chain accepts new messages (can be disabled for maintenance)
- `sequence_number`: Message ordering counter (increments with each message)
- `allowlist_enabled`: Whether sender allowlist is enforced for this chain
- `allowed_senders`: Vector of authorized sender addresses (when allowlist enabled)

### Message Structures

#### RampMessageHeader
- `message_id`: Cryptographically secure unique message identifier
- `source_chain_selector`: Sui chain selector
- `dest_chain_selector`: Destination chain selector
- `sequence_number`: Message ordering number
- `nonce`: Sender-specific nonce for replay protection

#### Sui2AnyRampMessage
Complete cross-chain message structure.

- `header`: Message metadata and routing information
- `sender`: Sui sender address
- `data`: Arbitrary message payload
- `receiver`: Destination address on target chain
- `extra_args`: Additional execution parameters
- `fee_token`: Token used for fee payment
- `fee_token_amount`: Amount of fee tokens paid
- `fee_value_juels`: Fee value in juels (for compatibility)
- `token_amounts`: Vector of token transfers included in message

#### Sui2AnyTokenTransfer
Token transfer information within a message.

- `source_pool_address`: Source token pool address
- `dest_token_address`: Token address on destination chain
- `extra_data`: Additional token-specific data
- `amount`: Token amount to transfer
- `dest_exec_data`: Destination execution data

## Security Features

### Access Control
The onramp implements multiple layers of access control:

1. **Owner Controls**: Configuration changes, fee withdrawal, chain enable/disable
2. **Allowlist Admin**: Sender allowlist management for specific chains
3. **Capability-Based**: Token transfers and nonce management through capabilities

### Validation Mechanisms
- **Sender Authorization**: Validates sender against allowlists when enabled
- **Destination Chain**: Ensures target chain is supported and enabled
- **Fee Validation**: Verifies sufficient fee payment before processing
- **RMN Integration**: Checks for cursed chains via Risk Management Network
- **Token Validation**: Validates token transfers through onramp_state_helper
- **Message Integrity**: Unique message IDs and sequence numbers prevent replay attacks

### Upgrade Registry Integration
The onramp integrates with the CCIP Upgrade Registry to:
- Block specific function versions if vulnerabilities are discovered
- Prevent execution of problematic code paths
- Provide emergency stop mechanisms for critical functions

## Integration with CCIP Core

The onramp integrates with several CCIP core components:

### OnRamp State Helper
- **Token Parameters**: Constructs token transfer parameters for outgoing messages
- **Validation**: Ensures only authorized token pools can add transfers
- **Type Safety**: Uses proof-based validation for token operations

### Fee Quoter
- **Fee Calculation**: Determines gas costs and token fees for cross-chain operations
- **Cross-Chain Rates**: Handles exchange rates between different chains
- **Execution Costs**: Calculates destination chain execution costs

### Nonce Manager
- **Sequence Control**: Manages message ordering and nonce assignment
- **Replay Protection**: Prevents duplicate message execution
- **Out-of-Order Support**: Handles messages that don't require strict ordering

### RMN Remote
- **Risk Management**: Checks for cursed chains and halted operations
- **Security Verification**: Validates against known security issues

## Events

The onramp emits several events for monitoring and indexing:

- **`CCIPMessageSent`**: Emitted when a message is successfully sent
  - Contains message ID, source/dest chain selectors, sequence number, nonce
- **`ConfigSet`**: Configuration updates (static and dynamic configs)
- **`DestChainConfigSet`**: Destination chain configuration changes
- **`AllowlistSendersAdded`**: Senders added to allowlist for specific chains
- **`AllowlistSendersRemoved`**: Senders removed from allowlist
- **`FeeTokenWithdrawn`**: Fee token withdrawal events
- **`OwnershipTransferred`**: Ownership transfer events

## MCMS Integration

The onramp supports Multi-Chain Multi-Signature (MCMS) governance through several functions:

```move
// MCMS-specific functions for governance
onramp::mcms_initialize(...);
onramp::mcms_set_dynamic_config(...);
onramp::mcms_apply_dest_chain_config_updates(...);
onramp::mcms_apply_allowlist_updates(...);
onramp::mcms_transfer_ownership(...);
onramp::mcms_withdraw_fee_tokens<T>(...);
```

## Dependencies

- **ChainlinkCCIP**: Core CCIP functionality
- **ChainlinkManyChainMultisig**: Governance support

## Package Information

- **Name**: ChainlinkCCIPOnramp
- **Version**: 1.6.0
- **Edition**: 2024

## Usage

### Initialization
1. **Deploy Package**: Publish the onramp package to Sui
2. **Initialize State**: Set up onramp with chain selector and fee aggregator
3. **Configure Chains**: Add supported destination chains and their configurations
4. **Set Allowlists**: Configure sender allowlists for specific chains (optional)

### Sending Messages
1. **Estimate Fees**: Use `get_fee` to calculate required fees
2. **Prepare Token Transfers**: Use `onramp_state_helper` to create token parameters
3. **Send in PTB**: Use `ccip_send` within a Programmable Transaction Block
4. **Monitor Events**: Track `CCIPMessageSent` events for delivery confirmation

### Management Operations
1. **Configuration Updates**: Use owner functions to update chain configs
2. **Allowlist Management**: Use allowlist admin functions to manage senders
3. **Fee Withdrawal**: Use owner functions to withdraw collected fees
4. **Emergency Controls**: Use RMN integration for emergency stops

### Best Practices
- Always use `ccip_send` within PTBs for atomic execution
- Validate fees before sending to avoid transaction failures
- Monitor events for message delivery status
- Use allowlists for production deployments
- Keep configurations up-to-date with network changes
- Use query functions for UI and monitoring purposes

For detailed implementation examples and advanced configuration, see the source code and integration tests.