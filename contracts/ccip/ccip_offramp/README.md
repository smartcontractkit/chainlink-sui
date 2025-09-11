# Chainlink CCIP OffRamp

The CCIP OffRamp handles incoming cross-chain messages and token transfers from other supported blockchain networks to Sui. It serves as the destination point for all inbound cross-chain operations on Sui.

## Overview

The OffRamp manages merkle root commitments, message execution, and token transfers from source chains to Sui. It provides secure, verifiable execution of cross-chain messages while maintaining proper access controls and validation through OCR3 (Off-Chain Reporting) integration.

**Important**: Many OffRamp functions, particularly `init_execute` and `finish_execute`, are designed to be used within Programmable Transaction Blocks (PTBs) to ensure atomic execution and proper state management.

## Key Features

- **Cross-Chain Message Execution**: Execute incoming messages from other blockchains
- **Token Transfer Processing**: Handle token transfers from source chains to Sui
- **Merkle Root Commitments**: Secure commitment and verification of message batches
- **OCR3 Integration**: Off-chain reporting for price updates and execution reports
- **Access Control**: Configurable permissions and owner controls
- **Security**: RMN integration and message integrity validation
- **Message Ordering**: Sequence number management for reliable delivery
- **Multi-Chain Support**: Configurable support for multiple source chains
- **Upgrade Registry**: Integration with CCIP upgrade management system

## Core Functions

### Message Execution
**⚠️ PTB Required**: These functions must be used within Programmable Transaction Blocks.

```move
// Step 1: Initialize message execution
let receiver_params = offramp::init_execute(
    &ccip_ref,
    &mut offramp_state,
    &clock,
    report_context,
    report,
    token_receiver,
    &mut ctx
);

// Step 2: Interact with token pools and/or call message receivers
// This step involves:
// - Token pool interactions (if token transfers are present)
// - Message receiver calls (if messages need to be delivered)
// - The receiver_params "hot potato" is passed between functions and updated

// Example: Call a message receiver
let updated_params = your_receiver::ccip_receive(
    &ccip_ref,
    &mut your_receiver_state,
    receiver_package_id,
    receiver_params
);

// Step 3: Finish message execution
offramp::finish_execute(
    &ccip_ref,
    updated_params
);
```

**Process Flow:**
1. **`init_execute`**: Validates and processes execution report, creates receiver parameters (hot potato)
2. **Middle Step**: Token pool interactions and/or message receiver calls that update the receiver parameters
3. **`finish_execute`**: Inspects the updated receiver parameters and completes execution

**Important**: The middle step is where the actual business logic happens - token transfers are processed and messages are delivered to receivers. The `finish_execute` function validates that all required interactions have been completed properly.

### Hot Potato Pattern

The OffRamp uses a "hot potato" pattern where the `ReceiverParams` object is passed between functions and updated to track the execution state:

1. **`init_execute`** creates the initial `ReceiverParams` (hot potato) containing:
   - Token transfer information (if any)
   - Message data (if any)
   - Source chain information
   - Execution state tracking

2. **Middle Step** updates the hot potato through:
   - **Token Pool Interactions**: If token transfers are present, the token pools will update the receiver parameters to mark transfers as completed
   - **Message Receiver Calls**: If messages need delivery, the receiver functions will process the message and update the parameters
   - **State Tracking**: Each interaction updates the hot potato to track completion status

3. **`finish_execute`** inspects the hot potato to ensure:
   - All required token transfers have been completed
   - All required message deliveries have been processed
   - The execution state is consistent and valid

This pattern ensures atomic execution where either all operations succeed or none do, maintaining the integrity of cross-chain message processing.

### Merkle Root Commitments
**⚠️ PTB Required**: This function must be used within a Programmable Transaction Block.

```move
// Commit merkle roots for message verification
offramp::commit(
    &mut ccip_ref,
    &mut offramp_state,
    &clock,
    report_context,
    report,
    signatures,
    &mut ctx
);
```

**Process:**
1. Verifies OCR3 signatures on the commit report
2. Updates merkle roots for message verification
3. Processes price updates for fee calculation
4. Updates blessed merkle roots from RMN

### Manual Execution
**⚠️ PTB Required**: For emergency or manual message execution.

```move
// Manually execute messages without OCR3 transmission
let receiver_params = offramp::manually_init_execute(
    &ccip_ref,
    &mut offramp_state,
    &clock,
    report,
    token_receiver,
    &mut ctx
);
```

### Configuration Management
**Owner-only functions**: These require owner capabilities and should be used in PTBs.

```move
// Update source chain configuration (owner only)
offramp::apply_source_chain_config_updates(
    &mut offramp_state,
    &owner_cap,
    source_chain_selectors,
    source_chain_enabled,
    source_chain_rmn_verification_disabled
);

// Update dynamic configuration (owner only)
offramp::set_dynamic_config(
    &mut offramp_state,
    &owner_cap,
    permissionless_execution_threshold_seconds
);

// Set OCR3 configuration (owner only)
offramp::set_ocr3_config(
    &mut offramp_state,
    &owner_cap,
    ocr_config
);
```

### Query Functions
**Safe for standalone calls**: These functions can be called outside of PTBs for information retrieval.

```move
// Get execution state for a message
let execution_state = offramp::get_execution_state(
    &offramp_state,
    source_chain_selector,
    sequence_number
);

// Get source chain configuration
let config = offramp::get_source_chain_config(
    &offramp_state,
    source_chain_selector
);

// Get latest price sequence number
let price_seq = offramp::get_latest_price_sequence_number(&offramp_state);

// Get merkle root timestamp
let timestamp = offramp::get_merkle_root(&offramp_state, merkle_root);
```

## State Structure

### OffRampState
The main state object that manages all offramp operations and configurations.

- `id`: Unique object identifier
- `package_ids`: Vector of package IDs for upgrade tracking
- `ocr3_base_state`: OCR3 base state for off-chain reporting
- `chain_selector`: Sui chain identifier (unique per network)
- `permissionless_execution_threshold_seconds`: Threshold for permissionless execution
- `source_chain_configs`: Map of source chain configurations
- `execution_states`: Table tracking execution states per chain and sequence
- `roots`: Table of merkle roots and their timestamps
- `latest_price_sequence_number`: Latest OCR sequence number for price updates
- `fee_quoter_cap`: Capability for fee quoter operations
- `dest_transfer_cap`: Capability for destination token transfers
- `ownable_state`: Ownership management state

### SourceChainConfig
Configuration for each supported source chain.

- `router`: Router address on the source chain
- `is_enabled`: Whether the source chain is enabled for message processing
- `min_seq_nr`: Minimum sequence number for processing
- `is_rmn_verification_disabled`: Whether RMN verification is disabled
- `on_ramp`: OnRamp address on the source chain

### Message Structures

#### RampMessageHeader
- `message_id`: Cryptographically secure unique message identifier
- `source_chain_selector`: Source chain selector
- `dest_chain_selector`: Destination chain selector (Sui)
- `sequence_number`: Message ordering number
- `nonce`: Sender-specific nonce for replay protection

#### Any2SuiRampMessage
Incoming cross-chain message structure.

- `header`: Message metadata and routing information
- `sender`: Sender address on source chain
- `data`: Arbitrary message payload
- `receiver`: Destination address on Sui
- `gas_limit`: Gas limit for execution
- `token_amounts`: Vector of token transfers included in message

#### Any2SuiTokenTransfer
Token transfer information within an incoming message.

- `source_pool_address`: Source token pool address
- `dest_token_address`: Token address on Sui (coin metadata)
- `dest_gas_amount`: Gas amount for destination execution
- `extra_data`: Additional token-specific data
- `amount`: Token amount to transfer

## Security Features

### Access Control
The offramp implements multiple layers of access control:

1. **Owner Controls**: Configuration changes, OCR3 config updates
2. **OCR3 Integration**: Off-chain reporting for secure message processing
3. **Capability-Based**: Token transfers and fee quoter operations

### Validation Mechanisms
- **Message Verification**: Merkle proof verification for message authenticity
- **Source Chain Validation**: Ensures source chain is supported and enabled
- **RMN Integration**: Risk management network validation for blessed roots
- **Sequence Number Validation**: Prevents replay attacks and ensures ordering
- **Signature Verification**: OCR3 signature verification for reports

### OCR3 Integration
The offramp integrates with OCR3 for:
- **Execution Reports**: Secure transmission of message execution data
- **Commit Reports**: Merkle root commitments and price updates
- **Signature Verification**: Multi-signature validation for reports
- **Price Updates**: Token and gas price updates for fee calculation

## Integration with CCIP Core

The offramp integrates with several CCIP core components:

### OffRamp State Helper
- **Receiver Parameters**: Manages receiver parameters for message processing
- **Token Transfers**: Handles destination token transfers
- **Message Extraction**: Extracts messages for receiver processing

### Fee Quoter
- **Price Updates**: Receives token and gas price updates
- **Fee Calculation**: Uses updated prices for fee calculations

### Receiver Registry
- **Receiver Validation**: Validates message receivers
- **Type Safety**: Ensures proper receiver configuration

### RMN Remote
- **Risk Management**: Validates blessed merkle roots
- **Security Verification**: Checks for cursed chains and halted operations

## Events

The offramp emits several events for monitoring and indexing:

- **`ExecutionReportTransmitted`**: OCR3 execution report transmission
- **`CommitReportTransmitted`**: OCR3 commit report transmission
- **`ConfigSet`**: Configuration updates
- **`SourceChainConfigSet`**: Source chain configuration changes
- **`OwnershipTransferred`**: Ownership transfer events

## MCMS Integration

The offramp supports Multi-Chain Multi-Signature (MCMS) governance through several functions:

```move
// MCMS-specific functions for governance
offramp::mcms_set_dynamic_config(...);
offramp::mcms_apply_source_chain_config_updates(...);
offramp::mcms_set_ocr3_config(...);
offramp::mcms_transfer_ownership(...);
offramp::mcms_execute_ownership_transfer(...);
```

## Dependencies

- **ChainlinkCCIP**: Core CCIP functionality
- **ChainlinkManyChainMultisig**: Governance support

## Package Information

- **Name**: ChainlinkCCIPOfframp
- **Version**: 1.6.0
- **Edition**: 2024

## Usage

### Initialization
1. **Deploy Package**: Publish the offramp package to Sui
2. **Initialize State**: Set up offramp with chain selector and configurations
3. **Configure Source Chains**: Add supported source chains and their configurations
4. **Set OCR3 Config**: Configure OCR3 for off-chain reporting

### Message Processing
1. **Commit Phase**: Use `commit` to process merkle root commitments and price updates
2. **Execution Phase**: Use the three-step execution process:
   - `init_execute`: Initialize execution and get receiver parameters
   - **Middle Step**: Interact with token pools and/or call message receivers
   - `finish_execute`: Complete execution and validate all interactions
3. **Manual Execution**: Use `manually_init_execute` for emergency processing
4. **Monitor Events**: Track execution and commit events for status

### Management Operations
1. **Configuration Updates**: Use owner functions to update source chain configs
2. **OCR3 Management**: Use owner functions to update OCR3 configuration
3. **Emergency Controls**: Use manual execution for emergency message processing

### Best Practices
- Always use execution functions within PTBs for atomic execution
- Monitor OCR3 reports for proper message processing
- Use query functions for status monitoring and debugging
- Keep source chain configurations up-to-date
- Monitor events for execution status and errors

For detailed implementation examples and advanced configuration, see the source code and integration tests.
