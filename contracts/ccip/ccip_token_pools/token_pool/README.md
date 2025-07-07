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
    allowlist_addresses,
    ctx
);

// Configure supported chains
token_pool::apply_chain_updates(
    &mut pool_state,
    chains_to_remove,
    chains_to_add,
    pool_addresses_to_add,
    token_addresses_to_add
);

// Validate cross-chain operations
let remote_token = token_pool::validate_lock_or_burn(
    ccip_ref,
    clock,
    &mut pool_state,
    sender,
    destination_chain,
    amount
);
```

## Integration

This module integrates with:
- **ChainlinkCCIP**: Core CCIP protocol functionality
- **RMN Remote**: Risk Management Network for security
- **Allowlist**: Access control mechanism
- **Sui Framework**: Clock, events, and standard library functions

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