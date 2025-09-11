# CCIP Token Pool

A foundational Move module for managing cross-chain token transfers within the Chainlink CCIP (Cross-Chain Interoperability Protocol) ecosystem on Sui.

## Overview

The Token Pool module provides core functionality for secure cross-chain token transfers, including validation, rate limiting, and ownership management. It serves as a base implementation for specific token pool types (burn/mint, lock/release, etc.).

## Key Components

### Core Module (`token_pool.move`)
- **Cross-chain validation**: Validates lock/burn and release/mint operations across chains
- **Chain management**: Supports multiple remote chains with configurable token addresses and pool addresses
- **Rate limiting**: Configurable inbound/outbound rate limits per chain to prevent abuse
- **Decimal handling**: Automatic conversion between different token decimal representations across chains
- **Allowlist support**: Optional sender allowlist for restricted access

### Ownership Management (`ownable.move`)
- **Two-step ownership transfer**: Secure ownership transfer requiring explicit acceptance
- **Multiple transfer modes**: Support for standard users, objects, and MCMS (Multi-Chain Multi-Sig)
- **Event emission**: Comprehensive ownership change tracking

### Rate Limiting (`rate_limiter.move`, `token_pool_rate_limiter.move`)
- **Per-chain configuration**: Independent rate limits for each supported chain
- **Inbound/outbound controls**: Separate rate limiting for incoming and outgoing transfers
- **Time-based buckets**: Token bucket algorithm for smooth rate limiting

## Key Features

- ✅ **Multi-chain support**: Configure multiple remote chains with their respective token and pool addresses
- ✅ **Rate limiting**: Prevent abuse with configurable per-chain rate limits
- ✅ **Decimal compatibility**: Handle tokens with different decimal precision across chains
- ✅ **Access control**: Owner-only operations with secure two-step ownership transfer
- ✅ **RMN integration**: Respects Risk Management Network curse status
- ✅ **Event emission**: Comprehensive event logging for monitoring and indexing

## Usage

This module is designed to be used as a foundation for specific token pool implementations:

```move
use ccip_token_pool::token_pool;

// Initialize token pool state
let pool_state = token_pool::initialize(
    coin_metadata_address,
    local_decimals,
    allowlist,  // vector<address>
    ctx
);

// Configure supported chains
token_pool::apply_chain_updates(
    &mut pool_state,
    remote_chain_selectors_to_remove,  // vector<u64>
    remote_chain_selectors_to_add,     // vector<u64>
    remote_pool_addresses_to_add,      // vector<vector<vector<u8>>>
    remote_token_addresses_to_add      // vector<vector<u8>>
);

// Validate cross-chain operations
let remote_token = token_pool::validate_lock_or_burn(
    ccip_ref,                    // &CCIPObjectRef
    clock,                       // &Clock
    &mut pool_state,
    sender,                      // address
    remote_chain_selector,       // u64 (destination chain)
    local_amount                 // u64
);
```

## Core Functions Reference

### Initialization & Configuration

| Function | Description | Parameters |
|----------|-------------|------------|
| `initialize()` | Initialize token pool state | coin_metadata_address, local_decimals, allowlist, ctx |
| `apply_chain_updates()` | Add/remove supported chains | state, chains_to_remove, chains_to_add, pool_addresses, token_addresses |
| `add_remote_pool()` | Add remote pool address | state, remote_chain_selector, remote_pool_address |
| `remove_remote_pool()` | Remove remote pool address | state, remote_chain_selector, remote_pool_address |

### Validation Functions

| Function | Description | Returns |
|----------|-------------|---------|
| `validate_lock_or_burn()` | Validate outbound operations | vector<u8> (remote token address) |
| `validate_release_or_mint()` | Validate inbound operations | vector<u8> (remote token address) |

### Query Functions

| Function | Description | Returns |
|----------|-------------|---------|
| `get_token()` | Get token metadata address | address |
| `get_supported_chains()` | Get list of supported chains | vector<u64> |
| `is_supported_chain()` | Check if chain is supported | bool |
| `get_remote_pools()` | Get remote pool addresses for chain | vector<vector<u8>> |
| `is_remote_pool()` | Check if address is valid remote pool | bool |
| `get_remote_token()` | Get remote token address for chain | vector<u8> |
| `get_local_decimals()` | Get local token decimals | u8 |

### Amount Calculation

| Function | Description | Returns |
|----------|-------------|---------|
| `calculate_local_amount()` | Convert remote amount to local | u64 |
| `calculate_release_or_mint_amount()` | Calculate release/mint amount | u64 |
| `encode_local_decimals()` | Encode decimals for cross-chain | vector<u8> |
| `parse_remote_decimals()` | Parse remote decimals from data | u8 |

### Rate Limiting

| Function | Description | Parameters |
|----------|-------------|------------|
| `set_chain_rate_limiter_config()` | Configure rate limits for chain | clock, state, chain_selector, config |

### Allowlist Management

| Function | Description | Returns |
|----------|-------------|---------|
| `get_allowlist_enabled()` | Check if allowlist is enabled | bool |
| `set_allowlist_enabled()` | Enable/disable allowlist | state, enabled |
| `get_allowlist()` | Get allowlist addresses | vector<address> |
| `apply_allowlist_updates()` | Update allowlist | state, addresses_to_remove, addresses_to_add |

### Event Emission

| Function | Description | Purpose |
|----------|-------------|---------|
| `emit_locked_or_burned()` | Emit lock/burn event | Track outbound operations |
| `emit_released_or_minted()` | Emit release/mint event | Track inbound operations |
| `emit_liquidity_added()` | Emit liquidity addition | Track liquidity changes |
| `emit_liquidity_removed()` | Emit liquidity removal | Track liquidity changes |
| `emit_rebalancer_set()` | Emit rebalancer change | Track rebalancer updates |

### Utility Functions

| Function | Description | Returns |
|----------|-------------|---------|
| `get_token_decimals()` | Get decimals from coin metadata | u8 |
| `destroy_token_pool()` | Destroy token pool state | void |

## Integration

This module integrates with:
- **ChainlinkCCIP**: Core CCIP protocol functionality
- **RMN Remote**: Risk Management Network for security
- **Allowlist**: Access control mechanism
- **Sui Framework**: Clock, events, and standard library functions

## Package Information

- **Package Name**: `CCIPTokenPool`
- **Edition**: 2024
- **Dependencies**: 
  - `ChainlinkCCIP` (Core CCIP functionality)

### Module Structure
- `token_pool.move` - Core token pool functionality and validation
- `ownable.move` - Ownership management with two-step transfer
- `rate_limiter.move` - Token bucket rate limiting implementation
- `token_pool_rate_limiter.move` - Per-chain rate limiting configuration

### Key Data Structures
- `TokenPoolState` - Main pool state with allowlist, metadata, and chain configs
- `RemoteChainConfig` - Configuration for each supported remote chain
- `OwnerCap` - Ownership capability for administrative functions
- `OwnableState` - Ownership state with pending transfer tracking

## Security Considerations

- All cross-chain operations validate RMN curse status
- Rate limiting prevents rapid token drainage
- Ownership changes require two-step confirmation
- Remote pool addresses are validated against configured lists
- Decimal overflow protection prevents arithmetic errors

## Testing

Run tests with:
```bash
sui move test
```

The test suite covers core functionality, edge cases, and security validations. 