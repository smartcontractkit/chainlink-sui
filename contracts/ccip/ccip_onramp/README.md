# Chainlink CCIP Onramp

The CCIP Onramp is a core component of Chainlink's Cross-Chain Interoperability Protocol (CCIP) that enables secure cross-chain message and token transfers from Sui to other supported blockchain networks.

## Overview

The onramp serves as the entry point for cross-chain operations on Sui, handling message construction, fee calculation, token transfers, and destination chain routing. It provides a secure, configurable interface for sending messages and tokens to other blockchains while maintaining proper access controls and validation.

## Key Components

### OnRampState

The main state object that manages all onramp operations and configurations.

**Core Fields:**
- `chain_selector`: Unique identifier for the Sui chain
- `fee_aggregator`: Address where collected fees are sent
- `allowlist_admin`: Administrator for managing sender allowlists
- `dest_chain_configs`: Configuration for each supported destination chain
- `fee_tokens`: Storage for collected fee tokens
- `nonce_manager_cap`: Capability for managing message nonces
- `source_transfer_cap`: Capability for token transfer operations

### Destination Chain Configuration

Each supported destination chain has its own configuration:

**DestChainConfig Fields:**
- `is_enabled`: Whether the destination chain accepts new messages
- `sequence_number`: Current sequence number for message ordering
- `allowlist_enabled`: Whether sender allowlist is enforced
- `allowed_senders`: List of authorized sender addresses (if allowlist enabled)

### Message Structure

The onramp constructs `Sui2AnyRampMessage` objects containing:

**Message Components:**
- `header`: Message metadata (ID, chain selectors, sequence number, nonce)
- `sender`: Address of the message sender on Sui
- `data`: Arbitrary message data
- `receiver`: Destination address on the target chain
- `fee_token`: Token used for fee payment
- `fee_token_amount`: Amount of fee tokens paid
- `token_amounts`: List of token transfers included in the message

## Core Functions

### Message Sending

#### `ccip_send<T>()`

The primary function for sending cross-chain messages with optional token transfers.

**Parameters:**
- `ref`: Reference to CCIP state object
- `state`: Onramp state object
- `clock`: Sui clock for timestamp validation
- `receiver`: Destination address on target chain
- `data`: Message payload
- `token_params`: Token transfer parameters (from dynamic dispatcher)
- `fee_token_metadata`: Metadata for fee payment token
- `fee_token`: Fee payment token
- `extra_args`: Additional execution parameters
- `ctx`: Transaction context

**Returns:**
- `vector<u8>`: Message ID for tracking

**Process:**
1. Validates destination chain and sender authorization
2. Calculates required fees based on message content
3. Collects fee payment from sender
4. Increments sequence number and manages nonce
5. Constructs and emits the cross-chain message
6. Returns unique message ID

### Fee Management

#### `get_fee<T>()`

Calculates the required fee for a cross-chain message.

**Parameters:**
- `ref`: CCIP state object reference
- `clock`: Sui clock
- `dest_chain_selector`: Target chain identifier
- `receiver`: Destination address
- `data`: Message payload
- `token_addresses`: List of token addresses for transfers
- `token_amounts`: List of token amounts
- `fee_token`: Fee token metadata
- `extra_args`: Additional parameters

**Returns:**
- `u64`: Required fee amount in the specified token

#### `withdraw_fee_tokens<T>()`

Allows the owner to withdraw collected fees to the configured fee aggregator.

**Access:** Owner only

### Configuration Management

#### `initialize()`

Initializes the onramp with essential configurations.

**Parameters:**
- `nonce_manager_cap`: Capability for nonce management
- `source_transfer_cap`: Capability for token transfers
- `chain_selector`: Sui chain identifier
- `fee_aggregator`: Fee collection address
- `allowlist_admin`: Allowlist administrator address
- `dest_chain_selectors`: List of supported destination chains
- `dest_chain_enabled`: Enable/disable flags for each chain
- `dest_chain_allowlist_enabled`: Allowlist enforcement flags

#### `apply_dest_chain_config_updates()`

Updates destination chain configurations.

**Access:** Owner only

#### `set_dynamic_config()`

Updates dynamic configuration parameters.

**Access:** Owner only

### Allowlist Management

#### `apply_allowlist_updates()`

Manages sender allowlists for specific destination chains.

**Access:** Allowlist admin only

**Parameters:**
- `dest_chain_selectors`: Target chains to update
- `dest_chain_allowlist_enabled`: Enable/disable allowlist per chain
- `dest_chain_add_allowed_senders`: Addresses to add to allowlist
- `dest_chain_remove_allowed_senders`: Addresses to remove from allowlist

## Security Features

### Access Control

The onramp implements multiple layers of access control:

1. **Owner Controls**: Configuration changes, fee withdrawal
2. **Allowlist Admin**: Sender allowlist management
3. **Capability-Based**: Token transfers and nonce management

### Validation

- **Sender Authorization**: Validates sender against allowlists when enabled
- **Destination Chain**: Ensures target chain is supported and enabled
- **Fee Validation**: Verifies sufficient fee payment
- **RMN Integration**: Checks for cursed chains via Risk Management Network
- **Token Validation**: Validates token transfers through dynamic dispatcher

### Message Integrity

- **Unique Message IDs**: Each message gets a cryptographically secure ID
- **Sequence Numbers**: Maintains ordered message delivery per destination
- **Nonce Management**: Prevents replay attacks and ensures ordering
- **Hash Verification**: Messages are hashed for integrity verification

## Integration with CCIP Core

The onramp integrates with several CCIP core components:

### Dynamic Dispatcher

- **Token Parameters**: Constructs token transfer parameters
- **Validation**: Ensures only authorized token pools can add transfers
- **Type Safety**: Uses proof-based validation for token operations

### Fee Quoter

- **Fee Calculation**: Determines gas costs and token fees
- **Cross-Chain Rates**: Handles exchange rates between chains
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
- **`ConfigSet`**: Configuration updates
- **`DestChainConfigSet`**: Destination chain configuration changes
- **`AllowlistSendersAdded/Removed`**: Allowlist modifications
- **`FeeTokenWithdrawn`**: Fee token withdrawal events

## Usage Examples

### Basic Message Sending

```move
use ccip_onramp::onramp;
use ccip::dynamic_dispatcher;

// Create token parameters (if sending tokens)
let token_params = dynamic_dispatcher::create_token_params(dest_chain_selector);

// Send message
let message_id = onramp::ccip_send<SUI>(
    &mut ccip_ref,
    &mut onramp_state,
    &clock,
    receiver_address,
    message_data,
    token_params,
    &fee_token_metadata,
    &mut fee_token_coin,
    extra_args,
    &mut ctx
);
```

### Fee Calculation

```move
// Get required fee before sending
let required_fee = onramp::get_fee<SUI>(
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

### Configuration Management

```move
// Update destination chain config (owner only)
onramp::apply_dest_chain_config_updates(
    &mut onramp_state,
    &owner_cap,
    dest_chain_selectors,
    dest_chain_enabled,
    dest_chain_allowlist_enabled
);
```

## Deployment and Operations

The onramp package includes operational tools for:

- **Deployment**: Automated package publishing and initialization
- **Configuration**: Chain configuration and allowlist management
- **Monitoring**: Chain support status and configuration queries
- **Maintenance**: Fee withdrawal and ownership management

## Ownership Model

The onramp uses a robust ownership model with:

- **Owner Capabilities**: Secure capability-based access control
- **Transfer Process**: Two-step ownership transfer with acceptance
- **Pending Transfers**: Support for pending ownership changes
- **MCMS Integration**: Multi-signature management support

## Package Information

- **Package Name**: ChainlinkCCIPOnramp
- **Version**: 1.6.0
- **Edition**: 2024.beta
- **Dependencies**: 
  - ChainlinkCCIP (core CCIP functionality)
  - ChainlinkManyChainMultisig (governance)

## Architecture

The onramp follows a modular design pattern:

1. **State Management**: Centralized state with capability-based access
2. **Message Construction**: Standardized message format with integrity checks
3. **Fee Processing**: Integrated fee calculation and collection
4. **Configuration**: Flexible destination chain and allowlist management
5. **Security**: Multi-layered validation and access control
6. **Monitoring**: Comprehensive event emission for observability

This architecture ensures secure, reliable cross-chain communication while maintaining flexibility for different use cases and operational requirements. 