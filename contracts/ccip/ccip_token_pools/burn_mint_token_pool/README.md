# Burn Mint Token Pool

A CCIP token pool implementation that burns tokens on the source chain and mints them on the destination chain. This pool maintains the treasury cap and controls all token minting/burning operations for cross-chain transfers.

## Overview

The Burn Mint Token Pool is designed for tokens that need to maintain a consistent total supply across all chains. When tokens are sent cross-chain, they are burned on the source chain. When tokens are received from other chains, they are minted on the destination chain. This ensures the total supply remains constant across all chains.

**⚠️ Important**: This token pool stores the treasury cap object for the coin, which means no other parties (whether objects or EOAs) can burn or mint the coin outside of this token pool on Sui. All minting and burning operations must go through this token pool.

**💡 Alternative**: If you need to burn/mint this token outside of this token pool, wrap your treasury cap in the managed token package and use the managed token pool package to deploy your token pool instead.

## Key Features

- ✅ **Token Burning**: Burns tokens on source chain for outbound cross-chain transfers
- ✅ **Token Minting**: Mints tokens on destination chain for inbound cross-chain transfers
- ✅ **Treasury Cap Control**: Maintains treasury cap to control all minting operations (exclusive control)
- ✅ **Multi-Chain Support**: Supports multiple remote chains with configurable token addresses
- ✅ **Rate Limiting**: Per-chain rate limiting to prevent abuse
- ✅ **Access Control**: Owner-only operations with MCMS integration
- ✅ **Allowlist Support**: Optional sender allowlist for restricted access
- ✅ **RMN Integration**: Respects Risk Management Network curse status

## Architecture

### State Structure

```move
public struct BurnMintTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,  // Base token pool functionality
    treasury_cap: TreasuryCap<T>,      // Treasury cap for minting control
    ownable_state: OwnableState,       // Ownership management
}
```

### Core Operations

1. **Burn Operation**: When tokens are sent cross-chain, they are burned on the source chain
2. **Mint Operation**: When tokens are received from other chains, they are minted on the destination chain
3. **Supply Management**: Total supply remains constant across all chains through burn/mint operations

## Usage

### Initialization

```move
use burn_mint_token_pool::burn_mint_token_pool;

// Initialize with treasury cap (standard deployment)
burn_mint_token_pool::initialize<MyToken>(
    &mut ccip_ref,
    &coin_metadata,
    treasury_cap,  // Treasury cap is consumed and stored in state
    token_pool_administrator,
    ctx
);

// Initialize by CCIP admin (admin deployment)
burn_mint_token_pool::initialize_by_ccip_admin<MyToken>(
    &mut ccip_ref,
    ccip_admin_proof,
    &coin_metadata,
    treasury_cap,  // Treasury cap is consumed and stored in state
    token_pool_administrator,
    ctx
);
```

### Cross-Chain Operations

```move
// Burn tokens for outbound transfer
burn_mint_token_pool::lock_or_burn<MyToken>(
    &mut ccip_ref,
    &clock,
    &mut pool_state,
    sender,
    destination_chain,
    amount,
    ctx
);

// Mint tokens for inbound transfer
burn_mint_token_pool::release_or_mint<MyToken>(
    &mut ccip_ref,
    &clock,
    &mut pool_state,
    recipient,
    source_chain,
    amount,
    ctx
);
```

## Core Functions Reference

### Initialization Functions

| Function | Description | Parameters |
|----------|-------------|------------|
| `initialize<T>()` | Initialize pool with treasury cap | ccip_ref, coin_metadata, treasury_cap, admin, ctx |
| `initialize_by_ccip_admin<T>()` | Initialize by CCIP admin | ccip_ref, admin_proof, coin_metadata, treasury_cap, admin, ctx |

### Cross-Chain Operations

| Function | Description | Parameters |
|----------|-------------|------------|
| `lock_or_burn<T>()` | Burn tokens for outbound transfer | ccip_ref, clock, state, sender, chain, amount, ctx |
| `release_or_mint<T>()` | Mint tokens for inbound transfer | ccip_ref, clock, state, recipient, chain, amount, ctx |

### Chain Management

| Function | Description | Parameters |
|----------|-------------|------------|
| `apply_chain_updates<T>()` | Add/remove supported chains | state, chains_to_remove, chains_to_add, pool_addresses, token_addresses |
| `add_remote_pool<T>()` | Add remote pool address | state, owner_cap, chain_selector, pool_address, ctx |
| `remove_remote_pool<T>()` | Remove remote pool address | state, owner_cap, chain_selector, pool_address, ctx |
| `get_supported_chains<T>()` | Get supported chain selectors | state |
| `is_supported_chain<T>()` | Check if chain is supported | state, chain_selector |

### Query Functions

| Function | Description | Returns |
|----------|-------------|---------|
| `get_token<T>()` | Get token metadata address | address |
| `get_token_decimals<T>()` | Get token decimals | u8 |
| `get_remote_pools<T>()` | Get remote pool addresses | vector<vector<u8>> |
| `is_remote_pool<T>()` | Check if address is remote pool | bool |
| `get_remote_token<T>()` | Get remote token address | vector<u8> |

### Allowlist Management

| Function | Description | Parameters |
|----------|-------------|------------|
| `get_allowlist_enabled<T>()` | Check if allowlist is enabled | state |
| `set_allowlist_enabled<T>()` | Enable/disable allowlist | state, owner_cap, enabled, ctx |
| `get_allowlist<T>()` | Get allowlist addresses | state |
| `apply_allowlist_updates<T>()` | Update allowlist | state, owner_cap, addresses_to_remove, addresses_to_add, ctx |

### Rate Limiting

| Function | Description | Parameters |
|----------|-------------|------------|
| `set_chain_rate_limiter_configs<T>()` | Set rate limits for multiple chains | clock, state, owner_cap, configs, ctx |
| `set_chain_rate_limiter_config<T>()` | Set rate limit for single chain | clock, state, owner_cap, chain_selector, config, ctx |

### Ownership Management

| Function | Description | Parameters |
|----------|-------------|------------|
| `owner<T>()` | Get current owner | state |
| `transfer_ownership<T>()` | Transfer ownership | state, owner_cap, new_owner, ctx |
| `accept_ownership<T>()` | Accept ownership transfer | state, ctx |
| `has_pending_transfer<T>()` | Check for pending transfer | state |
| `pending_transfer_from<T>()` | Get pending transfer from | state |
| `pending_transfer_to<T>()` | Get pending transfer to | state |
| `pending_transfer_accepted<T>()` | Check if transfer accepted | state |

### MCMS Integration Functions

| Function | Description | Purpose |
|----------|-------------|---------|
| `mcms_accept_ownership<T>()` | Accept ownership via MCMS | Multi-sig ownership transfer |
| `mcms_register_upgrade_cap()` | Register upgrade capability | Multi-sig upgrade management |
| `mcms_set_allowlist_enabled<T>()` | Set allowlist via MCMS | Multi-sig access control |
| `mcms_apply_allowlist_updates<T>()` | Update allowlist via MCMS | Multi-sig allowlist management |
| `mcms_apply_chain_updates<T>()` | Update chains via MCMS | Multi-sig chain management |
| `mcms_add_remote_pool<T>()` | Add remote pool via MCMS | Multi-sig pool management |
| `mcms_remove_remote_pool<T>()` | Remove remote pool via MCMS | Multi-sig pool management |
| `mcms_set_chain_rate_limiter_configs<T>()` | Set rate limits via MCMS | Multi-sig rate limiting |
| `mcms_set_chain_rate_limiter_config<T>()` | Set rate limit via MCMS | Multi-sig rate limiting |
| `mcms_transfer_ownership<T>()` | Transfer ownership via MCMS | Multi-sig ownership transfer |
| `mcms_execute_ownership_transfer<T>()` | Execute transfer via MCMS | Multi-sig ownership transfer |

### Utility Functions

| Function | Description | Parameters |
|----------|-------------|------------|
| `destroy_token_pool<T>()` | Destroy pool state | state, owner_cap, ctx |

## Integration

This module integrates with:
- **CCIPTokenPool**: Base token pool functionality
- **ChainlinkCCIP**: Core CCIP protocol functionality
- **MCMS**: Multi-Chain Multi-Sig governance
- **RMN Remote**: Risk Management Network for security
- **Sui Framework**: Clock, events, and standard library functions

## Security Considerations

- All cross-chain operations validate RMN curse status
- Rate limiting prevents rapid token drainage
- Ownership changes require two-step confirmation
- Remote pool addresses are validated against configured lists
- Treasury cap is stored in state, preventing external minting
- **Exclusive Control**: No other parties can mint or burn tokens outside this pool
- Decimal overflow protection prevents arithmetic errors

## Package Information

- **Package Name**: `BurnMintTokenPool`
- **Version**: 1.6.0
- **Edition**: 2024
- **Dependencies**: 
  - `CCIPTokenPool` (Base token pool functionality)
  - `ChainlinkManyChainMultisig` (MCMS integration)

### Module Structure
- `burn_mint_token_pool.move` - Core burn/mint functionality

## Testing

Run tests with:
```bash
sui move test
```

The test suite covers:
- Initialization and configuration
- Burn and mint operations
- Chain configuration
- Rate limiting
- Access control
- Ownership management
- MCMS integration
- Error conditions and edge cases

## Comparison with Lock Release Token Pool

| Feature | Burn Mint Token Pool | Lock Release Token Pool |
|---------|---------------------|------------------------|
| **Token Management** | Burns on source, mints on destination | Locks on source, releases on destination |
| **Supply Control** | Total supply constant across chains | Total supply varies by chain |
| **Liquidity Management** | No liquidity management needed | Requires rebalancer for liquidity |
| **Use Case** | Single token across multiple chains | Token exists on multiple chains |
| **Complexity** | Simpler, no liquidity concerns | More complex, requires rebalancing |
| **Treasury Cap** | Stored in state, controls all minting | Not stored, external minting possible |
| **External Minting** | Not possible (exclusive control) | Possible (use managed token + managed token pool) |

## Alternative: Managed Token + Managed Token Pool

If you need flexibility to burn/mint tokens outside of the token pool:

1. **Use Managed Token Package**: Wrap your treasury cap in the managed token package
2. **Use Managed Token Pool Package**: Deploy your token pool using the managed token pool package
3. **Benefits**: 
   - Maintains treasury cap control through managed token
   - Allows external minting/burning through managed token functions
   - Provides more flexibility for token management
   - Still supports cross-chain operations through CCIP
